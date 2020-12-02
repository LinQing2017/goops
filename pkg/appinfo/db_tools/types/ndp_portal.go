package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type ServicesId struct {
	ID   primitive.ObjectID `bson:"_id"`
	TYPE int                `bson:"type"`
	ENV  int                `bson:"env"`
}

type App struct {
	ID             primitive.ObjectID `bson:"_id"`
	Name           string             `bson:"name"`
	Desc           string             `bson:"desc"`
	Type           int                `bson:"type"`
	Creator        int                `bson:"creator"`
	CreatorName    string             `bson:"creatorName"`
	Members        []int              `bson:"members"`
	SingleInstance bool               `bson:"singleInstance"`
	Services       []ServicesId       `bson:"services"`
}

type Environment struct {
	ID             primitive.ObjectID `bson:"_id"`
	EnvDesc        string             `bson:"envDesc"`
	AppId          string             `bson:"appId"`
	EnvType        int                `bson:"envType"`
	EnvName        string             `bson:"envName"`
	EnvChineseName string             `bson:"envChineseName"`
	ServicesList   []ServicesId       `bson:"services"`
	ApplyReason    string             `bson:"applyReason"`
	ApplyState     int                `bson:"applyState"`
	PinPointStatus bool               `bson:"pinPointStatus"`
	IsQaCheck      bool               `bson:"isQaCheck"`
	IsPlanCheck    bool               `bson:"isPlanCheck"`
}

type Service struct {
	ID         primitive.ObjectID `bson:"_id"`
	ServerType int                `bson:"service_type"`
	JdkVersion string             `bson:"jdk_version"`
	Name       string             `bson:"name"`
	PreN       string             `bson:preN`
	ClusterId  string             `bson:"clusterId"`
	Type       int                `bson:"type"`
	Env        int                `bson:"env"`
	AppId      string             `bson:"appId`
}

type AppPortalInfo struct {
	APP            App
	Environments   []*Environment
	EWSServiceList []*Service
	K8SServiceList []*Service
}
