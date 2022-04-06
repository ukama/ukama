package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/common/rest"
	"github.com/ukama/openIoR/services/factory/nmr/cmd/version"
	"github.com/ukama/openIoR/services/factory/nmr/pkg"
	rs "github.com/ukama/openIoR/services/factory/nmr/pkg/router"
	"github.com/wI2L/fizz"
)

type Router struct {
	fizz *fizz.Fizz
	port int
	R    *rs.RouterServer
}

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		panic(err)
	}
}

func NewRouter(config *pkg.Config, rs *rs.RouterServer) *Router {

	f := rest.NewFizzRouter(&config.Server, pkg.ServiceName, version.Version, pkg.IsDebugMode)

	r := &Router{fizz: f,
		port: config.Server.Port,
		R:    rs,
	}

	r.init()
	return r
}

func (r *Router) init() {
	r.fizz.GET("/", nil, tonic.Handler(r.nmrGetHandler, http.StatusOK))
}

func (r *Router) nmrGetHandler(c *gin.Context, req *BootstrapRequest) error {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil
}
