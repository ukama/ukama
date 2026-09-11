/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <fcntl.h>
#include <jansson.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/file.h>
#include <unistd.h>

#include "state_store.h"

static bool sync_directory(const char *path) {

    int fd;
    int result;

    fd = open(path, O_RDONLY | O_DIRECTORY | O_CLOEXEC);
    if (fd < 0) {
        return false;
    }

    result = fsync(fd);
    if (close(fd) != 0) {
        result = -1;
    }
    return result == 0;
}

static bool prepare_directory(const char *path, char *directory, size_t size) {

    char *cursor;
    char *last;
    char parent[PATH_MAX];
    char *slash;
    bool finished;

    if (path[0] != '/' || snprintf(directory, size, "%s", path) >= (int)size) {
        return false;
    }

    last = strrchr(directory, '/');
    if (last == directory) {
        last[1] = '\0';
        return true;
    }
    *last = '\0';

    for (cursor = directory + 1; ; cursor++) {
        finished = *cursor == '\0';

        if (*cursor != '/' && !finished) {
            continue;
        }
        *cursor = '\0';
        if (mkdir(directory, 0700) == 0) {
            snprintf(parent, sizeof(parent), "%s", directory);
            slash = strrchr(parent, '/');
            if (slash == parent) {
                slash[1] = '\0';
            } else {
                *slash = '\0';
            }
            if (!sync_directory(parent)) {
                return false;
            }
        } else if (errno != EEXIST) {
            return false;
        }
        if (finished) {
            break;
        }
        *cursor = '/';
    }
    return true;
}

bool state_store_open(LifecycleContext *ctx) {

    char directory[PATH_MAX];
    char path[PATH_MAX];

    ctx->stateLockFd = -1;
    if (!prepare_directory(ctx->config->stateFile, directory, sizeof(directory)) ||
        snprintf(path, sizeof(path), "%s.lock", ctx->config->stateFile)
        >= (int)sizeof(path)) {
        return false;
    }

    ctx->stateLockFd = open(path, O_RDWR | O_CREAT | O_CLOEXEC, 0600);
    if (ctx->stateLockFd < 0) {
        return false;
    }
    if (flock(ctx->stateLockFd, LOCK_EX | LOCK_NB) != 0) {
        state_store_close(ctx);
        return false;
    }
    return true;
}

void state_store_close(LifecycleContext *ctx) {

    if (ctx->stateLockFd >= 0) {
        close(ctx->stateLockFd);
        ctx->stateLockFd = -1;
    }
}

static json_t *event_json(const LifecycleEvent *event) {

    return json_pack("{s:i,s:I,s:I,s:s,s:s,s:s,s:s,s:I}",
                     "state", event->state,
                     "sequence", (json_int_t)event->sequence,
                     "occurredAt", (json_int_t)event->occurredAt,
                     "reason", event->reason,
                     "bootId", event->bootId,
                     "requestId", event->requestId,
                     "configMode", event->configMode,
                     "configGeneration", (json_int_t)event->configGeneration);
}

static json_t *checkpoint_json(const LifecycleContext *ctx) {

    const LifecycleFsm *fsm = &ctx->fsm;
    json_t *root;
    json_t *events;
    json_t *event;
    size_t index;
    size_t slot;

    root = json_pack("{s:i,s:s,s:i,s:i,s:i,s:I,s:I,s:I,s:b,s:b,s:b,"
                     "s:s,s:s,s:s,s:I,s:I}",
                     "version", 2,
                     "bootId", ctx->bootId,
                     "state", fsm->state,
                     "faultReturnState", fsm->faultReturnState,
                     "fault", fsm->fault,
                     "sequence", (json_int_t)fsm->sequence,
                     "stateSince", (json_int_t)fsm->stateSince,
                     "checkInDeadlineMs", (json_int_t)fsm->checkInDeadlineMs,
                     "gateOpen", fsm->gateOpen,
                     "configurationSeen", fsm->configurationSeen,
                     "configurationApplied", fsm->configurationApplied,
                     "requestId", fsm->requestId,
                     "configMode", fsm->configMode,
                     "reason", fsm->reason,
                     "configGeneration", (json_int_t)fsm->configGeneration,
                     "confirmedGeneration", (json_int_t)fsm->confirmedGeneration);
    if (root == NULL) {
        return NULL;
    }

    events = json_array();
    if (events == NULL) {
        json_decref(root);
        return NULL;
    }

    for (index = 0; index < ctx->eventCount; index++) {
        slot = (ctx->eventHead + index) % LIFECYCLED_EVENT_QUEUE;
        event = event_json(&ctx->events[slot]);
        if (event == NULL || json_array_append_new(events, event) != 0) {
            json_decref(events);
            json_decref(root);
            return NULL;
        }
    }

    if (json_object_set_new(root, "events", events) != 0) {
        json_decref(root);
        return NULL;
    }
    return root;
}

static bool read_text(json_t *root, const char *key, char *value, size_t size) {

    json_t *entry = json_object_get(root, key);
    const char *text;

    if (!json_is_string(entry)) {
        return false;
    }
    text = json_string_value(entry);
    if (json_string_length(entry) >= size ||
        json_string_length(entry) != strlen(text)) {
        return false;
    }
    snprintf(value, size, "%s", text);
    return true;
}

static bool read_integer(json_t *root, const char *key, int64_t *value) {

    json_t *entry = json_object_get(root, key);

    if (!json_is_integer(entry) || json_integer_value(entry) < 0) {
        return false;
    }
    *value = json_integer_value(entry);
    return true;
}

