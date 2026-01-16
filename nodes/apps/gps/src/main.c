/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <pthread.h>

#include "config.h"
#include "gpsd.h"
#include "web_client.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_services.h"

#include "version.h"

/* Global */
GPSData *gData = NULL;

/* network.c */
extern int start_web_service(Config *config, UInst *serviceInst);

void handle_sigint(int signum) {
    usys_log_debug("Terminate signal.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "logs",            required_argument, 0, 'l' },
    { "gps-module-host", required_argument, 0, 'H' },
    { "help",            no_argument, 0, 'h' },
    { "version",         no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

static void set_log_level(char *slevel) {

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

static void usage() {

    usys_puts("Usage: gps.d [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-H, --gps-module-host         Host where GPS is running");
    usys_puts("-v, --version                 Software version");
}

static void init_gps_data() {

    gData = (GPSData *)calloc(1, sizeof(GPSData));
    if (gData == NULL) {
        usys_log_error("Unable to allocate memory of size: %d",
                       sizeof(GPSData));
        usys_exit(0);
    }

    gData->gpsLock   = USYS_FALSE;
    gData->time      = NULL;
    gData->latitude  = NULL;
    gData->longitude = NULL;

    pthread_mutex_init(&gData->mutex, NULL);
}

static void cleanup_gps_data() {

    if (gData == NULL) return;
    
    pthread_mutex_destroy(&gData->mutex);

    usys_free(gData->time);
    usys_free(gData->latitude);
    usys_free(gData->longitude);

    usys_free(gData);
    gData = NULL;
}

int main(int argc, char **argv) {

    int opt, optIdx;

    char *debug     = DEF_LOG_LEVEL;
    char *gpsHost   = DEF_GPS_MODULE_HOST;
    UInst serviceInst;
    Config serviceConfig = {0};

    pthread_t tid = 0;

    usys_log_set_service(SERVICE_NAME);
    //    usys_log_remote_init(SERVICE_NAME);

    init_gps_data();

    while (true) {

        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "vh:l:H", longOptions,
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

        case 'H':
            gpsHost = optarg;
            break;

        default:
            usage();
            usys_exit(0);
        }
    }

    /* Service config update */
    serviceConfig.serviceName  = usys_strdup(SERVICE_NAME);
    serviceConfig.servicePort  = usys_find_service_port(SERVICE_NAME);
    serviceConfig.nodedPort    = usys_find_service_port(SERVICE_NODE);
    serviceConfig.notifydPort  = usys_find_service_port(SERVICE_NOTIFY);
    serviceConfig.nodeID       = NULL;
    serviceConfig.nodeType     = NULL;
    serviceConfig.gpsHost      = strdup(gpsHost);
    serviceConfig.nodedHost    = DEF_NODED_HOST;
    serviceConfig.nodedEP      = DEF_NODED_EP;

    if (!serviceConfig.servicePort ||
        !serviceConfig.nodedPort   ||
        !serviceConfig.notifydPort) {
        usys_log_error("Unable to determine the port for services");
        usys_exit(1);
    }

    usys_log_debug("Starting %s ... ", SERVICE_NAME);

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* Read Node Info from node.d */
    if (getenv(ENV_DEVICED_DEBUG_MODE)) {
        serviceConfig.nodeID   = strdup(DEF_NODE_ID);
        serviceConfig.nodeType = strdup(DEF_NODE_TYPE);
        usys_log_debug("%s: using default Node ID: %s Type: %s",
                       SERVICE_NAME,
                       DEF_NODE_ID,
                       DEF_NODE_TYPE);
    } else {
        if (get_nodeid_from_noded(&serviceConfig) == STATUS_NOK) {
            usys_log_error("Unable to connect with node.d");
            usys_free(serviceConfig.serviceName);
            usys_free(serviceConfig.gpsHost);
            usys_exit(1);
        }
    }

    if (start_web_service(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_free(serviceConfig.serviceName);
        usys_free(serviceConfig.gpsHost);
        usys_log_error("Webservice failed to setup for clients. Exiting.");
        usys_exit(1);
    }

    if (start_gps_data_collection_and_processing(&serviceConfig, &tid) != USYS_TRUE) {
        usys_log_error("Unable to start pthread for GPS. Existing.");

        ulfius_stop_framework(&serviceInst);
        ulfius_clean_instance(&serviceInst);
        free(serviceConfig.serviceName);

        usys_exit(1);
    }

    pause();

    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);
    usys_free(serviceConfig.serviceName);
    usys_free(serviceConfig.gpsHost);

    stop_gps_data_collection_and_processing(tid);
    cleanup_gps_data();

    usys_log_debug("Exiting device.d ...");

    return USYS_TRUE;
}
