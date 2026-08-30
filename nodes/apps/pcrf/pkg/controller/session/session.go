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

	return s, nil
}

func (s *sessionManager) Status() Status {
	return Status{
		DataPath: s.d.Status(),
	}
}

func (s *sessionManager) storeStats(imsi string, lastStats bool) error {
	var err error
	sc, ok := s.cache[imsi]
	if ok {
		if sc.paused {
			if lastStats {
				sc.s.UpdatedAt = uint64(time.Now().Unix())

				return s.store.EndSession(sc.s)
			}

			tNow := time.Now().Unix()
			temp := int64(sc.s.UpdatedAt + uint64(s.idle.Seconds()))

			if tNow > temp {
				log.Infof("[SessionId %d ] Subscriber %s has been paused/idle for more than %s since %d. Ending session.",
					sc.s.ID, imsi, s.idle, sc.s.UpdatedAt)

				_ = s.endSessionLocked(sc.ctx, &store.Subscriber{Imsi: imsi})

				return fmt.Errorf("session idle timeout exceeded while paused")
			}

			return nil
		}

		var rx, tx uint64
		rx, _, tx, _, err = s.d.DataPathStats(sc.rxCookie, sc.txCookie)
		if err != nil {
			log.Errorf("[SessionId %d ] Failed to read final stats for data path of Imsi %s. Error: %v",
				sc.s.ID, sc.s.SubscriberID.Imsi, err)

			return fmt.Errorf("failed to read final stats for data path of Imsi %s. Error: %w",
				sc.s.SubscriberID.Imsi, err)
		}

		sc.s.RxBytes = sc.baseRxBytes + rx
		sc.s.TxBytes = sc.baseTxBytes + tx

		log.Infof("Rx Cookie 0x%x Rx Bytes %d Tx Cookie 0x%x TxBytes %d for imsi %s",
			sc.rxCookie, sc.s.RxBytes, sc.txCookie, sc.s.TxBytes, imsi)

		tNow := time.Now().Unix()
		lastUpdate := sc.s.UpdatedAt

		totalBytes := sc.s.TxBytes + sc.s.RxBytes

		if lastStats {
			sc.s.UpdatedAt = uint64(tNow)
			sc.s.TotalBytes = sc.s.TxBytes + sc.s.RxBytes

			err = s.store.EndSession(sc.s)
			if err != nil {
				log.Warnf("[SessionId %d ] Failed to update last session usage to db store for Imsi %s. Error: %s",
					sc.s.ID, sc.s.SubscriberID.Imsi, err.Error())
			}
		} else {
			if totalBytes != sc.s.TotalBytes {
				sc.s.UpdatedAt = uint64(tNow)
				sc.s.TotalBytes = sc.s.TxBytes + sc.s.RxBytes

				err = s.store.UpdateSessionUsage(sc.s)
				if err != nil {
					log.Warnf("[SessionId %d ] Failed to update session usage to db store for Imsi %s. Error: %s",
						sc.s.ID, sc.s.SubscriberID.Imsi, err.Error())
				}

				p, err := s.store.GetApplicablePolicyByImsi(imsi)
				if err != nil {
					log.Errorf("[SessionId %d ] failed to get policy by Imsi for subscriber %s. Error %v",
						sc.s.ID, imsi, err)

					return fmt.Errorf("failed to get policy by Imsi for subscriber %s. Error %w",
						imsi, err)
				}

				totalUsage := sc.InitUsage + sc.s.TotalBytes
				availableData := p.Data - p.Consumed
				if totalUsage >= availableData {
					log.Errorf("[SessionId %d ] Subscriber %s hit max data limit available=%d totalUsage=%d",
						sc.s.ID, imsi, availableData, totalUsage)

					_ = s.endSessionLocked(sc.ctx, &store.Subscriber{Imsi: imsi})

					return fmt.Errorf("max data cap limit exceeded")
				}
			}

			threshold := s.idle
			if sc.s.TotalBytes == 0 {
				threshold = s.newSessionGrace
			}

			temp := int64(lastUpdate + uint64(threshold.Seconds()))
			log.Debugf("[SessionId %d ] Subscriber %s time now %d timeout %d",
				sc.s.ID, imsi, tNow, temp)

			if tNow > temp {
				log.Infof("[SessionId %d ] Subscriber %s is idle for more than %s since %d. Ending session.",
					sc.s.ID, imsi, threshold, lastUpdate)

				_ = s.endSessionLocked(sc.ctx, &store.Subscriber{Imsi: imsi})

				return fmt.Errorf("session idle timeout exceeded")
			}
		}

		log.Debugf("[SessionId %d ] Updated stats for %s are %+v", sc.s.ID, imsi, sc.s)
	} else {
		log.Errorf("Session for Imsi %s not found.", imsi)

		return fmt.Errorf("session for imsi not found: %s", imsi)
	}

	return err
}

