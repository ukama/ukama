package session

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/datapath"
)

type Cache struct {
	imsi    string
	ip      string
	rxMeter uint32
	txMeter uint32
}

type sessionCache struct {
	s        *store.Session
	txCookie uint64
	rxCookie uint64
	cancel   context.CancelFunc
	ctx      context.Context
}
type sessionManager struct {
	period time.Duration `default:"2s"`
	store  *store.Store
	d      datapath.DataPath
	cache  map[string]*sessionCache
	rc     client.RemoteController
}

type SessionManager interface {
	CreateSesssion(ctx context.Context, sub *store.Subscriber, ns *store.Session, rxf *store.Flow, txf *store.Flow) error
	EndSesssion(ctx context.Context, sub *store.Subscriber) error
}

func NewSessionManager(rc client.RemoteController, store *store.Store, name, ip, netType string, period time.Duration) *sessionManager {
	d, err := datapath.InitDataPath(name, ip, netType)
	if err != nil {
		log.Fatalf("error initializing session manager. Error: %s", err.Error())
	}

	s := &sessionManager{
		d:      d,
		store:  store,
		period: period,
		rc:     rc,
	}

	return s
}

func (s *sessionManager) storeStats(imsi string, lastStats bool) error {
	var err error
	sc := s.cache[imsi]
	/* Read sats */
	sc.s.RxBytes, _, sc.s.TxBytes, _, err = s.d.DataPathStats(sc.rxCookie, sc.txCookie)
	if err != nil {
		log.Errorf("[SessionId %d ] Failed to read final stats for data path of Imsi %s. Error: %s", sc.s.ID, sc.s.SubscriberID.Imsi, err.Error())
		return err
	}

	/* Update to DB */
	if lastStats {
		err = s.store.EndSession(sc.s)
		if err != nil {
			log.Warnf("[SessionId %d ] Failed to update last session usage to db store for Imsi %s. Error: %s", sc.s.ID, sc.s.SubscriberID.Imsi, err.Error())
		}
	} else {

		err = s.store.UpdateSessionUsage(sc.s)
		if err != nil {
			log.Warnf("[SessionId %d ] Failed to update session usage to db store for Imsi %s. Error: %s", sc.s.ID, sc.s.SubscriberID.Imsi, err.Error())
		}

	}

	/* Update session */
	s.cache[imsi] = sc

	return err
}

func (s *sessionManager) CreateSesssion(ctx context.Context, sub *store.Subscriber, ns *store.Session, rxf *store.Flow, txf *store.Flow) error {

	sc := sessionCache{
		s:        ns,
		txCookie: txf.Cookie,
		rxCookie: rxf.Cookie,
	}

	/* Add new data path */
	err := s.d.AddNewDataPath(sc.s.UeIpAddr, uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID),
		uint32(sc.s.TxMeterID.Rate), uint32(sc.s.RxMeterID.Rate), uint32(sc.s.RxMeterID.Burst),
		sc.rxCookie, sc.txCookie)
	if err != nil {
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

func (s *sessionManager) EndSesssion(ctx context.Context, sub *store.Subscriber) error {

	sc := s.cache[sub.Imsi]

	/* Stop montioring*/
	err := s.StopSessionMonitor(ctx, sub.Imsi)
	if err != nil {
		log.Errorf("Failed to stop monitor for Imsi %s. Error: %s", sub.Imsi, err.Error())
		return err
	}

	/* Read sats */
	sc.s.RxBytes, _, sc.s.TxBytes, _, err = s.d.DataPathStats(sc.rxCookie, sc.txCookie)
	if err != nil {
		log.Errorf("Failed to read final stats for data path of Imsi %s. Error: %s", sub.Imsi, err.Error())
		return err
	}

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
	c := store.PrepareCDR(sc.s)
	err = s.rc.PushCdr(c)
	if err != nil {
		log.Warnf("Failed to push cdr to cloud for Imsi %s. Error: %s", sub.Imsi, err.Error())
	}

	/* Update sync state */
	sc.s.Sync = store.SessionSyncReady

	/* Update to DB */
	err = s.store.UpdateSessionEndUsage(sc.s)
	if err != nil {
		log.Warnf("Failed to update session to db store for Imsi %s. Error: %s", sub.Imsi, err.Error())
	}

	return nil

}

func (s *sessionManager) StartSessionMonitor(ctx context.Context, imsi string) error {
	sc := s.cache[imsi]
	log.Infof("[SessionId %d ] Starting session monitor for subscriber %s and IP address %s", sc.s.ID, imsi, sc.s.UeIpAddr)

	sc.ctx, sc.cancel = context.WithCancel(context.Background())
	s.cache[imsi] = sc

	go s.sessionMonitorRoutine(ctx, s.period, sc)

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
			_ = s.storeStats(sc.s.SubscriberID.Imsi, true)
			log.Infof("[SessionId %d ] Exiting montoring for subscriber %s", sc.s.ID, sc.s.SubscriberID.Imsi)
			return
		}
	}
}
