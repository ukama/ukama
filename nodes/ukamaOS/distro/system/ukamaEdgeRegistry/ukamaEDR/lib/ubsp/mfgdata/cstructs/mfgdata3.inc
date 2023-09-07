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

const UKDB ukdb3 = {
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
						.payload_size = 238, /*119*2*/
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
		.modinfo = {
				.uuid = "UK-2001-LTE-1001",
				.name = "MASK",
				.module = E_MASK,
				.partno = "MASK-WIFI-0XXXX",
				.hwver = "REV-A",
				.mac = "",
				.swver = { .major= 2, .minor = 1 },
				.pswver = { .major= 2, .minor = 2 },
				.mfg_date = "31-07-2020",
				.mfg_name = "UBQ",
				.dev_count = 2,
		},
		.modcfg = (ModuleCfg[]) {
			{
				.dev_name = "ADT7481",
				.dev_disc = "Wifi Controller",
				.dev_type = DEV_TYPE_TMP,
				.dev_class = DEV_CLASS_I2C,
				.sysfile = "/tmp/sys/class/hwmon/hwmon0/3/adt7481_1/",
				.cfg = &(DevI2cCfg){
					.bus = 0,
					.add = 48
				}
			},
			{
				.dev_name = "INA226",
				.dev_disc = "PCI",
				.dev_type = DEV_TYPE_PWR,
				.dev_class = DEV_CLASS_I2C,
				.sysfile = "/tmp/sys/class/hwmon/hwmon0/3/ina226_1/",
				.cfg = &(DevI2cCfg){
					.bus = 2,
					.add = 50
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

