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

	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/site/pb/gen"
	"google.golang.org/grpc"
)

type SiteRegistry struct {
	conn    *grpc.ClientConn
	client  pb.SiteServiceClient
	timeout time.Duration
	host    string
}

func NewSiteRegistry(siteHost string, timeout time.Duration) *SiteRegistry {

	conn, err := grpc.NewClient(siteHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
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

func (r *SiteRegistry) Close() {
	r.conn.Close()
}

func (r *SiteRegistry) GetSite(siteId string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.Get(ctx, &pb.GetRequest{SiteId: siteId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *SiteRegistry) GetSites(networkId string) (*pb.GetSitesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetSites(ctx, &pb.GetSitesRequest{NetworkId: networkId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *SiteRegistry) AddSite(networkId, name, backhaulId, powerId, accessId, switchId string, isDeactivated bool, latitude, longitude float64, installDate string) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	res, err := r.client.Add(ctx, &pb.AddRequest{
		Name:          name,
		NetworkId:     networkId,
		BackhaulId:    backhaulId,
		PowerId:       powerId,
		AccessId:      accessId,
		SwitchId:      switchId,
		IsDeactivated: isDeactivated,
		Latitude:      latitude,
		Longitude:     longitude,
		InstallDate:   installDate,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *SiteRegistry) UpdateSite(siteId, name, backhaulId, powerId, accessId, switchId string, isDeactivated bool, latitude, longitude float64, installDate string) (*pb.UpdateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	req := &pb.UpdateRequest{
		SiteId:        siteId,
		Name:          name,
		BackhaulId:    backhaulId,
		PowerId:       powerId,
		AccessId:      accessId,
		SwitchId:      switchId,
		IsDeactivated: isDeactivated,
		Latitude:      latitude,
		Longitude:     longitude,
		InstallDate:   installDate,
	}

	res, err := r.client.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
