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
#include <unistd.h>
#include <time.h>
#include <sys/time.h>

#include "safety_monitor.h"
#include "femd.h"

// Forward declarations
static void* safety_monitor_thread(void *arg);
static int safety_monitor_perform_checks(SafetyMonitor *monitor);
static int safety_monitor_handle_violation(SafetyMonitor *monitor, const SafetyViolation *violation);
static void safety_monitor_create_violation(SafetyViolation *violation, SafetyViolationType type, 
                                          FemUnit unit, float measured, float threshold, const char *desc);

int safety_monitor_init(SafetyMonitor *monitor, GpioController *gpio_ctrl, I2CController *i2c_ctrl) {
    if (!monitor || !gpio_ctrl || !i2c_ctrl) {
        usys_log_error("Invalid parameters for safety monitor initialization");
        return STATUS_NOK;
    }

    memset(monitor, 0, sizeof(SafetyMonitor));
    
    // Initialize mutex
    if (pthread_mutex_init(&monitor->mutex, NULL) != 0) {
        usys_log_error("Failed to initialize safety monitor mutex");
        return STATUS_NOK;
    }

    monitor->gpio_controller = gpio_ctrl;
    monitor->i2c_controller = i2c_ctrl;
    
    // Load YAML configuration
    if (safety_monitor_load_yaml_config(monitor, YAML_CONFIG_PATH) != STATUS_OK) {
        usys_log_warn("Failed to load YAML config, using defaults");
        yaml_config_set_defaults(&monitor->config.yaml_config);
    }
    
    // Initialize PA shutdown states
    for (int i = 0; i <= FEM_UNIT_2; i++) {
        monitor->config.pa_shutdown_state[i] = false;
    }
    
    monitor->initialized = true;
    usys_log_info("Safety monitor initialized");
    yaml_config_print(&monitor->config.yaml_config);
    
    return STATUS_OK;
}

int safety_monitor_start(SafetyMonitor *monitor) {
    if (!monitor || !monitor->initialized) {
        usys_log_error("Safety monitor not initialized");
        return STATUS_NOK;
    }

    if (monitor->running) {
        usys_log_warn("Safety monitor is already running");
        return STATUS_OK;
    }

    // Create monitoring thread
    if (pthread_create(&monitor->monitor_thread, NULL, safety_monitor_thread, monitor) != 0) {
        usys_log_error("Failed to create safety monitor thread");
        return STATUS_NOK;
    }

    monitor->running = true;
    usys_log_info("Safety monitor started (interval: %u ms)", monitor->config.yaml_config.check_interval_ms);
    
    return STATUS_OK;
}

void safety_monitor_stop(SafetyMonitor *monitor) {
    if (!monitor || !monitor->running) {
        return;
    }

    pthread_mutex_lock(&monitor->mutex);
    monitor->running = false;
    pthread_mutex_unlock(&monitor->mutex);

    // Wait for monitor thread to finish
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
    
    usys_log_info("Safety monitor thread started");
    
    while (monitor->running) {
        pthread_mutex_lock(&monitor->mutex);
        
        if (monitor->config.yaml_config.enabled) {
            safety_monitor_perform_checks(monitor);
            monitor->total_checks++;
        }
        
        pthread_mutex_unlock(&monitor->mutex);
        
        // Sleep for the configured interval
        usleep(monitor->config.yaml_config.check_interval_ms * 1000);
    }
    
    usys_log_info("Safety monitor thread stopped");
    return NULL;
}

static int safety_monitor_perform_checks(SafetyMonitor *monitor) {
    int violations = 0;
    
    // Check both FEM units
    for (FemUnit unit = FEM_UNIT_1; unit <= FEM_UNIT_2; unit++) {
        violations += safety_monitor_check_fem_unit(monitor, unit);
    }
    
    return violations;
}

int safety_monitor_check_fem_unit(SafetyMonitor *monitor, FemUnit unit) {
    if (!monitor || !monitor->initialized) {
        return 0;
    }
    
    int violations = 0;
    
    // Skip checks if PA is already shutdown for this unit
    if (monitor->config.pa_shutdown_state[unit]) {
        return 0;
    }
    
    // Check reverse power
    violations += safety_monitor_check_reverse_power(monitor, unit);
    
    // Check PA current
    violations += safety_monitor_check_pa_current(monitor, unit);
    
    // Check temperature
    violations += safety_monitor_check_temperature(monitor, unit);
    
    return violations;
}

int safety_monitor_check_reverse_power(SafetyMonitor *monitor, FemUnit unit) {
    float reverse_power;
    SafetyViolation violation;
    
    // Read reverse power
    if (adc_read_reverse_power(monitor->i2c_controller, unit, &reverse_power) != STATUS_OK) {
        usys_log_debug("Failed to read reverse power for FEM%d", unit);
        return 0;
    }
    
    // Check threshold
    if (reverse_power > monitor->config.yaml_config.max_reverse_power_dbm) {
        safety_monitor_create_violation(&violation, SAFETY_VIOLATION_REVERSE_POWER, unit,
                                      reverse_power, monitor->config.yaml_config.max_reverse_power_dbm,
                                      "Reverse power exceeded threshold");
        
        return safety_monitor_handle_violation(monitor, &violation);
    }
    
    return 0;
}

