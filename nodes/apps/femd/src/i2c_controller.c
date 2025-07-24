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
#include <math.h>
#include <stddef.h>
#include <fcntl.h>
#include <sys/ioctl.h>
#include <linux/i2c-dev.h>
#include <errno.h>

#include "i2c_controller.h"
#include "femd.h"

static const I2CDeviceInfo device_info[I2C_DEVICE_MAX] = {
    {"AD5667", I2C_ADDR_DAC_AD5667, "16-bit DAC"},
    {"LM75A", I2C_ADDR_TEMP_LM75A, "Temperature sensor"},
    {"ADS1015", I2C_ADDR_ADC_ADS1015, "12-bit ADC"},
    {"EEPROM", I2C_ADDR_EEPROM, "Serial number storage"}
};

int i2c_controller_init(I2CController *controller) {
    if (!controller) {
        usys_log_error("I2C controller pointer is NULL");
        return STATUS_NOK;
    }

    memset(controller, 0, sizeof(I2CController));
    
    controller->bus_fem1 = I2C_BUS_FEM1;
    controller->bus_fem2 = I2C_BUS_FEM2;
    
    controller->dac_state.carrier_voltage = 0.0f;
    controller->dac_state.peak_voltage = 0.0f;
    controller->dac_state.initialized = false;
    
    controller->temp_state.temperature = 0.0f;
    controller->temp_state.threshold = 85.0f;  // Default threshold
    controller->temp_state.alert_enabled = false;
    
    controller->adc_state.max_reverse_power = -10.0f;  // Default -10 dBm
    controller->adc_state.max_current = 5.0f;          // Default 5A
    controller->adc_state.safety_enabled = true;
    
    controller->eeprom_state.has_data = false;
    memset(controller->eeprom_state.serial_number, 0, sizeof(controller->eeprom_state.serial_number));
    
    controller->initialized = true;
    usys_log_info("I2C controller initialized");
    
    return STATUS_OK;
}

void i2c_controller_cleanup(I2CController *controller) {
    if (controller) {
        controller->initialized = false;
        usys_log_info("I2C controller cleaned up");
    }
}

int i2c_get_bus_for_fem(FemUnit unit) {
    return (unit == FEM_UNIT_1) ? I2C_BUS_FEM1 : I2C_BUS_FEM2;
}

int i2c_write_bytes(int bus, uint8_t device_addr, uint8_t reg, const uint8_t *data, size_t len) {
    char device_path[32];
    int fd;
    uint8_t buffer[256];
    
    if (len > 254) {  // Reserve space for register byte
        usys_log_error("I2C write data too long: %zu bytes", len);
        return STATUS_NOK;
    }
    
    snprintf(device_path, sizeof(device_path), "/dev/i2c-%d", bus);
    
    fd = open(device_path, O_RDWR);
    if (fd < 0) {
        usys_log_error("Failed to open I2C device %s: %s", device_path, strerror(errno));
        return STATUS_NOK;
    }
    
    if (ioctl(fd, I2C_SLAVE, device_addr) < 0) {
        usys_log_error("Failed to set I2C slave address 0x%02X: %s", device_addr, strerror(errno));
        close(fd);
        return STATUS_NOK;
    }
    
    buffer[0] = reg;
    memcpy(&buffer[1], data, len);
    
    if (write(fd, buffer, len + 1) != (ssize_t)(len + 1)) {
        usys_log_error("I2C write failed: bus=%d, addr=0x%02X, reg=0x%02X: %s", 
                       bus, device_addr, reg, strerror(errno));
        close(fd);
        return STATUS_NOK;
    }
    
    close(fd);
    usys_log_debug("I2C write successful: bus=%d, addr=0x%02X, reg=0x%02X, len=%zu", 
                   bus, device_addr, reg, len);
    
    return STATUS_OK;
}

