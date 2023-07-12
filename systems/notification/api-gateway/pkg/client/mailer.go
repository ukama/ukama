package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	emailPkg "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
)

type Mailer interface {
	SendEmail(*emailPkg.SendEmailRequest) (*emailPkg.SendEmailResponse, error)
	GetEmailById(*emailPkg.GetEmailByIdRequest) (*emailPkg.GetEmailByIdResponse, error)
}

type mailer struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.MailerServiceClient
	host    string
}

func NewMailer(host string, timeout time.Duration) (*mailer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
		return nil, err
	}

	client := pb.NewMailerServiceClient(conn)

	return &mailer{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}, nil
}

func NewMailerFromClient(mailerClient pb.MailerServiceClient) *mailer {
	return &mailer{
		host:    "localhost",
		timeout: 10 * time.Second,
		conn:    nil,
		client:  mailerClient,
	}
}

func (m *mailer) Close() {
	m.conn.Close()
}

func (m *mailer) SendEmail(req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	res, err := m.client.SendEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *mailer) GetEmailById(req *pb.GetEmailByIdRequest) (*pb.GetEmailByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	res, err := m.client.GetEmailById(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
