package core

import (
	"github.com/gin-gonic/gin"
)

type Engine struct {
	*gin.Engine
	container IContainer
}

func New() *Engine {
	return &Engine{
		Engine: gin.New(),
	}
}

func Default() *Engine {
	return &Engine{
		Engine: gin.Default(),
	}
}

func (e *Engine) Bind(provider IServiceProvider) error {
	return e.container.Bind(provider)
}

func (e *Engine) IsBind(key string) bool {
	return e.container.IsBind(key)
}

func (e *Engine) SetContainer(container IContainer) {
	e.container = container
}

func (e *Engine) GetContainer() IContainer {
	return e.container
}

func (e *Engine) Make(key string) (interface{}, error) {
	return e.container.Make(key)
}

func (e *Engine) MustMake(key string) interface{} {
	return e.container.MustMake(key)
}

func (e *Engine) MakeNew(key string, params []interface{}) (interface{}, error) {
	return e.container.MakeNew(key, params)
}
