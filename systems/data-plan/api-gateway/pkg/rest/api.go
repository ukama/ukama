package rest

type AddPackageRequest struct {
	Name        string `json:"name" validate:"required"`
	Duration    uint64   `json:"duration" validate:"required"`
	OrgId       uint64 `json:"org_id" validate:"required"`
	SimType     string `json:"sim_type" validate:"required"`
	SmsVolume   int64 `json:"sms_volume" validate:"required"`
	DataVolume  int64 `json:"data_volume" validate:"required"`
	Active bool `json:"active" validate:"required"`
	VoiceVolume int64 `json:"voice_volume" validate:"required"`
	OrgRatesId  uint64 `json:"org_rates_id" validate:"required"`
}

type UpdatePackageRequest struct {
	Id          uint64 `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Duration    bool   `json:"duration" validate:"required"`
	OrgId       uint64 `json:"org_id" validate:"required"`
	SimType     string `json:"sim_type" validate:"required"`
	SmsVolume   uint64 `json:"sms_volume" validate:"required"`
	DataVolume  uint64 `json:"data_volume" validate:"required"`
	VoiceVolume uint64 `json:"voice_volume" validate:"required"`
	OrgRatesId  string `json:"org_rates_id" validate:"required"`
}
type DeletePackageRequest struct {
	Id    uint64 `json:"id" validate:"required"`
	OrgId uint64 `json:"org_id" validate:"required"`
}

type GetPackagesRequest struct {
	Id    uint64 `json:"id" validate:"required"`
	OrgId uint64 `json:"org_id"`
}
