package rest

// org group

type AddOrgRequest struct {
	OrgName string `path:"org" validate:"required"`
	Owner   string `path:"org" validate:"required"`
}

// Users group

type GetUserRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
}

type UserRequest struct {
	Org      string `path:"org" validate:"required"`
	SimToken string `json:"simToken"`
	Name     string `json:"name,omitempty" validate:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
}

type UpdateUserRequest struct {
	OrgName       string `path:"org" validate:"required"`
	UserId        string `path:"user" validate:"required"`
	Name          string `json:"name,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	IsDeactivated bool   `json:"isDeactivated,omitempty"`
}

type SetSimStatusRequest struct {
	OrgName string       `path:"org" validate:"required"`
	UserId  string       `path:"user" validate:"required"`
	Iccid   string       `path:"iccid" validate:"required"`
	Carrier *SimServices `json:"carrier,omitempty"`
	Ukama   *SimServices `json:"ukama,omitempty"`
}

type GetSimQrRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
	Iccid   string `path:"iccid" validate:"required"`
}

type SimServices struct {
	Voice *bool `json:"voice,omitempty"`
	Sms   *bool `json:"sms,omitempty"`
	Data  *bool `json:"data,omitempty"`
}

type DeleteUserRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
}

// Nodes group

type NodesList struct {
	OrgName string  `json:"orgName"`
	Nodes   []*Node `json:"nodes"`
}

type Node struct {
	NodeId string `json:"nodeId,omitempty"`
	State  string `json:"state,omitempty"`
	Type   string `json:"type,omitempty"`
	Name   string `json:"name"`
}

type NodeExtended struct {
	Node
	Attached []*Node `json:"attached,omitempty"`
}

type GetNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

// struct for creating or updating node
type AddUpdateNodeRequest struct {
	OrgName string     `path:"org" validate:"required"`
	NodeId  string     `path:"node" validate:"required"`
	Node    NodeModify `json:"node" validate:"required"`
}

type NodeModify struct {
	Name     string        `json:"name,omitempty"`
	Attached []*NodeAttach `json:"attached,omitempty"`
}

type NodeAttach struct {
	NodeId string `json:"nodeId,omitempty" validate:"required"`
}

type DeleteNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

type DetachNodeRequest struct {
	OrgName        string `path:"org" validate:"required"`
	NodeId         string `path:"node" validate:"required"`
	AttachedNodeId string `path:"attachedId" validate:"required"`
}
