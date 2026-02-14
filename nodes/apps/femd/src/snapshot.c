/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "snapshot.h"
#include "usys_log.h"

static inline bool valid_unit(FemUnit u) {
    return u == FEM_UNIT_1 || u == FEM_UNIT_2;
}

int snapshot_init(SnapshotStore *store) {

    if (!store) return STATUS_NOK;

    memset(store, 0, sizeof(*store));
    if (pthread_rwlock_init(&store->lock, NULL) != 0) {
        usys_log_error("snapshot: rwlock init failed");
        return STATUS_NOK;
    }

    store->initialized = true;
    return STATUS_OK;
}

void snapshot_cleanup(SnapshotStore *store) {

    if (!store || !store->initialized) return;

    (void)pthread_rwlock_destroy(&store->lock);
    memset(store, 0, sizeof(*store));
}

int snapshot_set_fem_present(SnapshotStore *store, FemUnit unit, bool present, uint32_t tsMs) {

    if (!store || !store->initialized || !valid_unit(unit)) return STATUS_NOK;

    if (pthread_rwlock_wrlock(&store->lock) != 0) return STATUS_NOK;

    store->fem[unit].present    = present;
    store->fem[unit].sampleTsMs = tsMs;

    (void)pthread_rwlock_unlock(&store->lock);
    return STATUS_OK;
}

int snapshot_set_ctrl_present(SnapshotStore *store, bool present, uint32_t tsMs) {

    if (!store || !store->initialized) return STATUS_NOK;

    if (pthread_rwlock_wrlock(&store->lock) != 0) return STATUS_NOK;

    store->ctrl.present    = present;
    store->ctrl.sampleTsMs = tsMs;

    (void)pthread_rwlock_unlock(&store->lock);
    return STATUS_OK;
}

int snapshot_update_fem(SnapshotStore *store, FemUnit unit, const FemSnapshot *in) {

    if (!store || !store->initialized || !in || !valid_unit(unit)) return STATUS_NOK;

    if (pthread_rwlock_wrlock(&store->lock) != 0) return STATUS_NOK;

    store->fem[unit] = *in;

    (void)pthread_rwlock_unlock(&store->lock);
    return STATUS_OK;
}

int snapshot_update_ctrl(SnapshotStore *store, const CtrlSnapshot *in) {

    if (!store || !store->initialized || !in) return STATUS_NOK;

    if (pthread_rwlock_wrlock(&store->lock) != 0) return STATUS_NOK;

    store->ctrl = *in;

    (void)pthread_rwlock_unlock(&store->lock);
    return STATUS_OK;
}

int snapshot_get_fem(SnapshotStore *store, FemUnit unit, FemSnapshot *out) {

    if (!store || !store->initialized || !out || !valid_unit(unit)) return STATUS_NOK;

    if (pthread_rwlock_rdlock(&store->lock) != 0) return STATUS_NOK;

    *out = store->fem[unit];

    (void)pthread_rwlock_unlock(&store->lock);
    return STATUS_OK;
}

int snapshot_get_ctrl(SnapshotStore *store, CtrlSnapshot *out) {

    if (!store || !store->initialized || !out) return STATUS_NOK;

    if (pthread_rwlock_rdlock(&store->lock) != 0) return STATUS_NOK;

    *out = store->ctrl;

    (void)pthread_rwlock_unlock(&store->lock);
    return STATUS_OK;
}


#include <time.h>

uint32_t snapshot_now_ms(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (uint32_t)((uint64_t)ts.tv_sec * 1000ULL + (uint64_t)ts.tv_nsec / 1000000ULL);
}
