package pkg

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/cloud/net/pkg/metrics"
	"github.com/ukama/ukamaX/common/ukama"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const nodeIdKeyPrefix = "nodeId:"

// NNS service maintainse NodeID -> IP pamming and provides methods to work with it
type Nns struct {
	etcd  *clientv3.Client
	cache map[string]string
}

type NnsReader interface {
	Get(c context.Context, nodeId string) (string, error)
	List(ctx context.Context) (map[string]string, error)
}

type NnsWriter interface {
	Set(c context.Context, nodeId string, ip string) error
	Delete(ctx context.Context, nodeId string) error
}

func NewNns(config *Config) *Nns {
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: config.DialTimeoutSecond,
		Endpoints:   []string{config.EtcdHost},
	})
	if err != nil {
		logrus.Fatalf("Cannot connect to etcd: %v", err)
	}

	return &Nns{
		etcd:  client,
		cache: make(map[string]string),
	}
}

func (n *Nns) Get(c context.Context, nodeId string) (ip string, err error) {
	nodeId = strings.ToLower(nodeId)

	if _, err = ukama.ValidateNodeId(nodeId); err != nil {
		metrics.RecordIpRequestFailureMetric()
		return "", status.Error(codes.InvalidArgument, err.Error())
	}
	var ok bool

	nodeIdKey := formatNodeIdKey(nodeId)
	if ip, ok = n.cache[nodeIdKey]; !ok {
		if ip, err = n.getFromEtcd(c, nodeIdKey); err != nil {
			metrics.RecordIpRequestFailureMetric()
			return "", err
		}
		n.cache[nodeIdKey] = ip
	}

	metrics.RecordIpRequestSuccessMetric()
	return ip, nil
}

func (n *Nns) getFromEtcd(c context.Context, nodeId string) (string, error) {
	logrus.Infof("Getting ip from etcd for nodeId: %s", nodeId)
	val, err := n.etcd.Get(c, nodeId)
	if err != nil {
		return "", fmt.Errorf("failed to get record from db. Error: %v", err)
	}

	if val.Count == 0 {
		return "", status.Error(codes.NotFound, fmt.Sprintf("record %s not found", nodeId))
	}
	if val.Count > 1 {
		return "", status.Error(codes.Internal, fmt.Sprintf("more than one record %s found", nodeId))
	}

	return string(val.Kvs[0].Value), nil
}

func (n *Nns) Set(c context.Context, nodeId string, ip string) (err error) {
	nodeId = strings.ToLower(nodeId)
	metrics.RecordSetIpMetric()
	if _, err = ukama.ValidateNodeId(nodeId); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	i := net.ParseIP(ip)
	if i == nil {
		return fmt.Errorf("not valid ip")
	}
	nodeIdKey := formatNodeIdKey(nodeId)
	_, err = n.etcd.Put(c, nodeIdKey, i.String())
	if err != nil {
		return fmt.Errorf("failed to add record to db. Error: %v", err)
	}

	n.cache[nodeIdKey] = ip
	return nil
}

func (n *Nns) List(ctx context.Context) (map[string]string, error) {
	// list never uses local cache
	vals, err := n.etcd.Get(ctx, nodeIdKeyPrefix, clientv3.WithPrefix())

	if err != nil {
		return nil, fmt.Errorf("failed to get record from db. Error: %v", err)
	}

	ips := map[string]string{}
	for _, val := range vals.Kvs {
		ips[strings.TrimPrefix(string(val.Key), nodeIdKeyPrefix)] = string(val.Value)
	}

	return ips, nil
}

func (n *Nns) Delete(ctx context.Context, nodeId string) error {
	nd, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		return err
	}

	delete(n.cache, formatNodeIdKey(nodeId))

	_, err = n.etcd.Delete(ctx, formatNodeIdKey(nd.StringLowercase()))
	if err != nil {
		return fmt.Errorf("failed to delete record from db. Error: %v", err)
	}

	return nil
}

func formatNodeIdKey(nodeId string) string {
	return nodeIdKeyPrefix + nodeId
}
