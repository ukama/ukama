/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package refresh

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ukama/ukama/systems/analytics/collector/pkg/clients"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/db"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

const (
	SourceRegistry   = "registry"
	SourceSubscriber = "subscriber"
	SourceDataplan   = "dataplan"
	SourceMetrics    = "metrics"
	SourceNode       = "node"
	SourceInventory  = "inventory"
	SourceBilling    = "billing"

	StatusOk      = "ok"
	StatusRunning = "running"
	StatusFailed  = "failed"
)

var Sources = []string{SourceRegistry, SourceSubscriber, SourceDataplan,
	SourceMetrics, SourceNode, SourceInventory, SourceBilling}

// Refresher orchestrates per-source snapshot refreshes: it calls the source
// client, upserts snapshots, transitions the refresh state row
// (running -> ok/failed) and marks the affected rollups dirty.
type Refresher struct {
	stateRepo    db.StateRepo
	snapshotRepo db.SnapshotRepo
	factRepo     db.FactRepo

	registryClient   clients.RegistryClient
	subscriberClient clients.SubscriberClient
	dataplanClient   clients.DataplanClient
	metricsClient    clients.MetricsClient
	nodeClient       clients.NodeClient
	inventoryClient  clients.InventoryClient
	billingClient    clients.BillingClient
}

func NewRefresher(stateRepo db.StateRepo, snapshotRepo db.SnapshotRepo, factRepo db.FactRepo,
	registryClient clients.RegistryClient, subscriberClient clients.SubscriberClient,
	dataplanClient clients.DataplanClient, metricsClient clients.MetricsClient,
	nodeClient clients.NodeClient, inventoryClient clients.InventoryClient,
	billingClient clients.BillingClient) *Refresher {
	return &Refresher{
		stateRepo:        stateRepo,
		snapshotRepo:     snapshotRepo,
		factRepo:         factRepo,
		registryClient:   registryClient,
		subscriberClient: subscriberClient,
		dataplanClient:   dataplanClient,
		metricsClient:    metricsClient,
		nodeClient:       nodeClient,
		inventoryClient:  inventoryClient,
		billingClient:    billingClient,
	}
}

// Refresh runs a refresh for the given source and returns the final state.
func (r *Refresher) Refresh(source string) (*db.RefreshState, error) {
	now := time.Now().UTC()

	state := &db.RefreshState{
		Source:    source,
		Status:    StatusRunning,
		Detail:    "",
		LastRunAt: now,
	}

	if err := r.stateRepo.UpsertRefreshState(state); err != nil {
		return nil, err
	}

	var err error

	switch source {
	case SourceRegistry:
		err = r.refreshRegistry(now)
	case SourceSubscriber:
		err = r.refreshSubscriber(now)
	case SourceDataplan:
		err = r.refreshDataplan(now)
	case SourceMetrics:
		err = r.refreshMetrics(now)
	case SourceNode:
		err = r.refreshNode(now)
	case SourceInventory:
		err = r.refreshInventory(now)
	case SourceBilling:
		err = r.refreshBilling(now)
	default:
		err = fmt.Errorf("unknown refresh source: %s", source)
	}

	if err != nil {
		log.Errorf("refresh of source %s failed: %v", source, err)

		state.Status = StatusFailed
		state.Detail = err.Error()
	} else {
		state.Status = StatusOk
		state.Detail = ""
		state.LastSuccessAt = time.Now().UTC()
	}

	if uerr := r.stateRepo.UpsertRefreshState(state); uerr != nil {
		return nil, uerr
	}

	return state, err
}

func (r *Refresher) refreshRegistry(now time.Time) error {
	networks, err := r.registryClient.GetNetworks()
	if err != nil {
		return err
	}

	for _, n := range networks {
		id, perr := uuid.FromString(n.Id)
		if perr != nil {
			log.Warnf("skipping network with invalid id %q: %v", n.Id, perr)

			continue
		}

		status := "active"
		if n.IsDeactivated {
			status = "inactive"
		}

		if err := r.snapshotRepo.UpsertNetwork(&db.NetworkSnapshot{
			NetworkId: id,
			Name:      n.Name,
			Status:    status,
			UpdatedAt: now,
		}); err != nil {
			return err
		}
	}

	sites, err := r.registryClient.GetSites()
	if err != nil {
		return err
	}

	for _, s := range sites {
		id, perr := uuid.FromString(s.Id)
		if perr != nil {
			log.Warnf("skipping site with invalid id %q: %v", s.Id, perr)

			continue
		}

		netId, _ := uuid.FromString(s.NetworkId)
		lat, _ := strconv.ParseFloat(s.Latitude, 64)
		lng, _ := strconv.ParseFloat(s.Longitude, 64)

		status := "online"
		if s.IsDeactivated {
			status = "offline"
		}

		if err := r.snapshotRepo.UpsertSite(&db.SiteSnapshot{
			SiteId:    id,
			NetworkId: netId,
			Name:      s.Name,
			Status:    status,
			Latitude:  lat,
			Longitude: lng,
			UpdatedAt: now,
		}); err != nil {
			return err
		}
	}

	return r.stateRepo.MarkRollupDirty("network_health_hourly")
}

