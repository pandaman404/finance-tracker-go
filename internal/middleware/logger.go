package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pandaman404/finance-tracker-go/pkg/logger"
)

func Logger() gin.HandlerFunc {
	log := logger.New()

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
			"user_id", c.GetString(UserIDKey),
		)
	}
}
