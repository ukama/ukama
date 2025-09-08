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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/client"

	log "github.com/sirupsen/logrus"
	subRegPb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	simMangPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	simPoolPb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
)

const SUBS_URL_PARAMETER = "subscriber"

var REDIRECT_URI = "https://subscriber.dev.ukama.com/swagger/#/"

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
	auth          *config.Auth
}

type Clients struct {
	sp  simPool
	sm  simManager
	sub subscriber
}

type simPool interface {
	Get(iccid string) (*simPoolPb.GetByIccidResponse, error)
	GetStats(simType string) (*simPoolPb.GetStatsResponse, error)
	GetSims(simType string) (*simPoolPb.GetSimsResponse, error)
	AddSimsToSimPool(req *simPoolPb.AddRequest) (*simPoolPb.AddResponse, error)
	UploadSimsToSimPool(req *simPoolPb.UploadRequest) (*simPoolPb.UploadResponse, error)
	DeleteSimFromSimPool(id []uint64) (*simPoolPb.DeleteResponse, error)
}

type simManager interface {
	AllocateSim(req *simMangPb.AllocateSimRequest) (*simMangPb.AllocateSimResponse, error)
	GetSim(simId string) (*simMangPb.GetSimResponse, error)
	ListSims(iccid, imsi, subscriberId, networkId, simType, simStatus string, trafficPolicy uint32,
		isPhysical, sort bool, count uint32) (*simMangPb.ListSimsResponse, error)
	ToggleSimStatus(simId string, status string) (*simMangPb.ToggleSimStatusResponse, error)
	AddPackageToSim(req *simMangPb.AddPackageRequest) (*simMangPb.AddPackageResponse, error)
	RemovePackageForSim(req *simMangPb.RemovePackageRequest) (*simMangPb.RemovePackageResponse, error)
	TerminateSim(simId string) (*simMangPb.TerminateSimResponse, error)
	ListPackagesForSim(simId, dataPlanId, fromStartDate, toStartDate, fromEndDate,
		toEndDate string, isActive, asExpired, sort bool, count uint32) (*simMangPb.ListPackagesForSimResponse, error)
	SetActivePackageForSim(req *simMangPb.SetActivePackageRequest) (*simMangPb.SetActivePackageResponse, error)
	GetUsages(iccid, simType, cdrType, from, to, region string) (*simMangPb.UsageResponse, error)

	// Deprecated: Use pkg.client.SimManager.ListSims with subscriberId as filtering param instead.
	GetSimsBySub(subscriberId string) (*simMangPb.GetSimsBySubscriberResponse, error)
	// Deprecated: Use pkg.client.SimManager.ListSims with networkId as filtering param instead.
	GetSimsByNetwork(networkId string) (*simMangPb.GetSimsByNetworkResponse, error)
	// Deprecated: Use pkg.client.SimManager.ListPackagesForSim with simId as filtering param instead.
	GetPackagesForSim(simId string) (*simMangPb.GetPackagesForSimResponse, error)
}

