package accessor

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type RedisAccessor struct {
	Client *redis.Client
}

func NewRedisAccessor(addr string) *RedisAccessor {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisAccessor{Client: client}
}

func (ra *RedisAccessor) Ping(ctx context.Context) error {
	return ra.Client.Ping(ctx).Err()
}

// Sorted Set을 이용해 userID별로 postID를 저장 (만료 시간 설정 포함)
func (ra *RedisAccessor) AddPostToUserFeed(ctx context.Context, userID string, postID string, ttlSeconds int) error {
	key := "feed:" + userID
	pipe := ra.Client.Pipeline()
	now := time.Now().Unix()
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: postID})
	pipe.Expire(ctx, key, time.Duration(ttlSeconds)*time.Second)
	_, err := pipe.Exec(ctx)
	return err
}

// userID의 feed에서 최신순으로 최대 count개 조회
func (ra *RedisAccessor) GetUserFeed(ctx context.Context, userID string, count int64) ([]string, error) {
	key := "feed:" + userID
	return ra.Client.ZRevRange(ctx, key, 0, count-1).Result()
}

var (
	onceRedis     sync.Once
	redisInstance *RedisAccessor
)

func GetRedisAccessor() *RedisAccessor {
	onceRedis.Do(func() {
		redisInstance = NewRedisAccessor("localhost:6379")
	})
	return redisInstance
}
