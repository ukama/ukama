/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ukama/ukama/nodes/apps/pcrf/pkg"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/session"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/datapath"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

type Controller struct {
	store     *store.Store
	sm        session.SessionManager
	rc        client.RemoteController
	publisher *Publisher
	nodeId    string
}

type ControllerStatus struct {
	NodeID   string          `json:"nodeId"`
	Ready    bool            `json:"ready"`
	State    string          `json:"state"`
	Reason   string          `json:"reason"`
	DataPath datapath.Status `json:"datapath"`
	Sessions struct {
		Active uint32 `json:"active"`
	} `json:"sessions"`
}

type Publisher struct {
	ctx    context.Context
	cancel context.CancelFunc
	period time.Duration
}

func newPublisher(t time.Duration) *Publisher {
	p := &Publisher{}
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.period = t
	return p
}

func NewController(db string, br pkg.BrdigeConfig, rc client.RemoteController, period time.Duration, nodeId string, debug bool) (*Controller, error) {
	c := &Controller{}

	store, err := store.NewStore(db)
	if err != nil {
		log.Errorf("Failed to create db: %v", err)

		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	sm, err := session.NewSessionManager(store, br)
	if err != nil {
		log.Errorf("Failed to create session manager: %v", err)

		return nil, fmt.Errorf("failed to create session manager: %w", err)
	}

	c.nodeId = nodeId
	c.rc = rc
	c.sm = sm
	c.store = store
	c.publisher = newPublisher(period)

	c.startPublisher()

	return c, nil
}

func (c *Controller) Status() ControllerStatus {
	smStatus := c.sm.Status()

	status := ControllerStatus{
		NodeID:   c.nodeId,
		Ready:    smStatus.DataPath.Connected,
		State:    "ready",
		Reason:   "none",
		DataPath: smStatus.DataPath,
	}

	if !smStatus.DataPath.Connected {
		status.Ready = false
		status.State = "failed"
		status.Reason = "datapath_not_connected"
	}

	status.Sessions.Active = smStatus.DataPath.UECount

	return status
}

func sessionResponse(s *store.Session) *api.SessionResponse {
	return &api.SessionResponse{
		ID:         s.ID,
		NodeId:     s.NodeId,
		Imsi:       s.SubscriberID.Imsi,
		PolicyID:   s.PolicyID.ID.String(),
		ApnName:    s.ApnName,
		UeIpaddr:   s.UeIpAddr,
		StartTime:  s.StartTime,
		EndTime:    s.EndTime,
		TxBytes:    s.TxBytes,
		RxBytes:    s.RxBytes,
		TotalBytes: s.TotalBytes,
		TxMeterId:  uint32(s.TxMeterID.ID),
		RxMeterId:  uint32(s.RxMeterID.ID),
		State:      s.State.String(),
		Sync:       s.Sync.String(),
		UpdatedAt:  s.UpdatedAt,
	}
}

func policyResponse(p *store.Policy) *api.PolicyResponse {
	return &api.PolicyResponse{
		ID:        p.ID,
		Burst:     p.Burst,
		Data:      p.Data,
		Consumed:  p.Consumed,
		Dlbr:      p.Dlbr,
		Ulbr:      p.Ulbr,
		StartTime: p.StartTime,
		EndTime:   p.EndTime,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func flowResponse(flows []*store.Flow) []*api.FlowResponse {
	fr := make([]*api.FlowResponse, len(flows))
	for i, flow := range flows {
		fr[i] = &api.FlowResponse{
			ID:        flow.ID,
			Cookie:    flow.Cookie,
			Table:     flow.Tableid,
			Priority:  flow.Priority,
			UeIpaddr:  flow.UeIpAddr,
			ReRouting: flow.ReRouting.IpAddr,
			MeterID:   uint32(flow.MeterID.ID),
		}
	}
	return fr
}

func reRouteResponse(route *store.ReRoute) *api.ReRouteResponse {
	return &api.ReRouteResponse{
		ID: route.ID,
		Ip: route.IpAddr,
	}
}

func subscriberResponse(s *store.Subscriber) *api.SubscriberResponse {
	return &api.SubscriberResponse{
		ID:       s.ID,
		Imsi:     s.Imsi,
		PolicyID: s.PolicyID.ID,
		ReRoute:  s.ReRouteID.IpAddr,
	}
}

func (c *Controller) ExitController() error {
	err := c.sm.EndAllSessions()
	if err != nil {
		log.Errorf("Failed to end all sessions.Error: %v", err)
	}

	return c.stopPublisher()
}

func (c *Controller) validateSubscriber(imsi string) (*store.Subscriber, error) {
	s, err := c.store.GetSubscriber(imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber for %s.Error: %v", imsi, err)

		return nil, fmt.Errorf("failed to get subscriber for %s.Error: %w", imsi, err)
	}

	now := time.Now().Unix()
	if s.PolicyID.StartTime > now || s.PolicyID.EndTime <= now {
		return nil, fmt.Errorf("failed to get valid policy")
	}

	err = c.store.ValidateDataCapLimits(imsi, &s.PolicyID)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Controller) CreateSession(ctx *gin.Context, req *api.CreateSession) error {
	var sub *store.Subscriber
	var err error

	log.Infof("New session request received for subscriber %s and Ip address %s",
		req.ImsiStr, req.IpStr)

	sub, err = c.validateSubscriber(req.ImsiStr)
	if err != nil {
		log.Errorf("Subscriber %s not configured locally or policy invalid: %v",
			req.ImsiStr, err)

		return fmt.Errorf("subscriber %s not configured locally or policy invalid: %w",
			req.ImsiStr, err)
	}

	state := c.sm.IfSessionExist(ctx, sub.Imsi, req.IpStr)
	if state {
		log.Errorf("Session already exists for %s user with ip %s terminating it.",
			sub.Imsi, req.IpStr)

		err = c.sm.EndSession(ctx, sub)
		if err != nil {
			log.Errorf("Failed to end session on bridge for subscriber %s.Error: %v",
				req.ImsiStr, err)

			return fmt.Errorf("failed to end session on bridge for subscriber %s.Error: %w",
				req.ImsiStr, err)
		}

		return nil
	}

	s, rxF, txF, err := c.store.CreateSession(sub, req.IpStr, c.nodeId)
	if err != nil {
		log.Errorf("Failed to create a session for subscriber %s.Error: %v",
			req.ImsiStr, err)

		return fmt.Errorf("failed to create a session for subscriber %s.Error: %w",
			req.ImsiStr, err)
	}

	err = c.sm.CreateSesssion(ctx, sub, s, rxF, txF)
	if err != nil {
		log.Errorf("Failed to monitor session on bridge for subscriber %s.Error: %v",
			req.ImsiStr, err)

		return fmt.Errorf("failed to monitor session on bridge for subscriber %s.Error: %w",
			req.ImsiStr, err)
	}

	return nil
}

func (c *Controller) EndSession(ctx *gin.Context, req *api.EndSession) error {
	sub, err := c.store.GetSubscriber(req.ImsiStr)
	if err != nil {
		log.Errorf("Failed to get subscriber for imsi %s.Error: %v", req.ImsiStr, err)

		return fmt.Errorf("failed to get subscriber for imsi %s.Error: %w", req.ImsiStr, err)
	}

	err = c.sm.EndSession(ctx, sub)
	if err != nil {
		log.Errorf("Failed to end session on bridge for subscriber %s.Error: %v",
			req.ImsiStr, err)

		return fmt.Errorf("failed to end session on bridge for subscriber %s.Error: %w",
			req.ImsiStr, err)
	}

	return nil
}

func (c *Controller) GetSessionByID(ctx *gin.Context, req *api.GetSessionByID) (*api.SessionResponse, error) {
	s, err := c.store.GetSessionByID(int(req.ID))
	if err != nil {
		log.Errorf("Failed to get session with id %d.Error: %v", req.ID, err)

		return nil, fmt.Errorf("failed to get session with id %d.Error: %w", req.ID, err)
	}

	return sessionResponse(s), nil
}

func (c *Controller) GetActiveSessionByImsi(ctx *gin.Context, req *api.GetSessionByImsi) (*api.SessionResponse, error) {
	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get active session for Imsi %s.Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get active session for Imsi %s.Error: %w", req.Imsi, err)
	}
	return sessionResponse(s), nil
}

func (c *Controller) GetCDRBySessionId(ctx *gin.Context, req *api.GetCDRBySessionId) (*api.CDR, error) {
	s, err := c.store.GetSessionByID(int(req.ID))
	if err != nil {
		log.Errorf("Failed to get session with id %d.Error: %v", req.ID, err)

		return nil, fmt.Errorf("failed to get session with id %d.Error: %w", req.ID, err)
	}

	cdr := store.PrepareCDR(s)

	return cdr, nil
}

func (c *Controller) GetCDRByImsi(ctx *gin.Context, req *api.GetCDRByImsi) ([]*api.CDR, error) {
	sess, err := c.store.GetSessionsByImsi(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get session for Imsi %s.Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get session for Imsi %s.Error: %w", req.Imsi, err)
	}

	cdrs := make([]*api.CDR, len(sess))
	for i, s := range sess {
		cdrs[i] = store.PrepareCDR(&s)
	}

	return cdrs, nil
}

func (c *Controller) GetPolicyByImsi(ctx *gin.Context, req *api.GetPolicyByImsi) (*api.PolicyResponse, error) {
	p, err := c.store.GetApplicablePolicyByImsi(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get policy for Imsi %s.Error: %v", req.Imsi, err.Error())

		return nil, fmt.Errorf("failed to get policy for Imsi %s.Error: %w", req.Imsi, err)
	}

	return policyResponse(p), nil
}

func (c *Controller) GetPolicyByID(ctx *gin.Context, req *api.GetPolicyByID) (*api.PolicyResponse, error) {
	id, err := uuid.FromString(req.ID)
	if err != nil {
		log.Errorf("Invalid policy id.Error: %v", err.Error())

		return nil, fmt.Errorf("invalid policy id.Error: %w", err)
	}

	p, err := c.store.GetPolicyByID(id)
	if err != nil {
		log.Errorf("Failed to get policy with ID %s.Error: %v", req.ID, err.Error())

		return nil, fmt.Errorf("failed to get policy with ID %s.Error: %w", req.ID, err)

	}
	return policyResponse(p), nil
}

func (c *Controller) AddPolicy(ctx *gin.Context, req *api.Policy) error {
	_, err := c.store.CreatePolicy(req)
	if err != nil {
		log.Errorf("Failed to add policy %s. Error: %s", req.Uuid.String(), err.Error())

		return fmt.Errorf("failed to add policy %s. Error: %w", req.Uuid.String(), err)
	}

	return nil
}

func (c *Controller) GetFlowsForImsi(ctx *gin.Context, req *api.GetFlowsForImsi) ([]*api.FlowResponse, error) {
	var flows []*store.Flow

	_, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber with Imsi %s.Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get subscriber with Imsi %s.Error: %w", req.Imsi, err)
	}

	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get active session for Imsi %s.Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get active session for Imsi %s.Error: %w", req.Imsi, err)
	}

	fRx, err := c.store.GetFlowForMeter(s.RxMeterID.ID)
	if err != nil {
		log.Errorf("Failed to get RX flow for Imsi %s.Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get RX flow for Imsi %s.Error: %w", req.Imsi, err)
	}
	flows = append(flows, fRx)

	fTx, err := c.store.GetFlowForMeter(s.TxMeterID.ID)
	if err != nil {
		log.Errorf("Failed to get TX flow for Imsi %s.Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get TX flow for Imsi %s.Error: %w", req.Imsi, err)
	}
	flows = append(flows, fTx)

	return flowResponse(flows), nil
}

