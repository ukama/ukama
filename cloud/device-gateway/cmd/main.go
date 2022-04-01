package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/cloud/device-gateway/cmd/version"
	gen "github.com/ukama/ukamaX/cloud/device-gateway/pb/gen"
	"github.com/ukama/ukamaX/cloud/device-gateway/pkg"
	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var svcConf = pkg.NewConfig()
var ServiceName = "device-gateway"

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grcpMux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register HSS gRPC server endpoint
	err := gen.RegisterImsiServiceHandlerFromEndpoint(ctx, grcpMux, svcConf.Services.Hss, opts)
	if err != nil {
		return err
	}

	mux := registerExtraHandlers(grcpMux)
	logMiddlware := pkg.LoggingMiddleware(svcConf.DebugMode)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	logrus.Info("Listening on port ", svcConf.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", svcConf.Port), logMiddlware(mux))
}

func registerExtraHandlers(grcpMux *runtime.ServeMux) *http.ServeMux {
	// serve swagger-ui
	m := http.NewServeMux()
	m.Handle("/", grcpMux)

	swaggerUiPath := svcConf.SwaggerAssets
	const swaggerUrlPath = "/swagger/"

	if _, err := os.Stat(swaggerUiPath); os.IsNotExist(err) {
		logrus.Fatalf("Directory %s not found", swaggerUiPath)
	}

	m.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerUiPath+"/hss.swagger.json")
	})

	sh := http.StripPrefix(swaggerUrlPath, http.FileServer(http.Dir(swaggerUiPath)))
	m.Handle(swaggerUrlPath, sh)

	// handle ping
	m.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			logrus.Errorf("Error writing response")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	return m
}

func main() {
	ccmd.ProcessVersionArgument(ServiceName, os.Args, version.Version)
	initConfig()

	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(ServiceName, svcConf)
}
