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

func NewPrometheusExporter(metricsProvider *MetricsProvider,siteId string) *PrometheusExporter {
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
}

type SolarProvider struct {
	startTime      time.Time
	energyTotal    float64
	weatherPattern float64
	metricsProvider *MetricsProvider 
}

type MetricsProvider struct {
	backhaul       *BackhaulProvider
	battery        *BatteryProvider
	solar          *SolarProvider
	portStatus     map[int]bool  
	scenarioActive string        
	currentProfile cenums.Profile
}

func (b *BackhaulProvider) UpdateMetricsProvider(provider *MetricsProvider) {
	b.metricsProvider = provider
}

func (b *BatteryProvider) UpdateMetricsProvider(provider *MetricsProvider) {
	b.metricsProvider = provider
}

// Added method to update metrics provider reference for SolarProvider
func (s *SolarProvider) UpdateMetricsProvider(provider *MetricsProvider) {
	s.metricsProvider = provider
}

func NewMetricsProvider() *MetricsProvider {
	mp := &MetricsProvider{
		backhaul:       NewBackhaulProvider(),
		battery:        NewBatteryProvider(),
		solar:          NewSolarProvider(),
		portStatus:     map[int]bool{
			PORT_AMPLIFIER: true,
			PORT_TOWER:     true,
			PORT_SOLAR:     true,
			PORT_BACKHAUL:  true,
		},
		scenarioActive: "default",
		currentProfile: cenums.PROFILE_NORMAL,
	}
	
	mp.backhaul.UpdateMetricsProvider(mp)
	mp.battery.UpdateMetricsProvider(mp)  
	mp.solar.UpdateMetricsProvider(mp)   
	
	return mp
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
		// If backhaul port is down, zero out all backhaul metrics
		backhaulMetrics.Latency = 0
		backhaulMetrics.Status = 0
		backhaulMetrics.Speed = 0
	}

	if !m.portStatus[PORT_SOLAR] {
		// If solar port is down, zero out all solar metrics
		solarMetrics.PowerGeneration = 0
		solarMetrics.PanelPower = 0
		solarMetrics.PanelCurrent = 0
		solarMetrics.InverterStatus = 0
		// Keep energy total as it's cumulative
	}

	// Apply scenario effects
	switch m.scenarioActive {
	case "power_down":
		// Simulate power down - battery draining, solar off
		solarMetrics.PowerGeneration = 0
		solarMetrics.PanelPower = 0
		solarMetrics.PanelCurrent = 0
		solarMetrics.InverterStatus = 0
		
		// Battery discharging rapidly
		if batteryMetrics.Capacity > 5 {
			batteryMetrics.Capacity = math.Max(5, batteryMetrics.Capacity - 10)
		}
		batteryMetrics.Current = -2.0 // Heavy discharge
		
	case "switch_off":
		// Simulate switch being off - network port issues
		backhaulMetrics.SwitchStatus = 0
		backhaulMetrics.SwitchBandwidth = 0
		
	case "backhaul_down":
		// Simulate backhaul down - all backhaul metrics to zero
		backhaulMetrics.Latency = 0
		backhaulMetrics.Status = 0
		backhaulMetrics.Speed = 0
	}

	// Apply battery voltage effects
	// If battery voltage is below 12V, solar metrics should be affected
	if batteryMetrics.Voltage < 12.0 && solarMetrics.PowerGeneration < 50 {
		solarMetrics.InverterStatus = 0
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
	}
}

func (b *BackhaulProvider) GetMetrics() *BackhaulMetrics {
	status := 1.0
	
	// Profile-based backhaul metrics
	var baseLatency, baseSpeed float64
	
	// Get profile from metrics provider
	profile := b.metricsProvider.currentProfile
	
	switch profile {
	case cenums.PROFILE_MIN:
		baseLatency = 100.0  // Higher latency
		baseSpeed = 5.0      // Lower speed
		if rand.Float64() > 0.8 { // More frequent downtime
			status = 0.0
		}
	case cenums.PROFILE_MAX:
		baseLatency = 20.0   // Lower latency
		baseSpeed = 100.0    // Higher speed
		if rand.Float64() > 0.98 { // Less frequent downtime
			status = 0.0
		}
	default: // PROFILE_NORMAL
		baseLatency = 50.0   // Normal latency
		baseSpeed = 50.0     // Normal speed
		if rand.Float64() > 0.95 {
			status = 0.0
		}
	}

	var latency, speed float64
	if status == 1.0 {
		latency = baseLatency + (rand.Float64() * 20.0)
		speed = baseSpeed + (rand.Float64() * 20.0)
	}

	switchStatus := 1.0
	if rand.Float64() > 0.98 { 
		switchStatus = 0.0
	}

	var switchBandwidth float64
	if switchStatus == 1.0 {
		switchBandwidth = 100.0 + (rand.Float64() * 900.0)
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
	}
}

