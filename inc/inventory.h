/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_INVENTORY_H_
#define INC_INVENTORY_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "device.h"
#include "ledger.h"
#include "schema.h"

/**
 * @fn      int invt_create_db(char*)
 * @brief   Creates inventory database for module uuid from the parsed data by
 *          mfg parser from the schema provided.
 *
 * @param   pUuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_create_db(char *pUuid);;

/**
 * @fn      int invt_deserialize_node_cfg_data(NodeCfg**, char*,
 *          uint8_t, uint16_t*)
 * @brief   Creates the Unit config from the block of data read from the
 *          inventory database. Count specifies numbers of module present.
 *
 * @param   nodeCfg
 * @param   payload
 * @param   count
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_deserialize_node_cfg_data(NodeCfg **nodeCfg, char *payload,
                uint8_t count, uint16_t *size);

/**
 * @fn      int invt_deserialize_module_cfg_data(ModuleCfg**, char*,
 * uint8_t, uint16_t*)
 * @brief   Creates the module config from the block of data read from the
 *          inventory database.  Count specifies numbers of sensor devices
 *          present.
 *
 * @param   modCfg
 * @param   payload
 * @param   count
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_deserialize_module_cfg_data(ModuleCfg **modCfg, char *payload,
                                uint8_t count, uint16_t *size);
/**
 * @fn      int invt_erase_db(char*)
 * @brief   Erase the complete database for the module uuid.
 *
 * @param   pUuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_erase_db(char *pUuid);

/**
 * @fn      int invt_erase_idx(char*)
 * @brief   Erases the last index .i.e index present at the offset
 *          SCH_IDX_CUR_TPL_COUNT_OFFSET
 *
 * @param   pUuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_erase_idx(char *pUuid);

/**
 * @fn      int invt_erase_payload(char*, uint16_t)
 * @brief   Erases payload for the field id for module with uuid in argument.
 *
 * @param   pUuid
 * @param   fieldId
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_erase_payload(char *pUuid, uint16_t fieldId);

/**
 * @fn      int invt_get_field_id_idx(SchemaIdxTuple*, uint16_t,
 *          uint8_t, uint8_t*)
 * @brief   Searches for the field id in the index table.
 *
 * @param   index
 * @param   fid
 * @param   idxCount
 * @param   idxIter
 * @return  On success, positive value i.e index for field id in the index table
 *          On failure, -1
 */
int invt_get_field_id_idx(SchemaIdxTuple *index, uint16_t fid,
                uint8_t idxCount, uint8_t *idxIter);
