package MetricManage

import (
	"context"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/apps/ServiceManage"
)

const (
	AppName = "MetricManager"
)

type MetricService interface {
	//TODO 获取指定服务2小时的指标
	FetchMetricsForServiceLastTwoHour(ctx context.Context, service *ServiceManage.Service) (*ServiceInstanceMetrics, error)
	//TODO 取指定服务的6小时内的指标
	FetchMetricsForServiceLastSixHour(ctx context.Context, service *ServiceManage.Service) (*ServiceInstanceMetrics, error)
	//TODO  获取指定服务12小时内的指标
	FetchMetricsForServiceRecent(ctx context.Context, service *ServiceManage.Service) (*ServiceInstanceMetrics, error)
	//TODO 获取指定服务24小时内的指标
	FetchServiceMetricsForLastDay(ctx context.Context, service *ServiceManage.Service) (*ServiceInstanceMetrics, error)
	//GRPC server,上报各种指标
	MetricsServiceServer
}
