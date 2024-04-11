package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/MetricManage"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/ServiceManage"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"net"
	"time"
)

// TODO 获取指定服务2小时的指标
func (m *MetricManageImpl) FetchMetricsForServiceLastTwoHour(ctx context.Context, service *ServiceManage.Service) (*MetricManage.ServiceInstanceMetrics, error) {
	//定义MetricData实例用于接受数据
	metricsData := []*MetricManage.MetricsData{}
	//定义ServiceInstance用于组合数据
	ServiceInstanceMetrics := &MetricManage.ServiceInstanceMetrics{}
	// 获取集合
	collection := m.MongoClient.Database("ServiceManage").Collection("ProjectServiceMetrics")

	// 计算24小时前的时间戳
	twentyFourHoursAgo := time.Now().Add(-2 * time.Hour).Unix()

	// 使用 bson.D 构建查询条件
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{Key: "projectname", Value: service.ProjectName}},
			bson.D{{"servicename", service.ServiceName}},
			bson.D{{"servicehost", service.ServiceHost}},
			bson.D{{"timestamp", bson.D{{"$gte", twentyFourHoursAgo}}}},
		}},
	}

	// 查询操作，使用传入的 ctx
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历查询结果
	for cursor.Next(ctx) {
		var metric MetricManage.MetricsData
		if err := cursor.Decode(&metric); err != nil {
			return nil, err
		}
		metricsData = append(metricsData, &metric)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	ServiceInstanceMetrics.ServiceHost = service.ServiceHost
	ServiceInstanceMetrics.MetricsData = metricsData
	return ServiceInstanceMetrics, nil
}

// TODO 取指定服务的6小时内的指标
func (m *MetricManageImpl) FetchMetricsForServiceLastSixHour(ctx context.Context, service *ServiceManage.Service) (*MetricManage.ServiceInstanceMetrics, error) {
	//定义MetricData实例用于接受数据
	metricsData := []*MetricManage.MetricsData{}
	//定义ServiceInstance用于组合数据
	ServiceInstanceMetrics := &MetricManage.ServiceInstanceMetrics{}
	// 获取集合
	collection := m.MongoClient.Database("ServiceManage").Collection("ProjectServiceMetrics")

	// 计算1小时前的时间戳
	oneHourAgo := time.Now().Add(-6 * time.Hour).Unix()

	// 使用 bson.D 构建查询条件
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{Key: "projectname", Value: service.ProjectName}},
			bson.D{{"servicename", service.ServiceName}},
			bson.D{{"servicehost", service.ServiceHost}},
			bson.D{{"timestamp", bson.D{{"$gte", oneHourAgo}}}},
		}},
	}

	// 查询操作，使用传入的 ctx
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历查询结果
	for cursor.Next(ctx) {
		var metric MetricManage.MetricsData
		if err := cursor.Decode(&metric); err != nil {
			return nil, err
		}
		metricsData = append(metricsData, &metric)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	//组合数据
	ServiceInstanceMetrics.ServiceHost = service.ServiceHost
	ServiceInstanceMetrics.MetricsData = metricsData
	return ServiceInstanceMetrics, nil
}

// TODO  获取指定服务12小时内的指标
func (m *MetricManageImpl) FetchMetricsForServiceRecent(ctx context.Context, service *ServiceManage.Service) (*MetricManage.ServiceInstanceMetrics, error) {
	//定义MetricData实例用于接受数据
	metricsData := []*MetricManage.MetricsData{}
	//定义ServiceInstance用于组合数据
	ServiceInstanceMetrics := &MetricManage.ServiceInstanceMetrics{}
	// 获取集合
	collection := m.MongoClient.Database("ServiceManage").Collection("ProjectServiceMetrics")

	// 计算24小时前的时间戳
	twentyFourHoursAgo := time.Now().Add(-12 * time.Hour).Unix()

	// 使用 bson.D 构建查询条件
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{Key: "projectname", Value: service.ProjectName}},
			bson.D{{"servicename", service.ServiceName}},
			bson.D{{"servicehost", service.ServiceHost}},
			bson.D{{"timestamp", bson.D{{"$gte", twentyFourHoursAgo}}}},
		}},
	}

	// 查询操作，使用传入的 ctx
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历查询结果
	for cursor.Next(ctx) {
		var metric MetricManage.MetricsData
		if err := cursor.Decode(&metric); err != nil {
			return nil, err
		}
		metricsData = append(metricsData, &metric)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	ServiceInstanceMetrics.ServiceHost = service.ServiceHost
	ServiceInstanceMetrics.MetricsData = metricsData
	return ServiceInstanceMetrics, nil
}

