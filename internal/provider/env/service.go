package env

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/lackone/gin-scaffold/internal/contract"
	"os"
	"path"
	"strings"
)

type Env struct {
	folder string            // 代表.env所在的目录
	maps   map[string]string // 保存所有的环境变量
}

func NewEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewEnv params error")
	}

	folder := params[0].(string)

	env := &Env{
		folder: folder,
		maps:   map[string]string{"APP_ENV": contract.EnvDev},
	}

	appEnv := env.Get("APP_ENV")

	envFile := path.Join(folder, ".env")
	appEnvFile := path.Join(folder, ".env."+appEnv)

	envMap, _ := godotenv.Read(envFile, appEnvFile)
	for k, v := range envMap {
		env.maps[k] = v
	}

	//获取当前程序的环境变量，并且覆盖.env文件下的变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		env.maps[pair[0]] = pair[1]
	}

	return env, nil
}

func (e *Env) AppEnv() string {
	return e.Get("APP_ENV")
}

func (e *Env) IsExist(key string) bool {
	_, ok := e.maps[key]
	return ok
}

func (e *Env) Get(key string) string {
	if val, ok := e.maps[key]; ok {
		return val
	}
	return ""
}

func (e *Env) All() map[string]string {
	return e.maps
}
