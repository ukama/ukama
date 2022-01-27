package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/cloud/net/pkg"
	"github.com/ukama/ukamaX/common/ukama"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type Nns struct {
	etcd *clientv3.Client
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
		etcd: client,
	}
}

func (n *Nns) Get(c context.Context, nodeId string) (string, error) {
	nd, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		return "", err
	}

	val, err := n.etcd.Get(c, nd.StringLowercase())
	if err != nil {
		return "", fmt.Errorf("failed to get record from db. Error: %v", err)
	}

	if val.Count == 0 {
		return "", status.Error(codes.NotFound, fmt.Sprintf("record %s not found", nd))
	}
	if val.Count > 1 {
		return "", status.Error(codes.Internal, fmt.Sprintf("more than one record %s found", nd))
	}

	return string(val.Kvs[0].Value), nil
}

func (n *Nns) Set(c context.Context, nodeId string, ip string) error {
	nd, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		return err
	}

	i := net.ParseIP(ip)
	if i == nil {
		return fmt.Errorf("not valid ip")
	}

	_, err = n.etcd.Put(c, nd.StringLowercase(), i.String())
	if err != nil {
		return fmt.Errorf("failed to add record to db. Error: %v", err)
	}

	return nil
}

func (n *Nns) List(ctx context.Context) (map[string]string, error) {
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

	return nil
}
