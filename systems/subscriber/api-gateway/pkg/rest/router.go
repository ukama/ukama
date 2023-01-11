package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/cmd/version"
	pb "github.com/ukama/ukama/systems/subscriber/api-gateway/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/client"
	"github.com/wI2L/fizz/openapi"
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
	GetStats(simType string) (*pb.GetStatsResponse, error)
	AddSimsToSimPool(req *pb.AddRequest) (*pb.AddResponse, error)
	UploadSimsToSimPool(req *pb.UploadRequest) (*pb.UploadResponse, error)
	DeleteSimFromSimPool(id []uint64) (*pb.DeleteResponse, error)
}

type simManager interface {
	AllocateSim(req *pb.AllocateSimRequest) (*pb.AllocateSimResponse, error)
	GetSim(simId string) (*pb.GetSimResponse, error)
	GetSimsBySub(subscriberId string) (*pb.GetSimsBySubscriberResponse, error)
	ToggleSimStatus(simId string, status string) (*pb.ToggleSimStatusResponse, error)
	AddPackageToSim(req *pb.AddPackageRequest) (*pb.AddPackageResponse, error)
	RemovePackageForSim(req *pb.RemovePackageRequest) (*pb.RemovePackageResponse, error)
	GetSimUsage(simId string) (*pb.GetSimUsageResponse, error)
	DeleteSim(simId string) (*pb.DeleteSimResponse, error)
}

