package main

import (
	"github.com/spf13/cobra"
	"kube-tools/src/cmd"
	"kube-tools/src/config"
	cusPprof "kube-tools/src/pprof"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ktool [sub]",
		Short: "ND Kubernetes 运维违规小助手",
		Run:   runHelp,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return cusPprof.InitProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return cusPprof.FlushProfiling()
		},
	}

	flags := rootCmd.PersistentFlags()

	config.AddFlags(flags)

	cusPprof.AddProfilingFlags(flags)
	rootCmd.AddCommand(cmd.NewCmdNode())
	rootCmd.AddCommand(cmd.NewCmdDeploy())
	rootCmd.AddCommand(cmd.NewCmdSh())
	rootCmd.AddCommand(cmd.NewCmdDist())
	rootCmd.AddCommand(cmd.NewCmdCrawler())
	rootCmd.AddCommand(cmd.NewCmdDocker())

	if err := execute(rootCmd); err != nil {
		os.Exit(1)
	}
}

func execute(cmd *cobra.Command) error {
	err := cmd.Execute()
	return err
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
