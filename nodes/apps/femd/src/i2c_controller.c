/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <unistd.h>
#include <math.h>

#include "usys_log.h"

#include "i2c_controller.h"

static const I2CDeviceInfo deviceInfo[I2C_DEVICE_MAX] = {
    [I2C_DEVICE_DAC]    = { "AD5667",  I2C_ADDR_DAC_AD5667,  "16-bit DAC" },
    [I2C_DEVICE_TEMP]   = { "LM75A",   I2C_ADDR_TEMP_LM75A,  "Temperature sensor" },
    [I2C_DEVICE_ADC]    = { "ADS1015", I2C_ADDR_ADC_ADS1015, "12-bit ADC" },
    [I2C_DEVICE_EEPROM] = { "EEPROM",  I2C_ADDR_EEPROM,      "Serial number storage" }
};

const I2CDeviceInfo* i2c_get_device_info(I2CDevice device) {
    if (device >= I2C_DEVICE_MAX) return NULL;
    return &deviceInfo[device];
}

void i2c_print_device_scan(I2cBus *bus, const char *label) {
    if (!bus || !bus->initialized) return;

    usys_log_info("I2C scan: %s (bus=%d)", label ? label : "bus", bus->busNum);
    for (int i = 0; i < I2C_DEVICE_MAX; i++) {
        const I2CDeviceInfo *info = &deviceInfo[i];
        int ok = i2c_bus_detect(bus, info->address);
        usys_log_info("  %s (0x%02X): %s - %s",
                      info->name, info->address, info->description,
                      ok == STATUS_OK ? "DETECTED" : "NOT FOUND");
    }
}

static uint16_t voltage_to_dac_value(float voltage) {
    float v = voltage / 2.0f;
    float code = (v / DAC_VREF) * 65535.0f;
    if (code < 0.0f) code = 0.0f;
    if (code > 65535.0f) code = 65535.0f;
    return (uint16_t)(code + 0.5f);
}

static float dac_value_to_voltage(uint16_t dacValue) {
    return ((float)dacValue / 65535.0f) * DAC_VREF * 2.0f;
}

int dac_init(I2cBus *bus) {
    uint8_t resetData[2];

    if (!bus || !bus->initialized) return STATUS_NOK;

    if (i2c_bus_detect(bus, I2C_ADDR_DAC_AD5667) != STATUS_OK) {
        usys_log_error("DAC not detected (bus=%d)", bus->busNum);
        return STATUS_NOK;
    }

    resetData[0] = 0x06;
    resetData[1] = 0x00;

    if (i2c_bus_write_reg(bus, I2C_ADDR_DAC_AD5667, 0x40, resetData, 2) != STATUS_OK) {
        usys_log_error("DAC reset failed (bus=%d)", bus->busNum);
        return STATUS_NOK;
    }

    usleep(10000);
    bus->dacKnown = true;
    bus->dacCarrierV = 0.0f;
    bus->dacPeakV = 0.0f;

    return STATUS_OK;
}

static int dac_set_voltage(I2cBus *bus, float voltage, bool isCarrier) {
    uint16_t dacValue;
    uint8_t data[2];
    uint8_t reg;

    if (!bus || !bus->initialized) return STATUS_NOK;

    if (voltage < 0.0f || voltage > 5.0f) {
        usys_log_error("DAC voltage out of range: %.2f", voltage);
        return STATUS_NOK;
    }

    dacValue = voltage_to_dac_value(voltage);
    data[0] = (uint8_t)((dacValue >> 8) & 0xFF);
    data[1] = (uint8_t)(dacValue & 0xFF);
    reg = isCarrier ? 0x59 : 0x58;

    if (i2c_bus_write_reg(bus, I2C_ADDR_DAC_AD5667, reg, data, 2) != STATUS_OK) {
        usys_log_error("DAC write failed (bus=%d)", bus->busNum);
        return STATUS_NOK;
    }

    bus->dacKnown = true;
    if (isCarrier) bus->dacCarrierV = voltage;
    else bus->dacPeakV = voltage;

    return STATUS_OK;
}

int dac_set_carrier_voltage(I2cBus *bus, float voltage) {
    return dac_set_voltage(bus, voltage, true);
}

int dac_set_peak_voltage(I2cBus *bus, float voltage) {
    return dac_set_voltage(bus, voltage, false);
}

int dac_disable_pa(I2cBus *bus) {
    if (dac_set_carrier_voltage(bus, 0.0f) != STATUS_OK) return STATUS_NOK;
    if (dac_set_peak_voltage(bus, 0.0f) != STATUS_OK) return STATUS_NOK;
    return STATUS_OK;
}

int dac_get_cached(I2cBus *bus, float *carrierV, float *peakV) {
    if (!bus || !bus->initialized) return STATUS_NOK;
    if (!bus->dacKnown) return STATUS_NOK;
    if (carrierV) *carrierV = bus->dacCarrierV;
    if (peakV)    *peakV = bus->dacPeakV;
    return STATUS_OK;
}

static float lm75a_raw_to_celsius(uint8_t msb, uint8_t lsb) {
    uint16_t tempRaw = (uint16_t)((msb << 8) | lsb);
    int16_t temp9bit = (int16_t)(tempRaw >> 7);
    if (temp9bit > 255) temp9bit = (int16_t)(temp9bit - 512);
    return (float)temp9bit * 0.5f;
}

