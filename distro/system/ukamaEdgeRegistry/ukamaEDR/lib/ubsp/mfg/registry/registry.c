
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

#include "headers/errorcode.h"
#include "headers/ubsp/devices.h"
#include "headers/ubsp/property.h"
#include "headers/ubsp/ubsp.h"
#include "inc/devicedb.h"
#include "inc/globalheader.h"
#include "inc/ukdb.h"
#include "test/test.h"
#include "utils/crc32.h"
#include "headers/utils/log.h"
#include "ukdb/db/db.h"
#include "ukdb/db/file.h"
#include "ukdb/idb/cs.h"
#include "ukdb/idb/idb.h"

#include <getopt.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

//#define PROPERTYJSON "mfgdata/property/property.json"
#define PROPERTYJSON "lib/ubsp/mfgdata/property/property.json"
#define VERSION	"0.0.1"
#define MAX_BOARDS	5


/* Writing a newly parsed json to DB (sysfs/eeprom) */
int create_db_hook(char **puuid, char** name, char** schema, int count) {
    int ret = 0;

    JSONInput jip;
    jip.fname = schema;
    jip.count = count;
    jip.pname = PROPERTYJSON;

    UnitCfg *pcfg = NULL;
    DevI2cCfg *i2c_cfg = NULL;

    ret = ubsp_devdb_init(jip.pname);
	if (ret) {
		log_error("MFGUTIL:: UBSP DEVDB init failed %d", ret);
		goto cleanup;
	}

	ret = ubsp_idb_init(&jip);
	if (ret) {
		log_error("MFGUTIL:: UBSP IDB init failed %d", ret);
		goto cleanup;
	}

	ret = ubsp_ukdb_init(NULL); /* Will just initialize the db if NULL is passed*/
	if (ret) {
		log_info("MFGUTIL:: UBSP init failed %d (Expected -1)", ret);
	}

    for(int idx = 0; idx < count; idx++) {
    	log_debug("UUID[%d] = %24s Name[%d] = %24s Schema[%d] = %s \n", idx, puuid[idx], idx, name[idx], idx, jip.fname[idx]);


    	UnitCfg *udata = (UnitCfg[]){
    		{ .mod_uuid = "UK-5001-RFC-1101",
    			.mod_name = "RF_CTRL",
				.sysfs = "/tmp/sys/bus/i2c/devices/i2c-0/0-0051/eeprom",
				.eeprom_cfg = &(DevI2cCfg){ .bus = 1, .add = 0x50ul } },
				{ .mod_uuid = "UK-4001-RFA-1101",
						.mod_name = "RF_AMP",
						.sysfs = "/tmp/sys/bus/i2c/devices/i2c-1/1-0052/eeprom",
						.eeprom_cfg = &(DevI2cCfg){ .bus = 2, .add = 0x50ul } },
						{ .mod_uuid = "UK-1001-COM-1101",
								.mod_name = "COM",
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


    	for (int iter = 0; iter < MAX_BOARDS; iter++) {
    		if (!strcmp(name[idx], udata[iter].mod_name)) {
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

    				memcpy(pcfg->mod_uuid, puuid[idx], strlen(puuid[idx]));
    				memcpy(pcfg->mod_name, name[idx], strlen(name[idx]));

    				break;
    			} else {
    				log_error("MFGUTIL:: Err(%d): Memory exhausted while getting unit"
    						" config from Test data.",
							ERR_UBSP_MEMORY_EXHAUSTED);
    				goto cleanup;
    			}
    		}
    	}

    	/* register Module */
    	ret = ubsp_register_module(pcfg);
    	if (!ret) {
    		ret = ubsp_create_ukdb(pcfg->mod_uuid);
    		if (!ret) {
    			log_info("MFGUTIL:: Created registry for module %s.", name[idx]);
    		} else {
    			log_error("MFGUTIL:: UBSP registry creation failed %d", ret);
    			goto cleanup;
    		}
    	} else {
    		log_error("MFGUTIL:: Registering module failed %d", ret);
    		goto cleanup;
    	}
    }
cleanup:
    ubsp_idb_exit();
    ubsp_exit();
    UBSP_FREE(pcfg->eeprom_cfg);
    UBSP_FREE(pcfg);

    return ret;
}

static struct option long_options[] = {
    { "name", required_argument, 0, 'n' },
    { "uuid", required_argument, 0, 'u' },
    { "logs", required_argument, 0, 'l' },
    { "help", no_argument, 0, 'h' },
    { "version", no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

/* Set the verbosity level for logs. */
void set_log_level(char *slevel) {
    int ilevel = LOG_TRACE;
    if (!strcmp(slevel, "TRACE")) {
        ilevel = LOG_TRACE;
    } else if (!strcmp(slevel, "DEBUG")) {
        ilevel = LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        ilevel = LOG_INFO;
    }
    log_set_level(ilevel);
}

/* Usage options for the ukamaEDR */
void usage() {
    printf("Usage: ukamaEDR [options] \n");
    printf("Options:\n");
    printf("--h, --help                             Help menu.\n");
    printf(
        "--l, --logs <TRACE> <DEBUG> <INFO>       Log level for the process.\n");
    printf(
        "--n, --name <ModuleName>                 Unit Type RF or COM.\n");
    printf(
        "--u, --uuid <unique Id>                  Unit unique Id.\n");
    printf(
        "--s, --schema <unique Id>               Schema file.\n");
    printf("--v, --version                       Software Version.\n");
}

int main(int argc, char** argv) {
    char *name[MAX_BOARDS] = {"\0"};
    char *uuid[MAX_BOARDS] = {"\0"};
    char *schema[MAX_BOARDS] = {"\0"};
    char *debug = "TRACE";
    char *ip = "";
    set_log_level(debug);

    if (argc < 2 ) {
    	log_error("Not enough arguments.");
    	exit(1);
    }

    int uidx = 0;
    int nidx = 0;
    int sidx = 0;
    /* Parsing command line args. */
    while (true) {
        int opt = 0;
        int opdidx = 0;

        opt = getopt_long(argc, argv, "h:v:u:n:s:l:", long_options, &opdidx);
        if (opt == -1) {
        	break;
        }

        switch (opt) {
        case 'h':
        	usage();
        	exit(0);
        	break;

        case 'v':
        	puts(VERSION);
        	break;

        case 'n':
        	name[nidx] = optarg;
        	nidx++;
        	break;

        case 'u':
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

    if ((sidx != uidx) || (sidx != nidx) || (sidx > MAX_BOARDS)) {
    	log_error("MFGUTIL:: Name, schema and UUID entries have to match in count.");
    	exit(0);
    }

    for(int idx = 0; idx < uidx;idx++) {
    	log_debug("UUID[%d] = %24s Name[%d] = %24s Schema[%d] = %s \n", idx, uuid[idx], idx, name[idx], idx, schema[idx]);
    }

    int ret = create_db_hook(uuid, name, schema, uidx);
    if (ret) {
    	log_error("MFGUTIL:: Error:: Failed to create registry DB for %s device.", name);
    } else {
    	log_info("MFGUTIL:: Created registry DB for device.");
    	log_info("MFGUTIL:: Copy directory from /tmp/sys");
    }

    return 0;
}
