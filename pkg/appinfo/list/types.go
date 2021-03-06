package list

type PrintAppInfo struct {
	NAME       string
	APPID      string
	Creator    string
	EWSNum     int
	K8SNum     int
	MigrateMsg string
}

type PrintAppURL struct {
	NAME string
	URL  string
}

type PrintClusterInfo struct {
	APPNAME    string
	NAME       string
	ID         string
	Creator    string
	Memory     string
	K8SArea    string
	NodeType   string
	NodeName   string
	PackageURL string
	Instances  int
	Version    string
	MigrateMsg string
}
