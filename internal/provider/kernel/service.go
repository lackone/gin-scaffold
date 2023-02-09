package kernel

import (
	"github.com/lackone/gin-scaffold/internal/core"
	"net/http"
)

type Kernel struct {
	engine *core.Engine
}

func NewKernel(params ...interface{}) (interface{}, error) {
	engine := params[0].(*core.Engine)
	return &Kernel{engine: engine}, nil
}

func (k *Kernel) HttpEngine() http.Handler {
	return k.engine
}
