package cmd

import (
	"fmt"
	"github.com/erikdubbelboer/gspt"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/lackone/gin-scaffold/pkg/utils"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var cronDaemon = false

func initCronCommand() *cobra.Command {
	cronStartCommand.Flags().BoolVarP(&cronDaemon, "daemon", "d", false, "start Serve daemon")

	cronCommand.AddCommand(cronRestartCommand)
	cronCommand.AddCommand(cronStateCommand)
	cronCommand.AddCommand(cronStopCommand)
	cronCommand.AddCommand(cronListCommand)
	cronCommand.AddCommand(cronStartCommand)

	return cronCommand
}

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "定时任务相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var cronListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有的定时任务",
	RunE: func(cmd *cobra.Command, args []string) error {
		cronSpecs := global.RootCmd.GetCronSpecs()
		ps := make([][]string, 0)
		for _, cronSpec := range cronSpecs {
			line := []string{cronSpec.Type, cronSpec.Spec, cronSpec.Cmd.Use, cronSpec.Cmd.Short, cronSpec.ServiceName, "\r\n"}
			ps = append(ps, line)
		}
		fmt.Println(ps)
		return nil
	},
}

var cronStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := global.RootCmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		pidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")
		logFile := filepath.Join(appService.LogFolder(), "cron.log")
		curFolder := appService.BaseFolder()

		//守护进程模式
		if cronDaemon {
			//创建一个Context
			ctx := &daemon.Context{
				//设置pid文件
				PidFileName: pidFile,
				PidFilePerm: 0664,
				//设置日志文件
				LogFileName: logFile,
				LogFilePerm: 0640,
				//设置工作路径
				WorkDir: curFolder,
				//设置所有设置文件的mask，默认为750
				Umask: 027,
				//子进程的参数，按照这个参数设置，子进程的命令为 ./gin-scaffold cron start --daemon=true
				Args: []string{"", "cron", "start", "--daemon=true"},
			}
			//启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			d, err := ctx.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				fmt.Println("cron serve started, pid:", d.Pid)
				fmt.Println("log file:", logFile)
				return nil
			}

			//子进程执行Cron.Run
			defer ctx.Release()
			fmt.Println("cron daemon started")
			gspt.SetProcTitle("gin-scaffold cron")
			global.RootCmd.GetCron().Run()
			return nil
		}

		//无守护进程模式
		pid := strconv.Itoa(os.Getpid())
		fmt.Println("cron started")
		fmt.Println("cron serve started, pid:", pid)
		err := os.WriteFile(pidFile, []byte(pid), 0664)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("gin-scaffold cron")
		global.RootCmd.GetCron().Run()
		return nil
	},
}

var cronRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重启cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := global.RootCmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		pidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		pid, err := os.ReadFile(pidFile)
		if err != nil {
			return err
		}

		if pid != nil && len(pid) > 0 {
			pid, err := strconv.Atoi(string(pid))
			if err != nil {
				return err
			}
			if utils.CheckProcessExist(pid) {
				if err := utils.KillProcess(pid); err != nil {
					return err
				}
				//循环检查进程
				for i := 0; i < 10; i++ {
					if utils.CheckProcessExist(pid) == false {
						break
					}
					time.Sleep(1 * time.Second)
				}
				fmt.Println("kill cron process:" + strconv.Itoa(pid))
			}
		}

		cronDaemon = true
		return cronStartCommand.RunE(cmd, args)
	},
}

var cronStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止cron常驻进程",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := global.RootCmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		pidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		pid, err := os.ReadFile(pidFile)
		if err != nil {
			return err
		}

		if pid != nil && len(pid) > 0 {
			pid, err := strconv.Atoi(string(pid))
			if err != nil {
				return err
			}
			if err := utils.KillProcess(pid); err != nil {
				return err
			}
			if err := os.WriteFile(pidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("stop pid:", pid)
		}
		return nil
	},
}

var cronStateCommand = &cobra.Command{
	Use:   "state",
	Short: "cron常驻进程状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := global.RootCmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		pidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		pid, err := os.ReadFile(pidFile)
		if err != nil {
			return err
		}

		if pid != nil && len(pid) > 0 {
			pid, err := strconv.Atoi(string(pid))
			if err != nil {
				return err
			}
			if utils.CheckProcessExist(pid) {
				fmt.Println("cron server started, pid:", pid)
				return nil
			}
		}
		fmt.Println("no cron server start")
		return nil
	},
}
