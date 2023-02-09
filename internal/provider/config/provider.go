package config

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"path/filepath"
)

type ConfigProvider struct{}

func (cp *ConfigProvider) Register(c core.IContainer) core.NewInstance {
	return NewConfig
}

func (cp *ConfigProvider) Boot(c core.IContainer) error {
	return nil
}

func (cp *ConfigProvider) IsDefer() bool {
	return false
}

func (cp *ConfigProvider) Params(c core.IContainer) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()
	configFolder := appService.ConfigFolder()
	envConfigFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envConfigFolder, envService.All()}
}

func (cp *ConfigProvider) Name() string {
	return contract.ConfigKey
}
