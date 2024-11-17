package shared

import (
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRateLimiter(t *testing.T) {
	middleware := NewFixedWindowRateLimiter(Config{
		RequestPerTimeFrame: 5,
		TimeFrame:           1,
		Enabled:             true,
	}, nil)

	t.Run("should run when limit is not exceeded", func(t *testing.T) {
		// Arrange
		ip := "localhost"

		// Act
		allowed, _ := middleware.Allow(ip)

		// Assert
		assert.True(t, allowed)

	})

	t.Run("should return false when limit is exceeded", func(t *testing.T) {
		// Arrange
		ip := "localhost"
		for i := 0; i < 5; i++ {
			middleware.Allow(ip)
		}

		// Act
		allowed, duration := middleware.Allow(ip)

		// Assert
		assert.False(t, allowed)
		assert.Equal(t, middleware.window, duration)
	})
}
