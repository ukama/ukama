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

 /**
 * @fn      int store_init()
 * @brief   Creates a store list to hold modules present in unit.
 *
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_init();

/**
 * @fn      int store_deregister_module(char*)
 * @brief   Remove the module with uuid from the store list.
 *
 * @param   uuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_deregister_module(char *uuid);

/**
 * @fn      int store_erase_block(char*, off_t, uint16_t)
 * @brief   Erase the block of size bytes from the inventory db at offset
 *          and overwrite it with 0xff.
 *
 * @param   uuid
 * @param   offset
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_erase_block(char* uuid, off_t offset, uint16_t size);

/**
 * @fn      int store_read_block(char*, void*, off_t, uint16_t)
 * @brief   read the block of data of size bytes at offset from beginning of the
 *          inventory db.
 *
 * @param   uuid
 * @param   data
 * @param   offset
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_read_block(char* uuid, void* data, off_t offset , uint16_t size);

/**
 * @fn      int store_read_number(char*, void*, off_t, uint16_t, uint8_t)
 * @brief   read the count amount of numbers at offset from beginning of the
 *          inventory db.
 *
 * @param   uuid
 * @param   val
 * @param   offset
 * @param   count
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_read_number(char* uuid, void* val, off_t offset, uint16_t count,
                uint8_t size);
/**
 * @fn      int store_register_module(UnitCfg*)
 * @brief   Register a module to store by adding its detail to the store list.
 *
 * @param   pCfg
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_register_module(UnitCfg* pCfg);

/**
 * @fn      int store_register_update_module(char*, UnitCfg*, uint8_t)
 * @brief   Update the registered module in store list.
 *
 * @param   uuid
 * @param   uCfg
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_register_update_module(char* uuid, UnitCfg *uCfg, uint8_t count);

/**
 * @fn      int store_remove(char*)
 * @brief   removes the inventory db for the module uuid.
 *
 * @param   uuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_remove(char *uuid);

/**
 * @fn      int store_rename(char*, char*, char*)
 * @brief   renames the inventory db for the module uuid.
 *
 * @param   uuid
 * @param   oldName
 * @param   newName
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_rename(char* uuid, char* oldName, char* newName);

/**
 * @fn      int store_write_block(char*, void*, off_t, uint16_t)
 * @brief   write the block of data of size bytes at offset from beginning of the
 *          inventory db.
 *
 * @param   uuid
 * @param   data
 * @param   offset
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_write_block(char* uuid, void* data, off_t offset , uint16_t size);

/**
 * @fn      int store_write_number(char*, void*, off_t, uint16_t, uint8_t)
 * @brief   write the count amount of numbers at offset from beginning of the
 *          inventory db.
 *
 * @param   uuid
 * @param   val
 * @param   offset
 * @param   count
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_write_number(char* uuid, void* val, off_t offset, uint16_t count,
                uint8_t size);
/**
 * @fn      int store_write_protect(char*, void*)
 * @brief   Enable write protection for the inventory db of the module uuid.
 *
 * @param   uuid
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int store_write_protect (char* uuid, void* data);

/**
 * @fn      void store_deregister_all_module()
 * @brief   Remove all the module from the module list.
 *
 */
void store_deregister_all_module();

/**
 * @fn      ModuleMap store_choose_module*(char*)
 * @brief   Selects the inventory db based on the UUID of the module.
 *
 * @param   uuid
 * @return  On success, 0
 *          On failure, non zero value
 */
ModuleMap* store_choose_module(char* uuid);

#ifdef __cplusplus
}
#endif

#endif /* INC_STORE_H_ */
