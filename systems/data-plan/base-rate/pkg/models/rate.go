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
	Id          int64   `json:"id" gorm:"primaryKey"`
	Country     string  `json:"country"`
	Network     string  `json:"network"`
	Vpmn        string  `json:"vpmn"`
	Imsi        string  `json:"imsi"`
	SmsMo       string  `json:"smsMo"`
	SmsMt       string  `json:"smsMt"`
	Data        string  `json:"data"`
	X2g         string  `json:"_2g"`
	X3g         string  `json:"_3g"`
	X5g         string  `json:"_5g"`
	Lte         string  `json:"lte"`
	LteM        string  `json:"lteM"`
	Apn         string  `json:"apn"`
	CreatedAt   string  `json:"createdAt"`
	EffectiveAt string  `json:"effectiveAt"`
	EndAt       string  `json:"endAt"`
	SimType     SimType `gorm:"type:uint;not null"`
}
