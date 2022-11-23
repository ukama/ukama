package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/ukama/ukama/systems/common/rest"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	pborg "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"google.golang.org/grpc"
)

const DefaultNetworkName = "default"

type Registry struct {
	conn      *grpc.ClientConn
	orgConn   *grpc.ClientConn
	client    pb.NetworkServiceClient
	orgClient pborg.OrgServiceClient
	timeout   time.Duration
	host      string
}

func NewRegistry(networkHost string, orgHost string, timeout time.Duration) *Registry {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewNetworkServiceClient(conn)

	orgConn, err := grpc.DialContext(ctx, orgHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	orgClient := pborg.NewOrgServiceClient(orgConn)

	return &Registry{
		conn:      conn,
		client:    client,
		orgConn:   orgConn,
		orgClient: orgClient,
		timeout:   timeout,
		host:      networkHost,
	}
}

func NewRegistryFromClient(networkClient pb.NetworkServiceClient, orgClient pborg.OrgServiceClient) *Registry {
	return &Registry{
		host:      "localhost",
		timeout:   1 * time.Second,
		conn:      nil,
		client:    networkClient,
		orgClient: orgClient,
	}
}

func (r *Registry) Close() {
	r.conn.Close()
	r.orgConn.Close()
}

func (r *Registry) GetOrg(orgName string) (*pborg.Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.GetByName(ctx, &pborg.GetByNameRequest{Name: orgName})
	if err != nil {
		return nil, err
	}

	return res.Org, nil
}

func (r *Registry) GetOrgs(ownerUUID string) (*pborg.GetByOwnerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.GetByOwner(ctx, &pborg.GetByOwnerRequest{UserUuid: ownerUUID})
	if err != nil {
		return nil, err
	}

	if res.Orgs == nil {
		return &pborg.GetByOwnerResponse{Orgs: []*pborg.Organization{}, Owner: ownerUUID}, nil
	}

	return res, nil
}

func (r *Registry) AddOrg(orgName string, owner string) (*pborg.Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	organization := &pborg.Organization{Name: orgName, Owner: owner}
	res, err := r.orgClient.Add(ctx, &pborg.AddRequest{Org: organization})

	if err != nil {
		return nil, err
	}

	return res.Org, nil
}

func (r *Registry) IsAuthorized(userId string, org string) (bool, error) {
	orgResp, err := r.GetOrg(org)
	if err != nil {
		if gErr, ok := err.(rest.HttpError); ok {
			if gErr.HttpCode != http.StatusNotFound {
				return false, nil
			}

			return false, gErr
		} else {
			return false, fmt.Errorf(err.Error())
		}
	}
	if orgResp.Owner == userId {
		return true, nil
	}
	return false, nil
}
