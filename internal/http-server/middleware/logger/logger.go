package logger

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

func Middleware(log *slog.Logger) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		TimeStamp := time.Now()
		Latency := TimeStamp.Sub(start)

		ClientIP := c.ClientIP()
		Method := c.Request.Method
		StatusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		Path := path

		log.Info("[SLOG]",
			slog.String("method", Method),
			slog.String("path", Path),
			slog.String("clientIP", ClientIP),
			slog.String("time", Latency.String()),
			slog.Int("status", StatusCode),
		)
	}

	return fn
}
