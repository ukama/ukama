/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <getopt.h>
#include <signal.h>
#include <sqlite3.h>
#include <stdlib.h>
#include <string.h>
#include <ulfius.h>
#include <unistd.h>

#include "agent.h"
#include "agent/network.h"
#include "http_status.h"
#include "log.h"
#include "wimc.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_services.h"
#include "usys_types.h"

#include "version.h"

bool start_web_service(char *method, struct _u_instance *webInstance);

static volatile sig_atomic_t gTerminate = 0;

static UsysOption longOptions[] = {
    { "logs",    required_argument, 0, 'l' },
    { "method",  required_argument, 0, 'm' },
    { "help",    no_argument,       0, 'h' },
    { "version", no_argument,       0, 'v' },
    { 0, 0, 0, 0 }
};

static void handle_signal(int signum) {

    (void)signum;
    gTerminate = 1;
}

static void usage(void) {

    printf("Agent: WIMC agent for package fetches\n");
    printf("Supported methods: test, chunk, tar.gz\n");
    printf("Usage: agent [options]\n");
    printf("Options:\n");
    printf("-h, --help                         Help menu\n");
    printf("-m, --method <test|chunk|tar.gz>   Transfer method\n");
    printf("-l, --logs <TRACE|DEBUG|INFO>      Log level\n");
    printf("-v, --version                      Version\n");
}

static void set_log_level(char *slevel) {

    int ilevel;

    ilevel = USYS_LOG_TRACE;

    if (slevel == NULL) {
        return;
    }

    if (!strcmp(slevel, "DEBUG")) {
        ilevel = USYS_LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        ilevel = USYS_LOG_INFO;
    } else if (!strcmp(slevel, "ERROR")) {
        ilevel = USYS_LOG_ERROR;
    }

    usys_log_set_level(ilevel);
}

int main(int argc, char **argv) {

    int opt;
    int optIdx;
    int rc;
    char *method;
    UInst inst;

    rc = EXIT_FAILURE;
    method = WIMC_METHOD_TEST_STR;
    memset(&inst, 0, sizeof(inst));

    usys_log_set_service(SERVICE_WIMC_AGENT);
    usys_log_remote_init(SERVICE_WIMC_AGENT);

    if (usys_find_service_port(SERVICE_WIMC) == 0) {
        usys_log_error("Unable to find service port for %s", SERVICE_WIMC);
        return EXIT_FAILURE;
    }

    while (TRUE) {
        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "m:l:hv", longOptions, &optIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'm':
            method = optarg;
            break;

        case 'h':
            usage();
            return EXIT_SUCCESS;

        case 'l':
            set_log_level(optarg);
            break;

        case 'v':
            fprintf(stdout, "%s - Version: %s\n", SERVICE_WIMC_AGENT,
                    VERSION);
            return EXIT_SUCCESS;

        default:
            usage();
            return EXIT_FAILURE;
        }
    }

    signal(SIGINT, handle_signal);
    signal(SIGTERM, handle_signal);

    if (start_web_service(method, &inst) != TRUE) {
        usys_log_error("Failed to start webservice");
        return EXIT_FAILURE;
    }

    while (!gTerminate) {
        pause();
    }

    ulfius_stop_framework(&inst);
    ulfius_clean_instance(&inst);

    rc = EXIT_SUCCESS;
    return rc;
}
