package conf

import (
	"context"
	"fmt"
	"gitee.com/King_of_Universe_Sailing/MircoCenter/ServiceCenter/utils/Log"
	"github.com/BurntSushi/toml"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var Conf = &Config{}
var logger = Log.NewLogger("etc/log/conf.log", 2, 7, 100)

func Load(projectRoot string) error {
	Mongo := new(MongoDB)
	EtcdCli := new(Etcd)
	GrpcS := new(GrpcServer)
	RedisCli := new(Redis)
	client, err := Mongo.Load(projectRoot)
	if err != nil {
		logger.Err(err).Msg("MongoDB连接失败")
		return err
	}
	etcdcli, err := EtcdCli.Load(projectRoot)
	if err != nil {
		logger.Err(err).Msg("Etcd连接失败")
		return err
	}
	GrpcSvc, err := GrpcS.Load(projectRoot)
	if err != nil {
		logger.Err(err)
		return err
	}
	Conf.RedisClient = RedisCli.Load(projectRoot)
	Conf.GrpcServer = GrpcSvc
	Conf.MongoClient = client
	Conf.EtcdClient = etcdcli
	return nil
}

func (m *MongoDB) Load(projectRoot string) (*mongo.Client, error) {
	mo := struct {
		MongoDB MongoDB `toml:"MongoDB"`
	}{}
	_, err := toml.DecodeFile(fmt.Sprintf("%s/etc/config.toml", projectRoot), &mo)
	if err != nil {
		return nil, err
	}
	m = &mo.MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	connect, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", m.Username, m.Password, m.Ipaddr, m.Port)))
	if err != nil {
		return nil, err
	}
	// Ping the primary
	if err = connect.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	logger.Info().Msg(fmt.Sprintf("MongoDB:%s:%s,连接成功", m.Ipaddr, m.Port))
	return connect, nil
}

func (e *Etcd) Load(projectRoot string) (*clientv3.Client, error) {
	et := struct {
		Etcd Etcd `toml:"Etcd"`
	}{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := toml.DecodeFile(fmt.Sprintf("%s/etc/config.toml", projectRoot), &et)
	if err != nil {
		return nil, err
	}
	e = &et.Etcd
	client, err := clientv3.New(clientv3.Config{Endpoints: e.Endpoints, DialTimeout: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	_, err = client.Get(ctx, "etcd-ping-test")
	if err != nil {
		return nil, err
	}
	logger.Info().Msg(fmt.Sprintf("Etcd:%s,连接成功", e.Endpoints))
	return client, nil
}

func (g *GrpcServer) Load(projectRoot string) (*GrpcServer, error) {
	server := &GrpcServer{
		MetricManageSvc:  &MetricManageSvc{},
		ServiceManageSvc: &ServiceManageSvc{},
	}
	_, err := toml.DecodeFile(fmt.Sprintf("%s/etc/config.toml", projectRoot), server)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (r *Redis) Load(projectRoot string) *redis.Client {
	re := struct {
		Redis Redis `toml:"Redis"`
	}{}
	_, err := toml.DecodeFile(fmt.Sprintf("%s/etc/config.toml", projectRoot), &re)
	if err != nil {
		return nil
	}
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", re.Redis.Ipaddr, re.Redis.Port),
	})
}
