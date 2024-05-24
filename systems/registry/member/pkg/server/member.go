/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/registry/member/pkg"
	"github.com/ukama/ukama/systems/registry/member/pkg/db"

	log "github.com/sirupsen/logrus"
	metric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"
)

type MemberServer struct {
	pb.UnimplementedMemberServiceServer
	mRepo          db.MemberRepo
	orgClient      cnucl.OrgClient
	userClient     cnucl.UserClient
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
	OrgId          uuid.UUID
	OrgName        string
}

func NewMemberServer(orgName string, mRepo db.MemberRepo, orgClient cnucl.OrgClient, userClient cnucl.UserClient,
	msgBus mb.MsgBusServiceClient, pushGateway string, id uuid.UUID) *MemberServer {

	return &MemberServer{
		mRepo:          mRepo,
		orgClient:      orgClient,
		userClient:     userClient,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		pushGateway:    pushGateway,
		OrgId:          id,
		OrgName:        orgName,
	}
}

/* This will be called by nucleus/org service on every user invited to be member */
func (m *MemberServer) AddMember(ctx context.Context, req *pb.AddMemberRequest) (*pb.MemberResponse, error) {

	// Get the User
	userUUID, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	/* validate user uuid */
	/* Causing a loop when user is getting added a memeber by default */
	// _, err = m.nucleusSystem.GetUserById(userUUID.String())
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get user with id %s. Error %s", userUUID.String(), err.Error())
	// }

	log.Infof("Adding member")
	member := &db.Member{
		MemberId: uuid.NewV4(),
		UserId:   userUUID,
		Role:     db.RoleType(req.Role),
	}

	err = m.mRepo.AddMember(member, m.OrgId.String(), nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	route := m.baseRoutingKey.SetActionCreate().SetObject("member").MustBuild()
	err = m.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	_ = m.PushOrgMemberCountMetric(m.OrgId)

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

/* This is called when user already exists as a member of another org */
func (m *MemberServer) AddOtherMember(ctx context.Context, req *pb.AddMemberRequest) (*pb.MemberResponse, error) {

	// Get the User
	userUUID, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	/* validate user uuid */
	_, err = m.userClient.GetById(userUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user with id %s. Error %s", userUUID.String(), err.Error())
	}

	log.Infof("Adding member")
	member := &db.Member{
		UserId: userUUID,
		Role:   db.RoleType(req.Role),
	}

	err = m.mRepo.AddMember(member, m.OrgId.String(), func(orgId string, userId string) error {
		err := m.orgClient.AddUser(orgId, userId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	route := m.baseRoutingKey.SetActionCreate().SetObject("member").MustBuild()
	err = m.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	_ = m.PushOrgMemberCountMetric(m.OrgId)

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (m *MemberServer) GetMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.FromString(req.GetMemberId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of member uuid. Error %s", err.Error())
	}

	member, err := m.mRepo.GetMember(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (m *MemberServer) GetMemberByUserId(ctx context.Context, req *pb.GetMemberByUserIdRequest) (*pb.GetMemberByUserIdResponse, error) {
	uuid, err := uuid.FromString(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	member, err := m.mRepo.GetMemberByUserId(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	return &pb.GetMemberByUserIdResponse{Member: dbMemberToPbMember(member)}, nil
}

func (m *MemberServer) GetMembers(ctx context.Context, req *pb.GetMembersRequest) (*pb.GetMembersResponse, error) {

	members, err := m.mRepo.GetMembers()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "orgs")
	}

	resp := &pb.GetMembersResponse{
		Members: dbMembersToPbMembers(members),
	}

	return resp, nil
}

func (m *MemberServer) UpdateMember(ctx context.Context, req *pb.UpdateMemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.FromString(req.GetMember().GetMemberId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	member := &db.Member{
		MemberId:    uuid,
		Deactivated: req.Attributes.GetIsDeactivated(),
	}

	err = m.mRepo.UpdateMember(member)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	_ = m.PushOrgMemberCountMetric(m.OrgId)

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (m *MemberServer) RemoveMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.FromString(req.GetMemberId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of member uuid. Error %s", err.Error())
	}

	member, err := m.mRepo.GetMember(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	if !member.Deactivated {
		return nil, status.Errorf(codes.FailedPrecondition,
			"member must be deactivated first")
	}

	err = m.mRepo.RemoveMember(uuid, m.OrgId.String(), func(orgId string, userId string) error {
		err := m.orgClient.RemoveUser(orgId, userId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	route := m.baseRoutingKey.SetActionDelete().SetObject("member").MustBuild()
	err = m.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	_ = m.PushOrgMemberCountMetric(m.OrgId)

	return &pb.MemberResponse{}, nil
}

func (m *MemberServer) PushOrgMemberCountMetric(orgId uuid.UUID) error {

	actMemOrg, inactMemOrg, err := m.mRepo.GetMemberCount()
	if err != nil {
		log.Errorf("failed to get member count for org %s.Error: %s", orgId.String(), err.Error())
		return err
	}

	labels := make(map[string]string)
	labels["org"] = orgId.String()

	err = metric.CollectAndPushSimMetrics(m.pushGateway, pkg.MemberMetric, pkg.NumberOfActiveMembers, float64(actMemOrg), labels, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing active members of Org metric to pushgateway %s", err.Error())
		return err
	}

	err = metric.CollectAndPushSimMetrics(m.pushGateway, pkg.MemberMetric, pkg.NumberOfInactiveMembers, float64(inactMemOrg), labels, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing inactive members Org metric to pushgateway %s", err.Error())
		return err
	}
	return nil
}

func dbMemberToPbMember(member *db.Member) *pb.Member {
	return &pb.Member{
		MemberId:      member.MemberId.String(),
		UserId:        member.UserId.String(),
		IsDeactivated: member.Deactivated,
		Role:          pb.RoleType(member.Role),
		CreatedAt:     timestamppb.New(member.CreatedAt),
	}
}

func dbMembersToPbMembers(members []db.Member) []*pb.Member {
	res := []*pb.Member{}

	for _, m := range members {
		res = append(res, dbMemberToPbMember(&m))
	}

	return res
}
