package controller

import (
	"errors"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/apps/MetricManage"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/apps/ServiceManage"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/conf"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/ioc"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	ioc.Controller().Register(ServiceManage.AppName, &ServiceManageImpl{})
}

type ServiceManageImpl struct {
	//etcd客户端
	etcdClient *clientv3.Client
	//redis客户端
	redisClient *redis.Client
	//Mongo客户端
	MongoClient *mongo.Client
	//依赖指标
	MetricSvc MetricManage.MetricService
	//微服务端口和地址
	ServerIpaddr string
	ServerPort   string
}

func (s *ServiceManageImpl) Init() error {
	if config, ok := ioc.Conf().Get(conf.AppName).(*conf.Config); ok {
		s.etcdClient = config.EtcdClient
		s.redisClient = config.RedisClient
		s.MongoClient = config.MongoClient
		s.ServerIpaddr = config.GrpcServer.ServiceManageSvc.Ipaddr
		s.ServerPort = config.GrpcServer.ServiceManageSvc.Port
	} else {
		return errors.New("ServiceManage无法获取相应依赖")
	}

	if MetricSvc, ok := ioc.Controller().Get(MetricManage.AppName).(MetricManage.MetricService); ok {
		s.MetricSvc = MetricSvc
	}
	return nil
}

func (s *ServiceManageImpl) Destroy() error {
	return nil
}
