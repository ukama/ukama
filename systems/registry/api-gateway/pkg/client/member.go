package client

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	"google.golang.org/grpc"
)

type MemberRegistry struct {
	conn    *grpc.ClientConn
	client  pb.MemberServiceClient
	timeout time.Duration
	host    string
}

func NewMemberRegistry(networkHost string, timeout time.Duration) *MemberRegistry {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewMemberServiceClient(conn)

	return &MemberRegistry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    networkHost,
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

func (r *MemberRegistry) GetMember(userUUID string) (*pb.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetMember(ctx, &pb.MemberRequest{UserUuid: userUUID})
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

func (r *MemberRegistry) UpdateMember(userUUID string, isDeactivated bool, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.client.UpdateMember(ctx, &pb.UpdateMemberRequest{
		Member:     &pb.MemberRequest{UserUuid: userUUID},
		Attributes: &pb.MemberAttributes{IsDeactivated: isDeactivated}})

	return err
}

func (r *MemberRegistry) RemoveMember(userUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.client.RemoveMember(ctx, &pb.MemberRequest{UserUuid: userUUID})

	return err
}
