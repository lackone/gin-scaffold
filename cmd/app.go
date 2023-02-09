package cmd

import (
	"context"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/global"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func initAppCmd() *cobra.Command {
	appCmd.AddCommand(appStartCmd)

	return appCmd
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "app相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

var appStartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动一个app服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := global.RootCmd.GetContainer()
		kernel := container.MustMake(contract.KernelKey).(contract.Kernel)
		engine := kernel.HttpEngine()

		server := &http.Server{
			Addr:    ":8080",
			Handler: engine,
		}

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalln("server listen err:", err)
			}
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-quit

		timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
		defer cancelFunc()

		if err := server.Shutdown(timeout); err != nil {
			log.Fatalln(err)
		}

		return nil
	},
}
