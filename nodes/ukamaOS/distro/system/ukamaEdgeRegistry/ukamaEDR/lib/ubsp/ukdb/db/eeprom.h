/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
