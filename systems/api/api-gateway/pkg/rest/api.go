package rest

type GetNetworkReq struct {
	NetworkId string `json:"network_id" path:"network_id" validate:"required"`
}

type AddNetworkReq struct {
	OrgName          string   `example:"milky-way"  json:"org" validate:"required"`
	NetName          string   `example:"mesh-network" json:"network_name" validate:"required"`
	AllowedCountries []string `json:"allowed_countries"`
	AllowedNetworks  []string `json:"allowed_networks"`
	Budget           float64  `json:"budget"`
	Overdraft        float64  `json:"overdraft"`
	TrafficPolicy    uint     `json:"traffic_policy"`
	PaymentLinks     bool     `example:"true" json:"payment_links"`
}

type GetPackageReq struct {
	PackageId string `json:"package_id" path:"package_id" validate:"required"`
}

type AddPackageReq struct {
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

type GetSimReq struct {
	Iccid string `json:"iccid" path:"iccid" validate:"required"`
}

type AddSimReq struct {
	SubscriberId string `json:"subscriber_id" validate:"required"`
	NetworkId    string `json:"network_id" validate:"required"`
	PackageId    string `json:"package_id" validate:"required"`
	SimType      string `json:"sim_type"`
	SimToken     string `json:"sim_token"`
}
