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
#include <time.h>
#include <sys/time.h>

#include "safety_monitor.h"
#include "femd.h"

/* Forward statics */
static void* safety_monitor_thread(void *arg);
static int   safety_monitor_perform_checks(SafetyMonitor *monitor);
static int   safety_monitor_handle_violation(SafetyMonitor *monitor,
                                             const SafetyViolation *violation);
static void  safety_monitor_create_violation(SafetyViolation *violation, SafetyViolationType type,
                                             FemUnit unit, float measured, float threshold,
                                             const char *desc);

int safety_monitor_init(SafetyMonitor *monitor, GpioController *gpio_ctrl, I2CController *i2c_ctrl) {
    int i;

    if (!monitor || !gpio_ctrl || !i2c_ctrl) {
        usys_log_error("Invalid parameters for safety monitor initialization");
        return STATUS_NOK;
    }

    memset(monitor, 0, sizeof(SafetyMonitor));

    if (pthread_mutex_init(&monitor->mutex, NULL) != 0) {
        usys_log_error("Failed to initialize safety monitor mutex");
        return STATUS_NOK;
    }

    monitor->gpio_controller = gpio_ctrl;
    monitor->i2c_controller  = i2c_ctrl;

    /* Important: mark initialized BEFORE loading YAML (loader checks this) */
    monitor->initialized = true;

    if (safety_monitor_load_yaml_config(monitor, YAML_CONFIG_PATH) != STATUS_OK) {
        usys_log_warn("Failed to load YAML config, using defaults");
        yaml_config_set_defaults(&monitor->config.yaml_config);
    }

    /* Clear per-unit state (1-based FemUnit) */
    for (i = 0; i <= FEM_UNIT_2; ++i) {
        monitor->config.pa_shutdown_state[i] = false;
    }
    memset(monitor->config.violation_count, 0, sizeof(monitor->config.violation_count));

    monitor->running = false;
    monitor->total_checks = 0;
    monitor->total_violations = 0;
    memset(&monitor->last_violation, 0, sizeof(monitor->last_violation));

    usys_log_info("Safety monitor initialized");
    yaml_config_print(&monitor->config.yaml_config);

    return STATUS_OK;
}

