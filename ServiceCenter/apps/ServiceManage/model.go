package ServiceManage

import (
	"encoding/json"
)

type AllServices struct {
	ServiceInstances []map[string]map[string][]string `json:"service_instances"` //map[projectname: string]{map[servicename: string]:[{service: string...}]}
}

type Service struct {
	ProjectName string `json:"project_name"` //项目名
	ServiceName string `json:"service_name"` //服务名
	ServiceHost string `json:"service_host"` //服务地址
}

type IsPriority struct {
	Priority bool `json:"priority"`
}

func (a *AllServices) String() string {
	indent, _ := json.MarshalIndent(a, "", " ")
	return string(indent)
}

type TimeSeriesData struct {
	Timestamps   []string  `json:"timestamps"`
	CpuUsages    []float64 `json:"cpuUsages"`
	MemoryUsages []float64 `json:"memoryUsages"`
	IcmpDelays   []int64   `json:"icmpDelays"`
}

func (a *TimeSeriesData) String() string {
	indent, _ := json.MarshalIndent(a, "", " ")
	return string(indent)
}
