package metrics

import (
	"math"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)
type PrometheusExporter struct {
	// Solar metrics
	solarPowerGeneration prometheus.Gauge
	solarEnergyTotal    prometheus.Gauge
	solarPanelPower     prometheus.Gauge
	solarPanelCurrent   prometheus.Gauge
	solarPanelVoltage   prometheus.Gauge
	solarInverterStatus prometheus.Gauge

	// Battery metrics
	batteryChargeStatus prometheus.Gauge
	batteryVoltage      prometheus.Gauge
	batteryHealth       prometheus.Gauge
	batteryCurrent      prometheus.Gauge
	batteryTemperature  prometheus.Gauge

	// Network metrics
	backhaulLatency      prometheus.Gauge
	backhaulStatus       prometheus.Gauge
	backhaulSpeed        prometheus.Gauge
	switchPortStatus     prometheus.Gauge
	switchPortBandwidth  prometheus.Gauge

	metricsProvider *MetricsProvider
}





func NewPrometheusExporter(metricsProvider *MetricsProvider) *PrometheusExporter {
	exporter := &PrometheusExporter{
		metricsProvider: metricsProvider,

		// Solar metrics
		solarPowerGeneration: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "solar_power_generation",
			Help: "Current solar power generation in watts",
		}),
		solarEnergyTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "solar_energy_total",
			Help: "Total solar energy generated in kilowatt-hours",
		}),
		solarPanelPower: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "solar_panel_power",
			Help: "Current solar panel power in watts",
		}),
		solarPanelCurrent: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "solar_panel_current",
			Help: "Current solar panel current in amperes",
		}),
		solarPanelVoltage: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "solar_panel_voltage",
			Help: "Current solar panel voltage in volts",
		}),
		solarInverterStatus: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "solar_inverter_status",
			Help: "Solar inverter status (1 = working, 0 = not working)",
		}),

		// Battery metrics
		batteryChargeStatus: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "battery_charge_status",
			Help: "Battery charge status in percentage",
		}),
		batteryVoltage: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "battery_voltage_volts",
			Help: "Battery voltage in volts",
		}),
		batteryHealth: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "battery_health",
			Help: "Battery health status (1 = good, 0 = poor)",
		}),
		batteryCurrent: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "battery_current",
			Help: "Battery current in amperes",
		}),
		batteryTemperature: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "battery_temperature",
			Help: "Battery temperature in Celsius",
		}),

		// Network metrics
		backhaulLatency: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "backhaul_latency_ms",
			Help: "Backhaul latency in milliseconds",
		}),
		backhaulStatus: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "backhaul_status",
			Help: "Backhaul status (1 = up, 0 = down)",
		}),
		backhaulSpeed: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "backhaul_speed_mbps",
			Help: "Backhaul speed in Mbps",
		}),
		switchPortStatus: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "switch_port_status",
			Help: "Switch port status (1 = up, 0 = down)",
		}),
		switchPortBandwidth: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "switch_port_bandwidth",
			Help: "Switch port bandwidth in Mbps",
		}),
	}

	return exporter
}

// BackhaulMetrics represents network backhaul measurements
type BackhaulMetrics struct {
	Latency float64
	Status  float64
	Speed   float64
	SwitchStatus  float64  
	SwitchBandwidth float64 
}

// BatteryMetrics represents battery system measurements
type BatteryMetrics struct {
	Capacity    float64
	Voltage     float64
	Current     float64
	Temperature float64
	Status      string
	Health      string
}

// SolarMetrics represents solar system measurements
type SolarMetrics struct {
	PowerGeneration float64
	EnergyTotal    float64
	PanelPower     float64
	PanelCurrent   float64
	PanelVoltage   float64
	InverterStatus float64
}

// SystemMetrics combines all metrics into a single structure
type ControllerMetrics struct {
	Backhaul *BackhaulMetrics
	Battery  *BatteryMetrics
	Solar    *SolarMetrics
	Time     time.Time
}

// BackhaulProvider handles backhaul metrics collection
type BackhaulProvider struct {
	lastUpdate time.Time
}

// BatteryProvider handles battery metrics collection
type BatteryProvider struct {
	startTime    time.Time
	lastCapacity float64
	cycleCount   int
}

// SolarProvider handles solar metrics collection
type SolarProvider struct {
	startTime      time.Time
	energyTotal    float64
	weatherPattern float64
}

// MetricsProvider combines all providers into a single interface
type MetricsProvider struct {
	backhaul *BackhaulProvider
	battery  *BatteryProvider
	solar    *SolarProvider
}

// NewMetricsProvider creates a new instance of the combined metrics provider
func NewMetricsProvider() *MetricsProvider {
	return &MetricsProvider{
		backhaul: NewBackhaulProvider(),
		battery:  NewBatteryProvider(),
		solar:    NewSolarProvider(),
	}
}

