package metrics

import (
	"math"
	"math/rand"
	"time"
)

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