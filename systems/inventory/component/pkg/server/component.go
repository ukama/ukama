/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"encoding/json"
	"runtime/debug"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/common/gitClient"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/component/pkg"
	"github.com/ukama/ukama/systems/inventory/component/pkg/db"
	"github.com/ukama/ukama/systems/inventory/component/pkg/utils"
	"github.com/ukama/ukama/systems/inventory/component/scheduler"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cfactory "github.com/ukama/ukama/systems/common/rest/client/factory"
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
)

const uuidParsingError = "Error parsing UUID"

const (
	jobTag               = "node-sync-job"
	eventPublishErrorMsg = "Failed to publish message %+v with key %+v. Errors %v"
)

// Scheduler supervision tuning. The supervisor periodically reconciles the
// desired state (running) with the actual state and restarts the scheduler
// with capped exponential backoff if it is found dead.
const (
	schedulerHealthCheckInterval = 30 * time.Second
	schedulerRestartBackoffMin   = 1 * time.Second
	schedulerRestartBackoffMax   = 5 * time.Minute
)

type ComponentServer struct {
	pb.UnimplementedComponentServiceServer
	orgName            string
	gitClient          gitClient.GitClient
	componentRepo      db.ComponentRepo
	msgbus             mb.MsgBusServiceClient
	baseRoutingKey     msgbus.RoutingKeyBuilder
	componentScheduler scheduler.ComponentScheduler
	pushGateway        string
	gitDirPath         string
	factoryClient      cfactory.NodeFactoryClient
	config             *pkg.Config

	// schedulerMu guards schedulerCancel so the supervisor lifecycle can be
	// controlled safely from concurrent RPC calls and service startup.
	schedulerMu     sync.Mutex
	schedulerCancel context.CancelFunc
}

func NewComponentServer(orgName string, componentRepo db.ComponentRepo, msgBus mb.MsgBusServiceClient, pushGateway string, gc gitClient.GitClient, path string, factoryClient cfactory.NodeFactoryClient, config *pkg.Config) *ComponentServer {
	componentScheduler := scheduler.NewComponentScheduler(config.SchedulerInterval)
	c := &ComponentServer{
		gitClient:          gc,
		gitDirPath:         path,
		msgbus:             msgBus,
		config:             config,
		orgName:            orgName,
		pushGateway:        pushGateway,
		factoryClient:      factoryClient,
		componentRepo:      componentRepo,
		componentScheduler: componentScheduler,
		baseRoutingKey:     msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}

	return c
}

func (c *ComponentServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting component %v", req)

	cuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of component uuid. Error %s", err.Error())
	}
	component, err := c.componentRepo.Get(cuuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetResponse{
		Component: dbComponentToPbComponent(component),
	}, nil
}

// Deprecated: This function is deprecated and will be removed in a future version. Use List instead.
func (c *ComponentServer) GetByUser(ctx context.Context, req *pb.GetByUserRequest) (*pb.GetByUserResponse, error) {
	log.Infof("Getting components by user %v", req)

	components, err := c.componentRepo.GetByUser(req.GetUserId(), int32(ukama.ParseComponentCategory(req.GetCategory())))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.GetByUserResponse{
		Components: dbComponentsToPbComponents(components),
	}, nil
}

func (c *ComponentServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("Listing components %v", req)

	components, err := c.componentRepo.List(req.GetUserId(), req.GetPartNumber(), int32(ukama.ParseComponentCategory(req.GetCategory())))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "component")
	}

	return &pb.ListResponse{
		Components: dbComponentsToPbComponents(components),
	}, nil
}

func (c *ComponentServer) SyncComponents(ctx context.Context, req *pb.SyncComponentsRequest) (*pb.SyncComponentsResponse, error) {
	log.Infof("Syncing components %v", req)

	if err := c.gitClient.SetupDir(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to setup directory: %s", err.Error())
	}

	err := c.gitClient.CloneGitRepo("main")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clone git repo: %s", err.Error())
	}

	rootFileContent, err := c.gitClient.ReadFileJSON(c.gitDirPath + "/root.json")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read file. Error %s", err.Error())
	}

	var env gitClient.Environment
	err = json.Unmarshal(rootFileContent, &env)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal json. Error %s", err.Error())
	}

	var components []utils.Component
	environment := utils.GetEnvironmentField(env, c.config.ComponentEnvironment)

	for _, company := range environment {
		err := c.gitClient.BranchCheckout(company.GitBranchName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to checkout branch. Error %s", err.Error())
		}

		paths, _ := c.gitClient.GetFilesPath("components")
		for _, path := range paths {
			content, err := c.gitClient.ReadFileYML(path)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to read file. Error %s", err.Error())
			}
			var component utils.Component
			err = yaml.Unmarshal(content, &component)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal yaml. Error %s", err.Error())
			}

			components = append(components, component)
		}

		company.UserId = c.config.OwnerId
		userIdUUID, err := uuid.FromString(company.UserId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
		}

		cdb := utilComponentsToDbComponents(components, userIdUUID)

		err = c.componentRepo.Delete()
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "component")
		}
		log.Info("Deleted all component records")

		err = c.componentRepo.Add(cdb)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "component")
		}

		route := c.baseRoutingKey.SetAction("sync").SetObject("components").MustBuild()
		err = c.msgbus.PublishRequest(route, req)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
		}
	}

	return &pb.SyncComponentsResponse{}, nil
}

