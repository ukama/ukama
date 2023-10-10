package rest

type RestartNodeRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required" example:"{{NodeId}}"`
}

type RestartSiteRequest struct {
	SiteName  string `json:"site_name"  example:"site1" validate:"required"`
	NetworkId string `json:"network_id" example:"{{NetworkId}}" validate:"required"`
}

type ApplyConfigRequest struct {
	Commit string `json:"commit" path:"commit" example:"commit" validate:"required"`
}

type GetConfigVersionRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}
