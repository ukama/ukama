/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef I2C_CONTROLLER_H
#define I2C_CONTROLLER_H

#include <stdint.h>
#include <stdbool.h>
#include "gpio_controller.h"

#define STATUS_OK                 0
#define STATUS_NOK               -1

#define I2C_BUS_FEM1              1
#define I2C_BUS_FEM2              2

// I2C Device addresses
#define I2C_ADDR_DAC_AD5667       0x0C  // 16-bit DAC
#define I2C_ADDR_TEMP_LM75A       0x49  // Temperature sensor
#define I2C_ADDR_ADC_ADS1015      0x48  // 12-bit ADC
#define I2C_ADDR_EEPROM           0x50  // Serial number storage

// DAC voltage limits
#define DAC_MAX_CARRIER_VOLTAGE   3.1f
#define DAC_MAX_PEAK_VOLTAGE      2.5f
#define DAC_VREF                  2.5f

// Temperature limits
#define TEMP_MIN                  -55.0f
#define TEMP_MAX                  125.0f

// ADC channels
#define ADC_CHANNEL_REVERSE_POWER 0
#define ADC_CHANNEL_FORWARD_POWER 1
#define ADC_CHANNEL_PA_CURRENT    2
#define ADC_CHANNEL_TEMPERATURE   3

typedef enum {
    I2C_DEVICE_DAC = 0,
    I2C_DEVICE_TEMP,
    I2C_DEVICE_ADC,
    I2C_DEVICE_EEPROM,
    I2C_DEVICE_MAX
} I2CDevice;

typedef struct {
    char name[32];
    uint8_t address;
    char description[64];
} I2CDeviceInfo;

typedef struct {
    float carrier_voltage;
    float peak_voltage;
    bool initialized;
} DACState;

typedef struct {
    float temperature;
    float threshold;
    bool alert_enabled;
} TempSensorState;

typedef struct {
    float reverse_power_dbm;
    float forward_power_dbm;
    float pa_current_a;
    float temperature_c;
    float max_reverse_power;
    float max_current;
    bool safety_enabled;
} ADCState;

typedef struct {
    char serial_number[17];  // 16 chars + null terminator
    bool has_data;
} EEPROMState;

typedef struct {
    int bus_fem1;
    int bus_fem2;
    DACState dac_state;
    TempSensorState temp_state;
    ADCState adc_state;
    EEPROMState eeprom_state;
    bool initialized;
} I2CController;

// Controller management
int i2c_controller_init(I2CController *controller);
void i2c_controller_cleanup(I2CController *controller);
int i2c_get_bus_for_fem(FemUnit unit);

// Low-level I2C operations
int i2c_write_bytes(int bus, uint8_t device_addr, uint8_t reg, const uint8_t *data, size_t len);
int i2c_read_bytes(int bus, uint8_t device_addr, uint8_t reg, uint8_t *data, size_t len);
int i2c_detect_device(int bus, uint8_t device_addr);

// DAC operations (AD5667)
int dac_init(I2CController *controller, FemUnit unit);
int dac_set_carrier_voltage(I2CController *controller, FemUnit unit, float voltage);
int dac_set_peak_voltage(I2CController *controller, FemUnit unit, float voltage);
int dac_get_config(I2CController *controller, float *carrier, float *peak);
int dac_disable_pa(I2CController *controller, FemUnit unit);

// Temperature sensor operations (LM75A)
int temp_sensor_init(I2CController *controller, FemUnit unit);
int temp_sensor_read(I2CController *controller, FemUnit unit, float *temperature);
int temp_sensor_set_threshold(I2CController *controller, FemUnit unit, float threshold);
int temp_sensor_check_alert(I2CController *controller, FemUnit unit, bool *alert);

// ADC operations (ADS1015)
int adc_init(I2CController *controller, FemUnit unit);
int adc_read_channel(I2CController *controller, FemUnit unit, int channel, float *voltage);
int adc_read_reverse_power(I2CController *controller, FemUnit unit, float *power_dbm);
int adc_read_pa_current(I2CController *controller, FemUnit unit, float *current_a);
int adc_read_all_channels(I2CController *controller, FemUnit unit);
int adc_set_safety_thresholds(I2CController *controller, float max_reverse_power, float max_current);
int adc_check_safety(I2CController *controller, FemUnit unit, bool *safety_violation);

// EEPROM operations
int eeprom_write_serial(I2CController *controller, FemUnit unit, const char *serial);
int eeprom_read_serial(I2CController *controller, FemUnit unit, char *serial, size_t max_len);

// Utility functions
uint16_t voltage_to_dac_value(float voltage);
float dac_value_to_voltage(uint16_t dac_value);
float lm75a_raw_to_celsius(uint8_t msb, uint8_t lsb);
uint16_t celsius_to_lm75a_raw(float temperature);
float adc_raw_to_voltage(uint16_t raw_value);
float voltage_to_reverse_power(float voltage);
float voltage_to_current(float voltage);

// Device information
const I2CDeviceInfo* i2c_get_device_info(I2CDevice device);
void i2c_print_device_scan(FemUnit unit);

#endif /* I2C_CONTROLLER_H */