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
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/session"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/systems/common/uuid"
)

type Controller struct {
	store     *store.Store
	sm        session.SessionManager
	rc        client.RemoteController
	publisher *Publisher
	nodeId    string
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

func NewController(db string, br pkg.BrdigeConfig, remote string, period time.Duration, nodeId string, debug bool) (*Controller, error) {
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
	c.nodeId = nodeId
	c.sm = session.NewSessionManager(rc, store, br)
	c.store = store
	c.rc = rc
	c.publisher = newPublisher(period)

	c.startPublisher()

	return c, nil
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
		log.Errorf("failed to end all sessions.Error: %v", err)
	}

	return c.stopPublisher()

}

func (c *Controller) validateSubscriber(imsi string) (*store.Subscriber, error) {

	/* Get subscriber policy by imsi*/
	s, err := c.store.GetSubscriber(imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber for %s.Error: %v", imsi, err)
		return nil, err
	}

	now := time.Now().Unix()
	if s.PolicyID.StartTime > now || s.PolicyID.EndTime <= now {
		return nil, fmt.Errorf("failed to get valid policy")
	}

	/* TODO: Also get the usage from the remote cloud and add all the CDR usage after the last update time to it.
	and then compare the value to allowed data cap if it is less than allowed data cap let the user establish session otherwise not.
	*/
	err = c.store.ValidateDataCapLimits(imsi, &s.PolicyID)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *Controller) updateSubscriberProfile(imsi string, p *api.Policy, ip string, d *api.UsageDetails) (*store.Subscriber, error) {
	var sub *store.Subscriber

	sub, err := c.store.CreateSubscriber(imsi, p, &ip, d)
	if err != nil {
		log.Errorf("Failed to create subscriber %s.Error: %v", imsi, err)
		return nil, err
	}

	return sub, nil
}

func (c *Controller) CreateSession(ctx *gin.Context, req *api.CreateSession) error {
	var sub *store.Subscriber
	var err error
	log.Infof("New session request received for subscriber %s and Ip address %s", req.ImsiStr, req.IpStr)
	/* validate subscriber*/
	/* TODO: Validate subscriber should always get the values from the remote  server
	just to make sure the usage values are correct or we could have some timeouts
	if the new session is started like with in 60 secs then we can consider same values
	just to avoid makking duplicate requests
	*/
	sub, err = c.validateSubscriber(req.ImsiStr)
	if err != nil {
		/* Get subscriber policy from remote */
		spr, err := c.rc.GetSubscriberProfile(req.ImsiStr)
		if err != nil {
			log.Errorf("Failed to get subscriber %s policy.Error: %v", req.ImsiStr, err)
			return err
		}

		sub, err = c.updateSubscriberProfile(req.ImsiStr, &spr.Policy, spr.ReRoute, &spr.Usage)
		if err != nil {
			log.Errorf("Failed to update subscriber %s with policy %s.Error: %v", req.ImsiStr, spr.Policy.Uuid.String(), err)
			return err
		}
	}

	/* Check if session already exist wiath same ip */
	state := c.sm.IfSessionExist(ctx, sub.Imsi, req.IpStr)
	if state {
		/* We need to terminate session. may be policy change has happened which updated the data rates etc
		due to which we might need to reset meters */
		log.Errorf("Session already exist for %s user with ip %s terminating it.", sub.Imsi, req.IpStr)
		err = c.sm.EndSession(ctx, sub)
		if err != nil {
			log.Errorf("Failed to end session on bridge for subscriber %s.Error: %v", req.ImsiStr, err)
			return err
		}
		return nil
	}

	/* create session */
	s, rxF, txF, err := c.store.CreateSession(sub, req.IpStr, c.nodeId)
	if err != nil {
		log.Errorf("Failed to create a session for subscriber %s.Error: %v", req.ImsiStr, err)
		return err
	}

	/* create UE data path and monitoring session */
	err = c.sm.CreateSesssion(ctx, sub, s, rxF, txF)
	if err != nil {
		log.Errorf("Failed to monitor session on bridge for subscriber %s.Error: %v", req.ImsiStr, err)
		return err
	}

	return nil
}

func (c *Controller) EndSession(ctx *gin.Context, req *api.EndSession) error {

	sub, err := c.store.GetSubscriber(req.ImsiStr)
	if err != nil {
		log.Errorf("failed to get subscriber for imsi %s.Error: %v", req.ImsiStr, err)
		return err
	}

	err = c.sm.EndSession(ctx, sub)
	if err != nil {
		log.Errorf("Failed to end session on bridge for subscriber %s.Error: %v", req.ImsiStr, err)
		return err
	}
	return nil
}

func (c *Controller) GetSessionByID(ctx *gin.Context, req *api.GetSessionByID) (*api.SessionResponse, error) {
	s, err := c.store.GetSessionByID(int(req.ID))
	if err != nil {
		log.Errorf("failed to get session with id %d.Error: %v", req.ID, err)
		return nil, err
	}
	return sessionResponse(s), nil
}

func (c *Controller) GetActiveSessionByImsi(ctx *gin.Context, req *api.GetSessionByImsi) (*api.SessionResponse, error) {
	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("failed to get active session for Imsi %s.Error: %v", req.Imsi, err)
		return nil, err
	}
	return sessionResponse(s), nil
}

func (c *Controller) GetCDRBySessionId(ctx *gin.Context, req *api.GetCDRBySessionId) (*api.CDR, error) {
	s, err := c.store.GetSessionByID(int(req.ID))
	if err != nil {
		log.Errorf("failed to get session with id %d.Error: %v", req.ID, err)
		return nil, err
	}

	cdr := store.PrepareCDR(s)

	return cdr, nil
}

