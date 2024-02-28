package rest

type ReqData struct {
	Iccid     string `json:"iccid" path:"iccid" validate:"required"`
	Imsi      string `json:"imsi,omitempty"`
	SimId     string `json:"sim_id,omitempty"`
	PackageId string `json:"package_id,omitempty"`
	NetworkId string `json:"netwrok_id,omitempty"`
}

type ActivateReq ReqData

type DeactivateReq ReqData

type UpdatePackageReq ReqData

type ReadSubscriberReq struct {
	Iccid string `path:"iccid" validate:"required"`
}
