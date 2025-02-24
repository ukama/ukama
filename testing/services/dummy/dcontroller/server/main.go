package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ukama/ukama/testing/services/dummy/dcontroller/config"  // Replace with your actual import path
	"github.com/ukama/ukama/testing/services/dummy/dcontroller/metrics" // Replace with your actual import path
)

type Server struct {
	mu              sync.RWMutex
	activeExporters map[string]*metrics.PrometheusExporter
	metricsProvider *metrics.MetricsProvider
}

func NewServer() *Server {
	return &Server{
		activeExporters: make(map[string]*metrics.PrometheusExporter),
		metricsProvider: metrics.NewMetricsProvider(),
	}
}

func (s *Server) updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg config.WMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if msg.SiteId == "" {
		http.Error(w, "SiteId cannot be empty", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	exporter, exists := s.activeExporters[msg.SiteId]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, fmt.Sprintf("Site %s does not exist", msg.SiteId), http.StatusNotFound)
		return
	}

	// Stop existing metrics collection for this site
	exporter.Shutdown()
	
	go s.startMetricsCollection(msg.SiteId, msg.Profile, msg.Scenario)

	log.Printf("Updated site %s with profile %d and scenario %s", msg.SiteId, msg.Profile, msg.Scenario)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Site %s updated with profile %d and scenario %s", msg.SiteId, msg.Profile, msg.Scenario)))
}

func (s *Server) siteCreatedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg config.WMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if msg.SiteId == "" {
		http.Error(w, "SiteId cannot be empty", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	_, exists := s.activeExporters[msg.SiteId]
	s.mu.RUnlock()

	if exists {
		http.Error(w, fmt.Sprintf("Site %s already exists", msg.SiteId), http.StatusConflict)
		return
	}

	if msg.Profile == 0 {
		msg.Profile = config.PROFILE_NORMAL 
	}

	if msg.Scenario == "" {
		msg.Scenario = config.SCENARIO_DEFAULT 
	}

	go s.startMetricsCollection(msg.SiteId, msg.Profile, msg.Scenario)

	log.Printf("Created site %s with profile %d and scenario %s", msg.SiteId, msg.Profile, msg.Scenario)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Site %s created with profile %d and scenario %s", msg.SiteId, msg.Profile, msg.Scenario)))
}

func (s *Server) startMetricsCollection(siteId string, profile config.Profile, scenario config.SCENARIOS) {
	s.applyScenario(scenario, siteId)
	
	exporter := metrics.NewPrometheusExporter(s.metricsProvider, siteId)
	
	s.mu.Lock()
	s.activeExporters[siteId] = exporter
	s.mu.Unlock()
	
	var interval time.Duration
	switch profile {
	case config.PROFILE_MIN:
		interval = 30 * time.Second
	case config.PROFILE_MAX:
		interval = 5 * time.Second
	default: 
		interval = 10 * time.Second
	}
	
	ctx := context.Background()
	log.Printf("Starting metrics collection for site %s with profile %d and scenario %s (interval: %v)", 
		siteId, profile, scenario, interval)
	
	go func() {
		if err := exporter.StartMetricsCollection(ctx, interval); err != nil {
			log.Printf("Metrics collection for site %s stopped: %v", siteId, err)
			
			s.mu.Lock()
			delete(s.activeExporters, siteId)
			s.mu.Unlock()
		}
	}()
}

func (s *Server) applyScenario(scenario config.SCENARIOS, siteId string) {
    log.Printf("Applying scenario %s to site %s", scenario, siteId)
    
    provider := s.metricsProvider
    
    if provider.Backhaul() != nil {
        provider.Backhaul().SetForceBackhaulDown(false)
        provider.Backhaul().SetForceSwitchOff(false)
    }
    
    switch scenario {
    case config.SCENARIO_SOLAR_DOWN:
        if provider.Solar() != nil {
            provider.Solar().SetWeatherPattern(0.0)  
        }
        
    case config.SCENARIO_BATTERY_LOW:
        if provider.Battery() != nil {
            provider.Battery().SetLastCapacity(15.0) 
        }
        
    case config.SCENARIO_SWITCH_OFF:
        if provider.Backhaul() != nil {
            provider.Backhaul().SetForceSwitchOff(true)
        }
        
    case config.SCENARIO_BACKHAUL_DOWN:
        if provider.Backhaul() != nil {
            provider.Backhaul().SetForceBackhaulDown(true)
        }
        
    default: 
        if provider.Solar() != nil {
            provider.Solar().SetWeatherPattern(1.0) 
        }
        if provider.Battery() != nil {
            provider.Battery().SetLastCapacity(85.0)  
        }
    }
}


func (s *Server) listSitesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	s.mu.RLock()
	sites := make([]string, 0, len(s.activeExporters))
	for site := range s.activeExporters {
		sites = append(sites, site)
	}
	s.mu.RUnlock()
	
	response := struct {
		Sites     []string `json:"sites"`
		Count     int      `json:"count"`
		Timestamp int64    `json:"timestamp"`
	}{
		Sites:     sites,
		Count:     len(sites),
		Timestamp: time.Now().Unix(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	server := NewServer()
	
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/update", server.updateHandler)
	http.HandleFunc("/create", server.siteCreatedHandler)
	
	port := config.PORT
	log.Printf("Server starting on port %d", port)
	log.Printf("Dmetrics available at http://localhost:%d/metrics", port)
	
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}