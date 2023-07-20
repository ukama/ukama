package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type PDFServer struct {
	host   string
	port   int
	path   string
	prefix string
}

func NewPDFServer(host, path, prefix string, port int) *PDFServer {
	return &PDFServer{
		host:   host,
		port:   port,
		path:   path,
		prefix: prefix,
	}
}

func (p *PDFServer) Start() {
	fs := http.FileServer(http.Dir(p.path))
	http.Handle(p.prefix, http.StripPrefix(p.prefix, filter(fs, p.path)))

	log.Infof("Starting PDF file server on %s:%d ...", p.host, p.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", p.host, p.port), nil))
}

func filter(next http.Handler, root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasPdfExt(r.URL.Path) || isDirOrNotExists(root, r.URL.Path) {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isDirOrNotExists(root, fp string) bool {
	fp = filepath.Join(filepath.Clean(root), filepath.Clean(fp))

	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
	}

	if info.IsDir() {
		return true
	}

	return false
}

func hasPdfExt(s string) bool {
	return strings.HasSuffix(strings.ToLower(s), ".pdf")
}
