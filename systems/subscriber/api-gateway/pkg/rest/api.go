package rest

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SimInfo struct {
	Iccid          string `json:"iccid" query:"iccid" binding:"required" validate:"required"`
	SimType        string `json:"simType" query:"simType" binding:"required" validate:"required"`
	Msidn          string `json:"msidn" query:"msidn" binding:"required" validate:"required"`
	SmDpAddress    string `json:"smdpAddress" query:"smdpAddress" binding:"required" validate:"required"`
	ActivationCode string `json:"activationCode" query:"activationCode" binding:"required" validate:"required"`
	QrCode         string `json:"qrcode" query:"qrcode" binding:"required" validate:"required"`
	IsPhysicalSim  bool   `json:"isPhysicalSim" query:"isPhysicalSim" binding:"required" validate:"required"`
}

type SimPoolStats struct {
	Total     uint64 `json:"count"`
	Available uint64 `json:"available"`
	Consumed  uint64 `json:"consumed"`
	Failed    uint64 `json:"failed"`
}

type SIM struct {
	Id                string    `json:"id" validate:"required"`
	SubscriberId      string    `json:"subscriberId" validate:"required"`
	Iccid             string    `json:"iccid" validate:"required"`
	Msidn             string    `json:"msidn" validate:"required"`
	Package           Package   `json:"package" validate:"required"`
	FirstActivatedOn  time.Time `json:"firstActivedOn" validate:"required"`
	LastActivationOn  time.Time `json:"lastActivationOn" validate:"required"`
	DeactivationCount uint64    `json:"deactivationCount" validate:"required"`
	IsPhysical        string    `json:"type" validate:"required"`
	SimType           string    `json:"simType" validate:"required"`
	ActivationCount   uint64    `json:"activationCount" validate:"required"`
	AllocatedAt       time.Time `json:"allocatedAt" validate:"required"`
	Status            string    `json:"status" validate:"required"`
}

type Package struct {
	PackageId string    `json:"packageId" validate:"required"`
	StartDate time.Time `json:"startDate" validate:"required"`
	EndDate   time.Time `json:"endDate" validate:"required"`
}

type Subscriber struct {
	SubscriberId          uuid.UUID `json:"subscriberId" validate:"required"`
	FirstName             string    `json:"firstName" validate:"required"`
	LastName              string    `json:"lastName" validate:"required"`
	Email                 string    `json:"email" validate:"email,required"`
	Phone                 string    `json:"phone" validate:"required"`
	DOB                   time.Time `json:"dob" validate:"required"`
	ProofOfIdentification string    `json:"proofOfId" validate:"required"`
	IdSerial              string    `json:"idSerial" validate:"required"`
	Address               string    `json:"address" validate:"required"`
	SimList               []SIM     `json:"sims" validate:"required"`
}

type SimPoolStatByTypeReq struct {
	SimType string `form:"simType" json:"simType" query:"simType" binding:"required" validate:"required"`
}

type SimPoolRemoveSimReq struct {
	Id []uint64 `form:"id" json:"id" query:"id" binding:"required" validate:"required"`
}

type SimPoolUploadSimReq struct {
}

type SimPoolAddSimReq struct {
	SimInfo []SimInfo
}

type SubscriberAddReq struct {
	FirstName             string                 `json:"firstName" validate:"required"`
	LastName              string                 `json:"lastName" validate:"required"`
	Email                 string                 `json:"email" validate:"required"`
	Phone                 string                 `json:"phone" validate:"required"`
	DOB                   *timestamppb.Timestamp `json:"dob" validate:"required"`
	ProofOfIdentification string                 `json:"proofOfId" validate:"required"`
	IdSerial              string                 `json:"idSerial" validate:"required"`
	Address               string                 `json:"address" validate:"required"`
}

type SubscriberGetReq struct {
	SubscriberId string `form:"subscriberId" json:"subscriberId" path:"subscriberId" binding:"required" validate:"required"`
}

type SubscriberDeleteReq struct {
	SubscriberId string `form:"subscriberId" json:"subscriberId" path:"subscriberId" binding:"required" validate:"required"`
}

type SubscriberByNetworkReq struct {
	NetworkId string `form:"networkId" json:"networkId" path:"networkId" binding:"required" validate:"required"`
}

type SubscriberUpdateReq struct {
	SubscriberId          string `json:"subscriberId" validate:"required"`
	Email                 string `json:"email" validate:"required"`
	Phone                 string `json:"phone" validate:"required"`
	Address               string `json:"address" validate:"required"`
	ProofOfIdentification string `json:"proofOfIdentification" validate:"required"`
	IdSerial              string `json:"idSerial" validate:"required"`
}

type SimListReq struct {
	NetworkId uuid.UUID `path:"networkId" validate:"required"`
}

type SimListResp struct {
	Subscribers []Subscriber `json:"subscribers"`
}

type AllocateSimReq struct {
	SubscriberId string `form:"subscriberId" json:"subscriberId" path:"subscriberId" binding:"required" validate:"required`
	SimToken     string `form:"simToken" json:"simToken" path:"simToken" binding:"required" validate:"required`
	PackageId    string `form:"packageId" json:"packageId" path:"packageId" binding:"required" validate:"required`
	NetworkId    string `form:"networkId" json:"networkId" path:"networkId" binding:"required" validate:"required`
}
type SimReq struct {
	SimId string `form:"sim_id" json:"sim_id" query:"sim_id" binding:"required" validate:"required`
}

type GetSimsBySubReq struct {
	SubscriberId string `form:"subscriber_id" json:"subscriber_id" query:"subscriber_id" binding:"required" validate:"required`
}
type AddPkgToSimReq struct {
	SimId        string                 `form:"simId" json:"simId" path:"simId" binding:"required" validate:"required`
	SubscriberId string                 `form:"subscriberId" json:"subscriberId" path:"subscriberId" binding:"required" validate:"required`
	PackageId    string                 `form:"packageId" json:"packageId" path:"packageId" binding:"required" validate:"required`
	StartDate    *timestamppb.Timestamp `form:"startDate" json:"startDate" path:"startDate" binding:"required" validate:"required`
}

type RemovePkgFromSimReq struct {
	SimId     string `form:"simId" json:"simId" path:"simId" binding:"required" validate:"required`
	PackageId string `form:"packageId" json:"packageId" path:"packageId" binding:"required" validate:"required`
}
