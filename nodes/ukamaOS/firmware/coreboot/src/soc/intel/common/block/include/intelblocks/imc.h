/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 Facebook, Inc.
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

#include <device/pci.h>
#include <stdint.h>

#ifndef SOC_INTEL_COMMON_BLOCK_IMC_H
#define SOC_INTEL_COMMON_BLOCK_IMC_H

enum smbus_command { IMC_READ, IMC_WRITE };

enum access_width { IMC_DATA_BYTE, IMC_DATA_WORD };

enum memory_controller_id { IMC_CONTROLLER_ID0 = 0, IMC_CONTROLLER_ID1 };

enum device_type_id {
	IMC_DEVICE_TSOD = 0x3,
	IMC_DEVICE_WP_EEPROM = 0x6,
	IMC_DEVICE_EEPROM = 0xa
};

/* Initiate SMBus/I2C transaction to DIMM EEPROM */
int imc_smbus_spd_xfer(pci_devfn_t dev, uint8_t slave_addr, uint8_t bus_addr,
		       enum device_type_id dti, enum access_width width,
		       enum memory_controller_id mcid, enum smbus_command cmd, void *data);
#endif
