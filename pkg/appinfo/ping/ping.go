package ping

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/modood/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/db_tools/types"
	systools "goops/pkg/util/sys"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewCmdPing() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ping",
		Short:                 "对应用进行拨测",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			Main(cmd, args)
		},
	}
	AddPingFlag(cmd.Flags())
	return cmd
}

func Main(cmd *cobra.Command, args []string) {

	db_tools.InitDBClient()
	printPings := make([]*PrintPing, 0)
	apps := make(map[string]common.AppInformation)
	var app common.AppInformation
	if batch {
		for _, appname := range systools.ReadLine(nameFromFile) {
			appname = strings.TrimSpace(appname)
			app, printPings = GetPrintPings(appname, printPings)
			apps[appname] = app
		}
	} else {
		appname := strings.TrimSpace(args[0])
		app, printPings = GetPrintPings(args[0], printPings)
		apps[appname] = app
	}
	db_tools.CloseAllDBClient()

	for _, codeStr := range strings.Split(codeList, ",") {
		code, _ := strconv.Atoi(codeStr)
		codeAccept.Add(code)
	}

	for _, printPing := range printPings {
		ping(printPing)
	}
	table.Output(printPings)

	if batch {
		summary(printPings, apps)
	}
}

func GetPrintPings(appname string, printPings []*PrintPing) (common.AppInformation, []*PrintPing) {

	appInformation, err := common.GetAppInfo(appname, envType)
	if err != nil {
		logrus.Error("查询数据库失败，", err.Error())
		return appInformation, printPings
	}
	server := getServer(appInformation)
	for _, cluster := range server {
		binds := appInformation.ClusterBindDomains[cluster.ClusterId]
		if binds == nil || len(binds) < 1 {
			logrus.Info("集群", cluster.ClusterId, " ("+appname+") ", "没有绑定域名")
			continue
		}
		for _, domain := range binds {
			isTestDomain := false
			for _, label := range domain.Labels {
				if strings.EqualFold(label.Key, "role") && strings.EqualFold(label.Value, "test") {
					logrus.Debug("测试域名不进行拨测：", domain.Domain)
					isTestDomain = true
				}
			}
			if !isTestDomain {
				printPing := PrintPing{
					Domain:    domain.Domain,
					AppName:   appInformation.NAME,
					ClusterId: cluster.ClusterId,
					Code:      -999,
					WafCode:   -999,
				}
				printPings = append(printPings, &printPing)
			}

		}
	}
	return appInformation, printPings

}

func ping(printPing *PrintPing) {
	url := "http://" + printPing.Domain
	if resp, err := connect(url); err != nil {
		printPing.Message += "拨测异常"
	} else {
		printPing.Code = resp.StatusCode
	}
	wafUrl := "http://" + printPing.Domain + "/$waf/ping"

	if resp, err := connect(wafUrl); err != nil {
		printPing.Message += "waf拨测异常"
	} else {
		printPing.WafCode = resp.StatusCode
		if !codeAccept.Contains(printPing.WafCode) {
			printPing.Message += "WAF返回值异常"
		}
	}
}

func connect(url string) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(25 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("构造Request失败（%s）。", err.Error())
		return nil, err
	}

	if err != nil {

		return nil, err

	}

	resp, err := client.Do(req)
	return resp, err
}

func summary(printPings []*PrintPing, apps map[string]common.AppInformation) {

	fmt.Println()
	fmt.Println("==========================================")

	appPings := mapset.NewSet()            // 实际进行拨测的应用
	clusterPings := mapset.NewSet()        // 实际进行拨测的集群
	appInputs := 0                         // 文件输入的应用名称(不计重复)
	clusterFromProtalDB := mapset.NewSet() // 数据库查询到的集群记录
	for _, p := range printPings {
		appPings.Add(p.AppName)
		clusterPings.Add(p.ClusterId)
	}
	for _, app := range apps {
		appInputs++
		server := getServer(app)
		if server == nil || len(server) == 0 {
			fmt.Println(app.NAME, "没有查询到集群信息")
		} else {
			for _, cluster := range server {
				clusterFromProtalDB.Add(cluster.ClusterId)
			}
		}

	}

	fmt.Printf("拨测应用：%d个，输入应用：%d\n", len(appPings.ToSlice()), appInputs)
	fmt.Printf("拨测集群：%d个，门户查询到集群：%d\n", len(clusterPings.ToSlice()), len(clusterFromProtalDB.ToSlice()))
}

func getServer(app common.AppInformation) []*types.Service {
	var server []*types.Service
	switch clusterType {
	case "k8s":
		server = app.PortalInfo.K8SServiceList
	case "ews":
		server = app.PortalInfo.EWSServiceList
	default:
		logrus.Error("选择集群类型异常")
		os.Exit(-1)
	}
	return server
}
