/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <math.h>

#include "yaml_config.h"
#include "femd.h"

/* ---- small helpers for naive YAML line scanning ---- */

static int parse_float_value(const char *line, const char *key, float *value) {
    char search_key[64];
    char *key_pos;
    char *value_start;

    if (!line || !key || !value) return STATUS_NOK;

    snprintf(search_key, sizeof(search_key), "%s:", key);
    key_pos = strstr((char*)line, search_key);
    if (!key_pos) return STATUS_NOK;

    value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;

    *value = (float)strtod(value_start, NULL);
    return STATUS_OK;
}

static int parse_int_value(const char *line, const char *key, int *value) {
    char search_key[64];
    char *key_pos;
    char *value_start;

    if (!line || !key || !value) return STATUS_NOK;

    snprintf(search_key, sizeof(search_key), "%s:", key);
    key_pos = strstr((char*)line, search_key);
    if (!key_pos) return STATUS_NOK;

    value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;

    *value = atoi(value_start);
    return STATUS_OK;
}

static int parse_uint32_dec_or_hex(const char *line, const char *key, uint32_t *value) {
    char search_key[64];
    char *key_pos;
    char *value_start;
    unsigned long v;

    if (!line || !key || !value) return STATUS_NOK;

    snprintf(search_key, sizeof(search_key), "%s:", key);
    key_pos = strstr((char*)line, search_key);
    if (!key_pos) return STATUS_NOK;

    value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;

    if (value_start[0] == '0' && (value_start[1] == 'x' || value_start[1] == 'X')) {
        v = strtoul(value_start + 2, NULL, 16);
    } else {
        v = strtoul(value_start, NULL, 10);
    }
    *value = (uint32_t)v;
    return STATUS_OK;
}

static int parse_uint32_value(const char *line, const char *key, uint32_t *value) {
    /* keep the original decimal-only behavior where needed */
    char search_key[64];
    char *key_pos;
    char *value_start;

    if (!line || !key || !value) return STATUS_NOK;

    snprintf(search_key, sizeof(search_key), "%s:", key);
    key_pos = strstr((char*)line, search_key);
    if (!key_pos) return STATUS_NOK;

    value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;

    *value = (uint32_t)strtoul(value_start, NULL, 10);
    return STATUS_OK;
}

static int parse_bool_value(const char *line, const char *key, bool *value) {
    char search_key[64];
    char *key_pos;
    char *value_start;

    if (!line || !key || !value) return STATUS_NOK;

    snprintf(search_key, sizeof(search_key), "%s:", key);
    key_pos = strstr((char*)line, search_key);
    if (!key_pos) return STATUS_NOK;

    value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;

    if (strncmp(value_start, "true", 4) == 0) {
        *value = true;
        return STATUS_OK;
    } else if (strncmp(value_start, "false", 5) == 0) {
        *value = false;
        return STATUS_OK;
    }
    return STATUS_NOK;
}

static int parse_string_value(const char *line, const char *key, char *value, size_t max_len) {
    char search_key[64];
    char *key_pos;
    char *value_start;
    char *end;
    size_t len;

    if (!line || !key || !value || max_len == 0U) return STATUS_NOK;

    snprintf(search_key, sizeof(search_key), "%s:", key);
    key_pos = strstr((char*)line, search_key);
    if (!key_pos) return STATUS_NOK;

    value_start = key_pos + strlen(search_key);
    while (*value_start == ' ' || *value_start == '\t') value_start++;

    if (*value_start == '"') {
        value_start++;
        end = strchr(value_start, '"');
        if (!end) return STATUS_NOK;
        len = (size_t)(end - value_start);
        if (len >= max_len) return STATUS_NOK;
        strncpy(value, value_start, len);
        value[len] = '\0';
        return STATUS_OK;
    }

    end = strchr(value_start, '\n');
    if (!end) end = strchr(value_start, '#');
    if (!end) end = (char*)value_start + strlen(value_start);

    len = (size_t)(end - value_start);
    while (len > 0 &&
           (value_start[len - 1] == ' ' ||
            value_start[len - 1] == '\t' ||
            value_start[len - 1] == '\r')) {
        len--;
    }

    if (len >= max_len) return STATUS_NOK;
    strncpy(value, value_start, len);
    value[len] = '\0';
    return STATUS_OK;
}

