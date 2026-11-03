/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <pthread.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>
#include <ulfius.h>

#include "service_api.h"
#include "usys_file.h"
#include "usys_services.h"
#include "version.h"
#include "watchdog.h"

typedef struct {
    ServicePeer peers[SERVICE_MEMBER_COUNT];
    char stateUrl[SERVICE_URL_LEN + DEVICE_PATH_LEN];
    char reportUrl[SERVICE_URL_LEN + DEVICE_PATH_LEN];
    int port;
} WatchdogConfig;

static volatile sig_atomic_t stopping;
static pthread_mutex_t statusLock = PTHREAD_MUTEX_INITIALIZER;
static uint64_t lastReport;

static uint64_t watchdog_clock_msec(void) {

    struct timespec now;

    if (clock_gettime(CLOCK_MONOTONIC, &now) != 0) {
        return 0;
    }

    return (uint64_t)now.tv_sec * CLOCK_MSEC_PER_SEC +
           now.tv_nsec / CLOCK_NSEC_PER_MSEC;
}

static void stop(int signalNumber) {

    (void)signalNumber;
    stopping = 1;
}

static bool load_configuration(WatchdogConfig *config) {

    const char *base;
    char defaultBase[SERVICE_URL_LEN];
    int port;

    if (!service_peers_init(config->peers)) {
        return false;
    }

    base = getenv(ENV_DEVICE_URL);

    if (base == NULL) {
        port = usys_find_service_port(SERVICE_DEVICE);

        if (port < MIN_SERVICE_PORT || port > MAX_SERVICE_PORT) {
            fprintf(stderr, "Cannot resolve service port: %s\n",
                    SERVICE_DEVICE);

            return false;
        }

        snprintf(defaultBase, sizeof(defaultBase), DEF_DEVICE_URL_FORMAT,
                 port);
        base = defaultBase;
    }

    if (strncmp(base, HTTP_SCHEME, sizeof(HTTP_SCHEME) - 1) != 0 ||
        strlen(base) >= sizeof(config->stateUrl) - DEVICE_PATH_LEN) {
        return false;
    }

    port = usys_find_service_port(SERVICE_NAME);

    if (port < MIN_SERVICE_PORT || port > MAX_SERVICE_PORT) {
        fprintf(stderr, "Cannot resolve service port: %s\n", SERVICE_NAME);

        return false;
    }

    config->port = port;
    snprintf(config->stateUrl, sizeof(config->stateUrl), "%s%s", base,
             DEVICE_API_STATE);
    snprintf(config->reportUrl, sizeof(config->reportUrl), "%s%s", base,
             DEVICE_API_WATCHDOG);

    return true;
}

static int web_service_cb_ping(const struct _u_request *request,
                               struct _u_response *response, void *data) {

    (void)request;
    (void)data;

    ulfius_set_empty_body_response(response, HttpStatus_OK);

    return U_CALLBACK_CONTINUE;
}

static int web_service_cb_version(const struct _u_request *request,
                                  struct _u_response *response, void *data) {

    (void)request;
    (void)data;

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);

    return U_CALLBACK_CONTINUE;
}

static int web_service_cb_ready(const struct _u_request *request,
                                struct _u_response *response, void *data) {

    bool ready;
    json_t *body;
    int status;

    (void)request;
    (void)data;

    pthread_mutex_lock(&statusLock);
    ready = lastReport != 0 &&
            watchdog_clock_msec() - lastReport < WATCHDOG_REPORT_MSEC;
    pthread_mutex_unlock(&statusLock);

    status = ready ? HttpStatus_OK : HttpStatus_ServiceUnavailable;
    body   = json_pack("{s:b,s:s}", JTAG_READY, ready, JTAG_REASON,
                     ready ? "ready" : "controller unavailable");
    ulfius_set_json_body_response(response, status, body);
    json_decref(body);

    return U_CALLBACK_CONTINUE;
}

static bool start_web_service(struct _u_instance *server, int port) {

    if (ulfius_init_instance(server, port, NULL, NULL) != U_OK) {
        return false;
    }

    ulfius_add_endpoint_by_val(server, HTTP_GET, API_PREFIX, API_EP_PING, 0,
                               web_service_cb_ping, NULL);
    ulfius_add_endpoint_by_val(server, HTTP_GET, API_PREFIX, API_EP_READY, 0,
                               web_service_cb_ready, NULL);
    ulfius_add_endpoint_by_val(server, HTTP_GET, API_PREFIX, API_EP_VERSION, 0,
                               web_service_cb_version, NULL);

    if (ulfius_start_framework(server) != U_OK) {
        ulfius_clean_instance(server);

        return false;
    }

    return true;
}

static void report_group_state(const WatchdogConfig *config) {

    ServiceObservation members[SERVICE_MEMBER_COUNT];
    json_t *state = NULL;
    json_t *report;
    int code;
    size_t i;

    code = service_http(HTTP_GET, config->stateUrl, NULL, &state);

    if (code != HttpStatus_OK) {
        goto done;
    }

    for (i = 0; i < SERVICE_MEMBER_COUNT; i++) {
        members[i] = service_peer_read(&config->peers[i]);
    }

    report = watchdog_observation(json_object_get(state, JTAG_SERVICE_GROUP),
                                  members);

    if (report == NULL) {
        goto done;
    }

    code = service_http(HTTP_POST, config->reportUrl, report, NULL);

    if (code == HttpStatus_Accepted) {
        pthread_mutex_lock(&statusLock);
        lastReport = watchdog_clock_msec();
        pthread_mutex_unlock(&statusLock);
    }

    json_decref(report);

done:
    json_decref(state);
}

int main(int argc, char **argv) {

    WatchdogConfig config;
    struct _u_instance server;

    if (argc == 2 && (strcmp(argv[1], ARG_VERSION_LONG) == 0 ||
                      strcmp(argv[1], ARG_VERSION_SHORT) == 0)) {
        puts(VERSION);

        return EXIT_SUCCESS;
    }

    if (argc != 1 || !load_configuration(&config)) {
        return EXIT_FAILURE;
    }

    signal(SIGTERM, stop);
    signal(SIGINT, stop);

    if (!start_web_service(&server, config.port)) {
        return EXIT_FAILURE;
    }

    while (!stopping) {
        report_group_state(&config);
        sleep(WATCHDOG_POLL_SEC);
    }

    ulfius_stop_framework(&server);
    ulfius_clean_instance(&server);

    return EXIT_SUCCESS;
}
