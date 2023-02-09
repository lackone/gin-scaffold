package cmd

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/spf13/cobra"
)

func initConfigCommand() *cobra.Command {
	configCommand.AddCommand(configGetCommand)
	return configCommand
}

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "获取配置相关信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var configGetCommand = &cobra.Command{
	Use:   "get",
	Short: "获取某个配置信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := global.RootCmd.GetContainer()
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		if len(args) != 1 {
			fmt.Println("参数错误")
			return nil
		}
		configPath := args[0]
		val := configService.Get(configPath)
		if val == nil {
			fmt.Println("配置路径 ", configPath, " 不存在")
			return nil
		}
		fmt.Printf("%# v\n", pretty.Formatter(val))
		return nil
	},
}
