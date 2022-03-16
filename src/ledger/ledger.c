/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "ledger.h"

#include "errorcode.h"
#include "irqdb.h"
#include "devices/adc.h"
#include "devices/att.h"
#include "devices/gpio_type.h"
#include "devices/led.h"
#include "devices/pwr.h"
#include "devices/tmp.h"

#include "usys_list.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static const DevOps* get_dev_ops(DevObj *dobj) {
    const DevOps *devOps = NULL;
    switch (dobj->type) {
    case DEV_TYPE_TMP: {
        devOps = get_tmp_dev_ops(dobj->name);
        break;
    }
    case DEV_TYPE_PWR: {
        devOps = get_pwr_dev_ops(dobj->name);
        break;
    }
    case DEV_TYPE_GPIO: {
        devOps = get_gpiow_dev_ops(dobj->name);
        break;
    }
    case DEV_TYPE_LED: {
        devOps = get_led_dev_ops(dobj->name);
        break;
    }
    case DEV_TYPE_ADC: {
        devOps = get_adc_dev_ops(dobj->name);
        break;
    }
    case DEV_TYPE_ATT: {
        devOps = get_att_dev_ops(dobj->name);
        break;
    }
    default: {
    }
    }
    return devOps;
}

static ListInfo *get_dev_ldgr(DeviceType type) {
    ListInfo *devLdgr = NULL;
    switch (type) {
    case DEV_TYPE_TMP: {
        devLdgr = get_tmp_dev_ldgr();
        break;
    }
    case DEV_TYPE_PWR: {
        devLdgr = get_pwr_dev_ldgr();
        break;
    }
    case DEV_TYPE_GPIO: {
        devLdgr = get_gpiow_dev_ldgr();
        break;
    }
    case DEV_TYPE_LED: {
        devLdgr = get_led_dev_ldgr();
        break;
    }
    case DEV_TYPE_ADC: {
        devLdgr = get_adc_dev_ldgr();
        break;
    }
    case DEV_TYPE_ATT: {
        devLdgr = get_att_dev_ldgr();
        break;
    }
    default: {
    }
    }
    return devLdgr;
}

static void ldgr_usys_free(Device *dev) {
    if (dev) {
        if (dev->hwAttr) {
            usys_free(dev->hwAttr);
        }
        usys_free(dev);
    }
}
int compare_dev_node(void *ipt, void *sd) {
    Device *ip = (Device *)ipt;
    Device *op = (Device *)sd;
    int ret = 0;

    /* If module if  and device name, disc, type matches it
     * means devices is same.*/
    if (!usys_strcmp(ip->obj.modUuid, op->obj.modUuid) &&
        !usys_strcmp(ip->obj.name, op->obj.name) &&
        !usys_strcmp(ip->obj.desc, op->obj.desc) &&
        (ip->obj.type == op->obj.type)) {
        ret = 1;
    }

    return ret;
}

/* Searching device in the device list*/
static Device *search_device_object(DevObj *dev_obj) {
    Device *fdev = NULL;

    Device *sdev = usys_zmalloc(sizeof(Device));
    if (sdev) {
        usys_memcpy(&sdev->obj, dev_obj, sizeof(DevObj));

        /* Search return 1 for found.*/
        fdev = list_search(get_dev_ldgr(dev_obj->type), sdev);
        if (fdev) {
            usys_log_trace("Ledger:: Device Name %s, Disc: %s "
                            "Module UUID: %s found.",
                            dev_obj->name, dev_obj->desc, dev_obj->modUuid);
        } else {

            if (fdev) {
                usys_free(fdev);
                fdev = NULL;
            }

            usys_log_debug(
                            "Ledger:: Device Name %s, Disc: %s "
                            "Module UUID: %s not found.",
                            dev_obj->name, dev_obj->desc, dev_obj->modUuid);
        }

        if (sdev) {
            usys_free(sdev);
            sdev = NULL;
        }

    }

    return fdev;
}

static void ldgr_destory() {
    DeviceType type = DEV_TYPE_TMP;

    for (; type < DEV_TYPE_MAX; type++) {
        list_destroy(get_dev_ldgr(type));
        log_warn("Ledger:: Removing DB for 0x%d device type.", type);
    }
}

