package pkg

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const orgNetMappingKeyPrefix = "map:"

type NodeOrgMap struct {
	etcd *clientv3.Client
}
type OrgNet struct {
	Org      string
	Network  string
	Site     string
	NodePort int32
	NodeIp   string
	MeshPort int32
}

func NewNodeToOrgMap(config *Config) *NodeOrgMap {
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: config.DialTimeoutSecond,
		Endpoints:   []string{config.EtcdHost},
	})
	if err != nil {
		logrus.Fatalf("Cannot connect to etcd: %v", err)
	}

	return &NodeOrgMap{
		etcd: client,
	}
}

func (n *NodeOrgMap) Add(ctx context.Context, nodeId, org, network, site, nodeIp string, nodePort, meshPort int32) error {
	nodeIdKey := formatMappKey(nodeId)
	_, err := n.etcd.Put(ctx, nodeIdKey, org+"."+network+"."+site+"."+b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", nodeIp, nodePort))))
	if err != nil {
		return fmt.Errorf("failed to add record to db. Error: %v", err)
	}
	return nil
}

func (n *NodeOrgMap) List(ctx context.Context) (map[string]OrgNet, error) {
	vals, err := n.etcd.Get(ctx, orgNetMappingKeyPrefix, clientv3.WithPrefix())

	if err != nil {
		return nil, fmt.Errorf("failed to get record from db. Error: %v", err)
	}

	res := map[string]OrgNet{}
	var p int64
	for _, val := range vals.Kvs {
		c := strings.Split(string(val.Value), ".")
		if len(c) != 4 {
			logrus.Errorf("failed to parse org.net.site.ip:port structure for '%s' with value '%s'", string(val.Key), string(val.Value))
		}

		b64Add, err := b64.StdEncoding.DecodeString(c[3])
		add := strings.Split(string(b64Add), ":")
		if len(add) != 2 {
			logrus.Errorf("failed to parse ip:port structure for '%s'", add)
			return nil, err
		} else {
			p, err = strconv.ParseInt(add[1], 10, 32)
			if err != nil {
				logrus.Errorf("failed to parse covert port '%s' to int32", add[1])
				return nil, err
			}
		}

		res[strings.TrimPrefix(string(val.Key), orgNetMappingKeyPrefix)] = OrgNet{
			Org:      c[0],
			Network:  c[1],
			Site:     c[2],
			NodeIp:   add[0],
			NodePort: int32(p),
		}
	}

	return res, nil
}

func formatMappKey(nodeId string) string {
	return orgNetMappingKeyPrefix + nodeId
}
