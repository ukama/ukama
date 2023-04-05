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
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"github.com/wI2L/fizz"
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
	r rates
}

type rates interface {
	GetRate(req *rpb.GetRateRequest) (*rpb.GetRateResponse, error)
	UpdateDefaultMarkup(req *rpb.UpdateDefaultMarkupRequest) (*rpb.UpdateDefaultMarkupResponse, error)
	UpdateMarkup(req *rpb.UpdateMarkupRequest) (*rpb.UpdateMarkupResponse, error)
	GetDefaultMarkup(req *rpb.GetDefaultMarkupRequest) (*rpb.GetDefaultMarkupResponse, error)
	GetDefaultMarkupHistory(req *rpb.GetDefaultMarkupHistoryRequest) (*rpb.GetDefaultMarkupHistoryResponse, error)
	GetMarkupHistory(req *rpb.GetMarkupHistoryRequest) (*rpb.GetMarkupHistoryResponse, error)
	DeleteMarkup(req *rpb.DeleteMarkupRequest) (*rpb.DeleteMarkupResponse, error)
	GetMarkup(req *rpb.GetMarkupRequest) (*rpb.GetMarkupResponse, error)
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
	c.d = client.NewDataPlan(endpoints.Package, endpoints.Baserate, endpoints.Timeout)
	c.r = client.NewRateClient(endpoints.Rate, endpoints.Timeout)

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

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.config.serverConf.Port)
	err := r.f.Engine().Run(fmt.Sprint(":", r.config.serverConf.Port))
	if err != nil {
		logrus.Error(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "Data-plan system ", "Data-plan  system version v1")

	baseRates := v1.Group("/baseRates", "BaseRates", "BaseRates operations")
	baseRates.POST("", formatDoc("Upload baseRates", ""), tonic.Handler(r.uploadBaseRateHandler, http.StatusCreated))
	baseRates.GET("/:base_rate", formatDoc("Get BaseRate", ""), tonic.Handler(r.getBaseRateHandler, http.StatusOK))
	baseRates.GET("", formatDoc("Get BaseRates by country", ""), tonic.Handler(r.getBaseRatesHandler, http.StatusOK))

	packages := v1.Group("/packages", "Packages", "Packages operations")
	packages.POST("", formatDoc("Add Package", ""), tonic.Handler(r.AddPackageHandler, http.StatusCreated))
	packages.GET("/org/:org_id", formatDoc("Get packages of org", ""), tonic.Handler(r.getPackagesHandler, http.StatusOK))
	packages.GET("/:uuid", formatDoc("Get package", ""), tonic.Handler(r.getPackageHandler, http.StatusOK))
	packages.PATCH("/:uuid", formatDoc("Update Package", ""), tonic.Handler(r.UpdatePackageHandler, http.StatusOK))
	packages.DELETE("/:uuid", formatDoc("Delete Package", ""), tonic.Handler(r.deletePackageHandler, http.StatusOK))

	rates := v1.Group("/rates", "Rates", "Get rates for a user")
	rates.GET("/users/:user_id", formatDoc("Get Rate for user", ""), tonic.Handler(r.getRateHandler, http.StatusOK))

	rates.GET("/users/:user_id/markup", formatDoc("get markup rate for user", ""), tonic.Handler(r.getMarkupHandler, http.StatusOK))
	rates.DELETE("/users/:user_id/markup", formatDoc("delete markup rate for user", ""), tonic.Handler(r.deleteMarkupHandler, http.StatusOK))
	rates.POST("/users/:user_id/markup/:markup", formatDoc("set markup rate for user", ""), tonic.Handler(r.setMarkupHandler, http.StatusCreated))
	rates.GET("/users/:user_id/markup/history", formatDoc("get markup rate history", ""), tonic.Handler(r.getMarkupHistory, http.StatusOK))

	rates.POST("/default/markup/:markup", formatDoc("set default markup rate", ""), tonic.Handler(r.setDefaultMarkupHandler, http.StatusCreated))
	rates.GET("/default/markup/:markup", formatDoc("get default markup rate", ""), tonic.Handler(r.getDefaultMarkupHandler, http.StatusOK))
	rates.GET("/default/markup/history", formatDoc("get default markup rate history", ""), tonic.Handler(r.getDefaultMarkupHistory, http.StatusOK))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getPackagesHandler(c *gin.Context, req *GetPackageByOrgRequest) (*pb.GetByOrgPackageResponse, error) {
	resp, err := r.clients.d.GetPackageByOrg(req.OrgId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) getBaseRateHandler(c *gin.Context, req *GetBaseRateRequest) (*pbBaseRate.GetBaseRateResponse, error) {
	resp, err := r.clients.d.GetBaseRate(req.RateId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) uploadBaseRateHandler(c *gin.Context, req *UploadBaseRatesRequest) (*pbBaseRate.UploadBaseRatesResponse, error) {

	resp, err := r.clients.d.UploadBaseRates(&pbBaseRate.UploadBaseRatesRequest{
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

func (r *Router) getBaseRatesHandler(c *gin.Context, req *GetBaseRatesRequest) (*pbBaseRate.GetBaseRatesResponse, error) {

	resp, err := r.clients.d.GetBaseRates(&pbBaseRate.GetBaseRatesRequest{
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

func (r *Router) getPackageHandler(c *gin.Context, req *PackagesRequest) (*pb.GetPackageResponse, error) {
	resp, err := r.clients.d.GetPackage(req.Uuid)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) deletePackageHandler(c *gin.Context, req *PackagesRequest) (*pb.DeletePackageResponse, error) {
	resp, err := r.clients.d.DeletePackage(req.Uuid)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return nil, err

	}
	return resp, nil
}

func (r *Router) UpdatePackageHandler(c *gin.Context, req *UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	resp, err := r.clients.d.UpdatePackage(&pb.UpdatePackageRequest{
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

func (r *Router) AddPackageHandler(c *gin.Context, req *AddPackageRequest) (*pb.AddPackageResponse, error) {

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

	return r.clients.d.AddPackage(pack)
}

func (r *Router) getRateHandler(c *gin.Context, req *GetRateRequest) (*rpb.GetRateResponse, error) {
	resp, err := r.clients.r.GetRate(&rpb.GetRateRequest{
		OwnerId:     req.OwnerId,
		Country:     req.Country,
		Provider:    req.Provider,
		To:          req.To,
		From:        req.From,
		SimType:     req.SimType,
		EffectiveAt: req.EffectiveAt,
	})
	if err != nil {
		logrus.Errorf("Failed to get rate for user %s.Error %s", req.OwnerId, err.Error())
		return nil, err
	}

	return resp, nil

}

func (r *Router) deleteMarkupHandler(c *gin.Context, req *DeleteMarkupRequest) (*rpb.DeleteMarkupResponse, error) {
	resp, err := r.clients.r.DeleteMarkup(&rpb.DeleteMarkupRequest{
		OwnerId: req.OwnerId,
	})
	if err != nil {
		logrus.Errorf("Failed to delete markup for user %s. Error %s", req.OwnerId, err.Error())
		return nil, err
	}

	return resp, nil
}

func (r *Router) setMarkupHandler(c *gin.Context, req *SetMarkupRequest) (*rpb.UpdateMarkupResponse, error) {
	resp, err := r.clients.r.UpdateMarkup(&rpb.UpdateMarkupRequest{
		OwnerId: req.OwnerId,
		Markup:  req.Markup,
	})
	if err != nil {
		logrus.Errorf("Failed to update markup for user %s. Error %s", req.OwnerId, err.Error())
		return nil, err
	}

	return resp, nil
}

func (r *Router) getMarkupHandler(c *gin.Context, req *GetMarkupRequest) (*rpb.GetMarkupResponse, error) {
	resp, err := r.clients.r.GetMarkup(&rpb.GetMarkupRequest{
		OwnerId: req.OwnerId,
	})
	if err != nil {
		logrus.Errorf("Failed to get markup for user %s. Error %s", req.OwnerId, err.Error())
		return nil, err
	}

	return resp, nil
}

func (r *Router) getMarkupHistory(c *gin.Context, req *GetMarkupHistoryRequest) (*rpb.GetMarkupHistoryResponse, error) {
	resp, err := r.clients.r.GetMarkupHistory(&rpb.GetMarkupHistoryRequest{
		OwnerId: req.OwnerId,
	})
	if err != nil {
		logrus.Errorf("Failed to get markup history for user %s. Error %s", req.OwnerId, err.Error())
		return nil, err
	}

	return resp, nil
}

func (r *Router) setDefaultMarkupHandler(c *gin.Context, req *SetDefaultMarkupRequest) (*rpb.UpdateDefaultMarkupResponse, error) {
	resp, err := r.clients.r.UpdateDefaultMarkup(&rpb.UpdateDefaultMarkupRequest{
		Markup: req.Markup,
	})
	if err != nil {
		logrus.Errorf("Failed to update default markup. Error %s", err.Error())
		return nil, err
	}

	return resp, nil
}

func (r *Router) getDefaultMarkupHandler(c *gin.Context, req *GetDefaultMarkupRequest) (*rpb.GetDefaultMarkupResponse, error) {
	resp, err := r.clients.r.GetDefaultMarkup(&rpb.GetDefaultMarkupRequest{})
	if err != nil {
		logrus.Errorf("Failed to get default markup. Error %s", err.Error())
		return nil, err
	}

	return resp, nil
}

func (r *Router) getDefaultMarkupHistory(c *gin.Context, req *GetDefaultMarkupHistoryRequest) (*rpb.GetDefaultMarkupHistoryResponse, error) {
	resp, err := r.clients.r.GetDefaultMarkupHistory(&rpb.GetDefaultMarkupHistoryRequest{})
	if err != nil {
		logrus.Errorf("Failed to get default markup history. Error %s", err.Error())
		return nil, err
	}

	return resp, nil
}
