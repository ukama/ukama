package models

import "database/sql/driver"

type SIM_TYPE uint8

const (
	INTER_NONE      = "inter_none"
	INTER_MNO_DATA  = "inter_mno_data"
	INTER_MNO_ALL   = "inter_mno_all"
	INTER_UKAMA_ALL = "inter_ukama_all"
)

func (e *SIM_TYPE) Scan(value interface{}) error {
	*e = SIM_TYPE(uint8(value.(int64)))
	return nil
}

func (e SIM_TYPE) Value() (driver.Value, error) {
	return uint8(e), nil
}

type Rate struct {
	Id           int64  `json:"id" gorm:"primaryKey"`
	Country      string `json:"country"`
	Network      string `json:"network"`
	Vpmn         string `json:"vpmn"`
	Imsi         string `json:"imsi"`
	Sms_mo       string `json:"sms_mo"`
	Sms_mt       string `json:"sms_mt"`
	Data         string `json:"data"`
	X2g          string `json:"x2g"`
	X3g          string `json:"x3g"`
	X5g          string `json:"x5g"`
	Lte          string `json:"lte"`
	Lte_m        string `json:"lte_m"`
	Apn          string `json:"apn"`
	Created_at   string `json:"created_at"`
	Effective_at string `json:"effective_at"`
	End_at       string `json:"end_at"`
	Sim_type     string `json:"sim_type"`
}
