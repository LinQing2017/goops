package list

import "github.com/spf13/cobra"

func NewListGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "批量获取应用信息",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			Main(cmd, args)
		},
	}
	AddListFlag(cmd.Flags())
	return cmd
}

func Main(cmd *cobra.Command, args []string) {

}
