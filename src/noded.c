/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "service.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_types.h"

#define DEV_PROPERTY_JSON "lib/ubsp/mfgdata/property/property.json"
#define INVENTORY_DB "/tmp/sys/cnode_inevetory_db"
#define DEF_LOG_LEVEL "TRACE"

#define NODED_VERSION "0.0.0"

/* Starting NodeD service */
void noded_service() {
    service();
}

/* Startup for NodeD.*/
int noded_startup(char *invDb, char *pCfg) {
    int ret = 0;
    ret = service_init();
    return ret;
}

/* NodeD service exit */
void noded_exit() {
    service_at_exit();
}

/* Terminate signal handler for ukamaEDR */
void handle_sigint(int signum) {
    usys_log_debug("Caught terminate signal.\n");

    /* Exiting NodeD */
    noded_exit();

    usys_log_debug("Cleanup complete.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "INVENTORY_DB", required_argument, 0, 's' },
    { "propertyConfig", required_argument, 0, 'p' },
    { "logs", required_argument, 0, 'l' },
    { "help", no_argument, 0, 'h' },
    { "version", no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

/* Set the verbosity level for logs. */
void set_log_level(char *slevel) {
    int ilevel = USYS_LOG_TRACE;
    if (!strcmp(slevel, "TRACE")) {
        ilevel = USYS_LOG_TRACE;
    } else if (!strcmp(slevel, "DEBUG")) {
        ilevel = USYS_LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        ilevel = USYS_LOG_INFO;
    }
    usys_log_set_level(ilevel);
}

/* Check if args supplied config file exist and have read permissions. */
void verify_file(char *file) {
    if (!usys_file_exist(file)) {
        usys_log_error("NodeD: File %s is missing.", file);
        exit(0);
    }
}

/* Usage options for the ukamaEDR */
void usage() {
    usys_puts("Usage: noded [options] \n");
    usys_puts("Options:\n");
    usys_puts("--h, --help                       Help menu.\n");
    usys_puts(
        "--l, --logs <TRACE> <DEBUG> <INFO>      Log level for the process.\n");
    usys_puts("--p, --propertyConfig <path>            Property config for the "
              "System.\n");
    usys_puts(
        "--s, --INVENTORY_DB <path>              System database or EEPROM DB for"
        " the System.\n");
    usys_puts("--v, --version                    Software Version.\n");
}

int main(int argc, char **argv) {
    int ret = USYS_OK;
    char *pCfg = DEV_PROPERTY_JSON;
    char *invDb = INVENTORY_DB;
    char *debug = DEF_LOG_LEVEL;

    /* Parsing command line args. */
    while (true) {
        int opt = 0;
        int opdIdx = 0;

        opt = usys_getopt_long(argc, argv, "s:p:l:", longOptions, &opdIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            usys_exit(0);
            break;

        case 'v':
            usys_puts(NODED_VERSION);
            break;

        case 's':
            invDb = optarg;
            verify_file(invDb);
            break;

        case 'p':
            pCfg = optarg;
            verify_file(pCfg);
            break;

        case 'l':
            debug = optarg;
            set_log_level(debug);
            break;

        default:
            usage();
            usys_exit(0);
        }
    }

    usys_log_debug(
        "Starting NodeD service for monitoring and configuring node.");

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* NodeD startup routine. */
    ret = noded_startup(invDb, pCfg);
    if (!ret) {
        /* NodeD Service started.*/
        noded_service();

        while (1) {
            usys_sleep(30);
        }
    }

    /* Should never reach here */
    noded_exit();

    usys_log_debug("Exiting NodeD service.");
    return ret;
}
