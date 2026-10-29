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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ukama/ukama/nodes/apps/pcrf/pkg"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/session"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/datapath"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

const (
	defaultApnName = "internet"
)

type Controller struct {
	pdnMu           sync.Mutex
	apn             string
	pdnCancel       context.CancelFunc
	pdnDone         chan struct{}
	restartRequired bool
	stopping        bool
	store           *store.Store
	sm              session.SessionManager
	rc              client.RemoteController
	publisher       *Publisher
	nodeId          string
	serviceMu       sync.RWMutex
	serviceOn       bool
	serviceReason   string
}

type ControllerStatus struct {
	NodeID   string              `json:"nodeId"`
	Ready    bool                `json:"ready"`
	State    string              `json:"state"`
	Reason   string              `json:"reason"`
	Service  api.ServiceResponse `json:"service"`
	DataPath datapath.Status     `json:"datapath"`
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

func NewController(db string, br pkg.BrdigeConfig, rc client.RemoteController,
	period time.Duration, nodeId string, debug bool, apn string) (
	*Controller, error) {
	selectedAPN, err := api.NormalizeAPN(apn)
	if err != nil {
		return nil, fmt.Errorf("configured APN: %w", err)
	}
	c := &Controller{apn: selectedAPN}

	store, err := store.NewStore(db)
	if err != nil {
		log.Errorf("Failed to create db: %v", err)

		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	if err := store.RetireEPCs(); err != nil {
		return nil, err
	}
	sm, err := session.NewSessionManager(store, br)
	if err != nil {
		log.Errorf("Failed to create session manager: %v", err)

		return nil, fmt.Errorf("failed to create session manager: %w", err)
	}

	c.nodeId = nodeId
	c.serviceOn = false
	c.serviceReason = "startup_diable"
	c.rc = rc
	c.sm = sm
	c.store = store
	c.publisher = newPublisher(period)

	if err := c.cleanupPDNsLocked(context.Background()); err != nil {
		return nil, err
	}
	c.startPDNMonitor(br.Period)
	c.startPublisher()

	return c, nil
}

func (c *Controller) Status() ControllerStatus {
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	smStatus := c.sm.Status()

	status := ControllerStatus{
		NodeID:   c.nodeId,
		Ready:    smStatus.DataPath.Connected,
		State:    "ready",
		Reason:   "none",
		Service:  c.ServiceStatus(),
		DataPath: smStatus.DataPath,
	}

	if !smStatus.DataPath.Connected {
		status.Ready = false
		status.State = "failed"
		status.Reason = "datapath_not_connected"
	}

	if c.restartRequired || c.stopping ||
		smStatus.DataPath.ConnectedCount > 1 {
		status.Ready = false
		status.State = "failed"
		status.Reason = "restart_required"
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
		ID:        s.ID,
		Imsi:      s.Imsi,
		PolicyID:  s.PolicyID.ID,
		ReRoute:   s.ReRouteID.IpAddr,
		ServiceOn: s.ServiceOn,
	}
}

func (c *Controller) ExitController() error {
	if c.pdnCancel != nil {
		c.pdnCancel()
		<-c.pdnDone
	}
	c.pdnMu.Lock()
	c.stopping = true
	err := c.store.RetireEPCs()
	if err == nil {
		err = c.cleanupPDNsLocked(context.Background())
	}
	c.pdnMu.Unlock()
	if c.publisher != nil {
		_ = c.stopPublisher()
	}
	return err
}

func (c *Controller) ServiceStatus() api.ServiceResponse {
	c.serviceMu.RLock()
	defer c.serviceMu.RUnlock()

	state := "off"
	admission := "disabled"
	reason := c.serviceReason
	if c.serviceOn {
		state = "on"
		admission = "enabled"
	}
	if reason == "" {
		reason = "ok"
	}

	return api.ServiceResponse{
		State:     state,
		Admission: admission,
		Reason:    reason,
	}
}

func (c *Controller) SetService(ctx context.Context,
	req *api.ServiceRequest) error {
	state := strings.ToLower(strings.TrimSpace(req.State))
	if state != "on" && state != "off" {
		return rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "invalid service state"}
	}
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	c.serviceMu.Lock()
	c.serviceOn = state == "on"
	c.serviceReason = "ok"
	if !c.serviceOn {
		c.serviceReason = "operator_disabled"
	}
	c.serviceMu.Unlock()
	return c.reconcilePDNsLocked(ctx)
}

func (c *Controller) serviceEnabled() bool {
	c.serviceMu.RLock()
	defer c.serviceMu.RUnlock()
	return c.serviceOn
}

func (c *Controller) GetSessionByID(ctx *gin.Context, req *api.GetSessionByID) (*api.SessionResponse, error) {
	s, err := c.store.GetSessionByID(int(req.ID))
	if err != nil {
		log.Errorf("Failed to get session with id %d. Error: %v", req.ID, err)

		return nil, fmt.Errorf("failed to get session with id %d. Error: %w", req.ID, err)
	}

	return sessionResponse(s), nil
}

func (c *Controller) GetActiveSessionByImsi(ctx *gin.Context, req *api.GetSessionByImsi) (*api.SessionResponse, error) {
	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get active session for Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get active session for Imsi %s. Error: %w", req.Imsi, err)
	}
	return sessionResponse(s), nil
}

