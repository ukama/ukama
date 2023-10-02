//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

func CreateBaseRateClient() (*grpc.ClientConn, pb.BaseRatesServiceClient, error) {
	log.Infoln("Connecting to BaseRate service ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewBaseRatesServiceClient(conn)
	return conn, c, nil
}

func Test_UploadBaseRate(t *testing.T) {

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateBaseRateClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	// test data
	fileURL := "https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/base-rate/template/template.csv"
	effectiveAt := "2023-04-11T07:20:50.52Z"
	simType := "ukama_data"

	// create request
	req := &pb.UploadBaseRatesRequest{
		FileURL:     fileURL,
		EffectiveAt: effectiveAt,
		SimType:     simType,
	}

	res, err := c.UploadBaseRates(ctx, req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok && s.Code() == codes.InvalidArgument {
			assert.Equal(t, "Please supply valid fileURL: \"https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/base-rate/template/template.csv\", effectiveAt: \"2023-03-01T00:00:00Z\" & simType: WORLD", s.Message())
			return
		}
	}

	// validate response
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
