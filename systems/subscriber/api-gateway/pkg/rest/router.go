package rest

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/client"
	"github.com/wI2L/fizz/openapi"

	subRegPb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	simMangPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	simPoolPb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
)

const SUBS_URL_PARAMETER = "subscriber"

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
	sp  simPool
	sm  simManager
	sub subscriber
}

type simPool interface {
	Get(iccid string) (*simPoolPb.GetByIccidResponse, error)
	GetStats(simType string) (*simPoolPb.GetStatsResponse, error)
	AddSimsToSimPool(req *simPoolPb.AddRequest) (*simPoolPb.AddResponse, error)
	UploadSimsToSimPool(req *simPoolPb.UploadRequest) (*simPoolPb.UploadResponse, error)
	DeleteSimFromSimPool(id []uint64) (*simPoolPb.DeleteResponse, error)
}

type simManager interface {
	AllocateSim(req *simMangPb.AllocateSimRequest) (*simMangPb.AllocateSimResponse, error)
	GetSim(simId string) (*simMangPb.GetSimResponse, error)
	GetSimsBySub(subscriberId string) (*simMangPb.GetSimsBySubscriberResponse, error)
	GetSimsByNetwork(networkId string) (*simMangPb.GetSimsByNetworkResponse, error)
	ToggleSimStatus(simId string, status string) (*simMangPb.ToggleSimStatusResponse, error)
	AddPackageToSim(req *simMangPb.AddPackageRequest) (*simMangPb.AddPackageResponse, error)
	RemovePackageForSim(req *simMangPb.RemovePackageRequest) (*simMangPb.RemovePackageResponse, error)
	DeleteSim(simId string) (*simMangPb.DeleteSimResponse, error)
	GetPackagesForSim(simId string) (*simMangPb.GetPackagesBySimResponse, error)
	SetActivePackageForSim(req *simMangPb.SetActivePackageRequest) (*simMangPb.SetActivePackageResponse, error)
}

