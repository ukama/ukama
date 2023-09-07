/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/adc/adc.h"

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "inc/devicefxn.h"
#include "headers/utils/log.h"
#include "devdb/adc/ads1015.h"

#include <string.h>

static ListInfo adcdb;
static int adcdbflag = 0;

const DevFxnTable ads1015_fxn_table = { .init = ads1015_init,
                                        .registration = ads1015_registration,
                                        .read_prop_count =
                                            ads1015_read_prop_count,
                                        .read_prop = ads1015_read_properties,
                                        .configure = ads1015_configure,
                                        .read = ads1015_read,
                                        .write = ads1015_write,
                                        .enable = ads1015_enable,
                                        .disable = ads1015_disable,
                                        .register_cb = NULL,
                                        .dregister_cb = NULL,
                                        .enable_irq = NULL,
                                        .disable_irq = NULL,
                                        .confirm_irq = NULL,
                                        .irq_type = NULL };

DevFxnMap adc_dev_map[MAX_ADC_SENSOR_TYPE] = {
    { .name = "ADS1015", .fxn_table = &ads1015_fxn_table }
};

const DevFxnTable *get_dev_adc_fxn_tbl(char *name) {
    const DevFxnTable *fxn_tbl = NULL;
    for (uint8_t iter = 0; iter < MAX_ADC_SENSOR_TYPE; iter++) {
        if (!strcmp(name, adc_dev_map[iter].name)) {
            fxn_tbl = adc_dev_map[iter].fxn_table;
            break;
        }
    }
    return fxn_tbl;
}

static void free_adc_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            free(node->data);
        }
        free(node);
    }
}

int compare_adc_dev(void *ipt, void *sd) {
    Device *ip = (Device *)ipt;
    Device *op = (Device *)sd;
    int ret = 0;
    /* If module if  and device name, disc, type matches it means devices is same.*/
    if (!strcmp(ip->obj.mod_UUID, op->obj.mod_UUID) &&
        !strcmp(ip->obj.name, op->obj.name) &&
        !strcmp(ip->obj.disc, op->obj.disc) && (ip->obj.type == op->obj.type)) {
        ret = 1;
    }
    return ret;
}

ListInfo *get_dev_adc_db() {
    /* Initialize DB for the first time we try to access it.*/
    if (adcdbflag == 0) {
        list_new(&adcdb, sizeof(Device), free_adc_dev, compare_adc_dev, NULL);
        adcdbflag = 1;
        log_trace("ADC:: ADC DB initialized.");
    }
    return &adcdb;
}
