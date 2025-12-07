package pivot

import "github.com/google/uuid"

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
