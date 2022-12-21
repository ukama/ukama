package utils

import (
	pb "github.com/ukama/ukama/systems/subscriber/simPool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/simPool/pkg/db"
)

func SimPoolStats(slice []db.SimPool) *pb.GetStatsResponse {
	total := len(slice)
	failed := 0
	available := 0
	consumed := 0
	for _, value := range slice {
		if value.Is_allocated {
			consumed = consumed + 1
		} else {
			available = available + 1
		}
	}
	return &pb.GetStatsResponse{
		Total:     uint64(total),
		Failed:    uint64(failed),
		Available: uint64(available),
		Consumed:  uint64(consumed),
	}
}

func PbParseToModel(slice []*pb.AddSimPool) []db.SimPool {
	var simPool []db.SimPool
	for _, value := range slice {
		simPool = append(simPool, db.SimPool{
			Iccid:          value.Iccid,
			Msisdn:         value.Msisdn,
			SmDpAddress:    value.SmDpAddress,
			ActivationCode: value.ActivationCode,
			QrCode:         value.QrCode,
		})
	}
	return simPool
}
