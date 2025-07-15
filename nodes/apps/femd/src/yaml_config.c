/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

#include "yaml_config.h"
#include "femd.h"


static int parse_float_value(const char *line, const char *key, float *value) {
    char search_key[64];
    snprintf(search_key, sizeof(search_key), "%s:", key);
    
    char *key_pos = strstr(line, search_key);
    if (!key_pos) return STATUS_NOK;
    
    char *value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;
    
    *value = strtof(value_start, NULL);
    return STATUS_OK;
}

static int parse_int_value(const char *line, const char *key, int *value) {
    char search_key[64];
    snprintf(search_key, sizeof(search_key), "%s:", key);
    
    char *key_pos = strstr(line, search_key);
    if (!key_pos) return STATUS_NOK;
    
    char *value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;
    
    *value = atoi(value_start);
    return STATUS_OK;
}

static int parse_uint32_value(const char *line, const char *key, uint32_t *value) {
    char search_key[64];
    snprintf(search_key, sizeof(search_key), "%s:", key);
    
    char *key_pos = strstr(line, search_key);
    if (!key_pos) return STATUS_NOK;
    
    char *value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;
    
    *value = (uint32_t)strtoul(value_start, NULL, 10);
    return STATUS_OK;
}

static int parse_bool_value(const char *line, const char *key, bool *value) {
    char search_key[64];
    snprintf(search_key, sizeof(search_key), "%s:", key);
    
    char *key_pos = strstr(line, search_key);
    if (!key_pos) return STATUS_NOK;
    
    char *value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;
    
    if (strncmp(value_start, "true", 4) == 0) {
        *value = true;
    } else if (strncmp(value_start, "false", 5) == 0) {
        *value = false;
    } else {
        return STATUS_NOK;
    }
    
    return STATUS_OK;
}

static int parse_string_value(const char *line, const char *key, char *value, size_t max_len) {
    char search_key[64];
    snprintf(search_key, sizeof(search_key), "%s:", key);
    
    char *key_pos = strstr(line, search_key);
    if (!key_pos) return STATUS_NOK;
    
    char *value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;
    
    if (*value_start == '"') {
        value_start++;
        char *end_quote = strchr(value_start, '"');
        if (end_quote) {
            size_t len = end_quote - value_start;
            if (len < max_len) {
                strncpy(value, value_start, len);
                value[len] = '\0';
                return STATUS_OK;
            }
        }
    } else {
        char *end = strchr(value_start, '\n');
        if (!end) end = strchr(value_start, '#');
        if (!end) end = value_start + strlen(value_start);
        
        size_t len = end - value_start;
        while (len > 0 && (value_start[len-1] == ' ' || value_start[len-1] == '\t' || value_start[len-1] == '\r')) {
            len--;
        }
        
        if (len < max_len) {
            strncpy(value, value_start, len);
            value[len] = '\0';
            return STATUS_OK;
        }
    }
    
    return STATUS_NOK;
}

static int parse_temp_voltage_line(const char *line, TempVoltagePoint *point) {
    
    char *colon = strchr(line, ':');
    if (!colon) return STATUS_NOK;
    
    char temp_str[16];
    const char *temp_start = line;
    while (*temp_start == ' ' || *temp_start == '\t') temp_start++;
    
    size_t temp_len = colon - temp_start;
    if (temp_len >= sizeof(temp_str)) return STATUS_NOK;
    
    strncpy(temp_str, temp_start, temp_len);
    temp_str[temp_len] = '\0';
    point->temperature_c = strtof(temp_str, NULL);
    
    char *open_brace = strchr(colon, '{');
    char *close_brace = strchr(colon, '}');
    if (!open_brace || !close_brace) return STATUS_NOK;
    
    char *carrier_pos = strstr(open_brace, "carrier:");
    if (!carrier_pos || carrier_pos > close_brace) return STATUS_NOK;
    
    carrier_pos += 8; // Skip "carrier:"
    while (*carrier_pos == ' ' || *carrier_pos == '\t') carrier_pos++;
    point->carrier_voltage = strtof(carrier_pos, NULL);
    
    char *peak_pos = strstr(open_brace, "peak:");
    if (!peak_pos || peak_pos > close_brace) return STATUS_NOK;
    
    peak_pos += 5; // Skip "peak:"
    while (*peak_pos == ' ' || *peak_pos == '\t') peak_pos++;
    point->peak_voltage = strtof(peak_pos, NULL);
    
    return STATUS_OK;
}

