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
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"
)

type MemberRegistry struct {
	conn    *grpc.ClientConn
	client  pb.MemberServiceClient
	timeout time.Duration
	host    string
}

func NewMemberRegistry(memberHost string, timeout time.Duration) *MemberRegistry {
	conn, err := grpc.NewClient(memberHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Member Server: %v", err)
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

func (m *MemberRegistry) Close() {
	if m.conn != nil {
		err := m.conn.Close()
		if err != nil {
			log.Warnf("Failed to gracefully close Member Service connection: %v", err)
		}
	}
}

func (m *MemberRegistry) GetMember(memberId string) (*pb.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	return m.client.GetMember(ctx, &pb.MemberRequest{MemberId: memberId})
}

func (m *MemberRegistry) GetMemberByUserId(userId string) (*pb.GetMemberByUserIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	return m.client.GetMemberByUserId(ctx, &pb.GetMemberByUserIdRequest{MemberId: userId})
}

func (m *MemberRegistry) GetMembers() (*pb.GetMembersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	return m.client.GetMembers(ctx, &pb.GetMembersRequest{})
}

func (m *MemberRegistry) AddMember(userUUID string, role string) (*pb.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	member := &pb.AddMemberRequest{UserUuid: userUUID, Role: upb.RoleType(upb.RoleType_value[role])}
	return m.client.AddMember(ctx, member)
}

func (m *MemberRegistry) UpdateMember(memberId string, isDeactivated bool, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.client.UpdateMember(ctx, &pb.UpdateMemberRequest{
		MemberId:      memberId,
		IsDeactivated: isDeactivated,
	})

	return err
}

func (m *MemberRegistry) RemoveMember(memberId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.client.RemoveMember(ctx, &pb.MemberRequest{MemberId: memberId})

	return err
}