static int parse_temp_voltage_line(const char *line, TempVoltagePoint *point) {
    char *colon;
    char *open_brace;
    char *close_brace;
    char temp_str[16];
    const char *temp_start;
    size_t temp_len;
    char *carrier_pos;
    char *peak_pos;

    if (!line || !point) return STATUS_NOK;

    colon = strchr(line, ':');
    if (!colon) return STATUS_NOK;

    temp_start = line;
    while (*temp_start == ' ' || *temp_start == '\t') temp_start++;

    temp_len = (size_t)(colon - temp_start);
    if (temp_len == 0U || temp_len >= sizeof(temp_str)) return STATUS_NOK;

    strncpy(temp_str, temp_start, temp_len);
    temp_str[temp_len] = '\0';
    point->temperature_c = (float)strtod(temp_str, NULL);

    open_brace  = strchr(colon, '{');
    close_brace = strchr(colon, '}');
    if (!open_brace || !close_brace || close_brace <= open_brace) return STATUS_NOK;

    carrier_pos = strstr(open_brace, "carrier:");
    if (!carrier_pos || carrier_pos > close_brace) return STATUS_NOK;
    carrier_pos += 8;
    while (*carrier_pos == ' ' || *carrier_pos == '\t') carrier_pos++;
    point->carrier_voltage = (float)strtod(carrier_pos, NULL);

    peak_pos = strstr(open_brace, "peak:");
    if (!peak_pos || peak_pos > close_brace) return STATUS_NOK;
    peak_pos += 5;
    while (*peak_pos == ' ' || *peak_pos == '\t') peak_pos++;
    point->peak_voltage = (float)strtod(peak_pos, NULL);

    return STATUS_OK;
}

/* ---- defaults, validate, print ---- */

void yaml_config_set_defaults(YamlSafetyConfig *config) {
    if (!config) return;

    memset(config, 0, sizeof(YamlSafetyConfig));

    config->enabled              = true;
    config->check_interval_ms    = 1000;
    config->max_violations       = 3;

    config->max_reverse_power_dbm = -10.0f;
    config->max_pa_current_a      = 5.0f;
    config->max_temperature_c     = 85.0f;
    config->min_temperature_c     = -40.0f;
    config->max_forward_power_dbm = 30.0f;

    config->temp_critical_high = 85.0f;
    config->temp_warning_high  = 75.0f;
    config->temp_normal_high   = 65.0f;
    config->temp_normal_low    = 0.0f;
    config->temp_warning_low   = -20.0f;
    config->temp_critical_low  = -40.0f;

    config->dac_min_voltage     = 0.0f;
    config->dac_max_voltage     = 2.5f;
    config->dac_resolution_bits = 12;

    config->default_carrier_voltage = 1.2f;
    config->default_peak_voltage    = 2.0f;
    config->shutdown_voltage        = 0.0f;
    config->standby_voltage         = 0.5f;

    config->adc_sampling_rate_hz    = 1000;
    config->adc_averaging_samples   = 10;
    config->adc_calibration_offset_mv = 0;

    strcpy(config->temp_sensor_type, "LM75A");
    config->temp_i2c_addr_fem1      = 0x48;
    config->temp_i2c_addr_fem2      = 0x49;
    config->temp_resolution_bits    = 12;
    config->temp_update_interval_ms = 2000;

    config->current_shunt_resistance   = 0.01f;
    config->current_max_rating         = 10.0f;
    config->current_alarm_threshold_percent = 80;

    /* Global emergency flags; parsed later as well */
    config->emergency_immediate_shutdown = true;
    config->emergency_disable_tx_rf      = true;
    config->emergency_disable_pa_vds     = true;
    config->emergency_disable_28v_vds    = true;
    config->emergency_log_event          = true;

    config->fem1_temp_table.num_points = 0;
    config->fem2_temp_table.num_points = 0;
}