int i2c_read_bytes(int bus, uint8_t device_addr, uint8_t reg, uint8_t *data, size_t len) {
    char device_path[32];
    int fd;
    
    if (!data || len == 0) {
        usys_log_error("Invalid parameters for I2C read");
        return STATUS_NOK;
    }
    
    snprintf(device_path, sizeof(device_path), "/dev/i2c-%d", bus);
    
    fd = open(device_path, O_RDWR);
    if (fd < 0) {
        usys_log_error("Failed to open I2C device %s: %s", device_path, strerror(errno));
        return STATUS_NOK;
    }
    
    if (ioctl(fd, I2C_SLAVE, device_addr) < 0) {
        usys_log_error("Failed to set I2C slave address 0x%02X: %s", device_addr, strerror(errno));
        close(fd);
        return STATUS_NOK;
    }
    
    // Write register address
    if (write(fd, &reg, 1) != 1) {
        usys_log_error("Failed to write register address: %s", strerror(errno));
        close(fd);
        return STATUS_NOK;
    }
    
    // Read data
    if (read(fd, data, len) != (ssize_t)len) {
        usys_log_error("I2C read failed: bus=%d, addr=0x%02X, reg=0x%02X: %s", 
                       bus, device_addr, reg, strerror(errno));
        close(fd);
        return STATUS_NOK;
    }
    
    close(fd);
    usys_log_debug("I2C read successful: bus=%d, addr=0x%02X, reg=0x%02X, len=%zu", 
                   bus, device_addr, reg, len);
    
    return STATUS_OK;
}

int i2c_detect_device(int bus, uint8_t device_addr) {
    char device_path[32];
    int fd;
    
    snprintf(device_path, sizeof(device_path), "/dev/i2c-%d", bus);
    
    fd = open(device_path, O_RDWR);
    if (fd < 0) {
        usys_log_debug("Failed to open I2C device %s: %s", device_path, strerror(errno));
        return STATUS_NOK;
    }
    
    if (ioctl(fd, I2C_SLAVE, device_addr) < 0) {
        close(fd);
        return STATUS_NOK;
    }
    
    // Try a no-data write to detect device presence
    // This is more reliable than a read as it doesn't block on devices
    // that don't support reads or require specific command sequences
    int result = (write(fd, NULL, 0) >= 0) ? STATUS_OK : STATUS_NOK;
    
    close(fd);
    return result;
}

int dac_init(I2CController *controller, FemUnit unit) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    
    if (i2c_detect_device(bus, I2C_ADDR_DAC_AD5667) != STATUS_OK) {
        usys_log_error("DAC AD5667 not detected on bus %d", bus);
        return STATUS_NOK;
    }
    
    uint8_t reset_data[] = {0x06, 0x00};
    if (i2c_write_bytes(bus, I2C_ADDR_DAC_AD5667, 0x40, reset_data, 2) != STATUS_OK) {
        usys_log_error("DAC reset failed");
        return STATUS_NOK;
    }
    
    usleep(10000); // 10ms delay after reset
    
    controller->dac_state.initialized = true;
    usys_log_info("DAC initialized for FEM%d", unit);
    
    return STATUS_OK;
}

uint16_t voltage_to_dac_value(float voltage) {
    float compensated_voltage = voltage / 2.0f;
    uint16_t dac_value = (uint16_t)((compensated_voltage / DAC_VREF) * 65535.0f);
    
    if (dac_value > 65535) dac_value = 65535;
    
    return dac_value;
}

float dac_value_to_voltage(uint16_t dac_value) {
    return ((float)dac_value / 65535.0f) * DAC_VREF * 2.0f; // Account for 2x gain
}

int dac_set_carrier_voltage(I2CController *controller, FemUnit unit, float voltage) {
    if (!controller || !controller->initialized || !controller->dac_state.initialized) {
        return STATUS_NOK;
    }
    
    if (voltage < 0.0f || voltage > DAC_MAX_CARRIER_VOLTAGE) {
        usys_log_error("Carrier voltage out of range: %.2fV (max: %.2fV)", 
                       voltage, DAC_MAX_CARRIER_VOLTAGE);
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    uint16_t dac_value = voltage_to_dac_value(voltage);
    
    uint8_t data[] = {(dac_value >> 8) & 0xFF, dac_value & 0xFF};
    
    if (i2c_write_bytes(bus, I2C_ADDR_DAC_AD5667, 0x59, data, 2) != STATUS_OK) {
        usys_log_error("Failed to set carrier voltage");
        return STATUS_NOK;
    }
    
    controller->dac_state.carrier_voltage = voltage;
    usys_log_info("Carrier voltage set to %.2fV for FEM%d", voltage, unit);
    
    return STATUS_OK;
}

int dac_set_peak_voltage(I2CController *controller, FemUnit unit, float voltage) {
    if (!controller || !controller->initialized || !controller->dac_state.initialized) {
        return STATUS_NOK;
    }
    
    if (voltage < 0.0f || voltage > DAC_MAX_PEAK_VOLTAGE) {
        usys_log_error("Peak voltage out of range: %.2fV (max: %.2fV)", 
                       voltage, DAC_MAX_PEAK_VOLTAGE);
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    uint16_t dac_value = voltage_to_dac_value(voltage);
    
    uint8_t data[] = {(dac_value >> 8) & 0xFF, dac_value & 0xFF};
    
    if (i2c_write_bytes(bus, I2C_ADDR_DAC_AD5667, 0x58, data, 2) != STATUS_OK) {
        usys_log_error("Failed to set peak voltage");
        return STATUS_NOK;
    }
    
    controller->dac_state.peak_voltage = voltage;
    usys_log_info("Peak voltage set to %.2fV for FEM%d", voltage, unit);
    
    return STATUS_OK;
}

int dac_get_config(I2CController *controller, float *carrier, float *peak) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }
    
    if (carrier) *carrier = controller->dac_state.carrier_voltage;
    if (peak) *peak = controller->dac_state.peak_voltage;
    
    return STATUS_OK;
}

