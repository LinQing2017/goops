package main

import (
	"github.com/spf13/cobra"
	"goops/pkg/jvm_setter"
	sdpLogger "goops/pkg/logger"
	"os"
)

func main() {
	sdpLogger.InitLogger()
	command := NewCmdJvmSetter()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewCmdJvmSetter() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "jvm_set",
		Short:                 "配置环境的JVM以及tomcat。",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			jvm_setter.Main(cmd, args)
		},
	}
	jvm_setter.AddJVMSetter(cmd.Flags())
	return cmd
}
