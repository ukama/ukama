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
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
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
	p packageS
	r rates
	b baserate
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

type baserate interface {
	GetBaseRatesById(req *bpb.GetBaseRatesByIdRequest) (*bpb.GetBaseRatesByIdResponse, error)
	GetBaseRatesByCountry(req *bpb.GetBaseRatesByCountryRequest) (*bpb.GetBaseRatesResponse, error)
	GetBaseRatesHistoryByCountry(req *bpb.GetBaseRatesByCountryRequest) (*bpb.GetBaseRatesResponse, error)
	GetBaseRatesForPeriod(req *bpb.GetBaseRatesByPeriodRequest) (*bpb.GetBaseRatesResponse, error)
	GetBaseRatesForPackage(req *bpb.GetBaseRatesByPeriodRequest) (*bpb.GetBaseRatesResponse, error)
	UploadBaseRates(req *bpb.UploadBaseRatesRequest) (*bpb.UploadBaseRatesResponse, error)
}
type packageS interface {
	AddPackage(req *pb.AddPackageRequest) (*pb.AddPackageResponse, error)
	UpdatePackage(req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error)
	GetPackage(id string) (*pb.GetPackageResponse, error)
	GetPackageDetails(id string) (*pb.GetPackageResponse, error)
	GetPackageByOrg(orgId string) (*pb.GetByOrgPackageResponse, error)
	DeletePackage(id string) (*pb.DeletePackageResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.p = client.NewPackageClient(endpoints.Package, endpoints.Timeout)
	c.b = client.NewBaseRateClient(endpoints.Baserate, endpoints.Timeout)
	c.r = client.NewRateClient(endpoints.Rate, endpoints.Timeout)

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

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.config.serverConf.Port)
	err := r.f.Engine().Run(fmt.Sprint(":", r.config.serverConf.Port))
	if err != nil {
		logrus.Error(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect="+REDIRECT_URI)
	auth := r.f.Group("/v1", "API Gateway", "Data-plan system version v1", func(ctx *gin.Context) {
		s := fmt.Sprintf("%s, %s, %s", "ukama", pkg.SystemName, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
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
		baseRates := auth.Group("/baserates", "BaseRates", "BaseRates operations")
		baseRates.GET("/:base_rate", formatDoc("Get BaseRate", ""), tonic.Handler(r.getBaseRateHandler, http.StatusOK))
		baseRates.POST("/upload", formatDoc("Upload baseRates", ""), tonic.Handler(r.uploadBaseRateHandler, http.StatusCreated))
		baseRates.POST("/country/:country", formatDoc("Get BaseRate", ""), tonic.Handler(r.getBaseRateByCountryHandler, http.StatusOK))
		baseRates.POST("/country/:country/history", formatDoc("Get BaseRate", ""), tonic.Handler(r.getBaseRateHistoryByCountryHandler, http.StatusOK))
		baseRates.POST("/country/:country/period", formatDoc("Get BaseRate", ""), tonic.Handler(r.getBaseRateForPeriodHandler, http.StatusOK))
		baseRates.POST("/country/:country/package", formatDoc("Get BaseRate for package", ""), tonic.Handler(r.getBaseRateForPackageHandler, http.StatusOK))

		packages := auth.Group("/packages", "Packages", "Packages operations")
		packages.POST("", formatDoc("Add Package", ""), tonic.Handler(r.AddPackageHandler, http.StatusCreated))
		packages.GET("/org/:org_id", formatDoc("Get packages of org", ""), tonic.Handler(r.getPackagesHandler, http.StatusOK))
		packages.GET("/:uuid", formatDoc("Get package", ""), tonic.Handler(r.getPackageHandler, http.StatusOK))
		packages.GET("/:uuid/details", formatDoc("Get package details", ""), tonic.Handler(r.getPackageDetailsHandler, http.StatusOK))
		packages.PATCH("/:uuid", formatDoc("Update Package", ""), tonic.Handler(r.UpdatePackageHandler, http.StatusOK))
		packages.DELETE("/:uuid", formatDoc("Delete Package", ""), tonic.Handler(r.deletePackageHandler, http.StatusOK))

		rates := auth.Group("/rates", "Rates", "Get rates for a user")
		rates.POST("/users/:user_id/rate", formatDoc("Get Rate for user", ""), tonic.Handler(r.getRateHandler, http.StatusOK))

		markup := auth.Group("/markup", "Rates", "Get rates for a user and set markup percentages for user")
		markup.GET("/users/:user_id", formatDoc("get markup percentage for user", ""), tonic.Handler(r.getMarkupHandler, http.StatusOK))
		markup.DELETE("/users/:user_id", formatDoc("delete markup percentage for user", ""), tonic.Handler(r.deleteMarkupHandler, http.StatusOK))
		markup.POST(":markup/users/:user_id", formatDoc("set markup percentage for user", ""), tonic.Handler(r.setMarkupHandler, http.StatusCreated))
		markup.GET("/users/:user_id/history", formatDoc("get markup percentage history", ""), tonic.Handler(r.getMarkupHistory, http.StatusOK))

		markup.POST("/:markup/default", formatDoc("set default markup percentage", ""), tonic.Handler(r.setDefaultMarkupHandler, http.StatusCreated))
		markup.GET("/default", formatDoc("get default markup percentage", ""), tonic.Handler(r.getDefaultMarkupHandler, http.StatusOK))
		markup.GET("/default/history", formatDoc("get default markup percentage history", ""), tonic.Handler(r.getDefaultMarkupHistory, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getPackagesHandler(c *gin.Context, req *GetPackageByOrgRequest) (*pb.GetByOrgPackageResponse, error) {
	resp, err := r.clients.p.GetPackageByOrg(req.OrgId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) getBaseRateHandler(c *gin.Context, req *GetBaseRateRequest) (*bpb.GetBaseRatesByIdResponse, error) {
	resp, err := r.clients.b.GetBaseRatesById(&bpb.GetBaseRatesByIdRequest{
		Uuid: req.RateId,
	})

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) getBaseRateByCountryHandler(c *gin.Context, req *GetBaseRatesByCountryRequest) (*bpb.GetBaseRatesResponse, error) {
	resp, err := r.clients.b.GetBaseRatesByCountry(&bpb.GetBaseRatesByCountryRequest{
		Country:  req.Country,
		Provider: req.Provider,
		SimType:  req.SimType,
	})

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) getBaseRateHistoryByCountryHandler(c *gin.Context, req *GetBaseRatesByCountryRequest) (*bpb.GetBaseRatesResponse, error) {
	resp, err := r.clients.b.GetBaseRatesHistoryByCountry(&bpb.GetBaseRatesByCountryRequest{
		Country:  req.Country,
		Provider: req.Provider,
		SimType:  req.SimType,
	})

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) getBaseRateForPeriodHandler(c *gin.Context, req *GetBaseRatesForPeriodRequest) (*bpb.GetBaseRatesResponse, error) {
	resp, err := r.clients.b.GetBaseRatesForPeriod(&bpb.GetBaseRatesByPeriodRequest{
		Country:  req.Country,
		Provider: req.Provider,
		SimType:  req.SimType,
		From:     req.From,
		To:       req.To,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) getBaseRateForPackageHandler(c *gin.Context, req *GetBaseRatesForPeriodRequest) (*bpb.GetBaseRatesResponse, error) {
	resp, err := r.clients.b.GetBaseRatesForPackage(&bpb.GetBaseRatesByPeriodRequest{
		Country:  req.Country,
		Provider: req.Provider,
		SimType:  req.SimType,
		From:     req.From,
		To:       req.To,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) uploadBaseRateHandler(c *gin.Context, req *UploadBaseRatesRequest) (*bpb.UploadBaseRatesResponse, error) {

	resp, err := r.clients.b.UploadBaseRates(&bpb.UploadBaseRatesRequest{
		FileURL:     req.FileURL,
		EffectiveAt: req.EffectiveAt,
		EndAt:       req.EndAt,
		SimType:     req.SimType,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) getPackageHandler(c *gin.Context, req *PackagesRequest) (*pb.GetPackageResponse, error) {
	resp, err := r.clients.p.GetPackage(req.Uuid)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) getPackageDetailsHandler(c *gin.Context, req *PackagesRequest) (*pb.GetPackageResponse, error) {
	resp, err := r.clients.p.GetPackageDetails(req.Uuid)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}

func (r *Router) deletePackageHandler(c *gin.Context, req *PackagesRequest) (*pb.DeletePackageResponse, error) {
	resp, err := r.clients.p.DeletePackage(req.Uuid)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return nil, err

	}
	return resp, nil
}

func (r *Router) UpdatePackageHandler(c *gin.Context, req *UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	resp, err := r.clients.p.UpdatePackage(&pb.UpdatePackageRequest{
		Uuid:   req.Uuid,
		Name:   req.Name,
		Active: req.Active,
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
		OwnerId:     req.OwnerId,
		From:        req.From,
		To:          req.To,
		Baserate:    req.BaserateId,
		VoiceVolume: req.VoiceVolume,
		Active:      req.Active,
		DataVolume:  req.DataVolume,
		SmsVolume:   req.SmsVolume,
		DataUnit:    req.DataUnit,
		VoiceUnit:   req.VoiceUnit,
		SimType:     req.SimType,
		Apn:         req.Apn,
		Markup:      req.Markup,
		Type:        req.Type,
		Flatrate:    req.Flatrate,
		Amount:      req.Amount,
	}

	return r.clients.p.AddPackage(pack)
}

func (r *Router) getRateHandler(c *gin.Context, req *GetRateRequest) (*rpb.GetRateResponse, error) {
	resp, err := r.clients.r.GetRate(&rpb.GetRateRequest{
		OwnerId:  req.OwnerId,
		Country:  req.Country,
		Provider: req.Provider,
		To:       req.To,
		From:     req.From,
		SimType:  req.SimType,
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