/**
 * @fn      int invt_get_master_node_cfg(NodeCfg*, char*)
 * @brief   Read the node config pointed by the invtLnkDb soft link.
 *
 * @param   pcfg
 * @param   invtLnkDb
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_get_master_node_cfg(NodeCfg *pcfg, char *invtLnkDb);

/**
 * @fn      int invt_init(char*, RegisterDeviceCB)
 * @brief   Reads the master node config and then register the module to store
 *          so that its inventory database can be accessed.
 *
 * @param   invtDb
 * @param   regCb
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_init(char *invtDb, RegisterDeviceCB regCb);

/**
 * @fn      int invt_mfg_init(void*)
 * @brief   initializes the manufacturing module whose job is to parse the
 *          manufacturing schema for creating inventory database.
 *
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_mfg_init(void *data);

/**
 * @fn      int invt_pre_create_store_setup(char*)
 * @brief   This function is meant for testing function like creating new
 *          inventory database. It provide dummy values for the node configs.
 *          needs to be removed at some point.
 *
 * @param   puuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_pre_create_store_setup(char *puuid);

/**
 * @fn      int invt_read_bs_certs(char*, void*, uint16_t*)
 * @brief   read boot strap certs from the store of module.
 *          Store abstract the access to inventory database.
 *
 * @param   pUuid
 * @param   data
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_bs_certs(char *pUuid, void *data, uint16_t *size);

/**
 * @fn      int invt_read_cloud_certs(char*, void*, uint16_t*)
 * @brief   read cloud certs from the store of module.
 *          Store abstract the access to inventory database.
 *
 * @param   pUuid
 * @param   data
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_cloud_certs(char *pUuid, void *data, uint16_t *size);

/**
 * @fn      int invt_read_current_idx_count(char*, uint16_t*)
 * @brief   read current index count from the store of module.
 *          Store abstract the access to inventory database.
 *
 * @param   pUuid
 * @param   idxCount
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_current_idx_count(char *pUuid, uint16_t *idxCount);

/**
 * @fn      int invt_read_dbversion(char*, Version*)
 * @brief   read inventory db verison from the store of module.
 *          Store abstract the access to inventory database.
 *
 * @param   pUuid
 * @param   ver
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_dbversion(char *pUuid, Version *ver);

/**
 * @fn      int invt_read_fact_calib(char*, void*, uint16_t*)
 * @brief   read factory calibration from the store of module.
 *
 * @param   pUuid
 * @param   data
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_fact_calib(char *pUuid, void *data, uint16_t *size);

/**
 * @fn      int invt_read_fact_config(char*, void*, uint16_t*)
 * @brief   read factory config from the store of module.
 *
 * @param   pUuid
 * @param   data
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_fact_config(char *pUuid, void *data, uint16_t *size);

/**
 * @fn      int invt_read_header(char*, SchemaHeader*)
 * @brief   read header from the store of module.
 *
 * @param   pUuid
 * @param   header
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_header(char *pUuid, SchemaHeader *header);

/**
 * @fn      int invt_read_idx(char*, SchemaIdxTuple*, uint16_t)
 * @brief   read index pointed by idx from the store of module.
 *
 * @param   pUuid
 * @param   p_data
 * @param   idx
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_idx(char *pUuid, SchemaIdxTuple *p_data, uint16_t idx);

/**
 * @fn      int invt_read_module_cfg(char*, ModuleCfg*, uint8_t, uint16_t*)
 * @brief   read count amount of module config from the store of module.
 *
 * @param   pUuid
 * @param   modCfg
 * @param   count
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_module_cfg(char *pUuid, ModuleCfg *modCfg, uint8_t count,
                         uint16_t *size);
/**
 * @fn      int invt_read_module_info(char*, ModuleInfo*, uint16_t*)
 * @brief   read module info from the store of module.
 *
 * @param   pUuid
 * @param   p_info
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_module_info(char *pUuid, ModuleInfo *p_info, uint16_t *size);

/**
 * @fn      int invt_read_payload(char*, void*, uint16_t, uint16_t)
 * @brief   Inventory read payload from the store of size bytes at offset.
 *
 * @param   pUuid
 * @param   p_data
 * @param   offset
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_payload(char *pUuid, void *p_data, uint16_t offset,
                      uint16_t size);

/**
 * @fn      int invt_read_payload_for_field_id(char*, void**, uint16_t, uint16_t*)
 * @brief   Inventory read the payload pointed by the field id.
 *
 * @param   pUuid
 * @param   data
 * @param   fid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_payload_for_field_id(char *pUuid, void **data, uint16_t fid,
                                  uint16_t *size);
/**
 * @fn      int invt_read_payload_from_invt(char*, void*, uint16_t, uint16_t*)
 * @brief   read the payload from the store with field id.
 *
 * @param   pUuid
 * @param   p_data
 * @param   id
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_payload_from_store(char *pUuid, void *p_data, uint16_t id,
                                uint16_t *size);

/**
 * @fn      int invt_read_node_cfg(char*, NodeCfg*, uint8_t, uint16_t*)
 * @brief   reads the count amount of node configs from the store.
 *
 * @param   pUuid
 * @param   p_ucfg
 * @param   count
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_node_cfg(char *pUuid, NodeCfg *p_ucfg, uint8_t count,
                       uint16_t *size);
/**
 * @fn      int invt_read_node_info(char*, NodeInfo*, uint16_t*)
 * @brief   reads the node info from the store.
 *
 * @param   pUuid
 * @param   nodeInfo
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_node_info(char *pUuid, NodeInfo *nodeInfo, uint16_t *size);

/**
 * @fn      int invt_read_user_calib(char*, void*, uint16_t*)
 * @brief   reads the user calibration data from the store.
 *
 * @param   pUuid
 * @param   data
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_user_calib(char *pUuid, void *data, uint16_t *size);

/**
 * @fn      int invt_read_user_config(char*, void*, uint16_t*)
 * @brief   reads the user config data from the store.
 *
 * @param   pUuid
 * @param   data
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_read_user_config(char *pUuid, void *data, uint16_t *size);

/**
 * @fn      int invt_register_devices(char*, RegisterDeviceCB)
 * @brief   register the sensor devices present under the module config for
 *          module pUuid to ledger.
 *
 * @param   pUuid
 * @param   registerDev
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_register_devices(char *pUuid, RegisterDeviceCB registerDev);

/**
 * @fn      int invt_register_module(NodeCfg*)
 * @brief   register modules pointed by node config to the store.
 *
 * @param   cfg
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_register_module(NodeCfg *cfg);

/**
 * @fn      int invt_register_modules(char*, RegisterDeviceCB)
 * @brief   register the modules present in the master module node config to
 *          store and the sensor devices present under these modules to ledger.
 *
 * @param   pUuid
 * @param   registerDev
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_register_modules(char *pUuid, RegisterDeviceCB registerDev);

/**
 * @fn      int invt_remove_db(char*)
 * @brief   deletes the the inventory database and module from the store.
 *
 * @param   puuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_remove_db(char *puuid);

/**
 * @fn      int invt_search_field_id(char*, SchemaIdxTuple**, uint16_t*,
 *          uint16_t)
 * @brief   Search for the field id with valid data. Valid data is checked by
 *          looking at the valid member variable in the index tuple of the
 *          index table.
 *
 * @param   pUuid
 * @param   idx_data
 * @param   idx
 * @param   fieldId
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_search_field_id(char *pUuid, SchemaIdxTuple **idx_data, uint16_t *idx,
                        uint16_t fieldId);

/**
 * @fn      int invt_update_current_idx_count(char*, uint16_t*)
 * @brief   Update the current index count in the inventory header.
 *
 * @param   pUuid
 * @param   idxCount
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_update_current_idx_count(char *pUuid, uint16_t *idxCount);

/**
 * @fn      int invt_update_dbversion(char*, Version)
 * @brief   Update the version in the inventory header.
 *
 * @param   pUuid
 * @param   ver
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_update_dbversion(char *pUuid, Version ver);

/**
 * @fn      int invt_update_idx_at_id(char*, SchemaIdxTuple*, uint16_t)
 * @brief   Update the index present at the idx in the inventory index table.
 *
 * @param   pUuid
 * @param   p_data
 * @param   idx
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_update_idx_at_id(char *pUuid, SchemaIdxTuple *p_data, uint16_t idx);

/**
 * @fn      int invt_update_idx_for_field_id(char*, void*, uint16_t)
 * @brief   Update the index for the field id in the inventory index table.
 *
 * @param   pUuid
 * @param   p_data
 * @param   fieldId
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_update_idx_for_field_id(char *pUuid, void *p_data,
                                  uint16_t fieldId);
/**
 * @fn      int invt_update_payload(char*, void*, uint16_t, uint16_t)
 * @brief   Updates the payload for the field id and the index entry in the
 *          index table with the new values.
 *
 * @param   pUuid
 * @param   p_data
 * @param   fieldId
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_update_payload(char *pUuid, void *p_data, uint16_t fieldId,
                        uint16_t size );
/**
 * @fn      int invt_unregister_module(char*)
 * @brief   Remove the module from the store list.
 *
 * @param   puuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_unregister_module(char *puuid);

/**
 * @fn      int invt_validating_magic_word(char*)
 * @brief   Validate the magic word from the inventory header.
 *
 * @param   pUuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_validating_magic_word(char *pUuid);

/**
 * @fn      int invt_validate_payload(void*, uint32_t, uint16_t)
 * @brief   Validate the payload by comparing stored crc values and calculated
 *           crc value.
 *
 * @param   p_data
 * @param   crc
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_validate_payload(void *p_data, uint32_t crc, uint16_t size);

/**
 * @fn      int invt_write_generic_data(SchemaIdxTuple*, char*, uint16_t)
 * @brief   Write generic payload data to inventory data base (store).
 *          These API are used while creating inventory database for fields
 *          like certificates and configs etc. which are provided as file input
 *          in the mfg data.
 *
 * @param   index
 * @param   pUuid
 * @param   fid
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_write_generic_data(SchemaIdxTuple *index, char *pUuid, uint16_t fid);

/**
 * @fn      int invt_write_header(char*, SchemaHeader*)
 * @brief   Writes the header info for the inventory database.
 *
 * @param   pUuid
 * @param   header
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_write_header(char *pUuid, SchemaHeader *header);

/**
 * @fn      int invt_write_index(char*, SchemaIdxTuple*)
 * @brief   Writes the index tuple pointed by  tuple to inventory database of
 *          module uuid.
 *
 * @param   uuid
 * @param   idTtuple
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_write_index(char *pUuid, SchemaIdxTuple *tuple);

/**
 * @fn      int invt_write_magic_word(char*)
 * @brief   Writes the magic word to inventory database of module uuid.
 *
 * @param   pUuid
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_write_magic_word(char *pUuid);

/**
 * @fn      int invt_write_module_cfg_data(char*, ModuleInfo*, SchemaIdxTuple*)
 * @brief   Writes the module config to inventory database of module uuid and
 *          updates the index table for it.
 *
 * @param   pUuid
 * @param   minfo
 * @param   cfgIndex
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_write_module_cfg_data(char *pUuid, ModuleInfo *minfo,
                               SchemaIdxTuple *cfgIndex);
/**
 * @fn      int invt_write_module_info_data(char*, SchemaIdxTuple*, SchemaIdxTuple*, uint8_t*, uint8_t*)
 * @brief   Writes the module info to inventory database of module uuid and
 *          updates the index table for it.
 *
 * @param   pUuid
 * @param   info_index
 * @param   cfgIndex
 * @param   modCount
 * @param   idx
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_write_module_info_data(char *pUuid, SchemaIdxTuple *info_index,
                                SchemaIdxTuple *cfgIndex, uint8_t *modCount,
                                uint8_t *idx);
/**
 * @fn      int invt_write_module_payload(char*, void*, uint16_t, uint16_t)
 * @brief   Write the size bytes of  module payload to the store at some offset
 *          to module uuid inventory database.
 *
 * @param   pUuid
 * @param   data
 * @param   offset
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_write_module_payload(char *pUuid, void *data, uint16_t offset,
                              uint16_t size);

/**
 * @fn      int invt_write_node_cfg_data(char*, SchemaIdxTuple*, uint8_t)
 * @brief   Write Unit Config and update the index to index table of the
 *          inventory data base.
 *
 * @param   pUuid
 * @param   index
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_write_node_cfg_data(char *pUuid, SchemaIdxTuple *index, uint8_t count);

/**
 * @fn      int invt_write_node_info_data(char*, SchemaIdxTuple*,
 *          char*, uint8_t*)
 * @brief   Write Unit info and update the index to index table of the
 *          inventory data base.
 *
 * @param   p1_uuid
 * @param   index
 * @param   pUuid
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int invt_write_node_info_data(char *p1_uuid, SchemaIdxTuple *index, char *pUuid,
                              uint8_t *count);
/**
 * @fn      void invt_exit()
 * @brief   de-register all the modules from the store.
 *
 */
