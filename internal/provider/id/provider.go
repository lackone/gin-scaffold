package id

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
)

type IDProvider struct {
}

func (i *IDProvider) Register(c core.IContainer) core.NewInstance {
	return NewIDService
}

func (i *IDProvider) Boot(c core.IContainer) error {
	return nil
}

func (i *IDProvider) IsDefer() bool {
	return false
}

func (i *IDProvider) Params(c core.IContainer) []interface{} {
	return []interface{}{}
}

func (i *IDProvider) Name() string {
	return contract.IDKey
}
