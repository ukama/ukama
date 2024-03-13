/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/test-sys/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/test-sys/api-gateway/pkg"
	"github.com/wI2L/fizz/openapi"

	log "github.com/sirupsen/logrus"
)
 
 var REDIRECT_URI = "https://subscriber.dev.ukama.com/swagger/#/"
 
 type Router struct {
	 f      *fizz.Fizz
	 config *RouterConfig
 }
 
 type RouterConfig struct {
	 debugMode  bool
	 serverConf *rest.HttpConfig
 }
 
 func NewRouter(config *RouterConfig) *Router {
	 r := &Router{
		 config: config,
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
	 r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName,
		 version.Version, r.config.debugMode, "")
 
	 routes := r.f.Group("/v1", "test API GW ", "Test system version v1")
 
	 routes.Use()
	 {
		 // mailer routes
		 test := routes.Group("/test", "Test", "Test")
		 test.POST("/ping", formatDoc("get ping resp", ""), tonic.Handler(r.pingHandler, http.StatusOK))
	 }
 }
 
 func formatDoc(summary string, description string) []fizz.OperationOption {
	 return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		 info.Summary = summary
		 info.Description = description
	 }}
 }
 
 func (r *Router) pingHandler(c *gin.Context) {
	 c.String(http.StatusOK, "pong hello")
 }
 