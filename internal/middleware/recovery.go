package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		engine := c.MustGet("engine").(*core.Engine)
		container := engine.GetContainer()
		logService := container.MustMake(contract.LogKey).(contract.Log)

		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logService.Error(c, c.Request.URL.Path, map[string]interface{}{
						"error":   err,
						"request": string(httpRequest),
					})
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logService.Error(c, "[Recovery from panic]", map[string]interface{}{
					"error":   err,
					"request": string(httpRequest),
					"stack":   string(debug.Stack()),
					"path":    c.Request.URL.Path,
				})
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
