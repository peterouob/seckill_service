package config

import (
	"os"
)

type Config struct {
	DbDSN      string
	RDbDSN     string
	GrpcAddr   string
	HttpAddr   string
	Table      []any
	EtcdConfig *EtcdConfig
}

type EtcdConfig struct {
	ServiceName string
	Endpoints   []string
}

func ProvideConfig(grpcAddr, httpAddr string, etcdServiceName string) *Config {
	var dbDSN string
	var rDbDSN string
	var etcdEndpoint string
	if dbDSN = os.Getenv("DB_DSN"); dbDSN == "" {
		dbDSN = "root:123456@tcp(localhost:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local"
	}
	rDbDSN = os.Getenv("RDB_DSN")
	if etcdEndpoint = os.Getenv("ETCD_ENDPOINTS"); etcdEndpoint == "" {
		etcdEndpoint = "127.0.0.1:2379"
	}
	return &Config{
		DbDSN:    dbDSN,
		RDbDSN:   rDbDSN,
		GrpcAddr: grpcAddr,
		HttpAddr: httpAddr,
		EtcdConfig: &EtcdConfig{
			Endpoints:   []string{etcdEndpoint},
			ServiceName: etcdServiceName,
		},
	}
}
