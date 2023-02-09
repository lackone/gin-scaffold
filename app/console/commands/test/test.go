package test

import (
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/spf13/cobra"
	"log"
)

var TestCmd = &core.Command{
	Command: &cobra.Command{
		Use:   "test",
		Short: "测试",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("test....")
			return nil
		},
	},
}
