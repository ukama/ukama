/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package reconciler

import (
	"context"
	"fmt"

	"github.com/ukama/ukama/systems/node/site-controller/pkg/adapters"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/policy"
)

type Reconciler struct {
	intents    db.IntentRepo
	states     db.StateRepo
	ports      db.PortMapRepo
	components db.ComponentRepo
	tower      *adapters.TowerAdapter
	amplifier  *adapters.AmplifierAdapter
	cnode      *adapters.CNodeAdapter
}

func New(intents db.IntentRepo, states db.StateRepo, ports db.PortMapRepo, components db.ComponentRepo, tower *adapters.TowerAdapter, amp *adapters.AmplifierAdapter, cnode *adapters.CNodeAdapter) *Reconciler {
	return &Reconciler{intents: intents, states: states, ports: ports, components: components, tower: tower, amplifier: amp, cnode: cnode}
}

// GetSnapshot returns desired intent, derived state, component snapshot JSON, and static port map.
func (r *Reconciler) GetSnapshot(ctx context.Context, siteID string) (*SiteSnapshot, error) {
	st, intent, err := r.GetState(ctx, siteID)
	if err != nil {
		return nil, err
	}
	ports, err := r.ports.GetBySite(siteID)
	if err != nil {
		return nil, err
	}
	var componentsJSON string
	c, err := r.components.Get(siteID)
	if err != nil {
		return nil, err
	}
	if c != nil {
		componentsJSON = c.Components
	}
	return &SiteSnapshot{Intent: intent, DerivedState: st, ComponentsJSON: componentsJSON, Ports: ports}, nil
}
func (r *Reconciler) GetState(ctx context.Context, siteID string) (*db.SiteState, *db.SiteIntent, error) {
	intent, err := r.getIntent(siteID)
	if err != nil {
		return nil, nil, err
	}
	state, err := r.states.Get(siteID)
	if err != nil {
		return nil, nil, err
	}
	if state == nil {
		state = derive(intent)
		_ = r.states.Upsert(state)
	}
	return state, intent, nil
}
func (r *Reconciler) UpsertPortMap(ctx context.Context, siteID, cnodeID string, ports []db.SitePortMap) error {
	if err := policy.ValidatePortMap(ports); err != nil {
		return err
	}
	return r.ports.Upsert(siteID, cnodeID, ports)
}
func (r *Reconciler) GetPortMap(ctx context.Context, siteID string) ([]db.SitePortMap, error) {
	return r.ports.GetBySite(siteID)
}
func (r *Reconciler) ApplySwitchPolicy(ctx context.Context, siteID string) error {
	ports, err := r.ports.GetBySite(siteID)
	if err != nil {
		return err
	}
	sp, err := policy.BuildSwitchPolicy(siteID, ports)
	if err != nil {
		return err
	}
	cnodePort, err := policy.FindRole(ports, policy.RoleCNode)
	if err != nil {
		return err
	}
	if cnodePort.NodeID == "" {
		return fmt.Errorf("cnode node_id missing")
	}
	return r.cnode.ApplySwitchPolicy(ctx, cnodePort.NodeID, sp)
}
func (r *Reconciler) SetSite(ctx context.Context, siteID, state, reason, requestedBy string) (*db.SiteState, error) {
	if state != StateOn && state != StateOff {
		return nil, fmt.Errorf("invalid site state %s", state)
	}
	intent, err := r.getIntent(siteID)
	if err != nil {
		return nil, err
	}
	intent.Reason = reason
	if requestedBy != "" {
		intent.RequestedBy = requestedBy
	}
	if state == StateOn {
		intent.DesiredService = StateOn
		intent.DesiredRadio = StateOn
	} else {
		intent.DesiredService = StateOff
		intent.DesiredRadio = StateOff
	}
	if err := r.intents.Upsert(intent); err != nil {
		return nil, err
	}
	if err := r.reconcile(ctx, intent); err != nil {
		return nil, err
	}
	st := derive(intent)
	return st, r.states.Upsert(st)
}
func (r *Reconciler) SetService(ctx context.Context, siteID, state, reason, requestedBy string) (*db.SiteState, error) {
	if state != StateOn && state != StateOff {
		return nil, fmt.Errorf("invalid service state %s", state)
	}
	intent, err := r.getIntent(siteID)
	if err != nil {
		return nil, err
	}
	intent.DesiredService = state
	intent.Reason = reason
	if requestedBy != "" {
		intent.RequestedBy = requestedBy
	}
	if err := r.intents.Upsert(intent); err != nil {
		return nil, err
	}
	if err := r.applyService(ctx, siteID, state); err != nil {
		return nil, err
	}
	st := derive(intent)
	return st, r.states.Upsert(st)
}
func (r *Reconciler) SetRadio(ctx context.Context, siteID, state, reason, requestedBy string) (*db.SiteState, error) {
	if state != StateOn && state != StateOff {
		return nil, fmt.Errorf("invalid radio state %s", state)
	}
	intent, err := r.getIntent(siteID)
	if err != nil {
		return nil, err
	}
	intent.DesiredRadio = state
	intent.Reason = reason
	if requestedBy != "" {
		intent.RequestedBy = requestedBy
	}
	if err := r.intents.Upsert(intent); err != nil {
		return nil, err
	}
	if err := r.applyRadio(ctx, siteID, state); err != nil {
		return nil, err
	}
	st := derive(intent)
	return st, r.states.Upsert(st)
}
// PowerCycleNode looks up the port for role and forwards PoE cycle to the CNode. It does not reject
// the cnode role; switch.d on the node enforces never_off_remote and related policy.
func (r *Reconciler) PowerCycleNode(ctx context.Context, siteID, role, reason string) error {
	ports, err := r.ports.GetBySite(siteID)
	if err != nil {
		return err
	}
	target, err := policy.FindRole(ports, role)
	if err != nil {
		return err
	}
	cnodePort, err := policy.FindRole(ports, policy.RoleCNode)
	if err != nil {
		return err
	}
	return r.cnode.PowerCyclePort(ctx, cnodePort.NodeID, target.Port, reason)
}
func (r *Reconciler) getIntent(siteID string) (*db.SiteIntent, error) {
	intent, err := r.intents.Get(siteID)
	if err != nil {
		return nil, err
	}
	if intent == nil {
		intent = &db.SiteIntent{SiteID: siteID, DesiredService: StateOff, DesiredRadio: StateOff, Reason: "initial"}
	}
	return intent, nil
}
func (r *Reconciler) reconcile(ctx context.Context, intent *db.SiteIntent) error {
	if err := r.ApplySwitchPolicy(ctx, intent.SiteID); err != nil {
		return err
	}
	if intent.DesiredService == StateOn {
		if err := r.ensureCriticalPoe(ctx, intent.SiteID); err != nil {
			return err
		}
		if err := r.applyRadio(ctx, intent.SiteID, StateOn); err != nil {
			return err
		}
		return r.applyService(ctx, intent.SiteID, StateOn)
	}
	if err := r.applyRadio(ctx, intent.SiteID, StateOff); err != nil {
		return err
	}
	return r.applyService(ctx, intent.SiteID, StateOff)
}
func (r *Reconciler) ensureCriticalPoe(ctx context.Context, siteID string) error {
	ports, err := r.ports.GetBySite(siteID)
	if err != nil {
		return err
	}
	cnodePort, err := policy.FindRole(ports, policy.RoleCNode)
	if err != nil {
		return err
	}
	for _, p := range ports {
		if p.Role == policy.RoleTower || p.Role == policy.RoleAmplifier || p.Role == policy.RoleBackhaul || p.Role == policy.RoleUplink {
			if err := r.cnode.SetPortPoe(ctx, cnodePort.NodeID, p.Port, true, "site_on"); err != nil {
				return err
			}
		}
	}
	return nil
}
func (r *Reconciler) applyService(ctx context.Context, siteID, state string) error {
	ports, err := r.ports.GetBySite(siteID)
	if err != nil {
		return err
	}
	tower, err := policy.FindRole(ports, policy.RoleTower)
	if err != nil {
		return err
	}
	if tower.NodeID == "" {
		return fmt.Errorf("tower node_id missing")
	}
	return r.tower.SetService(ctx, tower.NodeID, state)
}
func (r *Reconciler) applyRadio(ctx context.Context, siteID, state string) error {
	ports, err := r.ports.GetBySite(siteID)
	if err != nil {
		return err
	}
	amp, err := policy.FindRole(ports, policy.RoleAmplifier)
	if err != nil {
		return err
	}
	if amp.NodeID == "" {
		return fmt.Errorf("amplifier node_id missing")
	}
	return r.amplifier.SetRadio(ctx, amp.NodeID, state)
}
