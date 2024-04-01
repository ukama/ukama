/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ukama/ukama/systems/common/rest"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/node/node-gateway/cmd/version"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg/client"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/node/node-gateway/pkg"

	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/node/health/pb/gen"
	healthPb "github.com/ukama/ukama/systems/node/health/pb/gen"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
	logger  *logrus.Logger 
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
}

type Clients struct {
	Health health
}

type health interface {
	StoreRunningAppsInfo(req *healthPb.StoreRunningAppsInfoRequest) (*healthPb.StoreRunningAppsInfoResponse, error)
	GetRunningAppsInfo(nodeId string) (*healthPb.GetRunningAppsResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Health = client.NewHealth(endpoints.Health, endpoints.Timeout)

	return c
}

func NewRouter(clients *Clients, config *RouterConfig) *Router {
    r := &Router{
        clients: clients,
        config:  config,
    }

    if !config.debugMode {
        gin.SetMode(gin.ReleaseMode)
    }

    r.init() // Remove the authfunc parameter
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
    r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode,"")
    
    endpoint := r.f.Group("/v1", "API gateway", "node system version v1")
    endpoint.GET("/ping", formatDoc("Ping the server", "Returns a response indicating that the server is running."), tonic.Handler(r.pingHandler, http.StatusOK))
    endpoint.POST("/health/nodes/:node_id/performance", formatDoc("Create system performance report", "This endpoint allows you to create and update system performance information."), tonic.Handler(r.postSystemPerformanceInfoHandler, http.StatusCreated))
    endpoint.GET("/health/nodes/:node_id/performance", formatDoc("Get system performance report", "Retrieve system performance information for analysis and monitoring."), tonic.Handler(r.getSystemPerformanceInfoHandler, http.StatusOK))
	endpoint.POST("/logger/nodes/:node_id", formatDoc("Log data", "Endpoint to log data"),tonic.Handler(r.logHandler,http.StatusCreated))

}

func (r *Router) logHandler(c *gin.Context, req *AddLogsRequest) (string, error) {
    if err := c.BindJSON(&req.Logs); err != nil {
        r.logger.Errorf("Failed to bind JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return "", err
    }

    // Ensure that the NodeID is provided
    if req.NodeId == "" {
        r.logger.Error("NodeID is required")
        c.JSON(http.StatusBadRequest, gin.H{"error": "NodeID is required"})
        return "", errors.New("NodeID is required")
    }

    // Log the NodeID
    r.logger.WithField("node_id", req.NodeId).Info("Received log data")

    // Define the log file path
    logFileName := "log_" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
    logFilePath := filepath.Join("logs", logFileName)

    // Create the "logs" directory if it doesn't exist
    if err := os.MkdirAll("logs", 0755); err != nil {
        r.logger.Errorf("Failed to create logs directory: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save log data"})
        return "", err
    }

    // Open the log file for appending
    file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
    if err != nil {
        r.logger.Errorf("Failed to open log file: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save log data"})
        return "", err
    }
    defer file.Close()

    // Write each log entry to the log file
    for _, logEntry := range req.Logs {
        logLine := fmt.Sprintf("%s %s %s %s %s\n", req.NodeId, logEntry.AppName, logEntry.Time, logEntry.Level, logEntry.Message)
        if _, err := file.WriteString(logLine); err != nil {
            r.logger.Errorf("Failed to write log entry to file: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save log data"})
            return "", err
        }
    }

    c.JSON(http.StatusOK, gin.H{"message": "Log data received and saved"})
    return "Logs received", nil
}



func (r *Router) postSystemPerformanceInfoHandler(c *gin.Context, req *StoreRunningAppsInfoRequest) (*healthPb.StoreRunningAppsInfoResponse, error) {
	var genSystems []*gen.System
	for _, sys := range req.System {
		genSystem := &gen.System{
			Name:  sys.Name,
			Value: sys.Value,
		}
		genSystems = append(genSystems, genSystem)
	}

	var genCapps []*gen.Capps
	for _, capp := range req.Capps {
		var genResources []*gen.Resource
		for _, resource := range capp.Resources {
			genResource := &gen.Resource{
				Name:  resource.Name,
				Value: resource.Value,
			}
			genResources = append(genResources, genResource)
		}
		genCapp := &gen.Capps{
			Space:     capp.Space,
			Name:      capp.Name,
			Tag:       capp.Tag,
			Status:    gen.Status(gen.Status_value[capp.Status]),
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
		logrus.Error(err)
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

    response["status"] = pkg.SystemName+" is running"
    response["version"] = version.Version

    c.JSON(http.StatusOK, response)

    return nil
}
