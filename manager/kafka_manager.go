package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"news-feed/accessor"
	"news-feed/dto"
	"sync"
)

type KafkaManager struct {
	Accessor *accessor.KafkaAccessor
}

var (
	kafkaManagerOnce     sync.Once
	kafkaManagerInstance *KafkaManager
)

func NewKafkaManager() (*KafkaManager, error) {
	ka, err := accessor.GetKafkaAccessor()
	if err != nil {
		return nil, err
	}
	return &KafkaManager{Accessor: ka}, nil
}

func GetKafkaManager() (*KafkaManager, error) {
	var err error
	kafkaManagerOnce.Do(func() {
		kafkaManagerInstance, err = NewKafkaManager()
	})
	return kafkaManagerInstance, err
}

// Produce 메시지 전송 (dto.KafkaMessage를 받아 직렬화 후 전송)
func (km *KafkaManager) Produce(ctx context.Context, msg dto.KafkaMessage) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return km.Accessor.ProduceSync(ctx, "feed", value)
}

// Consume 메시지 핸들러 (accessor의 ConsumeSync 사용)
func (km *KafkaManager) Consume(ctx context.Context) error {
	return km.Accessor.ConsumeSync(ctx, handleMessage)
}

func handleMessage(msg *kgo.Record) {
	fmt.Printf("Received message: %s\n", string(msg.Value))
	var kafkaMsg dto.KafkaMessage
	if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
		fmt.Printf("Error unmarshalling message: %v\n", err)
		return
	}
	redisAccessor := accessor.GetRedisAccessor()
	userID := fmt.Sprintf("%d", kafkaMsg.UserID)
	postID := fmt.Sprintf("%d", kafkaMsg.PostID)
	// 예시: TTL 24시간(86400초)
	err := redisAccessor.AddPostToUserFeed(context.Background(), userID, postID, kafkaMsg.Timestamp, 86400)
	if err != nil {
		fmt.Printf("Error adding post to user feed: %v\n", err)
		return
	}
}
