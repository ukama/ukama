package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/net/pb/gen"
	"github.com/ukama/ukamaX/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

type HttpServer struct {
	nnsClient       *Nns
	httpConf        config.Http
	grpcConf        config.Grpc
	nodeMetricsPort int
}

func NewHttpServer(httpConf config.Http, grpcConf config.Grpc, nodeMetricsPort int, nnsClient *Nns) *HttpServer {
	return &HttpServer{nnsClient: nnsClient,
		httpConf:        httpConf,
		grpcConf:        grpcConf,
		nodeMetricsPort: nodeMetricsPort}
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

	l, err := h.nnsClient.List(ctx)
	if err != nil {
		logrus.Error("Error getting list of namespaces. Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := marshallTargets(l, h.nodeMetricsPort)
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

func marshallTargets(t map[string]string, nodeMetricsPort int) ([]byte, error) {
	resp := make([]targets, 0, len(t))

	for k, v := range t {
		resp = append(resp, targets{
			Targets: []string{fmt.Sprintf("%s:%d", v, nodeMetricsPort)},
			Labels: map[string]string{
				"nodeid": k,
			},
		})
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}