static uint16_t celsius_to_lm75a_raw(float temperature) {
    int16_t temp9bit = (int16_t)(temperature / 0.5f);
    if (temp9bit < 0) temp9bit = (int16_t)(temp9bit + 512);
    return (uint16_t)((temp9bit << 7) & 0xFF80);
}

int temp_sensor_init(I2cBus *bus) {
    if (!bus || !bus->initialized) return STATUS_NOK;
    if (i2c_bus_detect(bus, I2C_ADDR_TEMP_LM75A) != STATUS_OK) return STATUS_NOK;
    return STATUS_OK;
}

int temp_sensor_read(I2cBus *bus, float *temperatureC) {
    uint8_t data[2];

    if (!bus || !bus->initialized || !temperatureC) return STATUS_NOK;

    if (i2c_bus_read_reg(bus, I2C_ADDR_TEMP_LM75A, 0x00, data, 2) != STATUS_OK) {
        return STATUS_NOK;
    }

    *temperatureC = lm75a_raw_to_celsius(data[0], data[1]);
    return STATUS_OK;
}

int temp_sensor_set_threshold(I2cBus *bus, float thresholdC) {
    uint16_t raw;
    uint8_t data[2];

    if (!bus || !bus->initialized) return STATUS_NOK;

    raw = celsius_to_lm75a_raw(thresholdC);
    data[0] = (uint8_t)((raw >> 8) & 0xFF);
    data[1] = (uint8_t)(raw & 0xFF);

    if (i2c_bus_write_reg(bus, I2C_ADDR_TEMP_LM75A, 0x03, data, 2) != STATUS_OK) {
        return STATUS_NOK;
    }

    return STATUS_OK;
}

static int adc_configure_channel(I2cBus *bus, ADCChannel channel) {
    uint16_t muxConfigs[ADC_CHANNEL_MAX] = { 0x4000, 0x5000, 0x6000, 0x7000 };
    uint16_t config;
    uint8_t data[2];

    if (!bus || !bus->initialized) return STATUS_NOK;
    if (channel < 0 || channel >= ADC_CHANNEL_MAX) return STATUS_NOK;

    config  = 0x8000;
    config |= muxConfigs[(int)channel];
    config |= 0x0200;
    config |= 0x0100;
    config |= 0x0080;
    config |= 0x0003;

    data[0] = (uint8_t)((config >> 8) & 0xFF);
    data[1] = (uint8_t)(config & 0xFF);

    return i2c_bus_write_reg(bus, I2C_ADDR_ADC_ADS1015, 0x01, data, 2);
}

static float adc_raw_to_voltage(uint16_t rawValue) {
    return ((int16_t)rawValue >> 4) * 4.096f / 2048.0f;
}

int adc_init(I2cBus *bus) {
    if (!bus || !bus->initialized) return STATUS_NOK;
    if (i2c_bus_detect(bus, I2C_ADDR_ADC_ADS1015) != STATUS_OK) return STATUS_NOK;
    return STATUS_OK;
}

int adc_read_channel(I2cBus *bus, ADCChannel channel, float *voltage) {
    uint8_t data[2];
    uint16_t adcRaw;

    if (!bus || !bus->initialized || !voltage) return STATUS_NOK;
    if (channel < 0 || channel >= ADC_CHANNEL_MAX) return STATUS_NOK;

    if (adc_configure_channel(bus, channel) != STATUS_OK) return STATUS_NOK;

    usleep(10000);

    if (i2c_bus_read_reg(bus, I2C_ADDR_ADC_ADS1015, 0x00, data, 2) != STATUS_OK) return STATUS_NOK;

    adcRaw = (uint16_t)((data[0] << 8) | data[1]);
    *voltage = adc_raw_to_voltage(adcRaw);

    return STATUS_OK;
}

static float voltage_to_reverse_power(float voltage) {
    return (voltage - 2.0f) * 20.0f - 30.0f;
}

static float voltage_to_current(float voltage) {
    return voltage;
}

int adc_read_reverse_power(I2cBus *bus, float *powerDbm) {
    float v;
    if (!powerDbm) return STATUS_NOK;
    if (adc_read_channel(bus, ADC_CHANNEL_REVERSE_POWER, &v) != STATUS_OK) return STATUS_NOK;
    *powerDbm = voltage_to_reverse_power(v);
    return STATUS_OK;
}

int adc_read_forward_power(I2cBus *bus, float *powerDbm) {
    float v;
    if (!powerDbm) return STATUS_NOK;
    if (adc_read_channel(bus, ADC_CHANNEL_FORWARD_POWER, &v) != STATUS_OK) return STATUS_NOK;
    *powerDbm = voltage_to_reverse_power(v);
    return STATUS_OK;
}

int adc_read_pa_current(I2cBus *bus, float *currentA) {
    float v;
    if (!currentA) return STATUS_NOK;
    if (adc_read_channel(bus, ADC_CHANNEL_PA_CURRENT, &v) != STATUS_OK) return STATUS_NOK;
    *currentA = voltage_to_current(v);
    return STATUS_OK;
}

int ctrl_temp_read_tmp10x(I2cBus *bus, float *tempC) {
    uint8_t data[2];
    int16_t raw;

    if (!bus || !bus->initialized || !tempC) return STATUS_NOK;

    if (i2c_bus_read_reg(bus, I2C_ADDR_CTRL_TEMP, 0x00, data, 2) != STATUS_OK) return STATUS_NOK;

    raw = (int16_t)((data[0] << 8) | data[1]);
    raw >>= 4;
    *tempC = (float)raw * 0.0625f;

    return STATUS_OK;
}
