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

type Registry struct {
	conn          *grpc.ClientConn
	orgConn       *grpc.ClientConn
	networkClient pb.NetworkServiceClient
	timeout       time.Duration
	host          string
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

	return &Registry{
		conn:          conn,
		networkClient: client,
		timeout:       timeout,
		host:          networkHost,
	}
}

func NewRegistryFromClient(networkClient pb.NetworkServiceClient) *Registry {
	return &Registry{
		host:          "localhost",
		timeout:       1 * time.Second,
		conn:          nil,
		networkClient: networkClient,
	}
}

func (r *Registry) Close() {
	r.conn.Close()
	r.orgConn.Close()
}

// func (r *Registry) GetMember(orgName string, userUUID string) (*orgpb.MemberResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
// 	defer cancel()

// 	res, err := r.orgClient.GetMember(ctx, &orgpb.MemberRequest{OrgName: orgName, UserUuid: userUUID})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }

// func (r *Registry) GetMembers(orgName string) (*orgpb.GetMembersResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
// 	defer cancel()

// 	res, err := r.orgClient.GetMembers(ctx, &orgpb.GetMembersRequest{OrgName: orgName})
// 	if err != nil {
// 		return nil, err
// 	}

// 	if res.Members == nil {
// 		return &orgpb.GetMembersResponse{Members: []*orgpb.OrgUser{}, Org: orgName}, nil
// 	}

// 	return res, nil
// }

// func (r *Registry) AddMember(orgName string, userUUID string, role string) (*orgpb.MemberResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
// 	defer cancel()

// 	member := &orgpb.AddMemberRequest{OrgName: orgName, UserUuid: userUUID, Role: orgpb.RoleType(orgpb.RoleType_value[role])}
// 	res, err := r.orgClient.AddMember(ctx, member)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }

// func (r *Registry) UpdateMember(orgName string, userUUID string, isDeactivated bool, role string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
// 	defer cancel()

// 	_, err := r.orgClient.UpdateMember(ctx, &orgpb.UpdateMemberRequest{
// 		Member:     &orgpb.MemberRequest{OrgName: orgName, UserUuid: userUUID},
// 		Attributes: &orgpb.OrgUserAttributes{IsDeactivated: isDeactivated}})

// 	return err
// }

// func (r *Registry) RemoveMember(orgName string, userUUID string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
// 	defer cancel()

// 	_, err := r.orgClient.RemoveMember(ctx, &orgpb.MemberRequest{OrgName: orgName, UserUuid: userUUID})

// 	return err
// }

func (r *Registry) AddNetwork(orgName string, netName string) (*netpb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.networkClient.Add(ctx, &netpb.AddRequest{OrgName: orgName, Name: netName})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Registry) GetNetwork(netID string) (*netpb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.networkClient.Get(ctx, &netpb.GetRequest{NetworkId: netID})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Registry) GetNetworks(orgID string) (*netpb.GetByOrgResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.networkClient.GetByOrg(ctx, &netpb.GetByOrgRequest{OrgId: orgID})
	if err != nil {
		return nil, err
	}

	if res.Networks == nil {
		return &netpb.GetByOrgResponse{Networks: []*netpb.Network{}, OrgId: orgID}, nil
	}

	return res, nil
}

func (r *Registry) AddSite(netID string, siteName string) (*netpb.AddSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.networkClient.AddSite(ctx, &netpb.AddSiteRequest{NetworkId: netID, SiteName: siteName})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Registry) GetSite(netID string, siteName string) (*netpb.GetSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.networkClient.GetSiteByName(ctx, &netpb.GetSiteByNameRequest{NetworkId: netID, SiteName: siteName})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Registry) GetSites(netID string) (*netpb.GetSitesByNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.networkClient.GetSitesByNetwork(ctx, &netpb.GetSitesByNetworkRequest{NetworkId: netID})
	if err != nil {
		return nil, err
	}

	if res.Sites == nil {
		return &netpb.GetSitesByNetworkResponse{Sites: []*netpb.Site{}, NetworkId: netID}, nil
	}

	return res, nil
}
