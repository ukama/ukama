package db

import (
	"time"

	"gorm.io/gorm"
)

type Rate struct {
	gorm.Model
	Country      string
	Network      string
	Vpmn         string
	Imsi         string
	Sms_mo       string
	Sms_mt       string
	Data         string
	X2g          string
	X3g          string
	X5g          string
	Lte          string
	Lte_m        string
	Apn          string
	Effective_at time.Time
	End_at       time.Time
	Sim_type     string
}
type RateModel struct {
	gorm.Model
	Id           int8
	Country      string
	Network      string
	Vpmn         string
	Imsi         string
	Sms_mo       string
	Sms_mt       string
	Data         string
	X2g          string
	X3g          string
	X5g          string
	Lte          string
	Lte_m        string
	Apn          string
	Created_at   time.Time
	Updated_at   time.Time
	Deleted_at   time.Time
	Effective_at time.Time
	End_at       time.Time
	Sim_type     string
}

func (r RateModel) ToObject() RateModel {
	rate := RateModel{
		Country:      r.Country,
		Network:      r.Network,
		Vpmn:         r.Vpmn,
		Imsi:         r.Imsi,
		Sms_mo:       r.Sms_mo,
		Sms_mt:       r.Sms_mt,
		Data:         r.Data,
		X2g:          r.X2g,
		X3g:          r.X3g,
		Lte:          r.Lte,
		Lte_m:        r.Lte_m,
		Apn:          r.Apn,
		Effective_at: r.Effective_at,
		End_at:       r.End_at,
		Sim_type:     r.Sim_type,
	}
	return rate
}

type RateList []*RateModel

func (r RateList) ToArray() []RateModel {
	var rates []RateModel
	for _, rate := range r {
		_rate := RateModel{
			Id:           rate.Id,
			X2g:          rate.X2g,
			X3g:          rate.X3g,
			Lte:          rate.Lte,
			Apn:          rate.Apn,
			Vpmn:         rate.Vpmn,
			Imsi:         rate.Imsi,
			Data:         rate.Data,
			Lte_m:        rate.Lte_m,
			Sms_mo:       rate.Sms_mo,
			Sms_mt:       rate.Sms_mt,
			End_at:       rate.End_at,
			Network:      rate.Network,
			Country:      rate.Country,
			Sim_type:     rate.Sim_type,
			Deleted_at:   rate.Deleted_at,
			Created_at:   rate.Created_at,
			Updated_at:   rate.Updated_at,
			Effective_at: rate.Effective_at,
		}
		rates = append(rates, _rate)
	}
	return rates
}
