/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "drivers/eeprom_wrapper.h"

int eeprom_wrapper_open(int mode) {
    return 0;
}
int eeprom_wrapper_read(void *data, void *buff, off_t offset, uint16_t size) {
    return 0;
}
int eeprom_wrapper_rename(char *old_name, char *new_name) {
    return 0;
}
int eeprom_wrapper_write(void *data, void *buff, off_t offset, uint16_t size) {
    return 0;
}
int eeprom_wrapper_init(void *data) {
    return 0;
}
int eeprom_wrapper_cleanup(void *data) {
    return 0;
}
int eeprom_wrapper_protect(void *data) {
    return 0;
}
int eeprom_wrapper_read_int(int *value, off_t offset, uint16_t count) {
    return 0;
}
int eeprom_wrapper_write_int(int *value, off_t offset, uint16_t count) {
    return 0;
}
int eeprom_wrapper_read_number(void *data, void *value, off_t offset, uint16_t count,
                       uint8_t size) {
    return 0;
}
int eeprom_wrapper_write_number(void *data, void *value, off_t offset, uint16_t count,
                        uint8_t size) {
    return 0;
}
int eeprom_wrapper_erase(void *data, off_t offset, uint16_t size) {
    return 0;
}
int eeprom_wrapper_remove(void *data) {
    return 0;
}
void eeprom_wrapper_close() {
}
