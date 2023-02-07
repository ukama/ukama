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
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/node-gateway/cmd/version"
	"github.com/ukama/ukama/systems/ukama-agent/node-gateway/pkg"
	"github.com/ukama/ukama/systems/ukama-agent/node-gateway/pkg/client"
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
	UpdateGuti(req *pb.UpdateGutiReq) (*pb.UpdateGutiResp, error)
	UpdateTai(req *pb.UpdateTaiReq) (*pb.UpdateTaiResp, error)
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
	v1 := r.f.Group("/v1", "ukama-agent-node-gateway ", "Ukama-agent system")

	asr := v1.Group("/subscriber", "Asr", "Active susbcriber registry")
	asr.GET("/:imsi", formatDoc("Get Subscriber", ""), tonic.Handler(r.getActiveSubscriber, http.StatusOK))
	asr.POST("/:imsi/guti", formatDoc("GUTI update for subscriber", ""), tonic.Handler(r.postGuti, http.StatusOK))
	asr.POST("/:imsi/tai", formatDoc("TAI update for subscriber", ""), tonic.Handler(r.postTai, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) postGuti(c *gin.Context, req *UpdateGutiReq) (*pb.UpdateGutiResp, error) {

	return r.clients.a.UpdateGuti(&pb.UpdateGutiReq{
		Imsi:      req.Imsi,
		UpdatedAt: req.UpdatedAt,
		Guti: &pb.Guti{
			PlmnId: req.Guti.PlmnId,
			Mmegi:  req.Guti.Mmegi,
			Mmec:   req.Guti.Mmec,
			Mtmsi:  req.Guti.Mtmsi,
		},
	})

}

func (r *Router) postTai(c *gin.Context, req *UpdateTaiReq) (*pb.UpdateTaiResp, error) {

	return r.clients.a.UpdateTai(&pb.UpdateTaiReq{
		Imsi:      req.Imsi,
		UpdatedAt: req.UpdatedAt,
		Tac:       req.Tac,
	})
}

func (r *Router) getActiveSubscriber(c *gin.Context, req *GetSubscriberReq) (*pb.ReadResp, error) {
	return r.clients.a.Read(&pb.ReadReq{
		Id: &pb.ReadReq_Imsi{
			Imsi: req.Imsi,
		},
	})
}
