package integration

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/ukama/ukama/systems/common/config"
// 	pb "github.com/ukama/ukama/systems/node/software-manager/pb/gen"

// 	rconf "github.com/num30/config"
// 	log "github.com/sirupsen/logrus"
// 	grpc "google.golang.org/grpc"
// )

// var tConfig *TestConfig

// func init() {
// 	// load config
// 	tConfig = &TestConfig{}

// 	reader := rconf.NewConfReader("integration")

// 	err := reader.Read(tConfig)
// 	if err != nil {
// 		log.Fatalf("Failed to read config: %v", err)
// 	}

// 	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
// 	log.Infof("Config: %+v\n", tConfig)
// }

// type TestConfig struct {
// 	ServiceHost string        `default:"localhost:9090"`
// 	Queue       *config.Queue `default:"{}"`
// 	OrgId       string
// 	OrgName     string
// }

// func Test_FullFlow(t *testing.T) {

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	log.Infoln("Connecting to sw-manager ", tConfig.ServiceHost)

// 	conn, err := grpc.DialContext(ctx, tConfig.ServiceHost, grpc.WithInsecure(), grpc.WithBlock())
// 	if err != nil {
// 		assert.NoError(t, err, "did not connect: %v", err)

// 		return
// 	}

// 	c := pb.NewSoftwareManagerServiceClient(conn)

// 	var r interface{}
// 	t.Run("CreateSoftwareUpdate", func(tt *testing.T) {
// 		r, err = c.CreateSoftwareUpdate(ctx, &pb.CreateSoftwareUpdateRequest{
// 			Name:        "test",
// 			Version:     "1.0.0",
// 			ReleaseDate: "2021-01-01",
// 			Status:      pb.Status(1),
// 		})

// 		handleResponse(tt, err, r)
// 	})

// 	t.Run("GetLatestSoftwareUpdate", func(tt *testing.T) {
// 		r, err = c.GetLatestSoftwareUpdate(ctx, &pb.GetLatestSoftwareUpdateRequest{
// 		})

// 		handleResponse(tt, err, r)
// 	})

// 	t.Run("ListSoftwareUpdates", func(tt *testing.T) {
// 		r, err = c.ListSoftwareUpdates(ctx, &pb.ListSoftwareUpdatesRequest{})

// 		handleResponse(tt, err, r)
// 	})

// }

// func handleResponse(t *testing.T, err error, r interface{}) {
// 	t.Helper()

// 	log.Printf("Response: %v\n", r)

// 	if err != nil {
// 		assert.FailNow(t, "Request failed: %v\n", err)
// 	}
// }