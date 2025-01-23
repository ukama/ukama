// go:build integration
//go:build integration
// +build integration

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/uuid"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/billing/report/pb/gen"
)

const raw = `{
	"lago_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"sequential_id": 2,
	"number": "LAG-1234-001-002",
	"issuing_date": "2022-04-30",
	"status": "finalized",
	"payment_status": "succeeded",
	"amount_cents": 100,
	"amount_currency": "EUR",
	"vat_amount_cents": 20,
	"vat_amount_currency": "EUR",
	"credit_amount_cents": 10,
	"credit_amount_currency": "EUR",
	"total_amount_cents": 110,
	"total_amount_currency": "EUR",
	"file_url": "https://getlago.com/invoice/file",
	"legacy": false,
	"customer": {
	"lago_id": "99a6094e-199b-4101-896a-54e927ce7bd7",
	"external_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"address_line1": "5230 Penfield Ave",
	"address_line2": null,
	"city": "Woodland Hills",
	"country": "US",
	"created_at": "2022-04-29T08:59:51Z",
	"email": "dinesh@piedpiper.test",
	"legal_name": "Coleman-Blair",
	"legal_number": "49-008-2965",
	"logo_url": "http://hooli.com/logo.png",
	"name": "Gavin Belson",
	"phone": "1-171-883-3711 x245",
	"state": "CA",
	"url": "http://hooli.com",
	"vat_rate": 20.0,
	"zipcode": "91364"
	},
	"subscriptions": [
	{
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"external_id": "susbcription_external_id",
	"lago_customer_id": "99a6094e-199b-4101-896a-54e927ce7bd7",
	"external_customer_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"canceled_at": "2022-04-29T08:59:51Z",
	"created_at": "2022-04-29T08:59:51Z",
	"plan_code": "new_code",
	"started_at": "2022-04-29T08:59:51Z",
	"status": "active",
	"terminated_at": null
	}
	],
	"fees": [
	{
	"lago_id": "6be23c42-47d2-45a3-9770-5b3572f225c3",
	"lago_group_id": null,
	"item": {
	"type": "subscription",
	"code": "plan_code",
	"name": "Plan"
	},
	"amount_cents": 100,
	"amount_currency": "EUR",
	"vat_amount_cents": 20,
	"vat_amount_currency": "EUR",
	"total_amount_cents": 120,
	"total_amount_currency": "EUR",
	"units": "0.32",
	"events_count": 23
	}
	],
	"credits": [
	{
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"item": {
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"type": "coupon",
	"code": "coupon_code",
	"name": "Coupon"
	},
	"amount_cents": 100,
	"amount_currency": "EUR"
	}
	],
	"metadata": [
	{
	"lago_id": "27f12d13-4ae0-437b-b822-8771bcd62e3a",
	"key": "digital_ref_id",
	"value": "INV-0123456-98765",
	"created_at": "2022-04-29T08:59:51Z"
	}
	]
	}`

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

func init() {
	// load config
	tConfig = &TestConfig{}

	reader := rconf.NewConfReader("integration")

	err := reader.Read(tConfig)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("Config: %+v\n", tConfig)
}

func Test_FullFlow(t *testing.T) {
	// we need real subscriberId from subscriber-registry or a real orgId
	ownerId := uuid.NewV4().String()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateReportServiceClient()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}
	defer conn.Close()

	repResp := &pb.ReportResponse{}

	t.Run("AddReport", func(t *testing.T) {
		var err error

		repResp, err = c.Add(ctx, &pb.AddRequest{
			RawReport: raw,
		})

		assert.NoError(t, err)
	})

	t.Run("GetReport", func(t *testing.T) {
		_, err := c.Get(ctx, &pb.GetRequest{
			ReportId: repResp.Report.Id,
		})

		assert.NoError(t, err)
	})

	t.Run("GetReportByOwner", func(t *testing.T) {
		_, err := c.List(ctx, &pb.ListRequest{
			OwnerId: ownerId,
		})

		assert.NoError(t, err)
	})

	t.Run("DeleteReport", func(t *testing.T) {
		_, err := c.Delete(ctx, &pb.DeleteRequest{
			ReportId: repResp.Report.Id,
		})

		assert.NoError(t, err)
	})

	t.Run("GetReport", func(t *testing.T) {
		_, err := c.Get(ctx, &pb.GetRequest{
			ReportId: repResp.Report.Id,
		})

		assert.Error(t, err)
	})
}

func CreateReportServiceClient() (*grpc.ClientConn, pb.ReportServiceClient, error) {
	log.Infoln("Connecting to Report Server ", tConfig.ServiceHost)

	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewReportServiceClient(conn)
	return conn, c, nil
}
