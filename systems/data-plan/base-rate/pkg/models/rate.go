package models

type Rate struct {
	Id          int64  `json:"id" gorm:"primaryKey"`
	Country     string `json:"country"`
	Network     string `json:"network"`
	Vpmn        string `json:"vpmn"`
	Imsi        string `json:"imsi"`
	SmsMo       string `json:"smsMo"`
	SmsMt       string `json:"smsMt"`
	Data        string `json:"data"`
	X2g         string `json:"_2g"`
	X3g         string `json:"_3g"`
	X5g         string `json:"_5g"`
	Lte         string `json:"lte"`
	LteM        string `json:"lteM"`
	Apn         string `json:"apn"`
	CreatedAt   string `json:"createdAt"`
	EffectiveAt string `json:"effectiveAt"`
	EndAt       string `json:"endAt"`
	SimType     string `json:"simType"`
}
