/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/*
 * Tool for creating EEPROM database for device.
 * This binary is only meant to create a dummy device database by reading a JSON configs.
 * This tool would take JSON files, node name and  serial number as command line arguments.
 *
 * Example usage is:
 * mfgutil --n ComV1 --m UK-8001-COM-1102 --s mfgdata/schema/com.json -n LTE --m UK-8001-LTE-1102 --s mfgdata/schema/lte.json --n MASK -m UK-8001-MSK-1102 --s mfgdata/schema/mask.json
 * mfgutil -n RF_CTRL -u UK-5001-RFC-1101 -s mfgdata/schema/rfctrl.json -n RF_AMP -u UK-4001-RFA-1101 -s mfgdata/schema/rffe.json
 */

/*
 * - If FEMD_SYSROOT is set, pcfg->sysFs is rewritten to point into that tree.
 * - Also updates your hardcoded udata mapping to match the new /tmp/sys mock you showed:
 *      ctrl  -> i2c-0/0-0051/eeprom
 *      fe1   -> i2c-1/1-0050/eeprom
 *      fe2   -> i2c-2/2-0050/eeprom
 *
 * Important:
 * - You must pass -n fe1 and -n fe2 (NOT just fe) to create both databases.
 * - If FEMD_SYSROOT is not set, your existing absolute paths still work.
 */

#include "device.h"
#include "errorcode.h"
#include "inventory.h"
#include "ledger.h"
#include "noded_macros.h"
#include "property.h"

#include "utils/crc32.h"
#include "utils/mfg_helper.h"

#include "usys_api.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

#define PROPERTY_JSON "mfgdata/property/property.json"
#define VERSION "0.0.1"

#define ENV_FEMD_SYSROOT "FEMD_SYSROOT"

static const char *femd_sysroot(void) {
    const char *v = getenv(ENV_FEMD_SYSROOT);
    return (v && v[0] != '\0') ? v : NULL;
}

static bool starts_with(const char *s, const char *pfx) {
    if (!s || !pfx) return false;
    size_t n = strlen(pfx);
    return strncmp(s, pfx, n) == 0;
}

/* Resolve a sysfs-ish path into an absolute path under FEMD_SYSROOT. */
static int resolve_with_sysroot(char *out, size_t outsz, const char *path) {
    const char *root = femd_sysroot();
    if (!out || outsz == 0 || !path || path[0] == '\0') return -1;

    if (!root) {
        int n = snprintf(out, outsz, "%s", path);
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }

    if (starts_with(path, root)) {
        int n = snprintf(out, outsz, "%s", path);
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }

    if (starts_with(path, "/sys/")) {
        int n = snprintf(out, outsz, "%s%s", root, path + 4); /* drop "/sys" */
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }

    if (path[0] == '/') {
        int n = snprintf(out, outsz, "%s%s", root, path);
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }

    /* relative: leave */
    int n = snprintf(out, outsz, "%s", path);
    return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
}

static int safe_copy_sysfs_path(NodeCfg *pcfg, const char *absPath) {
    if (!pcfg || !absPath) return -1;

    /* sysFs is an array in NodeCfg, so we must copy into it safely */
    size_t cap = sizeof(pcfg->sysFs);
    if (cap == 0) return -1;

    size_t n = strlen(absPath);
    if (n >= cap) {
        log_error("MFGUTIL:: sysFs path too long (%zu >= %zu): %s", n, cap, absPath);
        return -1;
    }

    usys_memset(pcfg->sysFs, 0, cap);
    usys_memcpy(pcfg->sysFs, absPath, n);
    pcfg->sysFs[n] = '\0';
    return 0;
}