// StartSchedulerSupervisor starts the node-sync scheduler together with a
// supervisor goroutine that keeps it alive, restarting it automatically if it
// dies unexpectedly. It is safe to call multiple times: a supervisor is only
// started once until it is stopped. The provided context bounds the lifetime
// of the supervisor (e.g. cancel it on graceful service shutdown).
func (c *ComponentServer) StartSchedulerSupervisor(ctx context.Context) {
	c.schedulerMu.Lock()
	if c.schedulerCancel != nil {
		c.schedulerMu.Unlock()
		log.Info("Scheduler supervisor already running, skipping start")

		return
	}

	supervisorCtx, cancel := context.WithCancel(ctx)
	c.schedulerCancel = cancel
	c.schedulerMu.Unlock()

	go c.superviseScheduler(supervisorCtx)
}

// stopSchedulerSupervisor cancels the supervisor goroutine if one is running.
func (c *ComponentServer) stopSchedulerSupervisor() {
	c.schedulerMu.Lock()
	defer c.schedulerMu.Unlock()

	if c.schedulerCancel != nil {
		c.schedulerCancel()
		c.schedulerCancel = nil
	}
}

// superviseScheduler is a self-healing control loop. It (re)starts the
// scheduler whenever it is found not running, applying capped exponential
// backoff on repeated start failures, and exits cleanly when ctx is cancelled.
func (c *ComponentServer) superviseScheduler(ctx context.Context) {
	log.Info("Scheduler supervisor started")

	backoff := schedulerRestartBackoffMin

	for {
		if !c.componentScheduler.IsRunning() {
			if err := c.startSchedulerOnce(); err != nil {
				log.Errorf("Failed to start scheduler: %v. Retrying in %s", err, backoff)

				if sleepOrDone(ctx, backoff) {
					break
				}

				backoff = nextBackoff(backoff)

				continue
			}

			log.Info("Scheduler is running")
			backoff = schedulerRestartBackoffMin
		}

		if sleepOrDone(ctx, schedulerHealthCheckInterval) {
			break
		}
	}

	log.Info("Scheduler supervisor stopping (context cancelled)")

	if err := c.componentScheduler.Stop(); err != nil {
		log.Errorf("Failed to stop scheduler during supervisor shutdown: %v", err)
	}
}

// startSchedulerOnce kicks off an immediate, non-blocking sync run and then
// registers the recurring job. NodeSyncJob is always run through the panic-safe
// wrapper so a failing run can never crash the process or the supervisor.
func (c *ComponentServer) startSchedulerOnce() error {
	go c.safeNodeSyncJob(context.Background())

	return c.componentScheduler.Start(jobTag, func() {
		c.safeNodeSyncJob(context.Background())
	})
}

// safeNodeSyncJob runs NodeSyncJob with panic recovery so a panic in a single
// run is logged and contained instead of taking down the service.
func (c *ComponentServer) safeNodeSyncJob(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic in node sync job: %v\n%s", r, debug.Stack())
		}
	}()

	c.NodeSyncJob(ctx)
}

// nextBackoff doubles the backoff duration up to the configured maximum.
func nextBackoff(d time.Duration) time.Duration {
	next := d * 2
	if next > schedulerRestartBackoffMax {
		return schedulerRestartBackoffMax
	}

	return next
}

// sleepOrDone waits for d or until ctx is cancelled. It returns true if ctx was
// cancelled (i.e. the caller should stop).
func sleepOrDone(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return true
	case <-timer.C:
		return false
	}
}

func (c *ComponentServer) StartScheduler(ctx context.Context, req *pb.StartSchedulerRequest) (*pb.StartSchedulerResponse, error) {
	log.Info("StartScheduler RPC invoked")

	// Start (or ensure) the supervised scheduler. This returns immediately;
	// the initial sync runs in the background so the RPC is never blocked.
	c.StartSchedulerSupervisor(context.Background())

	return &pb.StartSchedulerResponse{}, nil
}

