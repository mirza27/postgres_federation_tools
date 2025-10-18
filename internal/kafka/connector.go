package kafka

import (
	"context"
	"time"

	"github.com/tidwall/gjson"

	kfk "github.com/segmentio/kafka-go"
	// "github.com/yourorg/mapping-engine/internal/util"
)

type Consumer struct {
	reader *kfk.Reader
}

func NewConsumer(brokers []string, groupID string, topics []string) *Consumer {
	r := kfk.NewReader(kfk.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		GroupTopics:    topics,
		MinBytes:       1,
		MaxBytes:       10e6,
		MaxWait:        500 * time.Millisecond,
		CommitInterval: time.Second,
	})
	return &Consumer{reader: r}
}


func (c *Consumer) Close() error { return c.reader.Close() }

func (c *Consumer) Listen(ctx context.Context, out chan<- Event) error {
	defer close(out)
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		ev := Event{
			Topic: string(m.Topic),
			Key:   m.Key,
			Value: m.Value,
		}
		out <- ev
	}
}

func ExtractOpAndSource(ev *Event) (op, src string) {
	// gunakan gjson untuk membaca nilai
	v := gjson.ParseBytes(ev.Value)
	op = v.Get("op").String()
	src = v.Get("source.schema").String() + "." + v.Get("source.table").String()
	if src == "." { // fallback
		src = v.Get("payload.source.schema").String() + "." + v.Get("payload.source.table").String()
	}
	return
}