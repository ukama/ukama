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
	 // Backhaul metrics
	 backhaulLatency      *prometheus.GaugeVec
	 backhaulStatus       *prometheus.GaugeVec
	 backhaulSpeed        *prometheus.GaugeVec
 
	 // Ethernet switch metrics
	 switchPortStatus     *prometheus.GaugeVec
	 switchPortSpeed      *prometheus.GaugeVec
	 switchPortPower      *prometheus.GaugeVec
 
	 // Power metrics
	 batteryChargePercentage *prometheus.GaugeVec
	 solarPanelVoltage       *prometheus.GaugeVec
	 solarPanelCurrent       *prometheus.GaugeVec
	 solarPanelPower         *prometheus.GaugeVec
	 chargeControllerStatus  *prometheus.GaugeVec
	 chargeControllerMode    *prometheus.GaugeVec
	 chargeControllerCurrent *prometheus.GaugeVec
	 chargeControllerVoltage *prometheus.GaugeVec
 
	 metricsProvider *MetricsProvider
	 siteId          string
	 shutdown        chan struct{}
 }
 
 func NewPrometheusExporter(metricsProvider *MetricsProvider, siteId string) *PrometheusExporter {
	 exporter := &PrometheusExporter{
		 metricsProvider: metricsProvider,
		 siteId:          siteId,
		 shutdown:        make(chan struct{}),
 
		 // Backhaul metrics
		 backhaulLatency: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "main_backhaul_latency",
			 Help: "Backhaul latency in milliseconds",
		 }, []string{"unit", "site"}),
		 backhaulStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "backhaul_status",
			 Help: "Backhaul status (1 = up, 0 = down)",
		 }, []string{"unit", "site"}),
		 backhaulSpeed: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "backhaul_speed",
			 Help: "Backhaul speed in Mbps",
		 }, []string{"unit", "site"}),
 
		 // Ethernet switch metrics
		 switchPortStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "switch_port_status",
			 Help: "Switch port status (1 = up, 0 = down)",
		 }, []string{"unit", "site"}),
		 switchPortSpeed: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "switch_port_speed",
			 Help: "Switch port speed in Mbps",
		 }, []string{"unit", "site"}),
		 switchPortPower: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "switch_port_power",
			 Help: "Switch port power in watts",
		 }, []string{"unit", "site"}),
 
		 // Power metrics
		 batteryChargePercentage: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "battery_charge_percentage",
			 Help: "Battery charge percentage",
		 }, []string{"unit", "site"}),
		 solarPanelVoltage: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "solar_panel_voltage",
			 Help: "Solar panel voltage in volts",
		 }, []string{"unit", "site"}),
		 solarPanelCurrent: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "solar_panel_current",
			 Help: "Solar panel current in amperes",
		 }, []string{"unit", "site"}),
		 solarPanelPower: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "solar_panel_power",
			 Help: "Solar panel power in watts",
		 }, []string{"unit", "site"}),
		 chargeControllerStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_status",
			 Help: "Charge controller status (1 = working, 0 = not working)",
		 }, []string{"unit", "site"}),
		 chargeControllerMode: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_mode",
			 Help: "Charge controller mode (0 = Bulk, 1 = Absorption, 2 = Float, 3 = Equalization, -1 = Idle)",
		 }, []string{"unit", "site"}),
		 chargeControllerCurrent: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_current",
			 Help: "Charge controller current in amperes",
		 }, []string{"unit", "site"}),
		 chargeControllerVoltage: promauto.NewGaugeVec(prometheus.GaugeOpts{
			 Name: "charge_controller_voltage",
			 Help: "Charge controller voltage in volts",
		 }, []string{"unit", "site"}),
	 }
 
	 return exporter
 }
 
 type BackhaulMetrics struct {
	 Latency         float64
	 Status          float64
	 Speed           float64
	 SwitchStatus    float64
	 SwitchBandwidth float64
	 SwitchPortPower float64
 }
 
 type BatteryMetrics struct {
	 Voltage     float64
	 Current     float64
	 Power       float64
 }
 
 type SolarMetrics struct {
	 PanelPower          float64
	 PanelCurrent        float64
	 PanelVoltage        float64
	 ControllerStatus    float64
	 ControllerMode      string
	 ControllerModeValue int
	 ControllerCurrent   float64
	 ControllerVoltage   float64
 }
 
 type ControllerMetrics struct {
	 Backhaul *BackhaulMetrics
	 Battery  *BatteryMetrics
	 Solar    *SolarMetrics
	 Time     time.Time
 }
 
 type BackhaulProvider struct {
	 lastUpdate      time.Time
	 metricsProvider *MetricsProvider
	 jitterFactor    float64
	 noiseAmplitude  float64
 }
 
 const (
	 PORT_AMPLIFIER = 1
	 PORT_TOWER     = 2
	 PORT_SOLAR     = 3
	 PORT_BACKHAUL  = 4
 )
 
 type BatteryProvider struct {
	 lastUpdate             time.Time
	 voltage                float64
	 current                float64
	 metricsProvider        *MetricsProvider
	 microFluctuationFactor float64
 }
 
 type SolarProvider struct {
	 startTime           time.Time
	 metricsProvider     *MetricsProvider
	 cloudCoverFactor    float64
	 microWeatherChanges float64
 }
 
 type MetricsProvider struct {
	 backhaul               *BackhaulProvider
	 battery                *BatteryProvider
	 solar                  *SolarProvider
	 portStatus             map[int]bool
	 currentProfile         cenums.Profile
	 globalVariationFactor  float64
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
			 lastUpdate:     time.Now(),
			 jitterFactor:   rand.Float64() * 0.5,
			 noiseAmplitude: rand.Float64() * 0.3,
		 },
		 battery: &BatteryProvider{
			 lastUpdate:             time.Now(),
			 voltage:                12.5,
			 current:                2.0,
			 microFluctuationFactor: rand.Float64() * 0.2,
		 },
		 solar: &SolarProvider{
			 startTime:           time.Now(),
			 cloudCoverFactor:    rand.Float64() * 0.4,
			 microWeatherChanges: rand.Float64() * 0.3,
		 },
		 portStatus: map[int]bool{
			 PORT_AMPLIFIER: true,
			 PORT_TOWER:     true,
			 PORT_SOLAR:     true,
			 PORT_BACKHAUL:  true,
		 },
		 currentProfile:        cenums.PROFILE_NORMAL,
		 globalVariationFactor: rand.Float64(),
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
	 batteryMetrics := m.battery.GetMetrics()
	 solarMetrics := m.solar.GetMetrics()
 
	 if !m.portStatus[PORT_BACKHAUL] {
		 backhaulMetrics.Latency = 0
		 backhaulMetrics.Status = 0
		 backhaulMetrics.Speed = 0
		 backhaulMetrics.SwitchStatus = 0
		 backhaulMetrics.SwitchBandwidth = 0
		 backhaulMetrics.SwitchPortPower = 0
	 }
 
	 if !m.portStatus[PORT_SOLAR] {
		 solarMetrics.PanelPower = 0
		 solarMetrics.PanelCurrent = 0
		 solarMetrics.PanelVoltage = 0
		 solarMetrics.ControllerStatus = 0
		 solarMetrics.ControllerCurrent = 0
		 solarMetrics.ControllerVoltage = 0
	 }
 
	 if m.currentProfile == cenums.PROFILE_MIN && batteryMetrics.Voltage < 10.5 {
		 solarMetrics.ControllerStatus = 0
		 if rand.Float64() > 0.7 {
			 solarMetrics.PanelPower *= 0.5
		 }
	 }
 
	 return &ControllerMetrics{
		 Backhaul: backhaulMetrics,
		 Battery:  batteryMetrics,
		 Solar:    solarMetrics,
		 Time:     time.Now(),
	 }, nil
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
		 speed = math.Max(0.1, math.Min(150.0, speed)) 
 
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
 
	 var switchPortPower float64
	 if switchStatus == 1.0 {
		 switchPortPower = 5.0 + rand.Float64()*2.0 
	 } else {
		 switchPortPower = 0.0
	 }
 
	 return &BackhaulMetrics{
		 Latency:         latency,
		 Status:          status,
		 Speed:           speed,
		 SwitchStatus:    switchStatus,
		 SwitchBandwidth: switchBandwidth,
		 SwitchPortPower: switchPortPower,
	 }
 }
 
 func (b *BatteryProvider) GetMetrics() *BatteryMetrics {
	 currentTime := float64(time.Now().UnixNano()) / 1e9
 
	 microCycle := math.Sin(currentTime*2*math.Pi) * b.microFluctuationFactor
 
	 globalFactor := b.metricsProvider.globalVariationFactor
	 var baseVoltage, baseCurrent float64
 
	 switch b.metricsProvider.currentProfile {
	 case cenums.PROFILE_MIN:
		 baseVoltage = 10.0 + (1.5 * globalFactor)
		 baseCurrent = -1.5 - (0.5 * globalFactor) 
 
		 if baseVoltage < 10.2 && rand.Float64() > 0.7 {
			 baseVoltage = 0
			 baseCurrent = 0
		 }
	 case cenums.PROFILE_MAX:
		 baseVoltage = 12.0 + (0.8 * globalFactor)
 
		 // Determine if daytime for charging simulation
		 hourOfDay := float64(time.Now().Hour())
		 daytime := hourOfDay >= 6 && hourOfDay <= 18
 
		 if daytime {
			 baseCurrent = 3.0 + (1.5 * globalFactor) 
		 } else {
			 baseCurrent = -0.5 - (0.3 * globalFactor) 
		 }
	 default:
		 baseVoltage = 12.0 + (0.5 * globalFactor)
 
		 hourOfDay := float64(time.Now().Hour())
		 daytime := hourOfDay >= 6 && hourOfDay <= 18
 
		 if daytime {
			 baseCurrent = 2.0 + (0.8 * globalFactor)
		 } else {
			 baseCurrent = -1.0 - (0.4 * globalFactor)
		 }
	 }
 
	 voltage := baseVoltage + (microCycle * 0.1)
	 current := baseCurrent + (microCycle * 0.2)
 
	 power := voltage * math.Abs(current)
 
	 b.voltage = voltage
	 b.current = current
 
	 return &BatteryMetrics{
		 Voltage: voltage,
		 Current: current,
		 Power:   power,
	 }
 }
 
 func (s *SolarProvider) GetMetrics() *SolarMetrics {
	 currentTime := float64(time.Now().UnixNano()) / 1e9
	 hourOfDay := float64(time.Now().Hour())
 
	 baseGeneration := math.Sin((hourOfDay-6)*math.Pi/12) * 1000
	 if hourOfDay < 6 || hourOfDay > 18 {
		 baseGeneration = 0
	 }
 
	 powerGeneration := baseGeneration * (1 - s.cloudCoverFactor)
	 powerGeneration += math.Sin(currentTime*2)*s.microWeatherChanges*100
	 powerGeneration = math.Max(0, powerGeneration)
 
	 // Apply profile adjustments
	 switch s.metricsProvider.currentProfile {
	 case cenums.PROFILE_MIN:
		 powerGeneration *= 0.5
	 case cenums.PROFILE_MAX:
		 powerGeneration *= 1.2
	 }
 
	 controllerStatus := 0.0
	 if powerGeneration > 10.0 {
		 controllerStatus = 1.0
	 }
 
	 var controllerMode string
	 var controllerModeValue int
 
	 // Base controller mode on battery conditions and time of day
	 if s.metricsProvider.battery != nil {
		 batteryMetrics := s.metricsProvider.battery.GetMetrics()
		 batteryVoltage := batteryMetrics.Voltage
 
		 if hourOfDay < 6 || hourOfDay > 18 {
			 controllerMode = "Idle"
			 controllerModeValue = -1
		 } else if batteryVoltage < 11.5 && powerGeneration > 100 {
			 controllerMode = "Bulk"
			 controllerModeValue = 0
		 } else if batteryVoltage >= 11.5 && batteryVoltage < 13.0 && powerGeneration > 50 {
			 controllerMode = "Absorption"
			 controllerModeValue = 1
		 } else if batteryVoltage >= 13.0 && powerGeneration > 20 {
			 controllerMode = "Float"
			 controllerModeValue = 2
		 } else if time.Now().Weekday() == time.Sunday && powerGeneration > 300 {
			 controllerMode = "Equalization"
			 controllerModeValue = 3
		 } else {
			 controllerMode = "Idle"
			 controllerModeValue = -1
		 }
	 } else {
		 controllerMode = "Idle"
		 controllerModeValue = -1
	 }
 
	 var controllerVoltage float64
	 var batteryVoltage float64 = 12.0
 
	 if s.metricsProvider.battery != nil {
		 batteryVoltage = s.metricsProvider.battery.voltage
	 }
 
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
 
	 panelVoltage := 0.0
	 panelCurrent := 0.0
 
	 if powerGeneration > 0 {
		 panelVoltage = 18.0 + math.Sin(currentTime*3)*0.5
 
		 switch s.metricsProvider.currentProfile {
		 case cenums.PROFILE_MIN:
			 panelVoltage = math.Max(10.0, panelVoltage*0.7)
		 case cenums.PROFILE_MAX:
			 panelVoltage = math.Min(24.0, panelVoltage*1.2)
		 }
 
		 panelVoltage = math.Min(100.0, panelVoltage) 
		 panelCurrent = powerGeneration / panelVoltage
		 panelCurrent = math.Min(20.0, panelCurrent) 
	 }
 
	 return &SolarMetrics{
		 PanelPower:          powerGeneration,
		 PanelCurrent:        panelCurrent,
		 PanelVoltage:        panelVoltage,
		 ControllerStatus:    controllerStatus,
		 ControllerMode:      controllerMode,
		 ControllerModeValue: controllerModeValue,
		 ControllerCurrent:   controllerCurrent,
		 ControllerVoltage:   controllerVoltage,
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
 
	 log.Debugf("Collecting metrics for site %s: Solar power: %f, Battery power: %f",
		 e.siteId, metrics.Solar.PanelPower, metrics.Battery.Power)
 
	 // Backhaul metrics
	 e.backhaulLatency.WithLabelValues("ms", e.siteId).Set(metrics.Backhaul.Latency)
	 e.backhaulStatus.WithLabelValues("status", e.siteId).Set(metrics.Backhaul.Status)
	 e.backhaulSpeed.WithLabelValues("mbps", e.siteId).Set(metrics.Backhaul.Speed)
 
	 // Switch metrics
	 e.switchPortStatus.WithLabelValues("status", e.siteId).Set(metrics.Backhaul.SwitchStatus)
	 e.switchPortSpeed.WithLabelValues("mbps", e.siteId).Set(metrics.Backhaul.SwitchBandwidth)
	 e.switchPortPower.WithLabelValues("watts", e.siteId).Set(metrics.Backhaul.SwitchPortPower)
 
	 // Power metrics
	 voltage := metrics.Battery.Voltage
	 var percentage float64
	 if voltage <= 10.5 {
		 percentage = 0
	 } else if voltage >= 12.7 {
		 percentage = 100
	 } else {
		 percentage = (voltage - 10.5) / (12.7 - 10.5) * 100
	 }
	 e.batteryChargePercentage.WithLabelValues("percent", e.siteId).Set(percentage)
 
	 e.solarPanelVoltage.WithLabelValues("volts", e.siteId).Set(metrics.Solar.PanelVoltage)
	 e.solarPanelCurrent.WithLabelValues("amps", e.siteId).Set(metrics.Solar.PanelCurrent)
	 e.solarPanelPower.WithLabelValues("watts", e.siteId).Set(metrics.Solar.PanelPower)
 
	 // Charge controller metrics
	 e.chargeControllerStatus.WithLabelValues("status", e.siteId).Set(metrics.Solar.ControllerStatus)
	 e.chargeControllerMode.WithLabelValues("mode", e.siteId).Set(float64(metrics.Solar.ControllerModeValue))
	 e.chargeControllerCurrent.WithLabelValues("amps", e.siteId).Set(metrics.Solar.ControllerCurrent)
	 e.chargeControllerVoltage.WithLabelValues("volts", e.siteId).Set(metrics.Solar.ControllerVoltage)
 
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
	 return true, nil
 }