void yaml_config_set_defaults(YamlSafetyConfig *config) {
    memset(config, 0, sizeof(YamlSafetyConfig));
    
    config->enabled = true;
    config->check_interval_ms = 1000;
    config->max_violations = 3;
    
    config->max_reverse_power_dbm = -10.0f;
    config->max_pa_current_a = 5.0f;
    config->max_temperature_c = 85.0f;
    config->min_temperature_c = -40.0f;
    config->max_forward_power_dbm = 30.0f;
    
    config->temp_critical_high = 85.0f;
    config->temp_warning_high = 75.0f;
    config->temp_normal_high = 65.0f;
    config->temp_normal_low = 0.0f;
    config->temp_warning_low = -20.0f;
    config->temp_critical_low = -40.0f;
    
    config->dac_min_voltage = 0.0f;
    config->dac_max_voltage = 2.5f;
    config->dac_resolution_bits = 12;
    
    config->default_carrier_voltage = 1.2f;
    config->default_peak_voltage = 2.0f;
    config->shutdown_voltage = 0.0f;
    config->standby_voltage = 0.5f;
    
    config->adc_sampling_rate_hz = 1000;
    config->adc_averaging_samples = 10;
    config->adc_calibration_offset_mv = 0;
    
    strcpy(config->temp_sensor_type, "LM75A");
    config->temp_i2c_addr_fem1 = 0x48;
    config->temp_i2c_addr_fem2 = 0x49;
    config->temp_resolution_bits = 12;
    config->temp_update_interval_ms = 2000;
    
    config->current_shunt_resistance = 0.01f;
    config->current_max_rating = 10.0f;
    config->current_alarm_threshold_percent = 80;
    
    config->emergency_immediate_shutdown = true;
    config->emergency_disable_tx_rf = true;
    config->emergency_disable_pa_vds = true;
    config->emergency_disable_28v_vds = true;
    config->emergency_log_event = true;
}