static int yaml_clamp_validate(YamlSafetyConfig *c) {
    /* Basic sanity; do not be strict – keep current behavior */
    if (c->check_interval_ms < 100U)
        c->check_interval_ms = 100U;
    if (c->max_violations == 0U)
        c->max_violations    = 1U;
    if (c->dac_min_voltage < 0.0f)
        c->dac_min_voltage = 0.0f;
    if (c->dac_max_voltage < c->dac_min_voltage)
        c->dac_max_voltage = c->dac_min_voltage;
    if (c->default_carrier_voltage < c->dac_min_voltage)
        c->default_carrier_voltage = c->dac_min_voltage;
    if (c->default_carrier_voltage > c->dac_max_voltage)
        c->default_carrier_voltage = c->dac_max_voltage;
    if (c->default_peak_voltage < c->dac_min_voltage)
        c->default_peak_voltage = c->dac_min_voltage;
    if (c->default_peak_voltage > c->dac_max_voltage)
        c->default_peak_voltage = c->dac_max_voltage;
    if (c->temp_resolution_bits <= 0)
        c->temp_resolution_bits = 12;
    if (c->adc_averaging_samples == 0U)
        c->adc_averaging_samples = 1U;

    return STATUS_OK;
}

int yaml_config_validate(const YamlSafetyConfig *config) {
    /* Legacy signature in header expects an int; we keep simple OK */
    (void)config;
    return STATUS_OK;
}

int yaml_config_load(const char *filename, YamlSafetyConfig *config) {
    FILE *file;
    char line[512];
    int in_fem1_section = 0;
    int in_fem2_section = 0;
    int fem1_points = 0;
    int fem2_points = 0;

    if (!filename || !config) {
        return STATUS_NOK;
    }

    yaml_config_set_defaults(config);

    file = fopen(filename, "r");
    if (!file) {
        usys_log_warn("Could not open YAML config file %s, using defaults", filename);
        return STATUS_NOK;
    }

    usys_log_info("Loading YAML configuration from %s", filename);

    while (fgets(line, sizeof(line), file)) {
        /* Skip comments/blank lines early */
        if (line[0] == '#' || line[0] == '\n' || line[0] == '\r') continue;

        /* FEM temperature tables */
        if (strstr(line, "fem1:")) {
            in_fem1_section = 1;
            in_fem2_section = 0;
            continue;
        } else if (strstr(line, "fem2:")) {
            in_fem1_section = 0;
            in_fem2_section = 1;
            continue;
        } else if (strstr(line, "voltage_lookup:")) {
            /* stay in current fem section */
            continue;
        } else if (line[0] != ' ' && line[0] != '\t') {
            /* new top-level block; clear fem section flags */
            in_fem1_section = 0;
            in_fem2_section = 0;
        }

        if (in_fem1_section && fem1_points < MAX_TEMP_POINTS) {
            TempVoltagePoint p1;
            if (parse_temp_voltage_line(line, &p1) == STATUS_OK) {
                config->fem1_temp_table.points[fem1_points++] = p1;
            }
        } else if (in_fem2_section && fem2_points < MAX_TEMP_POINTS) {
            TempVoltagePoint p2;
            if (parse_temp_voltage_line(line, &p2) == STATUS_OK) {
                config->fem2_temp_table.points[fem2_points++] = p2;
            }
        }

        /* Global safety */
        (void)parse_bool_value(line,    "enabled", &config->enabled);
        (void)parse_uint32_value(line,  "check_interval_ms", &config->check_interval_ms);
        (void)parse_uint32_value(line,  "max_violations_before_shutdown", &config->max_violations);

        /* Thresholds */
        (void)parse_float_value(line, "max_reverse_power_dbm", &config->max_reverse_power_dbm);
        (void)parse_float_value(line, "max_pa_current_a",      &config->max_pa_current_a);
        (void)parse_float_value(line, "max_temperature_c",     &config->max_temperature_c);
        (void)parse_float_value(line, "min_temperature_c",     &config->min_temperature_c);
        (void)parse_float_value(line, "max_forward_power_dbm", &config->max_forward_power_dbm);

        /* Temperature zones */
        (void)parse_float_value(line, "critical_high", &config->temp_critical_high);
        (void)parse_float_value(line, "warning_high",  &config->temp_warning_high);
        (void)parse_float_value(line, "normal_high",   &config->temp_normal_high);
        (void)parse_float_value(line, "normal_low",    &config->temp_normal_low);
        (void)parse_float_value(line, "warning_low",   &config->temp_warning_low);
        (void)parse_float_value(line, "critical_low",  &config->temp_critical_low);

        /* DAC block */
        (void)parse_float_value(line, "min_voltage",      &config->dac_min_voltage);
        (void)parse_float_value(line, "max_voltage",      &config->dac_max_voltage);
        (void)parse_int_value  (line, "resolution_bits",  &config->dac_resolution_bits);

        /* Defaults/operating modes */
        (void)parse_float_value(line, "carrier_voltage",  &config->default_carrier_voltage);
        (void)parse_float_value(line, "peak_voltage",     &config->default_peak_voltage);
        (void)parse_float_value(line, "shutdown_voltage", &config->shutdown_voltage);
        (void)parse_float_value(line, "standby_voltage",  &config->standby_voltage);

        /* ADC monitoring */
        (void)parse_uint32_value(line, "sampling_rate_hz",     &config->adc_sampling_rate_hz);
        (void)parse_uint32_value(line, "averaging_samples",    &config->adc_averaging_samples);
        (void)parse_int_value  (line,  "calibration_offset_mv",&config->adc_calibration_offset_mv);

        /* Temperature monitoring */
        (void)parse_string_value(line, "sensor_type",          config->temp_sensor_type,
                                 sizeof(config->temp_sensor_type));
        (void)parse_uint32_value(line, "update_interval_ms",   &config->temp_update_interval_ms);

        /* If present in file (hex or dec), honor these i2c addresses */
        {
            uint32_t addr;
            if (parse_uint32_dec_or_hex(line, "i2c_address_fem1", &addr) == STATUS_OK) {
                config->temp_i2c_addr_fem1 = (uint8_t)(addr & 0xFFu);
            }
            if (parse_uint32_dec_or_hex(line, "i2c_address_fem2", &addr) == STATUS_OK) {
                config->temp_i2c_addr_fem2 = (uint8_t)(addr & 0xFFu);
            }
        }

        /* Current monitoring */
        (void)parse_float_value(line, "shunt_resistance_ohm", &config->current_shunt_resistance);
        (void)parse_float_value(line, "max_current_rating_a", &config->current_max_rating);
        (void)parse_int_value  (line, "alarm_threshold_percent", &config->current_alarm_threshold_percent);

        /* Emergency flags:
         * The simple parser treats any occurrence as global. Your sample sets all to true. */
        (void)parse_bool_value(line, "immediate_shutdown", &config->emergency_immediate_shutdown);
        (void)parse_bool_value(line, "disable_tx_rf",      &config->emergency_disable_tx_rf);
        (void)parse_bool_value(line, "disable_pa_vds",     &config->emergency_disable_pa_vds);
        (void)parse_bool_value(line, "disable_28v_vds",    &config->emergency_disable_28v_vds);
        (void)parse_bool_value(line, "log_event",          &config->emergency_log_event);
    }

    fclose(file);

    config->fem1_temp_table.num_points = fem1_points;
    config->fem2_temp_table.num_points = fem2_points;

    (void)yaml_clamp_validate(config);

    usys_log_info("YAML config loaded: FEM1 temp points=%d, FEM2 temp points=%d",
                  fem1_points, fem2_points);

    return yaml_config_validate(config);
}

