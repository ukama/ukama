/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package policy

import (
	"encoding/json"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/dataplan"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

type policyController struct {
	dp                   dataplan.PackageClient
	Rules                []Rule
	asrRepo              db.AsrRecordRepo
	period               time.Duration
	pR                   chan bool
	msgbus               mb.MsgBusServiceClient
	NodeFeederRoutingKey msgbus.RoutingKeyBuilder
	MsgBusRoutingKey     msgbus.RoutingKeyBuilder
	OrgName              string
	OrgId                string
	reroute              string
}

const (
	ADD    = "POST"
	UPDATE = "POST"
	DELETE = "DELETE"
)

type SimInfo struct {
	Imsi      string `path:"imsi" validate:"required" json:"-"`
	Iccid     string
	PackageId uuid.UUID
	NetworkId uuid.UUID
	Visitor   bool
	ID        uint
}

type SimPackageUpdate struct {
	Imsi      string `path:"imsi" validate:"required" json:"-"`
	PackageId uuid.UUID
}

type MsgPolicy struct {
	Uuid         string `json:"uuid"`
	Burst        uint64 `json:"burst"`
	TotalData    uint64 `json:"total_data"`
	ConsumedData uint64 `json:"consumed_data"`
	Dlbr         uint64 `json:"dlbr"`
	Ulbr         uint64 `json:"ulbr"`
	StartTime    uint64 `json:"start_time"`
	EndTime      uint64 `json:"end_time"`
}

type MsgSubscriber struct {
	Policy  MsgPolicy `json:"policy"`
	Reroute string    `json:"reroute"`
}

type Controller interface {
	InitPolicyController()
	NewPolicy(packageId uuid.UUID) (*db.Policy, error)
	SyncProfile(s *SimInfo, as *db.Asr, action string, object string, event bool) error
	RunPolicyControl(imsi string, event bool) (error, bool)
}

func NewPolicyController(asrRepo db.AsrRecordRepo, msgB mb.MsgBusServiceClient, dataplanHost string, orgName string, orgId string, reroute string, period time.Duration, monitor bool) *policyController {
	p := &policyController{
		dp:                   dataplan.NewPackageClient(dataplanHost),
		asrRepo:              asrRepo,
		msgbus:               msgB,
		NodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetRequestType().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName), //Need to have something same to other routes
		MsgBusRoutingKey:     msgbus.NewRoutingKeyBuilder().SetEventType().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		OrgName:              orgName,
		OrgId:                orgId,
		reroute:              reroute,
		period:               period,
	}
	p.InitPolicyController()

	p.pR = make(chan bool)

	if monitor {
		p.StartPolicyControllerRoutine()
	}

	return p
}

func createMessage(p *db.Policy, reroute string) *MsgSubscriber {
	return &MsgSubscriber{
		Policy: MsgPolicy{
			Uuid:         p.Id.String(),
			Burst:        p.Burst,
			Ulbr:         p.Ulbr,
			Dlbr:         p.Dlbr,
			TotalData:    p.TotalData,
			ConsumedData: p.ConsumedData,
			StartTime:    p.StartTime,
			EndTime:      p.EndTime,
		},
		Reroute: reroute,
	}
}

func (p *policyController) InitPolicyController() {
	// This could be populated as a part of config
	p.Rules = []Rule{
		{
			Name:   "DataCap",
			ID:     1,
			Check:  DataCapCheck,
			Action: RemoveProfile,
		},
		{
			Name:   "AllowedServiceTime",
			ID:     2,
			Check:  AllowedTimeOfServiceCheck,
			Action: RemoveProfile,
		},
		{
			Name:   "ValidityCheck",
			ID:     3,
			Check:  ValidityCheck,
			Action: RemoveProfile,
		},
	}
}

func (p *policyController) NewPolicy(packageId uuid.UUID) (*db.Policy, error) {
	log.Infof("Creating new policy based on package %s", packageId.String())

	pack, err := p.dp.Get(packageId.String())
	if err != nil {
		log.Errorf("Failed to get package %s.Error: %v", packageId.String(), err)

		return nil, err
	}

	st := uint64(time.Now().Unix())

	// st is in seconds and pack.Duration is in days
	et := uint64(st) + (pack.Duration * 24 * 3600)

	policy := db.Policy{
		Id:           uuid.NewV4(),
		Burst:        1500,
		TotalData:    pack.DataVolume,
		ConsumedData: 0,
		Dlbr:         pack.PackageDetails.Dlbr,
		Ulbr:         pack.PackageDetails.Ulbr,
		StartTime:    st,
		EndTime:      et,
	}

	return &policy, nil
}

