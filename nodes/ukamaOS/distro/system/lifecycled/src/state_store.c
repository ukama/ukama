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
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#include "state_store.h"

#define STORE_LINE_LEN 320
#define STORE_PATH_LEN 1024

static void trim_newline(char *value) {

    size_t length;

    if (!value) return;
    length = strlen(value);

    while (length > 0 &&
           (value[length - 1] == '\n' || value[length - 1] == '\r')) {
        value[--length] = '\0';
    }
}

static bool mkdir_parent(const char *path) {

    char copy[STORE_PATH_LEN];
    char *cursor;

    if (!path || strlen(path) >= sizeof(copy)) return false;

    snprintf(copy, sizeof(copy), "%s", path);
    cursor = strrchr(copy, '/');
    if (!cursor) return true;
    if (cursor == copy) return true;
    *cursor = '\0';

    for (cursor = copy + 1; *cursor; cursor++) {
        if (*cursor != '/') continue;
        *cursor = '\0';
        if (mkdir(copy, 0755) != 0 && errno != EEXIST) return false;
        *cursor = '/';
    }

    return mkdir(copy, 0755) == 0 || errno == EEXIST;
}

static void copy_value(char *dst, size_t size, const char *src) {

    if (!dst || size == 0) return;
    snprintf(dst, size, "%s", src ? src : "");
}

bool state_store_load(const char *path,
                      const char *bootId,
                      LifecycleFsm *fsm) {

    FILE *file;
    LifecycleFsm loaded;
    char storedBootId[LIFECYCLE_ID_LEN];
    char line[STORE_LINE_LEN];
    char *separator;
    char *key;
    char *value;
    LifecycleState state;

    if (!path || !bootId || !fsm) return false;

    file = fopen(path, "r");
    if (!file) return false;

    memset(&loaded, 0, sizeof(loaded));
    memset(storedBootId, 0, sizeof(storedBootId));

    while (fgets(line, sizeof(line), file)) {
        trim_newline(line);
        separator = strchr(line, '=');
        if (!separator) continue;

        *separator = '\0';
        key = line;
        value = separator + 1;

        if (strcmp(key, "boot_id") == 0) {
            copy_value(storedBootId, sizeof(storedBootId), value);
        } else if (strcmp(key, "state") == 0) {
            if (lifecycle_state_parse(value, &state)) loaded.state = state;
        } else if (strcmp(key, "fault_return_state") == 0) {
            if (lifecycle_state_parse(value, &state)) {
                loaded.faultReturnState = state;
            }
        } else if (strcmp(key, "fault") == 0) {
            loaded.fault = (LifecycleFault)strtol(value, NULL, 10);
        } else if (strcmp(key, "sequence") == 0) {
            loaded.sequence = strtoull(value, NULL, 10);
        } else if (strcmp(key, "state_since") == 0) {
            loaded.stateSince = strtoll(value, NULL, 10);
        } else if (strcmp(key, "check_in_deadline_ms") == 0) {
            loaded.checkInDeadlineMs = strtoll(value, NULL, 10);
        } else if (strcmp(key, "config_deadline_ms") == 0) {
            loaded.configDeadlineMs = strtoll(value, NULL, 10);
        } else if (strcmp(key, "gate_open") == 0) {
            loaded.gateOpen = strtol(value, NULL, 10) != 0;
        } else if (strcmp(key, "configuration_seen") == 0) {
            loaded.configurationSeen = strtol(value, NULL, 10) != 0;
        } else if (strcmp(key, "configuration_applied") == 0) {
            loaded.configurationApplied = strtol(value, NULL, 10) != 0;
        } else if (strcmp(key, "request_id") == 0) {
            copy_value(loaded.requestId, sizeof(loaded.requestId), value);
        } else if (strcmp(key, "assignment_id") == 0) {
            copy_value(loaded.assignmentId,
                       sizeof(loaded.assignmentId),
                       value);
        } else if (strcmp(key, "reason") == 0) {
            copy_value(loaded.reason, sizeof(loaded.reason), value);
        }
    }

    fclose(file);

    if (strcmp(storedBootId, bootId) != 0 || loaded.sequence == 0) {
        return false;
    }

    *fsm = loaded;
    return true;
}

bool state_store_save(const char *path,
                      const char *bootId,
                      const LifecycleFsm *fsm) {

    FILE *file;
    char temporary[STORE_PATH_LEN];
    bool ok;

    if (!path || !bootId || !fsm || strlen(path) + 5 >= sizeof(temporary)) {
        return false;
    }

    if (!mkdir_parent(path)) return false;

    snprintf(temporary, sizeof(temporary), "%s.tmp", path);
    file = fopen(temporary, "w");
    if (!file) return false;

    ok = fprintf(file,
                 "boot_id=%s\n"
                 "state=%s\n"
                 "fault_return_state=%s\n"
                 "fault=%d\n"
                 "sequence=%llu\n"
                 "state_since=%lld\n"
                 "check_in_deadline_ms=%lld\n"
                 "config_deadline_ms=%lld\n"
                 "gate_open=%d\n"
                 "configuration_seen=%d\n"
                 "configuration_applied=%d\n"
                 "request_id=%s\n"
                 "assignment_id=%s\n"
                 "reason=%s\n",
                 bootId,
                 lifecycle_state_str(fsm->state),
                 lifecycle_state_str(fsm->faultReturnState),
                 (int)fsm->fault,
                 (unsigned long long)fsm->sequence,
                 (long long)fsm->stateSince,
                 (long long)fsm->checkInDeadlineMs,
                 (long long)fsm->configDeadlineMs,
                 fsm->gateOpen ? 1 : 0,
                 fsm->configurationSeen ? 1 : 0,
                 fsm->configurationApplied ? 1 : 0,
                 fsm->requestId,
                 fsm->assignmentId,
                 fsm->reason) > 0;

    if (fflush(file) != 0 || fsync(fileno(file)) != 0) ok = false;
    if (fclose(file) != 0) ok = false;

    if (!ok || rename(temporary, path) != 0) {
        unlink(temporary);
        return false;
    }

    return true;
}
