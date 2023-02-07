package rest

type Guti struct {
	PlmnId string `json:"plmn_id" validate:"required"`
	Mmegi  uint32 `json:"mmegi" validate:"required"`
	Mmec   uint32 `json:"mmec" validate:"required"`
	Mtmsi  uint32 `json:"mtmsi" validate:"required"`
}

type UpdateGutiReq struct {
	Imsi      string `path:"imsi" validate:"required"`
	UpdatedAt uint32 `json:"updated_at" validate:"required"`
	Guti      Guti   `json:"guti" validate:"required"`
}

type GetSubscriberReq struct {
	Imsi string `path:"imsi" validate:"required"`
}

type UpdateTaiReq struct {
	Imsi      string `path:"imsi" validate:"required"`
	UpdatedAt uint32 `json:"updated_at" validate:"required"`
	PlmnId    string `json:"plmn_id" validate:"required"`
	Tac       uint32 `json:"tac" validate:"required"`
}
