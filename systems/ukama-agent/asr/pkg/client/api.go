package client

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type SimCardInfo struct {
	Imsi           string
	Iccid          string
	Op             []byte
	Amf            []byte
	Key            []byte
	AlgoType       uint32
	UeDlAmbrBps    uint32
	UeUlAmbrBps    uint32
	Sqn            uint64
	CsgIdPrsent    bool
	CsgId          uint32
	DefaultApnName string
}

type NetworkInfo struct {
	NetworkId     string    `json:"id"`
	Name          string    `json:"name"`
	OrgId         string    `json:"org_id"`
	IsDeactivated bool      `json:"is_deactivated"`
	CreatedAt     time.Time `json:"created_at"`
}


