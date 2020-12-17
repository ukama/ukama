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

const UKDB ukdb4 = {
		.magicword =  {
				.magic_word = UKDB_MAGICWORD
		},
		.header = {
				.dbversion = { .major = 0, .minor = 0 },
				.idx_tbl_offset = UKDB_IDX_TABLE_OFFSET,
				.idx_tpl_size = UKDB_IDX_TPL_SIZE,
				.idx_tpl_max_count = UKDB_IDX_MAX_TPL_COUNT,
				.idx_cur_tpl = 0,
				.mod_cap = MOD_CAP_DEPENDENT,
				.mod_mode = MOD_MODE_MASTER,
				.mod_devown = MOD_DEV_LENDER
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
						.payload_size = 1190, /*119*10*/
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
				.uuid = "UK4001-RFFE",
				.name = "RF BOARD",
				.module = E_GSMRF,
				.partno = "GSM-RFFE-BAND3",
				.hwver = "REV-A",
				.mac = "",
				.swver = { .major= 2, .minor = 1 },
				.pswver = { .major= 2, .minor = 2 },
				.mfg_date = "31-07-2020",
				.mfg_name = "AVALON",
				.dev_count = 10,
		},
		.modcfg = (ModuleCfg[]) {
			{
				.dev_name = "DAT-31R5A-PP",
				.dev_disc = "rx-att",
				.dev_type = DEV_TYPE_ATT,
				.dev_class = DEV_CLASS_I2C,
				.sysfile = "/tmp/sys/class/hwmon/hwmon0/4/att1/",

			},
			{
					.dev_name = "DAT-31R5A-PP",
					.dev_disc = "tx-att",
					.dev_type = DEV_TYPE_ATT,
					.dev_class = DEV_CLASS_I2C,
					.sysfile = "/tmp/sys/class/hwmon/hwmon0/4/att2/",

			},
			{
					.dev_name = "ADS1015",
					.dev_disc = "rf-power detector",
					.dev_type = DEV_TYPE_ADC,
					.dev_class = DEV_CLASS_I2C,
					.sysfile = "/tmp/sys/class/hwmon/hwmon0/4/adc1/",

			},
			{
				.dev_name = "TMP464",
				.dev_disc = "RFFE Board",
				.dev_type = DEV_TYPE_TMP,
				.dev_class = DEV_CLASS_I2C,
				.sysfile = "/tmp/sys/class/hwmon/hwmon0/4/tmp464_1/",
				.cfg = &(DevI2cCfg){
					.bus = 1,
					.add = 0x49ul
				}
			},
			{
				.dev_name = "GPIO",
				.dev_disc = "PGOOD 5V",
				.dev_type = DEV_TYPE_GPIO,
				.dev_class = DEV_CLASS_GPIO,
				.sysfile = "/tmp/sys/class/gpio/gpio38/",
			},
			{
				.dev_name = "GPIO",
				.dev_disc = "PGOOD 3.3V",
				.dev_type = DEV_TYPE_GPIO,
				.dev_class = DEV_CLASS_GPIO,
				.sysfile = "/tmp/sys/class/gpio/gpio35/",
			},
			{
				.dev_name = "GPIO",
				.dev_disc = "PGOOD 5.7V",
				.dev_type = DEV_TYPE_GPIO,
				.dev_class = DEV_CLASS_GPIO,
				.sysfile = "/tmp/sys/class/gpio/gpio34/",
			},
			{
				.dev_name = "GPIO",
				.dev_disc = "PA DISABLE",
				.dev_type = DEV_TYPE_GPIO,
				.dev_class = DEV_CLASS_GPIO,
				.sysfile = "/tmp/sys/class/gpio/gpio63/",
			},
			{
				.dev_name = "GPIO",
				.dev_disc = "PGA DISABLE",
				.dev_type = DEV_TYPE_GPIO,
				.dev_class = DEV_CLASS_GPIO,
				.sysfile = "/tmp/sys/class/gpio/gpio61/",
			},
			{
				.dev_name = "GPIO",
				.dev_disc = "RF POWER DISABLE",
				.dev_type = DEV_TYPE_GPIO,
				.dev_class = DEV_CLASS_GPIO,
				.sysfile = "/tmp/sys/class/gpio/gpio40/",
			},
		},
		.factcfg = "factcfgabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.usercfg = "usercfgabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.factcalib = "factcalibabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.usercalib = "usercalibabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.bscerts = "bscertsabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
		.lwm2mcerts = "lwm2mcertsabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789"
};

