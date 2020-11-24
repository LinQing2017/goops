package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type App struct {
	ID             primitive.ObjectID `bson:"_id"`
	Name           string             `bson:"name"`
	Type           int                `bson:"type"`
	Creator        int                `bson:"creator"`
	CreatorName    string             `bson:"creatorName"`
	Members        []int              `bson:"members"`
	SingleInstance bool               `bson:"singleInstance"`
}

type Environment struct {
	ID             string       `bson:"_id"`
	EnvDesc        string       `bson:"envDesc"`
	AppId          string       `bson:"appId"`
	EnvType        int          `bson:"envType"`
	EnvName        string       `bson:"envName"`
	EnvChineseName string       `bson:"envChineseName"`
	ServicesList   []ServicesId `bson:"services"`
}

type ServicesId struct {
	ID   string `bson:"_id"`
	TYPE int    `bson:"type"`
	ENV  int    `bson:"env"`
}

type Services struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	ClusterId string `bson:"clusterId"`
	Type      string `bson:"type"`
	Env       string `bson:"env"`
}

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
