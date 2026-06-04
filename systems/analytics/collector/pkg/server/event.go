/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ukama/ukama/systems/analytics/collector/pkg/db"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/datatypes"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

// CollectorEventServer consumes platform events and writes them to the
// analytics database as snapshots, facts and intervals.
//
// Every handler first records the event into analytics_event_logs using a
// deterministic message id derived from the routing key and the event's own
// identity; duplicate deliveries are acknowledged and skipped. Malformed
// payloads are recorded into analytics_event_errors and acknowledged: the
// collector never NACKs an event it cannot process, since redelivery would
// not fix it.
type CollectorEventServer struct {
	orgName      string
	eventRepo    db.EventRepo
	stateRepo    db.StateRepo
	snapshotRepo db.SnapshotRepo
	factRepo     db.FactRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewCollectorEventServer(orgName string, eventRepo db.EventRepo, stateRepo db.StateRepo,
	snapshotRepo db.SnapshotRepo, factRepo db.FactRepo) *CollectorEventServer {
	return &CollectorEventServer{
		orgName:      orgName,
		eventRepo:    eventRepo,
		stateRepo:    stateRepo,
		snapshotRepo: snapshotRepo,
		factRepo:     factRepo,
	}
}

type eventTransactionRunner interface {
	InTransaction(fn func(db.EventRepo, db.StateRepo, db.SnapshotRepo, db.FactRepo) error) error
}

func (es *CollectorEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, e.Msg)

	runner, ok := es.eventRepo.(eventTransactionRunner)
	if !ok {
		return es.dispatchEvent(ctx, e)
	}

	var resp *epb.EventResponse
	err := runner.InTransaction(func(eventRepo db.EventRepo, stateRepo db.StateRepo,
		snapshotRepo db.SnapshotRepo, factRepo db.FactRepo) error {
		txServer := *es
		txServer.eventRepo = eventRepo
		txServer.stateRepo = stateRepo
		txServer.snapshotRepo = snapshotRepo
		txServer.factRepo = factRepo

		r, err := txServer.dispatchEvent(ctx, e)
		if err != nil {
			return err
		}

		resp = r

		return nil
	})
	if err != nil {
		es.recordProcessingError(e.RoutingKey, e.Msg, err)

		return nil, err
	}

	return resp, nil
}

