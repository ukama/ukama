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

func (s *SiteController) SetSite(siteID, state, reason string) (*pb.SetSiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.SetSite(ctx, &pb.SetSiteRequest{SiteId: siteID, State: state, Reason: reason})
}

func (s *SiteController) SetService(siteID, state, reason string) (*pb.SetServiceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.SetService(ctx, &pb.SetServiceRequest{SiteId: siteID, State: state, Reason: reason})
}

func (s *SiteController) SetRadio(siteID, state, reason string) (*pb.SetRadioResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.SetRadio(ctx, &pb.SetRadioRequest{SiteId: siteID, State: state, Reason: reason})
}

func (s *SiteController) GetSiteState(siteID string) (*pb.GetSiteStateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.GetSiteState(ctx, &pb.GetSiteStateRequest{SiteId: siteID})
}

func (s *SiteController) GetSwitchPolicy(siteID string) (*pb.GetSwitchPolicyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.GetSwitchPolicy(ctx, &pb.GetSwitchPolicyRequest{SiteId: siteID})
}

func (s *SiteController) RefreshSwitchPolicy(siteID, cnodeID, reason string) (*pb.RefreshSwitchPolicyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.RefreshSwitchPolicy(ctx, &pb.RefreshSwitchPolicyRequest{SiteId: siteID, CnodeId: cnodeID, Reason: reason})
}

func (s *SiteController) ReportSwitchPolicy(siteID, cnodeID string, policy *pb.SwitchPolicy) (*pb.ReportSwitchPolicyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.ReportSwitchPolicy(ctx, &pb.ReportSwitchPolicyRequest{SiteId: siteID, CnodeId: cnodeID, Policy: policy})
}

func (s *SiteController) PowerCycleNode(siteID, role, reason string) (*pb.PowerCycleNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.PowerCycleNode(ctx, &pb.PowerCycleNodeRequest{SiteId: siteID, Role: role, Reason: reason})
}
