package cmd

import (
	"github.com/lackone/gin-scaffold/cmd/model"
	"github.com/lackone/gin-scaffold/internal/core"
)

func AddKernelCommands(rootCmd *core.Command) {
	//挂载app命令
	rootCmd.AddCommand(initAppCmd())
	//挂载cron命令
	rootCmd.AddCommand(initCronCommand())
	//挂载config命令
	rootCmd.AddCommand(initConfigCommand())
	//挂载env命令
	rootCmd.AddCommand(initEnvCommand())
	//挂载swag命令
	rootCmd.AddCommand(initSwaggerCommand())
	//挂载model命令
	rootCmd.AddCommand(model.InitModelCommand())
}
