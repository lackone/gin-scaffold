package app

import (
	"errors"
	"github.com/google/uuid"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type App struct {
	container  core.IContainer //服务容器
	baseFolder string          //基础路径
	appId      string          //应用唯一ID

	configMap map[string]string //配置
	envMap    map[string]string //环境变量
	argsMap   map[string]string //参数
}

func (app *App) AppID() string {
	return app.appId
}

func (app *App) Version() string {
	return "0.0.1"
}

func (app *App) BaseFolder() string {
	if app.baseFolder != "" {
		return app.baseFolder
	}
	baseFolder := app.getConfigBySequence("base_folder", "BASE_FOLDER", "app.path.base_folder")
	if baseFolder != "" {
		return baseFolder
	}
	return utils.GetExecDirectory()
}

func (app *App) ConfigFolder() string {
	folder := app.getConfigBySequence("config_folder", "CONFIG_FOLDER", "app.path.config_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.BaseFolder(), "configs")
}

func (app *App) LogFolder() string {
	folder := app.getConfigBySequence("log_folder", "LOG_FOLDER", "app.path.log_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.StorageFolder(), "log")
}

func (app *App) HttpFolder() string {
	folder := app.getConfigBySequence("http_folder", "HTTP_FOLDER", "app.path.http_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.BaseFolder(), "app", "http")
}

func (app *App) ConsoleFolder() string {
	folder := app.getConfigBySequence("console_folder", "CONSOLE_FOLDER", "app.path.console_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.BaseFolder(), "app", "console")
}

func (app *App) StorageFolder() string {
	folder := app.getConfigBySequence("storage_folder", "STORAGE_FOLDER", "app.path.storage_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.BaseFolder(), "storage")
}

func (app *App) ProviderFolder() string {
	folder := app.getConfigBySequence("provider_folder", "PROVIDER_FOLDER", "app.path.provider_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.BaseFolder(), "app", "provider")
}

func (app *App) MiddlewareFolder() string {
	folder := app.getConfigBySequence("middleware_folder", "MIDDLEWARE_FOLDER", "app.path.middleware_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.HttpFolder(), "middleware")
}

func (app *App) CommandFolder() string {
	folder := app.getConfigBySequence("command_folder", "COMMAND_FOLDER", "app.path.command_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.ConsoleFolder(), "command")
}

func (app *App) RuntimeFolder() string {
	folder := app.getConfigBySequence("runtime_folder", "RUNTIME_FOLDER", "app.path.runtime_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.StorageFolder(), "runtime")
}

func (app *App) TestFolder() string {
	folder := app.getConfigBySequence("test_folder", "TEST_FOLDER", "app.path.test_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.BaseFolder(), "test")
}

func (app *App) DeployFolder() string {
	folder := app.getConfigBySequence("deploy_folder", "DEPLOY_FOLDER", "app.path.deploy_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.BaseFolder(), "deploy")
}

func (app *App) AppFolder() string {
	folder := app.getConfigBySequence("app_folder", "APP_FOLDER", "app.path.app_folder")
	if folder != "" {
		return folder
	}
	return filepath.Join(app.BaseFolder(), "app")
}

func NewApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("NewApp params error")
	}

	container := params[0].(core.IContainer)
	baseFolder := params[1].(string)

	appId := uuid.New().String()
	configMap := map[string]string{}

	app := &App{
		baseFolder: baseFolder,
		container:  container,
		appId:      appId,
		configMap:  configMap,
	}

	_ = app.loadEnvMaps()
	_ = app.loadArgsMaps()

	return app, nil
}

// 配置优先级：参数 > 环境变量 > 配置文件
func (app *App) getConfigBySequence(argsKey string, envKey string, configKey string) string {
	if app.argsMap != nil {
		if val, ok := app.argsMap[argsKey]; ok {
			return val
		}
	}

	if app.envMap != nil {
		if val, ok := app.envMap[envKey]; ok {
			return val
		}
	}

	if app.configMap != nil {
		if val, ok := app.configMap[configKey]; ok {
			return val
		}
	}
	return ""
}

func (app *App) loadEnvMaps() error {
	if app.envMap == nil {
		app.envMap = map[string]string{}
	}
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		app.envMap[pair[0]] = pair[1]
	}
	return nil
}

func (app *App) loadArgsMaps() error {
	if app.argsMap == nil {
		app.argsMap = map[string]string{}
	}
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			pair := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
			if len(pair) == 2 {
				app.argsMap[pair[0]] = pair[1]
			}
		}
	}
	return nil
}

func (app *App) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		app.configMap[key] = val
	}
}
