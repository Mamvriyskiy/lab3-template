package rollback

import (
	"encoding/json"
	"log"
	"time"
)

func startRetryWorker() {
	go func() {
		for {
			// Берём первый элемент из очереди (блокирующе, timeout 5 секунд)
			res, err := rdb.BLPop(ctx, 5*time.Second, "retry_queue").Result()
			if err != nil {
				if err == redis.Nil {
					continue // очередь пустая
				}
				log.Printf("Redis BLPop error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			var req RetryRequest
			if err := json.Unmarshal([]byte(res[1]), &req); err != nil {
				log.Printf("Failed to unmarshal request: %v", err)
				continue
			}

			// Выполняем запрос
			status, _, _, err := forwardRequest(nil, req.Method, req.URL, req.Headers, req.Body)
			if err != nil || status >= 400 {
				log.Printf("Retry failed for %s, re-enqueueing", req.URL)
				EnqueueRetry(req) // возвращаем в очередь
			} else {
				log.Printf("Retry succeeded for %s", req.URL)
			}

			time.Sleep(500 * time.Millisecond) // небольшой таймаут между попытками
		}
	}()
}
