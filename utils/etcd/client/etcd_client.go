package client

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/peterouob/seckill_service/utils/logs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceHub struct {
	client        *clientv3.Client
	endPointCache sync.Map
	watch         sync.Map
}

var (
	serviceHub *ServiceHub
	hubOnce    sync.Once
)

func GetService(etcdAddr []string) *ServiceHub {
	hubOnce.Do(func() {
		if serviceHub == nil {
			client, err := clientv3.New(clientv3.Config{
				Endpoints:   etcdAddr,
				DialTimeout: 5 * time.Second,
			})
			logs.HandelError("new etcd client error", err)
			serviceHub = &ServiceHub{
				client:        client,
				endPointCache: sync.Map{},
				watch:         sync.Map{},
			}
		}
	})

	return serviceHub
}

func (s *ServiceHub) getServiceEndpoint(service string) []string {
	ctx := context.Background()
	prefix := fmt.Sprintf("%s/%s/",
		strings.TrimRight("/service/grpc", "/"),
		service)

	resp, err := s.client.Get(ctx, prefix, clientv3.WithPrefix())
	logs.HandelError("get etcd service endpoint error", err)
	endpoints := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		path := strings.Split(string(kv.Key), "/")
		endpoints = append(endpoints, path[len(path)-1])
	}
	logs.Log(fmt.Sprintf("get %s etcd service endpoints: %v\n", service, endpoints))
	return endpoints
}

func (s *ServiceHub) watchEndpoint(service string) {
	if _, exists := s.watch.LoadOrStore(service, true); exists {
		return
	}
	ctx := context.Background()
	prefix := fmt.Sprintf("%s/%s/",
		strings.TrimRight("/service/grpc", "/"),
		service)
	ch := s.client.Watch(ctx, prefix, clientv3.WithPrefix())
	logs.Log(fmt.Sprintf("watching %s node change ...", service))
	go func() {
		for resp := range ch {
			for _, event := range resp.Events {
				path := strings.Split(string(event.Kv.Key), "/")
				if len(path) > 2 {
					service := path[len(path)-2]
					endpoints := s.getServiceEndpoint(service)
					if len(endpoints) > 0 {
						s.endPointCache.Store(service, endpoints)
					} else {
						s.endPointCache.Delete(service)
					}
				}
			}
		}
	}()
}

func (s *ServiceHub) GetServiceEndPoint(service string) []string {
	s.watchEndpoint(service)
	if endpoints, exists := s.endPointCache.Load(service); exists {
		return endpoints.([]string)
	} else {
		endpoints := s.getServiceEndpoint(service)
		if len(endpoints) > 0 {
			s.endPointCache.Store(service, endpoints)
		}
		return endpoints
	}
}
