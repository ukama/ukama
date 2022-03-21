/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_STORE_H_
#define INC_STORE_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "noded_macros.h"
#include "schema.h"

typedef int (*StoreInitFxn)(void* data);
typedef int (*StoreReadBlockFxn)(void* fname, void* data, off_t offset , uint16_t size);
typedef int (*StoreWriteBlockFxn)(void* fname, void* data, off_t offset , uint16_t size);
typedef int (*StoreEraseBlockFxn)(void* fname, off_t offset, uint16_t size);
typedef int (*StoreReadNumberFxn)(void* fname, void* val, off_t offset, uint16_t count, uint8_t size);
typedef int (*StoreWriteNumberFxn)(void* fname, void* val, off_t offset, uint16_t count, uint8_t size);
typedef int (*StoreWriteProtect)(void* fname);
typedef int (*StoreRenameFxn)(char* old_name, char* new_name);
typedef int (*StoreCleanupFxn)(void* data);
typedef int (*StoreRemoveFxn)(void* data);

/* Basic read write operation to store */
typedef struct  {
  StoreInitFxn init;
  StoreReadBlockFxn readBlock;
  StoreWriteBlockFxn writeBlock;
  StoreEraseBlockFxn eraseBlock;
  StoreReadNumberFxn readNumber;
  StoreWriteNumberFxn writeNumber;
  StoreWriteProtect writeProtect;
  StoreRenameFxn rename;
  StoreCleanupFxn cleanup;
  StoreRemoveFxn remove;
} StoreOperations;

typedef union  __attribute__((__packed__)) {
    char sysFs[PATH_LENGTH];
    void* eepromCfg;
} StoreAttr;

typedef struct {
     char modUuid[UUID_LENGTH];
     char modName[NAME_LENGTH];
     void* storeAttr;
     const StoreOperations* storeOps;
 } ModuleMap;


int store_init();
int store_deregister_module(char *uuid);
int store_erase_block(char* uuid, off_t offset, uint16_t size);
int store_read_block(char* uuid, void* data, off_t offset , uint16_t size);
int store_read_number(char* uuid, void* val, off_t offset, uint16_t count,
                uint8_t size);
int store_register_module(UnitCfg* pCfg);
int store_register_update_module(char* uuid, UnitCfg *uCfg, uint8_t count);
int store_remove(char *uuid);
int store_rename(char* uuid, char* oldName, char* newName);
int store_write_block(char* uuid, void* data, off_t offset , uint16_t size);
int store_write_number(char* uuid, void* val, off_t offset, uint16_t count,
                uint8_t size);
int store_write_protect (char* uuid, void* data);

void store_deregister_all_module();

ModuleMap* store_choose_module(char* uuid);

#ifdef __cplusplus
}
#endif

#endif /* INC_STORE_H_ */
