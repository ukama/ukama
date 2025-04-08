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

type SiteConfig struct {
	AvgBackhaulSpeed float64 
	AvgLatency       float64 
	SolarEfficiency  float64 
}

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
	siteConfigs              map[string]SiteConfig 
	batteryCharge            map[string]float64   

	backhaulSwitchPortStatus *prometheus.GaugeVec
	backhaulSwitchPortSpeed  *prometheus.GaugeVec
	backhaulSwitchPortPower  *prometheus.GaugeVec
	solarSwitchPortStatus    *prometheus.GaugeVec
	solarSwitchPortSpeed     *prometheus.GaugeVec
	solarSwitchPortPower     *prometheus.GaugeVec
	nodeSwitchPortPowerVec   *prometheus.GaugeVec
	nodeSwitchPortSpeedVec   *prometheus.GaugeVec
	nodeSwitchPortStatusVec  *prometheus.GaugeVec

	portStatus map[string]map[int]bool
}

func New() *Metrics {
	log.Infof("Initializing Metrics struct")
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

		siteUptimeCounters: make(map[string]int),
		siteConfigs:        make(map[string]SiteConfig),
		batteryCharge:      make(map[string]float64),
		portStatus:         make(map[string]map[int]bool),
	}

	log.Println("DEBUG: Registering metrics with Prometheus")
	prometheus.MustRegister(m.backhaulLatency, m.backhaulSpeed, m.batteryChargePercentage,
		m.solarPanelVoltage, m.solarPanelCurrent, m.solarPanelPower, m.siteUptimeSeconds,
		m.backhaulSwitchPortStatus, m.backhaulSwitchPortSpeed, m.backhaulSwitchPortPower,
		m.solarSwitchPortStatus, m.solarSwitchPortSpeed, m.solarSwitchPortPower,
		m.nodeSwitchPortPowerVec, m.nodeSwitchPortSpeedVec, m.nodeSwitchPortStatusVec)

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
	defer m.mu.Unlock()

	if _, exists := m.portStatus[siteID]; !exists {
		m.portStatus[siteID] = make(map[int]bool)
	}

	m.portStatus[siteID][portNumber] = enabled
	nodeEnabled := m.portStatus[siteID][PORT_NODE]
	backhaulEnabled := m.portStatus[siteID][PORT_BACKHAUL]
	bothEnabled := nodeEnabled && backhaulEnabled

	if (portNumber == PORT_NODE || portNumber == PORT_BACKHAUL) && !enabled {
		m.siteUptimeCounters[siteID] = 0
		log.Printf("DEBUG: Resetting uptime counter to 0 for site %s", siteID)
	}

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
			m.siteUptimeSeconds.WithLabelValues(siteID).Set(0)
		} else {
			log.Printf("DEBUG: Setting node metrics to default values for site %s", siteID)
			m.nodeSwitchPortSpeedVec.WithLabelValues(siteID).Set(100)
			m.nodeSwitchPortPowerVec.WithLabelValues(siteID).Set(50)
			if bothEnabled {
				log.Printf("DEBUG: Both node and backhaul enabled for site %s", siteID)
			}
		}
	case PORT_SOLAR:
		m.solarSwitchPortStatus.WithLabelValues(siteID).Set(status)
		if !enabled {
			m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(0)
			m.solarSwitchPortPower.WithLabelValues(siteID).Set(0)
			m.solarPanelPower.WithLabelValues(siteID).Set(0)
			m.solarPanelCurrent.WithLabelValues(siteID).Set(0)
			m.solarPanelVoltage.WithLabelValues(siteID).Set(0)
		} else {
			m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(100)
			m.solarSwitchPortPower.WithLabelValues(siteID).Set(50)
		}
	case PORT_BACKHAUL:
		m.backhaulSwitchPortStatus.WithLabelValues(siteID).Set(status)
		if !enabled {
			log.Printf("DEBUG: Setting backhaul metrics to zero/high for site %s", siteID)
			m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(0)
			m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(0)
			m.backhaulSpeed.WithLabelValues(siteID).Set(0)
			m.backhaulLatency.WithLabelValues(siteID).Set(9999)
			m.siteUptimeSeconds.WithLabelValues(siteID).Set(0)
		} else {
			log.Printf("DEBUG: Setting backhaul metrics to default values for site %s", siteID)
			m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(100)
			m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(50)
			if bothEnabled {
				log.Printf("DEBUG: Both node and backhaul enabled for site %s", siteID)
			}
		}
	default:
		return fmt.Errorf("unknown port number: %d", portNumber)
	}
	return nil
}