void invt_exit();

/**
 * @fn      void invt_free_module_cfg(ModuleCfg*, uint8_t)
 * @brief   free the memory allocated for the module config.
 *
 * @param   cfg
 * @param   count
 */
void invt_free_module_cfg(ModuleCfg *cfg, uint8_t count);

/**
 * @fn      void invt_free_node_cfg(NodeCfg*, uint8_t)
 * @brief   free the memory allocated for the node config.
 *
 * @param   cfg
 * @param   count
 */
void invt_free_node_cfg(NodeCfg *cfg, uint8_t count);

/**
 * @fn      void invt_mfg_exit()
 * @brief   Clean the memory and data structs used by the mfg parser for
 *          schema parsing.
 *
 */
void invt_mfg_exit();

/**
 * @fn      void invt_print_dev(void*, DeviceClass)
 * @brief   Logs the sensor device details.
 *
 * @param   dev
 * @param   devClass
 */
void invt_print_dev(void *dev, DeviceClass devClass);

/**
 * @fn      void invt_print_dev_i2c_cfg(DevI2cCfg*)
 * @brief   Logs the I2C sensor device details.
 *
 * @param   pdev
 */
void invt_print_dev_i2c_cfg(DevI2cCfg *pdev);

/**
 * @fn      void invt_print_dev_gpio_cfg(DevGpioCfg*)
 * @brief   Logs the GPIO sensor device details.
 *
 * @param   pdev
 */
