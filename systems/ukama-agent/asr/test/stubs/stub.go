package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"

	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

type SimCardInfoReq struct {
	Iccid string `path:"iccid" validate:"required"`
}

// TODO: update
type NetworkValidationReq struct {
	Network string `path:"network" validate:"required,uuid4"`
	Org     string `path:"org" validate:"required,uuid4"`
}

type DeleteSimReq struct {
	Imsi string `path:"imsi" validate:"required"`
}

type Router struct {
	f      *fizz.Fizz
	config *RouterConfig
}

type HttpEndpoints struct {
	Timeout time.Duration
}
type RouterConfig struct {
	httpEndpoints *HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
}

func NewRouter(config *RouterConfig) *Router {
	r := &Router{
		config: config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()
	return r
}

func NewRouterConfig() *RouterConfig {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &RouterConfig{
		httpEndpoints: &HttpEndpoints{
			Timeout: 3 * time.Second,
		},
		serverConf: &rest.HttpConfig{
			Port: 8085,
			Cors: defaultCors,
		},
		debugMode: true,
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

	r.f = rest.NewFizzRouter(r.config.serverConf, "asr-stubs", "0,.0.0.", r.config.debugMode)
	v1 := r.f.Group("/v1", " ", " ASR service stubs version v1")

	f := v1.Group("/factory", "Sim Factory", "Sim Factory")
	f.GET("/simcards/:iccid", formatDoc("Read Sim card Information", ""), tonic.Handler(r.getSimCard, http.StatusOK))

	n := v1.Group("/networks", "Registry", "Network Factory")
	n.GET("/:network/orgs/:org", formatDoc("Validate Network", ""), tonic.Handler(r.getValidateNetwork, http.StatusOK))

	p := v1.Group("/pcrf", "PCRF", "Policy control")
	p.PUT("/sims/:imsi", formatDoc("Add Sim", ""), tonic.Handler(r.putSim, http.StatusOK))
	p.DELETE("/sims/:imsi", formatDoc("Delete Sim", ""), tonic.Handler(r.deleteSim, http.StatusOK))
	p.PATCH("/sims/:imsi", formatDoc("Update Sim Packege Info", ""), tonic.Handler(r.patchSim, http.StatusOK))

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

const letterBytes = "0123456789"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (r *Router) getSimCard(c *gin.Context, req *SimCardInfoReq) (*client.SimCardInfo, error) {
	s := client.SimCardInfo{
		Iccid:          req.Iccid,
		Imsi:           randStringBytes(15),
		Op:             []byte("0123456789012345"),
		Key:            []byte("0123456789012345"),
		Amf:            []byte("800"),
		AlgoType:       1,
		UeDlAmbrBps:    2000000,
		UeUlAmbrBps:    2000000,
		Sqn:            1,
		CsgIdPrsent:    false,
		CsgId:          0,
		DefaultApnName: "ukama",
	}

	return &s, nil
}

func (r *Router) getValidateNetwork(c *gin.Context, req *NetworkValidationReq) error {
	/* No implementaion always return success.*/
	return nil
}

func (r *Router) putSim(c *gin.Context, req *client.PolicyControlSimInfo) error {
	/* No implementaion always return success.*/
	return nil
}

func (r *Router) deleteSim(c *gin.Context, req *DeleteSimReq) error {
	/* No implementaion always return success.*/
	return nil
}

func (r *Router) patchSim(c *gin.Context, req *client.PolicyControlSimPackageUpdate) error {
	/* No implementaion always return success.*/
	return nil
}

func main() {
	r := NewRouter(NewRouterConfig())
	r.Run()
}
