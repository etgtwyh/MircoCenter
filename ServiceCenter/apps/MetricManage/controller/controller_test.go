package controller

import (
	_ "gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/conf"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/ioc"
	"testing"
)

func init() {
	err := ioc.Init()
	if err != nil {
		panic(err)
	}
}

func TestReportMetrics(t *testing.T) {
	select {}
}
