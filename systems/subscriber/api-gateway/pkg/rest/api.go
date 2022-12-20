package rest

import (
	"time"

	"github.com/google/uuid"
)

type SimInfo struct {
	Iccid          string `json:"iccid" validate:"required"`
	SimType        string `json:"simType" validate:"required"`
	Msidn          string `json:"msidn" validate:"required"`
	SmDpAddress    string `json:"smdpAddress" validate:"required"`
	ActivationCode string `json:"activationCode" validate:"required"`
	QrCode         string `json:"qrcode" validate:"required"`
	
}

type SimPoolStats struct {
	Total     uint64 `json:"count"`
	Available uint64 `json:"available"`
	Consumed  uint64 `json:"consumed"`
	Failed    uint64 `json:"failed"`
}

type SIM struct {
	SimId             uuid.UUID `json:"simId" validate:"required"`
	SubscriberId      uuid.UUID `json:"packageId" validate:"required"`
	Iccid             string    `json:"iccid" validate:"required"`
	SimType           string    `json:"simType" validate:"required"`
	SimManager        string    `json:"simManager" validate:"required"`
	OrgId             uuid.UUID `json:"orgId" validate:"required"`
	NetworkId         uuid.UUID `json:"networkId" validate:"required"`
	ActivationCount   uint64    `json:"activationCount" validate:"required"`
	DeactivationCount uint64    `json:"DeactivationCount" validate:"required"`
	FirstActivatedOn  time.Time `json:"firstActivedOn" validate:"required"`
	LastActivationOn  time.Time `json:"lastActivationOn" validate:"required"`
	Msidn             string    `json:"msidn" validate:"required"`
	State             string    `json:"state" validate:"required"`
	Package           []Package `json:"packages" validate:"required"`
	ActivePackageId   uuid.UUID `json:"activePackageId" validate:"required"`
}

type Package struct {
	PackageId uuid.UUID `json:"packageId" validate:"required"`
	StartDate time.Time `json:"startDate" validate:"required"`
	EndDate   time.Time `json:"endDate" validate:"required"`
}

type Subscriber struct {
	SubscriberId          uuid.UUID `json:"subscriberId" validate:"required"`
	Name                  string    `json:"name" validate:"required"`
	EMail                 string    `json:"email" validate:"email,required"`
	PhoneNumber           string    `json:"phone" validate:"required"`
	DOB                   time.Time `json:"dob" validate:"required"`
	ProofOfIdentification string    `json:"proofOfId" validate:"required"`
	IdSerial              string    `json:"idSerial" validate:"required"`
	Address               string    `json:"address" validate:"required"`
	SimList               []SIM     `json:"sims" validate:"required"`
}

type SimPoolStatByTypeReq struct {
	SimType string `json:"simType"`
}

type SimPoolRemoveSimReq struct {
	Iccids []string `json:"iccids"`
}

type SimPoolRemoveSimResp struct {
	Iccids []string `json:"iccids"`
}
type SimPoolUploadSimReq struct {
}

type SimPoolUploadSimResp struct {
	Iccids []string `json:"iccids"`
}

type SimPoolAddSimReq struct {
	SimInfo []SimInfo
}

type SimPoolAddSimResp struct {
	Iccids []string `json:"iccids"`
}

type SubscriberAddReq struct {
	Name                  string    `json:"name" validate:"required"`
	EMail                 string    `json:"email" validate:"required"`
	PhoneNumber           string    `json:"phone" validate:"required"`
	DOB                   time.Time `json:"dob" validate:"required"`
	ProofOfIdentification string    `json:"proofOfId" validate:"required"`
	IdSerial              string    `json:"idSerial" validate:"required"`
	Address               string    `json:"address" validate:"required"`
}

type SubscriberAddResp struct {
	SubscriberId uuid.UUID `path:"subscriberId" validate:"required"`
}

type SubscriberGetReq struct {
	SubscriberId uuid.UUID `json:"subscriberId" validate:"required"`
}

type SubscriberGetResp struct {
	Subscriber
}

type SubscriberDeleteReq struct {
	SubscriberId uuid.UUID `path:"subscriberId" validate:"required"`
}

type SubscriberListReq struct {
	NetworkId uuid.UUID `path:"networkId" validate:"required"`
}

type SubscriberListResp struct {
	Subscribers []Subscriber `json:"subscribers"`
}
type SubscriberSimAllocateReq struct {
	SubscriberId uuid.UUID `path:"subscriberId" validate:"required"`
	NetworkId    uuid.UUID `json:"networkId" validate:"required"`
	SimType      string    `json:"type" validate:"required"`
	Token        string    `json:"token" validate:"required"`
	PlanId       uuid.UUID `json:"planId" validate:"required"`
}

type SubscriberSimAllocateResp struct {
	SIM
}

type SubscriberSimUpdateStateReq struct {
	SubscriberId uuid.UUID `path:"subscriberId" validate:"required"`
	SimId        uuid.UUID `path:"simId" validate:"required"`
	State        string    `json:"state" validate:"eq=inactive|eq=INACTIVE|eq=active|eq=ACTIVE,required" `
}

type SubscriberSimDeleteReq struct {
	SubscriberId uuid.UUID `path:"subscriberId" validate:"required"`
	SimId        uuid.UUID `path:"simId" validate:"required"`
}

type SubscriberSimReadReq struct {
	SubscriberId uuid.UUID `path:"subscriberId" validate:"required"`
	SimId        uuid.UUID `path:"simId" validate:"required"`
}

type SubscriberSimReadResp struct {
	SIM
}
type SubscriberSimAddPackageReq struct {
	SubscriberId uuid.UUID `path:"subscriberId" validate:"required"`
	SimId        uuid.UUID `path:"simId" validate:"required"`
	PackageId    uuid.UUID `json:"packageId" validate:"required"`
	StartDate    time.Time `json:"startDate" validate:"required"`
}

type SubscriberSimRemovePackageReq struct {
	SubscriberId uuid.UUID `path:"subscriberId" validate:"required"`
	SimId        uuid.UUID `path:"simId" validate:"required"`
	PackageId    uuid.UUID `json:"packageId" validate:"required"`
}

type SimListReq struct {
	NetworkId uuid.UUID `path:"networkId" validate:"required"`
}

type SimListResp struct {
	Subscribers []Subscriber `json:"subscribers"`
}
