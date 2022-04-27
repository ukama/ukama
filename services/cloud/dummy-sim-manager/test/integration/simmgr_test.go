//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/ukama/ukamaX/common/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen/simmgr"
	"google.golang.org/grpc"
)

var testConf *TestConf

type TestConf struct {
	Iccid          string
	SimManagerHost string
}

func init() {
	testConf = &TestConf{
		Iccid:          "890000000000000001",
		SimManagerHost: "localhost:9090",
	}

	config.LoadConfig("integration", testConf)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SIMMANAGERHOST")
	logrus.Infof("%+v", testConf)
}

func Test_GetSimInfo(t *testing.T) {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	logrus.Infoln("Connecting to service ", testConf.SimManagerHost)
	conn, err := grpc.DialContext(ctx, testConf.SimManagerHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewSimManagerServiceClient(conn)

	t.Run("GetSimInfo", func(t *testing.T) {
		r, err := c.GetSimInfo(ctx, &pb.GetSimInfoRequest{Iccid: testConf.Iccid})
		if assert.NoError(t, err) {
			if r != nil {
				assert.NotEmpty(t, r.Imsi)
			}
		}
	})

	t.Run("GetSimStatus", func(t *testing.T) {
		r, err := c.GetSimStatus(ctx, &pb.GetSimStatusRequest{Iccid: testConf.Iccid})
		if assert.NoError(t, err) {
			fmt.Printf("%+v", r)
			if r != nil {
				assert.NotEqual(t, pb.GetSimStatusResponse_UNKNOWN, r.Status)
			}
		}
	})

	t.Run("SetServiceStatus", func(t *testing.T) {
		_, err := c.SetServiceStatus(ctx, &pb.SetServiceStatusRequest{Iccid: testConf.Iccid, Services: &pb.Services{
			Sms: wrapperspb.Bool(false),
		}})
		if assert.NoError(t, err) {
			r, err := c.GetSimStatus(ctx, &pb.GetSimStatusRequest{Iccid: testConf.Iccid})
			if assert.NoError(t, err) {
				if r != nil {
					assert.Equal(t, false, r.Services.GetSms().Value)
				}
			}
		}

		_, err = c.SetServiceStatus(ctx, &pb.SetServiceStatusRequest{Iccid: testConf.Iccid, Services: &pb.Services{
			Sms: wrapperspb.Bool(true),
		}})

		assert.NoError(t, err, "Error resetting service status")
	})

	t.Run("TerminateSim", func(t *testing.T) {
		_, err := c.TerminateSim(ctx, &pb.TerminateSimRequest{Iccid: testConf.Iccid})
		if assert.NoError(t, err) {
			_, err := c.GetSimStatus(ctx, &pb.GetSimStatusRequest{Iccid: testConf.Iccid})
			s := status.Code(err)
			assert.Equal(t, codes.NotFound, s)
		}
	})
}
