package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/gin-scaffold/internal/global"
)

func SetEngine() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("engine", global.Engine)
		c.Next()
	}
}
