package rest

type AddPackageRequest struct {
	Name        string `json:"name" validation:"required"`
	Duration    uint64 `json:"duration" validation:"required"`
	OrgId       string `json:"org_id" validation:"required"`
	SimType     string `json:"sim_type" validation:"required"`
	SmsVolume   int64  `json:"sms_volume" validation:"required"`
	DataVolume  int64  `json:"data_volume" validation:"required"`
	Active      bool   `json:"active" validation:"required"`
	VoiceVolume int64  `json:"voice_volume" validation:"required"`
	OrgRatesId  uint64 `json:"org_rates_id" validation:"required"`
}

type UpdatePackageRequest struct {
	Uuid        string `json:"uuid" validation:"required"`
	Name        string `json:"name" validation:"required"`
	Duration    uint64 `json:"duration" validation:"required"`
	Active      bool   `json:"active" validation:"required"`
	SimType     string `json:"sim_type" validation:"required"`
	SmsVolume   int64  `json:"sms_volume" validation:"required"`
	DataVolume  int64  `json:"data_volume" validation:"required"`
	VoiceVolume int64  `json:"voice_volume" validation:"required"`
	OrgRatesId  uint64 `json:"org_rates_id" validation:"required"`
}

type PackagesRequest struct {
	Id string `form:"package" json:"package" path:"package" binding:"required" validate:"required"`
}

type GetBaseRatesRequest struct {
	Country     string `json:"country" binding:"required" validate:"required"`
	Provider    string `json:"provider" binding:"required" validate:"required"`
	To          uint64 `json:"to" binding:"required" validate:"required"`
	From        uint64 `json:"from" binding:"required" validate:"required"`
	SimType     string `json:"sim_type" binding:"required" validate:"required"`
	EffectiveAt string `json:"effective_at" binding:"required" validate:"required"`
}
type GetBaseRateRequest struct {
	RateId string `form:"base_rate" json:"base_rate" path:"base_rate" binding:"required" validate:"required"`
}
type GetPackageByOrgRequest struct {
	OrgId string ` form:"org_id" json:"org_id" path:"org_id" binding:"required" validate:"required"`
}
type UploadBaseRatesRequest struct {
	FileURL     string `json:"file_url" binding:"required" validate:"required"`
	EffectiveAt string `json:"effective_at" binding:"required" validate:"required"`
	SimType     string `json:"sim_type" binding:"required" validate:"required"`
}
