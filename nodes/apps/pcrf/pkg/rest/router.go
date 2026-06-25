/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/nodes/apps/pcrf/cmd/version"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller"
	"github.com/ukama/ukama/systems/common/config"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
)

var REDIRECT_URI = "https://0.0.0.0:8080/swagger/#/"

type RuntimeStatus struct {
	InitNetworkURL   string `json:"url"`
	InitNetworkReady bool   `json:"ready"`
	UECidr           string `json:"ueCidr"`
	DBPath           string `json:"dbPath"`
}

type Router struct {
	f          *fizz.Fizz
	controller *controller.Controller
	config     *RouterConfig
	nodeId     string
	runtime    RuntimeStatus
}

type RouterConfig struct {
	debugMode  bool
	serverConf *crest.HttpConfig
	auth       *config.Auth
}

func NewRouter(ctr *controller.Controller, config *RouterConfig, nodeId string, runtime RuntimeStatus) *Router {
	r := &Router{
		controller: ctr,
		config:     config,
		nodeId:     nodeId,
		runtime:    runtime,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()

	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
		auth:       svcConf.Auth,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)

	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) ping(c *gin.Context) {
	status := r.controller.Status()
	if !status.Ready {
		c.String(http.StatusServiceUnavailable, "not ready")
		return
	}

	c.String(http.StatusOK, "OK")
}

func (r *Router) version(c *gin.Context) {
	c.String(http.StatusOK, version.Version)
}

func (r *Router) status(c *gin.Context) {
	ctrStatus := r.controller.Status()

	c.JSON(http.StatusOK, gin.H{
		"ready":   ctrStatus.Ready,
		"state":   ctrStatus.State,
		"reason":  ctrStatus.Reason,
		"service": ctrStatus.Service,
		"nodeId":  r.nodeId,
		"initNetwork": gin.H{
			"url":   r.runtime.InitNetworkURL,
			"ready": r.runtime.InitNetworkReady,
		},
		"datapath": ctrStatus.DataPath,
		"ue": gin.H{
			"cidr": r.runtime.UECidr,
		},
		"db": gin.H{
			"path": r.runtime.DBPath,
		},
		"sessions": ctrStatus.Sessions,
	})
}

func (r *Router) init() {
	r.f = crest.NewFizzRouter(r.config.serverConf, pkg.SystemName,
		version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")

	r.f.Engine().GET("/v1/ping", r.ping)
	r.f.Engine().GET("/v1/version", r.version)
	r.f.Engine().GET("/v1/status", r.status)

	auth := r.f.Group("", "PCRF for Node", "API system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")
			return
		}

		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
	})

	auth.Use()
	{
		v1 := auth.Group("/v1", "PCRF", "pcrf")

		svc := v1.Group("/service", "Service", "service")
		svc.GET("", formatDoc("Get service state", "Get service admission state"), tonic.Handler(r.getService, http.StatusOK))
		svc.POST("", formatDoc("Set service state", "Enable or disable service admission"), tonic.Handler(r.setService, http.StatusAccepted))

		s := v1.Group("/session", "session", "session")
		s.POST("", formatDoc("Create session", "Create a new session"), tonic.Handler(r.createSession, http.StatusAccepted))
		s.DELETE("", formatDoc("End session", "End a session"), tonic.Handler(r.endSession, http.StatusAccepted))
		s.GET("/:id", formatDoc("Get session", "Get a particular session"), tonic.Handler(r.getSessionByID, http.StatusOK))

		cdr := v1.Group("/cdr", "CDR", "cdr")
		cdr.GET("/session/:id", formatDoc("Get a CDR", "Get CDR for specific session"), tonic.Handler(r.getCDRBySessionId, http.StatusOK))
		cdr.GET("/imsi/:imsi", formatDoc("Get CDRs", "Get CDRs for imsi"), tonic.Handler(r.getCDRByImsi, http.StatusOK))

		policy := v1.Group("/policy", "Policy", "policy")
		policy.POST("", formatDoc("Configure Policy", "Configure a new policy"), tonic.Handler(r.addPolicy, http.StatusCreated))
		policy.GET("/imsi/:imsi", formatDoc("Get policy", "Get a policy for a specific sim"), tonic.Handler(r.getPolicyByImsi, http.StatusOK))
		policy.GET("/id/:id", formatDoc("Get policy", "Get a policy by specific ID"), tonic.Handler(r.getPolicyByID, http.StatusOK))

		route := v1.Group("/reroute", "Reroute", "rerouting IP address")
		route.GET("/imsi/:imsi", formatDoc("reroute address", "Get a rerouting IP Address for Imsi"), tonic.Handler(r.getRerouteForImsi, http.StatusOK))
		route.POST("/id/:id", formatDoc("Add or Update Node", "Update a rerouting IP Address"), tonic.Handler(r.updateReroute, http.StatusCreated))

		sub := v1.Group("/subscriber", "Subscriber", "subscriber")
		sub.GET("/imsi/:imsi", formatDoc("Get subscriber", "Get a subscriber"), tonic.Handler(r.getSubscriber, http.StatusOK))
		sub.POST("/imsi/:imsi", formatDoc("Add subscriber", "Add subscriber"), tonic.Handler(r.addSubscriber, http.StatusCreated))
		sub.GET("/imsi/:imsi/flow", formatDoc("Get flows", "Get a subscriber UE data path"), tonic.Handler(r.getSubscriberFlows, http.StatusOK))
		sub.DELETE("/imsi/:imsi", formatDoc("Delete subscriber", "Delete subscriber"), tonic.Handler(r.deleteSubscriber, http.StatusOK))
		sub.PATCH("/imsi/:imsi", formatDoc("Update subscriber policy", "Update subscriber policy"), tonic.Handler(r.updateSubscriber, http.StatusOK))
	}
}

