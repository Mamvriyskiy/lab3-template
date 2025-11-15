package worker

import (
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"io"
	"time"
	"github.com/redis/go-redis/v9"
	myRedis "github.com/Mamvriyskiy/lab3-template/src/gateway/rollback"
)

func DoRequest(method, url string, headers map[string]string, body []byte) (int, []byte, error) {
    req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
    if err != nil {
        return 0, nil, err
    }

    for k, v := range headers {
        req.Header.Set(k, v)
    }

    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil {
        return 0, nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    return resp.StatusCode, respBody, err
}


func StartRetryWorker() {
	go func() {
		for {
			// Берём первый элемент из очереди (блокирующе, timeout 5 секунд)
			res, err := myRedis.Rdb.BLPop(myRedis.Ctx, 5*time.Second, "retry_queue").Result()
			if err != nil {
				if err == redis.Nil {
					continue // очередь пустая
				}
				log.Printf("Redis BLPop error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			var req myRedis.RetryRequest
			if err := json.Unmarshal([]byte(res[1]), &req); err != nil {
				log.Printf("Failed to unmarshal request: %v", err)
				continue
			}

			// Выполняем запрос без gin.Context
			status, _, err := DoRequest(req.Method, req.URL, req.Headers, req.Body)
			if err != nil || status >= 400 {
				log.Printf("Retry failed for %s, re-enqueueing", req.URL)
				myRedis.EnqueueRetry(req) // возвращаем в очередь
			} else {
				log.Printf("Retry succeeded for %s", req.URL)
			}

			time.Sleep(500 * time.Millisecond) // небольшой таймаут между попытками
		}
	}()
}

