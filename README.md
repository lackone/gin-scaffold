# gin-scaffold

gin脚手架，集成了viper，gorm，godotenv，cron，cobra，redis，jwt，gomail，pongo2

### 一、简单使用

```go
package main

import (
	"github.com/lackone/gin-scaffold/app/console"
	"github.com/lackone/gin-scaffold/app/http"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/provider/app"
	"github.com/lackone/gin-scaffold/internal/provider/cache"
	"github.com/lackone/gin-scaffold/internal/provider/config"
	"github.com/lackone/gin-scaffold/internal/provider/env"
	"github.com/lackone/gin-scaffold/internal/provider/gorm"
	"github.com/lackone/gin-scaffold/internal/provider/id"
	"github.com/lackone/gin-scaffold/internal/provider/kernel"
	"github.com/lackone/gin-scaffold/internal/provider/log"
	"github.com/lackone/gin-scaffold/internal/provider/redis"
)

// @title gin脚手架
// @version 1.0
// @description gin脚手架
// @termsOfService https://github.com/lackone/gin-scaffold
func main() {
	//初始化服务容器
	container := core.NewContainer()
	// 绑定App服务提供者
	container.Bind(&app.AppProvider{})
	container.Bind(&env.EnvProvider{})
	container.Bind(&config.ConfigProvider{})
	container.Bind(&id.IDProvider{})
	container.Bind(&log.LogProvider{})
	container.Bind(&gorm.GormProvider{})
	container.Bind(&redis.RedisProvider{})
	container.Bind(&cache.CacheProvider{})

	//将HTTP引擎，绑定到服务容器中
	if engine, err := http.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.KernelProvider{HttpEngine: engine})
	}

	//运行命令
	console.RunCommand(container)
}
```

在命令行输入

```shell
go run main.go app start
```