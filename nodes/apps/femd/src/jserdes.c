/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "jserdes.h"
#include "json_types.h"

#include "usys_log.h"

static const char *op_state_str(OpState s) {

    switch (s) {
    case OpStateQueued:  return "queued";
    case OpStateRunning: return "running";
    case OpStateDone:    return "done";
    case OpStateFailed:  return "failed";
    default:             return "unknown";
    }
}

void json_log(json_t *json) {

    char *str;

    str = json_dumps(json, 0);
    if (str) {
        log_debug("json str: %s", str);
        free(str);
    }
}

static bool get_json_entry(json_t *json, char *key, json_type type,
                           char **strValue, int *intValue,
                           double *doubleValue) {

    json_t *jEntry;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        log_error("Missing %s key in json", key);
        return USYS_FALSE;
    }

    switch (type) {
    case (JSON_STRING):
        *strValue = strdup(json_string_value(jEntry));
        break;
    case (JSON_INTEGER):
        *intValue = (int)json_integer_value(jEntry);
        break;
    case (JSON_REAL):
        *doubleValue = json_real_value(jEntry);
        break;
    default:
        log_error("Invalid type for json key-value: %d", type);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

bool json_deserialize_node_info(char **data,
                                int  *iData,
                                char *tag,
                                json_type type,
                                json_t *json) {

    json_t *jNodeInfo;

    if (json == NULL) return USYS_FALSE;

    jNodeInfo = json_object_get(json, JTAG_NODE_INFO);
    if (jNodeInfo == NULL) {
        log_error("Missing mandatory %s from JSON", JTAG_NODE_INFO);
        return USYS_FALSE;
    }

    if (get_json_entry(jNodeInfo, tag, type, data, iData, NULL) == USYS_FALSE) {
        log_error("Error deserializing node info. tag: %s", tag);
        json_log(json);
        *data = NULL;
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

bool json_serialize_pa_alarm_notification(json_t **json,
                                          Config *config,
                                          int type) {

    *json = json_object();
    if (!*json) return USYS_FALSE;

    json_object_set_new(*json, JTAG_SERVICE_NAME, json_string(config->serviceName));

    if (type == ALARM_TYPE_PA_OFF) {
        json_object_set_new(*json, JTAG_SEVERITY, json_string(ALARM_HIGH));
        json_object_set_new(*json, JTAG_VALUE,    json_string(ALARM_PA_AUTO_OFF));
        json_object_set_new(*json, JTAG_DETAILS,  json_string(ALARM_PA_AUTO_OFF_DESCRP));
    } else if (type == ALARM_TYPE_PA_ON) {
        json_object_set_new(*json, JTAG_SEVERITY, json_string(ALARM_LOW));
        json_object_set_new(*json, JTAG_VALUE,    json_string(ALARM_PA_AUTO_ON));
        json_object_set_new(*json, JTAG_DETAILS,  json_string(ALARM_PA_AUTO_ON_DESCRP));
    }

    json_object_set_new(*json, JTAG_TIME,   json_integer(time(NULL)));
    json_object_set_new(*json, JTAG_MODULE, json_string(MODULE_FEM));
    json_object_set_new(*json, JTAG_NAME,   json_string(ALARM_NODE));
    json_object_set_new(*json, JTAG_UNITS,  json_string(EMPTY_STRING));

    return USYS_TRUE;
}

bool json_serialize_op_id(json_t **json, uint64_t opId) {

    *json = json_object();
    if (!*json) return USYS_FALSE;

    json_object_set_new(*json, JTAG_OP_ID, json_integer((json_int_t)opId));
    return USYS_TRUE;
}

bool json_serialize_op_status(json_t **json, const OpStatus *st) {

    json_t *j;

    if (!st) return USYS_FALSE;

    *json = json_object();
    if (!*json) return USYS_FALSE;

    j = json_object();
    if (!j) {
        json_decref(*json);
        *json = NULL;
        return USYS_FALSE;
    }

    json_object_set_new(j, JTAG_OP_ID,    json_integer((json_int_t)st->opId));
    json_object_set_new(j, JTAG_OP_STATE, json_string(op_state_str(st->state)));
#if 0
    xxx
    json_object_set_new(j, JTAG_RC,       json_integer(st->rc));
    json_object_set_new(j, JTAG_UPDATED,  json_integer(st->updatedTsMs));
#endif

    json_object_set_new(*json, JTAG_OP, j);
    return USYS_TRUE;
}

static json_t *gpio_status_json(const GpioStatus *g) {

    json_t *j;

    if (!g) return NULL;

    j = json_object();
    if (!j) return NULL;

    json_object_set_new(j, JTAG_TX_RF_ENABLE,   json_boolean(g->tx_rf_enable));
    json_object_set_new(j, JTAG_RX_RF_ENABLE,   json_boolean(g->rx_rf_enable));
    json_object_set_new(j, JTAG_PA_VDS_ENABLE,  json_boolean(g->pa_vds_enable));
    json_object_set_new(j, JTAG_RF_PAL_ENABLE,  json_boolean(g->rf_pal_enable));
    json_object_set_new(j, JTAG_VDS_28V_ENABLE, json_boolean(!g->pa_disable));
    json_object_set_new(j, JTAG_PSU_PGOOD,      json_boolean(g->psu_pgood));

    return j;
}

bool json_serialize_ctrl_snapshot(json_t **json, const CtrlSnapshot *s) {

    if (!s) return USYS_FALSE;

    *json = json_object();
    if (!*json) return USYS_FALSE;

    json_object_set_new(*json, JTAG_TIMESTAMP, json_integer(s->sampleTsMs));

    if (s->haveTemp) {
        json_object_set_new(*json, JTAG_TEMPERATURE, json_real(s->tempC));
    }

    return USYS_TRUE;
}

bool json_serialize_fem_snapshot(json_t **json, FemUnit unit, const FemSnapshot *s) {

    json_t *jGpio;
    json_t *jDac;
    json_t *jAdc;

    if (!s) return USYS_FALSE;

    *json = json_object();
    if (!*json) return USYS_FALSE;

    json_object_set_new(*json, JTAG_FEM_UNIT,  json_integer(unit));
    json_object_set_new(*json, JTAG_TIMESTAMP, json_integer(s->sampleTsMs));

    if (s->haveGpio) {
        jGpio = gpio_status_json(&s->gpio);
        if (!jGpio) goto fail;
        json_object_set_new(*json, JTAG_GPIO_STATUS, jGpio);
    }

    if (s->haveTemp) {
        json_object_set_new(*json, JTAG_TEMPERATURE, json_real(s->tempC));
    }

    if (s->haveAdc) {
        jAdc = json_object();
        if (!jAdc) goto fail;

        json_object_set_new(jAdc, JTAG_REVERSE_POWER, json_real(s->reversePowerDbm));
        json_object_set_new(jAdc, JTAG_FORWARD_POWER, json_real(s->forwardPowerDbm));
        json_object_set_new(jAdc, JTAG_PA_CURRENT,    json_real(s->paCurrentA));

        json_object_set_new(*json, JTAG_ADC_READING, jAdc);
    }

    if (s->haveDac) {
        jDac = json_object();
        if (!jDac) goto fail;

        json_object_set_new(jDac, JTAG_CARRIER_VOLTAGE, json_real(s->carrierVoltage));
        json_object_set_new(jDac, JTAG_PEAK_VOLTAGE,    json_real(s->peakVoltage));

        json_object_set_new(*json, JTAG_DAC_CONFIG, jDac);
    }

    if (s->haveSerial) {
        json_object_set_new(*json, JTAG_SERIAL, json_string(s->serial));
    }

    return USYS_TRUE;

fail:
    if (*json) {
        json_decref(*json);
        *json = NULL;
    }
    return USYS_FALSE;
}

void json_free(json_t **json) {
    if (*json) {
        json_decref(*json);
        *json = NULL;
    }
}
