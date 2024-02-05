package session

import (
	"context"
	"log"
	"time"

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

type sessionManager struct {
	period time.Duration `default:"1s"`
	store  *store.Store
	d      datapath.DataPath
	cache  []Cache
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

func (s *sessionManager) CreateSesssion(ctx context.Context, ip string, imsi string) {

	/* Add new data path */

	/* Start session monitor thread */

}
