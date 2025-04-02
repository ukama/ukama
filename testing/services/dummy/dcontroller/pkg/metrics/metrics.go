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
	"math/rand"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	cenums "github.com/ukama/ukama/testing/common/enums"
)

// MetricsRegistry is a singleton that holds all Prometheus metrics
 type MetricsRegistry struct {
	 backhaulSwitchPortStatus *prometheus.GaugeVec
	 backhaulSwitchPortSpeed  *prometheus.GaugeVec
	 backhaulSwitchPortPower  *prometheus.GaugeVec
 
	 solarSwitchPortStatus *prometheus.GaugeVec
	 solarSwitchPortSpeed  *prometheus.GaugeVec
	 solarSwitchPortPower  *prometheus.GaugeVec
 
	 nodeSwitchPortStatus *prometheus.GaugeVec
	 nodeSwitchPortSpeed  *prometheus.GaugeVec
	 nodeSwitchPortPower  *prometheus.GaugeVec
 
	 backhaulLatency         *prometheus.GaugeVec
	 backhaulSpeed           *prometheus.GaugeVec
	 switchPortStatus        *prometheus.GaugeVec
	 switchPortSpeed         *prometheus.GaugeVec
	 switchPortPower         *prometheus.GaugeVec
	 batteryChargePercentage *prometheus.GaugeVec
	 solarPanelVoltage       *prometheus.GaugeVec
	 solarPanelCurrent       *prometheus.GaugeVec
	 solarPanelPower         *prometheus.GaugeVec
	 siteUp                  *prometheus.GaugeVec
 
	 initialized bool
	 mutex       sync.Mutex
 }
 
 // Global singleton registry
 var globalRegistry = &MetricsRegistry{}
 
 // Initialize creates all Prometheus metrics once
 func (r *MetricsRegistry) Initialize() {
	 r.mutex.Lock()
	 defer r.mutex.Unlock()
 
	 if r.initialized {
		 return
	 }
 
	 // Initialize all metrics once
	 r.backhaulSwitchPortStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "backhaul_switch_port_status",
		 Help: "Backhaul switch port status (1 = up, 0 = down)",
	 }, []string{"unit", "site"})
	 r.backhaulSwitchPortSpeed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "backhaul_switch_port_speed",
		 Help: "Backhaul switch port speed in Mbps",
	 }, []string{"unit", "site"})
	 r.backhaulSwitchPortPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "backhaul_switch_port_power",
		 Help: "Backhaul switch port power in watts",
	 }, []string{"unit", "site"})
 
	 r.solarSwitchPortStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "solar_switch_port_status",
		 Help: "Solar controller switch port status (1 = up, 0 = down)",
	 }, []string{"unit", "site"})
	 r.solarSwitchPortSpeed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "solar_switch_port_speed",
		 Help: "Solar controller switch port speed in Mbps",
	 }, []string{"unit", "site"})
	 r.solarSwitchPortPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "solar_switch_port_power",
		 Help: "Solar controller switch port power in watts",
	 }, []string{"unit", "site"})
 
	 r.nodeSwitchPortStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "node_switch_port_status",
		 Help: "Node switch port status (1 = up, 0 = down)",
	 }, []string{"unit", "site"})
	 r.nodeSwitchPortSpeed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "node_switch_port_speed",
		 Help: "Node switch port speed in Mbps",
	 }, []string{"unit", "site"})
	 r.nodeSwitchPortPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "node_switch_port_power",
		 Help: "Node switch port power in watts",
	 }, []string{"unit", "site"})
 
	 r.backhaulLatency = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "main_backhaul_latency",
		 Help: "Backhaul latency in milliseconds",
	 }, []string{"unit", "site"})
	 r.backhaulSpeed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "backhaul_speed",
		 Help: "Backhaul speed in Mbps",
	 }, []string{"unit", "site"})
	 r.switchPortStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "switch_port_status",
		 Help: "Switch port status (1 = up, 0 = down)",
	 }, []string{"unit", "site"})
	 r.switchPortSpeed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "switch_port_speed",
		 Help: "Switch port speed in Mbps",
	 }, []string{"unit", "site"})
	 r.switchPortPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "switch_port_power",
		 Help: "Switch port power in watts",
	 }, []string{"unit", "site"})
	 r.batteryChargePercentage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "battery_charge_percentage",
		 Help: "Battery charge percentage",
	 }, []string{"unit", "site"})
	 r.solarPanelVoltage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "solar_panel_voltage",
		 Help: "Solar panel voltage in volts",
	 }, []string{"unit", "site"})
	 r.solarPanelCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "solar_panel_current",
		 Help: "Solar panel current in amperes",
	 }, []string{"unit", "site"})
	 r.solarPanelPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "solar_panel_power",
		 Help: "Solar panel power in watts",
	 }, []string{"unit", "site"})
	 r.siteUp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		 Name: "site_uptime_seconds",
		 Help: "Site uptime in seconds since last outage",
	 }, []string{"site"})
 
	 r.initialized = true
	 log.Info("Global metrics registry initialized")
 }
 
 type PrometheusExporter struct {
	 metricsProvider *MetricsProvider
	 siteId          string
	 shutdown        chan struct{}
	 registry        *MetricsRegistry
 }
 
 func (e *PrometheusExporter) IncrementUptimeCounter(seconds float64) {
	 e.registry.siteUp.WithLabelValues(e.siteId).Add(seconds)
 }
 
 func (e *PrometheusExporter) ResetUptimeCounter() {
	 e.registry.siteUp.WithLabelValues(e.siteId).Set(0)
	 log.Infof("Reset uptime counter for site %s", e.siteId)
 }
 
 func NewPrometheusExporter(metricsProvider *MetricsProvider, siteId string) *PrometheusExporter {
	 // Initialize the global registry if not already done
	 globalRegistry.Initialize()
	 
	 exporter := &PrometheusExporter{
		 metricsProvider: metricsProvider,
		 siteId:          siteId,
		 shutdown:        make(chan struct{}),
		 registry:        globalRegistry,
	 }
	 
	 log.Infof("Created new PrometheusExporter for site %s using shared registry", siteId)
	 return exporter
 }
 
 type BackhaulMetrics struct {
	 Latency         float64
	 Speed           float64
	 SwitchStatus    float64
	 SwitchBandwidth float64
	 SwitchPortPower float64
 }
 
 type BatteryMetrics struct {
	 Voltage float64
	 Current float64
	 Power   float64
 }
 
 type SolarMetrics struct {
	 PanelPower   float64
	 PanelVoltage float64
	 PanelCurrent float64
 }
 
 type ControllerMetrics struct {
	 Backhaul *BackhaulMetrics
	 Battery  *BatteryMetrics
	 Solar    *SolarMetrics
	 Time     time.Time
 }
 
 const (
	 PORT_NODE             = 1
	 PORT_SOLAR_CONTROLLER = 3
	 PORT_BACKHAUL         = 4
 )
 
 type BackhaulProvider struct {
	 metricsProvider *MetricsProvider
 }
 
 type BatteryProvider struct {
	 metricsProvider *MetricsProvider
 }
 
 type SolarProvider struct {
	 metricsProvider *MetricsProvider
 }
 
 type MetricsProvider struct {
	 backhaul       *BackhaulProvider
	 battery        *BatteryProvider
	 solar          *SolarProvider
	 portStatus     map[int]bool
	 currentProfile cenums.Profile
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
		 backhaul:  &BackhaulProvider{},
		 battery:   &BatteryProvider{},
		 solar:     &SolarProvider{},
		 portStatus: map[int]bool{
			 PORT_NODE:             true,
			 PORT_SOLAR_CONTROLLER: true,
			 PORT_BACKHAUL:         true,
		 },
		 currentProfile: cenums.PROFILE_NORMAL,
	 }
	 mp.backhaul.UpdateMetricsProvider(mp)
	 mp.battery.UpdateMetricsProvider(mp)
	 mp.solar.UpdateMetricsProvider(mp)
	 return mp
 }
 
 func (m *MetricsProvider) SetProfile(profile cenums.Profile) {
	 m.currentProfile = profile
 }
 
 func (m *MetricsProvider) GetPortStatus(port int) bool {
	 status, exists := m.portStatus[port]
	 if !exists {
		 return false
	 }
	 return status
 }
 
 func (m *MetricsProvider) GetMetrics(siteId string) (*ControllerMetrics, error) {
	 backhaulMetrics := m.backhaul.GetMetrics()
	 batteryMetrics := m.battery.GetMetrics()
	 solarMetrics := m.solar.GetMetrics()
 
	 log.Debugf("Current port statuses for site %s: %+v", siteId, m.portStatus)
 
	 if !m.portStatus[PORT_BACKHAUL] {
		 log.Infof("Backhaul port is OFF for site %s, zeroing backhaul metrics", siteId)
		 backhaulMetrics.Latency = 0
		 backhaulMetrics.Speed = 0
		 backhaulMetrics.SwitchStatus = 0
		 backhaulMetrics.SwitchBandwidth = 0
		 backhaulMetrics.SwitchPortPower = 0
	 }
 
	 if !m.portStatus[PORT_SOLAR_CONTROLLER] {
		 solarMetrics.PanelPower = 0
		 solarMetrics.PanelVoltage = 0
		 solarMetrics.PanelCurrent = 0
		 batteryMetrics.Voltage = 0
		 batteryMetrics.Current = 0
		 batteryMetrics.Power = 0
		 
	 }
 
	 return &ControllerMetrics{
		 Backhaul: backhaulMetrics,
		 Battery:  batteryMetrics,
		 Solar:    solarMetrics,
		 Time:     time.Now(),
	 }, nil
 }
 
 func (b *BackhaulProvider) GetMetrics() *BackhaulMetrics {
	 if (!b.metricsProvider.portStatus[PORT_BACKHAUL]) {
		 log.Debugf("Backhaul port is OFF in BackhaulProvider.GetMetrics")
		 return &BackhaulMetrics{
			 Latency:         0,
			 Speed:           0,
			 SwitchStatus:    0,
			 SwitchBandwidth: 0,
			 SwitchPortPower: 0,
		 }
	 }
 
	 var latency, speed, switchBandwidth, switchPortPower float64
	 profile := b.metricsProvider.currentProfile
 
	 switch profile {
	 case cenums.PROFILE_MIN:
		 latency = 150 + rand.Float64()*(250-150)
		 speed = 0.5 + rand.Float64()*(2.5-0.5)
		 switchBandwidth = 10 + rand.Float64()*(50-10)
	 case cenums.PROFILE_NORMAL:
		 latency = 30 + rand.Float64()*(50-30)
		 speed = 20 + rand.Float64()*(50-20)
		 switchBandwidth = 100 + rand.Float64()*(200-100)
	 case cenums.PROFILE_MAX:
		 latency = 5 + rand.Float64()*(15)
		 speed = 100 + rand.Float64()*(100)
		 switchBandwidth = 500 + rand.Float64()*(500)
	 }
 
	 if profile == cenums.PROFILE_MAX {
		 switchPortPower = 6 + rand.Float64()
	 } else {
		 switchPortPower = 5 + rand.Float64()*(7-5)
	 }
 
	 return &BackhaulMetrics{
		 Latency:         latency,
		 Speed:           speed,
		 SwitchStatus:    1.0,
		 SwitchBandwidth: switchBandwidth,
		 SwitchPortPower: switchPortPower,
	 }
 }
 
 func (b *BatteryProvider) GetMetrics() *BatteryMetrics {
	 var voltage, current float64
	 profile := b.metricsProvider.currentProfile
 
	 switch profile {
	 case cenums.PROFILE_MIN:
		 voltage = 10.5 + rand.Float64()*(12.0-10.5)
		 current = -2.0 + rand.Float64()*(0-(-2.0))
	 case cenums.PROFILE_NORMAL:
		 voltage = 12.0 + rand.Float64()*(12.5-12.0)
		 current = -1.4 + rand.Float64()*(2.8-(-1.4))
	 case cenums.PROFILE_MAX:
		 voltage = 12.0 + rand.Float64()
		 current = 0.5 + rand.Float64()*(4.0)
	 }
 
	 power := voltage * current
	 if current < 0 {
		 power = 0
	 }
 
	 return &BatteryMetrics{
		 Voltage: voltage,
		 Current: current,
		 Power:   power,
	 }
 }
 
 func (s *SolarProvider) GetMetrics() *SolarMetrics {
	 var panelPower, panelVoltage, panelCurrent float64
	 profile := s.metricsProvider.currentProfile
 
	 switch profile {
	 case cenums.PROFILE_MIN:
		 panelPower = 100 + rand.Float64()*(500-100)
		 panelVoltage = 16 + rand.Float64()*(20-16)
	 case cenums.PROFILE_NORMAL:
		 panelPower = 100 + rand.Float64()*(800-100)
		 panelVoltage = 21 + rand.Float64()*(27-21)
	 case cenums.PROFILE_MAX:
		 panelPower = 500 + rand.Float64()*(500)
		 panelVoltage = 30 + rand.Float64()*(10)
	 }
 
	 if profile == cenums.PROFILE_MAX {
		 panelCurrent = panelPower / panelVoltage
		 if panelCurrent < 5 {
			 panelCurrent = 5 + rand.Float64()*(5)
		 }
	 } else {
		 panelCurrent = 2 + rand.Float64()*(10-2)
		 if panelVoltage > 0 {
			 panelCurrent = panelPower / panelVoltage
			 if panelCurrent < 2 {
				 panelCurrent = 2
			 } else if panelCurrent > 10 {
				 panelCurrent = 10
			 }
		 } else {
			 panelCurrent = 0
		 }
	 }
 
	 return &SolarMetrics{
		 PanelPower:   panelPower,
		 PanelVoltage: panelVoltage,
		 PanelCurrent: panelCurrent,
	 }
 }
 
 func (e *PrometheusExporter) StartMetricsCollection(ctx context.Context, interval time.Duration) error {
	 ticker := time.NewTicker(1 * time.Second)
	 defer ticker.Stop()
 
	 if err := e.CollectMetrics(); err != nil {
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
			 if err := e.CollectMetrics(); err != nil {
				 log.Errorf("Error collecting metrics: %v", err)
			 }
		 }
	 }
 }
 
 func (e *PrometheusExporter) Shutdown() {
	 close(e.shutdown)
 }
 
 func (e *PrometheusExporter) CollectMetrics() error {
	 metrics, err := e.metricsProvider.GetMetrics(e.siteId)
	 if err != nil {
		 return fmt.Errorf("failed to get metrics: %w", err)
	 }
 
	 log.Debugf("Collecting metrics for site %s: Backhaul speed: %f, Backhaul port status: %v",
		 e.siteId, metrics.Backhaul.Speed, e.metricsProvider.GetPortStatus(PORT_BACKHAUL))
 
	 // Set backhaul metrics
	 if e.metricsProvider.GetPortStatus(PORT_BACKHAUL) {
		 e.registry.backhaulLatency.WithLabelValues("ms", e.siteId).Set(metrics.Backhaul.Latency)
		 e.registry.backhaulSpeed.WithLabelValues("mbps", e.siteId).Set(metrics.Backhaul.Speed)
		 e.registry.switchPortStatus.WithLabelValues("status", e.siteId).Set(metrics.Backhaul.SwitchStatus)
		 e.registry.switchPortSpeed.WithLabelValues("mbps", e.siteId).Set(metrics.Backhaul.SwitchBandwidth)
		 e.registry.switchPortPower.WithLabelValues("watts", e.siteId).Set(metrics.Backhaul.SwitchPortPower)
	 } else {
		 log.Infof("Backhaul port is OFF for site %s, setting Prometheus metrics to 0", e.siteId)
		 e.registry.backhaulLatency.WithLabelValues("ms", e.siteId).Set(0)
		 e.registry.backhaulSpeed.WithLabelValues("mbps", e.siteId).Set(0)
		 e.registry.switchPortStatus.WithLabelValues("status", e.siteId).Set(0)
		 e.registry.switchPortSpeed.WithLabelValues("mbps", e.siteId).Set(0)
		 e.registry.switchPortPower.WithLabelValues("watts", e.siteId).Set(0)
	 }
 
	 if e.metricsProvider.GetPortStatus(PORT_BACKHAUL) {
		 e.registry.backhaulSwitchPortStatus.WithLabelValues("status", e.siteId).Set(1)
		 e.registry.backhaulSwitchPortSpeed.WithLabelValues("mbps", e.siteId).Set(metrics.Backhaul.SwitchBandwidth)
		 e.registry.backhaulSwitchPortPower.WithLabelValues("watts", e.siteId).Set(metrics.Backhaul.SwitchPortPower)
	 } else {
		 e.registry.backhaulSwitchPortStatus.WithLabelValues("status", e.siteId).Set(0)
		 e.registry.backhaulSwitchPortSpeed.WithLabelValues("mbps", e.siteId).Set(0)
		 e.registry.backhaulSwitchPortPower.WithLabelValues("watts", e.siteId).Set(0)
	 }
 
	 if e.metricsProvider.GetPortStatus(PORT_SOLAR_CONTROLLER) {
		 e.registry.solarSwitchPortStatus.WithLabelValues("status", e.siteId).Set(1)
		 solarSpeed := 60 + rand.Float64()*(120-60)  
		 solarPower := 10 + rand.Float64()*(30-10)  
		 e.registry.solarSwitchPortSpeed.WithLabelValues("mbps", e.siteId).Set(solarSpeed)
		 e.registry.solarSwitchPortPower.WithLabelValues("watts", e.siteId).Set(solarPower)
	 } else {
		 e.registry.solarSwitchPortStatus.WithLabelValues("status", e.siteId).Set(0)
		 e.registry.solarSwitchPortSpeed.WithLabelValues("mbps", e.siteId).Set(0)
		 e.registry.solarSwitchPortPower.WithLabelValues("watts", e.siteId).Set(0)
	 }
 
	 if e.metricsProvider.GetPortStatus(PORT_NODE) {
		 e.registry.nodeSwitchPortStatus.WithLabelValues("status", e.siteId).Set(1)
		 nodeSpeed := 10 + rand.Float64()*(50-10)   
		 nodePower := 5 + rand.Float64()*(10-5)    
		 e.registry.nodeSwitchPortSpeed.WithLabelValues("mbps", e.siteId).Set(nodeSpeed)
		 e.registry.nodeSwitchPortPower.WithLabelValues("watts", e.siteId).Set(nodePower)
	 } else {
		 e.registry.nodeSwitchPortStatus.WithLabelValues("status", e.siteId).Set(0)
		 e.registry.nodeSwitchPortSpeed.WithLabelValues("mbps", e.siteId).Set(0)
		 e.registry.nodeSwitchPortPower.WithLabelValues("watts", e.siteId).Set(0)
	 }
 
	 var percentage float64
	 voltage := metrics.Battery.Voltage
 
	 if !e.metricsProvider.GetPortStatus(PORT_SOLAR_CONTROLLER) {
		 e.registry.solarPanelVoltage.WithLabelValues("volts", e.siteId).Set(0)
		 e.registry.solarPanelCurrent.WithLabelValues("amps", e.siteId).Set(0)
		 e.registry.solarPanelPower.WithLabelValues("watts", e.siteId).Set(0)
		 e.registry.batteryChargePercentage.WithLabelValues("percent", e.siteId).Set(0)
		 voltage = 0
	 }
 
	 switch e.metricsProvider.currentProfile {
	 case cenums.PROFILE_MIN:
		 if voltage <= 10.0 {
			 percentage = 0
		 } else if voltage >= 12.0 {
			 percentage = 100
		 } else {
			 percentage = (voltage - 10.0) / (12.0 - 10.0) * 100
		 }
	 case cenums.PROFILE_MAX:
		 if voltage <= 12.0 {
			 percentage = 70
		 } else if voltage >= 13.0 {
			 percentage = 100
		 } else {
			 percentage = 70 + (voltage-12.0)/(13.0-12.0)*30
		 }
	 default:
		 if voltage <= 10.5 {
			 percentage = 0
		 } else if voltage >= 12.7 {
			 percentage = 100
		 } else {
			 percentage = (voltage - 10.5) / (12.7 - 10.5) * 100
		 }
	 }
	 e.registry.batteryChargePercentage.WithLabelValues("percent", e.siteId).Set(percentage)
 
	 e.registry.solarPanelVoltage.WithLabelValues("volts", e.siteId).Set(metrics.Solar.PanelVoltage)
	 e.registry.solarPanelCurrent.WithLabelValues("amps", e.siteId).Set(metrics.Solar.PanelCurrent)
	 e.registry.solarPanelPower.WithLabelValues("watts", e.siteId).Set(metrics.Solar.PanelPower)
 
	 log.Debugf("Site %s - Voltage: %f, Battery percentage: %f, Backhaul speed: %f, Switch port power: %f",
		 e.siteId, voltage, percentage, metrics.Backhaul.Speed, metrics.Backhaul.SwitchPortPower)
 
	 return nil
 }
 
 func (m *MetricsProvider) SetPortStatus(port int, status bool) error {
	 if port < 1 || port > 4 {
		 return fmt.Errorf("invalid port number: %d", port)
	 }
	 oldStatus := m.portStatus[port]
	 m.portStatus[port] = status
	 log.Infof("Port %d status changed from %v to %v", port, oldStatus, status)
	 log.Infof("Current port statuses: %+v", m.portStatus)
	 return nil
 }
 
 func (m *MetricsProvider) GetPowerStatus() (bool, error) {
	 return true, nil
 }