func (c *ComponentServer) StopScheduler(ctx context.Context, req *pb.StopSchedulerRequest) (*pb.StopSchedulerResponse, error) {
	log.Info("StopScheduler RPC invoked")

	// Cancel the supervisor first so it does not immediately restart the
	// scheduler we are about to stop.
	c.stopSchedulerSupervisor()

	if err := c.componentScheduler.Stop(); err != nil {
		log.Errorf("Failed to stop scheduler. Error %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to stop scheduler: %s", err.Error())
	}

	return &pb.StopSchedulerResponse{}, nil
}

func (c *ComponentServer) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	log.Infof("Verifying component %v", req)

	_, err := c.componentRepo.Verify(req.GetPartNumber())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.VerifyResponse{}, nil
}

func (c *ComponentServer) NodeSyncJob(ctx context.Context) {
	log.Infof("Node sync job started")
	nodes, err := c.factoryClient.List("", c.orgName, true)
	if err != nil {
		log.Errorf("Failed to get nodes from factory. Error %s", err.Error())
		return
	}
	log.Infof("Nodes: %v", nodes)
	for _, node := range nodes.Nodes {
		log.Infof("Node: %s", node.Id)
		_, err := c.Verify(ctx, &pb.VerifyRequest{PartNumber: node.Id})
		if err != nil && status.Code(err) == codes.NotFound {
			log.Infof("Component not found, creating new component: %s", node.Id)
			ownerId, err := uuid.FromString(c.config.OwnerId)
			if err != nil {
				log.Errorf("Failed to parse owner ID. Error %s", err.Error())
				continue
			}
			component := &db.Component{
				Id:            uuid.NewV4(),
				PartNumber:    node.Id,
				Type:          node.Type,
				UserId:        ownerId,
				Description:   ukama.GetPlaceholderNameByType(node.Type),
				Category:      ukama.ParseComponentCategory(c.config.NodeComponentDetails.Category),
				Managed:       c.config.NodeComponentDetails.Managed,
				Warranty:      c.config.NodeComponentDetails.Warranty,
				ImagesURL:     c.config.NodeComponentDetails.ImagesURL,
				Inventory:     c.config.NodeComponentDetails.Inventory,
				DatasheetURL:  c.config.NodeComponentDetails.DatasheetURL,
				Manufacturer:  c.config.NodeComponentDetails.Manufacturer,
				Specification: c.config.NodeComponentDetails.Specification,
			}
			err = c.componentRepo.Add([]*db.Component{component})
			if err != nil {
				log.Errorf("Failed to add component. Error %s", err.Error())
				continue
			}
			route := c.baseRoutingKey.SetAction("added").SetObject("node").MustBuild()
			evt := &epb.EventInventoryNodeComponentAdd{
				Type:       component.Type,
				PartNumber: component.PartNumber,
				Id:         component.Id.String(),
				UserId:     component.UserId.String(),
			}

			err = c.msgbus.PublishRequest(route, evt)
			if err != nil {
				log.Errorf(eventPublishErrorMsg, evt, route, err)
			}
		} else if err != nil && status.Code(err) != codes.NotFound {
			log.Errorf("Failed to verify component. Error %s", err.Error())
			continue
		} else {
			log.Infof("Component already exists and verified: %s", node.Id)
		}
	}
}

func dbComponentToPbComponent(component *db.Component) *pb.Component {
	return &pb.Component{
		Id:            component.Id.String(),
		Type:          component.Type,
		Managed:       component.Managed,
		Warranty:      component.Warranty,
		Inventory:     component.Inventory,
		ImagesURL:     component.ImagesURL,
		PartNumber:    component.PartNumber,
		Description:   component.Description,
		DatasheetURL:  component.DatasheetURL,
		Manufacturer:  component.Manufacturer,
		Specification: component.Specification,
		UserId:        component.UserId.String(),
		Category:      ukama.ComponentCategory(component.Category).String(),
	}
}

func dbComponentsToPbComponents(components []*db.Component) []*pb.Component {
	res := []*pb.Component{}

	for _, i := range components {
		res = append(res, dbComponentToPbComponent(i))
	}

	return res
}

func utilComponentsToDbComponents(components []utils.Component, uId uuid.UUID) []*db.Component {
	res := []*db.Component{}

	for _, i := range components {
		res = append(res, &db.Component{
			UserId:        uId,
			Type:          i.Type,
			Managed:       i.Managed,
			Warranty:      i.Warranty,
			ImagesURL:     i.ImagesURL,
			Id:            uuid.NewV4(),
			PartNumber:    i.PartNumber,
			Description:   i.Description,
			Inventory:     i.InventoryID,
			DatasheetURL:  i.DatasheetURL,
			Manufacturer:  i.Manufacturer,
			Specification: i.Specification,
			Category:      ukama.ParseComponentCategory(i.Category),
		})
	}
	return res
}
