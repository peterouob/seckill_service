package server

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/peterouob/seckill_service/utils/logs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdService struct {
	client    *clientv3.Client
	heartbeat int64
}

var (
	serviceHub *EtcdService
	hubOnce    sync.Once
)

func RegisterETCD(etcdServers []string, heartbeat int64) *EtcdService {
	hubOnce.Do(func() {
	start:
		if serviceHub == nil {
			client, err := clientv3.New(clientv3.Config{
				Endpoints:   etcdServers,
				DialTimeout: 5 * time.Second,
			})

			if err != nil {
				time.Sleep(5 * time.Second)
				logs.Log("wait for etcd servers to be ready...")
				goto start
			}

			serviceHub = &EtcdService{client: client, heartbeat: heartbeat}
		} else {
			serviceHub = &EtcdService{
				client:    serviceHub.client,
				heartbeat: heartbeat,
			}
		}
	})

	return serviceHub
}

func (s *EtcdService) Register(service string, endpoint string, leaseID clientv3.LeaseID) clientv3.LeaseID {
	ctx := context.Background()
	if leaseID <= 0 {
		lease, err := s.client.Grant(ctx, s.heartbeat)
		logs.HandelError("grant lease error", err)
		key := fmt.Sprintf("%s/%s/%s",
			strings.TrimRight("/service/grpc", "/"),
			service,
			endpoint)

		_, err = s.client.Put(ctx, key, "", clientv3.WithLease(leaseID))
		logs.HandelError(fmt.Sprintf("puth in %s node %s on etcd error", service, endpoint), err)
		return lease.ID
	}
	keepAlive, err := s.client.KeepAlive(ctx, leaseID)

	if err != nil {
		logs.Error("error to keep etcd alive", err)
	}

	go func() {
		for keepResp := range keepAlive {
			if keepResp == nil {
				logs.Log("lease is unable")
				return
			}
		}
	}()

	return leaseID
}

func (s *EtcdService) UnRegister(service string, endpoint string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s/%s/%s",
		strings.TrimRight("/service/grpc", "/"),
		service,
		endpoint)
	resp, err := s.client.Get(ctx, key)
	if err != nil || len(resp.Kvs) == 0 {
		logs.Log(fmt.Sprintf("Key %s not found in etcd", key))
		return
	}

	leaseID := clientv3.LeaseID(resp.Kvs[0].Lease)
	logs.Log(fmt.Sprintf("Revoking lease %d for key %s", leaseID, key))

	_, err = s.client.Revoke(ctx, leaseID)
	logs.HandelError("revoke lease error", err)

	_, err = s.client.Delete(ctx, key)
	logs.HandelError("delete etcd node error", err)
	logs.Log(fmt.Sprintf("unregistered %s node %s from etcd", service, endpoint))
}
