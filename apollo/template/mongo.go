package template

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	op "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type Mongo struct {
	Addresses []string `toml:"addresses"`
	User      string   `toml:"user"`
	Password  string   `toml:"password"`
}

//支持单机和集群
func (d *Mongo) Engine() (*mongo.Client, error) {
	clientoptions := op.Client()
	clientoptions.SetConnectTimeout(time.Duration(int(time.Second) * 10))

	//连接参数采用拼接的方式  可用通用化
	host := d.Addresses[0]
	//集群版需要安装此格式添加所有节点地址
	for i := 1; i < len(d.Addresses); i++ {
		host += ","
		host += d.Addresses[i]
	}
	mongoUrl := "mongodb://" + host + "/"
	clientoptions.ApplyURI(mongoUrl)
	Log.Info("Mongo Connection : "+mongoUrl, zap.String("中间件", "Mongo"))
	client, err := mongo.Connect(context.Background(), clientoptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(nil, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
