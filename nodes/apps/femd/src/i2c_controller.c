
/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include <string.h>
#include <unistd.h>

#include "i2c_controller.h"
#include "usys_log.h"

static const I2CDeviceInfo deviceInfo[I2C_DEVICE_MAX] = {
    {"AD5667",  I2C_ADDR_DAC_AD5667,  "16-bit DAC"},
    {"LM75A",   I2C_ADDR_TEMP_LM75A,  "Temperature sensor"},
    {"ADS1015", I2C_ADDR_ADC_ADS1015, "12-bit ADC"},
    {"EEPROM",  I2C_ADDR_EEPROM,      "Serial number storage"}
};

/* Forward statics */
static int adc_configure_channel(I2cBus *bus, ADCChannel channel);
static int dac_set_voltage(I2cBus *bus, FemUnit unit, float voltage, bool isCarrier);

const I2CDeviceInfo* i2c_get_device_info(I2CDevice device) {

    if (device >= I2C_DEVICE_MAX) return NULL;
    return &deviceInfo[device];
}

void i2c_print_device_scan(I2cBus *bus, FemUnit unit) {

    int detected = STATUS_NOK;

    if (!bus || !bus->initialized) {
        usys_log_error("i2c_print_device_scan: bus not initialized");
        return;
    }

    usys_log_info("I2C Device Scan for FEM%d (bus %d):", unit, bus->busNum);

    for (int i = 0; i < I2C_DEVICE_MAX; i++) {
        const I2CDeviceInfo *info = &deviceInfo[i];
        detected = i2c_bus_detect(bus, info->address);
        usys_log_info("  %s (0x%02X): %s - %s",
                      info->name, info->address, info->description,
                      detected == STATUS_OK ? "DETECTED" : "NOT FOUND");
    }
}

/* DAC (AD5667) */

uint16_t voltage_to_dac_value(float voltage) {

    float compensatedVoltage = voltage / 2.0f;
    uint32_t dacValue = (uint32_t)((compensatedVoltage / DAC_VREF) * 65535.0f);

    if (dacValue > 65535U) dacValue = 65535U;
    return (uint16_t)dacValue;
}

float dac_value_to_voltage(uint16_t dacValue) {
    return ((float)dacValue / 65535.0f) * DAC_VREF * 2.0f;
}

int dac_init(I2cBus *bus, FemUnit unit) {

    uint8_t resetData[2];

    (void)unit; /* address is same on both FEM buses */

    if (!bus || !bus->initialized) {
        return STATUS_NOK;
    }

    if (i2c_bus_detect(bus, I2C_ADDR_DAC_AD5667) != STATUS_OK) {
        usys_log_error("DAC AD5667 not detected on bus %d", bus->busNum);
        return STATUS_NOK;
    }

    /* Reset sequence (same as your old code) */
    resetData[0] = 0x06;
    resetData[1] = 0x00;
    if (i2c_bus_write_reg(bus, I2C_ADDR_DAC_AD5667, 0x40, resetData, 2) != STATUS_OK) {
        usys_log_error("DAC reset failed (bus %d)", bus->busNum);
        return STATUS_NOK;
    }

    usleep(10000); /* 10ms */
    usys_log_info("DAC initialized (bus %d) for FEM%d", bus->busNum, unit);

    return STATUS_OK;
}

static int dac_set_voltage(I2cBus *bus, FemUnit unit, float voltage, bool isCarrier) {

    float maxAllowed = isCarrier ? DAC_MAX_CARRIER_VOLTAGE : DAC_MAX_PEAK_VOLTAGE;
    uint16_t dacValue = 0;
    uint8_t data[2];
    uint8_t reg = isCarrier ? 0x59 : 0x58;

    if (!bus || !bus->initialized) {
        return STATUS_NOK;
    }

    if (voltage < 0.0f || voltage > maxAllowed) {
        usys_log_error("%s voltage out of range: %.2fV (max %.2fV)",
                       isCarrier ? "Carrier" : "Peak", voltage, maxAllowed);
        return STATUS_NOK;
    }

    dacValue = voltage_to_dac_value(voltage);
    data[0] = (uint8_t)((dacValue >> 8) & 0xFF);
    data[1] = (uint8_t)(dacValue & 0xFF);

    if (i2c_bus_write_reg(bus, I2C_ADDR_DAC_AD5667, reg, data, 2) != STATUS_OK) {
        usys_log_error("Failed to set %s voltage (bus %d FEM%d)",
                       isCarrier ? "carrier" : "peak", bus->busNum, unit);
        return STATUS_NOK;
    }

    usys_log_info("%s voltage set to %.2fV (bus %d FEM%d)",
                  isCarrier ? "Carrier" : "Peak", voltage, bus->busNum, unit);

    return STATUS_OK;
}

