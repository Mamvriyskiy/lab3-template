package circuitbreaker

import (
    "sync"
    "time"
)

type CircuitBreaker struct {
    FailureCount     int
    FailureThreshold int
    State            int
    LastFailureTime  time.Time
    RetryTimeout     time.Duration
    mu               sync.Mutex
}

const (
    Closed = iota
    Open
    HalfOpen
)
