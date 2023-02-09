package services

import (
	"context"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/provider/log/formatter"
	"io"
	pkgLog "log"
	"time"
)

type Log struct {
	level      contract.LogLevel   // 日志级别
	formatter  contract.Formatter  // 日志格式化方法
	ctxFielder contract.CtxFielder // ctx获取上下文字段
	output     io.Writer           // 输出
	c          core.IContainer     // 容器
}

// IsLevelEnable 判断这个级别是否可以打印
func (l *Log) IsLevelEnable(level contract.LogLevel) bool {
	return level <= l.level
}

func (l *Log) logf(level contract.LogLevel, ctx context.Context, msg string, fields map[string]interface{}) error {
	if !l.IsLevelEnable(level) {
		return nil
	}

	//使用ctxFielder获取context中的信息
	fs := fields
	if l.ctxFielder != nil {
		t := l.ctxFielder(ctx)
		if t != nil {
			for k, v := range t {
				fs[k] = v
			}
		}
	}

	if l.formatter == nil {
		l.formatter = formatter.TextFormatter
	}

	ct, err := l.formatter(level, time.Now(), msg, fs)
	if err != nil {
		return err
	}

	// 如果是panic级别，则使用log进行panic
	if level == contract.PanicLevel {
		pkgLog.Panicln(string(ct))
		return nil
	}

	l.output.Write(ct)
	l.output.Write([]byte("\r\n"))
	return nil
}

func (l *Log) SetOutput(output io.Writer) {
	l.output = output
}

func (l *Log) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logf(contract.PanicLevel, ctx, msg, fields)
}

func (l *Log) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logf(contract.FatalLevel, ctx, msg, fields)
}

func (l *Log) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logf(contract.ErrorLevel, ctx, msg, fields)
}

func (l *Log) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logf(contract.WarnLevel, ctx, msg, fields)
}

func (l *Log) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logf(contract.InfoLevel, ctx, msg, fields)
}

func (l *Log) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logf(contract.DebugLevel, ctx, msg, fields)
}

func (l *Log) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logf(contract.TraceLevel, ctx, msg, fields)
}

func (l *Log) SetLevel(level contract.LogLevel) {
	l.level = level
}

func (l *Log) SetCtxFielder(handler contract.CtxFielder) {
	l.ctxFielder = handler
}

func (l *Log) SetFormatter(formatter contract.Formatter) {
	l.formatter = formatter
}
