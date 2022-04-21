/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/bsp_gpio.h"

#include "devhelper.h"
#include "errorcode.h"
#include "property.h"
#include "drivers/gpio_wrapper.h"
#include "drivers/sysfs_wrapper.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

const DrvrOps drvrGpioWrapperOps = { .init = gpio_wrapper_init,
                                     .configure = gpio_wrapper_configure,
                                     .read = gpio_wrapper_read,
                                     .write = gpio_wrapper_write,
                                     .enable = gpio_wrapper_enable,
                                     .disable = gpio_wrapper_disable,
                                     .registerCb = NULL,
                                     .dregisterCb = NULL,
                                     .enableIrq = NULL,
                                     .disableIrq = NULL };
static Property *gProperty = NULL;
static int gPropertyCount = 0;

static Property gpioProperty[MAXGPIOPROP] = {
    [DIRECTION] = { .name = "GPIO DIRECTION",
                    .dataType = TYPE_UINT8,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_AVAIL,
                    .propType = PROP_TYPE_CONFIG,
                    .units = "NA",
                    .sysFname = "direction",
                    .depProp = NULL },
    [VALUE] = { .name = "GPIO VALUE",
                .dataType = TYPE_UINT8,
                .perm = PERM_RD | PERM_WR,
                .available = PROP_AVAIL,
                .propType = PROP_TYPE_CONFIG,
                .units = "NA",
                .sysFname = "value",
                .depProp = NULL },
    [EDGE] = { .name = "GPIO EDGE",
               .dataType = TYPE_UINT8,
               .perm = PERM_RD | PERM_WR,
               .available = PROP_AVAIL,
               .propType = PROP_TYPE_CONFIG,
               .units = "NA",
               .sysFname = "edge",
               .depProp = NULL },
    [POLARITY] = { .name = "GPIO POLARITY",
                   .dataType = TYPE_UINT8,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysFname = "active_low",
                   .depProp = NULL }
};

static const DrvrOps* get_fxn_tbl(Device *pDev) {
    if (IF_SYSFS_SUPPORT(pDev->sysFile)) {
        return sysfs_wrapper_get_ops();
    } else {
        return &drvrGpioWrapperOps;
    }
}

int bsp_gpio_read_prop_count(Device *pDev, uint16_t *count) {
    int ret = 0;
    *count = gPropertyCount;
    return 0;
}

int bsp_gpio_read_properties(Device *pDev, void *prop) {
    int ret = 0;
    if (prop) {
        usys_memset(prop, '\0', sizeof(Property) * gPropertyCount);
        usys_memcpy(prop, gProperty, sizeof(Property) * gPropertyCount);
    } else {
        ret = -1;
    }
    return ret;
}

void gpio_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //gpio_read(void* pDev, void* prop, void* data );
}

int bsp_gpio_init(Device *pDev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(pDev, &gProperty, &gPropertyCount);
    if (ret) {
        gPropertyCount = MAXGPIOPROP;
        gProperty = gpioProperty;
        log_debug("GPIO: Using static property table with %d property.",
                  gPropertyCount);
    }

    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_init_driver(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_gpio_registration(Device *pDev) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_registration(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_gpio_configure(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_configure(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_gpio_read(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_read(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_gpio_write(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_write(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_gpio_enable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_gpio_disable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}
