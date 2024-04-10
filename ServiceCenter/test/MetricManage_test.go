package test

import (
	"context"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/apps/MetricManage"
	_ "gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/apps/MetricManage/controller"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/apps/ServiceManage"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/ioc"
	"testing"
)

func init() {
	err := ioc.Init()
	if err != nil {
		panic(err)
	}
}
func TestMetricManageImpl_FetchMetricsForServiceLastTwoHour(t *testing.T) {
	service := ioc.Controller().Get(MetricManage.AppName).(MetricManage.MetricService)
	hour, err := service.FetchMetricsForServiceLastTwoHour(context.Background(), &ServiceManage.Service{
		ProjectName: "ServiceTest",
		ServiceName: "test",
		ServiceHost: "127.0.0.1:8080",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(hour)
}
