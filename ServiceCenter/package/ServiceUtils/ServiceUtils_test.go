package ServiceUtils

import (
	"context"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/apps/ServiceManage"
	"testing"
)

func TestNewServiceUtils(t *testing.T) {
	utils, err := NewServiceUtils()
	if err != nil {
		t.Fatal(err)
	}
	err = utils.RegisterService(context.Background(), &ServiceManage.Service{
		ProjectName: "ServiceCenter",
		ServiceName: "MetricManage",
		ServiceHost: "127.0.0.1:8080",
	})
	if err != nil {
		t.Fatal(err)
	}
	select {}
}

func TestServiceUtilsImpl_FetchOptimalServiceMetricsForPrefix(t *testing.T) {
	utils, err := NewServiceUtils()
	if err != nil {
		t.Fatal(err)
	}
	//root/ServiceTest/test/127.0.0.1:8080
	prefix, err := utils.FetchOptimalServiceMetricsForPrefix(context.Background(), "ServiceTest", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(prefix)
}
