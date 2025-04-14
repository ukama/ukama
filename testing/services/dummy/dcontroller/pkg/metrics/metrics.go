package metrics

import (
	"fmt"
	"math"
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
	siteUptimePercentage     *prometheus.GaugeVec 
	switchPortStatus         *prometheus.GaugeVec 
	switchPortSpeed          *prometheus.GaugeVec 
	switchPortPower          *prometheus.GaugeVec 
	mu                       sync.Mutex
	siteUptimeStreakCounters map[string]int64       
	lastUptimeBeforeReset    map[string]int64       
	siteConfigs              map[string]SiteConfig
	batteryCharge            map[string]float64
	siteStartTimes           map[string]time.Time 
	cumulativeUptimeSeconds  map[string]float64   

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
			Name: "main_backhaul_latency", Help: "Backhaul latency in milliseconds",
		}, []string{"site"}),
		backhaulSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_speed", Help: "Backhaul speed in Mbps",
		}, []string{"site"}),
		batteryChargePercentage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "battery_charge_percentage", Help: "Battery charge percentage",
		}, []string{"site"}),
		solarPanelVoltage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_voltage", Help: "Solar panel voltage in volts",
		}, []string{"site"}),
		solarPanelCurrent: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_current", Help: "Solar panel current in amperes",
		}, []string{"site"}),
		solarPanelPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_power", Help: "Solar panel power in watts",
		}, []string{"site"}),

		siteUptimeSeconds: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "site_uptime_seconds",
			Help: "Site current continuous uptime streak in seconds (resets on critical port down)",
		}, []string{"site"}),
		siteUptimePercentage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "site_uptime_percentage",
			Help: "Site historical uptime percentage (0-100%) since monitoring started",
		}, []string{"site"}),

		backhaulSwitchPortStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_switch_port_status", Help: "Backhaul switch port status (1 = up, 0 = down)",
		}, []string{"site"}),
		backhaulSwitchPortSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_switch_port_speed", Help: "Backhaul switch port speed in Mbps",
		}, []string{"site"}),
		backhaulSwitchPortPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_switch_port_power", Help: "Backhaul switch port power in watts",
		}, []string{"site"}),
		solarSwitchPortStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_switch_port_status", Help: "Solar switch port status (1 = up, 0 = down)",
		}, []string{"site"}),
		solarSwitchPortSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_switch_port_speed", Help: "Solar switch port speed in Mbps",
		}, []string{"site"}),
		solarSwitchPortPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_switch_port_power", Help: "Solar switch port power in watts",
		}, []string{"site"}),
		nodeSwitchPortPowerVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_switch_port_power", Help: "Node switch port power in watts",
		}, []string{"site"}),
		nodeSwitchPortSpeedVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_switch_port_speed", Help: "Node switch port speed in Mbps",
		}, []string{"site"}),
		nodeSwitchPortStatusVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_switch_port_status", Help: "Node switch port status (1 = up, 0 = down)",
		}, []string{"site"}),

		siteUptimeStreakCounters: make(map[string]int64),   
		lastUptimeBeforeReset:    make(map[string]int64),  
		cumulativeUptimeSeconds:  make(map[string]float64),
		siteConfigs:              make(map[string]SiteConfig),
		batteryCharge:            make(map[string]float64),
		portStatus:               make(map[string]map[int]bool),
		siteStartTimes:           make(map[string]time.Time), 
	}

	log.Println("DEBUG: Registering metrics with Prometheus")
	prometheus.MustRegister(
		m.backhaulLatency, m.backhaulSpeed, m.batteryChargePercentage,
		m.solarPanelVoltage, m.solarPanelCurrent, m.solarPanelPower,
		m.siteUptimeSeconds, m.siteUptimePercentage, 
		m.backhaulSwitchPortStatus, m.backhaulSwitchPortSpeed, m.backhaulSwitchPortPower,
		m.solarSwitchPortStatus, m.solarSwitchPortSpeed, m.solarSwitchPortPower,
		m.nodeSwitchPortPowerVec, m.nodeSwitchPortSpeedVec, m.nodeSwitchPortStatusVec,
	)
	return m
}

