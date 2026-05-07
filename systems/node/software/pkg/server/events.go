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
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
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
	case msgbus.PrepareRoute(n.orgName, evt.NodeEventToEventConfig[evt.NodeAppChunkReady].RoutingKey):
		return n.handleNodeAppChunkReadyEvent(ctx, e)
	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline]):
		log.Infof(logMsgRoutingKey, e.RoutingKey, e.Msg)
		return n.handleNodeOnlineEvent(ctx, e)

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
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

	time.AfterFunc(120*time.Second, func() {
		if err := n.reconcileApps(msg.NodeId); err != nil {
			log.Errorf("failed to reconcile apps for node %s: %v", msg.NodeId, err)
		}
		if err := n.reconcileSoftware(msg.NodeId); err != nil {
			log.Errorf("failed to reconcile software for node %s: %v", msg.NodeId, err)
		}
	})

	return &epb.EventResponse{}, nil
}

func (n *SoftwareUpdateEventServer) handleNodeAppChunkReadyEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	p := &epb.EventArtifactChunkReady{}
	err := anypb.UnmarshalTo(e.Msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node app chunk ready event: %w", err)
	}

	log.Infof(logMsgRoutingKey, e.RoutingKey, p)

	err = n.reconcileCurrentAppVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to reconcile current app version: %w", err)
	}

	softwares, err := n.s.sRepo.List("", ukama.Unknown, p.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get softwares: %w", err)
	}

	for _, software := range softwares {
		software.DesiredVersion = p.Version
		software.ReleaseDate = time.Now()
		if validation.IsVersionMismatch(software.CurrentVersion, p.Version) {
			software.Status = ukama.SoftwareStatusType(ukama.UpdateAvailable)
			software.ChangeLogs = append(software.ChangeLogs, "New version "+p.Version+" available. Please update to the latest version.")
		} else {
			software.Status = ukama.SoftwareStatusType(ukama.UpToDate)
		}

		err = n.s.sRepo.Update(software)
		if err != nil {
			return nil, fmt.Errorf(errFailedUpdateSoftware, err)
		}
	}
	log.Infof("Updated software update for app %s and version %s", p.Name, p.Version)

	return &epb.EventResponse{}, nil
}

func (n *SoftwareUpdateEventServer) reconcileCurrentAppVersion() error {
	log.Infof("Reconciling current App version for all nodes")

	nodes, err := n.s.nodeRepo.List()
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	for _, node := range nodes {
		if err := n.reconcileCurrentVersionForNode(node.NodeId); err != nil {
			return err
		}
	}
	return nil
}

func (n *SoftwareUpdateEventServer) reconcileCurrentVersionForNode(nodeID string) error {
	cApps, err := n.listCApps(nodeID, "")
	if err != nil {
		return fmt.Errorf("failed to list capps for node %s: %w", nodeID, err)
	}

	for _, capp := range cApps.Capps {
		app, err := n.s.sRepo.List(nodeID, ukama.Unknown, capp.Name)
		if err != nil {
			return fmt.Errorf("failed to get app %s: %w", capp.Name, err)
		}
		if len(app) == 0 || app[0].CurrentVersion == capp.Tag {
			continue
		}

		app[0].CurrentVersion = capp.Tag
		err = n.s.sRepo.Update(app[0])
		if err != nil {
			return fmt.Errorf(errFailedUpdateSoftware, err)
		}
	}
	return nil
}

func (n *SoftwareUpdateEventServer) reconcileSoftware(nodeID string) error {
	log.Infof("Reconciling software for node %s", nodeID)
	healthReport, err := n.getHealthReport(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get health report for node %s: %w", nodeID, err)
	}
	if len(healthReport.Healths) == 0 {
		log.Infof("No health report found for node %s", nodeID)
		return nil
	}

	capps := healthReport.Healths[0].Capps

	for _, capp := range capps {
		listSoftware, err := n.s.sRepo.List(nodeID, ukama.Unknown, capp.Name)
		if err != nil {
			return fmt.Errorf("failed to list software for node %s: %w", nodeID, err)
		}

		if len(listSoftware) == 0 {
			err := n.s.sRepo.Create(&db.Software{
				Id:             uuid.NewV4(),
				NodeId:         nodeID,
				AppName:        capp.Name,
				CurrentVersion: capp.Tag,
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
				CurrentVersion: capp.Tag,
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

	capps, err := n.listCApps(nodeID, "")
	if err != nil {
		return fmt.Errorf("failed to list capps: %w", err)
	}

	if len(capps.Capps) == 0 {
		log.Infof("No apps found from health for node %s", nID.String())
		return nil
	}

	apps, err := n.s.appRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get apps: %w", err)
	}

	missingApps := findMissingApps(capps.Capps, apps, nodeType)
	if err := n.createMissingApps(missingApps); err != nil {
		return err
	}

	return nil
}

func findMissingApps(capps []*hpb.Capps, apps []db.App, nodeType *string) map[string]db.App {
	appNames := make(map[string]struct{}, len(apps))
	for _, app := range apps {
		appNames[strings.ToLower(app.Name)] = struct{}{}
	}

	missingByName := make(map[string]db.App)
	for _, capp := range capps {
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

func (n *SoftwareUpdateEventServer) getHealthReport(nodeID string) (*hpb.ListResponse, error) {
	healthClient, err := n.s.healthClient.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get health client: %w", err)
	}
	healthReport, err := healthClient.List(context.Background(), &hpb.ListRequest{
		NodeId:    nodeID,
		Timeframe: upb.FilterTimeframesType_LATEST,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get health report for node %s: %w", nodeID, err)
	}
	return healthReport, nil
}

func (n *SoftwareUpdateEventServer) listCApps(nodeId string, name string) (*hpb.ListAppsResponse, error) {
	healthClient, err := n.s.healthClient.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get health client: %w", err)
	}
	response, err := healthClient.ListApps(context.Background(), &hpb.ListAppsRequest{
		NodeId: nodeId,
		Name:   name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list capps: %w", err)
	}
	return response, nil
}