func (c *Controller) GetCDRBySessionId(ctx *gin.Context, req *api.GetCDRBySessionId) (*api.CDR, error) {
	s, err := c.store.GetSessionByID(int(req.ID))
	if err != nil {
		log.Errorf("Failed to get session with id %d. Error: %v", req.ID, err)

		return nil, fmt.Errorf("failed to get session with id %d. Error: %w", req.ID, err)
	}

	cdr := store.PrepareCDR(s)

	return cdr, nil
}

func (c *Controller) GetCDRByImsi(ctx *gin.Context, req *api.GetCDRByImsi) ([]*api.CDR, error) {
	sess, err := c.store.GetSessionsByImsi(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get session for Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get session for Imsi %s. Error: %w", req.Imsi, err)
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
		log.Errorf("Failed to get policy for Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get policy for Imsi %s. Error: %w", req.Imsi, err)
	}

	return policyResponse(p), nil
}

func (c *Controller) GetPolicyByID(ctx *gin.Context, req *api.GetPolicyByID) (*api.PolicyResponse, error) {
	id, err := uuid.FromString(req.ID)
	if err != nil {
		log.Errorf("Invalid policy id. Error: %v", err)

		return nil, rest.HttpError{HttpCode: http.StatusBadRequest, Message: "invalid policy id"}
	}

	p, err := c.store.GetPolicyByID(id)
	if err != nil {
		log.Errorf("Failed to get policy with ID %s. Error: %v", req.ID, err)

		return nil, fmt.Errorf("failed to get policy with ID %s. Error: %w", req.ID, err)

	}
	return policyResponse(p), nil
}

func (c *Controller) AddPolicy(ctx *gin.Context, req *api.Policy) error {
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	_, err := c.store.CreatePolicy(req)
	if err != nil {
		log.Errorf("Failed to add policy %s. Error: %v", req.Uuid.String(), err)

		return fmt.Errorf("failed to add policy %s. Error: %w", req.Uuid.String(), err)
	}

	return nil
}

func (c *Controller) GetFlowsForImsi(ctx *gin.Context, req *api.GetFlowsForImsi) ([]*api.FlowResponse, error) {
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	p, err := c.store.GetLivePDN(req.Imsi)
	if err != nil || p.State != store.PDNLive {
		if err != nil && !errors.Is(err, store.ErrPDNNotFound) {
			return nil, err
		}
		return nil, rest.HttpError{HttpCode: http.StatusNotFound,
			Message: "live PDN not found"}
	}
	var flows []*store.Flow

	_, err = c.store.GetSubscriber(req.Imsi)
	if err != nil {
		if errors.Is(err, store.ErrSubscriberNotFound) {
			log.Infof("No subscriber for Imsi %s; reporting flow as withdrawn", req.Imsi)

			return nil, rest.HttpError{HttpCode: http.StatusNotFound, Message: "subscriber not found"}
		}

		log.Errorf("Failed to get subscriber with Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get subscriber with Imsi %s. Error: %w", req.Imsi, err)
	}

	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		if errors.Is(err, store.ErrSessionNotFound) {
			log.Infof("No active session for Imsi %s; reporting flow as withdrawn", req.Imsi)

			return nil, rest.HttpError{HttpCode: http.StatusNotFound, Message: "active session not found"}
		}

		log.Errorf("Failed to get active session for Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get active session for Imsi %s. Error: %w", req.Imsi, err)
	}

	if s.FlowState == store.FlowsPaused {
		log.Infof("Session for Imsi %s is paused; reporting flow as withdrawn", req.Imsi)

		return nil, rest.HttpError{HttpCode: http.StatusNotFound, Message: "flow withdrawn (service paused)"}
	}

	if s.FlowState != store.FlowsActive {
		return nil, rest.HttpError{HttpCode: http.StatusServiceUnavailable,
			Message: "flow_update_pending"}
	}
	fRx, err := c.store.GetFlowForMeter(s.RxMeterID.ID)
	if err != nil {
		log.Errorf("Failed to get RX flow for Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get RX flow for Imsi %s. Error: %w", req.Imsi, err)
	}
	flows = append(flows, fRx)

	fTx, err := c.store.GetFlowForMeter(s.TxMeterID.ID)
	if err != nil {
		log.Errorf("Failed to get TX flow for Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get TX flow for Imsi %s. Error: %w", req.Imsi, err)
	}
	flows = append(flows, fTx)

	return flowResponse(flows), nil
}