int dac_disable_pa(I2CController *controller, FemUnit unit) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }
    
    int result1 = dac_set_carrier_voltage(controller, unit, 0.0f);
    int result2 = dac_set_peak_voltage(controller, unit, 0.0f);
    
    if (result1 == STATUS_OK && result2 == STATUS_OK) {
        usys_log_info("PA disabled - DAC values set to zero for FEM%d", unit);
        return STATUS_OK;
    }
    
    return STATUS_NOK;
}

float lm75a_raw_to_celsius(uint8_t msb, uint8_t lsb) {
    uint16_t temp_raw = (msb << 8) | lsb;
    
    int16_t temp_9bit = temp_raw >> 7;
    
    if (temp_9bit > 255) {
        temp_9bit = temp_9bit - 512;  // Subtract 2^9 for two's complement
    }
    
    return temp_9bit * 0.5f;
}

uint16_t celsius_to_lm75a_raw(float temperature) {
    int16_t temp_9bit = (int16_t)(temperature / 0.5f);
    
    if (temp_9bit < 0) {
        temp_9bit = temp_9bit + 512;  // Add 2^9 for two's complement
    }
    
    return (temp_9bit << 7) & 0xFF80;
}

int temp_sensor_init(I2CController *controller, FemUnit unit) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    
    if (i2c_detect_device(bus, I2C_ADDR_TEMP_LM75A) != STATUS_OK) {
        usys_log_error("Temperature sensor LM75A not detected on bus %d", bus);
        return STATUS_NOK;
    }
    
    usys_log_info("Temperature sensor initialized for FEM%d", unit);
    return STATUS_OK;
}

int temp_sensor_read(I2CController *controller, FemUnit unit, float *temperature) {
    if (!controller || !controller->initialized || !temperature) {
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    uint8_t data[2];
    
    if (i2c_read_bytes(bus, I2C_ADDR_TEMP_LM75A, 0x00, data, 2) != STATUS_OK) {
        usys_log_error("Failed to read temperature from FEM%d", unit);
        return STATUS_NOK;
    }
    
    *temperature = lm75a_raw_to_celsius(data[0], data[1]);
    controller->temp_state.temperature = *temperature;
    
    usys_log_debug("Temperature read: %.1f°C from FEM%d", *temperature, unit);
    return STATUS_OK;
}

int temp_sensor_set_threshold(I2CController *controller, FemUnit unit, float threshold) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }
    
    if (threshold < TEMP_MIN || threshold > TEMP_MAX) {
        usys_log_error("Temperature threshold out of range: %.1f°C", threshold);
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    uint16_t threshold_raw = celsius_to_lm75a_raw(threshold);
    
    uint8_t data[] = {(threshold_raw >> 8) & 0xFF, threshold_raw & 0xFF};
    
    if (i2c_write_bytes(bus, I2C_ADDR_TEMP_LM75A, 0x03, data, 2) != STATUS_OK) {
        usys_log_error("Failed to set temperature threshold");
        return STATUS_NOK;
    }
    
    controller->temp_state.threshold = threshold;
    controller->temp_state.alert_enabled = true;
    usys_log_info("Temperature threshold set to %.1f°C for FEM%d", threshold, unit);
    
    return STATUS_OK;
}

const I2CDeviceInfo* i2c_get_device_info(I2CDevice device) {
    if (device >= I2C_DEVICE_MAX) {
        return NULL;
    }
    return &device_info[device];
}

int adc_init(I2CController *controller, FemUnit unit) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    
    if (i2c_detect_device(bus, I2C_ADDR_ADC_ADS1015) != STATUS_OK) {
        usys_log_error("ADC ADS1015 not detected on bus %d", bus);
        return STATUS_NOK;
    }
    
    usys_log_info("ADC initialized for FEM%d", unit);
    return STATUS_OK;
}

