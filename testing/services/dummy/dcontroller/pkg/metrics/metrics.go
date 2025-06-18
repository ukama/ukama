/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package metrics

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/testing/services/dummy/dcontroller/pb/gen"
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
	
	siteNodeIDs    map[string]string
	siteNetworkIDs map[string]string
}

func New() *Metrics {
	log.Infof("Initializing Metrics struct")
	m := &Metrics{
		backhaulLatency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "main_backhaul_latency", Help: "Backhaul latency in milliseconds",
		}, []string{"site","node","network"}),
		backhaulSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_speed", Help: "Backhaul speed in Mbps",
		}, []string{"site","node","network"}),
		batteryChargePercentage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "battery_charge_percentage", Help: "Battery charge percentage",
		}, []string{"site","node","network"}),
		solarPanelVoltage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_voltage", Help: "Solar panel voltage in volts",
		}, []string{"site","node","network"}),
		solarPanelCurrent: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_current", Help: "Solar panel current in amperes",
		}, []string{"site","node","network"}),
		solarPanelPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_power", Help: "Solar panel power in watts",
		}, []string{"site","node","network"}),

		siteUptimeSeconds: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "site_uptime_seconds",
			Help: "Site current continuous uptime streak in seconds (resets on critical port down)",
		}, []string{"site","node","network"}),
		siteUptimePercentage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "site_uptime_percentage",
			Help: "Site historical uptime percentage (0-100%) since monitoring started",
		}, []string{"site","node","network"}),

		backhaulSwitchPortStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_switch_port_status", Help: "Backhaul switch port status (1 = up, 0 = down)",
		}, []string{"site","node","network"}),
		backhaulSwitchPortSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_switch_port_speed", Help: "Backhaul switch port speed in Mbps",
		}, []string{"site","node","network"}),
		backhaulSwitchPortPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_switch_port_power", Help: "Backhaul switch port power in watts",
		}, []string{"site","node","network"}),
		solarSwitchPortStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_switch_port_status", Help: "Solar switch port status (1 = up, 0 = down)",
		}, []string{"site","node","network"}),
		solarSwitchPortSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_switch_port_speed", Help: "Solar switch port speed in Mbps",
		}, []string{"site","node","network"}),
		solarSwitchPortPower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_switch_port_power", Help: "Solar switch port power in watts",
		}, []string{"site","node","network"}),
		nodeSwitchPortPowerVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_switch_port_power", Help: "Node switch port power in watts",
		}, []string{"site","node","network"}),
		nodeSwitchPortSpeedVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_switch_port_speed", Help: "Node switch port speed in Mbps",
		}, []string{"site","node","network"}),
		nodeSwitchPortStatusVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_switch_port_status", Help: "Node switch port status (1 = up, 0 = down)",
		}, []string{"site","node","network"}),

		siteUptimeStreakCounters: make(map[string]int64),   
		lastUptimeBeforeReset:    make(map[string]int64),  
		cumulativeUptimeSeconds:  make(map[string]float64),
		siteConfigs:              make(map[string]SiteConfig),
		batteryCharge:            make(map[string]float64),
		portStatus:               make(map[string]map[int]bool),
		siteStartTimes:           make(map[string]time.Time),
		
		siteNodeIDs:    make(map[string]string),
		siteNetworkIDs: make(map[string]string),
	}

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

