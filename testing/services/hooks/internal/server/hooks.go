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
	"github.com/ukama/ukama/systems/common/util/payments"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/testing/services/hooks/internal"
	"github.com/ukama/ukama/testing/services/hooks/internal/clients"
	"github.com/ukama/ukama/testing/services/hooks/internal/scheduler"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/testing/services/hooks/pb/gen"
)

const (
	HookTaskTag  = "Hook Response"
	queryPattern = "?status=processing"

	pawapayHookStatusPending          = "ACCEPTED"
	stripeStatusRequiresPaymentMethod = "requires_payment_method"
)

type HookServer struct {
	orgName        string
	pawapayClient  clients.PawapayClient
	stripeClient   clients.StripeClient
	paymentsClient clients.PaymentsClient
	webhooksClient clients.WebhooksClient
	cdrScheduler   scheduler.HookScheduler
	msgBus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedHookServiceServer
}

func NewHookServer(orgName string, pawapayClient clients.PawapayClient, stripeClient clients.StripeClient, paymentsClient clients.PaymentsClient,
	webhooksClient clients.WebhooksClient, cdrScheduler scheduler.HookScheduler, msgBus mb.MsgBusServiceClient) *HookServer {
	h := &HookServer{
		orgName:        orgName,
		pawapayClient:  pawapayClient,
		stripeClient:   stripeClient,
		paymentsClient: paymentsClient,
		webhooksClient: webhooksClient,
		cdrScheduler:   cdrScheduler,
		msgBus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(internal.SystemName).SetOrgName(orgName).SetService(internal.ServiceName),
	}

	_, err := h.startScheduler()
	if err != nil {
		log.Warnf("Failed to auto start webhook scheduler. You need to start RPC manually. Err: %v",
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
		if _, err = uuid.FromString(payment.ExternalId); err == nil {
			depositWebhook, err := fetchPawapayDeposit(payment.ExternalId, p)
			if err != nil {
				log.Errorf("Error while fetching deposit: %v", err)
				log.Warn("Skipping posting empty payload")

				continue
			}

			err = postPawapayDeposit(depositWebhook, p)
			if err != nil {
				log.Errorf("Error while making post deposit webhook: %v", err)
			}
		} else {
			intentWebhook, err := fetchStripePaymentIntent(payment.ExternalId, p)
			if err != nil {
				log.Errorf("Error while fetching intent: %v", err)
				log.Warn("Skipping posting empty payload")

				continue
			}

			err = postStripePaymentIntent(intentWebhook, p)
			if err != nil {
				log.Errorf("Error while making post payment intent webhook: %v", err)
			}
		}
	}

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

func fetchPawapayDeposit(externalId string, p *HookServer) (*payments.Deposit, error) {
	depositWebhook, err := p.pawapayClient.GetDeposit(externalId)
	if err != nil {
		return nil, err
	}

	return depositWebhook, nil
}

func fetchStripePaymentIntent(externalId string, p *HookServer) (*payments.Intent, error) {
	intentHook, err := p.stripeClient.GetPaymentIntent(externalId)
	if err != nil {
		return nil, err
	}

	return intentHook, nil
}

func postPawapayDeposit(depositWebhook *payments.Deposit, p *HookServer) error {
	if depositWebhook.Status != pawapayHookStatusPending {
		_, err := p.webhooksClient.PostDepositHook(depositWebhook)
		if err != nil {
			return err
		}
	}

	return nil
}

func postStripePaymentIntent(intentHook *payments.Intent, p *HookServer) error {
	if intentHook.Status != stripeStatusRequiresPaymentMethod {
		log.Infof("New status update (%v), will post payment hook...", intentHook.Status)
		_, err := p.webhooksClient.PostPaymentIntentHook(intentHook.PaymentIntent)
		if err != nil {
			return err
		}
	}

	return nil
}
