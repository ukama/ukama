package store

import "github.com/ukama/ukama/systems/common/uuid"

type PathType int

const (
	NONE    PathType = iota
	RX_PATH          = 1
	TX_PATH          = 2
)

type SessionState int
type SessionSync int

const (
	SessionUnkown     SessionState = iota
	SessionActive                  = 1
	SessionTerminated              = 2 /* Done by timeout if no changes in RX and Tx for 600 sec */
	SessionCompleted               = 3
)

const (
	SessionSyncUnkown    SessionSync = iota
	SessionSyncPending               = 1 /* Session in progress */
	SessionSyncReady                 = 2 /* Session completed now ready for sync */
	SessionSyncCompleted             = 3 /* Sync is compeleted */
)

type Subscriber struct {
	ID       uuid.UUID
	Imsi     string
	PolicyID Policy
}

type Policy struct {
	ID        uuid.UUID
	Data      uint64
	Dlbr      uint64
	Ulbr      uint64
	StartTime int64
	EndTime   int64
}

type Usage struct {
	ID           int
	SubscriberID Subscriber
	Epoch        uint64
	Data         uint64
}

type Session struct {
	ID           int
	SusbcriberID Subscriber
	ApnName      string
	UeIpaddr     string
	StartTime    uint64
	EndTime      uint64
	TxBytes      uint64
	RxBytes      uint64
	TotalBytes   uint64
	TXMeterId    Meter
	RXMeterId    Meter
	State        SessionState
	Sync         SessionSync
}

type Meter struct {
	ID        int
	Rate      uint64
	BurstSize uint64
	Type      PathType
}

type Flow struct {
	ID        int
	Cookie    uint64
	Table     uint64
	Priority  uint64
	UeIpaddr  string
	ReRouting ReRoute
	MeterID   Meter
}

type ReRoute struct {
	ID     int
	Ipaddr string
}
