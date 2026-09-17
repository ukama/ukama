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
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
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
	conn, err := grpc.NewClient(networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()),
		ugrpc.TracingDialOption())
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
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Network Service connection: %v", err)
		}
	}
}

func (n *NetworkRegistry) AddNetwork(ctx context.Context, netName string, allowedCountries, allowedNetworks []string,
	budget, overdraft float64, trafficPolicy uint32, paymentLinks bool) (*netpb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
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

func (n *NetworkRegistry) GetNetwork(ctx context.Context, netID string) (*netpb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.Get(ctx, &netpb.GetRequest{NetworkId: netID})
}

func (n *NetworkRegistry) SetNetworkDefault(ctx context.Context, netID string) (*netpb.SetDefaultResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.SetDefault(ctx, &netpb.SetDefaultRequest{NetworkId: netID})
}

func (n *NetworkRegistry) GetDefault(ctx context.Context) (*netpb.GetDefaultResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.GetDefault(ctx, &netpb.GetDefaultRequest{})
}

func (n *NetworkRegistry) GetNetworks(ctx context.Context) (*netpb.GetNetworksResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
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

func (n *NetworkRegistry) RemoveNetwork(ctx context.Context, netID string) (*netpb.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), n.timeout)
	defer cancel()

	return n.client.Delete(ctx, &netpb.DeleteRequest{NetworkId: netID})
}
