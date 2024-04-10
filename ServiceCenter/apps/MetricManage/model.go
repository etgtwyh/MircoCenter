package MetricManage

import "encoding/json"

type ServiceInstance struct {
	ServiceHost string       `json:"service_host"`
	MetricsData *MetricsData `json:"metrics_data"`
}

type ServiceInstanceMetrics struct {
	ServiceHost string         `json:"service_host"`
	MetricsData []*MetricsData `json:"metrics_data"`
}

func (s *ServiceInstance) String() string {
	indent, _ := json.MarshalIndent(s, "", " ")
	return string(indent)
}

func (s *ServiceInstanceMetrics) String() string {
	indent, _ := json.MarshalIndent(s, "", " ")
	return string(indent)
}
