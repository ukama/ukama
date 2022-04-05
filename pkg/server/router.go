package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/cmd/version"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg/client"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg/lookup"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg/nmr"
	"github.com/ukama/openIoR/services/common/rest"
	"github.com/wI2L/fizz"
)

type Router struct {
	fizz *fizz.Fizz
	port int
	rs   *client.Client
	ls   *lookup.LookUp
	fs   *nmr.Factory
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

	ls := lookup.NewLookUp(c)

	fs := nmr.NewFactory(c)

	r := &Router{fizz: f,
		port: config.Server.Port,
		rs:   c,
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
	var err error = nil
	logrus.Debugf("Handling bootstrap request %+v.", req)

	/* Validate the Node ID from DMR
	1. Check if nodeid exist in DMR
	2. Check for the status assigned or not
	*/
	valid, err := r.fs.NmrRequestNodeValidation(req.Nodeid)
	if err != nil {
		logrus.Errorf("Couldn't validate node %s validatin failed. Error %s", req.Nodeid, err.Error())
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

	_, err = io.Copy(c.Writer, bytes.NewReader(cred.OrgCred))
	if err != nil {

		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}

	}
	c.Header("Content-Type", "application/octet-stream")

	return nil

}
