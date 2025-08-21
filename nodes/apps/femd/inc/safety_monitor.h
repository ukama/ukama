/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */
#ifndef SAFETY_MONITOR_H
#define SAFETY_MONITOR_H

#include <stdint.h>
#include <stdbool.h>
#include <pthread.h>

#include "gpio_controller.h"
#include "i2c_controller.h"
#include "yaml_config.h"

#define SAFETY_MONITOR_INTERVAL_MS    1000  // Check every 1 second
#define SAFETY_SHUTDOWN_DELAY_MS      500   // Delay before shutdown
#define SAFETY_MAX_VIOLATIONS         3     // Max violations before action

typedef enum {
    SAFETY_VIOLATION_NONE = 0,
    SAFETY_VIOLATION_REVERSE_POWER,
    SAFETY_VIOLATION_PA_CURRENT,
    SAFETY_VIOLATION_TEMPERATURE,
    SAFETY_VIOLATION_MAX
} SafetyViolationType;

typedef struct {
    SafetyViolationType type;
    FemUnit             unit;
    float               measured_value;
    float               threshold;
    uint32_t            timestamp_ms;
    char                description[128];
} SafetyViolation;

typedef struct {
    YamlSafetyConfig yaml_config;
    uint32_t         violation_count[FEM_UNIT_2 + 1][SAFETY_VIOLATION_MAX];
    bool             pa_shutdown_state[FEM_UNIT_2 + 1];
} SafetyConfig;

typedef struct {
    bool            running;
    bool            initialized;
    pthread_t       monitor_thread;
    pthread_mutex_t mutex;
    
    SafetyConfig   config;
    GpioController *gpio_controller;
    I2CController  *i2c_controller;
    
    uint32_t        total_checks;
    uint32_t        total_violations;
    SafetyViolation last_violation;
    
    void (*violation_callback)(const SafetyViolation *violation);
    void (*shutdown_callback)(FemUnit unit, SafetyViolationType reason);
} SafetyMonitor;

int safety_monitor_init(SafetyMonitor *monitor, GpioController *gpio_ctrl, I2CController *i2c_ctrl);
int safety_monitor_start(SafetyMonitor *monitor);
void safety_monitor_stop(SafetyMonitor *monitor);
void safety_monitor_cleanup(SafetyMonitor *monitor);

int safety_monitor_load_yaml_config(SafetyMonitor *monitor, const char *yaml_file);
int safety_monitor_set_thresholds(SafetyMonitor *monitor, float max_reverse_power, float max_current, float max_temp);
int safety_monitor_set_interval(SafetyMonitor *monitor, uint32_t interval_ms);
int safety_monitor_enable(SafetyMonitor *monitor, bool enabled);
int safety_monitor_get_config(SafetyMonitor *monitor, SafetyConfig *config);
int safety_monitor_get_dac_voltages_for_temp(SafetyMonitor *monitor, FemUnit unit, float temperature, float *carrier_voltage, float *peak_voltage);

int safety_monitor_check_fem_unit(SafetyMonitor *monitor, FemUnit unit);
int safety_monitor_check_reverse_power(SafetyMonitor *monitor, FemUnit unit);
int safety_monitor_check_pa_current(SafetyMonitor *monitor, FemUnit unit);
int safety_monitor_check_temperature(SafetyMonitor *monitor, FemUnit unit);

int safety_monitor_shutdown_pa(SafetyMonitor *monitor, FemUnit unit, SafetyViolationType reason);
int safety_monitor_restore_pa(SafetyMonitor *monitor, FemUnit unit);
bool safety_monitor_is_pa_shutdown(SafetyMonitor *monitor, FemUnit unit);

int safety_monitor_get_status(SafetyMonitor *monitor, char *status_json, size_t max_len);
int safety_monitor_get_violations(SafetyMonitor *monitor, FemUnit unit, uint32_t *violation_counts);
void safety_monitor_reset_statistics(SafetyMonitor *monitor);

void safety_monitor_set_violation_callback(SafetyMonitor *monitor, void (*callback)(const SafetyViolation *));
void safety_monitor_set_shutdown_callback(SafetyMonitor *monitor, void (*callback)(FemUnit, SafetyViolationType));

const char* safety_violation_type_to_string(SafetyViolationType type);
uint32_t safety_monitor_get_timestamp_ms(void);

#endif /* SAFETY_MONITOR_H */
