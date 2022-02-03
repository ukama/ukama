package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/cloud/net/pkg"
	"github.com/ukama/ukamaX/cloud/net/pkg/metrics"
	"github.com/ukama/ukamaX/common/ukama"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"strings"
)

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

func NewNns(config *pkg.Config) *Nns {
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
	if ip, ok = n.cache[nodeId]; !ok {
		if ip, err = n.getFromEtcd(c, nodeId); err != nil {
			metrics.RecordIpRequestFailureMetric()
			return "", err
		}
		n.cache[nodeId] = ip
	}

	metrics.RecordIpRequestSuccessMetric()
	return ip, nil
}

func (n *Nns) getFromEtcd(c context.Context, nodeId string) (string, error) {
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

	_, err = n.etcd.Put(c, nodeId, i.String())
	if err != nil {
		return fmt.Errorf("failed to add record to db. Error: %v", err)
	}

	n.cache[nodeId] = ip
	return nil
}

func (n *Nns) List(ctx context.Context) (map[string]string, error) {
	// list is never use local cache
	vals, err := n.etcd.Get(ctx, "", clientv3.WithPrefix())

	if err != nil {
		return nil, fmt.Errorf("failed to get record from db. Error: %v", err)
	}

	// make sure that the list is unique
	ips := map[string]string{}
	for _, val := range vals.Kvs {
		ips[string(val.Key)] = string(val.Value)
	}

	return ips, nil
}

func (n *Nns) Delete(ctx context.Context, nodeId string) error {
	nd, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		return err
	}

	_, err = n.etcd.Delete(ctx, nd.StringLowercase())
	if err != nil {
		return fmt.Errorf("failed to delete record from db. Error: %v", err)
	}

	delete(n.cache, nodeId)
	return nil
}
