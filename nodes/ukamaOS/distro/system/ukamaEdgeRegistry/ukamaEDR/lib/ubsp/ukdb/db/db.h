/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DB_H_
#define DB_H_

#include "inc/globalheader.h"
#include "headers/ubsp/ukdblayout.h"

#include <stdio.h>
#include <stdint.h>

typedef int (*DbInitFxn)(void* data);
typedef int (*DbReadBlockFxn)(void* fname, void* data, off_t offset , uint16_t size);
typedef int (*DbWriteBlockFxn)(void* fname, void* data, off_t offset , uint16_t size);
typedef int (*DbEraseBlockFxn)(void* fname, off_t offset, uint16_t size);
typedef int (*DbReadNumberFxn)(void* fname, void* val, off_t offset, uint16_t count, uint8_t size);
typedef int (*DbWriteNumberFxn)(void* fname, void* val, off_t offset, uint16_t count, uint8_t size);
typedef int (*DbWriteProtect)(void* fname);
typedef int (*DBRenameFxn)(char* old_name, char* new_name);
typedef int (*DBCleanupFxn)(void* data);
typedef int (*DBRemoveFxn)(void* data);

/*basic read write operation to UKDB*/
typedef struct  {
	DbInitFxn init;
	DbReadBlockFxn read_block;
	DbWriteBlockFxn write_block;
	DbEraseBlockFxn erase_block;
	DbReadNumberFxn read_number;
	DbWriteNumberFxn write_number;
	DbWriteProtect write_protect;
	DBRenameFxn rename;
	DBCleanupFxn cleanup;
	DBRemoveFxn remove;
} DBFxnTable;

typedef union  __attribute__((__packed__)) {
    char sysfs[MAX_PATH_LENGTH];
    void* eeprom_cfg;
} DBAttr;

typedef struct {
     char mod_uuid[MAX_NAME_LENGTH];
     char mod_name[MAX_NAME_LENGTH];
     void* db_attr;
     const DBFxnTable* fxn_tbl;
 } ModuleDBMap;

int db_init();
ModuleDBMap* db_choose_module(char* p_uuid);
int db_register_module(UnitCfg* pcfg);
void db_unregister_all_module ();
int db_unregister_module(char *puuid);
int db_read_block(char* p_uuid, void* data, off_t offset , uint16_t size);
int db_rename(char* p_uuid, char* old_name, char* new_name);
int db_write_block(char* p_uuid, void* data, off_t offset , uint16_t size);
int db_erase_block(char* p_uuid, off_t offset, uint16_t size);
int db_read_number(char* p_uuid, void* val, off_t offset, uint16_t count, uint8_t size);
int db_write_number(char* p_uuid, void* val, off_t offset, uint16_t count, uint8_t size);
int db_write_protect (char* p_uuid, void* data);
int db_remove_database(char *p_uuid);
int db_register_update_module(char* p_uuid, UnitCfg *p_cfg, uint8_t count);
#endif /* DB_H_ */
