package get

type EWSInfo struct {
	AppId   string // 应用ID
	Name    string // 应用名称
	Creator string // 创建人
	PackURL string // 包地址
	Single  bool   // 是否是单实例应用
}