void invt_print_dev_gpio_cfg(DevGpioCfg *pdev);

/**
 * @fn      void invt_print_dev_spi_cfg(DevSpiCfg*)
 * @brief   Logs the SPI sensor device details.
 *
 * @param   pdev
 */
void invt_print_dev_spi_cfg(DevSpiCfg *pdev);

/**
 * @fn      void invt_print_dev_uart_cfg(DevUartCfg*)
 * @brief   Logs the UART sensor device details.
 *
 * @param   pdev
 */
void invt_print_dev_uart_cfg(DevUartCfg *pdev);

/**
 * @fn      void invt_print_header(SchemaHeader*)
 * @brief   Logs the header details for inventory database.
 *
 * @param   header
 */
void invt_print_header(SchemaHeader *header);

/**
 * @fn      void invt_print_index_table(SchemaIdxTuple*, uint8_t)
 * @brief   Logs the index table for inventory databse.
 *
 * @param   idxTbl
 * @param   count
 */
void invt_print_index_table(SchemaIdxTuple *idxTbl, uint8_t count);

/**
 * @fn      void invt_print_module_cfg(ModuleCfg*, uint8_t)
 * @brief   Logs the module config
 *
 * @param   modCfg
 * @param   count
 */
void invt_print_module_cfg(ModuleCfg *modCfg, uint8_t count);

