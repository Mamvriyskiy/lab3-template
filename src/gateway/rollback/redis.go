package rollback

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type RetryRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

var Ctx = context.Background()
var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // адрес Redis
		Password: "",               // если есть пароль
		DB:       0,                // база по умолчанию
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected")
}

func EnqueueRetry(req RetryRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	// Добавляем в конец списка "retry_queue"
	return Rdb.RPush(Ctx, "retry_queue", data).Err()
}