func (s *sessionManager) IfSessionExist(ctx context.Context, imsi, ip string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	sc, ok := s.cache[imsi]
	if ok {
		if sc.s.UeIpAddr == ip {
			return true
		}

		log.Errorf("Old session exists for subscriber %s with IP addr %s. Ending it.", imsi, sc.s.UeIpAddr)
		_ = s.endSessionLocked(ctx, &store.Subscriber{Imsi: imsi})
	}

	return false
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

	if ns.FlowState == store.FlowsPaused {
		sc.paused = true

		err = s.d.AddMetersOnly(
			uint32(sc.s.RxMeterID.ID),
			uint32(sc.s.TxMeterID.ID),
			uint32(sc.s.RxMeterID.Rate),
			uint32(sc.s.TxMeterID.Rate),
			uint32(sc.s.RxMeterID.Burst))
		if err != nil {
			log.Errorf("Failed to add meters for Imsi %s. Error: %v", sub.Imsi, err)

			return fmt.Errorf("failed to add meters for Imsi %s. Error: %w", sub.Imsi, err)
		}
	} else {
		err = s.d.AddNewDataPath(sc.s.UeIpAddr,
			uint32(sc.s.RxMeterID.ID),
			uint32(sc.s.TxMeterID.ID),
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

	s.cache[sub.Imsi] = &sc

	err = s.StartSessionMonitor(ctx, sub.Imsi)
	if err != nil {
		log.Errorf("Failed to start monitor for Imsi %s. Error: %v", sub.Imsi, err)

		return fmt.Errorf("failed to start monitor for Imsi %s. Error: %w", sub.Imsi, err)
	}

	return nil
}

func (s *sessionManager) EndAllSessions() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for imsi, session := range s.cache {
		err := s.endSessionLocked(context.Background(), &store.Subscriber{Imsi: imsi})
		if err != nil {
			log.Errorf("Failed to end session for Imsi %s. Error %v", imsi, err)
		}
		log.Infof("Ending session %+v for Imsi %s.", session, imsi)
	}

	return nil
}

func (s *sessionManager) EndSession(ctx context.Context, sub *store.Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.endSessionLocked(ctx, sub)
}

func (s *sessionManager) endSessionLocked(ctx context.Context, sub *store.Subscriber) error {
	sc, ok := s.cache[sub.Imsi]
	if !ok {
		log.Errorf("failed to find session for Imsi %s", sub.Imsi)
		return nil
	}

	err := s.StopSessionMonitor(ctx, sub.Imsi)
	if err != nil {
		log.Errorf("Failed to stop monitor for Imsi %s. Error: %v", sub.Imsi, err)

		return fmt.Errorf("failed to stop monitor for Imsi %s. Error: %w", sub.Imsi, err)
	}

	err = s.storeStats(sc.s.SubscriberID.Imsi, true)
	if err != nil {
		log.Warnf("Failed to store final stats for Imsi %s. Error: %v", sub.Imsi, err)
	}

	time.Sleep(1000 * time.Millisecond)

	if sc.paused {
		err = s.d.DeleteMetersOnly(uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID))
	} else {
		err = s.d.DeleteDataPath(sc.s.UeIpAddr, uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID))
	}
	if err != nil {
		log.Errorf("Failed to delete data path for Imsi %s. Error: %v", sub.Imsi, err)

		return fmt.Errorf("failed to delete data path for Imsi %s. Error: %w", sub.Imsi, err)
	}

	_ = s.SendCDR(sub.Imsi)

	delete(s.cache, sub.Imsi)

	return nil
}

