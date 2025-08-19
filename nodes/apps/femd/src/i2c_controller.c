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
#include <math.h>
#include <stddef.h>
#include <fcntl.h>
#include <sys/ioctl.h>
#include <linux/i2c-dev.h>
#include <errno.h>

#include "i2c_controller.h"
#include "femd.h"

/* Unified device metadata (variables in camelCase) */
static const I2CDeviceInfo deviceInfo[I2C_DEVICE_MAX] = {
    {"AD5667",  I2C_ADDR_DAC_AD5667,  "16-bit DAC"},
    {"LM75A",   I2C_ADDR_TEMP_LM75A,  "Temperature sensor"},
    {"ADS1015", I2C_ADDR_ADC_ADS1015, "12-bit ADC"},
    {"EEPROM",  I2C_ADDR_EEPROM,      "Serial number storage"}
};

/* Forward statics */
static int adc_configure_channel(int bus, ADCChannel channel);
static int dac_set_voltage(I2CController *controller, FemUnit unit, float voltage, bool isCarrier);

/* Controller lifecycle */

int i2c_controller_init(I2CController *controller) {
    int status = STATUS_OK;

    if (!controller) {
        usys_log_error("I2C controller pointer is NULL");
        return STATUS_NOK;
    }

    memset(controller, 0, sizeof(I2CController));

    controller->busFem1 = I2C_BUS_FEM1;
    controller->busFem2 = I2C_BUS_FEM2;

    controller->dacState.carrierVoltage = 0.0f;
    controller->dacState.peakVoltage    = 0.0f;
    controller->dacState.initialized    = false;

    controller->tempState.temperature  = 0.0f;
    controller->tempState.threshold    = 85.0f;
    controller->tempState.alertEnabled = false;

    controller->adcState.maxReversePower = -10.0f;
    controller->adcState.maxCurrent      = 5.0f;
    controller->adcState.safetyEnabled   = true;

    controller->eepromState.hasData = false;
    memset(controller->eepromState.serialNumber, 0, sizeof(controller->eepromState.serialNumber));

    controller->initialized = true;
    usys_log_info("I2C controller initialized");

    return status;
}

void i2c_controller_cleanup(I2CController *controller) {
    if (controller) {
        controller->initialized = false;
        usys_log_info("I2C controller cleaned up");
    }
}

int i2c_get_bus_for_fem(FemUnit unit) {
    int bus = I2C_BUS_FEM2;
    if (unit == FEM_UNIT_1) {
        bus = I2C_BUS_FEM1;
    }
    return bus;
}

/* Low-level helpers */

int i2c_write_bytes(int bus, uint8_t devAddr, uint8_t reg, const uint8_t *data, size_t len) {
    char devicePath[32] = {0};
    int fd = -1;
    uint8_t txBuffer[256] = {0};
    ssize_t wrote = 0;
    int status = STATUS_OK;

    if (len > 254) { /* 1 byte reserved for reg */
        usys_log_error("I2C write data too long: %zu bytes", len);
        return STATUS_NOK;
    }

    snprintf(devicePath, sizeof(devicePath), "/dev/i2c-%d", bus);

    fd = open(devicePath, O_RDWR);
    if (fd < 0) {
        usys_log_error("Failed to open I2C device %s: %s", devicePath, strerror(errno));
        status = STATUS_NOK;
        goto cleanup;
    }

    if (ioctl(fd, I2C_SLAVE, devAddr) < 0) {
        usys_log_error("Failed to set I2C slave address 0x%02X: %s", devAddr, strerror(errno));
        status = STATUS_NOK;
        goto cleanup;
    }

    txBuffer[0] = reg;
    if (len > 0 && data != NULL) {
        memcpy(&txBuffer[1], data, len);
    }

    wrote = write(fd, txBuffer, len + 1);
    if (wrote != (ssize_t)(len + 1)) {
        usys_log_error("I2C write failed: bus=%d, addr=0x%02X, reg=0x%02X: %s",
                       bus, devAddr, reg, strerror(errno));
        status = STATUS_NOK;
        goto cleanup;
    }

    usys_log_debug("I2C write OK: bus=%d, addr=0x%02X, reg=0x%02X, len=%zu",
                   bus, devAddr, reg, len);

cleanup:
    if (fd >= 0) {
        close(fd);
    }
    return status;
}

