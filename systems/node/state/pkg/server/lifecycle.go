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
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	stm "github.com/ukama/ukama/systems/common/stateMachine"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
	"github.com/ukama/ukama/systems/node/state/pkg/lifecycle"
)

func (n *StateEventServer) SetLifecycleRepo(repo *db.LifecycleRepo) {
	n.lifecycleRepo = repo
}

func (n *StateEventServer) processStoredEvent(ctx context.Context, event, nodeID string, message interface{}) error {
	id, err := ukama.ValidateNodeId(nodeID)
	if err != nil {
		return fmt.Errorf("invalid node ID: %w", err)
	}
	return n.lifecycleRepo.Update(ctx, id.String(), event, func(record *db.LifecycleRecord, state *db.State) error {
		if state.CurrentState == npb.NodeState_Offboarded && event != "online" && event != "offline" && event != "onboarding" {
			return nil
		}
		state.NodeType = id.GetNodeType()
		switch msg := message.(type) {
		case *epb.EventRegistryNodeAssign:
			return n.assignNode(record, state, msg)
		case *lifecycle.Event:
			return n.observeLifecycle(record, state, msg)
		case *epb.NodeOnlineEvent:
			state.Config = &db.NodeConfig{
				Id: uuid.NewV4(), NodeId: id.String(), NodeIp: msg.NodeIp,
				NodePort: int32(msg.NodePort), MeshIp: msg.MeshIp,
				MeshPort: int32(msg.MeshPort), MeshHostName: msg.MeshHostName,
			}
		}
		// Configuration completion is accepted only through lifecycle metadata.
		switch event {
		case "assignnoconfig", "configapplied", "operational", "configuring", "init", "platformready":
			return fmt.Errorf("event %s requires a lifecycle observation", event)
		}
		if err := n.transitionStoredState(state, event); err != nil {
			return err
		}
		if state.CurrentState == npb.NodeState_Offboarded {
			record.AwaitingOperational = false
		}
		return nil
	})
}

func (n *StateEventServer) assignNode(record *db.LifecycleRecord, state *db.State, msg *epb.EventRegistryNodeAssign) error {
	if !routingPart(msg.Network) || !routingPart(msg.Site) {
		return fmt.Errorf("assignment needs network and site routing identities")
	}
	if record.RequestID != "" {
		if record.Network != msg.Network || record.Site != msg.Site {
			return fmt.Errorf("node already has a different lifetime assignment")
		}
		return nil
	}
	record.Network = msg.Network
	record.Site = msg.Site
	record.RequestID = "assignment-" + uuid.NewV4().String()
	record.AwaitingOperational = true
	record.RetryAt = time.Now().UTC()
	return n.transitionStoredState(state, "assign")
}

func routingPart(value string) bool {
	return value != "" && !strings.ContainsAny(value, ".*# /\\\t\r\n")
}

func (n *StateEventServer) observeLifecycle(record *db.LifecycleRecord, state *db.State, event *lifecycle.Event) error {
	accepted, err := record.Cursor.Accept(event)
	if err != nil || !accepted {
		return err
	}
	switch event.State {
	case "INIT":
		if record.RequestID != "" {
			record.AwaitingOperational = true
			record.RetryAt = time.Now().UTC().Add(60 * time.Second)
		}
		return n.transitionStoredState(state, "init")
	case "READY":
		// READY never creates or resends an assignment.
		return n.transitionStoredState(state, "platformready")
	case "FAULTY":
		return n.transitionStoredState(state, "fault")
	}
	if record.RequestID == "" {
		return fmt.Errorf("configuration result arrived before assignment")
	}
	if event.ConfigMode == "NOCONFIG" && event.RequestID != record.RequestID {
		return fmt.Errorf("configuration request does not match assignment")
	}
	// This drop establishes the initial NOCONFIG decision. Later CONFIG results
	// may be observed only after the assignment has already been confirmed.
	if event.ConfigMode == "CONFIG" && !record.Completed {
		return fmt.Errorf("initial assignment requires NOCONFIG completion")
	}
	if event.State == "CONFIGURING" {
		return n.transitionStoredState(state, "configuring")
	}
	if err := n.transitionStoredState(state, "operational"); err != nil {
		return err
	}
	if state.CurrentState != npb.NodeState_Operational {
		return fmt.Errorf("OPERATIONAL cannot complete configuration from %s", state.CurrentState)
	}
	record.Completed = true
	record.AwaitingOperational = false
	return nil
}

