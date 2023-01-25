package integration

import (
	"context"
	"testing"
	"time"

	confr "github.com/num30/config"
	"github.com/ukama/ukama/systems/common/config"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{}
	r := confr.NewConfReader("integration")
	r.Read(tConfig)

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", tConfig)
}

func Test_FullFlow(t *testing.T) {
	const sysName = "sys"

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateSimClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}
	t.Run("Add", func(t *testing.T) {
		_, err := c.Add(ctx, &pb.AddRequest{
			Sim: []*pb.AddSim{
				{
					Iccid: "123456789", SimType: pb.SimType_INTER_MNO_DATA, Msisdn: "555-555-1234", SmDpAddress: "http://example.com", ActivationCode: "abc123", QrCode: "qr123", IsPhysical: true,
				},
				{
					Iccid: "12273", SimType: pb.SimType_INTER_MNO_DATA, Msisdn: "583-5343-0234", SmDpAddress: "http://example.com", ActivationCode: "abc123", QrCode: "qr123", IsPhysical: true,
				},
			},
		})
		assert.NoError(t, err)

	})


	t.Run("Get", func(t *testing.T) {
		r, err := c.Get(ctx, &pb.GetRequest{
			IsPhysicalSim:  true,
			SimType: pb.SimType_INTER_MNO_DATA,
		})

		if assert.NoError(t, err) {
            assert.Equal(t, true, r.Sim.IsPhysical)
		}
	})

	t.Run("GetByIccid", func(t *testing.T) {
		r, err := c.GetByIccid(ctx, &pb.GetByIccidRequest{
			Iccid: "123456789",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, "123456789", r.Sim.Iccid)
		}
	})
    t.Run("GetStats", func(t *testing.T) {
		r, err := c.GetStats(ctx, &pb.GetStatsRequest{
			SimType:pb.SimType_INTER_MNO_DATA,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, uint64(0x0), r.Total)
		}
	})

	t.Run("Delete", func(T *testing.T) {
		_, err := c.Delete(ctx, &pb.DeleteRequest{
			Id:[]uint64{
                1,
                2,
            },
			
		})
        
		assert.NoError(t, err)
	})

}

func CreateSimClient() (*grpc.ClientConn, pb.SimServiceClient, error) {
	log.Infoln("Connecting to SimPool ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewSimServiceClient(conn)
	return conn, c, nil
}
