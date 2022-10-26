package models

import "database/sql/driver"

type SimType uint8

const (
	INTER_NONE      = 0
	INTER_MNO_DATA  = 1
	INTER_MNO_ALL   = 2
	INTER_UKAMA_ALL = 3
)

func (e *SimType) Scan(value interface{}) error {
	*e = SimType(uint8(value.(int64)))
	return nil
}

func (e SimType) Value() (driver.Value, error) {
	return uint8(e), nil
}

type Rate struct {
	Id           int64   `json:"id" gorm:"primaryKey"`
	Country      string  `json:"country"`
	Network      string  `json:"network"`
	Vpmn         string  `json:"vpmn"`
	Imsi         string  `json:"imsi"`
	Sms_mo       string  `json:"sms_mo"`
	Sms_mt       string  `json:"sms_mt"`
	Data         string  `json:"data"`
	X2g          string  `json:"_2g"`
	X3g          string  `json:"_3g"`
	X5g          string  `json:"_5g"`
	Lte          string  `json:"lte"`
	Lte_m        string  `json:"lte_m"`
	Apn          string  `json:"apn"`
	SimType      SimType `gorm:"type:uint;not null"`
	Created_at   string  `json:"created_at"`
	Effective_at string  `json:"effective_at"`
	End_at       string  `json:"end_at"`
}
