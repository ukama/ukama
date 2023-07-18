package client

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"

	orgpb "github.com/ukama/ukama/systems/nucleus/orgs/pb/gen"
	"google.golang.org/grpc"
)

const DefaultNetworkName = "default"

type Registry struct {
	conn      *grpc.ClientConn
	orgClient orgpb.OrgServiceClient
	timeout   time.Duration
}

func NewRegistry(networkHost string, orgHost string, timeout time.Duration) *Registry {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	orgConn, err := grpc.DialContext(ctx, orgHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	orgClient := orgpb.NewOrgServiceClient(orgConn)

	return &Registry{
		conn:      orgConn,
		orgClient: orgClient,
		timeout:   timeout,
	}
}

func NewRegistryFromClient(orgClient orgpb.OrgServiceClient) *Registry {
	return &Registry{
		timeout:   1 * time.Second,
		conn:      nil,
		orgClient: orgClient,
	}
}

func (r *Registry) Close() {
	r.conn.Close()
}

func (r *Registry) GetOrg(orgName string) (*orgpb.GetByNameResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.GetByName(ctx, &orgpb.GetByNameRequest{Name: orgName})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Registry) GetOrgs(ownerUUID string) (*orgpb.GetByOwnerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.GetByOwner(ctx, &orgpb.GetByOwnerRequest{UserUuid: ownerUUID})
	if err != nil {
		return nil, err
	}

	if res.Orgs == nil {
		return &orgpb.GetByOwnerResponse{Orgs: []*orgpb.Organization{}, Owner: ownerUUID}, nil
	}

	return res, nil
}

func (r *Registry) AddOrg(orgName string, owner string, certificate string) (*orgpb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	organization := &orgpb.Organization{Name: orgName, Owner: owner, Certificate: certificate}
	res, err := r.orgClient.Add(ctx, &orgpb.AddRequest{Org: organization})

	if err != nil {
		return nil, err
	}

	return res, nil
}
