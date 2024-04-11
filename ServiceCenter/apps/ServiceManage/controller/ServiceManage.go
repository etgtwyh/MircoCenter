package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/MetricManage"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/ServiceManage"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

// TODO implement 获取服务信息
func (s *ServiceManageImpl) FetchServiceInfo(ctx context.Context) (*ServiceManage.AllServices, error) {
	// 定义返回数据实例
	allServices := ServiceManage.AllServices{}

	// 通过etcd client查找以"root"开头的所有服务
	resp, err := s.etcdClient.Get(ctx, "root", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	// 临时存储结构，用于构建最终的结构
	serviceInstancesMap := make(map[string]map[string][]string)

	// 遍历查询结果
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		// 分割键字符串以获取项目名、服务名和服务主机
		parts := strings.Split(key, "/")
		if len(parts) < 4 {
			continue // 如果分割后的部分不足以构成完整的项目名、服务名和服务主机，则跳过
		}

		projectName, serviceName, serviceHost := parts[1], parts[2], parts[3]

		// 确保数据结构的层级存在
		if _, ok := serviceInstancesMap[projectName]; !ok {
			serviceInstancesMap[projectName] = make(map[string][]string)
		}

		// 追加服务主机到对应的服务名下
		serviceInstancesMap[projectName][serviceName] = append(serviceInstancesMap[projectName][serviceName], serviceHost)
	}

	// 将临时存储结构转换为最终的AllServices结构
	for projectName, services := range serviceInstancesMap {
		projectMap := make(map[string]map[string][]string)
		projectMap[projectName] = services
		allServices.ServiceInstances = append(allServices.ServiceInstances, projectMap)
	}

	return &allServices, nil

}

