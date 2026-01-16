/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdio.h>

#include "lookout.h"
#include "config.h"
#include "web_client.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_mem.h"
#include "usys_services.h"

#include "version.h"

/* network.c */
extern int start_web_services(Config *config, UInst *serviceInst);

void handle_sigint(int signum) {
    usys_log_debug("Terminate signal.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "logs",          required_argument, 0, 'l' },
    { "help",          no_argument, 0, 'h' },
    { "version",       no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

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

void usage() {

    usys_puts("Usage: lookout.d [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-v, --version                 Software version");
}

int main(int argc, char **argv) {

    int opt, optIdx;
    char *debug        = DEF_LOG_LEVEL;
    UInst  serviceInst; 
    Config serviceConfig = {0};

    usys_log_set_service(SERVICE_NAME);
    //    usys_log_remote_init(SERVICE_NAME);

    /* Parsing command line args. */
    while (true) {

        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "vh:l:S", longOptions,
                               &optIdx);
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
            usys_exit(0);
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

    serviceConfig.servicePort    = usys_find_service_port(SERVICE_NAME);
    serviceConfig.nodedPort      = usys_find_service_port(SERVICE_NODE);
    serviceConfig.starterdPort   = usys_find_service_port(SERVICE_STARTER);
    serviceConfig.nodeID         = NULL;

    if (!usys_find_service_port(SERVICE_UKAMA)) {
        usys_log_error("Unable to determine the port for Ukama");
        usys_exit(1);
    }

    if (!serviceConfig.servicePort  ||
        !serviceConfig.nodedPort    ||
        !serviceConfig.starterdPort) {
        usys_log_error("Unable to determine the port for services");
        usys_exit(1);
    }

    usys_log_debug("Starting %s ... ", SERVICE_NAME);

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* Read Node Info from noded */
    if (getenv(ENV_LOOKOUT_DEBUG_MODE)) {
       serviceConfig.nodeID = strdup(DEF_NODE_ID);
       usys_log_debug("%s: Using default Node ID: %s", SERVICE_NAME, DEF_NODE_ID);
    } else {
        if (get_nodeid_from_noded(&serviceConfig) == STATUS_NOK) {
            usys_log_error("%s: Unable to connect with node.d", SERVICE_NAME);
            goto done;
        }
    }

    /* start web service */
    if (start_web_services(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("%s: unable to start webservice. Exiting.", SERVICE_NAME);
        exit(1);
    }

    /* until interrupted by SIG */
    while (USYS_TRUE) {
        if (send_health_report(&serviceConfig) == USYS_FALSE) {
            usys_log_error("Failed to send health report to ukama system");
        }
        sleep(DEF_REPORT_INTERVAL);
    }

done:
    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);

    usys_free(serviceConfig.nodeID);

    usys_log_debug("Exiting %s ...", SERVICE_NAME);

    return USYS_TRUE;
}
