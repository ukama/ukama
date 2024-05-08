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

type PostCDRReq struct {
	Session       uint64 `json:"session" validate:"required"`
	Imsi          string `json:"imsi" path:"imsi" validate:"required"`
	NodeId        string `json:"node_id" validate:"required"`
	Policy        string `json:"policy" validate:"required"`
	ApnName       string `json:"apn_name" validate:"required"`
	Ip            string `json:"ip" validate:"required"`
	StartTime     uint64 `json:"start_time" validate:"required"`
	EndTime       uint64 `json:"end_time"`
	LastUpdatedAt uint64 `json:"last_updated_at" validate:"required"`
	TxBytes       uint64 `json:"tx_bytes"`
	RxBytes       uint64 `json:"rx_bytes"`
	TotalBytes    uint64 `json:"total_bytes"`
}

type GetCDRReq struct {
	Imsi      string `path:"imsi" validate:"required"`
	StartTime uint64 `query:"start_time" validate:"required"`
	EndTime   uint64 `query:"end_time" validate:"required"`
	Policy    string `query:"policy"`
	SessionId uint64 `query:"session_id"`
}

type GetUsageReq struct {
	Imsi      string `json:"imsi" path:"imsi" validate:"required"`
	StartTime uint64 `json:"start_time"`
	EndTime   uint64 `json:"end_time"`
	Policy    string `json:"policy"`
	SessionId uint64 `json:"session_id"`
}
