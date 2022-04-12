/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inc/devicedb.h"

#include "headers/errorcode.h"
#include "utils/irqdb.h"
#include "headers/utils/log.h"
#include "utils/pparser.h"
#include "devdb/adc/adc.h"
#include "devdb/att/att.h"
#include "devdb/gpio/gpiowrapper.h"
#include "devdb/led/led.h"
#include "devdb/pwr/pwr.h"
#include "devdb/tmp/tmp.h"
#include "ukdb/db/eeprom.h"

static const DevFxnTable *get_fxn_tbl(DevObj *dobj) {
    const DevFxnTable *devdb_fxn_tbl = NULL;
    switch (dobj->type) {
    case DEV_TYPE_TMP: {
        devdb_fxn_tbl = get_dev_tmp_fxn_tbl(dobj->name);
        break;
    }
    case DEV_TYPE_PWR: {
        devdb_fxn_tbl = get_dev_pwr_fxn_tbl(dobj->name);
        break;
    }
    case DEV_TYPE_GPIO: {
        devdb_fxn_tbl = get_dev_gpiow_fxn_tbl(dobj->name);
        break;
    }
    case DEV_TYPE_LED: {
        devdb_fxn_tbl = get_dev_led_fxn_tbl(dobj->name);
        break;
    }
    case DEV_TYPE_ADC: {
        devdb_fxn_tbl = get_dev_adc_fxn_tbl(dobj->name);
        break;
    }
    case DEV_TYPE_ATT: {
        devdb_fxn_tbl = get_dev_att_fxn_tbl(dobj->name);
        break;
    }
    default: {
    }
    }
    return devdb_fxn_tbl;
}

static ListInfo *get_dev_db(DeviceType type) {
    ListInfo *devdb = NULL;
    switch (type) {
    case DEV_TYPE_TMP: {
        devdb = get_dev_tmp_db();
        break;
    }
    case DEV_TYPE_PWR: {
        devdb = get_dev_pwr_db();
        break;
    }
    case DEV_TYPE_GPIO: {
        devdb = get_dev_gpiow_db();
        break;
    }
    case DEV_TYPE_LED: {
        devdb = get_dev_led_db();
        break;
    }
    case DEV_TYPE_ADC: {
        devdb = get_dev_adc_db();
        break;
    }
    case DEV_TYPE_ATT: {
        devdb = get_dev_att_db();
        break;
    }
    default: {
    }
    }
    return devdb;
}

static void devdb_free(Device *dev) {
    if (dev) {
        if (dev->hw_attr) {
            free(dev->hw_attr);
        }
        free(dev);
    }
}
int compare_dev_node(void *ipt, void *sd) {
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

/* Searching device in the device list*/
static Device *search_device_object(DevObj *dev_obj) {
    Device *fdev = NULL;
    Device *sdev = malloc(sizeof(Device));
    if (sdev) {
        memcpy(&sdev->obj, dev_obj, sizeof(DevObj));
        /*Search return 1 for found.*/
        fdev = list_search(get_dev_db(dev_obj->type), sdev);
        if (fdev) {
            log_trace("DEVDB:: Device Name %s, Disc: %s Module UUID: %s found.",
                      dev_obj->name, dev_obj->disc, dev_obj->mod_UUID);
        } else {
            if (fdev) {
                free(fdev);
                fdev = NULL;
            }
            log_debug(
                "DEVDB:: Device Name %s, Disc: %s Module UUID: %s not found.",
                dev_obj->name, dev_obj->disc, dev_obj->mod_UUID);
        }

        if (sdev) {
            free(sdev);
            sdev = NULL;
        }
    }
    return fdev;
}

static void devdb_destory() {
    DeviceType type = DEV_TYPE_TMP;
    for (; type < DEV_TYPE_MAX; type++) {
        list_destroy(get_dev_db(type));
        log_warn("DEVDB:: Removing DB for 0x%d device type.", type);
    }
}

/* Searching device in the device list*/
static int update_device_db(Device *dev) {
    int ret = 0;
    if (dev) {
        /*Return 0 for update.*/
        ret = list_update(get_dev_db(dev->obj.type), dev);
        if (ret) {
            log_error(
                "DEVDB:: Device Name %s, Disc: %s Module UUID: %s update failed.",
                dev->obj.name, dev->obj.disc, dev->obj.mod_UUID);
        }
    }
    return ret;
}

void devdb_irq_callback(void *pcfg) {
    IRQCfg *cfg = pcfg;
    if (!cfg) {
        return;
    }
    DevObj *obj = &cfg->obj;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        /* if callback id registered */
        if (dev->dev_cb) {
            /* if fxn table is initialized.*/
            if (dev->fxn_tbl) {
                const DevFxnTable *fxn_tbl = get_fxn_tbl(&(dev->obj));
                /* if fxn is assigned.*/
                if (fxn_tbl->confirm_irq) {
                    AlertCallBackData *acbdata = NULL;
                    int acount = 0;
                    char fname[64] = { '\0' };
                    strcpy(fname, cfg->fname);
                    fxn_tbl->confirm_irq(dev, &acbdata, fname, &acount);
                    if (acount > 0) {
                        log_debug(
                            "DEVDB:: Calling callback fxn for device name %s, disc %s and module UIID %s with %d alerts.",
                            obj->name, obj->disc, obj->mod_UUID, acount);
                        dev->dev_cb(obj, &acbdata, &acount);
                    }
                }
            }
        } else {
            log_debug(
                "DEVDB:: No Alert callback registered for device name %s, disc %s and module UIID %s.",
                obj->name, obj->disc, obj->mod_UUID);
        }
    }
    if (dev) {
        free(dev);
        dev = NULL;
    }
}

