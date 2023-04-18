package rest

type AddPackageRequest struct {
	Name        string  `example:"Monthly-Data" json:"name" validation:"required"`
	From        string  `example:"2023-04-01T00:00:00Z" json:"from" validation:"required"`
	To          string  `example:"2023-05-01T00:00:00Z" json:"to" validation:"required"`
	OrgId       string  `example:"{{OrgUUID}}" json:"org_id" validation:"required"`
	OwnerId     string  `example:"{{OwnerUUID}}" json:"owner_id" validation:"required"`
	SimType     string  `example:"test" json:"sim_type" validation:"required"`
	SmsVolume   int64   `example:"0" json:"sms_volume" validation:"required"`
	DataVolume  int64   `example:"1024" json:"data_volume" validation:"required"`
	DataUnit    string  `example:"MegaBytes" json:"data_unit" validation:"required"`
	VoiceUnit   string  `example:"seconds" json:"voice_unit" validation:"required"`
	Type        string  `example:"postpaid" json:"type" validation:"required"`
	Flatrate    bool    `example:"false" json:"flat_rate" default:"false"`
	Amount      float64 `example:"0" json:"amount" default:"0.00"`
	Markup      float64 `example:"0" json:"markup" default:"0.00"`
	Apn         string  `example:"ukama.tel" json:"apn" default:"ukama.tel"`
	Active      bool    `example:"true" json:"active" validation:"required"`
	VoiceVolume int64   `example:"0" json:"voice_volume" default:"0"`
	BaserateId  string  `example:"{{baserate}}" json:"baserate_id" validation:"required"`
}

type UpdatePackageRequest struct {
	Uuid   string `example:"{{PackageUUID}}" json:"uuid" path:"uuid" validation:"required"`
	Name   string `example:"Monthly-Data-Updated" json:"name" validation:"required"`
	Active bool   `example:"true" json:"active" validation:"required"`
}

type PackagesRequest struct {
	Uuid string `example:"{{PackageUUID}}" form:"uuid" json:"uuid" path:"uuid" binding:"required" validate:"required"`
}

type GetBaseRatesByCountryRequest struct {
	Country  string `path:"country" validate:"required"`
	Provider string `json:"network"`
	SimType  string `json:"sim_type" binding:"required" validate:"required"`
}

type GetBaseRatesForPeriodRequest struct {
	Country  string `path:"country" validate:"required"`
	Provider string `json:"network" binding:"required" validate:"required"`
	To       string `json:"to" binding:"required" validate:"required"`
	From     string `json:"from" binding:"required" validate:"required"`
	SimType  string `json:"sim_type" binding:"required" validate:"required"`
}

type GetBaseRateRequest struct {
	RateId string `path:"base_rate" validate:"required"`
}

type GetPackageByOrgRequest struct {
	OrgId string `example:"{{OrgUUID}}" form:"org_id" json:"org_id" path:"org_id" binding:"required" validate:"required"`
}
type UploadBaseRatesRequest struct {
	FileURL     string `json:"file_url" binding:"required" validate:"required"`
	EffectiveAt string `json:"effective_at" binding:"required" validate:"required"`
	EndAt       string `json:"end_at" validate:"required"`
	SimType     string `json:"sim_type" binding:"required" validate:"required"`
}

type GetRateRequest struct {
	OwnerId     string `example:"{{UserUUID}}" path:"user_id" validate:"required"`
	Country     string `json:"country" binding:"required" validate:"required"`
	Provider    string `json:"provider" binding:"required" validate:"required"`
	To          uint64 `json:"to" binding:"required" validate:"required"`
	From        uint64 `json:"from" binding:"required" validate:"required"`
	SimType     string `json:"sim_type" binding:"required" validate:"required"`
	EffectiveAt string `json:"effective_at" binding:"required" validate:"required"`
}

type DeleteMarkupRequest struct {
	OwnerId string `example:"{{UserUUID}}" form:"user_id" json:"user_id" path:"user_id" validate:"required"`
}

type SetMarkupRequest struct {
	OwnerId string  `example:"{{UserUUID}}" form:"user_id" json:"user_id" path:"user_id"  validate:"required"`
	Markup  float64 `example:"10" json:"markup" path:"markup" validate:"required"`
}

type GetMarkupRequest struct {
	OwnerId string `example:"{{UserUUID}}" form:"user_id" json:"user_id" path:"user_id" validate:"required"`
}

type GetMarkupHistoryRequest struct {
	OwnerId string `example:"{{UserUUID}}" form:"user_id" json:"user_id" path:"user_id" validate:"required"`
}

type SetDefaultMarkupRequest struct {
	Markup float64 `example:"10" json:"markup" path:"markup" validate:"required"`
}

type GetDefaultMarkupRequest struct {
}

type GetDefaultMarkupHistoryRequest struct {
}
