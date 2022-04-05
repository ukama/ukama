package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/cmd/version"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg/client"
	"github.com/ukama/openIoR/services/common/rest"
	"github.com/wI2L/fizz"
)

type Router struct {
	fizz *fizz.Fizz
	port int
	sr   *client.Client
}

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		panic(err)
	}
}

func NewRouter(config *pkg.Config, c *client.Client) *Router {

	f := rest.NewFizzRouter(&config.Server, pkg.ServiceName, version.Version, pkg.IsDebugMode)

	r := &Router{fizz: f,
		port: config.Server.Port,
		sr:   c,
	}

	r.init()
	return r
}

func (r *Router) init() {
	r.fizz.GET("/", nil, tonic.Handler(r.bootstrapGetHandler, http.StatusOK))
}

func (r *Router) bootstrapGetHandler(c *gin.Context, req *BootstrapRequest) error {
	var err error = nil
	logrus.Debugf("Handling bootstrap request %+v.", req)

	/* Validate the Node ID from DMR
	1. Check if nodeid exist in DMR
	2. Check for the status assigned or not
	*/
	

	/* Lookup for node id details:
	1. Read organization assigned to node
	2. IP for the organization interface
	3. Certificates for connecting to orgabization.
	*/

	if err != nil {
		logrus.Errorf("Error while creating index file.")
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}
	return nil

}
