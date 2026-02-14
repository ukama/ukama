/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#include <jansson.h>

#include "femd.h"
#include "safety_config.h"
#include "usys_log.h"

static const char *env_band(void) {
    const char *b = getenv("ENV_FEM_BAND");
    if (!b || !*b) return "B41";
    return b;
}

static json_t *obj_get(json_t *obj, const char *key) {
    if (!obj || !json_is_object(obj) || !key) return NULL;
    return json_object_get(obj, key);
}

static void set_bool(json_t *obj, const char *key, bool *dst) {
    json_t *v;
    if (!dst) return;
    v = obj_get(obj, key);
    if (!v) return;
    if (json_is_true(v))  { *dst = true;  return; }
    if (json_is_false(v)) { *dst = false; return; }
}

static void set_u32(json_t *obj, const char *key, uint32_t *dst) {
    json_t *v;
    if (!dst) return;
    v = obj_get(obj, key);
    if (!v || !json_is_integer(v)) return;
    if (json_integer_value(v) < 0) return;
    *dst = (uint32_t)json_integer_value(v);
}

static void set_i32(json_t *obj, const char *key, int *dst) {
    json_t *v;
    if (!dst) return;
    v = obj_get(obj, key);
    if (!v || !json_is_integer(v)) return;
    *dst = (int)json_integer_value(v);
}

static void set_f32(json_t *obj, const char *key, float *dst) {
    json_t *v;
    if (!dst) return;
    v = obj_get(obj, key);
    if (!v) return;
    if (json_is_real(v)) {
        *dst = (float)json_real_value(v);
        return;
    }
    if (json_is_integer(v)) {
        *dst = (float)json_integer_value(v);
        return;
    }
}

static void set_str(json_t *obj, const char *key, char *dst, size_t dst_sz) {
    json_t *v;
    const char *s;
    if (!dst || dst_sz == 0) return;
    v = obj_get(obj, key);
    if (!v || !json_is_string(v)) return;
    s = json_string_value(v);
    if (!s) return;
    snprintf(dst, dst_sz, "%s", s);
}

/* Accept integer (preferred) OR string "0x48" (tolerant). */
static void set_u8_addr(json_t *obj, const char *key, uint8_t *dst) {
    json_t *v;
    if (!dst) return;
    v = obj_get(obj, key);
    if (!v) return;

    if (json_is_integer(v)) {
        json_int_t x = json_integer_value(v);
        if (x < 0 || x > 255) return;
        *dst = (uint8_t)x;
        return;
    }

    if (json_is_string(v)) {
        const char *s = json_string_value(v);
        unsigned long x;
        if (!s) return;
        x = strtoul(s, NULL, 0);
        if (x > 255) return;
        *dst = (uint8_t)x;
        return;
    }
}

typedef struct {
    float temp;
    float carrier;
    float peak;
} TempPoint;

static int cmp_temppoint(const void *a, const void *b) {
    const TempPoint *pa = (const TempPoint *)a;
    const TempPoint *pb = (const TempPoint *)b;
    if (pa->temp < pb->temp) return -1;
    if (pa->temp > pb->temp) return 1;
    return 0;
}

/*
 * voltage_lookup is an object:
 *   "0": { "carrier": 2.15, "peak": 1.25 },
 *   "10": { ... }
 */
static void parse_voltage_lookup(json_t *lookup, SafetyTempCompensationTable *out) {
    const char *k;
    json_t *v;

    TempPoint pts[SAFETY_TEMP_TABLE_MAX_POINTS];
    int n = 0;

    if (!out) return;
    out->num_points = 0;

    if (!lookup || !json_is_object(lookup)) return;

    json_object_foreach(lookup, k, v) {
        char *endp = NULL;
        float temp;
        json_t *cv, *pv;

        if (!k || !v || !json_is_object(v)) continue;
        if (n >= SAFETY_TEMP_TABLE_MAX_POINTS) break;

        temp = (float)strtod(k, &endp);
        if (endp == k) continue;

        cv = json_object_get(v, "carrier");
        pv = json_object_get(v, "peak");
        if (!cv || !pv) continue;

        if (!(json_is_real(cv) || json_is_integer(cv))) continue;
        if (!(json_is_real(pv) || json_is_integer(pv))) continue;

        pts[n].temp = temp;
        pts[n].carrier = json_is_real(cv) ? (float)json_real_value(cv) : (float)json_integer_value(cv);
        pts[n].peak   = json_is_real(pv) ? (float)json_real_value(pv) : (float)json_integer_value(pv);
        n++;
    }

    if (n <= 0) return;

    qsort(pts, (size_t)n, sizeof(pts[0]), cmp_temppoint);

    out->num_points = n;
    for (int i = 0; i < n; i++) {
        out->tempC[i]    = pts[i].temp;
        out->carrierV[i] = pts[i].carrier;
        out->peakV[i]    = pts[i].peak;
    }
}

