package server

import (
	"context"

	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"
)

type MaillingServer struct {
	maillingRepoRepo       db.MaillingRepo
pb.UnimplementedMaillingServiceServer
}

func NewMaillingServer(maillingRepoRepo db.MaillingRepo) *MaillingServer {
	return &MaillingServer{maillingRepoRepo: maillingRepoRepo,
}
}

func (s *MaillingServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {

	return &pb.SendEmailResponse{}, nil
}