/**
 * @fn      void invt_print_module_info(ModuleInfo*)
 * @brief   Logs the module info
 *
 * @param   modInfo
 */
void invt_print_module_info(ModuleInfo *modInfo);

/**
 * @fn      void invt_print_node_cfg(NodeCfg*, uint8_t)
 * @brief   Logs the node config
 *
 * @param   p_ucfg
 * @param   count
 */
void invt_print_node_cfg(NodeCfg *p_ucfg, uint8_t count);

/**
 * @fn      void invt_print_node_info(NodeInfo*)
 * @brief   Logs the node info
 *
 * @param   nodeInfo
 */
void invt_print_node_info(NodeInfo *nodeInfo);

/**
 * @fn      void invt_deserialize_devices*(const char*, int, DeviceClass, int*)
 * @brief   deserialize the data read from the inevntory data base to the
 *          device structs.
 *
 * @param   payload
 * @param   offset
 * @param   devClass
 * @param   size
 * @return  On success, pointer to type of the sensor device
 *          On failure, NULL
 */
void *invt_deserialize_devices(const char *payload, int offset, DeviceClass devClass,
                               int *size);
/**
 * @fn      char serialize_module_config_data*(ModuleCfg*, uint8_t, uint16_t*)
 * @brief   serialize the module config into the block of data bytes which
 *          can then be written to inventory database.
 *
 * @param   mcfg
 * @param   count
 * @param   size
 * @return  On success, pointer to block of data bytes
 *          On failure, NULL
 */
char *serialize_module_config_data(ModuleCfg *mcfg, uint8_t count,
                                   uint16_t *size);
/**
 * @fn      char serialize_unitcfg_payload*(NodeCfg*, uint8_t, uint16_t*)
 * @brief   serialize the node config into the block of data bytes which
 *          can then be written to inventory database.
 *
 * @param   ucfg
 * @param   count
 * @param   size
 * @return  On success, pointer to block of data bytes
 *          On failure, NULL
 */
char *serialize_unitcfg_payload(NodeCfg *ucfg, uint8_t count, uint16_t *size);

/**
 * @fn      ModuleCfg invt_alloc_module_cfg*(uint8_t)
 * @brief   Allocate the memory for the module config with count number of
 *          sensor devices.
 *
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
ModuleCfg *invt_alloc_module_cfg(uint8_t count);

/**
 * @fn      NodeCfg invt_alloc_node_cfg*(uint8_t)
 * @brief   Allocate the memory for the node config with count number of
 *          modules.
 *
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
NodeCfg *invt_alloc_node_cfg(uint8_t count);

#ifdef __cplusplus
}
#endif

#endif /* INC_INVENTORY_H_ */

