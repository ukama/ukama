/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/gpio_type.h"

#include "device.h"
#include "device_ops.h"
#include "errorcode.h"
#include "devices/bsp_gpio.h"

#include "usys_list.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static ListInfo gpioLdgr;
static int gpioLdgrflag = 0;

const DevOps gpioOps = { .init = bsp_gpio_init,
                                     .registration = bsp_gpio_registration,
                                     .readPropCount = bsp_gpio_read_prop_count,
                                     .readProp = bsp_gpio_read_properties,
                                     .configure = bsp_gpio_configure,
                                     .read = bsp_gpio_read,
                                     .write = bsp_gpio_write,
                                     .enable = bsp_gpio_enable,
                                     .disable = bsp_gpio_disable,
                                     .registerCb = NULL,
                                     .dregisterCb = NULL,
                                     .enableIrq = NULL,
                                     .disableIrq = NULL,
                                     .confirmIrq = NULL,
                                     .irqType = NULL };

DevOpsMap gpio_dev_map[MAX_GPIO_SENSOR_TYPE] = {
    { .name = "GPIO", .opsTable = &gpioOps }
};

const DevOps *get_dev_gpio_type_opsTbl(char *name) {
    const DevOps *opsTbl = NULL;
    for (uint8_t iter = 0; iter < MAX_GPIO_SENSOR_TYPE; iter++) {
        if (!usys_strcmp(name, gpio_dev_map[iter].name)) {
            opsTbl = gpio_dev_map[iter].opsTable;
            break;
        }
    }
    return opsTbl;
}

static void free_gpio_type_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            usys_free(node->data);
        }
        usys_free(node);
    }
}

int compare_gpio_type_dev(void *ipt, void *sd) {
    Device *ip = (Device *)ipt;
    Device *op = (Device *)sd;
    int ret = 0;
    /* If module if  and device name, disc, type matches it means devices is same.*/
    if (!usys_strcmp(ip->obj.modUuid, op->obj.modUuid) &&
        !usys_strcmp(ip->obj.name, op->obj.name) &&
        !usys_strcmp(ip->obj.desc, op->obj.desc) && (ip->obj.type == op->obj.type)) {
        ret = 1;
    }
    return ret;
}

ListInfo *get_gpio_type_dev_ldgr() {
    /* Initialize ledger for the first time we try to access it.*/
    if (gpioLdgrflag == 0) {
        usys_list_new(&gpioLdgr, sizeof(Device), free_gpio_type_dev, compare_gpio_type_dev,
                 NULL);
        gpioLdgrflag = 1;
        usys_log_trace("GPIO:: GPIO ledger initialized.");
    }
    return &gpioLdgr;
}
