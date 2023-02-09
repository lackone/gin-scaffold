package redis

import (
	"context"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// GetDefaultConfig 读取redis.yaml里公共的配置
func GetDefaultConfig(c core.IContainer) *contract.RedisConfig {
	logService := c.MustMake(contract.LogKey).(contract.Log)
	config := &contract.RedisConfig{Options: &redis.Options{}}
	opt := WithConfigPath("redis")
	err := opt(c, config)
	if err != nil {
		logService.Error(context.Background(), "parse cache config error", nil)
		return nil
	}
	return config
}

// WithConfigPath 加载某个配置
func WithConfigPath(key string) contract.RedisOption {
	return func(container core.IContainer, config *contract.RedisConfig) error {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		conf := configService.GetStringMapString(key)

		// 读取config配置
		if host, ok := conf["host"]; ok {
			if port, ok1 := conf["port"]; ok1 {
				config.Addr = host + ":" + port
			}
		}

		if db, ok := conf["db"]; ok {
			t, err := strconv.Atoi(db)
			if err != nil {
				return err
			}
			config.DB = t
		}

		if username, ok := conf["username"]; ok {
			config.Username = username
		}

		if password, ok := conf["password"]; ok {
			config.Password = password
		}

		if timeout, ok := conf["timeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.DialTimeout = t
		}

		if timeout, ok := conf["read_timeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.ReadTimeout = t
		}

		if timeout, ok := conf["write_timeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.WriteTimeout = t
		}

		if cnt, ok := conf["conn_min_idle"]; ok {
			t, err := strconv.Atoi(cnt)
			if err != nil {
				return err
			}
			config.MinIdleConns = t
		}

		if cnt, ok := conf["conn_max_idle"]; ok {
			t, err := strconv.Atoi(cnt)
			if err != nil {
				return err
			}
			config.MaxIdleConns = t
		}

		if max, ok := conf["conn_pool_size"]; ok {
			t, err := strconv.Atoi(max)
			if err != nil {
				return err
			}
			config.PoolSize = t
		}

		if timeout, ok := conf["conn_max_lifetime"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.ConnMaxLifetime = t
		}

		if timeout, ok := conf["conn_max_idletime"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.ConnMaxIdleTime = t
		}

		return nil
	}
}

// WithRedisConfig 表示自行配置redis的配置信息
func WithRedisConfig(f func(options *contract.RedisConfig)) contract.RedisOption {
	return func(container core.IContainer, config *contract.RedisConfig) error {
		f(config)
		return nil
	}
}
