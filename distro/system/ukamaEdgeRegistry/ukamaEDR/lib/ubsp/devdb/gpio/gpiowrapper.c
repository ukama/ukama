/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/gpio/gpiowrapper.h"

#include "inc/devicefxn.h"
#include "headers/utils/log.h"
#include "devdb/gpio/gpio.h"
#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"

#include <string.h>

static ListInfo gpiodb;
static int gpiodbflag = 0;

const DevFxnTable gpio_fxn_table = { .init = gpio_init,
                                     .registration = gpio_registration,
                                     .read_prop_count = gpio_read_prop_count,
                                     .read_prop = gpio_read_properties,
                                     .configure = gpio_configure,
                                     .read = gpio_read,
                                     .write = gpio_write,
                                     .enable = gpio_enable,
                                     .disable = gpio_disable,
                                     .register_cb = NULL,
                                     .dregister_cb = NULL,
                                     .enable_irq = NULL,
                                     .disable_irq = NULL,
                                     .confirm_irq = NULL,
                                     .irq_type = NULL };

DevFxnMap gpio_dev_map[MAX_GPIO_SENSOR_TYPE] = {
    { .name = "GPIO", .fxn_table = &gpio_fxn_table }
};

const DevFxnTable *get_dev_gpiow_fxn_tbl(char *name) {
    const DevFxnTable *fxn_tbl = NULL;
    for (uint8_t iter = 0; iter < MAX_GPIO_SENSOR_TYPE; iter++) {
        if (!strcmp(name, gpio_dev_map[iter].name)) {
            fxn_tbl = gpio_dev_map[iter].fxn_table;
            break;
        }
    }
    return fxn_tbl;
}

static void free_gpiow_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            free(node->data);
        }
        free(node);
    }
}

int compare_gpiow_dev(void *ipt, void *sd) {
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

ListInfo *get_dev_gpiow_db() {
    /* Initialize DB for the first time we try to access it.*/
    if (gpiodbflag == 0) {
        list_new(&gpiodb, sizeof(Device), free_gpiow_dev, compare_gpiow_dev,
                 NULL);
        gpiodbflag = 1;
        log_trace("GPIO:: GPIO DB initialized.");
    }
    return &gpiodb;
}
