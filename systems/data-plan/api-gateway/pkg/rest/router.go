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
	baseRates.POST("", formatDoc("Upload baseRates", ""), tonic.Handler(r.uploadBaseRateHandler, http.StatusCreated))
	baseRates.GET("/:baseRate", formatDoc("Get BaseRate", ""), tonic.Handler(r.getBaseRateHandler, http.StatusOK))
	baseRates.GET("", formatDoc("Get BaseRates by country", ""), tonic.Handler(r.getBaseRatesHandler, http.StatusOK))

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

	var reqBody GetBaseRateRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Error(err)
	}
	resp, err := p.clients.d.GetBaseRate(&pbBaseRate.GetBaseRateRequest{
		RateId: reqBody.RateId,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) uploadBaseRateHandler(c *gin.Context, req *UploadBaseRatesRequest) (*pbBaseRate.UploadBaseRatesResponse, error) {
	var reqBody UploadBaseRatesRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Error(err)
	}
	resp, err := p.clients.d.UploadBaseRates(&pbBaseRate.UploadBaseRatesRequest{
		FileURL:     reqBody.FileURL,
		EffectiveAt: reqBody.EffectiveAt,
		SimType:     pbBaseRate.SimType(pbBaseRate.SimType_value[reqBody.SimType]),
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) getBaseRatesHandler(c *gin.Context, req *GetBaseRatesRequest) (*pbBaseRate.GetBaseRatesResponse, error) {
	provider := c.Query("Provider")
	effectiveAt := c.Query("EffectiveAt")
	simType := c.Query("SimType")
	country, ok := c.GetQuery("country")
	if !ok {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "country is a mandatory query parameter"}
	}
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
	var reqBody GetPackagesRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Error(err)
	}
	resp, err := p.clients.d.GetPackage(&pb.GetPackagesRequest{
		Id:    reqBody.Id,
		OrgId: ORG_ID,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func (p *Router) deletePackageHandler(c *gin.Context, req *DeletePackageRequest) (*pb.DeletePackageResponse, error) {
	var reqBody DeletePackageRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Error(err)
	}
	resp, err := p.clients.d.DeletePackage(&pb.DeletePackageRequest{
		Id:    reqBody.Id,
		OrgId: ORG_ID,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return resp, nil
}
func (p *Router) UpdatePackageHandler(c *gin.Context, req *UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	var reqBody UpdatePackageRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Error(err)
	}
	resp, err := p.clients.d.UpdatePackage(&pb.UpdatePackageRequest{
		Id:          reqBody.Id,
		Name:        reqBody.Name,
		SimType:     pb.SimType(pb.SimType_value[reqBody.SimType]),
		Active:      reqBody.Active,
		Duration:    reqBody.Duration,
		SmsVolume:   reqBody.SmsVolume,
		DataVolume:  reqBody.DataVolume,
		VoiceVolume: reqBody.VoiceVolume,
		OrgRatesId:  reqBody.OrgRatesId,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	
	return resp , nil

}
func (p *Router) AddPackageHandler(c *gin.Context, req *AddPackageRequest) (*pb.AddPackageResponse, error) {
	var reqBody AddPackageRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Error(err)
	}
	return p.clients.d.AddPackage(&pb.AddPackageRequest{
		Name:        reqBody.Name,
		OrgId:       reqBody.OrgId,
		Duration:    reqBody.Duration,
		OrgRatesId:  reqBody.OrgRatesId,
		VoiceVolume: reqBody.VoiceVolume,
		Active:      reqBody.Active,
		DataVolume:  reqBody.DataVolume,
		SmsVolume:   reqBody.SmsVolume,
		SimType:     pb.SimType(pb.SimType_value[reqBody.SimType]),
	})
}
