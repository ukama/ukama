package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/messaging/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg"
	"github.com/ukama/ukama/systems/messaging/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"github.com/wI2L/fizz/openapi"
)

const ORG_URL_PARAMETER = "org"

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig  config.Metrics
	httpEndpoints  *pkg.HttpEndpoints
	debugMode      bool
	serverConf     *rest.HttpConfig
	auth           *config.Auth
	nodeMetricPort int32
}

type Clients struct {
	n nns
}

type nns interface {
	GetNodeIpRequest(req *pb.GetNodeIPRequest) (*pb.GetNodeIPResponse, error)
	SetNodeIpRequest(req *pb.SetNodeIPRequest) (*pb.SetNodeIPResponse, error)
	DeleteNodeIpRequest(req *pb.DeleteNodeIPRequest) (*pb.DeleteNodeIPResponse, error)
	ListNodeIpRequest(req *pb.ListNodeIPRequest) (*pb.ListNodeIPResponse, error)
	GetNodeOrgMapListRequest(req *pb.NodeOrgMapListRequest) (*pb.NodeOrgMapListResponse, error)
	GetNodeIPMapListRequest(req *pb.NodeIPMapListRequest) (*pb.NodeIPMapListResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.n = client.NewNns(endpoints.Nns, endpoints.Timeout)
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
		panic(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {

	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	auth := r.f.Group("/v1", "Messaging system API Gateway", "Messaging system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			logrus.Info("Bypassing auth")
			return
		}
		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
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
		nns := auth.Group("/nns", "Nns", "Looking for node Ip address")
		nns.GET("/node/:node_id", formatDoc("Get node Ip", ""), tonic.Handler(r.getNodeHandler, http.StatusOK))
		nns.PUT("/node/:node_id", formatDoc("Add node Ip", ""), tonic.Handler(r.putNodeHandler, http.StatusCreated))
		nns.DELETE("/node/:node_id", formatDoc("Remove node from dns", ""), tonic.Handler(r.deleteNodeIPHandler, http.StatusOK))
		nns.GET("/list", formatDoc("Get all Ip's", ""), tonic.Handler(r.getAllNodeIPHandler, http.StatusOK))
		nns.GET("/map", formatDoc("Node to Org map", ""), tonic.Handler(r.getNodeOrgMapHandler, http.StatusOK))

		prom := auth.Group("/prometheus", "Prometheus target", "Target discovery endpoint")
		prom.GET("", formatDoc("Get target to scrape", ""), tonic.Handler(r.prometheusHandler, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) prometheusHandler(c *gin.Context) error {
	w := c.Writer
	w.Header().Set("Content-Type", "application/json")

	m := make(chan bool)
	nodeToOrg := &pb.NodeOrgMapListResponse{}
	go func() {
		var errCh error
		if nodeToOrg, errCh = r.clients.n.GetNodeOrgMapListRequest(&pb.NodeOrgMapListRequest{}); errCh != nil {
			logrus.Error("Error getting node to org/network map. Error: ", errCh)
		}
		m <- true
	}()

	l, err := r.clients.n.GetNodeIPMapListRequest(&pb.NodeIPMapListRequest{})
	if err != nil {
		logrus.Error("Error getting list of namespaces. Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	// wait for nodeToOrgNetwork mapping to finish
	<-m

	b, err := marshallTargets(l, nodeToOrg, int(r.config.nodeMetricPort))
	if err != nil {
		logrus.Error("Error marshalling targets. Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		logrus.Error("Error writing response. Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

type targets struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func marshallTargets(l *pb.NodeIPMapListResponse, nodeToOrg *pb.NodeOrgMapListResponse, nodeMetricsPort int) ([]byte, error) {
	resp := make([]targets, 0, len(l.Map))
	var dname string
	for _, v := range l.Map {
		labels := map[string]string{
			"nodeid": v.NodeId,
		}

		if m, ok := func(m *pb.NodeOrgMapListResponse, id string) (*pb.NodeOrgMap, bool) {
			for _, k := range m.Map {
				if strings.EqualFold(k.NodeId, id) {
					return k, true
				}
			}
			return nil, false
		}(nodeToOrg, v.NodeId); ok {
			dname = m.GetDomainname()
			labels["org"] = m.Org
			labels["network"] = m.Network
		}

		resp = append(resp, targets{
			Targets: []string{v.NodeId + "." + dname},
			Labels:  labels,
		})
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *Router) getNodeHandler(c *gin.Context, req *GetNodeIPRequest) (*pb.GetNodeIPResponse, error) {

	return r.clients.n.GetNodeIpRequest(&pb.GetNodeIPRequest{
		NodeId: req.NodeId,
	})
}

func (r *Router) putNodeHandler(c *gin.Context, req *SetNodeIPRequest) (*pb.SetNodeIPResponse, error) {

	return r.clients.n.SetNodeIpRequest(&pb.SetNodeIPRequest{
		NodeId:   req.NodeId,
		NodeIp:   req.NodeIp,
		NodePort: req.NodePort,
		MeshIp:   req.MeshIp,
		MeshPort: req.MeshPort,
		Org:      req.Org,
		Network:  req.Network,
	})

}

func (r *Router) deleteNodeIPHandler(c *gin.Context, req *DeleteNodeIPRequest) (*pb.DeleteNodeIPResponse, error) {

	return r.clients.n.DeleteNodeIpRequest(&pb.DeleteNodeIPRequest{
		NodeId: req.NodeId,
	})
}

func (r *Router) getAllNodeIPHandler(c *gin.Context, req *ListNodeIPsRequest) (*pb.NodeIPMapListResponse, error) {

	return r.clients.n.GetNodeIPMapListRequest(&pb.NodeIPMapListRequest{})
}

func (r *Router) getNodeOrgMapHandler(c *gin.Context, req *NodeOrgMapListRequest) (*pb.NodeOrgMapListResponse, error) {

	return r.clients.n.GetNodeOrgMapListRequest(&pb.NodeOrgMapListRequest{})
}
