package cmd

import (
	"fmt"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/spf13/cobra"
)

func initEnvCommand() *cobra.Command {
	envCommand.AddCommand(envListCommand)
	envCommand.AddCommand(envGetCommand)
	return envCommand
}

var envCommand = &cobra.Command{
	Use:   "env",
	Short: "获取当前的环境变量",
	Run: func(cmd *cobra.Command, args []string) {
		container := global.RootCmd.GetContainer()
		envService := container.MustMake(contract.EnvKey).(contract.Env)
		fmt.Println("environment:", envService.AppEnv())
	},
}

var envListCommand = &cobra.Command{
	Use:   "list",
	Short: "获取所有的环境变量",
	Run: func(cmd *cobra.Command, args []string) {
		container := global.RootCmd.GetContainer()
		envService := container.MustMake(contract.EnvKey).(contract.Env)
		envs := envService.All()
		for k, v := range envs {
			fmt.Println(k, "=", v)
		}
	},
}

var envGetCommand = &cobra.Command{
	Use:   "get",
	Short: "获取某个环境变量",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := global.RootCmd.GetContainer()
		envService := container.MustMake(contract.EnvKey).(contract.Env)
		if len(args) != 1 {
			fmt.Println("参数错误")
			return nil
		}
		env := args[0]
		val := envService.Get(env)
		if val == "" {
			fmt.Println("环境变量 ", env, " 不存在")
			return nil
		}
		fmt.Println(env, "=", val)
		return nil
	},
}
