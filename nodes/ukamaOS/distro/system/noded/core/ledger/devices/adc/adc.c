/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/adc.h"

#include "device.h"
#include "device_ops.h"
#include "errorcode.h"
#include "devices/bsp_ads1015.h"

#include "usys_list.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static ListInfo adcLdgr;
static int adcLdgrflag = 0;

const DevOps ads1015Ops = { .init = bsp_ads1015_init,
                            .registration = bsp_ads1015_registration,
                            .readPropCount = bsp_ads1015_read_prop_count,
                            .readProp = bsp_ads1015_read_properties,
                            .configure = bsp_ads1015_configure,
                            .read = bsp_ads1015_read,
                            .write = bsp_ads1015_write,
                            .enable = bsp_ads1015_enable,
                            .disable = bsp_ads1015_disable,
                            .registerCb = NULL,
                            .dregisterCb = NULL,
                            .enableIrq = NULL,
                            .disableIrq = NULL,
                            .confirmIrq = NULL,
                            .irqType = NULL };

DevOpsMap adcDevMap[MAX_ADC_SENSOR_TYPE] = { { .name = "ADS1015",
                                                 .opsTable = &ads1015Ops } };

const DevOps *get_adc_dev_ops(char *name) {
    const DevOps *opsTbl = NULL;
    for (uint8_t iter = 0; iter < MAX_ADC_SENSOR_TYPE; iter++) {
        if (!usys_strcasecmp(name, adcDevMap[iter].name)) {
            opsTbl = adcDevMap[iter].opsTable;
            break;
        }
    }
    return opsTbl;
}

static void free_adc_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            usys_free(node->data);
        }
        usys_free(node);
    }
}

int compare_adc_dev(void *ipt, void *sd) {
    Device *ip = (Device *)ipt;
    Device *op = (Device *)sd;
    int ret = 0;
    /* If module if  and device name, disc, type matches it means devices is same.*/
    if (!usys_strcasecmp(ip->obj.modUuid, op->obj.modUuid) &&
        !usys_strcasecmp(ip->obj.name, op->obj.name) &&
        !usys_strcasecmp(ip->obj.desc, op->obj.desc) &&
        (ip->obj.type == op->obj.type)) {
        ret = 1;
    }
    return ret;
}

ListInfo *get_adc_dev_ldgr() {
    /* Initialize DB for the first time we try to access it.*/
    if (adcLdgrflag == 0) {
        usys_list_new(&adcLdgr, sizeof(Device), free_adc_dev, compare_adc_dev,
                      NULL);
        adcLdgrflag = 1;
        usys_log_trace("ADC:: ADC Ledger initialized.");
    }
    return &adcLdgr;
}
