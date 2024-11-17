package shared

import (
	"net/http"
)

func (limiter *FixedWindowRateLimiter) RateLimiterMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if limiter.enabled {
				allowed, duration := limiter.Allow(r.RemoteAddr)
				if !allowed {
					limiter.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
					w.Header().Set("Retry-After", duration.String())
					http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
