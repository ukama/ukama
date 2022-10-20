package services

import (
	"context"
	"net/http"

	"github.com/ukama/ukama/systems/data-plan/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/pkg/models"
)

type Server struct {
	H db.Handler
}



func (s *Server) GetRates(ctx context.Context, req *pb.RatesRequest) (*pb.RatesResponse, error) {
	var rate models.Rate

	if result := s.H.DB.First(&rate, req.GetCountry()); result.Error != nil {
		return &pb.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	data := &pb.FindOneData{
		Id:    rate.Id,
		Country:  rate.Country,
		Country_on_cronus: rate.Country_on_cronus,
		Network: rate.Network,
		Network_id_on_cronus: rate.Network_id_on_cronus,
   		Vpmn: rate.Vpmn,
   		Imsi: rate.Imsi,
 		Sms_mo: rate.Sms_mo,
 		Sms_mt: rate.Sms_mt,
   		Data: rate.Data,
   		X2g: rate.X2g,
        X3g: rate.X3g,
    	Lte: rate.Lte,
  		Lte_m: rate.Lte_m,
    	Apn: rate.Apn,
		Created_at :rate.Created_at,
		Effective_at :rate.Effective_at,
		End_at :rate.End_at,
	}

	return &pb.RatesResponse{
		Data:   data,
	}, nil
}


