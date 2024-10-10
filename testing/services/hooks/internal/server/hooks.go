/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/testing/services/hooks/internal"
	"github.com/ukama/ukama/testing/services/hooks/internal/clients"
	"github.com/ukama/ukama/testing/services/hooks/internal/scheduler"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/testing/services/hooks/pb/gen"
)

const (
	HookTaskTag          = "Hook Response"
	queryPattern         = "?status=processing"
	webhookStatusPending = "ACCEPTED"
)

type HookServer struct {
	orgName        string
	pawapayClient  clients.PawapayClient
	paymentsClient clients.PaymentsClient
	webhooksClient clients.WebhooksClient
	cdrScheduler   scheduler.HookScheduler
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedHookServiceServer
}

func NewHookServer(orgName string, pawapayClient clients.PawapayClient, paymentsClient clients.PaymentsClient,
	webhooksClient clients.WebhooksClient, cdrScheduler scheduler.HookScheduler) *HookServer {
	h := &HookServer{
		orgName:        orgName,
		pawapayClient:  pawapayClient,
		paymentsClient: paymentsClient,
		webhooksClient: webhooksClient,
		cdrScheduler:   cdrScheduler,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(internal.SystemName).SetOrgName(orgName).SetService(internal.ServiceName),
	}

	_, err := h.startScheduler()
	if err != nil {
		log.Warnf("failed to auto start webhook scheduler. You need to start RPC manually. Err: %v",
			err)
	}

	return h
}

func (p *HookServer) StartScheduler(ctx context.Context, req *pb.StartRequest) (*pb.StartResponse, error) {
	return p.startScheduler()
}

func (p *HookServer) StopScheduler(ctx context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
	err := p.cdrScheduler.Stop()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"an unexpected error has occured while stoping scheduler: %v", err)
	}

	return &pb.StopResponse{}, nil
}

func (p *HookServer) setActiveTaskFunc() (string, any, any) {
	return HookTaskTag, p.pullHooksResponse, ""
}

func (p *HookServer) pullHooksResponse(placeHolder string) error {
	log.Infof("Pulling webhooks from sandbox")

	payments, err := p.paymentsClient.ListPayments(queryPattern)
	if err != nil {
		return fmt.Errorf("error while listing pending payments: %w", err)
	}

	for _, payment := range payments {
		depositWebhook, err := p.pawapayClient.GetDeposit(payment.Id)
		if err != nil {
			log.Errorf("error while fetching deposit webhook: %v", err)
			log.Warn("skipping posting empty payload")

			continue
		}

		if depositWebhook.Status != webhookStatusPending {
			_, err = p.webhooksClient.PostDepositHook(depositWebhook)
			if err != nil {
				log.Errorf("error while posting deposit webhook: %v", err)
			}
		}
	}

	log.Infof("finished Pulling webhooks from sandbox")

	return nil
}

func (p *HookServer) startScheduler() (*pb.StartResponse, error) {
	err := p.cdrScheduler.Start(p.setActiveTaskFunc())
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"an unexpected error has occured while starting scheduler: %v", err)
	}

	return &pb.StartResponse{}, nil
}
