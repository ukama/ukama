package integration

import (
	"context"
	"testing"
	"time"

	confr "github.com/num30/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

func TestSendEmailAndGetEmail(t *testing.T) {

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	logrus.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateMailerClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}
	t.Run("SendEmail", func(t *testing.T) {
		sendEmailResponse, err := c.SendEmail(ctx, &pb.SendEmailRequest{
			To:      []string{"brackley@ukama.com"},
			Subject: "test",
			Body:    "test",
		})
		assert.NoError(t, err)
		r, err := c.GetEmailById(ctx, &pb.GetEmailByIdRequest{
			MailId: sendEmailResponse.MailId,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, sendEmailResponse.MailId, r.MailId)
		}

	})

	
	
	

}

func CreateMailerClient() (*grpc.ClientConn, pb.MailerServiceClient, error) {
	logrus.Infoln("Connecting to notification-mailer ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewMailerServiceClient(conn)
	return conn, c, nil
}