int dac_set_carrier_voltage(I2cBus *bus, FemUnit unit, float voltage) {
    return dac_set_voltage(bus, unit, voltage, true);
}

int dac_set_peak_voltage(I2cBus *bus, FemUnit unit, float voltage) {
    return dac_set_voltage(bus, unit, voltage, false);
}

int dac_disable_pa(I2cBus *bus, FemUnit unit) {

    int r1 = dac_set_carrier_voltage(bus, unit, 0.0f);
    int r2 = dac_set_peak_voltage(bus, unit, 0.0f);

    if (r1 == STATUS_OK && r2 == STATUS_OK) {
        usys_log_info("PA disabled - DAC values set to zero (bus %d FEM%d)",
                      bus->busNum, unit);
        return STATUS_OK;
    }

    return STATUS_NOK;
}

/* Temperature sensor (LM75A) */

float lm75a_raw_to_celsius(uint8_t msb, uint8_t lsb) {

    uint16_t tempRaw = (uint16_t)((msb << 8) | lsb);
    int16_t temp9bit = (int16_t)(tempRaw >> 7);

    if (temp9bit > 255) temp9bit = (int16_t)(temp9bit - 512);
    return (float)temp9bit * 0.5f;
}

uint16_t celsius_to_lm75a_raw(float temperature) {

    int16_t temp9bit = (int16_t)(temperature / 0.5f);
    if (temp9bit < 0) temp9bit = (int16_t)(temp9bit + 512);
    return (uint16_t)((temp9bit << 7) & 0xFF80);
}

int temp_sensor_init(I2cBus *bus, FemUnit unit) {

    (void)unit;

    if (!bus || !bus->initialized) {
        return STATUS_NOK;
    }

    if (i2c_bus_detect(bus, I2C_ADDR_TEMP_LM75A) != STATUS_OK) {
        usys_log_error("Temperature sensor LM75A not detected on bus %d", bus->busNum);
        return STATUS_NOK;
    }

    usys_log_info("Temperature sensor initialized (bus %d) for FEM%d", bus->busNum, unit);
    return STATUS_OK;
}

int temp_sensor_read(I2cBus *bus, FemUnit unit, float *temperature) {

    uint8_t data[2] = {0};

    if (!bus || !bus->initialized || !temperature) {
        return STATUS_NOK;
    }

    if (i2c_bus_read_reg(bus, I2C_ADDR_TEMP_LM75A, 0x00, data, 2) != STATUS_OK) {
        usys_log_error("Failed to read temperature (bus %d FEM%d)", bus->busNum, unit);
        return STATUS_NOK;
    }

    *temperature = lm75a_raw_to_celsius(data[0], data[1]);
    usys_log_debug("Temperature read: %.1fC (bus %d FEM%d)",
                   *temperature, bus->busNum, unit);

    return STATUS_OK;
}

int temp_sensor_set_threshold(I2cBus *bus, FemUnit unit, float threshold) {

    uint16_t thresholdRaw = 0;
    uint8_t data[2] = {0};

    if (!bus || !bus->initialized) {
        return STATUS_NOK;
    }

    if (threshold < TEMP_MIN || threshold > TEMP_MAX) {
        usys_log_error("Temperature threshold out of range: %.1fC", threshold);
        return STATUS_NOK;
    }

    thresholdRaw = celsius_to_lm75a_raw(threshold);
    data[0] = (uint8_t)((thresholdRaw >> 8) & 0xFF);
    data[1] = (uint8_t)(thresholdRaw & 0xFF);

    if (i2c_bus_write_reg(bus, I2C_ADDR_TEMP_LM75A, 0x03, data, 2) != STATUS_OK) {
        usys_log_error("Failed to set temperature threshold (bus %d FEM%d)",
                       bus->busNum, unit);
        return STATUS_NOK;
    }

    usys_log_info("Temperature threshold set to %.1fC (bus %d FEM%d)",
                  threshold, bus->busNum, unit);
    return STATUS_OK;
}

/* ADC (ADS1015) */

