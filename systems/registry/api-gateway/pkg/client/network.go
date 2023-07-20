package client

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	"google.golang.org/grpc"
)

const DefaultNetworkName = "default"

type NetworkRegistry struct {
	conn    *grpc.ClientConn
	client  pb.NetworkServiceClient
	timeout time.Duration
	host    string
}

func NewNetworkRegistry(networkHost string, timeout time.Duration) *NetworkRegistry {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, networkHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewNetworkServiceClient(conn)

	return &NetworkRegistry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    networkHost,
	}
}

func NewNetworkRegistryFromClient(networkClient pb.NetworkServiceClient) *NetworkRegistry {
	return &NetworkRegistry{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  networkClient,
	}
}

func (r *NetworkRegistry) Close() {
	r.conn.Close()
}

func (r *NetworkRegistry) AddNetwork(orgName string, netName string) (*netpb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.Add(ctx, &netpb.AddRequest{OrgName: orgName, Name: netName})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *NetworkRegistry) GetNetwork(netID string) (*netpb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.Get(ctx, &netpb.GetRequest{NetworkId: netID})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *NetworkRegistry) GetNetworks(orgID string) (*netpb.GetByOrgResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetByOrg(ctx, &netpb.GetByOrgRequest{OrgId: orgID})
	if err != nil {
		return nil, err
	}

	if res.Networks == nil {
		return &netpb.GetByOrgResponse{Networks: []*netpb.Network{}, OrgId: orgID}, nil
	}

	return res, nil
}

func (r *NetworkRegistry) AddSite(netID string, siteName string) (*netpb.AddSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.AddSite(ctx, &netpb.AddSiteRequest{NetworkId: netID, SiteName: siteName})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *NetworkRegistry) GetSite(netID string, siteName string) (*netpb.GetSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetSiteByName(ctx, &netpb.GetSiteByNameRequest{NetworkId: netID, SiteName: siteName})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *NetworkRegistry) GetSites(netID string) (*netpb.GetSitesByNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetSitesByNetwork(ctx, &netpb.GetSitesByNetworkRequest{NetworkId: netID})
	if err != nil {
		return nil, err
	}

	if res.Sites == nil {
		return &netpb.GetSitesByNetworkResponse{Sites: []*netpb.Site{}, NetworkId: netID}, nil
	}

	return res, nil
}
