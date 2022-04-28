package distribution

import (
	"context"
	"fmt"
	"io"

	"log"
	"net/http"
	"os"

	casync "github.com/folbricht/desync"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/hub/distributor/pkg"
	"github.com/ukama/ukama/services/hub/distributor/pkg/chunk"
)

var (
	stderr io.Writer = os.Stderr
)

func RunDistribution(ctx context.Context, serverCfg *pkg.DistributionConfig) error {

	addresses := serverCfg.Address
	if len(addresses) == 0 {
		addresses = []string{":http"}
	}
	logrus.Debugf("Starting distribution server at %+v", addresses)

	/* Store set up */
	s, err := chunkServerStore(serverCfg)
	if err != nil {
		logrus.Errorf("Error configuring distribution server store : %s", err.Error())
		return err
	}

	var converters casync.Converters
	if !serverCfg.StoreCfg.Uncompressed {
		converters = casync.Converters{casync.Compressor{}}
	}

	handler := casync.NewHTTPHandler(s, false, true, converters, "")

	// Wrap the handler in a logger if requested
	switch serverCfg.LogFile {
	case "": // No logging of requests
	case "-":
		handler = withLog(handler, log.New(stderr, "", log.LstdFlags))
	default:
		l, err := os.OpenFile(serverCfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logrus.Errorf("Error while setting log filr for distribution server : %s", err.Error())
			return err
		}
		defer l.Close()
		handler = withLog(handler, log.New(l, "", log.LstdFlags))
		logrus.Debugf("Distribution server logging at %s", serverCfg.LogFile)
	}

	http.Handle("/", handler)

	// Start the server
	return serve(ctx, &serverCfg.Security, addresses...)
}

// Wrapper for http.HandlerFunc to add logging for requests (and response codes)
func withLog(h http.Handler, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w}
		h.ServeHTTP(lrw, r)
		log.Printf("Client: %s, Request: %s %s, Response: %d", r.RemoteAddr, r.Method, r.RequestURI, lrw.statusCode)
	}
}

// Reads the store-related command line options and returns the appropriate store.
func chunkServerStore(serverCfg *pkg.DistributionConfig) (casync.Store, error) {
	stores := serverCfg.Chunk.Stores

	// Got to have at least one upstream store
	if len(stores) == 0 {
		return nil, errors.New("no store provided")
	}

	var s casync.Store
	// if false {
	// 	s, err := WritableStore(stores[0], serverCfg.StoreCfg)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// } else {
	var opt = &casync.StoreOptions{
		N:             10,
		ClientCert:    "",
		CACert:        "",
		ClientKey:     "",
		SkipVerify:    false,
		TrustInsecure: false,
		ErrorRetry:    3,
		Uncompressed:  false,
	}

	s, err := chunk.MultiStore(*opt, stores...)
	if err != nil {
		return nil, err
	}
	// We want to take the edge of a large number of requests coming in for the same chunk. No need
	// to hit the (potentially slow) upstream stores for duplicated requests.
	s = casync.NewDedupQueue(s)
	//	}

	return s, nil
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func serve(ctx context.Context, storeOptions *pkg.SecurityConfig, addresses ...string) error {

	logrus.Info("Starting Distribution server at", addresses)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, addr := range addresses {
		go func(a string) {
			server := &http.Server{
				Addr:     a,
				ErrorLog: log.New(stderr, "", log.LstdFlags),
			}
			err := server.ListenAndServe()

			fmt.Fprintln(stderr, err)
			cancel()
		}(addr)
	}
	// wait for either INT/TERM or an issue with the server
	<-ctx.Done()
	return nil
}
