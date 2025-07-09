/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef YAML_CONFIG_H
#define YAML_CONFIG_H

#include <stdint.h>
#include <stdbool.h>
#include "gpio_controller.h"

#define MAX_TEMP_POINTS 32
#define YAML_CONFIG_PATH "./config/safety_config.yaml"

// Temperature compensation lookup point
typedef struct {
    float temperature_c;
    float carrier_voltage;
    float peak_voltage;
} TempVoltagePoint;

// Temperature compensation table for one FEM unit
typedef struct {
    TempVoltagePoint points[MAX_TEMP_POINTS];
    int num_points;
} TempCompensationTable;

// Safety configuration from YAML
typedef struct {
    // Safety settings
    bool enabled;
    uint32_t check_interval_ms;
    uint32_t max_violations;
    
    // Safety thresholds
    float max_reverse_power_dbm;
    float max_pa_current_a;
    float max_temperature_c;
    float min_temperature_c;
    float max_forward_power_dbm;
    
    // Temperature zones
    float temp_critical_high;
    float temp_warning_high;
    float temp_normal_high;
    float temp_normal_low;
    float temp_warning_low;
    float temp_critical_low;
    
    // DAC configuration
    float dac_min_voltage;
    float dac_max_voltage;
    int dac_resolution_bits;
    
    // Default voltages
    float default_carrier_voltage;
    float default_peak_voltage;
    float shutdown_voltage;
    float standby_voltage;
    
    // Temperature compensation tables
    TempCompensationTable fem1_temp_table;
    TempCompensationTable fem2_temp_table;
    
    // Monitoring settings
    uint32_t adc_sampling_rate_hz;
    uint32_t adc_averaging_samples;
    int adc_calibration_offset_mv;
    
    // Temperature monitoring
    char temp_sensor_type[16];
    uint8_t temp_i2c_addr_fem1;
    uint8_t temp_i2c_addr_fem2;
    int temp_resolution_bits;
    uint32_t temp_update_interval_ms;
    
    // Current monitoring
    float current_shunt_resistance;
    float current_max_rating;
    int current_alarm_threshold_percent;
    
    // Emergency actions
    bool emergency_immediate_shutdown;
    bool emergency_disable_tx_rf;
    bool emergency_disable_pa_vds;
    bool emergency_disable_28v_vds;
    bool emergency_log_event;
    
} YamlSafetyConfig;

// Function declarations
int yaml_config_load(const char *filename, YamlSafetyConfig *config);
int yaml_config_get_dac_voltages_for_temp(const YamlSafetyConfig *config, FemUnit unit, 
                                          float temperature, float *carrier_voltage, float *peak_voltage);
void yaml_config_print(const YamlSafetyConfig *config);
int yaml_config_validate(const YamlSafetyConfig *config);
void yaml_config_set_defaults(YamlSafetyConfig *config);

#endif /* YAML_CONFIG_H */