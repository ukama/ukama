/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>
#include <pthread.h>

#include "ulfius.h"

#include "usys_api.h"
#include "usys_error.h"
#include "usys_log.h"
#include "usys_file.h"
#include "usys_mem.h"
#include "usys_services.h"

#include "nodeInfo.h"
#include "rlogd.h"

/* network.c */
extern int start_websocket_server(char *nodeID, int port, UInst *serviceInst);

/* Global */
ThreadData *gData = NULL;

static void usage() {

    printf("rlog.d: logging facility \n");
    printf("Usage: rlog.d [options] \n");
    printf("Options:\n");
    printf("--h, --help                         This help menu. \n");
    printf("--l, --level <ERROR | DEBUG | INFO> Log level for the process. \n");
    printf("--v, --version                      Version. \n");
}

void set_log_level(char *slevel) {

    int ilevel = LOG_TRACE;

    if (!strcmp(slevel, "DEBUG")) {
        ilevel = USYS_LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        ilevel = USYS_LOG_INFO;
    } else if (!strcmp(slevel, "ERROR")) {
        ilevel = USYS_LOG_ERROR;
    }

    log_set_level(ilevel);
}

void init_config_and_buffer() {

    gData = (ThreadData *)malloc(sizeof(ThreadData));

    gData->output        = DEF_OUTPUT;
    gData->level         = USYS_LOG_DEBUG;
    gData->flushTime     = DEF_FLUSH_TIME;
    gData->bufferSize    = 0;
    gData->jOutputBuffer = json_pack("{s:[]}", JTAG_LOGS);
    pthread_mutex_init(&gData->bufferMutex, NULL);
}

int main (int argc, char **argv) {

    char *debug=DEF_LOG_LEVEL;
    char *nodeID=NULL;
    int  opt, opdidx;
    int  nodedPort = 0, rlogdPort = 0;
    UInst serviceInst;

    log_set_service(SERVICE_NAME);
    init_config_and_buffer();

    while (USYS_TRUE) {
        opt    = 0;
        opdidx = 0;

        static struct option long_options[] = {
            { "level",   required_argument, 0, 'l'},
            { "help",    no_argument,       0, 'h'},
            { "version", no_argument,       0, 'v'},
            { 0,         0,                 0,  0}
        };

        opt = getopt_long(argc, argv, "l:hv:", long_options, &opdidx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            exit(0);
            break;

        case 'l':
            debug = optarg;
            set_log_level(debug);
            break;

        case 'v':
            fprintf(stdout, "rlog.d - Version: %s\n", RLOGD_VERSION);
            exit(0);

        default:
            usage();
            exit(0);
        }
    }

    nodedPort = usys_find_service_port(SERVICE_NODE);
    rlogdPort = usys_find_service_port(SERVICE_RLOG);

    if (nodedPort == 0 || rlogdPort == 0) {
        usys_log_error("Error getting noded/rlogd port from service db");
        exit(1);
    }

	if (get_nodeID_from_noded(&nodeID, DEF_NODED_HOST, nodedPort) != USYS_TRUE) {
	    usys_log_error("Error retreiving NodeID from noded.d at %s:%d",
                       DEF_NODED_HOST, nodedPort);
		goto done;
	}

    if (start_websocket_server(nodeID, rlogdPort, &serviceInst) != USYS_TRUE){
        usys_log_error("Unable to setup websocket on port: %d", rlogdPort);
        goto done;
    }

    pause();

done:
    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);
	usys_free(nodeID);

	return 0;
}
