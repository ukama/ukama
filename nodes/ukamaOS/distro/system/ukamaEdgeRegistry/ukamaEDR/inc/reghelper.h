/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_REGHELPER_H_
#define INC_REGHELPER_H_

#include "ifmsg.h"

#include "headers/globalheader.h"
#include "headers/ubsp/devices.h"
#include "headers/ubsp/property.h"
#include "registry/alarm.h"
#include "headers/utils/list.h"

int reg_register_devices();
int reg_register_misc();
int reg_add_dev(Device* dev);
int reg_exec_dev(MsgFrame *req);
int reg_read_dev(MsgFrame *req);
int reg_read_inst_count(MsgFrame *req);
int reg_read_max_instance(ObjectType type);
int reg_write_dev(MsgFrame *req);
int reg_update_dev(DRDBSchema *reg);
int reg_check_if_alarm_property_exist(Property* prop, int pcount);
int reg_enable_alarms(DevObj* obj, const AlarmSensorData* sdata);
int reg_read_inst_data_from_dev(DevObj* obj, PData* sd);
int reg_write_inst_data_to_dev(DevObj* obj, PData* sd);
void reg_exit();
void reg_init();
void reg_mod_init();
void reg_mod_exit();
void reg_unit_init();
void reg_unit_exit();
void reg_list_reg_devices();
void reg_exit_type(ObjectType type);
void reg_free_alarm_prop(AlarmPropertyData **pdata);
void reg_init_type(ObjectType type);
void reg_data_add_property(int pidx, Property *prop, PData *pdata);
void reg_data_copy_property(Property** destp, Property *srcp);
void reg_append_inst(ListInfo *regdb, void *node);
void reg_prepend_inst(ListInfo *regdb, void *node);
void reg_update_inst(ListInfo *regdb, void *node);
void reg_initialize_dev_idt(void* idt, uint16_t inst, uint16_t oid, uint16_t rid);
void reg_register_sensor_alarms(DevObj* obj, Property* prop, const AlarmSensorData* sdata,
		uint16_t inst, uint16_t objid, uint16_t rsrcid);
void free_reg(DRDBSchema **reg);
void free_reg_data(DRDBSchema *reg);
void free_sdata(Property **prop);

void* reg_data_value(PData* sd);
void* reg_initialize_alarm_prop(Property* prop, const AlarmSensorData* sdata);
DRDBSchema* reg_read_instance(int instance, DeviceType type);
DRDBSchema *reg_search_inst(int instance, uint16_t misc, ObjectType type);
ListInfo* reg_getdb(ObjectType type);
Property* reg_copy_pdata_prop(Property *sdata);

#endif /* INC_REGHELPER_H_ */
