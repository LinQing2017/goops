package docker_tools

import (
	"encoding/json"
	"github.com/beevik/etree"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type SDP_CONFIG struct {
	SdpJavaOpts        string `json:"sdp_java_opts"`
	SdpJavaOptsOthers  string `json:"sdp_java_opts_others"`
	SdpServerXmlArgs   string `json:"sdp_server_xml_args"`
	SdpServerAccessLog string `json:"sdp_server_access_log"`
	ServerArgsDict     map[string]string
	HostDomainConfig   []HostDomains
}

type HostDomains struct {
	APPNAME string
	APPID   string
	Domains []string
}

func initFromEnv(sdpConfig *SDP_CONFIG) {

	// 从环境变量读取JVM配置
	sdpConfig.SdpJavaOpts = os.Getenv("SDP_JAVA_OPTS")
	sdpConfig.SdpJavaOptsOthers = os.Getenv("SDP_JAVA_OPTS_OTHERS")
	sdpConfig.SdpServerXmlArgs = os.Getenv("SDP_SERVER_XML_ARGS")
	sdpConfig.SdpServerAccessLog = os.Getenv("SDP_SERVER_ACCESS_LOG")

	if strings.EqualFold(sdpConfig.SdpJavaOpts, "") ||
		strings.EqualFold(sdpConfig.SdpJavaOptsOthers, "") ||
		strings.EqualFold(sdpConfig.SdpServerXmlArgs, "") ||
		strings.EqualFold(sdpConfig.SdpServerAccessLog, "") {
		configStr, _ := json.Marshal(sdpConfig)
		logrus.Error("读取环境变量异常：" + string(configStr))
		os.Exit(-1)
	}

	serverArgs := strings.Split(sdpConfig.SdpServerXmlArgs, ",")
	sdpConfig.ServerArgsDict = make(map[string]string, 0)
	for _, arg := range serverArgs {
		t := strings.Split(arg, "=")
		if len(t) != 2 {
			logrus.Error("SDP_SERVER_XML_ARGS 配置异常", arg)
			continue
		}
		sdpConfig.ServerArgsDict[strings.TrimSpace(t[0])] = strings.TrimSpace(t[1])
	}

	// 从环境变量读取主机名配置
	sdpConfig.HostDomainConfig = make([]HostDomains, 0)
	for _, envStr := range os.Environ() {
		if strings.HasPrefix(envStr, "SDP_MERGE_APP_NAME_") {
			t := strings.Split(envStr, "=")
			key := t[0]
			seqNum := strings.Replace(key, "SDP_MERGE_APP_NAME_", "", -1)
			appName := t[1]
			appId := os.Getenv("SDP_MERGE_APP_ID_" + seqNum)
			appDomains := os.Getenv("SDP_MERGE_APP_DOMAINS_" + seqNum)
			if strings.EqualFold(appId, "") || strings.EqualFold(appDomains, "") {
				logrus.Error("Tomcat 虚拟机域名配置异常：", os.Environ())
				os.Exit(-1)
			}
			h := HostDomains{
				APPNAME: appName,
				APPID:   appId,
				Domains: strings.Split(appDomains, ","),
			}
			if len(h.Domains) < 1 {
				logrus.Error("应用没有指定域名：", appName)
				os.Exit(-1)
			}

			sdpConfig.HostDomainConfig = append(sdpConfig.HostDomainConfig, h)
		}
	}
	if len(sdpConfig.HostDomainConfig) < 1 {
		logrus.Error("Host域名信息配置异常配置异常")
		os.Exit(-1)
	}

}

func modifyConnectorEntity(sdpConfig *SDP_CONFIG, server *etree.Document) {
	connector := server.SelectElement("Server").SelectElement("Service").SelectElement("Connector")
	if connector == nil {
		logrus.Error("解析server.xml模板文件异常")
		os.Exit(-1)
	}
	for _, attr := range connector.Attr {
		if !reserveConnectorArgs.Contains(attr.Key) {
			logrus.Warn("Connector 配置不存在删除：", attr.Key)
			connector.RemoveAttr(attr.Key)
		}
	}
	for k, v := range sdpConfig.ServerArgsDict {
		connector.CreateAttr(k, v)
	}
}

func modifyHostEntity(sdpConfig *SDP_CONFIG, server *etree.Document) {
	hosts := server.SelectElement("Server").SelectElement("Service").SelectElement("Engine").SelectElements("Host")
	if hosts == nil {
		logrus.Error("解析server.xml模板文件异常")
		os.Exit(-1)
	}
	for _, host := range hosts {
		valves := host.SelectElements("Valve")
		for _, valve := range valves {
			className := valve.SelectAttr("className")
			if className != nil && strings.EqualFold(className.Value, "org.apache.catalina.valves.AccessLogValve") {

				pattern := valve.SelectAttr("pattern")
				pattern.Value = sdpConfig.SdpServerAccessLog
			}
		}
	}
}

func generateHostEntity(sdpConfig *SDP_CONFIG, server *etree.Document) {

	hosts := server.SelectElement("Server").SelectElement("Service").SelectElement("Engine").SelectElements("Host")
	engine := server.SelectElement("Server").SelectElement("Service").SelectElement("Engine")
	if len(hosts) == 0 {
		logrus.Error("没有找到Host配置模板")
		os.Exit(-1)
	}
	defaultHost := hosts[0]
	for _, hostConfig := range sdpConfig.HostDomainConfig {
		newHost := defaultHost.Copy()
		appBase := newHost.SelectAttr("appBase")
		appBase.Value = "webapps_base/" + hostConfig.APPNAME
		hostname := newHost.SelectAttr("name")
		hostname.Value = strings.ToLower(hostConfig.APPNAME)
		for _, aliasHostName := range hostConfig.Domains {
			newHost.CreateElement("Alias").SetText(aliasHostName)
		}
		engine.AddChild(newHost)
	}
	// 配置默认虚拟机
	engine.RemoveChild(defaultHost)
	engine.SelectAttr("defaultHost").Value = "localhost"
}

func modifyServerXml(sdpConfig *SDP_CONFIG) {

	server := etree.NewDocument()
	if err := server.ReadFromFile(templateFile); err != nil {
		logrus.Error("读取配置文件模板失败：", err)
		os.Exit(-1)
	}
	generateHostEntity(sdpConfig, server)
	modifyConnectorEntity(sdpConfig, server)
	modifyHostEntity(sdpConfig, server)
	server.WriteSettings.CanonicalAttrVal = true
	if err := server.WriteToFile(outputFile); err != nil {
		logrus.Error("保存配置文件失败：", err)
		os.Exit(-1)
	}
}

func Main(cmd *cobra.Command, args []string) {

	var sdpConfig SDP_CONFIG
	initFromEnv(&sdpConfig)
	modifyServerXml(&sdpConfig)

}
