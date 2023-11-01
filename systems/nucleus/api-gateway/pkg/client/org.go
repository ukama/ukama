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

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/credentials/insecure"

	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	"google.golang.org/grpc"
)

const DefaultNetworkName = "default"

type OrgRegistry struct {
	conn      *grpc.ClientConn
	orgClient orgpb.OrgServiceClient
	timeout   time.Duration
}

func NewOrgRegistry(orgHost string, timeout time.Duration) *OrgRegistry {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	orgConn, err := grpc.DialContext(ctx, orgHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
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

func (r *OrgRegistry) Close() {
	r.conn.Close()
}

func (r *OrgRegistry) GetOrg(orgName string) (*orgpb.GetByNameResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.GetByName(ctx, &orgpb.GetByNameRequest{Name: orgName})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *OrgRegistry) GetOrgs(ownerUUID string) (*orgpb.GetByUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.GetByUser(ctx, &orgpb.GetByOwnerRequest{UserUuid: ownerUUID})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *OrgRegistry) AddOrg(orgName string, owner string, certificate string) (*orgpb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	organization := &orgpb.Organization{Name: orgName, Owner: owner, Certificate: certificate}
	res, err := r.orgClient.Add(ctx, &orgpb.AddRequest{Org: organization})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *OrgRegistry) UpdateOrgToUser(orgId string, userId string) (*orgpb.UpdateOrgForUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.UpdateOrgForUser(ctx, &orgpb.UpdateOrgForUserRequest{
		UserId: userId,
		OrgId:  orgId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *OrgRegistry) RemoveOrgForUser(orgId string, userId string) (*orgpb.RemoveOrgForUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.RemoveOrgForUser(ctx, &orgpb.RemoveOrgForUserRequest{
		UserId: userId,
		OrgId:  orgId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