func (c *Controller) GetReroute(ctx *gin.Context, req *api.GetReRouteByImsi) (*api.ReRouteResponse, error) {
	_, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber with Imsi %s.Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get subscriber with Imsi %s.Error: %w", req.Imsi, err)
	}

	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get active session for Imsi %s.Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get active session for Imsi %s.Error: %w", req.Imsi, err)
	}

	flow, err := c.store.GetFlowForMeter(s.TxMeterID.ID)
	if err != nil {
		log.Errorf("failed to get TX flow for Imsi %s.Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get TX flow for Imsi %s.Error: %w", req.Imsi, err)
	}

	r, err := c.store.GetReRouteByID(flow.ReRouting.ID)
	if err != nil {
		log.Errorf("Failed to get reroute for imsi %s. Error %s", req.Imsi, err.Error())

		return nil, fmt.Errorf("failed to get reroute for imsi %s. Error %w", req.Imsi, err)
	}

	return reRouteResponse(r), nil
}

func (c *Controller) UpdateReroute(ctx *gin.Context, req *api.UpdateRerouteById) error {
	err := c.store.UpdateReroute(&store.ReRoute{
		ID:     int(req.Id),
		IpAddr: req.Ip,
	})
	if err != nil {
		log.Errorf("Failed to update route for Id %d. Error: %s", req.Id, err.Error())

		return fmt.Errorf("failed to update route for Id %d. Error: %w", req.Id, err)
	}
	return nil
}

