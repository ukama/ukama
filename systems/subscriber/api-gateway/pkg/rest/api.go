package rest

import (
	"time"

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
	SubscriberId          string `json:"subscriber_id" validate:"required"`
	FirstName             string `json:"first_name" validate:"required"`
	LastName              string `json:"last_name" validate:"required"`
	Email                 string `json:"email" validate:"email,required"`
	Phone                 string `json:"phone" validate:"required"`
	Dob                   string `json:"dob" validate:"required"`
	ProofOfIdentification string `json:"proof_of_Identification" validate:"required"`
	IdSerial              string `json:"id_serial" validate:"required"`
	Address               string `json:"address" validate:"required"`
	SimList               []SIM  `json:"sims" validate:"required"`
}

type SimByIccidReq struct {
	Iccid string `example:"8910300000003540855" path:"iccid" validate:"required"`
}

type SimPoolStatByTypeReq struct {
	SimType string `example:"test" form:"sim_type" json:"sim_type" path:"sim_type" binding:"required" validate:"required"`
}

type SimPoolRemoveSimReq struct {
	Id []uint64 `example:"[1]" form:"id" json:"id" path:"sim_id" binding:"required" validate:"required"`
}

type SimPoolUploadSimReq struct {
	SimType string `example:"test" form:"sim_type" json:"sim_type" binding:"required" validate:"required"`
	Data    string `example:"SUNDSUQsTVNJU0ROLFNtRHBBZGRyZXNzLEFjdGl2YXRpb25Db2RlLElzUGh5c2ljYWwsUXJDb2RlCjg5MTAzMDAwMDAwMDM1NDA4NTUsODgwMTcwMTI0ODQ3NTcxLDEwMDEuOS4wLjAuMSwxMDEwLFRSVUUsNDU5MDgxYQo4OTEwMzAwMDAwMDAzNTQwODQ1LDg4MDE3MDEyNDg0NzU3MiwxMDAxLjkuMC4wLjIsMTAxMCxUUlVFLDQ1OTA4MWIKODkxMDMwMDAwMDAwMzU0MDgzNSw4ODAxNzAxMjQ4NDc1NzMsMTAwMS45LjAuMC4zLDEwMTAsVFJVRSw0NTkwODFj" type:"array" format:"byte" form:"data" json:"data" binding:"required" validate:"required"`
}

type SimPoolAddSimReq struct {
	SimInfo []SimInfo `form:"sim_info" json:"sim_info" binding:"required"`
}

type SubscriberAddReq struct {
	FirstName             string `example:"John" json:"first_name" validate:"required"`
	LastName              string `example:"Doe" json:"last_name" validate:"required"`
	Email                 string `example:"john@example.com" json:"email" validate:"required"`
	Phone                 string `example:"4151231234" json:"phone" validate:"required"`
	Dob                   string `example:"Mon, 02 Jan 2006 15:04:05 MST" json:"dob" validate:"required"`
	ProofOfIdentification string `example:"passport" json:"proof_of_Identification" validate:"required"`
	IdSerial              string `example:"123456789" json:"id_serial" validate:"required"`
	Address               string `example:"Mr John Smith. 132, My Street, Kingston, New York 12401" json:"address" validate:"required"`
	NetworkId             string `example:"{{NetworkUUID}}" json:"network_id" validate:"required"`
	Gender                string `example:"male" json:"gender" validate:"required"`
	OrgId                 string `example:"{{OrgUUID}}" json:"org_id" validate:"required"`
}

type SubscriberGetReq struct {
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber_id" path:"subscriber_id" binding:"required" validate:"required"`
}

type SubscriberDeleteReq struct {
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber_id" path:"subscriber_id" binding:"required" validate:"required"`
}

type SubscriberByNetworkReq struct {
	NetworkId string `example:"{{NetworkUUID}}" form:"network_id" json:"network_id" path:"network_id" binding:"required" validate:"required"`
}

type SubscriberUpdateReq struct {
	SubscriberId          string `example:"{{SubscriberUUID}}" path:"subscriber_id" validate:"required"`
	Email                 string `example:"test@example.com" json:"email"`
	Phone                 string `example:"4151231234" json:"phone"`
	Address               string `example:"Mr John Smith. 132, My Street, Kingston, New York 12401" json:"address"`
	ProofOfIdentification string `example:"passport" json:"proof_of_Identification"`
	IdSerial              string `example:"123456789" json:"id_serial"`
}

type SimListReq struct {
	NetworkId string `example:"{{NetworkUUID}}" form:"network_id" json:"network_id" path:"network_id" binding:"required" validate:"required"`
}

type AllocateSimReq struct {
	SubscriberId string `example:"{{SubscriberUUID}}" json:"subscriber_id" validate:"required"`
	SimToken     string `example:"pj/9A5Hk8VkkZxOJyu0+9fWs7J6HCOjhmD5jEsIvOfZqmFFFMyStgC3Va4l1b6I5+2ibKOsJjR9KGug=" json:"sim_token"`
	PackageId    string `example:"{{PackageUUID}}" json:"package_id" validate:"required"`
	NetworkId    string `example:"{{NetworkUUID}}" json:"network_id" validate:"required"`
	SimType      string `example:"test" json:"sim_type" validate:"required"`
}

type SetActivePackageForSimReq struct {
	SimId     string `example:"{{SimUUID}}" path:"sim_id" validate:"required"`
	PackageId string `example:"{{PackageUUID}}" path:"package_id" validate:"required"`
}
type SimReq struct {
	SimId string `example:"{{SimUUID}}" form:"sim_id" json:"sim_id" path:"sim_id" binding:"required" validate:"required"`
}

type SimByNetworkReq struct {
	NetworkId string `example:"{{NetworkUUID}}" form:"network_id" json:"network_id" path:"network_id" binding:"required" validate:"required"`
}

type ActivateDeactivateSimReq struct {
	SimId  string `example:"{{SimUUID}}" path:"sim_id" validate:"required"`
	Status string `example:"active" json:"status" binding:"required" validate:"required"`
}

type GetSimsBySubReq struct {
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber_id" path:"subscriber_id" binding:"required" validate:"required"`
}
type AddPkgToSimReq struct {
	SimId     string                 `example:"{{SimUUID}}" json:"sim_id" validate:"required"`
	PackageId string                 `example:"{{PackageUUID}}" json:"package_id" validate:"required"`
	StartDate *timestamppb.Timestamp `example:"" json:"start_date" validate:"required"`
}

type RemovePkgFromSimReq struct {
	SimId     string `example:"{{SimUUID}}" form:"sim_id" json:"sim_id" path:"sim_id" binding:"required" validate:"required"`
	PackageId string `example:"{{PackageUUID}}" form:"package_id" json:"package_id" path:"package_id" binding:"required" validate:"required"`
}
