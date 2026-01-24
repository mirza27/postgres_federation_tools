package pivot

import (
	"time"

	"github.com/google/uuid"
)

// Row merepresentasikan baris yang dibaca executor dari _exec_queue. Struktur
// ini menampung payload SQL mentah, argumen dalam bentuk JSON, serta metadata
// yang memberitahu apakah eksekusi membutuhkan keymap, kolom RETURNING, atau
// merupakan pekerjaan split tambahan.
type Row struct {
	ID         int64
	QueueID    uuid.UUID
	Entity     string
	Op         string
	SQL        string
	ArgsJSON   []byte
	NeedKeymap bool
	KeymapID   *int64
	Returning  []string
	Status     string
	IsSplit    bool
	SplitName  string
}

// NeedJoinItem merepresentasikan baris pending join di tabel _need_join.
type NeedJoinItem struct {
	ID          int64
	QueueID     uuid.UUID
	Entity      string
	Op          string
	JoinKey     string
	JoinTopic   string
	JoinPayload []byte
	JoinFields  []byte
	Attempts    int
	NextAttempt time.Time
	Status      string
}

// SplitRow merepresentasikan baris split yang perlu dieksekusi setelah queue utama.
// _exec_split table representation
type SplitRow struct {
	ID            int64
	QueueID       uuid.UUID
	SQL           string
	ArgsJSON      []byte
	Returning     []string
	KeymapPayload []byte
	Status        string
}

// KeymapRequest merepresentasikan permintaan padanan kunci yang harus dicatat
// di database pivot ketika target menghasilkan nilai kunci baru. Struktur ini
// membantu checker dan executor saling bertukar konteks mengenai nama map,
// sumber data, serta kolom target yang nantinya akan diisi oleh hasil RETURNING.
type KeymapRequest struct {
	MapName   string `json:"map_name"`
	SrcKey    string `json:"src_key"`
	SrcColumn string `json:"src_column,omitempty"`
	TgtColumn string `json:"tgt_column,omitempty"`
	SrcTable  string `json:"src_table,omitempty"`
	TgtTable  string `json:"tgt_table,omitempty"`
}

// default config in pivot database configuration table
type Configuration struct {
	ConfigKey   string
	ConfigValue string
	UpdatedAt   time.Time
}

// MODEL HELPER

// CheckerItem adalah representasi kerja untuk worker checker. Di dalamnya
// tersimpan status terbaru baris _exec_queue yang sedang menunggu padanan
// key, termasuk payload permintaan keymap yang perlu diterjemahkan menjadi
// catatan baru pada tabel _keymap_generic.
type CheckerItem struct {
	ID            int64
	QueueID       uuid.UUID
	Entity        string
	NeedKeymap    bool
	KeymapID      *int64
	Status        string
	KeymapPayload []byte
}

// ExecItem menggambarkan satu pekerjaan yang akan dieksekusi di database
// target. Parser menuliskan informasi ini ke tabel _exec_queue dan memasukkan
// konteks tambahan seperti kebutuhan keymap atau kolom RETURNING sehingga
// worker lain tahu harus melakukan apa terhadap baris tersebut.
type ExecItem struct {
	QueueID    uuid.UUID
	Entity     string
	Op         string
	SQL        string
	Args       []interface{}
	NeedKeymap bool
	Keymap     *KeymapRequest
	Returning  []string

	// split flags untuk menandai pekerjaan tambahan
	IsSplit   bool
	SplitName string
}
