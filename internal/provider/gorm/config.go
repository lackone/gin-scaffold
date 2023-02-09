package gorm

import (
	"context"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
)

// GetDefaultConfig 读取database.yaml里公共的配置
func GetDefaultConfig(c core.IContainer) *contract.DBConfig {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	logService := c.MustMake(contract.LogKey).(contract.Log)
	config := &contract.DBConfig{}
	err := configService.Load("database", config)
	if err != nil {
		logService.Error(context.Background(), "parse database config error", nil)
		return nil
	}
	return config
}

// WithConfigPath 加载某个配置
func WithConfigPath(key string) contract.DBOption {
	return func(container core.IContainer, config *contract.DBConfig) error {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		if err := configService.Load(key, config); err != nil {
			return err
		}
		return nil
	}
}

// WithGormConfig 表示自行配置Gorm的配置信息
func WithGormConfig(f func(options *contract.DBConfig)) contract.DBOption {
	return func(container core.IContainer, config *contract.DBConfig) error {
		f(config)
		return nil
	}
}

// WithDryRun 设置空跑模式
func WithDryRun() contract.DBOption {
	return func(container core.IContainer, config *contract.DBConfig) error {
		config.DryRun = true
		return nil
	}
}

// WithFullSaveAssociations 设置保存时候关联
func WithFullSaveAssociations() contract.DBOption {
	return func(container core.IContainer, config *contract.DBConfig) error {
		config.FullSaveAssociations = true
		return nil
	}
}
