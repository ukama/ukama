/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include <curl/curl.h>

#include "data_plane.h"
#include "epcemu.h"
#include "init_network.h"
#include "pcrf.h"
#include "services.h"
#include "ue.h"
#include "web_service.h"

#include "usys_getopt.h"
#include "usys_log.h"

#include "version.h"

static volatile bool gTerminate = false;
DataPlane gDataPlane;

static UsysOption longOptions[] = {
    { "logs",    required_argument, 0, 'l' },
    { "help",    no_argument,       0, 'h' },
    { "version", no_argument,       0, 'v' },
    { 0, 0, 0, 0 }
};

static void handle_sigint(int signum) {

    (void)signum;

    usys_log_debug("Terminate signal");
    gTerminate = true;
}

static void set_log_level(char *slevel) {

    int ilevel = USYS_LOG_TRACE;

    if (slevel == NULL) return;

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

    printf("Usage: epcemu.d [options]\n");
    printf("Options:\n");
    printf("-h, --help                    Help menu\n");
    printf("-l, --logs <TRACE|DEBUG|INFO> Log level for the process\n");
    printf("-v, --version                 Software version\n");
}

static int detach_cb(const UeEntry *ue, void *arg) {

    EpcemuConfig *config;

    config = (EpcemuConfig *)arg;
    if (ue == NULL || config == NULL) return USYS_FALSE;

    usys_log_debug("detaching UE on shutdown imsi=%s ip=%s",
                   ue->imsi,
                   ue->ip);

    if (!pcrf_delete_session(config, ue->imsi)) {
        usys_log_error("PCRF session delete failed during shutdown imsi=%s",
                       ue->imsi);
    }

    ue_detach_complete(ue->imsi);
    return USYS_TRUE;
}

static void detach_all_ues(EpcemuConfig *config) {

    if (config == NULL) return;

    ue_for_each_attached(detach_cb, config);
}

int main(int argc, char **argv) {

    int opt;
    int optIdx;
    char *debug;
    UInst serviceInst;
    EpcemuConfig config;
    EpcemuStatus status;
    ServiceContext ctx;

    debug = EPCEMU_DEF_LOG_LEVEL;

    memset(&config,     0, sizeof(config));
    memset(&status,     0, sizeof(status));
    memset(&ctx,        0, sizeof(ctx));
    memset(&gDataPlane, 0, sizeof(gDataPlane));

    usys_log_set_service(EPCEMU_SERVICE_NAME);

    while (true) {
        opt    = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "l:hv", longOptions, &optIdx);
        if (opt == -1) break;

        switch (opt) {
        case 'h':
            usage();
            return 0;

        case 'v':
            printf("%s\n", VERSION);
            return 0;

        case 'l':
            debug = optarg;
            break;

        default:
            usage();
            return 1;
        }
    }

    set_log_level(debug);

    signal(SIGINT,  handle_sigint);
    signal(SIGTERM, handle_sigint);

    curl_global_init(CURL_GLOBAL_DEFAULT);
    ue_store_init();
    status_init(&status);

    ctx.config = &config;
    ctx.status = &status;

    usys_log_debug("Starting %s", EPCEMU_SERVICE_NAME);

    if (!services_resolve(&config, &status)) {
        goto failed;
    }

    if (!init_network_probe(&config, &status)) {
        goto failed;
    }

    (void)pcrf_probe(&config, &status);

    if (!data_plane_start(&gDataPlane, &config, &status)) {
        goto failed;
    }

    if (!init_network_reconcile(&config, &status)) {
        goto failed;
    }

    if (start_web_service(&ctx, &serviceInst) != USYS_TRUE) {
        status_fail(&status, "failed to start web service");
        goto failed;
    }

    status_set(&status, EpcemuStateReady, "none");

    while (!gTerminate) {
        pause();
    }

    detach_all_ues(&config);

    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);
    data_plane_stop(&gDataPlane);
    status_destroy(&status);
    ue_store_destroy();
    curl_global_cleanup();

    usys_log_debug("Exiting %s", EPCEMU_SERVICE_NAME);
    return 0;

failed:
    usys_log_error("%s failed to start", EPCEMU_SERVICE_NAME);

    data_plane_stop(&gDataPlane);
    status_destroy(&status);
    ue_store_destroy();
    curl_global_cleanup();

    return 1;
}
