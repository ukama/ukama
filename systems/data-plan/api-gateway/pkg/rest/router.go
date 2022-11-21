package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb"
	"github.com/wI2L/fizz/openapi"
)


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
	d dataPlan
}

type dataPlan interface {
	AddPackage(req *pb.AddPackageRequest) (*pb.AddPackageResponse, error)
	UpdatePackage(req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error)
	GetPackage(req *pb.GetPackagesRequest) (*pb.GetPackagesResponse, error)
	DeletePackage(req *pb.DeletePackageRequest) (*pb.DeletePackageResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.d = client.NewDataPlan(endpoints.Package, endpoints.Timeout)
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
		logrus.Error(err)
	}
}

func (r *Router) init() {
	const pack = "/packages"
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "Data-plan system ", "Data-plan  system version v1")

	packages := v1.Group(pack, "Packages", "looking for packages credentials")
	packages.GET("/:package", formatDoc("Get package", ""), tonic.Handler(r.getPackageHandler, http.StatusOK))
	packages.PUT("", formatDoc("Add Package", ""), tonic.Handler(r.AddPackageHandler, http.StatusCreated))
	packages.PATCH("", formatDoc("Update Package", ""), tonic.Handler(r.UpdatePackageHandler, http.StatusOK))
	packages.DELETE("/:package", formatDoc("Delete Package", ""), tonic.Handler(r.deletePackageHandler, http.StatusOK))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
func (p *Router) getPackageHandler(c *gin.Context, req *GetPackagesRequest) (*pb.GetPackagesResponse, error) {

	_id, err := strconv.ParseUint(c.Param("package"), 10, 64)
	if err != nil {
		logrus.Error(err)
	}
	
	resp,err:= p.clients.d.GetPackage(&pb.GetPackagesRequest{
		Id:    _id,
		OrgId: 12345,
	})
	if err != nil {
		logrus.Error(err)
	}

	return resp,nil
}
func (p *Router) deletePackageHandler(c *gin.Context, req *DeletePackageRequest) (*pb.DeletePackageResponse, error) {
	_id, err := strconv.ParseUint(c.Param("package"), 10, 64)
	if err != nil {
		logrus.Error(err)
	}
	resp,err:= p.clients.d.DeletePackage(&pb.DeletePackageRequest{
		Id:    _id,
		OrgId: 12345,
	})
	if err != nil {
		logrus.Error(err)
	}
	return resp,nil
}
func (p *Router) UpdatePackageHandler(c *gin.Context, req *UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	resp,err:= p.clients.d.UpdatePackage(&pb.UpdatePackageRequest{
		Id:req.Id,
		Name:         req.Name,
		SimType:     ReqStrTopb(req.SimType),
		Active:       req.Active,
		Duration:     req.Duration,
		SmsVolume:   req.SmsVolume,
		DataVolume:  req.DataVolume,
		VoiceVolume: req.VoiceVolume,
		OrgRatesId: req.OrgRatesId,
	})
	if err != nil {
		logrus.Error(err)
	}

	return resp,nil

	
}
func (p *Router) AddPackageHandler(c *gin.Context, req *AddPackageRequest) (*pb.AddPackageResponse, error) {
	return p.clients.d.AddPackage(&pb.AddPackageRequest{
		Name:   req.Name,
		OrgId:req.OrgId,
		Duration:req.Duration,
		OrgRatesId:req.OrgRatesId,
		VoiceVolume:req.VoiceVolume,
		Active:req.Active,
		DataVolume:req.DataVolume,
		SmsVolume:req.SmsVolume,
		SimType:     ReqStrTopb(req.SimType),
	})
}

func ReqStrTopb(e string) pb.SimType {
	switch e {
	case "inter_none":
		return pb.SimType_INTER_NONE
	case "inter_mno_data":
		return pb.SimType_INTER_MNO_DATA
	case "inter_ukama_all":
		return pb.SimType_INTER_UKAMA_ALL
	case "inter_mno_all":
		return pb.SimType_INTER_MNO_ALL
	default:
		return pb.SimType_INTER_NONE
	}
}