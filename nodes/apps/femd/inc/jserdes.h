/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef JSERDES_H
#define JSERDES_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>

#include "json_types.h"
#include "gpio_controller.h"
#include "i2c_controller.h"

#define ERR_FEMD_JSON_CREATION_ERR      -100
#define ERR_FEMD_JSON_NO_VAL_TO_ENCODE  -101
#define ERR_FEMD_JSON_PARSER            -102

typedef struct {
    char *name;
    char *value;
    char *units;
    int type;
} PropertyData;

typedef struct {
    char *service_name;
    char *version;
    int uptime;
} ServiceInfo;

typedef struct {
    float temperature;
    float threshold;
    int fem_unit;
    char *units;
} TemperatureData;

typedef struct {
    int channel;
    float voltage;
    int fem_unit;
    char *units;
} AdcData;

typedef struct {
    float reverse_power;
    float pa_current;
    int fem_unit;
} AdcAllData;

typedef struct {
    float carrier_voltage;
    float peak_voltage;
    int fem_unit;
} DacConfig;

typedef struct {
    char *serial;
    int fem_unit;
} EepromData;

typedef struct {
    float max_reverse_power;
    float max_current;
    float max_temperature;
} SafetyConfig;

int json_serialize_error(JsonObj **json, int code, const char *str);
int json_serialize_success(JsonObj **json, const char *message);
int json_serialize_gpio_status(JsonObj **json, const GpioStatus *status, int fem_unit);
int json_serialize_dac_config(JsonObj **json, const DacConfig *config);
int json_serialize_temperature_data(JsonObj **json, const TemperatureData *data);
int json_serialize_adc_data(JsonObj **json, const AdcData *data);
int json_serialize_adc_all_data(JsonObj **json, const AdcAllData *data);
int json_serialize_eeprom_data(JsonObj **json, const EepromData *data);
int json_serialize_safety_config(JsonObj **json, const SafetyConfig *config);
int json_serialize_service_info(JsonObj **json, const ServiceInfo *info);

int json_deserialize_gpio_control(JsonObj *json, char **gpio_name, bool *enabled);
int json_deserialize_dac_control(JsonObj *json, float *carrier_voltage, float *peak_voltage);
int json_deserialize_temperature_threshold(JsonObj *json, float *threshold);
int json_deserialize_eeprom_write(JsonObj *json, char **serial);
int json_deserialize_safety_config(JsonObj *json, SafetyConfig *config);

void json_free(JsonObj **json);
void json_log(JsonObj *json);

#ifdef __cplusplus
}
#endif

#endif /* JSERDES_H */