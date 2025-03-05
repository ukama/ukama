/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/testing/services/dummy/api-gateway/cmd/version"
	"github.com/ukama/ukama/testing/services/dummy/api-gateway/pkg"
	"github.com/ukama/ukama/testing/services/dummy/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	spb "github.com/ukama/ukama/testing/services/dummy/dsimfactory/pb/gen"
	pb "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
	logger  *logrus.Logger
}

type RouterConfig struct {
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
}

type Clients struct {
	Dsubscriber dsubscriber
	Dsimfactory dsimfactory
}

type dsubscriber interface {
	Update(req *pb.UpdateRequest) (*pb.UpdateResponse, error)
}

type dsimfactory interface {
	GetSims() (*spb.GetSimsResponse, error)
	GetSim(iccid string) (*spb.GetByIccidResponse, error)
	UploadSimsToSimPool(req *spb.UploadRequest) (*spb.UploadResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Dsubscriber = client.NewDsubscriber(endpoints.Dsubscriber, endpoints.Timeout)
	c.Dsimfactory = client.NewDsimfactory(endpoints.Dsimfactory, endpoints.Timeout)
	return c
}

func NewRouter(clients *Clients, config *RouterConfig) *Router {
	r := &Router{
		clients: clients,
		config:  config,
		logger:  log.New(),
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		httpEndpoints: &svcConf.HttpServices,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, "")
	endpoint := r.f.Group("/v1", "API gateway", "Dummy system version v1")
	endpoint.GET("/ping", formatDoc("Ping the server", "Returns a response indicating that the server is running."), tonic.Handler(r.pingHandler, http.StatusOK))

	dsub := endpoint.Group("/dsubscriber", "Dsubscriber", "Dummy subscriber service")
	dsub.PUT("/update", formatDoc("Update dsubscriber coroutine", "Update dsubscriber coroutine for specific subscriber."), tonic.Handler(r.updateHandler, http.StatusCreated))

	dsim := endpoint.Group("/factory", "Dsimfactory", "Dummy sim factory")
	dsim.GET("/simcards", formatDoc("Get SIMs", ""), tonic.Handler(r.getSims, http.StatusOK))
	dsim.GET("/simcards/:iccid", formatDoc("Get SIM by ICCID", ""), tonic.Handler(r.getSim, http.StatusOK))
	dsim.PUT("/simcards/upload", formatDoc("Upload CSV file to add new sim to SIM Pool", ""), tonic.Handler(r.uploadSimsToSimPool, http.StatusCreated))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) pingHandler(c *gin.Context) error {
	response := make(map[string]string)
	response["status"] = pkg.SystemName + " is running"
	response["version"] = version.Version
	c.JSON(http.StatusOK, response)

	return nil
}

func (r *Router) updateHandler(c *gin.Context, req *UpdateReq) (*pb.UpdateResponse, error) {
	return r.clients.Dsubscriber.Update(&pb.UpdateRequest{
		Dsubscriber: &pb.Dsubscriber{
			Iccid:   req.Iccid,
			Profile: req.Profile,
		}})
}

func (r *Router) getSims(c *gin.Context, req *GetSims) (*spb.GetSimsResponse, error) {
	resp, err := r.clients.Dsimfactory.GetSims()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Router) getSim(c *gin.Context, req *GetSimByIccid) (*spb.GetByIccidResponse, error) {
	resp, err := r.clients.Dsimfactory.GetSim(req.Iccid)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Router) uploadSimsToSimPool(c *gin.Context, req *SimPoolUploadSimReq) (*spb.UploadResponse, error) {

	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		log.Fatal("error:", err)
	}

	pbResp, err := r.clients.Dsimfactory.UploadSimsToSimPool(&spb.UploadRequest{
		SimData: data,
	})

	if err != nil {
		return nil, err
	}

	return pbResp, nil
}
