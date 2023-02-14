//go:build integration
// +build integration

package integration

import (
	"context"

	confr "github.com/num30/config"
	"github.com/ukama/ukama/systems/common/config"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/ukama-agent/profile/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig
var Imsi = "012345678912345"
var Iccid = "123456789123456789123"
var Network = "db081ef5-a8ae-4a95-bff3-a7041d52bb9b"
var Org = "abdc0cec-5553-46aa-b3a8-1e31b0ef58ad"
var Package = "fab4f98d-2e82-47e8-adb5-e516346880d8"
var NodePolicyPath = "v1/epc/pcrf"
var MonitoringPeriod time.Duration = 10 * time.Second
var profile = db.Profile{
	Iccid:                   Iccid,
	Imsi:                    Imsi,
	UeDlBps:                 10000000,
	UeUlBps:                 1000000,
	ApnName:                 "ukama",
	AllowedTimeOfService:    2592000,
	TotalDataBytes:          1024000,
	ConsumedDataBytes:       0,
	NetworkId:               uuid.FromStringOrNil(Network),
	PackageId:               uuid.FromStringOrNil(Package),
	LastStatusChangeReasons: db.ACTIVATION,
	LastStatusChangeAt:      time.Now(),
}

var pack = db.PackageDetails{
	PackageId:            uuid.FromStringOrNil(Package),
	UeDlBps:              10000000,
	UeUlBps:              1000000,
	ApnName:              "ukama",
	AllowedTimeOfService: time.Second * 2592000,
	TotalDataBytes:       1024000,
	ConsumedDataBytes:    0,
	LastStatusChangeAt:   time.Now(),
}

func init() {
	tConfig = &TestConfig{}
	r := confr.NewConfReader("integration")
	r.Read(tConfig)

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", tConfig)
}

func CreateProfileClient() (*grpc.ClientConn, pb.ProfileServiceClient, error) {
	log.Infoln("Connecting to Profile ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewProfileServiceClient(conn)
	return conn, c, nil
}

func Test_FullFlow(t *testing.T) {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateProfileClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	t.Run("Add", func(t *testing.T) {
		reqPb := pb.AddReq{
			Profile: &pb.Profile{
				Iccid:   profile.Iccid,
				Imsi:    profile.Imsi,
				UeDlBps: profile.UeDlBps,
				UeUlBps: profile.UeUlBps,
				Apn: &pb.Apn{
					Name: profile.ApnName,
				},
				NetworkId:            profile.NetworkId.String(),
				PackageId:            profile.PackageId.String(),
				AllowedTimeOfService: profile.AllowedTimeOfService,
				TotalDataBytes:       profile.TotalDataBytes,
				ConsumedDataBytes:    profile.ConsumedDataBytes,
				LastChange:           db.ACTIVATION.String(),
				LastChangeAt:         profile.LastStatusChangeAt.Unix(),
			},
		}
		_, err := c.Add(ctx, &reqPb)
		assert.NoError(t, err)
	})

	t.Run("ReadByIccid", func(t *testing.T) {
		resp, err := c.Read(ctx, &pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: Iccid,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, resp.Profile.Iccid, Iccid)
	})

	t.Run("UpdateUsage", func(t *testing.T) {
		_, err := c.UpdateUsage(ctx, &pb.UpdateUsageReq{
			Imsi:              profile.Imsi,
			ConsumedDataBytes: 1000,
		})

		assert.NoError(t, err)
	})

	t.Run("UpdatePackage", func(t *testing.T) {
		_, err := c.UpdatePackage(ctx, &pb.UpdatePackageReq{
			Iccid: Iccid,
			Package: &pb.Package{
				UeDlBps: profile.UeDlBps,
				UeUlBps: profile.UeUlBps,
				Apn: &pb.Apn{
					Name: profile.ApnName,
				},
				PackageId:            profile.PackageId.String(),
				AllowedTimeOfService: profile.AllowedTimeOfService,
				TotalDataBytes:       profile.TotalDataBytes,
				ConsumedDataBytes:    profile.ConsumedDataBytes,
			},
		})

		assert.NoError(t, err)
	})

	t.Run("RemoveByIccid", func(t *testing.T) {
		_, err := c.Remove(ctx, &pb.RemoveReq{
			Id: &pb.RemoveReq_Iccid{
				Iccid: Iccid,
			},
		})

		assert.NoError(t, err)
	})

	t.Run("Add", func(t *testing.T) {
		reqPb := pb.AddReq{
			Profile: &pb.Profile{
				Iccid:   profile.Iccid,
				Imsi:    profile.Imsi,
				UeDlBps: profile.UeDlBps,
				UeUlBps: profile.UeUlBps,
				Apn: &pb.Apn{
					Name: profile.ApnName,
				},
				NetworkId:            profile.NetworkId.String(),
				PackageId:            profile.PackageId.String(),
				AllowedTimeOfService: profile.AllowedTimeOfService,
				TotalDataBytes:       profile.TotalDataBytes,
				ConsumedDataBytes:    profile.ConsumedDataBytes,
				LastChange:           db.ACTIVATION.String(),
				LastChangeAt:         profile.LastStatusChangeAt.Unix(),
			},
		}
		_, err := c.Add(ctx, &reqPb)
		assert.NoError(t, err)
	})

	t.Run("ReadByImsi", func(t *testing.T) {
		resp, err := c.Read(ctx, &pb.ReadReq{
			Id: &pb.ReadReq_Imsi{
				Imsi: Imsi,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, resp.Profile.Imsi, Imsi)
	})

	t.Run("Sync", func(t *testing.T) {
		_, err := c.Sync(ctx, &pb.SyncReq{
			Iccid: []string{Iccid},
		})

		assert.NoError(t, err)
	})

	t.Run("RemoveByImsi", func(t *testing.T) {
		_, err := c.Remove(ctx, &pb.RemoveReq{
			Id: &pb.RemoveReq_Imsi{
				Imsi: Imsi,
			},
		})

		assert.NoError(t, err)
	})

}
