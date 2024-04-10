package ServiceManage

import (
	"context"
)

const (
	AppName = "ServiceManage"
)

type ServiceManage interface {
	//TODO implement 获取服务状态
	FetchServiceStatus(ctx context.Context, service *Service) (*IsPriority, error)
	//TODO implement 获取服务信息
	FetchServiceInfo(ctx context.Context) (*AllServices, error)
	//TODO implement http 服务注销
	UnRegisterService(ctx context.Context, Service *Service) error
	//TODO implement http 指定服务为最优
	SetBestService(ctx context.Context, Service *Service) error
	//TODO implement http 指定服务为不是优
	UnSetBestService(ctx context.Context, Service *Service) error
	//TODO 获取指定服务2小时的指标
	FetchMetricsForServiceLastTwoHour(ctx context.Context, service *Service) (*TimeSeriesData, error)
	//TODO 取指定服务的6小时内的指标
	FetchMetricsForServiceLastSixHour(ctx context.Context, service *Service) (*TimeSeriesData, error)
	//TODO  获取指定服务12小时内的指标
	FetchMetricsForServiceRecent(ctx context.Context, service *Service) (*TimeSeriesData, error)
	//TODO 获取指定服务24小时内的指标
	FetchServiceMetricsForLastDay(ctx context.Context, service *Service) (*TimeSeriesData, error)
}
