package pivot

import "github.com/google/uuid"

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
