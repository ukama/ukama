/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "lifecycled.h"
#include "state_store.h"
#include "web_client.h"

#include "usys_log.h"

#define BOOT_ID_FILE "/proc/sys/kernel/random/boot_id"

int64_t lifecycle_boottime_ms(void) {

    struct timespec now;

    if (clock_gettime(CLOCK_BOOTTIME, &now) != 0) return 0;

    return ((int64_t)now.tv_sec * 1000) +
           ((int64_t)now.tv_nsec / 1000000);
}

int64_t lifecycle_epoch_sec(void) {

    return (int64_t)time(NULL);
}

bool lifecycle_read_boot_id(char *buffer, size_t size) {

    FILE *file;
    size_t length;

    if (!buffer || size == 0) return false;

    file = fopen(BOOT_ID_FILE, "r");
    if (!file) return false;

    if (!fgets(buffer, (int)size, file)) {
        fclose(file);
        return false;
    }

    fclose(file);
    length = strlen(buffer);

    while (length > 0 &&
           (buffer[length - 1] == '\n' ||
            buffer[length - 1] == '\r')) {
        buffer[--length] = '\0';
    }

    return length > 0;
}

static void copy_error(char *error,
                       size_t errorSize,
                       const char *message) {

    if (!error || errorSize == 0) return;
    snprintf(error, errorSize, "%s", message ? message : "unknown error");
}

static void enqueue_event_locked(LifecycleContext *ctx) {

    LifecycleEvent *event;
    size_t index;

    if (ctx->eventCount == LIFECYCLED_EVENT_QUEUE) {
        ctx->eventHead =
            (ctx->eventHead + 1) % LIFECYCLED_EVENT_QUEUE;
        ctx->eventCount--;
        usys_log_warn("event: queue full; oldest transition dropped");
    }

    index = (ctx->eventHead + ctx->eventCount) %
        LIFECYCLED_EVENT_QUEUE;
    event = &ctx->events[index];
    memset(event, 0, sizeof(*event));

    event->state = ctx->fsm.state;
    event->sequence = ctx->fsm.sequence;
    event->occurredAt = ctx->fsm.stateSince;
    snprintf(event->reason,
             sizeof(event->reason),
             "%s",
             ctx->fsm.reason);
    ctx->eventCount++;
}

static void persist_locked(LifecycleContext *ctx) {

    if (!state_store_save(ctx->config->stateFile,
                          ctx->bootId,
                          &ctx->fsm)) {
        usys_log_warn("state: unable to save %s",
                      ctx->config->stateFile);
    }
}

static void capture_transition_locked(LifecycleContext *ctx,
                                      uint64_t previousSequence) {

    if (ctx->fsm.sequence == previousSequence) return;

    usys_log_info("state: %s sequence=%llu reason=%s",
                  lifecycle_state_str(ctx->fsm.state),
                  (unsigned long long)ctx->fsm.sequence,
                  ctx->fsm.reason);
    enqueue_event_locked(ctx);
}

static bool event_peek(LifecycleContext *ctx, LifecycleEvent *event) {

    bool present;

    pthread_mutex_lock(&ctx->mutex);
    present = ctx->eventCount > 0;
    if (present && event) {
        *event = ctx->events[ctx->eventHead];
    }
    pthread_mutex_unlock(&ctx->mutex);

    return present;
}

static void event_complete(LifecycleContext *ctx,
                           const LifecycleEvent *event) {

    LifecycleEvent *head;

    pthread_mutex_lock(&ctx->mutex);

    if (ctx->eventCount > 0) {
        head = &ctx->events[ctx->eventHead];
        if (head->sequence == event->sequence &&
            head->state == event->state) {
            ctx->eventHead =
                (ctx->eventHead + 1) % LIFECYCLED_EVENT_QUEUE;
            ctx->eventCount--;
        }
    }

    pthread_mutex_unlock(&ctx->mutex);
}

static bool wait_for_poll(LifecycleContext *ctx) {

    struct timespec deadline;
    int milliseconds;
    bool running;

    clock_gettime(CLOCK_REALTIME, &deadline);
    milliseconds = ctx->config->pollIntervalMs;
    deadline.tv_sec += milliseconds / 1000;
    deadline.tv_nsec += (long)(milliseconds % 1000) * 1000000L;

    if (deadline.tv_nsec >= 1000000000L) {
        deadline.tv_sec++;
        deadline.tv_nsec -= 1000000000L;
    }

    pthread_mutex_lock(&ctx->mutex);
    if (ctx->running) {
        pthread_cond_timedwait(&ctx->condition,
                               &ctx->mutex,
                               &deadline);
    }
    running = ctx->running;
    pthread_mutex_unlock(&ctx->mutex);

    return running;
}

