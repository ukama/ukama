%{
/*
 * sconfig, coreboot device tree compiler
 *
 * Copyright (C) 2010 coresystems GmbH
 *   written by Patrick Georgi <patrick@georgi-clan.de>
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

#include "sconfig.h"

int yylex();
void yyerror(const char *s);

static struct bus *cur_parent;
static struct chip_instance *cur_chip_instance;

%}
%union {
	struct device *dev;
	struct chip_instance *chip_instance;
	char *string;
	int number;
}

%token CHIP DEVICE REGISTER BOOL HIDDEN BUS RESOURCE END EQUALS HEX STRING PCI PNP I2C APIC CPU_CLUSTER CPU DOMAIN IRQ DRQ SLOT_DESC IO NUMBER SUBSYSTEMID INHERIT IOAPIC_IRQ IOAPIC PCIINT GENERIC SPI USB MMIO
%%
devtree: { cur_parent = root_parent; } chip;

chipchildren: chipchildren device | chipchildren chip | chipchildren registers | /* empty */ ;

devicechildren: devicechildren device | devicechildren chip | devicechildren resource | devicechildren subsystemid | devicechildren ioapic_irq | devicechildren smbios_slot_desc | /* empty */ ;

chip: CHIP STRING /* == path */ {
	$<chip_instance>$ = new_chip_instance($<string>2);
	chip_enqueue_tail(cur_chip_instance);
	cur_chip_instance = $<chip_instance>$;
}
	chipchildren END {
	cur_chip_instance = chip_dequeue_tail();
};

device: DEVICE BUS NUMBER /* == devnum */ status {
	$<dev>$ = new_device(cur_parent, cur_chip_instance, $<number>2, $<string>3, $<number>4);
	cur_parent = $<dev>$->last_bus;
}
	devicechildren END {
	cur_parent = $<dev>5->parent;
};

status: BOOL | HIDDEN;

resource: RESOURCE NUMBER /* == resnum */ EQUALS NUMBER /* == resval */
	{ add_resource(cur_parent, $<number>1, strtol($<string>2, NULL, 0), strtol($<string>4, NULL, 0)); } ;

registers: REGISTER STRING /* == regname */ EQUALS STRING /* == regval */
	{ add_register(cur_chip_instance, $<string>2, $<string>4); } ;

subsystemid: SUBSYSTEMID NUMBER NUMBER
	{ add_pci_subsystem_ids(cur_parent, strtol($<string>2, NULL, 16), strtol($<string>3, NULL, 16), 0); };

subsystemid: SUBSYSTEMID NUMBER NUMBER INHERIT
	{ add_pci_subsystem_ids(cur_parent, strtol($<string>2, NULL, 16), strtol($<string>3, NULL, 16), 1); };

ioapic_irq: IOAPIC_IRQ NUMBER PCIINT NUMBER
	{ add_ioapic_info(cur_parent, strtol($<string>2, NULL, 16), $<string>3, strtol($<string>4, NULL, 16)); };

smbios_slot_desc: SLOT_DESC STRING STRING STRING STRING
	{ add_slot_desc(cur_parent, $<string>2, $<string>3, $<string>4, $<string>5); };

smbios_slot_desc: SLOT_DESC STRING STRING STRING
	{ add_slot_desc(cur_parent, $<string>2, $<string>3, $<string>4, NULL); };

smbios_slot_desc: SLOT_DESC STRING STRING
	{ add_slot_desc(cur_parent, $<string>2, $<string>3, NULL, NULL); };

%%
