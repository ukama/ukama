/*
 * Example: How to use the YAML-based temperature compensation in I2C controller
 * This demonstrates how the temperature-based DAC voltage lookup tables work
 */

#include "safety_monitor.h"
#include "yaml_config.h"
#include "i2c_controller.h"

// Example function that would be integrated into i2c_controller.c
int dac_set_temperature_compensated_voltage(I2CController *i2c_ctrl, SafetyMonitor *safety_monitor, 
                                           FemUnit unit, float current_temperature) {
    float carrier_voltage, peak_voltage;
    
    // Get temperature-compensated voltages from YAML lookup tables
    if (safety_monitor_get_dac_voltages_for_temp(safety_monitor, unit, current_temperature, 
                                                &carrier_voltage, &peak_voltage) != STATUS_OK) {
        usys_log_error("Failed to get temperature-compensated voltages for FEM%d at %.1f°C", 
                       unit, current_temperature);
        return STATUS_NOK;
    }
    
    usys_log_info("FEM%d temp compensation: %.1f°C -> Carrier: %.3fV, Peak: %.3fV", 
                  unit, current_temperature, carrier_voltage, peak_voltage);
    
    // Apply the voltages to the DAC
    if (dac_set_voltage(i2c_ctrl, unit, DAC_CHANNEL_CARRIER, carrier_voltage) != STATUS_OK) {
        usys_log_error("Failed to set carrier voltage %.3fV for FEM%d", carrier_voltage, unit);
        return STATUS_NOK;
    }
    
    if (dac_set_voltage(i2c_ctrl, unit, DAC_CHANNEL_PEAK, peak_voltage) != STATUS_OK) {
        usys_log_error("Failed to set peak voltage %.3fV for FEM%d", peak_voltage, unit);
        return STATUS_NOK;
    }
    
    // Update the DAC state
    if (unit == FEM_UNIT_1) {
        i2c_ctrl->dac_state.carrier_voltage = carrier_voltage;
        i2c_ctrl->dac_state.peak_voltage = peak_voltage;
    } else if (unit == FEM_UNIT_2) {
        // Assuming FEM2 has separate DAC state or use same structure
        i2c_ctrl->dac_state.carrier_voltage = carrier_voltage; 
        i2c_ctrl->dac_state.peak_voltage = peak_voltage;
    }
    
    usys_log_info("Temperature compensation applied successfully for FEM%d", unit);
    return STATUS_OK;
}

// Example background monitoring function that would run periodically
int monitor_and_adjust_for_temperature(I2CController *i2c_ctrl, SafetyMonitor *safety_monitor) {
    // Read temperatures for both FEM units
    for (FemUnit unit = FEM_UNIT_1; unit <= FEM_UNIT_2; unit++) {
        float temperature;
        
        // Read current temperature
        if (temp_sensor_read(i2c_ctrl, unit, &temperature) != STATUS_OK) {
            usys_log_warn("Failed to read temperature for FEM%d", unit);
            continue;
        }
        
        // Apply temperature compensation
        if (dac_set_temperature_compensated_voltage(i2c_ctrl, safety_monitor, unit, temperature) != STATUS_OK) {
            usys_log_error("Failed to apply temperature compensation for FEM%d", unit);
        }
    }
    
    return STATUS_OK;
}

/*
 * Usage example in main daemon loop or in a dedicated temperature monitoring thread:
 *
 * while (daemon_running) {
 *     // Every 5 seconds, check temperature and adjust voltages
 *     monitor_and_adjust_for_temperature(&i2c_controller, &safety_monitor);
 *     sleep(5);
 * }
 */