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
	"github.com/ukama/ukama/systems/subscriber/api-gateway/cmd/version"
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

func (r *Router) getAllSubscribers(c *gin.Context) error {
	return nil
}

func (r *Router) getAllSims(c *gin.Context) error {
	return nil
}

func (r *Router) getSimPoolStats(c *gin.Context) error {
	return nil
}

func (r *Router) addSimsToSimPool(c *gin.Context) error {
	return nil
}

func (r *Router) uploadSimsToSimPool(c *gin.Context) error {
	return nil
}

func (r *Router) deleteSimFromSimPool(c *gin.Context) error {
	return nil
}

func (r *Router) getSubscriber(c *gin.Context) error {
	return nil
}

func (r *Router) putSubscriber(c *gin.Context) error {
	return nil
}

func (r *Router) deleteSubscriber(c *gin.Context) error {
	return nil
}

func (r *Router) getSim(c *gin.Context) error {
	return nil
}

func (r *Router) allocateSim(c *gin.Context) error {
	return nil
}

func (r *Router) postPackageToSim(c *gin.Context) error {
	return nil
}

func (r *Router) deletePackageForSim(c *gin.Context) error {
	return nil
}

func (r *Router) deleteSim(c *gin.Context) error {
	return nil
}

func (r *Router) patchSim(c *gin.Context) error {
	return nil
}
