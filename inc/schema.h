/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_SCHEMA_H_
#define INC_SCHEMA_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "noded_macros.h"

#include "usys_types.h"

#define SCH_START_OFFSET               0x0000
#define SCH_END_OFFSET                 0xFFFF

#define SCH_MAGIC_WORD_OFFSET          0x0000
#define SCH_MAGIC_WORD_SIZE            0x0004
#define SCH_MAGIC_WORD                 0xDEADBEEF
#define SCH_DEFVAL                     0xFFFF

#define SCH_HEADER_OFFSET              0x0010
#define SCH_HEADER_DBVER_OFFSET        0x0010
#define SCH_HEADER_SIZE                0x0018
#define SCH_MAJOR_VER                  0
#define SCH_MINOR_VER                  0



#define SCH_IDX_TABLE_OFFSET           0x0040
#define SCH_IDX_TPL_SIZE               0x0018
#define SCH_IDX_MAX_TPL_COUNT          0x0032  /* Max Index table entries 50 */
#define SCH_IDX_CUR_TPL_COUNT_OFFSET   0x0020
#define SCH_IDX_COUNT_SIZE             0x0002


#define SCH_FOOTER_OFFSET              0x09A0
#define SCH_FOOTER_SIZE                0x0050

/* MAX_PAYLOAD_SIZE */
#define SCH_MAX_PAYLOAD_SIZE           0x1000

/* Size */
#define SCH_UNITINFO_PAYLOADSIZE       0x009F
#define SCH_UNITCONFIG_PAYLOADSIZE     0x0073
#define SCH_MODULEINFO_PAYLOADSIZE     0x00A7
#define SCH_MODULECONFIG_PAYLOADSIZE   0x0077
#define SCH_FACT_CONFIG_SIZE           0x1000
#define SCH_USER_CONFIG_SIZE           0x1000
#define SCH_FACT_CALIB_SIZE            0x1000
#define SCH_USER_CALIB_SIZE            0x1000
#define SCH_BS_CERTS_SIZE              0x1000
#define SCH_CLOUD_CERTS_SIZE           0x1000


/* OFFSETS*/
#define MAX_NUMBER_MODULES_PER_UNIT     0x0008
#define MAX_NUMBER_DEVICES_PER_MODULE   0x0014
#define SCH_PAYLAOD_OFFSET             0x0A00

#define SCH_UNIT_INFO_OFFSET           0x0A00 //Offset: 2560 (64*40)       size = 192B     REQ  = 159B     MAX = 1 Entry
#define SCH_UNIT_CONFIG_OFFSET         0x0AC0 //Offset: 2752 (64*43)       size = 1024B    REQ: = 115B     MAX = 8 Entry   NM 7
#define SCH_MODULE_INFO_OFFSET         0x0EC0 //Offset: 3776 (64*59)       size = 192B     REQ  = 167B     MAX = 1 Entry
#define SCH_MODULE_CONFIG_OFFSET       0x0F80 //Offset: 3968 (64*62)       size = 2432B    REQ  = 119B     MAX = 20 Entry  NM 15
#define SCH_FACT_CONFIG_OFFSET         0x1900 //Offset: 6400 (64*100)      size = 5120B    REQ  = 159B     MAX = 1 Entry
#define SCH_USER_CONFIG_OFFSET         0x2900 //Offset: 10496 (64*164)     size = 4096B    REQ  = 159B     MAX = 1 Entry
#define SCH_FACT_CALIB_OFFSET          0x3900 //Offset: 14592 (64*228)     size = 4096B    REQ  = 159B     MAX = 1 Entry
#define SCH_USER_CALIB_OFFSET          0x4900 //Offset: 18688 (64*292)     size = 4096B    REQ  = 159B     MAX = 1 Entry
#define SCH_BS_CERTS_OFFSET            0x5900 //Offset: 22784 (64*256)     size = 4096B    REQ  = 159B     MAX = 1 Entry
#define SCH_CLOUD_CERTS_OFFSET         0x6900 //Offset: 26880 (64*320)     size = 4096B    REQ  = 159B     MAX = 1 Entry

# if NEW_OFFSETS
#define MAX_NUMBER_MODULES_PER_UNIT     0x0008
#define MAX_NUMBER_DEVICES_PER_MODULE   0x0014
#define SCH_PAYLAOD_OFFSET             0x0A00

