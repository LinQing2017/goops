package list

type ListAppInfo struct {
	NAME       string
	APPID      string
	Single     bool
	EWSPackage bool // 能否找到包地址
	K8SNum     int  // k8s集群数目
	EWSNum     int  // 弹性web数目
	URL        string
}
