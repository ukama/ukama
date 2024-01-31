package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller/store"
)

type Controller struct {
	store   *store.Store
	sw      string
	session string
}

func NewController(db string, sw string) (*Controller, error) {
	c := &Controller{}
	store, err := store.NewStore(db)
	if err != nil {
		log.Errorf("Failed to create db: %v", err)
		return nil, err
	}

	c.store = store
	c.sw = sw

	return c, nil
}

func (c *Controller) CreateSession(ctx *gin.Context, req *CreateSession) error {

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
