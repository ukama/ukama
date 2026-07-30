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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/emailTemplate"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/notification/mailer/mocks"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/storage"

	evt "github.com/ukama/ukama/systems/common/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

const (
	testEventOrgName    = "test-org"
	testSubscriberEmail = "subscriber@example.com"
	testSubscriberName  = "Test Subscriber"
	testInviteEmail     = "invitee@example.com"
	testNetworkName     = "test-network"
	testOwnerName       = "Test Owner"
	testQrCode          = "qr-code-payload"
	testPackageName     = "Monthly 5GB"
	testPayerEmail      = "brackley@ukama.com"
	testPayerName       = "Brackley"
	testReceiptNumber   = "RCPT-37522EA5"
	testReceiptBucket   = "report-ukama"
	testReceiptKey      = "37522ea5-3207-4d81-b19d-ab162adab725.pdf"
)

func setupEventServer(t *testing.T) (*MailerEventServer, *mocks.MailerRepo) {
	return setupEventServerWithStorage(t, nil)
}

func setupEventServerWithStorage(t *testing.T, store storage.Storage) (*MailerEventServer, *mocks.MailerRepo) {
	server, mockRepo := setupServer(t)

	return NewMailerEventServer(testEventOrgName, server, store), mockRepo
}

func eventFor(t *testing.T, key evt.EventId, msg proto.Message) *epb.Event {
	t.Helper()

	anyMsg, err := anypb.New(msg)
	require.NoError(t, err)

	return &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testEventOrgName, evt.EventRoutingKey[key]),
		Msg:        anyMsg,
	}
}

func expectQueuedEmail(repo *mocks.MailerRepo) *db.Mailing {
	created := &db.Mailing{}

	repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once().
		Run(func(args mock.Arguments) {
			*created = *args.Get(0).(*db.Mailing)
		})

	repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
		Status: ukama.MailStatusPending,
	}, nil).Maybe()
	repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()

	return created
}

func TestEventNotification_InviteCreate(t *testing.T) {
	t.Run("queues member invite email", func(t *testing.T) {
		es, repo := setupEventServer(t)
		created := expectQueuedEmail(repo)

		inviteId := "6a1b1a4e-0e1a-4f2e-8a1f-2f6b0b3c9d11"

		res, err := es.EventNotification(context.TODO(), eventFor(t, evt.EventInviteCreate,
			&epb.EventInvitationCreated{
				Id:        inviteId,
				Link:      "https://ukama.com/invite",
				Email:     testInviteEmail,
				Name:      "Invitee Name",
				Role:      upb.RoleType_ROLE_ADMIN,
				OrgName:   testEventOrgName,
				OwnerName: testOwnerName,
				ExpiresAt: "2026-07-25 14:30:00 +0000 UTC m=+86400.001234567",
			}))

		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, testInviteEmail, created.Email)
		assert.Equal(t, emailTemplate.EmailTemplateMemberInvite, created.TemplateName)
		assert.Equal(t, ukama.MailStatusPending, created.Status)
		assert.Equal(t, inviteId, created.Values[emailTemplate.EmailKeyInvitation])
		assert.Equal(t, "Invitee Name", created.Values[emailTemplate.EmailKeyName])
		assert.Equal(t, testOwnerName, created.Values[emailTemplate.EmailKeyOwner])
		assert.Equal(t, testEventOrgName, created.Values[emailTemplate.EmailKeyOrg])
		assert.Equal(t, "Administrator", created.Values[emailTemplate.EmailKeyRole])
		assert.Equal(t, "July 25, 2026 at 2:30 PM UTC", created.Values[emailTemplate.EmailKeyExpiration])
	})

	t.Run("keeps raw expiry when format is unknown", func(t *testing.T) {
		es, repo := setupEventServer(t)
		created := expectQueuedEmail(repo)

		res, err := es.EventNotification(context.TODO(), eventFor(t, evt.EventInviteCreate,
			&epb.EventInvitationCreated{
				Id:        "6a1b1a4e-0e1a-4f2e-8a1f-2f6b0b3c9d11",
				Email:     testInviteEmail,
				ExpiresAt: "tomorrow",
			}))

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "tomorrow", created.Values[emailTemplate.EmailKeyExpiration])
	})

	t.Run("skips when recipient is missing", func(t *testing.T) {
		es, repo := setupEventServer(t)

		res, err := es.EventNotification(context.TODO(), eventFor(t, evt.EventInviteCreate,
			&epb.EventInvitationCreated{
				Id:      "6a1b1a4e-0e1a-4f2e-8a1f-2f6b0b3c9d11",
				OrgName: testEventOrgName,
			}))

		assert.NoError(t, err)
		assert.NotNil(t, res)
		repo.AssertNotCalled(t, "CreateEmail", mock.Anything)
	})
}

