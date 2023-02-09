package env

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
)

type EnvProvider struct {
	Folder string
}

func (e *EnvProvider) Register(c core.IContainer) core.NewInstance {
	return NewEnv
}

func (e *EnvProvider) Boot(c core.IContainer) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	e.Folder = app.BaseFolder()
	return nil
}

func (e *EnvProvider) IsDefer() bool {
	return false
}

func (e *EnvProvider) Params(c core.IContainer) []interface{} {
	return []interface{}{e.Folder}
}

func (e *EnvProvider) Name() string {
	return contract.EnvKey
}
