//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/ukama/ukama/systems/common/config"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/registry/users/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TestConfig struct {
	UsersHost string
	Iccid     string
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{
		UsersHost: "localhost:9090",
	}

	config.LoadConfig("integration", testConf)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", testConf)

}

func Test_UserService(t *testing.T) {
	if testConf.Iccid == "" {
		assert.FailNow(t, "Missing ICCID. Set ICCID env var")
	}

	// One timeout for whole test
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	logrus.Infoln("Connecting to service ", testConf.UsersHost)
	conn, err := grpc.DialContext(ctx, testConf.UsersHost, grpc.WithInsecure(), grpc.WithBlock())
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
	defer cleanupUser(c, addResp)

	var esimUsr *pb.AddResponse
	t.Run("AddUserWithESim", func(tt *testing.T) {
		esimUsr, err = c.Add(ctx, &pb.AddRequest{
			User: &pb.User{
				Email: "test@example.com",
				Name:  "Joe Esim",
			},
			Org: testOrg})

		if handleResponse(tt, err, esimUsr) {
			logrus.Info("Failed to add")
			assert.NotEmpty(tt, esimUsr.User.Uuid)
		} else {
			t.FailNow()
		}
	})
	defer cleanupUser(c, esimUsr)

	// todo: test limit
	t.Run("list", func(tt *testing.T) {
		listResp, err := c.List(ctx, &pb.ListRequest{
			Org: testOrg,
		})

		if handleResponse(tt, err, listResp) {
			assert.Equal(tt, 2, len(listResp.Users))
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
			assert.Equal(tt, 1, len(listResp.Users))
		}
	})
}

func cleanupUser(c pb.UserServiceClient, addResp ...*pb.AddResponse) {
	logrus.Info("Cleaning up")
	for _, rsp := range addResp {
		if rsp != nil && rsp.User != nil {
			r := *rsp

			logrus.Info("Deleting user: ", r.User.Uuid, " Iccid: ", r.Iccid)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			_, err := c.Delete(ctx, &pb.DeleteRequest{UserId: rsp.User.Uuid})
			if err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
					return
				}
				logrus.Errorf("Failed to delete user %s: %v", r.User.Uuid, err)
			}
		}
	}

}

// return false if error is not nil
func handleResponse(t *testing.T, err error, r interface{}) bool {
	fmt.Printf("Response: %v\n", r)
	return assert.NoErrorf(t, err, "Request failed: %v\n", err)
}
