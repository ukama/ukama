//go:build integration
// +build integration

package integration

import (
	"context"

	confr "github.com/num30/config"
	"github.com/ukama/ukama/systems/common/config"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

var sim = client.SimCardInfo{
	Iccid:          "012345678901234567891",
	Imsi:           "012345678912345",
	Op:             []byte("0123456789012345"),
	Key:            []byte("0123456789012345"),
	Amf:            []byte("800"),
	AlgoType:       1,
	UeDlAmbrBps:    2000000,
	UeUlAmbrBps:    2000000,
	Sqn:            1,
	CsgIdPrsent:    false,
	CsgId:          0,
	DefaultApnName: "ukama",
}

var guti = db.Guti{
	Imsi:            "012345678912345",
	PlmnId:          "00101",
	Mmegi:           101,
	Mmec:            101,
	MTmsi:           101,
	DeviceUpdatedAt: time.Now(),
}

var tai = db.Tai{
	PlmnId:          "00101",
	Tac:             101,
	DeviceUpdatedAt: time.Now(),
}

func init() {
	tConfig = &TestConfig{}
	r := confr.NewConfReader("integration")
	r.Read(tConfig)

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", tConfig)
}

func CreateAsrClient() (*grpc.ClientConn, pb.AsrRecordServiceClient, error) {
	log.Infoln("Connecting to ASR ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewAsrRecordServiceClient(conn)
	return conn, c, nil
}
func Test_FullFlow(t *testing.T) {
	var Imsi string
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateAsrClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	t.Run("Activate", func(t *testing.T) {
		_, err := c.Activate(ctx, &pb.ActivateReq{
			Network:   "40987edb-ebb6-4f84-a27c-99db7c136127",
			Iccid:     sim.Iccid,
			PackageId: "40987edb-ebb6-4f84-a27c-99db7c136300",
		})
		assert.NoError(t, err)
	})

	t.Run("ReadByIccid", func(t *testing.T) {
		resp, err := c.Read(ctx, &pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: sim.Iccid,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		Imsi = resp.Record.Imsi
	})

	t.Run("ReadByImsi", func(t *testing.T) {
		resp, err := c.Read(ctx, &pb.ReadReq{
			Id: &pb.ReadReq_Imsi{
				Imsi: Imsi,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("UpdateGuti", func(t *testing.T) {
		_, err := c.UpdateGuti(ctx, &pb.UpdateGutiReq{
			Imsi: Imsi,
			Guti: &pb.Guti{
				PlmnId: guti.PlmnId,
				Mmegi:  guti.Mmegi,
				Mmec:   guti.Mmec,
				Mtmsi:  guti.MTmsi,
			},
			UpdatedAt: uint32(time.Now().Unix()),
		})
		assert.NoError(t, err)
	})

	t.Run("UpdateTai", func(t *testing.T) {
		_, err := c.UpdateTai(ctx, &pb.UpdateTaiReq{
			Imsi:      Imsi,
			PlmnId:    tai.PlmnId,
			Tac:       tai.Tac,
			UpdatedAt: uint32(time.Now().Unix()),
		})
		assert.NoError(t, err)
	})

	t.Run("UpdatePackage", func(t *testing.T) {
		_, err := c.UpdatePackage(ctx, &pb.UpdatePackageReq{
			Iccid:     sim.Iccid,
			PackageId: "40987edb-ebb6-4f84-a27c-99db7c136127",
		})
		assert.NoError(t, err)
	})

	t.Run("Inactivate", func(t *testing.T) {
		_, err := c.Inactivate(ctx, &pb.InactivateReq{
			Id: &pb.InactivateReq_Iccid{
				Iccid: sim.Iccid,
			},
		})
		assert.NoError(t, err)
	})

	t.Run("Activate", func(t *testing.T) {
		_, err := c.Activate(ctx, &pb.ActivateReq{
			Network:   "40987edb-ebb6-4f84-a27c-99db7c136127",
			Iccid:     sim.Iccid,
			PackageId: "40987edb-ebb6-4f84-a27c-99db7c136300",
		})
		assert.NoError(t, err)
	})

	t.Run("ReadByIccid", func(t *testing.T) {
		resp, err := c.Read(ctx, &pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: sim.Iccid,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		Imsi = resp.Record.Imsi
	})

	t.Run("InactivateByImsi", func(t *testing.T) {
		_, err := c.Inactivate(ctx, &pb.InactivateReq{
			Id: &pb.InactivateReq_Imsi{
				Imsi: Imsi,
			},
		})
		assert.NoError(t, err)
	})
}
