package main

import (
	"github.com/spf13/cobra"
	"goops/pkg/docker_tools"
	sdpLogger "goops/pkg/logger"
	"os"
)

func main() {
	sdpLogger.InitLogger()

	rootCmd := &cobra.Command{
		Use:                   "sdp-docker-cli",
		DisableFlagsInUseLine: true,
		Run:                   runHelp,
	}
	rootCmd.AddCommand(NewCmdTomcatSetter())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func NewCmdTomcatSetter() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "tomcat-set",
		Short:                 "配置环境的JVM以及tomcat。",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			docker_tools.Main(cmd, args)
		},
	}
	docker_tools.AddTomcatSetter(cmd.Flags())
	return cmd
}
