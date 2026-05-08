package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/reconciler"
)

type HTTPServer struct {
	reconciler *reconciler.Reconciler
}

func NewHTTPServer(r *reconciler.Reconciler) *HTTPServer {
	return &HTTPServer{reconciler: r}
}

func (s *HTTPServer) Register(mux *http.ServeMux) {
	mux.HandleFunc("/v1/sites/", s.handleSites)
	mux.HandleFunc("/v1/ping", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("/v1/version", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("site-controller")) })
}

func (s *HTTPServer) handleSites(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/sites/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		writeError(w, http.StatusBadRequest, "missing site id")
		return
	}
	siteId := parts[0]

	if len(parts) == 1 && r.Method == http.MethodGet {
		s.getState(w, r, siteId)
		return
	}
	if len(parts) < 2 {
		writeError(w, http.StatusNotFound, "unknown route")
		return
	}

	switch parts[1] {
	case "on":
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.setSite(w, r, siteId, db.SiteStateOn)
	case "off":
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.setSite(w, r, siteId, db.SiteStateOff)
	case "service":
		if len(parts) != 3 || r.Method != http.MethodPost {
			writeError(w, http.StatusNotFound, "bad service route")
			return
		}
		s.setService(w, r, siteId, parts[2])
	case "radio":
		if len(parts) != 3 || r.Method != http.MethodPost {
			writeError(w, http.StatusNotFound, "bad radio route")
			return
		}
		s.setRadio(w, r, siteId, parts[2])
	case "state":
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.getState(w, r, siteId)
	case "ports":
		s.handlePorts(w, r, siteId, parts)
	case "nodes":
		s.handleNodes(w, r, siteId, parts)
	case "switch-policy":
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := s.reconciler.ApplySwitchPolicy(siteId); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	default:
		writeError(w, http.StatusNotFound, "unknown route")
	}
}

func (s *HTTPServer) handlePorts(w http.ResponseWriter, r *http.Request, siteId string, parts []string) {
	if len(parts) != 2 {
		writeError(w, http.StatusNotFound, "bad ports route")
		return
	}
	switch r.Method {
	case http.MethodGet:
		ports, err := s.reconciler.GetPortMap(siteId)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"ports": ports})
	case http.MethodPut:
		req := struct {
			Ports []db.SitePortMap `json:"ports"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "bad json")
			return
		}
		if err := s.reconciler.UpsertPortMap(siteId, req.Ports); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *HTTPServer) handleNodes(w http.ResponseWriter, r *http.Request, siteId string, parts []string) {
	if len(parts) != 4 || parts[3] != "power-cycle" || r.Method != http.MethodPost {
		writeError(w, http.StatusNotFound, "bad node action route")
		return
	}
	role := parts[2]
	req := reasonRequest{}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Reason == "" {
		req.Reason = "operator_request"
	}
	if err := s.reconciler.PowerCycleNode(siteId, role, req.Reason); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}

func (s *HTTPServer) setSite(w http.ResponseWriter, r *http.Request, siteId string, state string) {
	req := reasonRequest{}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Reason == "" {
		req.Reason = "operator_request"
	}
	if err := s.reconciler.SetSite(siteId, state, req.Reason, req.RequestedBy); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}

func (s *HTTPServer) setService(w http.ResponseWriter, r *http.Request, siteId string, state string) {
	req := reasonRequest{}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Reason == "" {
		req.Reason = "operator_request"
	}
	if err := s.reconciler.SetService(siteId, state, req.Reason, req.RequestedBy); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}

func (s *HTTPServer) setRadio(w http.ResponseWriter, r *http.Request, siteId string, state string) {
	req := reasonRequest{}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Reason == "" {
		req.Reason = "operator_request"
	}
	if err := s.reconciler.SetRadio(siteId, state, req.Reason, req.RequestedBy); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}

func (s *HTTPServer) getState(w http.ResponseWriter, r *http.Request, siteId string) {
	snapshot, err := s.reconciler.GetSnapshot(siteId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

type reasonRequest struct {
	Reason      string `json:"reason"`
	RequestedBy string `json:"requestedBy"`
}

func writeJSON(w http.ResponseWriter, code int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Errorf("site-controller: failed to write response: %s", err.Error())
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg, "code": strconv.Itoa(code)})
}

func routeErr(format string, args ...interface{}) error { return fmt.Errorf(format, args...) }
