package battery

import (
	"math"
	"time"
)

type MockBatteryProvider struct {
    startTime    time.Time
    lastCapacity float64
    cycleCount   int
}

func NewMockBatteryProvider() *MockBatteryProvider {
    return &MockBatteryProvider{
        startTime:    time.Now(),
        lastCapacity: 85.0, // Start at 85% capacity
        cycleCount:   0,
    }
}

func (m *MockBatteryProvider) GetMetrics() (*BatteryMetrics, error) {
    elapsed := time.Since(m.startTime).Seconds()
    
    // Simulate day/night cycle (24 hour period)
    dayNightCycle := math.Sin(elapsed/(24*3600)*2*math.Pi)
    
    // Battery discharge/charge cycle (4 hour period)
    batteryCycle := math.Sin(elapsed/(4*3600)*2*math.Pi)
    
    // Calculate capacity with realistic patterns
    baseCapacity := 85.0 + (dayNightCycle * 10.0) + (batteryCycle * 5.0)
    capacity := math.Max(20.0, math.Min(100.0, baseCapacity))
    
    // Determine charging state
    charging := dayNightCycle > 0 && capacity < 95.0
    status := "Discharging"
    if charging {
        status = "Charging"
    }
    
    // Calculate voltage based on capacity
    voltage := 11.4 + (capacity/100.0 * 1.2)
    
    // Calculate current based on charging state
    current := 0.1
    if charging {
        current = 2.0 + (dayNightCycle * 0.5)
    } else {
        current = -(1.0 + math.Abs(batteryCycle*0.5))
    }
    
    // Temperature varies with current and ambient
    ambientTemp := 22.0 + (dayNightCycle * 3.0)
    temperature := ambientTemp + (math.Abs(current) * 0.5)
    
    // Health degrades very slowly over time
    timeBasedWear := elapsed / (365 * 24 * 3600) * 10 // 10% degradation per year
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