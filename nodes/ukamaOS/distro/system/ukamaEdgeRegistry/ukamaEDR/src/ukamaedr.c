/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inc/ukamadr.h"

#include "inc/reghelper.h"
#include "inc/registry.h"
#include "inc/alarmhandler.h"
#include "inc/telemhandler.h"
#include "headers/errorcode.h"
#include "msghandler.h"
#include "headers/globalheader.h"
#include "headers/ubsp/devices.h"
#include "headers/utils/file.h"
#include "headers/utils/log.h"
#include "headers/ubsp/ubsp.h"
#include "headers/ubsp/ukdblayout.h"
#include "dmt.h"

#include <getopt.h>
#include <pthread.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#define PROPERTYJSON "lib/ubsp/mfgdata/property/property.json"
#define SYSTEMDB "/tmp/sys/cnode-systemdb"
#define DEFLOGLEVEL "TRACE"

/* Starting ukamaEDR service */
void ukama_service() {
    /* start the msghandler process*/
    msghandler_start();

    /* Start alarm handler reporting*/
    alarmhandler_start();

    /* start the telemetry process*/
    telemhandler_start();
}

/* Setting registry for ukamaEDR */
int ukama_setup_drdb() {
    int ret = 0;
    /* Create Db */
    reg_init();

    /* Register Unit and  Module */
    ret |= reg_register_misc();

    /* Register devices */
    ret |= reg_register_devices();

    /* Print registered devices */
    reg_list_reg_devices();

    return ret;
}

/* Registry Exit.*/
void ukama_cleanup_drdb() {
    /* Destroy Db */
    reg_exit();
}

/* Startup for UkamaEDR.*/
int ukama_startup(char *sysdb, char *pcfg) {
    int ret = RET_SUCCESS;

    /* Parse device property from json and prepare IRQDB */
    ret = ubsp_devdb_init(pcfg);
    if (ret) {
        log_error("Err(%d): UKAMADR:: UBSP device DB initialization failed.");
        goto end;
    }

    /*Parse the systemDB from the EEPROM link. After that register module and devices in UBSP lib*/
    ret = ubsp_ukdb_init(sysdb);
    if (ret) {
        log_error("Err(%d): UKAMADR:: UBSP System DB initialization failed.");
        goto end;
    }

    /* Alarm handler Creates a thread and queue for alarms.*/
    alarmhandler_init();

    /* Create registry for Unit, Module and Devices */
    ret = ukama_setup_drdb();
    if (ret) {
        log_error("Err(%d): UKAMADR:: UBSP Driver DB initialization failed.");
        goto end;
    }

    /* Start Telemetry handler thread  */
    telemhandler_init();

    /* Start msghandler thread for UkamaEDR server. */
    msghandler_init();

end:
    return ret;
}

/* Ukama exit */
void ukama_exit() {
    /* Exit msg handler thread */
    msghandler_exit();

    /* Exit telemetry handler thread */
    telemhandler_exit();

    /*De-initilaize bsp lib*/
    ubsp_exit();

    /* Cleaning registry */
    ukama_cleanup_drdb();

    /* Exit alarm handler */
    alarmhandler_exit();
}

/* Terminate signal handler for ukamaEDR */
void handle_sigint(int signum) {
    log_debug("UkamaEDR: Caught terminate signal.\n");

    /* Exiting UkamaEDR */
    ukama_exit();

    log_debug("UkamaEDR: Cleanup complete.\n");
    exit(0);
}

/* UkamaEDR banner. */
void ukama_edr_banner() {
    log_trace(
        "|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||");
    log_trace(
        "|||  |||||  |||||  ||||  |||||         |||||    ||||||    |||||         |||||        |||||         |||||         ||||||||");
    log_trace(
        "|||  |||||  |||||  |||  ||||||  |||||  |||||  |  ||||  |  |||||  |||||  |||||  |||||||||||  |||||  |||||  |||||  ||||||||");
    log_trace(
        "|||  |||||  |||||  ||  |||||||  |||||  |||||  ||  ||  ||  |||||  |||||  |||||  |||||||||||  |||||  |||||  |||||  ||||||||");
    log_trace(
        "|||  |||||  |||||  |  ||||||||         |||||  ||||  ||||  |||||         |||||       ||||||  |||||  |||||         ||||||||");
    log_trace(
        "|||  |||||  |||||  ||  |||||||  |||||  |||||  ||||||||||  |||||  |||||  |||||  |||||||||||  |||||  |||||  ||  |||||||||||");
    log_trace(
        "|||  |||||  |||||  |||  ||||||  |||||  |||||  ||||||||||  |||||  |||||  |||||  |||||||||||  |||||  |||||  |||  ||||||||||");
    log_trace(
        "|||         |||||  ||||  |||||  |||||  |||||  ||||||||||  |||||  |||||  |||||        |||||         |||||  ||||  |||||||||");
    log_trace(
        "|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||");
    log_trace(
        "|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||");
}

static struct option long_options[] = {
    { "systemdb", required_argument, 0, 's' },
    { "propertyConfig", required_argument, 0, 'p' },
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

/* Check if args supplied config file exist and have read permissions. */
void verify_file(char *file) {
    if (!file_exist(file)) {
        log_error("UkamaEDR: File %s is missing.", file);
        exit(0);
    }
}

/* Usage options for the ukamaEDR */
void usage() {
    printf("Usage: ukamaEDR [options] \n");
    printf("Options:\n");
    printf("--h, --help                             Help menu.\n");
    printf(
        "--l, --logs <TRACE> <DEBUG> <INFO>      Log level for the process.\n");
    printf(
        "--p, --propertyConfig <path>            Property config for the System.\n");
    printf(
        "--s, --systemdb <path>                  System database or EEPROM DB for the System.\n");
    printf("--v, --version                          Software Version.\n");
}

int main(int argc, char **argv) {
    int ret = RET_SUCCESS;
    char *pcfg = PROPERTYJSON;
    char *sysdb = SYSTEMDB;
    char *debug = DEFLOGLEVEL;

    /* Parsing command line args. */
    while (true) {
        int opt = 0;
        int opdidx = 0;

        opt = getopt_long(argc, argv, "s:p:l:", long_options, &opdidx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            exit(0);
            break;

        case 'v':
            puts("option -b\n");
            break;

        case 's':
            sysdb = optarg;
            verify_file(sysdb);
            break;

        case 'p':
            pcfg = optarg;
            verify_file(pcfg);
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

    ukama_edr_banner();

    log_debug(
        "UKAMADR:: Started Ukama Driver Registry and Monitoring process.");

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* Start up for EDR. */
    ret = ukama_startup(sysdb, pcfg);

    if (!ret) {
        /* Serivce started.*/
        ukama_service();

        while (1) {
            sleep(30);
        }
    }

    /* Should never reach here */
    ukama_exit();
    
    dmt_dump();

    log_debug(
        "UKAMADR:: Exiting Ukama Driver Registry and Monitoring process.");
    return ret;
}