func (c *Controller) GetCDRByImsi(ctx *gin.Context, req *api.GetCDRByImsi) ([]*api.CDR, error) {

	sess, err := c.store.GetSessionsByImsi(req.Imsi)
	if err != nil {
		log.Errorf("failed to get session for Imsi %s.Error: %v", req.Imsi, err)
		return nil, err
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
		log.Errorf("failed to get policy for Imsi %s.Error: %v", req.Imsi, err.Error())
		return nil, err
	}
	return policyResponse(p), nil
}

func (c *Controller) GetPolicyByID(ctx *gin.Context, req *api.GetPolicyByID) (*api.PolicyResponse, error) {
	id, err := uuid.FromString(req.ID)
	if err != nil {
		log.Errorf("invalid policy id.Error: %v", err.Error())
		return nil, err
	}

	p, err := c.store.GetPolicyByID(id)
	if err != nil {
		log.Errorf("failed to get policy with ID %s.Error: %v", req.ID, err.Error())
		return nil, err
	}
	return policyResponse(p), nil
}

func (c *Controller) AddPolicy(ctx *gin.Context, req *api.Policy) error {
	_, err := c.store.CreatePolicy(req)
	if err != nil {
		log.Errorf("failed to add policy %s.Error: %s", req.Uuid.String(), err.Error())
		return err
	}

	return nil
}

func (c *Controller) GetFlowsForImsi(ctx *gin.Context, req *api.GetFlowsForImsi) ([]*api.FlowResponse, error) {
	var flows []*store.Flow
	_, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber with Imsi %s.Error: %v", req.Imsi, err)
		return nil, err
	}

	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("failed to get active session for Imsi %s.Error: %v", req.Imsi, err)
		return nil, err
	}

	fRx, err := c.store.GetFlowForMeter(s.RxMeterID.ID)
	if err != nil {
		log.Errorf("failed to get RX flow for Imsi %s.Error: %v", req.Imsi, err)
		return nil, err
	}
	flows = append(flows, fRx)

	fTx, err := c.store.GetFlowForMeter(s.TxMeterID.ID)
	if err != nil {
		log.Errorf("failed to get TX flow for Imsi %s.Error: %v", req.Imsi, err)
		return nil, err
	}
	flows = append(flows, fTx)

	return flowResponse(flows), nil
}

func (c *Controller) GetReroute(ctx *gin.Context, req *api.GetReRouteByImsi) (*api.ReRouteResponse, error) {

	_, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber with Imsi %s.Error: %v", req.Imsi, err)
		return nil, err
	}

	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("failed to get active session for Imsi %s.Error: %v", req.Imsi, err)
		return nil, err
	}

	flow, err := c.store.GetFlowForMeter(s.TxMeterID.ID)
	if err != nil {
		log.Errorf("failed to get TX flow for Imsi %s.Error: %v", req.Imsi, err)
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
		IpAddr: req.Ip,
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
		log.Errorf("failed to get subscriber with imsi %s.Error: %s", req.Imsi, err.Error())
		return nil, err
	}
	return subscriberResponse(s), nil
}

func (c *Controller) DeleteSubscriber(ctx *gin.Context, req *api.RequestSubscriber) error {
	s, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("failed to get subscriber with imsi %s.Error: %s", req.Imsi, err.Error())
		return err
	}

	err = c.store.DeleteSubscriber(s)
	if err != nil {
		log.Errorf("failed to delete subscriber with imsi %s.Error: %s", req.Imsi, err.Error())
		return err
	}
	return nil
}

func (c *Controller) AddSubscriber(ctx *gin.Context, req *api.CreateSubscriber) error {
	_, err := c.store.CreateSubscriber(req.Imsi, &req.Policy, &req.ReRoute, nil)
	if err != nil {
		log.Errorf("failed to delete subscriber with imsi %s. Error: %s", req.Imsi, err.Error())
		return err
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
		cdr := store.PrepareCDR(&session)
		err := c.rc.PushCdr(cdr)
		if err != nil {
			log.Errorf("[Publisher] Failed to publish session %+v for subscriber %s from store.Error %s", session, session.SubscriberID.Imsi, err.Error())
			continue
		}

		/* Update store if published successfully */
		err = c.store.UpdateSessionSyncState(session.ID, store.SessionSyncCompleted)
		if err != nil {
			log.Errorf("[Publisher] Failed to update session %+v for subscriber %s in store.Error %s", session, session.SubscriberID.Imsi, err.Error())
			continue
		}

		log.Infof("[Publisher] Published CDR for session %+v for subscriber %s from store. CDR data %+v", session, session.SubscriberID.Imsi, cdr)
	}
}

func handleTerminatedSession(c *Controller) {
	sessions, err := c.store.GetAllNonPublishedTerminatedSessions()
	if err != nil {
		log.Errorf("[Publisher] Failed to get non-published terminated sessions from store.Error %s", err.Error())
		return
	}
	for _, session := range sessions {

		/* Update store. Mark message as ready for sync */
		err = c.store.UpdateSessionSyncState(session.ID, store.SessionSyncReady)
		if err != nil {
			log.Errorf("[Publisher] Failed to update session %d for subscriber %s in store.Error %s", session.ID, session.SubscriberID.Imsi, err.Error())
			return
		}

		log.Infof("[Publisher] Session %d for subscriber %s from store updated to sync ready state.", session.ID, session.SubscriberID.Imsi)
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
			log.Infof("[Publisher] Ending routine to pusblish CDR's")
			return
		}
	}
}

func (c *Controller) startPublisher() {
	log.Infof("Starting publisher routine")
	go c.publishCDR()
}

func (c *Controller) stopPublisher() error {
	log.Infof("Stoping publisher routine.")
	c.publisher.cancel()
	time.Sleep(100 * time.Millisecond)
	return nil
}
