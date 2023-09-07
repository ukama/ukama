/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2019 HardenedLinux
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#ifndef CRC_BYTE_H
#define CRC_BYTE_H

#include <stdint.h>

/* This function is used to calculate crc7 byte by byte, with polynomial
 * x^7 + x^3 + 1.
 *
 * prev_crc: old crc result (0 for first)
 * data: new byte
 * return value: new crc result
 */
uint8_t crc7_byte(uint8_t prev_crc, uint8_t data);

/* This function is used to calculate crc16 byte by byte, with polynomial
 * x^16 + x^12 + x^5 + 1.
 *
 * prev_crc: old crc result (0 for first)
 * data: new byte
 * return value: new crc result
 */
uint16_t crc16_byte(uint16_t prev_crc, uint8_t data);


#endif /* CRC_BYTE_H */
