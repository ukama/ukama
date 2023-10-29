/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