int safety_monitor_start(SafetyMonitor *monitor) {
    if (!monitor || !monitor->initialized) {
        usys_log_error("Safety monitor not initialized");
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    if (monitor->running) {
        pthread_mutex_unlock(&monitor->mutex);
        usys_log_warn("Safety monitor is already running");
        return STATUS_OK;
    }
    monitor->running = true; /* set before creating the thread to avoid race */
    pthread_mutex_unlock(&monitor->mutex);

    if (pthread_create(&monitor->monitor_thread, NULL, safety_monitor_thread, monitor) != 0) {
        usys_log_error("Failed to create safety monitor thread");
        pthread_mutex_lock(&monitor->mutex);
        monitor->running = false;
        pthread_mutex_unlock(&monitor->mutex);
        return STATUS_NOK;
    }

    usys_log_info("Safety monitor started (interval: %u ms)",
                  monitor->config.yaml_config.check_interval_ms);
    return STATUS_OK;
}

void safety_monitor_stop(SafetyMonitor *monitor) {
    if (!monitor) {
        return;
    }

    pthread_mutex_lock(&monitor->mutex);
    if (!monitor->running) {
        pthread_mutex_unlock(&monitor->mutex);
        return;
    }
    monitor->running = false;
    pthread_mutex_unlock(&monitor->mutex);

    pthread_join(monitor->monitor_thread, NULL);
    usys_log_info("Safety monitor stopped");
}

void safety_monitor_cleanup(SafetyMonitor *monitor) {
    if (!monitor) {
        return;
    }

    safety_monitor_stop(monitor);

    if (monitor->initialized) {
        pthread_mutex_destroy(&monitor->mutex);
        monitor->initialized = false;
    }

    usys_log_info("Safety monitor cleaned up");
}

static void* safety_monitor_thread(void *arg) {
    SafetyMonitor *monitor = (SafetyMonitor*)arg;

    /* Local snapshots to minimize lock contention */
    int enabled;
    uint32_t interval_ms;

    usys_log_info("Safety monitor thread started");

    for (;;) {
        /* Snapshot control flags */
        pthread_mutex_lock(&monitor->mutex);
        if (!monitor->running) {
            pthread_mutex_unlock(&monitor->mutex);
            break;
        }
        enabled     = monitor->config.yaml_config.enabled ? 1 : 0;
        interval_ms = monitor->config.yaml_config.check_interval_ms;
        if (interval_ms < 100U) {
            interval_ms = 100U; /* floor to a sane minimum */
        }
        pthread_mutex_unlock(&monitor->mutex);

        /* Perform checks without holding the mutex */
        if (enabled) {
            (void)safety_monitor_perform_checks(monitor);

            /* Update totals under lock */
            pthread_mutex_lock(&monitor->mutex);
            monitor->total_checks += 1;
            pthread_mutex_unlock(&monitor->mutex);
        }

        usleep(interval_ms * 1000U);
    }

    usys_log_info("Safety monitor thread stopped");
    return NULL;
}

static int safety_monitor_perform_checks(SafetyMonitor *monitor) {
    int violations = 0;
    FemUnit unit;

    for (unit = FEM_UNIT_1; unit <= FEM_UNIT_2; unit++) {
        violations += safety_monitor_check_fem_unit(monitor, unit);
    }

    return violations;
}

int safety_monitor_check_fem_unit(SafetyMonitor *monitor, FemUnit unit) {
    int violations = 0;

    if (!monitor || !monitor->initialized) {
        return 0;
    }

    if (unit < FEM_UNIT_1 || unit > FEM_UNIT_2) {
        return 0;
    }

    /* If this PA is already shutdown, skip checks */
    if (monitor->config.pa_shutdown_state[unit]) {
        return 0;
    }

    violations += safety_monitor_check_reverse_power(monitor, unit);
    violations += safety_monitor_check_pa_current(monitor, unit);
    violations += safety_monitor_check_temperature(monitor, unit);

    return violations;
}

int safety_monitor_check_reverse_power(SafetyMonitor *monitor, FemUnit unit) {
    float reverse_power = 0.0f;
    SafetyViolation violation;

    if (adc_read_reverse_power(monitor->i2c_controller, unit, &reverse_power) != STATUS_OK) {
        usys_log_debug("Failed to read reverse power for FEM%d", unit);
        return 0;
    }

    if (reverse_power > monitor->config.yaml_config.max_reverse_power_dbm) {
        safety_monitor_create_violation(&violation,
                                        SAFETY_VIOLATION_REVERSE_POWER, unit,
                                        reverse_power,
                                        monitor->config.yaml_config.max_reverse_power_dbm,
                                        "Reverse power exceeded threshold");
        return safety_monitor_handle_violation(monitor, &violation);
    }

    return 0;
}

int safety_monitor_check_pa_current(SafetyMonitor *monitor, FemUnit unit) {
    float pa_current = 0.0f;
    SafetyViolation violation;

    if (adc_read_pa_current(monitor->i2c_controller, unit, &pa_current) != STATUS_OK) {
        usys_log_debug("Failed to read PA current for FEM%d", unit);
        return 0;
    }

    if (pa_current > monitor->config.yaml_config.max_pa_current_a) {
        safety_monitor_create_violation(&violation, SAFETY_VIOLATION_PA_CURRENT, unit,
                                        pa_current, monitor->config.yaml_config.max_pa_current_a,
                                        "PA current exceeded threshold");
        return safety_monitor_handle_violation(monitor, &violation);
    }

    return 0;
}

int safety_monitor_check_temperature(SafetyMonitor *monitor, FemUnit unit) {
    float temperature = 0.0f;
    SafetyViolation violation;

    if (temp_sensor_read(monitor->i2c_controller, unit, &temperature) != STATUS_OK) {
        usys_log_debug("Failed to read temperature for FEM%d", unit);
        return 0;
    }

    if (temperature > monitor->config.yaml_config.max_temperature_c) {
        safety_monitor_create_violation(&violation, SAFETY_VIOLATION_TEMPERATURE, unit,
                                        temperature, monitor->config.yaml_config.max_temperature_c,
                                        "Temperature exceeded threshold");
        return safety_monitor_handle_violation(monitor, &violation);
    }

    return 0;
}

static void safety_monitor_create_violation(SafetyViolation *violation,
                                            SafetyViolationType type,
                                            FemUnit unit,
                                            float measured,
                                            float threshold, const char *desc) {
    if (!violation) return;

    violation->type = type;
    violation->unit = unit;
    violation->measured_value = measured;
    violation->threshold = threshold;
    violation->timestamp_ms = safety_monitor_get_timestamp_ms();
    strncpy(violation->description, desc ? desc : "", sizeof(violation->description) - 1);
    violation->description[sizeof(violation->description) - 1] = '\0';
}

static int safety_monitor_handle_violation(SafetyMonitor *monitor, const SafetyViolation *violation) {
    uint32_t count;
    int immediate;

    if (!monitor || !violation) {
        return 0;
    }

    /* Update counters and last_violation under lock */
    pthread_mutex_lock(&monitor->mutex);
    monitor->config.violation_count[violation->unit][violation->type] += 1;
    monitor->total_violations += 1;
    monitor->last_violation = *violation;
    count = monitor->config.violation_count[violation->unit][violation->type];
    immediate = monitor->config.yaml_config.emergency_immediate_shutdown ? 1 : 0;
    pthread_mutex_unlock(&monitor->mutex);

    if (monitor->config.yaml_config.emergency_log_event) {
        usys_log_warn("SAFETY VIOLATION: FEM%d %s - Measured: %.2f, Threshold: %.2f (count=%u)",
                      violation->unit,
                      safety_violation_type_to_string(violation->type),
                      violation->measured_value,
                      violation->threshold,
                      (unsigned int)count);
    } else {
        usys_log_warn("SAFETY VIOLATION: FEM%d %s (count=%u)",
                      violation->unit,
                      safety_violation_type_to_string(violation->type),
                      (unsigned int)count);
    }

    /* Immediate shutdown policy or thresholded shutdown */
    if (immediate || count >= monitor->config.yaml_config.max_violations) {
        usys_log_error("SAFETY SHUTDOWN: FEM%d - %u violations of %s (immediate=%d)",
                       violation->unit, (unsigned int)count,
                       safety_violation_type_to_string(violation->type),
                       immediate ? 1 : 0);
        (void)safety_monitor_shutdown_pa(monitor, violation->unit, violation->type);
    }

    /* Return 1 to count this as a violation in the periodic summary */
    return 1;
}

int safety_monitor_shutdown_pa(SafetyMonitor *monitor, FemUnit unit, SafetyViolationType reason) {
    int rc;

    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }
    if (unit < FEM_UNIT_1 || unit > FEM_UNIT_2) {
        return STATUS_NOK;
    }

    usys_log_error("EXECUTING PA SHUTDOWN for FEM%d due to %s",
                   unit, safety_violation_type_to_string(reason));

    /* Always zero DAC drive */
    rc = dac_disable_pa(monitor->i2c_controller, unit);
    if (rc != STATUS_OK) {
        usys_log_error("Failed to disable DAC for FEM%d", unit);
    }

    /* Honor emergency flags for GPIO cuts */
    if (monitor->config.yaml_config.emergency_disable_pa_vds) {
        if (gpio_set(monitor->gpio_controller, unit, GPIO_PA_VDS, false) != STATUS_OK) {
            usys_log_error("Failed to disable PA_VDS for FEM%d", unit);
        }
    }

    if (monitor->config.yaml_config.emergency_disable_28v_vds) {
        if (gpio_set(monitor->gpio_controller, unit, GPIO_28V_VDS, false) != STATUS_OK) {
            usys_log_error("Failed to disable 28V_VDS for FEM%d", unit);
        }
    }

    if (monitor->config.yaml_config.emergency_disable_tx_rf) {
        if (gpio_set(monitor->gpio_controller, unit, GPIO_TX_RF, false) != STATUS_OK) {
            usys_log_error("Failed to disable TX_RF for FEM%d", unit);
        }
    }

    /* Mark PA as shutdown and invoke callback (under lock for coherence) */
    pthread_mutex_lock(&monitor->mutex);
    monitor->config.pa_shutdown_state[unit] = true;
    pthread_mutex_unlock(&monitor->mutex);

    if (monitor->shutdown_callback) {
        monitor->shutdown_callback(unit, reason);
    }

    usys_log_error("PA SHUTDOWN COMPLETE for FEM%d", unit);
    return STATUS_OK;
}