int safety_monitor_check_pa_current(SafetyMonitor *monitor, FemUnit unit) {
    float pa_current;
    SafetyViolation violation;
    
    // Read PA current
    if (adc_read_pa_current(monitor->i2c_controller, unit, &pa_current) != STATUS_OK) {
        usys_log_debug("Failed to read PA current for FEM%d", unit);
        return 0;
    }
    
    // Check threshold
    if (pa_current > monitor->config.yaml_config.max_pa_current_a) {
        safety_monitor_create_violation(&violation, SAFETY_VIOLATION_PA_CURRENT, unit,
                                      pa_current, monitor->config.yaml_config.max_pa_current_a,
                                      "PA current exceeded threshold");
        
        return safety_monitor_handle_violation(monitor, &violation);
    }
    
    return 0;
}

int safety_monitor_check_temperature(SafetyMonitor *monitor, FemUnit unit) {
    float temperature;
    SafetyViolation violation;
    
    // Read temperature
    if (temp_sensor_read(monitor->i2c_controller, unit, &temperature) != STATUS_OK) {
        usys_log_debug("Failed to read temperature for FEM%d", unit);
        return 0;
    }
    
    // Check threshold
    if (temperature > monitor->config.yaml_config.max_temperature_c) {
        safety_monitor_create_violation(&violation, SAFETY_VIOLATION_TEMPERATURE, unit,
                                      temperature, monitor->config.yaml_config.max_temperature_c,
                                      "Temperature exceeded threshold");
        
        return safety_monitor_handle_violation(monitor, &violation);
    }
    
    return 0;
}

static void safety_monitor_create_violation(SafetyViolation *violation, SafetyViolationType type, 
                                          FemUnit unit, float measured, float threshold, const char *desc) {
    violation->type = type;
    violation->unit = unit;
    violation->measured_value = measured;
    violation->threshold = threshold;
    violation->timestamp_ms = safety_monitor_get_timestamp_ms();
    strncpy(violation->description, desc, sizeof(violation->description) - 1);
    violation->description[sizeof(violation->description) - 1] = '\0';
}

static int safety_monitor_handle_violation(SafetyMonitor *monitor, const SafetyViolation *violation) {
    // Update violation count
    monitor->config.violation_count[violation->unit][violation->type]++;
    monitor->total_violations++;
    
    // Store last violation
    monitor->last_violation = *violation;
    
    // Log the violation
    usys_log_warn("SAFETY VIOLATION: FEM%d %s - Measured: %.2f, Threshold: %.2f", 
                  violation->unit,
                  safety_violation_type_to_string(violation->type),
                  violation->measured_value,
                  violation->threshold);
    
    // Call violation callback if set
    if (monitor->violation_callback) {
        monitor->violation_callback(violation);
    }
    
    // Check if we need to shut down PA
    uint32_t violation_count = monitor->config.violation_count[violation->unit][violation->type];
    if (violation_count >= monitor->config.yaml_config.max_violations) {
        usys_log_error("SAFETY SHUTDOWN: FEM%d - %u violations of %s", 
                       violation->unit, violation_count,
                       safety_violation_type_to_string(violation->type));
        
        safety_monitor_shutdown_pa(monitor, violation->unit, violation->type);
    }
    
    return 1;
}

int safety_monitor_shutdown_pa(SafetyMonitor *monitor, FemUnit unit, SafetyViolationType reason) {
    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }
    
    usys_log_error("EXECUTING PA SHUTDOWN for FEM%d due to %s", 
                   unit, safety_violation_type_to_string(reason));
    
    // 1. Disable PA by setting DAC values to zero
    if (dac_disable_pa(monitor->i2c_controller, unit) != STATUS_OK) {
        usys_log_error("Failed to disable DAC for FEM%d", unit);
    }
    
    // 2. Disable PA_VDS_Enable GPIO
    if (gpio_set_pa_vds(monitor->gpio_controller, unit, false) != STATUS_OK) {
        usys_log_error("Failed to disable PA_VDS GPIO for FEM%d", unit);
    }
    
    // 3. Disable 28V_VDS (pa_disable = 1, which means VDS disabled)
    if (gpio_set_28v_vds(monitor->gpio_controller, unit, false) != STATUS_OK) {
        usys_log_error("Failed to disable 28V_VDS for FEM%d", unit);
    }
    
    // 4. Disable TX_RF
    if (gpio_set_tx_rf(monitor->gpio_controller, unit, false) != STATUS_OK) {
        usys_log_error("Failed to disable TX_RF for FEM%d", unit);
    }
    
    // Mark PA as shutdown
    monitor->config.pa_shutdown_state[unit] = true;
    
    // Call shutdown callback if set
    if (monitor->shutdown_callback) {
        monitor->shutdown_callback(unit, reason);
    }
    
    usys_log_error("PA SHUTDOWN COMPLETE for FEM%d", unit);
    return STATUS_OK;
}

