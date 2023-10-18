package rest

type RestartNodeRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
}

type RestartSiteRequest struct {
	SiteName  string `json:"site_name"  example:"site1" validate:"required" path:"site_name"`
	NetworkId string `json:"network_id" example:"{{NetworkId}}" validate:"required" path:"network_id"`
}

type RestartNodesRequest struct {
	NetworkId string   `json:"network_id" example:"{{NetworkId}}" validate:"required" path:"network_id"`
	NodeIds   []string `json:"node_ids" example:"{{NodeIds}}" validate:"required"`
}

type ApplyConfigRequest struct {
	Commit string `json:"commit" path:"commit" example:"commit" validate:"required"`
}

type GetConfigVersionRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}
