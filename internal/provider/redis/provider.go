package redis

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
)

type RedisProvider struct {
}

func (r *RedisProvider) Register(c core.IContainer) core.NewInstance {
	return NewRedisService
}

func (r *RedisProvider) Boot(c core.IContainer) error {
	return nil
}

func (r *RedisProvider) IsDefer() bool {
	return true
}

func (r *RedisProvider) Params(c core.IContainer) []interface{} {
	return []interface{}{c}
}

func (r *RedisProvider) Name() string {
	return contract.RedisKey
}
