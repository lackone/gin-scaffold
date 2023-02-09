package services

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"os"
)

type ConsoleLog struct {
	Log
}

func NewConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(core.IContainer)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &ConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	log.SetOutput(os.Stdout)
	log.c = c

	return log, nil
}