#define SCH_UNIT_INFO_OFFSET           0x0A00 //Offset: 2560 (64*40)       size = 192B     REQ/CFG  = 159B   (192)    MAX = 1 Entry
#define SCH_UNIT_CONFIG_OFFSET         0x0AC0 //Offset: 2752 (64*43)       size = 1024B    REQ/CFG =  115B   (146)    MAX = 7 Entry
#define SCH_MODULE_INFO_OFFSET         0x0EC0 //Offset: 3776 (64*59)       size = 192B     REQ/CFG  = 167B   (192)    MAX = 1 Entry
#define SCH_MODULE_CONFIG_OFFSET       0x0F80 //Offset: 3968 (64*62)       size = 2432B    REQ/CFG  = 119B   (152)    MAX = 16 Entry
#define SCH_FACT_CONFIG_OFFSET         0x1900 //Offset: 6400 (64*100)      size = 5120B    REQ/CFG  = 159B   (192)    MAX = 1 Entry
#define SCH_USER_CONFIG_OFFSET         0x2900 //Offset: 10496 (64*164)     size = 4096B    REQ/CFG  = 159B   (192)    MAX = 1 Entry
#define SCH_FACT_CALIB_OFFSET          0x3900 //Offset: 14592 (64*228)     size = 4096B    REQ/CFG  = 159B   (192)    MAX = 1 Entry
#define SCH_USER_CALIB_OFFSET          0x4900 //Offset: 18688 (64*292)     size = 4096B    REQ/CFG  = 159B   (192)    MAX = 1 Entry
#define SCH_BS_CERTS_OFFSET            0x5900 //Offset: 22784 (64*256)     size = 4096B    REQ/CFG  = 159B   (192)    MAX = 1 Entry
#define SCH_LWM2M_CERTS_OFFSET         0x6900 //Offset: 26880 (64*320)     size = 4096B    REQ/CFG  = 159B   (192)    MAX = 1 Entry

#endif
#if 0
#define SCH_PAYLAOD_OFFSET             0x0A00
#define SCH_UNIT_INFO_OFFSET           0x0A00 //Offset: 2560 (64*40)       size = 192B     REQ  = 159B     MAX = 1 Entry
#define SCH_UNIT_CONFIG_OFFSET         0x0AC0 //Offset: 2752 (64*43)       size = 1024B    REQ: = 115B     MAX = 8 Entry
#define SCH_MODULE_INFO_OFFSET         0x0EC0 //Offset: 3776 (64*59)       size = 192B     REQ  = 167B     MAX = 1 Entry
#define SCH_MODULE_CONFIG_OFFSET       0x0F80 //Offset: 3968 (64*62)       size = 2432B    REQ  = 119B     MAX = 20 Entry
#define SCH_FACT_CONFIG_OFFSET         0x1900 //Offset: 6400 (64*100)      size = 5120B    REQ  = 159B     MAX = 1 Entry
#define SCH_USER_CONFIG_OFFSET         0x2900 //Offset: 10496 (64*164)     size = 4096B    REQ  = 159B     MAX = 1 Entry
#define SCH_FACT_CALIB_OFFSET          0x3900 //Offset: 14592 (64*228)     size = 4096B    REQ  = 159B     MAX = 1 Entry
#define SCH_USER_CALIB_OFFSET          0x4900 //Offset: 18688 (64*292)     size = 4096B    REQ  = 159B     MAX = 1 Entry
#define SCH_BS_CERTS_OFFSET            0x5900 //Offset: 22784 (64*256)     size = 4096B    REQ  = 159B     MAX = 1 Entry
#define SCH_LWM2M_CERTS_OFFSET         0x6900 //Offset: 26880 (64*320)     size = 4096B    REQ  = 159B     MAX = 1 Entry
#endif
#define SCH_FEAT_ENABLED               0x01
#define SCH_FEAT_DISABLED              0x00




/* Field Id */
#define FIELD_ID_UNIT_INFO               0x0001
#define FIELD_ID_UNIT_CFG                0x0002
#define FIELD_ID_MODULE_INFO             0x0003
#define FIELD_ID_MODULE_CFG              0x0004
#define FIELD_ID_FACT_CFG                0x0005
#define FIELD_ID_USER_CFG                0x0006
#define FIELD_ID_FACT_CALIB              0x0007
#define FIELD_ID_USER_CALIB              0x0008
#define FIELD_ID_BS_CERTS                0x0009
#define FIELD_ID_CLOUD_CERTS             0x000a



/*Module Capability*/
#define MOD_CAP_DEPENDENT               0x00  // Need power as well as instruction from the other module to bootup.
#define MOD_CAP_AUTONOMOUS              0x01  //Just need power. It can boot up and do all things on it's own.

/* Module Mode */
#define MOD_MODE_SLAVE                  0x00  //Controlled by other.
#define MOD_MODE_MASTER                 0x01  //Controls other module too.

/*Module devices ownership.*/
#define MOD_DEV_LENDER                  0x00 //All sensors are controlled by other device.
#define MOD_DEV_OWNER                   0x01 //All sensors are controlled by module itself.

/* State */
#define IDX_ENTRY_DISABLED              0x00
#define IDX_ENTRY_ENABLED               0x01

/* Validation */
#define VALIDATE_INDEX_COUNT(count)             ( ( (count >= 0) && (count < SCH_IDX_MAX_TPL_COUNT) )?1:0 )
#define VALIDATE_MODULE_COUNT(count)            ( ( (count >= 0) && (count < MAX_NUMBER_MODULES_PER_UNIT) )?1:0 )
#define VALIDATE_DEVICE_COUNT(count)            ( ( (count >= 0) && (count < MAX_NUMBER_DEVICES_PER_MODULE) )?1:0 )

