/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Tool for creating EEPROM database for device.
 * This binary is only meant to create a dummy device database by reading a JSON configs.
 * This tool would take JSON files, unit name and  serial number as command line arguments.
 *
 * Example usage is:
 * mfgutil -n COM -u UK-1001-COM-1101 -s mfgdata/schema/com.json -n LTE -u UK-2001-LTE-1101 -s mfgdata/schema/lte.json -n MASK -u UK-3001-MASK-1101 -s mfgdata/schema/mask.json
 * mfgutil -n RF_CTRL -u UK-5001-RFC-1101 -s mfgdata/schema/rfctrl.json -n RF_AMP -u UK-4001-RFA-1101 -s mfgdata/schema/rffe.json
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

#define PROPERTY_JSON "mfgdata/property/property.json"
#define VERSION	"0.0.1"

/* Writing a newly parsed json to DB (sysFs/eeprom) */
int create_db_hook(char **puuid, char** name, char** schema, int count) {
    int ret = 0;

    JSONInput jip;
    jip.fname = schema;
    jip.count = count;
    jip.pname = PROPERTY_JSON;

    UnitCfg *pcfg = NULL;
    DevI2cCfg *i2cCfg = NULL;

    /* Initializes for ledgers for devices */
    ret = ldgr_init(jip.pname);
    if (ret) {
        log_error("MFGUTIL:: ledger initialization failed %d", ret);
        goto cleanup;
    }

    /* Initializes manufacturing module.
     * Parses schema provided in JsonInput fname
     *  */
    ret = invt_mfg_init(&jip);
    if (ret) {
        log_error("MFGUTIL:: UBSP IDB init failed %d", ret);
        goto cleanup;
    }

  /* Will just initialize the db if NULL is passed*/
  ret = invt_init(NULL, &ldgr_register);
  if (ret) {
    usys_log_info("MFGUTIL:: UBSP init failed %d (Expected -1)", ret);
  }

    for(int idx = 0; idx < count; idx++) {
      log_debug("UUID[%d] = %24s Name[%d] = %24s Schema[%d] = %s \n", idx, puuid[idx], idx, name[idx], idx, jip.fname[idx]);

      /* Assumption Module Name in argument should match */
      UnitCfg *udata = (UnitCfg[]){
        { .modUuid = "UK-5001-RFC-1101",
          .modName = "RF CTRL BOARD",
        .sysFs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0051/eeprom",
        .eepromCfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
        { .modUuid = "UK-4001-RFA-1101",
            .modName = "RF BOARD",
            .sysFs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0052/eeprom",
            .eepromCfg = &(DevI2cCfg){ .bus = 2, .add = 0x50ul } },
            { .modUuid = "UK-1001-COM-1101",
                .modName = "ComV1",
                .sysFs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0050/eeprom",
                .eepromCfg = &(DevI2cCfg){ .bus = 0, .add = 0x50ul } },
                { .modUuid = "UK-2001-LTE-1101",
                    .modName = "LTE",
                    .sysFs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0050/eeprom",
                    .eepromCfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
                    { .modUuid = "UK-3001-MSK-1101",
                        .modName = "MASK",
                        .sysFs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0051/eeprom",
                        .eepromCfg = &(DevI2cCfg){ .bus = 1, .add = 0x51ul } },
      };

      /* Find and Read unitCfg of the module from above UnitCfg struct */
      for (int iter = 0; iter < MAX_BOARDS; iter++) {
        if (!usys_strcmp(name[idx], udata[iter].modName)) {

          pcfg = usys_zmalloc(sizeof(UnitCfg));
          if (pcfg) {

            usys_memset(pcfg, '\0', sizeof(UnitCfg));
            usys_memcpy(pcfg, &udata[iter], sizeof(UnitCfg));

            if (udata[iter].eepromCfg) {

              i2cCfg = usys_zmalloc(sizeof(DevI2cCfg));
              if (i2cCfg) {
                usys_memset(i2cCfg, '\0', sizeof(DevI2cCfg));
                usys_memcpy(i2cCfg, udata[iter].eepromCfg,
                    sizeof(DevI2cCfg));
              }

            }

            pcfg->eepromCfg = i2cCfg;
            usys_memcpy(pcfg->modUuid, puuid[idx], strlen(puuid[idx]));
            usys_memcpy(pcfg->modName, name[idx], strlen(name[idx]));

            break;

          } else {

            log_error("MFGUTIL:: Err(%d): Memory exhausted while getting unit"
                " config from Test data.",
              ERR_NODED_MEMORY_EXHAUSTED);
            goto cleanup;

          }
        }
      }

      /* Register Module */
      ret = invt_register_module(pcfg);
      if (!ret) {

        /* Create a EEPROM DB */
        ret = invt_create_db(pcfg->modUuid);
        if (!ret) {
          usys_log_info("MFGUTIL:: Created registry for module %s.", name[idx]);
        } else {
          log_error("MFGUTIL:: UBSP registry creation failed %d", ret);
          goto cleanup;
        }

      } else {
        log_error("MFGUTIL:: Registering module failed %d", ret);
        goto cleanup;
      }
    }

    /* Cleanup */
    cleanup:
    invt_mfg_exit();
    ldgr_exit();
    invt_exit();

    usys_free(pcfg->eepromCfg);
    pcfg->eepromCfg = NULL;
    usys_free(pcfg);
    pcfg = NULL;

    return ret;
}

static struct option longOptions[] = {
    { "name", required_argument, 0, 'n' },
    { "muuid", required_argument, 0, 'm' },
    { "logs", required_argument, 0, 'l' },
    { "help", no_argument, 0, 'h' },
    { "version", no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

/* Set the verbosity level for logs. */
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

/* Usage options for the ukamaEDR */
void usage() {
    printf("Usage: ukamaEDR [options] \n");
    printf("Options:\n");
    printf("--h, --help                                                      Help menu.\n");
    printf("--l, --logs <TRACE> <DEBUG> <INFO>                               Log level for the process.\n");
    printf("--n, --name <ComV1>|<LTE>|<MASK>|<RF CTRL BOARD>,<RF BOARD>      Name of module.\n");
    printf("--m, --muuid <Module UUID>                                       Module UUID.\n");
    printf("--s, --schema <json file path>                                   JSON Schema file.\n");
    printf("--v, --version                                                   Software Version.\n");
}

/* Utility to Create a EEPROM DB for devices.*/
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

    /* Parsing command line args. */
    while (true) {
        int opt = 0;
        int opIdx = 0;

        opt = usys_getopt_long(argc, argv, "h:v:m:n:s:l:", longOptions, &opIdx);
        if (opt == -1) {
          break;
        }

        switch (opt) {
        case 'h':
          usage();
          usys_exit(0);
          break;

        case 'v':
          usys_puts(VERSION);
          break;

        case 'n':
          name[nidx] = optarg;
          nidx++;
          break;

        case 'm':
          uuid[uidx] = optarg;
          uidx++;
          break;

        case 's':
          schema[sidx] = optarg;
          sidx++;
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

    /* Args check for schema info.*/
    if ((sidx != uidx) || (sidx != nidx) || (!sidx) || (sidx > MAX_BOARDS)  ) {
      log_error("MFGUTIL:: Name, schema and UUID entries have to match in count.");
      log_error("MFGUTIL:: At least one set of entries or %d set of entries can be made simultaneously.", MAX_BOARDS);
      exit(0);
    }

    /* Input args and their verification */
    for(int idx = 0; idx < uidx;idx++) {
      /* Verify module uuid and name */
      if (verify_uuid(uuid[idx]) || verify_board_name(name[idx]) ) {
        usage();
        exit(0);
      }
      log_trace("UUID[%d] = %24s Name[%d] = %24s Schema[%d] = %s \n", idx, uuid[idx], idx, name[idx], idx, schema[idx]);
    }

    /* Create EEPROM DB */
    int ret = create_db_hook(uuid, name, schema, uidx);
    if (ret) {
      log_error("MFGUTIL:: Error:: Failed to create registry DB for %s device.", name);
    } else {
      usys_log_info("MFGUTIL:: Created registry DB for device.");
      usys_log_info("MFGUTIL:: Copy directory from /tmp/sys");
    }

    return 0;
}
