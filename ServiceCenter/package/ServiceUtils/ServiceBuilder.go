package ServiceUtils

import (
	"google.golang.org/grpc/resolver"
	"strings"
)

const Scheme = "wcenter"

type ServiceBuilder struct{}

func parseTarget(endpoint string) (projectName, serviceName string) {
	// 这里需要根据实际的endpoint格式进行解析，下面是一个简单的示例
	parts := strings.Split(endpoint, "/")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

func (*ServiceBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 通过target.Endpoint获取服务名称和项目名称
	projectName, serviceName := parseTarget(target.Endpoint())
	// 创建ServiceUtilsImpl实例
	s := &ServiceUtilsImpl{
		MetricMonitor: util.MetricMonitor,
		EtcdClient:    util.EtcdClient,
		RedisClient:   util.RedisClient,
		MongoClient:   util.MongoClient,
		clientConn:    cc,
		serviceName:   serviceName,
		projectName:   projectName,
	}
	go s.resolve()
	return s, nil
}

func (*ServiceBuilder) Scheme() string {
	return Scheme
}
