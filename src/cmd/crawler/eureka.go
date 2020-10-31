package crawler

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kube-tools/src/util"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	eurekaNetFilter  string
	eurekaPortFilter string
	outputFilename   string
)

func addEurekaFlag(flags *pflag.FlagSet) {
	flags.StringVar(&eurekaNetFilter, "net", "", "过滤对应网段。")
	flags.StringVar(&eurekaPortFilter, "port", "", "过滤对应端口。")
	flags.StringVar(&outputFilename, "file", "", "输出结果数据到Excel文件。")
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
	Upstreams   []string
	AppId       []string
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
					make([]string, 0),
					make([]string, 0),
				}
				colum.Find("a").Each(
					func(i int, target *goquery.Selection) {
						href, isExit := target.Attr("href")
						if isExit {
							u, _ := url.Parse(href)
							isFilter := false
							if !strings.EqualFold(eurekaNetFilter, "") && !util.IpaddrBelong(u.Host, eurekaNetFilter) {
								isFilter = true
							}
							if !strings.EqualFold(eurekaPortFilter, "") && !strings.EqualFold(u.Port(), eurekaPortFilter) {
								isFilter = true
							}
							if !isFilter {
								instance.Upstreams = append(instance.Upstreams, u.Host)
								instance.AppId = append(instance.AppId, strings.Split(target.Text(), "-")[0])
							}
						}
					})
				if len(instance.Upstreams) > 0 {
					instanceList = append(instanceList, instance)
				}
			})
		}
	})

	instanceStr := table.Table(instanceList)
	fmt.Println(instanceStr)
	excelOutput(instanceList)
}

func excelOutput(instances []Instance) {
	if strings.EqualFold(outputFilename, "") {
		return
	}

	title := map[string]string{"A1": "编号", "B1": "应用名称", "C1": "服务端口", "D1": "APP-ID"}
	f := excelize.NewFile()
	for k, v := range title {
		f.SetCellValue("Sheet1", k, v)
	}
	id := 1
	for _, instance := range instances {

		for i, upstream := range instance.Upstreams {
			rowNum := strconv.Itoa(id + 1)
			row := map[string]string{
				"A" + rowNum: strconv.Itoa(id),
				"B" + rowNum: instance.Application,
				"C" + rowNum: upstream,
				"D" + rowNum: instance.AppId[i],
			}
			for k, v := range row {
				f.SetCellValue("Sheet1", k, v)
			}
			id += 1
		}

	}

	if err := f.SaveAs(outputFilename + ".xlsx"); err != nil {
		fmt.Println(err)
	}
}
