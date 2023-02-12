package integration

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUploadBaseRates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set up connection to server
	conn, err := grpc.DialContext(ctx, "localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBaseRatesServiceClient(conn)

	// test data
	fileURL := "https://raw.githubusercontent.com/ukama/ukama/upload-rates/systems/data-plan/base-rate/template/template.csv"
	effectiveAt := "2023-03-01T00:00:00Z"
	simType := pb.SimType_INTER_MNO_DATA

	// create request
	req := &pb.UploadBaseRatesRequest{
		FileURL:    fileURL,
		EffectiveAt: effectiveAt,
		SimType:    simType,
	}

	// call method
	res, err := client.UploadBaseRates(ctx, req)
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
