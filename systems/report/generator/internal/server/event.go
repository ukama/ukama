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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/report/generator/internal"
	"github.com/ukama/ukama/systems/report/generator/internal/pdf"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	invoiceItemType = "invoice"
	defaultTemplate = "templates/invoice.html.tmpl"
	pdfFolder       = "/srv/static/"
)

type GeneratorEventServer struct {
	orgName        string
	pdfEngine      pdf.PdfEngine
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	epb.UnimplementedEventNotificationServiceServer
}

func NewGeneratorEventServer(orgName string, pdfEngine pdf.PdfEngine, msgBus mb.MsgBusServiceClient) *GeneratorEventServer {
	return &GeneratorEventServer{
		orgName:   orgName,
		pdfEngine: pdfEngine,
		msgbus:    msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(internal.SystemName).SetOrgName(orgName).SetService(internal.ServiceName),
	}
}

func (g *GeneratorEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(g.orgName, "event.cloud.local.{{ .Org}}.billing.invoice.invoice.generate"):
		msg, err := unmarshalInvoiceGenerateEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = g.handleInvoiceGenerateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (g *GeneratorEventServer) handleInvoiceGenerateEvent(key string, msg *epb.Invoice) error {
	err := g.GeneratePDF(msg, defaultTemplate, filepath.Join(pdfFolder, msg.Id+".pdf"))
	if err != nil {
		log.Errorf("Failed to generate invoice PDF: %v", err)
	}

	return err
}

func unmarshalInvoiceGenerateEvent(msg *anypb.Any) (*epb.Invoice, error) {
	p := &epb.Invoice{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal invoice generated message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func (g *GeneratorEventServer) GeneratePDF(data any, templatePath, outputPath string) error {
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

	return nil
}
