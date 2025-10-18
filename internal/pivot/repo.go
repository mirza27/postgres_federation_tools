package pivot

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Repo, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil { return nil, err }
	return &Repo{DB: pool}, nil
}

func (r *Repo) Close() { r.DB.Close() }

func (r *Repo) EnsureSchema(ctx context.Context, schemaSQL string) error {
	_, err := r.DB.Exec(ctx, schemaSQL)
	return err
}

// KeyMap: upsert & lookup generik
func (r *Repo) LookupKey(ctx context.Context, mapName, srcKey string) (string, bool, error) {
	var tgt string
	err := r.DB.QueryRow(ctx,
		`select tgt_key::text from _keymap_generic where map_name=$1 and src_key=$2`,
		mapName, srcKey).Scan(&tgt)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return "", false, nil
		}
		return "", false, err
	}
	return tgt, true, nil
}

func (r *Repo) PutKeyIfAbsent(ctx context.Context, mapName, srcKey, genUUID string) (string, error) {
	// idempoten
	_, err := r.DB.Exec(ctx, `
		insert into _keymap_generic(map_name, src_key, tgt_key)
		values ($1,$2,$3::uuid)
		on conflict (map_name, src_key) do nothing`,
		mapName, srcKey, genUUID)
	if err != nil { return "", err }

	// baca balik
	tgt, _, err := r.LookupKey(ctx, mapName, srcKey)
	return tgt, err
}

// Exec Queue
type ExecItem struct {
	Entity string
	Op     string
	SQL    string
	Args   []interface{}
}

func (r *Repo) Enqueue(ctx context.Context, it ExecItem) error {
	argsJSON, _ := json.Marshal(it.Args)
	_, err := r.DB.Exec(ctx, `
		insert into _exec_queue (entity, op, sql_text, sql_args)
		values ($1,$2,$3,$4::jsonb)`,
		it.Entity, it.Op, it.SQL, string(argsJSON))
	return err
}

type Row struct {
	ID int64
	SQL string
	ArgsJSON []byte
	Entity string
	Op string
}

func (r *Repo) FetchPending(ctx context.Context, limit int) ([]Row, error) {
	rows, err := r.DB.Query(ctx, `
		select id, entity, op, sql_text, coalesce(sql_args,'[]'::jsonb)
		from _exec_queue
		where status='pending'
		order by id asc
		limit $1`, limit)
	if err != nil { return nil, err }
	defer rows.Close()

	var out []Row
	for rows.Next() {
		var id int64
		var sql string
		var args []byte
		var entity, op string
		if err := rows.Scan(&id, &entity, &op, &sql, &args); err != nil { return nil, err }
		out = append(out, Row{ID: id, Entity: entity, Op: op, SQL: sql, ArgsJSON: args})
	}
	return out, nil
}

func (r *Repo) MarkDone(ctx context.Context, id int64) error {
	_, err := r.DB.Exec(ctx, `update _exec_queue set status='done' where id=$1`, id)
	return err
}

func (r *Repo) MarkError(ctx context.Context, id int64, msg string) error {
	_, err := r.DB.Exec(ctx, `update _exec_queue set status='error', error=$2 where id=$1`, id, msg)
	return err
}

// Batch log
func (r *Repo) Log(ctx context.Context, entity, op string, keyValues any, payload any, status, errMsg string) {
	keyJSON, _ := json.Marshal(keyValues)
	plJSON, _ := json.Marshal(payload)
	_, _ = r.DB.Exec(ctx, `insert into _batch_log (entity, op, key_values, payload, status, error)
	values ($1,$2,$3::jsonb,$4::jsonb,$5,$6)`,
		entity, op, string(keyJSON), string(plJSON), status, errMsg)
}
