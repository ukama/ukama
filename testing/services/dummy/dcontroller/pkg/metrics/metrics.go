package metrics

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

const (
	PORT_NODE     = 1
	PORT_SOLAR    = 2
	PORT_BACKHAUL = 3
)

type Metrics struct {
	backhaulLatency          *prometheus.GaugeVec
	backhaulSpeed            *prometheus.GaugeVec
	batteryChargePercentage  *prometheus.GaugeVec
	solarPanelVoltage        *prometheus.GaugeVec
	solarPanelCurrent        *prometheus.GaugeVec
	solarPanelPower          *prometheus.GaugeVec
	siteUptimeSeconds        *prometheus.GaugeVec
	switchPortStatus         *prometheus.GaugeVec
	switchPortSpeed          *prometheus.GaugeVec
	switchPortPower          *prometheus.GaugeVec
	mu                       sync.Mutex
	siteUptimeCounters       map[string]int

	backhaulSwitchPortStatus *prometheus.GaugeVec
	backhaulSwitchPortSpeed  *prometheus.GaugeVec
	backhaulSwitchPortPower  *prometheus.GaugeVec

	solarSwitchPortStatus *prometheus.GaugeVec
	solarSwitchPortSpeed  *prometheus.GaugeVec
	solarSwitchPortPower  *prometheus.GaugeVec

	nodeSwitchPortPowerVec  *prometheus.GaugeVec
	nodeSwitchPortSpeedVec  *prometheus.GaugeVec
	nodeSwitchPortStatusVec *prometheus.GaugeVec

	portStatus map[string]map[int]bool 
}

func New() *Metrics {
	log.Println("DEBUG: Initializing Metrics struct")
	m := &Metrics{
		backhaulLatency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "main_backhaul_latency",
			Help: "Backhaul latency in milliseconds",
		}, []string{"site"}),

		backhaulSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_speed",
			Help: "Backhaul speed in Mbps",
		}, []string{"site"}),

		batteryChargePercentage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "battery_charge_percentage",
			Help: "Battery charge percentage",
		}, []string{"site"}),

		solarPanelVoltage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_voltage",
			Help: "Solar panel voltage in volts",
		}, []string{"site"}),

		solarPanelCurrent: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_current",
			Help: "Solar panel current in amperes",
		}, []string{"site"}),

		solarPanelPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_power",
			Help: "Solar panel power in watts",
		}, []string{"site"}),

		siteUptimeSeconds: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "site_uptime_seconds",
			Help: "Site uptime in seconds",
		}, []string{"site"}),

		siteUptimeCounters: make(map[string]int),

		backhaulSwitchPortStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_switch_port_status",
			Help: "Backhaul switch port status (1 = up, 0 = down)",
		}, []string{"site"}),

		backhaulSwitchPortSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_switch_port_speed",
			Help: "Backhaul switch port speed in Mbps",
		}, []string{"site"}),

		backhaulSwitchPortPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_switch_port_power",
			Help: "Backhaul switch port power in watts",
		}, []string{"site"}),
		solarSwitchPortStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_switch_port_status",
			Help: "Solar switch port status (1 = up, 0 = down)",
		}, []string{"site"}),

		solarSwitchPortSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_switch_port_speed",
			Help: "Solar switch port speed in Mbps",
		}, []string{"site"}),

		solarSwitchPortPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_switch_port_power",
			Help: "Solar switch port power in watts",
		}, []string{"site"}),

		nodeSwitchPortPowerVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_switch_port_power",
			Help: "Node switch port power in watts",
		}, []string{"site"}),

		nodeSwitchPortSpeedVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_switch_port_speed",
			Help: "Node switch port speed in Mbps",
		}, []string{"site"}),

		nodeSwitchPortStatusVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_switch_port_status",
			Help: "Node switch port status (1 = up, 0 = down)",
		}, []string{"site"}),

		// Initialize port status map
		portStatus: make(map[string]map[int]bool),
	}

	// Register all metrics with Prometheus
	log.Println("DEBUG: Registering metrics with Prometheus")
	prometheus.MustRegister(m.backhaulLatency)
	prometheus.MustRegister(m.backhaulSpeed)
	prometheus.MustRegister(m.batteryChargePercentage)
	prometheus.MustRegister(m.solarPanelVoltage)
	prometheus.MustRegister(m.solarPanelCurrent)
	prometheus.MustRegister(m.solarPanelPower)
	prometheus.MustRegister(m.siteUptimeSeconds)

	prometheus.MustRegister(m.backhaulSwitchPortStatus)
	prometheus.MustRegister(m.backhaulSwitchPortSpeed)
	prometheus.MustRegister(m.backhaulSwitchPortPower)

	prometheus.MustRegister(m.solarSwitchPortStatus)
	prometheus.MustRegister(m.solarSwitchPortSpeed)
	prometheus.MustRegister(m.solarSwitchPortPower)

	prometheus.MustRegister(m.nodeSwitchPortPowerVec)
	prometheus.MustRegister(m.nodeSwitchPortSpeedVec)
	prometheus.MustRegister(m.nodeSwitchPortStatusVec)

	return m
}