type subscriber interface {
	GetSubscriber(sid string) (*subRegPb.GetSubscriberResponse, error)
	AddSubscriber(req *subRegPb.AddSubscriberRequest) (*subRegPb.AddSubscriberResponse, error)
	DeleteSubscriber(sid string) (*subRegPb.DeleteSubscriberResponse, error)
	UpdateSubscriber(subscriber *subRegPb.UpdateSubscriberRequest) (*subRegPb.UpdateSubscriberResponse, error)
	GetByNetwork(networkId string) (*subRegPb.GetByNetworkResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.sp = client.NewSimPool(endpoints.SimPool, endpoints.Timeout)
	c.sm = client.NewSimManager(endpoints.SimManager, endpoints.Timeout)
	c.sub = client.NewRegistry(endpoints.Registry, endpoints.Timeout)
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
	v1 := r.f.Group("/v1", "subscriber system ", "Subscriber system version v1")
	/* These two API will be available based on RBAC */
	v1.GET("/subscribers/networks/:network_id", formatDoc("List all subscribers for a Network", ""), tonic.Handler(r.getSubscriberByNetwork, http.StatusOK))
	v1.GET("/sims/networks/:network_id", formatDoc("List all sims for a Network", ""), tonic.Handler(r.getSimsByNetwork, http.StatusOK))

	pool := v1.Group("/simpool", "SIM Pool", "SIM store for Org")
	pool.GET("/sim/:iccid", formatDoc("Get SIM by Iccid", ""), tonic.Handler(r.getSimByIccid, http.StatusOK))
	pool.GET("/stats/:sim_type", formatDoc("Get SIM Pool stats", ""), tonic.Handler(r.getSimPoolStats, http.StatusOK))
	pool.PUT("", formatDoc("Add new SIM to SIM pool", ""), tonic.Handler(r.addSimsToSimPool, http.StatusCreated))
	pool.PUT("/upload", formatDoc("Upload CSV file to add new sim to SIM Pool", ""), tonic.Handler(r.uploadSimsToSimPool, http.StatusCreated))
	pool.DELETE("/sim/:sim_id", formatDoc("Remove SIM from SIM Pool", ""), tonic.Handler(r.deleteSimFromSimPool, http.StatusOK))

	subscriber := v1.Group("/subscriber", "Subscriber", "Orgs Subscriber database")
	subscriber.GET("/:subscriber_id", formatDoc("Get subscriber by id", ""), tonic.Handler(r.getSubscriber, http.StatusOK))
	subscriber.PUT("", formatDoc("Add a new subscriber", ""), tonic.Handler(r.putSubscriber, http.StatusOK))
	subscriber.DELETE("/:subscriber_id", formatDoc("Delete a subscriber", ""), tonic.Handler(r.deleteSubscriber, http.StatusOK))
	subscriber.PATCH("/:subscriber_id", formatDoc("Update a subscriber", ""), tonic.Handler(r.updateSubscriber, http.StatusOK))

	sim := v1.Group("/sim", "SIM", "Orgs SIM data base")
	sim.GET("/:sim_id", formatDoc("Get SIM by Id", ""), tonic.Handler(r.getSim, http.StatusOK))
	sim.GET("/subscriber/:subscriber_id", formatDoc("Get a SIMs of the subscriber by Subscriber Id", ""), tonic.Handler(r.getSimsBySub, http.StatusOK))
	sim.GET("/packages/:sim_id", formatDoc("Get packages for sim", ""), tonic.Handler(r.getPackagesForSim, http.StatusOK))
	sim.POST("/package", formatDoc("Add a new package to the subscriber's sim", ""), tonic.Handler(r.addPkgForSim, http.StatusOK))
	sim.POST("/", formatDoc("Allocate a new sim to subscriber", ""), tonic.Handler(r.allocateSim, http.StatusOK))
	sim.PATCH("/:sim_id", formatDoc("Activate/Deactivate sim of subscriber", ""), tonic.Handler(r.updateSimStatus, http.StatusOK))
	sim.PATCH("/:sim_id/package/:package_id", formatDoc("Set active package for sim", ""), tonic.Handler(r.setActivePackageForSim, http.StatusOK))
	sim.DELETE("/:sim_id/package/:package_id", formatDoc("Delete a package from subscriber's sim", ""), tonic.Handler(r.removePkgForSim, http.StatusOK))
	sim.DELETE("/:sim_id", formatDoc("Delete the SIM for the subscriber", ""), tonic.Handler(r.deleteSim, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getSimByIccid(c *gin.Context, req *SimByIccidReq) (*simPoolPb.GetByIccidResponse, error) {
	resp, err := r.clients.sp.Get(req.Iccid)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Router) getSimPoolStats(c *gin.Context, req *SimPoolStatByTypeReq) (*simPoolPb.GetStatsResponse, error) {
	resp, err := r.clients.sp.GetStats(req.SimType)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Router) addSimsToSimPool(c *gin.Context, req *SimPoolAddSimReq) (*simPoolPb.AddResponse, error) {
	pbreq, err := addReqToAddSimReqPb(req)
	if err != nil {
		return nil, err
	}
	pbResp, err := r.clients.sp.AddSimsToSimPool(pbreq)
	if err != nil {
		return nil, err
	}
	return pbResp, nil
}

func (r *Router) uploadSimsToSimPool(c *gin.Context, req *SimPoolUploadSimReq) (*simPoolPb.UploadResponse, error) {

	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		log.Fatal("error:", err)
	}

	pbResp, err := r.clients.sp.UploadSimsToSimPool(&simPoolPb.UploadRequest{
		SimData: data,
		SimType: req.SimType,
	})

	if err != nil {
		return nil, err
	}

	return pbResp, nil
}

func (r *Router) deleteSimFromSimPool(c *gin.Context, req *SimPoolRemoveSimReq) (*simPoolPb.DeleteResponse, error) {
	res, err := r.clients.sp.DeleteSimFromSimPool(req.Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) getSubscriber(c *gin.Context, req *SubscriberGetReq) (*subRegPb.GetSubscriberResponse, error) {
	subsId := req.SubscriberId

	pbResp, err := r.clients.sub.GetSubscriber(subsId)
	if err != nil {
		return nil, err
	}

	return pbResp, nil
}

func (r *Router) putSubscriber(c *gin.Context, req *SubscriberAddReq) (*subRegPb.AddSubscriberResponse, error) {

	pbResp, err := r.clients.sub.AddSubscriber(&subRegPb.AddSubscriberRequest{
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		Email:                 req.Email,
		PhoneNumber:           req.Phone,
		Dob:                   req.Dob,
		Address:               req.Address,
		ProofOfIdentification: req.ProofOfIdentification,
		IdSerial:              req.IdSerial,
		NetworkId:             req.NetworkId,
		Gender:                req.Gender,
		OrgId:                 req.OrgId,
	})

	if err != nil {
		return nil, err
	}

	return pbResp, nil
}

func (r *Router) deleteSubscriber(c *gin.Context, req *SubscriberDeleteReq) (*subRegPb.DeleteSubscriberResponse, error) {
	res, err := r.clients.sub.DeleteSubscriber(req.SubscriberId)

	return res, err
}

func (r *Router) updateSubscriber(c *gin.Context, req *SubscriberUpdateReq) (*subRegPb.UpdateSubscriberResponse, error) {

	res, err := r.clients.sub.UpdateSubscriber(&subRegPb.UpdateSubscriberRequest{
		SubscriberId:          req.SubscriberId,
		Email:                 req.Email,
		PhoneNumber:           req.Phone,
		Address:               req.Address,
		ProofOfIdentification: req.ProofOfIdentification,
		IdSerial:              req.IdSerial,
	})

	return res, err
}

func (r *Router) getSubscriberByNetwork(c *gin.Context, req *SubscriberByNetworkReq) (*subRegPb.GetByNetworkResponse, error) {
	subs, err := r.clients.sub.GetByNetwork(req.NetworkId)

	return subs, err
}

func (r *Router) allocateSim(c *gin.Context, req *AllocateSimReq) (*simMangPb.AllocateSimResponse, error) {
	simReq := simMangPb.AllocateSimRequest{
		SubscriberId: req.SubscriberId,
		SimToken:     req.SimToken,
		PackageId:    req.PackageId,
		NetworkId:    req.NetworkId,
		SimType:      req.SimType,
	}
	res, err := r.clients.sm.AllocateSim(&simReq)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) getSim(c *gin.Context, req *SimReq) (*simMangPb.GetSimResponse, error) {
	res, err := r.clients.sm.GetSim(req.SimId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) getSimsBySub(c *gin.Context, req *GetSimsBySubReq) (*simMangPb.GetSimsBySubscriberResponse, error) {
	res, err := r.clients.sm.GetSimsBySub(req.SubscriberId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) getSimsByNetwork(c *gin.Context, req *SimByNetworkReq) (*simMangPb.GetSimsByNetworkResponse, error) {
	res, err := r.clients.sm.GetSimsByNetwork(req.NetworkId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) updateSimStatus(c *gin.Context, req *ActivateDeactivateSimReq) (*simMangPb.ToggleSimStatusResponse, error) {
	res, err := r.clients.sm.ToggleSimStatus(req.SimId, req.Status)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) addPkgForSim(c *gin.Context, req *AddPkgToSimReq) error {
	payload := simMangPb.AddPackageRequest{
		SimId:     req.SimId,
		PackageId: req.PackageId,
		StartDate: req.StartDate,
	}
	_, err := r.clients.sm.AddPackageToSim(&payload)
	if err != nil {
		return err
	}
	return nil
}

func (r *Router) removePkgForSim(c *gin.Context, req *RemovePkgFromSimReq) error {
	payload := simMangPb.RemovePackageRequest{
		SimId:     req.SimId,
		PackageId: req.PackageId,
	}
	_, err := r.clients.sm.RemovePackageForSim(&payload)
	if err != nil {
		return err
	}
	return nil
}

func (r *Router) deleteSim(c *gin.Context, req *SimReq) (*simMangPb.DeleteSimResponse, error) {
	res, err := r.clients.sm.DeleteSim(req.SimId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) getPackagesForSim(c *gin.Context, req *SimReq) (*simMangPb.GetPackagesBySimResponse, error) {

	res, err := r.clients.sm.GetPackagesForSim(req.SimId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) setActivePackageForSim(c *gin.Context, req *SetActivePackageForSimReq) (*simMangPb.SetActivePackageResponse, error) {
	payload := simMangPb.SetActivePackageRequest{
		SimId:     req.SimId,
		PackageId: req.PackageId,
	}
	resp, err := r.clients.sm.SetActivePackageForSim(&payload)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func addReqToAddSimReqPb(req *SimPoolAddSimReq) (*simPoolPb.AddRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid add request")
	}

	list := make([]*simPoolPb.AddSim, len(req.SimInfo))
	for i, iter := range req.SimInfo {
		list[i] = &simPoolPb.AddSim{
			Iccid:          iter.Iccid,
			Msisdn:         iter.Msidn,
			ActivationCode: iter.ActivationCode,
			IsPhysical:     iter.IsPhysicalSim,
			QrCode:         iter.QrCode,
			SmDpAddress:    iter.SmDpAddress,
			SimType:        iter.SimType,
		}
	}
	pbReq := &simPoolPb.AddRequest{
		Sim: list,
	}

	return pbReq, nil
}