static int devdb_registration(Device *p_dev) {
    int ret = 0;
    uint8_t idx = 0;
    /* Search in list if device already exist.*/
    ret = list_if_element_found(get_dev_db(p_dev->obj.type), p_dev);
    if (!ret) {
        if (!p_dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = get_fxn_tbl(&(p_dev->obj));
            if (fxn_tbl) {
                /* Updating fxn table for TMP */
                p_dev->fxn_tbl = fxn_tbl;
                if (fxn_tbl->init) {
                    ret = fxn_tbl->init(p_dev);
                    if (!ret) {
                        if (fxn_tbl->register_cb) {
                            ret = fxn_tbl->register_cb(p_dev,
                                                       &devdb_irq_callback);
                        }
                    }
                } else {
                    ret = ERR_UBSP_DEV_DRVR_MISSING;
                }
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            /* Should not come here*/
            const DevFxnTable *fxn_tbl = p_dev->fxn_tbl;
            if (fxn_tbl->init) {
                ret = fxn_tbl->init(p_dev);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        }

        if (!ret) {
            /* if doesn't exist add to list.*/
            list_append(get_dev_db(p_dev->obj.type), p_dev);
            log_debug(
                "DEVDB:: Device Name: %s, Disc: %s, Type: %d, Module Id %s is added to DB.",
                p_dev->obj.name, p_dev->obj.disc, p_dev->obj.type,
                p_dev->obj.mod_UUID);
        } else {
            log_debug(
                "Err(%d): Dev:: Device Name: %s, Disc: %s, Type: %d, Module Id %s is not added to DB.",
                ret, p_dev->obj.name, p_dev->obj.disc, p_dev->obj.type,
                p_dev->obj.mod_UUID);
        }

    } else {
        ret = 0;
        log_debug(
            "DEVDB:: Device Name: %s, Disc: %s, Type: %d, Module Id %s is already in DB.",
            p_dev->obj.name, p_dev->obj.disc, p_dev->obj.type,
            p_dev->obj.mod_UUID);
    }
    return ret;
}

static Device *devdb_create_device(char *p_uuid, ModuleCfg *p_mcfg) {
    Device *dev = malloc(sizeof(Device));
    if (dev) {
        dev->obj.type = p_mcfg->dev_type;
        memset(dev->obj.mod_UUID, '\0', 24);
        memcpy(dev->obj.mod_UUID, p_uuid, strlen(p_uuid));
        memset(dev->obj.name, '\0', 24);
        memcpy(dev->obj.name, p_mcfg->dev_name, strlen(p_mcfg->dev_name));
        memset(dev->obj.disc, '\0', 24);
        memcpy(dev->obj.disc, p_mcfg->dev_disc, strlen(p_mcfg->dev_disc));
        memset(dev->sysfile, '\0', 64);
        memcpy(dev->sysfile, p_mcfg->sysfile, strlen(p_mcfg->sysfile));
        dev->fxn_tbl = NULL;
        uint16_t dev_cfg_size = 0;
        SIZE_OF_DEVICE_CFG(dev_cfg_size, p_mcfg->dev_type);
        dev->hw_attr = malloc(dev_cfg_size);
        if (dev->hw_attr) {
            memset(dev->hw_attr, '\0', dev_cfg_size);
            memcpy(dev->hw_attr, &p_mcfg->cfg, dev_cfg_size);
        }
        dev->dev_cb = NULL;
    }
    return dev;
}

/* Init function for device db*/
int devdb_init(void *data) {
    int ret = 0;

    /*IRQDB Init*/
    irqdb_init();

    /* Property parser */
    ret = parser_property_init(data);

    log_debug("DEVDB:: Initializing device DB.");
    return ret;
}

void devdb_exit() {
    int ret = 0;
    irqdb_exit();
    parser_property_exit();
    devdb_destory();
    log_debug("DEVDB:: Cleaning process  completed for device DB.");
}
/* Registering devices to device db. */
//TODO: make sure ModuleCfg -> cfg is same to Device->deviceCfg
int devdb_register(char *p_uuid, char *name, uint8_t count, ModuleCfg *p_mcfg) {
    int ret = 0;
    for (uint8_t iter = 0; iter < count; iter++) {
        Device *dev = devdb_create_device(p_uuid, &p_mcfg[iter]);
        if (dev) {
            ret = devdb_registration(dev);
            if (ret) {
                log_debug(
                    "Err(%d): DEVDB:: Failed to register Device Name: %s, Disc: %s, Type: %d, Module Id %s is not added to DB.",
                    ret, dev->obj.name, dev->obj.disc, dev->obj.type,
                    dev->obj.mod_UUID);
            }
        }
        devdb_free(dev);
    }
    return ret;
}

int devdb_read_reg_dev_count(DeviceType type, uint16_t *count) {
    int ret = 0;
    *count = list_size(get_dev_db(type));
    return ret;
}

int devdb_read_reg_dev(DeviceType type, Device *p_dev) {
    int ret = 0;
    if (p_dev) {
        list_copy(get_dev_db(type), p_dev);
    } else {
        ret = ERR_UBSP_INVALID_POINTER;
    }
    return ret;
}

int devdb_read_prop_count(DevObj *obj, uint16_t *count) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = dev->fxn_tbl;
            if (fxn_tbl->read_prop) {
                ret = fxn_tbl->read_prop_count(dev, count);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_UBSP_DEV_DRVR_MISSING;
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
    }
    return ret;
}

int devdb_read_prop(DevObj *obj, void *prop) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = dev->fxn_tbl;
            if (fxn_tbl->read_prop) {
                ret = fxn_tbl->read_prop(dev, prop);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_UBSP_DEV_DRVR_MISSING;
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
    }
    return ret;
}

