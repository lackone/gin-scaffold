package log

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/provider/log/formatter"
	"github.com/lackone/gin-scaffold/internal/provider/log/services"
	"io"
	"strings"
)

type LogProvider struct {
	Driver     string              // 驱动
	Level      contract.LogLevel   // 日志级别
	Formatter  contract.Formatter  // 日志输出格式方法
	CtxFielder contract.CtxFielder // 日志context上下文信息获取函数
	Output     io.Writer           // 日志输出信息
}

func (l *LogProvider) Register(c core.IContainer) core.NewInstance {
	if l.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			return services.NewConsoleLog
		}

		cs := tcs.(contract.Config)
		l.Driver = strings.ToLower(cs.GetString("log.driver"))
	}

	switch l.Driver {
	case "single":
		return services.NewSingleLog
	case "rotate":
		return services.NewRotateLog
	case "console":
		return services.NewConsoleLog
	case "custom":
		return services.NewCustomLog
	case "zap":
		return services.NewZapLog
	default:
		return services.NewConsoleLog
	}
}

func (l *LogProvider) Boot(c core.IContainer) error {
	return nil
}

func (l *LogProvider) IsDefer() bool {
	return false
}

func (l *LogProvider) Params(c core.IContainer) []interface{} {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	if l.Formatter == nil {
		l.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				l.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				l.Formatter = formatter.TextFormatter
			}
		}
	}

	if l.Level == contract.UnknownLevel {
		l.Level = contract.InfoLevel
		if configService.IsExist("log.level") {
			l.Level = logLevel(configService.GetString("log.level"))
		}
	}

	return []interface{}{c, l.Level, l.CtxFielder, l.Formatter, l.Output}
}

func (l *LogProvider) Name() string {
	return contract.LogKey
}

func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnknownLevel
}
