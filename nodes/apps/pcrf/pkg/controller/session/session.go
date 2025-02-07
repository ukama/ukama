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
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/datapath"
)

// type Cache struct {
// 	imsi    string
// 	ip      string
// 	rxMeter uint32
// 	txMeter uint32
// }

type sessionCache struct {
	s              *store.Session
	txCookie       uint64
	rxCookie       uint64
	InitUsage      uint64
	idleReportSent bool
	cancel         context.CancelFunc
	ctx            context.Context
}
type sessionManager struct {
	period time.Duration `default:"2s"`
	idle   time.Duration `default:"60s"`
	store  *store.Store
	d      datapath.DataPath
	cache  map[string]*sessionCache
	rc     client.RemoteController
}

type SessionManager interface {
	CreateSesssion(ctx context.Context, sub *store.Subscriber, ns *store.Session, rxf *store.Flow, txf *store.Flow) error
	EndSession(ctx context.Context, sub *store.Subscriber) error
	IfSessionExist(ctx context.Context, imsi, ip string) bool
	EndAllSessions() error
}

func NewSessionManager(rc client.RemoteController, store *store.Store, br pkg.BrdigeConfig) *sessionManager {
	d, err := datapath.InitDataPath(br.Name, br.Ip, br.NetType, br.Management)
	if err != nil {
		log.Fatalf("error initializing session manager. Error: %s", err.Error())
	}

	s := &sessionManager{
		d:      d,
		store:  store,
		period: br.Period,
		idle:   br.SessionIdleTime,
		rc:     rc,
		cache:  make(map[string]*sessionCache),
	}

	return s
}

