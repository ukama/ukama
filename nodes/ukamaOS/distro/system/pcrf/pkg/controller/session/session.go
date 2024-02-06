package session

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/datapath"
)

type Cache struct {
	imsi     string
	ip       string
	rxMeter  uint32
	txMeter  uint32
	txCookie uint64
	rxCookie uint64
}

type sessionCache struct {
	s      *store.Session
	cancel context.CancelFunc
	ctx    context.Context
}
type sessionManager struct {
	period time.Duration `default:"1s"`
	store  *store.Store
	d      datapath.DataPath
	cache  map[string]*sessionCache
}

type SessionManger interface {
}

func NewSessionManager(store *store.Store, name, ip, netType string, period time.Duration) *sessionManager {
	d, err := datapath.InitDataPath(name, ip, netType)
	if err != nil {
		log.Fatalf("error initializing session manager. Error: %s", err.Error())
	}

	s := &sessionManager{
		d:      d,
		store:  store,
		period: period,
	}

	return s
}

func (s *sessionManager) CreateSesssion(ctx context.Context, sub *store.Subscriber) error {

	/* Add new data path */

	/* Start session monitor thread */

	/* Add session to list */

	return nil
}

func (s *sessionManager) DeleteSesssion(ctx context.Context, sub *store.Subscriber) error {

	/* Stop montioring*/

	/* Delete the UE Data path */

	/* Sync data to cloud */

	/* Update sync state */

	return nil

}

func (s *sessionManager) StartSessionMonitor(ctx context.Context, session *store.Session) error {
	log.Infof("[SessionId %d ] Starting session monitor for subscriber %s and IP address %s", session.ID, session.SusbcriberID.Imsi, session.UeIpaddr)
	sc := s.cache[session.SusbcriberID.Imsi]
	sc.ctx, sc.cancel = context.WithCancel(context.Background())
	s.cache[session.SusbcriberID.Imsi] = sc

	go s.sessionMonitorRoutine(ctx, s.period)

	return nil
}

func (s *sessionManager) StopSessionMonitor(ctx context.Context, session *store.Session) error {
	log.Infof("[SessionId %d ] Stop session monitor for subscriber %s and IP address %s", session.ID, session.SusbcriberID.Imsi, session.UeIpaddr)
	sc := s.cache[session.SusbcriberID.Imsi]

	sc.cancel()

	return nil
}

func (s *sessionManager) sessionMonitorRoutine(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Perform your periodic task here
			log.Infof("[SessionId %d ] Stat Collection")
		case <-ctx.Done():
			// Context canceled, exit the goroutine
			log.Infof("[SessionId %d ] Exiting montoring ")
			return
		}
	}
}
