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

#include "jserdes.h"
#include "femd.h"

int json_serialize_error(JsonObj **json, int code, const char *str) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_ERROR, json_object());
    JsonObj *jError = json_object_get(*json, JTAG_ERROR);
    if (jError) {
        json_object_set_new(jError, JTAG_ERROR_CODE, json_integer(code));
        json_object_set_new(jError, JTAG_ERROR_CSTRING, json_string(str));
    } else {
        json_decref(*json);
        *json = NULL;
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_success(JsonObj **json, const char *message) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_STATUS, json_string("success"));
    if (message) {
        json_object_set_new(*json, JTAG_MESSAGE, json_string(message));
    }

    return ret;
}

int json_serialize_gpio_status(JsonObj **json, const GpioStatus *status, int fem_unit) {
    int ret = JSON_ENCODING_OK;

    if (!status) {
        return ERR_FEMD_JSON_NO_VAL_TO_ENCODE;
    }

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_GPIO_STATUS, json_object());
    JsonObj *jGpio = json_object_get(*json, JTAG_GPIO_STATUS);
    if (jGpio) {
        json_object_set_new(jGpio, JTAG_TX_RF_ENABLE, json_boolean(status->tx_rf_enable));
        json_object_set_new(jGpio, JTAG_RX_RF_ENABLE, json_boolean(status->rx_rf_enable));
        json_object_set_new(jGpio, JTAG_PA_VDS_ENABLE, json_boolean(status->pa_vds_enable));
        json_object_set_new(jGpio, JTAG_RF_PAL_ENABLE, json_boolean(status->rf_pal_enable));
        json_object_set_new(jGpio, JTAG_28V_VDS_ENABLE, json_boolean(!status->pa_disable));
        json_object_set_new(jGpio, JTAG_PSU_PGOOD, json_boolean(status->pg_reg_5v));
        json_object_set_new(jGpio, JTAG_FEM_UNIT, json_integer(fem_unit));
    } else {
        json_decref(*json);
        *json = NULL;
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_dac_config(JsonObj **json, const DacConfig *config) {
    int ret = JSON_ENCODING_OK;

    if (!config) {
        return ERR_FEMD_JSON_NO_VAL_TO_ENCODE;
    }

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_DAC_CONFIG, json_object());
    JsonObj *jDac = json_object_get(*json, JTAG_DAC_CONFIG);
    if (jDac) {
        json_object_set_new(jDac, JTAG_CARRIER_VOLTAGE, json_real(config->carrier_voltage));
        json_object_set_new(jDac, JTAG_PEAK_VOLTAGE, json_real(config->peak_voltage));
        json_object_set_new(jDac, JTAG_FEM_UNIT, json_integer(config->fem_unit));
    } else {
        json_decref(*json);
        *json = NULL;
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_temperature_data(JsonObj **json, const TemperatureData *data) {
    int ret = JSON_ENCODING_OK;

    if (!data) {
        return ERR_FEMD_JSON_NO_VAL_TO_ENCODE;
    }

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_TEMP_READING, json_object());
    JsonObj *jTemp = json_object_get(*json, JTAG_TEMP_READING);
    if (jTemp) {
        json_object_set_new(jTemp, JTAG_TEMPERATURE, json_real(data->temperature));
        json_object_set_new(jTemp, JTAG_FEM_UNIT, json_integer(data->fem_unit));
        if (data->units) {
            json_object_set_new(jTemp, JTAG_UNITS, json_string(data->units));
        }
        if (data->threshold > 0) {
            json_object_set_new(jTemp, JTAG_THRESHOLD, json_real(data->threshold));
        }
    } else {
        json_decref(*json);
        *json = NULL;
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_adc_data(JsonObj **json, const AdcData *data) {
    int ret = JSON_ENCODING_OK;

    if (!data) {
        return ERR_FEMD_JSON_NO_VAL_TO_ENCODE;
    }

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_ADC_READING, json_object());
    JsonObj *jAdc = json_object_get(*json, JTAG_ADC_READING);
    if (jAdc) {
        json_object_set_new(jAdc, JTAG_CHANNEL, json_integer(data->channel));
        json_object_set_new(jAdc, JTAG_VOLTAGE, json_real(data->voltage));
        json_object_set_new(jAdc, JTAG_FEM_UNIT, json_integer(data->fem_unit));
        if (data->units) {
            json_object_set_new(jAdc, JTAG_UNITS, json_string(data->units));
        }
    } else {
        json_decref(*json);
        *json = NULL;
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_adc_all_data(JsonObj **json, const AdcAllData *data) {
    int ret = JSON_ENCODING_OK;

    if (!data) {
        return ERR_FEMD_JSON_NO_VAL_TO_ENCODE;
    }

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_ADC_READING, json_object());
    JsonObj *jAdc = json_object_get(*json, JTAG_ADC_READING);
    if (jAdc) {
        json_object_set_new(jAdc, JTAG_REVERSE_POWER, json_real(data->reverse_power));
        json_object_set_new(jAdc, JTAG_PA_CURRENT, json_real(data->pa_current));
        json_object_set_new(jAdc, JTAG_FEM_UNIT, json_integer(data->fem_unit));
        
        json_object_set_new(jAdc, JTAG_UNITS, json_object());
        JsonObj *jUnits = json_object_get(jAdc, JTAG_UNITS);
        if (jUnits) {
            json_object_set_new(jUnits, JTAG_REVERSE_POWER, json_string("dBm"));
            json_object_set_new(jUnits, JTAG_PA_CURRENT, json_string("A"));
        }
    } else {
        json_decref(*json);
        *json = NULL;
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_eeprom_data(JsonObj **json, const EepromData *data) {
    int ret = JSON_ENCODING_OK;

    if (!data) {
        return ERR_FEMD_JSON_NO_VAL_TO_ENCODE;
    }

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_EEPROM_DATA, json_object());
    JsonObj *jEeprom = json_object_get(*json, JTAG_EEPROM_DATA);
    if (jEeprom) {
        if (data->serial) {
            json_object_set_new(jEeprom, JTAG_SERIAL, json_string(data->serial));
        }
        json_object_set_new(jEeprom, JTAG_FEM_UNIT, json_integer(data->fem_unit));
    } else {
        json_decref(*json);
        *json = NULL;
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_safety_config(JsonObj **json, const SafetyConfig *config) {
    int ret = JSON_ENCODING_OK;

    if (!config) {
        return ERR_FEMD_JSON_NO_VAL_TO_ENCODE;
    }

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_SAFETY_CONFIG, json_object());
    JsonObj *jSafety = json_object_get(*json, JTAG_SAFETY_CONFIG);
    if (jSafety) {
        json_object_set_new(jSafety, JTAG_MAX_REVERSE_POWER, json_real(config->max_reverse_power));
        json_object_set_new(jSafety, JTAG_MAX_CURRENT, json_real(config->max_current));
        json_object_set_new(jSafety, JTAG_MAX_TEMPERATURE, json_real(config->max_temperature));
    } else {
        json_decref(*json);
        *json = NULL;
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_service_info(JsonObj **json, const ServiceInfo *info) {
    int ret = JSON_ENCODING_OK;

    if (!info) {
        return ERR_FEMD_JSON_NO_VAL_TO_ENCODE;
    }

    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_SERVICE_INFO, json_object());
    JsonObj *jService = json_object_get(*json, JTAG_SERVICE_INFO);
    if (jService) {
        if (info->service_name) {
            json_object_set_new(jService, JTAG_SERVICE_NAME, json_string(info->service_name));
        }
        if (info->version) {
            json_object_set_new(jService, JTAG_VERSION, json_string(info->version));
        }
        json_object_set_new(jService, JTAG_UPTIME, json_integer(info->uptime));
    } else {
        json_decref(*json);
        *json = NULL;
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    return ret;
}

int json_deserialize_gpio_control(JsonObj *json, char **gpio_name, bool *enabled) {
    int ret = JSON_DECODING_OK;

    if (!json || !gpio_name || !enabled) {
        return ERR_FEMD_JSON_PARSER;
    }

    JsonObj *jName = json_object_get(json, JTAG_NAME);
    if (jName && json_is_string(jName)) {
        const char *name_str = json_string_value(jName);
        if (name_str) {
            *gpio_name = strdup(name_str);
        }
    }

    JsonObj *jEnabled = json_object_get(json, JTAG_ENABLED);
    if (jEnabled && json_is_boolean(jEnabled)) {
        *enabled = json_boolean_value(jEnabled);
    }

    return ret;
}

int json_deserialize_dac_control(JsonObj *json, float *carrier_voltage, float *peak_voltage) {
    int ret = JSON_DECODING_OK;

    if (!json || !carrier_voltage || !peak_voltage) {
        return ERR_FEMD_JSON_PARSER;
    }

    JsonObj *jCarrier = json_object_get(json, JTAG_CARRIER_VOLTAGE);
    if (jCarrier && json_is_real(jCarrier)) {
        *carrier_voltage = (float)json_real_value(jCarrier);
    }

    JsonObj *jPeak = json_object_get(json, JTAG_PEAK_VOLTAGE);
    if (jPeak && json_is_real(jPeak)) {
        *peak_voltage = (float)json_real_value(jPeak);
    }

    return ret;
}

int json_deserialize_temperature_threshold(JsonObj *json, float *threshold) {
    int ret = JSON_DECODING_OK;

    if (!json || !threshold) {
        return ERR_FEMD_JSON_PARSER;
    }

    JsonObj *jThreshold = json_object_get(json, JTAG_THRESHOLD);
    if (jThreshold && json_is_real(jThreshold)) {
        *threshold = (float)json_real_value(jThreshold);
    }

    return ret;
}

int json_deserialize_eeprom_write(JsonObj *json, char **serial) {
    int ret = JSON_DECODING_OK;

    if (!json || !serial) {
        return ERR_FEMD_JSON_PARSER;
    }

    JsonObj *jSerial = json_object_get(json, JTAG_SERIAL);
    if (jSerial && json_is_string(jSerial)) {
        const char *serial_str = json_string_value(jSerial);
        if (serial_str) {
            *serial = strdup(serial_str);
        }
    }

    return ret;
}

int json_deserialize_safety_config(JsonObj *json, SafetyConfig *config) {
    int ret = JSON_DECODING_OK;

    if (!json || !config) {
        return ERR_FEMD_JSON_PARSER;
    }

    JsonObj *jMaxPower = json_object_get(json, JTAG_MAX_REVERSE_POWER);
    if (jMaxPower && json_is_real(jMaxPower)) {
        config->max_reverse_power = (float)json_real_value(jMaxPower);
    }

    JsonObj *jMaxCurrent = json_object_get(json, JTAG_MAX_CURRENT);
    if (jMaxCurrent && json_is_real(jMaxCurrent)) {
        config->max_current = (float)json_real_value(jMaxCurrent);
    }

    JsonObj *jMaxTemp = json_object_get(json, JTAG_MAX_TEMPERATURE);
    if (jMaxTemp && json_is_real(jMaxTemp)) {
        config->max_temperature = (float)json_real_value(jMaxTemp);
    }

    return ret;
}

void json_free(JsonObj **json) {
    if (*json) {
        json_decref(*json);
        *json = NULL;
    }
}

void json_log(JsonObj *json) {
    if (json) {
        char *json_str = json_dumps(json, JSON_COMPACT);
        if (json_str) {
            usys_log_debug("JSON: %s", json_str);
            free(json_str);
        }
    }
}