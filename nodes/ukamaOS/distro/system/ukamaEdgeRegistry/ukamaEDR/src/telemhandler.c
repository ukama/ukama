/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/tmp.h"
#include "registry/pwr.h"
#include "registry/gpio.h"
#include "registry/adc.h"
#include "registry/led.h"
#include "registry/atten.h"

#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/registry.h"
#include "inc/ukamadr.h"
#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "headers/objects/objects.h"
#include "headers/utils/file.h"
#include "headers/utils/list.h"
#include "headers/utils/log.h"
#include "utils/timer.h"
#include "dmt.h"

#include <stdint.h>
#include <string.h>
#include <time.h>

#define RTMPFILE "/history/tempdata.csv"
#define RPWRFILE "/history/pwrdata.csv"
#define RADCFILE "/history/adcdata.csv"
#define RATTFILE "/history/attdata.csv"
#define RGPIOFILE "/history/gpiodata.csv"

#define TELM_PERIODIC_TIMER 30000

static pthread_t gtelmthread_id = 0;
size_t gtelmtimer;

void telemhandler_init() {
    initialize(&gtelmthread_id);
}

int telm_read_tmp_inst_from_dev(DRDBSchema *inst) {
    int ret = 0;
    char rowdesc[256] = { '\0' };
    char rowdata[256] = { '\0' };
    int64_t seconds = time(NULL);
    TempData *data = (TempData *)inst->data;
    ret = db_read_inst_data_from_dev(&inst->obj, &data->value);
    if (!ret) {
        double temp = *(double *)reg_data_value(&data->value);
        drdb_update_tmp_inst_data(temp, &data->min, &data->max, &data->avg,
                                  &data->cumm, &data->counter);
        /* Row description*/
        sprintf(rowdesc, "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,\n", "Time",
                "Instance", "Name", "Module UUID", "Type", "Description",
                "Temp Value", "Min Temp", "Max Temp", " Temp Average",
                "Counter", "Unit");

        /* Temp Sensor Data */
        sprintf(rowdata, "%ld,%d,%s,%s,%d,%s,%f,%f,%f,%f,%d,%s,\n", seconds,
                inst->instance, inst->obj.name, inst->obj.mod_UUID,
                inst->obj.type, inst->obj.disc, data->value.value.doubleval,
                data->min.value.doubleval, data->max.value.doubleval,
                data->avg.value.doubleval, data->counter.value.intval,
                data->units.value.stringval);

        ret = file_add_record(RTMPFILE, rowdesc, rowdata);
    }
    return ret;
}

int telm_read_pwr_inst_from_dev(DRDBSchema *inst) {
    int ret = 0;
    char rowdesc[256] = { '\0' };
    char rowdata[256] = { '\0' };
    int64_t seconds = time(NULL);
    GenPwrData *data = (GenPwrData *)inst->data;
    ret = db_read_inst_data_from_dev(&inst->obj, &data->value);
    if (!ret) {
        double val = *(double *)reg_data_value(&data->value);
        drdb_update_pwr_inst_data(val, &data->min, &data->max, &data->avg,
                                  &data->cumm, &data->counter);
        /* Row description*/
        sprintf(rowdesc, "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,\n", "Time",
                "Instance", "Name", "Module UUID", "Type", "Description",
                "Sensor Value", "Min Value", "Max Value", "Average", "Counter",
                "Unit");

        /* Voltage/Power/Current Data */
        sprintf(rowdata, "%ld,%d,%s,%s,%d,%s,%f,%f,%f,%f,%d,%s,\n", seconds,
                inst->instance, inst->obj.name, inst->obj.mod_UUID,
                inst->obj.type, inst->obj.disc, data->value.value.doubleval,
                data->min.value.doubleval, data->max.value.doubleval,
                data->avg.value.doubleval, data->counter.value.intval,
                data->units.value.stringval);

        ret = file_add_record(RPWRFILE, rowdesc, rowdata);
    }
    return ret;
}

int telm_read_gpio_inst_from_dev(DRDBSchema *inst) {
    int ret = 0;
    char rowdesc[256] = { '\0' };
    char rowdata[256] = { '\0' };
    int64_t seconds = time(NULL);
    DigitalData *data = (DigitalData *)inst->data;
    ret = db_read_inst_data_from_dev(&inst->obj, &data->state);
    if (!ret) {
        /* Row description*/
        sprintf(rowdesc, "%s,%s,%s,%s,%s,%s,%s,%s,%s,\n", "Time", "Instance",
                "Name", "Module UUID", "Type", "Description", "Direction",
                "State", "Counter");

        /* GPIO Data */
        sprintf(rowdata, "%ld,%d,%s,%s,%d,%s,%d,%d,%d,\n", seconds,
                inst->instance, inst->obj.name, inst->obj.mod_UUID,
                inst->obj.type, inst->obj.disc, data->direction.value.intval,
                data->state.value.intval, data->counter.value.intval);

        ret = file_add_record(RGPIOFILE, rowdesc, rowdata);
    }
    return ret;
}

