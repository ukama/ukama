/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package session

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ukama/ukama/nodes/apps/pcrf/pkg"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/datapath"

	log "github.com/sirupsen/logrus"
)

type Status struct {
	DataPath datapath.Status `json:"datapath"`
}

type sessionCache struct {
	s           *store.Session
	txCookie    uint64
	rxCookie    uint64
	InitUsage   uint64
	cancel      context.CancelFunc
	ctx         context.Context
	statsDirty  bool
	ending      bool
	paused      bool
	baseRxBytes uint64 // sc.s.RxBytes at the moment flows were last removed
	baseTxBytes uint64 // sc.s.TxBytes at the moment flows were last removed
}

type sessionManager struct {
	period          time.Duration `default:"2s"`
	idle            time.Duration `default:"60s"`
	newSessionGrace time.Duration `default:"10m"`
	store           *store.Store
	d               datapath.DataPath
	mu              sync.Mutex
	cache           map[string]*sessionCache
}

type SessionManager interface {
	CreateSesssion(ctx context.Context, sub *store.Subscriber, ns *store.Session, rxf *store.Flow, txf *store.Flow) error
	EndSession(ctx context.Context, sub *store.Subscriber) error
	IfSessionExist(ctx context.Context, imsi, ip string) bool
	EndAllSessions() error
	Status() Status
	PauseSession(ctx context.Context, sub *store.Subscriber) error
	ResumeSession(ctx context.Context, sub *store.Subscriber) error
	HasActiveSession(imsi string) bool
}

func NewSessionManager(store *store.Store, br pkg.BrdigeConfig) (*sessionManager, error) {
	d, err := datapath.InitDataPath(br.Name, br.Ip, br.NetType, br.Management)
	if err != nil {
		log.Errorf("Error initializing session manager. Error: %v", err)

		return nil, fmt.Errorf("error initializing session manager. Error: %w", err)
	}

	s := &sessionManager{
		d:               d,
		store:           store,
		period:          br.Period,
		idle:            br.SessionIdleTime,
		newSessionGrace: br.NewSessionGrace,
		cache:           make(map[string]*sessionCache),
	}

	if err := s.recoverSessions(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *sessionManager) Status() Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	return Status{
		DataPath: s.d.Status(),
	}
}

func (s *sessionManager) storeStats(imsi string, lastStats bool) error {
	sc, ok := s.cache[imsi]
	if !ok {
		return fmt.Errorf("session for imsi not found: %s", imsi)
	}
	if sc.paused {
		return nil
	}
	rx, _, tx, _, err := s.d.DataPathStats(sc.rxCookie, sc.txCookie)
	if err != nil {
		return err
	}
	// A reset switch must not erase counters already recorded locally.
	if sc.baseRxBytes+rx < sc.s.RxBytes ||
		sc.baseTxBytes+tx < sc.s.TxBytes {
		return fmt.Errorf("flow counters decreased for imsi %s", imsi)
	}
	if !sc.statsDirty && sc.s.RxBytes == sc.baseRxBytes+rx &&
		sc.s.TxBytes == sc.baseTxBytes+tx {
		return nil
	}
	sc.s.RxBytes = sc.baseRxBytes + rx
	sc.s.TxBytes = sc.baseTxBytes + tx
	sc.s.TotalBytes = sc.s.RxBytes + sc.s.TxBytes
	sc.s.UpdatedAt = uint64(time.Now().Unix())
	sc.statsDirty = true
	if err := s.store.UpdateSessionUsage(sc.s); err != nil {
		return err
	}
	sc.statsDirty = false
	return nil
}

func (s *sessionManager) IfSessionExist(ctx context.Context,
	imsi, ip string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	sc, ok := s.cache[imsi]
	return ok && sc.s.UeIpAddr == ip
}

// A zero bandwidth is unspecified, not an OpenFlow meter with rate zero.
// Cookie-based flow counters still enforce the package byte allowance.

