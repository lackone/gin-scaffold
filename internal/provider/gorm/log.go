package gorm

import (
	"context"
	"github.com/lackone/gin-scaffold/internal/contract"
	"gorm.io/gorm/logger"
	"time"
)

// 实现了gorm.Logger.Interface
type GORMLogger struct {
	logger contract.Log // 存放log服务
}

func NewGORMLogger(logger contract.Log) *GORMLogger {
	return &GORMLogger{logger: logger}
}

// LogMode 什么都不实现，日志级别完全依赖log服务的日志定义
func (g *GORMLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g *GORMLogger) Info(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	g.logger.Info(ctx, s, fields)
}

func (g *GORMLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	g.logger.Warn(ctx, s, fields)
}

func (g *GORMLogger) Error(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	g.logger.Error(ctx, s, fields)
}

func (g *GORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	fields := map[string]interface{}{
		"begin": begin,
		"error": err,
		"sql":   sql,
		"rows":  rows,
		"time":  elapsed,
	}
	s := "gorm trace sql"
	g.logger.Trace(ctx, s, fields)
}
