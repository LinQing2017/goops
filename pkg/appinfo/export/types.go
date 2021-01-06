package export

type ClusterInfo struct {
	AppName     string
	ClusterName string
	Creator     string
	Replica     int
	MaxReplica  int
	MemRequest  string
	MemLimit    string
	AppUrl      string
}
