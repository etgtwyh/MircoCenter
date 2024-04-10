package controller

import (
	"errors"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/apps/MetricManage"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/conf"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/ioc"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/utils/Log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	var Impl = &MetricManageImpl{}
	ioc.Controller().Register(MetricManage.AppName, Impl)

}

var logger = Log.NewLogger("etc/log/MetricManage.log", 2, 7, 100)

type MetricManageImpl struct {
	//TODO etcd client
	EtcdClient *clientv3.Client
	//MongoDB client
	MongoClient *mongo.Client
	//Server addr
	IpAddr string
	//Server Port
	Port string
	//组合GRPC
	MetricManage.UnimplementedMetricsServiceServer
}

func (m *MetricManageImpl) Init() error {
	if config, ok := ioc.Conf().Get(conf.AppName).(*conf.Config); ok {
		m.MongoClient = config.MongoClient
		m.EtcdClient = config.EtcdClient
		m.Port = config.GrpcServer.MetricManageSvc.Port
		m.IpAddr = config.GrpcServer.ServiceManageSvc.Ipaddr
	} else {
		logger.Err(errors.New("指标中心获取MongoDB客户端失败"))
	}
	go m.StartServer()
	return nil
}

func (m *MetricManageImpl) Destroy() error {
	return nil
}
