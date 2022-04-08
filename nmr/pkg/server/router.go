package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/common/rest"
	"github.com/ukama/openIoR/services/common/sql"
	"github.com/ukama/openIoR/services/common/ukama"
	"github.com/ukama/openIoR/services/factory/nmr/cmd/version"
	"github.com/ukama/openIoR/services/factory/nmr/internal/db"
	"github.com/ukama/openIoR/services/factory/nmr/pkg"
	rs "github.com/ukama/openIoR/services/factory/nmr/pkg/router"
	"github.com/wI2L/fizz"
)

const (
	NodePath   = "/node"
	ModulePath = "/module"
)

type Router struct {
	fizz       *fizz.Fizz
	port       int
	R          *rs.RouterServer
	nodeRepo   db.NodeRepo
	moduleRepo db.ModuleRepo
}

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		panic(err)
	}
}

func NewRouter(config *pkg.Config, rs *rs.RouterServer, d sql.Db) *Router {

	f := rest.NewFizzRouter(&config.Server, pkg.ServiceName, version.Version, pkg.IsDebugMode)
	nodeRepo := db.NewNodeRepo(d)

	moduleRepo := db.NewModuleRepo(d)

	r := &Router{fizz: f,
		port:       config.Server.Port,
		R:          rs,
		nodeRepo:   nodeRepo,
		moduleRepo: moduleRepo,
	}

	r.init()
	return r
}

func (r *Router) init() {
	node := r.fizz.Group(NodePath, "Node", "Node related operations")
	node.GET("/", nil, tonic.Handler(r.GetNodeHandler, http.StatusOK))

	module := r.fizz.Group(ModulePath, "Module", "Module related operations")
	module.GET("/", nil, tonic.Handler(r.GetModuleHandler, http.StatusOK))
}

func (r *Router) GetNodeHandler(c *gin.Context, req *ReqGetNode) (*RespGetNode, error) {
	logrus.Debugf("Handling NMR get request %+v.", req)

	node, err := r.nodeRepo.GetNode(ukama.NodeID(req.NodeID))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	rnode := &RespGetNode{
		Node: node,
	}

	return rnode, nil
}

func (r *Router) GetModuleHandler(c *gin.Context, req *ReqGetNode) error {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil
}
