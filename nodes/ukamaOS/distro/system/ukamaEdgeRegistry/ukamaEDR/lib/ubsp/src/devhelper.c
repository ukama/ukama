/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "headers/errorcode.h"
#include "inc/globalheader.h"
#include "headers/ubsp/property.h"
#include "utils/irqdb.h"
#include "utils/irqhelper.h"
#include "headers/utils/log.h"
#include "devdb/sysfs/drvrsysfs.h"

/* dev ->hw_aatr type is still not known at this will can only be type casted by driver file.*/
#define HWATTR_GET(hwattr, dev, prop)                         \
    {                                                         \
        if (IF_SYSFS_SUPPORT(dev->sysfile)) {                 \
            char sysf[MAX_PATH_LENGTH] = { '\0' };            \
            memcpy(sysf, dev->sysfile, strlen(dev->sysfile)); \
            strcat(sysf, prop->sysfname);                     \
            hwattr = sysf;                                    \
        } else {                                              \
            hwattr = dev->hw_attr;                            \
        }                                                     \
    }

int dhelper_validate_property(Property *prop, int pidx) {
    int ret = 0;
    if (prop[pidx].available == PROP_AVAIL) {
        ret = 1;
    }
    return ret;
}

int dhelper_init_property_from_parser(Device *p_dev, Property **prop,
                                      int *count) {
    int ret = 0;
    *count = get_property_count(p_dev->obj.name);
    if (*count < 0) {
        ret = -1;
        log_debug("DEVHELEPR: No json property found for device name %s "
                  "desc %s module UUID %s.",
                  p_dev->obj.name, p_dev->obj.disc, p_dev->obj.mod_UUID);
    } else {
        *prop = get_property_table(p_dev->obj.name);
    }
    return ret;
}

/* Initializes driver for the device.*/
int dhelper_init_driver(const DrvDBFxnTable *drvr, Device *dev) {
    int ret = 0;
    if (drvr) {
        if (drvr->init) {
            ret = drvr->init(dev);
        }
    }
    return ret;
}

int dhelper_registration(const DrvDBFxnTable *drvr, Device *p_dev) {
    return 0;
}

void dhelper_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
}

/* Configuring device properties. */
int dhelper_configure(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
                      int pidx, void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (dhelper_validate_property(prop, pidx)) {
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->configure) {
                ret = drvr->configure(hwattr, pdata, data);
            }
        } else {
            ret = ERR_UBSP_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_UBSP_DEV_PROPERTY_MISSING;
    }
    return ret;
}

/* Performing read on device properties. */
int dhelper_read(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
                 int pidx, void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (dhelper_validate_property(prop, pidx)) {
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->read) {
                ret = drvr->read(hwattr, pdata, data);
            }
        } else {
            ret = ERR_UBSP_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_UBSP_DEV_PROPERTY_MISSING;
    }
    return ret;
}

/* Performing write to device properties. */
int dhelper_write(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
                  int pidx, void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (dhelper_validate_property(prop, pidx)) {
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->write) {
                ret = drvr->write(hwattr, pdata, data);
            }
        } else {
            ret = ERR_UBSP_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_UBSP_DEV_PROPERTY_MISSING;
    }
    return ret;
}

/* Enable device .*/
int dhelper_enable(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
                   int pidx, void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (dhelper_validate_property(prop, pidx)) {
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->enable) {
                ret = drvr->enable(hwattr, pdata, data);
            }
        } else {
            ret = ERR_UBSP_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_UBSP_DEV_PROPERTY_MISSING;
    }
    return ret;
}

/* Disable device .*/
int dhelper_disable(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
                    int pidx, void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (dhelper_validate_property(prop, pidx)) {
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->disable) {
                ret = drvr->disable(hwattr, pdata, data);
            }
        } else {
            ret = ERR_UBSP_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_UBSP_DEV_PROPERTY_MISSING;
    }
    return ret;
}

/* Enabling IRQ by registering callback and starting thread for IRQ.*/
int dhelper_enable_irq(const DrvDBFxnTable *drvr, SensorCallbackFxn sensor_cb,
                       Device *dev, Property *prop, int pidx, void *data) {
    //TODO: check if IRQ has to enable and disabled here or in driver layer below it.
    int ret = 0;
    void *hwattr = NULL;
    if (dhelper_validate_property(prop, pidx)) {
        /* IRQSrc config for interrupts.*/
        IRQSrcInfo *irq_src = malloc(sizeof(IRQSrcInfo));
        if (irq_src) {
            memset(irq_src, '\0', sizeof(IRQSrcInfo));
            Property *aprop = &prop[pidx];
            HWATTR_GET(hwattr, dev, aprop);
            if (hwattr) {
                strcpy(irq_src->src.sysfs_name, hwattr);
                irq_src->type = IRQ_SYSFS;
            } else {
                /* Still need to figure out about drivers hwattr */
            }
            memcpy(&irq_src->obj, &dev->obj, sizeof(DevObj));
            /* Enabling IRQ */
            ret = irqdb_register_for_device_irq(irq_src, sensor_cb, NULL);
            if (ret) {
                log_error(
                    "DEVHELPER(%d):: Failed to register IRQ for Device name: %s disc: %s module UUID %s",
                    ret, dev->obj.name, dev->obj.disc, dev->obj.mod_UUID);
            }

            if (irq_src) {
                free(irq_src);
                irq_src = NULL;
            }
        } else {
            ret = ERR_UBSP_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_UBSP_DEV_PROPERTY_MISSING;
    }
    return ret;
}

/* Disable IRQ by de-registering callback and stopping thread for IRQ.*/
int dhelper_disable_irq(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
                        int pidx, void *data) {
    int ret = 0;
    void *hwattr = NULL;
    if (dhelper_validate_property(prop, pidx)) {
        /* IRQSrc config for interrupts.*/
        IRQSrcInfo *irq_src = malloc(sizeof(IRQSrcInfo));
        if (irq_src) {
            memset(irq_src, '\0', sizeof(IRQSrcInfo));
            Property *aprop = &prop[pidx];
            HWATTR_GET(hwattr, dev, aprop);
            if (hwattr) {
                strcpy(irq_src->src.sysfs_name, hwattr);
                irq_src->type = IRQ_SYSFS;
            } else {
                /* Still need to figure out about drivers hwattr */
            }
            memcpy(&irq_src->obj, &dev->obj, sizeof(DevObj));
            /* Disable IRQ */
            ret = irqdb_deregister_for_device_irq(irq_src, NULL);
            if (ret) {
                log_error(
                    "DEVHELPER(%d):: Failed to register IRQ for Device name: %s disc: %s module UUID %s",
                    ret, dev->obj.name, dev->obj.disc, dev->obj.mod_UUID);
            }

            if (irq_src) {
                free(irq_src);
                irq_src = NULL;
            }
        } else {
            ret = ERR_UBSP_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_UBSP_DEV_PROPERTY_MISSING;
    }
    return ret;
}

/* Reading and confirming interrupts for ADT7481 device */
int dhelper_confirm_irq(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
                        AlertCallBackData **acbdata, char *fpath, int maxpcount,
                        int *evt) {
    int ret = 0;
    if (dev && fpath) {
        ret = irqhelper_confirm_irq(drvr, dev, acbdata, prop, maxpcount, fpath,
                                    evt);
    } else {
        ret = ERR_UBSP_INVALID_POINTER;
    }
    return ret;
}