/* Searching device in the device list*/
static int update_device_db(Device *dev) {
    int ret = 0;
    if (dev) {

        /*Return 0 for update.*/
        ret = list_update(get_dev_ldgr(dev->obj.type), dev);
        if (ret) {
           usys_log_error(
                "Ledger:: Device Name %s, Disc: %s "
                "Module UUID: %s update failed.",
                dev->obj.name, dev->obj.desc, dev->obj.modUuid);
        }

    }
    return ret;
}

void ldgr_irq_callback(void *pcfg) {
    IRQCfg *cfg = pcfg;
    if (!cfg) {
        return;
    }
    DevObj *obj = &cfg->obj;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {

        /* if callback id registered */
        if (dev->devCb) {

            /* if Operation available for device.*/
            if (dev->devOps) {

                const DevOps *devOps = get_dev_ops(&(dev->obj));
                if (devOps->confirmIrq) {
                    AlertCallBackData *acbdata = NULL;
                    int acount = 0;
                    char fName[64] = { '\0' };
                    usys_strcpy(fName, cfg->fName);
                    devOps->confirmIrq(dev, &acbdata, fName, &acount);

                    if (acount > 0) {
                       usys_log_debug(
                            "Ledger:: Calling callback fxn for device name %s,"
                            "disc %s and module UIID %s with %d alerts.",
                            obj->name, obj->desc, obj->modUuid, acount);
                        dev->devCb(obj, &acbdata, &acount);
                    }

                }

            }

        } else {
           usys_log_debug(
                "Ledger:: No Alert callback registered for device name %s, "
                "disc %s and module UIID %s.",
                obj->name, obj->desc, obj->modUuid);
        }
    }

    if (dev) {
        usys_free(dev);
        dev = NULL;
    }

}

static int ldgr_registration(Device *pDev) {
    int ret = 0;
    uint8_t idx = 0;

    /* Search in list if device already exist.*/
    ret = list_if_element_found(get_dev_ldgr(pDev->obj.type), pDev);
    if (ret) {
        ret = 0;
        usys_log_debug(
                        "Ledger:: Device Name: %s, Disc: %s, Type: %d, "
                        "Module Id %s is already in DB.",
                        pDev->obj.name, pDev->obj.desc, pDev->obj.type,
                        pDev->obj.modUuid);
        return 0;
    }

    /* If element not found. It means new device.*/
    if (!pDev->devOps) {

        const DevOps *devOps = get_dev_ops(&(pDev->obj));
        if (devOps) {

            /* Initializing device and registering callback for IRQ's */
            pDev->devOps = devOps;
            if (devOps->init) {
                ret = devOps->init(pDev);
                if (!ret) {
                    if (devOps->registerCb) {
                        ret = devOps->registerCb(pDev,
                                        &ldgr_irq_callback);
                    }
                }

            } else {
                ret = ERR_NODED_DEV_DRVR_MISSING;
            }

        } else {
            ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
        }

    } else {

        /* Should not come here*/
        const DevOps *devOps = pDev->devOps;
        if (devOps->init) {
            ret = devOps->init(pDev);
        } else {
            ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
        }
    }

    /* Add device to ledger */
    if (!ret) {
        list_append(get_dev_ldgr(pDev->obj.type), pDev);
        usys_log_debug(
                        "Ledger:: Device Name: %s, Disc: %s, Type: %d,"
                        " Module Id %s is added to DB.",
                        pDev->obj.name, pDev->obj.desc, pDev->obj.type,
                        pDev->obj.modUuid);
    } else {
        usys_log_debug(
                        "Err(%d): Dev:: Device Name: %s, Disc: %s, Type: %d,"
                        " Module Id %s is not added to DB.",
                        ret, pDev->obj.name, pDev->obj.desc, pDev->obj.type,
                        pDev->obj.modUuid);
    }

    return ret;

}


