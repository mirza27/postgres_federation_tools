package pivot

import "github.com/google/uuid"

// Row merepresentasikan baris yang dibaca executor dari _exec_queue. Struktur
// ini menampung payload SQL mentah, argumen dalam bentuk JSON, serta metadata
// yang memberitahu apakah eksekusi membutuhkan keymap atau kolom RETURNING.
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
}
