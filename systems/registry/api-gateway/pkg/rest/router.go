package rest

import (
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/errors"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/wI2L/fizz/openapi"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/registry/api-gateway/cmd/version"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/registry/api-gateway/pkg"
	"github.com/ukama/ukama/systems/registry/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	pbnode "github.com/ukama/ukama/systems/registry/node/pb/gen"
	pborg "github.com/ukama/ukama/systems/registry/org/pb/gen"
	userspb "github.com/ukama/ukama/systems/registry/users/pb/gen"
)

const USER_ID_KEY = "UserId"
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
	Registry registry
	User     *client.Users
}

type registry interface {
	GetOrg(orgName string) (*pborg.Organization, error)
	AddOrg(orgName string, owner string) (*pborg.Organization, error)
	Add(orgName string, nodeId string, name string, attachedNodes ...string) (node *pbnode.Node, err error)
	UpdateNode(orgName string, nodeId string, name string, attachedNodes ...string) (node *pbnode.Node, err error)
	GetNodes(orgName string) (*pb.GetNodesResponse, error)
	GetNode(nodeId string) (*pbnode.GetNodeResponse, error)
	AttachNode(towerNodeId string, amplNodeId ...string)
	IsAuthorized(userId string, org string) (bool, error)
	DeleteNode(nodeId string) (*pb.DeleteNodeResponse, error)
	DetachNode(nodeId string, attachedId string) (*pbnode.Node, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Registry = client.NewRegistry(endpoints.Network, endpoints.Org, endpoints.Node, endpoints.Timeout)
	c.User = client.NewUsers(endpoints.Users, endpoints.Timeout)
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

	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "API gateway", "Registry system version v1")

	// org handler
	orgs := v1.Group(org, "Orgs", "Operations on Orgs")
	orgs.GET("", []fizz.OperationOption{}, tonic.Handler(r.getOrgHandler, http.StatusOK))
	orgs.PUT("", formatDoc("Add Org", ""), tonic.Handler(r.putOrgHandler, http.StatusCreated))

	// network
	nodes := orgs.Group("/nodes", "Nodes", "Nodes operations")
	nodes.GET("", nil, tonic.Handler(r.getNodesHandler, http.StatusOK))
	nodes.PUT("/:node", nil, tonic.Handler(r.addNodeHandler, http.StatusCreated))
	nodes.PATCH("/:node", nil, tonic.Handler(r.updateNodeHandler, http.StatusOK))
	nodes.GET("/:node", nil, tonic.Handler(r.getNodeHandler, http.StatusOK))
	nodes.DELETE("/:node", nil, tonic.Handler(r.deleteNodeHandler, http.StatusOK))
	nodes.DELETE("/:node/attached/:attachedId", nil, tonic.Handler(r.detachNode, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) getOrgNameFromRoute(c *gin.Context) string {
	return c.Param("org")
}

// Org handlers

func (r *Router) getOrgHandler(c *gin.Context) (*pborg.Organization, error) {
	orgName := r.getOrgNameFromRoute(c)
	return r.clients.Registry.GetOrg(orgName)
}

func (r *Router) putOrgHandler(c *gin.Context, req *AddOrgRequest) (*pborg.Organization, error) {
	orgName := r.getOrgNameFromRoute(c)
	return r.clients.Registry.AddOrg(orgName, req.Owner)
}

// Node handlers

func (r *Router) getNodesHandler(c *gin.Context) (*NodesList, error) {
	orgName := r.getOrgNameFromRoute(c)
	nl, err := r.clients.Registry.GetNodes(orgName)
	if err != nil {
		return nil, err
	}

	return MapNodesList(nl), nil
}

func (r *Router) getNodeHandler(c *gin.Context, req *GetNodeRequest) (*NodeExtended, error) {
	nodeName := c.Param("node")
	resp, err := r.clients.Registry.GetNode(nodeName)
	if err != nil {
		return nil, err
	}
	return mapExtendeNode(resp.Node), nil
}

func (r *Router) getAttacheNodesIds(nodes []*NodeAttach) []string {
	nds := []string{}
	for _, n := range nodes {
		nds = append(nds, n.NodeId)
	}
	return nds
}

func (r *Router) updateNodeHandler(c *gin.Context, req *AddUpdateNodeRequest) (*pbnode.Node, error) {
	node, err := r.clients.Registry.UpdateNode(req.OrgName, req.NodeId, req.Node.Name, r.getAttacheNodesIds(req.Node.Attached)...)

	return node, err
}

func (r *Router) addNodeHandler(c *gin.Context, req *AddUpdateNodeRequest) (*pbnode.Node, error) {
	node, err := r.clients.Registry.Add(req.OrgName, req.NodeId, req.Node.Name, r.getAttacheNodesIds(req.Node.Attached)...)
	return node, err
}

func (r *Router) deleteNodeHandler(c *gin.Context, req *DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	return r.clients.Registry.DeleteNode(req.NodeId)

}

func (r *Router) detachNode(c *gin.Context, req *DetachNodeRequest) (*pbnode.Node, error) {
	node, err := r.clients.Registry.DetachNode(req.NodeId, req.AttachedNodeId)
	return node, err
}

// Users handlers: all these need to be updated.

func (r *Router) getUsersHandler(c *gin.Context) (*userspb.ListResponse, error) {
	orgName := r.getOrgNameFromRoute(c)
	return r.clients.User.GetUsers(orgName, c.GetString(USER_ID_KEY))
}

func (r *Router) postUsersHandler(c *gin.Context, req *UserRequest) (*userspb.AddResponse, error) {
	return r.clients.User.AddUser(req.Org, &userspb.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	},
		req.SimToken,
		c.GetString(USER_ID_KEY))
}

func (r *Router) updateUserHandler(c *gin.Context, req *UpdateUserRequest) (*userspb.User, error) {
	if req.IsDeactivated {
		err := r.clients.User.DeactivateUser(req.UserId, c.GetString(USER_ID_KEY))
		if err != nil {
			return nil, err
		}
	}

	_, err := r.clients.User.UpdateUser(req.UserId, &userspb.UserAttributes{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}, c.GetString(USER_ID_KEY))

	if err != nil {
		return nil, err
	}

	resUser, err := r.clients.User.Get(req.UserId, c.GetString(USER_ID_KEY))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated user")
	}

	return resUser.GetUser(), nil
}

