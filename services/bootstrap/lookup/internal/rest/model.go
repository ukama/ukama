package rest

type GetDeviceResponse struct {
	NodeId      string `json:"nodeId"`
	OrgName     string `json:"orgName"`
	Certificate string `json:"certificate" binding:"base64"`
	Ip          string `json:"ip" validate:"ip"`
}

type AddOrgRequest struct {
	Certificate string `json:"certificate" binding:"required,base64"`
	Ip          string `json:"ip" binding:"required,ip"`
}