static void poll_and_reduce(LifecycleContext *ctx) {

    StarterSnapshot snapshot;
    LifecycleFsm before;
    uint64_t previousSequence;

    memset(&snapshot, 0, sizeof(snapshot));
    snapshot.aggregate = STARTER_AGGREGATE_UNKNOWN;
    snapshot.configPhase = CONFIG_PHASE_UNKNOWN;

    if (!starter_client_get_status(ctx->config, &snapshot)) {
        snapshot.available = false;
    }

    pthread_mutex_lock(&ctx->mutex);
    ctx->starter = snapshot;
    before = ctx->fsm;
    previousSequence = ctx->fsm.sequence;

    lifecycle_fsm_tick(&ctx->fsm,
                       &ctx->starter,
                       ctx->config->starterUnavailableTimeoutSec,
                       lifecycle_boottime_ms(),
                       lifecycle_epoch_sec());

    capture_transition_locked(ctx, previousSequence);
    if (memcmp(&before, &ctx->fsm, sizeof(before)) != 0) {
        persist_locked(ctx);
    }
    pthread_mutex_unlock(&ctx->mutex);
}

static void *worker_main(void *arg) {

    LifecycleContext *ctx;
    LifecycleEvent event;
    bool running;

    ctx = (LifecycleContext *)arg;
    running = true;

    while (running) {
        if (event_peek(ctx, &event) &&
            notify_client_send_event(ctx->config, &event)) {
            event_complete(ctx, &event);
        }

        poll_and_reduce(ctx);
        running = wait_for_poll(ctx);
    }

    return NULL;
}

bool lifecycle_context_init(LifecycleContext *ctx, Config *config) {

    bool restored;

    if (!ctx || !config) return false;

    memset(ctx, 0, sizeof(*ctx));
    ctx->config = config;

    if (pthread_mutex_init(&ctx->mutex, NULL) != 0) return false;
    if (pthread_cond_init(&ctx->condition, NULL) != 0) {
        pthread_mutex_destroy(&ctx->mutex);
        return false;
    }

    if (!lifecycle_read_boot_id(ctx->bootId, sizeof(ctx->bootId))) {
        snprintf(ctx->bootId,
                 sizeof(ctx->bootId),
                 "fallback-%lld",
                 (long long)lifecycle_epoch_sec());
    }

    lifecycle_fsm_init(&ctx->fsm, lifecycle_epoch_sec());
    restored = state_store_load(config->stateFile,
                                ctx->bootId,
                                &ctx->fsm);

    if (restored) {
        usys_log_info("state: restored %s sequence=%llu",
                      lifecycle_state_str(ctx->fsm.state),
                      (unsigned long long)ctx->fsm.sequence);
    } else {
        persist_locked(ctx);
    }

    enqueue_event_locked(ctx);
    return true;
}

void lifecycle_context_free(LifecycleContext *ctx) {

    if (!ctx) return;

    lifecycle_context_stop(ctx);
    pthread_cond_destroy(&ctx->condition);
    pthread_mutex_destroy(&ctx->mutex);
    memset(ctx, 0, sizeof(*ctx));
}

bool lifecycle_context_start(LifecycleContext *ctx) {

    if (!ctx) return false;

    pthread_mutex_lock(&ctx->mutex);
    if (ctx->running) {
        pthread_mutex_unlock(&ctx->mutex);
        return true;
    }
    ctx->running = true;
    pthread_mutex_unlock(&ctx->mutex);

    if (pthread_create(&ctx->worker, NULL, worker_main, ctx) != 0) {
        pthread_mutex_lock(&ctx->mutex);
        ctx->running = false;
        pthread_mutex_unlock(&ctx->mutex);
        return false;
    }

    ctx->workerStarted = true;
    return true;
}

void lifecycle_context_stop(LifecycleContext *ctx) {

    bool join;

    if (!ctx) return;

    pthread_mutex_lock(&ctx->mutex);
    join = ctx->workerStarted;
    ctx->running = false;
    pthread_cond_broadcast(&ctx->condition);
    pthread_mutex_unlock(&ctx->mutex);

    if (join) {
        pthread_join(ctx->worker, NULL);
        ctx->workerStarted = false;
    }
}

