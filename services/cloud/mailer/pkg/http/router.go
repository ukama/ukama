package http

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/cloud/mailer/cmd/version"
	"github.com/ukama/ukama/services/cloud/mailer/pkg"
	"github.com/ukama/ukama/services/common/rest"
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

func NewRouter(config *rest.HttpConfig) *Router {
	f := rest.NewFizzRouter(config, pkg.ServiceName, version.Version, pkg.IsDebugMode)

	r := &Router{
		fizz: f,
		port: config.Port,
	}
	r.init()
	return r
}

func (r *Router) init() {
	// fooGroup := r.fizz.Group("/foos", "Foo list", "Foo operations")
	// fooGroup.GET("/:name", nil, tonic.Handler(r.fooGetHandler, http.StatusOK))
}
