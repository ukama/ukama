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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/datatypes"

	"github.com/ukama/ukama/systems/billing/report/pkg"
	"github.com/ukama/ukama/systems/billing/report/pkg/db"
	"github.com/ukama/ukama/systems/billing/report/pkg/util"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	receiptNumberPrefix  = "RCPT-"
	receiptStatus        = "finalized"
	receiptPaymentStatus = "succeeded"
	receiptFileURLFormat = "http://{API_ENDPOINT}/pdf/%s.pdf"
)

// TODO: We need to think about retry policies for failing interaction between
// TODO: We have unmarshal methods in common/pb/gen/events for all the event messages. We should use those.
// our backend and the upstream billing service provider.

type ReportEventServer struct {
	orgName        string
	orgId          string
	reportRepo     db.ReportRepo
	msgBus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	epb.UnimplementedEventNotificationServiceServer
}

func NewReportEventServer(orgName, orgId string, reportRepo db.ReportRepo,
	msgBus mb.MsgBusServiceClient) *ReportEventServer {
	return &ReportEventServer{
		orgName:    orgName,
		orgId:      orgId,
		reportRepo: reportRepo,
		msgBus:     msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (r *ReportEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(r.orgName, "event.cloud.local.{{ .Org}}.payments.processor.payment.success"):
		msg, err := unmarshalPayment(e.Msg)
		if err != nil {
			return nil, err
		}

		err = r.handlePaymentSuccessEvent(e.RoutingKey, msg, r)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(r.orgName, "event.cloud.local.{{ .Org}}.payments.processor.payment.update"):
		msg, err := unmarshalPayment(e.Msg)
		if err != nil {
			return nil, err
		}

		err = r.handlePaymentUpdateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (r *ReportEventServer) handlePaymentSuccessEvent(key string, msg *epb.Payment,
	b *ReportEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	if ukama.ParseItemType(msg.ItemType) == ukama.ItemTypePackage {
		return addReceipt(msg, r.orgId, r.reportRepo, r.msgBus, r.baseRoutingKey)
	}

	_, err := update(msg.ItemId, true, msg.TransactionId, r.reportRepo, r.msgBus, r.baseRoutingKey)

	return err
}

func (r *ReportEventServer) handlePaymentUpdateEvent(key string, msg *epb.Payment) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	if ukama.ParseStatusType(msg.Status) != ukama.StatusTypeCompleted {
		log.Infof("Skipping receipt for payment %s with non completed status: %s", msg.Id, msg.Status)

		return nil
	}

	if ukama.ParseItemType(msg.ItemType) != ukama.ItemTypePackage {
		log.Infof("Skipping receipt for payment %s with item type: %s", msg.Id, msg.ItemType)

		return nil
	}

	return addReceipt(msg, r.orgId, r.reportRepo, r.msgBus, r.baseRoutingKey)
}

func addReceipt(msg *epb.Payment, orgId string, reportRepo db.ReportRepo,
	msgBus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	paymentId, err := uuid.FromString(msg.Id)
	if err != nil {
		return fmt.Errorf("invalid format for payment uuid: %s. Error %w", msg.Id, err)
	}

	ownerId, err := uuid.FromString(orgId)
	if err != nil {
		return fmt.Errorf("invalid format for org uuid: %s. Error %w", orgId, err)
	}

	_, err = reportRepo.GetByTransactionId(paymentId.String())
	if err == nil {
		log.Infof("Receipt for payment %s already exists. Skipping", msg.Id)

		return nil
	}

	report := &db.Report{
		Id:            uuid.NewV4(),
		OwnerId:       ownerId,
		OwnerType:     ukama.OwnerTypeOrg,
		Type:          ukama.ReportTypeReceipt,
		Period:        time.Now().UTC(),
		IsPaid:        true,
		TransactionId: paymentId.String(),
	}

	rawReceipt := rawReceiptFromPayment(msg, report.Id.String())

	rawReceiptBytes, err := json.Marshal(rawReceipt)
	if err != nil {
		return fmt.Errorf("failed to marshal RawReport struct to RawReport JSON: %w", err)
	}

	report.RawReport = datatypes.JSON(rawReceiptBytes)

	log.Infof("Adding receipt for payment: %s", msg.Id)

	err = reportRepo.Add(report, nil)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "report")
	}

	val := &epb.RawReport{}

	m := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	err = m.Unmarshal(rawReceiptBytes, val)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RawReport JSON payload to epb.RawReport: %w", err)
	}

	evt := &epb.Report{
		Id:            report.Id.String(),
		OwnerId:       report.OwnerId.String(),
		OwnerType:     report.OwnerType.String(),
		Type:          report.Type.String(),
		Period:        report.Period.Format(time.RFC3339),
		RawReport:     val,
		IsPaid:        report.IsPaid,
		TransactionId: report.TransactionId,
		CreatedAt:     report.CreatedAt.Format(time.RFC3339),
	}

	route := baseRoutingKey.SetAction("generate").SetObject(report.Type.String()).MustBuild()

	err = msgBus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			evt, route, err.Error())
	}

	return nil
}

func rawReceiptFromPayment(msg *epb.Payment, reportId string) *util.RawReport {
	amountCents := int(msg.AmountCents)

	issuingDate := msg.PaidAt
	if issuingDate == "" {
		issuingDate = time.Now().UTC().Format(time.RFC3339)
	}

	return &util.RawReport{
		Number:                            receiptNumberPrefix + strings.ToUpper(reportId[:8]),
		IssuingDate:                       issuingDate,
		InvoiceType:                       ukama.ReportTypeReceipt.String(),
		Status:                            receiptStatus,
		PaymentStatus:                     receiptPaymentStatus,
		Currency:                          msg.Currency,
		FeesAmountCents:                   amountCents,
		SubTotalExcludingTaxesAmountCents: amountCents,
		SubTotalIncludingTaxesAmountCents: amountCents,
		TotalAmountCents:                  amountCents,
		FileURL:                           fmt.Sprintf(receiptFileURLFormat, reportId),
		Customer: &util.Customer{
			ExternalID: msg.Id,
			Name:       msg.PayerName,
			Email:      msg.PayerEmail,
			Phone:      msg.PayerPhone,
			Country:    msg.Country,
			Currency:   msg.Currency,
		},
		Fees: []util.Fee{
			{
				ExternalSubscriptionID: msg.TransactionId,
				AmountCents:            amountCents,
				TotalAmountCents:       amountCents,
				TotalAmountCurrency:    msg.Currency,
				Description:            msg.Description,
				Item: util.FeeItem{
					Type: msg.PaymentMethod,
					Code: msg.ItemId,
					Name: msg.ItemType,
				},
			},
		},
	}
}

func unmarshalPayment(msg *anypb.Any) (*epb.Payment, error) {
	p := &epb.Payment{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal payment success message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}
