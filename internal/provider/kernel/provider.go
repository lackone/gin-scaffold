package kernel

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
)

type KernelProvider struct {
	HttpEngine *core.Engine
}

func (k *KernelProvider) Register(c core.IContainer) core.NewInstance {
	return NewKernel
}

func (k *KernelProvider) Boot(c core.IContainer) error {
	if k.HttpEngine == nil {
		k.HttpEngine = core.Default()
	}
	k.HttpEngine.SetContainer(c)
	return nil
}

func (k *KernelProvider) IsDefer() bool {
	return false
}

func (k *KernelProvider) Params(c core.IContainer) []interface{} {
	return []interface{}{k.HttpEngine}
}

func (k *KernelProvider) Name() string {
	return contract.KernelKey
}
