/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef EEPROM_H_
#define EEPROM_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

int eeprom_wrapper_cleanup(void* data);
int eeprom_wrapper_erase(void* data, off_t offset, uint16_t size);
int eeprom_wrapper_init(void* data);
int eeprom_wrapper_open(int mode);
int eeprom_wrapper_read(void* data, void* buff, off_t offset, uint16_t size);
int eeprom_wrapper_read_int(int* value, off_t offset, uint16_t count);
int eeprom_wrapper_read_number(void* data, void* value, off_t offset, uint16_t count,
                uint8_t size);
int eeprom_wrapper_remove(void* data);
int eeprom_wrapper_rename(char* old_name, char* new_name);
int eeprom_wrapper_write(void* data, void* buff, off_t offset, uint16_t size) ;
int eeprom_wrapper_write_int(int* value, off_t offset, uint16_t count);
int eeprom_wrapper_write_number(void* data, void* value, off_t offset, uint16_t count,
                uint8_t size);
int eeprom_wrapper_protect(void* data);
void eeprom_close();

#ifdef __cplusplus
}
#endif

#endif /* EEPROM_H_ */
