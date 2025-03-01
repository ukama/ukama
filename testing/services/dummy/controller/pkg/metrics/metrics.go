/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package metrics

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	cenums "github.com/ukama/ukama/testing/common/enums"
)

type PrometheusExporter struct {
	// Solar metrics
	solarPowerGeneration *prometheus.GaugeVec
	solarEnergyTotal    *prometheus.GaugeVec
	solarPanelPower     *prometheus.GaugeVec
	solarPanelCurrent   *prometheus.GaugeVec
	solarPanelVoltage   *prometheus.GaugeVec
	solarInverterStatus *prometheus.GaugeVec

	// Battery metrics
	batteryChargeStatus *prometheus.GaugeVec
	batteryVoltage      *prometheus.GaugeVec
	batteryHealth       *prometheus.GaugeVec
	batteryCurrent      *prometheus.GaugeVec
	batteryTemperature  *prometheus.GaugeVec

	// Network metrics
	backhaulLatency      *prometheus.GaugeVec
	backhaulStatus       *prometheus.GaugeVec
	backhaulSpeed        *prometheus.GaugeVec
	switchPortStatus     *prometheus.GaugeVec
	switchPortBandwidth  *prometheus.GaugeVec

	metricsProvider *MetricsProvider
	siteId         string
	shutdown       chan struct{} 
}

func NewPrometheusExporter(metricsProvider *MetricsProvider, siteId string) *PrometheusExporter {
	exporter := &PrometheusExporter{
		metricsProvider: metricsProvider,
		siteId:          siteId,
		shutdown:       make(chan struct{}), 

		solarPowerGeneration: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_power_generation",
			Help: "Current solar power generation in watts",
		}, []string{"unit","site"}),

		solarEnergyTotal: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_energy_total",
			Help: "Total solar energy generated in kilowatt-hours",
		}, []string{"unit","site"}),

		solarPanelPower: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_power",
			Help: "Current solar panel power in watts",
		}, []string{"unit","site"}),
		solarPanelCurrent: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_current",
			Help: "Current solar panel current in amperes",
		}, []string{"unit","site"}),
		solarPanelVoltage: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_panel_voltage",
			Help: "Current solar panel voltage in volts",
		}, []string{"unit","site"}),
		solarInverterStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "solar_inverter_status",
			Help: "Solar inverter status (1 = working, 0 = not working)",
		}, []string{"unit","site"}),

		// Battery metrics
		batteryChargeStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "battery_charge_status",
			Help: "Battery charge status in percentage",
		}, []string{"unit","site"}),
		batteryVoltage: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "battery_voltage_volts",
			Help: "Battery voltage in volts",
		}, []string{"unit","site"}),
		batteryHealth: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "battery_health",
			Help: "Battery health status (1 = good, 0 = poor)",
		}, []string{"unit","site"}),
		batteryCurrent: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "battery_current",
			Help: "Battery current in amperes",
		}, []string{"unit","site"}),
		batteryTemperature: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "battery_temperature",
			Help: "Battery temperature in Celsius",
		}, []string{"unit","site"}),

		// Network metrics
		backhaulLatency: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_latency",
			Help: "Backhaul latency in milliseconds",
		}, []string{"unit","site"}),
		backhaulStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_status",
			Help: "Backhaul status (1 = up, 0 = down)",
		}, []string{"unit","site"}),
		backhaulSpeed: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "backhaul_speed",
			Help: "Backhaul speed in Mbps",
		}, []string{"unit","site"}),
		switchPortStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "switch_port_status",
			Help: "Switch port status (1 = up, 0 = down)",
		}, []string{"unit","site"}),
		switchPortBandwidth: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "switch_port_bandwidth",
			Help: "Switch port bandwidth in Mbps",
		}, []string{"unit","site"}),
	}

	return exporter
}

type BackhaulMetrics struct {
	Latency float64
	Status  float64
	Speed   float64
	SwitchStatus  float64  
	SwitchBandwidth float64 
}

type BatteryMetrics struct {
	Capacity    float64
	Voltage     float64
	Current     float64
	Temperature float64
	Status      string
	Health      string
}

type SolarMetrics struct {
	PowerGeneration float64
	EnergyTotal    float64
	PanelPower     float64
	PanelCurrent   float64
	PanelVoltage   float64
	InverterStatus float64
}

