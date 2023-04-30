package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type PDFServer struct {
	host string
	port int
	path string
}

func NewPDFServer(host, path string, port int) *PDFServer {
	return &PDFServer{
		host: host,
		port: port,
		path: path,
	}
}

func (p *PDFServer) Start() {
	fs := http.FileServer(http.Dir(p.path))
	http.Handle("/pdf/", http.StripPrefix("/pdf/", fs))

	log.Infof("Starting PDF file server on %s:%d ...", p.host, p.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", p.host, p.port), nil))
}
