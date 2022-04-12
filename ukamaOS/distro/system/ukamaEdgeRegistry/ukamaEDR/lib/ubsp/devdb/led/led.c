/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/led/led.h"

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "inc/devicefxn.h"
#include "headers/utils/log.h"
#include "devdb/led/ledtricol.h"

#include <string.h>

static ListInfo leddb;
static int leddbflag = 0;

const DevFxnTable led_tricol_fxn_table = {
    .init = led_tricol_init,
    .registration = led_tricol_registration,
    .read_prop_count = led_tricol_read_prop_count,
    .read_prop = led_tricol_read_properties,
    .configure = led_tricol_configure,
    .read = led_tricol_read,
    .write = led_tricol_write,
    .enable = led_tricol_enable,
    .disable = led_tricol_disable,
    .register_cb = NULL,
    .dregister_cb = NULL,
    .enable_irq = NULL,
    .disable_irq = NULL,
    .confirm_irq = NULL,
    .irq_type = NULL
};

DevFxnMap led_dev_map[MAX_LED_SENSOR_TYPE] = {
    { .name = "LED-TRICOLOR", .fxn_table = &led_tricol_fxn_table }
};

const DevFxnTable *get_dev_led_fxn_tbl(char *name) {
    const DevFxnTable *fxn_tbl = NULL;
    for (uint8_t iter = 0; iter < MAX_LED_SENSOR_TYPE; iter++) {
        if (!strcmp(name, led_dev_map[iter].name)) {
            fxn_tbl = led_dev_map[iter].fxn_table;
            break;
        }
    }
    return fxn_tbl;
}

static void free_led_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            free(node->data);
        }
        free(node);
    }
}

int compare_led_dev(void *ipt, void *sd) {
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

ListInfo *get_dev_led_db() {
    /* Initialize DB for the first time we try to access it.*/
    if (leddbflag == 0) {
        list_new(&leddb, sizeof(Device), free_led_dev, compare_led_dev, NULL);
        leddbflag = 1;
        log_trace("LED:: led DB initialized.");
    }
    return &leddb;
}
