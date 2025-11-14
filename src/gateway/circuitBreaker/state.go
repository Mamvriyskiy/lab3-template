package circuitbreaker

import (
	"github.com/gin-gonic/gin"
)

func (cb *CircuitBreaker) recordFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.FailureCount++
    cb.LastFailureTime = time.Now()

    if cb.FailureCount >= cb.FailureThreshold {
        cb.State = Open
    }
}

func (cb *CircuitBreaker) recordSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.State = Closed
    cb.FailureCount = 0
}

func (cb *CircuitBreaker) Execute(
    operation func() error,
    fallback func(c *gin.Context),
    c *gin.Context,
) {
    cb.mu.Lock()

    switch cb.State {
    case Open:
        if time.Since(cb.LastFailureTime) > cb.RetryTimeout {
            cb.State = HalfOpen
        } else {
            // сразу отдаём fallback
            cb.mu.Unlock()
            fallback(c)
            return
        }
    }
    cb.mu.Unlock()

    // Выполняем реальный запрос
    err := operation()
    if err != nil {
        cb.recordFailure()
        fallback(c)
        return
    }

    cb.recordSuccess()
}

