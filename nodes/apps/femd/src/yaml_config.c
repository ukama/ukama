/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#include <yaml.h>

#include "yaml_config.h"
#include "usys_log.h"

static int parse_bool(const char *s, bool *out) {

    if (!s || !out) return STATUS_NOK;

    if (!strcasecmp(s, "true") || !strcasecmp(s, "yes") || !strcmp(s, "1")) { *out = true; return STATUS_OK; }
    if (!strcasecmp(s, "false") || !strcasecmp(s, "no")  || !strcmp(s, "0")) { *out = false; return STATUS_OK; }
    return STATUS_NOK;
}

static int parse_u32(const char *s, uint32_t *out) {
    if (!s || !out) return STATUS_NOK;
    *out = (uint32_t)strtoul(s, NULL, 10);
    return STATUS_OK;
}

static int parse_i32(const char *s, int *out) {
    if (!s || !out) return STATUS_NOK;
    *out = (int)strtol(s, NULL, 10);
    return STATUS_OK;
}

static int parse_f32(const char *s, float *out) {
    if (!s || !out) return STATUS_NOK;
    *out = (float)strtod(s, NULL);
    return STATUS_OK;
}

void yaml_config_set_defaults(YamlSafetyConfig *c) {

    if (!c) return;

    memset(c, 0, sizeof(*c));

    c->enabled = true;
    c->check_interval_ms = 500;
    c->max_violations = 3;

    c->max_reverse_power_dbm = 10.0f;
    c->max_forward_power_dbm = 45.0f;
    c->max_pa_current_a = 10.0f;
    c->max_temperature_c = 85.0f;
    c->min_temperature_c = -20.0f;

    c->temp_critical_high = 90.0f;
    c->temp_warning_high  = 85.0f;
    c->temp_normal_high   = 80.0f;
    c->temp_normal_low    = 0.0f;
    c->temp_warning_low   = -10.0f;
    c->temp_critical_low  = -20.0f;

    c->dac_min_voltage = 0.0f;
    c->dac_max_voltage = 5.0f;
    c->dac_resolution_bits = 16;

    c->default_carrier_voltage = 1.0f;
    c->default_peak_voltage    = 1.0f;
    c->shutdown_voltage        = 0.0f;
    c->standby_voltage         = 0.2f;

    c->fem1_temp_table.num_points = 0;
    c->fem2_temp_table.num_points = 0;

    c->adc_sampling_rate_hz = 1600;
    c->adc_averaging_samples = 1;
    c->adc_calibration_offset_mv = 0;

    snprintf(c->temp_sensor_type, sizeof(c->temp_sensor_type), "%s", "LM75A");
    c->temp_i2c_addr_fem1 = 0x49;
    c->temp_i2c_addr_fem2 = 0x49;
    c->temp_resolution_bits = 9;
    c->temp_update_interval_ms = 2000;

    c->current_shunt_resistance = 0.01f;
    c->current_max_rating = 10.0f;
    c->current_alarm_threshold_percent = 80;

    c->emergency_immediate_shutdown = true;
    c->emergency_disable_tx_rf = true;
    c->emergency_disable_pa_vds = true;
    c->emergency_disable_28v_vds = true;
    c->emergency_log_event = true;

    c->auto_restore_enabled = true;
    c->restore_cooldown_ms = 30000;
    c->restore_ok_checks = 5;
    c->restore_reset_unit_stats = false;
}