int safety_monitor_restore_pa(SafetyMonitor *monitor, FemUnit unit) {
    int i;

    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }
    if (unit < FEM_UNIT_1 || unit > FEM_UNIT_2) {
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    if (!monitor->config.pa_shutdown_state[unit]) {
        pthread_mutex_unlock(&monitor->mutex);
        usys_log_info("PA for FEM%d is not shutdown", unit);
        return STATUS_OK;
    }

    usys_log_info("Restoring PA for FEM%d", unit);

    for (i = 0; i < SAFETY_VIOLATION_MAX; i++) {
        monitor->config.violation_count[unit][i] = 0;
    }
    monitor->config.pa_shutdown_state[unit] = false;
    pthread_mutex_unlock(&monitor->mutex);

    usys_log_info("PA restored for FEM%d (manual intervention required for re-enabling)", unit);
    return STATUS_OK;
}

bool safety_monitor_is_pa_shutdown(SafetyMonitor *monitor, FemUnit unit) {
    bool is_down;

    if (!monitor || !monitor->initialized) {
        return false;
    }
    if (unit < FEM_UNIT_1 || unit > FEM_UNIT_2) {
        return false;
    }

    pthread_mutex_lock(&monitor->mutex);
    is_down = monitor->config.pa_shutdown_state[unit] ? true : false;
    pthread_mutex_unlock(&monitor->mutex);

    return is_down;
}

