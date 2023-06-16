package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Mailer struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.MaillingServiceClient
	host    string
}

func NewMailer(host string, timeout time.Duration) *Mailer {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewMaillingServiceClient(conn)

	return &Mailer{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewMailerFromClient(MailerClient pb.MaillingServiceClient) *Mailer {
	return &Mailer{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  MailerClient,
	}
}

func (m *Mailer) Close() {
	m.conn.Close()
}

func (m *Mailer) SendEmail(req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	res, err := m.client.SendEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
