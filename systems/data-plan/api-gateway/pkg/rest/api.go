package rest

import "github.com/ukama/ukama/systems/common/uuid"

type AddPackageRequest struct {
	Name        string `json:"name" validation:"required"`
	Duration    uint64 `json:"duration" `
	OrgID       uuid.UUID `json:"org_id" validation:"required"`
	SimType     string `json:"sim_type" `
	SmsVolume   int64  `json:"sms_volume" `
	DataVolume  int64  `json:"data_volume" `
	Active      bool   `json:"active" `
	VoiceVolume int64  `json:"voice_volume"`
	OrgRatesId  uint64 `json:"org_rates_id" validation:"required"`
}

type UpdatePackageRequest struct {
	packageID uuid.UUID `path:"package" validate:"required"`
	Name        string `json:"name" `
	Duration    uint64 `json:"duration" `
	Active      bool   `json:"active"`
	SimType     string `json:"sim_type" `
	SmsVolume   int64  `json:"sms_volume" `
	DataVolume  int64  `json:"data_volume" `
	VoiceVolume int64  `json:"voice_volume" `
	OrgRatesId  uint64 `json:"org_rates_id"`
}
type DeletePackageRequest struct {
	packageID  uuid.UUID `path:"package" validate:"required"`
}

type GetPackagesRequest struct {
	packageID  uuid.UUID `path:"package" validate:"required"`
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
	RateId uuid.UUID `path:"baseRate" validate:"required"`
}
type GetPackageByOrgRequest struct {
	OrgId uuid.UUID `path:"org" validate:"required"`
}
type UploadBaseRatesRequest struct {
	FileURL     string `json:"file_url" validate:"required,url"`
	EffectiveAt string `json:"effective_at" validate:"required"`
	SimType     string `json:"sim_type" validate:"required"`
}