bool lifecycle_context_check_in(LifecycleContext *ctx,
                                const char *bootId,
                                bool bootHealthy,
                                char *error,
                                size_t errorSize) {

    LifecycleState state;
    uint64_t previousSequence;
    bool accepted;

    if (!ctx) return false;

    pthread_mutex_lock(&ctx->mutex);

    if (bootId && *bootId && strcmp(bootId, ctx->bootId) != 0) {
        copy_error(error, errorSize, "bootId does not match current boot");
        pthread_mutex_unlock(&ctx->mutex);
        return false;
    }

    state = ctx->fsm.state;
    accepted = state == LIFECYCLE_STATE_STARTING ||
        state == LIFECYCLE_STATE_CHECKING_IN ||
        state == LIFECYCLE_STATE_READY ||
        state == LIFECYCLE_STATE_CONFIGURING ||
        state == LIFECYCLE_STATE_OPERATIONAL ||
        (state == LIFECYCLE_STATE_FAULTY &&
         ctx->fsm.fault == LIFECYCLE_FAULT_BOOT);

    if (!accepted) {
        copy_error(error,
                   errorSize,
                   "check-in is not valid in current state");
        pthread_mutex_unlock(&ctx->mutex);
        return false;
    }

    previousSequence = ctx->fsm.sequence;
    lifecycle_fsm_begin_check_in(&ctx->fsm,
                                 bootHealthy,
                                 ctx->config->checkInTimeoutSec,
                                 lifecycle_boottime_ms(),
                                 lifecycle_epoch_sec());
    capture_transition_locked(ctx, previousSequence);
    persist_locked(ctx);
    pthread_cond_broadcast(&ctx->condition);
    pthread_mutex_unlock(&ctx->mutex);
    return true;
}

LifecycleConfigureResult lifecycle_context_configure(
    LifecycleContext *ctx,
    const char *requestId,
    const char *assignmentId,
    char *error,
    size_t errorSize) {

    LifecycleConfigureResult result;
    uint64_t previousSequence;

    if (!ctx) return LIFECYCLE_CONFIGURE_INVALID_REQUEST;

    pthread_mutex_lock(&ctx->mutex);

    if (requestId && *requestId && ctx->fsm.requestId[0] != '\0' &&
        strcmp(requestId, ctx->fsm.requestId) == 0) {
        pthread_mutex_unlock(&ctx->mutex);
        return LIFECYCLE_CONFIGURE_DUPLICATE;
    }

    if ((ctx->fsm.state == LIFECYCLE_STATE_READY ||
         ctx->fsm.state == LIFECYCLE_STATE_OPERATIONAL) &&
        (!ctx->starter.available ||
         ctx->starter.aggregate != STARTER_AGGREGATE_READY)) {
        copy_error(error,
                   errorSize,
                   "required applications are not currently ready");
        pthread_mutex_unlock(&ctx->mutex);
        return LIFECYCLE_CONFIGURE_INVALID_STATE;
    }

    previousSequence = ctx->fsm.sequence;
    result = lifecycle_fsm_configure(&ctx->fsm,
                                     requestId,
                                     assignmentId,
                                     ctx->config->configTimeoutSec,
                                     lifecycle_boottime_ms(),
                                     lifecycle_epoch_sec());

    if (result == LIFECYCLE_CONFIGURE_ACCEPTED) {
        capture_transition_locked(ctx, previousSequence);
        persist_locked(ctx);
        pthread_cond_broadcast(&ctx->condition);
    } else if (result == LIFECYCLE_CONFIGURE_BUSY) {
        copy_error(error, errorSize, "another configuration is active");
    } else if (result == LIFECYCLE_CONFIGURE_INVALID_STATE) {
        copy_error(error,
                   errorSize,
                   "configure command is not valid in current state");
    } else if (result == LIFECYCLE_CONFIGURE_INVALID_REQUEST) {
        copy_error(error, errorSize, "requestId is required");
    }

    pthread_mutex_unlock(&ctx->mutex);
    return result;
}

void lifecycle_context_snapshot(LifecycleContext *ctx,
                                LifecycleFsm *fsm,
                                StarterSnapshot *starter) {

    if (!ctx) return;

    pthread_mutex_lock(&ctx->mutex);
    if (fsm) *fsm = ctx->fsm;
    if (starter) *starter = ctx->starter;
    pthread_mutex_unlock(&ctx->mutex);
}