int i2c_read_bytes(int bus, uint8_t devAddr, uint8_t reg, uint8_t *data, size_t len) {
    char devicePath[32];
    int fd = -1;
    ssize_t wrote = 0;
    ssize_t readn = 0;
    int status = STATUS_OK;

    if (!data || len == 0) {
        usys_log_error("Invalid parameters for I2C read");
        return STATUS_NOK;
    }

    snprintf(devicePath, sizeof(devicePath), "/dev/i2c-%d", bus);

    fd = open(devicePath, O_RDWR);
    if (fd < 0) {
        usys_log_error("Failed to open I2C device %s: %s", devicePath, strerror(errno));
        status = STATUS_NOK;
        goto cleanup;
    }

    if (ioctl(fd, I2C_SLAVE, devAddr) < 0) {
        usys_log_error("Failed to set I2C slave address 0x%02X: %s", devAddr, strerror(errno));
        status = STATUS_NOK;
        goto cleanup;
    }

    wrote = write(fd, &reg, 1);
    if (wrote != 1) {
        usys_log_error("Failed to write register address: %s", strerror(errno));
        status = STATUS_NOK;
        goto cleanup;
    }

    readn = read(fd, data, len);
    if (readn != (ssize_t)len) {
        usys_log_error("I2C read failed: bus=%d, addr=0x%02X, reg=0x%02X: %s",
                       bus, devAddr, reg, strerror(errno));
        status = STATUS_NOK;
        goto cleanup;
    }

    usys_log_debug("I2C read OK: bus=%d, addr=0x%02X, reg=0x%02X, len=%zu",
                   bus, devAddr, reg, len);

cleanup:
    if (fd >= 0) {
        close(fd);
    }
    return status;
}

int i2c_detect_device(int bus, uint8_t devAddr) {
    char devicePath[32];
    int fd = -1;
    int status = STATUS_OK;
    ssize_t wrote = 0;

    snprintf(devicePath, sizeof(devicePath), "/dev/i2c-%d", bus);

    fd = open(devicePath, O_RDWR);
    if (fd < 0) {
        usys_log_debug("Failed to open I2C device %s: %s", devicePath, strerror(errno));
        status = STATUS_NOK;
        goto cleanup;
    }

    if (ioctl(fd, I2C_SLAVE, devAddr) < 0) {
        status = STATUS_NOK;
        goto cleanup;
    }

    /* No-data write to detect presence */
    wrote = write(fd, NULL, 0);
    if (wrote < 0) {
        status = STATUS_NOK;
    }

cleanup:
    if (fd >= 0) {
        close(fd);
    }
    return status;
}

/* Device info/scan */

const I2CDeviceInfo* i2c_get_device_info(I2CDevice device) {
    const I2CDeviceInfo *info = NULL;
    if (device < I2C_DEVICE_MAX) {
        info = &deviceInfo[device];
    }
    return info;
}

void i2c_print_device_scan(FemUnit unit) {
    int bus = 0;
    int i = 0;
    int detected = STATUS_NOK;
    const I2CDeviceInfo *info = NULL;

    bus = i2c_get_bus_for_fem(unit);
    usys_log_info("I2C Device Scan for FEM%d (bus %d):", unit, bus);

    for (i = 0; i < I2C_DEVICE_MAX; i++) {
        info = &deviceInfo[i];
        detected = i2c_detect_device(bus, info->address);
        usys_log_info("  %s (0x%02X): %s - %s",
                      info->name, info->address, info->description,
                      detected == STATUS_OK ? "DETECTED" : "NOT FOUND");
    }
}

/* DAC */

uint16_t voltage_to_dac_value(float voltage) {
    float compensatedVoltage = 0.0f;
    uint16_t dacValue = 0;

    compensatedVoltage = voltage / 2.0f;
    dacValue = (uint16_t)((compensatedVoltage / DAC_VREF) * 65535.0f);
    if (dacValue > 65535) {
        dacValue = 65535;
    }
    return dacValue;
}

float dac_value_to_voltage(uint16_t dacValue) {
    return ((float)dacValue / 65535.0f) * DAC_VREF * 2.0f;
}

