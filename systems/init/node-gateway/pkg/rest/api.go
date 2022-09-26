package rest

/* For Node bootstarping */
// rpc GetNode(GetNodeRequest) returns (GetNodeResponse);

// message GetNodeResponse {
// string nodeId = 1;
// string orgName = 2;
// string certificate = 3;
// string ip = 4;
// }

// message GetNodeRequest{
// string nodeId = 1 [(validator.field) = {string_not_empty: true}];
// }

type GetNodeRequest struct {
	NodeId string `path:"node" validate:"required"`
}

type GetNodeResponse struct {
	NodeId      string `path:"node" validate:"required"`
	OrgName     string `path:"org" validate:"required"`
	Certificate string `json:"certificate"`
	Ip          string `json:"ip"`
}
