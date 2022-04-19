/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "test/test.h"

#include "headers/errorcode.h"
#include "headers/ubsp/devices.h"
#include "headers/ubsp/property.h"
#include "headers/ubsp/ubsp.h"
#include "inc/devicedb.h"
#include "inc/globalheader.h"
#include "inc/ukdb.h"
#include "utils/crc32.h"
#include "headers/utils/log.h"
#include "ukdb/db/db.h"
#include "ukdb/db/file.h"
#include "ukdb/idb/cs.h"
#include "ukdb/idb/idb.h"

static uint16_t pass_count = 0;
static uint16_t fail_count = 0;
static uint16_t total_count = 0;

typedef struct {
    char uuid[24]; /* Module UUID*/
    char sysdb[64]; /* System DB (Softlink to the EEPROM sysfs) */
    uint8_t mod_cap; /* Module can work independently or not in unit.*/
    uint8_t
        mod_mode; /* Mode of the module Working Master(Will have unit info and cfg)/
                           Slave(Will not have unit info and cfg)*/
} SystemdbMap;

#define TEST_TIME(iter)                                     \
    int ctr = iter + 1;                                     \
    while (iter) {                                          \
        sleep(1);                                           \
        log_debug("(%d) Sleeping for 5 sec", (ctr - iter)); \
        iter--;                                             \
    }

#define TEST_RET(str, ret, exp_ret)                                            \
    if ((ret) != (exp_ret)) {                                                  \
        total_count++;                                                         \
        fail_count++;                                                          \
        log_debug("TEST[%d]:: Test Case Name %s Status: FAILED Reason: %d.",   \
                  total_count, str, ret);                                      \
        sleep(1);                                                              \
        if (ret) { /*sleep(10) ; exit(0);*/                                    \
        }                                                                      \
    } else {                                                                   \
        total_count++;                                                         \
        pass_count++;                                                          \
        log_debug("TEST[%d]:: Test Case Name %s Status: Passed.", total_count, \
                  str, ret);                                                   \
    }

#define PROPERTYJSON "lib/ubsp/mfgdata/property/property.json"
typedef struct {
    const char *input;
    uint32_t crc32;
} CRCCheck;

Property *g_prop;

static CRCCheck crc_test_data[] = {
    { "123456789", 0xCBF43926ul },
    { "abcdefghijklmnopqrstuvwxyz", 0x4C2750BDul },
    { "VishalThakur@ukama", 0xCA18C73B },
    { "", 0x00000000ul },
    { " ", 0xE96CCF45ul },
    { "factcfgabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
      0x40BE4BBCul },
    { "usercfgabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
      0x75A4AF35ul },
    { "factcalibabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
      0xF76620D2ul },
    { "usercalibabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
      0xBA4010ABul },
    { "bscertsabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
      0x13EE1DEEul },
    { "lwm2mcertsabcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789",
      0xDCA83D9Aul },
    {},
    { NULL, 0 },
};

int test_crc() {
    int ret = 0;
    uint8_t iter = 0;
    uint16_t len = 0;
    uint32_t crc32 = 0;
    const unsigned char *ptr;
    log_debug("TEST:: Starting CRC test.");
    while (crc_test_data[iter].input != NULL) {
        ptr = (const unsigned char *)crc_test_data[iter].input;
        len = strlen(crc_test_data[iter].input);
        crc32 = crc_32(ptr, len);
        if (crc32 != crc_test_data[iter].crc32) {
            log_debug(
                "TEST::CRC32:: FAIL: CRC32 \"%s\" of length %d bytes returns 0x%08X not 0x%08X",
                crc_test_data[iter].input, len, crc32,
                crc_test_data[iter].crc32);
            ret = -1;
        } else {
            log_debug(
                "TEST::CRC32:: PASS: CRC32 \"%s\" of length %d bytes returns 0x%08X matches stored CRC 0x%08X",
                crc_test_data[iter].input, len, crc32,
                crc_test_data[iter].crc32);
            ret++;
        }
        iter++;
    }

    if (iter != ret) {
        ret = -1;
    } else {
        ret = 0;
    }
    return ret;
}

int test_ukdb_header(char *p_uuid) {
    int ret = 0;
    UKDBHeader *header = malloc(sizeof(UKDBHeader));
    if (header) {
        ret = ubsp_read_header(p_uuid, header);
        if (!ret) {
            log_debug("TEST:: UKDB Read UKDB Header :: PASS");

        } else {
            log_debug("TEST:: UKDB Read UKDB Header :: FAIL");
        }
        UBSP_FREE(header);
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("TEST:: UKDB failed in memory allocation for header read.");
    }
    return ret;
}

/* Free the memory used by moduleCfg */
void test_free_module_cfg_dev(ModuleCfg *cfg, uint8_t count) {
    for (int iter = 0; iter < count; iter++) {
        void *dev_cfg = cfg[iter].cfg;
        UBSP_FREE(dev_cfg);
    }
    UBSP_FREE(cfg);
}

int test_ukdb_read_module(char *p_uuid) {
    int ret = 0;
    uint16_t size = 0;
    uint8_t count = 0;
    ModuleCfg *mcfg;
    ModuleInfo *minfo = ubsp_alloc(sizeof(ModuleInfo));
    if (minfo) {
        ret = ubsp_read_module_info(p_uuid, minfo, &size);
        if (!ret) {
            ukdb_print_module_info(minfo);
            log_debug("TEST:: UKDB Read UKDB Module Info for %s :: PASS",
                      p_uuid);
            count = minfo->dev_count;
            UBSP_FREE(minfo);

            /* Read Module Cfg */
            mcfg = ubsp_alloc_module_cfg(count);
            if (mcfg) {
                size = 0;
                ret = ubsp_read_module_cfg(p_uuid, mcfg, count, &size);
                if (!ret) {
                    ukdb_print_module_cfg(mcfg, count);
                    log_debug(
                        "TEST:: UKDB Read UKDB Module Config for %s :: PASS",
                        p_uuid);
                } else {
                    log_debug(
                        "TEST:: UKDB Read UKDB Module Config for %s :: FAIL",
                        p_uuid);
                    goto cleanmcfg;
                }
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_debug(
                    "TEST:: UKDB failed in memory allocation for module cfg.");
                log_debug("TEST:: UKDB Read UKDB Module Config for %s :: FAIL",
                          p_uuid);
                goto cleanminfo;
            }
        } else {
            log_debug("TEST:: UKDB Read UKDB Module Info for %s :: FAIL",
                      p_uuid);
            goto cleanminfo;
        }

    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("TEST:: UKDB failed in memory allocation for module info.");
        log_debug("TEST:: UKDB Read UKDB Module Info for %s :: FAIL", p_uuid);
        goto cleanminfo;
    }

cleanmcfg:
    ubsp_free_module_cfg(mcfg, count);

cleanminfo:
    ubsp_free(minfo);
    return ret;
}

/* May be write a free API too */
void test_free_unit_cfg(UnitCfg *cfg, uint8_t count) {
    for (int iter = 0; iter < count; iter++) {
        UBSP_FREE(cfg[iter].eeprom_cfg);
    }
    UBSP_FREE(cfg);
}