int dac_init(I2CController *controller, FemUnit unit) {
    int status = STATUS_OK;
    int bus = 0;
    uint8_t resetData[2];

    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }

    bus = i2c_get_bus_for_fem(unit);

    if (i2c_detect_device(bus, I2C_ADDR_DAC_AD5667) != STATUS_OK) {
        usys_log_error("DAC AD5667 not detected on bus %d", bus);
        return STATUS_NOK;
    }

    resetData[0] = 0x06;
    resetData[1] = 0x00;
    status = i2c_write_bytes(bus, I2C_ADDR_DAC_AD5667, 0x40, resetData, 2);
    if (status != STATUS_OK) {
        usys_log_error("DAC reset failed");
        return STATUS_NOK;
    }

    usleep(10000); /* 10ms */

    controller->dacState.initialized = true;
    usys_log_info("DAC initialized for FEM%d", unit);
    return STATUS_OK;
}

static int dac_set_voltage(I2CController *controller, FemUnit unit, float voltage, bool isCarrier) {
    int status = STATUS_OK;
    float maxAllowed = 0.0f;
    int bus = 0;
    uint16_t dacValue = 0;
    uint8_t data[2];
    uint8_t reg = 0;

    if (!controller || !controller->initialized || !controller->dacState.initialized) {
        return STATUS_NOK;
    }

    maxAllowed = isCarrier ? DAC_MAX_CARRIER_VOLTAGE : DAC_MAX_PEAK_VOLTAGE;
    if (voltage < 0.0f || voltage > maxAllowed) {
        usys_log_error("%s voltage out of range: %.2fV (max %.2fV)",
                       isCarrier ? "Carrier" : "Peak", voltage, maxAllowed);
        return STATUS_NOK;
    }

    bus      = i2c_get_bus_for_fem(unit);
    dacValue = voltage_to_dac_value(voltage);

    data[0] = (uint8_t)((dacValue >> 8) & 0xFF);
    data[1] = (uint8_t)(dacValue & 0xFF);
    reg     = isCarrier ? 0x59 : 0x58;

    status = i2c_write_bytes(bus, I2C_ADDR_DAC_AD5667, reg, data, 2);
    if (status != STATUS_OK) {
        usys_log_error("Failed to set %s voltage", isCarrier ? "carrier" : "peak");
        return STATUS_NOK;
    }

    if (isCarrier) {
        controller->dacState.carrierVoltage = voltage;
    } else {
        controller->dacState.peakVoltage = voltage;
    }

    usys_log_info("%s voltage set to %.2fV for FEM%d",
                  isCarrier ? "Carrier" : "Peak", voltage, unit);
    return STATUS_OK;
}

int dac_set_carrier_voltage(I2CController *controller, FemUnit unit, float voltage) {
    return dac_set_voltage(controller, unit, voltage, true);
}

int dac_set_peak_voltage(I2CController *controller, FemUnit unit, float voltage) {
    return dac_set_voltage(controller, unit, voltage, false);
}

int dac_get_config(I2CController *controller, float *carrier, float *peak) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }
    if (carrier) {
        *carrier = controller->dacState.carrierVoltage;
    }
    if (peak) {
        *peak = controller->dacState.peakVoltage;
    }
    return STATUS_OK;
}

int dac_disable_pa(I2CController *controller, FemUnit unit) {
    int result1 = STATUS_NOK;
    int result2 = STATUS_NOK;

    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }

    result1 = dac_set_carrier_voltage(controller, unit, 0.0f);
    result2 = dac_set_peak_voltage(controller, unit, 0.0f);

    if (result1 == STATUS_OK && result2 == STATUS_OK) {
        usys_log_info("PA disabled - DAC values set to zero for FEM%d", unit);
        return STATUS_OK;
    }
    return STATUS_NOK;
}

/* Temperature sensor (LM75A) */

float lm75a_raw_to_celsius(uint8_t msb, uint8_t lsb) {
    uint16_t tempRaw = 0;
    int16_t  temp9bit = 0;
    float    tempC = 0.0f;

    tempRaw  = (uint16_t)((msb << 8) | lsb);
    temp9bit = (int16_t)(tempRaw >> 7);
    if (temp9bit > 255) {
        temp9bit = (int16_t)(temp9bit - 512);
    }

    tempC = temp9bit * 0.5f;
    return tempC;
}

uint16_t celsius_to_lm75a_raw(float temperature) {
    int16_t  temp9bit = 0;
    uint16_t raw = 0;

    temp9bit = (int16_t)(temperature / 0.5f);
    if (temp9bit < 0) {
        temp9bit = (int16_t)(temp9bit + 512);
    }

    raw = (uint16_t)((temp9bit << 7) & 0xFF80);
    return raw;
}