func getSingleMetricValue(gauge *prometheus.GaugeVec, siteID string) float64 {
	metric, err := gauge.GetMetricWithLabelValues(siteID)
	if err != nil {
		log.Printf("DEBUG: Error getting metric: %v", err)
		return -999 
	}

	ch := make(chan prometheus.Metric, 1)
	metric.Collect(ch)
	m := <-ch

	var metricOut dto.Metric
	m.Write(&metricOut)
	return *metricOut.Gauge.Value
}

func (m *Metrics) UpdatePortStatus(siteID string, portNumber int, enabled bool) error {

	m.mu.Lock()

	if _, exists := m.portStatus[siteID]; !exists {
		m.portStatus[siteID] = make(map[int]bool)
	}

	m.portStatus[siteID][portNumber] = enabled
	m.mu.Unlock()

	status := 0.0
	if enabled {
		status = 1.0
	}

	switch portNumber {
	case PORT_NODE:
		m.nodeSwitchPortStatusVec.WithLabelValues(siteID).Set(status)

		if !enabled {
			log.Printf("DEBUG: Setting node metrics to zero for site %s", siteID)
			m.nodeSwitchPortSpeedVec.WithLabelValues(siteID).Set(0)
			m.nodeSwitchPortPowerVec.WithLabelValues(siteID).Set(0)
		} else {
			log.Printf("DEBUG: Setting node metrics to default values for site %s", siteID)
			m.nodeSwitchPortSpeedVec.WithLabelValues(siteID).Set(100) 
			m.nodeSwitchPortPowerVec.WithLabelValues(siteID).Set(50)  
		}

	case PORT_SOLAR:
		m.solarSwitchPortStatus.WithLabelValues(siteID).Set(status)

		if !enabled {
			m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(0)
			m.solarSwitchPortPower.WithLabelValues(siteID).Set(0)
			m.batteryChargePercentage.WithLabelValues(siteID).Set(0)
			m.solarPanelPower.WithLabelValues(siteID).Set(0)
			m.solarPanelCurrent.WithLabelValues(siteID).Set(0)
			m.solarPanelVoltage.WithLabelValues(siteID).Set(0)
		} else {
			m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(100) 
			m.solarSwitchPortPower.WithLabelValues(siteID).Set(50) 
		}

	case PORT_BACKHAUL:
		m.backhaulSwitchPortStatus.WithLabelValues(siteID).Set(status)
		log.Printf("DEBUG: Set backhaulSwitchPortStatus[%s] = %v", siteID, status)

		if !enabled {
			log.Printf("DEBUG: Setting backhaul metrics to zero/high for site %s", siteID)
			m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(0)
			m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(0)

			m.backhaulSpeed.WithLabelValues(siteID).Set(0)

			m.backhaulLatency.WithLabelValues(siteID).Set(9999) 
			log.Printf("DEBUG: Verifying backhaul metrics were set correctly:")
			log.Printf("DEBUG: backhaulSwitchPortSpeed=%v", getSingleMetricValue(m.backhaulSwitchPortSpeed, siteID))
			log.Printf("DEBUG: backhaulSwitchPortPower=%v", getSingleMetricValue(m.backhaulSwitchPortPower, siteID))
			log.Printf("DEBUG: backhaulSpeed=%v", getSingleMetricValue(m.backhaulSpeed, siteID))
			log.Printf("DEBUG: backhaulLatency=%v", getSingleMetricValue(m.backhaulLatency, siteID))
		} else {
			log.Printf("DEBUG: Setting backhaul metrics to default values for site %s", siteID)
			m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(100) 
			m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(50)  
		}

	default:
		return fmt.Errorf("unknown port number: %d", portNumber)
	}

	return nil
}

