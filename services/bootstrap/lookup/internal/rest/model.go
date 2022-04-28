package rest

type ReqAddNode struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"eq=add_node,required"`
	OrgName   string `query:"org" validate:"required"`
}

type ReqGetNode struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=org_credentials,required"`
	OrgName    string `query:"org"`
}

type RespGetNode struct {
	NodeId      string `json:"node"`
	OrgName     string `json:"org"`
	Certificate string `json:"certificate" binding:"base64"`
	Ip          string `json:"ip" validate:"ip"`
}

type ReqAddOrg struct {
	OrgName     string `query:"org" validate:"required"`
	LookingTo   string `query:"looking_to" validate:"eq=add_org,required"`
	Certificate string `json:"certificate" binding:"required,base64"`
	Ip          string `json:"ip" binding:"required,ip"`
}
