package ServiceUtils

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var (
	globalMetricManageSvc *MetricManageSvc
	globalRedisConf       *RedisConf
	globalMongoConf       *MongoConf
)

type RedisConf struct {
	Ipaddr string `toml:"ipaddr"`
	Port   string `toml:"port"`
}

type MongoConf struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Ipaddr   string `toml:"ipaddr"`
	Port     string `toml:"port"`
}

type MetricManageSvc struct {
	ServiceManageEndpoints []string `toml:"endpoints"`
	MetricManageIpAddr     string   `toml:"ipaddr"`
	MetricManagePort       string   `toml:"port"`
}

func Load(projectRoot string) error {
	MetricManageConfig := &MetricManageSvc{}
	RedisConfig := &RedisConf{}
	MongoConfig := &MongoConf{}
	mc, err := MetricManageConfig.load(projectRoot)
	if err != nil {
		return err
	}
	rc, err := RedisConfig.load(projectRoot)
	if err != nil {
		return err
	}
	moc, err := MongoConfig.load(projectRoot)
	if err != nil {
		return err
	}
	globalMetricManageSvc = mc
	globalRedisConf = rc
	globalMongoConf = moc
	return nil
}

func (e *MetricManageSvc) load(projectRoot string) (*MetricManageSvc, error) {
	conf := struct {
		MetricManageSvc MetricManageSvc `toml:"MetricManageSvc"`
	}{}
	_, err := toml.DecodeFile(fmt.Sprintf("%s/etc/ServiceCenter.toml", projectRoot), &conf)
	if err != nil {
		return nil, err
	}
	return &conf.MetricManageSvc, nil
}

func (r *RedisConf) load(projectRoot string) (*RedisConf, error) {
	conf := struct {
		RedisConf RedisConf `toml:"RedisConf"`
	}{}
	_, err := toml.DecodeFile(fmt.Sprintf("%s/etc/ServiceCenter.toml", projectRoot), &conf)
	if err != nil {
		return nil, err
	}
	return &conf.RedisConf, nil

}

func (m *MongoConf) load(projectRoot string) (*MongoConf, error) {
	conf := struct {
		MongoConf MongoConf `toml:"MongoConf"`
	}{}
	_, err := toml.DecodeFile(fmt.Sprintf("%s/etc/ServiceCenter.toml", projectRoot), &conf)
	if err != nil {
		return nil, err
	}

	return &conf.MongoConf, nil

}
