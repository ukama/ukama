/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <errno.h>
#include <getopt.h>
#include <signal.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "collector.h"
#include "file.h"
#include "web_service.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_services.h"
#include "usys_getopt.h"

#include "version.h"

#define METRIC_CONFIG "./config/config.toml"
#define DEF_LOG_LEVEL "TRACE"

/* define in network.c */
extern int start_admin_web_service(UInst *adminInst);

/* Terminate signal handler for Metrics collector */
void handle_sigint(int signum) {
  usys_log_info("Caught terminate signal.");

  /* Exiting Metrics */
  collector_exit(signum);
}

static struct option longOptions[] = {{"config", required_argument, 0, 'c'},
                                       {"logs", required_argument, 0, 'l'},
                                       {"help", no_argument, 0, 'h'},
                                       {"version", no_argument, 0, 'v'},
                                       {0, 0, 0, 0}};

void set_log_level(char *slevel) {

    int ilevel = USYS_LOG_TRACE;

    if (!strcmp(slevel, "TRACE")) {
        ilevel = USYS_LOG_TRACE;
    } else if (!strcmp(slevel, "DEBUG")) {
        ilevel = USYS_LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        ilevel = USYS_LOG_INFO;
    }

    log_set_level(ilevel);
}

/* Check if args supplied config file exist and have read permissions. */
void verify_file(char *file) {
  if (!file_exist(file)) {
    usys_log_error("Metrics: File %s is missing.", file);
    exit(0);
  }
}

void usage() {

    usys_puts("Usage: metrics.d [options]");
    usys_puts("Options:");
    usys_puts("--h, --help                         Help menu");
    usys_puts("--l, --logs <TRACE> <DEBUG> <INFO>  Log level for the process");
    usys_puts("--c, --config <path>                Config for the metrics collection");
    usys_puts("--v, --version                      Software version");
}

int main(int argc, char **argv) {

    int ret = 0;
    char *cfg = METRIC_CONFIG;
    char *debug = DEF_LOG_LEVEL;

    UInst adminInst;

    usys_log_set_service(SERVICE_NAME);
    //    usys_log_remote_init(SERVICE_NAME);

    /* Parsing command line args. */
    while (true) {

        int opt = 0;
        int opdIdx = 0;

        opt = getopt_long(argc, argv, "c:l", longOptions, &opdIdx);
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
            exit(0);

        case 'c':
            cfg = optarg;
            verify_file(cfg);
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

    usys_log_info("Starting metrics collector.");

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* start admin webservice */
    if(start_admin_web_service(&adminInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup for admin. Exiting");
        exit(1);
    }
  
    /* Start metrics collector. */
    ret = collector(cfg);

    usys_log_info("Stopping metrics collector.");
    return ret;
}
