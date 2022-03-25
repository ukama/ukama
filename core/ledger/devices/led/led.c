/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/led.h"

#include "device.h"
#include "device_ops.h"
#include "errorcode.h"
#include "devices/bsp_ledtricol.h"

#include "usys_list.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static ListInfo ledLdgr;
static int ledLdgrflag = 0;

const DevOps ledTriColOps = { .init = bsp_led_tricol_init,
                              .registration = bsp_led_tricol_registration,
                              .readPropCount = bsp_led_tricol_read_prop_count,
                              .readProp = bsp_led_tricol_read_properties,
                              .configure = bsp_led_tricol_configure,
                              .read = bsp_led_tricol_read,
                              .write = bsp_led_tricol_write,
                              .enable = bsp_led_tricol_enable,
                              .disable = bsp_led_tricol_disable,
                              .registerCb = NULL,
                              .dregisterCb = NULL,
                              .enableIrq = NULL,
                              .disableIrq = NULL,
                              .confirmIrq = NULL,
                              .irqType = NULL };

DevOpsMap ledDevMap[MAX_LED_SENSOR_TYPE] = { { .name = "LED-TRICOLOR",
                                                 .opsTable = &ledTriColOps } };

const DevOps *get_led_dev_ops(char *name) {
    const DevOps *opsTbl = NULL;
    for (uint8_t iter = 0; iter < MAX_LED_SENSOR_TYPE; iter++) {
        if (!usys_strcasecmp(name, ledDevMap[iter].name)) {
            opsTbl = ledDevMap[iter].opsTable;
            break;
        }
    }
    return opsTbl;
}

static void free_led_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            usys_free(node->data);
        }
        usys_free(node);
    }
}

static int compare_led_dev(void *ipt, void *sd) {
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

ListInfo *get_led_dev_ldgr() {
    /* Initialize ledger for the first time we try to access it.*/
    if (ledLdgrflag == 0) {
        usys_list_new(&ledLdgr, sizeof(Device), free_led_dev, compare_led_dev,
                      NULL);
        ledLdgrflag = 1;
        usys_log_trace("LED:: led ledger initialized.");
    }
    return &ledLdgr;
}
