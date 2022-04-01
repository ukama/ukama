//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/ukama/ukamaX/common/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
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

	// Contact the server and print out its response.
	testImsi := fmt.Sprintf("00000%010d", time.Now().Unix())
	testOrg := fmt.Sprintf("integration-test-org-imsi-service-%s", time.Now().Format("20060102150405"))
	t.Run("AddImis", func(t *testing.T) {
		addResp, err := c.Add(ctx, &pb.AddImsiRequest{Org: testOrg, Imsi: &pb.ImsiRecord{Imsi: testImsi, Apn: &pb.Apn{Name: "test-apn-name"}}})
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
		delResp, err := c.Delete(ctx, &pb.DeleteImsiRequest{Imsi: testImsi, Org: testOrg})
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

func Test_UserService(t *testing.T) {
	if testConf.Iccid == "" {
		assert.FailNow(t, "Missing ICCID")
	}

	// One timeout for whole test
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	logrus.Infoln("Connecting to service ", testConf.HssHost)
	conn, err := grpc.DialContext(ctx, testConf.HssHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)
	testOrg := fmt.Sprintf("integration-test-org-user-service-%s", time.Now().Format("20060102150405"))

	var addResp *pb.AddResponse

	testIccid := testConf.Iccid

	simToken := ""
	t.Run("GenerateSimToken", func(tt *testing.T) {
		r, err := c.GenerateSimToken(ctx, &pb.GenerateSimTokenRequest{Iccid: testIccid, FromPool: false})
		if !assert.NoError(tt, err) {
			assert.FailNow(t, "Error creating sim token")
		}
		simToken = r.SimToken
	})

	t.Run("AddWithCode", func(tt *testing.T) {
		addResp, err = c.Add(ctx, &pb.AddRequest{
			User: &pb.User{
				Email: "test@example.com",
				Name:  "Joe",
			},
			SimToken: simToken,
			Org:      testOrg})

		if handleResponse(tt, err, addResp) {
			logrus.Info("Failed to add")
			assert.NotEmpty(tt, addResp.User.Uuid)
		} else {
			t.FailNow()
		}
	})

	// todo: test limit
	t.Run("list", func(tt *testing.T) {
		listResp, err := c.List(ctx, &pb.ListUsersRequest{
			Org: testOrg,
		})

		if handleResponse(tt, err, listResp) {
			assert.Equal(tt, 1, len(listResp.Users))
		}
	})

	t.Run("get", func(tt *testing.T) {
		getResp, err := c.Get(ctx, &pb.GetUserRequest{Uuid: addResp.User.Uuid})
		if handleResponse(tt, err, getResp) {
			assert.NotNil(tt, getResp.Sim)
			assert.NotEqual(tt, getResp.Sim.Carrier.Status, pb.SimStatus_UNKNOWN)
		}
	})

	t.Run("Delete", func(tt *testing.T) {
		getResp, err := c.Delete(ctx, &pb.DeleteUserRequest{Uuid: addResp.User.Uuid})

		if !handleResponse(tt, err, getResp) {
			assert.FailNow(tt, "")
		}

		// make sure that user is deleted
		listResp, err := c.List(ctx, &pb.ListUsersRequest{
			Org: testOrg,
		})
		if handleResponse(tt, err, listResp) {
			assert.Equal(tt, 0, len(listResp.Users))
		}
	})
}

// return false if error is not nil
func handleResponse(t *testing.T, err error, r interface{}) bool {
	fmt.Printf("Response: %v\n", r)
	return assert.NoErrorf(t, err, "Request failed: %v\n", err)
}
