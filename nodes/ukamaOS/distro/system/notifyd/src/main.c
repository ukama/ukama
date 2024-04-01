/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include "config.h"
#include "notify_macros.h"
#include "service.h"
#include "web.h"
#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_file.h"
#include "usys_services.h"

/**
 * @fn      void handle_sigint(int)
 * @brief   Handle terminate signal for Noded
 *
 * @param   signum
 */
void handle_sigint(int signum) {
    usys_log_debug("Caught terminate signal.\n");
    usys_log_debug("Cleanup complete.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "logs",          required_argument, 0, 'l' },
    { "noded-host",    required_argument, 0, 'n' },
    { "noded-lep",     required_argument, 0, 'e' },
    { "status-file",   required_argument, 0, 'f' },
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
 * @fn      void usage()
 * @brief   Usage options for the ukamaEDR
 *
 */
void usage() {
    usys_puts("Usage: noded [options] \n");
    usys_puts("Options:\n");
    usys_puts(
        "-h, --help                             Help menu.\n");
    usys_puts(
        "-l, --logs <TRACE> <DEBUG> <INFO>      Log level for the process.\n");
    usys_puts(
        "-n, --noded-host <host>               Host at which noded service"
                  "will listen.\n");
    usys_puts(
        "-e, --noded-ep </node>                API EP at which noded service"
                       "will enquire for node info.\n");
    usys_puts(
        "-f, --map-file <file-name>         Status map file\n");

    usys_puts(
        "-v, --version                          Software Version.\n");
}

int main(int argc, char **argv) {

    int ret = USYS_OK;

    char *debug        = DEF_LOG_LEVEL;
    char *nodedHost    = DEF_NODED_HOST;
    char *nodedEP      = DEF_NODED_EP;
    char *mapFile      = DEF_MAP_FILE;
    UInst serviceInst;

    Config serviceConfig = {0};

    log_set_service(SERVICE_NAME);

    /* Parsing command line args. */
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
            usys_puts(NOTIFY_VERSION);
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
        usys_exit(1);
    }

    /* Service config update */
    serviceConfig.serviceName  = usys_strdup(SERVICE_NAME);
    serviceConfig.servicePort  = usys_find_service_port(SERVICE_NAME);
    serviceConfig.nodedHost    = usys_strdup(nodedHost);
    serviceConfig.nodedPort    = usys_find_service_port(SERVICE_NODE);
    serviceConfig.nodedEP      = usys_strdup(nodedEP);
    serviceConfig.numEntries   = readMapFile(serviceConfig.entries, mapFile);

    if (!serviceConfig.servicePort ||
        !serviceConfig.nodedPort) {
        usys_log_error("Unable to determine the port for services");
        usys_exit(1);
    }

    usys_log_debug("Starting notify.d ...");

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* Read Node Info from noded */
    if (getenv(ENV_NOTIFY_DEBUG_MODE)) {
       serviceConfig.nodeID = strdup(DEF_NODE_ID);
       usys_log_debug("notify.d: Using default Node ID: %s", DEF_NODE_ID);
    } else {
        if (get_nodeid_from_noded(&serviceConfig) == STATUS_NOK) {
            usys_log_error("notify.d: Unable to connect with node.d");
            goto done;
        }
    }

    if (start_web_services(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup for clients. Exiting.");
        exit(1);
    }

    pause();

done:
    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);

    free(serviceConfig.serviceName);
    free(serviceConfig.nodedHost);
    free(serviceConfig.nodedEP);
    free(serviceConfig.remoteServer);

    usys_log_debug("Exiting notify.d ...");
    return 1;
}