func (c *Controller) GetSubscriber(ctx *gin.Context, req *api.RequestSubscriber) (*api.SubscriberResponse, error) {
	s, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber with imsi %s.Error: %s", req.Imsi, err.Error())

		return nil, fmt.Errorf("failed to get subscriber with imsi %s. Error: %w", req.Imsi, err)
	}
	return subscriberResponse(s), nil
}

func (c *Controller) DeleteSubscriber(ctx *gin.Context, req *api.RequestSubscriber) error {
	s, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber with imsi %s.Error: %s", req.Imsi, err.Error())

		return fmt.Errorf("failed to get subscriber with imsi %s. Error: %w", req.Imsi, err)
	}

	err = c.store.DeleteSubscriber(s)
	if err != nil {
		log.Errorf("Failed to delete subscriber with imsi %s.Error: %s", req.Imsi, err.Error())

		return fmt.Errorf("failed to delete subscriber with imsi %s.Error: %w", req.Imsi, err)
	}
	return nil
}

func (c *Controller) AddSubscriber(ctx *gin.Context, req *api.CreateSubscriber) error {
	_, err := c.store.CreateSubscriber(req.Imsi, &req.Policy, &req.ReRoute, nil)
	if err != nil {
		log.Errorf("Failed to create subscriber with imsi %s. Error: %s", req.Imsi, err.Error())

		return fmt.Errorf("failed to create subscriber with imsi %s. Error: %w", req.Imsi, err)
	}
	return nil
}