static int apply_key_value(YamlSafetyConfig *c, const char *path, const char *val) {

    if (!c || !path || !val) return STATUS_NOK;

    if (!strcmp(path, "safety.enabled")) return parse_bool(val, &c->enabled);
    if (!strcmp(path, "safety.check_interval_ms")) return parse_u32(val, &c->check_interval_ms);
    if (!strcmp(path, "safety.max_violations")) return parse_u32(val, &c->max_violations);

    if (!strcmp(path, "thresholds.max_reverse_power_dbm")) return parse_f32(val, &c->max_reverse_power_dbm);
    if (!strcmp(path, "thresholds.max_forward_power_dbm")) return parse_f32(val, &c->max_forward_power_dbm);
    if (!strcmp(path, "thresholds.max_pa_current_a")) return parse_f32(val, &c->max_pa_current_a);
    if (!strcmp(path, "thresholds.max_temperature_c")) return parse_f32(val, &c->max_temperature_c);
    if (!strcmp(path, "thresholds.min_temperature_c")) return parse_f32(val, &c->min_temperature_c);

    if (!strcmp(path, "temp_zones.critical_high")) return parse_f32(val, &c->temp_critical_high);
    if (!strcmp(path, "temp_zones.warning_high")) return parse_f32(val, &c->temp_warning_high);
    if (!strcmp(path, "temp_zones.normal_high")) return parse_f32(val, &c->temp_normal_high);
    if (!strcmp(path, "temp_zones.normal_low")) return parse_f32(val, &c->temp_normal_low);
    if (!strcmp(path, "temp_zones.warning_low")) return parse_f32(val, &c->temp_warning_low);
    if (!strcmp(path, "temp_zones.critical_low")) return parse_f32(val, &c->temp_critical_low);

    if (!strcmp(path, "dac.vmin")) return parse_f32(val, &c->dac_min_voltage);
    if (!strcmp(path, "dac.vmax")) return parse_f32(val, &c->dac_max_voltage);
    if (!strcmp(path, "dac.res_bits")) return parse_i32(val, &c->dac_resolution_bits);
    if (!strcmp(path, "dac.defaults.carrier_voltage")) return parse_f32(val, &c->default_carrier_voltage);
    if (!strcmp(path, "dac.defaults.peak_voltage")) return parse_f32(val, &c->default_peak_voltage);
    if (!strcmp(path, "dac.defaults.shutdown_voltage")) return parse_f32(val, &c->shutdown_voltage);
    if (!strcmp(path, "dac.defaults.standby_voltage")) return parse_f32(val, &c->standby_voltage);

    if (!strcmp(path, "adc.sampling_rate_hz")) return parse_u32(val, &c->adc_sampling_rate_hz);
    if (!strcmp(path, "adc.averaging_samples")) return parse_u32(val, &c->adc_averaging_samples);
    if (!strcmp(path, "adc.calibration_offset_mv")) return parse_i32(val, &c->adc_calibration_offset_mv);

    if (!strcmp(path, "temp_sensor.type")) { snprintf(c->temp_sensor_type, sizeof(c->temp_sensor_type), "%s", val); return STATUS_OK; }
    if (!strcmp(path, "temp_sensor.addr_fem1")) { c->temp_i2c_addr_fem1 = (uint8_t)strtoul(val, NULL, 0); return STATUS_OK; }
    if (!strcmp(path, "temp_sensor.addr_fem2")) { c->temp_i2c_addr_fem2 = (uint8_t)strtoul(val, NULL, 0); return STATUS_OK; }
    if (!strcmp(path, "temp_sensor.res_bits")) return parse_i32(val, &c->temp_resolution_bits);
    if (!strcmp(path, "temp_sensor.update_interval_ms")) return parse_u32(val, &c->temp_update_interval_ms);

    if (!strcmp(path, "current.shunt_resistance_ohm")) return parse_f32(val, &c->current_shunt_resistance);
    if (!strcmp(path, "current.max_current_rating_a")) return parse_f32(val, &c->current_max_rating);
    if (!strcmp(path, "current.alarm_threshold_percent")) return parse_i32(val, &c->current_alarm_threshold_percent);

    if (!strcmp(path, "emergency.immediate_shutdown")) return parse_bool(val, &c->emergency_immediate_shutdown);
    if (!strcmp(path, "emergency.disable_tx_rf")) return parse_bool(val, &c->emergency_disable_tx_rf);
    if (!strcmp(path, "emergency.disable_pa_vds")) return parse_bool(val, &c->emergency_disable_pa_vds);
    if (!strcmp(path, "emergency.disable_28v_vds")) return parse_bool(val, &c->emergency_disable_28v_vds);
    if (!strcmp(path, "emergency.log_event")) return parse_bool(val, &c->emergency_log_event);

    if (!strcmp(path, "auto_restore.enabled")) return parse_bool(val, &c->auto_restore_enabled);
    if (!strcmp(path, "auto_restore.cooldown_ms")) return parse_u32(val, &c->restore_cooldown_ms);
    if (!strcmp(path, "auto_restore.ok_checks")) return parse_u32(val, &c->restore_ok_checks);
    if (!strcmp(path, "auto_restore.reset_stats")) return parse_bool(val, &c->restore_reset_unit_stats);

    return STATUS_OK;
}

