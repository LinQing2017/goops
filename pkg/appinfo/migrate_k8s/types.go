package migrate_k8s

type AppInformathion struct {
	AppId             string // 应用ID
	Name              string // 应用名称
	HostIP            string // 运行节点
	OtherIP           string // 其他实例的运行节点
	NodeSelectorLabel string // 运行标签
	Creator           string // 创建人
	URL               string // 连接url
	Metric            string
	Single            bool // 是否是单实例应用
}
