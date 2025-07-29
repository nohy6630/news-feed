package dto

type KafkaMessage struct {
	UserID    int64 `json:"user_id"`
	PostID    int64 `json:"post_id"`
	Timestamp int64 `json:"timestamp"`
}