/* Writing a newly parsed json to DB (sysFs/eeprom) */
int create_db_hook(char **puuid, char** name, char** schema, int count) {
    int ret = 0;

    JSONInput jip;
    jip.fname = schema;
    jip.count = count;
    jip.pname = PROPERTY_JSON;

    NodeCfg *pcfg = NULL;
    DevI2cCfg *i2cCfg = NULL;

    ret = ldgr_init(jip.pname);
    if (ret) {
        log_error("MFGUTIL:: ledger initialization failed %d", ret);
        goto cleanup;
    }

    ret = invt_mfg_init(&jip);
    if (ret) {
        log_error("MFGUTIL:: UBSP IDB init failed %d", ret);
        goto cleanup;
    }

    ret = invt_init(NULL, &ldgr_register);
    if (ret) {
        usys_log_warn("MFGUTIL:: Inventory init failed %d (Expected -1)", ret);
    }

    for(int idx = 0; idx < count; idx++) {

        log_debug("UUID[%d] = %24s Name[%d] = %24s Schema file[%d] = %s \n",
                  idx, puuid[idx], idx, name[idx], idx, jip.fname[idx]);

        /*
         * Updated mapping for your mock sysroot tree:
         *   ctrl : /bus/i2c/devices/i2c-0/0-0051/eeprom
         *   fe1  : /bus/i2c/devices/i2c-1/1-0050/eeprom
         *   fe2  : /bus/i2c/devices/i2c-2/2-0050/eeprom
         *
         * NOTE: sysFs paths are sysroot-relative.
         */
        NodeCfg *udata = (NodeCfg[]){
            {   .modUuid = "ukma-8001-ctrl-1102",
                .modName = "ctrl",
                .sysFs   = "/bus/i2c/devices/i2c-0/0-0051/eeprom",
                .eepromCfg = &(DevI2cCfg){ .bus = 0, .add = 0x51ul }
            },
            {   .modUuid = "ukma-8001-fe1-1102",
                .modName = "fe1",
                .sysFs   = "/bus/i2c/devices/i2c-1/1-0050/eeprom",
                .eepromCfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul }
            },
            {   .modUuid = "ukma-8001-fe2-1102",
                .modName = "fe2",
                .sysFs   = "/bus/i2c/devices/i2c-2/2-0050/eeprom",
                .eepromCfg = &(DevI2cCfg){ .bus = 2, .add = 0x50ul }
            },
        };

        for (int iter = 0; iter < MAX_BOARDS; iter++) {
            if (!udata[iter].modName) continue;

            if (!usys_strcmp(name[idx], udata[iter].modName)) {

                pcfg = usys_zmalloc(sizeof(NodeCfg));
                if (!pcfg) {
                    log_error("MFGUTIL:: Err(%d): Memory exhausted while getting node config.",
                              ERR_NODED_MEMORY_EXHAUSTED);
                    goto cleanup;
                }

                usys_memcpy(pcfg, &udata[iter], sizeof(NodeCfg));

                if (udata[iter].eepromCfg) {
                    i2cCfg = usys_zmalloc(sizeof(DevI2cCfg));
                    if (!i2cCfg) {
                        log_error("MFGUTIL:: Err(%d): Memory exhausted while copying eepromCfg.",
                                  ERR_NODED_MEMORY_EXHAUSTED);
                        goto cleanup;
                    }
                    usys_memcpy(i2cCfg, udata[iter].eepromCfg, sizeof(DevI2cCfg));
                }
                pcfg->eepromCfg = i2cCfg;

                /* Update Module UUID */
                usys_memset(pcfg->modUuid, 0, 32);
                usys_memcpy(pcfg->modUuid, puuid[idx], usys_strlen(puuid[idx]));

                /* Resolve sysFs into FEMD_SYSROOT (or fallback /tmp/sys) */
                {
                    char absPath[512] = {0};

                    if (!femd_sysroot()) {
                        /* preserve old behavior: write into /tmp/sys */
                        int n = snprintf(absPath, sizeof(absPath), "%s%s", "/tmp/sys", udata[iter].sysFs);
                        if (n < 0 || (size_t)n >= sizeof(absPath)) {
                            log_error("MFGUTIL:: sysFs path truncated for %s", udata[iter].modName);
                            goto cleanup;
                        }
                    } else {
                        if (resolve_with_sysroot(absPath, sizeof(absPath), udata[iter].sysFs) != 0) {
                            log_error("MFGUTIL:: Failed to resolve sysFs for %s: %s",
                                      udata[iter].modName, udata[iter].sysFs);
                            goto cleanup;
                        }
                    }

                    if (safe_copy_sysfs_path(pcfg, absPath) != 0) {
                        goto cleanup;
                    }
                }

                break;
            }
        }

        if (!pcfg) {
            log_debug("No module with name %s found.", name[idx]);
            continue;
        }

        ret = invt_register_module(pcfg);
        if (ret) {
            log_error("MFGUTIL:: Registering module failed %d", ret);
            goto cleanup;
        }

        ret = invt_create_db(pcfg->modUuid);
        if (ret) {
            log_error("MFGUTIL:: Failed while creating inventory Database "
                      "for module %s UUID %s.", name[idx], pcfg->modUuid);
            goto cleanup;
        }

        usys_log_info("MFGUTIL:: Created inventory Database for module %s UUID %s at %s",
                      name[idx], pcfg->modUuid, pcfg->sysFs);

        /* cleanup per-module allocations */
        if (pcfg) {
            usys_free(pcfg->eepromCfg);
            pcfg->eepromCfg = NULL;
            usys_free(pcfg);
            pcfg = NULL;
            i2cCfg = NULL;
        }
    }

