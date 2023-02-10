package cmd

import (
	"fmt"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/spf13/cobra"
	"github.com/swaggo/swag/gen"
	"path/filepath"
)

func initSwaggerCommand() *cobra.Command {
	swaggerCommand.AddCommand(swaggerGenCommand)
	return swaggerCommand
}

var swaggerCommand = &cobra.Command{
	Use:   "swagger",
	Short: "swagger对应命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var swaggerGenCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成对应的swagger文件,swagger.json swagger.yaml, docs.go",
	Run: func(c *cobra.Command, args []string) {
		container := global.RootCmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		outputDir := filepath.Join(appService.DocsFolder(), "swagger")
		baseFolder := filepath.Join(appService.AppFolder(), "http")

		conf := &gen.Config{
			// 遍历需要查询注释的目录
			SearchDir: baseFolder,
			// 不包含哪些文件
			Excludes: "",
			// 输出目录
			OutputDir: outputDir,
			// 整个swagger接口的说明文档注释
			MainAPIFile: "swagger.go",
			// 名字的显示策略，比如首字母大写等
			PropNamingStrategy: "",
			// 是否要解析vendor目录
			ParseVendor: false,
			// 是否要解析外部依赖库的包
			ParseDependency: false,
			// 是否要解析标准库的包
			ParseInternal: false,
			// 是否要查找markdown文件，这个markdown文件能用来为tag增加说明格式
			MarkdownFilesDir: "",
			// 是否应该在docs.go中生成时间戳
			GeneratedTime: false,
		}
		err := gen.New().Build(conf)
		if err != nil {
			fmt.Println(err)
		}
	},
}
