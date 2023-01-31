package rest

type ActivateReq struct {
	Iccid     string `path:"iccid" validate:"required"`
	Network   string `json:"network" validate:"required"`
	PackageId string `json:"packageId" validate:"required"`
}

type InactivateReq struct {
	Iccid string `path:"iccid, omitempty"`
}

type UpdatePackageReq struct {
	Iccid     string `path:"iccid" validate:"required"`
	PackageId string `json:"packageId" validate:"required"`
}

type ReadSubscriberReq struct {
	Iccid string `path:"iccid, omitempty"`
}


type ReadSubscriberResp struct {
	Iccid string `path:"iccid, omitempty"`
}