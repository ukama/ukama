package backhaul

import (
	"math/rand"
	"time"
)

type BackhaulMetrics struct {
    Latency float64
    Status  float64
    Speed   float64    // Add speed field
}

type BackhaulProvider struct {
    lastUpdate time.Time
}

func NewBackhaulProvider() *BackhaulProvider {
    return &BackhaulProvider{
        lastUpdate: time.Now(),
    }
}

func (b *BackhaulProvider) GetMetrics() *BackhaulMetrics {
    // Simulate latency between 10ms and 100ms
    latency := 10.0 + (rand.Float64() * 90.0)
    
    // Status 1 = up, 0 = down with 95% uptime
    status := 1.0
    if rand.Float64() > 0.95 {
        status = 0.0
    }

    // Simulate speed between 10 Mbps and 100 Mbps
    speed := 10.0 + (rand.Float64() * 90.0)

    return &BackhaulMetrics{
        Latency: latency,
        Status:  status,
        Speed:   speed,
    }
}