func (m *Metrics) StartMetricsGenerator(siteID string) {
	log.Printf("DEBUG: Starting metrics generator for site %s", siteID)

	m.mu.Lock()
	m.siteUptimeCounters[siteID] = rand.Intn(6) + 5
	m.portStatus[siteID] = map[int]bool{
		PORT_NODE:     true,
		PORT_SOLAR:    true,
		PORT_BACKHAUL: true,
	}
	log.Printf("DEBUG: Initialized port status: %v", m.portStatus[siteID])
	m.mu.Unlock()

	m.backhaulSwitchPortStatus.WithLabelValues(siteID).Set(1) 
	m.solarSwitchPortStatus.WithLabelValues(siteID).Set(1)   
	m.nodeSwitchPortStatusVec.WithLabelValues(siteID).Set(1)  

	m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(100) 
	m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(50)  

	m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(100) 
	m.solarSwitchPortPower.WithLabelValues(siteID).Set(50)  

	m.nodeSwitchPortSpeedVec.WithLabelValues(siteID).Set(100) 
	m.nodeSwitchPortPowerVec.WithLabelValues(siteID).Set(50)  

	go func() {
		tickCount := 0
		for {
			tickCount++

			debugLog := tickCount%10 == 0

			m.mu.Lock()
			sitePortStatus, exists := m.portStatus[siteID]
			if debugLog {
				log.Infof(" [%s]: Port status check (tick %d): %v", siteID, tickCount, sitePortStatus)
			}
			m.mu.Unlock()

			if !exists {
				return
			}

			m.mu.Lock()
			currentUptime := m.siteUptimeCounters[siteID]
			m.siteUptimeSeconds.WithLabelValues(siteID).Set(float64(currentUptime))
			m.siteUptimeCounters[siteID]++
			m.mu.Unlock()

			if portEnabled, ok := sitePortStatus[PORT_BACKHAUL]; ok && portEnabled {
				if debugLog {
					log.Infof("[%s]: Updating backhaul metrics (enabled)", siteID)
				}
				m.backhaulSpeed.WithLabelValues(siteID).Set(float64(rand.Intn(100) + 30))           // 30-130 Mbps
				m.backhaulLatency.WithLabelValues(siteID).Set(float64(rand.Intn(100) + 20))         // 20-120 ms
				m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(float64(rand.Intn(100) + 10)) // 10-110 Mbps
				m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(float64(rand.Intn(100) + 10)) // 10-110 watts
			} else if debugLog {
				log.Infof(" [%s]: Current backhaul metrics: speed=%v, latency=%v", 
				    siteID, 
				    getSingleMetricValue(m.backhaulSpeed, siteID),
				    getSingleMetricValue(m.backhaulLatency, siteID))
			}

			if portEnabled, ok := sitePortStatus[PORT_SOLAR]; ok && portEnabled {
				if debugLog {
					log.Infof(" [%s]: Updating solar metrics (enabled)", siteID)
				}
				m.batteryChargePercentage.WithLabelValues(siteID).Set(float64(rand.Intn(30) + 50)) 
				m.solarPanelPower.WithLabelValues(siteID).Set(float64(rand.Intn(300) + 200))      
				m.solarPanelCurrent.WithLabelValues(siteID).Set(float64(rand.Intn(13) + 2))        
				m.solarPanelVoltage.WithLabelValues(siteID).Set(float64(rand.Intn(50) + 50))       
				m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(float64(rand.Intn(100) + 10))   
				m.solarSwitchPortPower.WithLabelValues(siteID).Set(float64(rand.Intn(100) + 10))   
			} else if debugLog {
				log.Infof(" [%s]: Skipping solar metrics update (disabled or not found)", siteID)
			}

			if portEnabled, ok := sitePortStatus[PORT_NODE]; ok && portEnabled {
				if debugLog {
					log.Infof(" [%s]: Updating node metrics (enabled)", siteID)
				}
				m.nodeSwitchPortPowerVec.WithLabelValues(siteID).Set(float64(rand.Intn(100) + 10)) 
				m.nodeSwitchPortSpeedVec.WithLabelValues(siteID).Set(float64(rand.Intn(100) + 10)) 
			} else if debugLog {
				log.Infof("DEBUG [%s]: Skipping node metrics update (disabled or not found)", siteID)
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

type MetricsManager struct {
	Metrics     *Metrics
	mu          sync.Mutex 
	ActiveSites map[string]bool
}

func NewMetricsManager() *MetricsManager {
	log.Println("DEBUG: Creating new MetricsManager")
	rand.Seed(time.Now().UnixNano())
	return &MetricsManager{
		Metrics:     New(),
		ActiveSites: make(map[string]bool),
	}
}

func (mm *MetricsManager) StartSiteMetrics(siteID string) error {
	log.Infof(" StartSiteMetrics called for site %s", siteID)
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if mm.ActiveSites[siteID] {
		return fmt.Errorf("metrics already running for site: %s", siteID)
	}

	mm.ActiveSites[siteID] = true
	mm.Metrics.StartMetricsGenerator(siteID)
	log.Printf("DEBUG: Started metrics for site %s", siteID)
	return nil
}

func (mm *MetricsManager) StopSiteMetrics(siteID string) {
	log.Infof(" StopSiteMetrics called for site %s", siteID)
	mm.mu.Lock()
	delete(mm.ActiveSites, siteID)
	mm.mu.Unlock()

	mm.Metrics.mu.Lock()
	delete(mm.Metrics.siteUptimeCounters, siteID)
	delete(mm.Metrics.portStatus, siteID)
	mm.Metrics.mu.Unlock()
}

func (mm *MetricsManager) UpdatePortStatus(siteID string, portNumber int, enabled bool) error {
	mm.mu.Lock()
	isRunning := mm.ActiveSites[siteID]
	mm.mu.Unlock()

	if !isRunning {
		return fmt.Errorf("no metrics running for site: %s", siteID)
	}

	err := mm.Metrics.UpdatePortStatus(siteID, portNumber, enabled)
	if err != nil {
		log.Printf(" Error updating port status: %v", err)
	} 
	return err
}

func (mm *MetricsManager) IsMetricsRunning(siteID string) bool {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	isRunning := mm.ActiveSites[siteID]
	return isRunning
}

func (mm *MetricsManager) GetSiteMetrics(siteID string) (map[string]float64, error) {
	mm.mu.Lock()
	isRunning := mm.ActiveSites[siteID]
	mm.mu.Unlock()

	if !isRunning {
		return nil, fmt.Errorf("no metrics running for site: %s", siteID)
	}

	metrics := make(map[string]float64)

	mm.Metrics.mu.Lock()
	portStatus, exists := mm.Metrics.portStatus[siteID]
	mm.Metrics.mu.Unlock()

	if !exists {
		return nil, fmt.Errorf("port status not found for site: %s", siteID)
	}


	gathered, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return nil, fmt.Errorf("failed to gather metrics: %v", err)
	}

	for _, mf := range gathered {
		for _, m := range mf.GetMetric() {
			siteValue := ""
			for _, label := range m.GetLabel() {
				if label.GetName() == "site" {
					siteValue = label.GetValue()
					break
				}
			}
			if siteValue == siteID {
				metricName := *mf.Name
				metricValue := m.GetGauge().GetValue()
				metrics[metricName] = metricValue
				
				if metricName == "backhaul_speed" || metricName == "main_backhaul_latency" ||
				   metricName == "backhaul_switch_port_status" || metricName == "backhaul_switch_port_power" ||
				   metricName == "backhaul_switch_port_speed" {
					log.Infof("DEBUG: Found metric %s = %v for site %s", metricName, metricValue, siteID)
				}
			}
		}
	}

log.Infof("portStatus for site %s: %v", siteID, portStatus)
	return metrics, nil
}