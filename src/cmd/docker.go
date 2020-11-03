package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kube-tools/src/cmd/docker"
)

func addDockerFlag(flags *pflag.FlagSet) {
}

func NewCmdDocker() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "docker [sub]",
		Short:                 "Docker 工具",
		DisableFlagsInUseLine: true,
	}

	addDockerFlag(cmd.Flags())
	cmd.AddCommand(docker.NewCmdDisk())
	return cmd
}