static bool read_event(json_t *root, LifecycleEvent *event) {

    int64_t state;
    int64_t sequence;
    int64_t generation;

    if (!read_integer(root, "state", &state) || state > LIFECYCLE_STATE_FAULTY ||
        !read_integer(root, "sequence", &sequence) || sequence == 0 ||
        !read_integer(root, "occurredAt", &event->occurredAt) ||
        !read_integer(root, "configGeneration", &generation)) {
        return false;
    }

    event->state = state;
    event->sequence = sequence;
    event->configGeneration = generation;
    return read_text(root, "reason", event->reason, sizeof(event->reason)) &&
        read_text(root, "bootId", event->bootId, sizeof(event->bootId)) &&
        read_text(root, "requestId", event->requestId, sizeof(event->requestId)) &&
        read_text(root, "configMode", event->configMode, sizeof(event->configMode));
}

static bool read_fsm(json_t *root, LifecycleFsm *fsm) {

    int64_t state;
    int64_t returnState;
    int64_t fault;
    int64_t sequence;
    int64_t generation;
    int64_t confirmed;
    json_t *gate = json_object_get(root, "gateOpen");
    json_t *seen = json_object_get(root, "configurationSeen");
    json_t *applied = json_object_get(root, "configurationApplied");

    if (!read_integer(root, "state", &state) || state > LIFECYCLE_STATE_FAULTY ||
        !read_integer(root, "faultReturnState", &returnState) ||
        returnState > LIFECYCLE_STATE_FAULTY ||
        !read_integer(root, "fault", &fault) || fault > LIFECYCLE_FAULT_CONFIGURATION ||
        !read_integer(root, "sequence", &sequence) || sequence == 0 ||
        !read_integer(root, "configGeneration", &generation) ||
        !read_integer(root, "confirmedGeneration", &confirmed) || confirmed > generation ||
        !read_integer(root, "stateSince", &fsm->stateSince) ||
        !read_integer(root, "checkInDeadlineMs", &fsm->checkInDeadlineMs)) {
        return false;
    }
    if (!json_is_boolean(gate) || !json_is_boolean(seen) || !json_is_boolean(applied)) {
        return false;
    }

    fsm->state = state;
    fsm->faultReturnState = returnState;
    fsm->fault = fault;
    fsm->sequence = sequence;
    fsm->configGeneration = generation;
    fsm->confirmedGeneration = confirmed;
    fsm->gateOpen = json_is_true(gate);
    fsm->configurationSeen = json_is_true(seen);
    fsm->configurationApplied = json_is_true(applied);

    return read_text(root, "requestId", fsm->requestId, sizeof(fsm->requestId)) &&
        read_text(root, "configMode", fsm->configMode, sizeof(fsm->configMode)) &&
        read_text(root, "reason", fsm->reason, sizeof(fsm->reason));
}

bool state_store_load(LifecycleContext *ctx) {

    json_t *root;
    json_t *events;
    json_error_t error;
    LifecycleFsm fsm = {0};
    LifecycleEvent pending[LIFECYCLED_EVENT_QUEUE];
    char bootId[LIFECYCLED_ID_LEN];
    int64_t version;
    size_t index;
    size_t count;
    bool valid = false;

    errno = 0;
    root = json_load_file(ctx->config->stateFile, JSON_REJECT_DUPLICATES, &error);
    if (root == NULL) {
        if (errno != ENOENT) {
            errno = EINVAL;
        }
        return false;
    }
    errno = EINVAL;
    if (!read_integer(root, "version", &version) || version != 2 ||
        !read_text(root, "bootId", bootId, sizeof(bootId)) ||
        !read_fsm(root, &fsm)) {
        goto done;
    }

    if (strcmp(bootId, ctx->bootId) != 0) {
        errno = ESTALE;
        goto done;
    }

    events = json_object_get(root, "events");
    count = json_array_size(events);
    if (!json_is_array(events) || count > LIFECYCLED_EVENT_QUEUE) {
        goto done;
    }

    for (index = 0; index < count; index++) {
        if (!read_event(json_array_get(events, index), &pending[index]) ||
            strcmp(pending[index].bootId, bootId) != 0 ||
            pending[index].sequence > fsm.sequence ||
            (index > 0 && pending[index].sequence <= pending[index - 1].sequence)) {
            goto done;
        }
    }

    ctx->fsm = fsm;
    memcpy(ctx->events, pending, count * sizeof(*pending));
    ctx->eventHead = 0;
    ctx->eventCount = count;
    valid = true;
done:
    json_decref(root);
    return valid;
}

bool state_store_save(const LifecycleContext *ctx) {

    char directory[PATH_MAX];
    char temporary[PATH_MAX];
    json_t *root;
    int fd;
    bool saved = false;
    const char *path = ctx->config->stateFile;

    if (!prepare_directory(path, directory, sizeof(directory)) ||
        snprintf(temporary, sizeof(temporary), "%s.tmp.XXXXXX", path)
        >= (int)sizeof(temporary)) {
        return false;
    }

    root = checkpoint_json(ctx);
    if (root == NULL) {
        return false;
    }
    fd = mkstemp(temporary);
    if (fd < 0) {
        json_decref(root);
        return false;
    }

    if (json_dumpfd(root, fd, JSON_COMPACT) == 0 && fsync(fd) == 0) {
        saved = true;
    }
    json_decref(root);
    if (close(fd) != 0) {
        saved = false;
    }
    if (saved && rename(temporary, path) != 0) {
        saved = false;
    }
    if (saved && !sync_directory(directory)) {
        saved = false;
    }
    if (!saved) {
        unlink(temporary);
    }
    return saved;
}
