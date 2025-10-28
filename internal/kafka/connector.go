package kafka

import (
	"context"
	"time"

	"db_migrate_server/internal/util"

	kfk "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kfk.Reader
}

func NewConsumer(brokers []string, groupID string, topics []string) *Consumer {
	util.Info.Printf("kafka: NewConsumer brokers=%v groupID=%s topics=%v", brokers, groupID, topics)
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

// for testing
func (c *Consumer) ListenSync(ctx context.Context, handler func(Event) error) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		value, src, err := DecodeEventValue(m.Value)
		if err != nil {
			util.Error.Printf("kafka: decode message err=%v topic=%s partition=%d offset=%d", err, m.Topic, m.Partition, m.Offset)
			continue
		}
		ev := Event{
			Topic:  string(m.Topic),
			Key:    m.Key,
			Value:  value,
			Source: src,
		}
		if value != nil {
			ev.Op = value.Op
		}
		if err := handler(ev); err != nil {
			return err
		}
	}
}

func (c *Consumer) Listen(ctx context.Context, out chan<- Event) error {
	defer close(out)
	util.Info.Println("kafka: Listen loop started")
	for {
		m, err := c.reader.ReadMessage(ctx)

		if err != nil {
			if ctx.Err() != nil {
				util.Warn.Printf("kafka: listen canceled err=%v", err)
			} else {
				util.Error.Printf("kafka: read message err=%v", err)
			}
			return err
		}

		util.Debug.Printf("kafka: received message topic=%s partition=%d offset=%d key=%s len=%d",
			m.Topic, m.Partition, m.Offset, string(m.Key), len(m.Value))

		value, src, err := DecodeEventValue(m.Value)
		if err != nil {
			util.Error.Printf("kafka: decode message err=%v topic=%s partition=%d offset=%d", err, m.Topic, m.Partition, m.Offset)
			continue
		}

		ev := Event{
			Topic:  string(m.Topic),
			Key:    m.Key,
			Value:  value,
			Source: src,
		}
		if value != nil {
			ev.Op = value.Op
		}
		util.Debug.Printf("kafka: event parsed op=%s topic=%s", ev.Op, ev.Topic)
		out <- ev
	}
}