int safety_monitor_load_yaml_config(SafetyMonitor *monitor, const char *yaml_file) {
    int result;

    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    result = yaml_config_load(yaml_file, &monitor->config.yaml_config);
    pthread_mutex_unlock(&monitor->mutex);

    if (result == STATUS_OK) {
        usys_log_info("YAML safety configuration loaded from %s", yaml_file);
    }
    return result;
}

int safety_monitor_get_config(SafetyMonitor *monitor, SafetyConfig *config) {
    if (!monitor || !config) {
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    *config = monitor->config; /* struct copy (has small fixed arrays) */
    pthread_mutex_unlock(&monitor->mutex);

    return STATUS_OK;
}

int safety_monitor_get_dac_voltages_for_temp(SafetyMonitor *monitor,
                                             FemUnit unit,
                                             float temperature,
                                             float *carrier_voltage,
                                             float *peak_voltage) {
    int result;

    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    result = yaml_config_get_dac_voltages_for_temp(&monitor->config.yaml_config, unit,
                                                   temperature, carrier_voltage, peak_voltage);
    pthread_mutex_unlock(&monitor->mutex);

    return result;
}

int safety_monitor_set_thresholds(SafetyMonitor *monitor,
                                  float max_reverse_power,
                                  float max_current,
                                  float max_temp) {

    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    monitor->config.yaml_config.max_reverse_power_dbm = max_reverse_power;
    monitor->config.yaml_config.max_pa_current_a      = max_current;
    monitor->config.yaml_config.max_temperature_c     = max_temp;
    pthread_mutex_unlock(&monitor->mutex);

    usys_log_info("Safety thresholds updated: RP=%.1fdBm, Current=%.1fA, Temp=%.1fC",
                  max_reverse_power, max_current, max_temp);
    return STATUS_OK;
}

int safety_monitor_set_interval(SafetyMonitor *monitor, uint32_t interval_ms) {
    if (!monitor || !monitor->initialized || interval_ms < 100U) {
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    monitor->config.yaml_config.check_interval_ms = interval_ms;
    pthread_mutex_unlock(&monitor->mutex);

    usys_log_info("Safety monitor interval set to %u ms", interval_ms);
    return STATUS_OK;
}

int safety_monitor_enable(SafetyMonitor *monitor, bool enabled) {
    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    monitor->config.yaml_config.enabled = enabled ? true : false;
    pthread_mutex_unlock(&monitor->mutex);

    usys_log_info("Safety monitor %s", enabled ? "enabled" : "disabled");
    return STATUS_OK;
}

int safety_monitor_get_status(SafetyMonitor *monitor, char *status_json, size_t max_len) {
    int enabled, running;
    uint32_t total_checks, total_violations, interval_ms;
    float max_rp, max_i, max_t;
    int pa1, pa2;

    if (!monitor || !status_json || max_len == 0U) {
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    enabled          = monitor->config.yaml_config.enabled ? 1 : 0;
    running          = monitor->running ? 1 : 0;
    total_checks     = monitor->total_checks;
    total_violations = monitor->total_violations;
    max_rp           = monitor->config.yaml_config.max_reverse_power_dbm;
    max_i            = monitor->config.yaml_config.max_pa_current_a;
    max_t            = monitor->config.yaml_config.max_temperature_c;
    pa1              = monitor->config.pa_shutdown_state[FEM_UNIT_1] ? 1 : 0;
    pa2              = monitor->config.pa_shutdown_state[FEM_UNIT_2] ? 1 : 0;
    interval_ms      = monitor->config.yaml_config.check_interval_ms;
    pthread_mutex_unlock(&monitor->mutex);

    snprintf(status_json, max_len,
             "{"
             "\"enabled\":%s,"
             "\"running\":%s,"
             "\"total_checks\":%u,"
             "\"total_violations\":%u,"
             "\"thresholds\":{"
                "\"max_reverse_power\":%.1f,"
                "\"max_pa_current\":%.1f,"
                "\"max_temperature\":%.1f"
             "},"
             "\"pa_shutdown\":{"
                "\"fem1\":%s,"
                "\"fem2\":%s"
             "},"
             "\"check_interval_ms\":%u"
             "}",
             enabled ? "true" : "false",
             running ? "true" : "false",
             (unsigned int)total_checks,
             (unsigned int)total_violations,
             max_rp, max_i, max_t,
             pa1 ? "true" : "false",
             pa2 ? "true" : "false",
             (unsigned int)interval_ms);

    return STATUS_OK;
}

int safety_monitor_get_violations(SafetyMonitor *monitor,
                                  FemUnit unit,
                                  uint32_t *violation_counts) {
    int i;

    if (!monitor || !violation_counts) {
        return STATUS_NOK;
    }
    if (unit < FEM_UNIT_1 || unit > FEM_UNIT_2) {
        return STATUS_NOK;
    }

    pthread_mutex_lock(&monitor->mutex);
    for (i = 0; i < SAFETY_VIOLATION_MAX; ++i) {
        violation_counts[i] = monitor->config.violation_count[unit][i];
    }
    pthread_mutex_unlock(&monitor->mutex);

    return STATUS_OK;
}

void safety_monitor_reset_statistics(SafetyMonitor *monitor) {
    if (!monitor || !monitor->initialized) {
        return;
    }

    pthread_mutex_lock(&monitor->mutex);
    monitor->total_checks = 0;
    monitor->total_violations = 0;
    memset(monitor->config.violation_count, 0, sizeof(monitor->config.violation_count));
    memset(&monitor->last_violation, 0, sizeof(monitor->last_violation));
    pthread_mutex_unlock(&monitor->mutex);

    usys_log_info("Safety monitor statistics reset");
}

void safety_monitor_set_violation_callback(SafetyMonitor *monitor,
                                           void (*callback)(const SafetyViolation *)) {
    if (monitor) {
        monitor->violation_callback = callback;
    }
}

void safety_monitor_set_shutdown_callback(SafetyMonitor *monitor,
                                          void (*callback)(FemUnit, SafetyViolationType)) {
    if (monitor) {
        monitor->shutdown_callback = callback;
    }
}

const char* safety_violation_type_to_string(SafetyViolationType type) {
    switch (type) {
        case SAFETY_VIOLATION_REVERSE_POWER: return "Reverse Power";
        case SAFETY_VIOLATION_PA_CURRENT:    return "PA Current";
        case SAFETY_VIOLATION_TEMPERATURE:   return "Temperature";
        default:                             return "Unknown";
    }
}

uint32_t safety_monitor_get_timestamp_ms(void) {
    struct timeval tv;
    gettimeofday(&tv, NULL);
    return (uint32_t)(tv.tv_sec * 1000UL + tv.tv_usec / 1000UL);
}