type ControllerMetrics struct {
	Backhaul *BackhaulMetrics
	Battery  *BatteryMetrics
	Solar    *SolarMetrics
	Time     time.Time
}

type BackhaulProvider struct {
	lastUpdate    time.Time
	metricsProvider *MetricsProvider
	jitterFactor float64
	noiseAmplitude float64  
}

const (
	PORT_AMPLIFIER = 1
	PORT_TOWER     = 2
	PORT_SOLAR     = 3
	PORT_BACKHAUL  = 4
)

type BatteryProvider struct {
	startTime    time.Time
	lastCapacity float64
	cycleCount   int
	metricsProvider *MetricsProvider
	microFluctuationFactor float64
	temperatureVariance float64
}

type SolarProvider struct {
	startTime      time.Time
	energyTotal    float64
	weatherPattern float64
	metricsProvider *MetricsProvider
	cloudCoverFactor float64
	microWeatherChanges float64
	timeAcceleration float64
}

type MetricsProvider struct {
	backhaul       *BackhaulProvider
	battery        *BatteryProvider
	solar          *SolarProvider
	portStatus     map[int]bool  
	scenarioActive string        
	currentProfile cenums.Profile
	globalVariationFactor float64
	timeMultiplier float64
}

func (b *BackhaulProvider) UpdateMetricsProvider(provider *MetricsProvider) {
	b.metricsProvider = provider
}

func (b *BatteryProvider) UpdateMetricsProvider(provider *MetricsProvider) {
	b.metricsProvider = provider
}

func (s *SolarProvider) UpdateMetricsProvider(provider *MetricsProvider) {
	s.metricsProvider = provider
}

func NewMetricsProvider() *MetricsProvider {
	mp := &MetricsProvider{
		backhaul: &BackhaulProvider{
			lastUpdate: time.Now(),
			jitterFactor: rand.Float64() * 0.5,
			noiseAmplitude: rand.Float64() * 0.3,
		},
		battery: &BatteryProvider{
			startTime: time.Now(),
			lastCapacity: 85.0,
			cycleCount: 0,
			microFluctuationFactor: rand.Float64() * 0.2,
			temperatureVariance: rand.Float64() * 2.0,
		},
		solar: &SolarProvider{
			startTime: time.Now(),
			energyTotal: 0,
			weatherPattern: 1.0,
			cloudCoverFactor: rand.Float64() * 0.4,
			microWeatherChanges: rand.Float64() * 0.3,
			timeAcceleration: 1.0, 
		},
		portStatus: map[int]bool{
			PORT_AMPLIFIER: true,
			PORT_TOWER: true,
			PORT_SOLAR: true,
			PORT_BACKHAUL: true,
		},
		scenarioActive: "default",
		currentProfile: cenums.PROFILE_NORMAL,
		globalVariationFactor: rand.Float64(),
		timeMultiplier: 1.0, 
	}
	
	mp.backhaul.UpdateMetricsProvider(mp)
	mp.battery.UpdateMetricsProvider(mp)
	mp.solar.UpdateMetricsProvider(mp)
	
	go mp.updateGlobalVariation()
	
	return mp
}


func (m *MetricsProvider) updateGlobalVariation() {
	ticker := time.NewTicker(1 * time.Second) 
	defer ticker.Stop()
	
	for range ticker.C {
		targetVariation := rand.Float64()
		m.globalVariationFactor = m.globalVariationFactor*0.7 + targetVariation*0.3
	}
}

func (m *MetricsProvider) SetProfile(profile cenums.Profile) {
	m.currentProfile = profile
}

