/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
