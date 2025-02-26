func NewPrometheusExporter(metricsProvider *MetricsProvider, siteId string) *PrometheusExporter {
    exporter := &PrometheusExporter{
        metricsProvider: metricsProvider,
        siteId:         siteId,
        shutdown:       make(chan struct{}),

        // Update metric names to be more standardized
        batteryTemperature: promauto.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "dcontroller",  // Add namespace
            Name:      "battery_temperature_celsius",
            Help:     "Battery temperature in Celsius",
        }, []string{"unit", "site"}),

        // Update solar metrics registration
        solarPowerGeneration: promauto.NewGaugeVec(prometheus.GaugeOpts{
            Name: "solar_power_generation",
            Help: "Current solar power generation in watts",
        }, []string{"site"}),  // Remove unit label, just use site

        // Update battery metrics registration
        batteryCapacity: promauto.NewGaugeVec(prometheus.GaugeOpts{
            Name: "battery_capacity",
            Help: "Battery charge status in percentage",
        }, []string{"site"}),  // Remove unit label, just use site

        // ...existing code...
    }

    // Add debug logging for metric registration
    log.Printf("Registering metrics for site %s", siteId)
    log.Printf("Metrics registered for site %s", siteId)
    
    return exporter
}

func (e *PrometheusExporter) StartMetricsCollection(ctx context.Context, interval time.Duration) error {
    log.Printf("Starting metrics collection for site: %s", e.siteId)
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    // Collect metrics immediately on start
    if err := e.collectMetrics(); err != nil {
        log.Printf("Initial metrics collection failed: %v", err)
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-e.shutdown:
            return nil
        case <-ticker.C:
            if err := e.collectMetrics(); err != nil {
                log.Printf("Failed to collect metrics: %v", err)
                continue
            }
            log.Printf("Metrics collected successfully for site: %s", e.siteId)
        }
    }
}

func (e *PrometheusExporter) collectMetrics() error {
    metrics, err := e.metricsProvider.GetMetrics(e.siteId)
    if err != nil {
        return fmt.Errorf("failed to get metrics: %w", err)
    }

    // Debug print all metric values
    log.Printf("Collecting metrics for site %s:", e.siteId)
    log.Printf("Battery: temp=%.2f°C, charge=%.2f%%, health=%s", 
        metrics.Battery.Temperature,
        metrics.Battery.Capacity,
        metrics.Battery.Health)
    log.Printf("Solar: power=%.2fW, total=%.2fkWh", 
        metrics.Solar.PowerGeneration,
        metrics.Solar.EnergyTotal)
    log.Printf("Network: backhaul_status=%.0f, latency=%.2fms", 
        metrics.Backhaul.Status,
        metrics.Backhaul.Latency)

    // Set metrics with explicit logging for each
    e.batteryTemperature.WithLabelValues("celsius", e.siteId).Set(metrics.Battery.Temperature)
    log.Printf("Set dcontroller_battery_temperature_celsius{unit=\"celsius\",site=\"%s\"} = %v", 
        e.siteId, metrics.Battery.Temperature)

    // Update solar metrics with simpler label set
    e.solarPowerGeneration.WithLabelValues(e.siteId).Set(metrics.Solar.PowerGeneration)
    log.Printf("Set solar_power_generation{site=\"%s\"} = %v", 
        e.siteId, metrics.Solar.PowerGeneration)

    // Update battery metrics with simpler label set
    e.batteryCapacity.WithLabelValues(e.siteId).Set(metrics.Battery.Capacity)
    log.Printf("Set battery_capacity{site=\"%s\"} = %v", 
        e.siteId, metrics.Battery.Capacity)

    // ...existing code...
    return nil
}
