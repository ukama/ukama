package rest

type AddPackageRequest struct {
	Name        string `example:"Monthly-Data" json:"name" validation:"required"`
	Duration    uint64 `example:"36000" json:"duration" validation:"required"`
	OrgId       string `example:"{{OrgUUID}}" json:"org_id" validation:"required"`
	SimType     string `example:"test" json:"sim_type" validation:"required"`
	SmsVolume   int64  `example:"0" json:"sms_volume" validation:"required"`
	DataVolume  int64  `example:"1024" json:"data_volume" validation:"required"`
	Active      bool   `example:"true" json:"active" validation:"required"`
	VoiceVolume int64  `example:"0" json:"voice_volume" validation:"required"`
	BaserateId  string `example:"{{baserate}}" json:"baserate_id" validation:"required"`
}

type UpdatePackageRequest struct {
	Uuid        string `example:"{{PackageUUID}}" json:"uuid" path:"uuid" binding:"required" validation:"required"`
	Name        string `example:"Monthly-Data-Updated" json:"name" validation:"required"`
	Duration    uint64 `example:"36000" json:"duration" validation:"required"`
	OrgId       string `example:"{{OrgUUID}}" json:"org_id" validation:"required"`
	SimType     string `example:"test" json:"sim_type" validation:"required"`
	SmsVolume   int64  `example:"0" json:"sms_volume" validation:"required"`
	DataVolume  int64  `example:"1024" json:"data_volume" validation:"required"`
	Active      bool   `example:"true" json:"active" validation:"required"`
	VoiceVolume int64  `example:"0" json:"voice_volume" validation:"required"`
	OrgRatesId  uint64 `example:"1" json:"org_rates_id" validation:"required"`
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
	EndAt       string `json:"end_at" binding:"required" validate:"required" default:"2099-12-31T00:00:00+00:00"`
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
