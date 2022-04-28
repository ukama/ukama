//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ukama/ukama/services/common/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/services/cloud/hss/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	HssHost string
	Iccid   string
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{
		HssHost: "localhost:9090",
	}

	config.LoadConfig("integration", testConf)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", testConf)

}

func Test_ImsiService(t *testing.T) {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", testConf.HssHost)
	conn, err := grpc.DialContext(ctx, testConf.HssHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewImsiServiceClient(conn)

	userId := uuid.New()

	// Contact the server and print out its response.
	testImsi := fmt.Sprintf("00000%010d", time.Now().Unix())
	testOrg := fmt.Sprintf("integration-test-org-imsi-service-%s", time.Now().Format("20060102150405"))
	t.Run("AddImis", func(t *testing.T) {
		addResp, err := c.Add(ctx, &pb.AddImsiRequest{Org: testOrg, Imsi: &pb.ImsiRecord{Imsi: testImsi, UserId: userId.String(), Apn: &pb.Apn{Name: "test-apn-name"}}})
		handleResponse(t, err, addResp)
	})

	t.Run("GetImis", func(t *testing.T) {
		getResp, err := c.Get(ctx, &pb.GetImsiRequest{Imsi: testImsi})
		handleResponse(t, err, getResp)
	})

	t.Run("AddGuti", func(t *testing.T) {
		delResp, err := c.AddGuti(ctx, &pb.AddGutiRequest{Guti: &pb.Guti{
			PlmnId: "000001",
			Mmegi:  1,
			Mmec:   1,
			Mtmsi:  uint32(time.Now().Unix()),
		}, Imsi: testImsi,
			UpdatedAt: uint32(time.Now().Unix())})
		handleResponse(t, err, delResp)
	})

	t.Run("UpdateGutiAddedEarlier", func(t *testing.T) {
		delResp, err := c.AddGuti(ctx, &pb.AddGutiRequest{Guti: &pb.Guti{
			PlmnId: "000001",
			Mmegi:  1,
			Mmec:   1,
			Mtmsi:  uint32(time.Now().Unix()) + 1,
		}, Imsi: testImsi,
			UpdatedAt: uint32(time.Now().Unix() + 1)})
		handleResponse(t, err, delResp)
	})

	t.Run("AddTai", func(t *testing.T) {
		resp, err := c.UpdateTai(ctx, &pb.UpdateTaiRequest{Imsi: testImsi, Tac: 4654, PlmnId: "000001",
			UpdatedAt: uint32(time.Now().Unix())})
		handleResponse(t, err, resp)
	})

	t.Run("UpdateTaiAddedEarlier", func(t *testing.T) {
		resp, err := c.UpdateTai(ctx, &pb.UpdateTaiRequest{Imsi: testImsi, Tac: 4654, PlmnId: "000001",
			UpdatedAt: uint32(time.Now().Unix() + 1)})
		handleResponse(t, err, resp)
	})

	t.Run("DeleteImis", func(t *testing.T) {
		delResp, err := c.Delete(ctx, &pb.DeleteImsiRequest{IdOneof: &pb.DeleteImsiRequest_Imsi{
			Imsi: testImsi,
		}})
		handleResponse(t, err, delResp)
	})

	t.Run("UpdateTaiValidationFailure", func(t *testing.T) {
		_, err := c.UpdateTai(ctx, &pb.UpdateTaiRequest{Imsi: "000001111111111", Tac: 4654, PlmnId: "000001",
			UpdatedAt: uint32(time.Now().Unix())})
		s, ok := status.FromError(err)
		assert.True(t, ok, "should be a grpc error")
		assert.Equal(t, codes.NotFound, s.Code(), "should fail with not found")
	})
}

// return false if error is not nil
func handleResponse(t *testing.T, err error, r interface{}) bool {
	fmt.Printf("Response: %v\n", r)
	return assert.NoErrorf(t, err, "Request failed: %v\n", err)
}
