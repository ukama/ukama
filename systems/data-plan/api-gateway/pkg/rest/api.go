package rest

import (
	"github.com/ukama/ukama/systems/common/rest"
)

type AddPackageRequest struct {
	rest.BaseRequest
	Name          string   `example:"Monthly-Data" json:"name" validation:"required"`
	From          string   `example:"2023-04-01T00:00:00Z" json:"from" validation:"required"`
	To            string   `example:"2023-05-01T00:00:00Z" json:"to" validation:"required"`
	OrgId         string   `example:"{{OrgUUID}}" json:"org_id" validation:"required"`
	OwnerId       string   `example:"{{OwnerUUID}}" json:"owner_id" validation:"required"`
	SimType       string   `example:"test" json:"sim_type" validation:"required"`
	SmsVolume     int64    `example:"0" json:"sms_volume" validation:"required"`
	DataVolume    int64    `example:"1024" json:"data_volume" validation:"required"`
	DataUnit      string   `example:"MegaBytes" json:"data_unit" validation:"required"`
	VoiceUnit     string   `example:"seconds" json:"voice_unit" validation:"required"`
	Duration      uint64   `example:"1" json:"duration" validation:"required"`
	Type          string   `example:"postpaid" json:"type" validation:"required"`
	Flatrate      bool     `example:"false" json:"flat_rate" default:"false"`
	Amount        float64  `example:"0" json:"amount" default:"0.00"`
	Markup        float64  `example:"0" json:"markup" default:"0.00"`
	Apn           string   `example:"ukama.tel" json:"apn" default:"ukama.tel"`
	Active        bool     `example:"true" json:"active" validation:"required"`
	VoiceVolume   int64    `example:"0" json:"voice_volume" default:"0"`
	BaserateId    string   `example:"{{baserate}}" json:"baserate_id" validation:"required"`
	Overdraft     float64  `json:"overdraft"`
	TrafficPolicy uint32   `json:"traffic_policy"`
	Networks      []string `json:"networks"`
}

type UpdatePackageRequest struct {
	rest.BaseRequest
	Uuid   string `example:"{{PackageUUID}}" json:"uuid" path:"uuid" validation:"required"`
	Name   string `example:"Monthly-Data-Updated" json:"name" validation:"required"`
	Active bool   `example:"true" json:"active" validation:"required"`
}

type PackagesRequest struct {
	rest.BaseRequest
	Uuid string `example:"{{PackageUUID}}" form:"uuid" json:"uuid" path:"uuid" binding:"required" validate:"required"`
}

type GetBaseRatesByCountryRequest struct {
	rest.BaseRequest
	Country     string `json:"country" query:"country" binding:"required" validation:"required"`
	Provider    string `json:"provider" query:"provider" binding:"required"`
	SimType     string `json:"sim_type" query:"sim_type" binding:"required" validation:"required"`
	EffectiveAt string `json:"effective_at" query:"effective_at" binding:"required"`
}

type GetBaseRatesForPeriodRequest struct {
	rest.BaseRequest
	Country  string `query:"country" validate:"required"`
	Provider string `query:"provider" binding:"required" validate:"required"`
	To       string `query:"to" binding:"required" validate:"required"`
	From     string `query:"from" binding:"required" validate:"required"`
	SimType  string `query:"sim_type" binding:"required" validate:"required"`
}

type GetBaseRateRequest struct {
	// rest.BaseRequest
	RateId string `path:"base_rate" validate:"required"`
}

type GetPackageByOrgRequest struct {
	rest.BaseRequest
	OrgId string `example:"{{OrgUUID}}" form:"org_id" json:"org_id" path:"org_id" binding:"required" validate:"required"`
}
type UploadBaseRatesRequest struct {
	rest.BaseRequest
	FileURL     string `json:"file_url" binding:"required" validate:"required"`
	EffectiveAt string `json:"effective_at" binding:"required" validate:"required"`
	EndAt       string `json:"end_at" validate:"required"`
	SimType     string `json:"sim_type" binding:"required" validate:"required"`
}

type GetRateRequest struct {
	rest.BaseRequest
	UserId   string `json:"user_id" path:"user_id" binding:"required"`
	Country  string `json:"country" query:"country" binding:"required"`
	Provider string `json:"provider" query:"provider" binding:"required"`
	To       string `json:"to" query:"to" binding:"required" `
	From     string `json:"from" query:"from" binding:"required"`
	SimType  string `json:"sim_type" query:"sim_type" binding:"required"`
}

type DeleteMarkupRequest struct {
	rest.BaseRequest
	OwnerId string `example:"{{UserUUID}}" form:"user_id" json:"user_id" path:"user_id" validate:"required"`
}

type SetMarkupRequest struct {
	rest.BaseRequest
	OwnerId string  `example:"{{UserUUID}}" form:"user_id" json:"user_id" path:"user_id"  validate:"required"`
	Markup  float64 `example:"10" json:"markup" path:"markup" validate:"required"`
}

type GetMarkupRequest struct {
	rest.BaseRequest
	OwnerId string `example:"{{UserUUID}}" form:"user_id" json:"user_id" path:"user_id" validate:"required"`
}

type GetMarkupHistoryRequest struct {
	rest.BaseRequest
	OwnerId string `example:"{{UserUUID}}" form:"user_id" json:"user_id" path:"user_id" validate:"required"`
}

type SetDefaultMarkupRequest struct {
	rest.BaseRequest
	Markup float64 `example:"10" json:"markup" path:"markup" validate:"required"`
}

type GetDefaultMarkupRequest struct {
	rest.BaseRequest
}

type GetDefaultMarkupHistoryRequest struct {
	rest.BaseRequest
}
