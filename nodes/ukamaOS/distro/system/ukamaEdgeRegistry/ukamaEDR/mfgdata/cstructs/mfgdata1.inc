/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "headers/ubsp/devices.h"
#include "headers/ubsp/ukdblayout.h"

const UKDB ukdb1 = {
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
				.mod_mode = MOD_MODE_MASTER,
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
						.payload_size = 345,/*115*3 This has to proper otherwise read write error can happen*/
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
						.payload_size = 167,
						.payload_version =  { .major = 0, .minor = 0 },
						.payload_crc = 0xcbe853a8ul,
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
						.payload_size = 714, /* 119*6 This has to proper otherwise read write error can happen*/
						.payload_version =  { .major = 0, .minor = 0 },
						.payload_crc = 0x6fdc71e7ul,
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
						.payload_version =  { .major = 0, .minor = 0 },
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
						.payload_version =  { .major = 0, .minor = 0 },
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
						.payload_version =  { .major = 0, .minor = 0 },
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
						.payload_version =  { .major = 0, .minor = 0 },
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
						.payload_version =  { .major = 0, .minor = 0 },
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
						.payload_version =  { .major = 0, .minor = 0 },
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
			.uuid = "UK-0001-HNODE-SA03-1101",
			.name = "tNode",
			.unit = E_TNODELTE,
			.partno = "LTE-BAND-3-0XXXX",
			.skew = "UK_TNODE-LTE-0001",
			.mac = "10:20:30:20:50:60",
			.swver = { .major= 0, .minor = 1 },
			.pswver = { .major= 0, .minor = 1 },
			.assm_date = "30-07-2020",
			.oem_name = "SANMINA",
			.mod_count = 3,
		},
		.unitcfg = (UnitCfg[]){
			{
				.mod_uuid = "UK-1001-COM-1101",
				.mod_name = "COMv1",
				.sysfs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0050/eeprom",
				.eeprom_cfg = &(DevI2cCfg){
						.bus = 0,
						.add = 50
				}
			},
			{
				.mod_uuid = "UK-2001-LTE-1101",
				.mod_name = "LTE",
				.sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0050/eeprom",
				.eeprom_cfg = &(DevI2cCfg){
						.bus = 1,
						.add = 50
				}
			},
			{
				.mod_uuid = "UK-2001-LTE-1001",
				.mod_name = "MASK",
				.sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0051/eeprom",
				.eeprom_cfg = &(DevI2cCfg){
						.bus = 1,
						.add = 51
				}
			}


		},
		.modinfo = {
				.uuid = "UK-1001-COM-1101",
				.name = "COM-v1",
				.module = E_COMV1,
				.partno = "COMv1-X86-0XXXX",
				.hwver = "REV-A",
				.mac = "10:20:30:20:50:60",
				.swver = { .major= 0, .minor = 1 },
				.pswver = { .major= 0, .minor = 2 },
				.mfg_date = "31-07-2020",
				.mfg_name = "FOXCON",
				.dev_count = 6,
		},
		.modcfg = (ModuleCfg[]) {
			{
					.dev_name = "INA226",
					.dev_disc = "Pmic",
					.dev_type = DEV_TYPE_PWR,
					.dev_class = DEV_CLASS_I2C,
					.sysfile = "/tmp/sys/class/hwmon/hwmon0/1/ina226_1/",
					.cfg = &(DevI2cCfg){
							.bus = 0,
							.add = 71
					}
			},
			{
				.dev_name = "TMP464",
				.dev_disc = "Pmic",
				.dev_type = DEV_TYPE_TMP,
				.dev_class = DEV_CLASS_I2C,
				.sysfile = "/tmp/sys/class/hwmon/hwmon0/1/tmp464_1/",
				.cfg = &(DevI2cCfg){
					.bus = 0,
					.add = 48
				}
			},
			{
				.dev_name = "TMP464",
				.dev_disc = "DDR",
				.dev_type = DEV_TYPE_TMP,
				.dev_class = DEV_CLASS_I2C,
				.sysfile = "/tmp/sys/class/hwmon/hwmon0/1/tmp464_2/",
				.cfg = &(DevI2cCfg){
					.bus = 0,
					.add = 49
				}
			},
			{
				.dev_name = "SE98",
				.dev_disc = "X86",
				.dev_type = DEV_TYPE_TMP,
				.dev_class = DEV_CLASS_I2C,
				.sysfile = "/tmp/sys/class/hwmon/hwmon0/1/se98_1/",
				.cfg = &(DevI2cCfg){
					.bus = 0,
					.add = 41
				}
			},
			{
					.dev_name = "INA226",
					.dev_disc = "DDR",
					.dev_type = DEV_TYPE_PWR,
					.dev_class = DEV_CLASS_I2C,
					.sysfile = "/tmp/sys/class/hwmon/hwmon0/1/ina226_2/",
					.cfg = &(DevI2cCfg){
							.bus = 0,
							.add = 72
					}
			},
			{
					.dev_name = "INA226",
					.dev_disc = "Pmic",
					.dev_type = DEV_TYPE_PWR,
					.dev_class = DEV_CLASS_I2C,
					.sysfile = "/tmp/sys/class/hwmon/hwmon0/1/ina226_3/",
					.cfg = &(DevI2cCfg){
							.bus = 0,
							.add = 73
					}
			}
		},
		.factcfg = "factcfgabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.usercfg = "usercfgabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.factcalib = "factcalibabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.usercalib = "usercalibabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.bscerts = "bscertsabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.lwm2mcerts = "lwm2mcertsabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789"

};

