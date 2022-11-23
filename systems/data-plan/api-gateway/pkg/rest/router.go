package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/client"
	pbBaseRate "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

const ORG_ID = 12345

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
	UploadBaseRates(req *pbBaseRate.UploadBaseRatesRequest) (*pbBaseRate.UploadBaseRatesResponse, error)
	GetBaseRates(req *pbBaseRate.GetBaseRatesRequest) (*pbBaseRate.GetBaseRatesResponse, error)
	GetBaseRate(req *pbBaseRate.GetBaseRateRequest) (*pbBaseRate.GetBaseRateResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.d = client.NewDataPlan(endpoints.Package, endpoints.BaseRate, endpoints.Timeout)

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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "Data-plan system ", "Data-plan  system version v1")

	baseRates := v1.Group("/baseRates", "BaseRates", "BaseRates operations")
	baseRates.POST("", formatDoc("Upload baseRates", ""), tonic.Handler(r.uploadBaseRateHandler, http.StatusOK))
	baseRates.GET("/:baseRate", formatDoc("Get BaseRate", ""), tonic.Handler(r.getBaseRateHandler, http.StatusCreated))
	baseRates.GET("", formatDoc("Get BaseRates", ""), tonic.Handler(r.getBaseRatesHandler, http.StatusOK))

	packages := v1.Group("/packages", "Packages", "Packages operations")
	packages.PUT("", formatDoc("Add Package", ""), tonic.Handler(r.AddPackageHandler, http.StatusCreated))
	packages.GET("/:package", formatDoc("Get package", ""), tonic.Handler(r.getPackageHandler, http.StatusOK))
	packages.PATCH("", formatDoc("Update Package", ""), tonic.Handler(r.UpdatePackageHandler, http.StatusOK))
	packages.DELETE("/:package", formatDoc("Delete Package", ""), tonic.Handler(r.deletePackageHandler, http.StatusOK))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
func (p *Router) getBaseRateHandler(c *gin.Context, req *GetBaseRateRequest) (*pbBaseRate.GetBaseRateResponse, error) {

	RateId, err := strconv.ParseUint(c.Param("baseRate"), 10, 64)
	if err != nil {
		logrus.Error(err)
	}

	resp, err := p.clients.d.GetBaseRate(&pbBaseRate.GetBaseRateRequest{
		RateId: RateId,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) uploadBaseRateHandler(c *gin.Context, req *UploadBaseRatesRequest) (*pbBaseRate.UploadBaseRatesResponse, error) {

	resp, err := p.clients.d.UploadBaseRates(&pbBaseRate.UploadBaseRatesRequest{
		FileURL:     req.FileURL,
		EffectiveAt: req.EffectiveAt,
		SimType:     pbBaseRate.SimType(pbBaseRate.SimType_value[req.SimType]),
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) getBaseRatesHandler(c *gin.Context, req *GetBaseRatesRequest) (*pbBaseRate.GetBaseRatesResponse, error) {
	country := c.Param("Country")
	provider := c.Param("Provider")
	effectiveAt := c.Param("EffectiveAt")
	simType := c.Param("SimType")

	to, err := strconv.ParseUint(c.Param("To"), 10, 64)
	if err != nil {
		logrus.Error(err)
	}
	from, err := strconv.ParseUint(c.Param("From"), 10, 64)
	if err != nil {
		logrus.Error(err)
	}

	resp, err := p.clients.d.GetBaseRates(&pbBaseRate.GetBaseRatesRequest{
		Country:     country,
		Provider:    provider,
		To:          to,
		From:        from,
		SimType:     pbBaseRate.SimType(pbBaseRate.SimType_value[simType]),
		EffectiveAt: effectiveAt,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) getPackageHandler(c *gin.Context, req *GetPackagesRequest) (*pb.GetPackagesResponse, error) {

	_id, err := strconv.ParseUint(c.Param("package"), 10, 64)
	if err != nil {
		logrus.Error(err)
	}

	resp, err := p.clients.d.GetPackage(&pb.GetPackagesRequest{
		Id:    _id,
		OrgId: ORG_ID,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) deletePackageHandler(c *gin.Context, req *DeletePackageRequest) (*pb.DeletePackageResponse, error) {
	_id, err := strconv.ParseUint(c.Param("package"), 10, 64)
	if err != nil {
		logrus.Error(err)
	}
	resp, err := p.clients.d.DeletePackage(&pb.DeletePackageRequest{
		Id:    _id,
		OrgId: ORG_ID,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return resp, nil
}
func (p *Router) UpdatePackageHandler(c *gin.Context, req *UpdatePackageRequest) (*pb.GetPackagesResponse, error) {
	_, err := p.clients.d.UpdatePackage(&pb.UpdatePackageRequest{
		Id:          req.Id,
		Name:        req.Name,
		SimType:     pb.SimType(pb.SimType_value[req.SimType]),
		Active:      req.Active,
		Duration:    req.Duration,
		SmsVolume:   req.SmsVolume,
		DataVolume:  req.DataVolume,
		VoiceVolume: req.VoiceVolume,
		OrgRatesId:  req.OrgRatesId,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	res, err := p.clients.d.GetPackage(&pb.GetPackagesRequest{
		Id:    req.Id,
		OrgId: ORG_ID,
	})
	if err != nil {
		logrus.Errorf("package with %d does not exist %s", req.Id, err.Error())
		return nil, err
	}

	return res, nil

}
func (p *Router) AddPackageHandler(c *gin.Context, req *AddPackageRequest) (*pb.AddPackageResponse, error) {
	return p.clients.d.AddPackage(&pb.AddPackageRequest{
		Name:        req.Name,
		OrgId:       req.OrgId,
		Duration:    req.Duration,
		OrgRatesId:  req.OrgRatesId,
		VoiceVolume: req.VoiceVolume,
		Active:      req.Active,
		DataVolume:  req.DataVolume,
		SmsVolume:   req.SmsVolume,
		SimType:     pb.SimType(pb.SimType_value[req.SimType]),
	})
}
