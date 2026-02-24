package etcdregister

import (
	"fmt"

	"github.com/peterouob/seckill_service/utils/etcd/server"
	"github.com/peterouob/seckill_service/utils/logs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdRegister struct {
	client  *server.EtcdService
	leaseId clientv3.LeaseID
	heart   int64
}

func NewEtcdRegister(endpoints []string, heart int64) *EtcdRegister {
	c := server.RegisterETCD(endpoints, heart)
	e := &EtcdRegister{
		client: c,
		heart:  heart,
	}
	return e
}

func (e *EtcdRegister) Register(serviceName, value string) {
	e.leaseId = e.client.Register(serviceName, value, 0)
	//tools.Log(fmt.Sprintf("Registered service %s at %s", serviceName, addr))
	go func() {
		for {
			e.client.Register(serviceName, value, e.leaseId)
		}
	}()
}

func (e *EtcdRegister) UnRegister(serviceName, addr string) {
	e.client.UnRegister(serviceName, addr)
	logs.Log(fmt.Sprintf("unregiter service: %s from etcd, addr: %s", serviceName, addr))
}
