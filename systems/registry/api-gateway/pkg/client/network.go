/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
)

const DefaultNetworkName = "default"

type NetworkRegistry struct {
	conn    *grpc.ClientConn
	client  netpb.NetworkServiceClient
	timeout time.Duration
	host    string
}

func NewNetworkRegistry(networkHost string, timeout time.Duration) *NetworkRegistry {
	conn, err := grpc.NewClient(networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to registry's network service: %v", err)
	}
	client := netpb.NewNetworkServiceClient(conn)

	return &NetworkRegistry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    networkHost,
	}
}

func NewNetworkRegistryFromClient(networkClient netpb.NetworkServiceClient) *NetworkRegistry {
	return &NetworkRegistry{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  networkClient,
	}
}

func (n *NetworkRegistry) Close() {
	err := n.conn.Close()
	if err != nil {
		log.Warnf("Failed to gracefully close Network Service connection: %v", err)
	}
}

func (n *NetworkRegistry) AddNetwork(netName string, allowedCountries, allowedNetworks []string,
	budget, overdraft float64, trafficPolicy uint32, paymentLinks bool) (*netpb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Add(ctx, &netpb.AddRequest{
		Name:             netName,
		AllowedCountries: allowedCountries,
		AllowedNetworks:  allowedNetworks,
		Budget:           budget,
		Overdraft:        overdraft,
		TrafficPolicy:    trafficPolicy,
		PaymentLinks:     paymentLinks,
	})
}

func (n *NetworkRegistry) GetNetwork(netID string) (*netpb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Get(ctx, &netpb.GetRequest{NetworkId: netID})
}

func (n *NetworkRegistry) SetNetworkDefault(netID string) (*netpb.SetDefaultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.SetDefault(ctx, &netpb.SetDefaultRequest{NetworkId: netID})
}

func (n *NetworkRegistry) GetDefault() (*netpb.GetDefaultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.GetDefault(ctx, &netpb.GetDefaultRequest{})
}

func (n *NetworkRegistry) GetNetworks() (*netpb.GetNetworksResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	res, err := n.client.GetAll(ctx, &netpb.GetNetworksRequest{})
	if err != nil {
		return nil, err
	}

	if res.Networks == nil {
		return &netpb.GetNetworksResponse{Networks: []*netpb.Network{}}, nil
	}

	return res, nil
}
