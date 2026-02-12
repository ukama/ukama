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
#include <stddef.h>

#include "femd.h"
#include "i2c_bus.h"
#include "gpio_controller.h"

#define I2C_ADDR_DAC_AD5667   0x0C
#define I2C_ADDR_TEMP_LM75A   0x49
#define I2C_ADDR_ADC_ADS1015  0x48
#define I2C_ADDR_EEPROM       0x50

/* Controller temp sensor (bus0) */
#define I2C_ADDR_CTRL_TMP10X  0x48

#define DAC_VREF                2.5f
#define DAC_MAX_CARRIER_VOLTAGE 5.0f
#define DAC_MAX_PEAK_VOLTAGE    5.0f

#define TEMP_MIN  -40.0f
#define TEMP_MAX  125.0f

typedef enum {
    I2C_DEVICE_DAC = 0,
    I2C_DEVICE_TEMP,
    I2C_DEVICE_ADC,
    I2C_DEVICE_EEPROM,
    I2C_DEVICE_MAX
} I2CDevice;

typedef struct {
    const char *name;
    uint8_t     address;
    const char *description;
} I2CDeviceInfo;

/* ADC channels xxx */
typedef enum {
    ADC_CHANNEL_REVERSE_POWER = 0,
    ADC_CHANNEL_FORWARD_POWER = 1,
    ADC_CHANNEL_PA_CURRENT    = 2,
    ADC_CHANNEL_TEMPERATURE   = 3,
    ADC_CHANNEL_MAX
} ADCChannel;

const I2CDeviceInfo* i2c_get_device_info(I2CDevice device);
void     i2c_print_device_scan(I2cBus *bus, FemUnit unit);
int      dac_init(I2cBus *bus, FemUnit unit);
int      dac_set_carrier_voltage(I2cBus *bus, FemUnit unit, float voltage);
int      dac_set_peak_voltage(I2cBus *bus, FemUnit unit, float voltage);
int      dac_disable_pa(I2cBus *bus, FemUnit unit);
uint16_t voltage_to_dac_value(float voltage);
float    dac_value_to_voltage(uint16_t dacValue);

int      temp_sensor_init(I2cBus *bus, FemUnit unit);
int      temp_sensor_read(I2cBus *bus, FemUnit unit, float *temperature);
int      temp_sensor_set_threshold(I2cBus *bus, FemUnit unit, float threshold);
float    lm75a_raw_to_celsius(uint8_t msb, uint8_t lsb);
uint16_t celsius_to_lm75a_raw(float temperature);

int      adc_init(I2cBus *bus, FemUnit unit);
int      adc_read_channel(I2cBus *bus, FemUnit unit, ADCChannel channel, float *voltage);
int      adc_read_reverse_power(I2cBus *bus, FemUnit unit, float *powerDbm);
int      adc_read_pa_current(I2cBus *bus, FemUnit unit, float *currentA);
int      adc_read_all_channels(I2cBus *bus, FemUnit unit);
float    adc_raw_to_voltage(uint16_t rawValue);
float    voltage_to_reverse_power(float voltage);
float    voltage_to_current(float voltage);

int      eeprom_write_serial(I2cBus *bus, FemUnit unit, const char *serial);
int      eeprom_read_serial(I2cBus *bus, FemUnit unit, char *serial, size_t maxLen);

int      ctrl_temp_read_tmp10x(I2cBus *bus, float *temperatureC);

#endif /* I2C_CONTROLLER_H */