func (c *Controller) UpdateSubscriber(ctx *gin.Context, req *api.UpdateSubscriber) error {
	_, err := c.store.UpdateSubscriber(req.Imsi, &req.Policy)
	if err != nil {
		log.Errorf("Failed to update subscriber with imsi %s. Error: %s", req.Imsi, err.Error())

		return fmt.Errorf("failed to update subscriber with imsi %s. Error: %w", req.Imsi, err)
	}
	return nil
}

func handlePendingSyncSession(c *Controller) {
	sessions, err := c.store.GetAllNonPublishedSessions()
	if err != nil {
		log.Errorf("[Publisher] Failed to get unpublished sessions from store.Error %s", err.Error())
		return
	}

	for _, session := range sessions {
		log.Infof("[Publisher] Session %d for subscriber %s is ready for backend sync.",
			session.ID, session.SubscriberID.Imsi)

		err = c.publishCDRToRemoteController(session)
		if err != nil {
			log.Warnf("error while pushing CDR to remote backend controller: %v", err)
		}
	}
}

func handleTerminatedSession(c *Controller) {
	sessions, err := c.store.GetAllNonPublishedTerminatedSessions()
	if err != nil {
		log.Errorf("[Publisher] Failed to get non-published terminated sessions from store.Error %s", err.Error())
		return
	}

	for _, session := range sessions {
		err = c.store.UpdateSessionSyncState(session.ID, store.SessionSyncReady)
		if err != nil {
			log.Errorf("[Publisher] Failed to update session %d for subscriber %s in store.Error %s", session.ID, session.SubscriberID.Imsi, err.Error())
			return
		}

		log.Infof("[Publisher] Session %d for subscriber %s updated to sync ready state.",
			session.ID, session.SubscriberID.Imsi)
	}
}

func (c *Controller) publishCDR() {
	ticker := time.NewTicker(c.publisher.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			handlePendingSyncSession(c)
			handleTerminatedSession(c)

		case <-c.publisher.ctx.Done():
			log.Infof("[Publisher] Ending routine to publish CDRs")
			return
		}
	}
}

func (c *Controller) publishCDRToRemoteController(session store.Session) error {
	cdr := &api.CDR{
		Session:       session.ID,
		NodeId:        session.NodeId,
		Imsi:          session.SubscriberID.Imsi,
		Policy:        session.PolicyID.ID.String(),
		ApnName:       session.ApnName,
		Ip:            session.UeIpAddr,
		StartTime:     session.StartTime,
		EndTime:       session.EndTime,
		LastUpdatedAt: session.UpdatedAt,
		TxBytes:       session.TxBytes,
		RxBytes:       session.RxBytes,
		TotalBytes:    session.TotalBytes,
	}

	return c.rc.PushCdr(cdr)
}

func (c *Controller) startPublisher() {
	log.Infof("Starting publisher routine")
	go c.publishCDR()
}

func (c *Controller) stopPublisher() error {
	log.Infof("Stopping publisher routine.")
	c.publisher.cancel()
	time.Sleep(100 * time.Millisecond)
	return nil
}
