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
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/messaging/nns/pkg/metrics"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var separator = "|"

type Nns struct {
	etcd    *clientv3.Client
	orgName string
}

type NnsReader interface {
	Get(ctx context.Context, nodeId string) (*OrgMap, error)
	GetAll(ctx context.Context) ([]OrgMap, error)
}

func NewNns(config *Config) *Nns {
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: config.DialTimeoutSecond,
		Endpoints:   []string{config.EtcdHost},
	})
	if err != nil {
		log.Fatalf("Cannot connect to etcd: %v", err)
	}

	return &Nns{
		etcd:    client,
		orgName: config.OrgName,
	}
}

type OrgMesh struct {
	OrgName  string
	MeshIp   string
	MeshPort int32
}

func (o *OrgMesh) string() string {
	return fmt.Sprintf("%s%s%d", o.MeshIp, separator, o.MeshPort)
}

func (o *OrgMesh) parse(value string) error {
	parts := strings.Split(value, separator)
	if len(parts) != 2 {
		return fmt.Errorf("invalid org mesh string: %s", value)
	}
	o.MeshIp = parts[0]
	port, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse mesh port: %v", err)
	}
	o.MeshPort = int32(port)
	return nil
}

func (o *OrgMesh) constructKeyAndValue(obj OrgMesh) (string, string) {
	return obj.OrgName + separator + "mesh", obj.string()
}

func (n *Nns) SetMesh(ctx context.Context, ip string, port int32) error {
	obj := OrgMesh{
		MeshIp:   ip,
		MeshPort: port,
		OrgName:  n.orgName,
	}
	key, value := obj.constructKeyAndValue(obj)
	_, err := n.etcd.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to set mesh IP and port. Error: %v", err)
	}
	metrics.RecordSetIpMetric()
	return nil
}

func (n *Nns) GetMesh(ctx context.Context) (*OrgMesh, error) {
	mesh, err := n.etcd.Get(ctx, n.orgName+separator+"mesh")
	if err != nil {
		metrics.RecordIpRequestFailureMetric()
		return nil, fmt.Errorf("failed to get mesh. Error: %v", err)
	}
	if len(mesh.Kvs) == 0 {
		metrics.RecordIpRequestFailureMetric()
		return nil, fmt.Errorf("mesh not found for org: %s", n.orgName)
	}
	orgMesh := OrgMesh{}
	err = orgMesh.parse(string(mesh.Kvs[0].Value))
	if err != nil {
		metrics.RecordIpRequestFailureMetric()
		return nil, fmt.Errorf("failed to parse mesh. Error: %v", err)
	}
	metrics.RecordIpRequestSuccessMetric()
	return &orgMesh, nil
}

type OrgMap struct {
	Org          string
	Network      string
	Site         string
	MeshIp       string
	MeshHostName string
	MeshPort     int32
	NodeId       string
	NodeIp       string
	NodePort     int32
}

func (o *OrgMap) string() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%d%s%s%s%s%s%d", o.Org, separator, o.Network, separator, o.Site, separator, o.MeshIp, separator, o.MeshHostName, separator, o.MeshPort, separator, o.NodeId, separator, o.NodeIp, separator, o.NodePort)
}

func (o *OrgMap) parse(value string) error {
	parts := strings.Split(value, separator)
	if len(parts) != 9 {
		return fmt.Errorf("invalid org net string: %s", value)
	}

	o.Org = parts[0]
	o.Network = parts[1]
	o.Site = parts[2]
	o.MeshIp = parts[3]
	o.MeshHostName = parts[4]
	meshPort, err := strconv.ParseInt(parts[5], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse mesh port: %v", err)
	}
	o.MeshPort = int32(meshPort)
	o.NodeId = parts[6]
	o.NodeIp = parts[7]
	nodePort, err := strconv.ParseInt(parts[8], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse node port: %v", err)
	}
	o.NodePort = int32(nodePort)

	return nil
}

func (o *OrgMap) constructKeyAndValue(obj OrgMap) (string, string) {
	return obj.NodeId, obj.string()
}

