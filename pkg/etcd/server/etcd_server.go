package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/peterouob/seckill_service/pkg/logger"
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
				logger.Log("wait for etcd servers to be ready...")
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

func (s *EtcdService) Register(service string, endpoint string, leaseID clientv3.LeaseID, expired chan<- struct{}) clientv3.LeaseID {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if leaseID <= 0 {
		lease, err := s.client.Grant(ctx, s.heartbeat)
		if err != nil {
			logger.HandelError("grant lease error", err)
			return 0
		}

		key := fmt.Sprintf("%s/%s/%s",
			strings.TrimRight("/service/grpc", "/"),
			service,
			endpoint,
		)
		_, err = s.client.Put(ctx, key, "", clientv3.WithLease(lease.ID))
		if err != nil {
			logger.HandelError("put lease error", err)
			return 0
		}

		go s.keepAlive(lease.ID, expired)
		return lease.ID
	}

	return leaseID
}

func (s *EtcdService) keepAlive(leaseID clientv3.LeaseID, expired chan<- struct{}) {
	keepAlive, err := s.client.KeepAlive(context.Background(), leaseID)
	if err != nil {
		logger.Error("start keepalive failed: %v", err)
		expired <- struct{}{}
		return
	}

	for keepResp := range keepAlive {
		if keepResp == nil {
			logger.Log("etcd lease expired or revoked")
			expired <- struct{}{}
			return
		}
	}
}

func (s *EtcdService) UnRegister(service string, endpoint string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s/%s/%s",
		strings.TrimRight("/service/grpc", "/"),
		service,
		endpoint)
	resp, err := s.client.Get(ctx, key)
	if err != nil || len(resp.Kvs) == 0 {
		return errors.New("etcd service unregister fail")
	}

	leaseID := clientv3.LeaseID(resp.Kvs[0].Lease)
	logger.Log(fmt.Sprintf("Revoking lease %d for key %s", leaseID, key))

	if _, err := s.client.Revoke(ctx, leaseID); err != nil {
		logger.HandelError("revoke lease error", err)
		return err
	}

	if _, err = s.client.Delete(ctx, key); err != nil {
		logger.HandelError("delete etcd node error", err)
		return err
	}
	logger.Log(fmt.Sprintf("unregistered %s node %s from etcd", service, endpoint))
	return nil
}

func (s *EtcdService) Client() *clientv3.Client {
	return s.client
}