int telm_read_adc_inst_from_dev(DRDBSchema *inst) {
    int ret = 0;
    char rowdesc[256] = { '\0' };
    char rowdata[256] = { '\0' };
    int64_t seconds = time(NULL);
    AdcData *data = (AdcData *)inst->data;
    ret = db_read_inst_data_from_dev(&inst->obj, &data->outputcurr);
    if (!ret) {
        /* Row description*/
        sprintf(rowdesc, "%s,%s,%s,%s,%s,%s,%s,\n", "Time", "Instance", "Name",
                "Module UUID", "Type", "Description", "CurrValue");

        /* ADC Data */
        sprintf(rowdata, "%ld,%d,%s,%s,%d,%s,%f,\n", seconds, inst->instance,
                inst->obj.name, inst->obj.mod_UUID, inst->obj.type,
                inst->obj.disc, data->outputcurr.value.doubleval);

        ret = file_add_record(RADCFILE, rowdesc, rowdata);
    }
    return ret;
}

int telm_read_att_inst_from_dev(DRDBSchema *inst) {
    int ret = 0;
    char rowdesc[256] = { '\0' };
    char rowdata[256] = { '\0' };
    int64_t seconds = time(NULL);
    AttData *data = (AttData *)inst->data;
    ret = db_read_inst_data_from_dev(&inst->obj, &data->attvalue);
    ret |= db_read_inst_data_from_dev(&inst->obj, &data->latchenable);
    if (!ret) {
        /* Row description*/
        sprintf(rowdesc, "%s,%s,%s,%s,%s,%s,%s,%s,\n", "Time", "Instance",
                "Name", "Module UUID", "Type", "Description", "Attenuation",
                "Latch");

        /* ADC Data */
        sprintf(rowdata, "%ld,%d,%s,%s,%d,%s,%d,%d,\n", seconds, inst->instance,
                inst->obj.name, inst->obj.mod_UUID, inst->obj.type,
                inst->obj.disc, data->attvalue.value.intval,
                data->latchenable.value.intval);

        ret = file_add_record(RATTFILE, rowdesc, rowdata);
    }
    return ret;
}

int telm_read_inst_from_dev(DRDBSchema *instr) {
    int ret = 0;
    switch (instr->type) {
    case OBJ_TYPE_TMP:
        ret = telm_read_tmp_inst_from_dev(instr);
        break;
    case OBJ_TYPE_VOLT:
    case OBJ_TYPE_CURR:
    case OBJ_TYPE_PWR:
        ret = telm_read_pwr_inst_from_dev(instr);
        break;
    case OBJ_TYPE_DIP:
    case OBJ_TYPE_DOP:
        ret = telm_read_gpio_inst_from_dev(instr);
        break;
    case OBJ_TYPE_LED:
        //ret = telm_read_led_inst_from_dev(instr);
        break;
    case OBJ_TYPE_ADC:
        ret = telm_read_adc_inst_from_dev(instr);
        break;
    case OBJ_TYPE_ATT:
        ret = telm_read_att_inst_from_dev(instr);
        break;
    default: {
        ret = -1;
    }
    }
    return ret;
}

void telemhandler_service(size_t timer_id, void *data) {
    int ret = 0;
    /* Read data from every device type */
    log_debug("TELM:: Recording the Telemetry for Unit.\n");
    for (ObjectType type = OBJ_TYPE_UNIT; type < OBJ_TYPE_ALARM; type++) {
        ListInfo *db = reg_getdb(type);
        ListNode *node = NULL;
        int instance = 0;
        /* Read data from every instance */
        while (TRUE) {
            list_next(db, &node);
            if (node) {
                DRDBSchema *instr = node->data;
                ret = telm_read_inst_from_dev(instr);
                if (ret) {
                    log_debug(
                        "TELM:: Failed to record data for inst %d, Type 0x%x, name %s, Module ID %s\n",
                        instr->instance, type, instr->obj.name,
                        instr->obj.mod_UUID);
                }
            } else {
                break;
            }
        }
    }
}

void telemhandler_start() {
    gtelmtimer = start_timer(TELM_PERIODIC_TIMER, &telemhandler_service,
                             TIMER_PERIODIC, NULL);
    if (gtelmtimer) {
        log_debug("TELM:: Periodic timer started with period of %d millisec.\n",
                  (TELM_PERIODIC_TIMER));
    }
}

void telemhandler_stop() {
    stop_timer(gtelmtimer);
}

void telemhandler_exit() {
    /* TODO: Check if we need time stop. */
    log_debug("TELM: Stopping telemetry handler.");
    telemhandler_stop(gtelmtimer);
    finalize(gtelmthread_id);
}