int test_ukdb_read_unit(char *p_uuid) {
    int ret = 0;
    uint16_t size = 0;
    UnitCfg *ucfg;
    uint8_t mod_count = 0;
    UnitInfo *uinfo = ubsp_alloc(sizeof(UnitInfo));
    if (uinfo) {
        ret = ubsp_read_unit_info(p_uuid, uinfo, &size);
        if (!ret) {
            log_debug("TEST:: UKDB Read UKDB Unit Info for %s :: PASS", p_uuid);
            ukdb_print_unit_info(uinfo);
            mod_count = uinfo->mod_count;
            UBSP_FREE(uinfo);
            ucfg = ubsp_alloc_unit_cfg(mod_count);
            if (ucfg) {
                size = 0;
                ret = ubsp_read_unit_cfg(p_uuid, ucfg, mod_count, &size);
                if (!ret) {
                    log_debug(
                        "TEST:: UKDB Read UKDB Unit Config for %s :: PASS",
                        p_uuid);
                    ukdb_print_unit_cfg(ucfg, mod_count);
                } else {
                    log_debug("TEST:: UKDB Read UKDB Unit Configfor %s :: FAIL",
                              p_uuid);
                    goto cleanunitcfg;
                }
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_debug("TEST:: UKDB failed in memory allocation.");
                log_debug("TEST:: UKDB Read UKDB Unit Config for %s :: FAIL",
                          p_uuid);
                goto cleanunitinfo;
            }
        } else {
            log_debug("TEST:: UKDB Read UKDB Unit Info for %s :: FAIL", p_uuid);
            goto cleanunitinfo;
        }
    }

cleanunitcfg:
    ubsp_free_unit_cfg(ucfg, mod_count);

cleanunitinfo:
    ubsp_free(uinfo);

    return ret;
}

