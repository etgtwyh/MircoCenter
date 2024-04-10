package MetricMonitor

import "context"

type MetricMonitor interface {
	//上送监控指标
	ReportMetrics(ctx context.Context, ProjectName, ServiceName, ServiceHost string, LeaseID int64)
}
