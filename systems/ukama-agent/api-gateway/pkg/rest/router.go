package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	"github.com/wI2L/fizz/openapi"
)

const ORG_URL_PARAMETER = "org"

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
}

type Clients struct {
	a asr
}

type asr interface {
	Activate(req *pb.ActivateReq) (*pb.ActivateResp, error)
	Inactivate(req *pb.InactivateReq) (*pb.InactivateResp, error)
	UpdatePackage(req *pb.UpdatePackageReq) (*pb.UpdatePackageResp, error)
	Read(req *pb.ReadReq) (*pb.ReadResp, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.a = client.NewAsr(endpoints.Asr, endpoints.Timeout)
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
	logrus.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init() {

	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "ukama-agent ", "Ukama-agent system")

	asr := v1.Group("asr/subscriber/", "Asr", "Active susbcriber registry")
	asr.GET(":iccid", formatDoc("Get Orgs Credential", ""), tonic.Handler(r.getActiveSubscriber, http.StatusOK))
	asr.PUT(":iccid", formatDoc("Activate: Add a new subscriber", ""), tonic.Handler(r.putSubscriber, http.StatusCreated))
	asr.DELETE(":iccid", formatDoc("Inactivate: Remove a susbcriber", ""), tonic.Handler(r.deleteSubscriber, http.StatusOK))
	asr.PATCH(":iccid", formatDoc("Update package id", ""), tonic.Handler(r.patchPackageUpdate, http.StatusOK))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) putSubscriber(c *gin.Context, req *ActivateReq) (*pb.ActivateResp, error) {

	return r.clients.a.Activate(&pb.ActivateReq{
		Iccid:     req.Iccid,
		Network:   req.Network,
		PackageId: req.PackageId,
	})
}

func (r *Router) deleteSubscriber(c *gin.Context, req *InactivateReq) (*pb.InactivateResp, error) {

	return r.clients.a.Inactivate(&pb.InctivateReq{
		Id: &pb.InactivateReq_Iccid{
			Iccid: req.Iccid,
		},
	})
}

func (r *Router) patchPackageUpdate(c *gin.Context, req *UpdatePackageReq) (*pb.UpdatePackageResp, error) {

	return r.clients.a.Activate(&pb.ActivateReq{
		Iccid:     req.Iccid,
		PackageId: req.PackageId,
	})
}

func (r *Router) getActiveSubscriber(c *gin.Context, req *ReadSusbscriberReq) (*pb.ReadSubscriberResp, error) {
	return r.clients.a.Read(&pb.ReadReq{
		Id: &pb.ReadReq_Iccid{
			Iccid: req.Iccid,
		},
	})
}
