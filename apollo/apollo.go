package apollo

import (
	"bytes"
	"github.com/shima-park/agollo"
	"github.com/spf13/viper"
	"github.com/swhc/utils/template"
)

const (
	PublicNamespace = "develop.abroad"
)

type Apollo struct {
	Endpoint    string //节点
	AppID       string //服务ID
	ClusterName string //集群名
}

func NewApollo(endpoint, appID, cluster string) (*Apollo, error) {
	tem := &Apollo{
		Endpoint:    endpoint,
		AppID:       appID,
		ClusterName: cluster,
	}
	return tem, nil
}

// 获取到所有mysql地址
func (a *Apollo) Mysql() (map[string]*template.Mysql, error) {
	ap, err := agollo.New(a.Endpoint, a.AppID, agollo.Cluster(a.ClusterName), agollo.AutoFetchOnCacheMiss())
	if err != nil {
		return nil, err
	}
	content := ap.Get("mysql", agollo.WithNamespace(PublicNamespace))
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewReader([]byte(content))); err != nil {
		return nil, err
	}
	// 此时存在多个 key =》 mysql
	var result = make(map[string]*template.Mysql)
	if err := v.Unmarshal(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// 获取到所有redis的地址
func (a *Apollo) Redis() (map[string]*template.Redis, error) {
	ap, err := agollo.New(a.Endpoint, a.AppID, agollo.Cluster(a.ClusterName), agollo.AutoFetchOnCacheMiss())
	if err != nil {
		return nil, err
	}
	content := ap.Get("redis", agollo.WithNamespace(PublicNamespace))
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewReader([]byte(content))); err != nil {
		return nil, err
	}
	// 此时存在多个 key =》 mysql
	var result = make(map[string]*template.Redis)
	if err := v.Unmarshal(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// 获取到所有redis的地址
func (a *Apollo) Mongo() (map[string]*template.Mongo, error) {
	ap, err := agollo.New(a.Endpoint, a.AppID, agollo.Cluster(a.ClusterName), agollo.AutoFetchOnCacheMiss())
	if err != nil {
		return nil, err
	}
	content := ap.Get("mongo", agollo.WithNamespace(PublicNamespace))
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewReader([]byte(content))); err != nil {
		return nil, err
	}
	// 此时存在多个 key =》 mysql
	var result = make(map[string]*template.Mongo)
	if err := v.Unmarshal(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// 获取到所有redis的地址
func (a *Apollo) RabbitMq() (map[string]*template.RabbitMq, error) {
	ap, err := agollo.New(a.Endpoint, a.AppID, agollo.Cluster(a.ClusterName), agollo.AutoFetchOnCacheMiss())
	if err != nil {
		return nil, err
	}
	content := ap.Get("rabbitmq", agollo.WithNamespace(PublicNamespace))
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewReader([]byte(content))); err != nil {
		return nil, err
	}
	// 此时存在多个 key =》 mysql
	var result = make(map[string]*template.RabbitMq)
	if err := v.Unmarshal(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// mysql 配置获取
func (a *Apollo) LoadMatchMysql(key string) (map[string]string, error) {
	ap, err := agollo.New(a.Endpoint, a.AppID, agollo.Cluster(a.ClusterName), agollo.AutoFetchOnCacheMiss())
	if err != nil {
		return nil, err
	}
	result := ap.Get(key)
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewReader([]byte(result))); err != nil {
		return nil, err
	}
	//	具体服务名称
	tem := v.GetStringMapString("maoti_app")
	// 读取到第一个服务名:这里和
	return tem, nil
}

//导入一下非模板内容
func (a *Apollo) LoadConf(key string) (string, error) {
	ap, err := agollo.New(a.Endpoint, a.AppID, agollo.Cluster(a.ClusterName), agollo.AutoFetchOnCacheMiss())
	if err != nil {
		return "", err
	}
	result := ap.Get(key)
	// 读取到第一个服务名:这里和
	return result, nil
}

// 获取所有grpc 地址
func (a *Apollo) Grpc() (map[string]string, error) {
	// todo
	ap, err := agollo.New(a.Endpoint, a.AppID, agollo.Cluster(a.ClusterName), agollo.AutoFetchOnCacheMiss())
	if err != nil {
		return nil, err
	}
	grpc := ap.Get("grpc")
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewReader([]byte(grpc))); err != nil {
		return nil, err
	}
	//	具体服务名称
	tem := v.GetStringMapString("grpc")
	return tem, nil

}

func (a *Apollo) CustomConf(key string) (map[string]string, error) {
	ap, err := agollo.New(a.Endpoint, a.AppID, agollo.Cluster(a.ClusterName), agollo.AutoFetchOnCacheMiss())
	if err != nil {
		return nil, err
	}
	conf := ap.Get(key)
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewReader([]byte(conf))); err != nil {
		return nil, err
	}
	//	具体服务名称
	tem := v.GetStringMapString(key)
	return tem, nil
}

func (a *Apollo) Http() (map[string]string, error) {
	ap, err := agollo.New(a.Endpoint, a.AppID, agollo.Cluster(a.ClusterName), agollo.AutoFetchOnCacheMiss())
	if err != nil {
		return nil, err
	}
	port := ap.Get("port")
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewReader([]byte(port))); err != nil {
		return nil, err
	}
	//	具体服务名称
	tem := v.GetStringMapString("http")
	return tem, nil
}
