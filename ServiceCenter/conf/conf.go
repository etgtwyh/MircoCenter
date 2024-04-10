package conf

import (
	"encoding/json"
	"errors"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/ioc"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

const (
	AppName = "Conf"
)

func init() {
	//TODO 后续使用cli将项目目录传递进去
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		logger.Err(errors.New("项目根目录未找到"))
	}
	err := Load(projectRoot)
	if err != nil {
		panic(err)
	}
	ioc.Conf().Register(AppName, Conf)
}

type Config struct {
	GrpcServer  *GrpcServer
	MongoClient *mongo.Client
	EtcdClient  *clientv3.Client
	RedisClient *redis.Client
}

func (c *Config) String() string {
	indent, _ := json.MarshalIndent(c, "", " ")
	return string(indent)
}

type GrpcServer struct {
	MetricManageSvc  *MetricManageSvc  `toml:"MetricManageSvc"`
	ServiceManageSvc *ServiceManageSvc `toml:"ServiceManageSvc"`
}

type MetricManageSvc struct {
	Ipaddr string `toml:"ipaddr"`
	Port   string `toml:"port"`
}

type ServiceManageSvc struct {
	Ipaddr string `toml:"ipaddr"`
	Port   string `toml:"port"`
}

type MongoDB struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Ipaddr   string `toml:"ipaddr"`
	Port     string `toml:"port"`
}

type Etcd struct {
	Endpoints []string `toml:"endpoints"`
}

type Redis struct {
	Ipaddr string `toml:"ipaddr"`
	Port   string `toml:"port"`
}

func (c Config) Init() error {
	return nil
}
func (c Config) Destroy() error {
	return nil
}
