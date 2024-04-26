package pcrf

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
)

type policyFunction struct {
	r                    db.PolicyRepo
	msgbus               mb.MsgBusServiceClient
	NodeFeederRoutingKey msgbus.RoutingKeyBuilder
	OrgName              string
	reroute              string
}

type MsgPoilcy struct {
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
	Policy  MsgPoilcy `json:"policy"`
	Reroute string    `json:"reroute"`
}

type PolicyFunctionController interface {
	GetPolicy(id uuid.UUID) (*db.Policy, error)
	CreatePolicy(p *db.Policy) error
	DeletePolicy(id uuid.UUID) error
	DeletePolicyByAsrID(id uint) error
	UpdatePolicy(id uint, p *db.Policy) error
	ApplyPolicy(method string, imsi string, network string, p *db.Policy) error
}

func NewPolicyFunctionController(msgB mb.MsgBusServiceClient, db db.PolicyRepo, orgName string, reroute string) *policyFunction {
	return &policyFunction{
		r:                    db,
		msgbus:               msgB,
		NodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetRequestType().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName), //Need to have something same to other routes
		OrgName:              orgName,
		reroute:              reroute,
	}
}

func createMessage(p *db.Policy, reroute string) *MsgSubscriber {

	return &MsgSubscriber{
		Policy: MsgPoilcy{
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

func (pf *policyFunction) GetPolicy(id uuid.UUID) (*db.Policy, error) {
	policy, err := pf.r.Get(id)
	if err != nil {
		log.Errorf("Error creating policy %v.Error: %v", policy, err)
		return nil, err
	}

	return policy, nil
}

func (pf *policyFunction) CreatePolicy(p *db.Policy) error {
	err := pf.r.Add(p)
	if err != nil {
		log.Errorf("Error creating policy %v.Error: %v", p, err)
		return err
	}
	return nil
}

func (pf *policyFunction) DeletePolicy(id uuid.UUID) error {
	err := pf.r.Delete(id)
	if err != nil {
		log.Errorf("Error deleting policy %s.Error: %v", id.String(), err)
		return err
	}
	return nil
}

func (pf *policyFunction) DeletePolicyByAsrID(id uint) error {

	policy, err := pf.r.GetByAsrId(id)
	if err != nil {
		log.Errorf("Error deleting policy %v for ASR Id %d.Error: %v", policy, id, err)
		return err
	}

	err = pf.r.Delete(policy.Id)
	if err != nil {
		log.Errorf("Error deleting policy %s.Error: %v", policy.Id.String(), err)
		return err
	}
	return nil
}

func (pf *policyFunction) UpdatePolicy(id uint, p *db.Policy) error {
	err := pf.r.Update(id, p)
	if err != nil {
		log.Errorf("Error deleting policy for ASR id %d.Error: %v", id, err)
		return err
	}

	return nil
}

func (pf *policyFunction) MonitorPolicy() error {

	//TODO: May be have repo gorm.default, imsi, usage
	/*
		=>	This will be a periodic routine

		=>	Update usage from the event sent by CDR service on reciving a new CDR report for imsi
			Compare the usage to the policy data limit in periodin oand on events
			As soons as it hits the cap remove the subscriber

		=>	Important: need to make sure when the updates comes from the diffrent nodes we generate a new policy and update the max data limit
		 	available to subscriber and push them to all nodes

		=> may be CDR also contains the nodeId from which it was genrated
		   This might help us to resolve lot of issues like quick update , figure out if the user is moving , roaming, locations etc
	*/
	return nil
}

func (pf *policyFunction) ApplyPolicy(method string, imsi string, network string, p *db.Policy) error {

	route := pf.NodeFeederRoutingKey.SetObject("node").SetAction("publish").MustBuild()
	pMsg := createMessage(p, pf.reroute)

	jd, err := json.Marshal(pMsg)
	if err != nil {
		log.Errorf("Failed to marshal policy %+v for subscriber %s. Errors %s", pMsg, imsi, err.Error())
		return err
	}

	path := "/pcrf/v1/subscriber/imsi/" + imsi

	msg := &pb.NodeFeederMessage{
		Target:     pf.OrgName + "." + network + "." + "*" + "." + "*",
		HTTPMethod: method,
		Path:       path,
		Msg:        jd,
	}

	err = pf.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", pMsg, route, err.Error())
		return err
	}

	log.Infof("Published Policy %s  for imsi %s on route %s.", msg, imsi, route)
	return nil
}
