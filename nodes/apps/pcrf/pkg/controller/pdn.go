/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package controller

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

func pdnError(code int, reason string) error {
	return rest.HttpError{HttpCode: code, Message: reason}
}

func validIdentity(id string) bool {
	parsed, err := uuid.FromString(id)
	return err == nil && parsed != uuid.Nil && parsed.String() == id
}

func validIMSI(imsi string) bool {
	if len(imsi) < 5 || len(imsi) > 15 {
		return false
	}
	for _, digit := range imsi {
		if digit < '0' || digit > '9' {
			return false
		}
	}
	return true
}

func bindingResponse(p *store.PDN, data, reason string) *api.PDNResponse {
	return &api.PDNResponse{
		InstanceID: p.InstanceID, PDNID: p.ID, Imsi: p.Imsi,
		ApnName: p.ApnName, IP: p.IP, State: p.State,
		DataState: data, Reason: reason,
	}
}

// RegisterEPC is called once per EPC/EPCEMU process, before any PDN POST.
// A new process retires the old instance and cleans its permissions first.
func (c *Controller) RegisterEPC(ctx *gin.Context,
	req *api.RegisterEPC) (*api.EPCResponse, error) {
	if !validIdentity(req.InstanceID) ||
		(req.ExpectedInstanceID != "" &&
			!validIdentity(req.ExpectedInstanceID)) {
		return nil, pdnError(http.StatusBadRequest, "invalid_identity")
	}
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	if err := c.checkDataPathLocked(); err != nil {
		return nil, err
	}
	err := c.store.BeginEPC(req.InstanceID, req.ExpectedInstanceID)
	if err != nil {
		if errors.Is(err, store.ErrInstanceConflict) {
			return nil, pdnError(http.StatusConflict, err.Error())
		}
		if errors.Is(err, store.ErrInstanceRetired) {
			return nil, pdnError(http.StatusGone, err.Error())
		}
		return nil, err
	}
	if err := c.cleanupPDNsLocked(ctx); err != nil {
		return nil, pdnError(http.StatusServiceUnavailable, "cleanup_pending")
	}
	return &api.EPCResponse{InstanceID: req.InstanceID, State: "ready"}, nil
}

