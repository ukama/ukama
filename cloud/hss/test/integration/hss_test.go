//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ukama/ukamaX/common/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func Test_UserService(t *testing.T) {
	if testConf.Iccid == "" {
		assert.FailNow(t, "Missing ICCID. Set ICCID env var")
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
	defer cleanupUser(addResp, c)

	// todo: test limit
	t.Run("list", func(tt *testing.T) {
		listResp, err := c.List(ctx, &pb.ListRequest{
			Org: testOrg,
		})

		if handleResponse(tt, err, listResp) {
			assert.Equal(tt, 1, len(listResp.Users))
		}
	})

	getResp := &pb.GetResponse{}
	t.Run("get", func(tt *testing.T) {
		getResp, err = c.Get(ctx, &pb.GetRequest{UserId: addResp.User.Uuid})
		if handleResponse(tt, err, getResp) {
			assert.NotNil(tt, getResp.Sim)
			assert.NotEqual(tt, getResp.Sim.Carrier.Status, pb.SimStatus_UNKNOWN)
			assert.Equal(tt, false, getResp.User.IsDeactivated)
		}
	})

	t.Run("update ", func(tt *testing.T) {
		_, err := c.Update(ctx, &pb.UpdateRequest{
			UserId: addResp.User.Uuid,
			User: &pb.UserAttributes{
				Name:  "changed",
				Email: "changed@example.com",
				Phone: "1231223132",
			},
		})
		if !assert.NoError(tt, err) {
			assert.FailNow(tt, "update test failed")
			return
		}

		getResp, err := c.Get(ctx, &pb.GetRequest{UserId: addResp.User.Uuid})
		if handleResponse(tt, err, getResp) {
			assert.Equal(tt, "changed", getResp.User.Name)
			assert.Equal(tt, "1231223132", getResp.User.Phone)
			assert.Equal(tt, "changed@example.com", getResp.User.Email)
		}
	})

	t.Run("setCarrierServiceStatuses", func(tt *testing.T) {
		targetData := getResp.Sim.Carrier.Services.Data
		setResp, err := c.SetSimStatus(ctx, &pb.SetSimStatusRequest{
			Iccid: getResp.Sim.Iccid,
			Carrier: &pb.SetSimStatusRequest_SetServices{
				Data: wrapperspb.Bool(targetData),
			},
		})

		if handleResponse(tt, err, setResp) {
			getResp, err = c.Get(ctx, &pb.GetRequest{UserId: addResp.User.Uuid})
			assert.NoError(tt, err)
			assert.Equal(tt, targetData, getResp.Sim.Carrier.Services.Data)
		}
	})

	t.Run("DeactivateUser", func(tt *testing.T) {
		_, err = c.DeactivateUser(ctx, &pb.DeactivateUserRequest{
			UserId: addResp.User.Uuid,
		})
		if !assert.NoError(tt, err) {
			assert.FailNow(tt, "DeactivateUser test failed")
			return
		}

		getResp, err = c.Get(ctx, &pb.GetRequest{UserId: addResp.User.Uuid})
		if handleResponse(tt, err, getResp) {
			assert.Equal(tt, true, getResp.User.IsDeactivated)
		}
	})

	t.Run("Delete", func(tt *testing.T) {
		_, err = c.Delete(ctx, &pb.DeleteRequest{UserId: addResp.User.Uuid})

		if !handleResponse(tt, err, getResp) {
			assert.FailNow(tt, "")
		}

		// make sure that user is deleted
		listResp, err := c.List(ctx, &pb.ListRequest{
			Org: testOrg,
		})
		if handleResponse(tt, err, listResp) {
			assert.Equal(tt, 0, len(listResp.Users))
		}
	})
}

func cleanupUser(addResp *pb.AddResponse, c pb.UserServiceClient) {
	if addResp != nil && addResp.User != nil {
		r := *addResp

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err := c.Delete(ctx, &pb.DeleteRequest{UserId: addResp.User.Uuid})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				return
			}
			logrus.Errorf("Failed to delete user %s: %v", r.User.Uuid, err)
		}
	}
}

// return false if error is not nil
func handleResponse(t *testing.T, err error, r interface{}) bool {
	fmt.Printf("Response: %v\n", r)
	return assert.NoErrorf(t, err, "Request failed: %v\n", err)
}
