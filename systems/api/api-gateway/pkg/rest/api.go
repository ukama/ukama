package rest

type GetNetworkReq struct {
	NetworkId string `json:"network_id" path:"network_id" validate:"required"`
}

type AddNetworkReq struct {
	OrgName          string   `example:"milky-way"  json:"org" validate:"required"`
	NetName          string   `example:"mesh-network" json:"network_name" validate:"required"`
	AllowedCountries []string `json:"allowed_countries"`
	AllowedNetworks  []string `json:"allowed_networks"`
	PaymentLinks     bool     `example:"true" json:"payment_links"`
}

type GetPackageReq struct {
	PackageId string `json:"package_id" path:"package_id" validate:"required"`
}

type AddPackageReq struct {
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
