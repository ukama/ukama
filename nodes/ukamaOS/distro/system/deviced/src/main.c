/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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

void handle_sigint(int signum) {
    usys_log_debug("Terminate signal.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "port",        required_argument, 0, 'p' },
    { "logs",        required_argument, 0, 'l' },
    { "notify-port", required_argument, 0, 'n' },
    { "noded-port",  required_argument, 0, 'd' },
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
    usys_puts("--h, --help                    Help menu");
    usys_puts("--l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("--p, --port <port>             Local listening port");
    usys_puts("--n, --notify-port <port>      Notify.d port");
    usys_puts("--d, --noded-port <port>       Node.d port");
    usys_puts("--c, --client-mode             Run as client");
    usys_puts("--v, --version                 Software version");
}

int main(int argc, char **argv) {

    int opt, optIdx, clientMode=USYS_FALSE;
    char *debug        = DEF_LOG_LEVEL;
    char *port         = DEF_SERVICE_PORT;
    char *notifyPort   = DEF_NOTIFY_PORT;
    char *nodedPort    = DEF_NODED_PORT;
    UInst serviceInst;
    Config serviceConfig = {0};

    /* Parsing command line args. */
    while (true) {
        
        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "h:p:l:v:n:d", longOptions,
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
            usys_puts(DEVICED_VERSION);
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

        case 'd':
            notifyPort = optarg;
            if (!notifyPort) {
                usage();
                usys_exit(0);
            }
            break;

        case 'm':
            clientMode = USYS_TRUE;
            break;

        default:
            usage();
            usys_exit(0);
        }
    }

    /* Service config update */
    serviceConfig.serviceName  = usys_strdup(SERVICE_NAME);
    serviceConfig.servicePort  = usys_atoi(port);
    serviceConfig.nodedPort    = usys_atoi(nodedPort);
    serviceConfig.notifydPort  = usys_atoi(notifyPort);
    serviceConfig.nodeID       = NULL;
    serviceConfig.nodeType     = NULL;
    serviceConfig.clientMode   = clientMode;

    usys_log_debug("Starting %s ...", SERVICE_NAME);

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* Read Node Info from noded */
    if (getenv(ENV_DEVICED_DEBUG_MODE)) {
       serviceConfig.nodeID = strdup(DEF_NODE_ID);
       usys_log_debug("%s: using default Node ID: %s",
                      SERVICE_NAME,
                      DEF_NODE_ID);
    } else {
        if (get_nodeid_and_type_from_noded(&serviceConfig) == STATUS_NOK) {
            usys_log_error("%s: unable to connect with node.d", SERVICE_NAME);
            goto done;
        }
    }

    if (start_web_service(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup for clients. Exiting.");
        exit(1);
    }

    pause();

done:
    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);
    free(serviceConfig.serviceName);

    usys_log_debug("Exiting device.d ...");

    return USYS_TRUE;
}