void safety_config_set_defaults(SafetyConfig *c) {
    if (!c) return;
    memset(c, 0, sizeof(*c));

    c->enabled = true;
    c->check_interval_ms = 1000;
    c->max_violations_before_shutdown = 3;

    c->auto_restore.enabled          = true;
    c->auto_restore.cooldown_ms      = 30000;
    c->auto_restore.ok_checks        = 5;
    c->auto_restore.reset_unit_stats = true;

    c->thresholds.max_reverse_power_dbm = -10.0f;
    c->thresholds.max_forward_power_dbm = 30.0f;
    c->thresholds.max_pa_current_a      = 5.0f;
    c->thresholds.max_temperature_c     = 85.0f;
    c->thresholds.min_temperature_c     = -40.0f;

    c->temperature_zones.critical_high = 85.0f;
    c->temperature_zones.warning_high  = 75.0f;
    c->temperature_zones.normal_high   = 65.0f;
    c->temperature_zones.normal_low    = 0.0f;
    c->temperature_zones.warning_low   = -20.0f;
    c->temperature_zones.critical_low  = -40.0f;

    c->dac.min_voltage             = 0.0f;
    c->dac.max_voltage             = 2.5f;
    c->dac.resolution_bits         = 12;
    c->dac.default_carrier_voltage = 1.2f;
    c->dac.default_peak_voltage    = 2.0f;
    c->dac.shutdown_voltage        = 0.0f;
    c->dac.standby_voltage         = 0.5f;

    c->fem1_temp_table.num_points = 0;
    c->fem2_temp_table.num_points = 0;

    c->adc.sampling_rate_hz      = 1000;
    c->adc.averaging_samples     = 10;
    c->adc.calibration_offset_mv = 0;

    snprintf(c->temperature.sensor_type, sizeof(c->temperature.sensor_type), "%s", "LM75A");
    c->temperature.i2c_address_fem1   = 72;
    c->temperature.i2c_address_fem2   = 73;
    c->temperature.resolution_bits    = 12;
    c->temperature.update_interval_ms = 2000;

    c->current.shunt_resistance_ohm    = 0.01f;
    c->current.max_current_rating_a    = 10.0f;
    c->current.alarm_threshold_percent = 80;

    c->emergency.immediate_shutdown = true;
    c->emergency.disable_tx_rf      = true;
    c->emergency.disable_pa_vds     = true;
    c->emergency.disable_28v_vds    = true;
    c->emergency.reduce_dac_voltage = true;
    c->emergency.log_event          = true;

    c->logging.safety_events       = true;
    c->logging.temperature_logging = true;
    c->logging.voltage_adjustments = true;
    c->logging.current_measurements = true;
    snprintf(c->logging.log_level, sizeof(c->logging.log_level), "%s", "INFO");
    c->logging.max_log_file_size_mb = 100;
    c->logging.log_rotation_count = 5;
}

int safety_config_validate(const SafetyConfig *c) {
    if (!c) return STATUS_NOK;

    if (c->check_interval_ms == 0) return STATUS_NOK;
    if (c->auto_restore.ok_checks == 0) return STATUS_NOK;

    if (c->dac.max_voltage <= 0.0f) return STATUS_NOK;
    if (c->dac.min_voltage < 0.0f) return STATUS_NOK;
    if (c->dac.min_voltage > c->dac.max_voltage) return STATUS_NOK;

    return STATUS_OK;
}

