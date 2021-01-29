package portal

type DomainSwitch struct {
	Env       string   `json:"env"`
	ClusterId string   `json:"clusterId"`
	Domains   []Domain `json:"domains"`
}

type Domain struct {
	Protocol    string `json:"protocol"`
	Domain      string `json:"domain"`
	ClusterName string `json:"clusterName"`
	ClusterId   string `json:"clusterId"`
	Type        string `json:"type"`
	DomainId    string `json:"domainId"`
}

type Req4Portal struct {
	Type   string      `json:"type"`
	Params interface{} `json:"params"`
}