func datapathMeterID(m store.Meter) uint32 {
	if m.Rate == 0 {
		return 0
	}
	return uint32(m.ID)
}

func (s *sessionManager) CreateSesssion(ctx context.Context, sub *store.Subscriber, ns *store.Session, rxf *store.Flow, txf *store.Flow) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sc := sessionCache{
		s:        ns,
		txCookie: txf.Cookie,
		rxCookie: rxf.Cookie,
	}

	u, err := s.store.GetUsageByImsi(sub.Imsi)
	if err != nil {
		log.Errorf("Error getting usage for Imsi %s.Error: %v", sub.Imsi, err)

		return fmt.Errorf("error getting usage for Imsi %s.Error: %w", sub.Imsi, err)
	}
	sc.InitUsage = u.Data

	if _, exists := s.cache[sub.Imsi]; exists {
		return fmt.Errorf("session already exists for imsi %s", sub.Imsi)
	}
	s.cache[sub.Imsi] = &sc
	sc.ending = true
	// Also retain failed installations so cleanup can be retried.
	if err := s.StartSessionMonitor(ctx, sub.Imsi); err != nil {
		return err
	}
	if ns.FlowState == store.FlowsPaused {
		sc.paused = true

		err = s.d.AddMetersOnly(
			datapathMeterID(sc.s.RxMeterID),
			datapathMeterID(sc.s.TxMeterID),
			uint32(sc.s.RxMeterID.Rate),
			uint32(sc.s.TxMeterID.Rate),
			uint32(sc.s.RxMeterID.Burst))
		if err != nil {
			log.Errorf("Failed to add meters for Imsi %s. Error: %v", sub.Imsi, err)

			return fmt.Errorf("failed to add meters for Imsi %s. Error: %w", sub.Imsi, err)
		}
	} else {
		err = s.d.AddNewDataPath(sc.s.UeIpAddr,
			datapathMeterID(sc.s.RxMeterID),
			datapathMeterID(sc.s.TxMeterID),
			uint32(sc.s.RxMeterID.Rate),
			uint32(sc.s.TxMeterID.Rate),
			uint32(sc.s.RxMeterID.Burst),
			sc.rxCookie,
			sc.txCookie)
		if err != nil {
			log.Errorf("Failed to add data path for Imsi %s. Error: %v", sub.Imsi, err)

			return fmt.Errorf("failed to add data path for Imsi %s. Error: %w", sub.Imsi, err)
		}
	}

	sc.ending = false

	return nil
}

func (s *sessionManager) EndAllSessions() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var first error
	for imsi := range s.cache {
		if err := s.endSessionLocked(context.Background(),
			&store.Subscriber{Imsi: imsi}); err != nil && first == nil {
			first = err
		}
	}
	return first
}

func (s *sessionManager) EndSession(ctx context.Context,
	sub *store.Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.endSessionLocked(ctx, sub)
}

