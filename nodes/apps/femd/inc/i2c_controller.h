/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef I2C_CONTROLLER_H
#define I2C_CONTROLLER_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

#include "femd.h"
#include "i2c_bus.h"

#define I2C_ADDR_DAC_AD5667   0x0C
#define I2C_ADDR_TEMP_LM75A   0x49
#define I2C_ADDR_ADC_ADS1015  0x48
#define I2C_ADDR_EEPROM       0x50
#define I2C_ADDR_CTRL_TEMP    0x48

#define DAC_VREF  2.5f

typedef enum {
    ADC_CHANNEL_REVERSE_POWER = 0,
    ADC_CHANNEL_FORWARD_POWER = 1,
    ADC_CHANNEL_PA_CURRENT    = 2,
    ADC_CHANNEL_TEMPERATURE   = 3,
    ADC_CHANNEL_MAX
} ADCChannel;

typedef struct {
    const char *name;
    uint8_t     address;
    const char *description;
} I2CDeviceInfo;

typedef enum {
    I2C_DEVICE_DAC = 0,
    I2C_DEVICE_TEMP,
    I2C_DEVICE_ADC,
    I2C_DEVICE_EEPROM,
    I2C_DEVICE_MAX
} I2CDevice;

const I2CDeviceInfo* i2c_get_device_info(I2CDevice device);
void i2c_print_device_scan(I2cBus *bus, const char *label);

/* DAC (cached in I2cBus; AD5667 is effectively write-only in this design) */
int  dac_init(I2cBus *bus);
int  dac_set_carrier_voltage(I2cBus *bus, float voltage);
int  dac_set_peak_voltage(I2cBus *bus, float voltage);
int  dac_disable_pa(I2cBus *bus);
int  dac_get_cached(I2cBus *bus, float *carrierV, float *peakV);

/* Temp sensor (LM75A) */
int  temp_sensor_init(I2cBus *bus);
int  temp_sensor_read(I2cBus *bus, float *temperatureC);
int  temp_sensor_set_threshold(I2cBus *bus, float thresholdC);

/* ADC (ADS1015) */
int  adc_init(I2cBus *bus);
int  adc_read_channel(I2cBus *bus, ADCChannel channel, float *voltage);
int  adc_read_reverse_power(I2cBus *bus, float *powerDbm);
int  adc_read_forward_power(I2cBus *bus, float *powerDbm);
int  adc_read_pa_current(I2cBus *bus, float *currentA);

/* Controller temp sensor (TMP10x @ 0x48) */
int  ctrl_temp_read_tmp10x(I2cBus *bus, float *tempC);

#endif /* I2C_CONTROLLER_H */
