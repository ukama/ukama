package store

import "github.com/ukama/ukama/systems/common/uuid"

type PathType int

const (
	NONE    PathType = iota
	RX_PATH          = 1
	TX_PATH          = 2
)

type SessionState int

const (
	Unkown     SessionState = iota
	Active                  = 1
	Terminated              = 2 /* Done by timeout if no changes in RX and Tx for 600 sec */
	Completed               = 3
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
	RXBytes      uint64
	TotalBytes   uint64
	TXMeterId    Meter
	RXMeterId    Meter
	State        SessionState
}

type Meter struct {
	ID   int
	Rate uint64
	Type PathType
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
