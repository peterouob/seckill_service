package etcdregister

import (
	"fmt"
	"time"

	"github.com/peterouob/seckill_service/utils/etcd/server"
	"github.com/peterouob/seckill_service/utils/logs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdRegister struct {
	Client  *server.EtcdService
	leaseId clientv3.LeaseID
	heart   int64
	expire  chan struct{}
	stopCh  chan struct{}
}

func NewEtcdRegister(endpoints []string, heart int64) *EtcdRegister {
	c := server.RegisterETCD(endpoints, heart)
	e := &EtcdRegister{
		Client: c,
		heart:  heart,
		expire: make(chan struct{}, 1),
		stopCh: make(chan struct{}),
	}
	return e
}

func (e *EtcdRegister) Register(serviceName, value string) error {
	leaseID := e.Client.Register(serviceName, value, 0, e.expire)
	if leaseID == 0 {
		return fmt.Errorf("register service %s failed", serviceName)
	}
	e.leaseId = leaseID

	go e.watchExpire(serviceName, value)
	return nil
}

func (e *EtcdRegister) watchExpire(serviceName, value string) {
	ticker := time.NewTicker(time.Duration(e.heart/3) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-e.stopCh:
			logs.Log("etcd register loop stopped")
			return

		case <-e.expire:
			logs.Log("lease expired, re-registering immediately...")
			newID := e.Client.Register(serviceName, value, 0, e.expire)
			if newID == 0 {
				logs.ErrorMsgF("re-register failed for %s", serviceName)
				continue
			}
			e.leaseId = newID
			logs.Logf("re-registered %s with new leaseID %d", serviceName, newID)
		}
	}
}

func (e *EtcdRegister) UnRegister(serviceName, addr string) error {
	close(e.stopCh)
	err := e.Client.UnRegister(serviceName, addr)
	logs.Logf("unregister service: %s from etcd, addr: %s", serviceName, addr)
	return err
}