func (c *Controller) GetReroute(ctx *gin.Context, req *api.GetReRouteByImsi) (*api.ReRouteResponse, error) {
	_, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber with Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get subscriber with Imsi %s. Error: %w", req.Imsi, err)
	}

	s, err := c.store.GetActiveSessionByImsi(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get active session for Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get active session for Imsi %s. Error: %w", req.Imsi, err)
	}

	flow, err := c.store.GetFlowForMeter(s.TxMeterID.ID)
	if err != nil {
		log.Errorf("failed to get TX flow for Imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get TX flow for Imsi %s. Error: %w", req.Imsi, err)
	}

	r, err := c.store.GetReRouteByID(flow.ReRouting.ID)
	if err != nil {
		log.Errorf("Failed to get reroute for imsi %s. Error %v", req.Imsi, err)

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
		log.Errorf("Failed to update route for Id %d. Error: %v", req.Id, err)

		return fmt.Errorf("failed to update route for Id %d. Error: %w", req.Id, err)
	}
	return nil
}

func (c *Controller) GetSubscriber(ctx *gin.Context, req *api.RequestSubscriber) (*api.SubscriberResponse, error) {
	s, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber with imsi %s. Error: %v", req.Imsi, err)

		return nil, fmt.Errorf("failed to get subscriber with imsi %s. Error: %w", req.Imsi, err)
	}
	return subscriberResponse(s), nil
}

func (c *Controller) DeleteSubscriber(ctx *gin.Context,
	req *api.RequestSubscriber) error {
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	sub, err := c.store.GetSubscriber(req.Imsi)
	if err != nil {
		return err
	}
	// Persist the block before cleanup so a failed removal remains retryable.
	if err := c.store.UpdateSubscriberServiceStatus(sub, false); err != nil {
		return err
	}
	if err := c.sm.EndSession(ctx, sub); err != nil {
		return err
	}
	return c.store.DeleteSubscriber(sub)
}

func (c *Controller) AddSubscriber(ctx *gin.Context,
	req *api.CreateSubscriber) error {
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	sub, err := c.store.CreateSubscriber(req.Imsi, &req.Policy,
		&req.ReRoute, nil)
	if err != nil {
		return err
	}
	if req.IsServiceOn != nil {
		if err := c.store.UpdateSubscriberServiceStatus(sub,
			*req.IsServiceOn); err != nil {
			return err
		}
	}
	return c.reconcileImsiLocked(ctx, req.Imsi)
}

func (c *Controller) UpdateSubscriber(ctx *gin.Context,
	req *api.UpdateSubscriber) error {
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	sub, err := c.store.UpdateSubscriber(req.Imsi, &req.Policy)
	if err != nil {
		return err
	}
	if req.IsServiceOn != nil {
		if err := c.store.UpdateSubscriberServiceStatus(sub,
			*req.IsServiceOn); err != nil {
			return err
		}
	}
	return c.reconcileImsiLocked(ctx, req.Imsi)
}

func handlePendingSyncSession(c *Controller) {
	sessions, err := c.store.GetAllNonPublishedSessions()
	if err != nil {
		log.Errorf("[Publisher] Failed to get unpublished sessions from store. Error %v", err)
		return
	}

	for _, session := range sessions {
		log.Infof("[Publisher] Session %d for subscriber %s is ready for backend sync.",
			session.ID, session.SubscriberID.Imsi)

		err = c.sendCDRToRemoteController(&session)
		if err != nil {
			log.Warnf("error while pushing CDR to remote backend controller: %v", err)
		}
	}
}

func handleTerminatedSession(c *Controller) {
	sessions, err := c.store.GetAllNonPublishedTerminatedSessions()
	if err != nil {
		log.Errorf("[Publisher] Failed to get non-published terminated sessions from store. Error %v",
			err)

		return
	}

	for _, session := range sessions {
		err = c.store.UpdateSessionSyncState(session.ID, store.SessionSyncReady)
		if err != nil {
			log.Errorf("[Publisher] Failed to update session %d for subscriber %s in store. Error %v",
				session.ID, session.SubscriberID.Imsi, err)

			return
		}

		log.Infof("[Publisher] Session %d for subscriber %s updated to sync ready state.",
			session.ID, session.SubscriberID.Imsi)
	}
}

func safeRun(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("[Publisher] Recovered from panic in %s: %v", name, r)
		}
	}()

	fn()
}

func (c *Controller) publishCDR() {
	ticker := time.NewTicker(c.publisher.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			safeRun("handlePendingSyncSession", func() { handlePendingSyncSession(c) })
			safeRun("handleTerminatedSession", func() { handleTerminatedSession(c) })

		case <-c.publisher.ctx.Done():
			log.Infof("[Publisher] Ending routine to publish CDRs")
			return
		}
	}
}

func (c *Controller) sendCDRToRemoteController(session *store.Session) error {
	cdr := store.PrepareCDR(session)
	if cdr.ApnName == "" {
		cdr.ApnName = defaultApnName
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
