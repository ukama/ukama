/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/gpio/gpio.h"

#include "headers/errorcode.h"
#include "headers/ubsp/property.h"
#include "inc/devhelper.h"
#include "headers/utils/log.h"
#include "devdb/gpio/drvrgpio.h"
#include "devdb/sysfs/drvrsysfs.h"

const DrvDBFxnTable drvr_gpio_fxn_table = { .init = drvr_gpio_init,
                                            .configure = drvr_gpio_configure,
                                            .read = drvr_gpio_read,
                                            .write = drvr_gpio_write,
                                            .enable = drvr_gpio_enable,
                                            .disable = drvr_gpio_disable,
                                            .register_cb = NULL,
                                            .dregister_cb = NULL,
                                            .enable_irq = NULL,
                                            .disable_irq = NULL };
static Property *g_property = NULL;
static int g_property_count = 0;

static Property gpio_property[MAXGPIOPROP] = {
    [DIRECTION] = { .name = "GPIO DIRECTION",
                    .data_type = TYPE_UINT8,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_AVAIL,
                    .prop_type = PROP_TYPE_CONFIG,
                    .units = "NA",
                    .sysfname = "direction",
                    .dep_prop = NULL },
    [VALUE] = { .name = "GPIO VALUE",
                .data_type = TYPE_UINT8,
                .perm = PERM_RD | PERM_WR,
                .available = PROP_AVAIL,
                .prop_type = PROP_TYPE_CONFIG,
                .units = "NA",
                .sysfname = "value",
                .dep_prop = NULL },
    [EDGE] = { .name = "GPIO EDGE",
               .data_type = TYPE_UINT8,
               .perm = PERM_RD | PERM_WR,
               .available = PROP_AVAIL,
               .prop_type = PROP_TYPE_CONFIG,
               .units = "NA",
               .sysfname = "edge",
               .dep_prop = NULL },
    [POLARITY] = { .name = "GPIO POLARITY",
                   .data_type = TYPE_UINT8,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "active_low",
                   .dep_prop = NULL }
};

static const DrvDBFxnTable *get_fxn_tbl(Device *p_dev) {
    if (IF_SYSFS_SUPPORT(p_dev->sysfile)) {
        return drvr_sysfs_get_fxn_tbl();
    } else {
        return &drvr_gpio_fxn_table;
    }
}

int gpio_read_prop_count(Device *p_dev, uint16_t *count) {
    int ret = 0;
    *count = g_property_count;
    return 0;
}

int gpio_read_properties(Device *p_dev, void *prop) {
    int ret = 0;
    if (prop) {
        memset(prop, '\0', sizeof(Property) * g_property_count);
        memcpy(prop, g_property, sizeof(Property) * g_property_count);
    } else {
        ret = -1;
    }
    return ret;
}

void gpio_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //gpio_read(void* p_dev, void* prop, void* data );
}

int gpio_init(Device *p_dev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(p_dev, &g_property,
                                            &g_property_count);
    if (ret) {
        g_property_count = MAXGPIOPROP;
        g_property = gpio_property;
        log_debug("GPIO: Using static property table with %d property.",
                  g_property_count);
    }

    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_init_driver(drvr, p_dev);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int gpio_registration(Device *p_dev) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_registration(drvr, p_dev);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int gpio_configure(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_configure(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int gpio_read(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_read(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int gpio_write(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_write(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int gpio_enable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_enable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int gpio_disable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}
