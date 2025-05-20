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

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)


type RegistryEventServer struct {
	orgName string
	s       *SubcriberServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewRegistryEventServer(orgName string, s *SubcriberServer) *RegistryEventServer {
	return &RegistryEventServer{
		orgName: orgName,
		s:       s,
	}
}

func (es *RegistryEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sims.deletion_completed"):
		msg, err := es.unmarshalEventSubscriberAsrCleanupCompleted(e.Msg)
		if err != nil {
			return nil, err
		}

		err = es.handleEventCompleteSubscriberDeletion(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (es *RegistryEventServer) handleEventCompleteSubscriberDeletion(routingKey string, msg *epb.EventSubscriberDeletionCompleted) error {
	log.Infof("Received ASR cleanup completion message with Routing key %s and Message %+v", routingKey, msg)
	
	subscriberId, err := uuid.FromString(msg.SubscriberId)
	if err != nil {
		log.Errorf("Invalid subscriber ID format: %s, Error: %v", msg.SubscriberId, err)
		return status.Errorf(codes.InvalidArgument, "invalid subscriber ID format: %v", err)
	}

	subscriber, err := es.s.subscriberRepo.Get(subscriberId)
	if err != nil {
		log.Errorf("Error while getting subscriber %s: %v", msg.SubscriberId, err)
		return err
	}

	log.Infof("ASR cleanup completed for subscriber %s, proceeding with final deletion", subscriber.SubscriberId.String())

	if err := es.s.subscriberRepo.Delete(subscriber.SubscriberId); err != nil {
		log.Errorf("Error while deleting subscriber %s: %v", subscriber.SubscriberId.String(), err)
		return err
	}

	log.Infof("Successfully completed subscriber deletion for %s after ASR cleanup", subscriber.SubscriberId.String())

	route := es.s.subscriberRoutingKey.SetAction("delete").SetObject("subscriber").MustBuild()
	log.Infof("Publishing subscriber deletion completed event to %v", route)
	
	completionEvent := &epb.EventSubscriberDeleted{
		SubscriberId: subscriber.SubscriberId.String(),
		Name:         subscriber.Name,
	}

	err = es.s.PublishEventMessage(route, completionEvent)
	if err != nil {
		log.Errorf("Failed to publish deletion completed event: %v", err)
	}

	return nil
}

func (es *RegistryEventServer) unmarshalEventSubscriberAsrCleanupCompleted(msg *anypb.Any) (*epb.EventSubscriberDeletionCompleted, error) {
	p := &epb.EventSubscriberDeletionCompleted{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to unmarshal EventSubscriberAsrCleanupCompleted message: %+v. Error: %s", msg, err.Error())
		return nil, err
	}
	return p, nil
}