/* TODO: Check if this is really required */
int devdb_configure(DevObj *obj, void *prop, void *data) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = dev->fxn_tbl;
            if (fxn_tbl->configure) {
                ret = fxn_tbl->configure(dev, prop, data);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_UBSP_DEV_DRVR_MISSING;
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
    }
    return ret;
}

int devdb_read(DevObj *obj, void *prop, void *data) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = dev->fxn_tbl;
            if (fxn_tbl->read) {
                ret = fxn_tbl->read(dev, prop, data);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_UBSP_DEV_DRVR_MISSING;
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
    }
    return ret;
}

int devdb_write(DevObj *obj, void *prop, void *data) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = dev->fxn_tbl;
            if (fxn_tbl->write) {
                ret = fxn_tbl->write(dev, prop, data);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_UBSP_DEV_DRVR_MISSING;
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
    }
    return ret;
}

int devdb_enable(DevObj *obj, void *prop, void *data) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = dev->fxn_tbl;
            if (fxn_tbl->enable) {
                ret = fxn_tbl->enable(dev, prop, data);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_UBSP_DEV_DRVR_MISSING;
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
    }
    return ret;
}

int devdb_disable(DevObj *obj, void *prop, void *data) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = dev->fxn_tbl;
            if (fxn_tbl->disable) {
                ret = fxn_tbl->disable(dev, prop, data);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_UBSP_DEV_DRVR_MISSING;
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
    }
    return ret;
}

//TODO: remove prop
/* Register the callback provided by user app. */
int devdb_reg_app_cb(DevObj *obj, void *prop, CallBackFxn fun) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (fun) {
            dev->dev_cb = fun;
            ret = update_device_db(dev);
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
        log_debug(
            "Err(%d):: TMP callback for device name %s, device disc %s and module UIID %s failed.",
            ret, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}
//TODO: remove prop
/* Unregister the callback provided by user app. */
int devdb_dereg_app_cb(DevObj *obj, void *prop, CallBackFxn fun) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (fun) {
            dev->dev_cb = NULL;
            ret = update_device_db(dev);
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
        log_debug(
            "Err(%d):: TMP callback for device name %s, device disc %s and module UIID %s failed.",
            ret, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int devdb_enable_irq(DevObj *obj, void *prop, void *data) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = dev->fxn_tbl;
            if (fxn_tbl->enable_irq) {
                ret = fxn_tbl->enable_irq(dev, prop, data);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_UBSP_DEV_DRVR_MISSING;
        }
        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
    }
    return ret;
}

int devdb_disable_irq(DevObj *obj, void *prop, void *data) {
    int ret = 0;
    /*search device in device list*/
    Device *dev = search_device_object(obj);
    if (dev) {
        if (dev->fxn_tbl) {
            const DevFxnTable *fxn_tbl = dev->fxn_tbl;
            if (fxn_tbl->disable_irq) {
                ret = fxn_tbl->disable_irq(dev, prop, data);
            } else {
                ret = ERR_UBSP_DEV_API_NOT_SUPPORTED;
            }
        } else {
            ret = ERR_UBSP_DEV_DRVR_MISSING;
        }

        if (dev) {
            free(dev);
            dev = NULL;
        }
    } else {
        ret = ERR_UBSP_DEV_MISSING;
    }
    return ret;
}
