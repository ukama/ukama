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

	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	"google.golang.org/grpc"
)

type MemberRegistry struct {
	conn    *grpc.ClientConn
	client  pb.MemberServiceClient
	timeout time.Duration
	host    string
}

func NewMemberRegistry(memberHost string, timeout time.Duration) *MemberRegistry {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, memberHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewMemberServiceClient(conn)

	return &MemberRegistry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    memberHost,
	}
}

func NewRegistryFromClient(mClient pb.MemberServiceClient) *MemberRegistry {
	return &MemberRegistry{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *MemberRegistry) Close() {
	r.conn.Close()
}

func (r *MemberRegistry) GetMember(memberId string) (*pb.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetMember(ctx, &pb.MemberRequest{MemberId: memberId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MemberRegistry) GetMembers() (*pb.GetMembersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetMembers(ctx, &pb.GetMembersRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MemberRegistry) AddMember(userUUID string, role string) (*pb.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	member := &pb.AddMemberRequest{UserUuid: userUUID, Role: pb.RoleType(pb.RoleType_value[role])}
	res, err := r.client.AddMember(ctx, member)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MemberRegistry) AddOtherMember(userUUID string, role string) (*pb.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	member := &pb.AddMemberRequest{UserUuid: userUUID, Role: pb.RoleType(pb.RoleType_value[role])}
	res, err := r.client.AddOtherMember(ctx, member)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MemberRegistry) UpdateMember(memberId string, isDeactivated bool, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.client.UpdateMember(ctx, &pb.UpdateMemberRequest{
		Member:     &pb.MemberRequest{MemberId: memberId},
		Attributes: &pb.MemberAttributes{IsDeactivated: isDeactivated}})

	return err
}

func (r *MemberRegistry) RemoveMember(memberId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.client.RemoveMember(ctx, &pb.MemberRequest{MemberId: memberId})

	return err
}
