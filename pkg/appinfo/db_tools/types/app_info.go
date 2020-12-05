package types

type AppPortalInfo struct {
	APP            App
	Environments   []*Environment
	EWSServiceList []*Service
	K8SServiceList []*Service
}