func (p *policyController) SyncProfile(s *SimInfo, as *db.Asr, action string, object string, event bool) error {
	log.Infof("Syncing profile for subscriber %s based on action %s", as.Imsi, action)

	var httpMethod string

	subscriber := &epb.Subscriber{
		Imsi:       as.Imsi,
		Iccid:      as.Iccid,
		Network:    as.NetworkId.String(),
		Package:    as.PackageId.String(),
		SimPackage: as.SimPackageId.String(),
		Org:        p.OrgId,
		Policy:     as.Policy.Id.String(),
	}

	var msg protoreflect.ProtoMessage
	/* Create event */
	switch action {

	case "create":
		e := &epb.AsrActivated{
			Subscriber: subscriber,
		}
		httpMethod = "POST"
		msg = e
	case "delete":
		e := &epb.AsrInactivated{
			Subscriber: subscriber,
		}
		msg = e
		httpMethod = "DELETE"
	case "update":
		e := &epb.AsrUpdated{
			Subscriber: subscriber,
		}
		msg = e
		httpMethod = "PATCH"
	default:
		log.Errorf("invalid action %s to sync subscriber profile called.", action)
		return nil
	}

	err := p.syncSubscriberPolicy(httpMethod, s.Imsi, s.NetworkId.String(), &as.Policy)
	if err != nil {
		return err
	}

	if event {
		return p.publishEvent(action, object, msg)
	} else {
		return nil
	}
}

/*
For now all the policies are by default applicable for the profiles.
There might be more policies which are applicable for certain profiles
that can be easily managed by adding policy db and adding applicable policy id for each susbcriber.
*/
func (p *policyController) RunPolicyControl(imsi string, event bool) (error, bool) {
	log.Infof("Running policy control for subscriber %s", imsi)

	removed := false
	pf, err := p.asrRepo.GetByImsi(imsi)
	if err != nil {
		log.Errorf("failed to read profile for %s. Error %s", imsi, err.Error())
		return err, removed
	}

	for _, pt := range p.Rules {
		if pt.Check != nil {

			valid := pt.Check(*pf)
			if valid {
				continue
			}
			log.Infof("Policy Controller found profile %s has failed to comply policy type %s", pf.Imsi, pt.Name)
			/* if policy check failed, try the action */
			if pt.Action != nil {
				err, removed := pt.Action(p, *pf, event)
				if err != nil {
					log.Errorf("Error while applying action for failing policy compliance (%s, %s). Error: %v",
						pf.Imsi, pt.Name, err)

					return err, removed
				}

				/* if profile is removed then just stop checking polices further*/
				if removed {
					break
				}
			}
		}
	}

	return nil, removed
}

func (p *policyController) syncSubscriberPolicy(method string, imsi string, network string, policy *db.Policy) error {
	log.Infof("Syncing policy for subscriber %s", imsi)

	route := p.NodeFeederRoutingKey.SetObject("nodefeeder").SetAction("publish").MustBuild()
	pMsg := createMessage(policy, p.reroute)

	jd, err := json.Marshal(pMsg)
	if err != nil {
		log.Errorf("Failed to marshal policy %+v for subscriber %s. Errors %s", pMsg, imsi, err.Error())
		return err
	}

	path := "/pcrf/v1/subscriber/imsi/" + imsi

	msg := &pb.NodeFeederMessage{
		Target:     p.OrgName + "." + network + "." + "*" + "." + "*",
		HTTPMethod: method,
		Path:       path,
		Msg:        jd,
	}

	err = p.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", pMsg, route, err.Error())
		return err
	}

	log.Infof("Published Policy %s  for imsi %s on route %s.", msg, imsi, route)
	return nil
}

func (p *policyController) publishEvent(action string, object string, msg protoreflect.ProtoMessage) error {
	var err error
	if p.msgbus != nil {
		route := p.MsgBusRoutingKey.SetObject(object).SetAction(action).MustBuild()
		err = p.msgbus.PublishRequest(route, msg)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, route, err.Error())
			return err
		}
	}
	return err
}

func (p *policyController) StartPolicyControllerRoutine() {
	log.Infof("Starting policy check routine with period %s.", p.period)
	p.monitor()
}

func (p *policyController) StopPolicyControllerRoutine() {
	log.Infof("Stoping policy check routine with period %s.", p.period)
	p.pR <- true
}

func (p *policyController) doPolicyCheck() error {

	pf, err := p.asrRepo.List()
	log.Infof("Policy check routine started at %s for %d profiles.", time.Now().String(), len(pf))
	if err != nil {
		log.Errorf("Failed to list profiles: %s.", err.Error())
		return err
	}

	for _, profile := range pf {
		_, _ = p.RunPolicyControl(profile.Imsi, true)
	}
	log.Infof("Policy check routine ended at %s.", time.Now().String())
	return nil
}

func (p *policyController) monitor() {

	t := time.NewTicker(p.period)

	go func() {
		for {
			select {
			case <-t.C:
				_ = p.doPolicyCheck()
			case <-p.pR:
				t.Stop()
				return
			}
		}
	}()
}
