package metrics

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
    lastUpdate time.Time
    forceBackhaulDown bool
    forceSwitchOff bool
}


type BatteryProvider struct {
	startTime    time.Time
	lastCapacity float64
	cycleCount   int
}

type SolarProvider struct {
	startTime      time.Time
	energyTotal    float64
	weatherPattern float64
}

type MetricsProvider struct {
	backhaul *BackhaulProvider
	battery  *BatteryProvider
	solar    *SolarProvider
}

func NewMetricsProvider() *MetricsProvider {
	return &MetricsProvider{
		backhaul: NewBackhaulProvider(),
		battery:  NewBatteryProvider(),
		solar:    NewSolarProvider(),
	}
}

func (m *MetricsProvider) GetMetrics(siteId string) (*ControllerMetrics, error) {
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

func (m *MetricsProvider) Backhaul() *BackhaulProvider {
    return m.backhaul
}

func (m *MetricsProvider) Battery() *BatteryProvider {
    return m.battery
}

func (m *MetricsProvider) Solar() *SolarProvider {
    return m.solar
}

func NewBackhaulProvider() *BackhaulProvider {
    return &BackhaulProvider{
        lastUpdate: time.Now(),
        forceBackhaulDown: false,
        forceSwitchOff: false,
    }
}

func (b *BackhaulProvider) GetMetrics() *BackhaulMetrics {
    if b.forceBackhaulDown {
        return &BackhaulMetrics{
            Latency:         0.0,
            Status:          0.0,  
            Speed:           0.0,
            SwitchStatus:    1.0,  
            SwitchBandwidth: 0.0, 
        }
    }
    
    if b.forceSwitchOff {
        return &BackhaulMetrics{
            Latency:         0.0,
            Status:          0.0, 
            Speed:           0.0,
            SwitchStatus:    0.0,  
            SwitchBandwidth: 0.0,
        }
    }
    
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

func (b *BackhaulProvider) SetForceSwitchOff(value bool) {
    b.forceSwitchOff = value
}

func (b *BackhaulProvider) SetForceBackhaulDown(value bool) {
    b.forceBackhaulDown = value
}

func (b *BackhaulProvider) IsSwitchOff() bool {
    return b.forceSwitchOff
}

func (b *BackhaulProvider) IsBackhaulDown() bool {
    return b.forceBackhaulDown
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
	
	baseCapacity := 85.0 + (dayNightCycle * 10.0) + (batteryCycle * 5.0)
	capacity := math.Max(20.0, math.Min(100.0, baseCapacity))
	
	charging := dayNightCycle > 0 && capacity < 95.0
	status := "Discharging"
	if charging {
		status = "Charging"
	}
	
	voltage := 11.4 + (capacity/100.0 * 1.2)
	
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

func (b *BatteryProvider) SetLastCapacity(value float64) {
    b.lastCapacity = value
}

func (b *BatteryProvider) GetLastCapacity() float64 {
    return b.lastCapacity
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

func (s *SolarProvider) SetWeatherPattern(value float64) {
    s.weatherPattern = value
}

func (s *SolarProvider) GetWeatherPattern() float64 {
    return s.weatherPattern
}

func (e *PrometheusExporter) StartMetricsCollection(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping metrics collection due to context cancellation")
			return ctx.Err()
		case <-e.shutdown:
			log.Printf("Stopping metrics collection due to shutdown signal")
			return nil
		case <-ticker.C:
			if err := e.collectMetrics(); err != nil {
				log.Printf("Error collecting metrics: %v", err)
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