/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
)

type Mailer interface {
	SendEmail(req *pb.SendEmailRequest) (*pb.SendEmailResponse, error)
	GetEmailById(mailId string) (*pb.GetEmailByIdResponse, error)
}

type mailer struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.MailerServiceClient
	host    string
}

func NewMailer(host string, timeout time.Duration) (*mailer, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Failed to connect to Mailer Service: %v", err)
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
	if m.conn != nil {
		if err := m.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Mailer Service connection: %v", err)
		}
	}
}

func (m *mailer) SendEmail(req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	return m.client.SendEmail(ctx, req)
}

func (m *mailer) GetEmailById(mailerId string) (*pb.GetEmailByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	return m.client.GetEmailById(ctx, &pb.GetEmailByIdRequest{
		MailId: mailerId,
	})
}