int yaml_config_load(const char *filename, YamlSafetyConfig *config) {
    FILE *file;
    char line[512];
    int in_fem1_section = 0;
    int in_fem2_section = 0;
    int fem1_points = 0;
    int fem2_points = 0;
    
    yaml_config_set_defaults(config);
    
    file = fopen(filename, "r");
    if (!file) {
        usys_log_warn("Could not open YAML config file %s, using defaults", filename);
        return STATUS_NOK;
    }
    
    usys_log_info("Loading YAML configuration from %s", filename);
    
    while (fgets(line, sizeof(line), file)) {
        if (line[0] == '#' || line[0] == '\n' || line[0] == '\r') continue;
        
        if (strstr(line, "fem1:")) {
            in_fem1_section = 1;
            in_fem2_section = 0;
            continue;
        } else if (strstr(line, "fem2:")) {
            in_fem1_section = 0;
            in_fem2_section = 1;
            continue;
        } else if (strstr(line, "voltage_lookup:")) {
            continue; // Stay in current FEM section
        } else if (line[0] != ' ' && line[0] != '\t') {
            in_fem1_section = 0;
            in_fem2_section = 0;
        }
        
        if (in_fem1_section && fem1_points < MAX_TEMP_POINTS) {
            TempVoltagePoint point;
            if (parse_temp_voltage_line(line, &point) == STATUS_OK) {
                config->fem1_temp_table.points[fem1_points] = point;
                fem1_points++;
            }
        } else if (in_fem2_section && fem2_points < MAX_TEMP_POINTS) {
            TempVoltagePoint point;
            if (parse_temp_voltage_line(line, &point) == STATUS_OK) {
                config->fem2_temp_table.points[fem2_points] = point;
                fem2_points++;
            }
        }
        
        parse_bool_value(line, "enabled", &config->enabled);
        parse_uint32_value(line, "check_interval_ms", &config->check_interval_ms);
        parse_uint32_value(line, "max_violations_before_shutdown", &config->max_violations);
        
        parse_float_value(line, "max_reverse_power_dbm", &config->max_reverse_power_dbm);
        parse_float_value(line, "max_pa_current_a", &config->max_pa_current_a);
        parse_float_value(line, "max_temperature_c", &config->max_temperature_c);
        parse_float_value(line, "min_temperature_c", &config->min_temperature_c);
        parse_float_value(line, "max_forward_power_dbm", &config->max_forward_power_dbm);
        
        parse_float_value(line, "critical_high", &config->temp_critical_high);
        parse_float_value(line, "warning_high", &config->temp_warning_high);
        parse_float_value(line, "normal_high", &config->temp_normal_high);
        parse_float_value(line, "normal_low", &config->temp_normal_low);
        parse_float_value(line, "warning_low", &config->temp_warning_low);
        parse_float_value(line, "critical_low", &config->temp_critical_low);
        
        parse_float_value(line, "min_voltage", &config->dac_min_voltage);
        parse_float_value(line, "max_voltage", &config->dac_max_voltage);
        parse_int_value(line, "resolution_bits", &config->dac_resolution_bits);
        
        parse_float_value(line, "carrier_voltage", &config->default_carrier_voltage);
        parse_float_value(line, "peak_voltage", &config->default_peak_voltage);
        parse_float_value(line, "shutdown_voltage", &config->shutdown_voltage);
        parse_float_value(line, "standby_voltage", &config->standby_voltage);
        
        parse_uint32_value(line, "sampling_rate_hz", &config->adc_sampling_rate_hz);
        parse_uint32_value(line, "averaging_samples", &config->adc_averaging_samples);
        parse_int_value(line, "calibration_offset_mv", &config->adc_calibration_offset_mv);
        
        parse_string_value(line, "sensor_type", config->temp_sensor_type, sizeof(config->temp_sensor_type));
        parse_uint32_value(line, "update_interval_ms", &config->temp_update_interval_ms);
        
        parse_float_value(line, "shunt_resistance_ohm", &config->current_shunt_resistance);
        parse_float_value(line, "max_current_rating_a", &config->current_max_rating);
        parse_int_value(line, "alarm_threshold_percent", &config->current_alarm_threshold_percent);
        
        parse_bool_value(line, "immediate_shutdown", &config->emergency_immediate_shutdown);
        parse_bool_value(line, "disable_tx_rf", &config->emergency_disable_tx_rf);
        parse_bool_value(line, "disable_pa_vds", &config->emergency_disable_pa_vds);
        parse_bool_value(line, "disable_28v_vds", &config->emergency_disable_28v_vds);
        parse_bool_value(line, "log_event", &config->emergency_log_event);
    }
    
    fclose(file);
    
    config->fem1_temp_table.num_points = fem1_points;
    config->fem2_temp_table.num_points = fem2_points;
    
    usys_log_info("YAML config loaded: FEM1 temp points=%d, FEM2 temp points=%d", fem1_points, fem2_points);
    
    return yaml_config_validate(config);
}

int yaml_config_get_dac_voltages_for_temp(const YamlSafetyConfig *config, FemUnit unit, 
                                          float temperature, float *carrier_voltage, float *peak_voltage) {
    if (!config || !carrier_voltage || !peak_voltage) {
        return STATUS_NOK;
    }
    
    const TempCompensationTable *table;
    if (unit == FEM_UNIT_1) {
        table = &config->fem1_temp_table;
    } else if (unit == FEM_UNIT_2) {
        table = &config->fem2_temp_table;
    } else {
        return STATUS_NOK;
    }
    
    if (table->num_points == 0) {
        *carrier_voltage = config->default_carrier_voltage;
        *peak_voltage = config->default_peak_voltage;
        return STATUS_OK;
    }
    
    int lower_idx = -1;
    int upper_idx = -1;
    
    for (int i = 0; i < table->num_points; i++) {
        if (table->points[i].temperature_c <= temperature) {
            lower_idx = i;
        }
        if (table->points[i].temperature_c >= temperature && upper_idx == -1) {
            upper_idx = i;
        }
    }
    
    if (lower_idx == -1) {
        *carrier_voltage = table->points[0].carrier_voltage;
        *peak_voltage = table->points[0].peak_voltage;
        return STATUS_OK;
    }
    
    if (upper_idx == -1) {
        *carrier_voltage = table->points[table->num_points - 1].carrier_voltage;
        *peak_voltage = table->points[table->num_points - 1].peak_voltage;
        return STATUS_OK;
    }
    
    if (lower_idx == upper_idx) {
        *carrier_voltage = table->points[lower_idx].carrier_voltage;
        *peak_voltage = table->points[lower_idx].peak_voltage;
        return STATUS_OK;
    }
    
    float temp_lower = table->points[lower_idx].temperature_c;
    float temp_upper = table->points[upper_idx].temperature_c;
    float temp_ratio = (temperature - temp_lower) / (temp_upper - temp_lower);
    
    *carrier_voltage = table->points[lower_idx].carrier_voltage + 
                      temp_ratio * (table->points[upper_idx].carrier_voltage - table->points[lower_idx].carrier_voltage);
    
    *peak_voltage = table->points[lower_idx].peak_voltage + 
                   temp_ratio * (table->points[upper_idx].peak_voltage - table->points[lower_idx].peak_voltage);
    
    return STATUS_OK;
}

