package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	distNodename string
	distNodefile string
)

func RunDist(cmd *cobra.Command, args []string) {

}

func addDistFlag(flags *pflag.FlagSet) {
	flags.StringVar(&distNodename, "node", "", "分发文件到指定宿主机。")
	flags.StringVar(&distNodename, "nodefile", "", "通过文件指定要分发文件的宿主机。")
}

func NewCmdDist() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "dist [local file path] [host file path]",
		Short:                 "分发文件到宿主机目录",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunDist(cmd, args)
		},
	}
	addDistFlag(cmd.Flags())
	return cmd
}
