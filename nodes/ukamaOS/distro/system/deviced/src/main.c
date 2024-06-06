/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "config.h"
#include "deviced.h"
#include "web_client.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_services.h"

#include "version.h"

void handle_sigint(int signum) {
    usys_log_debug("Terminate signal.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "logs",        required_argument, 0, 'l' },
    { "client-host", required_argument, 0, 'H' },
    { "client-mode", no_argument, 0, 'c' },
    { "help",        no_argument, 0, 'h' },
    { "version",     no_argument, 0, 'v' },
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

    usys_puts("Usage: device.d [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-c, --client-mode             Run as client");
    usys_puts("-H, --client-host             Host where client is running");
    usys_puts("-v, --version                 Software version");
}

int main(int argc, char **argv) {

    int opt, optIdx, clientMode=USYS_FALSE;

    char *debug        = DEF_LOG_LEVEL;
    char *clientHost   = DEF_SERVICE_CLIENT_HOST;
    UInst serviceInst;
    Config serviceConfig = {0};

    usys_log_set_service(SERVICE_NAME);
    usys_log_remote_init(SERVICE_NAME);
    /* Parsing command line args. */
    while (true) {
        
        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "cvh:l:H", longOptions,
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

        case 'c':
            clientMode = USYS_TRUE;
            break;

        case 'H':
            clientHost = optarg;
            break;

        default:
            usage();
            usys_exit(0);
        }
    }

    /* Service config update */
    serviceConfig.serviceName  = usys_strdup(SERVICE_NAME);
    if (!clientMode)  {
        serviceConfig.servicePort  = usys_find_service_port(SERVICE_NAME);
    } else {
        serviceConfig.servicePort  =
            usys_find_service_port(SERVICE_DEVICE_CLIENT);
    }
    serviceConfig.nodedPort    = usys_find_service_port(SERVICE_NODE);
    serviceConfig.notifydPort  = usys_find_service_port(SERVICE_NOTIFY);
    serviceConfig.nodeID       = NULL;
    serviceConfig.nodeType     = NULL;
    serviceConfig.clientMode   = clientMode;
    serviceConfig.clientHost   = strdup(clientHost);
    serviceConfig.clientPort   = usys_find_service_port(SERVICE_DEVICE_CLIENT);

    if (!serviceConfig.servicePort ||
        !serviceConfig.nodedPort   ||
        !serviceConfig.notifydPort ||
        !serviceConfig.clientPort) {
        usys_log_error("Unable to determine the port for services");
        usys_exit(1);
    }

    usys_log_debug("Starting %s ... [client-mode:%d]",
                   SERVICE_NAME, clientMode);

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* Read Node Info from node.d */
    if (serviceConfig.clientMode == USYS_FALSE) {
        if (getenv(ENV_DEVICED_DEBUG_MODE)) {
            serviceConfig.nodeID   = strdup(DEF_NODE_ID);
            serviceConfig.nodeType = strdup(DEF_NODE_TYPE);
            usys_log_debug("%s: using default Node ID: %s Type: %s",
                           SERVICE_NAME,
                           DEF_NODE_ID,
                           DEF_NODE_TYPE);
        } else {
            if (get_nodeid_and_type_from_noded(&serviceConfig) == STATUS_NOK) {
                usys_log_error(
                    "%s: unable to connect with node.d", SERVICE_NAME);
                goto done;
            }
        }
    }

    if (start_web_service(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup for clients. Exiting.");
        usys_exit(1);
    }

    pause();

done:
    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);
    free(serviceConfig.serviceName);

    usys_log_debug("Exiting device.d ...");

    return USYS_TRUE;
}
