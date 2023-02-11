package config

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	c        core.IContainer         //容器
	folder   string                  //文件夹
	keyDelim string                  //路径的分隔符，默认为点
	envMaps  map[string]string       //所有的环境变量
	viper    map[string]*viper.Viper //viper实例
}

func NewConfig(params ...interface{}) (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	container := params[0].(core.IContainer)
	envConfigFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	config := &Config{
		c:        container,
		folder:   envConfigFolder,
		envMaps:  envMaps,
		keyDelim: ".",
		viper:    make(map[string]*viper.Viper),
	}

	//检查文件夹是否存在
	if _, err := os.Stat(envConfigFolder); os.IsNotExist(err) {
		return config, nil
	}

	filepath.Walk(envConfigFolder, func(p string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			config.loadConfigFile(filepath.Dir(p), filepath.Base(p))
		}
		return nil
	})

	for _, v := range config.viper {
		go func(v *viper.Viper) {
			v.OnConfigChange(func(in fsnotify.Event) {
				config.loadConfigFile(filepath.Dir(in.Name), filepath.Base(in.Name))
			})
			v.WatchConfig()
		}(v)
	}

	return config, nil
}

func (c *Config) loadConfigFile(folder string, fileName string) error {
	data, err := os.ReadFile(filepath.Join(folder, fileName))
	if err != nil {
		return err
	}

	data = replace(data, c.envMaps)
	name := strings.TrimSuffix(fileName, path.Ext(fileName))

	if _, ok := c.viper[name]; !ok {
		v := viper.NewWithOptions(viper.KeyDelimiter(c.keyDelim))
		v.SetConfigName(name)
		v.SetConfigType("yaml")
		v.AddConfigPath(c.folder)

		c.viper[name] = v
	}

	err = c.viper[name].ReadConfig(bytes.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if name == "app" && c.c.IsBind(contract.AppKey) {
		appService := c.c.MustMake(contract.AppKey).(contract.App)
		path := c.viper["app"].GetStringMapString("path")
		appService.LoadAppConfig(path)
	}

	return nil
}

// replace表示使用环境变量maps替换context中的env(xxx)的环境变量
func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}

	return content
}

func (c *Config) find(key string) interface{} {
	split := strings.SplitN(key, c.keyDelim, 2)
	if len(split) < 1 {
		return nil
	}
	if _, ok := c.viper[split[0]]; !ok {
		return nil
	}
	if len(split) == 1 {
		return c.viper[split[0]].AllSettings()
	}
	return c.viper[split[0]].Get(split[1])
}

func (c *Config) IsExist(key string) bool {
	return c.find(key) != nil
}

func (c *Config) Get(key string) interface{} {
	return c.find(key)
}

func (c *Config) GetBool(key string) bool {
	return cast.ToBool(c.find(key))
}

func (c *Config) GetInt(key string) int {
	return cast.ToInt(c.find(key))
}

func (c *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.find(key))
}

func (c *Config) GetTime(key string) time.Time {
	return cast.ToTime(c.find(key))
}

func (c *Config) GetString(key string) string {
	return cast.ToString(c.find(key))
}

func (c *Config) GetIntSlice(key string) []int {
	return cast.ToIntSlice(c.find(key))
}

func (c *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.find(key))
}

func (c *Config) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(c.find(key))
}

func (c *Config) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(c.find(key))
}

func (c *Config) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(c.find(key))
}

func (c *Config) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(c.find(key))
}
