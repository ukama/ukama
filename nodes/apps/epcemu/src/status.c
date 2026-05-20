/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "epcemu.h"

static void copy_str(char *dst, size_t size, const char *src) {

    if (dst == NULL || size == 0) return;

    if (src == NULL) {
        dst[0] = '\0';
        return;
    }

    snprintf(dst, size, "%s", src);
}

const char *status_state_str(EpcemuState state) {

    switch (state) {
    case EpcemuStateStarting:            return "starting";
    case EpcemuStateResolvingServices:   return "resolving-services";
    case EpcemuStateCheckingInitNetwork: return "checking-init-network";
    case EpcemuStateCheckingPcrf:        return "checking-pcrf";
    case EpcemuStateReady:               return "ready";
    case EpcemuStateFailed:              return "failed";
    default:                             return "unknown";
    }
}

void status_init(EpcemuStatus *status) {

    if (status == NULL) return;

    memset(status, 0, sizeof(EpcemuStatus));
    pthread_mutex_init(&status->mutex, NULL);

    status->state = EpcemuStateStarting;
    status->ready = false;
    copy_str(status->reason, sizeof(status->reason), "starting");
}

void status_destroy(EpcemuStatus *status) {

    if (status == NULL) return;

    pthread_mutex_destroy(&status->mutex);
}

void status_set(EpcemuStatus *status, EpcemuState state,
                const char *reason) {

    if (status == NULL) return;

    pthread_mutex_lock(&status->mutex);

    status->state = state;
    status->ready = (state == EpcemuStateReady);
    if (reason != NULL) {
        copy_str(status->reason, sizeof(status->reason), reason);
    }

    pthread_mutex_unlock(&status->mutex);
}

void status_fail(EpcemuStatus *status, const char *reason) {

    status_set(status, EpcemuStateFailed, reason);
}

bool status_is_ready(EpcemuStatus *status) {

    bool ready;

    if (status == NULL) return false;

    pthread_mutex_lock(&status->mutex);
    ready = status->ready;
    pthread_mutex_unlock(&status->mutex);

    return ready;
}
