//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

func CreateRateClient() (*grpc.ClientConn, pb.RateServiceClient, error) {
	log.Infoln("Connecting to Rate service ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewRateServiceClient(conn)
	return conn, c, nil
}

func Test_FullFlow(t *testing.T) {

	var markupVal float64 = 10
	ownerId := uuid.NewV4()
	var userMarkup float64 = 5

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateRateClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	t.Run("UpdateDefaultMarkup", func(t *testing.T) {
		_, err := c.UpdateDefaultMarkup(ctx, &pb.UpdateDefaultMarkupRequest{
			Markup: markupVal,
		})
		assert.NoError(t, err)
	})

	t.Run("GetDefaultMarkup", func(t *testing.T) {
		resp, err := c.GetDefaultMarkup(ctx, &pb.GetDefaultMarkupRequest{})
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, markupVal, resp.Markup)
		}
	})

	t.Run("GetMarkup", func(t *testing.T) {
		resp, err := c.GetMarkup(ctx, &pb.GetMarkupRequest{
			OwnerId: ownerId.String(),
		})
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, markupVal, resp.Markup)
		}
	})

	t.Run("UpdateMarkup", func(t *testing.T) {
		_, err := c.UpdateMarkup(ctx, &pb.UpdateMarkupRequest{
			OwnerId: ownerId.String(),
			Markup:  userMarkup,
		})
		assert.NoError(t, err)
	})

	t.Run("GetUserMarkup", func(t *testing.T) {
		resp, err := c.GetMarkup(ctx, &pb.GetMarkupRequest{
			OwnerId: ownerId.String(),
		})
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, userMarkup, resp.Markup)
		}
	})

	t.Run("DeleteMarkup", func(t *testing.T) {
		_, err := c.DeleteMarkup(ctx, &pb.DeleteMarkupRequest{
			OwnerId: ownerId.String(),
		})
		assert.NoError(t, err)
	})

	t.Run("GetMarkupHistory", func(t *testing.T) {
		resp, err := c.GetMarkupHistory(ctx, &pb.GetMarkupHistoryRequest{
			OwnerId: ownerId.String(),
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetDefaultMarkupHistory", func(t *testing.T) {
		resp, err := c.GetDefaultMarkupHistory(ctx, &pb.GetDefaultMarkupHistoryRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

}
