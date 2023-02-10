package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/gin-scaffold/app/http/controllers/test"
	"github.com/lackone/gin-scaffold/internal/core"
)

func ApiRoutes(r *core.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"test": "test"})
	})

	test := &test.Test{}
	r.GET("/test_orm", test.TestORM)
	r.GET("/test_swag", test.TestSwag)
	r.GET("/test_swag2", test.TestSwag2)
}
