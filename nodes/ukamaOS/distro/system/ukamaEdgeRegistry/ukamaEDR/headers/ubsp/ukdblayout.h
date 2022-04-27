/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INCLUDERR_UBSP_LAYOUT_H_
#define INCLUDERR_UBSP_LAYOUT_H_

#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#define UKDB_START_OFFSET				0x0000
#define UKDB_END_OFFSET					0xFFFF

#define UKDB_MAGICWORD_OFFSET 			0x0000
#define UKDB_MAGCIWORD_SIZE				0x0004
#define UKDB_MAGICWORD 					0xDEADBEEF
#define UKDB_DEFVAL						0xFFFF

#define UKDB_HEADER_OFFSET           	0x0010
#define UKDB_HEADER_DBVER_OFFSET        0x0010
#define UKDB_HEADER_SIZE             	0x0018
#define UKDB_MAJOR_VER					0
#define UKDB_MINOR_VER 					0



#define UKDB_IDX_TABLE_OFFSET 			0x0040
#define UKDB_IDX_TPL_SIZE           	0x0018
#define UKDB_IDX_MAX_TPL_COUNT          0x0032  /* Max Index table entries 50 */
#define UKDB_IDX_CUR_TPL_COUNT_OFFSET   0x0020
#define UKDB_IDX_COUNT_SIZE				0x0002


#define UKDB_FOOTER_OFFSET              0x09A0
#define UKDB_FOOTER_SIZE            	0x0050

/* MAX_PAYLOAD_SIZE */
#define UKDB_MAX_PAYLOAD_SIZE			0x1000

/* Size */
#define UKDB_UNITINFO_PAYLOADSIZE		0x009F
#define UKDB_UNITCONFIG_PAYLOADSIZE		0x0073
#define UKDB_MODULEINFO_PAYLOADSIZE		0x00A7
#define UKDB_MODULECONFIG_PAYLOADSIZE	0x0077
#define UKDB_FACT_CONFIG_SIZE			0x1000
#define UKDB_USER_CONFIG_SIZE			0x1000
#define UKDB_FACT_CALIB_SIZE			0x1000
#define UKDB_USER_CALIB_SIZE			0x1000
#define UKDB_BS_CERTS_SIZE				0x1000
#define UKDB_LWM2M_CERTS_SIZE			0x1000


/* OFFSETS*/
#define MAX_NUMBER_MODULES_PER_UNIT		0x0008
#define MAX_NUMBER_DEVICES_PER_MODULE	0x0014
#define UKDB_PAYLAOD_OFFSET             0x0A00

#define UKDB_UNIT_INFO_OFFSET			0x0A00 //Offset: 2560 (64*40)  		size = 192B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_UNIT_CONFIG_OFFSET			0x0AC0 //Offset: 2752 (64*43)  		size = 1024B 	REQ: = 115B		MAX = 8 Entry
#define UKDB_MODULE_INFO_OFFSET			0x0EC0 //Offset: 3776 (64*59)  		size = 192B 	REQ  = 167B		MAX = 1 Entry
#define UKDB_MODULE_CONFIG_OFFSET		0x0F80 //Offset: 3968 (64*62)  		size = 2432B 	REQ  = 119B		MAX = 20 Entry
#define UKDB_FACT_CONFIG_OFFSET			0x1900 //Offset: 6400 (64*100) 		size = 5120B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_USER_CONFIG_OFFSET			0x2900 //Offset: 10496 (64*164) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_FACT_CALIB_OFFSET			0x3900 //Offset: 14592 (64*228) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_USER_CALIB_OFFSET			0x4900 //Offset: 18688 (64*292) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_BS_CERTS_OFFSET			0x5900 //Offset: 22784 (64*256) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_LWM2M_CERTS_OFFSET			0x6900 //Offset: 26880 (64*320) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry

#if 0
#define UKDB_PAYLAOD_OFFSET             0x0A00
#define UKDB_UNIT_INFO_OFFSET			0x0A00 //Offset: 2560 (64*40)  		size = 192B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_UNIT_CONFIG_OFFSET			0x0AC0 //Offset: 2752 (64*43)  		size = 1024B 	REQ: = 115B		MAX = 8 Entry
#define UKDB_MODULE_INFO_OFFSET			0x0EC0 //Offset: 3776 (64*59)  		size = 192B 	REQ  = 167B		MAX = 1 Entry
#define UKDB_MODULE_CONFIG_OFFSET		0x0F80 //Offset: 3968 (64*62)  		size = 2432B 	REQ  = 119B		MAX = 20 Entry
#define UKDB_FACT_CONFIG_OFFSET			0x1900 //Offset: 6400 (64*100) 		size = 5120B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_USER_CONFIG_OFFSET			0x2900 //Offset: 10496 (64*164) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_FACT_CALIB_OFFSET			0x3900 //Offset: 14592 (64*228) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_USER_CALIB_OFFSET			0x4900 //Offset: 18688 (64*292) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_BS_CERTS_OFFSET			0x5900 //Offset: 22784 (64*256) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry
#define UKDB_LWM2M_CERTS_OFFSET			0x6900 //Offset: 26880 (64*320) 	size = 4096B 	REQ  = 159B		MAX = 1 Entry
#endif
#define UKDB_FEAT_ENABLED				0x01
#define UKDB_FEAT_DISABLED				0x00




/* Field Id */
#define FIELDID_UNIT_INFO				0x0001
#define FIELDID_UNIT_CONFIG				0x0002
#define FIELDID_MODULE_INFO				0x0003
#define FIELDID_MODULE_CONFIG			0x0004
#define FIELDID_FACT_CONFIG				0x0005
#define FIELDID_USER_CONFIG				0x0006
#define FIELDID_FACT_CALIB				0x0007
#define FIELDID_USER_CALIB				0x0008
#define FIELDID_BS_CERTS				0x0009
#define FIELDID_LWM2M_CERTS				0x000a



