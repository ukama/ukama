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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	pm "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type AsrEventServer struct {
	asrRepo  db.AsrRecordRepo
	gutiRepo db.GutiRepo
	s        *AsrRecordServer
	orgName  string
	epb.UnimplementedEventNotificationServiceServer
}

func NewAsrEventServer(asrRepo db.AsrRecordRepo, s *AsrRecordServer, gutiRepo db.GutiRepo, org string) *AsrEventServer {
	return &AsrEventServer{
		asrRepo:  asrRepo,
		gutiRepo: gutiRepo,
		orgName:  org,
		s:        s,
	}
}

func (l *AsrEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create"):
		msg, err := l.unmarshalCDRCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleEventCDRCreate(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sims.cleanup_requested"):
		msg, err := l.unmarshalSimAsrCleanupRequested(e.Msg)
		if err != nil {
			return nil, err
		}
		err = l.handleSimAsrCleanupRequested(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (l *AsrEventServer) unmarshalCDRCreate(msg *anypb.Any) (*epb.CDRReported, error) {
	p := &epb.CDRReported{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddSystemRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}


func (l *AsrEventServer) handleEventCDRCreate(key string, msg *epb.CDRReported) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	err := l.s.UpdateandSyncAsrProfile(msg.GetImsi())
	if err != nil {
		log.Errorf("Failed to update the active subscriber %+s.Error: %+v", msg.Imsi, err)
		return err
	}
	return nil
}

func (l *AsrEventServer) unmarshalSimAsrCleanupRequested(msg *anypb.Any) (*epb.EventSimAsrCleanupRequested, error) {
	p := &epb.EventSimAsrCleanupRequested{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal simAsrCleanupRequested message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}
func (l *AsrEventServer) handleSimAsrCleanupRequested(ctx context.Context, key string, msg *epb.EventSimAsrCleanupRequested) error {
	log.Infof("Processing SIM ASR cleanup request from SIM Manager. SubscriberId: %s, SimDetails: %d", 
		msg.SubscriberId, len(msg.SimDetails))

	var cleanupErrors []string
	var successfulCleanups []string

	for _, simDetail := range msg.SimDetails {
		log.Infof("Processing ASR cleanup for SIM - ID: %s, ICCID: %s, IMSI: %s, Status: %s", 
			simDetail.SimId, simDetail.Iccid, simDetail.Imsi, simDetail.Status)

		if simDetail.Status == ukama.SimStatusTerminated.String() {
			log.Infof("SIM %s (ICCID: %s) already terminated, skipping ASR cleanup", 
				simDetail.SimId, simDetail.Iccid)
			successfulCleanups = append(successfulCleanups, simDetail.Iccid)
			continue
		}

		err := l.deleteAsrRecordByIccid(simDetail.Iccid)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to cleanup ASR for ICCID %s: %v", simDetail.Iccid, err)
			log.Errorf(errorMsg)
			cleanupErrors = append(cleanupErrors, errorMsg)
			continue
		}

		log.Infof("Successfully cleaned up ASR record for ICCID: %s", simDetail.Iccid)
		successfulCleanups = append(successfulCleanups, simDetail.Iccid)
	}

	log.Infof("ASR cleanup summary for subscriber %s: %d successful, %d failed", 
		msg.SubscriberId, len(successfulCleanups), len(cleanupErrors))

	success := len(cleanupErrors) == 0

	err := l.publishAsrCleanupCompletionToSimManager(msg, success)
	if err != nil {
		log.Errorf("Failed to publish ASR cleanup completion: %v", err)
		return err
	}

	log.Infof("Completed ASR cleanup processing and notified for subscriber: %s", msg.SubscriberId)
	return nil
}
func (l *AsrEventServer) publishAsrCleanupCompletionToSimManager(originalEvent *epb.EventSimAsrCleanupRequested, success bool) error {
	completionEvent := &epb.EventSimAsrCleanupCompleted{
		SubscriberId: originalEvent.SubscriberId,
		NetworkId:    originalEvent.NetworkId,
		Success:      success,
	}
	route := l.s.baseRoutingKey.SetAction("cleanup_completed").SetObject("sims").MustBuild()

	log.Infof("Publishing ASR cleanup completion at %s: %+v", route, completionEvent)

	if l.s.msgbus != nil {
		err := l.s.msgbus.PublishRequest(route, completionEvent)
		if err != nil {
			log.Errorf("Failed to publish ASR cleanup completion: %v", err)
			return err
		}
		log.Infof("Successfully published ASR cleanup completion for subscriber: %s", 
			originalEvent.SubscriberId)
	} else {
		log.Warnf("Message bus client not available, cannot publish ASR cleanup completion")
	}

	return nil
}
func (l *AsrEventServer) deleteAsrRecordByIccid(iccid string) error {
	asrRecord, err := l.asrRepo.GetByIccid(iccid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("ASR record not found for ICCID %s, might already be deleted", iccid)
			return nil 
		}
		return fmt.Errorf("error checking ASR record for ICCID %s: %w", iccid, err)
	}

	pcrfData := &pm.SimInfo{
		ID:        asrRecord.ID,
		Imsi:      asrRecord.Imsi,
		Iccid:     asrRecord.Iccid,
		NetworkId: asrRecord.NetworkId,
	}

	err = l.asrRepo.Delete(asrRecord.Imsi, db.DEACTIVATION)
	if err != nil {
		return fmt.Errorf("error deleting ASR record for ICCID %s: %w", iccid, err)
	}

	err = l.s.pc.SyncProfile(pcrfData, asrRecord, msgbus.ACTION_CRUD_DELETE, "activesubscriber", true)
	if err != nil {
		log.Errorf("Error syncing PCRF for deleted ASR record (ICCID: %s): %v", iccid, err)
	}

	return nil
}

