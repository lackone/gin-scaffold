package gorm

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
)

type GormProvider struct {
}

func (g *GormProvider) Register(c core.IContainer) core.NewInstance {
	return NewGorm
}

func (g *GormProvider) Boot(c core.IContainer) error {
	return nil
}

func (g *GormProvider) IsDefer() bool {
	return true
}

func (g *GormProvider) Params(c core.IContainer) []interface{} {
	return []interface{}{c}
}

func (g *GormProvider) Name() string {
	return contract.GORMKey
}