int yaml_config_get_dac_voltages_for_temp(const YamlSafetyConfig *config,
                                          FemUnit unit,
                                          float temperature,
                                          float *carrier_voltage,
                                          float *peak_voltage) {
    const TempCompensationTable *table;
    int lower_idx, upper_idx;
    float temp_lower, temp_upper, temp_ratio;

    if (!config || !carrier_voltage || !peak_voltage) {
        return STATUS_NOK;
    }

    if (unit == FEM_UNIT_1) {
        table = &config->fem1_temp_table;
    } else if (unit == FEM_UNIT_2) {
        table = &config->fem2_temp_table;
    } else {
        return STATUS_NOK;
    }

    if (table->num_points == 0) {
        *carrier_voltage = config->default_carrier_voltage;
        *peak_voltage    = config->default_peak_voltage;
        return STATUS_OK;
    }

    lower_idx = -1;
    upper_idx = -1;

    {
        int i;
        for (i = 0; i < table->num_points; i++) {
            if (table->points[i].temperature_c <= temperature) {
                lower_idx = i;
            }
            if (table->points[i].temperature_c >= temperature && upper_idx == -1) {
                upper_idx = i;
            }
        }
    }

    if (lower_idx == -1) {
        *carrier_voltage = table->points[0].carrier_voltage;
        *peak_voltage    = table->points[0].peak_voltage;
        return STATUS_OK;
    }

    if (upper_idx == -1) {
        *carrier_voltage = table->points[table->num_points - 1].carrier_voltage;
        *peak_voltage    = table->points[table->num_points - 1].peak_voltage;
        return STATUS_OK;
    }

    if (lower_idx == upper_idx) {
        *carrier_voltage = table->points[lower_idx].carrier_voltage;
        *peak_voltage    = table->points[lower_idx].peak_voltage;
        return STATUS_OK;
    }

    temp_lower = table->points[lower_idx].temperature_c;
    temp_upper = table->points[upper_idx].temperature_c;
    if (temp_upper <= temp_lower) {
        *carrier_voltage = table->points[lower_idx].carrier_voltage;
        *peak_voltage    = table->points[lower_idx].peak_voltage;
        return STATUS_OK;
    }

    temp_ratio = (temperature - temp_lower) / (temp_upper - temp_lower);

    *carrier_voltage = table->points[lower_idx].carrier_voltage +
                       temp_ratio * (table->points[upper_idx].carrier_voltage -
                                     table->points[lower_idx].carrier_voltage);

    *peak_voltage = table->points[lower_idx].peak_voltage +
                    temp_ratio * (table->points[upper_idx].peak_voltage -
                                  table->points[lower_idx].peak_voltage);

    return STATUS_OK;
}

