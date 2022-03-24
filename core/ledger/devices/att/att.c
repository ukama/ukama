/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/att.h"

#include "device.h"
#include "device_ops.h"
#include "errorcode.h"
#include "devices/bsp_dat31r5a.h"

#include "usys_list.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"


static ListInfo attLdgr;
static int attLdgrflag = 0;

const DevOps dat31r5aOps = { .init = bsp_dat31r5a_init,
                                         .registration = bsp_dat31r5a_registration,
                                         .readPropCount =
                                             bsp_dat31r5a_read_prop_count,
                                         .readProp = bsp_dat31r5a_read_properties,
                                         .configure = bsp_dat31r5a_configure,
                                         .read = bsp_dat31r5a_read,
                                         .write = bsp_dat31r5a_write,
                                         .enable = bsp_dat31r5a_enable,
                                         .disable = bsp_dat31r5a_disable,
                                         .registerCb = NULL,
                                         .dregisterCb = NULL,
                                         .enableIrq = NULL,
                                         .disableIrq = NULL,
                                         .confirmIrq = NULL,
                                         .irqType = NULL };

DevOpsMap att_dev_map[MAX_ATT_SENSOR_TYPE] = {
    { .name = "DAT-31R5A-PP", .opsTable = &dat31r5aOps }
};

const DevOps *get_att_dev_ops(char *name) {
    const DevOps *opsTbl = NULL;
    for (uint8_t iter = 0; iter < MAX_ATT_SENSOR_TYPE; iter++) {
        if (!usys_strcasecmp(name, att_dev_map[iter].name)) {
            opsTbl = att_dev_map[iter].opsTable;
            break;
        }
    }
    return opsTbl;
}

static void free_att_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            usys_free(node->data);
        }
        usys_free(node);
    }
}

int compare_att_dev(void *ipt, void *sd) {
    Device *ip = (Device *)ipt;
    Device *op = (Device *)sd;
    int ret = 0;
    /* If module if  and device name, desc, type matches it means devices is same.*/
    if (!usys_strcasecmp(ip->obj.modUuid, op->obj.modUuid) &&
        !usys_strcasecmp(ip->obj.name, op->obj.name) &&
        !usys_strcasecmp(ip->obj.desc, op->obj.desc) && (ip->obj.type == op->obj.type)) {
        ret = 1;
    }
    return ret;
}

ListInfo *get_att_dev_ldgr() {
    /* Initialize DB for the first time we try to access it.*/
    if (attLdgrflag == 0) {
        usys_list_new(&attLdgr, sizeof(Device), free_att_dev, compare_att_dev, NULL);
        attLdgrflag = 1;
        usys_log_trace("ATT:: Attenuation ledger initialized.");
    }
    return &attLdgr;
}