int temp_sensor_init(I2CController *controller, FemUnit unit) {
    int bus    = 0;
    int status = STATUS_OK;

    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }

    bus    = i2c_get_bus_for_fem(unit);
    status = i2c_detect_device(bus, I2C_ADDR_TEMP_LM75A);
    if (status != STATUS_OK) {
        usys_log_error("Temperature sensor LM75A not detected on bus %d", bus);
        return STATUS_NOK;
    }

    usys_log_info("Temperature sensor initialized for FEM%d", unit);
    return STATUS_OK;
}

int temp_sensor_read(I2CController *controller, FemUnit unit, float *temperature) {
    int bus = 0;
    uint8_t data[2];
    int status = STATUS_OK;

    if (!controller || !controller->initialized || !temperature) {
        return STATUS_NOK;
    }

    bus    = i2c_get_bus_for_fem(unit);
    status = i2c_read_bytes(bus, I2C_ADDR_TEMP_LM75A, 0x00, data, 2);
    if (status != STATUS_OK) {
        usys_log_error("Failed to read temperature from FEM%d", unit);
        return STATUS_NOK;
    }

    *temperature = lm75a_raw_to_celsius(data[0], data[1]);
    controller->tempState.temperature = *temperature;

    usys_log_debug("Temperature read: %.1f°C from FEM%d", *temperature, unit);
    return STATUS_OK;
}

int temp_sensor_set_threshold(I2CController *controller, FemUnit unit, float threshold) {
    int bus = 0;
    uint16_t thresholdRaw = 0;
    uint8_t data[2];
    int status = STATUS_OK;

    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }

    if (threshold < TEMP_MIN || threshold > TEMP_MAX) {
        usys_log_error("Temperature threshold out of range: %.1f°C", threshold);
        return STATUS_NOK;
    }

    bus          = i2c_get_bus_for_fem(unit);
    thresholdRaw = celsius_to_lm75a_raw(threshold);
    data[0]      = (uint8_t)((thresholdRaw >> 8) & 0xFF);
    data[1]      = (uint8_t)(thresholdRaw & 0xFF);

    status = i2c_write_bytes(bus, I2C_ADDR_TEMP_LM75A, 0x03, data, 2);
    if (status != STATUS_OK) {
        usys_log_error("Failed to set temperature threshold");
        return STATUS_NOK;
    }

    controller->tempState.threshold    = threshold;
    controller->tempState.alertEnabled = true;
    usys_log_info("Temperature threshold set to %.1f°C for FEM%d", threshold, unit);
    return STATUS_OK;
}

/* ADC (ADS1015) */

static int adc_configure_channel(int bus, ADCChannel channel) {
    uint16_t muxConfigs[4];
    uint16_t config = 0;
    uint8_t data[2];
    int status = STATUS_OK;

    muxConfigs[0] = 0x4000; /* AIN0 vs GND */
    muxConfigs[1] = 0x5000; /* AIN1 vs GND */
    muxConfigs[2] = 0x6000; /* AIN2 vs GND */
    muxConfigs[3] = 0x7000; /* AIN3 vs GND */

    if (channel < 0 || channel >= ADC_CHANNEL_MAX) {
        return STATUS_NOK;
    }

    config = 0x8000;                   /* Start conversion */
    config |= muxConfigs[(int)channel];/* MUX */
    config |= 0x0200;                  /* PGA ±4.096V */
    config |= 0x0100;                  /* Single-shot mode */
    config |= 0x0080;                  /* 1600 SPS */
    config |= 0x0003;                  /* Disable comparator */

    data[0] = (uint8_t)((config >> 8) & 0xFF);
    data[1] = (uint8_t)(config & 0xFF);

    status = i2c_write_bytes(bus, I2C_ADDR_ADC_ADS1015, 0x01, data, 2);
    return status;
}

float adc_raw_to_voltage(uint16_t rawValue) {
    return ((int16_t)rawValue >> 4) * 4.096f / 2048.0f;
}

int adc_init(I2CController *controller, FemUnit unit) {
    int bus = 0;
    int status = STATUS_OK;

    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }

    bus    = i2c_get_bus_for_fem(unit);
    status = i2c_detect_device(bus, I2C_ADDR_ADC_ADS1015);
    if (status != STATUS_OK) {
        usys_log_error("ADC ADS1015 not detected on bus %d", bus);
        return STATUS_NOK;
    }

    usys_log_info("ADC initialized for FEM%d", unit);
    return STATUS_OK;
}

