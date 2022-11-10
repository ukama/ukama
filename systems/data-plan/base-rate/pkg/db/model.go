package db

import (
	"time"

	"github.com/ukama/ukama/systems/data-plan/base-rate/pb"
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

type RateList []*Rate

func (r Rate) ToPbRate() *pb.Rate {

	rate := &pb.Rate{
		Id:          int64(r.ID),
		X2G:         r.X2g,
		X3G:         r.X3g,
		X5G:         r.X5g,
		Lte:         r.Lte,
		Apn:         r.Apn,
		Vpmn:        r.Vpmn,
		Imsi:        r.Imsi,
		Data:        r.Data,
		LteM:        r.Lte_m,
		SmsMo:       r.Sms_mo,
		SmsMt:       r.Sms_mt,
		Network:     r.Network,
		Country:     r.Country,
		SimType:     r.Sim_type,
		EndAt:       r.End_at.String(),
		CreatedAt:   r.CreatedAt.String(),
		UpdatedAt:   r.UpdatedAt.String(),
		EffectiveAt: r.Effective_at.String(),
		DeletedAt:   r.DeletedAt.Time.String(),
	}

	return rate
}

func (r RateList) ToPbRates() []*pb.Rate {
	var rates []*pb.Rate
	for _, rate := range r {
		_rate := &pb.Rate{
			Id:          int64(rate.ID),
			X2G:         rate.X2g,
			X3G:         rate.X3g,
			X5G:         rate.X5g,
			Lte:         rate.Lte,
			Apn:         rate.Apn,
			Vpmn:        rate.Vpmn,
			Imsi:        rate.Imsi,
			Data:        rate.Data,
			LteM:        rate.Lte_m,
			SmsMo:       rate.Sms_mo,
			SmsMt:       rate.Sms_mt,
			Network:     rate.Network,
			Country:     rate.Country,
			SimType:     rate.Sim_type,
			EndAt:       rate.End_at.String(),
			CreatedAt:   rate.CreatedAt.String(),
			UpdatedAt:   rate.UpdatedAt.String(),
			EffectiveAt: rate.Effective_at.String(),
			DeletedAt:   rate.DeletedAt.Time.String(),
		}
		rates = append(rates, _rate)
	}
	return rates
}
