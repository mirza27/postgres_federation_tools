package pivot

import "github.com/google/uuid"

// NeedJoinItem merepresentasikan baris pending join di tabel _need_join.
type NeedJoinItem struct {
	ID            int64
	QueueID       uuid.UUID
	Entity        string
	Op            string
	JoinKey       string
	JoinTopic     string
	JoinSourceKey string
	JoinPayload   []byte
	JoinFields    []byte
	Status        string
}
