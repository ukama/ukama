package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/bootstrap/bootstrap/cmd/version"
	"github.com/ukama/ukama/services/bootstrap/bootstrap/pkg"
	"github.com/ukama/ukama/services/bootstrap/bootstrap/pkg/lookup"
	"github.com/ukama/ukama/services/bootstrap/bootstrap/pkg/nmr"
	"github.com/ukama/ukama/services/common/rest"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/wI2L/fizz"
)

type Router struct {
	fizz *fizz.Fizz
	port int
	ls   *lookup.LookUp
	fs   *nmr.Factory
}

func (r *Router) Run(close chan error) {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		close <- err
	}
	close <- nil
}

func NewRouter(config *pkg.Config, svcR *sr.ServiceRouter) *Router {

	f := rest.NewFizzRouter(&config.Server, pkg.ServiceName, version.Version, pkg.IsDebugMode)

	ls := lookup.NewLookUp(svcR)

	fs := nmr.NewFactory(svcR)

	r := &Router{fizz: f,
		port: config.Server.Port,
		ls:   ls,
		fs:   fs,
	}

	r.init()
	return r
}

func (r *Router) init() {
	r.fizz.GET("/", nil, tonic.Handler(r.bootstrapGetHandler, http.StatusOK))
}

func (r *Router) bootstrapGetHandler(c *gin.Context, req *BootstrapRequest) error {
	logrus.Debugf("Handling bootstrap request %+v.", req)

	/* Validate the Node ID from DMR
	1. Check if nodeid exist in DMR
	2. Check for the status assigned or not
	*/
	valid, err := r.fs.NmrRequestNodeValidation(req.Nodeid)
	if err != nil {
		logrus.Errorf("Couldn't validate node %s validation failed. Error %s", req.Nodeid, err.Error())
		/* handle failure */
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	if !valid {
		logrus.Errorf("Node %s validation failed.", req.Nodeid)
		/* handle failure */
		return rest.HttpError{
			HttpCode: http.StatusNotAcceptable,
			Message:  "Node validation failure.",
		}
	}

	/* Lookup for node id details:
	1. Read organization assigned to node
	2. IP for the organization interface
	3. Certificates for connecting to orgabization.
	*/
	avail, cred, err := r.ls.LookupRequestOrgCredentialForNode(req.Nodeid)
	if err != nil {
		logrus.Errorf("Couldn't fetch credential for the nodeid %s. Error %s", req.Nodeid, err.Error())
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	if !avail {
		logrus.Errorf("No credential found for the nodeid %s.", req.Nodeid)
		return rest.HttpError{
			HttpCode: http.StatusNotAcceptable,
			Message:  "No credetials for node",
		}
	}

	c.Header("Content-Type", "application/json")
	c.IndentedJSON(http.StatusOK, cred)

	return nil
}
