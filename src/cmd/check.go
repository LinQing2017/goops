package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kube-tools/src/cmd/check"
)

func addCheckFlag(flags *pflag.FlagSet) {
}

func NewCmdCheck() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "check [sub]",
		Short:                 "检查工具",
		DisableFlagsInUseLine: true,
	}

	addCheckFlag(cmd.Flags())
	cmd.AddCommand(check.NewCmdIngressTLS())
	cmd.AddCommand(check.NewCmdAppInfo())
	return cmd
}
