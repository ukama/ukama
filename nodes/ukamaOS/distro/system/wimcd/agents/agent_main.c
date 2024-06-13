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
#include <uuid/uuid.h>

#include "log.h"
#include "wimc.h"
#include "agent.h"
#include "http_status.h"
#include "agent/network.h"

#include "usys_types.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_getopt.h"

#include "version.h"

#define DEF_LOG_LEVEL          "TRACE"
#define DEF_WIMC_URL           "http://localhost:8087"
#define DEF_AGENT_SERVICE_PORT "8088"

static UsysOption longOptions[] = {
    { "logs",          required_argument, 0, 'l' },
    { "port",          required_argument, 0, 'p' },
    { "wimc",          required_argument, 0, 'w' },
    { "method",        required_argument, 0, 'm' },
    { "help",          no_argument, 0, 'h' },
    { "version",       no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

void usage() {

    printf("Agent: WIMC.d Agent to speak with service provider for contents\n");
    printf("Supported methods: Test, Chunk\n");
    printf("Usage: Agent [options] \n");
    printf("Options:\n");
    printf("--h, --help                         This help menu. \n");
    printf("--w, --wimc                         WIMC port \n");
    printf("--p, --port                         Local listening port. \n");
    printf("--m, --method                       Tx Method <Test | Chunk>.\n");
    printf("--l, --level <ERROR | DEBUG | INFO> Log level for the process. \n");
    printf("--v, --version                      Version. \n");
}

void set_log_level(char *slevel) {

    int ilevel = LOG_TRACE;

    if (!strcmp(slevel, "DEBUG")) {
        ilevel = LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        ilevel = LOG_INFO;
    } else if (!strcmp(slevel, "ERROR")) {
        ilevel = LOG_ERROR;
    }

    log_set_level(ilevel);
}

int main(int argc, char **argv) {

    int opt=0, opdix=0;
    uuid_t uuid;
    long code;
    char *url=DEF_WIMC_URL, *port=DEF_AGENT_SERVICE_PORT, *dbg=DEF_LOG_LEVEL;
    char *method="test";

    char wimcURL[WIMC_MAX_URL_LEN] = {0};
    char servicePort[WIMC_MAX_URL_LEN] = {0};
    char debug[WIMC_MAX_URL_LEN] = {0};
    char agentMethod[WIMC_MAX_URL_LEN] = {0};

    UInst inst;
    
    /* Prase command line args. */
    while (TRUE) {

        opt   = 0;
        opdix = 0;

        opt = usys_getopt_long(argc, argv, "w:p:m:l:hV:", longOptions, &opdix);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'w':
            url = optarg;
            break;

        case 'p':
            port = optarg;
            break;
            
        case 'm':
            method = optarg;
            break;
            
        case 'h':
            usage();
            exit(0);
            break;
            
        case 'l':
            dbg = optarg;
            set_log_level(debug);
            break;
            
        case 'v':
            fprintf(stdout, "Agent - Version: %s\n", VERSION);
            exit(0);
            
        default:
            usage();
            exit(0);
        }
    } /* while */
    
    uuid_clear(uuid);
    uuid_generate(uuid);
    strncpy(&wimcURL[0],     url,    strlen(url));
    strncpy(&servicePort[0], port,   strlen(port));
    strncpy(&debug[0],       dbg,    strlen(dbg));
    strncpy(&agentMethod[0], method, strlen(method));
    
    if (start_web_service(&servicePort[0], &wimcURL[0], &inst) != TRUE) {
        log_error("Failed to start webservice. Exiting.");
        exit(0);
    }

    if (!register_agent_with_wimc(&wimcURL[0],
                                  &servicePort[0],
                                  &agentMethod[0],
                                  uuid)) {
        usys_log_error("Failed to register to wimc.d. Exiting");
        goto cleanup;
    }

    pause();

cleanup:
    ulfius_stop_framework(&inst);
    ulfius_clean_instance(&inst);

    unregister_agent_with_wimc(&wimcURL[0], uuid);
    
    return 1;
}
