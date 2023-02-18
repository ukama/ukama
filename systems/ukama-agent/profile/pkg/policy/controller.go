package policy

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type PolicyController struct {
	Policy         []Policy
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	Org            string
	profileRepo    db.ProfileRepo
	nodePolicyPath string
	period         time.Duration
	pR             chan bool
}

func (p *PolicyController) InitPolicyController() {
	// This could be populated as apart of config
	p.Policy = []Policy{
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
	}
}

func NewPolicyController(pRepo db.ProfileRepo, org string, msgBus mb.MsgBusServiceClient, path string, monitor bool, period time.Duration) *PolicyController {
	p := &PolicyController{
		profileRepo:    pRepo,
		Org:            org,
		msgbus:         msgBus,
		nodePolicyPath: path,
		period:         period,
	}
	p.InitPolicyController()

	if msgBus != nil {
		p.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)
	}

	p.pR = make(chan bool)

	if monitor {
		p.StartPolicyRoutine()
	}

	return p
}

/*
For now all the policies are by default applicable for the profiles.
There might be more policies which are applicablee for certain profiles
that can be easily managed by adding policy db and adding applicable policy id for each susbcriber.
*/
func (p *PolicyController) RunPolicyControl(imsi string) (error, bool) {
	removed := false
	pf, err := p.profileRepo.GetByImsi(imsi)
	if err != nil {
		log.Errorf("failed to read profile for %s. Error %s", imsi, err.Error())
		return err, removed
	}

	for _, pl := range p.Policy {
		if pl.Check != nil {

			valid := pl.Check(*pf)
			if valid {
				continue
			}
			log.Infof("Policy Controller found profile %s failed to comply policy %s", pf.Imsi, pl.Name)
			/* if policy check failed, try the action */
			if pl.Action != nil {
				err, removed := pl.Action(p, *pf)
				if err != nil {
					log.Errorf("Error while checking policies: %s", err.Error())
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

func (p *PolicyController) syncProfile(method string, pf db.Profile) error {

	body, err := json.Marshal(pf)
	if err != nil {
		log.Errorf("error marshaling profile: %s", err.Error())
		return err
	}

	if p.msgbus != nil {
		route := p.baseRoutingKey.SetAction("node-feed").SetObject("policy").MustBuild()
		err = p.msgbus.PublishToNodeFeeder(route, "*", p.Org, p.nodePolicyPath, method, body)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", body, route, err.Error())
			return err
		}
	}

	return nil
}

func (p *PolicyController) publishEvent(action string, object string, msg protoreflect.ProtoMessage) error {
	var err error
	if p.msgbus != nil {
		route := p.baseRoutingKey.SetAction(action).SetObject(object).MustBuild()
		err = p.msgbus.PublishRequest(route, msg)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, route, err.Error())
			return err
		}
	}
	return err
}

func (p *PolicyController) StartPolicyRoutine() {
	log.Infof("Starting policy check routine with period %s.", p.period)
	p.monitor()
}

func (p *PolicyController) StopPolicyRoutine() {
	log.Infof("Stoping policy check routine with period %s.", p.period)
	p.pR <- true
}

func (p *PolicyController) doPolicyCheck() error {

	pf, err := p.profileRepo.List()
	log.Infof("Policy check routine started at %s for %d profiles.", time.Now().String(), len(pf))
	if err != nil {
		log.Errorf("Failed to list profiles: %s.", err.Error())
		return err
	}

	for _, profile := range pf {
		_, _ = p.RunPolicyControl(profile.Imsi)
	}
	log.Infof("Policy check routine ended at %s.", time.Now().String())
	return nil
}

func (p *PolicyController) monitor() {

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