void test_schema_creation(char *p_uuid) {
    int ret = 0;
    uint16_t size = 0;
    UnitCfg *ucfg = NULL;
    ModuleCfg *mcfg = NULL;
    ModuleInfo *minfo = NULL;
    char *schema = NULL;
    uint8_t mod_count =
        1; /* Unit info is only present for master module in other case module count set to 1 by default */
    UnitInfo *uinfo = ubsp_alloc(sizeof(UnitInfo));
    if (uinfo) {
        ret = ubsp_read_unit_info(p_uuid, uinfo, &size);
        if (!ret) {
            log_debug("TEST:: UKDB Read UKDB Unit Info for %s :: PASS", p_uuid);
            ukdb_print_unit_info(uinfo);
            mod_count = uinfo->mod_count;
            ucfg = ubsp_alloc_unit_cfg(mod_count);
            if (ucfg) {
                size = 0;
                ret = ubsp_read_unit_cfg(p_uuid, ucfg, mod_count, &size);
                if (!ret) {
                    log_debug(
                        "TEST:: UKDB Read UKDB Unit Config for %s :: PASS",
                        p_uuid);
                    ukdb_print_unit_cfg(ucfg, mod_count);
                } else {
                    log_debug("TEST:: UKDB Read UKDB Unit Configfor %s :: FAIL",
                              p_uuid);
                    goto cleanunitcfg;
                }
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
                log_debug("TEST:: UKDB failed in memory allocation.");
                log_debug("TEST:: UKDB Read UKDB Unit Config for %s :: FAIL",
                          p_uuid);
                goto cleanunitinfo;
            }
        } else {
            log_debug("TEST:: UKDB Read UKDB Unit Info for %s :: FAIL", p_uuid);
            ubsp_free(uinfo);
            uinfo = NULL;
        }
    }

    uint8_t count = 0;
    minfo = ubsp_alloc(sizeof(ModuleInfo) * mod_count);
    if (minfo) {
        for (int iter = 0; iter < mod_count; iter++) {
            char *uuid = NULL;
            if (ucfg) {
                uuid = ucfg[iter].mod_uuid;
            } else {
                uuid = p_uuid;
            }
            ret = ubsp_read_module_info(uuid, &minfo[iter], &size);
            if (!ret) {
                ukdb_print_module_info(&minfo[iter]);
                log_debug("TEST:: UKDB Read UKDB Module Info for %s :: PASS",
                          uuid);
                count = minfo[iter].dev_count;

                /* Read Module Cfg */
                mcfg = ubsp_alloc_module_cfg(count);
                if (mcfg) {
                    size = 0;
                    ret = ubsp_read_module_cfg(uuid, mcfg, count, &size);
                    if (!ret) {
                        ukdb_print_module_cfg(mcfg, count);
                        log_debug(
                            "TEST:: UKDB Read UKDB Module Config for %s :: PASS",
                            uuid);
                    } else {
                        log_debug(
                            "TEST:: UKDB Read UKDB Module Config for %s :: FAIL",
                            uuid);
                        goto cleanmcfg;
                    }
                    minfo[iter].module_cfg = mcfg;
                } else {
                    ret = ERR_UBSP_MEMORY_EXHAUSTED;
                    log_debug(
                        "TEST:: UKDB failed in memory allocation for module cfg.");
                    log_debug(
                        "TEST:: UKDB Read UKDB Module Config for %s :: FAIL",
                        uuid);
                    goto cleanminfo;
                }
            } else {
                log_debug("TEST:: UKDB Read UKDB Module Info for %s :: FAIL",
                          uuid);
                goto cleanminfo;
            }
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("TEST:: UKDB failed in memory allocation for module info.");
        log_debug("TEST:: UKDB Read UKDB Module Info for %s :: FAIL", p_uuid);
        goto cleanminfo;
    }

    /* Create schema */
    ret = ubsp_create_schema(uinfo, ucfg, minfo, &schema);
    TEST_RET("ubsp_create_schema", ret, 0);

cleanmcfg:
    for (int iter = 0; iter < mod_count; iter++) {
        ubsp_free_module_cfg(minfo[iter].module_cfg, minfo[iter].dev_count);
    }
cleanminfo:
    ubsp_free(minfo);
cleanunitcfg:
    ubsp_free_unit_cfg(ucfg, mod_count);
cleanunitinfo:
    ubsp_free(uinfo);
    UBSP_FREE(schema);
}

void print_test_data(char *info, uint16_t size) {
    uint16_t itr = 0;
    while (itr < size) {
        putchar(*(info + itr));
        itr++;
    }
    fflush(stdout);
}
int test_ukdb_read_fact_config(char *puuid) {
    int ret = 0;
    uint16_t size = 0;
    //TODO: May be implement a API to get size.*/
    char *data = malloc(UKDB_MAX_PAYLOAD_SIZE);
    if (data) {
        ret = ubsp_read_fact_config(puuid, data, &size);
        if (!ret) {
            log_debug("TEST:: UKDB Read UKDB fact config for %s :: PASS",
                      puuid);
            log_debug("TEST:: Fact Config ::\n");
            print_test_data(data, size);
            UBSP_FREE(data);
        } else {
            log_debug("TEST:: UKDB Read UKDB fact config for %s :: FAIL",
                      puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("TEST:: UKDB failed in memory allocation.");
        log_debug("TEST:: UKDB Read UKDB fact config for %s :: FAIL", puuid);
    }
    return ret;
}

int test_ukdb_read_user_config(char *puuid) {
    int ret = 0;
    uint16_t size = 0;
    char *data = malloc(UKDB_MAX_PAYLOAD_SIZE);
    if (data) {
        ret = ubsp_read_user_config(puuid, data, &size);
        if (!ret) {
            log_debug("TEST:: UKDB Read UKDB user config for %s :: PASS",
                      puuid);
            log_debug("TEST:: User Config  ::\n");
            print_test_data(data, size);
            UBSP_FREE(data);
        } else {
            log_debug("TEST: :UKDB Read UKDB user config for %s :: FAIL",
                      puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("TEST:: UKDB failed in memory allocation.");
        log_debug("TEST:: UKDB Read UKDB user config for %s :: FAIL", puuid);
    }
    return ret;
}

int test_ukdb_read_fact_calib(char *puuid) {
    int ret = 0;
    uint16_t size = 0;
    char *data = malloc(UKDB_MAX_PAYLOAD_SIZE);
    if (data) {
        ret = ubsp_read_fact_calib(puuid, data, &size);
        if (!ret) {
            log_debug("TEST:: UKDB Read UKDB fact calib for %s :: PASS", puuid);
            log_debug("TEST:: Fact Calibration  ::\n");
            print_test_data(data, size);
            UBSP_FREE(data);
        } else {
            log_debug("TEST:: UKDB Read UKDB fact calib for %s :: FAIL", puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("TEST:: UKDB failed in memory allocation.");
        log_debug("TEST:: UKDB Read UKDB fact calib for %s :: FAIL", puuid);
    }
    return ret;
}

int test_ukdb_read_user_calib(char *puuid) {
    int ret = 0;
    uint16_t size = 0;
    char *data = malloc(UKDB_MAX_PAYLOAD_SIZE);
    if (data) {
        ret = ubsp_read_user_calib(puuid, data, &size);
        if (!ret) {
            log_debug("TEST:: UKDB Read UKDB user calib for %s :: PASS", puuid);
            log_debug("TEST:: User Calibration ::\n");
            print_test_data(data, size);
            UBSP_FREE(data);
        } else {
            log_debug("TEST:: UKDB Read UKDB user calib for %s :: FAIL", puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("TEST:: UKDB failed in memory allocation.");
        log_debug("TEST:: UKDB Read UKDB user calib for %s :: FAIL", puuid);
    }
    return ret;
}

int test_ukdb_read_bs_certs(char *puuid) {
    int ret = 0;
    uint16_t size = 0;
    char *data = malloc(UKDB_MAX_PAYLOAD_SIZE);
    if (data) {
        ret = ubsp_read_bs_certs(puuid, data, &size);
        if (!ret) {
            log_debug("TEST:: UKDB Read UKDB bs certs for %s :: PASS", puuid);
            log_debug("TEST:: Boot Strap certs  ::\n");
            print_test_data(data, size);
            UBSP_FREE(data);
        } else {
            log_debug("TEST:: UKDB Read UKDB bs certs for %s :: FAIL", puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("TEST:: UKDB failed in memory allocation.");
        log_debug("TEST:: UKDB Read UKDB bs certs for %s :: FAIL", puuid);
    }
    return ret;
}

int test_ukdb_read_lwm2m_certs(char *puuid) {
    int ret = 0;
    uint16_t size = 0;
    char *data = malloc(UKDB_MAX_PAYLOAD_SIZE);
    if (data) {
        ret = ubsp_read_lwm2m_certs(puuid, data, &size);
        if (!ret) {
            log_debug("TEST:: UKDB Read UKDB lwm2m certs for %s :: PASS",
                      puuid);
            log_debug("TEST:: Lwm2m certs  ::\n");
            print_test_data(data, size);
            UBSP_FREE(data);
        } else {
            log_debug("TEST:: UKDB Read UKDB lwm2m certs for %s :: FAIL",
                      puuid);
        }
    } else {
        ret = ERR_UBSP_MEMORY_EXHAUSTED;
        log_debug("TEST:: UKDB failed in memory allocation.");
        log_debug("TEST:: UKDB Read UKDB lwm2m certs for %s :: FAIL", puuid);
    }
    return ret;
}

void test_master_module() {
    int ret = 0;
    char puuid[1] = "\0";
    log_debug(
        "****************************************************************");
    log_debug(
        "*********** Starting UKDB testing for Master Module ************");
    log_debug(
        "****************************************************************");
    ret = test_ukdb_header(puuid);
    TEST_RET("test_ukdb_header", ret, 0);
    log_debug(
        "****************************************************************");

    ret = test_ukdb_read_unit(puuid);
    TEST_RET("test_ukdb_read_unit", ret, 0);
    log_debug(
        "****************************************************************");

    ret = test_ukdb_read_module(puuid);
    TEST_RET("test_ukdb_read_module", ret, 0);
    log_debug(
        "****************************************************************");
    ret = test_ukdb_read_fact_config(puuid);
    TEST_RET("test_ukdb_read_fact_config", ret, 0);
    log_debug(
        "****************************************************************");
    ret = test_ukdb_read_user_config(puuid);
    TEST_RET("test_ukdb_read_user_config", ret, 0);
    log_debug(
        "****************************************************************");
    ret = test_ukdb_read_fact_calib(puuid);
    TEST_RET("test_ukdb_read_fact_calib", ret, 0);
    log_debug(
        "****************************************************************");
    ret = test_ukdb_read_user_calib(puuid);
    TEST_RET("test_ukdb_read_user_calib", ret, 0);
    log_debug(
        "****************************************************************");
    ret = test_ukdb_read_bs_certs(puuid);
    TEST_RET("test_ukdb_read_bs_certs", ret, 0);
    log_debug(
        "***************************************************************");
    ret = test_ukdb_read_lwm2m_certs(puuid);
    TEST_RET("test_ukdb_read_lwm2m_certs", ret, 0);
    log_debug(
        "***************************************************************");

    log_debug(
        "***************** UKDB testing for Module %s over *************",
	puuid);
    log_debug(
        "***************************************************************");
}

void test_ukdb(SystemdbMap *sysdbmap) {
    char *puuid = sysdbmap->uuid;
    uint8_t mode = sysdbmap->mod_mode;
    uint8_t cap = sysdbmap->mod_cap;

    int ret = 0;
    log_debug(
        "********************************************************************************");
    log_debug(
        "************* Starting UKDB testing for Module %s ***************",
        puuid);
    log_debug(
        "********************************************************************************");
    ret = test_ukdb_header(puuid);
    TEST_RET("test_ukdb_header", ret, 0);
    log_debug(
        "********************************************************************************");

    /* Only For master modules.*/
    if (mode) {
        ret = test_ukdb_read_unit(puuid);
        TEST_RET("test_ukdb_read_unit", ret, 0);
        log_debug(
            "********************************************************************************");
    }

    ret = test_ukdb_read_module(puuid);
    TEST_RET("test_ukdb_read_module", ret, 0);
    log_debug(
        "********************************************************************************");

    /*Only for modules  who can work independently of master */
    if (cap) {
        ret = test_ukdb_read_fact_config(puuid);
        TEST_RET("test_ukdb_read_fact_config", ret, 0);
        log_debug(
            "********************************************************************************");
        ret = test_ukdb_read_user_config(puuid);
        TEST_RET("test_ukdb_read_user_config", ret, 0);
        log_debug(
            "********************************************************************************");
        ret = test_ukdb_read_fact_calib(puuid);
        TEST_RET("test_ukdb_read_fact_calib", ret, 0);
        log_debug(
            "********************************************************************************");
        ret = test_ukdb_read_user_calib(puuid);
        TEST_RET("test_ukdb_read_user_calib", ret, 0);
        log_debug(
            "********************************************************************************");
        ret = test_ukdb_read_bs_certs(puuid);
        TEST_RET("test_ukdb_read_bs_certs", ret, 0);
        log_debug(
            "********************************************************************************");
        ret = test_ukdb_read_lwm2m_certs(puuid);
        TEST_RET("test_ukdb_read_lwm2m_certs", ret, 0);
        log_debug(
            "********************************************************************************");
    }
    log_debug(
        "***************** UKDB testing for Module %s over ********************",
        puuid);
    log_debug(
        "********************************************************************************");
}

uint16_t test_read_device_prop_count(DevObj *obj) {
    uint16_t count = 0;
    int ret = 0;
    ret = upsb_read_dev_prop_count(obj, &count);
    TEST_RET("upsb_read_dev_prop_count", ret, 0);
    if (count <= 1) {
        count = 0;
    }
    return count;
}

/* Test device properties */
int test_read_device_props(DevObj *obj, Property *prop, uint16_t count) {
    int ret = 0;
    ret = ubsp_read_dev_props(obj, prop);
    if (!ret) {
        print_properties(prop, count);
        log_debug(
            "TEST:: Property read for Device name %s Disc: %s Module UUID %s : PASS.",
            obj->name, obj->disc, obj->mod_UUID);
    } else {
        log_debug(
            "TEST:: Property read for Device name %s Disc: %s Module UUID %s : FAIL.",
            obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int test_read_dev_props(DevObj *obj) {
    int ret = -1;
    uint16_t count = test_read_device_prop_count(obj);
    if (count > 0) {
        Property *prop = malloc(sizeof(Property) * count);
        if (prop) {
            ret = test_read_device_props(obj, prop, count);
            TEST_RET("test_read_device_props", ret, 0);
        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_debug("TEST::Failed in memory allocation for Property read.");
        }
        UBSP_FREE(prop);
    }
    return ret;
}

/* TODO find the right location*/
void print_registered_devices(Device *dev, uint16_t count) {
    log_debug(
        "********************************************************************************");
    log_debug(
        "************************ Registered Devices ************************************");
    for (uint16_t iter = 0; iter < count; iter++) {
        log_debug(
            "********************************************************************************");
        log_debug("* Name                      : %s", dev[iter].obj.name);
        log_debug("* Disc               	   : %s", dev[iter].obj.disc);
        log_debug("* Module UUID               : %s", dev[iter].obj.mod_UUID);
        log_debug("* Type                      : %d", dev[iter].obj.type);
        log_debug("* SysFile Name:             : %s", dev[iter].sysfile);
        log_debug(
            "********************************************************************************");
    }
}

int test_read_registered_devices(DeviceType type) {
    int ret = 0;
    uint16_t count = 0;
    log_debug("TEST:: Read registered devices for device type 0x%x.", type);
    ret = ubsp_read_registered_dev_count(type, &count);
    if (ret) {
        return ret;
    }
    if (count > 0) {
        Device *dev = malloc(sizeof(Device) * count);
        if (dev) {
            ret = ubsp_read_registered_dev(type, dev);
            if (!ret) {
                print_registered_devices(dev, count);
                log_debug(
                    "TEST:: Read registered devices for device type 0x%x : PASS",
                    type);
            } else {
                log_debug(
                    "TEST:: Read registered devices for device type 0x%x : FAIL",
                    type);
            }
            UBSP_FREE(dev);
        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
            log_debug(
                "TEST::Failed in memory allocation for registered devices read.");
        }
    }

    return ret;
}

int test_read_from_dev_prop(DevObj *obj, int *prop, void *value) {
    int ret = 0;
    ret = ubsp_read_from_prop(obj, prop, value);
    if (!ret) {
        log_debug(
            "TEST:: Property[%d] read for Device name %s Disc: %s Module UUID %s : PASS.",
            *prop, obj->name, obj->disc, obj->mod_UUID);
    } else {
        log_debug(
            "TEST:: Property[%d] read for Device name %s Disc: %s Module UUID %s : FAIL.",
            *prop, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int test_write_to_dev_prop(DevObj *obj, int *prop, void *value) {
    int ret = 0;
    ret = ubsp_write_to_prop(obj, prop, value);
    if (!ret) {
        log_debug(
            "TEST:: Property[%d] write for Device name %s Disc: %s Module UUID %s : PASS.",
            *prop, obj->name, obj->disc, obj->mod_UUID);
    } else {
        log_debug(
            "TEST:: Property[%d] write for Device name %s Disc: %s Module UUID %s : FAIL.",
            *prop, obj->name, obj->disc, obj->mod_UUID);
    }
    return ret;
}

int test_read_write_to_dev_props(SystemdbMap *sysdbmap, uint8_t count) {
    int ret = 0;
    DevObj devobj[24] = {
        { .name = "TMP464",
          .disc = "Pmic",
          .mod_UUID = "UK-1001-COM-1101",
          .type = DEV_TYPE_TMP },
        { .name = "SE98",
          .disc = "X86",
          .mod_UUID = "UK-1001-COM-1101",
          .type = DEV_TYPE_TMP },
        { .name = "SE98",
          .disc = "ADI",
          .mod_UUID = "UK-2001-LTE-1101",
          .type = DEV_TYPE_TMP },
        { .name = "ADT7481",
          .disc = "Wifi Controller",
          .mod_UUID = "UK-3001-MSK-1101",
          .type = DEV_TYPE_TMP },
        { .name = "INA226",
          .disc = "Pmic",
          .mod_UUID = "UK-1001-COM-1101",
          .type = DEV_TYPE_PWR },
        { .name = "INA226",
          .disc = "DDR",
          .mod_UUID = "UK-1001-COM-1101",
          .type = DEV_TYPE_PWR },
        { .name = "INA226",
          .disc = "DDR",
          .mod_UUID = "UK-2001-LTE-1101",
          .type = DEV_TYPE_PWR },
        { .name = "INA226",
          .disc = "CNX",
          .mod_UUID = "UK-2001-LTE-1101",
          .type = DEV_TYPE_PWR },
        { .name = "INA226",
          .disc = "PCI",
          .mod_UUID = "UK-3001-MSK-1101",
          .type = DEV_TYPE_PWR },
        /* RF UP Board*/
        { .name = "SE98",
          .disc = "RF MicroProcessor",
          .mod_UUID = "UK-5001-RFC-1101",
          .type = DEV_TYPE_TMP },
        { .name = "LED-TRICOLOR",
          .disc = "RF LED 0",
          .mod_UUID = "UK-5001-RFC-1101",
          .type = DEV_TYPE_LED },
        { .name = "LED-TRICOLOR",
          .disc = "RF LED 1",
          .mod_UUID = "UK-5001-RFC-1101",
          .type = DEV_TYPE_LED },
        { .name = "LED-TRICOLOR",
          .disc = "RF LED 2",
          .mod_UUID = "UK-5001-RFC-1101",
          .type = DEV_TYPE_LED },
        { .name = "LED-TRICOLOR",
          .disc = "RF LED 3",
          .mod_UUID = "UK-5001-RFC-1101",
          .type = DEV_TYPE_LED },
        /*RF-FE*/
        { .name = "DAT-31R5A-PP",
          .disc = "rx-att",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_ATT },
        { .name = "DAT-31R5A-PP",
          .disc = "tx-att",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_ATT },
        { .name = "ADS1015",
          .disc = "rf-power detector",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_ADC },
        { .name = "TMP464",
          .disc = "RFFE Board",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_TMP },
        { .name = "GPIO",
          .disc = "PGOOD 5V",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_GPIO },
        { .name = "GPIO",
          .disc = "PGOOD 3.3V",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_GPIO },
        { .name = "GPIO",
          .disc = "PGOOD 5.7V",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_GPIO },
        { .name = "GPIO",
          .disc = "PA DISABLE",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_GPIO },
        { .name = "GPIO",
          .disc = "PGA DISABLE",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_GPIO },
        { .name = "GPIO",
          .disc = "RF POWER DISABLE",
          .mod_UUID = "UK-4001-RFA-1101",
          .type = DEV_TYPE_GPIO },

    };

    for (int mapiter = 0; mapiter < count; mapiter++) {
        for (uint8_t dev_c = 0; dev_c < 23; dev_c++) {
            if (!strcmp(devobj[dev_c].mod_UUID, sysdbmap[mapiter].uuid)) {
                log_debug(
                    "*********************************************************");
                log_debug(
                    "Starting Read/Write Property testing for Device Name %s "
                    "Disc: %s Module UUID %s",
                    devobj[dev_c].name, devobj[dev_c].disc,
                    devobj[dev_c].mod_UUID);
                log_debug(
                    "**********************************************************");
                uint16_t count = test_read_device_prop_count(&devobj[dev_c]);
                Property *prop = malloc(sizeof(Property) * count);
                if (prop) {
                    ret = test_read_device_props(&devobj[dev_c], prop, count);
                    TEST_RET("test_read_device_props", ret, 0);
                } else {
                    ret = ERR_UBSP_MEMORY_EXHAUSTED;
                    log_debug(
                        "TEST::Failed in memory allocation for Property read.");
                }

                for (int iter = 0; iter < count; iter++) {
                    if (prop[iter].available == PROP_NOTAVAIL) {
                        continue;
                    }
                    /* Data based on property type */
                    int64_t orgvalue = 0;
                    void *rvalue = 0;
                    void *wvalue = 0;

                    double wdvalue = 65000;
                    double rdvalue = 0;
                    int32_t wivalue = 55000;
                    int32_t rivalue = 0;
                    uint32_t wuivalue = 55000;
                    uint32_t ruivalue = 0;
                    bool wbvalue = true;
                    bool rbvalue = false;
                    uint8_t wussivalue = 1;
                    uint8_t russivalue = 0;
                    int8_t wssivalue = 1;
                    int8_t rssivalue = 0;
                    uint16_t wusivalue = 200;
                    uint16_t rusivalue = 0;
                    int16_t wsivalue = 200;
                    int16_t rsivalue = 0;
                    char wstrvalue[64] = "GREEN";
                    char rstrvalue[64] = { '\0' };

                    switch (prop[iter].data_type) {
                    case TYPE_DOUBLE: {
                        wvalue = &wdvalue;
                        rvalue = &rdvalue;
                        break;
                    }
                    case TYPE_BOOL: {
                        wvalue = &wbvalue;
                        rvalue = &rbvalue;
                        break;
                    }
                    case TYPE_UINT8:
                    {
                        wvalue = &wussivalue;
                        rvalue = &russivalue;
                        break;
                    }
                    case TYPE_INT8: {
                        wvalue = &wssivalue;
                        rvalue = &rssivalue;
                        break;
                    }
                    case TYPE_UINT16:
                        wvalue = &wusivalue;
                        wvalue = &rusivalue;
                        break;
                    case TYPE_INT16: {
                        wvalue = &wsivalue;
                        wvalue = &rsivalue;
                        break;
                    }
                    case TYPE_UINT32: {
                        wvalue = &wuivalue;
                        rvalue = &ruivalue;
                        break;
                    }
                    case TYPE_INT32: {
                        wvalue = &wivalue;
                        rvalue = &rivalue;
                        break;
                    }
                    case TYPE_STRING: {
                        wvalue = wstrvalue;
                        rvalue = rstrvalue;
                        break;
                    }
                    default: {
                        wvalue = &wivalue;
                        rvalue = &rivalue;
                    }
                    }
                    ret = test_read_from_dev_prop(&devobj[dev_c], &iter,
                                                  &orgvalue);
                    TEST_RET("test_read_from_dev_prop", ret, 0);
                    if (ret) {
                        log_debug(
                            "TEST:: Error %d while reading device property[%d] %s. Moving to next.",
                            ret, iter, prop[iter].name);
                        continue;
                    }
                    /* If property has Write Permission */
                    if (prop[iter].perm & PERM_WR) {
                        ret = test_write_to_dev_prop(&devobj[dev_c], &iter,
                                                     wvalue);
                        TEST_RET("test_write_to_dev_prop", ret, 0);
                        ret = test_read_from_dev_prop(&devobj[dev_c], &iter,
                                                      rvalue);
                        TEST_RET("test_read_from_dev_prop", ret, 0);
                        ret = test_write_to_dev_prop(&devobj[dev_c], &iter,
                                                     &orgvalue);
                        TEST_RET("test_write_to_dev_prop", ret, 0);
                        /* In case of bool read and write value will be either 1 or 0 */
                        if (prop[iter].data_type != TYPE_BOOL) {
                            if (memcmp(wvalue, rvalue,
                                       get_sizeof(prop[iter].data_type))) {
                                ret = -1;
                                TEST_RET("test_write_to_dev_prop", ret, 0);
                                break;
                            } else {
                                ret = 0;
                            }
                        }
                        log_debug(
                            "TEST:: Read/write test to a device property[%d] %s.",
                            iter, prop[iter].name);
                        TEST_RET("test_read_write_to_dev_props", ret, 0);
                    }
                }
                UBSP_FREE(prop);
                log_debug(
                    "**********************************************************");
            }
        }
    }
    return ret;
}

int test_enable_alert(DevObj *obj, int *prop) {
    int ret = 0;
    log_debug(
        "TEST:: Property[%d] enable alert for Device name %s Disc: %s Module UUID %s.",
        *prop, obj->name, obj->disc, obj->mod_UUID);
    ret = ubsp_enable_irq(obj, prop);
    return ret;
}

int test_disable_alert(DevObj *obj, int *prop) {
    int ret = 0;
    log_debug(
        "TEST:: Property[%d] disable alert for Device name %s Disc: %s Module UUID %s : PASS.",
        *prop, obj->name, obj->disc, obj->mod_UUID);
    ret = ubsp_disable_irq(obj, prop);
    return ret;
}

/* This is called from the ISR context should be released as soon as possible to start monitor again.
 *  We could have a thread waiting in the app. Once we get this callback release a flag/semaphore and
 *  pass the required info for processing. Here I have just used ISR context to print info which is for demo only.
 */
void test_app_cb(DevObj *obj, AlertCallBackData **acbdata, int *count) {
    if (*acbdata) {
        AlertCallBackData *adata = *acbdata;
        log_debug(
            "TEST:: App Callback function for Device name %s Disc: %s Module UUID %s called "
            "from ubsp with %d alerts.",
            obj->name, obj->disc, obj->mod_UUID, *count);
        int pidx = adata->pidx;
        int didx = g_prop[pidx].dep_prop->curr_idx;
        uint8_t alertstate = adata->alertstate;
        /* Considered double but need to be read based on type in g_prop[pidx].data_type
		 * int size = get_sizeof(prop[dep_idx].data_type); */
        double value = *(double *)adata->svalue;

        log_debug(
            "TEST:: Alert %d received for Property[%d], Name: %s , Value %lf %s.",
            alertstate, pidx, g_prop[pidx].name, value, g_prop[didx].units);
        UBSP_FREE(adata->svalue);
        UBSP_FREE(adata);
    }
}

int test_register_app_cb(DevObj *obj, int *prop, CallBackFxn fn) {
    int ret = 0;
    log_debug(
        "TEST:: Registering Callback function for Device name %s Disc: %s Module UUID %s.",
        obj->name, obj->disc, obj->mod_UUID);
    ret = ubsp_register_app_cb(obj, prop, fn);
    return ret;
}

int test_deregister_app_cb(DevObj *obj, int *prop, CallBackFxn fn) {
    int ret = 0;
    log_debug(
        "TEST:: De-Registering Callback function for Device name %s Disc: %s Module UUID %s.",
        obj->name, obj->disc, obj->mod_UUID);
    ret = ubsp_deregister_app_cb(obj, prop, fn);
    return ret;
}

int test_enable_disable_alerts(SystemdbMap *sysdbmap, uint8_t count) {
    int ret = 0;
    int rec = 1; /*Timer for loop*/
    DevObj devobj[10] = { { .name = "TMP464",
                            .disc = "Pmic",
                            .mod_UUID = "UK-1001-COM-1101",
                            .type = DEV_TYPE_TMP },
                          { .name = "SE98",
                            .disc = "ADI",
                            .mod_UUID = "UK-2001-LTE-1101",
                            .type = DEV_TYPE_TMP },
                          { .name = "TMP464",
                            .disc = "RFFE Board",
                            .mod_UUID = "UK-4001-RFA-1101",
                            .type = DEV_TYPE_TMP },
                          { .name = "SE98",
                            .disc = "RF MicroProcessor",
                            .mod_UUID = "UK-5001-RFC-1101",
                            .type = DEV_TYPE_TMP },
                          { .name = "ADT7481",
                            .disc = "Wifi Controller",
                            .mod_UUID = "UK-3001-MSK-1101",
                            .type = DEV_TYPE_TMP },
                          { .name = "INA226",
                            .disc = "Pmic",
                            .mod_UUID = "UK-1001-COM-1101",
                            .type = DEV_TYPE_PWR },
                          { .name = "INA226",
                            .disc = "DDR",
                            .mod_UUID = "UK-1001-COM-1101",
                            .type = DEV_TYPE_PWR },
                          { .name = "INA226",
                            .disc = "DDR",
                            .mod_UUID = "UK-2001-LTE-1101",
                            .type = DEV_TYPE_PWR },
                          { .name = "INA226",
                            .disc = "CNX",
                            .mod_UUID = "UK-2001-LTE-1101",
                            .type = DEV_TYPE_PWR },
                          { .name = "INA226",
                            .disc = "PCI",
                            .mod_UUID = "UK-3001-MSK-1101",
                            .type = DEV_TYPE_PWR } };

    for (int mapiter = 0; mapiter < count; mapiter++) {
        for (uint8_t dev_c = 0; dev_c < 8; dev_c++) {
            if (!strcmp(devobj[dev_c].mod_UUID, sysdbmap[mapiter].uuid)) {
                log_debug(
                    "*********************************************************");
                log_debug("Starting Alert testing for Device Name %s "
                          "Disc: %s Module UUID %s",
                          devobj[dev_c].name, devobj[dev_c].disc,
                          devobj[dev_c].mod_UUID);
                log_debug(
                    "**********************************************************");
                uint16_t count = test_read_device_prop_count(&devobj[dev_c]);

                g_prop = malloc(sizeof(Property) * count);
                if (g_prop) {
                    ret = test_read_device_props(&devobj[dev_c], g_prop, count);
                    TEST_RET("test_read_device_props", ret, 0);
                } else {
                    ret = ERR_UBSP_MEMORY_EXHAUSTED;
                    log_debug(
                        "TEST::Failed in memory allocation for Property read.");
                }

                for (int iter = 0; iter < count; iter++) {
                    /* Enable alert */
                    if (g_prop[iter].prop_type & PROP_TYPE_ALERT) {
                        /* Enable alert
			            TODO: Remove iter as it per device not per property.*/
                        ret = test_enable_alert(&devobj[dev_c], &iter);
                        TEST_RET("test_enable_alert", ret, 0);

                        /* Register Callback function
                         * TODO: Remove iter as it per device not per property.
                         * */
                        ret = test_register_app_cb(&devobj[dev_c], &iter,
                                                   &test_app_cb);
                        TEST_RET("test_register_app_cb", ret, 0);

                        TEST_TIME(rec);

                        /* De-register Callback function
			            TODO: Remove iter as it per device not per property.*/
                        ret = test_deregister_app_cb(&devobj[dev_c], &iter,
                                                     &test_app_cb);
                        TEST_RET("test_deregister_app_cb", ret, 0);

                        /* Disable alert
                         * TODO: Remove iter as it per device not per property.*/
                        ret = test_disable_alert(&devobj[dev_c], &iter);
                        TEST_RET("test_disable_alert", ret, 0);

                        rec = 1;
                    }
                }
                UBSP_FREE(g_prop);
                log_debug(
                    "**********************************************************");
            }
        }
    }
    return ret;
}

/* Test device db*/
int test_devicedb(SystemdbMap *sysdbmap, uint8_t count) {
    int ret = 0;
    log_debug(
        "********************************************************************************");
    log_debug(
        "******************** Starting DeviceDB testing *********************************");
    log_debug(
        "********************************************************************************");
    DevObj dev_obj[17] = { { .name = "TMP464",
                             .disc = "Pmic",
                             .mod_UUID = "UK-1001-COM-1101",
                             .type = DEV_TYPE_TMP },
                           { .name = "TMP464",
                             .disc = "DDR",
                             .mod_UUID = "UK-1001-COM-1101",
                             .type = DEV_TYPE_TMP },
                           { .name = "TMP464",
                             .disc = "PMIC",
                             .mod_UUID = "UK-2001-LTE-1101",
                             .type = DEV_TYPE_TMP },
                           { .name = "INA226",
                             .disc = "DDR",
                             .mod_UUID = "UK-1001-COM-1101",
                             .type = DEV_TYPE_PWR },
                           { .name = "SE98",
                             .disc = "X86",
                             .mod_UUID = "UK1001-LTE",
                             .type = DEV_TYPE_TMP },
                           { .name = "INA226",
                             .disc = "CNX",
                             .mod_UUID = "UK-2001-LTE-1101",
                             .type = DEV_TYPE_PWR },
                           { .name = "ADT7481",
                             .disc = "Wifi Controller",
                             .mod_UUID = "UK-3001-MSK-1101",
                             .type = DEV_TYPE_TMP },
                           { .name = "INA226",
                             .disc = "PCI",
                             .mod_UUID = "UK-3001-MSK-1101",
                             .type = DEV_TYPE_PWR },
                           { .name = "DAT-31R5A-PP",
                             .disc = "rx-att",
                             .mod_UUID = "UK-4001-RFA-1101",
                             .type = DEV_TYPE_ATT },
                           { .name = "DAT-31R5A-PP",
                             .disc = "tx-att",
                             .mod_UUID = "UK-4001-RFA-1101",
                             .type = DEV_TYPE_ATT },
                           { .name = "ADS1015",
                             .disc = "rf-power detector",
                             .mod_UUID = "UK-4001-RFA-1101",
                             .type = DEV_TYPE_ADC },
                           { .name = "GPIO",
                             .disc = "PGOOD 5V",
                             .mod_UUID = "UK-4001-RFA-1101",
                             .type = DEV_TYPE_GPIO },
                           { .name = "GPIO",
                             .disc = "PA DISABLE",
                             .mod_UUID = "UK-4001-RFA-1101",
                             .type = DEV_TYPE_GPIO },
                           { .name = "GPIO",
                             .disc = "RF POWER DISABLE",
                             .mod_UUID = "UK-4001-RFA-1101",
                             .type = DEV_TYPE_GPIO },
                           { .name = "SE98",
                             .disc = "RF MicroProcessor",
                             .mod_UUID = "UK-5001-RFC-1101",
                             .type = DEV_TYPE_TMP },
                           { .name = "LED-TRICOLOR",
                             .disc = "RF LED 0",
                             .mod_UUID = "UK-5001-RFC-1101",
                             .type = DEV_TYPE_LED },
                           { .name = "LED-TRICOLOR",
                             .disc = "RF LED 3",
                             .mod_UUID = "UK-5001-RFC-1101",
                             .type = DEV_TYPE_LED } };

    /* Read registered devices. */
    DeviceType type = DEV_TYPE_TMP;
    for (int iter = DEV_TYPE_TMP; iter <= DEV_TYPE_ATT; iter++) {
        ret = test_read_registered_devices(iter);
        TEST_RET("test_read_registered_devices", ret, 0);
    }

    for (int mapiter = 0; mapiter < count; mapiter++) {
        for (int iter = 0; iter < 10; iter++) {
            if (!strcmp(dev_obj[iter].mod_UUID, sysdbmap[mapiter].uuid)) {
                ret = test_read_dev_props(&dev_obj[iter]);
                TEST_RET("test_read_dev_props", ret, 0);
            }
        }
    }

    ret = test_read_write_to_dev_props(sysdbmap, count);
    TEST_RET("test_read_write_to_dev_props", ret, 0);

    ret = test_enable_disable_alerts(sysdbmap, count);
    TEST_RET("test_enable_disable_alerts", ret, 0);

    log_debug(
        "**************************** DeviceDB testing over *****************************");
    return ret;
}

/* Creates a db from json schema.*/
int test_ukdb_with_db_from_json_schema(JSONInput *ip, SystemdbMap *sysdbmap) {
    int ret = 0;
    int count = ip->count;
    ret = ubsp_devdb_init(ip->pname);
    TEST_RET("ubsp_devdb_init", ret, 0);

    ret = ubsp_idb_init(ip);
    TEST_RET("ubsp_idb_init", ret, 0);

    /* If DB exist it will prepare environment for reading./writing to DB.*/
    ret = ubsp_ukdb_init(sysdbmap[0].sysdb);
    TEST_RET("ubsp_ukdb_init", ret, (ERR_UBSP_DB_MISSING));
    if (ret) {
        /* When no database is present ukdb_init would just return with error.
		 * Preparing a environment for creating DB. This step is just for testing purpose .
		 * On Target we never have to go through this step.*/
        ret = ubsp_pre_create_ukdb_hook(sysdbmap[0].uuid);
        TEST_RET("ubsp_pre_create_ukdb_hook", ret, 0);
    }

    if (!ret) {
        //TODO: Enable only for Virtual Node.
        for (uint8_t iter = 0; iter < count; iter++) {
            log_debug(
                "********************************************************************************");
            log_debug(
                "******************* Creating UKDB for Module Id: %s *****************************",
                sysdbmap[iter].uuid);
            log_debug(
                "********************************************************************************");
            ret = ubsp_create_ukdb(sysdbmap[iter].uuid);
            TEST_RET("ubsp_create_ukdb", ret, 0);

            test_ukdb(&sysdbmap[iter]);

            log_debug(
                "********************************************************************************");
            if (ret) {
                break;
            }
        }
    }
    ubsp_idb_exit();
    ubsp_exit();
    return ret;
}

/* Writing a newly parsed json to DB (sysfs/eeprom) */
int test_create_db_hook(char *puuid, JSONInput *ip) {
    int ret = 0;
    UnitCfg *udata = (UnitCfg[]){
        { .mod_uuid = "UK-5001-RFC-1101",
          .mod_name = "RF CTRL BOARD",
          .sysfs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0051/eeprom",
          .eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
        { .mod_uuid = "UK-4001-RFA-1101",
          .mod_name = "RF BOARD",
          .sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0052/eeprom",
          .eeprom_cfg = &(DevI2cCfg){ .bus = 2, .add = 0x50ul } },
        { .mod_uuid = "UK-1001-COM-1101",
          .mod_name = "COMv1",
          .sysfs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0050/eeprom",
          .eeprom_cfg = &(DevI2cCfg){ .bus = 0, .add = 0x50ul } },
        { .mod_uuid = "UK-2001-LTE-1101",
          .mod_name = "LTE",
          .sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0050/eeprom",
          .eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
        { .mod_uuid = "UK-3001-MSK-1101",
          .mod_name = "MASK",
          .sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0051/eeprom",
          .eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x51ul } },
    };

    UnitCfg *pcfg = NULL;
    DevI2cCfg *i2c_cfg = NULL;
    for (int iter = 0; iter < 5; iter++) {
        if (!strcmp(puuid, udata[iter].mod_uuid)) {
            pcfg = malloc(sizeof(UnitCfg));
            if (pcfg) {
                memset(pcfg, '\0', sizeof(UnitCfg));
                memcpy(pcfg, &udata[iter], sizeof(UnitCfg));
                if (udata[iter].eeprom_cfg) {
                    i2c_cfg = malloc(sizeof(DevI2cCfg));
                    if (i2c_cfg) {
                        memset(i2c_cfg, '\0', sizeof(DevI2cCfg));
                        memcpy(i2c_cfg, udata[iter].eeprom_cfg,
                               sizeof(DevI2cCfg));
                    }
                }
                pcfg->eeprom_cfg = i2c_cfg;
                break;
            } else {
                log_error("Err(%d): TEST: Memory exhausted while getting unit"
                          " config from Test data.",
                          ERR_UBSP_MEMORY_EXHAUSTED);
            }
        }
    }

    ret = ubsp_devdb_init(ip->pname);
    TEST_RET("ubsp_devdb_init", ret, 0);
    if (ret) {
        goto cleanup;
    }
    
    ret = ubsp_idb_init(&ip);
    TEST_RET("ubsp_idb_init", ret, 0);
    if (ret) {
        goto cleanup;
    }

    ret =
        ubsp_ukdb_init(NULL); /* Will just initialize the db if NULL is passed*/

    /* register Module */
    ret = ubsp_register_module(pcfg);
    TEST_RET("ubsp_register_module", ret, 0);
    if (!ret) {
        ret = ubsp_create_ukdb(pcfg->mod_uuid);
        TEST_RET("ubsp_create_ukdb", ret, 0);
    }

cleanup:
    ubsp_idb_exit();
    ubsp_exit();
    UBSP_FREE(pcfg->eeprom_cfg);
    UBSP_FREE(pcfg);

    return ret;
}

/* Creates a db from c- structs.*/
int test_ukdb_with_db_creation_from_c_structs(JSONInput *ip,
                                              SystemdbMap *sysdbmap) {
    int ret = 0;
    int count = ip->count;
    ip = NULL; /* This forces to use c -structs.*/

    if (ip) {
        ret = ubsp_devdb_init(ip->pname);
    } else {
        ret = ubsp_devdb_init(NULL);
    }
    TEST_RET("ubsp_devdb_init", ret, 0);

    ret = ubsp_idb_init(ip);
    TEST_RET("ubsp_idb_init", ret, 0);

    /* Id DB exist it will prepare environment for reading./writing to DB.*/
    ret = ubsp_ukdb_init(sysdbmap[0].sysdb);
    TEST_RET("ubsp_ukdb_init", ret, (ERR_UBSP_DB_MISSING));
    if (ret) {
        /* When no database is present ukdb_init would just return with error.
         * Preparing a environment for creating DB. This step is just for testing purpose .
         * On Target we never have to go through this step.*/
        ret = ubsp_pre_create_ukdb_hook(sysdbmap[0].uuid);
        TEST_RET("ubsp_pre_create_ukdb_hook", ret, 0);
    }

    if (!ret) {
        for (uint8_t iter = 0; iter < count; iter++) {
            log_debug(
                "********************************************************************************");
            log_debug(
                "******************* Creating UKDB for Module Id: %s *****************************",
                sysdbmap[iter].uuid);
            log_debug(
                "********************************************************************************");

            ret = ubsp_create_ukdb(sysdbmap[iter].uuid);
            TEST_RET("ubsp_create_ukdb", ret, 0);

            test_ukdb(&sysdbmap[iter]);

            log_debug(
                "********************************************************************************");
            if (ret) {
                break;
            }
        }
    }
    /* Cleans the mfg data structs from memory Use this after ubsb_create*/
    ubsp_idb_exit();

    /* Remove DB*/
    for (uint8_t iter = 0; iter < count; iter++) {
        ret = ubsp_remove_ukdb(sysdbmap[iter].uuid);
        TEST_RET("ubsp_remove_ukdb", ret, 0);
    }

    ubsp_exit();
    return ret;
}

/* Test Unit/module info and cfg from existing db */
int test_ukdb_with_existing_db(SystemdbMap *sysdbmap, uint8_t count) {
    int ret = 0;
    ret = ubsp_devdb_init(PROPERTYJSON);
    TEST_RET("ubsp_devdb_init", ret, 0);

    /* We have database available at this point. ukdb will register modules to db and devices to devdb
     * This would only happen for master module which has unit info for others it will just return with Master info missing. */

    ret = ubsp_ukdb_init(sysdbmap->sysdb);
    TEST_RET("ubsp_ukdb_init", ret, 0);

    for (int iter = 0; iter < count; iter++) {
        log_debug(
            "********************************************************************************");
        log_debug(
            "****************** Starting UKDB testing From Existing DB for %s ***************",
            sysdbmap->uuid);
        log_debug(
            "********************************************************************************");
        /* Read module data with UUID */
        test_ukdb(&sysdbmap[iter]);

        /* Read module data without module UUID only done for master module*/
        if (sysdbmap[iter].mod_mode) {
            test_master_module();
        }

        /* Create JSON schema from EEPROM data */
        test_schema_creation(sysdbmap[iter].uuid);

        log_debug(
            "********************************************************************************");
    }

    /* Testing DeviceDB */
    ret = test_devicedb(sysdbmap, count);
    TEST_RET("test_devicedb", ret, 0);

    return ret;
}

/* Negative test cases for UKDB*/
int test_ukdb_with_existing_db_neg(JSONInput *ip, SystemdbMap *sysdbmap) {
    int ret = 0;
    int count = ip->count;
    ret = ubsp_devdb_init(PROPERTYJSON);
    TEST_RET("ubsp_devdb_init", ret, 0);
    for (uint8_t iter = 0; iter < count; iter++) {
        log_debug(
            "********************************************************************************");
        log_debug(
            "****************** Starting UKDB -ve test cases for Existing DB for %s ***************",
            sysdbmap[iter].uuid);
        log_debug(
            "********************************************************************************");
        /* No data base is available at this point. Link are set by prepare_env.sh for master module thats it.*/
        ret = ubsp_ukdb_init(sysdbmap[iter].sysdb);
        if (iter == 0) {
            /* DB Link exist but db doesn't exist*/
            TEST_RET("ubsp_ukdb_init", ret, (ERR_UBSP_DB_MISSING));
        } else {
            /* DB Link doen't exist for slave modules*/
            TEST_RET("ubsp_ukdb_init", ret, (ERR_UBSP_DB_LNK_MISSING));
        }
        log_debug(
            "********************************************************************************");
    }
    return ret;
}

int self_test() {
    int ret = 0;
    ret = test_crc();
    TEST_RET("test_crc", ret, 0);
    return ret;
}

void test_cases_call(JSONInput *ip, SystemdbMap *sysdbmap) {
    int ret = 0;

    test_ukdb_with_existing_db_neg(ip, sysdbmap);

    ret = test_ukdb_with_db_creation_from_c_structs(ip, sysdbmap);
    TEST_RET("test_ukdb_with_db_creation_from_c_structs", ret, 0);

    log_debug("Total Test Cases: %d Passed: %d Failed %d.", total_count,
              pass_count, fail_count);

    ret = test_ukdb_with_db_from_json_schema(ip, sysdbmap);
    TEST_RET("test_ukdb_with_db_from_json_schema", ret, 0);

    log_debug("Total Test Cases: %d Passed: %d Failed %d.", total_count,
              pass_count, fail_count);

    ret = test_ukdb_with_existing_db(sysdbmap, ip->count);
    TEST_RET("test_ukdb_with_existing_db", ret, 0);

    log_debug("Total Test Cases: %d Passed: %d Failed %d.", total_count,
              pass_count, fail_count);
}

void test_anode() {
    JSONInput ip = { .fname = (char *[]){ "mfgdata/schema/rfctrl.json",
                                          "mfgdata/schema/rffe.json" },
                     .pname = PROPERTYJSON,
                     .count = 2 };

    SystemdbMap sysdbmap[2] = { { "UK-5001-RFC-1101", "/tmp/sys/anode-systemdb",
                                  MOD_CAP_AUTONOMOUS, MOD_MODE_MASTER },
                                { "UK-4001-RFA-1101", "", MOD_CAP_DEPENDENT,
                                  MOD_MODE_SLAVE } };

    log_debug(
        "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$");
    log_debug(
        "$$$$$$$$$$$$$$$$$$ Started Testing for ANode $$$$$$$$$$$$$$$$$$");
    log_debug(
        "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$");

    test_cases_call(&ip, sysdbmap);

    log_debug(
        "###############################################################");

    log_debug("Total Test Cases: %d Passed: %d Failed %d.", total_count,
              pass_count, fail_count);

    log_debug(
        "###############################################################");

    log_debug(
        "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$");
    log_debug(
        "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$");
}

void test_tnode() {
    JSONInput ip = { .fname =
                         (char *[]){
                             "mfgdata/schema/com.json",
                             "mfgdata/schema/mask.json",
                             "mfgdata/schema/lte.json",
                         },
                     .pname = PROPERTYJSON,
                     .count = 3 };

    SystemdbMap sysdbmap[3] = {
        { "UK-1001-COM-1101", "/tmp/sys/cnode-systemdb", MOD_CAP_AUTONOMOUS,
          MOD_MODE_MASTER },
        { "UK-3001-MSK-1101", "", MOD_CAP_DEPENDENT, MOD_MODE_SLAVE },
        { "UK-2001-LTE-1101", "", MOD_CAP_AUTONOMOUS, MOD_MODE_SLAVE }
    };

    log_debug(
        "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$");
    log_debug(
        "$$$$$$$$$$$$$$$$$$ Started Testing for TNode $$$$$$$$$$$$$$$$$$");
    log_debug(
        "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$");

    test_cases_call(&ip, sysdbmap);

    log_debug(
        "###############################################################");

    log_debug("Total Test Cases: %d Passed: %d Failed %d.", total_count,
              pass_count, fail_count);

    log_debug(
        "###############################################################");

    log_debug(
        "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$");
    log_debug(
        "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$");
}

int main(int argc, char **argv) {
    int ret = 0;

    ret = self_test();
    TEST_RET("self_test", ret, 0);

    /* Testing anode*/
    test_anode();

    /* Testing tnode*/
    test_tnode();

     JSONInput ip = { .fname = (char *[]){ "mfgdata/schema/rfctrl.json",
                                          "mfgdata/schema/rffe.json" },
                     .pname = PROPERTYJSON,
                     .count = 2 };

    /* testing creation of db from new schema */
    test_create_db_hook("UK-5001-RFC-1101", &ip);
    TEST_RET("test_create_db_hook", ret, 0);

    /* testing creation of db from new schema */
    test_create_db_hook("UK-4001-RFA-1101", &ip);
    TEST_RET("test_create_db_hook", ret, 0);

    log_debug(
        "################################################################################");
    log_debug("Total Test Cases: %d Passed: %d Failed %d.", total_count,
              pass_count, fail_count);
    log_debug(
        "################################################################################");

    return 0;
}