static Device *ldgr_create_device(char *pUuid, ModuleCfg *pModCfg) {
    Device *dev = usys_zmalloc(sizeof(Device));
    if (dev) {
        dev->obj.type = pModCfg->devType;
        usys_memset(dev->obj.modUuid, '\0', 24);
        usys_memcpy(dev->obj.modUuid, pUuid, strlen(pUuid));
        usys_memset(dev->obj.name, '\0', 24);
        usys_memcpy(dev->obj.name, pModCfg->devName, strlen(pModCfg->devName));
        usys_memset(dev->obj.desc, '\0', 24);
        usys_memcpy(dev->obj.desc, pModCfg->devDesc, strlen(pModCfg->devDesc));
        usys_memset(dev->sysFile, '\0', 64);
        usys_memcpy(dev->sysFile, pModCfg->sysFile, strlen(pModCfg->sysFile));
        dev->devOps = NULL;

        uint16_t devCfgSize = 0;
        SIZE_OF_DEVICE_CFG(devCfgSize, pModCfg->devType);
        dev->hwAttr = usys_malloc(devCfgSize);
        if (dev->hwAttr) {
            usys_memset(dev->hwAttr, '\0', devCfgSize);
            usys_memcpy(dev->hwAttr, &pModCfg->cfg, devCfgSize);
        }

        dev->devCb = NULL;
    }
    return dev;
}

/* Init function for device db*/
int ldgr_init(void *data) {
    int ret = 0;

    /*IRQDB Init*/
    irqdb_init();

    /* Property parser */
    ret = parser_property_init(data);

   usys_log_debug("Ledger:: Initializing device DB.");
    return ret;
}

void ldgr_exit() {
    int ret = 0;
    irqdb_exit();
    parser_property_exit();
    ldgr_destory();
   usys_log_debug("Ledger:: Cleaning process  completed for device DB.");
}
/* Registering devices to device db. */
//TODO: make sure ModuleCfg -> cfg is same to Device->deviceCfg
int ldgr_register(char *pUuid, char *name, uint8_t count, ModuleCfg *pModCfg) {
    int ret = 0;
    for (uint8_t iter = 0; iter < count; iter++) {

        Device *dev = ldgr_create_device(pUuid, &pModCfg[iter]);
        if (dev) {

            ret = ldgr_registration(dev);
            if (ret) {
               usys_log_debug(
                    "Err(%d): Ledger:: Failed to register Device Name: %s,"
                    " Disc: %s, Type: %d, Module Id %s is not added to DB.",
                    ret, dev->obj.name, dev->obj.desc, dev->obj.type,
                    dev->obj.modUuid);
            }

        }

        ldgr_usys_free(dev);
    }
    return ret;
}

int ldgr_read_reg_dev_count(DeviceType type, uint16_t *count) {
    int ret = 0;
    *count = list_size(get_dev_ldgr(type));
    return ret;
}

int ldgr_read_reg_dev(DeviceType type, Device *pDev) {
    int ret = 0;
    if (pDev) {
        list_copy(get_dev_ldgr(type), pDev);
    } else {
        ret = ERR_NODED_INVALID_POINTER;
    }
    return ret;
}

int ldgr_read_prop_count(DevObj *obj, uint16_t *count) {
    int ret = 0;

    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {

        if (dev->devOps) {
            const DevOps *devOps = dev->devOps;
            if (devOps->readPropCount) {
                ret = devOps->readPropCount(dev, count);
            } else {
                ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
            }

        } else {
            ret = ERR_NODED_DEV_DRVR_MISSING;
        }

        if (dev) {
            usys_free(dev);
            dev = NULL;
        }

    } else {
        ret = ERR_NODED_DEV_MISSING;
    }

    return ret;
}

