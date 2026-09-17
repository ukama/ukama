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
	pb "github.com/ukama/ukama/systems/registry/site/pb/gen"
)

type SiteRegistry struct {
	conn    *grpc.ClientConn
	client  pb.SiteServiceClient
	timeout time.Duration
	host    string
}

func NewSiteRegistry(siteHost string, timeout time.Duration) *SiteRegistry {
	conn, err := grpc.NewClient(siteHost, grpc.WithTransportCredentials(insecure.NewCredentials()),
		ugrpc.TracingDialOption())
	if err != nil {
		log.Fatalf("Failed to connect to Site Service: %v", err)
	}
	client := pb.NewSiteServiceClient(conn)

	return &SiteRegistry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    siteHost,
	}
}

func NewSiteRegistryFromClient(mClient pb.SiteServiceClient) *SiteRegistry {
	return &SiteRegistry{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (s *SiteRegistry) Close() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Site Service connection: %v", err)
		}
	}
}

func (i *SiteRegistry) GetSite(ctx context.Context, siteId string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), i.timeout)
	defer cancel()

	return i.client.Get(ctx, &pb.GetRequest{SiteId: siteId})
}

func (i *SiteRegistry) List(ctx context.Context, networkId string, isDeactivate bool) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), i.timeout)
	defer cancel()

	return i.client.List(ctx, &pb.ListRequest{NetworkId: networkId, IsDeactivated: isDeactivate})
}

func (i *SiteRegistry) AddSite(ctx context.Context, networkId, name, backhaulId, powerId, accessId, switchId, location, spectrumId string,
	isDeactivated bool, latitude, longitude string, installDate string) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), i.timeout)
	defer cancel()

	return i.client.Add(ctx, &pb.AddRequest{
		Name:          name,
		NetworkId:     networkId,
		Location:      location,
		BackhaulId:    backhaulId,
		PowerId:       powerId,
		AccessId:      accessId,
		SwitchId:      switchId,
		SpectrumId:    spectrumId,
		IsDeactivated: isDeactivated,
		Latitude:      latitude,
		Longitude:     longitude,
		InstallDate:   installDate,
	})
}

// Site creation waits for the group result, or for the HTTP caller to leave.
// Its durable coordinator continues even when this context is cancelled.
func (i *SiteRegistry) AddSiteContext(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	return i.client.Add(ctx, req)
}

func (i *SiteRegistry) UpdateSite(ctx context.Context, siteId, name string) (*pb.UpdateResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), i.timeout)
	defer cancel()

	return i.client.Update(ctx, &pb.UpdateRequest{
		SiteId: siteId,
		Name:   name,
	})
}

func (i *SiteRegistry) RemoveSite(ctx context.Context, siteId string) (*pb.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), i.timeout)
	defer cancel()

	return i.client.Delete(ctx, &pb.DeleteRequest{SiteId: siteId})
}
