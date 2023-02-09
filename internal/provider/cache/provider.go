package cache

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/provider/cache/services"
	"strings"
)

type CacheProvider struct {
	Driver string // 驱动
}

func (l *CacheProvider) Register(c core.IContainer) core.NewInstance {
	if l.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			return services.NewMemoryService
		}

		cs := tcs.(contract.Config)
		l.Driver = strings.ToLower(cs.GetString("cache.driver"))
	}

	switch l.Driver {
	case "redis":
		return services.NewRedisService
	case "memory":
		return services.NewMemoryService
	default:
		return services.NewMemoryService
	}
}

func (l *CacheProvider) Boot(c core.IContainer) error {
	return nil
}

func (l *CacheProvider) IsDefer() bool {
	return true
}

func (l *CacheProvider) Params(c core.IContainer) []interface{} {
	return []interface{}{c}
}

func (l *CacheProvider) Name() string {
	return contract.CacheKey
}
