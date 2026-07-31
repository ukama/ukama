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
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/util/payments"
	"github.com/ukama/ukama/systems/report/generator/internal"
	"github.com/ukama/ukama/systems/report/generator/internal/pdf"
	"github.com/ukama/ukama/systems/report/generator/internal/storage"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	defaultTemplate   = "templates/invoice.html.tmpl"
	receiptTemplate   = "templates/receipt.html.tmpl"
	pdfFolder         = "/home/ukama/srv/static/"
	receiptDateFormat = "January 2, 2006"
	timeStringLayout  = "2006-01-02 15:04:05.999999999 -0700 MST"
)

type GeneratorEventServer struct {
	orgName        string
	pdfEngine      pdf.PdfEngine
	storage        storage.Storage
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	epb.UnimplementedEventNotificationServiceServer
}

func NewGeneratorEventServer(orgName string, pdfEngine pdf.PdfEngine, store storage.Storage, msgBus mb.MsgBusServiceClient) *GeneratorEventServer {
	return &GeneratorEventServer{
		orgName:   orgName,
		pdfEngine: pdfEngine,
		storage:   store,
		msgbus:    msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(internal.SystemName).SetOrgName(orgName).SetService(internal.ServiceName),
	}
}

func (g *GeneratorEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(g.orgName, "event.cloud.local.{{ .Org}}.billing.report.invoice.generate"),
		msgbus.PrepareRoute(g.orgName, "event.cloud.local.{{ .Org}}.billing.report.invoice.update"):
		msg, err := unmarshalInvoiceGenerateEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = g.handleInvoiceGenerateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(g.orgName, "event.cloud.local.{{ .Org}}.payments.processor.payment.success"):
		msg, err := unmarshalPaymentEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = g.handlePaymentSuccessEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (g *GeneratorEventServer) handleInvoiceGenerateEvent(key string, msg *epb.Report) error {
	err := g.GeneratePDF(msg, defaultTemplate, filepath.Join(pdfFolder, msg.Id+".pdf"))
	if err != nil {
		log.Errorf("Failed to generate invoice PDF: %v", err)
	}

	return err
}

func (g *GeneratorEventServer) handlePaymentSuccessEvent(key string, msg *epb.Payment) error {
	if ukama.ParseItemType(msg.ItemType) != ukama.ItemTypePackage {
		log.Infof("Skipping receipt for payment %s: item type %q is not a one-off package", msg.Id, msg.ItemType)

		return nil
	}

	report := buildReceiptReport(msg)
	objectName := msg.Id + ".pdf"

	err := g.GeneratePDF(report, receiptTemplate, filepath.Join(pdfFolder, objectName))
	if err != nil {
		log.Errorf("Failed to generate receipt PDF: %v", err)

		return err
	}

	return g.publishReceiptGenerated(msg, report, objectName)
}

func (g *GeneratorEventServer) publishReceiptGenerated(msg *epb.Payment, report *epb.Report, objectName string) error {
	if g.storage == nil {
		log.Warnf("Skipping receipt generated event for payment %s: storage is not configured", msg.Id)

		return nil
	}

	evt := &epb.EventReceiptGenerated{
		Id:            msg.Id,
		ReceiptNumber: report.RawReport.Number,
		PayerName:     msg.PayerName,
		PayerEmail:    msg.PayerEmail,
		Amount:        payments.ToAmount(msg.AmountCents),
		Currency:      msg.Currency,
		PaidAt:        report.RawReport.IssuingDate,
		PaymentMethod: msg.PaymentMethod,
		Description:   msg.Description,
		OrgName:       g.orgName,
		Bucket:        g.storage.Bucket(),
		ObjectKey:     objectName,
		FileName:      report.RawReport.Number + ".pdf",
	}

	route := g.baseRoutingKey.SetObject("receipt").SetAction("generate").MustBuild()

	err := g.msgbus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish receipt generated event for payment %s on route %s: %v", msg.Id, route, err)

		return fmt.Errorf("failed to publish receipt generated event: %w", err)
	}

	log.Infof("Published receipt generated event for payment %s on route %s", msg.Id, route)

	return nil
}

func buildReceiptReport(p *epb.Payment) *epb.Report {
	number := "RCPT-" + strings.ToUpper(shortID(p.Id))

	return &epb.Report{
		Id:            p.Id,
		Type:          "receipt",
		IsPaid:        true,
		TransactionId: p.Id,
		CreatedAt:     p.PaidAt,
		RawReport: &epb.RawReport{
			Number:                            number,
			IssuingDate:                       formatReceiptDate(p.PaidAt),
			InvoiceType:                       "receipt",
			Status:                            "finalized",
			PaymentStatus:                     "succeeded",
			Currency:                          p.Currency,
			FeesAmountCents:                   p.AmountCents,
			SubTotalIncludingTaxesAmountCents: p.AmountCents,
			TotalAmountCents:                  p.AmountCents,
			FileURL:                           fmt.Sprintf("http://{API_ENDPOINT}/pdf/%s.pdf", p.Id),
			Customer: &epb.Customer{
				ExternalId: p.Id,
				Name:       p.PayerName,
				Email:      p.PayerEmail,
				Phone:      p.PayerPhone,
				Country:    p.Country,
			},
			Fees: []*epb.Fee{
				{
					ExternalSubscriptionId: p.TransactionId,
					AmountCents:            p.AmountCents,
					AmountCurrency:         p.Currency,
					TotalAmountCents:       p.AmountCents,
					TotalAmountCurrency:    p.Currency,
					Description:            p.Description,
					Item: &epb.FeeItem{
						Type: p.PaymentMethod,
						Code: p.ItemId,
						Name: p.ItemType,
					},
				},
			},
		},
	}
}

func shortID(id string) string {
	if len(id) >= 8 {
		return id[:8]
	}

	return id
}

func formatReceiptDate(paidAt string) string {
	if paidAt == "" {
		return time.Now().UTC().Format(receiptDateFormat)
	}

	if i := strings.Index(paidAt, " m="); i != -1 {
		paidAt = paidAt[:i]
	}

	t, err := time.Parse(time.RFC3339, paidAt)
	if err != nil {
		t, err = time.Parse(timeStringLayout, paidAt)
	}

	if err != nil {
		return paidAt
	}

	return t.Format(receiptDateFormat)
}

func unmarshalInvoiceGenerateEvent(msg *anypb.Any) (*epb.Report, error) {
	p := &epb.Report{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal invoice generated message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalPaymentEvent(msg *anypb.Any) (*epb.Payment, error) {
	p := &epb.Payment{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal payment message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func (g *GeneratorEventServer) GeneratePDF(data any, templatePath, outputPath string) error {
	log.Info("Generating new PDF... ")

	pdf := pdf.NewPDFObject("", g.pdfEngine)

	err := pdf.ParseTemplate(templatePath, data)
	if err != nil {
		log.Errorf("failed to parse PDF template: %v", err)

		return fmt.Errorf("failed to parse PDF template: %w", err)

	}

	err = pdf.GenerateFile(outputPath)
	if err != nil {
		log.Errorf("failed to generate PDF file: %v", err)

		return fmt.Errorf("failed to generate PDF file: %w", err)
	}

	log.Info("PDF generated successfully")

	if g.storage != nil {
		objectName := filepath.Base(outputPath)

		location, err := g.storage.Upload(context.Background(), objectName, outputPath)
		if err != nil {
			log.Errorf("failed to upload PDF %s to storage: %v", objectName, err)

			return fmt.Errorf("failed to upload PDF to storage: %w", err)
		}

		log.Infof("Uploaded PDF %s to storage at %s", objectName, location)
	}

	return nil
}