func (m *MetricsProvider) GetMetrics(siteId string) (*ControllerMetrics, error) {
	backhaulMetrics := m.backhaul.GetMetrics()
	batteryMetrics, err := m.battery.GetMetrics()
	if err != nil {
		return nil, err
	}
	solarMetrics := m.solar.GetMetrics()

	if !m.portStatus[PORT_BACKHAUL] {
		backhaulMetrics.Latency = 0
		backhaulMetrics.Status = 0
		backhaulMetrics.Speed = 0
		backhaulMetrics.SwitchStatus = 0
		backhaulMetrics.SwitchBandwidth = 0
	}

	if !m.portStatus[PORT_SOLAR] {
		solarMetrics.PowerGeneration = 0
		solarMetrics.PanelPower = 0
		solarMetrics.PanelCurrent = 0
		solarMetrics.InverterStatus = 0
	}

	switch m.scenarioActive {
	case "power_down":
		solarMetrics.PowerGeneration = 0
		solarMetrics.PanelPower = 0
		solarMetrics.PanelCurrent = 0
		solarMetrics.PanelVoltage = 0
		solarMetrics.InverterStatus = 0
		
		batteryMetrics.Capacity = math.Max(0, batteryMetrics.Capacity - 0.5) 
		batteryMetrics.Current = -3.0 
		
		if batteryMetrics.Capacity < 5 {
			batteryMetrics.Voltage = math.Max(0, batteryMetrics.Voltage - 0.1) 
		}
		
	case "switch_off":
		backhaulMetrics.SwitchStatus = 0
		backhaulMetrics.SwitchBandwidth = 0
		
	case "backhaul_down":
		backhaulMetrics.Latency = 0
		backhaulMetrics.Status = 0
		backhaulMetrics.Speed = 0
	}

	if m.currentProfile == cenums.PROFILE_MIN && batteryMetrics.Voltage < 10.5 {
		solarMetrics.InverterStatus = 0
		if rand.Float64() > 0.7 {
			solarMetrics.PowerGeneration *= 0.5 
		}
	} else if batteryMetrics.Voltage < 11.0 && solarMetrics.PowerGeneration < 100 {
		solarMetrics.InverterStatus = 0
	}

	if rand.Float64() > 0.99 {
		scenarios := []string{"default", "power_down", "switch_off", "backhaul_down"}
		m.scenarioActive = scenarios[rand.Intn(len(scenarios))]
		
		log.Infof("Switching to scenario: %s", m.scenarioActive)
		
		go func(m *MetricsProvider) {
			time.Sleep(5 * time.Second)
			m.scenarioActive = "default"
		}(m)
	}

	return &ControllerMetrics{
		Backhaul: backhaulMetrics,
		Battery:  batteryMetrics,
		Solar:    solarMetrics,
		Time:     time.Now(),
	}, nil
}
func NewBackhaulProvider() *BackhaulProvider {
	return &BackhaulProvider{
		lastUpdate: time.Now(),
		jitterFactor: rand.Float64() * 0.5,
		noiseAmplitude: rand.Float64() * 0.3,
	}
}


func (b *BackhaulProvider) GetMetrics() *BackhaulMetrics {
	status := 1.0
	
	microTime := float64(time.Now().UnixNano()) / 1e9 
	microOscillation := math.Sin(microTime*2*math.Pi) * b.noiseAmplitude
	
	var baseLatency, baseSpeed float64
	
	profile := b.metricsProvider.currentProfile
	globalFactor := b.metricsProvider.globalVariationFactor
	
	switch profile {
	case cenums.PROFILE_MIN:
		baseLatency = 150.0 + 100.0*globalFactor  
		baseSpeed = 0.5 + 2.0*globalFactor    
		if rand.Float64() > 0.7 { 
			status = 0.0
		}
	case cenums.PROFILE_MAX:
		baseLatency = 5.0 + 15.0*globalFactor  
		baseSpeed = 80.0 + 120.0*globalFactor  
		if rand.Float64() > 0.98 { 
			status = 0.0
		}
	default: 
		baseLatency = 30.0 + 20.0*globalFactor 
		baseSpeed = 20.0 + 30.0*globalFactor     
		if rand.Float64() > 0.95 {
			status = 0.0
		}
	}

	var latency, speed, switchBandwidth float64
	var switchStatus float64 = 1.0

	if status == 1.0 {
		jitter := math.Sin(microTime*5) * 5
		latency = baseLatency + jitter + microOscillation*10
		speed = baseSpeed - (jitter*0.2) + microOscillation*2
		
		latency = math.Max(1.0, latency)
		speed = math.Max(0.1, speed)
		
		if profile == cenums.PROFILE_MIN && rand.Float64() > 0.8 { 
			switchStatus = 0.0
		} else if rand.Float64() > 0.98 { 
			switchStatus = 0.0
		}
		
		if switchStatus == 1.0 {
			switch profile {
			case cenums.PROFILE_MIN:
				switchBandwidth = 10.0 + (rand.Float64() * 40.0)
			case cenums.PROFILE_MAX:
				switchBandwidth = 500.0 + (rand.Float64() * 500.0) 
			default: 
				switchBandwidth = 100.0 + (rand.Float64() * 100.0) 
			}
		}
	} else {
		latency = 0
		speed = 0
		switchStatus = 0
		switchBandwidth = 0
	}

	return &BackhaulMetrics{
		Latency:         latency,
		Status:          status,
		Speed:           speed,
		SwitchStatus:    switchStatus,
		SwitchBandwidth: switchBandwidth,
	}
}