/*Module Capability*/
#define MOD_CAP_DEPENDENT	   			0x00  // Need power as well as instruction from the other module to bootup.
#define MOD_CAP_AUTONOMOUS 				0x01  //Just need power. It can boot up and do all things on it's own.

/* Module Mode */
#define MOD_MODE_SLAVE 					0x00  //Controlled by other.
#define MOD_MODE_MASTER 				0x01  //Controls other module too.

/*Module devices ownership.*/
#define MOD_DEV_LENDER				    0x00 //All sensors are controlled by other device.
#define MOD_DEV_OWNER			        0x01 //All sensors are controlled by module itself.

/* Validation */
#define VALIDATE_INDEX_COUNT(count)             ( ( (count >= 0) && (count < UKDB_IDX_MAX_TPL_COUNT) )?1:0 )
#define VALIDATE_MODULE_COUNT(count)            ( ( (count >= 0) && (count < MAX_NUMBER_MODULES_PER_UNIT) )?1:0 )
#define VALIDATE_DEVICE_COUNT(count)            ( ( (count >= 0) && (count < MAX_NUMBER_DEVICES_PER_MODULE) )?1:0 )


typedef enum {
	E_TNODESDR = 1,
	E_TNODELTE = 2,
	E_HNODE    = 3,
	E_ANODE	 = 4,
	E_PSNODE	 = 5
} UnitType;

typedef enum {
	E_COMV1 	= 0,
	E_SDR	   ,
	E_CNTRL	   ,
	E_GSMRF	   ,
	E_LTE	   ,
    E_LTERF	   ,
	E_MASK
} ModuleType;

typedef struct __attribute__((__packed__)){
       uint32_t magic_word;
       uint16_t resv1;
       uint16_t resv2;
} UKDBMagicWord;

typedef struct __attribute__((__packed__)){
	uint8_t major;
	uint8_t minor;
} Version;

typedef struct __attribute__((__packed__)){
	Version dbversion;
	uint16_t idx_tbl_offset;
    uint16_t idx_tpl_size;
	uint16_t idx_tpl_max_count;
	uint16_t idx_cur_tpl;
	uint8_t  mod_cap; // Self sustainable or not.
	uint8_t  mod_mode; // Like master or slave in unit.
	uint8_t  mod_devown; //Does module own the devices in it's schema or controlled by secondary module.
	uint8_t  resv1;
	uint16_t resv2;
	uint16_t resv3;
	uint16_t resv4;
	uint16_t resv5;
	uint16_t resv6;
} UKDBHeader;

typedef struct __attribute__((__packed__)){
	uint16_t fieldid;
	uint16_t payload_offset;
	uint16_t payload_size;
	Version  payload_version;
	uint32_t payload_crc;
	uint8_t	 state; // Enabled/disabled. Could based on the power switch or license or sw based like for knowing if fact config should be used or user config or HW capability reduced for low cost thing.
	bool  valid; // Mostly related to entries if marked deleted or error by user. or in simple words it tells data is usable.
	uint16_t resv1;
	uint16_t resv2;
	uint16_t resv3;
	uint16_t resv4;
	uint16_t resv5;
} UKDBIdxTuple;


/* TODO: TMP: As almost all the devices are I2C this is good start.*/
typedef struct  __attribute__((__packed__)) {
	uint8_t bus;
	uint16_t add;
} DeviceCfg;	

typedef struct __attribute__((__packed__)) {
	char dev_name[24]; //TODO: Check if this could be replaces by device object.
	char dev_disc[24];
	uint16_t dev_type;
	uint16_t dev_class;
	char sysfile[64];
	void* cfg; // TODO: Try union of the DevXXXXCfg
} ModuleCfg; //#124

typedef struct __attribute__((__packed__)) {
	char mod_uuid[24];
	char mod_name[24];
	char sysfs[64];
	void* eeprom_cfg;
} UnitCfg; //#120

typedef struct __attribute__((__packed__)){
	char uuid[24];
	char name[24];
	UnitType unit;
	char partno[24];
	char skew[24];
	char mac[18];
	Version swver;
	Version pswver;
	char assm_date[12];
	char oem_name[24];
	uint8_t mod_count;
} UnitInfo; //159

typedef struct __attribute__((__packed__)) {
	char uuid[24];
	char name[24];
	ModuleType module;
	char partno[24];
	char hwver[24];
	char mac[18];
	Version swver;
	Version pswver;
	char mfg_date[12];
	char mfg_name[24];
	uint8_t dev_count;
	ModuleCfg* module_cfg;
} ModuleInfo; //167

typedef struct __attribute__((__packed__)) {
	UKDBMagicWord magicword;
	UKDBHeader header;
	UKDBIdxTuple *indextable;
	UnitInfo unitinfo;
	UnitCfg* unitcfg;  //Contain list of modules lTE -1 /54
	ModuleInfo modinfo;
	ModuleCfg* modcfg; //Contains list of devices.
	void* factcfg;
	void* usercfg;
	void* factcalib;
	void* usercalib;
	void* bscerts;
	void* lwm2mcerts;
	void* resv1;
	void* resv2;
	void* resv3;
	void* resv4;
	void* resv5;
	void* resv6;
} UKDB;

typedef struct {
    char** fname;
    char* pname;
    uint8_t count;          /* Max 5 files to be allowed for now. Best to pass first json for master*/
} JSONInput;

#endif /*INCLUDERR_UBSP_LAYOUT_H_*/
