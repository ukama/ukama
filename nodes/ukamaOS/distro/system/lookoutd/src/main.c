/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

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

static bool str_eq(const char *a, const char *b) {

    if (a == NULL || b == NULL) return false;

    return strcmp(a, b) == 0;
}

static bool str_contains(const char *s, const char *needle) {

    if (s == NULL || needle == NULL) {
        return false;
    }

    return strstr(s, needle) != NULL;
}

static LookoutNodeType detect_node_type(const char *nodeID) {

    if (nodeID == NULL || nodeID[0] == '\0') {
        return LOOKOUT_NODE_UNKNOWN;
    }

    if (str_contains(nodeID, "tnode") ||
        str_contains(nodeID, "tower")) {
        return LOOKOUT_NODE_TOWER;
    }

    if (str_contains(nodeID, "anode") ||
        str_contains(nodeID, "amplifier")) {
        return LOOKOUT_NODE_AMPLIFIER;
    }

    if (str_contains(nodeID, "cnode") ||
        str_contains(nodeID, "control")) {
        return LOOKOUT_NODE_CONTROL;
    }

    return LOOKOUT_NODE_UNKNOWN;
}

static const char *node_type_to_string(LookoutNodeType nodeType) {

    switch (nodeType) {
    case LOOKOUT_NODE_TOWER:
        return "tower";

    case LOOKOUT_NODE_AMPLIFIER:
        return "amplifier";

    case LOOKOUT_NODE_CONTROL:
        return "control";

    default:
        return "unknown";
    }
}

static void configure_app_manager(Config *config) {

    char *mgr = NULL;

    config->appManager = LOOKOUT_APP_MANAGER_STARTERD;

    mgr = getenv(ENV_LOOKOUT_APP_MANAGER);
    if (str_eq(mgr, LOOKOUT_MANAGER_SUPERVISOR)) {
        config->appManager = LOOKOUT_APP_MANAGER_SUPERVISORD;
    }
}

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

    int opt, optIdx, ret=-1;
    char *debug = DEF_LOG_LEVEL;

    UInst  serviceInst;
    Config serviceConfig = {0};

    memset(&serviceInst, 0, sizeof(UInst));

    usys_log_set_service(SERVICE_NAME);

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

    configure_app_manager(&serviceConfig);

    serviceConfig.servicePort  = usys_find_service_port(SERVICE_NAME);
    serviceConfig.nodedPort    = usys_find_service_port(SERVICE_NODE);
    serviceConfig.starterdPort = usys_find_service_port(SERVICE_STARTER);

    serviceConfig.nodeType        = LOOKOUT_NODE_UNKNOWN;
    serviceConfig.isTowerNode     = false;
    serviceConfig.isAmplifierNode = false;
    serviceConfig.isControlNode   = false;

    if (!usys_find_service_port(SERVICE_UKAMA)) {
        usys_log_error("Unable to determine the port for Ukama");
        goto done;
    }

    if (!serviceConfig.servicePort || !serviceConfig.nodedPort) {
        usys_log_error("Unable to determine the port for required services");
        goto done;
    }

    if (serviceConfig.appManager == LOOKOUT_APP_MANAGER_STARTERD &&
        !serviceConfig.starterdPort) {
        usys_log_error("Unable to determine the port for starter.d");
        goto done;
    }

    usys_log_debug("Starting %s ... ", SERVICE_NAME);

    signal(SIGINT, handle_sigint);

    if (getenv(ENV_LOOKOUT_DEBUG_MODE)) {
        serviceConfig.nodeID = strdup(DEF_NODE_ID);
        usys_log_debug("%s: Using default Node ID: %s",
                       SERVICE_NAME,
                       DEF_NODE_ID);
    } else {
        if (get_nodeid_from_noded(&serviceConfig) == STATUS_NOK) {
            usys_log_error("%s: Unable to connect with node.d",
                           SERVICE_NAME);
            goto done;
        }
    }

    serviceConfig.nodeType = detect_node_type(serviceConfig.nodeID);

    serviceConfig.isTowerNode =
        serviceConfig.nodeType == LOOKOUT_NODE_TOWER;
    serviceConfig.isAmplifierNode =
        serviceConfig.nodeType == LOOKOUT_NODE_AMPLIFIER;
    serviceConfig.isControlNode =
        serviceConfig.nodeType == LOOKOUT_NODE_CONTROL;

    usys_log_info("%s: nodeID=%s nodeType=%s appManager=%s",
                  SERVICE_NAME,
                  serviceConfig.nodeID,
                  node_type_to_string(serviceConfig.nodeType),
                  serviceConfig.appManager == LOOKOUT_APP_MANAGER_STARTERD ?
                  "starter" : "supervisor");

    if (start_web_services(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("%s: unable to start webservice. Exiting.",
                       SERVICE_NAME);
        goto done;
    }

    while (USYS_TRUE) {
        if (send_health_report(&serviceConfig) == USYS_FALSE) {
            usys_log_error("Failed to send health report to ukama system");
        }

        sleep(DEF_REPORT_INTERVAL);
    }

    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);

    ret = 0;
done:
    usys_free(serviceConfig.nodeID);

    usys_log_debug("Exiting %s ...", SERVICE_NAME);

    return ret;
}
