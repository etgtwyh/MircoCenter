package ServiceUtils

import (
	"context"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/ServiceManage"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceUtils interface {
	//TODO implement 用户调用 服务注册
	RegisterService(ctx context.Context, service *ServiceManage.Service) error
	//TODO implement 用户调用 服务续约
	RenewService(ctx context.Context, leaseId clientv3.LeaseID)
}
