package rest

type AddPackageRequest struct {
	Name        string `json:"name" validate:"required"`
	Duration    uint64 `json:"duration" validate:"required"`
	OrgId       uint64 `json:"org_id" validate:"required"`
	SimType     string `json:"sim_type" validate:"required"`
	SmsVolume   int64  `json:"sms_volume"`
	DataVolume  int64  `json:"data_volume" `
	Active      bool   `json:"active"`
	VoiceVolume int64  `json:"voice_volume" `
	OrgRatesId  uint64 `json:"org_rates_id" validate:"required"`
}

type UpdatePackageRequest struct {
	Id          uint64 `json:"id" validate:"required"`
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
	Id    uint64 `path:"package" validate:"required"`
	OrgId uint64 `json:"org_id"`
}

type GetPackagesRequest struct {
	Id    uint64 `path:"package" validate:"required"`
	OrgId uint64 `json:"org_id"`
}
