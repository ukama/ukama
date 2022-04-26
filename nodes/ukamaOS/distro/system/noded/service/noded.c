/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <service.h>
#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_types.h"

#define DEV_PROPERTY_JSON "mfgdata/property/property.json"
#define INVENTORY_DB "/tmp/sys/cnode_inventory_db"
#define NOTIF_SERVER_URL "http://localhost:8090/"
#define DEF_LOG_LEVEL "TRACE"

#define NODED_VERSION "0.0.0"

/**
 * @fn      void noded_service()
 * @brief   Start Noded web service (REST server)
 *
 */
void noded_service() {
    service();
}

/**
 * @fn      int noded_startup(char*, char*, char*)
 * @brief   Do service initialization. Parse the required configs and
 *          initialize web frameworks.
 *
 * @param   invDb
 * @param   pCfg
 * @param   notifServer
 * @return  On success 0,
 *          On failure -1
 */
int noded_startup(char *invDb, char *pCfg, char* notifServer) {
    int ret = 0;
    ret = service_init(invDb, pCfg, notifServer);
    return ret;
}

/**
 * @fn      void noded_exit()
 * @brief   Service exit procedure. Release the data structure used.
 */
void noded_exit() {
    service_at_exit();
}

/**
 * @fn      void handle_sigint(int)
 * @brief   Handle terminate signal for Noded
 *
 * @param   signum
 */
void handle_sigint(int signum) {
    usys_log_debug("Caught terminate signal.\n");

    /* Exiting NodeD */
    noded_exit();

    usys_log_debug("Cleanup complete.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "inventoryDb", required_argument, 0, 'i' },
    { "propertyConfig", required_argument, 0, 'p' },
    { "notifyServer", required_argument, 0, 'n' },
    { "logs", required_argument, 0, 'l' },
    { "help", no_argument, 0, 'h' },
    { "version", no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

/**
 * @fn      void set_log_level(char*)
 * @brief   Set the verbosity level for logs.
 *
 * @param   slevel
 */
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



/**
 * @fn      void verify_file(char*)
 * @brief   Check if args supplied config file exist and have read permissions.
 *
 * @param   file
 */
void verify_file(char *file) {
    if (!usys_file_exist(file)) {
        usys_log_error("NodeD: File %s is missing.", file);
        exit(0);
    }
}


/**
 * @fn      void usage()
 * @brief   Usage options for the ukamaEDR
 *
 */
void usage() {
    usys_puts("Usage: noded [options] \n");
    usys_puts("Options:\n");
    usys_puts(
        "--h, --help                             Help menu.\n");
    usys_puts(
        "--l, --logs <TRACE> <DEBUG> <INFO>      Log level for the process.\n");
    usys_puts(
        "--p, --propertyConfig <path>            Property config for the "
              "System.\n");
    usys_puts(
        "--i, --inventoryDb <path>               Inventory database or EEPROM DB"
        " for the System.\n");
    usys_puts(
        "--n, --notifyServer <url>               Notification server for "
        "alerts\n");
    usys_puts(
        "--v, --version                          Software Version.\n");
}

/**
 * @fn      int main(int, char**)
 * @brief
 *
 * @param   argc
 * @param   argv
 * @return  Should stay in main function entire time.
 */
int main(int argc, char **argv) {
    int ret = USYS_OK;
    char *pCfg = DEV_PROPERTY_JSON;
    char *invDb = INVENTORY_DB;
    char *debug = DEF_LOG_LEVEL;
    char *notifServer = NOTIF_SERVER_URL;

    /* Parsing command line args. */
    while (true) {
        int opt = 0;
        int opdIdx = 0;

        opt = usys_getopt_long(argc, argv, "h:i:p:l:v:n:", longOptions, &opdIdx);
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
            usys_exit(0);
            break;

        case 'i':
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

        case 'n':
            notifServer = optarg;
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
    ret = noded_startup(invDb, pCfg, notifServer);
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
