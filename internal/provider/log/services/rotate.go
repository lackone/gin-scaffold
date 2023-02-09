package services

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/pkg/utils"
	"github.com/natefinch/lumberjack"
	"os"
	"path/filepath"
)

type RotateLog struct {
	Log
	folder string // 日志文件存储目录
	file   string // 日志文件名
}

func NewRotateLog(params ...interface{}) (interface{}, error) {
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

	log := &RotateLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.folder = folder
	log.file = file

	w := &lumberjack.Logger{
		Filename:   filepath.Join(folder, file),
		MaxSize:    rotateSize,
		MaxBackups: rotateCount,
		MaxAge:     maxAge,
		Compress:   false,
	}

	log.SetOutput(w)
	log.c = c

	return log, nil
}
