package rest

type GetNodeRequest struct {
	NodeId string `path:"node" validate:"required"`
}

type GetNodeResponse struct {
	NodeId      string `path:"node" validate:"required"`
	OrgName     string `path:"org" validate:"required"`
	Certificate string `json:"certificate"`
	Ip          string `json:"ip"`
}
