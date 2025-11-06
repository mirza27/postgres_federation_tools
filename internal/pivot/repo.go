package pivot

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"db_migrate_server/internal/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo menyatukan seluruh interaksi dengan database pivot sehingga worker lain
// dapat fokus pada orkestrasi event. Struktur ini bertindak layaknya façade
// yang menyembunyikan detail tabel pivot (_exec_queue, _keymap_generic).
type Repo struct {
	DB *pgxpool.Pool
}

// New membuka koneksi pool ke database pivot dan siap dipakai oleh worker.
// Jika koneksi gagal, worker tidak akan bisa melanjutkan proses ingest.
func New(ctx context.Context, dsn string) (*Repo, error) {
	util.Info.Printf("pivot: connecting dsn=%s", dsn)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		// Koneksi gagal dibuka, bubble up agar caller bisa menangani (misalnya retry).
		return nil, err
	}
	return &Repo{DB: pool}, nil
}

// Close menutup pool pivot ketika aplikasi berhenti agar resource database
// dirilis dengan rapi.
func (r *Repo) Close() { r.DB.Close() }

// EnsureSchema memastikan seluruh tabel pivot tersedia sebelum worker mulai
// membaca event. Fungsi ini idempoten sehingga aman dipanggil berkali-kali.
func (r *Repo) EnsureSchema(ctx context.Context, schemaSQL string) error {
	util.Info.Println("pivot: ensuring schema")
	_, err := r.DB.Exec(ctx, schemaSQL)
	return err
}

// KeyMap helpers ------------------------------------------------------------

// LookupKey mencari apakah sudah ada padanan key yang selesai diproses.
// Ini dipakai parser dan executor untuk menghindari pembuatan key map yang
// sebenarnya sudah tersedia.
func (r *Repo) LookupKey(ctx context.Context, mapName, srcKey string) (string, bool, error) {
	util.Debug.Printf("pivot: LookupKey map=%s src=%s", mapName, srcKey)

	var tgt sql.NullString
	err := r.DB.QueryRow(ctx, `
        select tgt_key
        from _keymap_generic
        where map_name=$1 and src_key=$2 and tgt_key_status='fulfilled'`,
		mapName, srcKey).Scan(&tgt)

	if err != nil {
		// Bila tidak ada baris, kita informasikan ke caller bahwa padanan belum tersedia.
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		// Untuk error lain, kembalikan agar lapisan di atas bisa mengambil tindakan.
		return "", false, err
	}

	// Jika baris ada tapi nilai target belum diisi, treat sebagai belum siap.
	if !tgt.Valid {
		return "", false, nil
	}

	return tgt.String, true, nil
}

// EnsureKeymapRequest membuat atau memperbarui permintaan padanan key. Checker
// memanggil fungsi ini setiap kali mendeteksi baris queue yang membutuhkan
// key baru dari database target.
func (r *Repo) EnsureKeymapRequest(ctx context.Context, queueID uuid.UUID, km KeymapRequest) (int64, error) {
	util.Debug.Printf("pivot: EnsureKeymapRequest map=%s src=%s queue=%s", km.MapName, km.SrcKey, queueID)

	// melakukan upsert pada tabel _keymap_generic sehingga setiap baris _exec_queue yang dibaca checker apapun statusnya
	var keymapID int64
	err := r.DB.QueryRow(ctx, `
        insert into _keymap_generic (
            map_name, src_key, src_column, tgt_column, src_table, tgt_table,
            queue_id, tgt_key_status, requested_at, attempts
        )
        values ($1,$2,$3,$4,$5,$6,$7,'requested',now(),1)
        on conflict (map_name, src_key, src_table, tgt_table)
        do update set
            queue_id = excluded.queue_id,
            tgt_key_status = case when _keymap_generic.tgt_key_status = 'fulfilled'
                then _keymap_generic.tgt_key_status else 'requested' end,
            requested_at = case when _keymap_generic.tgt_key_status = 'fulfilled'
                then _keymap_generic.requested_at else excluded.requested_at end,
            attempts = _keymap_generic.attempts + 1,
            last_error = null
        returning keymap_id`,
		km.MapName, km.SrcKey, nullIfEmpty(km.SrcColumn), nullIfEmpty(km.TgtColumn),
		nullIfEmpty(km.SrcTable), nullIfEmpty(km.TgtTable), queueID).Scan(&keymapID)

	return keymapID, err
}

// FulfillKeymap dipanggil executor setelah berhasil mendapatkan nilai RETURNING
// dari database target. Nilai tersebut disimpan sehingga entitas lain bisa
// menggunakan padanan yang sama.
func (r *Repo) FulfillKeymap(ctx context.Context, keymapID int64, tgtKey string) error {
	util.Debug.Printf("pivot: FulfillKeymap keymap_id=%d tgt=%s", keymapID, truncateSQL(tgtKey))

	_, err := r.DB.Exec(ctx, `
        update _keymap_generic
        set tgt_key=$2, tgt_key_status='fulfilled', fulfilled_at=now(), last_error=null
        where keymap_id=$1`, keymapID, tgtKey)
	return err
}

