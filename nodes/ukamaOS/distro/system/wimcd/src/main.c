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

/* init.c */
extern void open_db(sqlite3 **db, char *dbFile, int flag);

void handle_sigint(int signum) {

    usys_log_debug("Terminate signal.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "logs",          required_argument, 0, 'l' },
    { "dbFile",        required_argument, 0, 'd' },
    { "url",           required_argument, 0, 'u' },
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

    usys_puts("Usage: wimc.d [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-d, --dbFile                  dB file path");
    usys_puts("-u, --url                     Hub URL");
    usys_puts("-v, --version                 Software version");
}

int main (int argc, char **argv) {

    int opt, optIdx;
    
    Agent  *agents = NULL;
    WTasks *tasks  = NULL;
    char   *debug  = DEF_LOG_LEVEL;
    char   *dbFile = DEF_DB_FILE;
    char   hubURL[WIMC_MAX_URL_LEN] = {0};

    UInst  serviceInst;
    Config serviceConfig = {0};

    usys_log_set_service(SERVICE_NAME);
    usys_log_remote_init(SERVICE_NAME);

    if (usys_find_service_port(SERVICE_NAME) == 0) {
        usys_log_error("Unable to find service port for %s", SERVICE_NAME);
        usys_exit(1);
    }

    if (usys_find_service_port(SERVICE_UKAMA) == 0) {
        usys_log_error("Unable to find service port for %s", SERVICE_UKAMA);
        usys_exit(1);
    }

    sprintf(hubURL, "http://localhost:%d",
            usys_find_service_port(SERVICE_UKAMA));

    while (true) {

        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "vh:p:l:d:u", longOptions,
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

        case 'd':
            dbFile = optarg;
            if (!dbFile) {
                usage();
                usys_exit(0);
            }
            break;
          
        case 'l':
            debug = optarg;
            set_log_level(debug);
            break;

        case 'u':
            strcpy(&hubURL[0], optarg);
            if (strlen(hubURL) == 0) {
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
    serviceConfig.servicePort  = usys_find_service_port(SERVICE_NAME);
    serviceConfig.dbFile       = strdup(dbFile);
    serviceConfig.hubURL       = strdup(hubURL);

    if (!serviceConfig.servicePort) {
        usys_log_error("Unable to determine the port for %s", SERVICE_NAME);
        usys_exit(1);
    }

    /* Signal handler */
    signal(SIGINT, handle_sigint);
  
    usys_log_debug("Starting %s ... ", SERVICE_NAME);
  
    agents = (Agent *)calloc(MAX_AGENTS, sizeof(Agent));
    if (!agents) {
        usys_log_error("Memory failure. Exiting");
        exit(1);
    }
    serviceConfig.agents = &agents;

    /*
      tasks = (WTasks *)calloc(1, sizeof(WTasks));
      if (!tasks) {
      log_error("Memory failure. Exiting");
      exit(1);
      }
    */
    serviceConfig.tasks = &tasks;

    /* Step-1: open the local db */
    open_db(&serviceConfig.db, serviceConfig.dbFile, WIMC_FLAG_CREATE_DB);
    if (serviceConfig.db == NULL) {
        usys_log_error("Error creating db at: %s", serviceConfig.dbFile);
        usys_exit(0);
    }

    /* Step-2: setup all endpoints, cb and run webservice */
    if (start_web_service(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup. Exiting");
        usys_exit(0);
    }

    pause();

    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);
    sqlite3_close(serviceConfig.db);

    clear_agents(agents);
    clear_tasks(&tasks);

    free(agents);
  
    return 1;
}
