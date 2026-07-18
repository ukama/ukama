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

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	hpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const logMsgRoutingKey = "Received a message with Routing key %s and Message %+v"
const errFailedUpdateSoftware = "failed to update software: %w"

type SoftwareUpdateEventServer struct {
	s       *SoftwareServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewSoftwareEventServer(orgName string, s *SoftwareServer) *SoftwareUpdateEventServer {
	return &SoftwareUpdateEventServer{
		s:       s,
		orgName: orgName,
	}
}
func (n *SoftwareUpdateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof(logMsgRoutingKey, e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, "event.cloud.global.{{ .Org}}.hub.artifactmanager.app.uploaded"):
		return n.handleArtifactUploadedEvent(ctx, e)
	case msgbus.PrepareRoute(n.orgName, "event.cloud.global.{{ .Org}}.hub.distributor.app.chunkready"):
		return n.handleNodeAppChunkReadyEvent(ctx, e)
	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline]):
		return n.handleNodeOnlineEvent(ctx, e)

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

// objectTypeFromRoutingKey extracts the artifact type from ...hub.<service>.<object>.<action>.
func objectTypeFromRoutingKey(rk string) string {
	parts := strings.Split(rk, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return "app"
}

func (n *SoftwareUpdateEventServer) handleNodeOnlineEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, e.RoutingKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node online event: %w", err)
	}

	nodeType := ukama.GetNodeType(msg.NodeId)
	if nodeType == nil {
		return nil, fmt.Errorf("failed to get node type for node %s: %w", msg.NodeId, err)
	}

	err = n.s.nodeRepo.Create(db.Node{
		NodeId:   msg.NodeId,
		NodeType: ukama.NodeType(*nodeType),
	})
	if err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, fmt.Errorf("failed to create node %s: %w", msg.NodeId, err)
		}
		log.Infof("Node %s already exists, reusing existing record", msg.NodeId)
	}

	// Health-driven reconcile of installed apps/versions. Safe when health isn't ready
	// yet (returns empty); the periodic reconcile and watchers converge later.
	if err := n.reconcileApps(msg.NodeId); err != nil {
		log.Errorf("failed to reconcile apps for node %s: %v", msg.NodeId, err)
	}
	if err := n.reconcileSoftware(msg.NodeId); err != nil {
		log.Errorf("failed to reconcile software for node %s: %v", msg.NodeId, err)
	}

	// Catalog-driven: a newly-onboarded node immediately learns of available updates
	// from the persistent desired-release store — the fix for the availability/timing break.
	if err := n.applyDesiredToNode(msg.NodeId); err != nil {
		log.Errorf("failed to apply desired to node %s: %v", msg.NodeId, err)
	}

	return &epb.EventResponse{}, nil
}

// applyDesiredToNode sets each of the node's software rows to the fleet-wide desired
// version (from the catalog) and flags UpdateAvailable where the node lags.
func (n *SoftwareUpdateEventServer) applyDesiredToNode(nodeID string) error {
	rows, err := n.s.sRepo.List(nodeID, ukama.Unknown, "")
	if err != nil {
		return err
	}
	for _, sw := range rows {
		d, err := n.s.releaseRepo.GetDesired(sw.AppName, "app")
		if err != nil || d == nil {
			continue
		}
		sw.DesiredVersion = d.DesiredVersion
		if validation.IsVersionMismatch(sw.CurrentVersion, d.DesiredVersion) {
			sw.Status = ukama.SoftwareStatusType(ukama.UpdateAvailable)
		} else {
			sw.Status = ukama.SoftwareStatusType(ukama.UpToDate)
		}
		if err := n.s.sRepo.Update(sw); err != nil {
			log.Errorf("applyDesiredToNode: update %s: %v", sw.Id, err)
		}
	}
	return nil
}

// handleArtifactUploadedEvent records availability only — it does NOT set desired
// or touch any node's status. Promotion is an explicit, separate action.
func (n *SoftwareUpdateEventServer) handleArtifactUploadedEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	p := &epb.EventArtifactUploaded{}
	if err := anypb.UnmarshalTo(e.Msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artifact uploaded event: %w", err)
	}
	aType := objectTypeFromRoutingKey(e.RoutingKey)
	if err := n.s.releaseRepo.Upsert(&db.ReleaseCatalog{
		Name: p.Name, Type: aType, Version: p.Version, Available: true, UploadedAt: time.Now(),
	}); err != nil {
		return nil, fmt.Errorf("failed to upsert release catalog: %w", err)
	}
	log.Infof("catalog: %s/%s@%s recorded available", aType, p.Name, p.Version)
	return &epb.EventResponse{}, nil
}