int safety_monitor_restore_pa(SafetyMonitor *monitor, FemUnit unit) {
    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }
    
    if (!monitor->config.pa_shutdown_state[unit]) {
        usys_log_info("PA for FEM%d is not shutdown", unit);
        return STATUS_OK;
    }
    
    usys_log_info("Restoring PA for FEM%d", unit);
    
    // Reset violation counts for this unit
    for (int i = 0; i < SAFETY_VIOLATION_MAX; i++) {
        monitor->config.violation_count[unit][i] = 0;
    }
    
    // Mark PA as not shutdown (manual restoration)
    monitor->config.pa_shutdown_state[unit] = false;
    
    usys_log_info("PA restored for FEM%d (manual intervention required for re-enabling)", unit);
    return STATUS_OK;
}

bool safety_monitor_is_pa_shutdown(SafetyMonitor *monitor, FemUnit unit) {
    if (!monitor || !monitor->initialized) {
        return false;
    }
    
    return monitor->config.pa_shutdown_state[unit];
}

// Configuration functions
int safety_monitor_load_yaml_config(SafetyMonitor *monitor, const char *yaml_file) {
    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }
    
    pthread_mutex_lock(&monitor->mutex);
    int result = yaml_config_load(yaml_file, &monitor->config.yaml_config);
    pthread_mutex_unlock(&monitor->mutex);
    
    if (result == STATUS_OK) {
        usys_log_info("YAML safety configuration loaded from %s", yaml_file);
    }
    
    return result;
}

int safety_monitor_get_dac_voltages_for_temp(SafetyMonitor *monitor, FemUnit unit, float temperature, 
                                            float *carrier_voltage, float *peak_voltage) {
    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }
    
    pthread_mutex_lock(&monitor->mutex);
    int result = yaml_config_get_dac_voltages_for_temp(&monitor->config.yaml_config, unit, 
                                                      temperature, carrier_voltage, peak_voltage);
    pthread_mutex_unlock(&monitor->mutex);
    
    return result;
}

int safety_monitor_set_thresholds(SafetyMonitor *monitor, float max_reverse_power, float max_current, float max_temp) {
    if (!monitor || !monitor->initialized) {
        return STATUS_NOK;
    }
    
    pthread_mutex_lock(&monitor->mutex);
    
    monitor->config.yaml_config.max_reverse_power_dbm = max_reverse_power;
    monitor->config.yaml_config.max_pa_current_a = max_current;
    monitor->config.yaml_config.max_temperature_c = max_temp;
    
    pthread_mutex_unlock(&monitor->mutex);
    
    usys_log_info("Safety thresholds updated: RP=%.1fdBm, Current=%.1fA, Temp=%.1f°C",
                  max_reverse_power, max_current, max_temp);
    
    return STATUS_OK;
}

int safety_monitor_set_interval(SafetyMonitor *monitor, uint32_t interval_ms) {
    if (!monitor || !monitor->initialized || interval_ms < 100) {
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
    monitor->config.yaml_config.enabled = enabled;
    pthread_mutex_unlock(&monitor->mutex);
    
    usys_log_info("Safety monitor %s", enabled ? "enabled" : "disabled");
    return STATUS_OK;
}

int safety_monitor_get_status(SafetyMonitor *monitor, char *status_json, size_t max_len) {
    if (!monitor || !status_json || max_len == 0) {
        return STATUS_NOK;
    }
    
    pthread_mutex_lock(&monitor->mutex);
    
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
             monitor->config.yaml_config.enabled ? "true" : "false",
             monitor->running ? "true" : "false",
             monitor->total_checks,
             monitor->total_violations,
             monitor->config.yaml_config.max_reverse_power_dbm,
             monitor->config.yaml_config.max_pa_current_a,
             monitor->config.yaml_config.max_temperature_c,
             monitor->config.pa_shutdown_state[FEM_UNIT_1] ? "true" : "false",
             monitor->config.pa_shutdown_state[FEM_UNIT_2] ? "true" : "false",
             monitor->config.yaml_config.check_interval_ms);
    
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

// Callback functions
void safety_monitor_set_violation_callback(SafetyMonitor *monitor, void (*callback)(const SafetyViolation *)) {
    if (monitor) {
        monitor->violation_callback = callback;
    }
}

void safety_monitor_set_shutdown_callback(SafetyMonitor *monitor, void (*callback)(FemUnit, SafetyViolationType)) {
    if (monitor) {
        monitor->shutdown_callback = callback;
    }
}

// Utility functions
const char* safety_violation_type_to_string(SafetyViolationType type) {
    switch (type) {
        case SAFETY_VIOLATION_REVERSE_POWER: return "Reverse Power";
        case SAFETY_VIOLATION_PA_CURRENT: return "PA Current";
        case SAFETY_VIOLATION_TEMPERATURE: return "Temperature";
        default: return "Unknown";
    }
}

uint32_t safety_monitor_get_timestamp_ms(void) {
    struct timeval tv;
    gettimeofday(&tv, NULL);
    return (uint32_t)(tv.tv_sec * 1000 + tv.tv_usec / 1000);
}