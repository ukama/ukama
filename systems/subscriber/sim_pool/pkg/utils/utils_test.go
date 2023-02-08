package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/subscriber/sim_pool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim_pool/pkg/db"
)

func TestSimPoolStats_Success(t *testing.T) {
	sims := PoolStats([]db.Sim{
		{
			Msisdn:         "2345678901",
			ActivationCode: "123456",
			IsAllocated:    false,
			SimType:        "inter_ukama_all",
			Iccid:          "1234567890123456789",
			SmDpAddress:    "http://localhost:8080",
		},
		{
			Msisdn:         "2345678901",
			ActivationCode: "123456",
			IsAllocated:    true,
			SimType:        "inter_ukama_all",
			Iccid:          "1234567890123456789",
			SmDpAddress:    "http://localhost:8080",
		},
	})
	assert.NotNil(t, sims)
	assert.Equal(t, sims.Available, uint64(1))
	assert.Equal(t, sims.Consumed, uint64(1))
	assert.Equal(t, sims.Total, uint64(2))
}

func TestPbParseToModel_Success(t *testing.T) {
	sims := PbParseToModel([]*pb.AddSim{{
		Iccid:          "1234567890123456789",
		Msisdn:         "2345678901",
		SimType:        pb.SimType(pb.SimType_value["inter_ukama_all"]),
		SmDpAddress:    "http://localhost:8080",
		ActivationCode: "123456",
	}})
	assert.NotNil(t, sims)
	assert.Equal(t, sims[0].Iccid, "1234567890123456789")
}

func TestBytesToRawSim_Success(t *testing.T) {
	rsims := RawSimToPb([]RawSim{{
		Iccid:          "1234567890123456789",
		Msisdn:         "2345678901",
		SmDpAddress:    "http://localhost:8080",
		ActivationCode: "123456",
		IsPhysical:     "true",
	}}, "inter_ukama_all")
	assert.NotNil(t, rsims)
	assert.Equal(t, rsims[0].Iccid, "1234567890123456789")
}
