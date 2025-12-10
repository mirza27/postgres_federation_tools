package pivot

import "github.com/google/uuid"

// SplitRow merepresentasikan baris split yang perlu dieksekusi setelah queue utama.
type SplitRow struct {
	ID            int64
	QueueID       uuid.UUID
	SQL           string
	ArgsJSON      []byte
	Returning     []string
	KeymapPayload []byte
	Status        string
}
