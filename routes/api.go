package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/gin-scaffold/app/http/controllers/test"
	_ "github.com/lackone/gin-scaffold/docs/swagger" //千万别忘了导入生成的docs
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ApiRoutes(r *core.Engine) {
	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	// 如果配置了swagger，则显示swagger的中间件
	if configService.GetBool("app.swagger") == true {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"test": "test"})
	})

	test := &test.Test{}
	r.GET("/test_orm", test.TestORM)
	r.GET("/test_swag", test.TestSwag)
	r.GET("/test_swag2", test.TestSwag2)
}