static int adc_configure_channel(int bus, int channel) {
    uint16_t mux_configs[] = {0x4000, 0x5000, 0x6000, 0x7000}; // AIN0-3 vs GND
    
    if (channel < 0 || channel > 3) {
        return STATUS_NOK;
    }
    
    uint16_t config = 0x8000;  // Start conversion
    config |= mux_configs[channel];
    config |= 0x0200;  // PGA ±4.096V
    config |= 0x0100;  // Single-shot mode
    config |= 0x0080;  // 1600 SPS
    config |= 0x0003;  // Disable comparator
    
    uint8_t data[] = {(config >> 8) & 0xFF, config & 0xFF};
    
    return i2c_write_bytes(bus, I2C_ADDR_ADC_ADS1015, 0x01, data, 2);
}

float adc_raw_to_voltage(uint16_t raw_value) {
    return ((int16_t)raw_value >> 4) * 4.096f / 2048.0f;
}

int adc_read_channel(I2CController *controller, FemUnit unit, int channel, float *voltage) {
    if (!controller || !controller->initialized || !voltage) {
        return STATUS_NOK;
    }
    
    if (channel < 0 || channel > 3) {
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    
    if (adc_configure_channel(bus, channel) != STATUS_OK) {
        return STATUS_NOK;
    }
    
    usleep(10000); // 10ms delay for conversion
    
    uint8_t data[2];
    if (i2c_read_bytes(bus, I2C_ADDR_ADC_ADS1015, 0x00, data, 2) != STATUS_OK) {
        usys_log_error("Failed to read ADC channel %d from FEM%d", channel, unit);
        return STATUS_NOK;
    }
    
    uint16_t adc_raw = (data[0] << 8) | data[1];
    *voltage = adc_raw_to_voltage(adc_raw);
    
    usys_log_debug("ADC channel %d: %.3fV from FEM%d", channel, *voltage, unit);
    return STATUS_OK;
}

float voltage_to_reverse_power(float voltage) {
    return (voltage - 2.0f) * 20.0f - 30.0f;
}

float voltage_to_current(float voltage) {
    return voltage; // Example: 1V = 1A
}

int adc_read_reverse_power(I2CController *controller, FemUnit unit, float *power_dbm) {
    if (!power_dbm) {
        return STATUS_NOK;
    }
    
    float voltage;
    if (adc_read_channel(controller, unit, ADC_CHANNEL_REVERSE_POWER, &voltage) != STATUS_OK) {
        return STATUS_NOK;
    }
    
    *power_dbm = voltage_to_reverse_power(voltage);
    controller->adc_state.reverse_power_dbm = *power_dbm;
    
    return STATUS_OK;
}

int adc_read_pa_current(I2CController *controller, FemUnit unit, float *current_a) {
    if (!current_a) {
        return STATUS_NOK;
    }
    
    float voltage;
    if (adc_read_channel(controller, unit, ADC_CHANNEL_PA_CURRENT, &voltage) != STATUS_OK) {
        return STATUS_NOK;
    }
    
    *current_a = voltage_to_current(voltage);
    controller->adc_state.pa_current_a = *current_a;
    
    return STATUS_OK;
}

int adc_read_all_channels(I2CController *controller, FemUnit unit) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }
    
    usys_log_info("Reading all ADC channels for FEM%d:", unit);
    
    struct {
        int channel;
        const char *name;
        const char *unit;
        float (*convert)(float);
    } channels[] = {
        {0, "Reverse Power", "dBm", voltage_to_reverse_power},
        {1, "Forward Power", "dBm", voltage_to_reverse_power}, // Adjust conversion as needed
        {2, "PA Current", "A", voltage_to_current},
        {3, "Temperature", "°C", NULL}  // Direct voltage reading
    };
    
    for (int i = 0; i < 4; i++) {
        float voltage;
        if (adc_read_channel(controller, unit, channels[i].channel, &voltage) == STATUS_OK) {
            if (channels[i].convert) {
                float converted = channels[i].convert(voltage);
                usys_log_info("  Channel %d (%s): %.2f %s", 
                              channels[i].channel, channels[i].name, converted, channels[i].unit);
            } else {
                usys_log_info("  Channel %d (%s): %.3f V", 
                              channels[i].channel, channels[i].name, voltage);
            }
        } else {
            usys_log_error("  Channel %d (%s): Read error", 
                           channels[i].channel, channels[i].name);
        }
    }
    
    return STATUS_OK;
}

