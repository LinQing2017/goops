package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kube-tools/src/util"
	"net/http"
	"net/url"
	"strings"
)

var (
	eurekaNetFilter string
)

func addEurekaFlag(flags *pflag.FlagSet) {
	flags.StringVar(&eurekaNetFilter, "net", "10.187.0.0/16", "过滤对应网段。")
}

func NewCmdEureka() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "eureka",
		Short:                 "抓取Spring Eureka页面的target地址。",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunEureka(cmd, args)
		},
	}
	addEurekaFlag(cmd.Flags())
	return cmd
}

type Instance struct {
	Application string
	Upstreams   string
}

func RunEureka(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		panic("请填写对应URL")
	}
	baseUrl := args[0]
	instanceList := make([]Instance, 0)
	res, _ := http.Get(baseUrl)
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	tableNodes := doc.Find("table")
	tableNodes.Each(func(i int, tableNode *goquery.Selection) {
		tableId, _ := tableNode.Attr("id")
		tableClass, _ := tableNode.Attr("class")
		if strings.EqualFold(tableId, "instances") && strings.EqualFold(tableClass, "table table-striped table-hover") {
			tableNode.Find("tbody").Find("tr").Each(func(i int, row *goquery.Selection) {
				colum := row.Find("td")
				instance := Instance{
					colum.First().Text(),
					"",
				}
				colum.Find("a").Each(
					func(i int, target *goquery.Selection) {
						href, isExit := target.Attr("href")
						if isExit {

							u, _ := url.Parse(href)
							instance.Upstreams = u.Host
						}
					})
				if util.IpaddrBelong(instance.Upstreams, eurekaNetFilter) {
					instanceList = append(instanceList, instance)
				}

			})
		}
	})

	instanceStr := table.Table(instanceList)
	fmt.Println(instanceStr)
}