func (m *BatteryProvider) GetMetrics() (*BatteryMetrics, error) {
	elapsed := time.Since(m.startTime).Seconds()
	
	dayNightCycle := math.Sin(elapsed/(24*3600)*2*math.Pi)
	batteryCycle := math.Sin(elapsed/(4*3600)*2*math.Pi)
	
	// Base values adjusted by profile
	var baseCapacity, baseVoltage float64
	switch m.metricsProvider.currentProfile { // Changed from b.provider to m.metricsProvider
	case cenums.PROFILE_MIN:
		baseCapacity = 45.0 // Lower base capacity
		baseVoltage = 11.2  // Lower voltage
	case cenums.PROFILE_MAX:
		baseCapacity = 95.0 // Higher base capacity
		baseVoltage = 12.8  // Higher voltage
	default: // PROFILE_NORMAL
		baseCapacity = 85.0 // Normal base capacity
		baseVoltage = 12.3  // Normal voltage
	}
	
	capacity := baseCapacity + (dayNightCycle * 10.0) + (batteryCycle * 5.0)
	capacity = math.Max(20.0, math.Min(100.0, capacity))
	
	voltage := baseVoltage + (capacity/100.0 * 0.5)
	
	charging := dayNightCycle > 0 && capacity < 95.0
	status := "Discharging"
	if charging {
		status = "Charging"
	}
	
	var current float64
	if charging {
		current = 2.0 + (dayNightCycle * 0.5)
	} else {
		current = -(1.0 + math.Abs(batteryCycle*0.5))
	}
	
	ambientTemp := 22.0 + (dayNightCycle * 3.0)
	temperature := ambientTemp + (math.Abs(current) * 0.5)
	
	timeBasedWear := elapsed / (365 * 24 * 3600) * 10
	health := "Good"
	if timeBasedWear > 20 {
		health = "Fair"
	} else if timeBasedWear > 50 {
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
	}
}

func (s *SolarProvider) GetMetrics() *SolarMetrics {
	elapsed := time.Since(s.startTime).Seconds()
	
	hourOfDay := float64(time.Now().Hour())
	daylight := math.Max(0, math.Sin((hourOfDay-6)*math.Pi/12))
	
	// Profile-based solar generation
	var maxPower float64
	switch s.metricsProvider.currentProfile { 
	case cenums.PROFILE_MIN:
		maxPower = 500.0  // Lower maximum power
		s.weatherPattern = 0.4 + (math.Sin(elapsed/14400)*0.2)  // Worse weather conditions
	case cenums.PROFILE_MAX:
		maxPower = 2000.0 // Higher maximum power
		s.weatherPattern = 0.9 + (math.Sin(elapsed/14400)*0.1)  // Better weather conditions
	default: // PROFILE_NORMAL
		maxPower = 1000.0 // Normal maximum power
		s.weatherPattern = 0.7 + (math.Sin(elapsed/14400)*0.3)  // Normal weather conditions
	}
	
	s.weatherPattern = math.Max(0.1, math.Min(1.0, s.weatherPattern))
	
	powerGeneration := maxPower * daylight * s.weatherPattern
	
	panelVoltage := 24.0 + (daylight * 4.0)
	panelCurrent := powerGeneration / panelVoltage
	
	intervalHours := 1.0 / 3600.0
	s.energyTotal += powerGeneration * intervalHours / 1000.0
	
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
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

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
	log.Infof("Collecting metrics for site %s: Solar power: %f, Battery capacity: %f", 
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