func (m *Metrics) UpdatePortStatus(siteID string, portNumber int, enabled bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sitePorts, exists := m.portStatus[siteID]
	if !exists {
		log.Warnf("Attempted to update port status for unknown or stopped site: %s", siteID)
		return fmt.Errorf("port status map not initialized for site: %s", siteID)
	}

	nodeID, nodeExists := m.siteNodeIDs[siteID]
	networkID, networkExists := m.siteNetworkIDs[siteID]
	
	if !nodeExists || !networkExists {
		log.Errorf("Missing nodeID or networkID for site %s", siteID)
		return fmt.Errorf("missing nodeID or networkID for site: %s", siteID)
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
		m.nodeSwitchPortStatusVec.WithLabelValues(siteID, nodeID, networkID).Set(statusValue)
		if !enabled { 
			log.Infof("INFO: [%s] NODE port DOWN. Saving uptime streak before reset.", siteID)
			m.nodeSwitchPortSpeedVec.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.nodeSwitchPortPowerVec.WithLabelValues(siteID, nodeID, networkID).Set(0)
			
			m.lastUptimeBeforeReset[siteID] = m.siteUptimeStreakCounters[siteID]
			
			m.siteUptimeSeconds.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.siteUptimeStreakCounters[siteID] = 0
		} else { 
			log.Infof("INFO: [%s] NODE port UP.", siteID)
			if isBackhaulNowUp {
				if lastUptime, exists := m.lastUptimeBeforeReset[siteID]; exists && lastUptime > 0 {
					m.siteUptimeStreakCounters[siteID] = lastUptime
					m.siteUptimeSeconds.WithLabelValues(siteID, nodeID, networkID).Set(float64(lastUptime))
					log.Infof("INFO: [%s] Resuming uptime streak from %d seconds", siteID, lastUptime)
				}
			}
		}
	case PORT_SOLAR: 
		m.solarSwitchPortStatus.WithLabelValues(siteID, nodeID, networkID).Set(statusValue)
		if !enabled {
			log.Infof("[%s] SOLAR port down. Resetting related solar metrics.", siteID)
			m.solarSwitchPortSpeed.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.solarSwitchPortPower.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.solarPanelPower.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.solarPanelCurrent.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.solarPanelVoltage.WithLabelValues(siteID, nodeID, networkID).Set(0)
		} else {
			log.Infof("[%s] SOLAR port up.", siteID)
		}
	case PORT_BACKHAUL:
		m.backhaulSwitchPortStatus.WithLabelValues(siteID, nodeID, networkID).Set(statusValue)
		if !enabled { 
			log.Infof("INFO: [%s] BACKHAUL port DOWN. Saving uptime streak before reset.", siteID)
			m.backhaulSwitchPortSpeed.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.backhaulSwitchPortPower.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.backhaulSpeed.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.backhaulLatency.WithLabelValues(siteID, nodeID, networkID).Set(999) 
			m.lastUptimeBeforeReset[siteID] = m.siteUptimeStreakCounters[siteID]
			
			m.siteUptimeSeconds.WithLabelValues(siteID, nodeID, networkID).Set(0)
			m.siteUptimeStreakCounters[siteID] = 0
		} else { 
			log.Infof("INFO: [%s] BACKHAUL port UP.", siteID)
			if isNodeNowUp {
				if lastUptime, exists := m.lastUptimeBeforeReset[siteID]; exists && lastUptime > 0 {
					m.siteUptimeStreakCounters[siteID] = lastUptime
					m.siteUptimeSeconds.WithLabelValues(siteID, nodeID, networkID).Set(float64(lastUptime))
					log.Infof("INFO: [%s] Resuming uptime streak from %d seconds", siteID, lastUptime)
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

func (m *Metrics) StartMetricsGenerator(siteID string, config SiteConfig, nodeID string, networkID string) {
	log.Infof("DEBUG: Starting metrics generator for site %s with config %+v", siteID, config)
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
	m.siteNodeIDs[siteID] = nodeID
	m.siteNetworkIDs[siteID] = networkID
	startTime := m.siteStartTimes[siteID] 
	log.Infof(" Initialized site %s node %s network %s. Start Time: %s", siteID, nodeID, networkID, startTime.Format(time.RFC3339))
	m.mu.Unlock()

	m.backhaulSwitchPortStatus.WithLabelValues(siteID, nodeID, networkID).Set(1)
	m.solarSwitchPortStatus.WithLabelValues(siteID, nodeID, networkID).Set(1)
	m.nodeSwitchPortStatusVec.WithLabelValues(siteID, nodeID, networkID).Set(1)
	m.backhaulSwitchPortSpeed.WithLabelValues(siteID, nodeID, networkID).Set(config.AvgBackhaulSpeed)
	m.backhaulSwitchPortPower.WithLabelValues(siteID, nodeID, networkID).Set(5.0)
	m.solarSwitchPortSpeed.WithLabelValues(siteID, nodeID, networkID).Set(100)
	m.solarSwitchPortPower.WithLabelValues(siteID, nodeID, networkID).Set(5.0)
	m.nodeSwitchPortSpeedVec.WithLabelValues(siteID, nodeID, networkID).Set(100)
	m.nodeSwitchPortPowerVec.WithLabelValues(siteID, nodeID, networkID).Set(5.0)
	m.siteUptimeSeconds.WithLabelValues(siteID, nodeID, networkID).Set(0)   
	m.siteUptimePercentage.WithLabelValues(siteID, nodeID, networkID).Set(100.0) 
	m.batteryChargePercentage.WithLabelValues(siteID, nodeID, networkID).Set(m.batteryCharge[siteID])
	m.backhaulSpeed.WithLabelValues(siteID, nodeID, networkID).Set(config.AvgBackhaulSpeed)
	m.backhaulLatency.WithLabelValues(siteID, nodeID, networkID).Set(config.AvgLatency)

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
				m.siteUptimeSeconds.DeleteLabelValues(siteID, nodeID, networkID)
				m.siteUptimePercentage.DeleteLabelValues(siteID, nodeID, networkID)
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
				
				m.siteUptimeSeconds.WithLabelValues(siteID, nodeID, networkID).Set(float64(currentUptimeStreak))

				if debugLog {
					log.Infof("[%s] Site UP. Streak: %ds (Test Multiplier: %d), Cumulative: %.0fs",
						siteID, currentUptimeStreak, testUptimeMultiplier, currentCumulativeUptime)
				}
			} else {
				if currentUptimeStreak != 0 {
					log.Warnf("[%s] Site DOWN, but uptime streak counter was %d. Forcing reset.", siteID, currentUptimeStreak)
					currentUptimeStreak = 0
					m.siteUptimeStreakCounters[siteID] = 0
					m.siteUptimeSeconds.WithLabelValues(siteID, nodeID, networkID).Set(0)
					
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
			m.siteUptimePercentage.WithLabelValues(siteID, nodeID, networkID).Set(uptimePercentage)

			if debugLog {
				log.Infof("[%s] Uptime Calculation: TotalDuration=%.2fs, CumulativeUptime=%.2fs, Percentage=%.4f%%",
					siteID, totalDurationSeconds, currentCumulativeUptime, uptimePercentage)
			}

			if backhaulEnabled {
				speed := config.AvgBackhaulSpeed + float64(r.Intn(21)-10) 
				if speed < 0 { speed = 0 }
				latency := config.AvgLatency + float64(r.Intn(11)-5)     
				if latency < 1 { latency = 1 }

				m.backhaulSpeed.WithLabelValues(siteID, nodeID, networkID).Set(speed)
				m.backhaulLatency.WithLabelValues(siteID, nodeID, networkID).Set(latency)
				m.backhaulSwitchPortSpeed.WithLabelValues(siteID, nodeID, networkID).Set(speed * (0.95 + r.Float64()*0.1)) 
				m.backhaulSwitchPortPower.WithLabelValues(siteID, nodeID, networkID).Set(4.5 + r.Float64()*1.0)
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

				m.solarPanelPower.WithLabelValues(siteID, nodeID, networkID).Set(solarPower)
				voltage := 48.0 + float64(r.Intn(11)-5) 
				current := 0.0
				if voltage > 1 { current = solarPower / voltage }
				m.solarPanelVoltage.WithLabelValues(siteID, nodeID, networkID).Set(voltage)
				m.solarPanelCurrent.WithLabelValues(siteID, nodeID, networkID).Set(current)
				m.solarSwitchPortSpeed.WithLabelValues(siteID, nodeID, networkID).Set(100.0) 
				m.solarSwitchPortPower.WithLabelValues(siteID, nodeID, networkID).Set(4.5 + r.Float64()*1.0)
			} 

            chargeDelta := (netPower / 3600.0) * 0.1 
            currentBatteryCharge += chargeDelta
            currentBatteryCharge = math.Max(0.0, math.Min(100.0, currentBatteryCharge))

            m.batteryCharge[siteID] = currentBatteryCharge 
            m.batteryChargePercentage.WithLabelValues(siteID, nodeID, networkID).Set(currentBatteryCharge)

            if debugLog {
                 log.Infof("[%s] Battery Update: NetPower=%.2fW, ChargeDelta=%.4f%%, CurrentCharge=%.2f%%", siteID, netPower, chargeDelta*100, currentBatteryCharge)
            }

			if nodeEnabled {
				m.nodeSwitchPortPowerVec.WithLabelValues(siteID, nodeID, networkID).Set(consumption * (0.9 + r.Float64()*0.2)) 
				m.nodeSwitchPortSpeedVec.WithLabelValues(siteID, nodeID, networkID).Set(100.0)
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
	return &MetricsManager{
		Metrics:     New(),
		ActiveSites: make(map[string]bool),
	}
}

func (mm *MetricsManager) StartSiteMetrics(siteID string, config SiteConfig, nodeID string, networkID string) error {
	log.Infof("StartSiteMetrics called for site %s", siteID)
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if mm.ActiveSites[siteID] {
		log.Warnf("Attempted to start metrics for already active site: %s", siteID)
		return fmt.Errorf("metrics already running for site: %s", siteID)
	}

	mm.ActiveSites[siteID] = true
	log.Infof("Marked site %s as active in manager.", siteID)

	mm.Metrics.StartMetricsGenerator(siteID, config, nodeID, networkID)
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

	nodeID, nodeExists := mm.Metrics.siteNodeIDs[siteID]
	networkID, networkExists := mm.Metrics.siteNetworkIDs[siteID]
	
	delete(mm.Metrics.siteUptimeStreakCounters, siteID)
	delete(mm.Metrics.lastUptimeBeforeReset, siteID)
	delete(mm.Metrics.cumulativeUptimeSeconds, siteID)
	delete(mm.Metrics.portStatus, siteID)
	delete(mm.Metrics.siteConfigs, siteID)
	delete(mm.Metrics.batteryCharge, siteID)
	delete(mm.Metrics.siteStartTimes, siteID)
	delete(mm.Metrics.siteNodeIDs, siteID)      
	delete(mm.Metrics.siteNetworkIDs, siteID)  

	if nodeExists && networkExists {
		mm.Metrics.backhaulLatency.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.backhaulSpeed.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.batteryChargePercentage.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.solarPanelVoltage.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.solarPanelCurrent.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.solarPanelPower.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.siteUptimeSeconds.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.siteUptimePercentage.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.backhaulSwitchPortStatus.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.backhaulSwitchPortSpeed.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.backhaulSwitchPortPower.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.solarSwitchPortStatus.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.solarSwitchPortSpeed.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.solarSwitchPortPower.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.nodeSwitchPortPowerVec.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.nodeSwitchPortSpeedVec.DeleteLabelValues(siteID, nodeID, networkID)
		mm.Metrics.nodeSwitchPortStatusVec.DeleteLabelValues(siteID, nodeID, networkID)
	} else {
		log.Warnf("Could not find nodeID/networkID for site %s during cleanup. Some Prometheus labels may not be deleted properly.", siteID)
	}

	log.Infof("INFO: Cleaned up internal state and Prometheus labels for site %s.", siteID)
}


func (mm *MetricsManager) UpdateMetricsProfile(siteID string, profile pb.Profile) error {
    mm.mu.Lock()
    isActive := mm.ActiveSites[siteID]
    mm.mu.Unlock()

    if !isActive {
        return fmt.Errorf("site %s is not actively managed", siteID)
    }

    mm.Metrics.mu.Lock()
    defer mm.Metrics.mu.Unlock()

    config, exists := mm.Metrics.siteConfigs[siteID]
    if !exists {
        return fmt.Errorf("no config found for site %s", siteID)
    }

    oldConfig := config

    switch profile {
    case pb.Profile_PROFILE_MINI:
        config.AvgBackhaulSpeed = math.Max(5, config.AvgBackhaulSpeed * 0.5)
        config.AvgLatency = math.Min(100, config.AvgLatency * 2)
        config.SolarEfficiency = math.Max(0.1, config.SolarEfficiency * 0.5)
        
    case pb.Profile_PROFILE_MAX:
        config.AvgBackhaulSpeed = math.Min(1000, config.AvgBackhaulSpeed * 2)
        config.AvgLatency = math.Max(5, config.AvgLatency * 0.5)
        config.SolarEfficiency = math.Min(1.0, config.SolarEfficiency * 1.5)
        
    case pb.Profile_PROFILE_NORMAL:
        if config.AvgBackhaulSpeed < 20 || config.AvgBackhaulSpeed > 200 {
            config.AvgBackhaulSpeed = 100
        }
        if config.AvgLatency < 5 || config.AvgLatency > 100 {
            config.AvgLatency = 20
        }
        if config.SolarEfficiency < 0.3 || config.SolarEfficiency > 0.9 {
            config.SolarEfficiency = 0.7
        }
    }

    mm.Metrics.siteConfigs[siteID] = config
    
    log.Infof("[%s] Updated metrics profile to %s. Config changed from %+v to %+v", 
        siteID, profile.String(), oldConfig, config)
    
    return nil
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