func (s *sessionManager) HasActiveSession(imsi string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.cache[imsi]
	return ok
}

func (s *sessionManager) PauseSession(ctx context.Context, sub *store.Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sc, ok := s.cache[sub.Imsi]
	if !ok {
		log.Errorf("failed to find session for Imsi %s", sub.Imsi)
		return nil
	}

	if sc.paused {
		return nil
	}

	if err := s.storeStats(sub.Imsi, false); err != nil {
		log.Warnf("Failed to flush stats before pausing session for Imsi %s. Error: %v", sub.Imsi, err)
	}

	sc, ok = s.cache[sub.Imsi]
	if !ok || sc.paused {
		return nil
	}

	sc.baseRxBytes = sc.s.RxBytes
	sc.baseTxBytes = sc.s.TxBytes

	err := s.d.DeleteFlowOnly(sc.s.UeIpAddr)
	if err != nil {
		log.Errorf("Failed to delete flow for Imsi %s. Error: %v", sub.Imsi, err)

		return fmt.Errorf("failed to delete flow for Imsi %s. Error: %w", sub.Imsi, err)
	}

	sc.paused = true
	sc.s.FlowState = store.FlowsPaused

	if err := s.store.UpdateSessionFlowState(sc.s.ID, store.FlowsPaused); err != nil {
		log.Errorf("Failed to persist paused flow state for Imsi %s. Error: %v", sub.Imsi, err)

		return fmt.Errorf("failed to persist paused flow state for Imsi %s. Error: %w", sub.Imsi, err)
	}

	log.Infof("[SessionId %d ] Paused session for Imsi %s", sc.s.ID, sub.Imsi)

	return nil
}

func (s *sessionManager) ResumeSession(ctx context.Context, sub *store.Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sc, ok := s.cache[sub.Imsi]
	if !ok {
		log.Errorf("failed to find session for Imsi %s", sub.Imsi)
		return nil
	}

	if !sc.paused {
		return nil
	}

	err := s.d.AddFlowOnly(sc.s.UeIpAddr,
		uint32(sc.s.RxMeterID.ID),
		uint32(sc.s.TxMeterID.ID),
		sc.rxCookie,
		sc.txCookie)
	if err != nil {
		log.Errorf("Failed to add flow for Imsi %s. Error: %v", sub.Imsi, err)

		return fmt.Errorf("failed to add flow for Imsi %s. Error: %w", sub.Imsi, err)
	}

	sc.paused = false
	sc.s.FlowState = store.FlowsActive

	sc.s.UpdatedAt = uint64(time.Now().Unix())

	if err := s.store.UpdateSessionUsage(sc.s); err != nil {
		log.Errorf("Failed to persist resumed session state for Imsi %s. Error: %v", sub.Imsi, err)

		return fmt.Errorf("failed to persist resumed session state for Imsi %s. Error: %w", sub.Imsi, err)
	}

	if err := s.store.UpdateSessionFlowState(sc.s.ID, store.FlowsActive); err != nil {
		log.Errorf("Failed to persist active flow state for Imsi %s. Error: %v", sub.Imsi, err)

		return fmt.Errorf("failed to persist active flow state for Imsi %s. Error: %w", sub.Imsi, err)
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
			_ = s.storeStats(sc.s.SubscriberID.Imsi, false)
			s.mu.Unlock()

		case <-ctx.Done():
			log.Infof("[SessionId %d ] Exiting monitoring for subscriber %s",
				sc.s.ID, sc.s.SubscriberID.Imsi)

			return
		}
	}
}