func (r *Router) deleteUserHandler(c *gin.Context, req *DeleteUserRequest) error {
	return r.clients.User.Delete(req.UserId, c.GetString(USER_ID_KEY))
}

func (r *Router) getUserHandler(c *gin.Context, req *GetUserRequest) (*userspb.GetResponse, error) {
	return r.clients.User.Get(req.UserId, c.GetString(USER_ID_KEY))
}

func (r *Router) setSimStatusHandler(c *gin.Context, req *SetSimStatusRequest) (*userspb.Sim, error) {
	return r.clients.User.SetSimStatus(&userspb.SetSimStatusRequest{
		Iccid:   req.Iccid,
		Carrier: simServicesToPbService(req.Carrier),
		Ukama:   simServicesToPbService(req.Ukama),
	}, c.GetString(USER_ID_KEY))
}

func (r *Router) getSimQrHandler(c *gin.Context, req *GetSimQrRequest) (*userspb.GetQrCodeResponse, error) {
	return r.clients.User.GetQr(req.Iccid, c.GetString(USER_ID_KEY))
}

func boolToPbBool(data *bool) *wrapperspb.BoolValue {
	if data == nil {
		return nil
	}
	return &wrapperspb.BoolValue{Value: *data}
}

func simServicesToPbService(data *SimServices) *userspb.SetSimStatusRequest_SetServices {
	if data == nil {
		return nil
	}

	return &userspb.SetSimStatusRequest_SetServices{
		Data:  boolToPbBool(data.Data),
		Voice: boolToPbBool(data.Voice),
		Sms:   boolToPbBool(data.Sms),
	}
}
