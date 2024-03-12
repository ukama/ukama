/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
 
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
	SessionUnknown    SessionState = iota
	SessionActive                  = 1
	SessionTerminated              = 2 /* Done by timeout if no changes in RX and Tx for 600 sec */
	SessionCompleted               = 3
)

const (
	SessionSyncUnknown   SessionSync = iota
	SessionSyncPending               = 1 /* Session in progress */
	SessionSyncReady                 = 2 /* Session completed now ready for sync */
	SessionSyncCompleted             = 3 /* Sync is compeleted */
)

type Subscriber struct {
	ID        int
	Imsi      string
	PolicyID  Policy
	ReRouteID ReRoute
}

type Policy struct {
	ID        uuid.UUID
	Data      uint64
	Dlbr      uint64
	Ulbr      uint64
	Burst     uint64
	StartTime int64
	EndTime   int64
}

type Usage struct {
	ID           int
	SubscriberID Subscriber
	Updatedat    uint64
	Data         uint64
}

type Session struct {
	ID           int
	SubscriberID Subscriber
	PolicyID     Policy
	ApnName      string
	UeIpAddr     string
	StartTime    uint64
	EndTime      uint64
	TxBytes      uint64
	RxBytes      uint64
	TotalBytes   uint64
	TxMeterID    Meter
	RxMeterID    Meter
	State        SessionState
	Sync         SessionSync
	UpdatedAt    uint64
}

type Meter struct {
	ID    int
	Rate  uint64
	Burst uint64
	Type  PathType
}

type Flow struct {
	ID        int
	Cookie    uint64
	Tableid   uint64
	Priority  uint64
	UeIpAddr  string
	ReRouting ReRoute
	MeterID   Meter
}

type ReRoute struct {
	ID     int
	IpAddr string
}

func (s SessionSync) String() string {
	switch s {
	case SessionSyncPending:
		return "SessionSyncPending"
	case SessionSyncReady:
		return "SessionSyncReady"
	case SessionSyncCompleted:
		return "SessionSyncCompleted"
	default:
		return "SessionSyncUnkown"
	}
}

func ParseSessionSync(s string) SessionSync {
	switch s {
	case "SessionSyncPending":
		return SessionSyncPending
	case "SessionSyncReady":
		return SessionSyncReady
	case "SessionSyncCompleted":
		return SessionSyncCompleted
	default:
		return SessionSyncUnknown
	}
}

func (s SessionState) String() string {
	switch s {
	case SessionActive:
		return "SessionActive"
	case SessionTerminated:
		return "SessionTerminated"
	case SessionCompleted:
		return "SessionCompleted"
	default:
		return "SessionUnknown"
	}
}

func ParseSessionState(s string) SessionState {
	switch s {
	case "SessionActive":
		return SessionActive
	case "SessionTerminated":
		return SessionTerminated
	case "SessionCompleted":
		return SessionCompleted
	default:
		return SessionUnknown
	}
}
