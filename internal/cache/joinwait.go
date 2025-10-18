package cache

import (
	"sync"
	"time"
)

type key = string

type payload struct {
	Value   []byte
	Arrived time.Time
}

type JoinWait struct {
	mu   sync.Mutex
	data map[key]map[string]payload // joinKey -> topic -> payload
	ttl  time.Duration
}

func NewJoinWait(ttlSec int) *JoinWait {
	return &JoinWait{
		data: make(map[key]map[string]payload),
		ttl:  time.Duration(ttlSec) * time.Second,
	}
}

// Add stores an event for joinKey & topic.
func (j *JoinWait) Add(joinKey, topic string, value []byte) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if _, ok := j.data[joinKey]; !ok {
		j.data[joinKey] = make(map[string]payload)
	}
	j.data[joinKey][topic] = payload{Value: value, Arrived: time.Now()}
	j.gcLocked()
}

// GetSet returns all topic payloads for a joinKey if completeTopics are all present.
func (j *JoinWait) GetSet(joinKey string, completeTopics []string) (map[string][]byte, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	m, ok := j.data[joinKey]
	if !ok { return nil, false }
	for _, t := range completeTopics {
		if _, exists := m[t]; !exists {
			return nil, false
		}
	}
	out := map[string][]byte{}
	for k, v := range m { out[k] = v.Value }
	delete(j.data, joinKey) // consume
	return out, true
}

func (j *JoinWait) gcLocked() {
	now := time.Now()
	for k, topics := range j.data {
		for _, p := range topics {
			if now.Sub(p.Arrived) > j.ttl {
				delete(j.data, k)
				break
			}
		}
	}
}