void yaml_config_print(const YamlSafetyConfig *c) {
    if (!c) return;

    usys_log_info("=== YAML Safety Configuration ===");
    usys_log_info("Safety: enabled=%s, interval=%u ms, max_violations=%u",
                  c->enabled ? "true" : "false",
                  (unsigned int)c->check_interval_ms,
                  (unsigned int)c->max_violations);

    usys_log_info("Thresholds: reverse=%.1f dBm, current=%.1f A, temp=%.1f C (min=%.1f C, fwd_max=%.1f dBm)",
                  c->max_reverse_power_dbm,
                  c->max_pa_current_a,
                  c->max_temperature_c,
                  c->min_temperature_c,
                  c->max_forward_power_dbm);

    usys_log_info("Temp zones: crit_hi=%.1f, warn_hi=%.1f, norm_hi=%.1f, norm_lo=%.1f, warn_lo=%.1f, crit_lo=%.1f",
                  c->temp_critical_high, c->temp_warning_high, c->temp_normal_high,
                  c->temp_normal_low, c->temp_warning_low, c->temp_critical_low);

    usys_log_info("DAC: vmin=%.2f V, vmax=%.2f V, res=%d bits; defaults: carrier=%.2f, peak=%.2f, shutdown=%.2f, standby=%.2f",
                  c->dac_min_voltage, c->dac_max_voltage, c->dac_resolution_bits,
                  c->default_carrier_voltage, c->default_peak_voltage,
                  c->shutdown_voltage, c->standby_voltage);

    usys_log_info("ADC: fs=%u Hz, avg=%u, cal_off=%d mV",
                  (unsigned int)c->adc_sampling_rate_hz,
                  (unsigned int)c->adc_averaging_samples,
                  c->adc_calibration_offset_mv);

    usys_log_info("Temp sensor: type=%s, addr_fem1=0x%02X, addr_fem2=0x%02X, res=%d bits, update=%u ms",
                  c->temp_sensor_type,
                  (unsigned int)c->temp_i2c_addr_fem1,
                  (unsigned int)c->temp_i2c_addr_fem2,
                  c->temp_resolution_bits,
                  (unsigned int)c->temp_update_interval_ms);

    usys_log_info("Current: Rsh=%.4f ohm, Imax=%.1f A, alarm=%d%%",
                  c->current_shunt_resistance,
                  c->current_max_rating,
                  c->current_alarm_threshold_percent);

    usys_log_info("Emergency: immediate=%s, cut_tx=%s, cut_pa=%s, cut_28v=%s, log=%s",
                  c->emergency_immediate_shutdown ? "true" : "false",
                  c->emergency_disable_tx_rf      ? "true" : "false",
                  c->emergency_disable_pa_vds     ? "true" : "false",
                  c->emergency_disable_28v_vds    ? "true" : "false",
                  c->emergency_log_event          ? "true" : "false");

    usys_log_info("Tables: fem1_points=%d, fem2_points=%d",
                  c->fem1_temp_table.num_points, c->fem2_temp_table.num_points);
}
