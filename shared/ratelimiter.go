package shared

import (
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type RateLimiter interface {
	Allow(key string) (bool, time.Duration)
	RateLimiterMiddleware() func(http.Handler) http.Handler
}

type Config struct {
	RequestPerTimeFrame int
	TimeFrame           time.Duration
	Enabled             bool
}

type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients map[string]int
	limit   int
	window  time.Duration
	enabled bool
	logger  *zap.SugaredLogger
}

func NewFixedWindowRateLimiter(config Config, logger *zap.SugaredLogger) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   config.RequestPerTimeFrame,
		window:  config.TimeFrame,
		enabled: config.Enabled,
		logger:  logger,
	}
}

func (limiter *FixedWindowRateLimiter) Allow(key string) (bool, time.Duration) {
	limiter.Lock()
	defer limiter.Unlock()

	if !limiter.enabled {
		return true, 0
	}

	_, exist := limiter.clients[key]
	if !exist {
		go limiter.reset(key)
	}

	limiter.clients[key]++
	if limiter.clients[key] > limiter.limit {
		return false, limiter.window
	}

	return true, 0
}

func (limiter *FixedWindowRateLimiter) reset(key string) {
	time.Sleep(limiter.window)
	limiter.Lock()
	defer limiter.Unlock()

	delete(limiter.clients, key)
}