typedef enum {
    UNIT_TNODESDR = 1,
    UNIT_TNODELTE = 2,
    UNIT_HNODE    = 3,
    UNIT_ANODE    = 4,
    UNIT_PSNODE   = 5
} UnitType;

typedef enum {
    MOD_COMV1  = 0,
    MOD_SDR,
    MOD_CNTRL,
    MOD_RFFE,
    MOD_MASK
} ModuleType;

typedef struct __attribute__((__packed__, aligned(1))){
    uint32_t magicWord;
    uint16_t resv1;
    uint16_t resv2;
} SchemaMagicWord;

typedef struct __attribute__((__packed__, aligned(1))){
    uint8_t major;
    uint8_t minor;
} Version;

typedef struct __attribute__((__packed__, aligned(1))){
    Version  version;
    uint16_t idxTblOffset;
    uint16_t idxTplSize;
    uint16_t idxTplMaxCount;
    uint16_t idxCurTpl;
    uint8_t  modCap; // Self sustainable or not.
    uint8_t  modMode; // Like master or slave in unit.
    uint8_t  modDevOwn; //Does module own the devices in it's schema or controlled by secondary module.
    uint8_t  resv1;
    uint16_t resv2;
    uint16_t resv3;
    uint16_t resv4;
    uint16_t resv5;
    uint16_t resv6;
} SchemaHeader;

typedef struct __attribute__((__packed__, aligned(1))){
    uint16_t fieldId;
    uint16_t payloadOffset;
    uint16_t payloadSize;
    Version  payloadVer;
    uint32_t payloadCrc;
    uint8_t  state; // Enabled/disabled. Could based on the power switch or license or sw based like for knowing if fact config should be used or user config or HW capability reduced for low cost thing.
    bool  valid; // Mostly related to entries if marked deleted or error by user. or in simple words it tells data is usable.
    uint16_t resv1;
    uint16_t resv2;
    uint16_t resv3;
    uint16_t resv4;
    uint16_t resv5;
} SchemaIdxTuple;


/* TODO: TMP: As almost all the devices are I2C this is good start.*/
typedef struct __attribute__((__packed__, aligned(1))) {
    uint8_t bus;
    uint16_t add;
} DeviceCfg;

typedef struct __attribute__((__packed__, aligned(1))) {
    char devName[NAME_LENGTH]; //TODO: Check if this could be replaces by device object.
    char devDesc[DESC_LENGTH];
    uint16_t devType;
    uint16_t devClass;
    char sysFile[PATH_LENGTH];
    void* cfg; // TODO: Try union of the DevXXXXCfg
} ModuleCfg; //#124

typedef struct __attribute__((__packed__, aligned(1))) {
    char modUuid[UUID_LENGTH];
    char modName[NAME_LENGTH];
    char sysFs[64];
    void* eepromCfg;
} UnitCfg; //#120

typedef struct __attribute__((__packed__, aligned(1))){
    char uuid[UUID_LENGTH];
    char name[NAME_LENGTH];
    UnitType unit;
    char partNo[NAME_LENGTH];
    char skew[NAME_LENGTH];
    char mac[MAC_LENGTH];
    Version swVer;
    Version pswVer;
    char assmDate[DATE_LENGTH];
    char oemName[NAME_LENGTH];
    uint8_t modCount;
} UnitInfo; //159

typedef struct __attribute__((__packed__, aligned(1))){
    char uuid[UUID_LENGTH];
    char name[NAME_LENGTH];
    ModuleType module;
    char partNo[NAME_LENGTH];
    char hwVer[NAME_LENGTH];
    char mac[MAC_LENGTH];
    Version swVer;
    Version pSwVer;
    char mfgDate[DATE_LENGTH];
    char mfgName[NAME_LENGTH];
    uint8_t devCount;
    ModuleCfg* modCfg;
} ModuleInfo; //167

typedef struct __attribute__((__packed__, aligned(1))){
    SchemaMagicWord magicWord;
    SchemaHeader header;
    SchemaIdxTuple *indexTable;
    UnitInfo unitInfo;
    UnitCfg* unitCfg;  //Contain list of modules lTE -1 /54
    ModuleInfo modInfo;
    ModuleCfg* modCfg; //Contains list of devices.
    void* factCfg;
    void* userCfg;
    void* factCalib;
    void* userCalib;
    void* bsCerts;
    void* cloudCerts;
    void* resv1;
    void* resv2;
    void* resv3;
    void* resv4;
    void* resv5;
    void* resv6;
} StoreSchema;

typedef struct {
    char** fname;
    char* pname;
    uint8_t count;          /* Max 5 files to be allowed for now. Best to pass first json for master*/
} JSONInput;

#ifdef __cplusplus
}
#endif

#endif /* INC_SCHEMA_H_ */