// MarkKeymapError menyimpan catatan kegagalan saat mencoba memenuhi keymap,
// sehingga worker lain bisa melihat penyebab terakhir sebelum melakukan retry.
func (r *Repo) MarkKeymapError(ctx context.Context, keymapID int64, msg string) error {
	util.Warn.Printf("pivot: MarkKeymapError keymap_id=%d msg=%s", keymapID, msg)

	_, err := r.DB.Exec(ctx, `
        update _keymap_generic
        set tgt_key_status='error', last_error=$2, attempts=attempts+1
        where keymap_id=$1`, keymapID, msg)
	return err
}

// Exec queue ----------------------------------------------------------------

// Enqueue dipanggil parser untuk menyimpan rencana eksekusi ke _exec_queue.
// Informasi inilah yang nantinya dikonsumsi checker maupun executor.
func (r *Repo) Enqueue(ctx context.Context, it ExecItem) error {
	// Pastikan setiap item memiliki queue_id agar statusnya mudah ditelusuri.
	if it.QueueID == uuid.Nil {
		it.QueueID = uuid.New()
	}

	// Argumen SQL disimpan sebagai JSON agar bisa dideserialisasi oleh executor.
	argsJSON, err := json.Marshal(it.Args)
	if err != nil {
		// Tidak bisa menyimpan argumen ketika serialisasi gagal.
		return err
	}

	var keymapArg interface{}
	// Jika pekerjaan membutuhkan padanan key, sertakan payload permintaan.
	if it.Keymap != nil {
		payload, err := json.Marshal(it.Keymap)
		if err != nil {
			// Tanpa payload keymap, checker tidak punya informasi untuk menulis map.
			return err
		}
		keymapArg = string(payload)
	}

	util.Debug.Printf("pivot: Enqueue queue=%s entity=%s op=%s sql=%s needKeymap=%t",
		it.QueueID, it.Entity, it.Op, truncateSQL(it.SQL), it.NeedKeymap)

	_, err = r.DB.Exec(ctx, `
        insert into _exec_queue (
            queue_id, entity, op, sql_text, sql_args,
            returning_cols, keymap_payload, need_keymap, status
        )
        values ($1,$2,$3,$4,$5::jsonb,$6,$7,$8,'pending')`,
		it.QueueID, it.Entity, it.Op, it.SQL, string(argsJSON), it.Returning, keymapArg, it.NeedKeymap)
	return err
}

// FetchReady memberikan daftar pekerjaan yang siap dijalankan executor. Hanya
// baris dengan status 'ready' yang dikembalikan sehingga executor tidak perlu
// menyaring ulang di sisi aplikasi.
func (r *Repo) FetchReady(ctx context.Context, limit int) ([]Row, error) {
	util.Debug.Printf("pivot: FetchReady limit=%d", limit)

	rows, err := r.DB.Query(ctx, `
        select id, queue_id, entity, op, sql_text,
               coalesce(sql_args,'[]'::jsonb), need_keymap, keymap_id,
               coalesce(returning_cols,'{}'::text[]), status
        from _exec_queue
        where status='ready'
        order by created_at asc
        limit $1`, limit)
	if err != nil {
		// Ketika query gagal, executor tidak boleh melanjutkan.
		return nil, err
	}
	defer rows.Close()

	var out []Row
	for rows.Next() {
		var (
			id         int64
			entity     string
			op         string
			sqlText    string
			args       []byte
			needKeymap bool
			keymapID   pgtype.Int8
			returning  []string
			status     string
			queueID    uuid.UUID
		)

		if err := rows.Scan(&id, &queueID, &entity, &op, &sqlText, &args,
			&needKeymap, &keymapID, &returning, &status); err != nil {
			// Jika ada format data yang tidak sesuai, hentikan supaya tidak mengembalikan data korup.
			return nil, err
		}

		var kmPtr *int64
		// Simpan pointer hanya jika ada relasi dengan baris keymap.
		if keymapID.Valid {
			val := keymapID.Int64
			kmPtr = &val
		}

		out = append(out, Row{
			ID:         id,
			QueueID:    queueID,
			Entity:     entity,
			Op:         op,
			SQL:        sqlText,
			ArgsJSON:   args,
			NeedKeymap: needKeymap,
			KeymapID:   kmPtr,
			Returning:  returning,
			Status:     status,
		})
	}
	return out, rows.Err()
}

