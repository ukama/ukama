package solar

import (
	"math"
	"time"
)

type SolarMetrics struct {
    PowerGeneration float64 
    EnergyTotal    float64 
    PanelPower     float64 
    PanelCurrent   float64
    PanelVoltage   float64    // New field
    InverterStatus float64 
}

type SolarProvider struct {
    startTime      time.Time
    energyTotal    float64
    weatherPattern float64
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
        PanelVoltage:   panelVoltage,    // Add panel voltage
        InverterStatus: inverterStatus,
    }
}
