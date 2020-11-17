package main

import (
	"github.com/spf13/cobra"
	"kube-tools/pkg/appinfo"
	sdpLogger "kube-tools/pkg/logger"
	"os"
)

func main() {
	sdpLogger.InitLogger()
	command := NewCmdAppInfo()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewCmdAppInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sdpappinfo",
		Short:                 "获取集群应用信息。",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			appinfo.Main(cmd, args)
		},
	}
	appinfo.AddAppInfoFlag(cmd.Flags())
	return cmd
}