func (c *Controller) CreateSession(ctx *gin.Context,
	req *api.CreateSession) (*api.PDNResponse, error) {
	if !validIdentity(req.InstanceID) || !validIdentity(req.PDNID) {
		return nil, pdnError(http.StatusBadRequest, "invalid_identity")
	}
	if !validIMSI(req.ImsiStr) {
		return nil, pdnError(http.StatusBadRequest, "invalid_imsi")
	}
	ip := net.ParseIP(req.IpStr)
	if ip == nil || ip.To4() == nil || !ip.IsGlobalUnicast() {
		return nil, pdnError(http.StatusBadRequest, "invalid_address")
	}
	apn, err := api.NormalizeAPN(req.ApnName)
	if err != nil {
		return nil, pdnError(http.StatusBadRequest, "invalid_apn")
	}
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	if apn != c.apn {
		return nil, pdnError(http.StatusForbidden, "unsupported_apn")
	}
	if err := c.checkDataPathLocked(); err != nil {
		return nil, err
	}
	active, err := c.store.EPCActive(req.InstanceID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, pdnError(http.StatusGone, "epc_instance_not_active")
	}
	p, err := c.store.GetPDN(req.InstanceID, req.PDNID)
	switch {
	case err == nil:
		if p.Imsi == req.ImsiStr && p.State != store.PDNLive {
			return nil, pdnError(http.StatusGone, "pdn_closed")
		}
		if p.Imsi != req.ImsiStr || p.ApnName != apn || p.IP != ip.String() {
			return nil, pdnError(http.StatusConflict, "pdn_conflict")
		}
		if p.State != store.PDNLive {
			return nil, pdnError(http.StatusGone, "pdn_closed")
		}
	case errors.Is(err, store.ErrPDNNotFound):
		if err := c.cleanupPDNsLocked(ctx); err != nil {
			return nil, pdnError(http.StatusServiceUnavailable,
				"cleanup_pending")
		}
		old, err := c.store.GetLivePDN(req.ImsiStr)
		if err == nil && old != nil {
			return nil, pdnError(http.StatusConflict, "pdn_conflict")
		}
		if !errors.Is(err, store.ErrPDNNotFound) {
			return nil, err
		}
		old, err = c.store.GetPDNByIP(ip.String())
		if err == nil && old != nil {
			return nil, pdnError(http.StatusConflict, "address_in_use")
		}
		if !errors.Is(err, store.ErrPDNNotFound) {
			return nil, err
		}
		p = &store.PDN{InstanceID: req.InstanceID, ID: req.PDNID,
			Imsi: req.ImsiStr, ApnName: apn, IP: ip.String(),
			State: store.PDNLive}
		if err := c.store.InsertPDN(p); err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	result, err := c.applyPDNLocked(ctx, p)
	if err != nil {
		// An unsuccessful attach must not later enable data in the monitor.
		if markErr := c.store.SetPDNState(p, store.PDNClosing); markErr != nil {
			c.restartRequired = true
			return nil, markErr
		}
		if closeErr := c.closePDNLocked(ctx, p); closeErr != nil {
			log.Warnf("PDN %s cleanup pending: %v", p.ID, closeErr)
		}
		return nil, pdnError(http.StatusServiceUnavailable,
			"pdn_setup_failed")
	}
	return result, nil
}

func (c *Controller) EndSession(ctx *gin.Context, req *api.EndSession) error {
	if !validIdentity(req.InstanceID) || !validIdentity(req.PDNID) {
		return pdnError(http.StatusBadRequest, "invalid_identity")
	}
	if !validIMSI(req.ImsiStr) {
		return pdnError(http.StatusBadRequest, "invalid_imsi")
	}
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	p, err := c.store.GetPDN(req.InstanceID, req.PDNID)
	if errors.Is(err, store.ErrPDNNotFound) {
		// DELETE may overtake POST. Keep a tombstone for the delayed POST.
		return c.store.InsertPDN(&store.PDN{InstanceID: req.InstanceID,
			ID: req.PDNID, Imsi: req.ImsiStr, State: store.PDNClosed})
	}
	if err != nil {
		return err
	}
	if p.Imsi != req.ImsiStr {
		return pdnError(http.StatusConflict, "pdn_conflict")
	}
	if p.State == store.PDNClosed {
		return nil
	}
	if err := c.store.SetPDNState(p, store.PDNClosing); err != nil {
		return err
	}
	if err := c.closePDNLocked(ctx, p); err != nil {
		return pdnError(http.StatusServiceUnavailable, "cleanup_pending")
	}
	return nil
}

func (c *Controller) closePDNLocked(ctx context.Context, p *store.PDN) error {
	if err := c.sm.EndSession(ctx,
		&store.Subscriber{Imsi: p.Imsi}); err != nil {
		return err
	}
	return c.store.SetPDNState(p, store.PDNClosed)
}

func (c *Controller) cleanupPDNsLocked(ctx context.Context) error {
	pdns, err := c.store.GetLivePDNs()
	if err != nil {
		return err
	}
	var first error
	for i := range pdns {
		if pdns[i].State != store.PDNClosing {
			continue
		}
		if err := c.closePDNLocked(ctx, &pdns[i]); err != nil && first == nil {
			first = err
		}
	}
	return first
}

// policyReason does no network lookup. Missing policy is a normal blocked
// PDN; the existing subscriber push endpoint supplies policy independently.
func (c *Controller) policyReason(imsi string, se *store.Session) (
	*store.Subscriber, string, error) {
	sub, err := c.store.GetSubscriber(imsi)
	if errors.Is(err, store.ErrSubscriberNotFound) {
		return nil, "policy_missing", nil
	}
	if err != nil {
		return nil, "", err
	}
	if !c.serviceEnabled() {
		return sub, "service_disabled", nil
	}
	if !sub.ServiceOn {
		return sub, "subscriber_disabled", nil
	}
	p := &sub.PolicyID
	now := time.Now().Unix()
	if p.StartTime > now {
		return sub, "package_not_started", nil
	}
	if p.EndTime <= now {
		return sub, "package_expired", nil
	}
	u, err := c.store.GetUsageByImsi(imsi)
	if err != nil {
		return nil, "", err
	}
	// Subtract only after bounds checks; consumed may exceed the allowance.
	if p.Consumed >= p.Data || u.Data >= p.Data-p.Consumed {
		return sub, "quota_exhausted", nil
	}
	left := p.Data - p.Consumed - u.Data
	if se != nil && se.PolicyID.ID == p.ID && se.TotalBytes >= left {
		return sub, "quota_exhausted", nil
	}
	return sub, "none", nil
}

func sameSessionPolicy(se *store.Session, sub *store.Subscriber) bool {
	p := &sub.PolicyID
	return se.PolicyID.ID == p.ID &&
		se.RxMeterID.Rate == p.Ulbr && se.TxMeterID.Rate == p.Dlbr &&
		se.RxMeterID.Burst == p.Burst && se.TxMeterID.Burst == p.Burst
}

func (c *Controller) applyPDNLocked(ctx context.Context,
	p *store.PDN) (*api.PDNResponse, error) {
	if p.State == store.PDNClosing {
		err := c.closePDNLocked(ctx, p)
		return bindingResponse(p, "unknown", "cleanup_pending"), err
	}
	se, err := c.store.GetActiveSessionByImsi(p.Imsi)
	if err != nil && !errors.Is(err, store.ErrSessionNotFound) {
		return nil, err
	}
	sub, reason, err := c.policyReason(p.Imsi, se)
	if err != nil {
		return nil, err
	}
	if se != nil && (sub == nil || !sameSessionPolicy(se, sub) ||
		se.UeIpAddr != p.IP || se.ApnName != p.ApnName ||
		!c.sm.HasActiveSession(p.Imsi)) {
		if err := c.sm.EndSession(ctx,
			&store.Subscriber{Imsi: p.Imsi}); err != nil {
			return nil, err
		}
		se = nil
		// Ending a segment commits its usage to the existing allowance.
		sub, reason, err = c.policyReason(p.Imsi, nil)
		if err != nil {
			return nil, err
		}
	}
	if reason != "none" {
		if se != nil {
			if err := c.sm.PauseSession(ctx, sub); err != nil {
				return nil, err
			}
		}
		return bindingResponse(p, "blocked", reason), nil
	}
	if err := c.checkDataPathLocked(); err != nil {
		return nil, err
	}
	if se == nil {
		ns, rx, tx, err := c.store.CreatePDNSession(sub, p.IP,
			c.nodeId, p.ApnName)
		if err != nil {
			return nil, err
		}
		if err := c.sm.CreateSesssion(ctx, sub, ns, rx, tx); err != nil {
			return nil, err
		}
	}
	if err := c.sm.ResumeSession(ctx, sub); err != nil {
		return nil, err
	}
	return bindingResponse(p, "allowed", "none"), nil
}

func (c *Controller) reconcileImsiLocked(ctx context.Context,
	imsi string) error {
	p, err := c.store.GetLivePDN(imsi)
	if errors.Is(err, store.ErrPDNNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = c.applyPDNLocked(ctx, p)
	return err
}

func (c *Controller) reconcilePDNsLocked(ctx context.Context) error {
	pdns, err := c.store.GetLivePDNs()
	if err != nil {
		return err
	}
	var first error
	for i := range pdns {
		_, applyErr := c.applyPDNLocked(ctx, &pdns[i])
		if applyErr != nil && first == nil {
			first = applyErr
		}
	}
	return first
}

// GetPDN only observes state. It never creates or repairs subscriber flows.
func (c *Controller) GetPDN(ctx *gin.Context,
	req *api.GetPDN) (*api.PDNResponse, error) {
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	p, err := c.store.GetLivePDN(req.Imsi)
	if errors.Is(err, store.ErrPDNNotFound) {
		return nil, pdnError(http.StatusNotFound, "pdn_not_found")
	}
	if err != nil {
		return nil, err
	}
	if p.State != store.PDNLive {
		return bindingResponse(p, "unknown", "cleanup_pending"), nil
	}
	if c.restartRequired || !c.sm.Status().DataPath.Connected {
		return bindingResponse(p, "unknown", "datapath_unavailable"), nil
	}
	se, err := c.store.GetActiveSessionByImsi(p.Imsi)
	if err != nil && !errors.Is(err, store.ErrSessionNotFound) {
		return nil, err
	}
	sub, reason, err := c.policyReason(p.Imsi, se)
	if err != nil {
		return nil, err
	}
	if se != nil && se.FlowState != store.FlowsPaused {
		if reason == "none" && se.FlowState == store.FlowsActive &&
			sameSessionPolicy(se, sub) && c.sm.HasActiveSession(p.Imsi) {
			return bindingResponse(p, "allowed", reason), nil
		}
		return bindingResponse(p, "unknown", "policy_update_pending"), nil
	}
	if reason == "none" {
		reason = "policy_update_pending"
	}
	return bindingResponse(p, "blocked", reason), nil
}

func (c *Controller) checkDataPathLocked() error {
	dp := c.sm.Status().DataPath
	if !dp.Connected || dp.ConnectedCount > 1 {
		c.restartRequired = true
	}
	if c.restartRequired || c.stopping {
		return pdnError(http.StatusServiceUnavailable, "restart_required")
	}
	return nil
}

func (c *Controller) startPDNMonitor(period time.Duration) {
	if period <= 0 {
		period = 2 * time.Second
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.pdnCancel = cancel
	c.pdnDone = make(chan struct{})
	go func() {
		defer close(c.pdnDone)
		ticker := time.NewTicker(period)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.pdnMu.Lock()
				if err := c.checkDataPathLocked(); err != nil {
					if err := c.store.RetireEPCs(); err != nil {
						log.Errorf("Retire EPC after disconnect: %v", err)
					}
				}
				if err := c.reconcilePDNsLocked(ctx); err != nil {
					log.Warnf("PDN reconciliation pending: %v", err)
				}
				c.pdnMu.Unlock()
			}
		}
	}()
}

// GetEPC supplies the compare-and-replace value used on process startup.
// This prevents a delayed, previously unseen boot request replacing a newer
// process. Reading the owner never changes or restores PDNs.
func (c *Controller) GetEPC(ctx *gin.Context) (*api.EPCResponse, error) {
	c.pdnMu.Lock()
	defer c.pdnMu.Unlock()
	id, err := c.store.CurrentEPC()
	if err != nil {
		return nil, err
	}
	state := "ready"
	if c.restartRequired || c.stopping {
		state = "restart_required"
	}
	return &api.EPCResponse{InstanceID: id, State: state}, nil
}
