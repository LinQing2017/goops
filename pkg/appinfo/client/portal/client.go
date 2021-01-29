package portal

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools/types"
	"strconv"
	"strings"
)

func SwitchDomain(appname, domain, clusterName string, env int) {
	appInfo, err := common.GetAppInfo(appname, env)
	if err != nil {
		errw := errors.Wrapf(err, "%s没有查询到应用信息", appname)
		logrus.Error(errw.Error())
		return
	}

	clusterId := appInfo.GetFullClusterId(clusterName)
	if strings.EqualFold(clusterId, "") {
		logrus.Error("目标集群不存在")
		return
	}
	var bdomain *types.RMDomains
	for _, bdomain = range appInfo.GetAllDomains() {
		if strings.EqualFold(bdomain.Domain, domain) {
			break
		}
	}
	if bdomain == nil {
		logrus.Error("没有查询到域名")
		return
	}
	if strings.EqualFold(bdomain.ClusterId, clusterId) {
		logrus.Error("该域名绑定在当当前集群")
		return
	}
	params := DomainSwitch{
		Env:       strconv.Itoa(env),
		ClusterId: clusterId,
		Domains: []Domain{
			{
				Protocol:    strings.Join(bdomain.Protocols, "/"),
				Domain:      bdomain.Domain,
				ClusterName: clusterName,
				ClusterId:   bdomain.ClusterId,
				Type:        bdomain.GetLabel("type"),
				DomainId:    bdomain.ID,
			},
		},
	}

	reqbody := Req4Portal{
		Type:   "WEB_DOMAIN_SWITCH_MANAGE",
		Params: params,
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-loginname", XLoginName).
		SetHeader("x-password", XPassword).
		ForceContentType("application/json").
		SetPathParam("url", fmt.Sprintf("%s/v0.2/domain/switch/%s", ServerUrl, appInfo.APPID)).
		SetBody(reqbody).
		Post(UcProxyURL)

	if err != nil {
		logrus.Errorf(errors.Wrapf(err, "请求请求失败。", err).Error())
		return
	}
	if resp.StatusCode() != 200 {
		logrus.Error("请求请求失败。")
		return
	}
}
