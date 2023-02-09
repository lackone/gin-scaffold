package contract

import (
	"fmt"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/redis/go-redis/v9"
)

const RedisKey = "app:redis"

type RedisOption func(container core.IContainer, config *RedisConfig) error

type IRedis interface {
	GetClient(option ...RedisOption) (*redis.Client, error)
}

type RedisConfig struct {
	*redis.Options
}

// UniqKey 用来唯一标识一个RedisConfig配置
func (config *RedisConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}