static int adc_configure_channel(I2cBus *bus, ADCChannel channel) {

    uint16_t muxConfigs[4];
    uint16_t config = 0;
    uint8_t data[2];

    if (!bus || !bus->initialized) {
        return STATUS_NOK;
    }

    muxConfigs[0] = 0x4000; /* AIN0 vs GND */
    muxConfigs[1] = 0x5000; /* AIN1 vs GND */
    muxConfigs[2] = 0x6000; /* AIN2 vs GND */
    muxConfigs[3] = 0x7000; /* AIN3 vs GND */

    if (channel < 0 || channel >= ADC_CHANNEL_MAX) {
        return STATUS_NOK;
    }

    config  = 0x8000;                    /* Start conversion */
    config |= muxConfigs[(int)channel];  /* MUX */
    config |= 0x0200;                    /* PGA ±4.096V */
    config |= 0x0100;                    /* Single-shot mode */
    config |= 0x0080;                    /* 1600 SPS */
    config |= 0x0003;                    /* Disable comparator */

    data[0] = (uint8_t)((config >> 8) & 0xFF);
    data[1] = (uint8_t)(config & 0xFF);

    return i2c_bus_write_reg(bus, I2C_ADDR_ADC_ADS1015, 0x01, data, 2);
}

float adc_raw_to_voltage(uint16_t rawValue) {
    return ((int16_t)rawValue >> 4) * 4.096f / 2048.0f;
}

int adc_init(I2cBus *bus, FemUnit unit) {

    (void)unit;

    if (!bus || !bus->initialized) {
        return STATUS_NOK;
    }

    if (i2c_bus_detect(bus, I2C_ADDR_ADC_ADS1015) != STATUS_OK) {
        usys_log_error("ADC ADS1015 not detected on bus %d", bus->busNum);
        return STATUS_NOK;
    }

    usys_log_info("ADC initialized (bus %d) for FEM%d", bus->busNum, unit);
    return STATUS_OK;
}

int adc_read_channel(I2cBus *bus, FemUnit unit, ADCChannel channel, float *voltage) {

    uint8_t data[2] = {0};
    uint16_t adcRaw = 0;

    if (!bus || !bus->initialized || !voltage) {
        return STATUS_NOK;
    }

    if (channel < 0 || channel >= ADC_CHANNEL_MAX) {
        return STATUS_NOK;
    }

    if (adc_configure_channel(bus, channel) != STATUS_OK) {
        return STATUS_NOK;
    }

    usleep(10000); /* 10ms conversion settle; keep for now */

    if (i2c_bus_read_reg(bus, I2C_ADDR_ADC_ADS1015, 0x00, data, 2) != STATUS_OK) {
        usys_log_error("Failed to read ADC channel %d (bus %d FEM%d)",
                       (int)channel, bus->busNum, unit);
        return STATUS_NOK;
    }

    adcRaw = (uint16_t)((data[0] << 8) | data[1]);
    *voltage = adc_raw_to_voltage(adcRaw);

    usys_log_debug("ADC channel %d: %.3fV (bus %d FEM%d)",
                   (int)channel, *voltage, bus->busNum, unit);

    return STATUS_OK;
}

float voltage_to_reverse_power(float voltage) {
    return (voltage - 2.0f) * 20.0f - 30.0f;
}

float voltage_to_current(float voltage) {
    /* TODO: use shunt + calibration; keep placeholder compatible with current logic */
    return voltage; /* Example: 1V = 1A */
}

int adc_read_reverse_power(I2cBus *bus, FemUnit unit, float *powerDbm) {

    float voltage = 0.0f;

    if (!powerDbm) return STATUS_NOK;

    if (adc_read_channel(bus, unit, ADC_CHANNEL_REVERSE_POWER, &voltage) != STATUS_OK) {
        return STATUS_NOK;
    }

    *powerDbm = voltage_to_reverse_power(voltage);
    return STATUS_OK;
}

int adc_read_pa_current(I2cBus *bus, FemUnit unit, float *currentA) {

    float voltage = 0.0f;

    if (!currentA) return STATUS_NOK;

    if (adc_read_channel(bus, unit, ADC_CHANNEL_PA_CURRENT, &voltage) != STATUS_OK) {
        return STATUS_NOK;
    }

    *currentA = voltage_to_current(voltage);
    return STATUS_OK;
}

