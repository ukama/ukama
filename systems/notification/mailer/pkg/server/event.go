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
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/emailTemplate"
	"github.com/ukama/ukama/systems/common/msgbus"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const dateLayout = "January 2, 2006"

type MailerEventServer struct {
	orgName string
	s       *MailerServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewMailerEventServer(orgName string, s *MailerServer) *MailerEventServer {
	return &MailerEventServer{
		orgName: orgName,
		s:       s,
	}
}

func (es *MailerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, e.Msg)

	switch e.RoutingKey {

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInviteCreate]):
		c := evt.EventToEventConfig[evt.EventInviteCreate]
		msg, err := epb.UnmarshalEventInvitationCreated(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		return es.handleEventInviteCreate(ctx, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAllocate]):
		c := evt.EventToEventConfig[evt.EventSimAllocate]
		msg, err := epb.UnmarshalEventSimAllocation(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		return es.handleEventSimAllocate(ctx, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAddPackage]):
		c := evt.EventToEventConfig[evt.EventSimAddPackage]
		msg, err := epb.UnmarshalEventSimAddPackage(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		return es.handleEventSimAddPackage(ctx, msg)

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)

		return nil, fmt.Errorf("no handler for routing key %s", e.RoutingKey)
	}
}

func (es *MailerEventServer) handleEventInviteCreate(ctx context.Context, msg *epb.EventInvitationCreated) (*epb.EventResponse, error) {
	if msg.Email == "" {
		log.Warnf("Skipping %s email for invitation %s: no recipient in event",
			emailTemplate.EmailTemplateMemberInvite, msg.Id)

		return &epb.EventResponse{}, nil
	}

	return es.queue(ctx, msg.Email, emailTemplate.EmailTemplateMemberInvite, map[string]string{
		emailTemplate.EmailKeyInvitation: msg.Id,
		emailTemplate.EmailKeyLink:       msg.Link,
		emailTemplate.EmailKeyName:       msg.Name,
		emailTemplate.EmailKeyOwner:      msg.OwnerName,
		emailTemplate.EmailKeyOrg:        msg.OrgName,
		emailTemplate.EmailKeyRole:       msg.Role.String(),
		emailTemplate.EmailKeyExpiration: formatExpiry(msg.ExpiresAt),
	})
}

func formatExpiry(expiresAt string) string {
	if i := strings.Index(expiresAt, " m="); i != -1 {
		expiresAt = expiresAt[:i]
	}

	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", expiresAt)
	if err != nil {
		return expiresAt
	}

	return t.Format("January 2, 2006 at 3:04 PM MST")
}

func (es *MailerEventServer) handleEventSimAllocate(ctx context.Context, msg *epb.EventSimAllocation) (*epb.EventResponse, error) {
	if msg.QrCode == "" || msg.IsPhysical {
		log.Infof("Skipping %s email for sim %s: physical sim or no qr code",
			emailTemplate.EmailTemplateSimAllocation, msg.Id)

		return &epb.EventResponse{}, nil
	}

	if msg.SubscriberEmail == "" {
		log.Warnf("Skipping %s email for sim %s: no recipient in event",
			emailTemplate.EmailTemplateSimAllocation, msg.Id)

		return &epb.EventResponse{}, nil
	}

	endDate := ""
	if msg.PackageEndDate != nil {
		endDate = msg.PackageEndDate.AsTime().Format(dateLayout)
	}

	return es.queue(ctx, msg.SubscriberEmail, emailTemplate.EmailTemplateSimAllocation, map[string]string{
		emailTemplate.EmailKeySubscriber: msg.SubscriberName,
		emailTemplate.EmailKeyNetwork:    msg.NetworkName,
		emailTemplate.EmailKeyName:       msg.OwnerName,
		emailTemplate.EmailKeyQRCode:     msg.QrCode,
		emailTemplate.EmailKeyVolume:     msg.PackageDataVolume,
		emailTemplate.EmailKeyUnit:       msg.PackageDataUnit,
		emailTemplate.EmailKeyOrg:        msg.OrgName,
		emailTemplate.EmailKeyEndDate:    endDate,
		emailTemplate.EmailKeyPackage:    msg.PackageName,
		emailTemplate.EmailKeyDuration:   msg.PackageDuration,
		emailTemplate.EmailKeyAmount:     msg.PackageAmount,
	})
}

func (es *MailerEventServer) handleEventSimAddPackage(ctx context.Context, msg *epb.EventSimAddPackage) (*epb.EventResponse, error) {
	if msg.SubscriberEmail == "" {
		log.Warnf("Skipping %s email for sim %s: no recipient in event",
			emailTemplate.EmailTemplatePackageAddition, msg.Id)

		return &epb.EventResponse{}, nil
	}

	return es.queue(ctx, msg.SubscriberEmail, emailTemplate.EmailTemplatePackageAddition, map[string]string{
		emailTemplate.EmailKeySubscriber:      msg.SubscriberName,
		emailTemplate.EmailKeyNetwork:         msg.NetworkName,
		emailTemplate.EmailKeyName:            msg.OwnerName,
		emailTemplate.EmailKeyOrg:             msg.OrgName,
		emailTemplate.EmailKeyPackagesCount:   msg.PackagesCount,
		emailTemplate.EmailKeyPackagesDetails: msg.PackagesDetails,
		emailTemplate.EmailKeyExpiration:      msg.PackageEndDate,
		emailTemplate.EmailKeyPackage:         msg.PackageName,
	})
}

func (es *MailerEventServer) queue(ctx context.Context, to string, templateName string, values map[string]string) (*epb.EventResponse, error) {
	mailId, err := es.s.QueueEmail(ctx, []string{to}, templateName, values, nil)
	if err != nil {
		log.Errorf("Failed to queue %s email. Error %v", templateName, err)

		return nil, err
	}

	log.Infof("Queued %s email with mail id %s", templateName, mailId)

	return &epb.EventResponse{}, nil
}