func (n *StateEventServer) transitionStoredState(state *db.State, event string) error {
	machine := stm.NewStateMachine(nil)
	instance, err := machine.NewInstance(n.configPath, state.NodeId, state.CurrentState.String())
	if err != nil {
		return err
	}
	instance.CurrentSubstate = "off"
	if len(state.SubState) != 0 {
		instance.CurrentSubstate = state.SubState[len(state.SubState)-1]
	}
	if err := instance.Transition(event); err != nil {
		return err
	}
	value, ok := npb.NodeState_value[instance.CurrentState]
	if !ok {
		return fmt.Errorf("unknown state %s", instance.CurrentState)
	}
	state.CurrentState = npb.NodeState(value)
	state.SubState = db.StringArray{instance.CurrentSubstate}
	return nil
}

func (n *StateEventServer) sendNoConfig(record *db.LifecycleRecord) {
	body, err := json.Marshal(struct {
		Mode      string `json:"mode"`
		RequestID string `json:"requestId"`
	}{Mode: "NOCONFIG", RequestID: record.RequestID})
	if err != nil {
		log.Errorf("Cannot encode NOCONFIG: %v", err)
		return
	}
	route := n.baseRoutingKey.SetRequestType().SetObject("nodefeeder").SetAction("publish").MustBuild()
	message := &epb.NodeFeederMessage{
		NodeId:     record.NodeID,
		Target:     strings.Join([]string{n.orgName, record.Network, record.Site, strings.ToLower(record.NodeID)}, "."),
		HttpMethod: "POST", Path: "configd/v1/config", Msg: body,
	}
	if err := n.msgbus.PublishRequest(route, message); err != nil {
		log.Errorf("NOCONFIG for node %s was not published; retrying in 60 seconds: %v", record.NodeID, err)
		return
	}
}

func (n *StateEventServer) publishLifecycle(item *db.LifecyclePublication) error {
	route := n.baseRoutingKey.SetAction("transition").SetObject("node").MustBuild()
	return n.msgbus.PublishRequest(route, &epb.NodeStateChangeEvent{
		NodeId: item.NodeID, State: item.State, Substate: item.Substate, Events: []string{item.Event},
	})
}

func (n *StateEventServer) StartLifecycleWorker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			if err := n.lifecycleRepo.RetryAssignments(ctx, time.Now().UTC(), n.sendNoConfig); err != nil {
				log.Errorf("Assignment retry failed: %v", err)
			}
			if err := n.lifecycleRepo.PublishPending(ctx, n.publishLifecycle); err != nil {
				log.Errorf("Lifecycle publication failed: %v", err)
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}

// Read the current state inside the same transaction used for lifecycle events.
// A timeout sweep's earlier snapshot must never overwrite a newer observation.
func (n *StateEventServer) runStoredTimeout(ctx context.Context, nodeID string, now time.Time) (bool, error) {
	moved := false
	err := n.lifecycleRepo.Update(ctx, nodeID, TimeoutEventName, func(_ *db.LifecycleRecord, state *db.State) error {
		if state.CurrentState != npb.NodeState_Updating {
			return nil
		}
		machine := stm.NewStateMachine(nil)
		instance, err := machine.NewInstance(n.configPath, nodeID, state.CurrentState.String())
		if err != nil {
			return err
		}
		if len(state.SubState) != 0 {
			instance.CurrentSubstate = state.SubState[len(state.SubState)-1]
		}
		moved, err = instance.TimeoutTransition(state.CreatedAt, now)
		if err != nil || !moved {
			return err
		}
		state.CurrentState = npb.NodeState(npb.NodeState_value[instance.CurrentState])
		state.SubState = db.StringArray{instance.CurrentSubstate}
		return nil
	})
	return moved && err == nil, err
}
