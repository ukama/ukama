/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devhelper.h"

#include "errorcode.h"
#include "irqdb.h"
#include "irqhelper.h"
#include "noded_macros.h"
#include "property.h"
#include "drivers/sysfs_wrapper.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

/* dev ->hw_aatr type is still not known at this will can only be type casted by driver file.*/
#define HWATTR_GET(hwattr, dev, prop)                                   \
    {                                                                   \
        if (IF_SYSFS_SUPPORT(dev->sysFile)) {                           \
            char sysf[MAX_PATH_LENGTH] = { '\0' };                      \
            usys_memcpy(sysf, dev->sysFile, usys_strlen(dev->sysFile)); \
            usys_strcat(sysf, prop->sysFname);                          \
            hwattr = sysf;                                              \
        } else {                                                        \
            hwattr = dev->hwAttr;                                       \
        }                                                               \
    }

int dhelper_validate_property(Property *prop, int pidx) {
    int ret = -1;
    if (prop[pidx].available == PROP_AVAIL) {
        ret = 0;
    }
    return ret;
}

int dhelper_validate_permissions(Property *prop, int pidx, uint16_t perm) {
    int ret = -1;
    if (prop[pidx].perm & perm) {
        ret = 0;
    }
    return ret;
}

int dhelper_validate_property_type_alert(Property *prop, int pidx) {
    int ret = -1;
    if (prop[pidx].propType == PROP_TYPE_ALERT) {
        ret = 0;
    }
    return ret;
}

int dhelper_init_property_from_parser(Device *p_dev, Property **prop,
                                      int *count) {
    int ret = 0;
    *count = get_property_count(p_dev->obj.name);
    if (*count < 0) {
        ret = -1;
        usys_log_debug("DEVHELEPR: No json property found for device name %s "
                       "desc %s module UUID %s.",
                       p_dev->obj.name, p_dev->obj.desc, p_dev->obj.modUuid);
    } else {
        *prop = get_property_table(p_dev->obj.name);
    }
    return ret;
}

/* Initializes driver for the device.*/
int dhelper_init_driver(const DrvrOps *drvr, Device *dev) {
    int ret = 0;
    if (drvr) {
        if (drvr->init) {
            ret = drvr->init(dev);
        }
    }
    return ret;
}

int dhelper_registration(const DrvrOps *drvr, Device *p_dev) {
    return 0;
}

void dhelper_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
}

/* Configuring device properties. */
int dhelper_configure(const DrvrOps *drvr, Device *dev, Property *prop,
                      int pidx, void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (!(dhelper_validate_property(prop, pidx))) {
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->configure) {
                ret = drvr->configure(hwattr, pdata, data);
            }
        } else {
            ret = ERR_NODED_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_NODED_DEV_PROPERTY_MARKED_NOT_AVAILABLE;
    }
    return ret;
}

/* Performing read on device properties. */
int dhelper_read(const DrvrOps *drvr, Device *dev, Property *prop, int pidx,
                 void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (!(dhelper_validate_property(prop, pidx))) {
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->read) {
                ret = drvr->read(hwattr, pdata, data);
            }
        } else {
            ret = ERR_NODED_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_NODED_DEV_PROPERTY_MARKED_NOT_AVAILABLE;
    }
    return ret;
}

/* Performing write to device properties. */
int dhelper_write(const DrvrOps *drvr, Device *dev, Property *prop, int pidx,
                  void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (!(dhelper_validate_property(prop, pidx))) {
        /* check for write permissions */
        if (dhelper_validate_permissions(prop, pidx, PERM_WR)) {
            return ERR_NODED_DEV_PERMISSION_DENIED;
        }
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->write) {
                ret = drvr->write(hwattr, pdata, data);
            }
        } else {
            ret = ERR_NODED_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_NODED_DEV_PROPERTY_MARKED_NOT_AVAILABLE;
    }
    return ret;
}

/* Enable device .*/
int dhelper_enable(const DrvrOps *drvr, Device *dev, Property *prop, int pidx,
                   void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (!(dhelper_validate_property(prop, pidx))) {
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->enable) {
                ret = drvr->enable(hwattr, pdata, data);
            }
        } else {
            ret = ERR_NODED_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_NODED_DEV_PROPERTY_MARKED_NOT_AVAILABLE;
    }
    return ret;
}