func NewBatteryProvider() *BatteryProvider {
	return &BatteryProvider{
		startTime:    time.Now(),
		lastCapacity: 85.0,
		cycleCount:   0,
		microFluctuationFactor: rand.Float64() * 0.2,
		temperatureVariance: rand.Float64() * 2.0,
	}
}

func (m *BatteryProvider) GetMetrics() (*BatteryMetrics, error) {
	// Use current time for second-by-second changes
	currentTime := float64(time.Now().UnixNano()) / 1e9
	
	// Add micro-oscillations for every second change
	microCycle := math.Sin(currentTime*2*math.Pi) * m.microFluctuationFactor
	
	// Base values adjusted by profile and global variation
	globalFactor := m.metricsProvider.globalVariationFactor
	var baseCapacity, baseVoltage float64
	
	switch m.metricsProvider.currentProfile {
	case cenums.PROFILE_MIN:
		baseCapacity = 10.0 + (10.0 * globalFactor) // 10-20% capacity
		baseVoltage = 10.0 + (1.5 * globalFactor)   // 10.0-11.5V
		
		// For very low profiles, we might want metrics to show critical condition
		if baseCapacity < 15 && rand.Float64() > 0.7 {
			baseCapacity = 0 // Complete battery drain
			baseVoltage = 0  // System shutdown
		}
	case cenums.PROFILE_MAX:
		baseCapacity = 80.0 + (20.0 * globalFactor) // 80-100% capacity
		baseVoltage = 12.0 + (0.8 * globalFactor)   // 12.0-12.8V
	default: // PROFILE_NORMAL
		baseCapacity = 50.0 + (30.0 * globalFactor) // 50-80% capacity
		baseVoltage = 12.0 + (0.5 * globalFactor)   // 12.0-12.5V
	}
	
	// Add second-by-second small variations
	capacity := baseCapacity + (microCycle * 2.0)
	capacity = math.Max(0.0, math.Min(100.0, capacity))
	
	// Add small voltage fluctuations per second
	voltage := baseVoltage + (microCycle * 0.1)
	
	// Determine charging state based on profile and time of day
	hourOfDay := float64(time.Now().Hour())
	daytime := hourOfDay >= 6 && hourOfDay <= 18
	
	charging := daytime && m.metricsProvider.currentProfile != cenums.PROFILE_MIN
	status := "Discharging"
	if charging {
		status = "Charging"
	}
	
	// Current varies by second and depends on charging state
	var current float64
	if charging {
		switch m.metricsProvider.currentProfile {
		case cenums.PROFILE_MIN:
			current = 0.5 + (microCycle * 0.2) // 0.3-0.7A charging
		case cenums.PROFILE_MAX:
			current = 5.0 + (microCycle * 1.0) // 4-6A charging
		default: // PROFILE_NORMAL
			current = 2.0 + (microCycle * 0.5) // 1.5-2.5A charging
		}
	} else {
		switch m.metricsProvider.currentProfile {
		case cenums.PROFILE_MIN:
			current = -(1.5 + microCycle * 0.5) // 1-2A discharge
		case cenums.PROFILE_MAX:
			current = -(0.5 + microCycle * 0.2) // 0.3-0.7A discharge
		default: // PROFILE_NORMAL
			current = -(1.0 + microCycle * 0.3) // 0.7-1.3A discharge
		}
	}
	
	// Temperature varies by profile and with seconds
	var temperature float64
	switch m.metricsProvider.currentProfile {
	case cenums.PROFILE_MIN:
		temperature = 30.0 + (microCycle * 5.0) // 25-35C (running hot)
	case cenums.PROFILE_MAX:
		temperature = 20.0 + (microCycle * 2.0) // 18-22C (optimal)
	default: // PROFILE_NORMAL
		temperature = 25.0 + (microCycle * 3.0) // 22-28C (normal)
	}
	
	// Health status
	var health string
	switch m.metricsProvider.currentProfile {
	case cenums.PROFILE_MIN:
		health = "Poor"
		if rand.Float64() > 0.7 {
			health = "Fair"
		}
	case cenums.PROFILE_MAX:
		health = "Good"
	default: // PROFILE_NORMAL
		health = "Good"
		if rand.Float64() > 0.8 {
			health = "Fair"
		}
	}

	// If battery is completely drained, zero out all metrics
	if voltage == 0 || capacity == 0 {
		current = 0
		temperature = 0
		health = "Poor"
	}

	m.lastCapacity = capacity
	
	return &BatteryMetrics{
		Capacity:    capacity,
		Voltage:     voltage,
		Current:     current,
		Temperature: temperature,
		Status:      status,
		Health:      health,
	}, nil
}