int adc_read_channel(I2CController *controller, FemUnit unit, ADCChannel channel, float *voltage) {
    int bus = 0;
    int status = STATUS_OK;
    uint8_t data[2];
    uint16_t adcRaw = 0;

    if (!controller || !controller->initialized || !voltage) {
        return STATUS_NOK;
    }

    if (channel < 0 || channel >= ADC_CHANNEL_MAX) {
        return STATUS_NOK;
    }

    bus = i2c_get_bus_for_fem(unit);

    status = adc_configure_channel(bus, channel);
    if (status != STATUS_OK) {
        return STATUS_NOK;
    }

    usleep(10000); /* 10ms */

    status = i2c_read_bytes(bus, I2C_ADDR_ADC_ADS1015, 0x00, data, 2);
    if (status != STATUS_OK) {
        usys_log_error("Failed to read ADC channel %d from FEM%d", (int)channel, unit);
        return STATUS_NOK;
    }

    adcRaw = (uint16_t)((data[0] << 8) | data[1]);
    *voltage = adc_raw_to_voltage(adcRaw);

    usys_log_debug("ADC channel %d: %.3fV from FEM%d", (int)channel, *voltage, unit);
    return STATUS_OK;
}

float voltage_to_reverse_power(float voltage) {
    return (voltage - 2.0f) * 20.0f - 30.0f;
}

float voltage_to_current(float voltage) {
    float current = 0.0f;
    current = voltage; /* Example: 1V = 1A */
    return current;
}

int adc_read_reverse_power(I2CController *controller, FemUnit unit, float *powerDbm) {
    float voltage = 0.0f;
    int status    = STATUS_OK;

    if (!powerDbm) {
        return STATUS_NOK;
    }

    status = adc_read_channel(controller, unit, ADC_CHANNEL_REVERSE_POWER, &voltage);
    if (status != STATUS_OK) {
        return STATUS_NOK;
    }

    *powerDbm = voltage_to_reverse_power(voltage);
    controller->adcState.reversePowerDbm = *powerDbm;

    return STATUS_OK;
}

int adc_read_pa_current(I2CController *controller, FemUnit unit, float *currentA) {
    float voltage = 0.0f;
    int status    = STATUS_OK;

    if (!currentA) {
        return STATUS_NOK;
    }

    status = adc_read_channel(controller, unit, ADC_CHANNEL_PA_CURRENT, &voltage);
    if (status != STATUS_OK) {
        return STATUS_NOK;
    }

    *currentA = voltage_to_current(voltage);
    controller->adcState.paCurrentA = *currentA;

    return STATUS_OK;
}

int adc_read_all_channels(I2CController *controller, FemUnit unit) {
    int i = 0;
    float voltage = 0.0f;
    float converted = 0.0f;
    int status = STATUS_OK;

    struct ChannelInfo {
        ADCChannel channel;
        const char *name;
        const char *unit;
        float      (*convert)(float);
    };

    struct ChannelInfo channelsInfo[4];

    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }

    channelsInfo[0].channel = ADC_CHANNEL_REVERSE_POWER;
    channelsInfo[0].name    = "Reverse Power";
    channelsInfo[0].unit    = "dBm";
    channelsInfo[0].convert = voltage_to_reverse_power;

    channelsInfo[1].channel = ADC_CHANNEL_FORWARD_POWER;
    channelsInfo[1].name    = "Forward Power";
    channelsInfo[1].unit    = "dBm";
    channelsInfo[1].convert = voltage_to_reverse_power;

    channelsInfo[2].channel = ADC_CHANNEL_PA_CURRENT;
    channelsInfo[2].name    = "PA Current";
    channelsInfo[2].unit    = "A";
    channelsInfo[2].convert = voltage_to_current;

    channelsInfo[3].channel = ADC_CHANNEL_TEMPERATURE;
    channelsInfo[3].name    = "Temperature";
    channelsInfo[3].unit    = "V";
    channelsInfo[3].convert = NULL; /* direct voltage */

    usys_log_info("Reading all ADC channels for FEM%d:", unit);

    for (i = 0; i < 4; i++) {
        status = adc_read_channel(controller, unit, channelsInfo[i].channel, &voltage);
        if (status == STATUS_OK) {
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
                           (int)channelsInfo[i].channel,
                           channelsInfo[i].name);
        }
    }

    return STATUS_OK;
}