// GetMetrics returns all system metrics
func (m *MetricsProvider) GetMetrics() (*ControllerMetrics, error) {
	backhaulMetrics := m.backhaul.GetMetrics()
	batteryMetrics, err := m.battery.GetMetrics()
	if err != nil {
		return nil, err
	}
	solarMetrics := m.solar.GetMetrics()

	return &ControllerMetrics{
		Backhaul: backhaulMetrics,
		Battery:  batteryMetrics,
		Solar:    solarMetrics,
		Time:     time.Now(),
	}, nil
}

// NewBackhaulProvider creates a new backhaul provider
func NewBackhaulProvider() *BackhaulProvider {
	return &BackhaulProvider{
		lastUpdate: time.Now(),
	}
}

// GetMetrics returns backhaul metrics
func (b *BackhaulProvider) GetMetrics() *BackhaulMetrics {
	status := 1.0
	if rand.Float64() > 0.95 {
		status = 0.0
	}

	var latency, speed float64
	
	if status == 1.0 {
		latency = 10.0 + (rand.Float64() * 90.0)
		speed = 10.0 + (rand.Float64() * 90.0)
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

// NewBatteryProvider creates a new battery provider
func NewBatteryProvider() *BatteryProvider {
	return &BatteryProvider{
		startTime:    time.Now(),
		lastCapacity: 85.0,
		cycleCount:   0,
	}
}

// GetMetrics returns battery metrics
func (m *BatteryProvider) GetMetrics() (*BatteryMetrics, error) {
	elapsed := time.Since(m.startTime).Seconds()
	
	dayNightCycle := math.Sin(elapsed/(24*3600)*2*math.Pi)
	batteryCycle := math.Sin(elapsed/(4*3600)*2*math.Pi)
	
	baseCapacity := 85.0 + (dayNightCycle * 10.0) + (batteryCycle * 5.0)
	capacity := math.Max(20.0, math.Min(100.0, baseCapacity))
	
	charging := dayNightCycle > 0 && capacity < 95.0
	status := "Discharging"
	if charging {
		status = "Charging"
	}
	
	voltage := 11.4 + (capacity/100.0 * 1.2)
	
	current := 0.1
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

// NewSolarProvider creates a new solar provider
func NewSolarProvider() *SolarProvider {
	return &SolarProvider{
		startTime:      time.Now(),
		energyTotal:    0,
		weatherPattern: 1.0,
	}
}

// GetMetrics returns solar metrics
func (s *SolarProvider) GetMetrics() *SolarMetrics {
	elapsed := time.Since(s.startTime).Seconds()
	
	hourOfDay := float64(time.Now().Hour())
	daylight := math.Max(0, math.Sin((hourOfDay-6)*math.Pi/12))
	
	s.weatherPattern = 0.7 + (math.Sin(elapsed/14400)*0.3) + (math.Sin(elapsed/3600)*0.1)
	s.weatherPattern = math.Max(0.1, math.Min(1.0, s.weatherPattern))
	
	maxPower := 1000.0
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
func (e *PrometheusExporter) StartMetricsCollection(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			e.collectMetrics()
		}
	}()
}

func (e *PrometheusExporter) collectMetrics() {
	metrics, err := e.metricsProvider.GetMetrics()
	if err != nil {
		return
	}

	// Update Solar metrics
	e.solarPowerGeneration.Set(metrics.Solar.PowerGeneration)
	e.solarEnergyTotal.Set(metrics.Solar.EnergyTotal)
	e.solarPanelPower.Set(metrics.Solar.PanelPower)
	e.solarPanelCurrent.Set(metrics.Solar.PanelCurrent)
	e.solarPanelVoltage.Set(metrics.Solar.PanelVoltage)
	e.solarInverterStatus.Set(metrics.Solar.InverterStatus)

	// Update Battery metrics
	e.batteryChargeStatus.Set(metrics.Battery.Capacity)
	e.batteryVoltage.Set(metrics.Battery.Voltage)
	e.batteryHealth.Set(map[string]float64{
		"Good": 1.0,
		"Fair": 0.5,
		"Poor": 0.0,
	}[metrics.Battery.Health])
	e.batteryCurrent.Set(metrics.Battery.Current)
	e.batteryTemperature.Set(metrics.Battery.Temperature)

	// Update Network metrics
	e.backhaulLatency.Set(metrics.Backhaul.Latency)
	e.backhaulStatus.Set(metrics.Backhaul.Status)
	e.backhaulSpeed.Set(metrics.Backhaul.Speed)
	e.switchPortStatus.Set(metrics.Backhaul.SwitchStatus)
	e.switchPortBandwidth.Set(metrics.Backhaul.SwitchBandwidth)
}