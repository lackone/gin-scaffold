package http

import (
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/lackone/gin-scaffold/routes"
)

func NewHttpEngine(container core.IContainer) (*core.Engine, error) {
	//gin.SetMode(gin.ReleaseMode)
	global.Engine = core.Default()
	global.Engine.SetContainer(container)

	routes.ApiRoutes(global.Engine)
	routes.WebRoutes(global.Engine)

	return global.Engine, nil
}