func NewSolarProvider() *SolarProvider {
	return &SolarProvider{
		startTime:      time.Now(),
		energyTotal:    0,
		weatherPattern: 1.0,
		cloudCoverFactor: rand.Float64() * 0.4,
		microWeatherChanges: rand.Float64() * 0.3,
		timeAcceleration: 5.0 + rand.Float64() * 10.0,
	}
}

func (s *SolarProvider) GetMetrics() *SolarMetrics {
	// Use current time for second-by-second changes
	currentTime := float64(time.Now().UnixNano()) / 1e9
	
	// Determine current hour
	hourOfDay := float64(time.Now().Hour())
	daylight := 0.0
	
	// Only generate solar power during daylight hours (6am-6pm)
	if hourOfDay >= 6 && hourOfDay <= 18 {
		// Peak at noon, lower at edges of daylight
		daylight = math.Sin((hourOfDay-6)*math.Pi/12)
		
		// Add small second-by-second variations
		daylight += math.Sin(currentTime*2*math.Pi) * 0.05
	}
	
	// Add cloud cover simulation with second-by-second changes
	cloudCover := math.Sin(currentTime*0.5) * s.cloudCoverFactor
	
	// Profile-based solar generation with global variation factor
	globalFactor := s.metricsProvider.globalVariationFactor
	var maxPower float64
	
	switch s.metricsProvider.currentProfile { 
	case cenums.PROFILE_MIN:
		maxPower = 100.0 + 50.0*globalFactor  // 100-150W maximum
		// More frequent cloud cover and poor conditions
		daylight *= math.Max(0.1, 0.3 - cloudCover)
		
		// For PROFILE_MIN, sometimes all solar is offline
		if rand.Float64() > 0.7 {
			daylight = 0
			maxPower = 0
		}
	case cenums.PROFILE_MAX:
		maxPower = 1500.0 + 500.0*globalFactor // 1500-2000W maximum
		// Less cloud impact
		daylight *= math.Max(0.7, 0.8 - cloudCover*0.5)
	default: // PROFILE_NORMAL
		maxPower = 800.0 + 200.0*globalFactor // 800-1000W maximum
		// Normal cloud impact
		daylight *= math.Max(0.5, 0.6 - cloudCover*0.7)
	}
	
	// Calculate power generation with second-by-second variations
	powerGeneration := maxPower * daylight
	powerGeneration += math.Sin(currentTime*5*math.Pi) * 10 // Small oscillations
	
	// Keep it positive
	powerGeneration = math.Max(0, powerGeneration)
	
	// Panel voltage depends on daylight and varies by second
	var panelVoltage float64
	if daylight > 0 {
		switch s.metricsProvider.currentProfile {
		case cenums.PROFILE_MIN:
			panelVoltage = 18.0 + (daylight * 5.0) + math.Sin(currentTime*2*math.Pi)*0.5
		case cenums.PROFILE_MAX:
			panelVoltage = 24.0 + (daylight * 12.0) + math.Sin(currentTime*2*math.Pi)*1.0
		default: // PROFILE_NORMAL
			panelVoltage = 20.0 + (daylight * 8.0) + math.Sin(currentTime*2*math.Pi)*0.8
		}
	}
	
	// Calculate current with a more dynamic relationship
	var panelCurrent float64
	if panelVoltage > 0 {
		panelCurrent = powerGeneration / panelVoltage
		// Add small second-by-second current variations
		panelCurrent += math.Sin(currentTime*3*math.Pi) * 0.1
		panelCurrent = math.Max(0, panelCurrent)
	}
	
	// Update energy total
	intervalHours := 1.0 / 3600.0
	s.energyTotal += powerGeneration * intervalHours / 1000.0
	
	// Inverter status
	inverterStatus := 0.0
	if powerGeneration > 50.0 {
		inverterStatus = 1.0
	}
	
	return &SolarMetrics{
		PowerGeneration: powerGeneration,
		EnergyTotal:    s.energyTotal,
		PanelPower:     powerGeneration,
		PanelCurrent:   panelCurrent,
		PanelVoltage:   panelVoltage,
		InverterStatus: inverterStatus,
	}
}
func (e *PrometheusExporter) StartMetricsCollection(ctx context.Context, interval time.Duration) error {
	// Use 1 second interval for real-time updates
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Initialize metrics with a first collection immediately
	if err := e.collectMetrics(); err != nil {
		log.Warnf("Initial metrics collection failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Infof("Stopping metrics collection due to context cancellation")
			return ctx.Err()
		case <-e.shutdown:
			log.Infof("Stopping metrics collection due to shutdown signal")
			return nil
		case <-ticker.C:
			if err := e.collectMetrics(); err != nil {
				log.Errorf("Error collecting metrics: %v", err)
			}
		}
	}
}


