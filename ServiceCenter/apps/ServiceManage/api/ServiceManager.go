package api

import (
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/ServiceManage"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/execption"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/reponse"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/utils/Log"
	"github.com/gin-gonic/gin"
	"net/http"
)

var logger = Log.NewLogger("etc/log/ServiceManageApi.log", 2, 7, 100)

// //获取服务信息
func (s *ServiceManagerApi) GetServiceInfo(ctx *gin.Context) {
	Serviceinfo, err := s.ServiceManageSvc.FetchServiceInfo(ctx)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "获取服务信息失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	reponse.RepSuccess(ctx, struct {
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}{
		Msg:  "获取服务信息成功",
		Data: Serviceinfo,
	})
}

// 获取指定服务2小时的指标
func (s *ServiceManagerApi) MetricsForServiceLastTwoHour(ctx *gin.Context) {
	//初始化实例
	var service ServiceManage.Service
	//绑定数据
	service.ProjectName = ctx.Query("projectname")
	service.ServiceName = ctx.Query("servicename")
	service.ServiceHost = ctx.Query("servicehost")
	//查询指标
	hourMetrics, err := s.ServiceManageSvc.FetchMetricsForServiceLastTwoHour(ctx, &service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "获取指标信息失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	reponse.RepSuccess(ctx, struct {
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}{
		Msg:  "获取指标信息成功",
		Data: hourMetrics,
	})
}

// 获取指定服务的6小时内的指标
func (s *ServiceManagerApi) MetricsForServiceLastSixHour(ctx *gin.Context) {
	//初始化实例
	var service ServiceManage.Service
	//绑定数据
	service.ProjectName = ctx.Query("projectname")
	service.ServiceName = ctx.Query("servicename")
	service.ServiceHost = ctx.Query("servicehost")
	//查询指标
	hourMetrics, err := s.ServiceManageSvc.FetchMetricsForServiceLastSixHour(ctx, &service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "获取指标信息失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	reponse.RepSuccess(ctx, struct {
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}{
		Msg:  "获取指标信息成功",
		Data: hourMetrics,
	})
}

// 获取指定服务12小时内的指标
func (s *ServiceManagerApi) MetricsForServiceTwelveHour(ctx *gin.Context) {
	//初始化实例
	var service ServiceManage.Service
	//绑定数据
	service.ProjectName = ctx.Query("projectname")
	service.ServiceName = ctx.Query("servicename")
	service.ServiceHost = ctx.Query("servicehost")
	//查询指标
	hourMetrics, err := s.ServiceManageSvc.FetchMetricsForServiceRecent(ctx, &service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "获取指标信息失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	reponse.RepSuccess(ctx, struct {
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}{
		Msg:  "获取指标信息成功",
		Data: hourMetrics,
	})
}

// 获取指定服务24小时内的指标
func (s *ServiceManagerApi) MetricsForServiceLastDay(ctx *gin.Context) {
	//初始化实例
	var service ServiceManage.Service
	//绑定数据
	service.ProjectName = ctx.Query("projectname")
	service.ServiceName = ctx.Query("servicename")
	service.ServiceHost = ctx.Query("servicehost")
	//查询指标
	hourMetrics, err := s.ServiceManageSvc.FetchServiceMetricsForLastDay(ctx, &service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "获取指标信息失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	reponse.RepSuccess(ctx, struct {
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}{
		Msg:  "获取指标信息成功",
		Data: hourMetrics,
	})
}

// 服务注销
func (s *ServiceManagerApi) UnRegisterService(ctx *gin.Context) {
	//初始化实例
	var service ServiceManage.Service
	//绑定数据
	err := ctx.Bind(&service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "服务注销失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	//服务注销
	err = s.ServiceManageSvc.UnRegisterService(ctx, &service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "服务注销失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	reponse.RepSuccess(ctx, struct {
		Msg string `json:"msg"`
	}{Msg: "服务注销成功"})
}

// 指定服务为最优
func (s *ServiceManagerApi) SetBestService(ctx *gin.Context) {
	//初始化实例
	var service ServiceManage.Service
	//绑定数据
	err := ctx.Bind(&service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "优先级设置失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	//设置该服务为最优
	err = s.ServiceManageSvc.SetBestService(ctx, &service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "优先级设置失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	reponse.RepSuccess(ctx, struct {
		Msg string `json:"msg"`
	}{Msg: "优先级设置成功"})
}

// 还原服务优先级
func (s *ServiceManagerApi) UnSetBestService(ctx *gin.Context) {
	//初始化实例
	var service ServiceManage.Service
	//绑定数据
	err := ctx.Bind(&service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "优先级还原失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	//还原服务优先级
	err = s.ServiceManageSvc.UnSetBestService(ctx, &service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "优先级还原失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	reponse.RepSuccess(ctx, struct {
		Msg string `json:"msg"`
	}{Msg: "优先级还原成功"})
}

func (s *ServiceManagerApi) GetServiceStatus(ctx *gin.Context) {
	//初始化实例
	var service ServiceManage.Service
	//绑定数据
	service.ProjectName = ctx.Query("projectname")
	service.ServiceName = ctx.Query("servicename")
	service.ServiceHost = ctx.Query("servicehost")
	status, err := s.ServiceManageSvc.FetchServiceStatus(ctx, &service)
	if err != nil {
		logger.Err(err)
		exception := execption.NewApiException(http.StatusBadGateway, "获取服务状态失败:", err.Error())
		reponse.RepFailed(ctx, exception)
		return
	}
	reponse.RepSuccess(ctx, struct {
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}{
		Msg:  "获取指标信息成功",
		Data: status,
	})
}
