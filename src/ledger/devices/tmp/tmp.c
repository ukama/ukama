/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/tmp.h"

#include "device.h"
#include "device_ops.h"
#include "errorcode.h"
#include "devices/bsp_adt7481.h"
#include "devices/bsp_se98.h"
#include "devices/bsp_tmp464.h"

#include "usys_list.h"
#include "usys_log.h"
#include "usys_string.h"

static ListInfo tmpLdgr;
static int tmpLdgrflag = 0;

const DevOps tmp464Ops = { .init =bsp_tmp464_init,
                                       .registration =bsp_tmp464_registration,
                                       .read_prop_count =
                                          bsp_tmp464_read_prop_count,
                                       .read_prop =bsp_tmp464_read_properties,
                                       .configure =bsp_tmp464_configure,
                                       .read =bsp_tmp464_read,
                                       .write =bsp_tmp464_write,
                                       .enable =bsp_tmp464_enable,
                                       .disable =bsp_tmp464_disable,
                                       .registerCb =bsp_tmp464_reg_cb,
                                       .dregisterCb =bsp_tmp464_dreg_cb,
                                       .enableIrq =bsp_tmp464_enable_irq,
                                       .disableIrq =bsp_tmp464_disable_irq,
                                       .confirmIrq =bsp_tmp464_confirm_irq,
                                       .irqType =bsp_tmp464_get_irq_type };

const DevOps se98Ops = { .init = bsp_se98_init,
                                     .registration = bsp_se98_registration,
                                     .read_prop_count = bsp_se98_read_prop_count,
                                     .read_prop = bsp_se98_read_properties,
                                     .configure = bsp_se98_configure,
                                     .read = bsp_se98_read,
                                     .write = bsp_se98_write,
                                     .enable = bsp_se98_enable,
                                     .disable = bsp_se98_disable,
                                     .registerCb = bsp_se98_reg_cb,
                                     .dregisterCb = bsp_se98_dreg_cb,
                                     .enableIrq = bsp_se98_enable_irq,
                                     .disableIrq = bsp_se98_disable_irq,
                                     .confirmIrq = bsp_se98_confirm_irq,
                                     .irqType = bsp_se98_get_irq_type };

const DevOps adt7481Ops = { .init = bsp_adt7481_init,
                                        .registration = bsp_adt7481_registration,
                                        .read_prop_count =
                                            bsp_adt7481_read_prop_count,
                                        .read_prop = bsp_adt7481_read_properties,
                                        .configure = bsp_adt7481_configure,
                                        .read = bsp_adt7481_read,
                                        .write = bsp_adt7481_write,
                                        .enable = bsp_adt7481_enable,
                                        .disable = bsp_adt7481_disable,
                                        .registerCb = bsp_adt7481_reg_cb,
                                        .dregisterCb = bsp_adt7481_dreg_cb,
                                        .enableIrq = bsp_adt7481_enable_irq,
                                        .disableIrq = bsp_adt7481_disable_irq,
                                        .confirmIrq = bsp_adt7481_confirm_irq,
                                        .irqType = bsp_adt7481_get_irq_type };

DevOpsMap tmp_dev_map[MAX_TEMP_SENSOR_TYPE] = {
    { .name = "TMP464", .dev_ops = &tmp464Ops },
    { .name = "ADT7481", .dev_ops = &adt7481Ops },
    { .name = "SE98", .dev_ops = &se98Ops }
};

const DevOps *get_tmp_dev_ops(char *name) {
    const DevOps *devOps = NULL;
    for (uint8_t iter = 0; iter < MAX_TEMP_SENSOR_TYPE; iter++) {
        if (!usys_strcmp(name, tmp_dev_map[iter].name)) {
            devOps = tmp_dev_map[iter].opsTable;
            break;
        }
    }
    return devOps;
}

static void free_tmp_dev(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            usys_free(node->data);
        }
        usys_free(node);
    }
}

int compare_tmp_dev(void *ipt, void *sd) {
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

ListInfo *get_tmp_dev_ldgr() {
    if (tmpLdgrflag == 0) {
        usys_list_new(&tmpLdgr, sizeof(Device), free_tmp_dev, compare_tmp_dev, NULL);
        tmpLdgrflag = 1;
        usys_log_trace("TMP:: TEMP ledger initialized.");
    }
    return &tmpLdgr;
}
