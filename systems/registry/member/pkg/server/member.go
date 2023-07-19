package server

import (
	"context"
	"fmt"

	"github.com/ukama/ukama/systems/common/grpc"
	metric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/member/pkg"
	"github.com/ukama/ukama/systems/registry/member/pkg/db"
	"github.com/ukama/ukama/systems/registry/member/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"

	log "github.com/sirupsen/logrus"
)

const uuidParsingError = "Error parsing UUID"

type MemberServer struct {
	pb.UnimplementedMemberServiceServer
	mRepo          db.MemberRepo
	orgService     providers.OrgClientProvider
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pushGateway    string
	OrgId          uuid.UUID
	OrgName        string
}

func NewMemberServer(mRepo db.MemberRepo, orgService providers.OrgClientProvider, msgBus mb.MsgBusServiceClient, pushGateway string, id uuid.UUID, name string) *MemberServer {

	return &MemberServer{
		mRepo:          mRepo,
		orgService:     orgService,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		pushGateway:    pushGateway,
		OrgId:          id,
		OrgName:        name,
	}
}

func (m *MemberServer) AddMember(ctx context.Context, req *pb.AddMemberRequest) (*pb.MemberResponse, error) {

	// Get the User
	userUUID, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	/* validate user uuid */
	_, err = m.orgService.GetUserById(userUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user with id %s. Error %s", userUUID.String(), err.Error())
	}

	log.Infof("Adding member")
	member := &db.Member{
		UserId: userUUID,
		Role:   db.RoleType(req.Role),
	}

	err = m.mRepo.AddMember(member)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	route := m.baseRoutingKey.SetActionCreate().SetObject("member").MustBuild()
	err = m.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	_ = m.PushOrgMemberCountMetric(m.OrgId)

	return &pb.MemberResponse{Member: dbMemberToPbMember(member, m.OrgId.String())}, nil
}

func (m *MemberServer) GetMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	member, err := m.mRepo.GetMember(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	return &pb.MemberResponse{Member: dbMemberToPbMember(member, m.OrgId.String())}, nil
}

func (m *MemberServer) GetMembers(ctx context.Context, req *pb.GetMembersRequest) (*pb.GetMembersResponse, error) {

	members, err := m.mRepo.GetMembers()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "orgs")
	}

	resp := &pb.GetMembersResponse{
		Members: dbMembersToPbMembers(members, m.OrgId.String()),
	}

	return resp, nil
}

func (m *MemberServer) UpdateMember(ctx context.Context, req *pb.UpdateMemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.FromString(req.GetMember().GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	member := &db.Member{
		UserId:      uuid,
		Deactivated: req.GetAttributes().IsDeactivated,
	}

	err = m.mRepo.UpdateMember(member)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	_ = m.PushOrgMemberCountMetric(m.OrgId)

	return &pb.MemberResponse{Member: dbMemberToPbMember(member, m.OrgId.String())}, nil
}

func (m *MemberServer) RemoveMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	member, err := m.mRepo.GetMember(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	if !member.Deactivated {
		return nil, status.Errorf(codes.FailedPrecondition,
			"member must be deactivated first")
	}

	err = m.mRepo.RemoveMember(uuid)
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

func dbMemberToPbMember(member *db.Member, orgId string) *pb.Member {
	return &pb.Member{
		OrgId:         orgId,
		UserId:        member.UserId.String(),
		IsDeactivated: member.Deactivated,
		CreatedAt:     timestamppb.New(member.CreatedAt),
	}
}

func dbMembersToPbMembers(members []db.Member, orgId string) []*pb.Member {
	res := []*pb.Member{}

	for _, m := range members {
		res = append(res, dbMemberToPbMember(&m, orgId))
	}

	return res
}
