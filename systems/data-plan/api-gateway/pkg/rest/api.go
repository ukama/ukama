package rest

type AddPackageRequest struct {
	Name        string `json:"name" binding:"required"`
	Duration    uint64 `json:"duration" `
	OrgId       uint64 `json:"org_id" binding:"required"`
	SimType     string `json:"sim_type" `
	SmsVolume   int64  `json:"sms_volume" `
	DataVolume  int64  `json:"data_volume" `
	Active      bool   `json:"active" `
	VoiceVolume int64  `json:"voice_volume"`
	OrgRatesId  uint64 `json:"org_rates_id" binding:"required"`
}

type UpdatePackageRequest struct {
	Id          uint64 `json:"id" binding:"required"`
	Name        string `json:"name" `
	Duration    uint64 `json:"duration" `
	OrgId       uint64 `json:"org_id" `
	Active      bool   `json:"active"`
	SimType     string `json:"sim_type" `
	SmsVolume   int64  `json:"sms_volume" `
	DataVolume  int64  `json:"data_volume" `
	VoiceVolume int64  `json:"voice_volume" `
	OrgRatesId  uint64 `json:"org_rates_id"`
}
type DeletePackageRequest struct {
	Id uint64 `path:"package" uri:"package" validate:"required"`
}

type GetPackagesRequest struct {
	Id uint64 `path:"package" uri:"package" validate:"required"`
}

type GetBaseRatesRequest struct {
	Country     string `json:"country"`
	Provider    string `json:"provider"`
	To          uint64 `json:"to"`
	From        uint64 `json:"from"`
	SimType     string `json:"sim_type"`
	EffectiveAt string `json:"effective_at"`
}
type GetBaseRateRequest struct {
	RateId uint64 `path:"baseRate" uri:"baseRate" validate:"required"`
}
type GetPackageByOrgRequest struct {
	OrgId uint64 `json:"org_id"`
}
type UploadBaseRatesRequest struct {
	FileURL     string `json:"file_url" binding:"required,url"`
	EffectiveAt string `json:"effective_at" binding:"required"`
	SimType     string `json:"sim_type" binding:"required"`
}
