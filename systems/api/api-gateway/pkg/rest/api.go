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
	TrafficPolicy    uint32   `json:"traffic_policy"`
	PaymentLinks     bool     `example:"true" json:"payment_links"`
}

type GetPackageReq struct {
	PackageId string `json:"package_id" path:"package_id" validate:"required"`
}

type AddPackageReq struct {
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

type GetSimReq struct {
	Id string `json:"id" path:"id" validate:"required"`
}

type AddSimReq struct {
	SubscriberId          string `json:"subscriber_id"`
	OrgId                 string `json:"org_id"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Email                 string `json:"email"`
	PhoneNumber           string `json:"phone_number"`
	Address               string `json:"address"`
	Dob                   string `json:"date_of_birth"`
	ProofOfIdentification string `json:"proof_of_identification"`
	IdSerial              string `json:"id_serial"`
	NetworkId             string `json:"network_id" validate:"required"`
	PackageId             string `json:"package_id" validate:"required"`
	SimType               string `json:"sim_type" validate:"required"`
	SimToken              string `json:"sim_token"`
	TrafficPolicy         uint32 `json:"traffic_policy"`
}
