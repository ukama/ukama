/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/att/att.h"

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "inc/devicefxn.h"
#include "headers/utils/log.h"
#include "devdb/att/dat31r5a.h"

#include <string.h>

static ListInfo attdb;
static int attdbflag = 0;

const DevFxnTable dat31r5a_fxn_table = { .init = dat31r5a_init,
                                         .registration = dat31r5a_registration,
                                         .read_prop_count =
                                             dat31r5a_read_prop_count,
                                         .read_prop = dat31r5a_read_properties,
                                         .configure = dat31r5a_configure,
                                         .read = dat31r5a_read,
                                         .write = dat31r5a_write,
                                         .enable = dat31r5a_enable,
                                         .disable = dat31r5a_disable,
                                         .register_cb = NULL,
                                         .dregister_cb = NULL,
                                         .enable_irq = NULL,
                                         .disable_irq = NULL,
                                         .confirm_irq = NULL,
                                         .irq_type = NULL };

DevFxnMap att_dev_map[MAX_ATT_SENSOR_TYPE] = {
    { .name = "DAT-31R5A-PP", .fxn_table = &dat31r5a_fxn_table }
};

const DevFxnTable *get_dev_att_fxn_tbl(char *name) {
    const DevFxnTable *fxn_tbl = NULL;
    for (uint8_t iter = 0; iter < MAX_ATT_SENSOR_TYPE; iter++) {
        if (!strcmp(name, att_dev_map[iter].name)) {
            fxn_tbl = att_dev_map[iter].fxn_table;
            break;
        }
    }
    return fxn_tbl;
}

static void free_att_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            free(node->data);
        }
        free(node);
    }
}

int compare_att_dev(void *ipt, void *sd) {
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

ListInfo *get_dev_att_db() {
    /* Initialize DB for the first time we try to access it.*/
    if (attdbflag == 0) {
        list_new(&attdb, sizeof(Device), free_att_dev, compare_att_dev, NULL);
        attdbflag = 1;
        log_trace("ATT:: Attenuation DB initialized.");
    }
    return &attdb;
}