type subscriber interface {
	GetSubscriber(sid string) (*subRegPb.GetSubscriberResponse, error)
	GetSubscriberByEmail(sEmail string) (*subRegPb.GetSubscriberByEmailResponse, error)
	AddSubscriber(req *subRegPb.AddSubscriberRequest) (*subRegPb.AddSubscriberResponse, error)
	DeleteSubscriber(sid string) (*subRegPb.DeleteSubscriberResponse, error)
	UpdateSubscriber(subscriber *subRegPb.UpdateSubscriberRequest) (*subRegPb.UpdateSubscriberResponse, error)
	GetByNetwork(networkId string) (*subRegPb.GetByNetworkResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.sp = client.NewSimPool(endpoints.SimPool, endpoints.Timeout)
	c.sm = client.NewSimManager(endpoints.SimManager, endpoints.Timeout)
	c.sub = client.NewRegistry(endpoints.Registry, endpoints.Timeout)
	return c
}

func NewRouter(clients *Clients, config *RouterConfig, authfunc func(*gin.Context, string) error) *Router {

	r := &Router{
		clients: clients,
		config:  config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init(authfunc)

	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		httpEndpoints: &svcConf.HttpServices,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
		auth:          svcConf.Auth,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	auth := r.f.Group("/v1", "Subscriber API GW ", "Subs system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")
			return
		}
		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
		err := f(ctx, r.config.auth.AuthServerUrl)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
	})
	auth.Use()
	{
		/* These two API will be available based on RBAC */
		auth.GET("/subscribers/networks/:network_id", formatDoc("List all subscribers for a Network", ""), tonic.Handler(r.getSubscriberByNetwork, http.StatusOK))
		auth.GET("/sims/networks/:network_id", formatDoc("List all sims for a Network", ""), tonic.Handler(r.getSimsByNetwork, http.StatusOK))

		pool := auth.Group("/simpool", "SIM Pool", "SIM store for Org")
		pool.GET("/sim/:iccid", formatDoc("Get SIM by Iccid", ""), tonic.Handler(r.getSimByIccid, http.StatusOK))
		pool.GET("/sims/:sim_type", formatDoc("Get SIMs by type", ""), tonic.Handler(r.getSims, http.StatusOK))
		pool.GET("/stats/:sim_type", formatDoc("Get SIM Pool stats", ""), tonic.Handler(r.getSimPoolStats, http.StatusOK))
		pool.PUT("", formatDoc("Add new SIM to SIM pool", ""), tonic.Handler(r.addSimsToSimPool, http.StatusCreated))
		pool.PUT("/upload", formatDoc("Upload CSV file to add new sim to SIM Pool", ""), tonic.Handler(r.uploadSimsToSimPool, http.StatusCreated))
		pool.DELETE("/sim/:sim_id", formatDoc("Remove SIM from SIM Pool", ""), tonic.Handler(r.deleteSimFromSimPool, http.StatusOK))

		subscriber := auth.Group("/subscriber", "Subscriber", "Orgs Subscriber database")
		subscriber.GET("/:subscriber_id", formatDoc("Get subscriber by id", ""), tonic.Handler(r.getSubscriber, http.StatusOK))
		subscriber.GET("/email/:subscriber_email", formatDoc("Get subscriber by email", ""), tonic.Handler(r.getSubscriberByEmail, http.StatusOK))
		subscriber.PUT("", formatDoc("Add a new subscriber", ""), tonic.Handler(r.putSubscriber, http.StatusCreated))
		subscriber.DELETE("/:subscriber_id", formatDoc("Delete a subscriber", ""), tonic.Handler(r.deleteSubscriber, http.StatusOK))
		subscriber.PATCH("/:subscriber_id", formatDoc("Update a subscriber", ""), tonic.Handler(r.updateSubscriber, http.StatusOK))

		sim := auth.Group("/sim", "SIM", "Orgs SIM data base")
		sim.GET("", formatDoc("List SIMs with various query params as filters", ""), tonic.Handler(r.listSims, http.StatusOK))
		sim.GET("/:sim_id", formatDoc("Get SIM by Id", ""), tonic.Handler(r.getSim, http.StatusOK))
		sim.POST("/", formatDoc("Allocate a new SIM to given subscriber", ""), tonic.Handler(r.allocateSim, http.StatusCreated))
		sim.PATCH("/:sim_id", formatDoc("Activate/Deactivate a given SIM", ""), tonic.Handler(r.updateSimStatus, http.StatusOK))
		sim.DELETE("/:sim_id", formatDoc("Terminate a given SIM", ""), tonic.Handler(r.terminateSim, http.StatusOK))
		sim.GET("/:sim_id/package", formatDoc("Get packages for a given SIM", ""), tonic.Handler(r.listPackagesForSim, http.StatusOK))
		sim.POST("/:sim_id/package", formatDoc("Add a new package to the given SIM", ""), tonic.Handler(r.addPackageForSim, http.StatusCreated))
		sim.PATCH("/:sim_id/package/:package_id", formatDoc("Set active package for a given SIM", ""), tonic.Handler(r.setActivePackageForSim, http.StatusOK))
		sim.DELETE("/:sim_id/package/:package_id", formatDoc("Delete a package from a given SIM", ""), tonic.Handler(r.removePkgForSim, http.StatusOK))
		// Deprecated: Use GET /v1/sim with subscriberId as query param instead.
		sim.GET("/subscriber/:subscriber_id", formatDoc("Get the list of SIMs for a given subscriber", ""), tonic.Handler(r.getSimsBySub, http.StatusOK))
		// Deprecated: Use GET /v1/sim/:sim_id/package with query params  for filtering instead.
		sim.GET("/packages/:sim_id", formatDoc("Get packages for a given SIM", ""), tonic.Handler(r.getPackagesForSim, http.StatusOK))
		// Deprecated: Use POST /v1/sim/:sim_id/package instead.
		sim.POST("/package", formatDoc("Add a new package to the given subscriber's SIM", ""), tonic.Handler(r.postPkgForSim, http.StatusCreated))

		usage := auth.Group("usages", "Usages", "Operator sims usages endpoints")
		usage.GET("", formatDoc("Get Usages", "Get sim usages with filters"), tonic.Handler(r.getUsages, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getSimByIccid(c *gin.Context, req *SimByIccidReq) (*simPoolPb.GetByIccidResponse, error) {
	return r.clients.sp.Get(req.Iccid)
}

func (r *Router) getSims(c *gin.Context, req *SimPoolStatReq) (*simPoolPb.GetSimsResponse, error) {
	return r.clients.sp.GetSims(req.SimType)
}

func (r *Router) getSimPoolStats(c *gin.Context, req *SimPoolTypeReq) (*simPoolPb.GetStatsResponse, error) {
	return r.clients.sp.GetStats(req.SimType)
}

func (r *Router) addSimsToSimPool(c *gin.Context, req *SimPoolAddSimReq) (*simPoolPb.AddResponse, error) {
	pbreq, err := addReqToAddSimReqPb(req)
	if err != nil {
		return nil, err
	}
	return r.clients.sp.AddSimsToSimPool(pbreq)
}

func (r *Router) uploadSimsToSimPool(c *gin.Context, req *SimPoolUploadSimReq) (*simPoolPb.UploadResponse, error) {

	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: fmt.Sprintf("failed to decode base64 data: %v", err)}
	}

	return r.clients.sp.UploadSimsToSimPool(&simPoolPb.UploadRequest{
		SimData: data,
		SimType: req.SimType,
	})
}

func (r *Router) deleteSimFromSimPool(c *gin.Context, req *SimPoolRemoveSimReq) (*simPoolPb.DeleteResponse, error) {
	return r.clients.sp.DeleteSimFromSimPool(req.Id)
}

func (r *Router) getSubscriberByEmail(c *gin.Context, req *SubscriberGetReqByEmail) (*subRegPb.GetSubscriberByEmailResponse, error) {
	return r.clients.sub.GetSubscriberByEmail(strings.ToLower(req.Email))
}

func (r *Router) getSubscriber(c *gin.Context, req *SubscriberGetReq) (*subRegPb.GetSubscriberResponse, error) {
	subsId := req.SubscriberId

	return r.clients.sub.GetSubscriber(subsId)
}

func (r *Router) putSubscriber(c *gin.Context, req *SubscriberAddReq) (*subRegPb.AddSubscriberResponse, error) {

	return r.clients.sub.AddSubscriber(&subRegPb.AddSubscriberRequest{
		Name:                  req.Name,
		Email:                 strings.ToLower(req.Email),
		PhoneNumber:           req.Phone,
		Dob:                   req.Dob,
		Address:               req.Address,
		ProofOfIdentification: req.ProofOfIdentification,
		IdSerial:              req.IdSerial,
		NetworkId:             req.NetworkId,
		Gender:                req.Gender,
	})
}

func (r *Router) deleteSubscriber(c *gin.Context, req *SubscriberDeleteReq) (*subRegPb.DeleteSubscriberResponse, error) {
	return r.clients.sub.DeleteSubscriber(req.SubscriberId)
}

func (r *Router) updateSubscriber(c *gin.Context, req *SubscriberUpdateReq) (*subRegPb.UpdateSubscriberResponse, error) {
	return r.clients.sub.UpdateSubscriber(&subRegPb.UpdateSubscriberRequest{
		SubscriberId:          req.SubscriberId,
		Name:                  req.Name,
		PhoneNumber:           req.Phone,
		Address:               req.Address,
		ProofOfIdentification: req.ProofOfIdentification,
		IdSerial:              req.IdSerial,
	})
}

func (r *Router) getSubscriberByNetwork(c *gin.Context, req *SubscriberByNetworkReq) (*subRegPb.GetByNetworkResponse, error) {
	return r.clients.sub.GetByNetwork(req.NetworkId)
}

func (r *Router) allocateSim(c *gin.Context, req *AllocateSimReq) (*simMangPb.AllocateSimResponse, error) {
	simReq := simMangPb.AllocateSimRequest{
		SubscriberId:  req.SubscriberId,
		SimToken:      req.SimToken,
		PackageId:     req.PackageId,
		NetworkId:     req.NetworkId,
		SimType:       req.SimType,
		TrafficPolicy: req.TrafficPolicy,
	}
	return r.clients.sm.AllocateSim(&simReq)
}

func (r *Router) getSim(c *gin.Context, req *SimReq) (*simMangPb.GetSimResponse, error) {
	return r.clients.sm.GetSim(req.SimId)
}

func (r *Router) listSims(c *gin.Context, req *ListSimsReq) (*simMangPb.ListSimsResponse, error) {
	return r.clients.sm.ListSims(req.Iccid, req.Imsi, req.SubscriberId, req.NetworkId,
		req.SimType, req.SimStatus, req.TrafficPolicy, req.IsPhysical, req.Sort, req.Count)
}

// Deprecated: Use pkg.rest.Router.ListSims with subscriberId as filtering param instead.
func (r *Router) getSimsBySub(c *gin.Context, req *GetSimsBySubReq) (*simMangPb.GetSimsBySubscriberResponse, error) {
	return r.clients.sm.GetSimsBySub(req.SubscriberId)
}

// Deprecated: Use pkg.rest.Router.ListSims with networkId as filtering param instead.
func (r *Router) getSimsByNetwork(c *gin.Context, req *SimByNetworkReq) (*simMangPb.GetSimsByNetworkResponse, error) {
	return r.clients.sm.GetSimsByNetwork(req.NetworkId)
}

func (r *Router) updateSimStatus(c *gin.Context, req *ActivateDeactivateSimReq) (*simMangPb.ToggleSimStatusResponse, error) {
	return r.clients.sm.ToggleSimStatus(req.SimId, req.Status)
}

func (r *Router) terminateSim(c *gin.Context, req *SimReq) (*simMangPb.TerminateSimResponse, error) {
	return r.clients.sm.TerminateSim(req.SimId)
}

func (r *Router) addPackageForSim(c *gin.Context, req *AddPkgToSimReq) (*simMangPb.AddPackageResponse, error) {
	payload := simMangPb.AddPackageRequest{
		SimId:     req.SimId,
		PackageId: req.PackageId,
		StartDate: req.StartDate,
	}
	return r.clients.sm.AddPackageToSim(&payload)
}

// Deprecated: Use pkg.rest.Router.addPkgForSim instead.
func (r *Router) postPkgForSim(c *gin.Context, req *PostPkgToSimReq) error {
	payload := simMangPb.AddPackageRequest{
		SimId:     req.SimId,
		PackageId: req.PackageId,
		StartDate: req.StartDate,
	}
	_, err := r.clients.sm.AddPackageToSim(&payload)
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) listPackagesForSim(c *gin.Context, req *ListPackagesForSimReq) (*simMangPb.ListPackagesForSimResponse, error) {
	return r.clients.sm.ListPackagesForSim(req.SimId, req.DataPlanId, req.FromStartDate, req.ToStartDate, req.FromEndDate,
		req.ToEndDate, req.IsActive, req.AsExpired, req.Sort, req.Count)
}

// Deprecated: Use pkg.rest.Router.listPackagesForSim instead.
func (r *Router) getPackagesForSim(c *gin.Context, req *SimReq) (*simMangPb.GetPackagesForSimResponse, error) {
	return r.clients.sm.GetPackagesForSim(req.SimId)
}

func (r *Router) setActivePackageForSim(c *gin.Context, req *SetActivePackageForSimReq) (*simMangPb.SetActivePackageResponse, error) {
	payload := simMangPb.SetActivePackageRequest{
		SimId:     req.SimId,
		PackageId: req.PackageId,
	}

	return r.clients.sm.SetActivePackageForSim(&payload)
}

func (r *Router) removePkgForSim(c *gin.Context, req *RemovePkgFromSimReq) (*simMangPb.RemovePackageResponse, error) {
	payload := simMangPb.RemovePackageRequest{
		SimId:     req.SimId,
		PackageId: req.PackageId,
	}

	return r.clients.sm.RemovePackageForSim(&payload)
}

func (r *Router) getUsages(c *gin.Context, req *GetUsagesReq) (*simMangPb.UsageResponse, error) {
	cdrType, ok := c.GetQuery("cdr_type")
	if !ok {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "cdr_type is a mandatory query parameter"}
	}

	return r.clients.sm.GetUsages(req.SimId, req.SimType, cdrType, req.From, req.To, req.Region)

}

func addReqToAddSimReqPb(req *SimPoolAddSimReq) (*simPoolPb.AddRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid add request")
	}

	list := make([]*simPoolPb.AddSim, len(req.SimInfo))
	for i, iter := range req.SimInfo {
		list[i] = &simPoolPb.AddSim{
			Iccid:          iter.Iccid,
			Msisdn:         iter.Msisdn,
			ActivationCode: iter.ActivationCode,
			IsPhysical:     iter.IsPhysicalSim,
			QrCode:         iter.QrCode,
			SmDpAddress:    iter.SmDpAddress,
			SimType:        iter.SimType,
		}
	}
	pbReq := &simPoolPb.AddRequest{
		Sim: list,
	}

	return pbReq, nil
}
