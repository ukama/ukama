package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/cmd/version"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg"
	"github.com/ukama/openIoR/services/common/rest"
	"github.com/wI2L/fizz"
)

type Router struct {
	fizz *fizz.Fizz
	port int
}

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		panic(err)
	}
}

func NewRouter(config *pkg.Config) *Router {

	f := rest.NewFizzRouter(&config.Server, pkg.ServiceName, version.Version, pkg.IsDebugMode)

	r := &Router{fizz: f,
		port: config.Server.Port,
	}

	r.init()
	return r
}

func (r *Router) init() {
	r.fizz.GET("/node", nil, tonic.Handler(r.bootstrapGetHandler, http.StatusOK))
}

func (r *Router) bootstrapGetHandler(c *gin.Context, req *BootstrapRequest) error {
	var err error = nil

	logrus.Debugf("Handling bootstrap request %+v.", req)

	if err != nil {
		logrus.Errorf("Error while creating index file.")
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}
	return nil

}
