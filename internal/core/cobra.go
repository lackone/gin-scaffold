package core

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"log"
)

type Command struct {
	*cobra.Command
	container IContainer
	cron      *cron.Cron
	cronSpecs []CronSpec
	root      *Command
}

type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

func (c *Command) SetContainer(container IContainer) {
	c.container = container
}

func (c *Command) GetContainer() IContainer {
	return c.container
}

func (c *Command) GetRoot() *Command {
	if c.root == nil {
		return c
	}
	return c.root.GetRoot()
}

func (c *Command) GetCronSpecs() []CronSpec {
	return c.cronSpecs
}

func (c *Command) GetCron() *cron.Cron {
	return c.cron
}

func (c *Command) AddCronCommand(spec string, cmd *Command) {
	root := c.GetRoot()
	if root.cron == nil {
		//初始化cron
		root.cron = cron.New(cron.WithSeconds())
		root.cronSpecs = []CronSpec{}
	}

	root.cronSpecs = append(root.cronSpecs, CronSpec{
		Type: "normal",
		Cmd:  cmd,
		Spec: spec,
	})

	//创建一个command
	ctx := root.Context()
	cronCmd := *cmd
	cronCmd.SetArgs([]string{})
	cronCmd.ResetCommands()
	cronCmd.SetContainer(root.GetContainer())

	//添加调用函数
	root.cron.AddFunc(spec, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})
}
