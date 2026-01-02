/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <pthread.h>

#include "config.h"
#include "notify_macros.h"
#include "service.h"
#include "web.h"
#include "web_client.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_file.h"
#include "usys_mem.h"
#include "usys_services.h"

#include "version.h"

/* Global */
ThreadData *gData = NULL;

void handle_sigint(int signum) {
    usys_log_debug("Caught terminate signal.\n");
    usys_log_debug("Cleanup complete.\n");
    usys_exit(0);
}

/* network.c */
extern int start_admin_web_service(Config *config, UInst *adminInst);
extern int start_web_services(Config *config, UInst *serviceInst);

static UsysOption longOptions[] = {
    { "logs",          required_argument, 0, 'l' },
    { "noded-host",    required_argument, 0, 'n' },
    { "noded-lep",     required_argument, 0, 'e' },
    { "map-file",      required_argument, 0, 'f' },
    { "help",          no_argument,       0, 'h' },
    { "version",       no_argument,       0, 'v' },
    { 0,               0,                 0,  0 }
};

static int readMapFile(Entry* entries, char *fileName) {

    FILE *file=NULL;
    char line[MAX_LINE_LENGTH];
    int numEntries=0;

    file = fopen(fileName, "r");
    if (file == NULL) {
        usys_log_error("Failed to open the status map file: %s", fileName);
        return 0;
    }

    while (fgets(line, sizeof(line), file) != NULL &&
           numEntries < MAX_ENTRIES) {

        if (line[0] != '#') {
            sscanf(line, "%s %s %s %s %s %d",
                   entries[numEntries].serviceName,
                   entries[numEntries].moduleName,
                   entries[numEntries].propertyName,
                   entries[numEntries].type,
                   entries[numEntries].severity,
                   &entries[numEntries].code);
            numEntries++;
        }
    }

    fclose(file);
    return numEntries;
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
    usys_puts("Usage: noded [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                           help menu");
    usys_puts("-l, --logs <TRACE> <DEBUG> <INFO>    log level for the process.\n");
    usys_puts("-n, --noded-host <host>              noded host");
    usys_puts("-e, --noded-ep </node>               API EP at which noded service "
              "will enquire for node info");
    usys_puts("-f, --map-file <file-name>           status map file");
    usys_puts("-v, --version                        version.\n");
}

void init_global_data() {

    gData = (ThreadData *)malloc(sizeof(ThreadData));
    if (gData == NULL) {
        usys_log_error("Unable to allocate memory of size: %d", sizeof(ThreadData));
        usys_exit(1);
    }

    gData->output = DEF_OUTPUT;
    gData->count  = 0;
    pthread_mutex_init(&gData->mutex, NULL);
}

void clean_memory_and_exit(int stage, UInst *service, UInst *admin, Config *config) {

    usys_free(gData);
    usys_free(config->serviceName);
    usys_free(config->nodedHost);
    usys_free(config->nodedEP);
    usys_free(config->remoteServer);

    switch (stage) {
    case NORMAL_EXIT:
        ulfius_stop_framework(service);
        ulfius_clean_instance(service);
        ulfius_stop_framework(admin);
        ulfius_stop_framework(admin);
        usys_free(config->nodeID);
        usys_exit(0);
        break;

    case WEB_ADMIN_FAIL:
        ulfius_stop_framework(service);
        ulfius_clean_instance(service);
        usys_free(config->nodeID);
        usys_exit(1);
        break;

    case WEB_SERVICE_FAIL:
        usys_free(config->nodeID);
        usys_exit(1);

    case NODED_FAIL:
        usys_exit(1);
        break;

    case ADDR_FAIL:
        usys_exit(1);
        break;

    default:
        usys_exit(1);
    }
}

int main(int argc, char **argv) {

    char *debug        = DEF_LOG_LEVEL;
    char *nodedHost    = DEF_NODED_HOST;
    char *nodedEP      = DEF_NODED_EP;
    char *mapFile      = DEF_MAP_FILE;
    UInst serviceInst;
    UInst adminInst;

    Config serviceConfig = {0};

    usys_log_set_service(SERVICE_NAME);
    usys_log_remote_init(SERVICE_NAME);
    init_global_data();

    while (true) {
        int opt = 0;
        int opdIdx = 0;

        opt = getopt_long(argc, argv, "f:l:n:e:hv", longOptions, &opdIdx);
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

        case 'n':
            nodedHost = optarg;
            if (!nodedHost) {
                usage();
                usys_exit(0);
            }
            break;

        case 'e':
            nodedEP = optarg;
            if (!nodedEP) {
                usage();
                usys_exit(0);
            }
            break;

        case 'f':
            mapFile = optarg;
            break;

        default:
            usage();
            usys_exit(0);
        }
    }

    usys_find_ukama_service_address(&serviceConfig.remoteServer);
    if (serviceConfig.remoteServer == NULL) {
        usys_log_error("Ukama not configured in /etc/services");
        clean_memory_and_exit(ADDR_FAIL, &serviceInst, &adminInst, &serviceConfig);
    }

    /* Service config update */
    serviceConfig.serviceName  = usys_strdup(SERVICE_NOTIFY);
    serviceConfig.servicePort  = usys_find_service_port(SERVICE_NOTIFY);
    serviceConfig.adminPort    = usys_find_service_port(SERVICE_NOTIFY_ADMIN);
    serviceConfig.nodedHost    = usys_strdup(nodedHost);
    serviceConfig.nodedPort    = usys_find_service_port(SERVICE_NODE);
    serviceConfig.nodedEP      = usys_strdup(nodedEP);
    serviceConfig.numEntries   = readMapFile(serviceConfig.entries, mapFile);

    if (!serviceConfig.servicePort ||
        !serviceConfig.nodedPort ||
        !serviceConfig.adminPort) {
        usys_log_error("Unable to determine the port for services");
        clean_memory_and_exit(ADDR_FAIL, &serviceInst, &adminInst, &serviceConfig);
    }

    usys_log_debug("Starting notify.d ...");

    signal(SIGINT, handle_sigint);

    /* Read Node Info from noded */
    if (getenv(ENV_NOTIFY_DEBUG_MODE)) {
       serviceConfig.nodeID = strdup(DEF_NODE_ID);
       usys_log_debug("notify.d: Using default Node ID: %s", DEF_NODE_ID);
    } else {
        if (get_nodeid_from_noded(&serviceConfig) == STATUS_NOK) {
            usys_log_error("notify.d: Unable to connect with node.d");
            clean_memory_and_exit(NODED_FAIL, &serviceInst, &adminInst, &serviceConfig);
        }
    }

    if (start_web_services(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup for clients. Exiting.");
        clean_memory_and_exit(WEB_SERVICE_FAIL, &serviceInst, &adminInst, &serviceConfig);
    }

    if (start_admin_web_service(&serviceConfig, &adminInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup for clients. Exiting.");
        clean_memory_and_exit(WEB_SERVICE_FAIL, &serviceInst, &adminInst, &serviceConfig);
    }

    pause();

    usys_log_debug("Exiting $s", SERVICE_NAME);
    clean_memory_and_exit(NORMAL_EXIT, &serviceInst, &adminInst, &serviceConfig);
}
