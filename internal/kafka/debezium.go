package kafka

// Struktur minimal untuk envelope Debezium JSON
// (disederhanakan—pakai gjson untuk fleksibilitas).
type Event struct {
	Op   string // c,u,d,r
	// payload "after"/"before" akan diambil via gjson dari raw []byte
	Source string // ex: "public.users" / topic-key if needed
	Topic  string // kafka topic
	Key    []byte
	Value  []byte
}