func TestEventNotification_SimAllocate(t *testing.T) {
	endDate := time.Date(2026, time.March, 4, 12, 0, 0, 0, time.UTC)

	allocation := func() *epb.EventSimAllocation {
		return &epb.EventSimAllocation{
			Id:                "9f6b1a4e-0e1a-4f2e-8a1f-2f6b0b3c9d22",
			SubscriberName:    testSubscriberName,
			SubscriberEmail:   testSubscriberEmail,
			NetworkName:       testNetworkName,
			OrgName:           testEventOrgName,
			OwnerName:         testOwnerName,
			QrCode:            testQrCode,
			PackageName:       testPackageName,
			PackageDataVolume: "5",
			PackageDataUnit:   "GB",
			PackageAmount:     "10",
			PackageDuration:   "30",
			PackageEndDate:    timestamppb.New(endDate),
		}
	}

	t.Run("queues sim allocation email", func(t *testing.T) {
		es, repo := setupEventServer(t)
		created := expectQueuedEmail(repo)

		res, err := es.EventNotification(context.TODO(),
			eventFor(t, evt.EventSimAllocate, allocation()))

		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, testSubscriberEmail, created.Email)
		assert.Equal(t, emailTemplate.EmailTemplateSimAllocation, created.TemplateName)
		assert.Equal(t, testSubscriberName, created.Values[emailTemplate.EmailKeySubscriber])
		assert.Equal(t, testNetworkName, created.Values[emailTemplate.EmailKeyNetwork])
		assert.Equal(t, testOwnerName, created.Values[emailTemplate.EmailKeyName])
		assert.Equal(t, testQrCode, created.Values[emailTemplate.EmailKeyQRCode])
		assert.Equal(t, testPackageName, created.Values[emailTemplate.EmailKeyPackage])
		assert.Equal(t, "March 4, 2026", created.Values[emailTemplate.EmailKeyEndDate])
	})

	t.Run("skips physical sim", func(t *testing.T) {
		es, repo := setupEventServer(t)

		msg := allocation()
		msg.IsPhysical = true

		res, err := es.EventNotification(context.TODO(), eventFor(t, evt.EventSimAllocate, msg))

		assert.NoError(t, err)
		assert.NotNil(t, res)
		repo.AssertNotCalled(t, "CreateEmail", mock.Anything)
	})

	t.Run("skips when qr code is missing", func(t *testing.T) {
		es, repo := setupEventServer(t)

		msg := allocation()
		msg.QrCode = ""

		res, err := es.EventNotification(context.TODO(), eventFor(t, evt.EventSimAllocate, msg))

		assert.NoError(t, err)
		assert.NotNil(t, res)
		repo.AssertNotCalled(t, "CreateEmail", mock.Anything)
	})

	t.Run("skips when recipient is missing", func(t *testing.T) {
		es, repo := setupEventServer(t)

		msg := allocation()
		msg.SubscriberEmail = ""

		res, err := es.EventNotification(context.TODO(), eventFor(t, evt.EventSimAllocate, msg))

		assert.NoError(t, err)
		assert.NotNil(t, res)
		repo.AssertNotCalled(t, "CreateEmail", mock.Anything)
	})

	t.Run("tolerates missing end date", func(t *testing.T) {
		es, repo := setupEventServer(t)
		created := expectQueuedEmail(repo)

		msg := allocation()
		msg.PackageEndDate = nil

		res, err := es.EventNotification(context.TODO(), eventFor(t, evt.EventSimAllocate, msg))

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "", created.Values[emailTemplate.EmailKeyEndDate])
	})
}

