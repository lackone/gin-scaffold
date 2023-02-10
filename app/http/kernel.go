package http

import (
	_ "github.com/lackone/gin-scaffold/docs/swagger" //千万别忘了导入生成的docs
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/lackone/gin-scaffold/internal/middleware"
	"github.com/lackone/gin-scaffold/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func NewHttpEngine(container core.IContainer) (*core.Engine, error) {
	//gin.SetMode(gin.ReleaseMode)
	global.Engine = core.Default()
	global.Engine.SetContainer(container)

	global.Engine.Use(middleware.SetEngine())
	global.Engine.Use(middleware.Translations())

	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	// 如果配置了swagger，则显示swagger的中间件
	if configService.GetBool("app.swagger") == true {
		global.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	global.Engine.StaticFS("/uploads", http.Dir(configService.GetString("app.upload_save_path")))

	routes.ApiRoutes(global.Engine)
	routes.WebRoutes(global.Engine)

	return global.Engine, nil
}
