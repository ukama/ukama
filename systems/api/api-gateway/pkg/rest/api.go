package rest

type GetNetworkReq struct {
	NetworkId string `json:"network_id" path:"network_id" validate:"required"`
}

type AddNetworkReq struct {
	OrgName string `example:"milky-way"  json:"org" validate:"required"`
	NetName string `example:"mesh-network" json:"network_name" validate:"required"`
}
