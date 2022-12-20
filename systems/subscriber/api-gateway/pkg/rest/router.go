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
	sp simPool
	sm simManager
}

type simPool interface {
	GetSimPoolStats(req *pb.GetStatsRequest) (*pb.GetStatsResponse, error)
	AddSimsToSimPool(req *pb.AddRequest) (*pb.AddResponse, error)
	UploadSimsToSimPool(req *pb.UploadRequest) (*pb.UploadResponse, error)
	DeleteSimFromSimPool(req *pb.DeleteRequest) (*pb.DeleteResponse, error)
}

type simManager interface {
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.sp = client.NewSimPool(endpoints.SimPool, endpoints.Timeout)
	c.sm = client.NewSimManager(endpoints.SimManager, endpoints.Timeout)
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
	const subs = "/subscriber/" + ":" + SUBS_URL_PARAMETER

	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "subscriber system ", "subscriber system version v1")
	/* These two API will be available based on RBAC */
	v1.GET("/susbscribers/networks/:networkid", formatDoc("List all subscibers for a Network", ""), tonic.Handler(r.getAllSubscribers, http.StatusOK))
	v1.GET("/sims/networks/:networkid", formatDoc("List all sims for a Network", ""), tonic.Handler(r.getAllSims, http.StatusOK))

	pool := v1.Group("/simpool", "SIM Pool", "SIM store for Org")
	pool.GET("/stats", formatDoc("Get SIM Pool stats", ""), tonic.Handler(r.getSimPoolStats, http.StatusOK))
	pool.PUT("/sim", formatDoc("Add new SIM to SIM pool", ""), tonic.Handler(r.addSimsToSimPool, http.StatusCreated))
	pool.PUT("/sim/upload", formatDoc("Upload CSV file to add new sim to SIM Pool", ""), tonic.Handler(r.uploadSimsToSimPool, http.StatusCreated))
	pool.DELETE("/sim", formatDoc("Remove SIM from SIM Pool", ""), tonic.Handler(r.deleteSimFromSimPool, http.StatusOK))

	subscriber := v1.Group(subs, "Subscriber", "Orgs Subscriber database")
	subscriber.GET("", formatDoc("Get System credential for Org", ""), tonic.Handler(r.getSubscriber, http.StatusOK))
	subscriber.PUT("", formatDoc("Add a new subsciber", ""), tonic.Handler(r.putSubscriber, http.StatusOK))
	subscriber.DELETE("", formatDoc("Delete a subsciber", ""), tonic.Handler(r.deleteSubscriber, http.StatusOK))

	sim := subscriber.Group("/sim", "SIM", "Orgs SIM data base")
	sim.GET("/:sim", formatDoc("Get a SIM of the subscriber by SIM Id", ""), tonic.Handler(r.getSim, http.StatusOK))
	sim.POST("/:sim", formatDoc("Allocate  anew sim to subsciber", ""), tonic.Handler(r.allocateSim, http.StatusOK))
	sim.POST("/:sim/package", formatDoc("Add a new package to the subscriber's sim", ""), tonic.Handler(r.postPackageToSim, http.StatusOK))
	sim.DELETE("/:sim/package", formatDoc("Delete a package from subscriber's sim", ""), tonic.Handler(r.deletePackageForSim, http.StatusOK))
	sim.DELETE("/:sim", formatDoc("Delete the SIM for the subscriber", ""), tonic.Handler(r.deleteSim, http.StatusOK))
	sim.PATCH("/:sim", formatDoc("Update the SIM state to active/inactive for the subscriber's sim", ""), tonic.Handler(r.patchSim, http.StatusOK))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getAllSubscribers(c *gin.Context, req *SubscriberListReq) (*SubscriberListResp, error) {
	return nil, nil
}

func (r *Router) getAllSims(c *gin.Context, req *SimListReq) (*SimListResp, error) {
	return nil, nil
}

func (r *Router) getSimPoolStats(c *gin.Context, req *SimPoolStatByTypeReq) (*SimPoolStats, error) {
	resp, err := r.clients.sp.GetSimPoolStats(&pb.GetStatsRequest{
		SimType: pb.SimType_inter_mno_all,
	})
	if err != nil {
		return nil, err
	}

	return &SimPoolStats{
		Total:     uint64(resp.Total),
		Available: uint64(resp.Available),
		Consumed:  uint64(resp.Available),
		Failed:    uint64(resp.Failed),
	}, nil
}

