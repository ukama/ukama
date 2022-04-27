/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INCLUDERR_UBSP_H_
#define INCLUDERR_UBSP_H_

#include "headers/ubsp/ukdblayout.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

UnitCfg* ukdb_alloc_unit_cfg(uint8_t module_count);
void ukdb_free_unit_cfg(UnitCfg *cfg, uint8_t count);
ModuleCfg* ukdb_alloc_module_cfg(uint8_t dev_count);
void ukdb_free_module_cfg(ModuleCfg *cfg, uint8_t count);

/* Moduledb */
int get_fieldid_index(UKDBIdxTuple* index,uint16_t fid, uint8_t idx_count, uint8_t *idx_iter);
int ukdb_idb_init(void* data);
int ukdb_env_setup(char *name);
int ukdb_validating_magicword(char* p_uuid);
int ukdb_write_magicword(char* p_uuid);
int ukdb_write_header(char* p_uuid, UKDBHeader* header);
int ukdb_read_header(char* p_uuid, UKDBHeader* header);
int ukdb_read_dbversion(char* p_uuid, Version* ver);
int ukdb_register_modules(char* p_uuid);
int ukdb_register_devices(char* p_uuid);
int ukdb_update_dbversion(char* p_uuid, Version ver);
int ukdb_update_current_idx_count(char* p_uuid, uint16_t* idx_count);
int ukdb_read_current_idx_count(char* p_uuid, uint16_t* idx_count);
int ukdb_erase_db(char* p_uuid);
int ukdb_erase_index(char* p_uuid);
int ukdb_update_index_atid(char* p_uuid, UKDBIdxTuple* p_data, uint16_t idx);
int ukdb_search_fieldid(char* p_uuid, UKDBIdxTuple** idx_data, uint16_t *idx, uint16_t fieldid);
int ukdb_update_index_for_fieldid(char* p_uuid, void* p_data, uint16_t fieldid);
int ukdb_write_index(char* p_uuid, UKDBIdxTuple* p_data);
int ukdb_read_index(char* p_uuid, UKDBIdxTuple* p_data, uint16_t idx);
int ukdb_erase_payload(char* p_uuid, uint16_t fieldid);
int ukdb_write_payload(char* p_uuid, void* p_data, uint16_t offset, uint16_t size);
int ukdb_update_payload(char* p_uuid, void* p_data, uint16_t fieldid, uint16_t size, uint8_t state, Version version);
int ukdb_write_unit_info_data(char* p1_uuid, UKDBIdxTuple* index, char* p_uuid, uint8_t* count);
int ukdb_write_unit_cfg_data(char* p1_uuid, UKDBIdxTuple* index, uint8_t count);
int ukdb_write_module_cfg_data(char* p_uuid, ModuleInfo* minfo, UKDBIdxTuple* cfg_index);
int ukdb_write_module_info_data(char* p_uuid, UKDBIdxTuple* info_index, UKDBIdxTuple* cfg_index, uint8_t *mcount, uint8_t *idx);
int ukdb_write_generic_data(UKDBIdxTuple* index, char* p_uuid, uint16_t fid);
int ukdb_create_db(char *p_uuid);
int ukdb_remove_db(char *p_uuid);
int ukdb_init(char* systemdb);
int ukdb_register_module(UnitCfg* cfg);
int ukdb_unregister_module(char* puuid);
int ukdb_pre_create_db_setup(char* mod_uuid);
void ukdb_idb_exit();
void ukdb_exit();
/*UKDB read operations */

/* Direct read*/
int ukdb_validate_payload( void* p_data, uint32_t crc, uint16_t size);
int ukdb_read_payload(char *p_uuid,  void *p_data, uint16_t offset, uint16_t size);
int ukdb_read_payload_for_fieldid(char *p_uuid, void* p_payload, uint16_t fid, uint16_t* size);
int ukdb_read_payload_from_ukdb(char *p_uuid, void* p_data, uint16_t id, uint16_t* size);
/*Higher level read */
int ukdb_read_module(char *p_uuid, ModuleInfo* p_minfo);
int ukdb_read_unit(char *p_uuid, UnitInfo* p_uinfo, UnitCfg* p_ucfg);

int ukdb_read_unit_info(char *p_uuid, UnitInfo* data, uint16_t* size);
int ukdb_read_unit_cfg(char *p_uuid, UnitCfg* p_ucfg, uint8_t count, uint16_t* size);
int ukdb_read_module_info(char *p_uuid, ModuleInfo* p_info, uint16_t* size);
int ukdb_read_module_cfg(char *p_uuid, ModuleCfg* p_cfg, uint8_t count, uint16_t* size);
int ukdb_read_fact_config(char *p_uuid, void* data, uint16_t* size);
int ukdb_read_user_config(char *p_uuid, void* data, uint16_t* size);
int ukdb_read_fact_calib(char *p_uuid, void* data, uint16_t* size);
int ukdb_read_user_calib(char *p_uuid, void* data, uint16_t* size);
int ukdb_read_bs_certs(char *p_uuid, void* data, uint16_t* size);
int ukdb_read_lwm2m_certs(char *p_uuid, void* data, uint16_t* size);

/*logs*/
void ukdp_print_header(UKDBHeader* header);
void ukdb_print_index_table(UKDBIdxTuple* idx_tbl, uint8_t count);
void ukdb_print_unit_info(UnitInfo* p_uinfo);
void ukdb_print_module_info(ModuleInfo* p_minfo);
void ukdb_print_unit_cfg(UnitCfg* p_ucfg, uint8_t count);
void ukdb_print_module_cfg(ModuleCfg* p_mcfg, uint8_t count);

/*TODO*/
/* First search if module uuid is is present and then return the info.*/
int ukdb_read_module_info_by_uuid();
int ukdb_read_module_cfg_by_uuid();
int ukdb_read_device_cfg_by_uuid(char* mod_uuid, char* devname);

#endif /*INCLUDERR_UBSP_H_*/