// FetchForChecker menyiapkan data yang perlu diproses worker checker. Hanya
// baris pending yang dikembalikan agar checker fokus pada tugas screening.
func (r *Repo) FetchForChecker(ctx context.Context, limit int) ([]CheckerItem, error) {
	util.Debug.Printf("pivot: FetchForChecker limit=%d", limit)

	rows, err := r.DB.Query(ctx, `
        select id, queue_id, entity, need_keymap, keymap_id, status,
               coalesce(keymap_payload,'{}'::jsonb)
        from _exec_queue
        where status='pending'
        order by created_at asc
        limit $1`, limit)
	if err != nil {
		// Checker tidak punya sumber lain; laporkan error ke caller.
		return nil, err
	}
	defer rows.Close()

	var out []CheckerItem
	for rows.Next() {
		var (
			id         int64
			entity     string
			needKeymap bool
			status     string
			keymapID   pgtype.Int8
			payload    []byte
			queueID    uuid.UUID
		)

		if err := rows.Scan(&id, &queueID, &entity, &needKeymap, &keymapID, &status, &payload); err != nil {
			// Data yang tidak bisa diparsing lebih baik gagal lebih awal daripada memproses info salah.
			return nil, err
		}

		var kmPtr *int64
		// Hanya set pointer bila keymap_id sudah diketahui.
		if keymapID.Valid {
			val := keymapID.Int64
			kmPtr = &val
		}

		out = append(out, CheckerItem{
			ID:            id,
			QueueID:       queueID,
			Entity:        entity,
			NeedKeymap:    needKeymap,
			KeymapID:      kmPtr,
			Status:        status,
			KeymapPayload: payload,
		})
	}
	return out, rows.Err()
}

// MarkReady mengubah status queue menjadi siap dieksekusi oleh executor.
func (r *Repo) MarkReady(ctx context.Context, id int64) error {
	util.Debug.Printf("pivot: MarkReady id=%d", id)
	_, err := r.DB.Exec(ctx, `
        update _exec_queue
        set status='ready', updated_at=now(), last_error=null
        where id=$1`, id)
	return err
}

// AttachKeymap menautkan baris queue dengan catatan keymap tertentu dan
// mengatur status sesuai kebutuhan checker (misal ready atau awaiting).
func (r *Repo) AttachKeymap(ctx context.Context, id int64, keymapID int64, status string) error {
	util.Debug.Printf("pivot: AttachKeymap id=%d keymap_id=%d status=%s", id, keymapID, status)
	_, err := r.DB.Exec(ctx, `
        update _exec_queue
        set keymap_id=$2, need_keymap=true, status=$3, updated_at=now(), last_error=null
        where id=$1`, id, keymapID, status)
	return err
}

// MarkExecuting memberi tanda bahwa executor tertentu sedang memproses baris
// queue, sehingga worker lain tidak mengambil pekerjaan yang sama.
func (r *Repo) MarkExecuting(ctx context.Context, id int64, worker string) error {
	util.Debug.Printf("pivot: MarkExecuting id=%d worker=%s", id, worker)
	_, err := r.DB.Exec(ctx, `
        update _exec_queue
        set status='executing', locked_by=$2, locked_at=now(), updated_at=now()
        where id=$1`, id, worker)
	return err
}

// MarkDone menandai bahwa eksekusi berhasil sehingga baris queue bisa
// dianggap selesai dan lock dibersihkan.
func (r *Repo) MarkDone(ctx context.Context, id int64) error {
	util.Debug.Printf("pivot: MarkDone id=%d", id)
	_, err := r.DB.Exec(ctx, `
        update _exec_queue
        set status='done', last_error=null, updated_at=now(), locked_at=null, locked_by=null
        where id=$1`, id)
	return err
}

// MarkError mencatat kegagalan eksekusi dan menambah retry_count sebagai bahan
// monitoring agar error berulang tidak dibiarkan.
func (r *Repo) MarkError(ctx context.Context, id int64, msg string) error {
	util.Warn.Printf("pivot: MarkError id=%d msg=%s", id, msg)
	_, err := r.DB.Exec(ctx, `
        update _exec_queue
        set status='error', last_error=$2, retry_count=retry_count+1,
            updated_at=now(), locked_at=null, locked_by=null
        where id=$1`, id, msg)
	return err
}

// Batch log -----------------------------------------------------------------

// Log mencatat konteks batch ke tabel audit sederhana agar alur ETL mudah ditelusuri.
func (r *Repo) Log(ctx context.Context, entity, op string, keyValues any, payload any, status, errMsg string) {
	keyJSON, _ := json.Marshal(keyValues)
	plJSON, _ := json.Marshal(payload)
	util.Debug.Printf("pivot: Log entity=%s op=%s status=%s", entity, op, status)
	_, _ = r.DB.Exec(ctx, `insert into _batch_log (entity, op, key_values, payload, status, error)
    values ($1,$2,$3::jsonb,$4::jsonb,$5,$6)`,
		entity, op, string(keyJSON), string(plJSON), status, errMsg)
}

// Helpers -------------------------------------------------------------------

// truncateSQL memotong log SQL agar output tidak terlalu panjang ketika dicatat.
func truncateSQL(sql string) string {
	// Bila query melebihi 64 karakter, tampilkan ringkasannya saja.
	if len(sql) > 64 {
		return sql[:64] + "..."
	}
	return sql
}

// nullIfEmpty mengubah string kosong menjadi NULL supaya kolom opsional
// di database tetap bersih (tidak terisi string kosong).
func nullIfEmpty(s string) any {
	if strings.TrimSpace(s) == "" {
		// String kosong diperlakukan sebagai NULL biar tidak mengotori kolom opsional.
		return nil
	}
	return s
}
