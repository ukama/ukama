package client

import uuid "github.com/satori/go.uuid"

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
	Name      string
	NetworkId uuid.UUID
	OrgID     uuid.UUID
	OrgName   string
}

type PolicyControlSimInfo struct {
	Imsi      string `path:"imsi" validate:"required" json:"-"`
	Iccid     string
	PackageId uuid.UUID
	NetworkId uuid.UUID
	Visitor   bool
}

type PolicyControlSimPackageUpdate struct {
	Imsi      string `path:"imsi" validate:"required" json:"-"`
	PackageId uuid.UUID
}