// TODO implement http 服务注销---删除服务
func (s *ServiceManageImpl) UnRegisterService(ctx context.Context, Service *ServiceManage.Service) error {
	//通过etcd client查找是否存在该服务
	result, err := s.etcdClient.Get(ctx, fmt.Sprintf("root/%s/%s/%s", Service.ProjectName, Service.ServiceName, Service.ServiceHost))
	if err != nil {
		return err
	}
	if len(result.Kvs) == 0 {
		return errors.New("该服务不存在")
	}
	//存在则删除该服务
	//先向redis发送消息,停止上报数据和续约
	err = s.MsgPublish(ctx, fmt.Sprintf("root/%s/%s/%s", Service.ProjectName, Service.ServiceName, Service.ServiceHost))
	if err != nil {
		return err
	}
	//再删除etcd中的数据
	_, err = s.etcdClient.Delete(ctx, fmt.Sprintf("root/%s/%s/%s", Service.ProjectName, Service.ServiceName, Service.ServiceHost))
	if err != nil {
		//如果失败应该也向redis发送消息让微服务取消续约
		err = s.MsgPublish(ctx, fmt.Sprintf("root/%s/%s/%s", Service.ProjectName, Service.ServiceName, Service.ServiceHost))
		if err != nil {
			return err
		}
		return err
	}
	//最后删除mongodb中相应的数据
	collection := s.MongoClient.Database("ServiceManage").Collection("ProjectServiceMetrics")
	filter := bson.D{{
		"$and", bson.A{
			bson.D{{"projectname", Service.ProjectName}},
			bson.D{{"servicename", Service.ServiceName}},
			bson.D{{"servicehost", Service.ServiceHost}},
		}}}
	_, err = collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// TODO implement http 指定服务为最优
func (s *ServiceManageImpl) SetBestService(ctx context.Context, Service *ServiceManage.Service) error {
	type etcdValue struct {
		Op       int                       `json:"Op"`
		Addr     string                    `json:"Addr"`
		Metadata *MetricManage.MetricsData `json:"Metadata"`
	}
	var value etcdValue
	//先根据service获取到etcd的元数据
	result, err := s.etcdClient.Get(ctx, fmt.Sprintf("root/%s/%s/%s", Service.ProjectName, Service.ServiceName, Service.ServiceHost))
	if err != nil {
		return err
	}
	for _, kvs := range result.Kvs {
		err = json.Unmarshal(kvs.Value, &value)
		if err != nil {
			return err
		}
		if value.Metadata.IsPriority == "yes" {
			return errors.New("该服务已经是最优")
		} else {
			value.Metadata.IsPriority = "yes"
		}
	}
	//再更新数据
	grant, err := s.etcdClient.Grant(ctx, 10)
	if err != nil {
		return err
	}
	manager, err := endpoints.NewManager(s.etcdClient, "root")
	if err != nil {
		return err
	}
	err = manager.Update(ctx, []*endpoints.UpdateWithOpts{
		endpoints.NewAddUpdateOpts(fmt.Sprintf("root/%s/%s/%s", Service.ProjectName, Service.ServiceName, Service.ServiceHost), endpoints.Endpoint{
			Addr:     Service.ServiceHost,
			Metadata: value.Metadata,
		}, clientv3.WithLease(grant.ID)),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceManageImpl) UnSetBestService(ctx context.Context, Service *ServiceManage.Service) error {
	type etcdValue struct {
		Op       int                       `json:"Op"`
		Addr     string                    `json:"Addr"`
		Metadata *MetricManage.MetricsData `json:"Metadata"`
	}
	var value etcdValue
	//先根据service获取到etcd的元数据
	result, err := s.etcdClient.Get(ctx, fmt.Sprintf("root/%s/%s/%s", Service.ProjectName, Service.ServiceName, Service.ServiceHost))
	if err != nil {
		return err
	}
	for _, kvs := range result.Kvs {
		err = json.Unmarshal(kvs.Value, &value)
		if err != nil {
			return err
		}
		if value.Metadata.IsPriority == "no" {
			return errors.New("该服务已经不是最优")
		} else {
			value.Metadata.IsPriority = "no"
		}
	}
	//再更新数据
	grant, err := s.etcdClient.Grant(ctx, 10)
	if err != nil {
		return err
	}
	manager, err := endpoints.NewManager(s.etcdClient, "root")
	if err != nil {
		return err
	}
	err = manager.Update(ctx, []*endpoints.UpdateWithOpts{
		endpoints.NewAddUpdateOpts(fmt.Sprintf("root/%s/%s/%s", Service.ProjectName, Service.ServiceName, Service.ServiceHost), endpoints.Endpoint{
			Addr:     Service.ServiceHost,
			Metadata: value.Metadata,
		}, clientv3.WithLease(grant.ID)),
	})
	if err != nil {
		return err
	}
	return nil
}

// TODO 获取指定服务2小时的指标
func (s *ServiceManageImpl) FetchMetricsForServiceLastTwoHour(ctx context.Context, service *ServiceManage.Service) (*ServiceManage.TimeSeriesData, error) {
	// 实例化返回数据
	metrics := ServiceManage.TimeSeriesData{}

	// 调用MetricSvc获取过去两小时的指标数据
	hourMetrics, err := s.MetricSvc.FetchMetricsForServiceLastTwoHour(ctx, service)
	if err != nil {
		return nil, err
	}

	// 遍历hourMetrics中的MetricsData
	for _, md := range hourMetrics.MetricsData {
		// 将时间戳（秒）转换为time.Time
		timestamp := time.Unix(md.TimeStamp, 0)
		// 格式化时间戳并添加到结果中
		metrics.Timestamps = append(metrics.Timestamps, timestamp.Format("2006-01-02 15:04:05"))

		// 计算CPU使用率的平均值
		var cpuUsageSum float64
		for _, usage := range md.CpuUsage {
			cpuUsageSum += usage
		}
		averageCpuUsage := cpuUsageSum / float64(len(md.CpuUsage)) // 注意避免除以零
		metrics.CpuUsages = append(metrics.CpuUsages, averageCpuUsage)

		metrics.MemoryUsages = append(metrics.MemoryUsages, md.MemoryUsage)
		metrics.IcmpDelays = append(metrics.IcmpDelays, md.IcmpDelay)
	}
	if metrics.Timestamps == nil || metrics.CpuUsages == nil || metrics.IcmpDelays == nil {
		return nil, errors.New("请确认服务是否存在")
	}
	return &metrics, nil
}

// TODO 取指定服务的6小时内的指标
func (s *ServiceManageImpl) FetchMetricsForServiceLastSixHour(ctx context.Context, service *ServiceManage.Service) (*ServiceManage.TimeSeriesData, error) {
	// 实例化返回数据
	metrics := ServiceManage.TimeSeriesData{}

	hourMetrics, err := s.MetricSvc.FetchMetricsForServiceLastSixHour(ctx, service)
	if err != nil {
		return nil, err
	}

	// 遍历hourMetrics中的MetricsData
	for _, md := range hourMetrics.MetricsData {
		// 将时间戳（秒）转换为time.Time
		timestamp := time.Unix(md.TimeStamp, 0)
		// 格式化时间戳并添加到结果中
		metrics.Timestamps = append(metrics.Timestamps, timestamp.Format("2006-01-02 15:04:05"))

		// 计算CPU使用率的平均值
		var cpuUsageSum float64
		for _, usage := range md.CpuUsage {
			cpuUsageSum += usage
		}
		averageCpuUsage := cpuUsageSum / float64(len(md.CpuUsage)) // 注意避免除以零
		metrics.CpuUsages = append(metrics.CpuUsages, averageCpuUsage)

		metrics.MemoryUsages = append(metrics.MemoryUsages, md.MemoryUsage)
		metrics.IcmpDelays = append(metrics.IcmpDelays, md.IcmpDelay)
	}
	if metrics.Timestamps == nil || metrics.CpuUsages == nil || metrics.IcmpDelays == nil {
		return nil, errors.New("请确认服务是否存在")
	}
	return &metrics, nil

}

// TODO  获取指定服务12小时内的指标
func (s *ServiceManageImpl) FetchMetricsForServiceRecent(ctx context.Context, service *ServiceManage.Service) (*ServiceManage.TimeSeriesData, error) {
	// 实例化返回数据
	metrics := ServiceManage.TimeSeriesData{}

	// 调用MetricSvc获取过去两小时的指标数据
	hourMetrics, err := s.MetricSvc.FetchMetricsForServiceRecent(ctx, service)
	if err != nil {
		return nil, err
	}

	// 遍历hourMetrics中的MetricsData
	for _, md := range hourMetrics.MetricsData {
		// 将时间戳（秒）转换为time.Time
		timestamp := time.Unix(md.TimeStamp, 0)
		// 格式化时间戳并添加到结果中
		metrics.Timestamps = append(metrics.Timestamps, timestamp.Format("2006-01-02 15:04:05"))

		// 计算CPU使用率的平均值
		var cpuUsageSum float64
		for _, usage := range md.CpuUsage {
			cpuUsageSum += usage
		}
		averageCpuUsage := cpuUsageSum / float64(len(md.CpuUsage)) // 注意避免除以零
		metrics.CpuUsages = append(metrics.CpuUsages, averageCpuUsage)

		metrics.MemoryUsages = append(metrics.MemoryUsages, md.MemoryUsage)
		metrics.IcmpDelays = append(metrics.IcmpDelays, md.IcmpDelay)
	}
	if metrics.Timestamps == nil || metrics.CpuUsages == nil || metrics.IcmpDelays == nil {
		return nil, errors.New("请确认服务是否存在")
	}
	return &metrics, nil

}

// TODO 获取指定服务24小时内的指标
func (s *ServiceManageImpl) FetchServiceMetricsForLastDay(ctx context.Context, service *ServiceManage.Service) (*ServiceManage.TimeSeriesData, error) {
	// 实例化返回数据
	metrics := ServiceManage.TimeSeriesData{}

	hourMetrics, err := s.MetricSvc.FetchServiceMetricsForLastDay(ctx, service)
	if err != nil {
		return nil, err
	}

	// 遍历hourMetrics中的MetricsData
	for _, md := range hourMetrics.MetricsData {
		// 将时间戳（秒）转换为time.Time
		timestamp := time.Unix(md.TimeStamp, 0)
		// 格式化时间戳并添加到结果中
		metrics.Timestamps = append(metrics.Timestamps, timestamp.Format("2006-01-02 15:04:05"))

		// 计算CPU使用率的平均值
		var cpuUsageSum float64
		for _, usage := range md.CpuUsage {
			cpuUsageSum += usage
		}
		averageCpuUsage := cpuUsageSum / float64(len(md.CpuUsage)) // 注意避免除以零
		metrics.CpuUsages = append(metrics.CpuUsages, averageCpuUsage)

		metrics.MemoryUsages = append(metrics.MemoryUsages, md.MemoryUsage)
		metrics.IcmpDelays = append(metrics.IcmpDelays, md.IcmpDelay)
	}
	if metrics.Timestamps == nil || metrics.CpuUsages == nil || metrics.IcmpDelays == nil {
		return nil, errors.New("请确认服务是否存在")
	}
	return &metrics, nil
}

func (s *ServiceManageImpl) FetchServiceStatus(ctx context.Context, service *ServiceManage.Service) (*ServiceManage.IsPriority, error) {
	//定义实例用于接受数据
	type EtcdValue struct {
		Op       int                       `json:"Op"`
		Addr     string                    `json:"Addr"`
		Metadata *MetricManage.MetricsData `json:"Metadata"`
	}
	var etcdValue EtcdValue
	resp, err := s.etcdClient.Get(ctx, fmt.Sprintf("root/%s/%s/%s", service.ProjectName, service.ServiceName, service.ServiceHost))
	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		if err := json.Unmarshal(kv.Value, &etcdValue); err != nil {
			return nil, err
		}
	}
	if etcdValue.Metadata == nil {
		return nil, errors.New("未查询到服务")
	}
	// 如果IsPriority为true，立即返回这个实例
	if etcdValue.Metadata.IsPriority == "yes" {
		return &ServiceManage.IsPriority{Priority: true}, nil
	} else {
		return &ServiceManage.IsPriority{Priority: false}, nil
	}
}

func (s *ServiceManageImpl) MsgPublish(ctx context.Context, msg string) error {
	err := s.redisClient.Publish(ctx, "ServiceUnRegister", msg).Err()
	if err != nil {
		return err
	}
	return nil
}
