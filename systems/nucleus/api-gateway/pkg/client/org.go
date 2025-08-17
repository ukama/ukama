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
	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
)

const DefaultNetworkName = "default"

type OrgRegistry struct {
	conn      *grpc.ClientConn
	orgClient orgpb.OrgServiceClient
	timeout   time.Duration
}

func NewOrgRegistry(orgHost string, timeout time.Duration) *OrgRegistry {
	orgConn, err := grpc.NewClient(orgHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Org Service: %v", err)
	}
	orgClient := orgpb.NewOrgServiceClient(orgConn)

	return &OrgRegistry{
		conn:      orgConn,
		orgClient: orgClient,
		timeout:   timeout,
	}
}

func NewOrgRegistryFromClient(orgClient orgpb.OrgServiceClient) *OrgRegistry {
	return &OrgRegistry{
		timeout:   1 * time.Second,
		conn:      nil,
		orgClient: orgClient,
	}
}

func (o *OrgRegistry) Close() {
	if o.conn != nil {
		if err := o.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Org Service connection: %v", err)
		}
	}
}

func (o *OrgRegistry) GetOrg(orgName string) (*orgpb.GetByNameResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.orgClient.GetByName(ctx, &orgpb.GetByNameRequest{Name: orgName})
}

func (o *OrgRegistry) GetOrgs(ownerUUID string) (*orgpb.GetByUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.orgClient.GetByUser(ctx, &orgpb.GetByOwnerRequest{UserUuid: ownerUUID})
}

func (o *OrgRegistry) AddOrg(orgName string, owner string, certificate string, country string,
	currency string) (*orgpb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	organization := &orgpb.Organization{Name: orgName, Owner: owner,
		Certificate: certificate, Country: country, Currency: currency}

	return o.orgClient.Add(ctx, &orgpb.AddRequest{Org: organization})
}

func (o *OrgRegistry) UpdateOrgToUser(orgId string, userId string) (*orgpb.UpdateOrgForUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.orgClient.UpdateOrgForUser(ctx, &orgpb.UpdateOrgForUserRequest{
		UserId: userId,
		OrgId:  orgId,
	})
}

func (o *OrgRegistry) RemoveOrgForUser(orgId string, userId string) (*orgpb.RemoveOrgForUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.orgClient.RemoveOrgForUser(ctx, &orgpb.RemoveOrgForUserRequest{
		UserId: userId,
		OrgId:  orgId,
	})
}
