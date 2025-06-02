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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/db"
	simMangerPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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
	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.delete"):
		msg, err := es.unmarshalAsrInactivated(e.Msg)
		if err != nil {
			return nil, err
		}
		err = es.handleSubscriberDeactivated(ctx, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activate"):
		msg, err := es.unmarshalSimActivation(e.Msg)
		if err != nil {
			return nil, err
		}
		err = es.handleSimActivation(ctx, msg)
		if err != nil {
			return nil, err
		}

	// NEW: SIM Deactivation Event  
	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.deactivate"):
		msg, err := es.unmarshalSimDeactivation(e.Msg)
		if err != nil {
			return nil, err
		}
		err = es.handleSimDeactivation(ctx, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.create"):
		msg, err := es.unmarshalAsrActivated(e.Msg)
		if err != nil {
			return nil, err
		}
		err = es.handleSubscriberActivated(ctx, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
func (es *RegistryEventServer) handleSimDeactivation(ctx context.Context, event *epb.EventSimDeactivation) error {
	log.Infof("SIM deactivated - ID: %s, ICCID: %s, Subscriber: %s", event.Id, event.Iccid, event.SubscriberId)
	
	subscriberId, err := uuid.FromString(event.SubscriberId)
	if err != nil {
		log.Errorf("Invalid subscriber ID %s: %v", event.SubscriberId, err)
		return err
	}
	
	subscriber := db.Subscriber{
		SubscriberStatus: ukama.SubscriberStatusInactive,
	}
	
	err = es.s.subscriberRepo.Update(subscriberId, subscriber)
	if err != nil {
		log.Errorf("Failed to update subscriber %s status to inactive: %v", subscriberId, err)
		return err
	}
	
	log.Infof("Updated subscriber %s status to inactive due to SIM deactivation", subscriberId)
	return nil
}
func (es *RegistryEventServer) handleSimActivation(ctx context.Context, event *epb.EventSimActivation) error {
	log.Infof("SIM activated - ID: %s, ICCID: %s, Subscriber: %s", event.Id, event.Iccid, event.SubscriberId)
	
	subscriberId, err := uuid.FromString(event.SubscriberId)
	if err != nil {
		log.Errorf("Invalid subscriber ID %s: %v", event.SubscriberId, err)
		return err
	}
	
	subscriber := db.Subscriber{
		SubscriberStatus: ukama.SubscriberStatusActive,
	}
	
	err = es.s.subscriberRepo.Update(subscriberId, subscriber)
	if err != nil {
		log.Errorf("Failed to update subscriber %s status to active: %v", subscriberId, err)
		return err
	}
	
	log.Infof("Updated subscriber %s status to active due to SIM activation", subscriberId)
	return nil
}
func (es *RegistryEventServer) handleSubscriberDeactivated(ctx context.Context, event *epb.AsrInactivated) error {
	log.Infof("SIM deactivated - IMSI: %s, ICCID: %s", event.Subscriber.Imsi, event.Subscriber.Iccid)
	return es.updateSubscriberStatusBySim(ctx, event.Subscriber.Iccid, ukama.SubscriberStatusInactive)
}

func (es *RegistryEventServer) handleSubscriberActivated(ctx context.Context, event *epb.AsrActivated) error {
	log.Infof("SIM activated - IMSI: %s, ICCID: %s", event.Subscriber.Imsi, event.Subscriber.Iccid)
	return es.updateSubscriberStatusBySim(ctx, event.Subscriber.Iccid, ukama.SubscriberStatusActive)
}

func (es *RegistryEventServer) updateSubscriberStatusBySim(ctx context.Context, iccid string, status ukama.SubscriberStatus) error {
	simManagerClient, err := es.s.simManagerService.GetSimManagerService()
	if err != nil {
		log.Errorf("Failed to get SIM manager client: %v", err)
		return err
	}
	
	simResp, err := simManagerClient.ListSims(ctx, &simMangerPb.ListSimsRequest{
		Iccid: iccid,
	})
	if err != nil {
		log.Errorf("Failed to get SIM by ICCID %s: %v", iccid, err)
		return err
	}
	
	subscriberId, err := uuid.FromString(simResp.Sims[0].SubscriberId)
	if err != nil {
		log.Errorf("Invalid subscriber ID %s: %v", simResp.Sims[0].SubscriberId, err)
		return err
	}
	
	subscriber := db.Subscriber{
		SubscriberStatus: status,
	}
	
	err = es.s.subscriberRepo.Update(subscriberId, subscriber)
	if err != nil {
		log.Errorf("Failed to update subscriber %s status: %v", subscriberId, err)
		return err
	}
	
	log.Infof("Updated subscriber %s status to %s", subscriberId, status)
	return nil
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
func (es *RegistryEventServer) unmarshalAsrInactivated(msg *anypb.Any) (*epb.AsrInactivated, error) {
	p := &epb.AsrInactivated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to unmarshal AsrInactivated: %v", err)
		return nil, err
	}
	return p, nil
}
func (es *RegistryEventServer) unmarshalAsrActivated(msg *anypb.Any) (*epb.AsrActivated, error) {
	p := &epb.AsrActivated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to unmarshal AsrActivated: %v", err)
		return nil, err
	}
	return p, nil
}
func (es *RegistryEventServer) unmarshalSimActivation(msg *anypb.Any) (*epb.EventSimActivation, error) {
	p := &epb.EventSimActivation{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to unmarshal EventSimActivation message: %+v. Error: %s", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (es *RegistryEventServer) unmarshalSimDeactivation(msg *anypb.Any) (*epb.EventSimDeactivation, error) {
	p := &epb.EventSimDeactivation{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to unmarshal EventSimDeactivation message: %+v. Error: %s", msg, err.Error())
		return nil, err
	}
	return p, nil
}