func (es *CollectorEventServer) dispatchEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	switch e.RoutingKey {

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPaymentSuccess]):
		c := evt.EventToEventConfig[evt.EventPaymentSuccess]
		msg, err := epb.UnmarshalPayment(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handlePayment(e.RoutingKey, msg, "success")

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPaymentFailed]):
		c := evt.EventToEventConfig[evt.EventPaymentFailed]
		msg, err := epb.UnmarshalPayment(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handlePayment(e.RoutingKey, msg, "failed")

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberCreate]):
		c := evt.EventToEventConfig[evt.EventSubscriberCreate]
		msg, err := epb.UnmarshalEventSubscriberAdded(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSubscriberAdded(e.RoutingKey, msg, "create")

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberUpdate]):
		c := evt.EventToEventConfig[evt.EventSubscriberUpdate]
		msg, err := epb.UnmarshalEventSubscriberAdded(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSubscriberAdded(e.RoutingKey, msg, "update")

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSubscriberDelete]):
		c := evt.EventToEventConfig[evt.EventSubscriberDelete]
		msg, err := epb.UnmarshalEventSubscriberDeleted(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSubscriberDeleted(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAllocate]):
		c := evt.EventToEventConfig[evt.EventSimAllocate]
		msg, err := epb.UnmarshalEventSimAllocation(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSimAllocate(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimActivate]):
		c := evt.EventToEventConfig[evt.EventSimActivate]
		msg, err := epb.UnmarshalEventSimActivation(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSimActivate(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimAddPackage]):
		c := evt.EventToEventConfig[evt.EventSimAddPackage]
		msg, err := epb.UnmarshalEventSimAddPackage(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSimAddPackage(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimActivePackage]):
		c := evt.EventToEventConfig[evt.EventSimActivePackage]
		msg, err := epb.UnmarshalEventSimActivePackage(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSimActivePackage(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimRemovePackage]):
		c := evt.EventToEventConfig[evt.EventSimRemovePackage]
		msg, err := epb.UnmarshalEventSimRemovePackage(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSimRemovePackage(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimDelete]):
		c := evt.EventToEventConfig[evt.EventSimDelete]
		msg, err := epb.UnmarshalEventSimTermination(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSimDelete(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSimsUpload]):
		c := evt.EventToEventConfig[evt.EventSimsUpload]
		msg, err := epb.UnmarshalEventSimsUploaded(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSimsUpload(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageCreate]):
		c := evt.EventToEventConfig[evt.EventPackageCreate]
		msg, err := epb.UnmarshalCreatePackageEvent(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handlePackageCreate(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageUpdate]):
		c := evt.EventToEventConfig[evt.EventPackageUpdate]
		msg, err := epb.UnmarshalUpdatePackageEvent(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handlePackageUpdate(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventPackageDelete]):
		c := evt.EventToEventConfig[evt.EventPackageDelete]
		msg, err := epb.UnmarshalDeletePackageEvent(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handlePackageDelete(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNetworkAdd]):
		c := evt.EventToEventConfig[evt.EventNetworkAdd]
		msg, err := epb.UnmarshalEventNetworkCreate(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleNetworkAdd(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSiteCreate]):
		c := evt.EventToEventConfig[evt.EventSiteCreate]
		msg, err := epb.UnmarshalEventAddSite(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSiteCreate(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventSiteUpdate]):
		c := evt.EventToEventConfig[evt.EventSiteUpdate]
		msg, err := epb.UnmarshalEventUpdateSite(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleSiteUpdate(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeCreate]):
		c := evt.EventToEventConfig[evt.EventNodeCreate]
		msg, err := epb.UnmarshalEventRegistryNodeCreate(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleNodeCreate(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeUpdate]):
		c := evt.EventToEventConfig[evt.EventNodeUpdate]
		msg, err := epb.UnmarshalEventRegistryNodeUpdate(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleNodeUpdate(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeAssign]):
		c := evt.EventToEventConfig[evt.EventNodeAssign]
		msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleNodeAssign(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeRelease]):
		c := evt.EventToEventConfig[evt.EventNodeRelease]
		msg, err := epb.UnmarshalEventRegistryNodeRelease(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleNodeRelease(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeOnline]):
		c := evt.EventToEventConfig[evt.EventNodeOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleNodeOnline(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeOffline]):
		c := evt.EventToEventConfig[evt.EventNodeOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleNodeOffline(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventNodeStateTransition]):
		c := evt.EventToEventConfig[evt.EventNodeStateTransition]
		msg, err := epb.UnmarshalNodeStateChangeEvent(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleNodeStateTransition(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventHealthReportStore]):
		c := evt.EventToEventConfig[evt.EventHealthReportStore]
		msg, err := epb.UnmarshalHealthReportEvent(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleHealthReportStore(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventComponentsSync]):
		c := evt.EventToEventConfig[evt.EventComponentsSync]
		msg, err := epb.UnmarshalEventInventoryNodeComponentAdd(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleComponentsSync(e.RoutingKey, msg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventInvoiceGenerate]):
		c := evt.EventToEventConfig[evt.EventInvoiceGenerate]
		msg, err := epb.UnmarshalReport(e.Msg, c.Name)
		if err != nil {
			return es.recordMalformed(e.RoutingKey, err)
		}
		return es.handleInvoiceGenerate(e.RoutingKey, msg)

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)

		return &epb.EventResponse{}, nil
	}
}

// logEvent records the event in analytics_event_logs. It returns true when
// the event is new and must be processed, false on duplicate delivery.
// MsgId is derived from the routing key plus a stable identifier carried by
// the event itself (the platform Event envelope has no message id).
func (es *CollectorEventServer) logEvent(routingKey, stableId string, msg interface{}, occurredAt time.Time) (bool, error) {
	payload, err := json.Marshal(msg)
	if err != nil {
		payload = []byte("{}")
	}

	return es.eventRepo.LogEvent(&db.EventLog{
		RoutingKey: routingKey,
		MsgId:      routingKey + ":" + stableId,
		Payload:    datatypes.JSON(payload),
		OccurredAt: occurredAt,
	})
}

// recordMalformed records an unprocessable payload and ACKs the event:
// redelivering a malformed payload would never succeed.
func (es *CollectorEventServer) recordMalformed(routingKey string, cause error) (*epb.EventResponse, error) {
	log.Errorf("Failed to unmarshal message for %s. Error %+v", routingKey, cause)

	if err := es.eventRepo.RecordError(&db.EventError{
		RoutingKey: routingKey,
		Reason:     cause.Error(),
	}); err != nil {
		log.Errorf("failed to record event error: %v", err)
	}

	return &epb.EventResponse{}, nil
}