func (r *Refresher) refreshSubscriber(now time.Time) error {
	subs, err := r.subscriberClient.GetSubscribers()
	if err != nil {
		return err
	}

	for _, s := range subs {
		id, perr := uuid.FromString(s.SubscriberId)
		if perr != nil {
			log.Warnf("skipping subscriber with invalid id %q: %v", s.SubscriberId, perr)

			continue
		}

		netId, _ := uuid.FromString(s.NetworkId)

		snap := &db.CustomerSnapshot{
			CustomerId: id,
			NetworkId:  netId,
			Name:       s.Name,
			Email:      s.Email,
			Status:     "active",
			UpdatedAt:  now,
		}

		if t, terr := time.Parse(time.RFC3339, s.CreatedAt); terr == nil {
			snap.SourceCreatedAt = &t
		}

		if err := r.snapshotRepo.UpsertCustomer(snap); err != nil {
			return err
		}
	}

	return r.stateRepo.MarkRollupDirty("customer_state_daily")
}

func (r *Refresher) refreshDataplan(now time.Time) error {
	packages, err := r.dataplanClient.GetPackages()
	if err != nil {
		return err
	}

	for _, p := range packages {
		id, perr := uuid.FromString(p.Uuid)
		if perr != nil {
			log.Warnf("skipping package with invalid id %q: %v", p.Uuid, perr)

			continue
		}

		status := "inactive"
		if p.IsActive {
			status = "active"
		}

		if err := r.snapshotRepo.UpsertPackage(&db.PackageSnapshot{
			PackageId:    id,
			Name:         p.Name,
			Price:        p.Amount,
			Currency:     p.Currency,
			DurationDays: uint32(p.Duration),
			DataQuotaMb:  float64(p.DataVolume),
			Status:       status,
			UpdatedAt:    now,
		}); err != nil {
			return err
		}
	}

	return r.stateRepo.MarkRollupDirty("business_package_daily")
}

func (r *Refresher) refreshMetrics(now time.Time) error {
	metrics, err := r.metricsClient.GetLatestMetrics()
	if err != nil {
		return err
	}

	for _, m := range metrics {
		sampledAt := now
		if m.Timestamp > 0 {
			sampledAt = time.Unix(m.Timestamp, 0).UTC()
		}

		if err := r.factRepo.AddMetricSample(&db.MetricSample{
			Metric:       m.Metric,
			ResourceType: m.ResourceType,
			ResourceId:   m.ResourceId,
			Value:        m.Value,
			Unit:         m.Unit,
			SampledAt:    sampledAt,
		}); err != nil {
			return err
		}
	}

	return r.stateRepo.MarkRollupDirty("metric_hourly")
}

func (r *Refresher) refreshNode(now time.Time) error {
	nodes, err := r.nodeClient.GetNodes()
	if err != nil {
		return err
	}

	for _, n := range nodes {
		siteId, _ := uuid.FromString(n.SiteId)
		netId, _ := uuid.FromString(n.NetworkId)

		if err := r.snapshotRepo.UpsertNode(&db.NodeSnapshot{
			NodeId:       n.Id,
			SiteId:       siteId,
			NetworkId:    netId,
			Name:         n.Name,
			Type:         n.Type,
			Status:       n.State,
			Connectivity: n.Connectivity,
			UpdatedAt:    now,
		}); err != nil {
			return err
		}
	}

	return r.stateRepo.MarkRollupDirty("node_health_hourly")
}

func (r *Refresher) refreshInventory(now time.Time) error {
	components, err := r.inventoryClient.GetComponents()
	if err != nil {
		return err
	}

	for _, c := range components {
		if err := r.snapshotRepo.UpsertInventory(&db.InventorySnapshot{
			ComponentId: c.Id,
			Type:        c.Type,
			State:       c.Inventory,
			UpdatedAt:   now,
		}); err != nil {
			return err
		}
	}

	return r.stateRepo.MarkRollupDirty("business_inventory_daily")
}

func (r *Refresher) refreshBilling(now time.Time) error {
	account, err := r.billingClient.GetAccount()
	if err != nil {
		return err
	}

	snap := &db.BillingSnapshot{
		Id:                  1,
		Balance:             account.Balance,
		PaymentMethodStatus: account.PaymentMethodStatus,
		UpdatedAt:           now,
	}

	if t, terr := time.Parse(time.RFC3339, account.LastInvoiceAt); terr == nil {
		snap.LastInvoiceAt = &t
	}

	if err := r.snapshotRepo.UpsertBilling(snap); err != nil {
		return err
	}

	return r.stateRepo.MarkRollupDirty("business_billing_daily")
}