int adc_set_safety_thresholds(I2CController *controller, float maxReversePower, float maxCurrent) {
    if (!controller || !controller->initialized) {
        return STATUS_NOK;
    }

    controller->adcState.maxReversePower = maxReversePower;
    controller->adcState.maxCurrent      = maxCurrent;
    controller->adcState.safetyEnabled   = true;

    usys_log_info("Safety thresholds set: reverse power %.1f dBm, current %.1f A",
                  maxReversePower, maxCurrent);
    return STATUS_OK;
}

int adc_check_safety(I2CController *controller, FemUnit unit, bool *safetyViolation) {
    float reversePower = 0.0f;
    float paCurrent    = 0.0f;

    if (!controller || !controller->initialized || !safetyViolation) {
        return STATUS_NOK;
    }

    *safetyViolation = false;

    if (!controller->adcState.safetyEnabled) {
        return STATUS_OK;
    }

    if (adc_read_reverse_power(controller, unit, &reversePower) == STATUS_OK) {
        if (reversePower > controller->adcState.maxReversePower) {
            usys_log_warn("Safety violation: reverse power %.1f dBm exceeds threshold %.1f dBm",
                          reversePower,
                          controller->adcState.maxReversePower);
            *safetyViolation = true;
        }
    }

    if (adc_read_pa_current(controller, unit, &paCurrent) == STATUS_OK) {
        if (paCurrent > controller->adcState.maxCurrent) {
            usys_log_warn("Safety violation: PA current %.1f A exceeds threshold %.1f A",
                          paCurrent,
                          controller->adcState.maxCurrent);
            *safetyViolation = true;
        }
    }

    return STATUS_OK;
}

/* EEPROM */

int eeprom_write_serial(I2CController *controller, FemUnit unit, const char *serial) {
    int bus = 0;
    size_t len = 0U;
    size_t i = 0U;
    uint8_t byteData = 0U;
    uint8_t nullTerm = 0x00;
    int status = STATUS_OK;

    if (!controller || !controller->initialized || !serial) {
        return STATUS_NOK;
    }

    bus = i2c_get_bus_for_fem(unit);
    len = strlen(serial);

    if (len > 16U) {
        usys_log_error("Serial number too long (max 16 characters)");
        return STATUS_NOK;
    }

    for (i = 0U; i < len; i++) {
        byteData = (uint8_t)serial[i];
        status   = i2c_write_bytes(bus, I2C_ADDR_EEPROM, (uint8_t)i, &byteData, 1);
        if (status != STATUS_OK) {
            usys_log_error("Failed to write EEPROM at position %zu", i);
            return STATUS_NOK;
        }
        usleep(10000); /* 10ms write cycle */
    }

    (void)i2c_write_bytes(bus, I2C_ADDR_EEPROM, (uint8_t)len, &nullTerm, 1);

    strncpy(controller->eepromState.serialNumber, serial,
            sizeof(controller->eepromState.serialNumber) - 1);
    controller->eepromState.hasData = true;

    usys_log_info("Serial number written to EEPROM: %s", serial);
    return STATUS_OK;
}

int eeprom_read_serial(I2CController *controller, FemUnit unit, char *serial, size_t maxLen) {
    int bus = 0;
    size_t readLen = 0U;
    size_t i = 0U;
    uint8_t byteData = 0U;
    int status = STATUS_OK;

    if (!controller || !controller->initialized || !serial || maxLen == 0U) {
        return STATUS_NOK;
    }

    bus     = i2c_get_bus_for_fem(unit);
    readLen = (maxLen - 1U < 16U) ? (maxLen - 1U) : 16U;

    serial[0] = '\0';

    for (i = 0U; i < readLen; i++) {
        status = i2c_read_bytes(bus, I2C_ADDR_EEPROM, (uint8_t)i, &byteData, 1);
        if (status != STATUS_OK) {
            usys_log_error("Failed to read EEPROM at position %zu", i);
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
        strncpy(controller->eepromState.serialNumber, serial,
                sizeof(controller->eepromState.serialNumber) - 1);
        controller->eepromState.hasData = true;
        usys_log_info("Serial number read from EEPROM: %s", serial);
        return STATUS_OK;
    } else {
        usys_log_info("No serial number found in EEPROM");
        return STATUS_NOK;
    }
}