func (s *sessionManager) endSessionLocked(ctx context.Context,
	sub *store.Subscriber) error {
	sc, ok := s.cache[sub.Imsi]
	if !ok {
		// Creation may have persisted a session before caching it failed.
		se, err := s.store.GetActiveSessionByImsi(sub.Imsi)
		if errors.Is(err, store.ErrSessionNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := s.d.DeleteDataPath(se.UeIpAddr,
			datapathMeterID(se.RxMeterID),
			datapathMeterID(se.TxMeterID)); err != nil {
			return err
		}
		return s.store.EndSession(se)
	}
	sc.ending = true
	if !sc.paused {
		if err := s.storeStats(sub.Imsi, false); err != nil {
			log.Warnf("Finalizing %s with recorded counters: %v", sub.Imsi, err)
		}
	}
	if err := s.store.UpdateSessionFlowState(sc.s.ID,
		store.FlowStateUnknown); err != nil {
		log.Warnf("Pending flow state for %s: %v", sub.Imsi, err)
	}
	sc.s.FlowState = store.FlowStateUnknown
	var err error
	if sc.paused {
		err = s.d.DeleteMetersOnly(datapathMeterID(sc.s.RxMeterID),
			datapathMeterID(sc.s.TxMeterID))
	} else {
		err = s.d.DeleteDataPath(sc.s.UeIpAddr,
			datapathMeterID(sc.s.RxMeterID), datapathMeterID(sc.s.TxMeterID))
	}
	if err != nil {
		return err
	}
	sc.paused = true
	if err := s.store.EndSession(sc.s); err != nil {
		return err
	}
	if sc.cancel != nil {
		sc.cancel()
	}
	delete(s.cache, sub.Imsi)
	return nil
}

func (s *sessionManager) HasActiveSession(imsi string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	sc, ok := s.cache[imsi]
	return ok && !sc.ending
}

func (s *sessionManager) PauseSession(ctx context.Context,
	sub *store.Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sc, ok := s.cache[sub.Imsi]
	if !ok {
		log.Errorf("failed to find session for Imsi %s", sub.Imsi)
		return nil
	}

	if sc.ending {
		return fmt.Errorf("session cleanup pending")
	}
	if sc.paused {
		return s.store.UpdateSessionFlowState(sc.s.ID, store.FlowsPaused)
	}

	if err := s.storeStats(sub.Imsi, false); err != nil {
		log.Warnf("Failed to flush stats before pausing %s: %v",
			sub.Imsi, err)
	}

	sc, ok = s.cache[sub.Imsi]
	if !ok || sc.paused {
		return nil
	}

	if err := s.store.UpdateSessionFlowState(sc.s.ID,
		store.FlowStateUnknown); err != nil {
		log.Warnf("Pending flow state for %s: %v", sub.Imsi, err)
	}
	sc.s.FlowState = store.FlowStateUnknown
	err := s.d.DeleteFlowOnly(sc.s.UeIpAddr)
	if err != nil {
		sc.ending = true
		log.Errorf("Failed to delete flow for Imsi %s: %v", sub.Imsi, err)

		return fmt.Errorf("failed to delete flow for Imsi %s: %w",
			sub.Imsi, err)
	}

	sc.baseRxBytes = sc.s.RxBytes
	sc.baseTxBytes = sc.s.TxBytes
	sc.paused = true
	sc.s.FlowState = store.FlowsPaused

	if err := s.store.UpdateSessionFlowState(sc.s.ID,
		store.FlowsPaused); err != nil {
		log.Errorf("Failed to persist paused flow for %s: %v", sub.Imsi, err)

		return fmt.Errorf("failed to persist paused flow for %s: %w",
			sub.Imsi, err)
	}

	log.Infof("[SessionId %d ] Paused session for Imsi %s", sc.s.ID, sub.Imsi)

	return nil
}

func (s *sessionManager) ResumeSession(ctx context.Context,
	sub *store.Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sc, ok := s.cache[sub.Imsi]
	if !ok {
		log.Errorf("failed to find session for Imsi %s", sub.Imsi)
		return nil
	}

	if sc.ending {
		return fmt.Errorf("session cleanup pending")
	}
	if !sc.paused {
		return nil
	}

	if err := s.store.UpdateSessionFlowState(sc.s.ID,
		store.FlowStateUnknown); err != nil {
		return err
	}
	sc.s.FlowState = store.FlowStateUnknown
	err := s.d.AddFlowOnly(sc.s.UeIpAddr,
		datapathMeterID(sc.s.RxMeterID),
		datapathMeterID(sc.s.TxMeterID),
		sc.rxCookie,
		sc.txCookie)
	if err != nil {
		log.Errorf("Failed to add flow for Imsi %s. Error: %v", sub.Imsi, err)

		sc.ending = true
		sc.paused = false
		return fmt.Errorf("failed to add flow for Imsi %s: %w", sub.Imsi, err)
	}

	sc.paused = false
	sc.s.FlowState = store.FlowsActive

	sc.s.UpdatedAt = uint64(time.Now().Unix())

	if err := s.store.UpdateSessionUsage(sc.s); err != nil {
		log.Errorf("Failed to persist resumed session state for Imsi %s. Error: %v", sub.Imsi, err)

		sc.ending = true
		return fmt.Errorf("failed to persist resumed session for %s: %w",
			sub.Imsi, err)
	}

	if err := s.store.UpdateSessionFlowState(sc.s.ID, store.FlowsActive); err != nil {
		log.Errorf("Failed to persist active flow state for Imsi %s. Error: %v", sub.Imsi, err)

		sc.ending = true
		return fmt.Errorf("failed to persist active flow state for %s: %w",
			sub.Imsi, err)
	}

	log.Infof("[SessionId %d ] Resumed session for Imsi %s", sc.s.ID, sub.Imsi)

	return nil
}

func (s *sessionManager) SendCDR(imsi string) error {
	sc := s.cache[imsi]
	if sc == nil {
		return fmt.Errorf("session for imsi %s not found", imsi)
	}

	log.Infof("[ SessionId %d ] Marking CDR ready for subscriber %s and IP address %s",
		sc.s.ID, imsi, sc.s.UeIpAddr)

	return s.store.UpdateSessionSyncState(sc.s.ID, store.SessionSyncReady)
}

func (s *sessionManager) StartSessionMonitor(ctx context.Context, imsi string) error {
	sc := s.cache[imsi]
	if sc == nil {
		return fmt.Errorf("session for imsi %s not found", imsi)
	}

	log.Infof("[SessionId %d ] Starting session monitor for subscriber %s and IP address %s",
		sc.s.ID, imsi, sc.s.UeIpAddr)

	sc.ctx, sc.cancel = context.WithCancel(context.Background())
	s.cache[imsi] = sc

	go s.sessionMonitorRoutine(sc.ctx, s.period, sc)

	return nil
}

func (s *sessionManager) StopSessionMonitor(ctx context.Context, imsi string) error {
	sc := s.cache[imsi]
	if sc == nil {
		return fmt.Errorf("session for imsi %s not found", imsi)
	}

	log.Infof("[SessionId %d ] Stop session monitor for subscriber %s and IP address %s",
		sc.s.ID, imsi, sc.s.UeIpAddr)

	if sc.cancel != nil {
		sc.cancel()
	}

	return nil
}

func (s *sessionManager) sessionMonitorRoutine(ctx context.Context, interval time.Duration, sc *sessionCache) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Infof("[SessionId %d ] Stat Collection", sc.s.ID)
			s.mu.Lock()
			// A canceled ticker must not touch a replacement session.
			if s.cache[sc.s.SubscriberID.Imsi] == sc {
				if sc.ending {
					_ = s.endSessionLocked(ctx, &sc.s.SubscriberID)
				} else {
					_ = s.storeStats(sc.s.SubscriberID.Imsi, false)
				}
			}
			s.mu.Unlock()

		case <-ctx.Done():
			log.Infof("[SessionId %d ] Exiting monitoring for subscriber %s",
				sc.s.ID, sc.s.SubscriberID.Imsi)

			return
		}
	}
}

// recoverSessions is called only during startup, before admitting PDNs.
func (s *sessionManager) recoverSessions() error {
	old, err := s.store.GetAllActiveSessions()
	if err != nil {
		return err
	}
	for i := range old {
		se := &old[i]
		if err := s.d.DeleteDataPath(se.UeIpAddr,
			datapathMeterID(se.RxMeterID),
			datapathMeterID(se.TxMeterID)); err != nil {
			return err
		}
		if err := s.store.EndSession(se); err != nil {
			return err
		}
	}
	return nil
}
