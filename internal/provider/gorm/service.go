package gorm

import (
	"context"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"sync"
	"time"
)

type GormService struct {
	container core.IContainer     // 服务容器
	dbs       map[string]*gorm.DB // key为dsn,value为gorm.DB
	lock      *sync.RWMutex       //读写锁
}

func NewGorm(params ...interface{}) (interface{}, error) {
	return &GormService{
		container: params[0].(core.IContainer),
		dbs:       make(map[string]*gorm.DB),
		lock:      &sync.RWMutex{},
	}, nil
}

func (g *GormService) GetDB(option ...contract.DBOption) (*gorm.DB, error) {
	logger := g.container.MustMake(contract.LogKey).(contract.Log)
	logService := g.container.MustMake(contract.LogKey).(contract.Log)

	// 读取默认配置
	config := GetDefaultConfig(g.container)

	// 设置Logger
	gormLogger := NewGORMLogger(logService)
	config.Config = &gorm.Config{
		Logger: gormLogger,
	}

	// 遍历然后设置参数
	for _, opt := range option {
		if err := opt(g.container, config); err != nil {
			return nil, err
		}
	}

	// 如果最终的config没有设置dsn,就生成dsn
	if config.Dsn == "" {
		dsn, err := config.FormatDSN()
		if err != nil {
			return nil, err
		}
		config.Dsn = dsn
	}

	// 判断是否已经实例化了gorm.DB
	g.lock.RLock()
	if db, ok := g.dbs[config.Dsn]; ok {
		g.lock.RUnlock()
		return db, nil
	}
	g.lock.RUnlock()

	// 没有实例化gorm.DB，那么就要进行实例化操作
	g.lock.Lock()
	defer g.lock.Unlock()

	// 实例化gorm.DB
	var db *gorm.DB
	var err error

	switch config.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Dsn), config)
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Dsn), config)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Dsn), config)
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(config.Dsn), config)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(config.Dsn), config)
	}

	// 设置对应的连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	if config.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
	}
	if config.ConnMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
	}
	if config.ConnMaxLifetime != "" {
		liftTime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			logger.Error(context.Background(), "conn max lift time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxLifetime(liftTime)
		}
	}
	if config.ConnMaxIdletime != "" {
		idleTime, err := time.ParseDuration(config.ConnMaxIdletime)
		if err != nil {
			logger.Error(context.Background(), "conn max idle time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxIdleTime(idleTime)
		}
	}

	// 挂载到map中，结束配置
	if err == nil {
		g.dbs[config.Dsn] = db
	}

	return db, err
}
