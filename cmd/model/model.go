package model

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/lackone/gin-scaffold/internal/provider/gorm"
	"github.com/lackone/gin-scaffold/pkg/utils"
	"github.com/spf13/cobra"
	"gorm.io/gen"
	"os"
	"path/filepath"
)

// 代表输出路径
var output string

// 代表数据库连接
var database string

func InitModelCommand() *cobra.Command {
	modelGenCommand.Flags().StringVarP(&output, "output", "o", "", "模型输出地址")
	modelGenCommand.MarkFlagRequired("output")
	modelGenCommand.Flags().StringVarP(&database, "database", "d", "database.default", "模型连接的数据库")

	modelCommand.AddCommand(modelGenCommand)
	return modelCommand
}

var modelCommand = &cobra.Command{
	Use:   "model",
	Short: "数据库模型相关的命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var modelGenCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成模型",
	RunE: func(c *cobra.Command, args []string) error {

		// 确认output路径是绝对路径
		if !filepath.IsAbs(output) {
			absOutput, err := filepath.Abs(output)
			if err != nil {
				return err
			}
			output = absOutput
		}

		// 获取env环境
		container := global.RootCmd.GetContainer()
		logger := container.MustMake(contract.LogKey).(contract.Log)
		logger.SetLevel(contract.ErrorLevel)

		// 加载数据库配直
		gormService := container.MustMake(contract.GORMKey).(contract.IGORM)
		db, err := gormService.GetDB(gorm.WithConfigPath(database))
		if err != nil {
			fmt.Println("数据库连接：" + database + "失败，请检查配置")
			return err
		}

		// 获取所有表
		allTables, err := db.Migrator().GetTables()
		if err != nil {
			return err
		}
		tables := make([]string, 0, len(allTables)+1)
		tables = append(tables, "*")
		tables = append(tables, allTables...)

		// 第一步询问要生成的表
		genTables := make([]string, 0, len(allTables)+1)
		survey.AskOne(&survey.MultiSelect{
			Message: "请选择要生成模型的表：",
			Options: tables,
		}, &genTables)
		if utils.SliceStrContains(genTables, "*") {
			genTables = allTables
		}
		if len(genTables) == 0 {
			return errors.New("未选择任何需要生成模型的表")
		}

		// 第二步确认要生成的目录和文件，以及覆盖提示：
		existFileNames := make([]string, 0)
		if utils.Exists(output) {
			files, err := os.ReadDir(output)
			if err != nil {
				return err
			}
			for _, file := range files {
				existFileNames = append(existFileNames, file.Name())
			}
		}

		genFileNames := make([]string, 0, len(genTables))
		for _, genTable := range genTables {
			genFileNames = append(genFileNames, genTable+".gen.go")
		}

		createFileNames := utils.SliceStrDiff(genFileNames, existFileNames)
		replaceFileNames := utils.SliceStrInter(genFileNames, existFileNames)

		fmt.Println("继续下列操作会在目录（" + output + "）生成下列文件：")
		for _, genFileName := range genFileNames {
			if utils.SliceStrContains(createFileNames, genFileName) {
				fmt.Println(genFileName + "（新文件）")
			} else if utils.SliceStrContains(replaceFileNames, genFileName) {
				fmt.Println(genFileName + "（覆盖）")
			} else {
				fmt.Println(genFileName)
			}
		}

		genContinue := false
		survey.AskOne(&survey.Confirm{
			Message: "请确认是否继续？",
		}, &genContinue)

		if genContinue == false {
			fmt.Println("操作暂停")
			return nil
		}

		// 第三步选择后是一个生成模型的选项：
		selectRuleTips := make([]string, 0)
		ruleTips := map[string]string{
			"FieldNullable":     "FieldNullable, 对于数据库的可null字段设置指针",
			"FieldCoverable":    "FieldCoverable, 根据数据库的Default设置字段的默认值",
			"FieldWithIndexTag": "FieldWithIndexTag, 根据数据库的索引关系设置索引标签",
			"FieldWithTypeTag":  "FieldWithTypeTag, 生成类型字段",
		}
		tips := make([]string, 0, len(ruleTips))
		for _, val := range ruleTips {
			tips = append(tips, val)
		}
		survey.AskOne(&survey.MultiSelect{
			Message: "请选择生成的模型规则：",
			Options: tips,
		}, &selectRuleTips)

		isSelectRule := func(key string, selectRuleTips []string, allRuleTips map[string]string) bool {
			tip := allRuleTips[key]
			return utils.SliceStrContains(selectRuleTips, tip)
		}

		// 生成模型文件
		g := gen.NewGenerator(gen.Config{
			ModelPkgPath: output,

			FieldNullable:     isSelectRule("FieldNullable", selectRuleTips, ruleTips),
			FieldCoverable:    isSelectRule("FieldCoverable", selectRuleTips, ruleTips),
			FieldWithIndexTag: isSelectRule("FieldWithIndexTag", selectRuleTips, ruleTips),
			FieldWithTypeTag:  isSelectRule("FieldWithTypeTag", selectRuleTips, ruleTips),

			Mode: gen.WithDefaultQuery,
		})

		g.UseDB(db)

		for _, table := range genTables {
			g.GenerateModel(table)
		}

		g.Execute()

		fmt.Println("生成模型成功")
		return nil
	},
}