func (s *sessionManager) storeStats(imsi string, lastStats bool) error {
	var err error
	sc, ok := s.cache[imsi]
	if ok {

		/* Read sats */
		sc.s.RxBytes, _, sc.s.TxBytes, _, err = s.d.DataPathStats(sc.rxCookie, sc.txCookie)
		if err != nil {
			log.Errorf("[SessionId %d ] Failed to read final stats for data path of Imsi %s. Error: %s", sc.s.ID, sc.s.SubscriberID.Imsi, err.Error())
			return err
		}
		log.Infof("Rx Cookie 0x%x Rx Bytes %d Tx Cookie 0x%x TxBytes %d for imsi %s", sc.rxCookie, sc.s.RxBytes, sc.txCookie, sc.s.TxBytes, imsi)

		tNow := time.Now().Unix()
		lastUpdate := sc.s.UpdatedAt

		totalBytes := sc.s.TxBytes + sc.s.RxBytes

		/* Update to DB */
		if lastStats {
			sc.s.UpdatedAt = uint64(tNow)
			/* This adds the stats for TX and RX and store them*/
			err = s.store.EndSession(sc.s)
			if err != nil {
				log.Warnf("[SessionId %d ] Failed to update last session usage to db store for Imsi %s. Error: %s", sc.s.ID, sc.s.SubscriberID.Imsi, err.Error())
			}

		} else {

			/* Only do this when RX or TX value changes */
			if totalBytes != sc.s.TotalBytes {
				sc.idleReportSent = false // Reset the idleReportSent flag on change in stats
				sc.s.UpdatedAt = uint64(tNow)
				sc.s.TotalBytes = sc.s.TxBytes + sc.s.RxBytes
				err = s.store.UpdateSessionUsage(sc.s)
				if err != nil {
					log.Warnf("[SessionId %d ] Failed to update session usage to db store for Imsi %s. Error: %s", sc.s.ID, sc.s.SubscriberID.Imsi, err.Error())
				}

				p, err := s.store.GetApplicablePolicyByImsi(imsi)
				if err != nil {
					log.Errorf("[SessionId %d ] failed to get policy by Imsi for subscriber %s. Error %v", sc.s.ID, imsi, err)
					return err
				}

				totalUsage := sc.InitUsage + sc.s.TotalBytes
				availableData := p.Data - p.Consumed
				if totalUsage >= availableData {
					/* this means we need to terminates session */
					log.Errorf("[SessionId %d ] Subscriber %s hit the max data CapLimits of %d with previusly used data %d current used data %d.", sc.s.ID, imsi, p.Consumed, p.Data, totalUsage)
					_ = s.EndSession(sc.ctx, &store.Subscriber{Imsi: imsi})
					return fmt.Errorf("max data cap limit exceeded")
				}
			}

			/* If report is laready sent no need to send again */
			temp := (int64)(lastUpdate + uint64(s.idle.Seconds()))
			log.Debugf("[SessionId %d ] Subscriber %s Idle session report flag %v, time now %d and timeout Val %d", sc.s.ID, imsi, sc.idleReportSent, tNow, temp)
			if !sc.idleReportSent && tNow > temp {
				log.Infof("[SessionId %d ] Subscriber %s is idle for more than %d secconds from %d.", sc.s.ID, imsi, s.idle, lastUpdate)
				/* Sync data to cloud */
				err = s.SendCDR(imsi)
				if err == nil {
					sc.idleReportSent = true
				}

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
	sc, ok := s.cache[imsi]
	if ok {
		if sc.s.UeIpAddr == ip {
			return true
		} else {
			log.Errorf("Old Session exist for subscriber %s with IP addr %s. Ending it.", imsi, sc.s.UeIpAddr)
			_ = s.EndSession(ctx, &store.Subscriber{Imsi: imsi})
		}
	}
	return false
}

func (s *sessionManager) CreateSesssion(ctx context.Context, sub *store.Subscriber, ns *store.Session, rxf *store.Flow, txf *store.Flow) error {

	sc := sessionCache{
		s:              ns,
		txCookie:       txf.Cookie,
		rxCookie:       rxf.Cookie,
		idleReportSent: false,
	}

	/* Get usage for before session creation */
	u, err := s.store.GetUsageByImsi(sub.Imsi)
	if err != nil {
		log.Errorf("Error getting usage for Imsi %s.Error: %s", sub.Imsi, err.Error())
		return err
	}
	sc.InitUsage = u.Data

	/* Add new data path */
	err = s.d.AddNewDataPath(sc.s.UeIpAddr, uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID),
		uint32(sc.s.TxMeterID.Rate), uint32(sc.s.RxMeterID.Rate), uint32(sc.s.RxMeterID.Burst),
		sc.rxCookie, sc.txCookie)
	if err != nil {
		/* TODO: mark session as failure for KPI purposes in DB */
		log.Errorf("Failed to add data path for Imsi %s. Error: %s", sub.Imsi, err.Error())
		return err
	}

	/* Add session to list */
	s.cache[sub.Imsi] = &sc

	/* Start session monitor thread */
	err = s.StartSessionMonitor(ctx, sub.Imsi)
	if err != nil {
		log.Errorf("Failed to start monitor for Imsi %s. Error: %s", sub.Imsi, err.Error())
		return err
	}

	return nil
}

func (s *sessionManager) EndAllSessions() error {
	for imsi, session := range s.cache {
		err := s.EndSession(context.Background(), &store.Subscriber{Imsi: imsi})
		if err != nil {
			log.Errorf("Failed to end session for Imsi %s.Error %s", imsi, err.Error())
		}
		log.Infof("Ending session %+v for Imsi %s.", session, imsi)
	}
	return nil
}

func (s *sessionManager) EndSession(ctx context.Context, sub *store.Subscriber) error {

	sc, ok := s.cache[sub.Imsi]
	if !ok {
		log.Errorf("failed to find session for Imsi %s", sub.Imsi)
		return nil
	}

	/* Stop montioring*/
	err := s.StopSessionMonitor(ctx, sub.Imsi)
	if err != nil {
		log.Errorf("Failed to stop monitor for Imsi %s. Error: %s", sub.Imsi, err.Error())
		return err
	}

	_ = s.storeStats(sc.s.SubscriberID.Imsi, true)

	time.Sleep(1000 * time.Millisecond)

	/* Delete the UE Data path */
	err = s.d.DeleteDataPath(sc.s.UeIpAddr, uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID))
	if err != nil {
		log.Errorf("Failed to delete data path for Imsi %s. Error: %s", sub.Imsi, err.Error())
		/* TODO: Need to figure out way to stop traffic for UE
		Another poin is the command to stop session comes from the EPC whihc means the connection is dropped
		so it might be ok. TBU based on the test result of this case */
		return err
	}

	/* Sync data to cloud */
	_ = s.SendCDR(sub.Imsi)

	/* Update usage to DB */
	err = s.store.EndSession(sc.s)
	if err != nil {
		log.Warnf("Failed to update session to db store for Imsi %s. Error: %s", sub.Imsi, err.Error())
	}

	delete(s.cache, sub.Imsi)

	return nil

}

func (s *sessionManager) SendCDR(imsi string) error {
	sc := s.cache[imsi]
	log.Infof("[ SessionId %d ] Sending session CDR for subscriber %s  and IP address %s", sc.s.ID, imsi, sc.s.UeIpAddr)

	/* Sync data to cloud */
	c := store.PrepareCDR(sc.s)
	err := s.rc.PushCdr(c)
	if err != nil {
		log.Warnf("Failed to push cdr %+v to cloud for Imsi %s. Error: %s", c, imsi, err.Error())
		return err
	}

	return nil
}

func (s *sessionManager) StartSessionMonitor(ctx context.Context, imsi string) error {
	sc := s.cache[imsi]
	log.Infof("[SessionId %d ] Starting session monitor for subscriber %s and IP address %s", sc.s.ID, imsi, sc.s.UeIpAddr)

	sc.ctx, sc.cancel = context.WithCancel(context.Background())
	s.cache[imsi] = sc

	go s.sessionMonitorRoutine(sc.ctx, s.period, sc)

	return nil
}

func (s *sessionManager) StopSessionMonitor(ctx context.Context, imsi string) error {
	sc := s.cache[imsi]
	log.Infof("[SessionId %d ] Stop session monitor for subscriber %s and IP address %s", sc.s.ID, imsi, sc.s.UeIpAddr)

	sc.cancel()

	return nil
}

/* For now we are starting session for ach active session */
func (s *sessionManager) sessionMonitorRoutine(ctx context.Context, interval time.Duration, sc *sessionCache) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Perform your periodic task here
			log.Infof("[SessionId %d ] Stat Collection", sc.s.ID)
			_ = s.storeStats(sc.s.SubscriberID.Imsi, false)

		case <-ctx.Done():
			// Context canceled, exit the goroutine
			log.Infof("[SessionId %d ] Exiting montoring for subscriber %s", sc.s.ID, sc.s.SubscriberID.Imsi)
			return
		}
	}
}
