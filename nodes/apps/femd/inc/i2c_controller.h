/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */
#ifndef I2C_CONTROLLER_H
#define I2C_CONTROLLER_H

#include <stdint.h>
#include <stdbool.h>

#include "gpio_controller.h"

/* FEM Bus IDs */
#define I2C_BUS_FEM1  1
#define I2C_BUS_FEM2  2

/* Device addresses */
#define I2C_ADDR_DAC_AD5667   0x0C
#define I2C_ADDR_TEMP_LM75A   0x49
#define I2C_ADDR_ADC_ADS1015  0x48
#define I2C_ADDR_EEPROM       0x50

/* Limits */
#define DAC_MAX_CARRIER_VOLTAGE  3.1f
#define DAC_MAX_PEAK_VOLTAGE     2.5f
#define DAC_VREF                 2.5f
#define TEMP_MIN               -55.0f
#define TEMP_MAX               125.0f

/* ADC channels */
typedef enum {
    ADC_CHANNEL_REVERSE_POWER = 0,
    ADC_CHANNEL_FORWARD_POWER,
    ADC_CHANNEL_PA_CURRENT,
    ADC_CHANNEL_TEMPERATURE,
    ADC_CHANNEL_MAX
} ADCChannel;

typedef enum {
    I2C_DEVICE_DAC = 0,
    I2C_DEVICE_TEMP,
    I2C_DEVICE_ADC,
    I2C_DEVICE_EEPROM,
    I2C_DEVICE_MAX
} I2CDevice;

typedef struct {
    const char *name;
    uint8_t    address;
    const char *description;
} I2CDeviceInfo;

/* State structs */
typedef struct {
    float carrierVoltage;
    float peakVoltage;
    bool  initialized;
} DacState;

typedef struct {
    float temperature;
    float threshold;
    bool  alertEnabled;
} TempSensorState;

typedef struct {
    float reversePowerDbm;
    float forwardPowerDbm;
    float paCurrentA;
    float temperatureC;
    float maxReversePower;
    float maxCurrent;
    bool  safetyEnabled;
} ADCState;

typedef struct {
    char serialNumber[17];
    bool hasData;
} EEPROMState;

typedef struct {
    int busFem1;
    int busFem2;
    DacState dacState;
    TempSensorState tempState;
    ADCState adcState;
    EEPROMState eepromState;
    bool initialized;
} I2CController;

/* Controller lifecycle */
int i2c_controller_init(I2CController *controller);
void i2c_controller_cleanup(I2CController *controller);
int i2c_get_bus_for_fem(FemUnit unit);

/* Generic helpers */
int i2c_write_bytes(int bus, uint8_t devAddr, uint8_t reg, const uint8_t *data, size_t len);
int i2c_read_bytes(int bus,  uint8_t devAddr, uint8_t reg, uint8_t *data, size_t len);
int i2c_detect_device(int bus, uint8_t devAddr);
const I2CDeviceInfo* i2c_get_device_info(I2CDevice device);
void i2c_print_device_scan(FemUnit unit);

/* DAC */
int dac_init(I2CController *controller, FemUnit unit);
int dac_set_carrier_voltage(I2CController *controller, FemUnit unit, float voltage);
int dac_set_peak_voltage(I2CController *controller, FemUnit unit, float voltage);
int dac_get_config(I2CController *controller, float *carrier, float *peak);
int dac_disable_pa(I2CController *controller, FemUnit unit);
uint16_t voltage_to_dac_value(float voltage);
float dac_value_to_voltage(uint16_t dacValue);

/* Temp Sensor */
int temp_sensor_init(I2CController *controller, FemUnit unit);
int temp_sensor_read(I2CController *controller, FemUnit unit, float *temperature);
int temp_sensor_set_threshold(I2CController *controller, FemUnit unit, float threshold);

float lm75a_raw_to_celsius(uint8_t msb, uint8_t lsb);
uint16_t celsius_to_lm75a_raw(float temperature);

/* ADC */
int adc_init(I2CController *controller, FemUnit unit);
int adc_read_channel(I2CController *controller, FemUnit unit, ADCChannel channel, float *voltage);
int adc_read_reverse_power(I2CController *controller, FemUnit unit, float *powerDbm);
int adc_read_pa_current(I2CController *controller, FemUnit unit, float *currentA);
int adc_read_all_channels(I2CController *controller, FemUnit unit);
int adc_set_safety_thresholds(I2CController *controller, float maxReversePower, float maxCurrent);
int adc_check_safety(I2CController *controller, FemUnit unit, bool *safetyViolation);

float adc_raw_to_voltage(uint16_t rawValue);
float voltage_to_reverse_power(float voltage);
float voltage_to_current(float voltage);

/* EEPROM */
int eeprom_write_serial(I2CController *controller, FemUnit unit, const char *serial);
int eeprom_read_serial(I2CController *controller, FemUnit unit, char *serial, size_t maxLen);

#endif /* I2C_CONTROLLER_H */
