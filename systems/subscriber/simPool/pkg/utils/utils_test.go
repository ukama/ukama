package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/subscriber/simPool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/simPool/pkg/db"
)

func TestSimPoolStats_Success(t *testing.T) {
	simPool := SimPoolStats([]db.SimPool{
		{
			Msisdn:         "2345678901",
			ActivationCode: "123456",
			Is_allocated:   false,
			Sim_type:       "inter_ukama_all",
			Iccid:          "1234567890123456789",
			SmDpAddress:    "http://localhost:8080",
			QrCode:         "http://localhost:8080",
		},
		{
			Msisdn:         "2345678901",
			ActivationCode: "123456",
			Is_allocated:   true,
			Sim_type:       "inter_ukama_all",
			Iccid:          "1234567890123456789",
			SmDpAddress:    "http://localhost:8080",
			QrCode:         "http://localhost:8080",
		},
	})
	assert.NotNil(t, simPool)
	assert.Equal(t, simPool.Available, uint64(1))
	assert.Equal(t, simPool.Consumed, uint64(1))
	assert.Equal(t, simPool.Total, uint64(2))
}

// Fetch data success case
func TestPbParseToModel_Success(t *testing.T) {
	simPool := PbParseToModel([]*pb.AddSimPool{{
		Iccid:          "1234567890123456789",
		Msisdn:         "2345678901",
		QrCode:         "http://localhost:8080",
		SimType:        pb.SimType(pb.SimType_value["inter_ukama_all"]),
		SmDpAddress:    "http://localhost:8080",
		ActivationCode: "123456",
	}})
	assert.NotNil(t, simPool)
	assert.Equal(t, simPool[0].Iccid, "1234567890123456789")
}
