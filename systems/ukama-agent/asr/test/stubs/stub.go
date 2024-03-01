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
	"github.com/ukama/ukama/systems/common/uuid"
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
}

type GetPackageRequest struct {
	Package string `path:"id" validate:"required"`
}

type DeleteSimReq struct {
	Imsi string `path:"imsi" validate:"required"`
}

type PackageMarkup struct {
	PackageID  string  `json:"package_id"`
	BaseRateId string  `json:"base_rate_id"`
	Markup     float64 `json:"markup"`
}

type PackageDetails struct {
	Dlbr uint64 `json:"dlbr"`
	Ulbr uint64 `json:"ulbr"`
	Apn  string
}

type PackageInfo struct {
	Id             string         `json:"uuid"`
	Name           string         `json:"name"`
	From           string         `json:"from" validation:"required"`
	To             string         `json:"to" validation:"required"`
	OrgId          string         `json:"org_id" validation:"required"`
	OwnerId        string         `json:"owner_id" validation:"required"`
	SimType        string         `json:"sim_type" validation:"required"`
	SmsVolume      uint64         `json:"sms_volume,string" validation:"required"`
	VoiceVolume    uint64         `json:"voice_volume,string" default:"0"`
	DataVolume     uint64         `json:"data_volume,string" validation:"required"`
	VoiceUnit      string         `json:"voice_unit" validation:"required"`
	DataUnit       string         `json:"data_unit" validation:"required"`
	Type           string         `json:"type" validation:"required"`
	Flatrate       bool           `json:"flat_rate" default:"false"`
	Amount         float64        `json:"amount" default:"0.00"`
	Markup         PackageMarkup  `json:"markup" default:"0.00"`
	PackageDetails PackageDetails `json:"package_details"`
	Apn            string         `json:"apn" default:"ukama.tel"`
	BaserateId     string         `json:"baserate_id" validation:"required"`
	IsActive       bool           `json:"active"`
	Duration       uint64         `json:"duration,string"`
	Overdraft      float64        `json:"overdraft"`
	TrafficPolicy  uint32         `json:"traffic_policy"`
	Networks       []string       `json:"networks"`
	SyncStatus     string         `json:"sync_status,omitempty"`
}

type Package struct {
	PackageInfo *PackageInfo `json:"package"`
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

	r.f = rest.NewFizzRouter(r.config.serverConf, "asr-stubs", "0,.0.0.", r.config.debugMode, "")
	v1 := r.f.Group("/v1", " ", " ASR service stubs version v1")

	f := v1.Group("/factory", "Sim Factory", "Sim Factory")
	f.GET("/simcards/:iccid", formatDoc("Read Sim card Information", ""), tonic.Handler(r.getSimCard, http.StatusOK))

	n := v1.Group("/networks", "Registry", "Network Factory")
	n.GET("/:network/", formatDoc("Validate Network", ""), tonic.Handler(r.getValidateNetwork, http.StatusOK))

	p := v1.Group("/packages", "Packages", "Package")
	p.GET("/:id", formatDoc("Get package", ""), tonic.Handler(r.getPackage, http.StatusOK))

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

func (r *Router) getValidateNetwork(c *gin.Context, req *NetworkValidationReq) (*client.NetworkInfo, error) {
	/* No implementaion always return success.*/
	if req.Network == "40987edb-ebb6-4f84-a27c-99db7c136127" {
		return &client.NetworkInfo{
			NetworkId:     "40987edb-ebb6-4f84-a27c-99db7c136127",
			Name:          "ukama",
			OrgId:         "40987edb-ebb6-4f84-a27c-99db7c136100",
			IsDeactivated: false,
			CreatedAt:     time.Now(),
		}, nil
	}
	return nil, fmt.Errorf("network %s not found", req.Network)
}

func (r *Router) getPackage(c *gin.Context, req *GetPackageRequest) (*Package, error) {
	/* No implementaion always return success.*/
	if req.Package == "40987edb-ebb6-4f84-a27c-99db7c136127" {
		return &Package{
			PackageInfo: &PackageInfo{
				Name:        "Monthly Data",
				OrgId:       uuid.NewV4().String(),
				OwnerId:     uuid.NewV4().String(),
				From:        "2023-04-01T00:00:00Z",
				To:          "2023-04-01T00:00:00Z",
				BaserateId:  uuid.NewV4().String(),
				VoiceVolume: 0,
				IsActive:    true,
				DataVolume:  1024000000,
				SmsVolume:   0,
				DataUnit:    "bytes",
				VoiceUnit:   "seconds",
				SimType:     "test",
				Apn:         "ukama.tel",
				PackageDetails: PackageDetails{
					Dlbr: 15000,
					Ulbr: 2000,
					Apn:  "xyz",
				},
				Type:     "postpaid",
				Flatrate: false,
				Amount:   0,
			},
		}, nil
	} else if req.Package == "40987edb-ebb6-4f84-a27c-99db7c136128" {
		return &Package{
			PackageInfo: &PackageInfo{
				Name:        "Monthly Data 2",
				OrgId:       uuid.NewV4().String(),
				OwnerId:     uuid.NewV4().String(),
				From:        "2023-04-01T00:00:00Z",
				To:          "2023-04-01T00:00:00Z",
				BaserateId:  uuid.NewV4().String(),
				VoiceVolume: 0,
				IsActive:    true,
				DataVolume:  2024000000,
				SmsVolume:   0,
				DataUnit:    "bytes",
				VoiceUnit:   "seconds",
				SimType:     "test",
				Apn:         "ukama.tel",
				PackageDetails: PackageDetails{
					Dlbr: 25000,
					Ulbr: 4000,
					Apn:  "xyz",
				},
				Type:     "postpaid",
				Flatrate: false,
				Amount:   0,
			},
		}, nil
	}
	return nil, fmt.Errorf("package %s not found", req.Package)
}

func main() {
	r := NewRouter(NewRouterConfig())
	r.Run()
}
