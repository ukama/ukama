package rest

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SimInfo struct {
	Iccid          string `json:"iccid" binding:"required" validate:"required"`
	SimType        string `json:"sim_type" binding:"required" validate:"required"`
	Msidn          string `json:"msidn"  binding:"required" validate:"required"`
	SmDpAddress    string `json:"smdp_address"  binding:"required" validate:"required"`
	ActivationCode string `json:"activation_code"  binding:"required" validate:"required"`
	QrCode         string `json:"qr_code"  binding:"required" validate:"required"`
	IsPhysicalSim  bool   `json:"is_physical_sim" binding:"required" validate:"required"`
}

type SimPoolStats struct {
	Total     uint64 `json:"count"`
	Available uint64 `json:"available"`
	Consumed  uint64 `json:"consumed"`
	Failed    uint64 `json:"failed"`
}

type SIM struct {
	Id                string    `json:"id" validate:"required"`
	SubscriberId      string    `json:"subscriber_id" validate:"required"`
	Iccid             string    `json:"iccid" validate:"required"`
	Msidn             string    `json:"msidn" validate:"required"`
	Package           Package   `json:"package" validate:"required"`
	FirstActivatedOn  time.Time `json:"first_actived_on" validate:"required"`
	LastActivationOn  time.Time `json:"last_activation_on" validate:"required"`
	DeactivationCount uint64    `json:"deactivation_count" validate:"required"`
	IsPhysical        string    `json:"type" validate:"required"`
	SimType           string    `json:"sim_type" validate:"required"`
	ActivationCount   uint64    `json:"activation_count" validate:"required"`
	AllocatedAt       time.Time `json:"allocated_at" validate:"required"`
	Status            string    `json:"status" validate:"required"`
}

type Package struct {
	PackageId string    `json:"package_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

type Subscriber struct {
	SubscriberId          uuid.UUID `json:"subscriber_id" validate:"required"`
	FirstName             string    `json:"first_name" validate:"required"`
	LastName              string    `json:"last_name" validate:"required"`
	Email                 string    `json:"email" validate:"email,required"`
	Phone                 string    `json:"phone" validate:"required"`
	DOB                   time.Time `json:"dob" validate:"required"`
	ProofOfIdentification string    `json:"proof_of_Identification" validate:"required"`
	IdSerial              string    `json:"id_serial" validate:"required"`
	Address               string    `json:"address" validate:"required"`
	SimList               []SIM     `json:"sims" validate:"required"`
}

type SimByIccidReq struct {
	Iccid string `path:"iccid" validate:"required"`
}

type SimPoolStatByTypeReq struct {
	SimType string `form:"sim_type" json:"sim_type" path:"sim_type" binding:"required" validate:"required"`
}

type SimPoolRemoveSimReq struct {
	Id []uint64 `form:"id" json:"id" path:"sim_id" binding:"required" validate:"required"`
}

type SimPoolUploadSimReq struct {
	SimType string `form:"sim_type" json:"sim_type" binding:"required" validate:"required"`
}

type SimPoolAddSimReq struct {
	SimInfo []SimInfo
}

type SubscriberAddReq struct {
	FirstName             string `json:"first_name" validate:"required"`
	LastName              string `json:"last_name" validate:"required"`
	Email                 string `json:"email" validate:"required"`
	Phone                 string `json:"phone" validate:"required"`
	DOB                   string `json:"dob" validate:"required"`
	ProofOfIdentification string `json:"proof_of_Identification" validate:"required"`
	IdSerial              string `json:"id_serial" validate:"required"`
	Address               string `json:"address" validate:"required"`
	NetworkID             string `json:"network_id" validate:"required"`
	Gender                string `json:"gender" validate:"required"`
	OrgID                 string `json:"org_id" validate:"required"`
}

type SubscriberGetReq struct {
	SubscriberId string `form:"subscriber_id" json:"subscriber_id" path:"subscriber_id" binding:"required" validate:"required"`
}

type SubscriberDeleteReq struct {
	SubscriberId string `form:"subscriber_id" json:"subscriber_id" path:"subscriber_id" binding:"required" validate:"required"`
}

type SubscriberByNetworkReq struct {
	NetworkId string `form:"network_id" json:"network_id" path:"network_id" binding:"required" validate:"required"`
}

type SubscriberUpdateReq struct {
	SubscriberId          string `json:"subscriber_id" validate:"required"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	Address               string `json:"address"`
	ProofOfIdentification string `json:"proof_of_Identification"`
	IdSerial              string `json:"id_serial"`
}

type SimListReq struct {
	NetworkId uuid.UUID `form:"network_id" json:"network_id" path:"network_id" binding:"required" validate:"required"`
}

type SimListResp struct {
	Subscribers []Subscriber `json:"subscribers"`
}

type AllocateSimReq struct {
	SubscriberId string `json:"subscriber_id" validate:"required`
	SimToken     string `json:"sim_token" validate:"required`
	PackageId    string `json:"package_id" validate:"required`
	NetworkId    string `json:"network_id" validate:"required`
}

type SetActivePackageForSimReq struct {
	SimId     string `json:"sim_id" validate:"required`
	PackageId string `json:"package_id" validate:"required`
}
type SimReq struct {
	SimId string `form:"sim_id" json:"sim_id" path:"sim_id" binding:"required" validate:"required`
}

type SimByNetworkReq struct {
	NetworkId string `form:"network_id" json:"network_id" path:"network_id" binding:"required" validate:"required"`
}

type ActivateDeactivateSimReq struct {
	SimId  string `json:"sim_id" binding:"required" validate:"required`
	Status string `json:"status" binding:"required" validate:"required`
}

type GetSimsBySubReq struct {
	SubscriberId string `form:"subscriber_id" json:"subscriber_id" path:"subscriber_id" binding:"required" validate:"required`
}
type AddPkgToSimReq struct {
	SimId     string                 `json:"sim_id" validate:"required"`
	PackageId string                 `json:"package_id" validate:"required"`
	StartDate *timestamppb.Timestamp `json:"start_date" validate:"required"`
}

type RemovePkgFromSimReq struct {
	SimId     string `form:"sim_id" json:"sim_id" path:"sim_id" binding:"required" validate:"required`
	PackageId string `form:"package_id" json:"package_id" path:"package_id" binding:"required" validate:"required`
}