int yaml_config_load(const char *filename, YamlSafetyConfig *c) {

    FILE *fh;
    yaml_parser_t parser;
    yaml_event_t event;

    char stackKeys[16][64];
    int depth = 0;

    char curKey[64] = {0};
    bool haveKey = false;

    if (!c || !filename) return STATUS_NOK;

    yaml_config_set_defaults(c);

    fh = fopen(filename, "rb");
    if (!fh) {
        usys_log_warn("safety yaml not found: %s (using defaults)", filename);
        return STATUS_OK;
    }

    if (!yaml_parser_initialize(&parser)) {
        fclose(fh);
        return STATUS_NOK;
    }

    yaml_parser_set_input_file(&parser, fh);

    while (yaml_parser_parse(&parser, &event)) {

        if (event.type == YAML_MAPPING_START_EVENT) {
            if (haveKey) {
                if (depth < 16) {
                    snprintf(stackKeys[depth], sizeof(stackKeys[depth]), "%s", curKey);
                    depth++;
                }
                haveKey = false;
                curKey[0] = '\0';
            }
        } else if (event.type == YAML_MAPPING_END_EVENT) {
            if (depth > 0) depth--;
        } else if (event.type == YAML_SCALAR_EVENT) {

            const char *s = (const char *)event.data.scalar.value;

            if (!haveKey) {
                snprintf(curKey, sizeof(curKey), "%s", s);
                haveKey = true;
            } else {
                char path[256] = {0};
                int i;
                size_t off = 0;

                for (i = 0; i < depth; i++) {
                    off += (size_t)snprintf(path + off, sizeof(path) - off, "%s%s", i ? "." : "", stackKeys[i]);
                    if (off >= sizeof(path)) break;
                }
                if (off < sizeof(path)) {
                    (void)snprintf(path + off, sizeof(path) - off, "%s%s", depth ? "." : "", curKey);
                    (void)apply_key_value(c, path, s);
                }

                haveKey = false;
                curKey[0] = '\0';
            }
        }

        if (event.type == YAML_STREAM_END_EVENT) {
            yaml_event_delete(&event);
            break;
        }

        yaml_event_delete(&event);
    }

    yaml_parser_delete(&parser);
    fclose(fh);

    return yaml_config_validate(c);
}

int yaml_config_validate(const YamlSafetyConfig *c) {

    if (!c) return STATUS_NOK;

    if (c->restore_ok_checks == 0) return STATUS_NOK;
    if (c->check_interval_ms == 0) return STATUS_NOK;

    if (c->dac_max_voltage <= 0.0f) return STATUS_NOK;
    if (c->dac_min_voltage < 0.0f) return STATUS_NOK;
    if (c->dac_min_voltage > c->dac_max_voltage) return STATUS_NOK;

    return STATUS_OK;
}

static int interpolate_table(const TempCompensationTable *t, float tempC, float *cv, float *pv) {

    int i;

    if (!t || t->num_points <= 0 || !cv || !pv) return STATUS_NOK;

    if (tempC <= t->tempC[0]) {
        *cv = t->carrierV[0];
        *pv = t->peakV[0];
        return STATUS_OK;
    }

    if (tempC >= t->tempC[t->num_points - 1]) {
        *cv = t->carrierV[t->num_points - 1];
        *pv = t->peakV[t->num_points - 1];
        return STATUS_OK;
    }

    for (i = 0; i < t->num_points - 1; i++) {
        float t0 = t->tempC[i];
        float t1 = t->tempC[i + 1];
        if (tempC >= t0 && tempC <= t1) {
            float a = (tempC - t0) / (t1 - t0);
            *cv = t->carrierV[i] + a * (t->carrierV[i + 1] - t->carrierV[i]);
            *pv = t->peakV[i] + a * (t->peakV[i + 1] - t->peakV[i]);
            return STATUS_OK;
        }
    }

    return STATUS_NOK;
}

int yaml_config_get_dac_voltages_for_temp(const YamlSafetyConfig *c,
                                          FemUnit unit,
                                          float temperature,
                                          float *carrier_voltage,
                                          float *peak_voltage) {

    const TempCompensationTable *t;

    if (!c || !carrier_voltage || !peak_voltage) return STATUS_NOK;

    t = (unit == FEM_UNIT_1) ? &c->fem1_temp_table : &c->fem2_temp_table;

    if (t->num_points <= 0) {
        *carrier_voltage = c->default_carrier_voltage;
        *peak_voltage    = c->default_peak_voltage;
        return STATUS_OK;
    }

    return interpolate_table(t, temperature, carrier_voltage, peak_voltage);
}

void yaml_config_print(const YamlSafetyConfig *c) {

    if (!c) return;

    usys_log_info("Safety: enabled=%s interval_ms=%u max_violations=%u",
                  c->enabled ? "true" : "false",
                  (unsigned)c->check_interval_ms,
                  (unsigned)c->max_violations);

    usys_log_info("Thresholds: reverse=%.1f dBm, current=%.1f A, temp=%.1f C (min=%.1f C, fwd_max=%.1f dBm)",
                  c->max_reverse_power_dbm,
                  c->max_pa_current_a,
                  c->max_temperature_c,
                  c->min_temperature_c,
                  c->max_forward_power_dbm);

    usys_log_info("DAC defaults: carrier=%.2f peak=%.2f shutdown=%.2f standby=%.2f",
                  c->default_carrier_voltage,
                  c->default_peak_voltage,
                  c->shutdown_voltage,
                  c->standby_voltage);

    usys_log_info("Auto-restore: enabled=%s cooldown_ms=%u ok_checks=%u reset_stats=%s",
                  c->auto_restore_enabled ? "true":"false",
                  (unsigned)c->restore_cooldown_ms,
                  (unsigned)c->restore_ok_checks,
                  c->restore_reset_unit_stats ? "true":"false");
}
