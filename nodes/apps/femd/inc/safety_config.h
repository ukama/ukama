/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SAFETY_CONFIG_H
#define SAFETY_CONFIG_H

#include <stdint.h>
#include <stdbool.h>

#include "gpio_controller.h" /* for FemUnit */

#define SAFETY_TEMP_TABLE_MAX_POINTS 16
#define SAFETY_SENSOR_TYPE_MAX_LEN   16

typedef struct {
    float carrier_voltage;
    float peak_voltage;
} SafetyDacVoltages;

/*
 * Temperature compensation table:
 * - num_points: number of valid points in arrays
 * - tempC[i] corresponds to carrierV[i] and peakV[i]
 *
 * NOTE: points should be sorted by tempC ascending for predictable lookup.
 */
typedef struct {
    int   num_points;
    float tempC[SAFETY_TEMP_TABLE_MAX_POINTS];
    float carrierV[SAFETY_TEMP_TABLE_MAX_POINTS];
    float peakV[SAFETY_TEMP_TABLE_MAX_POINTS];
} SafetyTempCompensationTable;

typedef struct {
    float max_reverse_power_dbm;
    float max_forward_power_dbm;
    float max_pa_current_a;
    float max_temperature_c;
    float min_temperature_c;
} SafetyThresholds;

typedef struct {
    float critical_high;
    float warning_high;
    float normal_high;
    float normal_low;
    float warning_low;
    float critical_low;
} SafetyTemperatureZones;

typedef struct {
    float min_voltage;
    float max_voltage;
    int   resolution_bits;

    float default_carrier_voltage;
    float default_peak_voltage;
    float shutdown_voltage;
    float standby_voltage;
} SafetyDacConfig;

typedef struct {
    uint32_t sampling_rate_hz;
    uint32_t averaging_samples;
    int      calibration_offset_mv;
} SafetyAdcConfig;

typedef struct {
    char     sensor_type[SAFETY_SENSOR_TYPE_MAX_LEN];
    uint8_t  i2c_address_fem1;
    uint8_t  i2c_address_fem2;
    int      resolution_bits;
    uint32_t update_interval_ms;
} SafetyTemperatureSensorConfig;

typedef struct {
    float shunt_resistance_ohm;
    float max_current_rating_a;
    int   alarm_threshold_percent;
} SafetyCurrentConfig;

/*
 * Emergency policy:
 * The YAML had multiple emergency_* blocks, but your old struct collapsed them.
 * We keep it simple for now: one policy applied to all emergency types.
 * (We can expand later if you want per-type policies.)
 */
typedef struct {
    bool immediate_shutdown;
    bool disable_tx_rf;
    bool disable_pa_vds;
    bool disable_28v_vds;
    bool reduce_dac_voltage;
    bool log_event;
} SafetyEmergencyPolicy;

/* Auto-restore behavior */
typedef struct {
    bool     enabled;
    uint32_t cooldown_ms;
    uint32_t ok_checks;
    bool     reset_unit_stats;
} SafetyAutoRestoreConfig;

typedef struct {
    /* core enable + cadence */
    bool     enabled;
    uint32_t check_interval_ms;
    uint32_t max_violations_before_shutdown;

    /* safety behavior */
    SafetyAutoRestoreConfig auto_restore;

    /* limits + zones */
    SafetyThresholds       thresholds;
    SafetyTemperatureZones temperature_zones;

    /* DAC configuration */
    SafetyDacConfig dac;

    /* temperature compensation tables (selected band resolved at load time) */
    SafetyTempCompensationTable fem1_temp_table;
    SafetyTempCompensationTable fem2_temp_table;

    /* monitoring */
    SafetyAdcConfig               adc;
    SafetyTemperatureSensorConfig temperature;
    SafetyCurrentConfig           current;

    /* emergency policy */
    SafetyEmergencyPolicy emergency;

    /* logging (kept minimal; if you need full logging config, add later) */
    struct {
        bool safety_events;
        bool temperature_logging;
        bool voltage_adjustments;
        bool current_measurements;
        char log_level[8];             /* "DEBUG","INFO","WARN","ERROR" */
        uint32_t max_log_file_size_mb;
        uint32_t log_rotation_count;
    } logging;
} SafetyConfig;

int  safety_config_load_json(const char *filename, SafetyConfig *cfg);
int  safety_config_validate(const SafetyConfig *cfg);
void safety_config_set_defaults(SafetyConfig *cfg);
void safety_config_print(const SafetyConfig *cfg);

/*
 * Lookup DAC voltages for temperature using the selected unit table.
 * Expected behavior:
 * - clamp outside range to nearest endpoint (or return error; define in C)
 * - linear interpolation between points (recommended)
 */
int  safety_config_get_dac_voltages_for_temp(const SafetyConfig *cfg,
                                            FemUnit unit,
                                            float temperature_c,
                                            float *carrier_voltage,
                                            float *peak_voltage);

#endif /* SAFETY_CONFIG_H */