func (r *Router) addSimsToSimPool(c *gin.Context, req *SimPoolAddSimReq) (*SimPoolAddSimResp, error) {

	pbreq, err := addReqToAddSimReqPb(req)
	if err != nil {
		return nil, err
	}

	pbResp, err := r.clients.sp.AddSimsToSimPool(pbreq)
	if err != nil {
		return nil, err
	}

	resp, err := addSimRespTojsonSimArray(pbResp.Iccid)
	if err != nil {
		return nil, err
	}

	return &SimPoolAddSimResp{
		Iccids: resp,
	}, nil
}

func (r *Router) uploadSimsToSimPool(c *gin.Context, req *SimPoolUploadSimReq) (*SimPoolUploadSimResp, error) {

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	if err != nil {
		return nil, err
	}

	// CSV file itself not URL
	pbResp, err := r.clients.sp.UploadSimsToSimPool(&pb.UploadRequest{
		SimPoolData: bodyBytes,
	})
	if err != nil {
		return nil, err
	}

	resp, err := addSimRespTojsonSimArray(pbResp.Iccid)
	if err != nil {
		return nil, err
	}

	return &SimPoolUploadSimResp{
		Iccids: resp,
	}, nil

}

func (r *Router) deleteSimFromSimPool(c *gin.Context, req *SimPoolRemoveSimReq) (*SimPoolRemoveSimResp, error) {

	pbResp, err := r.clients.sp.DeleteSimFromSimPool(&pb.DeleteRequest{
		Iccid: req.Iccids,
	})
	if err != nil {
		return nil, err
	}

	return &SimPoolRemoveSimResp{
		Iccids: pbResp.Iccid,
	}, nil
}

func (r *Router) getSubscriber(c *gin.Context, req *SubscriberGetReq) (*SubscriberGetResp, error) {
	return nil, nil
}

func (r *Router) putSubscriber(c *gin.Context, req *SubscriberAddReq) (*SubscriberAddResp, error) {
	return nil, nil
}

func (r *Router) deleteSubscriber(c *gin.Context, req *SubscriberDeleteReq) error {
	return nil
}

func (r *Router) getSim(c *gin.Context, req *SubscriberSimReadReq) (*SubscriberSimReadResp, *error) {
	return nil, nil
}

func (r *Router) allocateSim(c *gin.Context, req *SubscriberSimAllocateReq) (*SubscriberSimAllocateResp, error) {
	return nil, nil
}

func (r *Router) postPackageToSim(c *gin.Context, req *SubscriberSimAddPackageReq) error {
	return nil
}

func (r *Router) deletePackageForSim(c *gin.Context, req *SubscriberSimRemovePackageReq) error {
	return nil
}

func (r *Router) deleteSim(c *gin.Context, req *SubscriberSimDeleteReq) error {
	return nil
}

func (r *Router) patchSim(c *gin.Context, req *SubscriberSimUpdateStateReq) error {
	return nil
}

func addReqToAddSimReqPb(req *SimPoolAddSimReq) (*pb.AddRequest, error) {

	if req == nil {
		return nil, fmt.Errorf("invalid add request")
	}

	reqSimInfo := make([]*pb.SimInfo, len(req.SimInfo))

	for idx, iter := range req.SimInfo {
		reqSimInfo[idx].Iccid = iter.Iccid
		reqSimInfo[idx].Msisdn = iter.Msidn
		// simInfo := pb.AddSimPool{
		// 	Iccid: iter.Iccid,
		// 	Msisdn: iter.Msidn,
		// }
		//reqSimInfo[idx].SimType= iter.SimType //TODO: tarnslate to enum or use some util string to enum
		//reqSimInfo[idx] = &simInfo
	}

	pbReq := &pb.AddRequest{
		SimInfo: reqSimInfo,
	}

	return pbReq, nil
}

func addSimRespTojsonSimArray(req []string) ([]string, error) {

	resp := make([]string, len(req))

	for idx, iter := range req {
		resp[idx] = iter
	}

	return resp, nil

}