/* Disable device .*/
int dhelper_disable(const DrvrOps *drvr, Device *dev, Property *prop, int pidx,
                    void *data) {
    int ret = 0;
    void *hwattr = NULL;
    Property *pdata = &prop[pidx];
    if (!(dhelper_validate_property(prop, pidx))) {
        HWATTR_GET(hwattr, dev, pdata);
        if (hwattr) {
            if (drvr->disable) {
                ret = drvr->disable(hwattr, pdata, data);
            }
        } else {
            ret = ERR_NODED_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_NODED_DEV_PROPERTY_MARKED_NOT_AVAILABLE;
    }
    return ret;
}

/* Enabling IRQ by registering callback and starting thread for IRQ.*/
int dhelper_enable_irq(const DrvrOps *drvr, SensorCallbackFxn sensorCb,
                       Device *dev, Property *prop, int pidx, void *data) {
    //TODO: check if IRQ has to enable and disabled here or in driver layer below it.
    int ret = 0;
    void *hwattr = NULL;
    if (!(dhelper_validate_property(prop, pidx))) {

        /* Check if property is alert type */
        if (dhelper_validate_property_type_alert(prop, pidx)) {
            return ERR_NODED_DEV_PROPERTY_IS_NOT_ALERT_TYPE;
        }

        /* IRQSrc config for interrupts.*/
        IRQSrcInfo *irqSrc = usys_zmalloc(sizeof(IRQSrcInfo));
        if (irqSrc) {
            Property *aprop = &prop[pidx];
            HWATTR_GET(hwattr, dev, aprop);
            if (hwattr) {
                usys_strcpy(irqSrc->src.sysFsName, hwattr);
                irqSrc->type = IRQ_SYSFS;
            } else {
                /* Still need to figure out about drivers hwattr */
            }
            usys_memcpy(&irqSrc->obj, &dev->obj, sizeof(DevObj));
            /* Enabling IRQ */
            ret = irqdb_register_for_device_irq(irqSrc, sensorCb, NULL);
            if (ret) {
                log_error(
                    "DEVHELPER(%d):: Failed to register IRQ for Device name: %s desc: %s module UUID %s",
                    ret, dev->obj.name, dev->obj.desc, dev->obj.modUuid);
            }

            if (irqSrc) {
                usys_free(irqSrc);
                irqSrc = NULL;
            }
        } else {
            ret = ERR_NODED_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_NODED_DEV_PROPERTY_MARKED_NOT_AVAILABLE;
    }
    return ret;
}

/* Disable IRQ by de-registering callback and stopping thread for IRQ.*/
int dhelper_disable_irq(const DrvrOps *drvr, Device *dev, Property *prop,
                        int pidx, void *data) {
    int ret = 0;
    void *hwattr = NULL;
    if (!(dhelper_validate_property(prop, pidx))) {
        /* IRQSrc config for interrupts.*/
        IRQSrcInfo *irqSrc = usys_zmalloc(sizeof(IRQSrcInfo));
        if (irqSrc) {
            Property *aprop = &prop[pidx];
            HWATTR_GET(hwattr, dev, aprop);
            if (hwattr) {
                usys_strcpy(irqSrc->src.sysFsName, hwattr);
                irqSrc->type = IRQ_SYSFS;
            } else {
                /* Still need to figure out about drivers hwattr */
            }
            usys_memcpy(&irqSrc->obj, &dev->obj, sizeof(DevObj));
            /* Disable IRQ */
            ret = irqdb_deregister_for_device_irq(irqSrc, NULL);
            if (ret) {
                log_error(
                    "DEVHELPER(%d):: Failed to deregister IRQ for Device name: %s desc: %s module UUID %s",
                    ret, dev->obj.name, dev->obj.desc, dev->obj.modUuid);
            }

            if (irqSrc) {
                usys_free(irqSrc);
                irqSrc = NULL;
            }
        } else {
            ret = ERR_NODED_DEV_HWATTR_MISSING;
        }
    } else {
        ret = ERR_NODED_DEV_PROPERTY_MARKED_NOT_AVAILABLE;
    }
    return ret;
}

/* Reading and confirming interrupts for device */
int dhelper_confirm_irq(const DrvrOps *drvr, Device *dev, Property *prop,
                        AlertCallBackData **acbdata, char *fpath, int maxPCount,
                        int *evt) {
    int ret = 0;
    if (dev && fpath) {
        ret = irqhelper_confirm_irq(drvr, dev, acbdata, prop, maxPCount, fpath,
                                    evt);
    } else {
        ret = ERR_NODED_INVALID_POINTER;
    }
    return ret;
}
