package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/session"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/store"
)

type Controller struct {
	store *store.Store
	sm    session.SessionManager
	rc    client.RemoteController
}

func NewController(db string, br pkg.BrdigeConfig, remote string, debug bool) (*Controller, error) {
	c := &Controller{}
	store, err := store.NewStore(db)
	if err != nil {
		log.Errorf("Failed to create db: %v", err)
		return nil, err
	}

	rc, err := client.NewRemoteControllerClient(remote, debug)
	if err != nil {
		log.Errorf("Failed to create client: %v", err)
		return nil, err
	}

	c.sm = session.NewSessionManager(rc, store, br.Name, br.Ip, br.NetType, br.Period)
	c.store = store
	c.rc = rc

	return c, nil
}

func sessionResponse(s *store.Session) *api.SessionResponse {
	return &api.SessionResponse{
		ID:         s.ID,
		Imsi:       s.SusbcriberID.Imsi,
		ApnName:    s.ApnName,
		UeIpaddr:   s.UeIpaddr,
		StartTime:  s.StartTime,
		EndTime:    s.EndTime,
		TxBytes:    s.TxBytes,
		RxBytes:    s.RxBytes,
		TotalBytes: s.TotalBytes,
		TxMeterId:  uint32(s.TXMeterId.ID),
		RxMeterId:  uint32(s.RXMeterId.ID),
		State:      s.State.String(),
		Sync:       s.Sync.String(),
	}
}

func policyResponse(p *store.Policy) *api.PolicyResponse {
	return &api.PolicyResponse{
		ID:        p.ID,
		Data:      p.Data,
		Dlbr:      p.Dlbr,
		Ulbr:      p.Ulbr,
		StartTime: p.StartTime,
		EndTime:   p.EndTime,
	}
}

func flowResponse(flows []*store.Flow) []*api.FlowResponse {
	fr := make([]*api.FlowResponse, len(flows))
	for i, flow := range flows {
		fr[i] = &api.FlowResponse{
			ID:        flow.ID,
			Cookie:    flow.Cookie,
			Table:     flow.Table,
			Priority:  flow.Priority,
			UeIpaddr:  flow.UeIpaddr,
			ReRouting: flow.ReRouting.Ipaddr,
			MeterID:   uint32(flow.MeterID.ID),
		}
	}
	return fr
}

func reRouteResponse(route *store.ReRoute) *api.ReRouteResponse {
	return &api.ReRouteResponse{
		ID: route.ID,
		Ip: route.Ipaddr,
	}
}

func subscriberResponse(s *store.Subscriber) *api.SubscriberResponse {
	return &api.SubscriberResponse{
		ID:       s.ID,
		Imsi:     s.Imsi,
		PolicyID: s.PolicyID.ID,
	}
}

func (c *Controller) validateSusbcriber(imsi string) error {
	/* Get subscriber policy by imsi*/
	s, err := c.store.GetSubscriber(imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber for %s:Error: %v", imsi, err)
		return err
	}

	/* store policy */
	p, err := c.store.GetPolicyByID(s.PolicyID.ID)
	if err != nil {
		log.Errorf("Failed to get subscriber policy %d:Error: %v", s.PolicyID, err)
		return err
	}

	now := time.Now().Unix()
	if p.StartTime > now && p.EndTime <= now {
		return fmt.Errorf("failed to get valid policy")
	}

	return nil
}

func (c *Controller) updateSubscriberPolicy(imsi string, p *api.Policy) (*store.Subscriber, error) {
	var sub *store.Subscriber

	sub, err := c.store.CreateSubscriber(imsi, p)
	if err != nil {
		log.Errorf("Failed to create subscriber %s:Error: %v", imsi, err)
		return nil, err
	}

	return sub, nil
}

func (c *Controller) CreateSession(ctx *gin.Context, req *api.CreateSession) error {
	var sub *store.Subscriber
	var err error
	/* validate subscriber*/
	err = c.validateSusbcriber(req.Imsi)
	if err != nil {
		/* Get subscriber policy from remote */
		p, err := c.rc.GetPolicy(req.Imsi)
		if err != nil {
			log.Errorf("Failed to get subscriber %s policy.Error: %v", req.Imsi, err)
			return err
		}

		sub, err = c.updateSubscriberPolicy(req.Imsi, p)
		if err != nil {
			log.Errorf("Failed to update subscriber %s with policy %s.Error: %v", req.Imsi, p.Uuid.String(), err)
			return err
		}
	}

	/* create session */
	s, rxF, txF, err := c.store.CreateSession(sub, req.Ip)
	if err != nil {
		log.Errorf("Failed to create a session for subscriber %s:Error: %v", req.Imsi, err)
		return err
	}

	/* create UE data path and monitoring session */
	err = c.sm.CreateSesssion(ctx, sub, s, rxF, txF)
	if err != nil {
		log.Errorf("Failed to monitor session on bridge for subscriber %s:Error: %v", req.Imsi, err)
		return err
	}

	return nil
}

func (c *Controller) EndSession(ctx *gin.Context, req *api.EndSession) error {

	sub, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber for imsi %s:Error: %v", req.Imsi, err)
		return err
	}

	err = c.sm.EndSesssion(ctx, sub)
	if err != nil {
		log.Errorf("Failed to end session on bridge for subscriber %s:Error: %v", req.Imsi, err)
		return err
	}
	return nil
}

func (c *Controller) GetSessionByID(ctx *gin.Context, req *api.GetSessionByID) (*api.SessionResponse, error) {
	s, err := c.store.GetSessionByID(int(req.ID))
	if err != nil {
		log.Errorf("failed to get session with id %d:Error: %v", req.ID, err)
		return nil, err
	}
	return sessionResponse(s), nil
}

