package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	battery "metrics-generator/internal/battery"
	"metrics-generator/internal/solar"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	solarPowerGeneration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solar_power_generation",
			Help: "Solar power generation in watts",
		},
		[]string{"unit"},
	)
	solarEnergyTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solar_energy_total",
			Help: "Solar energy total in kWh",
		},
		[]string{"unit"},
	)
	solarPanelPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solar_panel_power",
			Help: "Solar panel power in watts",
		},
		[]string{"unit"},
	)
	solarPanelCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solar_panel_current",
			Help: "Solar panel current in amperes",
		},
		[]string{"unit"},
	)
	batteryChargeStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "battery_charge_status",
			Help: "Battery charge status in percentage",
		},
		[]string{"unit"},
	)
	batteryVoltage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "battery_voltage_volts",
			Help: "Battery voltage in volts",
		},
		[]string{"unit"},
	)
	batteryHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "battery_health_percent",
			Help: "Battery health in percentage",
		},
		[]string{"unit"},
	)
	batteryCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "battery_current_amperes",
			Help: "Battery current in amperes",
		},
		[]string{"unit"},
	)
	batteryTemperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "battery_temperature",
			Help: "Battery temperature in degrees Celsius",
		},
		[]string{"unit"},
	)
	solarInverterStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solar_inverter_status",
			Help: "Solar inverter status (0 = off, 1 = on)",
		},
		[]string{"unit"},
	)
	switchPortStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "switch_port_status",
			Help: "Switch port status (0 = down, 1 = up)",
		},
		[]string{"unit"},
	)
	switchPortBandwidth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "switch_port_bandwidth_usage",
			Help: "Switch port bandwidth usage in Mbps",
		},
		[]string{"unit"},
	)
)

func init() {
	prometheus.MustRegister(
		solarPowerGeneration,
		solarEnergyTotal,
		solarPanelPower,
		solarPanelCurrent,
		batteryChargeStatus,
		batteryVoltage,
		batteryHealth,
		batteryCurrent,
		batteryTemperature,
		solarInverterStatus,
		switchPortStatus,
		switchPortBandwidth,
	)
}

func main() {
	logger := log.New(os.Stdout, "Dummy Contoller : ", log.Ldate|log.Ltime|log.Lshortfile)

	_, err := battery.GetBatteryMetrics()
	if err != nil {
		logger.Printf("Failed to access real battery metrics: %v", err)
		logger.Printf("Falling back to mock battery metrics")
		os.Setenv("MOCK_BATTERY", "true")
	}

	solarProvider := solar.NewSolarProvider()

	go func() {
		for {
			 batteryMetrics, err := battery.GetBatteryMetrics()
			 if err != nil {
				 logger.Printf("Error getting battery metrics: %v", err)
			 } else {
				 logger.Printf("Battery Status: %s, Capacity: %.1f%%, Voltage: %.2fV, Current: %.2fA, Temp: %.1f°C, Health: %s",
					 batteryMetrics.Status,
					 batteryMetrics.Capacity,
					 batteryMetrics.Voltage,
					 batteryMetrics.Current,
					 batteryMetrics.Temperature,
					 batteryMetrics.Health)
				 
				 batteryChargeStatus.WithLabelValues("percentage").Set(batteryMetrics.Capacity)
				 batteryVoltage.WithLabelValues("volts").Set(batteryMetrics.Voltage)
				 batteryCurrent.WithLabelValues("amperes").Set(batteryMetrics.Current)
				 batteryTemperature.WithLabelValues("degrees celsius").Set(batteryMetrics.Temperature)
				 healthValue := 100.0
				 if batteryMetrics.Health != "Good" {
					 healthValue = 50.0
				 }
				 batteryHealth.WithLabelValues("percentage").Set(healthValue)
			 }

			 solarMetrics := solarProvider.GetMetrics()
			 
			 logger.Printf("Solar Status: Generation: %.1fW, Total: %.2fkWh, Panel Current: %.2fA, Inverter: %v",
				 solarMetrics.PowerGeneration,
				 solarMetrics.EnergyTotal,
				 solarMetrics.PanelCurrent,
				 solarMetrics.InverterStatus == 1.0)
			 
			 solarPowerGeneration.WithLabelValues("watts").Set(solarMetrics.PowerGeneration)
			 solarEnergyTotal.WithLabelValues("kwh").Set(solarMetrics.EnergyTotal)
			 solarPanelPower.WithLabelValues("watts").Set(solarMetrics.PanelPower)
			 solarPanelCurrent.WithLabelValues("amperes").Set(solarMetrics.PanelCurrent)
			 solarInverterStatus.WithLabelValues("status").Set(solarMetrics.InverterStatus)

			switchPortStatusValue := float64(rand.Intn(2)) 
			switchPortStatus.WithLabelValues("status").Set(switchPortStatusValue)

			switchPortBandwidthValue := rand.Float64() * 1000 
			switchPortBandwidth.WithLabelValues("mbps").Set(switchPortBandwidthValue)

			time.Sleep(time.Second)
		}
	}()

	port := 2112
	address := fmt.Sprintf(":%d", port)

	http.Handle("/metrics", promhttp.Handler())

	logger.Printf("Starting Dummy controller Prometheus exporter on port %d", port)

	err = http.ListenAndServe(address, nil)
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}