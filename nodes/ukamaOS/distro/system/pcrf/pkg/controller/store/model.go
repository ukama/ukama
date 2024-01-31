package store

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
	ID     int
	Imsi   string
	Policy []Policy
}

type Policy struct {
	ID   int /* Define by package ID */
	Data uint64
	Dlbr uint64
	Ulbr uint64
}

type Usage struct {
	ID           int
	SubscriberID Subscriber
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