// recordProcessingError records a processing failure outside the event
// transaction. The original event processing error is still returned to the
// message bus caller so the delivery can be retried or dead-lettered by the
// platform.
func (es *CollectorEventServer) recordProcessingError(routingKey string, msg interface{}, cause error) {
	log.Errorf("Failed to process message for %s. Error %+v", routingKey, cause)

	payload, err := json.Marshal(msg)
	if err != nil {
		payload = []byte("{}")
	}

	if err := es.eventRepo.RecordError(&db.EventError{
		RoutingKey: routingKey,
		Reason:     cause.Error(),
		Payload:    datatypes.JSON(payload),
	}); err != nil {
		log.Errorf("failed to record event error: %v", err)
	}
}

// recordFailure returns the processing error to the caller. EventNotification
// wraps handlers in one DB transaction, so returning the error rolls back the
// event log, snapshot/fact writes and dirty rollup marks as one unit.
func (es *CollectorEventServer) recordFailure(routingKey string, msg interface{}, cause error) (*epb.EventResponse, error) {
	log.Errorf("Failed to process message for %s. Error %+v", routingKey, cause)

	return nil, cause
}

func (es *CollectorEventServer) handlePayment(routingKey string, msg *epb.Payment, outcome string) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	paidAt := now
	if t, err := time.Parse(time.RFC3339, msg.PaidAt); err == nil {
		paidAt = t
	}

	fresh, err := es.logEvent(routingKey, msg.Id, msg, paidAt)
	if err != nil {
		return nil, err
	}

	if !fresh {
		log.Infof("duplicate delivery of payment event %s; skipping", msg.Id)

		return &epb.EventResponse{}, nil
	}

	pe := &db.PaymentEvent{
		ExternalId:  msg.Id,
		Amount:      float64(msg.AmountCents) / 100.0,
		AmountCents: int64(msg.AmountCents),
		Currency:    msg.Currency,
		Status:      outcome,
		PaidAt:      paidAt,
	}

	/* Payment metadata may carry the target ids. */
	metadata := map[string]string{}
	if merr := json.Unmarshal(msg.Metadata, &metadata); merr == nil {
		if v, ok := metadata["targetId"]; ok {
			if id, perr := uuid.FromString(v); perr == nil {
				pe.CustomerId = id
			}
		}

		if v, ok := metadata["packageId"]; ok {
			if id, perr := uuid.FromString(v); perr == nil {
				pe.PackageId = id
			}
		}

		if v, ok := metadata["networkId"]; ok {
			if id, perr := uuid.FromString(v); perr == nil {
				pe.NetworkId = id
			}
		}

		if v, ok := metadata["siteId"]; ok {
			if id, perr := uuid.FromString(v); perr == nil {
				pe.SiteId = id
			}
		}
	}

	if err := es.factRepo.AddPaymentEvent(pe); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	for _, rollup := range []string{"business_sales_daily", "business_package_daily",
		"business_billing_daily"} {
		if err := es.stateRepo.MarkRollupDirty(rollup); err != nil {
			return es.recordFailure(routingKey, msg, err)
		}
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSubscriberAdded(routingKey string, msg *epb.EventSubscriberAdded, kind string) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.SubscriberId+":"+kind+":"+msg.CreatedAt, msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	customerId, perr := uuid.FromString(msg.SubscriberId)
	if perr != nil {
		return es.recordFailure(routingKey, msg, perr)
	}

	networkId, _ := uuid.FromString(msg.NetworkId)

	snap := &db.CustomerSnapshot{
		CustomerId: customerId,
		NetworkId:  networkId,
		Name:       msg.Name,
		Email:      msg.Email,
		Status:     "active",
		UpdatedAt:  now,
	}

	if t, terr := time.Parse(time.RFC3339, msg.CreatedAt); terr == nil {
		snap.SourceCreatedAt = &t
	}

	if err := es.snapshotRepo.UpsertCustomer(snap); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddCustomerEvent(&db.CustomerEvent{
		CustomerId: customerId,
		Kind:       kind,
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.stateRepo.MarkRollupDirty("customer_state_daily"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSubscriberDeleted(routingKey string, msg *epb.EventSubscriberDeleted) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.SubscriberId+":delete", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	customerId, perr := uuid.FromString(msg.SubscriberId)
	if perr != nil {
		return es.recordFailure(routingKey, msg, perr)
	}

	if err := es.snapshotRepo.DeleteCustomer(msg.SubscriberId); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddCustomerEvent(&db.CustomerEvent{
		CustomerId: customerId,
		Kind:       "delete",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.CloseCustomerPackageInterval(customerId, now); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.stateRepo.MarkRollupDirty("customer_state_daily"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSimAllocate(routingKey string, msg *epb.EventSimAllocation) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id+":allocate", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	customerId, _ := uuid.FromString(msg.SubscriberId)

	if err := es.snapshotRepo.UpsertSim(&db.SimSnapshot{
		SimId:       msg.Id,
		Iccid:       msg.Iccid,
		Status:      "assigned",
		CustomerId:  customerId,
		AllocatedAt: &now,
		UpdatedAt:   now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddSimEvent(&db.SimEvent{
		SimId:      msg.Id,
		Kind:       "allocate",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.TransitionSimState(msg.Id, "assigned", now); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSimActivate(routingKey string, msg *epb.EventSimActivation) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id+":activate", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	customerId, _ := uuid.FromString(msg.SubscriberId)

	if err := es.snapshotRepo.UpsertSim(&db.SimSnapshot{
		SimId:      msg.Id,
		Iccid:      msg.Iccid,
		Status:     "active",
		CustomerId: customerId,
		UpdatedAt:  now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddSimEvent(&db.SimEvent{
		SimId:      msg.Id,
		Kind:       "activate",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.TransitionSimState(msg.Id, "active", now); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSimAddPackage(routingKey string, msg *epb.EventSimAddPackage) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id+":addpackage:"+msg.PackageId, msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	if err := es.factRepo.AddSimEvent(&db.SimEvent{
		SimId:      msg.Id,
		Kind:       "add_package",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSimActivePackage(routingKey string, msg *epb.EventSimActivePackage) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id+":activepackage:"+msg.PackageId, msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	if err := es.factRepo.AddSimEvent(&db.SimEvent{
		SimId:      msg.Id,
		Kind:       "active_package",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	customerId, perr := uuid.FromString(msg.SubscriberId)
	if perr == nil {
		packageId, _ := uuid.FromString(msg.PackageId)

		startAt := now
		if msg.PackageStartDate != nil {
			startAt = msg.PackageStartDate.AsTime()
		}

		if err := es.factRepo.OpenCustomerPackageInterval(customerId, packageId,
			"active", startAt); err != nil {
			return es.recordFailure(routingKey, msg, err)
		}
	}

	if err := es.stateRepo.MarkRollupDirty("customer_state_daily"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSimRemovePackage(routingKey string, msg *epb.EventSimRemovePackage) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id+":removepackage:"+msg.PackageId, msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	if err := es.factRepo.AddSimEvent(&db.SimEvent{
		SimId:      msg.Id,
		Kind:       "remove_package",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	customerId, perr := uuid.FromString(msg.SubscriberId)
	if perr == nil {
		if err := es.factRepo.CloseCustomerPackageInterval(customerId, now); err != nil {
			return es.recordFailure(routingKey, msg, err)
		}
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSimDelete(routingKey string, msg *epb.EventSimTermination) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id+":delete", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	if err := es.snapshotRepo.DeleteSim(msg.Id); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddSimEvent(&db.SimEvent{
		SimId:      msg.Id,
		Kind:       "delete",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.TransitionSimState(msg.Id, "deleted", now); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSimsUpload(routingKey string, msg *epb.EventSimsUploaded) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.SimType+":"+now.Format(time.RFC3339), msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	/* The upload event only carries the sim type: use it as the batch id and
	let the next subscriber refresh fill quantities. */
	if err := es.snapshotRepo.UpsertSimBatch(&db.SimBatchSnapshot{
		BatchId:    msg.SimType,
		UploadedAt: &now,
		UpdatedAt:  now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddSimEvent(&db.SimEvent{
		Kind:       "upload",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.stateRepo.MarkRollupDirty("business_inventory_daily"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handlePackageCreate(routingKey string, msg *epb.CreatePackageEvent) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Uuid+":create", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	packageId, perr := uuid.FromString(msg.Uuid)
	if perr != nil {
		return es.recordFailure(routingKey, msg, perr)
	}

	if err := es.snapshotRepo.UpsertPackage(&db.PackageSnapshot{
		PackageId:   packageId,
		Price:       msg.Amount,
		DataQuotaMb: float64(msg.DataVolume),
		Status:      "active",
		UpdatedAt:   now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddPackageEvent(&db.PackageEvent{
		PackageId:  packageId,
		Kind:       "create",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.stateRepo.MarkRollupDirty("business_package_daily"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handlePackageUpdate(routingKey string, msg *epb.UpdatePackageEvent) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Uuid+":update:"+now.Format(time.RFC3339), msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	packageId, perr := uuid.FromString(msg.Uuid)
	if perr != nil {
		return es.recordFailure(routingKey, msg, perr)
	}

	if err := es.factRepo.AddPackageEvent(&db.PackageEvent{
		PackageId:  packageId,
		Kind:       "update",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.stateRepo.MarkRollupDirty("business_package_daily"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handlePackageDelete(routingKey string, msg *epb.DeletePackageEvent) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Uuid+":delete", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	packageId, perr := uuid.FromString(msg.Uuid)
	if perr != nil {
		return es.recordFailure(routingKey, msg, perr)
	}

	if err := es.snapshotRepo.DeletePackage(msg.Uuid); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddPackageEvent(&db.PackageEvent{
		PackageId:  packageId,
		Kind:       "delete",
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.stateRepo.MarkRollupDirty("business_package_daily"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleNetworkAdd(routingKey string, msg *epb.EventNetworkCreate) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id+":add", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	networkId, perr := uuid.FromString(msg.Id)
	if perr != nil {
		return es.recordFailure(routingKey, msg, perr)
	}

	networkStatus := "active"
	if msg.IsDeactivated {
		networkStatus = "inactive"
	}

	if err := es.snapshotRepo.UpsertNetwork(&db.NetworkSnapshot{
		NetworkId: networkId,
		Name:      msg.Name,
		Status:    networkStatus,
		UpdatedAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSiteCreate(routingKey string, msg *epb.EventAddSite) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.SiteId+":create", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	siteId, perr := uuid.FromString(msg.SiteId)
	if perr != nil {
		return es.recordFailure(routingKey, msg, perr)
	}

	networkId, _ := uuid.FromString(msg.NetworkId)

	siteStatus := "online"
	if msg.IsDeactivated {
		siteStatus = "offline"
	}

	if err := es.snapshotRepo.UpsertSite(&db.SiteSnapshot{
		SiteId:    siteId,
		NetworkId: networkId,
		Name:      msg.Name,
		Status:    siteStatus,
		Latitude:  parseCoordinate(msg.Latitude),
		Longitude: parseCoordinate(msg.Longitude),
		UpdatedAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddSiteStateEvent(&db.SiteStateEvent{
		SiteId:     siteId,
		State:      siteStatus,
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.TransitionSiteState(siteId, siteStatus, now); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleSiteUpdate(routingKey string, msg *epb.EventUpdateSite) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.SiteId+":update:"+now.Format(time.RFC3339), msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	siteId, perr := uuid.FromString(msg.SiteId)
	if perr != nil {
		return es.recordFailure(routingKey, msg, perr)
	}

	networkId, _ := uuid.FromString(msg.NetworkId)

	siteStatus := "online"
	if msg.IsDeactivated {
		siteStatus = "offline"
	}

	if err := es.snapshotRepo.UpsertSite(&db.SiteSnapshot{
		SiteId:    siteId,
		NetworkId: networkId,
		Name:      msg.Name,
		Status:    siteStatus,
		Latitude:  parseCoordinate(msg.Latitude),
		Longitude: parseCoordinate(msg.Longitude),
		UpdatedAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddSiteStateEvent(&db.SiteStateEvent{
		SiteId:     siteId,
		State:      siteStatus,
		OccurredAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleNodeCreate(routingKey string, msg *epb.EventRegistryNodeCreate) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.NodeId+":create", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	if err := es.snapshotRepo.UpsertNode(&db.NodeSnapshot{
		NodeId:    msg.NodeId,
		Name:      msg.Name,
		Type:      msg.Type,
		Status:    "configuring",
		UpdatedAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleNodeUpdate(routingKey string, msg *epb.EventRegistryNodeUpdate) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.NodeId+":update:"+now.Format(time.RFC3339), msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	if err := es.snapshotRepo.UpsertNode(&db.NodeSnapshot{
		NodeId:    msg.NodeId,
		Name:      msg.Name,
		UpdatedAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleNodeAssign(routingKey string, msg *epb.EventRegistryNodeAssign) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.NodeId+":assign:"+msg.Site, msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	siteId, _ := uuid.FromString(msg.Site)
	networkId, _ := uuid.FromString(msg.Network)

	if err := es.snapshotRepo.UpsertNode(&db.NodeSnapshot{
		NodeId:    msg.NodeId,
		SiteId:    siteId,
		NetworkId: networkId,
		Type:      msg.Type,
		UpdatedAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleNodeRelease(routingKey string, msg *epb.EventRegistryNodeRelease) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.NodeId+":release:"+msg.Site, msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	if err := es.snapshotRepo.UpsertNode(&db.NodeSnapshot{
		NodeId:    msg.NodeId,
		Type:      msg.Type,
		UpdatedAt: now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleNodeOnline(routingKey string, msg *epb.NodeOnlineEvent) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.NodeId+":online:"+now.Format(time.RFC3339), msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	return es.recordNodeState(routingKey, msg, msg.NodeId, "online", now)
}

func (es *CollectorEventServer) handleNodeOffline(routingKey string, msg *epb.NodeOfflineEvent) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.NodeId+":offline:"+now.Format(time.RFC3339), msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	return es.recordNodeState(routingKey, msg, msg.NodeId, "offline", now)
}

func (es *CollectorEventServer) recordNodeState(routingKey string, msg interface{}, nodeId, state string, at time.Time) (*epb.EventResponse, error) {
	if err := es.snapshotRepo.UpdateNodeStatus(nodeId, state, at); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddNodeStateEvent(&db.NodeStateEvent{
		NodeId:     nodeId,
		State:      state,
		OccurredAt: at,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.TransitionNodeState(nodeId, state, at); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.stateRepo.MarkRollupDirty("node_health_hourly"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleNodeStateTransition(routingKey string, msg *epb.NodeStateChangeEvent) (*epb.EventResponse, error) {
	at := time.Now().UTC()
	if msg.Timestamp != nil {
		at = msg.Timestamp.AsTime()
	}

	fresh, err := es.logEvent(routingKey,
		msg.NodeId+":"+msg.State+":"+msg.Substate+":"+at.Format(time.RFC3339Nano), msg, at)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	if err := es.factRepo.AddNodeStateEvent(&db.NodeStateEvent{
		NodeId:     msg.NodeId,
		State:      msg.State,
		OccurredAt: at,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.TransitionNodeState(msg.NodeId, msg.State, at); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleHealthReportStore(routingKey string, msg *epb.HealthReportEvent) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id, msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	reportedAt := now
	if t, terr := time.Parse(time.RFC3339, msg.ReportedAt); terr == nil {
		reportedAt = t
	}

	payload := msg.Payload
	if len(payload) == 0 {
		payload = []byte("{}")
	}

	if err := es.snapshotRepo.UpsertHealthReport(&db.HealthReportSnapshot{
		NodeId:     msg.NodeId,
		ReportedAt: reportedAt,
		Payload:    datatypes.JSON(payload),
		UpdatedAt:  now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleComponentsSync(routingKey string, msg *epb.EventInventoryNodeComponentAdd) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id+":sync", msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	if err := es.snapshotRepo.UpsertInventory(&db.InventorySnapshot{
		ComponentId: msg.Id,
		Type:        msg.Type,
		State:       "available",
		UpdatedAt:   now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.factRepo.AddInventoryEvent(&db.InventoryEvent{
		ComponentId: msg.Id,
		Kind:        "sync",
		OccurredAt:  now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.stateRepo.MarkRollupDirty("business_inventory_daily"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func (es *CollectorEventServer) handleInvoiceGenerate(routingKey string, msg *epb.Report) (*epb.EventResponse, error) {
	now := time.Now().UTC()

	fresh, err := es.logEvent(routingKey, msg.Id, msg, now)
	if err != nil {
		return nil, err
	}

	if !fresh {
		return &epb.EventResponse{}, nil
	}

	invoicedAt := now
	if t, terr := time.Parse(time.RFC3339, msg.CreatedAt); terr == nil {
		invoicedAt = t
	}

	if err := es.snapshotRepo.UpsertBilling(&db.BillingSnapshot{
		Id:            1,
		LastInvoiceAt: &invoicedAt,
		UpdatedAt:     now,
	}); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	if err := es.stateRepo.MarkRollupDirty("business_billing_daily"); err != nil {
		return es.recordFailure(routingKey, msg, err)
	}

	return &epb.EventResponse{}, nil
}

func parseCoordinate(s string) float64 {
	var v float64

	if s == "" {
		return 0
	}

	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return 0
	}

	return v
}
