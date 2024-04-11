package controller

import (
	"context"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/ServiceManage"

	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/ioc"

	"testing"
)

//func init() {
//	err := ioc.Init()
//	if err != nil {
//		panic(err)
//	}
//}

func TestServiceManageImpl_UnRegisterService(t *testing.T) {
	manage := ioc.Controller().Get(ServiceManage.AppName).(ServiceManage.ServiceManage)
	err := manage.UnRegisterService(context.Background(), &ServiceManage.Service{
		ProjectName: "ServiceTest",
		ServiceName: "test",
		ServiceHost: "127.0.0.1:8082",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestServiceManageImpl_SetBestService(t *testing.T) {
	manage := ioc.Controller().Get(ServiceManage.AppName).(ServiceManage.ServiceManage)
	err := manage.SetBestService(context.Background(), &ServiceManage.Service{
		ProjectName: "ServiceTest",
		ServiceName: "test",
		ServiceHost: "127.0.0.1:8082",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestServiceManageImpl_FetchMetricsForServiceLastTwoHour(t *testing.T) {
	err := ioc.Init()
	if err != nil {
		panic(err)
	}
	manage := ioc.Controller().Get(ServiceManage.AppName).(ServiceManage.ServiceManage)
	hourMetric, err := manage.FetchMetricsForServiceLastTwoHour(context.Background(), &ServiceManage.Service{
		ProjectName: "ServiceTest",
		ServiceName: "test",
		ServiceHost: "127.0.0.1:8082",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(hourMetric)
}

func TestServiceManageImpl_FetchServiceInfo(t *testing.T) {
	err := ioc.Init()
	if err != nil {
		panic(err)
	}
	manage := ioc.Controller().Get(ServiceManage.AppName).(ServiceManage.ServiceManage)
	info, err := manage.FetchServiceInfo(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
