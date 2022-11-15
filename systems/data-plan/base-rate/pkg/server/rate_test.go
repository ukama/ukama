package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/data-plan/base-rate/mocks"
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"

	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
)

var mockCountry = "The lunar maria"

func TestRateService_Get(t *testing.T) {
	baseRateRepo := &mocks.BaseRateRepo{}
	baseRateRepo.On("GetBaseRate", uint64(1)).Return(&db.Rate{
		Country: mockCountry,
	}, nil)
	s := NewBaseRateServer(baseRateRepo)
	rate, err := s.GetBaseRate(context.TODO(), &pb.GetBaseRateRequest{RateId: uint64(1)})

	assert.NoError(t, err)
	assert.Equal(t, mockCountry, rate.Rate.Country)
	baseRateRepo.AssertExpectations(t)
}
func TestRateService_GetRates(t *testing.T) {
	var mockFilters = pb.GetBaseRatesRequest{
		Country:     "The lunar maria",
		Provider:    "ABC Tel",
		EffectiveAt: "2022-12-01T00:00:00Z",
		SimType:     validations.ReqStrTopb("inter_mno_data"),
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	baseRateRepo.On("GetBaseRates").Return(&db.Rate{
		Country:      "The lunar maria",
		Network:      "ABC Tel",
		Effective_at: "2022-12-01T00:00:00Z",
		Sim_type:     "inter_mno_data",
	}, nil)
	s := NewBaseRateServer(baseRateRepo)
	rate, err := s.GetBaseRates(context.TODO(), &pb.GetBaseRatesRequest{Country: mockFilters.Country, Provider: mockFilters.Provider, EffectiveAt: mockFilters.EffectiveAt, SimType: mockFilters.SimType})
	fmt.Println(rate)
	assert.NoError(t, err)
	// assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
	// baseRateRepo.AssertExpectations(t)
}
