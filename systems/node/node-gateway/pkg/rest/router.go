/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/node/node-gateway/cmd/version"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg/client"

	log "github.com/sirupsen/logrus"
	healthPb "github.com/ukama/ukama/systems/node/health/pb/gen"
	npb "github.com/ukama/ukama/systems/node/notify/pb/gen"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
	logger  *log.Logger
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
}

type Clients struct {
	Health health
	Notify notify
}

type notify interface {
	Add(nodeId, severity, ntype, serviceName string, details json.RawMessage, status, epochTime uint32) (*npb.AddResponse, error)
	Get(id string) (*npb.GetResponse, error)
	List(nodeId, serviceName, nType string, count uint32, sort bool) (*npb.ListResponse, error)
	Delete(id string) (*npb.DeleteResponse, error)
	Purge(nodeId, serviceName, nType string) (*npb.ListResponse, error)
}

type health interface {
	StoreRunningAppsInfo(req *healthPb.StoreRunningAppsInfoRequest) (*healthPb.StoreRunningAppsInfoResponse, error)
	GetRunningAppsInfo(nodeId string) (*healthPb.GetRunningAppsResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Health = client.NewHealth(endpoints.Health, endpoints.Timeout)
	c.Notify = client.NewNotify(endpoints.Notify, endpoints.Timeout)
	return c
}

func NewRouter(clients *Clients, config *RouterConfig) *Router {
	r := &Router{
		clients: clients,
		config:  config,
		logger:  log.New(), // Initialize the logger
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
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
	endpoint := r.f.Group("/v1", "API gateway", "node system version v1")
	endpoint.GET("/ping", formatDoc("Ping the server", "Returns a response indicating that the server is running."), tonic.Handler(r.pingHandler, http.StatusOK))

	health := endpoint.Group("/health", "Health", "Health service for the node")
	health.POST("/nodes/:node_id/performance", formatDoc("Create system performance report", "This endpoint allows you to create and update system performance information."), tonic.Handler(r.postSystemPerformanceInfoHandler, http.StatusCreated))
	health.GET("/nodes/:node_id/performance", formatDoc("Get system performance report", "Retrieve system performance information for analysis and monitoring."), tonic.Handler(r.getSystemPerformanceInfoHandler, http.StatusOK))
	health.POST("/logger/node/:node_id", formatDoc("Log data", "Endpoint to log data"), tonic.Handler(r.logHandler, http.StatusCreated))

	notif := endpoint.Group("/notify", "Node Notify", "Notify service for the node")
	notif.POST("", formatDoc("Insert Notification", "Insert a new notification"), tonic.Handler(r.postNotification, http.StatusCreated))
	notif.GET("", formatDoc("Get Notifications", "Get a list of notifications"), tonic.Handler(r.getNotifications, http.StatusOK))
	notif.GET("/:notification_id", formatDoc("Get Notification", "Get a specific notification"), tonic.Handler(r.getNotification, http.StatusOK))
	notif.DELETE("", formatDoc("Delete Notifications", "Delete matching notifications"), tonic.Handler(r.deleteNotifications, http.StatusOK))
	notif.DELETE("/:notification_id", formatDoc("Delete Notification", "Delete a specific notification"), tonic.Handler(r.deleteNotification, http.StatusOK))
}

func (r *Router) logHandler(c *gin.Context, req *AddLogsRequest) (string, error) {
	// Ensure that the NodeID is provided
	if req.NodeId == "" {
		r.logger.Error("NodeID is required")
		return "", errors.New("NodeID is required")
	}

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	// Log the NodeID
	r.logger.WithField("node_id", nId).Info("Received log data")

	// Define the log file path
	logFileName := "log_" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
	logFilePath := filepath.Join("logs", logFileName)

	if err := os.MkdirAll("logs", 0755); err != nil {
		r.logger.Errorf("Failed to create logs directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save log data"})
		return "", err
	}

	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		r.logger.Errorf("Failed to open log file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save log data"})
		return "", err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Warnf("Failed to gracefully close log file: %v", err)
		}
	}()

	for _, logEntry := range req.Logs {
		logLine := fmt.Sprintf("%s %s %s %s %s\n", nId, logEntry.AppName, logEntry.Time, logEntry.Level, logEntry.Message)
		if _, err := file.WriteString(logLine); err != nil {
			r.logger.Errorf("Failed to write log entry to file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save log data"})
			return "", err
		}
		if err := file.Sync(); err != nil {
			r.logger.Errorf("Failed to flush log data to file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save log data"})
			return "", err
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Log data received and saved"})
	return "Logs received", nil
}

func (r *Router) postSystemPerformanceInfoHandler(c *gin.Context, req *StoreRunningAppsInfoRequest) (*healthPb.StoreRunningAppsInfoResponse, error) {
	var genSystems []*healthPb.System
	for _, sys := range req.System {
		genSystem := &healthPb.System{
			Name:  sys.Name,
			Value: sys.Value,
		}
		genSystems = append(genSystems, genSystem)
	}

	var genCapps []*healthPb.Capps
	for _, capp := range req.Capps {
		var genResources []*healthPb.Resource
		for _, resource := range capp.Resources {
			genResource := &healthPb.Resource{
				Name:  resource.Name,
				Value: resource.Value,
			}
			genResources = append(genResources, genResource)
		}
		genCapp := &healthPb.Capps{
			Space:     capp.Space,
			Name:      capp.Name,
			Tag:       capp.Tag,
			Status:    healthPb.Status(healthPb.Status_value[capp.Status]),
			Resources: genResources,
		}
		genCapps = append(genCapps, genCapp)
	}

	return r.clients.Health.StoreRunningAppsInfo(&healthPb.StoreRunningAppsInfoRequest{
		NodeId:    req.NodeId,
		Timestamp: req.Timestamp,
		System:    genSystems,
		Capps:     genCapps,
	})
}

func (r *Router) getSystemPerformanceInfoHandler(c *gin.Context, req *GetRunningAppsRequest) (*healthPb.GetRunningAppsResponse, error) {
	resp, err := r.clients.Health.GetRunningAppsInfo(req.NodeId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resp, nil
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

func (r *Router) postNotification(c *gin.Context, req *AddNotificationReq) (*npb.AddResponse, error) {
	return r.clients.Notify.Add(req.NodeId, req.Severity,
		req.Type, req.ServiceName, req.Details, req.Status, req.Time)
}

func (r *Router) getNotification(c *gin.Context, req *GetNotificationReq) (*npb.GetResponse, error) {
	return r.clients.Notify.Get(req.NotificationId)
}

func (r *Router) getNotifications(c *gin.Context, req *GetNotificationsReq) (*npb.ListResponse, error) {
	return r.clients.Notify.List(req.NodeId, req.ServiceName, req.Type, req.Count, req.Sort)
}

func (r *Router) deleteNotification(c *gin.Context, req *GetNotificationReq) (*npb.DeleteResponse, error) {
	return r.clients.Notify.Delete(req.NotificationId)
}

func (r *Router) deleteNotifications(c *gin.Context, req *DelNotificationsReq) (*npb.ListResponse, error) {
	return r.clients.Notify.Purge(req.NodeId, req.ServiceName, req.Type)
}
