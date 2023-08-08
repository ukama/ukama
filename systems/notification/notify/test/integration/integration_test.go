//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
	jdb "gorm.io/datatypes"
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
	node := ukama.NewVirtualHomeNodeId().String()

	nodeAlert := NewTestDbNotification(node, "alert")
	nodeAlertResp := &pb.AddResponse{}

	nodeEvent := NewTestDbNotification(node, "event")
	nodeEventResp := &pb.AddResponse{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateInvoiceServiceClient()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}
	defer conn.Close()

	t.Run("AddAlertNotification", func(t *testing.T) {
		var err error

		nodeAlertResp, err = c.Add(ctx, &pb.AddRequest{
			NodeId:      nodeAlert.NodeId,
			Severity:    nodeAlert.Severity.String(),
			Type:        nodeAlert.Type.String(),
			ServiceName: nodeAlert.ServiceName,
			Status:      nodeAlert.Status,
			EpochTime:   nodeAlert.Time,
			Description: nodeAlert.Description,
			Details:     nodeAlert.Details.String(),
		})

		assert.NoError(t, err)
	})

	t.Run("GetAlertNotification", func(t *testing.T) {
		nt, err := c.Get(ctx, &pb.GetRequest{
			NotificationId: nodeAlertResp.Notification.Id})

		assert.NoError(t, err)
		assert.NotNil(t, nt)
	})

	t.Run("AddEventNotification", func(t *testing.T) {
		var err error

		nodeEventResp, err = c.Add(ctx, &pb.AddRequest{
			NodeId:      nodeEvent.NodeId,
			Severity:    nodeEvent.Severity.String(),
			Type:        nodeEvent.Type.String(),
			ServiceName: nodeEvent.ServiceName,
			Status:      nodeEvent.Status,
			EpochTime:   nodeEvent.Time,
			Description: nodeEvent.Description,
			Details:     nodeEvent.Details.String(),
		})

		assert.NoError(t, err)
	})

	t.Run("GetEventNotification", func(t *testing.T) {
		nt, err := c.Get(ctx, &pb.GetRequest{
			NotificationId: nodeEventResp.Notification.Id})

		assert.NoError(t, err)
		assert.NotNil(t, nt)
	})

	t.Run("ListAll", func(t *testing.T) {
		list, err := c.List(ctx, &pb.ListRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, list)
	})

	t.Run("ListAlertsForNode", func(t *testing.T) {
		ntype := "alert"

		list, err := c.List(ctx, &pb.ListRequest{
			NodeId: node,
			Type:   ntype,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
	})

	t.Run("ListAlertsForService", func(t *testing.T) {
		service := "noded"
		ntype := "alert"

		list, err := c.List(ctx, &pb.ListRequest{
			ServiceName: service,
			Type:        ntype,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
	})

	t.Run("ListEventsForNode", func(t *testing.T) {
		ntype := "event"

		list, err := c.List(ctx, &pb.ListRequest{
			NodeId: node,
			Type:   ntype,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
	})

	t.Run("ListEventsForService", func(t *testing.T) {
		service := "noded"
		ntype := "event"

		list, err := c.List(ctx, &pb.ListRequest{
			ServiceName: service,
			Type:        ntype,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
	})

	t.Run("DeleteAlertNotification", func(tt *testing.T) {
		resp, err := c.Delete(ctx,
			&pb.GetRequest{NotificationId: nodeAlertResp.Notification.Id})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetAlertNotification", func(t *testing.T) {
		nt, err := c.Get(ctx, &pb.GetRequest{
			NotificationId: nodeAlertResp.Notification.Id})

		assert.Error(t, err)
		assert.Nil(t, nt)
	})

	t.Run("DeleteEventNotification", func(tt *testing.T) {
		resp, err := c.Delete(ctx,
			&pb.GetRequest{NotificationId: nodeEventResp.Notification.Id})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetEventNotification", func(t *testing.T) {
		nt, err := c.Get(ctx, &pb.GetRequest{
			NotificationId: nodeEventResp.Notification.Id})

		assert.Error(t, err)
		assert.Nil(t, nt)
	})
}

func CreateInvoiceServiceClient() (*grpc.ClientConn, pb.NotifyServiceClient, error) {
	log.Infoln("Connecting to Invoice Server ", tConfig.ServiceHost)

	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewNotifyServiceClient(conn)

	return conn, c, nil
}

func NewTestDbNotification(nodeId string, ntype string) db.Notification {
	return db.Notification{
		Id:          uuid.NewV4(),
		NodeId:      nodeId,
		NodeType:    *ukama.GetNodeType(nodeId),
		Severity:    db.SeverityType("high"),
		Type:        db.NotificationType(ntype),
		ServiceName: "noded",
		Status:      8200,
		Time:        uint32(time.Now().Unix()),
		Description: "Some random alert",
		Details:     jdb.JSON(`{"reason": "testing", "component":"router_test"}`),
	}
}
