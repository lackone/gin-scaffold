package services

import (
	"errors"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/pkg/utils"
	"os"
	"path/filepath"
)

type SingleLog struct {
	Log
	folder string   // 日志文件存储目录
	file   string   // 日志文件名
	fd     *os.File //文件指针
}

func NewSingleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(core.IContainer)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	appService := c.MustMake(contract.AppKey).(contract.App)
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	log := &SingleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	folder := appService.LogFolder()

	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}
	log.folder = folder

	if !utils.Exists(folder) {
		os.MkdirAll(folder, os.ModePerm)
	}

	log.file = "app.log"
	if configService.IsExist("log.file") {
		log.file = configService.GetString("log.file")
	}

	fd, err := os.OpenFile(filepath.Join(log.folder, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.New("open log file err")
	}

	log.fd = fd
	log.SetOutput(fd)
	log.c = c

	return log, nil
}
