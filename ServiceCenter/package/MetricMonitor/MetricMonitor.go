package MetricMonitor

import (
	"context"
	"fmt"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/MetricManage"
	"github.com/go-ping/ping"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type MetricMonitorImpl struct {
	trick      *time.Ticker
	ServerAddr string
	ServerPort string
}

// 暴露给外部用户使用,
func NewMetricMonitor(ServerPort string, ServerAddr string) *MetricMonitorImpl {
	m := &MetricMonitorImpl{}
	m.trick = time.NewTicker(5 * time.Second)
	m.ServerPort = ServerPort
	m.ServerAddr = ServerAddr
	return m
}

func (m *MetricMonitorImpl) ReportMetrics(ctx context.Context, ProjectName, ServiceName, ServiceHost string, LeaseID int64) {
	dial, err := grpc.Dial(fmt.Sprintf("%s:%s", m.ServerAddr, m.ServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := MetricManage.NewMetricsServiceClient(dial)

	metricsStream, err := client.ReportMetrics(ctx)
	if err != nil {
		panic(err)
	}

	for now := range m.trick.C {
		select {
		case <-ctx.Done():
			return
		default:
			log.Printf("%s:正在获取服务指标", now.Format("2006-01-02 15:04:05"))
			//获取cpu指标
			percent, err := cpu.Percent(time.Second, false)
			if err != nil {
				log.Printf("%s:获取CPU指标失败:%s", now.Format("2006-01-02 15:04:05"), err.Error())
				return
			}
			//获取内存指标
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("%s:获取内存指标失败:%s", now.Format("2006-01-02 15:04:05"), err.Error())
				return
			}
			//icmp延时检测
			delay, err := m.Delay()
			if err != nil {
				log.Printf("%s:获取延时指标失败:%s", now.Format("2006-01-02 15:04:05"), err.Error())
				return
			}
			//组装指标数据
			MetricData := &MetricManage.MetricsData{
				ProjectName: ProjectName,
				ServiceName: ServiceName,
				ServiceHost: ServiceHost,
				CpuUsage:    percent,
				MemoryUsage: memInfo.UsedPercent,
				IcmpDelay:   delay,
				TimeStamp:   now.Unix(),
				IsPriority:  "no",
				LeaseId:     LeaseID,
			}
			//发送到服务端
			err = metricsStream.Send(MetricData)
			if err != nil {
				log.Printf("%s:发送指标到服务端失败:%s", now.Format("2006-01-02 15:04:05"), err.Error())
				return
			}
		}

	}
}

func (m *MetricMonitorImpl) Delay() (int64, error) {
	pinger, err := ping.NewPinger(m.ServerAddr)
	pinger.SetPrivileged(true)
	if err != nil {
		return 0, err
	}
	pinger.Count = 5
	err = pinger.Run()
	if err != nil {
		return 0, err
	}
	statistics := pinger.Statistics()
	return int64(statistics.AvgRtt), nil
}
