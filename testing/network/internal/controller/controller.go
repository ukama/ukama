package controller

import "github.com/ukama/ukama/testing/network/internal/db"

type Controller struct {
	repo db.VNodeRepo
}

func NewController(d db.VNodeRepo) *Controller {
	return &Controller{
		repo: d,
	}
}

func (c *Controller) ControllerInit() {

}
