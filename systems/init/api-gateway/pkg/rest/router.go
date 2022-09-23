package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/config"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/services/common/rest"
	"github.com/ukama/ukama/systems/init/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"github.com/wI2L/fizz/openapi"
)

const ORG_URL_PARAMETER = "org"

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
	l lookup
}

type lookup interface {
	AddOrg(req *pb.AddOrgRequest) (*pb.AddOrgResponse, error)
	UpdateOrg(req *pb.UpdateOrgRequest) (*pb.UpdateOrgResponse, error)
	GetOrg(req *pb.GetOrgRequest) (*pb.GetOrgResponse, error)
	AddNodeForOrg(req *pb.AddNodeRequest) (*pb.AddNodeResponse, error)
	GetNodeForOrg(req *pb.GetNodeForOrgRequest) (*pb.GetNodeResponse, error)
	DeleteNodeForOrg(req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error)
	AddSystemForOrg(req *pb.AddSystemRequest) (*pb.AddSystemResponse, error)
	UpdateSystemForOrg(req *pb.UpdateSystemRequest) (*pb.UpdateSystemResponse, error)
	GetSystemForOrg(req *pb.GetSystemRequest) (*pb.GetSystemResponse, error)
	DeleteSystemForOrg(req *pb.DeleteSystemRequest) (*pb.DeleteSystemResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.l = client.Newlookup(endpoints.Lookup, endpoints.Timeout)
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
	const org = "/orgs/" + ":" + ORG_URL_PARAMETER

	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.ServiceName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "Init system ", "Init system version v1")

	orgs := v1.Group(org, "Orgs", "looking for orgs credentials")
	orgs.GET("", formatDoc("Get Orgs Credential", ""), tonic.Handler(r.getOrgHandler, http.StatusOK))
	orgs.PUT("", formatDoc("Add Org and Credential", ""), tonic.Handler(r.putOrgHandler, http.StatusCreated))
	orgs.PATCH("", formatDoc("Update Orgs Credential", ""), tonic.Handler(r.patchOrgHandler, http.StatusOK))

	nodes := orgs.Group("/nodes", "Nodes", "Orgs credentials for Node")
	nodes.GET("/:node", formatDoc("Get Orgs credential for Node", ""), tonic.Handler(r.getNodeHandler, http.StatusOK))
	nodes.PUT("/:node", formatDoc("Add Node to Org", ""), tonic.Handler(r.putNodeHandler, http.StatusCreated))
	nodes.DELETE("/:node", formatDoc("Delete Node from Org", ""), tonic.Handler(r.deleteNodeHandler, http.StatusOK))

	systems := orgs.Group("/systems", "Systems", "Orgs System credentials")
	systems.GET("/:system", formatDoc("Get System credential for Org", ""), tonic.Handler(r.getSystemHandler, http.StatusOK))
	systems.PUT("/:system", formatDoc("Add or Update System credential for Org", ""), tonic.Handler(r.putSystemHandler, http.StatusCreated))
	systems.DELETE("/:system", formatDoc("Delete System credential for Org", ""), tonic.Handler(r.deleteSystemHandler, http.StatusOK))
	systems.PATCH("/:system", formatDoc("Update System Credential", ""), tonic.Handler(r.patchSystemHandler, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getOrgHandler(c *gin.Context, req *GetOrgRequest) (*pb.GetOrgResponse, error) {
	org := c.Param("org")

	return r.clients.l.GetOrg(&pb.GetOrgRequest{
		OrgName: org,
	})
}

func (r *Router) putOrgHandler(c *gin.Context, req *AddOrgRequest) (*pb.AddOrgResponse, error) {
	org := c.Param("org")

	return r.clients.l.AddOrg(&pb.AddOrgRequest{
		OrgName:     org,
		Certificate: req.Certificate,
		Ip:          req.Ip,
	})

}

func (r *Router) patchOrgHandler(c *gin.Context, req *UpdateOrgRequest) (*pb.UpdateOrgResponse, error) {
	org := c.Param("org")

	return r.clients.l.UpdateOrg(&pb.UpdateOrgRequest{
		OrgName:     org,
		Certificate: req.Certificate,
		Ip:          req.Ip,
	})
}

func (r *Router) putNodeHandler(c *gin.Context, req *AddNodeRequest) (*pb.AddNodeResponse, error) {
	org := c.Param("org")
	node := c.Param("node")

	return r.clients.l.AddNodeForOrg(&pb.AddNodeRequest{
		OrgName: org,
		NodeId:  node,
	})
}

func (r *Router) getNodeHandler(c *gin.Context, req *GetNodeRequest) (*pb.GetNodeResponse, error) {
	org := c.Param("org")
	node := c.Param("node")

	return r.clients.l.GetNodeForOrg(&pb.GetNodeForOrgRequest{
		OrgName: org,
		NodeId:  node,
	})
}

func (r *Router) deleteNodeHandler(c *gin.Context, req *DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	org := c.Param("org")
	node := c.Param("node")

	return r.clients.l.DeleteNodeForOrg(&pb.DeleteNodeRequest{
		OrgName: org,
		NodeId:  node,
	})
}

func (r *Router) putSystemHandler(c *gin.Context, req *AddSystemRequest) (*pb.AddSystemResponse, error) {
	org := c.Param("org")
	sys := c.Param("system")

	return r.clients.l.AddSystemForOrg(&pb.AddSystemRequest{
		OrgName:     org,
		SystemName:  sys,
		Certificate: req.Certificate,
		Ip:          req.Ip,
		Port:        req.Port,
	})

}

func (r *Router) patchSystemHandler(c *gin.Context, req *UpdateSystemRequest) (*pb.UpdateSystemResponse, error) {
	org := c.Param("org")
	sys := c.Param("system")

	return r.clients.l.UpdateSystemForOrg(&pb.UpdateSystemRequest{
		OrgName:     org,
		SystemName:  sys,
		Certificate: req.Certificate,
		Ip:          req.Ip,
		Port:        req.Port,
	})
}

func (r *Router) getSystemHandler(c *gin.Context, req *GetSystemRequest) (*pb.GetSystemResponse, error) {
	org := c.Param("org")
	sys := c.Param("system")

	return r.clients.l.GetSystemForOrg(&pb.GetSystemRequest{
		OrgName:    org,
		SystemName: sys,
	})
}

func (r *Router) deleteSystemHandler(c *gin.Context, req *DeleteSystemRequest) (*pb.DeleteSystemResponse, error) {
	org := c.Param("org")
	sys := c.Param("system")

	return r.clients.l.DeleteSystemForOrg(&pb.DeleteSystemRequest{
		OrgName:    org,
		SystemName: sys,
	})
}
