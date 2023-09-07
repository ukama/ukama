/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "ukdb/db/eeprom.h"

int eeprom_open(int mode) {
    return 0;
}
int eeprom_read(void *data, void *buff, off_t offset, uint16_t size) {
    return 0;
}
int eeprom_rename(char *old_name, char *new_name) {
    return 0;
}
int eeprom_write(void *data, void *buff, off_t offset, uint16_t size) {
    return 0;
}
int eeprom_init(void *data) {
    return 0;
}
int eeprom_cleanup(void *data) {
    return 0;
}
int eeprom_protect(void *data) {
    return 0;
}
int eeprom_read_int(int *value, off_t offset, uint16_t count) {
    return 0;
}
int eeprom_write_int(int *value, off_t offset, uint16_t count) {
    return 0;
}
int eeprom_read_number(void *data, void *value, off_t offset, uint16_t count,
                       uint8_t size) {
    return 0;
}
int eeprom_write_number(void *data, void *value, off_t offset, uint16_t count,
                        uint8_t size) {
    return 0;
}
int eeprom_erase(void *data, off_t offset, uint16_t size) {
    return 0;
}
int eeprom_remove(void *data) {
    return 0;
}
void eeprom_close() {
}
