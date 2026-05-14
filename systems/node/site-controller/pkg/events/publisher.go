/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package events

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"google.golang.org/protobuf/types/known/structpb"
)

type Publisher interface {
	IntentChanged(siteId string, target string, state string, reason string) error
	StateChanged(siteId string, state *db.SiteState) error
	SwitchPolicyApplied(siteId string) error
	ReconcileFailed(siteId string, reason string) error
}

type MsgBusPublisher struct {
	orgName string
	msgbus  mb.MsgBusServiceClient
}

func NewMsgBusPublisher(orgName string, msgBus mb.MsgBusServiceClient) *MsgBusPublisher {
	return &MsgBusPublisher{orgName: orgName, msgbus: msgBus}
}

func (p *MsgBusPublisher) IntentChanged(siteId string, target string, state string, reason string) error {
	return p.publish(EventSiteIntentChanged, map[string]string{
		"site_id": siteId,
		"target":  target,
		"state":   state,
		"reason":  reason,
	})
}

func (p *MsgBusPublisher) StateChanged(siteId string, state *db.SiteState) error {
	return p.publish(EventSiteStateChanged, map[string]interface{}{
		"site_id": siteId,
		"power":   state.PowerState,
		"service": state.ServiceState,
		"radio":   state.RadioState,
		"access":  state.AccessState,
		"reason":  state.Reason,
	})
}

func (p *MsgBusPublisher) SwitchPolicyApplied(siteId string) error {
	return p.publish(EventSiteSwitchPolicyApplied, map[string]string{"site_id": siteId})
}

func (p *MsgBusPublisher) ReconcileFailed(siteId string, reason string) error {
	return p.publish(EventSiteReconcileFailed, map[string]string{
		"site_id": siteId,
		"reason":  reason,
	})
}

func (p *MsgBusPublisher) publish(event string, payload interface{}) error {
	if p == nil || p.msgbus == nil {
		return nil
	}
	msg, err := payloadToStruct(payload)
	if err != nil {
		return err
	}
	route := fmt.Sprintf("event.cloud.local.%s.node.site-controller.%s", p.orgName, event)
	log.Infof("site-controller: event=%s payload=%v", event, msg)
	return p.msgbus.PublishRequest(route, msg)
}

func payloadToStruct(payload interface{}) (*structpb.Struct, error) {
	switch p := payload.(type) {
	case map[string]string:
		m := make(map[string]interface{}, len(p))
		for k, v := range p {
			m[k] = v
		}
		return structpb.NewStruct(m)
	case map[string]interface{}:
		return structpb.NewStruct(p)
	default:
		return nil, fmt.Errorf("unsupported payload type")
	}
}
