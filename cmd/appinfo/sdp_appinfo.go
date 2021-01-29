package main

import (
	"github.com/spf13/cobra"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/export"
	"goops/pkg/appinfo/get"
	"goops/pkg/appinfo/list"
	"goops/pkg/appinfo/migrate_k8s"
	"goops/pkg/appinfo/ops"
	"goops/pkg/appinfo/ping"
	sdpLogger "goops/pkg/logger"
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
		Use:                   "sdp-app-cli",
		DisableFlagsInUseLine: false,
		Run:                   runHelp,
	}
	db_tools.AddDBFlag(cmd.Flags())
	cmd.AddCommand(migrate_k8s.NewCmdMigrateK8s())
	cmd.AddCommand(get.NewCmdGet())
	cmd.AddCommand(list.NewListCmd())
	cmd.AddCommand(ping.NewCmdPing())
	cmd.AddCommand(export.NewCmdExport())
	cmd.AddCommand(ops.NewOpsCmd())
	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