int adc_read_all_channels(I2cBus *bus, FemUnit unit) {

    typedef struct {
        ADCChannel channel;
        const char *name;
        const char *unit;
        float (*convert)(float);
    } ChannelInfo;

    ChannelInfo channelsInfo[4];
    float voltage = 0.0f;
    float converted = 0.0f;

    if (!bus || !bus->initialized) {
        return STATUS_NOK;
    }

    channelsInfo[0] = (ChannelInfo){ADC_CHANNEL_REVERSE_POWER,
                      "Reverse Power", "dBm", voltage_to_reverse_power};
    channelsInfo[1] = (ChannelInfo){ADC_CHANNEL_FORWARD_POWER,
                      "Forward Power", "dBm", voltage_to_reverse_power};
    channelsInfo[2] = (ChannelInfo){ADC_CHANNEL_PA_CURRENT,
                      "PA Current",    "A",   voltage_to_current};
    channelsInfo[3] = (ChannelInfo){ADC_CHANNEL_TEMPERATURE,
                      "Temperature",   "V",   NULL};

    usys_log_info("Reading all ADC channels (bus %d FEM%d):", bus->busNum, unit);

    for (int i = 0; i < 4; i++) {
        if (adc_read_channel(bus, unit, channelsInfo[i].channel, &voltage) == STATUS_OK) {
            if (channelsInfo[i].convert) {
                converted = channelsInfo[i].convert(voltage);
                usys_log_info("  Channel %d (%s): %.2f %s",
                              (int)channelsInfo[i].channel,
                              channelsInfo[i].name,
                              converted,
                              channelsInfo[i].unit);
            } else {
                usys_log_info("  Channel %d (%s): %.3f %s",
                              (int)channelsInfo[i].channel,
                              channelsInfo[i].name,
                              voltage,
                              channelsInfo[i].unit);
            }
        } else {
            usys_log_error("  Channel %d (%s): Read error",
                           (int)channelsInfo[i].channel, channelsInfo[i].name);
        }
    }

    return STATUS_OK;
}

/* EEPROM (24xx @ 0x50) */

int eeprom_write_serial(I2cBus *bus, FemUnit unit, const char *serial) {

    size_t len = 0U;
    uint8_t byteData = 0U;
    uint8_t nullTerm = 0x00;

    if (!bus || !bus->initialized || !serial) {
        return STATUS_NOK;
    }

    len = strlen(serial);
    if (len > 16U) {
        usys_log_error("Serial number too long (max 16 characters)");
        return STATUS_NOK;
    }

    for (size_t i = 0U; i < len; i++) {
        byteData = (uint8_t)serial[i];
        if (i2c_bus_write_reg(bus, I2C_ADDR_EEPROM, (uint8_t)i, &byteData, 1) != STATUS_OK) {
            usys_log_error("Failed to write EEPROM at position %zu (bus %d FEM%d)",
                           i, bus->busNum, unit);
            return STATUS_NOK;
        }
        usleep(10000); /* EEPROM write cycle */
    }

    (void)i2c_bus_write_reg(bus, I2C_ADDR_EEPROM, (uint8_t)len, &nullTerm, 1);

    usys_log_info("Serial number written to EEPROM (bus %d FEM%d): %s",
                  bus->busNum, unit, serial);

    return STATUS_OK;
}

int eeprom_read_serial(I2cBus *bus, FemUnit unit, char *serial, size_t maxLen) {

    size_t readLen = 0U;
    uint8_t byteData = 0U;

    if (!bus || !bus->initialized || !serial || maxLen == 0U) {
        return STATUS_NOK;
    }

    readLen = (maxLen - 1U < 16U) ? (maxLen - 1U) : 16U;
    serial[0] = '\0';

    for (size_t i = 0U; i < readLen; i++) {
        if (i2c_bus_read_reg(bus, I2C_ADDR_EEPROM, (uint8_t)i, &byteData, 1) != STATUS_OK) {
            usys_log_error("Failed to read EEPROM at position %zu (bus %d FEM%d)",
                           i, bus->busNum, unit);
            return STATUS_NOK;
        }

        if (byteData == 0U) {
            serial[i] = '\0';
            break;
        }

        serial[i] = (char)byteData;
        if (i == readLen - 1U) {
            serial[i + 1U] = '\0';
        }
    }

    if (strlen(serial) > 0U) {
        usys_log_info("Serial number read from EEPROM (bus %d FEM%d): %s",
                      bus->busNum, unit, serial);
        return STATUS_OK;
    }

    usys_log_info("No serial number found in EEPROM (bus %d FEM%d)", bus->busNum, unit);
    return STATUS_NOK;
}

/* Controller temp (TMP10x minimal read) */
/*
 * TMP102/TMP101 style:
 * reg 0x00: temp MSB/LSB, 12-bit typically.
 * We'll do a conservative decode: assume 12-bit, right-shift 4.
 */
int ctrl_temp_read_tmp10x(I2cBus *bus, float *temperatureC) {

    uint8_t data[2] = {0};
    int16_t raw = 0;

    if (!bus || !bus->initialized || !temperatureC) {
        return STATUS_NOK;
    }

    if (i2c_bus_read_reg(bus, I2C_ADDR_CTRL_TMP10X, 0x00, data, 2) != STATUS_OK) {
        usys_log_error("Failed to read controller temp (bus %d)", bus->busNum);
        return STATUS_NOK;
    }

    raw = (int16_t)((data[0] << 8) | data[1]);
    raw >>= 4; /* 12-bit */

    /* sign extend if negative */
    if (raw & 0x0800) {
        raw |= 0xF000;
    }

    *temperatureC = (float)raw * 0.0625f;
    return STATUS_OK;
}
