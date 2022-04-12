/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_DBHANDLER_H_
#define INC_DBHANDLER_H_

#include "inc/registry.h"
#include "headers/ubsp/ukdblayout.h"


int db_enable_alarm(DevObj *obj, int prop);
int db_read_inst_data_from_dev(DevObj *obj, PData *sd);
int db_register_alarm_callback(DevObj *obj, int *prop, CallBackFxn fn);
int db_read_boolval_prop(DevObj *obj, uint8_t req, bool *mdata, PData *pdata, int *size);
int db_read_doubleval_prop(DevObj* obj, uint8_t req, double* mdata, PData* pdata, int *size);
int db_read_intval_prop(DevObj* obj, uint8_t req, int* mdata, PData* pdata, int *size);
int db_read_longintval_prop(DevObj *obj, uint8_t req, int64_t *mdata, PData *pdata, int *size);
int db_read_shortintval_prop(DevObj *obj, uint8_t req, uint16_t *mdata, PData *pdata, int *size);
int db_read_strval_prop(DevObj* obj, uint8_t req, char* mdata, PData* pdata, int *size);
int db_set_prop_val(PData *pdata, int type, void *val);
int db_write_boolval_prop(DevObj *obj, uint8_t req, bool mdata, PData *pdata);
int db_write_doubleval_prop(DevObj *obj, uint8_t req, double mdata, PData *pdata);
int db_write_intval_prop(DevObj *obj, uint8_t req, int mdata, PData *pdata);
int db_write_longintval_prop(DevObj *obj, uint8_t req, int64_t mdata, PData *pdata);
int db_write_shortintval_prop(DevObj *obj, uint8_t req, uint16_t mdata, PData *pdata);
int db_write_strval_prop(DevObj *obj, uint8_t req, char *mdata, PData *pdata);
int db_write_inst_data_to_dev(DevObj *obj, PData *sd);
void db_versiontostr(Version ver, char* str);
void db_free_unit_cfg(UnitCfg** cfg, uint8_t count);
ModuleCfg* db_read_module_cfg (char *puuid, uint8_t count );
ModuleInfo* db_read_module_info (char *puuid );
Property *db_read_dev_property(DevObj* obj, int* pcount);
UnitCfg* db_read_unit_cfg (char *puuid, uint8_t count );
UnitInfo* db_read_unit_info (char *puuid );

#endif /* INC_DBHANDLER_H_ */
