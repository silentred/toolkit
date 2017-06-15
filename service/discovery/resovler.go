package discovery

import (
	"errors"
	"fmt"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
	"google.golang.org/grpc/naming"
)

// EtcdResolver is the implementaion of grpc.naming.Resolver
type EtcdResolver struct {
	prefix      string
	ServiceName string // service name to resolve
}

// NewResolver return EtcdResolver with service name
func NewResolver(serviceName, prefix string) *EtcdResolver {
	if len(prefix) == 0 {
		prefix = DefaultPrefix
	}
	return &EtcdResolver{ServiceName: serviceName, prefix: prefix}
}

// Resolve to resolve the service from etcd, target is the dial address of etcd
// target example: "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379"
func (er *EtcdResolver) Resolve(target string) (naming.Watcher, error) {
	if er.ServiceName == "" {
		return nil, errors.New("wonaming: no service name provided")
	}

	// generate etcd client, return if error
	endpoints := strings.Split(target, ",")
	conf := etcd.Config{
		Endpoints:               endpoints,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	client, err := etcd.New(conf)
	if err != nil {
		return nil, fmt.Errorf("wonaming: creat etcd error: %s", err.Error())
	}

	kapi := etcd.NewKeysAPI(client)

	// Return EtcdWatcher
	watcher := &EtcdWatcher{
		prefix:  er.prefix,
		srvName: er.ServiceName,
		client:  &client,
		kapi:    kapi,
		addrs:   nil,
	}
	return watcher, nil
}
