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

func (r *InvitationRegistry) RemoveInvitation(invitationId string) (*pb.DeleteInvitationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.Delete(ctx, &pb.DeleteInvitationRequest{Id: invitationId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InvitationRegistry) GetInvitationById(id string) (*pb.GetInvitationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.Get(ctx, &pb.GetInvitationRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InvitationRegistry) AddInvitation(org, name, email, role string) (*pb.AddInvitationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	res, err := r.client.Add(ctx, &pb.AddInvitationRequest{
		Org:   org,
		Name:  name,
		Email: email,
		Role:  pb.RoleType(pb.RoleType_value[role]),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InvitationRegistry) GetInvitationByOrg(org string) (*pb.GetInvitationByOrgResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	invitation := &pb.GetInvitationByOrgRequest{
		Org: org,
	}
	res, err := r.client.GetByOrg(ctx, invitation)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *InvitationRegistry) UpdateInvitation(id, status string) (*pb.UpdateInvitationStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.UpdateStatus(ctx, &pb.UpdateInvitationStatusRequest{
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

	_, err := r.client.GetInvitationByEmail(ctx, &pb.GetInvitationByEmailRequest{
		Email: email,
	},
	)
	return err
}