func TestEventNotification_SimAddPackage(t *testing.T) {
	addPackage := func() *epb.EventSimAddPackage {
		return &epb.EventSimAddPackage{
			Id:              "3c6b1a4e-0e1a-4f2e-8a1f-2f6b0b3c9d33",
			SubscriberName:  testSubscriberName,
			SubscriberEmail: testSubscriberEmail,
			NetworkName:     testNetworkName,
			OrgName:         testEventOrgName,
			OwnerName:       testOwnerName,
			PackageName:     testPackageName,
			PackagesCount:   "2",
			PackagesDetails: "$10.00 / 5 GB / 30 days",
			PackageEndDate:  "March 4, 2026",
		}
	}

	t.Run("queues package addition email", func(t *testing.T) {
		es, repo := setupEventServer(t)
		created := expectQueuedEmail(repo)

		res, err := es.EventNotification(context.TODO(),
			eventFor(t, evt.EventSimAddPackage, addPackage()))

		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, testSubscriberEmail, created.Email)
		assert.Equal(t, emailTemplate.EmailTemplatePackageAddition, created.TemplateName)
		assert.Equal(t, testSubscriberName, created.Values[emailTemplate.EmailKeySubscriber])
		assert.Equal(t, testOwnerName, created.Values[emailTemplate.EmailKeyName])
		assert.Equal(t, "2", created.Values[emailTemplate.EmailKeyPackagesCount])
		assert.Equal(t, "$10.00 / 5 GB / 30 days", created.Values[emailTemplate.EmailKeyPackagesDetails])
		assert.Equal(t, "March 4, 2026", created.Values[emailTemplate.EmailKeyExpiration])
	})

	t.Run("skips when recipient is missing", func(t *testing.T) {
		es, repo := setupEventServer(t)

		msg := addPackage()
		msg.SubscriberEmail = ""

		res, err := es.EventNotification(context.TODO(), eventFor(t, evt.EventSimAddPackage, msg))

		assert.NoError(t, err)
		assert.NotNil(t, res)
		repo.AssertNotCalled(t, "CreateEmail", mock.Anything)
	})
}