// handleNodeAppChunkReadyEvent only flags the catalog row as chunked — it no longer
// drives any node's desired version.
func (n *SoftwareUpdateEventServer) handleNodeAppChunkReadyEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	p := &epb.EventArtifactChunkReady{}
	if err := anypb.UnmarshalTo(e.Msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal node app chunk ready event: %w", err)
	}
	aType := objectTypeFromRoutingKey(e.RoutingKey)
	// chunkready can arrive before uploaded — ensure the row exists, then flag it.
	if err := n.s.releaseRepo.Upsert(&db.ReleaseCatalog{Name: p.Name, Type: aType, Version: p.Version, Available: true}); err != nil {
		return nil, fmt.Errorf("failed to upsert release catalog: %w", err)
	}
	if err := n.s.releaseRepo.SetChunked(p.Name, aType, p.Version); err != nil {
		return nil, fmt.Errorf("failed to mark chunked: %w", err)
	}
	log.Infof("catalog: %s/%s@%s chunk index ready", aType, p.Name, p.Version)
	return &epb.EventResponse{}, nil
}

func (n *SoftwareUpdateEventServer) reconcileSoftware(nodeID string) error {
	log.Infof("Reconciling software for node %s", nodeID)
	hrApps, err := n.listApps(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get apps for node %s: %w", nodeID, err)
	}
	if len(hrApps.Apps) == 0 {
		log.Infof("No apps found for node %s", nodeID)
		return nil
	}

	for _, capp := range hrApps.Apps {
		listSoftware, err := n.s.sRepo.List(nodeID, ukama.Unknown, capp.Name)
		if err != nil {
			return fmt.Errorf("failed to list software for node %s: %w", nodeID, err)
		}

		if len(listSoftware) == 0 {
			err := n.s.sRepo.Create(&db.Software{
				Id:             uuid.NewV4(),
				NodeId:         nodeID,
				AppName:        capp.Name,
				CurrentVersion: capp.Version,
				DesiredVersion: "",
				ReleaseDate:    time.Now(),
				Status:         ukama.UpToDate,
				ChangeLogs:     []string{},
			})
			if err != nil {
				return fmt.Errorf("failed to create software: %w", err)
			}

		} else {
			err := n.s.sRepo.Update(&db.Software{
				Id:             listSoftware[0].Id,
				NodeId:         nodeID,
				AppName:        capp.Name,
				CurrentVersion: capp.Version,
				DesiredVersion: listSoftware[0].DesiredVersion,
				ReleaseDate:    listSoftware[0].ReleaseDate,
				Status:         listSoftware[0].Status,
				ChangeLogs:     listSoftware[0].ChangeLogs,
			})
			if err != nil {
				return fmt.Errorf(errFailedUpdateSoftware, err)
			}
		}
	}
	log.Infof("Reconciled software for node %s", nodeID)
	return nil
}

func (n *SoftwareUpdateEventServer) reconcileApps(nodeID string) error {
	log.Infof("Reconciling apps for node %s", nodeID)

	nID, nodeType, err := ukama.ValidateNodeIdAndType(nodeID)
	if err != nil {
		return err
	}

	hrApps, err := n.listApps(nodeID)
	if err != nil {
		return fmt.Errorf("failed to list apps: %w", err)
	}

	if len(hrApps.Apps) == 0 {
		log.Infof("No apps found from health for node %s", nID.String())
		return nil
	}

	apps, err := n.s.appRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get apps: %w", err)
	}

	missingApps := findMissingApps(hrApps.Apps, apps, nodeType)
	if err := n.createMissingApps(missingApps); err != nil {
		return err
	}

	return nil
}

func findMissingApps(hrApps []*hpb.App, apps []db.App, nodeType *string) map[string]db.App {
	appNames := make(map[string]struct{}, len(apps))
	for _, app := range apps {
		appNames[strings.ToLower(app.Name)] = struct{}{}
	}

	missingByName := make(map[string]db.App)
	for _, capp := range hrApps {
		name := strings.TrimSpace(capp.Name)
		if name == "" {
			continue
		}

		lowerName := strings.ToLower(name)
		if _, found := appNames[lowerName]; found {
			continue
		}
		if _, exists := missingByName[lowerName]; exists {
			continue
		}
		uid := uuid.NewV4()
		missingByName[lowerName] = db.App{
			Id:          uid,
			Name:        name,
			Space:       "system",
			MetricsKeys: []string{name + "_cpu", name + "_memory", name + "_disk"},
			Notes:       "App is installed on " + ukama.GetPlaceholderNameByType(*nodeType),
		}
	}

	return missingByName
}

func (n *SoftwareUpdateEventServer) createMissingApps(missingApps map[string]db.App) error {
	for _, app := range missingApps {
		if err := n.s.appRepo.Create(app); err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				continue
			}
			return fmt.Errorf("failed to create app: %w", err)
		}
	}
	return nil
}

func (n *SoftwareUpdateEventServer) listApps(nodeId string) (*hpb.ListAppsResponse, error) {
	healthClient, err := n.s.healthClient.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get health client: %w", err)
	}
	response, err := healthClient.ListApps(context.Background(), &hpb.ListAppsRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get apps: %w", err)
	}
	return response, nil
}
