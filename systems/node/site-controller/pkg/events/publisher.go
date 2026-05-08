package events

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
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
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	route := fmt.Sprintf("event.cloud.local.%s.node.site-controller.%s", p.orgName, event)
	log.Infof("site-controller: event=%s payload=%s", event, string(body))
	return p.msgbus.PublishRequest(route, body)
}
