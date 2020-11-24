package main

import (
	"github.com/spf13/cobra"
	"goops/pkg/eureka"
	sdpLogger "goops/pkg/logger"
	"os"
)

func main() {
	sdpLogger.InitLogger()
	command := NewCmdEureka()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewCmdEureka() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sdpeureka",
		Short:                 "抓取Spring Eureka页面的target地址。",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			eureka.Main(cmd, args)
		},
	}
	eureka.AddEurekaFlag(cmd.Flags())
	return cmd
}
