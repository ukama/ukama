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
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/config"
	ukamaPb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
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
	StoreHealthReport(req *healthPb.StoreHealthReportRequest) (*healthPb.StoreHealthReportResponse, error)
	ListReports(request *healthPb.ListReportsRequest) (*healthPb.ListReportsResponse, error)
	ListApps(request *healthPb.ListAppsRequest) (*healthPb.ListAppsResponse, error)
	ListInterfaces(request *healthPb.ListInterfacesRequest) (*healthPb.ListInterfacesResponse, error)
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

func (r *Router) Run() {
	r.logger.Info("Listening on port ", r.config.serverConf.Port)
	err := r.f.Engine().Run(fmt.Sprint(":", r.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, "")
	endpoint := r.f.Group("/v1", "API gateway", "node system version v1")
	endpoint.GET("/ping", formatDoc("Ping the server", "Returns a response indicating that the server is running."), tonic.Handler(r.pingHandler, http.StatusOK))

	health := endpoint.Group("/health", "Health", "Health service for the node")
	health.POST("/nodes/:node_id/performance",
		append(formatDoc("Store health report", "Path: node id. Body: full health JSON (nodeType, reportedAt RFC3339, schemaVersion, capabilities, system, interfaces, apps, events, etc.)."),
			fizz.InputModel(&StoreHealthReportOpenAPIInput{}),
		),
		tonic.Handler(r.postHealthReportHandler, http.StatusCreated),
	)
	health.POST("/logger/node/:node_id", formatDoc("Log data", "Endpoint to log data"), tonic.Handler(r.logHandler, http.StatusCreated))
	health.GET("/reports", formatDoc("List health reports", "Retrieve the health reports for the node."), tonic.Handler(r.listHealthReportsHandler, http.StatusOK))
	health.GET("/apps", formatDoc("List apps", "Retrieve the apps for the node."), tonic.Handler(r.listAppsHandler, http.StatusOK))
	health.GET("/interfaces", formatDoc("List interfaces", "Retrieve the interfaces for the node."), tonic.Handler(r.listInterfacesHandler, http.StatusOK))

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
		if err := file.Close(); err != nil {
			r.logger.Warnf("Failed to gracefully close log file: %v", err)
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

func (r *Router) postHealthReportHandler(c *gin.Context, req *StoreHealthReportRequest) (*healthPb.StoreHealthReportResponse, error) {
	nID, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	payload := req.HealthPayloadBytes()

	return r.clients.Health.StoreHealthReport(&healthPb.StoreHealthReportRequest{
		NodeId:  nID.StringLowercase(),
		Payload: payload,
	})
}

func (r *Router) listHealthReportsHandler(c *gin.Context, req *ListHealthReportsRequest) (*healthPb.ListReportsResponse, error) {
	var reportedAtPb *timestamppb.Timestamp
	if req.ReportedAt != "" {
		t, err := time.Parse(time.RFC3339, req.ReportedAt)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "reportedAt must be RFC3339: %v", err)
		}
		reportedAtPb = timestamppb.New(t)
	}

	return r.clients.Health.ListReports(&healthPb.ListReportsRequest{
		ReportId:   req.ReportId,
		NodeId:     req.NodeId,
		ReportedAt: reportedAtPb,
		Timeframe:  ukamaPb.FilterTimeframesType(ukama.ReturnFilterTimeframesType(ukama.ParseFilterTimeframesType(req.Timeframe))),
	})
}

func (r *Router) listAppsHandler(c *gin.Context, req *ListAppsRequest) (*healthPb.ListAppsResponse, error) {
	return r.clients.Health.ListApps(&healthPb.ListAppsRequest{
		NodeId: req.NodeId,
		AppName: req.AppName,
	})
}

func (r *Router) listInterfacesHandler(c *gin.Context, req *ListInterfacesRequest) (*healthPb.ListInterfacesResponse, error) {
	return r.clients.Health.ListInterfaces(&healthPb.ListInterfacesRequest{
		NodeId: req.NodeId,
		InterfaceName: req.InterfaceName,
	})
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
