package client

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"google.golang.org/grpc"
)

const DefaultNetworkName = "default"

type Registry struct {
	conn          *grpc.ClientConn
	orgConn       *grpc.ClientConn
	networkClient pb.NetworkServiceClient
	orgClient     orgpb.OrgServiceClient
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

	orgConn, err := grpc.DialContext(ctx, orgHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	orgClient := orgpb.NewOrgServiceClient(orgConn)

	return &Registry{
		conn:          conn,
		networkClient: client,
		orgConn:       orgConn,
		orgClient:     orgClient,
		timeout:       timeout,
		host:          networkHost,
	}
}

func NewRegistryFromClient(networkClient pb.NetworkServiceClient, orgClient orgpb.OrgServiceClient) *Registry {
	return &Registry{
		host:          "localhost",
		timeout:       1 * time.Second,
		conn:          nil,
		networkClient: networkClient,
		orgClient:     orgClient,
	}
}

func (r *Registry) Close() {
	r.conn.Close()
	r.orgConn.Close()
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

func (r *Registry) GetMember(orgName string, userUUID string) (*orgpb.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.GetMember(ctx, &orgpb.MemberRequest{OrgName: orgName, UserUuid: userUUID})
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (r *Registry) GetMemberRole(orgId string, userUUID string) (*orgpb.GetMemberRoleResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	res, err :=r.orgClient.GetMemberRole(ctx, &orgpb.MemberRoleRequest{OrgId: orgId, UserUuid: userUUID})
	if err != nil {
		return nil, err
	}
	return res, nil
}


func (r *Registry) GetMembers(orgName string) (*orgpb.GetMembersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.orgClient.GetMembers(ctx, &orgpb.GetMembersRequest{OrgName: orgName})
	if err != nil {
		return nil, err
	}

	if res.Members == nil {
		return &orgpb.GetMembersResponse{Members: []*orgpb.OrgUser{}, Org: orgName}, nil
	}

	return res, nil
}

func (r *Registry) AddMember(orgName string, userUUID string,role string) (*orgpb.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	member := &orgpb.MemberRequest{OrgName: orgName, UserUuid: userUUID}
	res, err := r.orgClient.AddMember(ctx, member)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Registry) UpdateMember(orgName string, userUUID string, isDeactivated bool,role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.orgClient.UpdateMember(ctx, &orgpb.UpdateMemberRequest{
		Member:     &orgpb.MemberRequest{OrgName: orgName, UserUuid: userUUID},
		Attributes: &orgpb.OrgUserAttributes{IsDeactivated: isDeactivated}})

	return err
}

func (r *Registry) RemoveMember(orgName string, userUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.orgClient.RemoveMember(ctx, &orgpb.MemberRequest{OrgName: orgName, UserUuid: userUUID})

	return err
}

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
