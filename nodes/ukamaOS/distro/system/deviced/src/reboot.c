/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <unistd.h>
#include <pthread.h>
#include <sys/reboot.h>

#include "config.h"
#include "deviced.h"
#include "web_client.h"
#include "http_status.h"

static inline bool is_debug_mode(void) {
    return getenv(ENV_DEVICED_DEBUG_MODE) != NULL;
}

static void do_local_reboot(void) {
    if (is_debug_mode()) return;

    sync();
    (void)setuid(0);
    reboot(RB_AUTOBOOT);
}

static bool send_reboot_alarm_and_wait(Config *config,
                                       ThreadArgs *threadArgs) {
    if (wc_send_alarm_to_notifyd(config, &threadArgs->retCode) == USYS_NOK) {
        usys_log_error("Unable to send notification to notify.d");
        return false;
    }

    usys_log_debug("Reboot alarm sent to notify.d %ld", time(NULL));

    /* Always wait to allow notification to go out */
    sleep(WAIT_BEFORE_REBOOT);

    return true;
}

static bool reboot_remote_client_if_tower(Config *config,
                                         ThreadArgs *threadArgs) {

    if (strcmp(config->nodeType, UKAMA_TOWER_NODE) != 0) {
        /* amplifier (and others) have no remote client reboot */
        return true;
    }

    if (wc_send_reboot_to_client(config, &threadArgs->retCode) == USYS_NOK) {
        usys_log_error("Unable to send reboot to client device.d");
        return false;
    }

    return true;
}

void* reboot_node(void *args) {

    ThreadArgs *threadArgs = (ThreadArgs *)args;
    Config     *config     = threadArgs ? threadArgs->config : NULL;

    if (!threadArgs || !config || !config->nodeType) {
        usys_log_error("reboot_node: invalid args/config/nodeType");
        pthread_exit(NULL);
    }

    /* Only tower + amplifier participate in this flow */
    if (strcmp(config->nodeType, UKAMA_TOWER_NODE) != 0 &&
        strcmp(config->nodeType, UKAMA_AMPLIFIER_NODE) != 0) {
        usys_log_debug("reboot_node: ignoring nodeType=%s", config->nodeType);
        pthread_exit(NULL);
    }

    /* 1) alarm + wait (both tower + amplifier) */
    if (!send_reboot_alarm_and_wait(config, threadArgs)) {
        usys_log_error("Reboot not processed");
        pthread_exit(NULL);
    }

    /* 2) tower only: reboot client-mode device.d */
    if (!reboot_remote_client_if_tower(config, threadArgs)) {
        usys_log_error("Reboot not processed");
        pthread_exit(NULL);
    }

    /* 3) reboot this node */
    do_local_reboot();

    pthread_exit(NULL);
}

void process_reboot(Config *config) {

    pthread_t thread;
    ThreadArgs threadArgs = {0};

    threadArgs.config  = config;
    threadArgs.retCode = -1;

    if (config->clientMode) {
        usys_log_debug("Rebooting in client mode");
        if (getenv(ENV_DEVICED_DEBUG_MODE) == NULL) {
            sync();
            setuid(0);
            reboot(RB_AUTOBOOT);
            return;
        }

        return;
    }

    pthread_create(&thread, NULL, reboot_node, (void *)&threadArgs);
    pthread_join(thread, 0);

    usys_log_debug("Process reboot exit with code: %d Str: %s",
                   threadArgs.retCode,
                   HttpStatusStr(threadArgs.retCode));
}
