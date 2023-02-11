package http

import (
	"github.com/gin-contrib/cors"
	_ "github.com/lackone/gin-scaffold/docs/swagger" //千万别忘了导入生成的docs
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/lackone/gin-scaffold/internal/middleware"
	"github.com/lackone/gin-scaffold/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

func NewHttpEngine(container core.IContainer) (*core.Engine, error) {
	//gin.SetMode(gin.ReleaseMode)
	global.Engine = core.Default()
	global.Engine.SetContainer(container)

	global.Engine.Use(middleware.SetEngine())
	global.Engine.Use(middleware.Translations())

	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	//cors中间件
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = configService.GetBool("cors.allow_all_origins")
	corsConfig.AllowOrigins = configService.GetStringSlice("cors.allow_origins")
	corsConfig.AllowMethods = configService.GetStringSlice("cors.allow_methods")
	corsConfig.AllowHeaders = configService.GetStringSlice("cors.allow_headers")
	corsConfig.AllowCredentials = configService.GetBool("cors.allow_credentials")
	corsConfig.ExposeHeaders = configService.GetStringSlice("cors.expose_headers")
	corsMaxAge, _ := time.ParseDuration(configService.GetString("cors.max_age"))
	corsConfig.MaxAge = corsMaxAge
	corsConfig.AllowWildcard = configService.GetBool("cors.allow_wildcard")
	corsConfig.AllowBrowserExtensions = configService.GetBool("cors.allow_browser_extensions")
	corsConfig.AllowWebSockets = configService.GetBool("cors.allow_web_sockets")
	corsConfig.AllowFiles = configService.GetBool("cors.allow_files")
	global.Engine.Use(cors.New(corsConfig))

	// 如果配置了swagger，则显示swagger的中间件
	if configService.GetBool("app.swagger") == true {
		global.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	global.Engine.StaticFS("/uploads", http.Dir(configService.GetString("app.upload_save_path")))

	routes.ApiRoutes(global.Engine)
	routes.WebRoutes(global.Engine)

	return global.Engine, nil
}
