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
 
	 //solar controller metrics
	 chargeControllerStatus    *prometheus.GaugeVec
	 chargeControllerMode      *prometheus.GaugeVec
	 chargeControllerCurrent   *prometheus.GaugeVec
	 chargeControllerVoltage   *prometheus.GaugeVec
	 chargeControllerTemp      *prometheus.GaugeVec
	 chargeControllerEfficiency *prometheus.GaugeVec
	 
 
	 metricsProvider *MetricsProvider
	 siteId         string
	 shutdown       chan struct{} 
 }
 
 func NewPrometheusExporter(metricsProvider *MetricsProvider, siteId string) *PrometheusExporter {
	 exporter := &PrometheusExporter{
		 metricsProvider: metricsProvider,
		 siteId:          siteId,
		 shutdown:       make(chan struct{}), 
 
		 chargeControllerStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_status",
			 Help: "Charge controller status (1 = working, 0 = not working)",
		 }, []string{"unit","site"}),
		 
		 chargeControllerMode: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_mode",
			 Help: "Charge controller mode (0 = Bulk, 1 = Absorption, 2 = Float, 3 = Equalization)",
		 }, []string{"unit","site"}),
		 
		 chargeControllerCurrent: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_current",
			 Help: "Charge controller current in amperes",
		 }, []string{"unit","site"}),
		 
		 chargeControllerVoltage: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_voltage",
			 Help: "Charge controller voltage in volts",
		 }, []string{"unit","site"}),
		 
		 chargeControllerTemp: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_temperature",
			 Help: "Charge controller temperature in Celsius",
		 }, []string{"unit","site"}),
		 
		 chargeControllerEfficiency: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_efficiency",
			 Help: "Charge controller efficiency in percentage",
		 }, []string{"unit","site"}),
	 
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
	 ControllerStatus    float64
	 ControllerMode      string
	 ControllerCurrent   float64
	 ControllerVoltage   float64
	 ControllerTemp      float64
	 ControllerEfficiency float64
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
		 solarMetrics.ControllerStatus = 0
	 }
 
	 if m.currentProfile == cenums.PROFILE_MIN && batteryMetrics.Voltage < 10.5 {
		 solarMetrics.InverterStatus = 0
		 solarMetrics.ControllerStatus = 0
		 if rand.Float64() > 0.7 {
			 solarMetrics.PowerGeneration *= 0.5 
		 }
	 } else if batteryMetrics.Voltage < 11.0 && solarMetrics.PowerGeneration < 100 {
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
	 currentTime := float64(time.Now().UnixNano()) / 1e9
	 
	 microCycle := math.Sin(currentTime*2*math.Pi) * m.microFluctuationFactor
	 
	 globalFactor := m.metricsProvider.globalVariationFactor
	 var baseCapacity, baseVoltage float64
	 
	 switch m.metricsProvider.currentProfile {
	 case cenums.PROFILE_MIN:
		 baseCapacity = 10.0 + (10.0 * globalFactor) 
		 baseVoltage = 10.0 + (1.5 * globalFactor)   
		 
		 if baseCapacity < 15 && rand.Float64() > 0.7 {
			 baseCapacity = 0 
			 baseVoltage = 0  
		 }
	 case cenums.PROFILE_MAX:
		 baseCapacity = 80.0 + (20.0 * globalFactor) 
		 baseVoltage = 12.0 + (0.8 * globalFactor)   
	 default: 
		 baseCapacity = 50.0 + (30.0 * globalFactor) 
		 baseVoltage = 12.0 + (0.5 * globalFactor)   
	 }
	 
	 capacity := baseCapacity + (microCycle * 2.0)
	 capacity = math.Max(0.0, math.Min(100.0, capacity))
	 
	 voltage := baseVoltage + (microCycle * 0.1)
	 
	 hourOfDay := float64(time.Now().Hour())
	 daytime := hourOfDay >= 6 && hourOfDay <= 18
	 
	 charging := daytime && m.metricsProvider.currentProfile != cenums.PROFILE_MIN
	 status := "Discharging"
	 if charging {
		 status = "Charging"
	 }
	 
	 var current float64
	 if charging {
		 switch m.metricsProvider.currentProfile {
		 case cenums.PROFILE_MIN:
			 current = 0.5 + (microCycle * 0.2) 
		 case cenums.PROFILE_MAX:
			 current = 5.0 + (microCycle * 1.0) 
		 default: 
			 current = 2.0 + (microCycle * 0.5) 
		 }
	 } else {
		 switch m.metricsProvider.currentProfile {
		 case cenums.PROFILE_MIN:
			 current = -(1.5 + microCycle * 0.5) 
		 case cenums.PROFILE_MAX:
			 current = -(0.5 + microCycle * 0.2)
		 default: 
			 current = -(1.0 + microCycle * 0.3) 
		 }
	 }
	 
	 var temperature float64
	 switch m.metricsProvider.currentProfile {
	 case cenums.PROFILE_MIN:
		 temperature = 30.0 + (microCycle * 5.0) 
	 case cenums.PROFILE_MAX:
		 temperature = 20.0 + (microCycle * 2.0) 
	 default: 
		 temperature = 25.0 + (microCycle * 3.0) 
	 }
	 
	 var health string
	 switch m.metricsProvider.currentProfile {
	 case cenums.PROFILE_MIN:
		 health = "Poor"
		 if rand.Float64() > 0.7 {
			 health = "Fair"
		 }
	 case cenums.PROFILE_MAX:
		 health = "Good"
	 default: 
		 health = "Good"
		 if rand.Float64() > 0.8 {
			 health = "Fair"
		 }
	 }
 
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
	 currentTime := float64(time.Now().UnixNano()) / 1e9
	 hourOfDay := float64(time.Now().Hour())
	 
	 baseGeneration := math.Sin((hourOfDay-6)*math.Pi/12) * 1000 
	 if hourOfDay < 6 || hourOfDay > 18 {
		 baseGeneration = 0 
	 }
	 
	 powerGeneration := baseGeneration * s.weatherPattern * (1 - s.cloudCoverFactor)
	 powerGeneration += math.Sin(currentTime*2)*s.microWeatherChanges*100
	 powerGeneration = math.Max(0, powerGeneration)
	 
	 controllerStatus := 0.0
	 if powerGeneration > 10.0 {
		 controllerStatus = 1.0
	 }
	 
	 var controllerMode string
	 batteryCap := s.metricsProvider.battery.lastCapacity 
	 
	 if batteryCap < 70 && powerGeneration > 200 {
		 controllerMode = "Bulk" 
	 } else if batteryCap >= 70 && batteryCap < 90 && powerGeneration > 50 {
		 controllerMode = "Absorption" 
	 } else if batteryCap >= 90 && powerGeneration > 20 {
		 controllerMode = "Float" 
	 } else if dayOfWeek := time.Now().Weekday(); dayOfWeek == time.Sunday && powerGeneration > 300 {
			 controllerMode = "Equalization"
	 } else {
			 controllerMode = "Idle"
	 }
	 batteryVoltage := 12.0 
	 if s.metricsProvider.battery != nil {
		 batteryMetrics, _ := s.metricsProvider.battery.GetMetrics()
		 if batteryMetrics != nil {
			 batteryVoltage = batteryMetrics.Voltage
		 }
	 }
	 
	 var controllerVoltage float64
	 switch controllerMode {
	 case "Bulk":
		 controllerVoltage = batteryVoltage + 1.0 + (math.Sin(currentTime*2.5)*0.2)
	 case "Absorption":
		 controllerVoltage = batteryVoltage + 0.5 + (math.Sin(currentTime*2.5)*0.1)
	 case "Float":
		 controllerVoltage = batteryVoltage + 0.2 + (math.Sin(currentTime*2.5)*0.05)
	 case "Equalization":
		 controllerVoltage = batteryVoltage + 1.5 + (math.Sin(currentTime*2.5)*0.3)
	 default:
		 controllerVoltage = batteryVoltage
	 }
	 
	 var controllerCurrent float64
	 if controllerStatus > 0 {
		 switch controllerMode {
		 case "Bulk":
			 controllerCurrent = math.Min(powerGeneration/14.0, 20.0) 
		 case "Absorption":
			 controllerCurrent = math.Min(powerGeneration/14.5, 10.0) 
		 case "Float":
			 controllerCurrent = math.Min(powerGeneration/15.0, 2.0)  
		 case "Equalization":
			 controllerCurrent = math.Min(powerGeneration/14.0, 15.0) 
		 default:
			 controllerCurrent = 0
		 }
		 
		 controllerCurrent += math.Sin(currentTime*4)*0.2
		 controllerCurrent = math.Max(0, controllerCurrent)
	 }
	 
	 controllerTemp := 25.0 + (controllerCurrent * 0.5) + math.Sin(currentTime*0.5)*2
	 
	 controllerEfficiency := 97.0 - (math.Abs(controllerCurrent-10)/20)*5 - (math.Max(0, controllerTemp-30)/10)*2
	 controllerEfficiency = math.Max(80, math.Min(99, controllerEfficiency))
	 
	 panelVoltage := 24.0 + math.Sin(currentTime*3)*0.5 
	 panelCurrent := powerGeneration / panelVoltage      
	 inverterStatus := 1.0
	 if powerGeneration < 10.0 {
		 inverterStatus = 0.0
	 }
 
	 return &SolarMetrics{
		 PowerGeneration: powerGeneration,
		 EnergyTotal:    s.energyTotal,
		 PanelPower:     powerGeneration,
		 PanelCurrent:   panelCurrent,
		 PanelVoltage:   panelVoltage,
		 InverterStatus: inverterStatus,
		 ControllerStatus:    controllerStatus,
		 ControllerMode:      controllerMode,
		 ControllerCurrent:   controllerCurrent,
		 ControllerVoltage:   controllerVoltage,
		 ControllerTemp:      controllerTemp,
		 ControllerEfficiency: controllerEfficiency,
	 }
 }
 
 func (e *PrometheusExporter) StartMetricsCollection(ctx context.Context, interval time.Duration) error {
	 ticker := time.NewTicker(1 * time.Second)
	 defer ticker.Stop()
 
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
 
	 e.solarPowerGeneration.WithLabelValues("watts", e.siteId).Set(metrics.Solar.PowerGeneration)
	 e.solarEnergyTotal.WithLabelValues("kwh", e.siteId).Set(metrics.Solar.EnergyTotal)
	 e.solarPanelPower.WithLabelValues("watts", e.siteId).Set(metrics.Solar.PanelPower)
	 e.solarPanelCurrent.WithLabelValues("amps", e.siteId).Set(metrics.Solar.PanelCurrent)
	 e.solarPanelVoltage.WithLabelValues("volts", e.siteId).Set(metrics.Solar.PanelVoltage)
	 e.solarInverterStatus.WithLabelValues("status", e.siteId).Set(metrics.Solar.InverterStatus)
	 
	 e.chargeControllerStatus.WithLabelValues("status", e.siteId).Set(metrics.Solar.ControllerStatus)
	 e.chargeControllerCurrent.WithLabelValues("amps", e.siteId).Set(metrics.Solar.ControllerCurrent)
	 e.chargeControllerVoltage.WithLabelValues("volts", e.siteId).Set(metrics.Solar.ControllerVoltage)
	 e.chargeControllerTemp.WithLabelValues("celsius", e.siteId).Set(metrics.Solar.ControllerTemp)
	 e.chargeControllerEfficiency.WithLabelValues("percent", e.siteId).Set(metrics.Solar.ControllerEfficiency)
	 
	 modeValue := map[string]float64{
		 "Idle": -1,
		 "Bulk": 0,
		 "Absorption": 1,
		 "Float": 2,
		 "Equalization": 3,
	 }[metrics.Solar.ControllerMode]
	 
	 e.chargeControllerMode.WithLabelValues("mode", e.siteId).Set(modeValue)
 
	 e.batteryChargeStatus.WithLabelValues("capacity", e.siteId).Set(metrics.Battery.Capacity)
	 e.batteryVoltage.WithLabelValues("volts", e.siteId).Set(metrics.Battery.Voltage)
	 e.batteryHealth.WithLabelValues("status", e.siteId).Set(map[string]float64{
		 "Good": 1.0,
		 "Fair": 0.5,
		 "Poor": 0.0,
	 }[metrics.Battery.Health])
	 e.batteryCurrent.WithLabelValues("amps", e.siteId).Set(metrics.Battery.Current)
	 e.batteryTemperature.WithLabelValues("celsius", e.siteId).Set(metrics.Battery.Temperature)
 
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
 
 func (m *MetricsProvider) GetPowerStatus() (bool, error) {
	 // Always return true since scenarios are removed
	 return true, nil
 }