int ldgr_read_prop(DevObj *obj, void *prop) {
    int ret = 0;

    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {

        if (dev->devOps) {
            const DevOps *devOps = dev->devOps;
            if (devOps->readProp) {
                ret = devOps->readProp(dev, prop);
            } else {
                ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_NODED_DEV_DRVR_MISSING;
        }

        if (dev) {
            usys_free(dev);
            dev = NULL;
        }

    } else {
        ret = ERR_NODED_DEV_MISSING;
    }

    return ret;
}

/* TODO: Check if this is really required */
int ldgr_configure(DevObj *obj, void *prop, void *data) {
    int ret = 0;

    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {

        if (dev->devOps) {
            const DevOps *devOps = dev->devOps;
            if (devOps->configure) {
                ret = devOps->configure(dev, prop, data);
            } else {
                ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_NODED_DEV_DRVR_MISSING;
        }

        if (dev) {
            usys_free(dev);
            dev = NULL;
        }

    } else {
        ret = ERR_NODED_DEV_MISSING;
    }

    return ret;
}

int ldgr_read(DevObj *obj, void *prop, void *data) {
    int ret = 0;

    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {

        if (dev->devOps) {
            const DevOps *devOps = dev->devOps;
            if (devOps->read) {
                ret = devOps->read(dev, prop, data);
            } else {
                ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_NODED_DEV_DRVR_MISSING;
        }

        if (dev) {
            usys_free(dev);
            dev = NULL;
        }

    } else {
        ret = ERR_NODED_DEV_MISSING;
    }

    return ret;
}

int ldgr_write(DevObj *obj, void *prop, void *data) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->devOps) {
            const DevOps *devOps = dev->devOps;
            if (devOps->write) {
                ret = devOps->write(dev, prop, data);
            } else {
                ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_NODED_DEV_DRVR_MISSING;
        }
        if (dev) {
            usys_free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_NODED_DEV_MISSING;
    }
    return ret;
}

int ldgr_enable(DevObj *obj, void *prop, void *data) {
    int ret = 0;
    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->devOps) {
            const DevOps *devOps = dev->devOps;
            if (devOps->enable) {
                ret = devOps->enable(dev, prop, data);
            } else {
                ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_NODED_DEV_DRVR_MISSING;
        }
        if (dev) {
            usys_free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_NODED_DEV_MISSING;
    }
    return ret;
}

int ldgr_disable(DevObj *obj, void *prop, void *data) {
    int ret = 0;

    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {

        if (dev->devOps) {
            const DevOps *devOps = dev->devOps;
            if (devOps->disable) {
                ret = devOps->disable(dev, prop, data);
            } else {
                ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_NODED_DEV_DRVR_MISSING;
        }

        if (dev) {
            usys_free(dev);
            dev = NULL;
        }

    } else {
        ret = ERR_NODED_DEV_MISSING;
    }

    return ret;
}

//TODO: remove prop
/* Register the callback provided by user app. */
int ldgr_reg_app_cb(DevObj *obj, void *prop, CallBackFxn fun) {
    int ret = 0;

    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {

        if (fun) {
            dev->devCb = fun;
            ret = update_device_db(dev);
        }

        if (dev) {
            usys_free(dev);
            dev = NULL;
        }

    } else {

        ret = ERR_NODED_DEV_MISSING;
        usys_log_debug(
                        "Err(%d):: TMP callback for device name %s, "
                        "device disc %s and module UIID %s failed.",
                        ret, obj->name, obj->desc, obj->modUuid);

    }

    return ret;
}
//TODO: remove prop
/* Unregister the callback provided by user app. */
int ldgr_dereg_app_cb(DevObj *obj, void *prop, CallBackFxn fun) {
    int ret = 0;

    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {
        if (fun) {
            dev->devCb = NULL;
            ret = update_device_db(dev);
        }

        if (dev) {
            usys_free(dev);
            dev = NULL;
        }

    } else {

        ret = ERR_NODED_DEV_MISSING;
       usys_log_debug(
            "Err(%d):: Registering callback for device name %s, device disc %s "
            "and module UIID %s failed.",
            ret, obj->name, obj->desc, obj->modUuid);

    }

    return ret;
}

int ldgr_enable_irq(DevObj *obj, void *prop, void *data) {
    int ret = 0;

    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {

        if (dev->devOps) {
            const DevOps *devOps = dev->devOps;
            if (devOps->enableIrq) {
                ret = devOps->enableIrq(dev, prop, data);
            } else {
                ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_NODED_DEV_DRVR_MISSING;
        }

        if (dev) {
            usys_free(dev);
            dev = NULL;
        }

    } else {
        ret = ERR_NODED_DEV_MISSING;
    }

    return ret;
}

int ldgr_disable_irq(DevObj *obj, void *prop, void *data) {
    int ret = 0;

    /* Search device in device list */
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->devOps) {

            const DevOps *devOps = dev->devOps;
            if (devOps->disableIrq) {
                ret = devOps->disableIrq(dev, prop, data);
            } else {
                ret = ERR_NODED_DEV_API_NOT_SUPPORTED;
            }

        } else {
            ret = ERR_NODED_DEV_DRVR_MISSING;
        }

        if (dev) {
            usys_free(dev);
            dev = NULL;
        }

    } else {
        ret = ERR_NODED_DEV_MISSING;
    }

    return ret;
}