int adc_set_safety_thresholds(I2CController *controller, float max_reverse_power, float max_current) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }
    
    controller->adc_state.max_reverse_power = max_reverse_power;
    controller->adc_state.max_current = max_current;
    controller->adc_state.safety_enabled = true;
    
    usys_log_info("Safety thresholds set: reverse power %.1f dBm, current %.1f A", 
                  max_reverse_power, max_current);
    
    return STATUS_OK;
}

int adc_check_safety(I2CController *controller, FemUnit unit, bool *safety_violation) {
    if (!controller || !controller->initialized || !safety_violation) {
        return STATUS_NOK;
    }
    
    *safety_violation = false;
    
    if (!controller->adc_state.safety_enabled) {
        return STATUS_OK;
    }
    
    float reverse_power;
    if (adc_read_reverse_power(controller, unit, &reverse_power) == STATUS_OK) {
        if (reverse_power > controller->adc_state.max_reverse_power) {
            usys_log_warn("Safety violation: reverse power %.1f dBm exceeds threshold %.1f dBm", 
                          reverse_power, controller->adc_state.max_reverse_power);
            *safety_violation = true;
        }
    }
    
    float pa_current;
    if (adc_read_pa_current(controller, unit, &pa_current) == STATUS_OK) {
        if (pa_current > controller->adc_state.max_current) {
            usys_log_warn("Safety violation: PA current %.1f A exceeds threshold %.1f A", 
                          pa_current, controller->adc_state.max_current);
            *safety_violation = true;
        }
    }
    
    return STATUS_OK;
}

int eeprom_write_serial(I2CController *controller, FemUnit unit, const char *serial) {
    if (!controller || !controller->initialized || !serial) {
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    size_t len = strlen(serial);
    
    if (len > 16) {
        usys_log_error("Serial number too long (max 16 characters)");
        return STATUS_NOK;
    }
    
    for (size_t i = 0; i < len; i++) {
        uint8_t data = (uint8_t)serial[i];
        if (i2c_write_bytes(bus, I2C_ADDR_EEPROM, (uint8_t)i, &data, 1) != STATUS_OK) {
            usys_log_error("Failed to write EEPROM at position %zu", i);
            return STATUS_NOK;
        }
        usleep(10000); // 10ms delay for EEPROM write
    }
    
    uint8_t null_term = 0x00;
    i2c_write_bytes(bus, I2C_ADDR_EEPROM, (uint8_t)len, &null_term, 1);
    
    strncpy(controller->eeprom_state.serial_number, serial, sizeof(controller->eeprom_state.serial_number) - 1);
    controller->eeprom_state.has_data = true;
    
    usys_log_info("Serial number written to EEPROM: %s", serial);
    return STATUS_OK;
}

int eeprom_read_serial(I2CController *controller, FemUnit unit, char *serial, size_t max_len) {
    if (!controller || !controller->initialized || !serial || max_len == 0) {
        return STATUS_NOK;
    }
    
    int bus = i2c_get_bus_for_fem(unit);
    
    size_t read_len = (max_len - 1 < 16) ? max_len - 1 : 16;
    
    for (size_t i = 0; i < read_len; i++) {
        uint8_t data;
        if (i2c_read_bytes(bus, I2C_ADDR_EEPROM, (uint8_t)i, &data, 1) != STATUS_OK) {
            usys_log_error("Failed to read EEPROM at position %zu", i);
            return STATUS_NOK;
        }
        
        if (data == 0) {
            serial[i] = '\0';
            break;
        }
        
        serial[i] = (char)data;
        
        if (i == read_len - 1) {
            serial[i + 1] = '\0';
        }
    }
    
    if (strlen(serial) > 0) {
        strncpy(controller->eeprom_state.serial_number, serial, sizeof(controller->eeprom_state.serial_number) - 1);
        controller->eeprom_state.has_data = true;
        usys_log_info("Serial number read from EEPROM: %s", serial);
        return STATUS_OK;
    } else {
        usys_log_info("No serial number found in EEPROM");
        return STATUS_NOK;
    }
}

void i2c_print_device_scan(FemUnit unit) {
    int bus = i2c_get_bus_for_fem(unit);
    
    usys_log_info("I2C Device Scan for FEM%d (bus %d):", unit, bus);
    
    for (int i = 0; i < I2C_DEVICE_MAX; i++) {
        const I2CDeviceInfo *info = &device_info[i];
        int detected = i2c_detect_device(bus, info->address);
        
        usys_log_info("  %s (0x%02X): %s - %s", 
                      info->name, info->address, info->description,
                      detected == STATUS_OK ? "DETECTED" : "NOT FOUND");
    }
}