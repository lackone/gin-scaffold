package http

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/lackone/gin-scaffold/internal/middleware"
	"github.com/lackone/gin-scaffold/routes"
	"net/http"
)

func NewHttpEngine(container core.IContainer) (*core.Engine, error) {
	//gin.SetMode(gin.ReleaseMode)
	global.Engine = core.Default()
	global.Engine.SetContainer(container)

	global.Engine.Use(middleware.SetEngine())
	global.Engine.Use(middleware.Translations())

	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	global.Engine.StaticFS("/uploads", http.Dir(configService.GetString("app.upload_save_path")))

	routes.ApiRoutes(global.Engine)
	routes.WebRoutes(global.Engine)

	return global.Engine, nil
}