cleanup:
    if (pcfg) {
        usys_free(pcfg->eepromCfg);
        pcfg->eepromCfg = NULL;
        usys_free(pcfg);
        pcfg = NULL;
    }
    invt_mfg_exit();
    ldgr_exit();
    invt_exit();
    return ret;
}

static struct option longOptions[] = {
                { "name",    required_argument, 0, 'n' },
                { "muuid",   required_argument, 0, 'm' },
                { "file",    required_argument, 0, 'f' },
                { "logs",    required_argument, 0, 'l' },
                { "help",    no_argument, 0, 'h' },
                { "version", no_argument, 0, 'v' },
                { 0, 0, 0, 0 }
};

void set_log_level(char *slevel) {
    int ilevel = LOG_TRACE;
    if (!usys_strcmp(slevel, "TRACE")) {
        ilevel = LOG_TRACE;
    } else if (!usys_strcmp(slevel, "DEBUG")) {
        ilevel = LOG_DEBUG;
    } else if (!usys_strcmp(slevel, "INFO")) {
        ilevel = LOG_INFO;
    }
    log_set_level(ilevel);
}

void usage() {
    printf("Usage: genInventory [options] \n");
    printf("Options:\n");
    printf("--h, --help                                                      Help menu.\n");
    printf("--l, --logs <TRACE> <DEBUG> <INFO>                               Log level for the process.\n");
    printf("--n, --name <...>                                                Name of module.\n");
    printf("--m, --muuid <Module UUID>                                       Module UUID.\n");
    printf("--f, --file <json file>                                          JSON Schema file.\n");
    printf("--v, --version                                                   Software Version.\n");
}

int main(int argc, char** argv) {
    char *name[MAX_BOARDS] = {"\0"};
    char *uuid[MAX_BOARDS] = {"\0"};
    char *schema[MAX_BOARDS] = {"\0"};
    char *debug = "TRACE";

    set_log_level(debug);

    if (argc < 2 ) {
        log_error("Not enough arguments.");
        usage();
        usys_exit(1);
    }

    int uidx = 0;
    int nidx = 0;
    int sidx = 0;

    while (true) {
        int opt = 0;
        int opIdx = 0;

        opt = usys_getopt_long(argc, argv, "h:v:m:n:f:l:", longOptions, &opIdx);
        if (opt == -1) break;

        switch (opt) {
            case 'h':
                usage();
                usys_exit(0);
                break;
            case 'v':
                usys_puts(VERSION);
                break;
            case 'n':
                name[nidx++] = optarg;
                break;
            case 'm':
                uuid[uidx++] = optarg;
                break;
            case 'f':
                schema[sidx++] = optarg;
                break;
            case 'l':
                debug = optarg;
                set_log_level(debug);
                break;
            default:
                usage();
                exit(0);
        }
    }

    if ((sidx != uidx) || (sidx != nidx) || (!sidx) || (sidx > MAX_BOARDS)) {
        log_error("MFGUTIL:: Name, schema and UUID entries have to match in count.");
        log_error("MFGUTIL:: At least one set of entries or %d set of entries can be made simultaneously.", MAX_BOARDS);
        exit(0);
    }

    for(int idx = 0; idx < uidx;idx++) {
        if (verify_uuid(uuid[idx]) || verify_board_name(name[idx])) {
            usage();
            exit(0);
        }
        log_trace("UUID[%d] = %24s Name[%d] = %24s Schema[%d] = %s \n",
                  idx, uuid[idx], idx, name[idx], idx, schema[idx]);
    }

    int ret = create_db_hook(uuid, name, schema, uidx);
    if (ret) {
        log_error("MFGUTIL:: Error:: Failed to create registry DB.");
    } else {
        usys_log_info("MFGUTIL:: Created registry DB for device.");
        usys_log_info("MFGUTIL:: Sysroot: %s", femd_sysroot() ? femd_sysroot() : "/tmp/sys");
    }

    return 0;
}