func (m *Metrics) StartMetricsGenerator(siteID string, config SiteConfig) {
	log.Printf("DEBUG: Starting metrics generator for site %s with config %+v", siteID, config)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	m.mu.Lock()
	m.siteUptimeCounters[siteID] = r.Intn(6) + 5
	m.siteConfigs[siteID] = config
	m.batteryCharge[siteID] = 50.0 
	m.portStatus[siteID] = map[int]bool{
		PORT_NODE:     true,
		PORT_SOLAR:    true,
		PORT_BACKHAUL: true,
	}
	log.Printf("DEBUG: Initialized port status for site %s: %v", siteID, m.portStatus[siteID])
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
			config := m.siteConfigs[siteID]
			if debugLog {
				log.Infof(" [%s]: Port status check (tick %d): %v", siteID, tickCount, sitePortStatus)
			}
			m.mu.Unlock()

			if !exists {
				log.Printf("DEBUG: Site %s no longer exists, stopping metrics generator", siteID)
				return
			}

			nodeEnabled, nodeOk := sitePortStatus[PORT_NODE]
			backhaulEnabled, backhaulOk := sitePortStatus[PORT_BACKHAUL]
			if nodeOk && backhaulOk && nodeEnabled && backhaulEnabled {
				m.mu.Lock()
				currentUptime := m.siteUptimeCounters[siteID]
				m.siteUptimeSeconds.WithLabelValues(siteID).Set(float64(currentUptime))
				m.siteUptimeCounters[siteID]++
				if debugLog {
					log.Infof(" [%s]: Uptime: %d seconds", siteID, m.siteUptimeCounters[siteID])
				}
				m.mu.Unlock()
			}

			if portEnabled, ok := sitePortStatus[PORT_BACKHAUL]; ok && portEnabled {
				speed := config.AvgBackhaulSpeed + float64(r.Intn(20)-10) 
				latency := config.AvgLatency + float64(r.Intn(10)-5)    
				m.backhaulSpeed.WithLabelValues(siteID).Set(speed)
				m.backhaulLatency.WithLabelValues(siteID).Set(latency)
				m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(speed * 0.9) 
                m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(5.0 + float64(r.Intn(20))/10.0) 
                   
				if debugLog {
					log.Infof(" [%s]: Backhaul speed: %.2f Mbps, latency: %.2f ms", siteID, speed, latency)
				}
			}

			if portEnabled, ok := sitePortStatus[PORT_SOLAR]; ok && portEnabled {
				solarPower := (float64(r.Intn(300)+200) * config.SolarEfficiency) 
				consumption := 100.0                                              
				netPower := solarPower - consumption
				m.mu.Lock()
				m.batteryCharge[siteID] += netPower / 3600.0 
				if m.batteryCharge[siteID] > 100 {
					m.batteryCharge[siteID] = 100
				} else if m.batteryCharge[siteID] < 0 {
					m.batteryCharge[siteID] = 0
				}
				m.batteryChargePercentage.WithLabelValues(siteID).Set(m.batteryCharge[siteID])
				m.mu.Unlock()

				m.solarPanelPower.WithLabelValues(siteID).Set(solarPower)
				voltage := float64(r.Intn(50) + 50) 
				current := solarPower / voltage
				m.solarPanelVoltage.WithLabelValues(siteID).Set(voltage)
				m.solarPanelCurrent.WithLabelValues(siteID).Set(current)
				m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(float64(r.Intn(100) + 10))
                m.solarSwitchPortPower.WithLabelValues(siteID).Set(5.0 + float64(r.Intn(20))/10.0) 

				if debugLog {
					log.Infof(" [%s]: Solar power: %.2f W, Battery: %.2f%%", siteID, solarPower, m.batteryCharge[siteID])
				}
			} else {
				m.mu.Lock()
				m.batteryCharge[siteID] -= 100.0 / 3600.0 
				if m.batteryCharge[siteID] < 0 {
					m.batteryCharge[siteID] = 0
				}
				m.batteryChargePercentage.WithLabelValues(siteID).Set(m.batteryCharge[siteID])
				m.mu.Unlock()
			}

			if portEnabled, ok := sitePortStatus[PORT_NODE]; ok && portEnabled {
                m.nodeSwitchPortPowerVec.WithLabelValues(siteID).Set(5.0 + float64(r.Intn(20))/10.0) 
				m.nodeSwitchPortSpeedVec.WithLabelValues(siteID).Set(float64(r.Intn(100) + 10))
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
	return &MetricsManager{
		Metrics:     New(),
		ActiveSites: make(map[string]bool),
	}
}

func (mm *MetricsManager) StartSiteMetrics(siteID string, config SiteConfig) error {
	log.Infof("DEBUG: StartSiteMetrics called for site %s", siteID)
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if mm.ActiveSites[siteID] {
		return fmt.Errorf("metrics already running for site: %s", siteID)
	}

	mm.ActiveSites[siteID] = true
	mm.Metrics.StartMetricsGenerator(siteID, config)
	log.Printf("DEBUG: Started metrics for site %s, active sites: %v", siteID, mm.ActiveSites)
	return nil
}

func (mm *MetricsManager) StopSiteMetrics(siteID string) {
	log.Infof("DEBUG: StopSiteMetrics called for site %s", siteID)
	mm.mu.Lock()
	delete(mm.ActiveSites, siteID)
	mm.mu.Unlock()

	mm.Metrics.mu.Lock()
	delete(mm.Metrics.siteUptimeCounters, siteID)
	delete(mm.Metrics.portStatus, siteID)
	delete(mm.Metrics.siteConfigs, siteID)
	delete(mm.Metrics.batteryCharge, siteID)
	mm.Metrics.mu.Unlock()
	log.Printf("DEBUG: Stopped metrics for site %s", siteID)
}

func (mm *MetricsManager) UpdatePortStatus(siteID string, portNumber int, enabled bool) error {
	mm.mu.Lock()
	isRunning := mm.ActiveSites[siteID]
	mm.mu.Unlock()

	if !isRunning {
		return fmt.Errorf("no metrics running for site: %s", siteID)
	}
	return mm.Metrics.UpdatePortStatus(siteID, portNumber, enabled)
}

func (mm *MetricsManager) IsMetricsRunning(siteID string) bool {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	return mm.ActiveSites[siteID]
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
			}
		}
	}

	log.Infof("DEBUG: Port status for site %s: %v", siteID, portStatus)
	return metrics, nil
}