static int interpolate_table(const SafetyTempCompensationTable *t,
                             float tempC,
                             float *cv,
                             float *pv) {
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

int safety_config_get_dac_voltages_for_temp(const SafetyConfig *c,
                                            FemUnit unit,
                                            float temperature_c,
                                            float *carrier_voltage,
                                            float *peak_voltage) {
    const SafetyTempCompensationTable *t;

    if (!c || !carrier_voltage || !peak_voltage) return STATUS_NOK;

    t = (unit == FEM_UNIT_1) ? &c->fem1_temp_table : &c->fem2_temp_table;

    if (t->num_points <= 0) {
        *carrier_voltage = c->dac.default_carrier_voltage;
        *peak_voltage    = c->dac.default_peak_voltage;
        return STATUS_OK;
    }

    return interpolate_table(t, temperature_c, carrier_voltage, peak_voltage);
}

static void apply_json(SafetyConfig *c, json_t *root) {
    json_t *safety, *th, *tz;
    json_t *dac, *vr, *dv;
    json_t *mon, *adc, *temp, *cur;
    json_t *emerg, *logging;
    json_t *tc, *bands, *band_obj, *fem1, *fem2, *vl;

    if (!c || !root || !json_is_object(root)) return;

    safety = obj_get(root, "safety");
    if (safety) {
        set_bool(safety, "enabled", &c->enabled);
        set_u32(safety, "check_interval_ms", &c->check_interval_ms);
        set_u32(safety, "max_violations_before_shutdown", &c->max_violations_before_shutdown);

        set_bool(safety, "auto_restore_enabled",     &c->auto_restore.enabled);
        set_u32(safety, "restore_cooldown_ms",       &c->auto_restore.cooldown_ms);
        set_u32(safety, "restore_ok_checks",         &c->auto_restore.ok_checks);
        set_bool(safety, "restore_reset_unit_stats", &c->auto_restore.reset_unit_stats);

        th = obj_get(safety, "thresholds");
        if (th) {
            set_f32(th, "max_reverse_power_dbm", &c->thresholds.max_reverse_power_dbm);
            set_f32(th, "max_forward_power_dbm", &c->thresholds.max_forward_power_dbm);
            set_f32(th, "max_pa_current_a",      &c->thresholds.max_pa_current_a);
            set_f32(th, "max_temperature_c",     &c->thresholds.max_temperature_c);
            set_f32(th, "min_temperature_c",     &c->thresholds.min_temperature_c);
        }

        tz = obj_get(safety, "temperature_zones");
        if (tz) {
            set_f32(tz, "critical_high", &c->temperature_zones.critical_high);
            set_f32(tz, "warning_high",  &c->temperature_zones.warning_high);
            set_f32(tz, "normal_high",   &c->temperature_zones.normal_high);
            set_f32(tz, "normal_low",    &c->temperature_zones.normal_low);
            set_f32(tz, "warning_low",   &c->temperature_zones.warning_low);
            set_f32(tz, "critical_low",  &c->temperature_zones.critical_low);
        }
    }

    dac = obj_get(root, "dac");
    if (dac) {
        vr = obj_get(dac, "voltage_range");
        if (vr) {
            set_f32(vr, "min_voltage",     &c->dac.min_voltage);
            set_f32(vr, "max_voltage",     &c->dac.max_voltage);
            set_i32(vr, "resolution_bits", &c->dac.resolution_bits);
        }

        dv = obj_get(dac, "default_voltages");
        if (dv) {
            set_f32(dv, "carrier_voltage",  &c->dac.default_carrier_voltage);
            set_f32(dv, "peak_voltage",     &c->dac.default_peak_voltage);
            set_f32(dv, "shutdown_voltage", &c->dac.shutdown_voltage);
            set_f32(dv, "standby_voltage",  &c->dac.standby_voltage);
        }
    }

    mon = obj_get(root, "monitoring");
    if (mon) {
        adc = obj_get(mon, "adc");
        if (adc) {
            set_u32(adc, "sampling_rate_hz",      &c->adc.sampling_rate_hz);
            set_u32(adc, "averaging_samples",     &c->adc.averaging_samples);
            set_i32(adc, "calibration_offset_mv", &c->adc.calibration_offset_mv);
        }

        temp = obj_get(mon, "temperature");
        if (temp) {
            set_str(temp, "sensor_type",
                    c->temperature.sensor_type,
                    sizeof(c->temperature.sensor_type));
            set_u8_addr(temp, "i2c_address_fem1", &c->temperature.i2c_address_fem1);
            set_u8_addr(temp, "i2c_address_fem2", &c->temperature.i2c_address_fem2);
            set_i32(temp, "resolution_bits",      &c->temperature.resolution_bits);
            set_u32(temp, "update_interval_ms",   &c->temperature.update_interval_ms);
        }

        cur = obj_get(mon, "current");
        if (cur) {
            set_f32(cur, "shunt_resistance_ohm",    &c->current.shunt_resistance_ohm);
            set_f32(cur, "max_current_rating_a",    &c->current.max_current_rating_a);
            set_i32(cur, "alarm_threshold_percent", &c->current.alarm_threshold_percent);
        }
    }

    /* Emergency: we have 3 blocks in JSON; keep one unified policy in config for now.
       We prefer overtemperature if present, else reverse_power, else overcurrent. */
    emerg = obj_get(root, "emergency");
    if (emerg) {
        json_t *e = obj_get(emerg, "overtemperature_emergency");
        if (!e) e = obj_get(emerg, "reverse_power_emergency");
        if (!e) e = obj_get(emerg, "overcurrent_emergency");

        if (e) {
            set_bool(e, "immediate_shutdown", &c->emergency.immediate_shutdown);
            set_bool(e, "disable_tx_rf",      &c->emergency.disable_tx_rf);
            set_bool(e, "disable_pa_vds",     &c->emergency.disable_pa_vds);
            set_bool(e, "disable_28v_vds",    &c->emergency.disable_28v_vds);
            set_bool(e, "reduce_dac_voltage", &c->emergency.reduce_dac_voltage);
            set_bool(e, "log_event",          &c->emergency.log_event);
        }
    }

    logging = obj_get(root, "logging");
    if (logging) {
        set_bool(logging, "safety_events",        &c->logging.safety_events);
        set_bool(logging, "temperature_logging",  &c->logging.temperature_logging);
        set_bool(logging, "voltage_adjustments",  &c->logging.voltage_adjustments);
        set_bool(logging, "current_measurements", &c->logging.current_measurements);
        set_str(logging, "log_level",             c->logging.log_level, sizeof(c->logging.log_level));
        set_u32(logging, "max_log_file_size_mb",  &c->logging.max_log_file_size_mb);
        set_u32(logging, "log_rotation_count",    &c->logging.log_rotation_count);
    }

    /* Temperature compensation by band */
    tc = obj_get(root, "temperature_compensation");
    bands = tc ? obj_get(tc, "bands") : NULL;
    if (bands && json_is_object(bands)) {
        const char *band = env_band();
        band_obj = json_object_get(bands, band);
        if (!band_obj) {
            usys_log_warn("safety config: band '%s' missing under temperature_compensation.bands",
                          band);
            return;
        }

        fem1 = obj_get(band_obj, "fem1");
        fem2 = obj_get(band_obj, "fem2");

        if (fem1) {
            vl = obj_get(fem1, "voltage_lookup");
            parse_voltage_lookup(vl, &c->fem1_temp_table);
        }
        if (fem2) {
            vl = obj_get(fem2, "voltage_lookup");
            parse_voltage_lookup(vl, &c->fem2_temp_table);
        }
    }
}

int safety_config_load_json(const char *filename, SafetyConfig *c) {
    json_t *root;
    json_error_t err;

    if (!c || !filename) return STATUS_NOK;

    safety_config_set_defaults(c);

    root = json_load_file(filename, 0, &err);
    if (!root) {
        usys_log_warn("safety json not found/invalid: %s (line %d: %s) (using defaults)",
                      filename, err.line, err.text);
        return STATUS_OK; /* match old behavior: defaults if file absent/invalid */
    }

    apply_json(c, root);
    json_decref(root);

    return safety_config_validate(c);
}

void safety_config_print(const SafetyConfig *c) {
    if (!c) return;

    usys_log_info("Safety: enabled=%s interval_ms=%u max_violations=%u",
                  c->enabled ? "true" : "false",
                  (unsigned)c->check_interval_ms,
                  (unsigned)c->max_violations_before_shutdown);

    usys_log_info("Thresholds: reverse=%.1f dBm, current=%.1f A, temp=%.1f C (min=%.1f C, fwd_max=%.1f dBm)",
                  c->thresholds.max_reverse_power_dbm,
                  c->thresholds.max_pa_current_a,
                  c->thresholds.max_temperature_c,
                  c->thresholds.min_temperature_c,
                  c->thresholds.max_forward_power_dbm);

    usys_log_info("DAC defaults: carrier=%.2f peak=%.2f shutdown=%.2f standby=%.2f",
                  c->dac.default_carrier_voltage,
                  c->dac.default_peak_voltage,
                  c->dac.shutdown_voltage,
                  c->dac.standby_voltage);

    usys_log_info("Auto-restore: enabled=%s cooldown_ms=%u ok_checks=%u reset_stats=%s",
                  c->auto_restore.enabled ? "true" : "false",
                  (unsigned)c->auto_restore.cooldown_ms,
                  (unsigned)c->auto_restore.ok_checks,
                  c->auto_restore.reset_unit_stats ? "true" : "false");

    usys_log_info("Temp tables: fem1_points=%d fem2_points=%d band=%s",
                  c->fem1_temp_table.num_points,
                  c->fem2_temp_table.num_points,
                  env_band());
}
