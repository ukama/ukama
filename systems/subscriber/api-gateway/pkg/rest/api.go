package rest

import (
	"time"
)

type SimPoolStats struct {
	Count     int32 `json:"count"`
	Available int32 `json:"available"`
	Consumed  int32 `json:"consumed"`
	Failed    int32 `json:"failed"`
}

type SIM struct {
}

type SimType struct {
}

type SimPoolStatByTypeReq struct {
	stype SimType `json:"simType`
}

type SimPoolRemoveSimReq struct {
	Sims []UUID `json:"sims"`
}

type SimPoolUplaodSimReq struct {
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
	SubscriberId UUID `path:"subscriberId" validate:"required"`
}

type SubscriberGetReq struct {
	SubscriberId UUID `json:"subscriberId" validate:"required"`
}

type Subscriber struct {
	SubscriberId          UUID      `json:"subscriberId" validate:"required"`
	Name                  string    `json:"name" validate:"required"`
	EMail                 string    `json:"email" validate:"email,required"`
	PhoneNumber           string    `json:"phone" validate:"required"`
	DOB                   time.Time `json:"dob" validate:"required"`
	ProofOfIdentification string    `json:"proofOfId" validate:"required"`
	IdSerial              string    `json:"idSerial" validate:"required"`
	Address               string    `json:"address" validate:"required"`
	SimList               []SIM     `json:"sims" validate:"required"`
}

type SubscriberGetResp struct {
	Subscriber
}

type SubscriberDeleteReq struct {
	SubscriberId UUID `path:"subscriberId" validate:"required"`
}

type SubscriberListReq struct {
	NetworkId UUID `path:"networkId" validate:"required"`
}

type SubscriberListResp struct {
	Subscribers []Subscriber `json:"subscribers"`
}
type SubscriberSimAllocateReq struct {
	SubscriberId UUID    `path:"subscriberId" validate:"required"`
	NetworkId    UUID    `json:"networkId" validate:"required"`
	SType        SimType `json:"type" validate:"required"`
	Token        string  `json:"token" validate:"required"`
	PlanId       UUID    `json:"planId" validate:"required"`
}

type SubscriberSimAllocateResp struct {
	SIM
}

type SubscriberSimUpdateStateReq struct {
	SubscriberId UUID   `path:"subscriberId" validate:"required"`
	SimId        UUID   `path:"simId" validate:"required"`
	State        string `json:"state" validate:"eq=inactive|eq=INACTIVE|eq=active|eq=ACTIVE,required" `
}

type SubscriberSimDeleteReq struct {
	SubscriberId UUID `path:"subscriberId" validate:"required"`
	SimId        UUID `path:"simId" validate:"required"`
}

type SubscriberSimReadReq struct {
	SIM
}

type SubscriberSimAddPackageReq struct {
	SubscriberId UUID      `path:"subscriberId" validate:"required"`
	SimId        UUID      `path:"simId" validate:"required"`
	PackageID    UUID      `json: "packageId" validate:"required"`
	StartDate    time.Time `json: "startDate" validate:"required"`
}

type SubscriberSimRemovePackageReq struct {
	SubscriberId UUID `path:"subscriberId" validate:"required"`
	SimId        UUID `path:"simId" validate:"required"`
	PackageID    UUID `json: "packageId" validate:"required"`
}
