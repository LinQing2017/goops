package docker

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func addDiskFlag(flags *pflag.FlagSet) {

}

func NewCmdDisk() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "disk",
		Short:                 "检查Docker容器的磁盘使用情况。",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunDisk(cmd, args)
		},
	}
	addDiskFlag(cmd.Flags())
	return cmd
}

func RunDisk(cmd *cobra.Command, args []string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithHost("http://192.168.56.104:4243"))
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}