type subscriber interface {
	GetSubscriber(sid string) (*pb.GetSubscriberResponse, error)
	AddSubscriber(req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error)
	DeleteSubscriber(sid string) (*pb.DeleteSubscriberResponse, error)
	UpdateSubscriber(subscriber *pb.UpdateSubscriberRequest) (*pb.UpdateSubscriberResponse, error)
	GetByNetwork(networkId string) (*pb.GetByNetworkResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.sp = client.NewSimPool(endpoints.SimPool, endpoints.Timeout)
	c.sm = client.NewSimManager(endpoints.SimManager, endpoints.Timeout)
	c.sub = client.NewSubscriberRegistry(endpoints.SubscriberRegistry, endpoints.Timeout)
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
	const subs = "/subscriber"

	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "subscriber system ", "subscriber system version v1")
	/* These two API will be available based on RBAC */
	// v1.GET("/susbscribers/networks/:networkid", formatDoc("List all subscibers for a Network", ""), tonic.Handler(r.getSubscriberByNetwork, http.StatusOK))
	// v1.GET("/sims/networks/:networkid", formatDoc("List all sims for a Network", ""), tonic.Handler(r.getAllSims, http.StatusOK))

	pool := v1.Group("/simpool", "SIM Pool", "SIM store for Org")
	pool.GET("/stats", formatDoc("Get SIM Pool stats", ""), tonic.Handler(r.getSimPoolStats, http.StatusOK))
	pool.PUT("/sim", formatDoc("Add new SIM to SIM pool", ""), tonic.Handler(r.addSimsToSimPool, http.StatusCreated))
	pool.PUT("/upload", formatDoc("Upload CSV file to add new sim to SIM Pool", ""), tonic.Handler(r.uploadSimsToSimPool, http.StatusCreated))
	pool.DELETE("/sim", formatDoc("Remove SIM from SIM Pool", ""), tonic.Handler(r.deleteSimFromSimPool, http.StatusOK))

	subscriber := v1.Group(subs, "Subscriber", "Orgs Subscriber database")
	subscriber.GET("/", formatDoc("Get System credential for Org", ""), tonic.Handler(r.getSubscriber, http.StatusOK))
	subscriber.PUT("", formatDoc("Add a new subsciber", ""), tonic.Handler(r.putSubscriber, http.StatusOK))
	subscriber.DELETE("", formatDoc("Delete a subsciber", ""), tonic.Handler(r.deleteSubscriber, http.StatusOK))
	subscriber.PATCH("", formatDoc("Update a subsciber", ""), tonic.Handler(r.updateSubscriber, http.StatusOK))
	subscriber.GET("/network/:networkId", formatDoc("List all subscibers for a Network", ""), tonic.Handler(r.getSubscriberByNetwork, http.StatusOK))

	sim := v1.Group("/sim", "SIM", "Orgs SIM data base")
	sim.GET("/", formatDoc("Get a SIM of the subscriber by SIM Id", ""), tonic.Handler(r.getSim, http.StatusOK))
	sim.GET("/subscriber", formatDoc("Get a SIMs of the subscriber by Subscriber Id", ""), tonic.Handler(r.getSimsBySub, http.StatusOK))
	sim.GET("/usage", formatDoc("Get SIM usage", ""), tonic.Handler(r.getSimUsage, http.StatusOK))
	sim.PATCH("/", formatDoc("Activate/Deactivate sim of subscriber", ""), tonic.Handler(r.updateSimStatus, http.StatusOK))
	sim.POST("/package", formatDoc("Add a new package to the subscriber's sim", ""), tonic.Handler(r.addPkgForSim, http.StatusOK))
	sim.DELETE("/package", formatDoc("Delete a package from subscriber's sim", ""), tonic.Handler(r.removePkgForSim, http.StatusOK))
	sim.POST("/", formatDoc("Allocate a new sim to subscriber", ""), tonic.Handler(r.allocateSim, http.StatusOK))
	sim.DELETE("/:sim", formatDoc("Delete the SIM for the subscriber", ""), tonic.Handler(r.deleteSim, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getAllSims(c *gin.Context, req *SimListReq) (*SimListResp, error) {
	return nil, nil
}

func (r *Router) getSimPoolStats(c *gin.Context, req *SimPoolStatByTypeReq) (*pb.GetStatsResponse, error) {
	simType, ok := c.GetQuery("simType")
	if !ok {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "simType is a mandatory query parameter"}
	}
	resp, err := r.clients.sp.GetStats(simType)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Router) addSimsToSimPool(c *gin.Context, req *SimPoolAddSimReq) (*pb.AddResponse, error) {
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

func (r *Router) uploadSimsToSimPool(c *gin.Context, req *SimPoolUploadSimReq) (*pb.UploadResponse, error) {
	data, err := ioutil.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	if err != nil {
		return nil, err
	}

	pbResp, err := r.clients.sp.UploadSimsToSimPool(&pb.UploadRequest{
		SimData: data,
	})
	if err != nil {
		return nil, err
	}

	return pbResp, nil
}

func (r *Router) deleteSimFromSimPool(c *gin.Context, req *SimPoolRemoveSimReq) (*pb.DeleteResponse, error) {
	res, err := r.clients.sp.DeleteSimFromSimPool(req.Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) getSubscriber(c *gin.Context, req *SubscriberGetReq) (*pb.GetSubscriberResponse, error) {
	subsId := req.SubscriberId

	pbResp, err := r.clients.sub.GetSubscriber(subsId)
	if err != nil {
		return nil, err
	}

	return pbResp, nil
}

func (r *Router) putSubscriber(c *gin.Context, req *SubscriberAddReq) (*pb.AddSubscriberResponse, error) {

	pbResp, err := r.clients.sub.AddSubscriber(&pb.AddSubscriberRequest{
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		Email:                 req.Email,
		PhoneNumber:           req.Phone,
		DateOfBirth:           req.DOB,
		Address:               req.Address,
		ProofOfIdentification: req.ProofOfIdentification,
		IdSerial:              req.IdSerial})
	if err != nil {
		return nil, err
	}

	return pbResp, nil
}

func (r *Router) deleteSubscriber(c *gin.Context, req *SubscriberDeleteReq) (*pb.DeleteSubscriberResponse, error) {
	res, err := r.clients.sub.DeleteSubscriber(req.SubscriberId)

	return res, err
}

func (r *Router) updateSubscriber(c *gin.Context, req *SubscriberUpdateReq) (*pb.UpdateSubscriberResponse, error) {

	res, err := r.clients.sub.UpdateSubscriber(&pb.UpdateSubscriberRequest{
		SubscriberID:          req.SubscriberId,
		Email:                 req.Email,
		PhoneNumber:           req.Phone,
		Address:               req.Address,
		ProofOfIdentification: req.ProofOfIdentification,
		IdSerial:              req.IdSerial,
	})

	return res, err
}

func (r *Router) getSubscriberByNetwork(c *gin.Context, req *SubscriberByNetworkReq) (*pb.GetByNetworkResponse, error) {
	subs, err := r.clients.sub.GetByNetwork(req.NetworkId)

	return subs, err
}

func (r *Router) allocateSim(c *gin.Context, req *AllocateSimReq) (*pb.AllocateSimResponse, error) {
	simReq := pb.AllocateSimRequest{
		SubscriberID: req.SubscriberId,
		SimToken:     req.SimToken,
		PackageID:    req.PackageId,
		NetworkID:    req.NetworkId,
	}
	res, err := r.clients.sm.AllocateSim(&simReq)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) getSim(c *gin.Context, req *SimReq) (*pb.GetSimResponse, error) {
	res, err := r.clients.sm.GetSim(req.SimId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) getSimsBySub(c *gin.Context, req *GetSimsBySubReq) (*pb.GetSimsBySubscriberResponse, error) {
	res, err := r.clients.sm.GetSimsBySub(req.SubscriberId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) updateSimStatus(c *gin.Context, req *ActivateDeactivateSimReq) (*pb.ToggleSimStatusResponse, error) {
	res, err := r.clients.sm.ToggleSimStatus(req.SimId, req.Status)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Router) addPkgForSim(c *gin.Context, req *AddPkgToSimReq) error {
	payload := pb.AddPackageRequest{
		SimID:     req.SimId,
		PackageID: req.PackageId,
		StartDate: req.StartDate,
	}
	_, err := r.clients.sm.AddPackageToSim(&payload)
	if err != nil {
		return err
	}
	return nil
}

func (r *Router) removePkgForSim(c *gin.Context, req *RemovePkgFromSimReq) error {
	payload := pb.RemovePackageRequest{
		SimID:     req.SimId,
		PackageID: req.PackageId,
	}
	_, err := r.clients.sm.RemovePackageForSim(&payload)
	if err != nil {
		return err
	}
	return nil
}

func (r *Router) getSimUsage(c *gin.Context, req *SimReq) (*pb.GetSimUsageResponse, error) {
	res, err := r.clients.sm.GetSimUsage(req.SimId)
	if err != nil {
		return nil, err
	}
	return res, nil

}
func (r *Router) deleteSim(c *gin.Context, req *SimReq) (*pb.DeleteSimResponse, error) {
	res, err := r.clients.sm.DeleteSim(req.SimId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func addReqToAddSimReqPb(req *SimPoolAddSimReq) (*pb.AddRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid add request")
	}

	list := make([]*pb.AddSim, len(req.SimInfo))
	for i, iter := range req.SimInfo {
		list[i] = &pb.AddSim{
			Iccid:          iter.Iccid,
			Msisdn:         iter.Msidn,
			ActivationCode: iter.ActivationCode,
			IsPhysicalSim:  iter.IsPhysicalSim,
			QrCode:         iter.QrCode,
			SmDpAddress:    iter.SmDpAddress,
			SimType:        pb.SimType(pb.SimType_value[iter.SimType]),
		}
	}
	pbReq := &pb.AddRequest{
		Sim: list,
	}

	return pbReq, nil
}
