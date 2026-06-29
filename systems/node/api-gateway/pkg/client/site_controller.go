package client

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SiteController struct {
	conn    *grpc.ClientConn
	client  pb.SiteControllerServiceClient
	timeout time.Duration
}

func NewSiteController(host string, timeout time.Duration) *SiteController {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to SiteController service: %v", err)
	}
	return &SiteController{conn: conn, client: pb.NewSiteControllerServiceClient(conn), timeout: timeout}
}

func NewSiteControllerFromClient(mClient pb.SiteControllerServiceClient) *SiteController {
	return &SiteController{conn: nil, client: mClient, timeout: 1 * time.Second}
}

func (s *SiteController) SetSite(siteID, state, reason, requestedBy string) (*pb.SetSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.SetSite(ctx, &pb.SetSiteRequest{SiteId: siteID, State: state, Reason: reason, RequestedBy: requestedBy})
}

func (s *SiteController) SetService(siteID, state string) (*pb.SetServiceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.SetService(ctx, &pb.SetServiceRequest{SiteId: siteID, State: state})
}

func (s *SiteController) SetRadio(siteID, state string) (*pb.SetRadioResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.SetRadio(ctx, &pb.SetRadioRequest{SiteId: siteID, State: state})
}

func (s *SiteController) GetSiteState(siteID string) (*pb.GetSiteStateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.GetSiteState(ctx, &pb.GetSiteStateRequest{SiteId: siteID})
}

func (s *SiteController) UpsertPortMap(siteID, cnodeID string, ports []*pb.PortMapEntry) (*pb.UpsertPortMapResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.UpsertPortMap(ctx, &pb.UpsertPortMapRequest{SiteId: siteID, CnodeId: cnodeID, Ports: ports})
}

func (s *SiteController) GetPortMap(siteID string) (*pb.GetPortMapResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.GetPortMap(ctx, &pb.GetPortMapRequest{SiteId: siteID})
}

func (s *SiteController) ApplySwitchPolicy(siteID string) (*pb.ApplySwitchPolicyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.ApplySwitchPolicy(ctx, &pb.ApplySwitchPolicyRequest{SiteId: siteID})
}

func (s *SiteController) PowerCycleNode(siteID, role, reason, requestedBy string) (*pb.PowerCycleNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.PowerCycleNode(ctx, &pb.PowerCycleNodeRequest{SiteId: siteID, Role: role, Reason: reason, RequestedBy: requestedBy})
}

func (s *SiteController) RestartSite(siteID string) (*pb.RestartSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.RestartSite(ctx, &pb.RestartSiteRequest{SiteId: siteID, NetworkId: ""})
}

func (s *SiteController) ToggleInternetSwitch(siteID string, status bool, port int32) (*pb.ToggleInternetSwitchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.ToggleInternetSwitch(ctx, &pb.ToggleInternetSwitchRequest{SiteId: siteID, Status: status, Port: port})
}
