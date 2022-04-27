/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "utils/irqhelper.h"

#include "utils/irqdb.h"
#include "headers/utils/log.h"
#include "devdb/sysfs/drvrsysfs.h"
#include "devdb/tmp/adt7481.h"
#include "devdb/pwr/ina226.h"

/* Return 1 for active, 0 for inactive and -ve for error.*/
static int irqhelper_validate_irq(const DrvDBFxnTable *drvr_db_fx_tbl,
                                  Device *dev, Property *prop, int iter,
                                  double *curr_val) {
    int ret = 0;
    int valid = -1;
    void *hwattr;
    char sysf[64] = { '\0' };
    int cuprop = prop[iter].dep_prop->curr_idx;
    memset(sysf, '\0', 64);
    memcpy(sysf, dev->sysfile, strlen(dev->sysfile));
    strcat(sysf, prop[cuprop].sysfname);
    hwattr = sysf;
    ret = drvr_db_fx_tbl->read(hwattr, &prop[cuprop], curr_val);
    if (!ret) {
        log_debug("ALERTHELPER:: Alert for property %d for "
                  "Device name: %s disc: %s module UUID %s "
                  " received with current value %lf.",
                  iter, dev->obj.name, dev->obj.disc, dev->obj.mod_UUID,
                  *curr_val);
    } else {
        log_debug("Err(%d): ALERTHELPER:: Alert for property %d for "
                  "Device name: %s disc: %s module UUID %s "
                  " failed to read current value.",
                  ret, iter, dev->obj.name, dev->obj.disc, dev->obj.mod_UUID);
        goto retvalidate;
    }

    /* TODO: Compare the value to limit set.*/
    /*if true increment the alert count otherwise continue.*/
    if (prop[iter].dep_prop->lmt_idx >= 0) {
        int lmtprop = prop[iter].dep_prop->lmt_idx;
        double lmtvalue = 0;
        memset(sysf, '\0', 64);
        memcpy(sysf, dev->sysfile, strlen(dev->sysfile));
        strcat(sysf, prop[lmtprop].sysfname);
        hwattr = sysf;
        ret = drvr_db_fx_tbl->read(hwattr, &prop[lmtprop], &lmtvalue);
        if (!ret) {
            log_debug("ALERTHELPER:: property[%d] %s alert limit value is %lf "
                      "Device name: %s disc: %s module UUID %s.",
                      iter, prop[iter].name, lmtvalue, dev->obj.name,
                      dev->obj.disc, dev->obj.mod_UUID);
        } else {
            log_debug(
                "Err(%d): ALERTHELPER:: property[%d] %s alert limit value read"
                " failed for Device name: %s disc: %s module UUID %s.",
                ret, iter, prop[iter].name, dev->obj.name, dev->obj.disc,
                dev->obj.mod_UUID);
            goto retvalidate;
        }

        valid =
            validate_irq_limits(*curr_val, lmtvalue, prop[iter].dep_prop->cond);
        if (valid == 1) {
            log_debug("ALERTHELPER:: Alert reported property[%d] %s for"
                      " Device name: %s disc: %s module UUID %s "
                      " received with current value %lf limit value %d.",
                      iter, prop[iter].name, dev->obj.name, dev->obj.disc,
                      dev->obj.mod_UUID, *curr_val, lmtvalue);
        } else if (valid == 0) {
            log_debug(
                "ALERTHELPER:: Alert condition not meet for property[%d] %s"
                " for Device name: %s disc: %s module UUID %s"
                " received with current value %lf limit value %d.",
                iter, prop[iter].name, dev->obj.name, dev->obj.disc,
                dev->obj.mod_UUID, *curr_val, lmtvalue);
        } else {
            log_debug("ALERTHELPER:: No alert condition found.");
        }
    }

retvalidate:
    if (valid == 1) {
        ret = 1;
    } else if (valid == 0) {
        ret = 0;
    }

    return ret;
}

static AlertCallBackData *irqhelper_prepare_alertcb_data(uint8_t alertstate,
                                                         int pidx, void *data,
                                                         int szdata) {
    /* Prepare the alert callback data */
    int sz = sizeof(AlertCallBackData);
    AlertCallBackData *cbdata = malloc(sz);
    if (cbdata) {
        memset(cbdata, '\0', sz);
        cbdata->alertstate = alertstate;
        cbdata->pidx = pidx;
        cbdata->svalue = malloc(szdata);
        if (cbdata) {
            memset(cbdata->svalue, '\0', szdata);
            memcpy(cbdata->svalue, data, szdata);
        } else {
            free(cbdata);
            cbdata = NULL;
        }
    } else {
        cbdata = NULL;
    }
    return cbdata;
}

/* Reading and confirming interrupts for TMP464 device */
int irqhelper_confirm_irq(const DrvDBFxnTable *drvr_db_fx_tbl, Device *p_dev,
                          AlertCallBackData **acbdata, Property *prop,
                          int max_prop, char *fpath, int *evt) {
    int ret = 0;
    int alertevt = 0;
    int a_status = 0;
    uint8_t alertstate = 0;
    double value = 0;
    int prop_idx = -1;
    char *fname = get_sysfs_name(fpath);
    //fname = basename(fpath);

    /* Scan through properties and check all alert related one.*/
    for (int iter = 0; iter < max_prop; iter++) {
        if ((prop[iter].prop_type & PROP_TYPE_ALERT) &&
            !(strcmp(prop[iter].sysfname, fname))) {
            if (drvr_db_fx_tbl) {
                alertevt = 1;
                void *hwattr;
                char sysf[64] = { '\0' };
                Device *dev = (Device *)p_dev;
                if (IF_SYSFS_SUPPORT(dev->sysfile)) {
                    memcpy(sysf, dev->sysfile, strlen(dev->sysfile));
                    strcat(sysf, prop[iter].sysfname);
                    hwattr = sysf;
                } else {
                    hwattr = dev->hw_attr;
                }

                /* Verifying if it's true alert by reading sysfs */
                int valid = 0;
                ret = drvr_db_fx_tbl->read(hwattr, &prop[iter], &a_status);

                /* Just to check further compare value with thresholds of a sensor */
                int lmtcheck = 0;
                if (prop[iter].dep_prop) {
                    if (prop[iter].dep_prop[0].curr_idx >= 0) {
                        lmtcheck = irqhelper_validate_irq(drvr_db_fx_tbl, dev,
                                                          prop, iter, &value);
                        prop_idx = iter;
                    }
                }

                /* Check which alert need to be raised */
                if (lmtcheck) {
                    const DevFxnTable *devdb_fxn_tbl = dev->fxn_tbl;
                    ret = devdb_fxn_tbl->irq_type(prop_idx, &alertstate);
                } else if (!a_status && !lmtcheck) {
                    /* Clear alert */
                    alertstate = ALARM_STATE_NO_ALARM_ACTIVE;
                } else {
                    /* looks like false alert. */
                    /* TODO: if astatus is 1 but value is with in limits should we send clear or just report as false interrupt.
					 * check if we can add active_alert to property json. */
                    alertevt = 0;
                }
            }
            break; /* Already read the file which raised interrupt.*/
        } /* If property was alert*/
    }

    /* In case of true alert event prepared the callback data */
    if (alertevt) {
        *acbdata = irqhelper_prepare_alertcb_data(alertstate, prop_idx, &value,
                                                  sizeof(double));
    } else {
        log_trace("ALERTHELPER:: ** False Alert reported **.");
    }
    *evt = alertevt;
    return ret;
}
