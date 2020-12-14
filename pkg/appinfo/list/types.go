package list

type PrintAppInfo struct {
	NAME   string
	APPID  string
	Single string
	K8SNum int // k8s集群数目
	EWSNum int // 弹性web数目
	URL    string
}

type PrintEWSInfo struct {
	NAME       string
	Num        int // k8s集群数目
	PackageURL string
}

type PrintK8SInfo struct {
	NAME      string
	Num       int // k8s集群数目
	PodStatus string
	NodeName  []string
	NodeType  []string
	K8SAREA   []string
}
