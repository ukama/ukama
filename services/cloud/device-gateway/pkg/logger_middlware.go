package pkg

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"time"
)

func LoggingMiddleware(isDebug bool) func(http.Handler) http.Handler {
	if isDebug {
		return func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						logrus.Error(
							"err:", err,
							" trace:", debug.Stack(),
						)
					}
				}()

				start := time.Now()
				wrapped := wrapResponseWriter(w)
				next.ServeHTTP(wrapped, r)

				logrus.Info(
					"status: ", wrapped.status,
					" method:", r.Method,
					" path:", r.URL.EscapedPath(),
					" duration:", time.Since(start),
				)
			}

			return http.HandlerFunc(fn)
		}
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status int
}