func getSingleMetricValue(gauge *prometheus.GaugeVec, siteID string) float64 {
    metric, err := gauge.GetMetricWithLabelValues(siteID)
    if err != nil {
        log.Warnf("Error getting metric for site %s: %v", siteID, err)
        return -999
    }
    ch := make(chan prometheus.Metric, 1)
    metric.Collect(ch)
    select {
    case m := <-ch:
        var metricOut dto.Metric
        if err := m.Write(&metricOut); err != nil {
            log.Errorf("Error writing metric to DTO for site %s: %v", siteID, err)
            return -999
        }
        if metricOut.Gauge != nil && metricOut.Gauge.Value != nil {
            return *metricOut.Gauge.Value
        }
        log.Warnf("Metric collected for site %s is not a gauge or has no value", siteID)
        return -999
    case <-time.After(1 * time.Second):
        log.Warnf("Timeout getting metric value for site %s", siteID)
        return -999
    }
}


func (m *Metrics) UpdatePortStatus(siteID string, portNumber int, enabled bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sitePorts, exists := m.portStatus[siteID]
	if !exists {
		log.Warnf("Attempted to update port status for unknown or stopped site: %s", siteID)
		return fmt.Errorf("port status map not initialized for site: %s", siteID)
	}

	wasNodeUp := sitePorts[PORT_NODE]
	wasBackhaulUp := sitePorts[PORT_BACKHAUL]
	wasSiteUp := wasNodeUp && wasBackhaulUp

	m.portStatus[siteID][portNumber] = enabled

	isNodeNowUp := m.portStatus[siteID][PORT_NODE]
	isBackhaulNowUp := m.portStatus[siteID][PORT_BACKHAUL]
	isSiteNowUp := isNodeNowUp && isBackhaulNowUp

	statusValue := 0.0
	if enabled {
		statusValue = 1.0
	}

	switch portNumber {
	case PORT_NODE:
		m.nodeSwitchPortStatusVec.WithLabelValues(siteID).Set(statusValue)
		if !enabled { 
			log.Printf("INFO: [%s] NODE port DOWN. Saving uptime streak before reset.", siteID)
			m.nodeSwitchPortSpeedVec.WithLabelValues(siteID).Set(0)
			m.nodeSwitchPortPowerVec.WithLabelValues(siteID).Set(0)
			
			m.lastUptimeBeforeReset[siteID] = m.siteUptimeStreakCounters[siteID]
			
			m.siteUptimeSeconds.WithLabelValues(siteID).Set(0)
			m.siteUptimeStreakCounters[siteID] = 0
		} else { 
			log.Printf("INFO: [%s] NODE port UP.", siteID)
			if isBackhaulNowUp {
				if lastUptime, exists := m.lastUptimeBeforeReset[siteID]; exists && lastUptime > 0 {
					m.siteUptimeStreakCounters[siteID] = lastUptime
					m.siteUptimeSeconds.WithLabelValues(siteID).Set(float64(lastUptime))
					log.Printf("INFO: [%s] Resuming uptime streak from %d seconds", siteID, lastUptime)
				}
			}
		}
	case PORT_SOLAR: 
		m.solarSwitchPortStatus.WithLabelValues(siteID).Set(statusValue)
		if !enabled {
			log.Infof("[%s] SOLAR port down. Resetting related solar metrics.", siteID)
			m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(0)
			m.solarSwitchPortPower.WithLabelValues(siteID).Set(0)
			m.solarPanelPower.WithLabelValues(siteID).Set(0)
			m.solarPanelCurrent.WithLabelValues(siteID).Set(0)
			m.solarPanelVoltage.WithLabelValues(siteID).Set(0)
		} else {
			log.Infof("[%s] SOLAR port up.", siteID)
		}
	case PORT_BACKHAUL:
		m.backhaulSwitchPortStatus.WithLabelValues(siteID).Set(statusValue)
		if !enabled { 
			log.Printf("INFO: [%s] BACKHAUL port DOWN. Saving uptime streak before reset.", siteID)
			m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(0)
			m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(0)
			m.backhaulSpeed.WithLabelValues(siteID).Set(0)
			m.backhaulLatency.WithLabelValues(siteID).Set(999) 
			m.lastUptimeBeforeReset[siteID] = m.siteUptimeStreakCounters[siteID]
			
			m.siteUptimeSeconds.WithLabelValues(siteID).Set(0)
			m.siteUptimeStreakCounters[siteID] = 0
		} else { 
			log.Printf("INFO: [%s] BACKHAUL port UP.", siteID)
			if isNodeNowUp {
				if lastUptime, exists := m.lastUptimeBeforeReset[siteID]; exists && lastUptime > 0 {
					m.siteUptimeStreakCounters[siteID] = lastUptime
					m.siteUptimeSeconds.WithLabelValues(siteID).Set(float64(lastUptime))
					log.Printf("INFO: [%s] Resuming uptime streak from %d seconds", siteID, lastUptime)
				}
			}
		}
	default:
		return fmt.Errorf("unknown port number: %d for site: %s", portNumber, siteID)
	}

	if wasSiteUp && !isSiteNowUp {
		log.Infof("INFO: [%s] Site considered DOWN due to port status change.", siteID)
	} else if !wasSiteUp && isSiteNowUp {
		log.Infof("INFO: [%s] Site considered UP. Uptime streak resumed from %d seconds.", 
			siteID, m.siteUptimeStreakCounters[siteID])
	}

	return nil
}

