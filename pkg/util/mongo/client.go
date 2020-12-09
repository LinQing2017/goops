package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func MongoClient(uri string) *mongo.Client {
	// 设置客户端连接配置
	//clientOptions := options.Client().SetHosts([]string{"172.24.135.10:34002"}).SetAuth( options.Credential{
	//	AuthMechanism:           "SCRAM-SHA-1",
	//	AuthSource:              "ndp_portal",
	//	Username:                "ndp_portal",
	//	Password:                "okvlZuvSwsnZ",
	//}).SetDirect(true)
	clientOptions := options.Client().ApplyURI(uri).SetDirect(true)
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func MongoDisconnect(client *mongo.Client) {
	// 断开连接
	if client != nil {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}
}
