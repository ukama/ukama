/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>

#include "lookout.h"
#include "config.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_mem.h"

void handle_sigint(int signum) {
    usys_log_debug("Terminate signal.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "port",          required_argument, 0, 'p' },
    { "logs",          required_argument, 0, 'l' },
    { "noded-port",    required_argument, 0, 'd' },
    { "starter-port",  required_argument, 0, 's' },
    { "system-port",   required_argument, 0, 'S' },
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

    usys_puts("Usage: starter.d [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-p, --port <port>             Local listening port");
    usys_puts("-d, --noded-port    <port>    Node.d port");
    usys_puts("-s, --starter-port  <port>    Starter.d port");
    usys_puts("-S, --system-port   <port>    Node system port");
    usys_puts("-v, --version                 Software version");
}

int main(int argc, char **argv) {

    int opt, optIdx;
    char *debug        = DEF_LOG_LEVEL;
    char *port         = DEF_SERVICE_PORT;
    char *nodedPort    = DEF_NODED_PORT;
    char *starterPort  = DEF_STARTERD_PORT;
    char *systemPort   = DEF_NODE_SYSTEM_PORT;
    UInst  serviceInst; 
    Config serviceConfig = {0};

    /* Parsing command line args. */
    while (true) {

        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "vh:p:l:n:s:S", longOptions,
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
            usys_puts(LOOKOUT_VERSION);
            usys_exit(0);
            break;

        case 'p':
            port = optarg;
            if (!port) {
                usage();
                usys_exit(0);
            }
            break;

        case 'l':
            debug = optarg;
            set_log_level(debug);
            break;

        case 'n':
            nodedPort = optarg;
            if (!nodedPort) {
                usage();
                usys_exit(0);
            }
            break;

        case 's':
            starterPort = optarg;
            if (!starterPort) {
                usage();
                usys_exit(0);
            }
            break;

        case 'S':
            systemPort = optarg;
            if (!systemPort) {
                usage();
                usys_exit(0);
            }
            break;

        default:
            usage();
            usys_exit(0);
        }
    }

    /* Service config update */
    serviceConfig.servicePort    = usys_atoi(port);
    serviceConfig.nodedPort      = usys_atoi(nodedPort);
    serviceConfig.starterdPort   = usys_atoi(starterPort);
    serviceConfig.nodeSystemPort = usys_atoi(systemPort);
    serviceConfig.nodeID         = NULL;

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
            usys_log_error("Failed to send health report to system at %s:%d",
                           DEF_NODE_SYSTEM_HOST,
                           serviceConfig.nodeSystemPort);
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