func (n *Nns) Add(ctx context.Context, obj OrgMap) error {
	mesh, err := n.GetMesh(ctx)
	if mesh == nil {
		return fmt.Errorf("failed to get mesh. Error: %v", err)
	}
	obj.MeshIp = mesh.MeshIp
	obj.MeshPort = mesh.MeshPort

	key, value := obj.constructKeyAndValue(obj)
	_, err = n.etcd.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to add record %s with value %s. Error: %v", key, value, err)
	}
	log.Infof("Added node %s with value %s to etcd", key, value)
	return nil
}

func (n *Nns) Get(ctx context.Context, key string) (*OrgMap, error) {
	val, err := n.etcd.Get(ctx, key)
	if err != nil {
		metrics.RecordIpRequestFailureMetric()
		return nil, fmt.Errorf("failed to get record from db. Error: %v", err)
	}
	orgNet := OrgMap{}
	err = orgNet.parse(string(val.Kvs[0].Value))
	if err != nil {
		metrics.RecordIpRequestFailureMetric()
		return nil, fmt.Errorf("failed to parse stored node org map  for %s. Error: %v", val.Kvs[0].Key, err)
	}
	log.Infof("Got node %s from etcd", key)
	metrics.RecordIpRequestSuccessMetric()
	return &orgNet, nil
}

func (n *Nns) GetAll(ctx context.Context) ([]OrgMap, error) {
	vals, err := n.etcd.Get(ctx, "", clientv3.WithPrefix())
	if err != nil {
		metrics.RecordIpRequestFailureMetric()
		return nil, fmt.Errorf("failed to get record from db. Error: %v", err)
	}

	obj := make([]OrgMap, 0)
	for _, val := range vals.Kvs {
		orgMap := OrgMap{}
		err = orgMap.parse(string(val.Value))
		if err != nil {
			metrics.RecordIpRequestFailureMetric()
			return nil, fmt.Errorf("failed to parse stored node org map  for %s. Error: %v", val.Key, err)
		}
		obj = append(obj, orgMap)
	}
	log.Infof("Got %d nodes from etcd", len(obj))
	metrics.RecordIpRequestSuccessMetric()
	return obj, nil
}

func (n *Nns) DeleteAll(ctx context.Context) error {
	_, err := n.etcd.Delete(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to delete all nodes from etcd. Error: %v", err)
	}
	log.Infof("Deleted all nodes from etcd")
	return nil
}

func (n *Nns) Delete(ctx context.Context, key string) error {
	_, err := n.etcd.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete node %s from etcd. Error: %v", key, err)
	}
	log.Infof("Deleted node %s from etcd", key)
	return nil
}

func (n *Nns) UpdateNodeMesh(ctx context.Context, ip string, port int32) error {
	items, err := n.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get node org map records. Error: %v", err)
	}

	for _, item := range items {
		obj := OrgMap{
			MeshIp:       ip,
			MeshPort:     port,
			Org:          item.Org,
			Site:         item.Site,
			NodeId:       item.NodeId,
			NodeIp:       item.NodeIp,
			Network:      item.Network,
			NodePort:     item.NodePort,
			MeshHostName: item.MeshHostName,
		}
		_, err = n.etcd.Put(ctx, obj.NodeId, obj.string())
		if err != nil {
			return fmt.Errorf("failed to update mesh IP and port for %s. Error: %v", obj.NodeId, err)
		}
	}

	log.Infof("Updated mesh IP and port for %d nodes", len(items))

	return nil
}

func (n *Nns) UpdateNode(ctx context.Context, nodeId string, nodeIp string, nodePort int32) error {
	item, err := n.Get(ctx, nodeId)
	if err != nil {
		return fmt.Errorf("failed to get node record. Error: %v", err)
	}

	obj := OrgMap{
		NodeId:       item.NodeId,
		NodeIp:       nodeIp,
		NodePort:     nodePort,
		MeshIp:       item.MeshIp,
		MeshHostName: item.MeshHostName,
		MeshPort:     item.MeshPort,
		Org:          item.Org,
		Network:      item.Network,
		Site:         item.Site,
	}
	_, err = n.etcd.Put(ctx, obj.NodeId, obj.string())
	if err != nil {
		return fmt.Errorf("failed to update node IP and port for %s. Error: %v", obj.NodeId, err)
	}
	log.Infof("Updated node IP and port for %s", obj.NodeId)
	return nil
}
