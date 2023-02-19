package services

import (
	"fmt"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/pkg/utils"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
	"io"
	"os"
	"path/filepath"
)

type ZapLog struct {
	*zap.Logger
	folder     string              // 日志文件存储目录
	file       string              // 日志文件名
	level      contract.LogLevel   // 日志级别
	formatter  contract.Formatter  // 日志格式化方法
	ctxFielder contract.CtxFielder // ctx获取上下文字段
	output     io.Writer           // 输出
	c          core.IContainer     // 容器
}

func NewZapLog(params ...interface{}) (interface{}, error) {
	c := params[0].(core.IContainer)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	appService := c.MustMake(contract.AppKey).(contract.App)
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	folder := appService.LogFolder()
	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}

	if !utils.Exists(folder) {
		os.MkdirAll(folder, os.ModePerm)
	}

	file := "app.log"
	if configService.IsExist("log.file") {
		file = configService.GetString("log.file")
	}

	rotateCount := 10
	if configService.IsExist("log.rotate_count") {
		rotateCount = configService.GetInt("log.rotate_count")
	}

	rotateSize := 10
	if configService.IsExist("log.rotate_size") {
		rotateSize = configService.GetInt("log.rotate_size")
	}

	maxAge := 10
	if configService.IsExist("log.max_age") {
		maxAge = configService.GetInt("log.max_age")
	}

	w := &lumberjack.Logger{
		Filename:   filepath.Join(folder, file),
		MaxSize:    rotateSize,
		MaxBackups: rotateCount,
		MaxAge:     maxAge,
		Compress:   false,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	writeSyncer := zapcore.AddSync(w)
	zapLevel := ConvertZapLevel(level)
	var stackTraceLevel zap.LevelEnablerFunc = func(level zapcore.Level) bool {
		return level >= zapcore.DPanicLevel
	}
	zl := zap.New(zapcore.NewCore(encoder, writeSyncer, zapLevel), zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(stackTraceLevel))
	zap.ReplaceGlobals(zl)

	log := &ZapLog{
		Logger: zl,
	}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.folder = folder
	log.file = file

	log.SetOutput(w)
	log.c = c

	return log, nil
}

func (l *ZapLog) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	fs := l.ConvertZapFields(ctx, fields)

	l.Logger.Panic(msg, fs...)
}

func (l *ZapLog) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	fs := l.ConvertZapFields(ctx, fields)

	l.Logger.Fatal(msg, fs...)
}

func (l *ZapLog) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	fs := l.ConvertZapFields(ctx, fields)

	l.Logger.Error(msg, fs...)
}

func (l *ZapLog) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	fs := l.ConvertZapFields(ctx, fields)

	l.Logger.Warn(msg, fs...)
}

func (l *ZapLog) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	fs := l.ConvertZapFields(ctx, fields)

	l.Logger.Info(msg, fs...)
}

func (l *ZapLog) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	fs := l.ConvertZapFields(ctx, fields)

	l.Logger.Debug(msg, fs...)
}

func (l *ZapLog) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	fs := l.ConvertZapFields(ctx, fields)

	l.Logger.Debug(msg, fs...)
}

func (l *ZapLog) ConvertZapFields(ctx context.Context, fields map[string]interface{}) []zap.Field {
	//使用ctxFielder获取context中的信息
	fs := fields
	zfs := make([]zap.Field, 0)
	if l.ctxFielder != nil {
		t := l.ctxFielder(ctx)
		if t != nil {
			for k, v := range t {
				fs[k] = v
			}
		}
	}
	for k, v := range fs {
		zfs = append(zfs, zap.String(k, fmt.Sprint(v)))
	}
	return zfs
}

func (l *ZapLog) SetOutput(output io.Writer) {
	l.output = output
}

func (l *ZapLog) SetLevel(level contract.LogLevel) {
	l.level = level
}

func (l *ZapLog) SetCtxFielder(handler contract.CtxFielder) {
	l.ctxFielder = handler
}

func (l *ZapLog) SetFormatter(formatter contract.Formatter) {
	l.formatter = formatter
}

func ConvertZapLevel(level contract.LogLevel) *zapcore.Level {
	var l zapcore.Level
	switch level {
	case contract.UnknownLevel:
		l = zapcore.InfoLevel
	case contract.PanicLevel:
		l = zapcore.PanicLevel
	case contract.FatalLevel:
		l = zapcore.FatalLevel
	case contract.ErrorLevel:
		l = zapcore.ErrorLevel
	case contract.WarnLevel:
		l = zapcore.WarnLevel
	case contract.InfoLevel:
		l = zapcore.InfoLevel
	case contract.DebugLevel:
		l = zapcore.DebugLevel
	case contract.TraceLevel:
		l = zapcore.DebugLevel
	}
	return &l
}