func (m *MetricManageImpl) FetchServiceMetricsForLastDay(ctx context.Context, service *ServiceManage.Service) (*MetricManage.ServiceInstanceMetrics, error) {
	//定义MetricData实例用于接受数据
	metricsData := []*MetricManage.MetricsData{}
	//定义ServiceInstance用于组合数据
	ServiceInstanceMetrics := &MetricManage.ServiceInstanceMetrics{}
	// 获取集合
	collection := m.MongoClient.Database("ServiceManage").Collection("ProjectServiceMetrics")

	// 计算24小时前的时间戳
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour).Unix()

	// 使用 bson.D 构建查询条件
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{Key: "projectname", Value: service.ProjectName}},
			bson.D{{"servicename", service.ServiceName}},
			bson.D{{"servicehost", service.ServiceHost}},
			bson.D{{"timestamp", bson.D{{"$gte", twentyFourHoursAgo}}}},
		}},
	}

	// 查询操作，使用传入的 ctx
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历查询结果
	for cursor.Next(ctx) {
		var metric MetricManage.MetricsData
		if err := cursor.Decode(&metric); err != nil {
			return nil, err
		}
		metricsData = append(metricsData, &metric)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	ServiceInstanceMetrics.ServiceHost = service.ServiceHost
	ServiceInstanceMetrics.MetricsData = metricsData
	return ServiceInstanceMetrics, nil
}

// ReportMetrics(MetricsService_ReportMetricsServer) error
func (m *MetricManageImpl) ReportMetrics(report MetricManage.MetricsService_ReportMetricsServer) error {
	collection := m.MongoClient.Database("ServiceManage").Collection("ProjectServiceMetrics")
	//添加TTL索引
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bson.M{"timestamp": 1}, Options: options.Index().SetExpireAfterSeconds(88200)})
	if err != nil {
		logger.Err(err).Msg("添加TTL索引失败")
		return err
	}
	for {
		recv, err := report.Recv()
		if err != nil {
			logger.Err(err).Msg("获取客户端指标失败")
			return err
		}
		logger.Info().Msg(recv.String())
		err = m.UpdateMetrics(context.Background(), recv, clientv3.LeaseID(recv.LeaseId))
		if err != nil {
			logger.Err(err).Msg("etcd更新元数据失败")
			return err
		}
		_, err = collection.InsertOne(context.Background(), recv)
		if err != nil {
			logger.Err(err).Msg("ServiceManage插入服务指标失败")
			return err
		}
		err = report.SendMsg(&MetricManage.ReportResponse{
			Success: true,
			Message: "report success",
		})
		if err != nil {
			logger.Err(err).Msg("ServiceManage返回消息失败")
			return err
		}
	}
}

func (m *MetricManageImpl) UpdateMetrics(ctx context.Context, Metrics *MetricManage.MetricsData, id clientv3.LeaseID) error {
	type etcdValue struct {
		Op       int                       `json:"Op"`
		Addr     string                    `json:"Addr"`
		Metadata *MetricManage.MetricsData `json:"Metadata"`
	}

	// 构建查询路径
	path := fmt.Sprintf("root/%s/%s/%s", Metrics.ProjectName, Metrics.ServiceName, Metrics.ServiceHost)

	// 先查询是否已存在该路径的数据
	resp, err := m.EtcdClient.Get(ctx, path)
	if err != nil {
		return err
	}
	//如果有数据,覆盖默认优先级
	if len(resp.Kvs) > 0 {
		var value etcdValue
		err := json.Unmarshal(resp.Kvs[0].Value, &value)
		if err != nil {
			return err
		}
		if value.Metadata != nil {
			Metrics.IsPriority = value.Metadata.IsPriority
		}
	}

	manager, err := endpoints.NewManager(m.EtcdClient, "root")
	if err != nil {
		return err
	}
	err = manager.Update(ctx, []*endpoints.UpdateWithOpts{
		endpoints.NewAddUpdateOpts(fmt.Sprintf("root/%s/%s/%s", Metrics.ProjectName, Metrics.ServiceName, Metrics.ServiceHost), endpoints.Endpoint{
			Addr:     Metrics.ServiceHost,
			Metadata: Metrics,
		}, clientv3.WithLease(id)),
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *MetricManageImpl) StartServer() {
	logger.Info().Msg("ServiceManage正在启动")
	server := grpc.NewServer()
	MetricManage.RegisterMetricsServiceServer(server, m)
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", m.IpAddr, m.Port))
	defer listen.Close()
	defer server.GracefulStop()
	if err != nil {
		logger.Err(err).Msg("ServiceManage启动失败")
		panic(err)
	}
	logger.Info().Msg("ServiceManage启动成功....")
	err = server.Serve(listen)
	if err != nil {
		logger.Err(err).Msg("ServiceManage启动失败")
		panic(err)
	}
}