func (m *Metrics) StartMetricsGenerator(siteID string, config SiteConfig) {
	log.Printf("DEBUG: Starting metrics generator for site %s with config %+v", siteID, config)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	m.mu.Lock()
	m.siteUptimeStreakCounters[siteID] = 0       
	m.lastUptimeBeforeReset[siteID] = 0           
	m.cumulativeUptimeSeconds[siteID] = 0.0      
	m.siteConfigs[siteID] = config
	m.batteryCharge[siteID] = 50.0 + float64(r.Intn(21)-10) 
	m.portStatus[siteID] = map[int]bool{
		PORT_NODE:     true, 
		PORT_SOLAR:    true,
		PORT_BACKHAUL: true,
	}
	m.siteStartTimes[siteID] = time.Now()
	startTime := m.siteStartTimes[siteID] 
	log.Printf(" Initialized site %s. Start Time: %s", siteID, startTime.Format(time.RFC3339))
	log.Printf("DEBUG: Initial port status for site %s: Node:%t, Solar:%t, Backhaul:%t",
		siteID, m.portStatus[siteID][PORT_NODE], m.portStatus[siteID][PORT_SOLAR], m.portStatus[siteID][PORT_BACKHAUL])
	m.mu.Unlock()

	m.backhaulSwitchPortStatus.WithLabelValues(siteID).Set(1)
	m.solarSwitchPortStatus.WithLabelValues(siteID).Set(1)
	m.nodeSwitchPortStatusVec.WithLabelValues(siteID).Set(1)
	m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(config.AvgBackhaulSpeed)
	m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(5.0)
	m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(100)
	m.solarSwitchPortPower.WithLabelValues(siteID).Set(5.0)
	m.nodeSwitchPortSpeedVec.WithLabelValues(siteID).Set(100)
	m.nodeSwitchPortPowerVec.WithLabelValues(siteID).Set(5.0)
	m.siteUptimeSeconds.WithLabelValues(siteID).Set(0)   
	m.siteUptimePercentage.WithLabelValues(siteID).Set(100.0) 
	m.batteryChargePercentage.WithLabelValues(siteID).Set(m.batteryCharge[siteID])
	m.backhaulSpeed.WithLabelValues(siteID).Set(config.AvgBackhaulSpeed)
	m.backhaulLatency.WithLabelValues(siteID).Set(config.AvgLatency)

	go func() {
		ticker := time.NewTicker(1 * time.Second) 
		defer ticker.Stop()

		tickCount := 0                            
		const logIntervalSeconds = 60            
		const testUptimeMultiplier = 1 

		if testUptimeMultiplier > 1 {
			log.Infof("[%s] TEST MODE ACTIVE: Uptime streak accumulating %d times faster!", siteID, testUptimeMultiplier)
		}

		for range ticker.C {
			tickCount++
			debugLog := (tickCount % logIntervalSeconds) == 1 

			m.mu.Lock() 

			sitePortStatus, exists := m.portStatus[siteID]
			if !exists {
				m.mu.Unlock()
				log.Infof("INFO: [%s] Site no longer active, stopping metrics generator.", siteID)
				return 
			}
			config, configExists := m.siteConfigs[siteID]
			startTime, startTimeExists := m.siteStartTimes[siteID]
			currentCumulativeUptime := m.cumulativeUptimeSeconds[siteID]
			currentUptimeStreak := m.siteUptimeStreakCounters[siteID]
			currentBatteryCharge := m.batteryCharge[siteID]

			if !configExists || !startTimeExists {
				m.mu.Unlock()
				log.Errorf("ERROR: [%s] Missing config or start time in generator loop. Stopping.", siteID)
				m.siteUptimeSeconds.DeleteLabelValues(siteID)
				m.siteUptimePercentage.DeleteLabelValues(siteID)
				return
			}

			nodeEnabled, nodeOk := sitePortStatus[PORT_NODE]
			backhaulEnabled, backhaulOk := sitePortStatus[PORT_BACKHAUL]
			isSiteUp := nodeOk && backhaulOk && nodeEnabled && backhaulEnabled

			if isSiteUp {
				
				streakIncrement := int64(1 * testUptimeMultiplier)
				currentUptimeStreak += streakIncrement 

				currentCumulativeUptime += 1.0 

				m.siteUptimeStreakCounters[siteID] = currentUptimeStreak    
				m.cumulativeUptimeSeconds[siteID] = currentCumulativeUptime 
				
				m.siteUptimeSeconds.WithLabelValues(siteID).Set(float64(currentUptimeStreak))

				if debugLog {
					log.Infof("[%s] Site UP. Streak: %ds (Test Multiplier: %d), Cumulative: %.0fs",
						siteID, currentUptimeStreak, testUptimeMultiplier, currentCumulativeUptime)
				}
			} else {
				
				if currentUptimeStreak != 0 {
					log.Warnf("[%s] Site DOWN, but uptime streak counter was %d. Forcing reset.", siteID, currentUptimeStreak)
					currentUptimeStreak = 0
					m.siteUptimeStreakCounters[siteID] = 0
					m.siteUptimeSeconds.WithLabelValues(siteID).Set(0)
					
					if _, exists := m.lastUptimeBeforeReset[siteID]; !exists {
						m.lastUptimeBeforeReset[siteID] = currentUptimeStreak
					}
				}
				if debugLog {
					log.Infof("[%s] Site DOWN. Streak: %ds, Cumulative: %.0fs, Last Saved: %ds",
						siteID, currentUptimeStreak, currentCumulativeUptime, m.lastUptimeBeforeReset[siteID])
				}
			}

			totalDurationSeconds := time.Since(startTime).Seconds()
			var uptimePercentage float64

			if totalDurationSeconds > 0.001 { 
				uptimePercentage = (currentCumulativeUptime / totalDurationSeconds) * 100.0
			} else if currentCumulativeUptime > 0 { 
				 uptimePercentage = 100.0
			} else { 
                 uptimePercentage = 0.0 
            }


			uptimePercentage = math.Max(0.0, math.Min(100.0, uptimePercentage))

			m.siteUptimePercentage.WithLabelValues(siteID).Set(uptimePercentage)

			if debugLog {
				log.Infof("[%s] Uptime Calculation: TotalDuration=%.2fs, CumulativeUptime=%.2fs, Percentage=%.4f%%",
					siteID, totalDurationSeconds, currentCumulativeUptime, uptimePercentage)
			}

			if backhaulEnabled {
				speed := config.AvgBackhaulSpeed + float64(r.Intn(21)-10) 
				if speed < 0 { speed = 0 }
				latency := config.AvgLatency + float64(r.Intn(11)-5)     
				if latency < 1 { latency = 1 }

				m.backhaulSpeed.WithLabelValues(siteID).Set(speed)
				m.backhaulLatency.WithLabelValues(siteID).Set(latency)
				m.backhaulSwitchPortSpeed.WithLabelValues(siteID).Set(speed * (0.95 + r.Float64()*0.1)) 
				m.backhaulSwitchPortPower.WithLabelValues(siteID).Set(4.5 + r.Float64()*1.0)

			
			} 

			solarEnabled, solarOk := sitePortStatus[PORT_SOLAR]
			netPower := 0.0 
            consumption := 0.0
			if nodeEnabled { 
                consumption = 75.0 + float64(r.Intn(51)-25) 
                netPower -= consumption
			}


			if solarOk && solarEnabled {
				solarPower := math.Max(0, (200.0 + float64(r.Intn(301))) * config.SolarEfficiency) 
				netPower += solarPower

				m.solarPanelPower.WithLabelValues(siteID).Set(solarPower)
				voltage := 48.0 + float64(r.Intn(11)-5) 
				current := 0.0
				if voltage > 1 { current = solarPower / voltage }
				m.solarPanelVoltage.WithLabelValues(siteID).Set(voltage)
				m.solarPanelCurrent.WithLabelValues(siteID).Set(current)
				m.solarSwitchPortSpeed.WithLabelValues(siteID).Set(100.0) 
				m.solarSwitchPortPower.WithLabelValues(siteID).Set(4.5 + r.Float64()*1.0)

			} 

            chargeDelta := (netPower / 3600.0) * 0.1 
            currentBatteryCharge += chargeDelta

            currentBatteryCharge = math.Max(0.0, math.Min(100.0, currentBatteryCharge))

            m.batteryCharge[siteID] = currentBatteryCharge 
            m.batteryChargePercentage.WithLabelValues(siteID).Set(currentBatteryCharge)

            if debugLog {
                 log.Infof("[%s] Battery Update: NetPower=%.2fW, ChargeDelta=%.4f%%, CurrentCharge=%.2f%%", siteID, netPower, chargeDelta*100, currentBatteryCharge)
            }


			if nodeEnabled {
				m.nodeSwitchPortPowerVec.WithLabelValues(siteID).Set(consumption * (0.9 + r.Float64()*0.2)) 
				m.nodeSwitchPortSpeedVec.WithLabelValues(siteID).Set(100.0)
				
			} 

			m.mu.Unlock() 

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
	log.Infof("StartSiteMetrics called for site %s", siteID)
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if mm.ActiveSites[siteID] {
		log.Warnf("Attempted to start metrics for already active site: %s", siteID)
		return fmt.Errorf("metrics already running for site: %s", siteID)
	}

	mm.ActiveSites[siteID] = true
	log.Infof("Marked site %s as active in manager.", siteID)

	
	mm.Metrics.StartMetricsGenerator(siteID, config)
	log.Infof("INFO: Started metrics generator for site %s.", siteID)
	return nil
}

func (mm *MetricsManager) StopSiteMetrics(siteID string) {
	log.Infof("INFO: Attempting to stop metrics and clean up for site %s", siteID)

	mm.mu.Lock()
	isActive := mm.ActiveSites[siteID]
	if isActive {
		delete(mm.ActiveSites, siteID)
		log.Infof("INFO: Marked site %s as inactive in manager.", siteID)
	} else {
		log.Warnf("Attempted to stop metrics for already inactive site: %s. Proceeding with cleanup.", siteID)
	}
	mm.mu.Unlock() 
	mm.Metrics.mu.Lock()
	defer mm.Metrics.mu.Unlock()

	delete(mm.Metrics.siteUptimeStreakCounters, siteID)
	delete(mm.Metrics.lastUptimeBeforeReset, siteID)    
	delete(mm.Metrics.cumulativeUptimeSeconds, siteID) 
	delete(mm.Metrics.portStatus, siteID)               
	delete(mm.Metrics.siteConfigs, siteID)
	delete(mm.Metrics.batteryCharge, siteID)
	delete(mm.Metrics.siteStartTimes, siteID)          

	mm.Metrics.backhaulLatency.DeleteLabelValues(siteID)
	mm.Metrics.backhaulSpeed.DeleteLabelValues(siteID)
	mm.Metrics.batteryChargePercentage.DeleteLabelValues(siteID)
	mm.Metrics.solarPanelVoltage.DeleteLabelValues(siteID)
	mm.Metrics.solarPanelCurrent.DeleteLabelValues(siteID)
	mm.Metrics.solarPanelPower.DeleteLabelValues(siteID)
	mm.Metrics.siteUptimeSeconds.DeleteLabelValues(siteID)     
	mm.Metrics.siteUptimePercentage.DeleteLabelValues(siteID)   
	mm.Metrics.backhaulSwitchPortStatus.DeleteLabelValues(siteID)
	mm.Metrics.backhaulSwitchPortSpeed.DeleteLabelValues(siteID)
	mm.Metrics.backhaulSwitchPortPower.DeleteLabelValues(siteID)
	mm.Metrics.solarSwitchPortStatus.DeleteLabelValues(siteID)
	mm.Metrics.solarSwitchPortSpeed.DeleteLabelValues(siteID)
	mm.Metrics.solarSwitchPortPower.DeleteLabelValues(siteID)
	mm.Metrics.nodeSwitchPortPowerVec.DeleteLabelValues(siteID)
	mm.Metrics.nodeSwitchPortSpeedVec.DeleteLabelValues(siteID)
	mm.Metrics.nodeSwitchPortStatusVec.DeleteLabelValues(siteID)

	log.Infof("INFO: Cleaned up internal state and Prometheus labels for site %s.", siteID)
}


func (mm *MetricsManager) UpdatePortStatus(siteID string, portNumber int, enabled bool) error {
	mm.mu.Lock()
	isActive := mm.ActiveSites[siteID]
	mm.mu.Unlock()

	if !isActive {
		log.Infof("UpdatePortStatus called for site %s which is not actively managed.", siteID)
	
	}

	err := mm.Metrics.UpdatePortStatus(siteID, portNumber, enabled)
	if err != nil {
		log.Errorf("Error updating port status for site %s, port %d: %v", siteID, portNumber, err)
	}
	return err 
}


func (mm *MetricsManager) IsMetricsRunning(siteID string) bool {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	isRunning, exists := mm.ActiveSites[siteID]
	return exists && isRunning
}

func (mm *MetricsManager) GetSiteMetrics(siteID string) (map[string]float64, error) {

	metrics := make(map[string]float64)
	gathered, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		log.Errorf("[%s] Failed to gather Prometheus metrics: %v", siteID, err)
		return nil, fmt.Errorf("failed to gather metrics: %v", err)
	}

	foundSiteMetrics := false
	for _, mf := range gathered {
		metricName := mf.GetName()
		for _, m := range mf.GetMetric() {
			isTargetSite := false
			for _, label := range m.GetLabel() {
				if label.GetName() == "site" && label.GetValue() == siteID {
					isTargetSite = true
					break
				}
			}

			if isTargetSite {
				foundSiteMetrics = true 
				var value float64
				var ok bool

				if gg := m.GetGauge(); gg != nil {
					value = gg.GetValue()
					ok = true
				} else if ct := m.GetCounter(); ct != nil {
					value = ct.GetValue()
					ok = true
				} 

				if ok {
					metrics[metricName] = value
				} else {
					log.Warnf("[%s] Metric '%s' found but has an unsupported type or nil value", siteID, metricName)
				}
			}
		}
	}

	if !foundSiteMetrics {

		log.Infof("[%s] No Prometheus metrics found with the site label.", siteID)
	}

	return metrics, nil
}