func (e *PrometheusExporter) Shutdown() {
	close(e.shutdown)
}

func (e *PrometheusExporter) collectMetrics() error {
	metrics, err := e.metricsProvider.GetMetrics(e.siteId)
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}
	log.Debugf("Collecting metrics for site %s: Solar power: %f, Battery capacity: %f", 
	e.siteId, metrics.Solar.PowerGeneration, metrics.Battery.Capacity)

	// Update solar metrics
	
	e.solarPowerGeneration.WithLabelValues("watts", e.siteId).Set(metrics.Solar.PowerGeneration)
	e.solarEnergyTotal.WithLabelValues("kwh", e.siteId).Set(metrics.Solar.EnergyTotal)
	e.solarPanelPower.WithLabelValues("watts", e.siteId).Set(metrics.Solar.PanelPower)
	e.solarPanelCurrent.WithLabelValues("amps", e.siteId).Set(metrics.Solar.PanelCurrent)
	e.solarPanelVoltage.WithLabelValues("volts", e.siteId).Set(metrics.Solar.PanelVoltage)
	e.solarInverterStatus.WithLabelValues("status", e.siteId).Set(metrics.Solar.InverterStatus)

	// Update battery metrics
	e.batteryChargeStatus.WithLabelValues("capacity", e.siteId).Set(metrics.Battery.Capacity)
	e.batteryVoltage.WithLabelValues("volts", e.siteId).Set(metrics.Battery.Voltage)
	e.batteryHealth.WithLabelValues("status", e.siteId).Set(map[string]float64{
		"Good": 1.0,
		"Fair": 0.5,
		"Poor": 0.0,
	}[metrics.Battery.Health])
	e.batteryCurrent.WithLabelValues("amps", e.siteId).Set(metrics.Battery.Current)
	e.batteryTemperature.WithLabelValues("celsius", e.siteId).Set(metrics.Battery.Temperature)

	// Update network metrics
	e.backhaulLatency.WithLabelValues("ms", e.siteId).Set(metrics.Backhaul.Latency)
	e.backhaulStatus.WithLabelValues("status", e.siteId).Set(metrics.Backhaul.Status)
	e.backhaulSpeed.WithLabelValues("mbps", e.siteId).Set(metrics.Backhaul.Speed)
	e.switchPortStatus.WithLabelValues("status", e.siteId).Set(metrics.Backhaul.SwitchStatus)
	e.switchPortBandwidth.WithLabelValues("mbps", e.siteId).Set(metrics.Backhaul.SwitchBandwidth)

	return nil
}

func (m *MetricsProvider) SetPortStatus(port int, status bool) error {
	if port < 1 || port > 4 {
		return fmt.Errorf("invalid port number: %d", port)
	}
	m.portStatus[port] = status
	return nil
}

func (m *MetricsProvider) SetScenario(scenario string) {
	m.scenarioActive = scenario
}
func (m *MetricsProvider) GetPowerStatus() (bool, error) {
	return m.scenarioActive != string(cenums.SCENARIO_POWER_DOWN), nil
}