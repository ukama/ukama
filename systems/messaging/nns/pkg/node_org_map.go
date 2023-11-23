/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const orgNetMappingKeyPrefix = "map:"

const (
	E_NODE_IP_IDX   = 0
	E_NODE_PORT_IDX = 1
	E_MESH_HOST_IDX = 2
	E_MESH_PORT_IDX = 3
	E_MAX_IDX       = 4
)

/* org.net.site.b64Encoded(nodeIp:Nodeport:meshPort) */
const (
	MAP_ORG_IDX  = 0
	MAP_NW_IDX   = 1
	MAP_SITE_IDX = 2
	MAP_ENC_IDX  = 3
	MAP_MAX_IDX  = 4
)

type NodeOrgMap struct {
	etcd *clientv3.Client
}

type OrgNet struct {
	Org          string
	Network      string
	Site         string
	NodePort     int32
	NodeIp       string
	MeshHostName string
	MeshPort     int32
}

func NewNodeToOrgMap(config *Config) *NodeOrgMap {
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: config.DialTimeoutSecond,
		Endpoints:   []string{config.EtcdHost},
	})
	if err != nil {
		log.Fatalf("Cannot connect to etcd: %v", err)
	}

	return &NodeOrgMap{
		etcd: client,
	}
}

func (n *NodeOrgMap) UpdateMesh(ctx context.Context, hostname string, ip string) error {

	_, err := n.etcd.Put(ctx, hostname, ip)
	if err != nil {
		return fmt.Errorf("failed to add mesh IP to db. Error: %v", err)
	}
	return nil
}

func (n *NodeOrgMap) GetMesh(ctx context.Context, hostname string) (*string, error) {
	var meshIp string
	val, err := n.etcd.Get(ctx, hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to get mesh IP from db. Error: %v", err)
	}

	for _, val := range val.Kvs {
		meshIp = string(val.Value)
		ip := net.ParseIP(meshIp)
		if ip == nil {
			log.Errorf("failed to get mesh IP from db string %s", meshIp)
			return nil, fmt.Errorf("failed to get mesh IP from db string %s", meshIp)
		}
	}

	return &meshIp, nil
}

func (n *NodeOrgMap) Add(ctx context.Context, nodeId, org, network, site, nodeIp, meshHostName string, nodePort, meshPort int32) error {
	nodeIdKey := formatMapKey(nodeId)
	_, err := n.etcd.Put(ctx, nodeIdKey, org+"."+network+"."+site+"."+b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d:%s:%d", nodeIp, nodePort, meshHostName, meshPort))))
	if err != nil {
		return fmt.Errorf("failed to add record to db. Error: %v", err)
	}
	return nil
}

func (n *NodeOrgMap) Get(ctx context.Context, nodeId string) (OrgNet, error) {
	nodeIdKey := formatMapKey(nodeId)
	val, err := n.etcd.Get(ctx, nodeIdKey)
	if err != nil {
		return OrgNet{}, fmt.Errorf("failed to get record from db. Error: %v", err)
	}

	res := map[string]OrgNet{}

	for _, val := range val.Kvs {
		orgNet, err := parseMapValue(val.Value)
		if err != nil {
			return OrgNet{}, fmt.Errorf("failed to parse stored node org map  for %s. Error: %v", val.Key, err)
		}

		res[strings.TrimPrefix(string(val.Key), orgNetMappingKeyPrefix)] = *orgNet
	}

	return res[nodeId], nil
}

func (n *NodeOrgMap) List(ctx context.Context) (map[string]OrgNet, error) {

	vals, err := n.etcd.Get(ctx, orgNetMappingKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get record from db. Error: %v", err)
	}

	res := map[string]OrgNet{}
	for _, val := range vals.Kvs {
		orgNet, err := parseMapValue(val.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse stored node org map  for %s. Error: %v", val.Key, err)
		}

		res[strings.TrimPrefix(string(val.Key), orgNetMappingKeyPrefix)] = *orgNet
	}

	return res, nil
}

func formatMapKey(nodeId string) string {
	return orgNetMappingKeyPrefix + nodeId
}

func parseMapValue(data []byte) (*OrgNet, error) {
	var p1, p2 int64
	c := strings.Split(string(data), ".")
	if len(c) != MAP_MAX_IDX {
		log.Errorf("failed to parse org.net.site.ip:port structure for value '%s'", string(data))
	}

	b64Add, err := b64.StdEncoding.DecodeString(c[MAP_ENC_IDX])
	add := strings.Split(string(b64Add), ":")
	if len(add) != E_MAX_IDX {
		log.Errorf("failed to parse ip:port:meshHostName:meshport structure for '%s'", add)
		return nil, err
	} else {
		p1, err = strconv.ParseInt(add[E_NODE_PORT_IDX], 10, 32)
		if err != nil {
			log.Errorf("failed to convert port '%s' to int32", add[E_NODE_IP_IDX])
			return nil, err
		}

		p2, err = strconv.ParseInt(add[E_MESH_PORT_IDX], 10, 32)
		if err != nil {
			log.Errorf("failed to parse covert port '%s' to int32", add[E_MESH_PORT_IDX])
			return nil, err
		}
	}

	return &OrgNet{
		Org:          c[0],
		Network:      c[1],
		Site:         c[2],
		NodeIp:       add[E_NODE_IP_IDX],
		NodePort:     int32(p1),
		MeshHostName: add[E_MESH_HOST_IDX],
		MeshPort:     int32(p2),
	}, nil

}
