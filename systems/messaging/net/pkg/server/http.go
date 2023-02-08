package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ukama/ukama/systems/messaging/net/pkg"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	pb "github.com/ukama/ukama/systems/messaging/net/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HttpServer struct {
	nnsClient       pkg.NnsReader
	httpConf        rest.HttpConfig
	grpcConf        config.Grpc
	nodeMetricsPort int
	nodeOrgNetMap   *pkg.NodeOrgMap
}

func NewHttpServer(httpConf rest.HttpConfig, grpcConf config.Grpc, nodeMetricsPort int, nnsClient pkg.NnsReader, nodeOrgNetMap *pkg.NodeOrgMap) *HttpServer {
	return &HttpServer{nnsClient: nnsClient,
		httpConf:        httpConf,
		grpcConf:        grpcConf,
		nodeMetricsPort: nodeMetricsPort,
		nodeOrgNetMap:   nodeOrgNetMap,
	}
}

func (h *HttpServer) RunHttpServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterNnsHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", h.grpcConf.Port), opts)
	if err != nil {
		logrus.Fatal("Error registering Grpc endpoint. Error: ", err)
	}
	// handle prometheus
	err = mux.HandlePath(http.MethodGet, "/prometheus", h.prometheusHandler)
	if err != nil {
		logrus.Fatal("Error registering prometheus handler. Error: ", err)
	}

	// handle ping
	err = mux.HandlePath(http.MethodGet, "/ping", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			logrus.Error("Error writing response. Error: ", err)
		}
	})
	if err != nil {
		logrus.Fatal("Error registering prometheus handler. Error: ", err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	logrus.Infof("starting HTTP/1.1 server on %d", h.httpConf.Port)
	err = http.ListenAndServe(fmt.Sprint(":", h.httpConf.Port), mux)
	if err != nil {
		logrus.Fatal("Error starting HTTP server. Error: ", err)
	}
}

func (h *HttpServer) prometheusHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	const ETCD_TIMEOUT = 3
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TIMEOUT*time.Second)
	defer cancel()

	m := make(chan bool)
	nodeToOrg := make(map[string]pkg.OrgNet)
	go func() {
		var errCh error
		if nodeToOrg, errCh = h.nodeOrgNetMap.List(ctx); errCh != nil {
			logrus.Error("Error getting node to org/network map. Error: ", errCh)
		}
		m <- true
	}()

	l, err := h.nnsClient.List(ctx)
	if err != nil {
		logrus.Error("Error getting list of namespaces. Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// wait for nodeToOrgNetwork mapping to finish
	<-m

	b, err := marshallTargets(l, nodeToOrg, h.nodeMetricsPort)
	if err != nil {
		logrus.Error("Error marshalling targets. Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		logrus.Error("Error writing response. Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type targets struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func marshallTargets(t map[string]string, nodeToOrg map[string]pkg.OrgNet, nodeMetricsPort int) ([]byte, error) {
	resp := make([]targets, 0, len(t))

	for k, v := range t {
		labels := map[string]string{
			"nodeid": k,
		}
		if m, ok := nodeToOrg[k]; ok {
			labels["org"] = m.Org
			labels["network"] = m.Network
		}

		resp = append(resp, targets{
			Targets: []string{fmt.Sprintf("%s:%d", v, nodeMetricsPort)},
			Labels:  labels,
		})
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}
