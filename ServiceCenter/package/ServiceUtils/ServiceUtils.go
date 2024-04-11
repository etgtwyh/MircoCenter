package ServiceUtils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/MetricManage"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/apps/ServiceManage"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/package/MetricMonitor"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/utils/Log"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc/resolver"
	"log"
	"math"
	"os"
	"time"
)

var util = &ServiceUtilsImpl{}

var logger = Log.NewLogger("etc/log/ServiceUtils.log", 2, 7, 100)

func init() {
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		panic(errors.New("项目根目录未找到"))
	}
	err := Load(projectRoot)
	if err != nil {
		panic(err)
	}
	utils, err := NewServiceUtils()
	if err != nil {
		panic(err)
	}
	util = utils
	resolver.Register(&ServiceBuilder{})
}

type ServiceUtilsImpl struct {
	MetricMonitor MetricMonitor.MetricMonitor
	EtcdClient    *clientv3.Client
	RedisClient   *redis.Client
	MongoClient   *mongo.Client
	clientConn    resolver.ClientConn
	serviceName   string
	projectName   string
}

type EtcdValue struct {
	Op       int                       `json:"Op"`
	Addr     string                    `json:"Addr"`
	Metadata *MetricManage.MetricsData `json:"Metadata"`
}

func (s *ServiceUtilsImpl) resolve() {
	instance, err := s.FetchOptimalServiceMetricsForPrefix(context.Background(), s.projectName, s.serviceName)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch optimal service metrics")
		return
	}

	addr := []resolver.Address{{Addr: instance.ServiceHost}}
	s.clientConn.UpdateState(resolver.State{Addresses: addr})
}

func NewServiceUtils() (*ServiceUtilsImpl, error) {
	client, err := clientv3.New(clientv3.Config{Endpoints: globalMetricManageSvc.ServiceManageEndpoints, DialTimeout: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	//测试etcd是否可用
	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	_, err = client.Get(timeout, "etcd-ping-test")
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", globalRedisConf.Ipaddr, globalRedisConf.Port),
	})
	//测试redis是否可用
	if err := redisClient.Ping(timeout).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s:%s, error: %v", globalRedisConf.Ipaddr, globalRedisConf.Port, err)
	}
	mongoCli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", globalMongoConf.Username, globalMongoConf.Password, globalMongoConf.Ipaddr, globalMongoConf.Port)))
	if err != nil {
		return nil, err
	}
	// 测试Mongodb是否可用
	if err = mongoCli.Ping(timeout, readpref.Primary()); err != nil {
		return nil, err
	}
	return &ServiceUtilsImpl{
		MetricMonitor: MetricMonitor.NewMetricMonitor(globalMetricManageSvc.MetricManagePort, globalMetricManageSvc.MetricManageIpAddr),
		EtcdClient:    client,
		RedisClient:   redisClient,
		MongoClient:   mongoCli,
	}, nil
}

// TODO implement 用户调用 服务注册
func (s *ServiceUtilsImpl) RegisterService(ctx context.Context, service *ServiceManage.Service) error {
	//创建一个服务管理器
	manager, err := endpoints.NewManager(s.EtcdClient, "root")
	if err != nil {
		return err
	}
	//创建一个10s租约
	grant, err := s.EtcdClient.Grant(ctx, 10)
	if err != nil {
		return err
	}
	//添加一个服务
	err = manager.AddEndpoint(ctx, fmt.Sprintf("root/%s/%s/%s", service.ProjectName, service.ServiceName, service.ServiceHost), endpoints.Endpoint{
		Addr: service.ServiceHost,
	}, clientv3.WithLease(grant.ID))
	if err != nil {
		return err
	}
	//生产关闭上下文
	cancelctx, cancelFunc := context.WithCancel(ctx)
	//使用redis监听服务注销主题
	ServiceUnRegister := s.RedisClient.Subscribe(ctx, "ServiceUnRegister")
	//defer ServiceUnRegister.Close()
	//启动自动续约
	go s.RenewService(cancelctx, grant.ID)
	//启动指标上报
	go s.MetricMonitor.ReportMetrics(cancelctx, service.ProjectName, service.ServiceName, service.ServiceHost, int64(grant.ID))
	//启动一个监听器,检测是否需要注销服务
	go s.UnRegisterService(cancelFunc, ServiceUnRegister, service)
	return nil
}

func (s *ServiceUtilsImpl) UnRegisterService(cancelFunc context.CancelFunc, ServiceUnRegister *redis.PubSub, service *ServiceManage.Service) {
	for {
		msg, err := ServiceUnRegister.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}
		//如果接收到消息,判断内容是否是该服务,如果是则取消续约,如果不是则继续监听
		if msg.Payload == fmt.Sprintf("root/%s/%s/%s", service.ProjectName, service.ServiceName, service.ServiceHost) {
			cancelFunc()
			return
		} else {
			continue
		}
	}
}

// TODO implement 用户调用 服务续约
// TODO implement 使用消息队列订阅服务注销的主题,监听到后应不再续约
func (s *ServiceUtilsImpl) RenewService(ctx context.Context, leaseId clientv3.LeaseID) {
	alive, err := s.EtcdClient.KeepAlive(ctx, leaseId)
	if err != nil {
		return
	}
	for resp := range alive {
		logger.Info().Msg(resp.String())
	}
}

func (s *ServiceUtilsImpl) FetchOptimalServiceMetricsForPrefix(ctx context.Context, ProjectName, ServiceName string) (*MetricManage.ServiceInstance, error) {

	//定义ServiceInstance用于组合数据
	ServiceInstance := &MetricManage.ServiceInstance{}
	//定义评分
	var minScore float64 = math.MaxFloat64
	resp, err := s.EtcdClient.Get(ctx, "root/"+ProjectName+"/"+ServiceName, clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("Failed to get data from etcd: %v", err)
	}
	for _, kv := range resp.Kvs {
		//定义实例用于接受数据
		var etcdValue EtcdValue
		if err := json.Unmarshal(kv.Value, &etcdValue); err != nil {
			log.Printf("Failed to unmarshal metric data: %v", err)
			continue
		}
		// 如果IsPriority为true，立即返回这个实例
		if etcdValue.Metadata.IsPriority == "yes" {
			return &MetricManage.ServiceInstance{
				ServiceHost: etcdValue.Metadata.ServiceHost,
				MetricsData: etcdValue.Metadata,
			}, nil
		}
		// 计算得分
		var totalCPUUsage float64
		for _, usage := range etcdValue.Metadata.CpuUsage {
			totalCPUUsage += usage
		}
		avgCPUUsage := totalCPUUsage / float64(len(etcdValue.Metadata.CpuUsage))
		score := avgCPUUsage + etcdValue.Metadata.MemoryUsage + float64(etcdValue.Metadata.IcmpDelay)
		if score < minScore {
			minScore = score
			ServiceInstance = &MetricManage.ServiceInstance{
				ServiceHost: etcdValue.Metadata.ServiceHost,
				MetricsData: etcdValue.Metadata, // 同上
			}
		}
	}
	if ServiceInstance.MetricsData == nil {
		return nil, fmt.Errorf("no suitable instance found")
	}

	return ServiceInstance, nil
}

func (s *ServiceUtilsImpl) ResolveNow(opts resolver.ResolveNowOptions) {
	// 实现留空，因为resolve方法已经启动了服务解析的逻辑
}

func (s *ServiceUtilsImpl) Close() {
	// 清理资源
}
