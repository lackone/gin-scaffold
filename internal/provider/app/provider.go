package app

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
)

type AppProvider struct {
	BaseFolder string
}

func (a *AppProvider) Register(c core.IContainer) core.NewInstance {
	return NewApp
}

func (a *AppProvider) Boot(c core.IContainer) error {
	return nil
}

func (a *AppProvider) IsDefer() bool {
	return false
}

func (a *AppProvider) Params(c core.IContainer) []interface{} {
	return []interface{}{c, a.BaseFolder}
}

func (a *AppProvider) Name() string {
	return contract.AppKey
}
