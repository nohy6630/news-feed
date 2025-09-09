package accessor

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"news-feed/config"
	"sync"
)

type KafkaAccessor struct {
	Client *kgo.Client
}

var (
	kafkaOnce     sync.Once
	kafkaInstance *KafkaAccessor
)

func NewKafkaAccessor(brokers []string, topic string) (*KafkaAccessor, error) {
	var opts []kgo.Opt
	opts = append(opts, kgo.SeedBrokers(brokers...))
	opts = append(opts, kgo.ConsumerGroup("group"))
	opts = append(opts, kgo.ConsumeTopics(topic))

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &KafkaAccessor{Client: cl}, nil
}

func GetKafkaAccessor() (*KafkaAccessor, error) {
	var err error
	kafkaOnce.Do(func() {
		addr := fmt.Sprintf("%s:9092", config.GetKafkaAddress())
		kafkaInstance, err = NewKafkaAccessor([]string{addr}, "feed")
	})
	return kafkaInstance, err
}

// Produce 함수 (동기)
func (ka *KafkaAccessor) ProduceSync(ctx context.Context, topic string, value []byte) error {
	fmt.Printf("produce sync topic=%s value=%s\n", topic, value)
	record := &kgo.Record{Topic: topic, Value: value}
	err := ka.Client.ProduceSync(ctx, record).FirstErr()
	fmt.Printf("produce sync topic=%s value=%s\n", topic, value)
	return err
}

// Consume 함수 (동기)
func (ka *KafkaAccessor) ConsumeSync(ctx context.Context, handler func(*kgo.Record)) error {
	for {
		fetches := ka.Client.PollFetches(ctx)
		for _, ferr := range fetches.Errors() {
			if ferr.Err != nil {
				return ferr.Err
			}
		}
		fetches.EachRecord(func(record *kgo.Record) {
			handler(record)
		})
	}
}
