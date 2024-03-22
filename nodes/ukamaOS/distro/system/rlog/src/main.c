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

#include "ulfius.h"

#include "usys_api.h"
#include "usys_error.h"
#include "usys_log.h"
#include "usys_file.h"
#include "usys_services.h"

#include "nodeInfo.h"
#include "rlogd.h"

/* network.c */
extern int start_web_service(int port, UInst *serviceInst);

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
        ilevel = LOG_INFO;
    } else if (!strcmp(slevel, "ERROR")) {
        ilevel = LOG_ERROR;
    }

    log_set_level(ilevel);
}

int main (int argc, char **argv) {

    char *debug=DEF_LOG_LEVEL;
    char *nodeID=NULL;
    int  opt, opdidx;
    int  nodedPort = 0, rlogdPort = 0;
    UInst serviceInst;

    log_set_service(SERVICE_NAME);

    while (1) {
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

    /* find node-ID from node.d */
	if (get_nodeID_from_noded(&nodeID, DEF_NODED_HOST, nodedPort) != USYS_TRUE) {
	    usys_log_error("Error retreiving NodeID from noded.d at %s:%s",
                       DEF_NODED_HOST, nodedPort);
		goto done;
	}

    /* start web-service */
    if (start_web_service(nodedPort, &serviceInst) != USYS_TRUE){

    }

done:
	free(nodeID);

	return 0;
}
