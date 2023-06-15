package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"
)

type MaillingServer struct {
	maillingRepoRepo       db.MaillingRepo
pb.UnimplementedMaillingServiceServer
mailer     *config.Mailer
}

func NewMaillingServer(maillingRepoRepo db.MaillingRepo) *MaillingServer {
	return &MaillingServer{maillingRepoRepo: maillingRepoRepo,
}
}

func (s *MaillingServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {

	return &pb.SendEmailResponse{}, nil
}
