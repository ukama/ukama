package session

import (
	"context"
	"time"

	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/datapath"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/store"
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
	d := datapath.InitDataPath(name, ip, netType)

	s := &sessionManager{
		d:      d,
		store:  store,
		period: period,
	}

	return s
}

func (s *sessionManager) CreateSesssion(ctx context.Context, ip string, imsi string) {

}