void yaml_config_print(const YamlSafetyConfig *config) {
    if (!config) return;
    
    usys_log_info("=== YAML Safety Configuration ===");
    usys_log_info("Safety enabled: %s", config->enabled ? "true" : "false");
    usys_log_info("Check interval: %u ms", config->check_interval_ms);
    usys_log_info("Max violations: %u", config->max_violations);
    
    usys_log_info("Thresholds:");
    usys_log_info("  Max reverse power: %.1f dBm", config->max_reverse_power_dbm);
    usys_log_info("  Max PA current: %.1f A", config->max_pa_current_a);
    usys_log_info("  Max temperature: %.1f°C", config->max_temperature_c);
    usys_log_info("  Min temperature: %.1f°C", config->min_temperature_c);
    
    usys_log_info("DAC voltages:");
    usys_log_info("  Range: %.1f - %.1f V", config->dac_min_voltage, config->dac_max_voltage);
    usys_log_info("  Default carrier: %.2f V", config->default_carrier_voltage);
    usys_log_info("  Default peak: %.2f V", config->default_peak_voltage);
    
    usys_log_info("Temp compensation tables:");
    usys_log_info("  FEM1 points: %d", config->fem1_temp_table.num_points);
    usys_log_info("  FEM2 points: %d", config->fem2_temp_table.num_points);
}

int yaml_config_validate(const YamlSafetyConfig *config) {
    if (!config) return STATUS_NOK;
    
    if (config->max_temperature_c <= config->min_temperature_c) {
        usys_log_error("Invalid temperature range: max=%.1f <= min=%.1f", 
                       config->max_temperature_c, config->min_temperature_c);
        return STATUS_NOK;
    }
    
    if (config->dac_max_voltage <= config->dac_min_voltage) {
        usys_log_error("Invalid DAC voltage range: max=%.1f <= min=%.1f", 
                       config->dac_max_voltage, config->dac_min_voltage);
        return STATUS_NOK;
    }
    
    if (config->default_carrier_voltage < config->dac_min_voltage || 
        config->default_carrier_voltage > config->dac_max_voltage) {
        usys_log_error("Default carrier voltage %.2fV outside DAC range [%.1f, %.1f]", 
                       config->default_carrier_voltage, config->dac_min_voltage, config->dac_max_voltage);
        return STATUS_NOK;
    }
    
    if (config->default_peak_voltage < config->dac_min_voltage || 
        config->default_peak_voltage > config->dac_max_voltage) {
        usys_log_error("Default peak voltage %.2fV outside DAC range [%.1f, %.1f]", 
                       config->default_peak_voltage, config->dac_min_voltage, config->dac_max_voltage);
        return STATUS_NOK;
    }
    
    if (config->check_interval_ms < 100) {
        usys_log_error("Check interval %u ms too short (minimum 100ms)", config->check_interval_ms);
        return STATUS_NOK;
    }
    
    if (config->max_violations == 0) {
        usys_log_error("Max violations cannot be zero");
        return STATUS_NOK;
    }
    
    usys_log_info("YAML configuration validation passed");
    return STATUS_OK;
}