func (c *Controller) GetActiveSessionByImsi(ctx *gin.Context, req *api.GetSessionByImsi) (*api.SessionResponse, error) {
	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("failed to get active session for Imsi %s:Error: %v", req.Imsi, err)
		return nil, err
	}
	return sessionResponse(s), nil
}

func (c *Controller) GetCDRBySessionId(ctx *gin.Context, req *api.GetCDRBySessionId) (*api.CDR, error) {
	s, err := c.store.GetSessionByID(int(req.ID))
	if err != nil {
		log.Errorf("failed to get session with id %d:Error: %v", req.ID, err)
		return nil, err
	}

	cdr := store.PrepareCDR(s)

	return cdr, nil
}

func (c *Controller) GetCDRByImsi(ctx *gin.Context, req *api.GetCDRByImsi) ([]*api.CDR, error) {
	cdrs := []*api.CDR{}
	sess, err := c.store.GetSessionsByImsi(req.Imsi)
	if err != nil {
		log.Errorf("failed to get session for Imsi %s:Error: %v", req.Imsi, err)
		return nil, err
	}

	for i, s := range sess {
		cdrs[i] = store.PrepareCDR(&s)
	}

	return cdrs, nil
}

func (c *Controller) GetPolicyByImsi(ctx *gin.Context, req *api.PolicyByImsi) (*api.PolicyResponse, error) {
	p, err := c.store.GetApplicablePolicyByImsi(req.Imsi)
	if err != nil {
		log.Errorf("failed to get policy for Imsi %s:Error: %v", req.Imsi, err.Error())
		return nil, err
	}
	return policyResponse(p), nil
}

func (c *Controller) AddPolicy(ctx *gin.Context, req *api.AddPolicyByImsi) error {
	p, err := c.store.CreatePolicy(&req.Policy)
	if err != nil {
		log.Errorf("failed to add policy %s", req.Policy)
		return err
	}

	sub, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber with Imsi %s:Error: %v")
		return err
	}

	/* Update policy for subscriber */
	err = c.store.UpdateSubscriber(sub, p.ID)
	if err != nil {
		log.Errorf("failed to update policy for subscriber with Imsi %s:Error: %v")
		return err
	}
	return nil
}

func (c *Controller) GetFlowForImsi(ctx *gin.Context, req *api.GetFlowsForImsi) ([]*api.FlowResponse, error) {
	var flows []*store.Flow
	_, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber with Imsi %s:Error: %v")
		return nil, err
	}

	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("failed to get active session for Imsi %s:Error: %v", req.Imsi, err)
		return nil, err
	}

	fRx, err := c.store.GetFlowForMeter(s.RXMeterId.ID)
	if err != nil {
		log.Errorf("failed to get RX flow for Imsi %s:Error: %v", req.Imsi, err)
		return nil, err
	}
	flows = append(flows, fRx)

	fTx, err := c.store.GetFlowForMeter(s.TXMeterId.ID)
	if err != nil {
		log.Errorf("failed to get TX flow for Imsi %s:Error: %v", req.Imsi, err)
		return nil, err
	}
	flows = append(flows, fTx)

	return flowResponse(flows), nil
}

func (c *Controller) GetReroute(ctx *gin.Context, req *api.GetReRouteByImsi) (*api.ReRouteResponse, error) {

	_, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber with Imsi %s:Error: %v")
		return nil, err
	}

	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("failed to get active session for Imsi %s:Error: %v", req.Imsi, err)
		return nil, err
	}

	flow, err := c.store.GetFlowForMeter(s.TXMeterId.ID)
	if err != nil {
		log.Errorf("failed to get TX flow for Imsi %s:Error: %v", req.Imsi, err)
		return nil, err
	}

	r, err := c.store.GetReRouteByID(flow.ReRouting.ID)
	if err != nil {
		log.Errorf("failed to get reroute for imsi %s. Error %s", req.Imsi, err.Error())
		return nil, err
	}

	return reRouteResponse(r), nil
}

func (c *Controller) UpdateReroute(ctx *gin.Context, req *api.UpdateRerouteById) error {

	err := c.store.UpdateReroute(&store.ReRoute{
		ID:     int(req.Id),
		Ipaddr: req.Ip,
	})
	if err != nil {
		log.Errorf("failed to update route for Id %d. Error: %s", req.Id, err.Error())
		return err
	}
	return nil
}

func (c *Controller) GetSubscriber(ctx *gin.Context, req *api.RequestSubscriber) (*api.SubscriberResponse, error) {
	s, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber with imsi %s. Error: %s", req.Imsi, err.Error)
		return nil, err
	}
	return subscriberResponse(s), nil
}

func (c *Controller) DeleteSubscriber(ctx *gin.Context, req *api.RequestSubscriber) error {
	s, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber with imsi %s. Error: %s", req.Imsi, err.Error)
		return err
	}

	err = c.store.DeleteSubscriber(s)
	if err != nil {
		log.Errorf("failed to delete subscriber with imsi %s. Error: %s", req.Imsi, err.Error)
		return err
	}
	return nil
}

func (c *Controller) AddSubscriber(ctx *gin.Context, req *api.CreateSubscriber) error {
	_, err := c.store.CreateSubscriber(req.Imsi, &req.Policy)
	if err != nil {
		log.Errorf("failed to delete subscriber with imsi %s. Error: %s", req.Imsi, err.Error())
		return err
	}
	return nil
}
