/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <sqlite3.h>
#include <getopt.h>
#include <ulfius.h>
#include <curl/curl.h>
#include <signal.h>
#include <pthread.h>

#include "db.h"
#include "log.h"
#include "wimc.h"
#include "agent.h"
#include "common/utils.h"
#include "network.h"
#include "tasks.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"

#include "version.h"

static volatile sig_atomic_t gTerminate = 0;

static void handle_sigint(int signum) {

    (void)signum;
    gTerminate = 1;
}

static UsysOption longOptions[] = {
    { "logs",    required_argument, 0, 'l' },
    { "dbFile",  required_argument, 0, 'd' },
    { "url",     required_argument, 0, 'u' },
    { "help",    no_argument,       0, 'h' },
    { "version", no_argument,       0, 'v' },
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

static void usage(void) {

    usys_puts("Usage: wimc.d [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-d, --dbFile                  DB file path");
    usys_puts("-u, --url                     Hub URL");
    usys_puts("-v, --version                 Software version");
}

int main(int argc, char **argv) {

    int opt, optIdx;
    int rc = EXIT_FAILURE;
    int taskMutexInit = 0;
    int dbMutexInit   = 0;
    int curlInit      = 0;
    int webStarted    = 0;

    Agent  *agents = NULL;
    WTasks *tasks  = NULL;
    char   *debug  = DEF_LOG_LEVEL;
    char   *dbFile = NULL;
    char    hubURL[WIMC_MAX_URL_LEN] = {0};

    UInst  serviceInst;
    Config serviceConfig;

    memset(&serviceInst, 0, sizeof(serviceInst));
    memset(&serviceConfig, 0, sizeof(serviceConfig));

    usys_log_set_service(SERVICE_NAME);
    usys_log_remote_init(SERVICE_NAME);

    if (usys_find_service_port(SERVICE_NAME) == 0) {
        usys_log_error("Unable to find service port for %s", SERVICE_NAME);
        goto cleanup;
    }

    if (usys_find_service_port(SERVICE_UKAMA) == 0) {
        usys_log_error("Unable to find service port for %s", SERVICE_UKAMA);
        goto cleanup;
    }

    snprintf(hubURL, sizeof(hubURL), "http://localhost:%d",
             usys_find_service_port(SERVICE_UKAMA));

    dbFile = DEF_DB_FILE;

    while (USYS_TRUE) {
        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "hvl:d:u:", longOptions, &optIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            rc = EXIT_SUCCESS;
            goto cleanup;

        case 'v':
            usys_puts(VERSION);
            rc = EXIT_SUCCESS;
            goto cleanup;

        case 'd':
            if (optarg == NULL || *optarg == '\0') {
                usage();
                goto cleanup;
            }
            dbFile = optarg;
            break;

        case 'l':
            if (optarg != NULL && *optarg != '\0') {
                debug = optarg;
                set_log_level(debug);
            }
            break;

        case 'u':
            if (optarg == NULL || *optarg == '\0') {
                usage();
                goto cleanup;
            }

            snprintf(hubURL, sizeof(hubURL), "%s", optarg);
            if (strlen(hubURL) >= sizeof(hubURL)) {
                usys_log_error("Hub URL too long");
                goto cleanup;
            }
            break;

        default:
            usage();
            goto cleanup;
        }
    }

    serviceConfig.servicePort = usys_find_service_port(SERVICE_NAME);
    serviceConfig.dbFile      = strdup(dbFile ? dbFile : WIMC_DB_PATH);
    serviceConfig.hubURL      = strdup(hubURL);

    if (serviceConfig.dbFile == NULL || serviceConfig.hubURL == NULL) {
        usys_log_error("Memory allocation failure");
        goto cleanup;
    }

    if (!serviceConfig.servicePort) {
        usys_log_error("Unable to determine the port for %s", SERVICE_NAME);
        goto cleanup;
    }

    signal(SIGINT,  handle_sigint);
    signal(SIGTERM, handle_sigint);

    usys_log_debug("Starting %s ... ", SERVICE_NAME);

    agents = (Agent *)calloc(MAX_AGENTS, sizeof(Agent));
    if (agents == NULL) {
        usys_log_error("Memory failure. Exiting");
        goto cleanup;
    }

    serviceConfig.agents = &agents;
    serviceConfig.tasks  = &tasks;

    if (db_open_or_create(serviceConfig.dbFile, &serviceConfig.db) != 0) {
        usys_log_error("Unable to open/create DB file: %s",
                       serviceConfig.dbFile);
        goto cleanup;
    }

    if (pthread_mutex_init(&serviceConfig.taskMutex, NULL) != 0) {
        usys_log_error("taskMutex init failed");
        goto cleanup;
    }
    taskMutexInit = 1;

    if (pthread_mutex_init(&serviceConfig.dbMutex, NULL) != 0) {
        usys_log_error("dbMutex init failed");
        goto cleanup;
    }
    dbMutexInit = 1;

    if (curl_global_init(CURL_GLOBAL_ALL) != 0) {
        usys_log_error("curl_global_init failed");
        goto cleanup;
    }
    curlInit = 1;

    db_mark_old_downloads_failed(serviceConfig.db);

    if (start_web_service(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup. Exiting");
        goto cleanup;
    }
    webStarted = 1;

    while (!gTerminate) {
        pause();
    }

    rc = EXIT_SUCCESS;

cleanup:
    if (webStarted) {
        ulfius_stop_framework(&serviceInst);
        ulfius_clean_instance(&serviceInst);
    }

    if (serviceConfig.db != NULL) {
        sqlite3_close(serviceConfig.db);
        serviceConfig.db = NULL;
    }

    if (dbMutexInit) {
        pthread_mutex_destroy(&serviceConfig.dbMutex);
    }

    if (taskMutexInit) {
        pthread_mutex_destroy(&serviceConfig.taskMutex);
    }

    if (curlInit) {
        curl_global_cleanup();
    }

    clear_tasks(&tasks);
    clear_agents(agents);

    if (agents != NULL) {
        free(agents);
        agents = NULL;
    }

    if (serviceConfig.dbFile != NULL) {
        free(serviceConfig.dbFile);
        serviceConfig.dbFile = NULL;
    }

    if (serviceConfig.hubURL != NULL) {
        free(serviceConfig.hubURL);
        serviceConfig.hubURL = NULL;
    }

    return rc;
}
