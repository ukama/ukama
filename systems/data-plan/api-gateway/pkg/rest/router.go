package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/client"
	pbBaseRate "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

var REDIRECT_URI = "https://data-plan.dev.ukama.com/swagger/#/"

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
	auth          *config.Auth
}

type Clients struct {
	d dataPlan
}

type dataPlan interface {
	AddPackage(req *pb.AddPackageRequest) (*pb.AddPackageResponse, error)
	UpdatePackage(req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error)
	GetPackage(id string) (*pb.GetPackageResponse, error)
	GetPackageByOrg(orgId string) (*pb.GetByOrgPackageResponse, error)
	DeletePackage(id string) (*pb.DeletePackageResponse, error)
	UploadBaseRates(req *pbBaseRate.UploadBaseRatesRequest) (*pbBaseRate.UploadBaseRatesResponse, error)
	GetBaseRates(req *pbBaseRate.GetBaseRatesRequest) (*pbBaseRate.GetBaseRatesResponse, error)
	GetBaseRate(id string) (*pbBaseRate.GetBaseRateResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.d = client.NewDataPlan(endpoints.Package, endpoints.Rate, endpoints.Timeout)

	return c
}

func NewRouter(clients *Clients, config *RouterConfig, authfunc func(*gin.Context, string) error) *Router {

	r := &Router{
		clients: clients,
		config:  config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init(authfunc)
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		httpEndpoints: &svcConf.HttpServices,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
		auth:          svcConf.Auth,
	}
}

func (rt *Router) Run() {
	logrus.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		logrus.Error(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect="+REDIRECT_URI)
	auth := r.f.Group("/v1", "Data-plan system ", "Data-plan  system version v1", func(ctx *gin.Context) {
		err := f(ctx, r.config.auth.AuthAPIGW)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		if err == nil {
			return
		}
	})
	auth.Use()
	{
		baseRates := auth.Group("/baseRates", "BaseRates", "BaseRates operations")
		baseRates.POST("", formatDoc("Upload baseRates", ""), tonic.Handler(r.uploadBaseRateHandler, http.StatusCreated))
		baseRates.GET("/:base_rate", formatDoc("Get BaseRate", ""), tonic.Handler(r.getBaseRateHandler, http.StatusOK))
		baseRates.GET("", formatDoc("Get BaseRates by country", ""), tonic.Handler(r.getBaseRatesHandler, http.StatusOK))

		packages := auth.Group("/packages", "Packages", "Packages operations")
		packages.POST("", formatDoc("Add Package", ""), tonic.Handler(r.AddPackageHandler, http.StatusCreated))
		packages.GET("/org/:org_id", formatDoc("Get packages of org", ""), tonic.Handler(r.getPackagesHandler, http.StatusOK))
		packages.GET("/:uuid", formatDoc("Get package", ""), tonic.Handler(r.getPackageHandler, http.StatusOK))
		packages.PATCH("/:uuid", formatDoc("Update Package", ""), tonic.Handler(r.UpdatePackageHandler, http.StatusOK))
		packages.DELETE("/:uuid", formatDoc("Delete Package", ""), tonic.Handler(r.deletePackageHandler, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
func (p *Router) getPackagesHandler(c *gin.Context, req *GetPackageByOrgRequest) (*pb.GetByOrgPackageResponse, error) {
	resp, err := p.clients.d.GetPackageByOrg(req.OrgId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) getBaseRateHandler(c *gin.Context, req *GetBaseRateRequest) (*pbBaseRate.GetBaseRateResponse, error) {
	resp, err := p.clients.d.GetBaseRate(req.RateId)
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
		SimType:     req.SimType,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) getBaseRatesHandler(c *gin.Context, req *GetBaseRatesRequest) (*pbBaseRate.GetBaseRatesResponse, error) {

	resp, err := p.clients.d.GetBaseRates(&pbBaseRate.GetBaseRatesRequest{
		Country:     req.Country,
		Provider:    req.Provider,
		To:          req.To,
		From:        req.From,
		SimType:     req.SimType,
		EffectiveAt: req.EffectiveAt,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) getPackageHandler(c *gin.Context, req *PackagesRequest) (*pb.GetPackageResponse, error) {
	resp, err := p.clients.d.GetPackage(req.Uuid)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) deletePackageHandler(c *gin.Context, req *PackagesRequest) (*pb.DeletePackageResponse, error) {
	resp, err := p.clients.d.DeletePackage(req.Uuid)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return nil, err

	}
	return resp, nil
}
func (p *Router) UpdatePackageHandler(c *gin.Context, req *UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	resp, err := p.clients.d.UpdatePackage(&pb.UpdatePackageRequest{
		Uuid:        req.Uuid,
		Name:        req.Name,
		SimType:     req.SimType,
		Active:      req.Active,
		Duration:    req.Duration,
		SmsVolume:   req.SmsVolume,
		DataVolume:  req.DataVolume,
		VoiceVolume: req.VoiceVolume,
		OrgRatesId:  req.OrgRatesId,
	})
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return nil, err
	}

	return resp, nil

}
func (p *Router) AddPackageHandler(c *gin.Context, req *AddPackageRequest) (*pb.AddPackageResponse, error) {
	pack := &pb.AddPackageRequest{
		Name:        req.Name,
		OrgId:       req.OrgId,
		Duration:    req.Duration,
		OrgRatesId:  req.OrgRatesId,
		VoiceVolume: req.VoiceVolume,
		Active:      req.Active,
		DataVolume:  req.DataVolume,
		SmsVolume:   req.SmsVolume,
		SimType:     req.SimType,
	}

	return p.clients.d.AddPackage(pack)
}
