package utils

import (
	"github.com/jszwec/csvutil"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
)

func PoolStats(slice []db.Sim) *pb.GetStatsResponse {
	total := len(slice)
	failed := 0
	available := 0
	consumed := 0
	physical := 0
	esim := 0
	for _, value := range slice {
		if value.IsAllocated {
			consumed = consumed + 1
		} else if value.IsFailed {
			failed = failed + 1
		} else {
			available = available + 1
			if value.IsPhysical {
				physical = physical + 1
			} else {
				esim = esim + 1
			}
		}
	}
	return &pb.GetStatsResponse{
		Total:     uint64(total),
		Failed:    uint64(failed),
		Available: uint64(available),
		Consumed:  uint64(consumed),
		Physical:  uint64(physical),
		Esim:      uint64(esim),
	}
}

func PbParseToModel(slice []*pb.AddSim) []db.Sim {
	var sims []db.Sim
	for _, value := range slice {
		sims = append(sims, db.Sim{
			Iccid:          value.Iccid,
			Msisdn:         value.Msisdn,
			SmDpAddress:    value.SmDpAddress,
			ActivationCode: value.ActivationCode,
			QrCode:         value.QrCode,
			SimType:        db.ParseType(value.SimType),
			IsPhysical:     value.IsPhysical,
		})
	}
	return sims
}

type RawSim struct {
	Iccid          string `csv:"ICCID"`
	Msisdn         string `csv:"MSISDN"`
	SmDpAddress    string `csv:"SmDpAddress"`
	ActivationCode string `csv:"ActivationCode"`
	QrCode         string `csv:"QrCode"`
	IsPhysical     string `csv:"IsPhysical"`
}

func ParseBytesToRawSim(b []byte) ([]RawSim, error) {
	var r []RawSim
	err := csvutil.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func RawSimToPb(r []RawSim, simType string) []db.Sim {
	var s []db.Sim
	for _, value := range r {
		s = append(s, db.Sim{
			Iccid:          value.Iccid,
			Msisdn:         value.Msisdn,
			SmDpAddress:    value.SmDpAddress,
			ActivationCode: value.ActivationCode,
			IsPhysical:     value.IsPhysical == "TRUE",
			SimType:        db.ParseType(simType),
			QrCode:         value.QrCode,
		})
	}
	return s
}
