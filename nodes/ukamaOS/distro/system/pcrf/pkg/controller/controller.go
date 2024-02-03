package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/store"
)

type Controller struct {
	store   *store.Store
	sw      string
	session string
	rc      *client.RemoteController
}

func NewController(db string, sw string, remote string) (*Controller, error) {
	c := &Controller{}
	store, err := store.NewStore(db)
	if err != nil {
		log.Errorf("Failed to create db: %v", err)
		return nil, err
	}

	rc, err := client.NewRemoteControllerClient(remote)
	if err != nil {
		log.Errorf("Failed to create client: %v", err)
		return nil, err
	}

	c.store = store
	c.sw = sw
	c.rc = rc

	return c, nil
}

func (c *Controller) validateSusbcriber(imsi string) error {
	/* Get subscriber policy by imsi*/
	s, err := c.store.GetSubscriber(imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber for %s:Error: %v", imsi, err)
		return err
	}

	/* store policy */
	p, err := c.store.GetPolicyByID(s.PolicyID.ID)
	if err != nil {
		log.Errorf("Failed to get subscriber policy %d:Error: %v", s.PolicyID, err)
		return err
	}

	now := time.Now().Unix()
	if p.StartTime > now && p.EndTime <= now {
		return fmt.Errorf("failed to get valid policy")
	}

	return nil
}

func (c *Controller) updateSubscriberPolicy(imsi string, p *store.Policy) error {

	err := c.store.CreatePolicy(p)
	if err != nil {
		return err
	}

	err = c.store.CreateSubscriber(imsi, p)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) CreateSession(ctx *gin.Context, req *CreateSession) error {
	/* validate subscriber*/
	err := c.validateSusbcriber(req.Imsi)
	if err != nil {
		/* Get subscriber policy from remote */
		p, err := c.rc.GetPolicy(req.Imsi)
		if err != nil {
			log.Errorf("Failed to get subscriber policy %d:Error: %v", s.PolicyID, err)
			return err
		}

		err = c.updateSubscriberPolicy(req.Imsi, p)
		if err != nil {
			log.Errorf("Failed to update subscriber %s with policy %d:Error: %v", req.Imsi, p.ID, err)
			return err
		}
	}

	/* setup bridge */
	err = c.setUpBridge()
	if err != nil {
		log.Errorf("Failed to setup bridge for subscriber %s:Error: %v", req.Imsi, err)
		return err
	}

	/* create session */
	session, err := c.store.CreateSession(req.Imsi, req.Ip)
	if err != nil {
		log.Errorf("Failed to create a session for subscriber %s:Error: %v", req.Imsi, err)
		return err
	}

	/* start monitoring session */
	err = c.session.StartMonitoring()
	if err != nil {
		log.Errorf("Failed to monitor session on bridge for subscriber %s:Error: %v", req.Imsi, err)
		return err
	}

	return nil
}

func (c *Controller) EndSession(ctx *gin.Context, req *EndSession) error {
	return nil
}

func (c *Controller) GetSessionByID(ctx *gin.Context, req *GetSessionByID) (*db.Session, error) {
	return nil, nil
}

func (c *Controller) GetCDRById(ctx *gin.Context, req *GetCDRById) (*CDR, error) {
	return nil, nil
}

func (c *Controller) GetCDRByImsi(ctx *gin.Context, req *GetCDRByImsi) (*CDR, error) {
	return nil, nil
}

func (c *Controller) GetPolicy(ctx *gin.Context, req *PolicyByImsi) (*Policy, error) {
	return nil, nil
}

func (c *Controller) AddPolicy(ctx *gin.Context, req *AddPolicyByImsi) error {
	return nil
}

func (c *Controller) RemovePolicy(ctx *gin.Context, req *PolicyByImsi) error {
	return nil
}

func (c *Controller) GetReroute(ctx *gin.Context) (*Reroute, error) {
	return nil, nil
}

func (c *Controller) UpdateReroute(ctx *gin.Context, req *Reroute) error {
	return nil
}
