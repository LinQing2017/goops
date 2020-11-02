package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kube-tools/src/cmd/crawler"
)

var (
	CrawlerOutput string
)

func addCrawlerFlag(flags *pflag.FlagSet) {
}

func NewCmdCrawler() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "crawler [sub]",
		Short:                 "爬虫工具",
		DisableFlagsInUseLine: true,
	}

	addCrawlerFlag(cmd.Flags())
	cmd.AddCommand(crawler.NewCmdEureka())
	return cmd
}
