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
#include "schema.h"

int invt_create_db(char *pUuid);;

int invt_deserialize_unit_cfg_data(UnitCfg **unitCfg, char *payload, uint8_t count,
                              uint16_t *size);

int invt_deserialize_module_cfg_data(ModuleCfg **modCfg, char *payload,
                                uint8_t count, uint16_t *size);

int invt_erase_db(char *pUuid);

int invt_erase_idx(char *pUuid);

int invt_erase_payload(char *pUuid, uint16_t fieldId);

int invt_get_fieldid_index(SchemaIdxTuple *index, uint16_t fid, uint8_t idxCount,
                      uint8_t *idxIter);

int invt_idb_init(void *data);

int invt_pre_create_db_setup(char *puuid);

int invt_read_bs_certs(char *pUuid, void *data, uint16_t *size);

int invt_read_cloud_certs(char *pUuid, void *data, uint16_t *size);

int invt_read_current_idx_count(char *pUuid, uint16_t *idxCount);

int invt_read_dbversion(char *pUuid, Version *ver);

int invt_read_fact_calib(char *pUuid, void *data, uint16_t *size);

int invt_read_fact_config(char *pUuid, void *data, uint16_t *size);

int invt_read_header(char *pUuid, SchemaHeader *header);

int invt_read_idx(char *pUuid, SchemaIdxTuple *p_data, uint16_t idx);

int invt_read_module(char *pUuid, ModuleInfo *modInfo);

int invt_read_module_cfg(char *pUuid, ModuleCfg *modCfg, uint8_t count,
                         uint16_t *size);

int invt_read_module_info(char *pUuid, ModuleInfo *p_info, uint16_t *size);

int invt_read_payload(char *pUuid, void *p_data, uint16_t offset,
                      uint16_t size);

int invt_read_payload_for_fieldid(char *pUuid, void *p_payload, uint16_t fid,
                                  uint16_t *size);

int invt_read_payload_from_ukdb(char *pUuid, void *p_data, uint16_t id,
                                uint16_t *size);

int invt_read_unit(char *pUuid, UnitInfo *unitInfo, UnitCfg *unitCfg);

int invt_read_unit_cfg(char *pUuid, UnitCfg *p_ucfg, uint8_t count,
                       uint16_t *size);

int invt_read_unit_info(char *pUuid, UnitInfo *unitInfo, uint16_t *size);

int invt_read_user_calib(char *pUuid, void *data, uint16_t *size);

int invt_read_user_config(char *pUuid, void *data, uint16_t *size);

int invt_register_devices(char *pUuid);

int invt_register_module(UnitCfg *cfg);

int invt_register_modules(char *pUuid);

int invt_remove_db(char *puuid);

int invt_search_field_id(char *pUuid, SchemaIdxTuple **idx_data, uint16_t *idx,
                        uint16_t fieldId);

int invt_update_current_idx_count(char *pUuid, uint16_t *idxCount);

int invt_update_dbversion(char *pUuid, Version ver);

int invt_update_idx_at_id(char *pUuid, SchemaIdxTuple *p_data, uint16_t idx);

int invt_update_idx_for_field_id(char *pUuid, void *p_data,
                                  uint16_t fieldId);

int invt_update_payload(char *pUuid, void *p_data, uint16_t fieldId,
                        uint16_t size, uint8_t state, Version version);

int invt_unregister_module(char *puuid);

int invt_validating_magic_word(char *pUuid);

int invt_validate_payload(void *p_data, uint32_t crc, uint16_t size);

int invt_write_generic_data(SchemaIdxTuple *index, char *pUuid, uint16_t fid);

int invt_write_header(char *pUuid, SchemaHeader *header);

int invt_write_index(char *pUuid, SchemaIdxTuple *p_data);

int invt_write_magic_word(char *pUuid);

int invt_write_module_cfg_data(char *pUuid, ModuleInfo *minfo,
                               SchemaIdxTuple *cfgIndex);

int invt_write_module_info_data(char *pUuid, SchemaIdxTuple *info_index,
                                SchemaIdxTuple *cfgIndex, uint8_t *modCount,
                                uint8_t *idx);

int invt_write_module_payload(char *pUuid, void *p_data, uint16_t offset,
                              uint16_t size);

int invt_write_payload(char *pUuid, void *p_data, uint16_t offset,
                       uint16_t size);

int invt_write_unit_cfg_data(char *pUuid, SchemaIdxTuple *index, uint8_t count);

int invt_write_unit_info_data(char *p1_uuid, SchemaIdxTuple *index, char *pUuid,
                              uint8_t *count);

void invt_exit();

void invt_free_module_cfg(ModuleCfg *cfg, uint8_t count);

void invt_free_unit_cfg(UnitCfg *cfg, uint8_t count);

void invt_idb_exit();

void invt_print_dev(void *dev, DeviceClass devClass);

void invt_print_dev_i2c_cfg(DevI2cCfg *pdev);

void invt_print_dev_gpio_cfg(DevGpioCfg *pdev);

void invt_print_dev_spi_cfg(DevSpiCfg *pdev);

void invt_print_dev_uart_cfg(DevUartCfg *pdev);

void invt_print_header(SchemaHeader *header);

void invt_print_index_table(SchemaIdxTuple *idxTbl, uint8_t count);

void invt_print_module_cfg(ModuleCfg *modCfg, uint8_t count);

void invt_print_module_info(ModuleInfo *modInfo);

void invt_print_unit_cfg(UnitCfg *p_ucfg, uint8_t count);

void invt_print_unit_info(UnitInfo *unitInfo);

void *invt_deserialize_devices(const char *payload, int offset, DeviceClass devClass,
                               int *size);

char *serialize_module_config_data(ModuleCfg *mcfg, uint8_t count,
                                   uint16_t *size);

char *serialize_unitcfg_payload(UnitCfg *ucfg, uint8_t count, uint16_t *size);

ModuleCfg *ukdb_alloc_module_cfg(uint8_t count);

UnitCfg *ukdb_alloc_unit_cfg(uint8_t count);

#ifdef __cplusplus
}
#endif

#endif /* INC_INVENTORY_H_ */

