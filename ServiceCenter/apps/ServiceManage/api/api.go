package api

import (
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/apps/ServiceManage"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/ioc"
	"github.com/gin-gonic/gin"
)

const AppName = "ServiceManageApi"

func init() {
	ioc.Api().Register(AppName, &ServiceManagerApi{})
}

type ServiceManagerApi struct {
	ServiceManageSvc ServiceManage.ServiceManage
}

func (s *ServiceManagerApi) RegisterApi(router gin.IRouter) error {
	group := router.Group("/ServiceManager")
	//获取服务信息
	group.GET("", s.GetServiceInfo)
	//获取服务状态
	group.GET("/ServiceStatus", s.GetServiceStatus)
	//	获取指定服务2小时的指标 query: /Metrics/TwoHour?projectname=xx&servicename=xx&servicehost=xx
	group.GET("/Metrics/TwoHour", s.MetricsForServiceLastTwoHour)
	//	获取指定服务的6小时内的指标
	group.GET("/Metrics/SixHour", s.MetricsForServiceLastSixHour)
	//	获取指定服务12小时内的指标
	group.GET("/Metrics/TwelveHour", s.MetricsForServiceTwelveHour)
	//	获取指定服务24小时内的指标
	group.GET("/Metrics/LastDay", s.MetricsForServiceLastDay)
	//  服务注销
	group.DELETE("/UnRegisterService", s.UnRegisterService)
	//  指定服务为最优
	group.POST("/SetBestService", s.SetBestService)
	// 指定服务不是最优
	group.DELETE("/UnSetBestService", s.UnSetBestService)
	return nil
}

func (s *ServiceManagerApi) Init() error {
	if manage, ok := ioc.Controller().Get(ServiceManage.AppName).(ServiceManage.ServiceManage); ok {
		s.ServiceManageSvc = manage
	}
	return nil
}

func (s *ServiceManagerApi) Destroy() error {
	return nil
}