func TestEventNotification_Errors(t *testing.T) {
	t.Run("unknown routing key", func(t *testing.T) {
		es, repo := setupEventServer(t)

		res, err := es.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: "event.cloud.local.test-org.registry.member.member.create",
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		repo.AssertNotCalled(t, "CreateEmail", mock.Anything)
	})

	for _, tc := range []struct {
		name string
		key  evt.EventId
		msg  proto.Message
	}{
		{"invite create", evt.EventInviteCreate, &epb.EventSimAddPackage{Id: "sim"}},
		{"sim allocate", evt.EventSimAllocate, &epb.EventInvitationCreated{Id: "invite"}},
		{"sim add package", evt.EventSimAddPackage, &epb.EventInvitationCreated{Id: "invite"}},
	} {
		t.Run("payload does not match routing key: "+tc.name, func(t *testing.T) {
			es, repo := setupEventServer(t)

			res, err := es.EventNotification(context.TODO(), eventFor(t, tc.key, tc.msg))

			assert.Error(t, err)
			assert.Nil(t, res)
			repo.AssertNotCalled(t, "CreateEmail", mock.Anything)
		})
	}

	t.Run("propagates queue failure", func(t *testing.T) {
		es, repo := setupEventServer(t)

		repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).
			Return(errTestDatabaseError).Once()

		res, err := es.EventNotification(context.TODO(), eventFor(t, evt.EventInviteCreate,
			&epb.EventInvitationCreated{
				Id:    "6a1b1a4e-0e1a-4f2e-8a1f-2f6b0b3c9d11",
				Email: testInviteEmail,
			}))

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestFriendlyRole(t *testing.T) {
	cases := map[upb.RoleType]string{
		upb.RoleType_ROLE_OWNER:         "Owner",
		upb.RoleType_ROLE_ADMIN:         "Administrator",
		upb.RoleType_ROLE_NETWORK_OWNER: "Network Owner",
		upb.RoleType_ROLE_VENDOR:        "Vendor",
		upb.RoleType_ROLE_USER:          "User",
		upb.RoleType_ROLE_SUBSCRIBER:    "Subscriber",
		upb.RoleType_ROLE_INVALID:       "Member",
	}
	for role, want := range cases {
		assert.Equal(t, want, friendlyRole(role))
	}
}


func receiptEvent() *epb.EventReceiptGenerated {
	return &epb.EventReceiptGenerated{
		Id:            "37522ea5-3207-4d81-b19d-ab162adab725",
		ReceiptNumber: testReceiptNumber,
		PayerName:     testPayerName,
		PayerEmail:    testPayerEmail,
		Amount:        "15.00",
		Currency:      "USD",
		PaidAt:        "July 30, 2026",
		PaymentMethod: "cash",
		Description:   "Prepaid 15GB",
		OrgName:       testEventOrgName,
		Bucket:        testReceiptBucket,
		ObjectKey:     testReceiptKey,
		FileName:      testReceiptNumber + ".pdf",
	}
}

func TestEventNotification_ReceiptGenerate(t *testing.T) {
	t.Run("queues receipt email with the pdf fetched from storage", func(t *testing.T) {
		store := mocks.NewStorage(t)
		store.On("Get", mock.Anything, testReceiptBucket, testReceiptKey).
			Return([]byte("%PDF-1.4 receipt"), nil).Once()

		es, repo := setupEventServerWithStorage(t, store)
		created := expectQueuedEmail(repo)

		res, err := es.EventNotification(context.TODO(),
			eventFor(t, evt.EventReceiptGenerate, receiptEvent()))

		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, testPayerEmail, created.Email)
		assert.Equal(t, emailTemplate.EmailTemplatePaymentReceipt, created.TemplateName)
		assert.Equal(t, ukama.MailStatusPending, created.Status)
		assert.Equal(t, testPayerName, created.Values[emailTemplate.EmailKeyName])
		assert.Equal(t, testEventOrgName, created.Values[emailTemplate.EmailKeyOrg])
		assert.Equal(t, "USD 15.00", created.Values[emailTemplate.EmailKeyAmount])
		assert.Equal(t, testReceiptNumber, created.Values[emailTemplate.EmailKeyReceiptNumber])
		assert.Equal(t, "July 30, 2026", created.Values[emailTemplate.EmailKeyPaymentDate])
		assert.Equal(t, "cash", created.Values[emailTemplate.EmailKeyPaymentMethod])
	})

	t.Run("skips when the event carries no payer email", func(t *testing.T) {
		es, _ := setupEventServerWithStorage(t, nil)

		msg := receiptEvent()
		msg.PayerEmail = ""

		res, err := es.EventNotification(context.TODO(),
			eventFor(t, evt.EventReceiptGenerate, msg))

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("returns error so the event is retried when the pdf cannot be fetched", func(t *testing.T) {
		store := mocks.NewStorage(t)
		store.On("Get", mock.Anything, testReceiptBucket, testReceiptKey).
			Return(nil, assert.AnError).Once()

		es, _ := setupEventServerWithStorage(t, store)

		res, err := es.EventNotification(context.TODO(),
			eventFor(t, evt.EventReceiptGenerate, receiptEvent()))

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
