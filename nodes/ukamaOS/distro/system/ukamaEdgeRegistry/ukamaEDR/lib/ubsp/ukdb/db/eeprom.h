/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef EEPROM_H_
#define EEPROM_H_

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <stdbool.h>
#include <stdint.h>
#include <string.h>

int eeprom_open(int mode);
int eeprom_read(void* data, void* buff, off_t offset, uint16_t size);
int eeprom_rename(char* old_name, char* new_name);
int eeprom_write(void* data, void* buff, off_t offset, uint16_t size) ;
int eeprom_init(void* data);
int eeprom_cleanup(void* data);
int eeprom_protect(void* data);
int eeprom_read_int(int* value, off_t offset, uint16_t count);
int eeprom_write_int(int* value, off_t offset, uint16_t count);
int eeprom_read_number(void* data, void* value, off_t offset, uint16_t count, uint8_t size);
int eeprom_write_number(void* data, void* value, off_t offset, uint16_t count, uint8_t size);
int eeprom_erase(void* data, off_t offset, uint16_t size);
int eeprom_remove(void* data);
void eeprom_close();



#endif /* EEPROM_H_ */
