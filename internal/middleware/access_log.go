package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"time"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		engine := c.MustGet("engine").(*core.Engine)
		container := engine.GetContainer()
		logService := container.MustMake(contract.LogKey).(contract.Log)

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)

		logService.Trace(c, path, map[string]interface{}{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"errors":     c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"cost":       cost,
		})
	}
}
