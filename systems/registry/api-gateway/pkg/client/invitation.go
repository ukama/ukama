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

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	"google.golang.org/grpc"
)

type InvitationRegistry struct {
	conn    *grpc.ClientConn
	client  pb.InvitationServiceClient
	timeout time.Duration
	host    string
}

func NewInvitationRegistry(invitationHost string, timeout time.Duration) *InvitationRegistry {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, invitationHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewInvitationServiceClient(conn)

	return &InvitationRegistry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    invitationHost,
	}
}

func NewInvitationRegistryFromClient(mClient pb.InvitationServiceClient) *InvitationRegistry {
	return &InvitationRegistry{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *InvitationRegistry) Close() {
	r.conn.Close()
}

func (r *InvitationRegistry) RemoveInvitation(invitationId string) (*pb.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.Delete(ctx, &pb.DeleteRequest{Id: invitationId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InvitationRegistry) GetInvitationById(id string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.Get(ctx, &pb.GetRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InvitationRegistry) AddInvitation(name, email, role string) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	res, err := r.client.Add(ctx, &pb.AddRequest{
		Name:  name,
		Email: email,
		Role:  pb.RoleType(pb.RoleType_value[role]),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InvitationRegistry) GetAllInvitations() (*pb.GetAllResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	invitation := &pb.GetAllRequest{}
	res, err := r.client.GetAll(ctx, invitation)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InvitationRegistry) UpdateInvitation(id, status string) (*pb.UpdateStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.UpdateStatus(ctx, &pb.UpdateStatusRequest{
		Id:     id,
		Status: pb.StatusType(pb.StatusType_value[status]),
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InvitationRegistry) GetInvitationByEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.client.GetByEmail(ctx, &pb.GetByEmailRequest{
		Email: email,
	},
	)
	return err
}
