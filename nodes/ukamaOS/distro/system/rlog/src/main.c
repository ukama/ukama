/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <getopt.h>
#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <unistd.h>

#include "ulfius.h"

#include "usys_api.h"
#include "usys_mem.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_services.h"

#include "rlogd.h"
#include "version.h"

extern int start_web_services(int port, UInst *serviceInst);
extern int get_nodeID_from_noded(char **nodeID, char *host, int port);

ThreadData *gData = NULL;

static volatile sig_atomic_t gTerminate = 0;

static void usage(void) {
    printf("rlog.d: local node logging facility\n");
    printf("Usage: rlog.d [options]\n");
    printf("Options:\n");
    printf("  -l, --level LEVEL  trace, debug, info, warn, error, critical\n");
    printf("  -v, --version      Show version\n");
    printf("  -h, --help         Show this help\n");
}

static void on_signal(int signalNumber) {
    (void)signalNumber;
    gTerminate = 1;
}

static int install_signal_handlers(void) {
    struct sigaction action;

    memset(&action, 0, sizeof(action));
    action.sa_handler = on_signal;
    sigemptyset(&action.sa_mask);

    if (sigaction(SIGTERM, &action, NULL) != 0 ||
        sigaction(SIGINT, &action, NULL) != 0) {
        return -1;
    }

    return 0;
}

static int parse_log_level(const char *value) {
    if (!value) return -1;

    if (strcasecmp(value, "trace") == 0) return USYS_LOG_TRACE;
    if (strcasecmp(value, "debug") == 0) return USYS_LOG_DEBUG;
    if (strcasecmp(value, "info") == 0) return USYS_LOG_INFO;
    if (strcasecmp(value, "warn") == 0) return USYS_LOG_WARN;
    if (strcasecmp(value, "error") == 0) return USYS_LOG_ERROR;
    if (strcasecmp(value, "critical") == 0 ||
        strcasecmp(value, "fatal") == 0) {
        return USYS_LOG_CRITICAL;
    }

    return -1;
}

static int parse_options(int argc, char **argv) {
    int option;
    int optionIndex;

    static struct option longOptions[] = {
        {"level", required_argument, NULL, 'l'},
        {"help", no_argument, NULL, 'h'},
        {"version", no_argument, NULL, 'v'},
        {NULL, 0, NULL, 0}
    };

    while (true) {
        optionIndex = 0;
        option = getopt_long(argc, argv, "l:hv", longOptions,
                             &optionIndex);
        if (option == -1) break;

        switch (option) {
        case 'l': {
            int level = parse_log_level(optarg);

            if (level < 0) {
                fprintf(stderr, "rlog.d: invalid log level: %s\n",
                        optarg);
                return -1;
            }
            gData->level = level;
            usys_log_set_level(level);
            break;
        }
        case 'h':
            usage();
            exit(0);
        case 'v':
            printf("rlog.d - Version: %s\n", VERSION);
            exit(0);
        default:
            return -1;
        }
    }

    return optind == argc ? 0 : -1;
}

static char *resolve_node_id(int nodedPort) {
    const char *configured;
    char *nodeId;

    configured = getenv("RLOG_NODE_ID");
    if (configured && *configured) {
        return strdup(configured);
    }

    nodeId = NULL;
    if (nodedPort > 0 &&
        get_nodeID_from_noded(&nodeId, DEF_NODED_HOST,
                              nodedPort) == USYS_TRUE) {
        return nodeId;
    }

    usys_log_warn("Unable to retrieve node ID; using default");
    return strdup(DEF_NODE_ID);
}

static void cleanup(UInst *serviceInst, bool serviceStarted,
                    char *nodeId) {
    if (serviceStarted) {
        ulfius_stop_framework(serviceInst);
        ulfius_clean_instance(serviceInst);
    }

    if (gData) {
        if (gData->ingest) {
            ingest_stop(gData->ingest);
            gData->ingest = NULL;
        }
        if (gData->store) {
            log_store_close(gData->store);
            gData->store = NULL;
        }
        free(gData);
        gData = NULL;
    }

    usys_free(nodeId);
}

int main(int argc, char **argv) {
    UInst serviceInst;
    LogStoreConfig storeConfig;
    const char *socketPath;
    char *nodeId;
    int nodedPort;
    int adminPort;
    bool serviceStarted;
    int exitCode;

    memset(&serviceInst, 0, sizeof(serviceInst));
    memset(&storeConfig, 0, sizeof(storeConfig));
    serviceStarted = false;
    nodeId = NULL;
    exitCode = 1;

    usys_log_set_service(SERVICE_NAME);

    gData = calloc(1, sizeof(*gData));
    if (!gData) {
        return 1;
    }
    gData->level = USYS_LOG_DEBUG;
    usys_log_set_level(gData->level);

    if (parse_options(argc, argv) != 0) {
        usage();
        goto done;
    }

    if (install_signal_handlers() != 0) {
        usys_log_error("Unable to install signal handlers");
        goto done;
    }

    nodedPort = usys_find_service_port(SERVICE_NODE);
    adminPort = usys_find_service_port(SERVICE_RLOG_ADMIN);
    if (adminPort <= 0) {
        usys_log_error("Unable to find rlog admin service port");
        goto done;
    }

    nodeId = resolve_node_id(nodedPort);
    if (!nodeId) {
        usys_log_error("Unable to allocate node ID");
        goto done;
    }

    storeConfig.logDir = getenv("RLOG_LOG_DIR");
    storeConfig.nodeId = nodeId;
    gData->store = log_store_open(&storeConfig);
    if (!gData->store) {
        usys_log_error("Unable to open canonical log store");
        goto done;
    }

    socketPath = getenv("RLOG_INGEST_SOCKET");
    if (!socketPath || !*socketPath) {
        socketPath = DEF_INGEST_SOCKET;
    }

    gData->ingest = ingest_start(socketPath, gData->store);
    if (!gData->ingest) {
        usys_log_error("Unable to start local ingest socket: %s",
                       socketPath);
        goto done;
    }

    if (start_web_services(adminPort, &serviceInst) != USYS_TRUE) {
        usys_log_error("Unable to start admin service on port %d",
                       adminPort);
        goto done;
    }
    serviceStarted = true;

    usys_log_info("rlog.d ready: node=%s ingest=%s admin=%d",
                  nodeId, socketPath, adminPort);

    while (!gTerminate) {
        pause();
    }

    usys_log_info("rlog.d terminating");
    exitCode = 0;

done:
    cleanup(&serviceInst, serviceStarted, nodeId);
    return exitCode;
}