func (r *Router) getService(c *gin.Context, req *api.GetServiceRequest) (*api.ServiceResponse, error) {
	status := r.controller.ServiceStatus()
	return &status, nil
}

func (r *Router) setService(c *gin.Context, req *api.ServiceRequest) error {
	return r.controller.SetService(c, req)
}

func (r *Router) createSession(c *gin.Context, req *api.CreateSession) error {
	log.Infof("Received request for create session: %v", req)
	req.ImsiStr = uintArrayToString(req.Imsi)
	req.IpStr = uintToIp(req.Ip)
	return r.controller.CreateSession(c, req)
}

func (r *Router) endSession(c *gin.Context, req *api.EndSession) error {
	log.Infof("Received request for end session: %v", req)
	req.ImsiStr = uintArrayToString(req.Imsi)
	return r.controller.EndSession(c, req)
}

func (r *Router) getSessionByID(c *gin.Context, req *api.GetSessionByID) (*api.SessionResponse, error) {
	return r.controller.GetSessionByID(c, req)
}

func (r *Router) getCDRBySessionId(c *gin.Context, req *api.GetCDRBySessionId) (*api.CDR, error) {
	return r.controller.GetCDRBySessionId(c, req)
}

func (r *Router) getCDRByImsi(c *gin.Context, req *api.GetCDRByImsi) ([]*api.CDR, error) {
	return r.controller.GetCDRByImsi(c, req)
}

func (r *Router) getPolicyByImsi(c *gin.Context, req *api.GetPolicyByImsi) (*api.PolicyResponse, error) {
	return r.controller.GetPolicyByImsi(c, req)
}

func (r *Router) getPolicyByID(c *gin.Context, req *api.GetPolicyByID) (*api.PolicyResponse, error) {
	return r.controller.GetPolicyByID(c, req)
}

func (r *Router) addPolicy(c *gin.Context, req *api.Policy) error {
	return r.controller.AddPolicy(c, req)
}

func (r *Router) getRerouteForImsi(c *gin.Context, req *api.GetReRouteByImsi) (*api.ReRouteResponse, error) {
	return r.controller.GetReroute(c, req)
}

func (r *Router) updateReroute(c *gin.Context, req *api.UpdateRerouteById) error {
	return r.controller.UpdateReroute(c, req)
}

func (r *Router) addSubscriber(c *gin.Context, req *api.CreateSubscriber) error {
	return r.controller.AddSubscriber(c, req)
}

func (r *Router) updateSubscriber(c *gin.Context, req *api.UpdateSubscriber) error {
	return r.controller.UpdateSubscriber(c, req)
}

func (r *Router) deleteSubscriber(c *gin.Context, req *api.RequestSubscriber) error {
	return r.controller.DeleteSubscriber(c, req)
}

func (r *Router) getSubscriber(c *gin.Context, req *api.RequestSubscriber) (*api.SubscriberResponse, error) {
	return r.controller.GetSubscriber(c, req)
}

func (r *Router) getSubscriberFlows(c *gin.Context, req *api.GetFlowsForImsi) ([]*api.FlowResponse, error) {
	return r.controller.GetFlowsForImsi(c, req)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func uintArrayToString(array []uint8) string {
	str := ""
	for _, i := range array {
		str = str + strconv.Itoa(int(i))
	}
	log.Debugf("Byte %v to string %s", array, str)
	return str
}

func uintToIp(nn uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip.String()
}
