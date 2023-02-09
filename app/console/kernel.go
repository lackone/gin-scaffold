package console

import (
	"github.com/lackone/gin-scaffold/app/console/commands/test"
	"github.com/lackone/gin-scaffold/cmd"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/spf13/cobra"
)

func RunCommand(container core.IContainer) error {
	global.RootCmd = &core.Command{
		Command: &cobra.Command{
			Use:   "gin-scaffold",
			Short: "gin-scaffold命令",
			Long:  "gin-scaffold框架提供的命令行工具",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.InitDefaultHelpFlag()
				return cmd.Help()
			},
			CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		},
	}

	//为根命令设置容器
	global.RootCmd.SetContainer(container)
	//绑定框架的命令
	cmd.AddKernelCommands(global.RootCmd)
	//绑定业务的命令
	AddAppCommand(global.RootCmd)

	return global.RootCmd.Execute()
}

// 业务的命令
func AddAppCommand(rootCmd *core.Command) {
	rootCmd.AddCronCommand("* * * * * *", test.TestCmd)
}
