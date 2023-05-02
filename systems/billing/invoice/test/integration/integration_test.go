//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

func init() {
	// load config
	tConfig = &TestConfig{}

	reader := rconf.NewConfReader("integration")

	err := reader.Read(tConfig)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("Config: %+v\n", tConfig)
}

func Test_FullFlow(t *testing.T) {
	// we need real subscriberId from subscriber-registry
	subscriberId := uuid.NewV4().String()

	var period = time.Now().UTC()
	var raw = "{}"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateInvoiceServiceClient()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}
	defer conn.Close()

	invResp := &pb.AddResponse{}

	t.Run("AddInvoice", func(t *testing.T) {
		var err error

		invResp, err = c.Add(ctx, &pb.AddRequest{
			SubscriberId: subscriberId,
			Period:       timestamppb.New(period),
			RawInvoice:   raw,
		})

		assert.NoError(t, err)
	})

	t.Run("GetInvoice", func(t *testing.T) {
		_, err := c.Get(ctx, &pb.GetRequest{
			InvoiceId: invResp.Invoice.Id,
		})

		assert.NoError(t, err)
	})

	t.Run("GetInvoiceBySubscriber", func(t *testing.T) {
		_, err := c.GetBySubscriber(ctx, &pb.GetBySubscriberRequest{
			SubscriberId: subscriberId,
		})

		assert.NoError(t, err)
	})

	t.Run("DeleteInvoice", func(t *testing.T) {
		_, err := c.Delete(ctx, &pb.DeleteRequest{
			InvoiceId: invResp.Invoice.Id,
		})

		assert.NoError(t, err)
	})

	t.Run("GetInvoice", func(t *testing.T) {
		_, err := c.Get(ctx, &pb.GetRequest{
			InvoiceId: invResp.Invoice.Id,
		})

		assert.Error(t, err)
	})
}

func CreateInvoiceServiceClient() (*grpc.ClientConn, pb.InvoiceServiceClient, error) {
	log.Infoln("Connecting to Invoice Server ", tConfig.ServiceHost)

	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewInvoiceServiceClient(conn)
	return conn, c, nil
}
