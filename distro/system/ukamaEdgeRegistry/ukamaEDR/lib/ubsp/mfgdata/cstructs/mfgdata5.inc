/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* No unit info in this just module info and other data.*/

#include "headers/ubsp/devices.h"
#include "headers/ubsp/ukdblayout.h"

const UKDB ukdb5 = {
		.magicword =  {
				.magic_word = UKDB_MAGICWORD
		},
		.header = {
				.dbversion = { .major = 0, .minor = 0 },
				.idx_tbl_offset = UKDB_IDX_TABLE_OFFSET,
				.idx_tpl_size = UKDB_IDX_TPL_SIZE,
				.idx_tpl_max_count = UKDB_IDX_MAX_TPL_COUNT,
				.idx_cur_tpl = 0,
				.mod_cap = MOD_CAP_AUTONOMOUS,
				.mod_mode = MOD_MODE_SLAVE,
				.mod_devown = MOD_DEV_OWNER
		},
		.indextable = (UKDBIdxTuple[]) {
			{
						.fieldid = FIELDID_UNIT_INFO,
						.payload_offset = UKDB_UNIT_INFO_OFFSET,
						.payload_size = 159,
						.payload_version =  { .major = 0, .minor = 0 },
						.payload_crc = 0x75567b84ul,
						.state = UKDB_FEAT_ENABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				},
				{
						.fieldid = FIELDID_UNIT_CONFIG,
						.payload_offset = UKDB_UNIT_CONFIG_OFFSET,
						.payload_size = 230,/*115*2*/
						.payload_version =  { .major = 0, .minor = 0 },
						.payload_crc = 0x6aa02949ul,
						.state = UKDB_FEAT_ENABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				},
				{
						.fieldid = FIELDID_MODULE_INFO,
						.payload_offset = UKDB_MODULE_INFO_OFFSET,
						.payload_size = 167,//TODO
						.payload_version =  { .major = 2, .minor = 1 },
						.payload_crc = 10,//TODO
						.state = UKDB_FEAT_ENABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				},
				{
						.fieldid = FIELDID_MODULE_CONFIG,
						.payload_offset = UKDB_MODULE_CONFIG_OFFSET,
						.payload_size = 595, /*119*5*/
						.payload_version =  { .major = 2, .minor = 2 },
						.payload_crc = 10,//TODO
						.state = UKDB_FEAT_ENABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				},
				{
						.fieldid = FIELDID_FACT_CONFIG,
						.payload_offset = UKDB_FACT_CONFIG_OFFSET,
						.payload_size = 151,
						.payload_version =  { .major = 1, .minor = 3 },
						.payload_crc = 0x40BE4BBCul,
						.state = UKDB_FEAT_DISABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				},
				{
						.fieldid = FIELDID_USER_CONFIG,
						.payload_offset = UKDB_USER_CONFIG_OFFSET,
						.payload_size = 151,
						.payload_version =  { .major = 2, .minor = 4 },
						.payload_crc = 0x75A4AF35ul,
						.state = UKDB_FEAT_ENABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				},
				{
						.fieldid = FIELDID_FACT_CALIB,
						.payload_offset = UKDB_FACT_CALIB_OFFSET,
						.payload_size = 153,
						.payload_version =  { .major = 2, .minor = 5 },
						.payload_crc = 0xF76620D2ul,
						.state = UKDB_FEAT_ENABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				},
				{
						.fieldid = FIELDID_USER_CALIB,
						.payload_offset = UKDB_USER_CALIB_OFFSET,
						.payload_size = 153,
						.payload_version =  { .major = 2, .minor = 6 },
						.payload_crc = 0xBA4010ABul,
						.state = UKDB_FEAT_DISABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				},
				{
						.fieldid = FIELDID_BS_CERTS,
						.payload_offset = UKDB_BS_CERTS_OFFSET,
						.payload_size = 151,
						.payload_version =  { .major = 2, .minor = 7 },
						.payload_crc = 0x13EE1DEEul,
						.state = UKDB_FEAT_ENABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				},
				{
						.fieldid = FIELDID_LWM2M_CERTS,
						.payload_offset = UKDB_LWM2M_CERTS_OFFSET,
						.payload_size = 154,
						.payload_version =  { .major = 2, .minor = 8 },
						.payload_crc = 0xDCA83D9Aul,
						.state = UKDB_FEAT_ENABLED,
						.valid = true,
						.resv1 = UKDB_DEFVAL,
						.resv2 = UKDB_DEFVAL,
						.resv3 = UKDB_DEFVAL,
						.resv4 = UKDB_DEFVAL,
						.resv5 = UKDB_DEFVAL
				}

		},
		.unitinfo = {
			.uuid = "UK5001-ANODE-0X1001",
			.name = "aNode",
			.unit = E_ANODE,
			.partno = "RFCTRL-BOARD--0XXXX",
			.skew = "UK_ANODE-GSM-0001",
			.mac = "10:20:30:20:50:60",
			.swver = { .major= 0, .minor = 1 },
			.pswver = { .major= 0, .minor = 1 },
			.assm_date = "30-08-2020",
			.oem_name = "SANMINA",
			.mod_count = 2,
		},
		.unitcfg = (UnitCfg[]){
			{
				.mod_uuid = "UK5001-RFCTRL",
				.mod_name = "RF CTRL BOARD",
				.sysfs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0051/eeprom",
				.eeprom_cfg = &(DevI2cCfg){
						.bus = 2,
						.add = 0x50ul
				}
			},
			{
				.mod_uuid = "UK4001-RFFE",
				.mod_name = "RF BOARD",
				.sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0052/eeprom",
				.eeprom_cfg = &(DevI2cCfg){
						.bus = 1,
						.add = 0x50ul
				}
			}
		},
		.modinfo = {
				.uuid = "UK5001-RFCTRL",
				.name = "RF CTRL BOARD",
				.module = E_CNTRL,
				.partno = "RF-CTRL-1001",
				.hwver = "REV-A",
				.mac = "",
				.swver = { .major= 2, .minor = 1 },
				.pswver = { .major= 2, .minor = 2 },
				.mfg_date = "31-07-2020",
				.mfg_name = "AVALON",
				.dev_count = 5,
		},
		.modcfg = (ModuleCfg[]) {
			{
				.dev_name = "SE98",
				.dev_disc = "RF MicroProcessor",
				.dev_type = DEV_TYPE_TMP,
				.dev_class = DEV_CLASS_I2C,
				.sysfile = "/tmp/sys/class/hwmon/hwmon0/5/se98_1/",
				.cfg = &(DevI2cCfg){
					.bus = 1,
					.add = 0x18ul
				}
			},
			{
					.dev_name = "LED-TRICOLOR",
					.dev_disc = "RF LED 0",
					.dev_type = DEV_TYPE_LED,
					.dev_class = DEV_CLASS_GPIO,
					.sysfile = "/tmp/sys/class/led/led0/",

			},
			{
					.dev_name = "LED-TRICOLOR",
					.dev_disc = "RF LED 1",
					.dev_type = DEV_TYPE_LED,
					.dev_class = DEV_CLASS_GPIO,
					.sysfile = "/tmp/sys/class/led/led1/",

			},
			{
					.dev_name = "LED-TRICOLOR",
					.dev_disc = "RF LED 2",
					.dev_type = DEV_TYPE_LED,
					.dev_class = DEV_CLASS_GPIO,
					.sysfile = "/tmp/sys/class/led/led2/",

			},
			{
					.dev_name = "LED-TRICOLOR",
					.dev_disc = "RF LED 3",
					.dev_type = DEV_TYPE_LED,
					.dev_class = DEV_CLASS_GPIO,
					.sysfile = "/tmp/sys/class/led/led3/",

			},

		},
		.factcfg = "factcfgabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.usercfg = "usercfgabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.factcalib = "factcalibabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.usercalib = "usercalibabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.bscerts = "bscertsabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.lwm2mcerts = "lwm2mcertsabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789"
};

