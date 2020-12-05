package main

import (
	"github.com/spf13/cobra"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/get"
	"goops/pkg/appinfo/list"
	"goops/pkg/appinfo/migrate_k8s"
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
	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
