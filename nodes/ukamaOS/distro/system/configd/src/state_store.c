/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef _GNU_SOURCE
#define _GNU_SOURCE
#endif

#include <errno.h>
#include <fcntl.h>
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/file.h>
#include <sys/stat.h>
#include <unistd.h>

#include "state_store.h"

#define STATE_FILE_VERSION "config-state-v1"

/* Preserve the public result codes used by the existing HTTP handlers. */
enum {
    STORE_OK          = 200,
    STORE_BAD_REQUEST = 400,
    STORE_CONFLICT    = 409,
    STORE_UNAVAILABLE = 503
};

static int sync_directory(const char *path) {

    int fd;
    int result;

    fd = open(path, O_RDONLY | O_DIRECTORY | O_CLOEXEC);
    if (fd < 0) {
        return -1;
    }

    result = fsync(fd);
    if (close(fd) != 0) {
        result = -1;
    }

    return result;
}

static int get_parent_directory(const char *path, char *parent, size_t size) {

    char *slash;

    if (snprintf(parent, size, "%s", path) >= (int)size) {
        return -1;
    }

    slash = strrchr(parent, '/');
    if (slash == NULL) {
        return -1;
    }

    if (slash == parent) {
        slash[1] = '\0';
    } else {
        *slash = '\0';
    }

    return 0;
}

static int create_directory(const char *path) {

    char parent[PATH_MAX];
    struct stat info;

    if (mkdir(path, 0700) == 0) {
        if (get_parent_directory(path, parent, sizeof(parent)) != 0) {
            return -1;
        }

        return sync_directory(parent);
    }

    if (errno != EEXIST) {
        return -1;
    }

    if (stat(path, &info) != 0 || !S_ISDIR(info.st_mode)) {
        return -1;
    }

    return 0;
}

static int make_directories(const char *path) {

    char directory[PATH_MAX];
    char *cursor;

    snprintf(directory, sizeof(directory), "%s", path);

    for (cursor = directory + 1; *cursor != '\0'; cursor++) {
        if (*cursor != '/') {
            continue;
        }

        *cursor = '\0';
        if (create_directory(directory) != 0) {
            return -1;
        }
        *cursor = '/';
    }

    return create_directory(directory);
}

int config_store_valid_id(const char *value) {

    const unsigned char *cursor;
    size_t length;

    if (value == NULL) {
        return 0;
    }

    length = strlen(value);
    if (length == 0 || length >= CONFIG_STATE_ID_SIZE) {
        return 0;
    }

    for (cursor = (const unsigned char *)value; *cursor; cursor++) {
        if ((*cursor >= 'a' && *cursor <= 'z') ||
            (*cursor >= 'A' && *cursor <= 'Z') ||
            (*cursor >= '0' && *cursor <= '9') ||
            *cursor == '-' || *cursor == '_' || *cursor == '.' ||
            *cursor == ':') {
            continue;
        }
        return 0;
    }

    return 1;
}

/* Write and sync a same-directory temporary file, rename, then sync parent.
 * Memory is advanced only after the complete durability boundary succeeds. */
static int write_record(FILE *file, const char *nodeId,
                        const ConfigRecord *record) {

    int index;

    if (fprintf(file, "%s\n%s\n", STATE_FILE_VERSION, nodeId) < 0) {
        return -1;
    }

    if (fprintf(file, "%d %d %" PRIu64 " %d %d\n",
                record->mode, record->phase, record->generation,
                record->revision, record->appCount) < 0) {
        return -1;
    }

    if (fprintf(file, "%s\n", record->requestId) < 0) {
        return -1;
    }

    for (index = 0; index < record->appCount; index++) {
        if (fprintf(file, "%s\n", record->apps[index]) < 0) {
            return -1;
        }
    }

    return 0;
}

static int save_record(ConfigStateStore *store, const ConfigRecord *record) {

    char temporary[PATH_MAX];
    FILE *file;
    int fd;
    int result = -1;

    if (snprintf(temporary, sizeof(temporary), "%s.tmp.XXXXXX", store->path)
        >= (int)sizeof(temporary)) {
        return -1;
    }

    fd = mkstemp(temporary);
    if (fd < 0) {
        return -1;
    }

    file = fdopen(fd, "w");
    if (file == NULL) {
        close(fd);
        unlink(temporary);
        return -1;
    }

    if (write_record(file, store->nodeId, record) != 0) {
        goto cleanup;
    }

    if (fflush(file) != 0 || fsync(fd) != 0) {
        goto cleanup;
    }

    result = fclose(file);
    file = NULL;
    if (result != 0) {
        goto cleanup;
    }

    result = rename(temporary, store->path);
    if (result != 0) {
        goto cleanup;
    }

    result = sync_directory(store->directory);

cleanup:
    if (file != NULL) {
        fclose(file);
    }

    if (result != 0) {
        unlink(temporary);
    }

    return result;
}

static int read_line(FILE *file, char *buffer, size_t size) {

    char *newline;

    if (fgets(buffer, (int)size, file) == NULL) {
        return -1;
    }

    newline = strchr(buffer, '\n');
    if (newline != NULL) {
        *newline = '\0';
        return 0;
    }

    /* A maximum-length value leaves its newline for the next read. */
    if (strlen(buffer) == size - 1 && fgetc(file) == '\n') {
        return 0;
    }

    return -1;
}

static int valid_app(const char *app) {

    return config_store_valid_id(app) && !strchr(app, ':') &&
        strcmp(app, ".") != 0 && strcmp(app, "..") != 0;
}

static int read_record_header(FILE *file, ConfigRecord *record) {

    char line[256];
    char canonical[256];
    char extra;
    int mode;
    int phase;

    if (read_line(file, line, sizeof(line)) != 0) {
        return -1;
    }

    if (sscanf(line, "%d %d %" SCNu64 " %d %d %c",
               &mode, &phase, &record->generation,
               &record->revision, &record->appCount, &extra) != 5) {
        return -1;
    }

    snprintf(canonical, sizeof(canonical), "%d %d %" PRIu64 " %d %d",
             mode, phase, record->generation,
             record->revision, record->appCount);
    if (strcmp(line, canonical) != 0) {
        return -1;
    }

    if (mode != CONFIG_MODE_NOCONFIG && mode != CONFIG_MODE_CONFIG) {
        return -1;
    }

    if (phase < CONFIG_PHASE_PENDING || phase > CONFIG_PHASE_FAILED) {
        return -1;
    }

    if (record->generation == 0 || record->generation > INT64_MAX) {
        return -1;
    }

    if (record->appCount < 0 || record->appCount > CONFIG_STATE_APPS) {
        return -1;
    }

    record->mode = mode;
    record->phase = phase;

    if (mode == CONFIG_MODE_NOCONFIG) {
        if (record->revision != 0 || record->appCount != 0) {
            return -1;
        }
    } else {
        if (record->revision <= 0) {
            return -1;
        }
        if (phase == CONFIG_PHASE_COMPLETED && record->appCount == 0) {
            return -1;
        }
    }

    return 0;
}

static int find_app(const ConfigRecord *record, const char *app) {

    int index;

    for (index = 0; index < record->appCount; index++) {
        if (strcmp(record->apps[index], app) == 0) {
            return index;
        }
    }

    return -1;
}

static int read_record_apps(FILE *file, ConfigRecord *record) {

    char app[CONFIG_STATE_APP_SIZE];
    int expectedCount;

    expectedCount = record->appCount;
    record->appCount = 0;

    while (record->appCount < expectedCount) {
        if (read_line(file, app, sizeof(app)) != 0 || !valid_app(app)) {
            return -1;
        }
        if (find_app(record, app) >= 0) {
            return -1;
        }

        snprintf(record->apps[record->appCount], CONFIG_STATE_APP_SIZE,
                 "%s", app);
        record->appCount++;
    }

    return 0;
}

static int read_record(FILE *file, const char *nodeId, ConfigRecord *record) {

    char line[256];

    if (read_line(file, line, sizeof(line)) != 0 ||
        strcmp(line, STATE_FILE_VERSION) != 0) {
        return -1;
    }

    if (read_line(file, line, sizeof(line)) != 0 ||
        strcmp(line, nodeId) != 0) {
        return -1;
    }

    if (read_record_header(file, record) != 0) {
        return -1;
    }

    if (read_line(file, record->requestId, sizeof(record->requestId)) != 0 ||
        !config_store_valid_id(record->requestId)) {
        return -1;
    }

    if (read_record_apps(file, record) != 0) {
        return -1;
    }

    if (fgetc(file) != EOF || ferror(file)) {
        return -1;
    }

    return 0;
}

/* Return 1 for a loaded record, 0 for no file, and -1 for invalid state. */
static int load_record(ConfigStateStore *store) {

    ConfigRecord record = {0};
    FILE *file;
    int result;

    file = fopen(store->path, "r");
    if (file == NULL) {
        if (errno == ENOENT) {
            return 0;
        }
        return -1;
    }

    result = read_record(file, store->nodeId, &record);
    fclose(file);
    if (result != 0) {
        return -1;
    }

    store->record = record;
    return 1;
}

static int app_config_matches(ConfigStateStore *store, const char *app) {

    char active[PATH_MAX];
    char expected[PATH_MAX];
    char resolvedActive[PATH_MAX];
    char resolvedExpected[PATH_MAX];
    struct stat info;

    if (snprintf(active, sizeof(active), "%s/%s/active", store->configRoot, app)
        >= (int)sizeof(active)) {
        return 0;
    }

    if (snprintf(expected, sizeof(expected), "%s/%s/archive/%d",
                 store->configRoot, app, store->record.revision)
        >= (int)sizeof(expected)) {
        return 0;
    }

    if (realpath(active, resolvedActive) == NULL ||
        realpath(expected, resolvedExpected) == NULL) {
        return 0;
    }

    if (strcmp(resolvedActive, resolvedExpected) != 0) {
        return 0;
    }

    if (stat(resolvedActive, &info) != 0) {
        return 0;
    }

    return S_ISDIR(info.st_mode);
}

static int active_config_matches(ConfigStateStore *store) {

    int index;

    if (store->record.appCount == 0) {
        return 0;
    }

    for (index = 0; index < store->record.appCount; index++) {
        if (!app_config_matches(store, store->record.apps[index])) {
            return 0;
        }
    }

    return 1;
}

static void set_error(ConfigStateStore *store, const char *error) {

    store->record.phase = CONFIG_PHASE_FAILED;
    snprintf(store->record.error, sizeof(store->record.error), "%s", error);
}

static int commit_record(ConfigStateStore *store, ConfigRecord *record) {

    if (save_record(store, record) != 0) {
        set_error(store, "state_write_failed");
        return STORE_UNAVAILABLE;
    }

    store->record = *record;
    return STORE_OK;
}

static void restore_record(ConfigStateStore *store) {

    ConfigRecord record;

    if (load_record(store) < 0) {
        set_error(store, "state_invalid_or_unreadable");
        return;
    }

    record = store->record;
    if (record.phase == CONFIG_PHASE_PENDING) {
        if (record.mode == CONFIG_MODE_NOCONFIG) {
            /* NOCONFIG has no application side effects to recover. */
            record.phase = CONFIG_PHASE_COMPLETED;
            commit_record(store, &record);
        } else {
            set_error(store, "configuration_interrupted");
        }
        return;
    }

    if (record.mode == CONFIG_MODE_CONFIG &&
        record.phase == CONFIG_PHASE_COMPLETED &&
        !active_config_matches(store)) {
        set_error(store, "active_configuration_mismatch");
    }
}

static int lock_store(ConfigStateStore *store) {

    char lockPath[PATH_MAX];

    if (snprintf(lockPath, sizeof(lockPath), "%s.lock", store->path)
        >= (int)sizeof(lockPath)) {
        return -1;
    }

    store->lockFd = open(lockPath, O_RDWR | O_CREAT | O_CLOEXEC, 0600);
    if (store->lockFd < 0) {
        return -1;
    }

    return flock(store->lockFd, LOCK_EX | LOCK_NB);
}

int config_store_open(ConfigStateStore *store, const char *path,
                      const char *configRoot, const char *nodeId) {

    if (store == NULL || path == NULL || configRoot == NULL || nodeId == NULL) {

        return -1;
    }

    if (path[0] != '/' || nodeId[0] == '\0') {
        return -1;
    }

    if (strchr(nodeId, '\n') != NULL || strchr(nodeId, '\r') != NULL) {
        return -1;
    }

    if (strlen(path) >= sizeof(store->path) - 16 ||
        strlen(configRoot) >= sizeof(store->configRoot) ||
        strlen(nodeId) >= sizeof(store->nodeId)) {
        return -1;
    }

    memset(store, 0, sizeof(*store));
    store->lockFd = -1;
    if (pthread_mutex_init(&store->mutex, NULL) != 0) {
        return -1;
    }

    snprintf(store->path, sizeof(store->path), "%s", path);
    snprintf(store->configRoot, sizeof(store->configRoot), "%s", configRoot);
    snprintf(store->nodeId, sizeof(store->nodeId), "%s", nodeId);

    if (get_parent_directory(path, store->directory,
                             sizeof(store->directory)) != 0) {
        goto fail;
    }

    if (make_directories(store->directory) != 0) {
        goto fail;
    }

    if (lock_store(store) != 0) {
        goto fail;
    }

    restore_record(store);
    return 0;

fail:
    config_store_close(store);
    return -1;
}

void config_store_close(ConfigStateStore *store) {

    if (store->lockFd >= 0) {
        close(store->lockFd);
    }

    store->lockFd = -1;

    pthread_mutex_destroy(&store->mutex);
}

void config_store_snapshot(ConfigStateStore *store, ConfigRecord *record) {

    pthread_mutex_lock(&store->mutex);
    *record = store->record;
    pthread_mutex_unlock(&store->mutex);
}

/* A transient write failure must not require a daemon restart. Re-read only
 * after syncing the directory; never assume whether the failed rename landed. */
static int recover_write(ConfigStateStore *store) {

    ConfigRecord previous;

    if (strcmp(store->record.error, "state_write_failed") != 0) {
        return 0;
    }

    previous = store->record;
    if (sync_directory(store->directory) != 0) {
        return -1;
    }

    memset(&store->record, 0, sizeof(store->record));
    if (load_record(store) < 0) {
        store->record = previous;
        return -1;
    }

    return 0;
}

int config_store_noconfig(ConfigStateStore *store, const char *requestId) {

    ConfigRecord record;
    int result = STORE_CONFLICT;

    if (!config_store_valid_id(requestId)) {
        return STORE_BAD_REQUEST;
    }

    pthread_mutex_lock(&store->mutex);
    if (recover_write(store) != 0) {
        result = STORE_UNAVAILABLE;
        goto done;
    }

    record = store->record;
    if (record.error[0] != '\0') {
        result = STORE_UNAVAILABLE;
        goto done;
    }

    if (record.mode == CONFIG_MODE_CONFIG ||
        (record.mode == CONFIG_MODE_NOCONFIG &&
         strcmp(record.requestId, requestId) != 0)) {
        goto done;
    }

    if (record.generation >= INT64_MAX) {
        result = STORE_UNAVAILABLE;
        goto done;
    }

    record.mode = CONFIG_MODE_NOCONFIG;
    record.generation++;
    snprintf(record.requestId, sizeof(record.requestId), "%s", requestId);
    if (record.phase != CONFIG_PHASE_COMPLETED) {
        record.phase = CONFIG_PHASE_PENDING;
        result = commit_record(store, &record);
        if (result != STORE_OK) {
            goto done;
        }
    }

    record.phase = CONFIG_PHASE_COMPLETED;
    result = commit_record(store, &record);

done:
    pthread_mutex_unlock(&store->mutex);
    return result;
}

int config_store_begin(ConfigStateStore *store, const char *requestId,
                       int revision) {

    ConfigRecord record = {0};
    int result = STORE_UNAVAILABLE;

    if (!config_store_valid_id(requestId) || revision <= 0) {
        return STORE_BAD_REQUEST;
    }

    pthread_mutex_lock(&store->mutex);

    if (recover_write(store) != 0) {
        goto done;
    }

    /* Explicit CONFIG retries may replace a failed application attempt,
     * but cannot overwrite corrupt or uncertain persistent state. */
    if (strcmp(store->record.error, "state_invalid_or_unreadable") == 0 ||
        strcmp(store->record.error, "state_write_failed") == 0) {
        goto done;
    }

    if (store->record.generation >= INT64_MAX) {
        goto done;
    }

    record.mode       = CONFIG_MODE_CONFIG;
    record.phase      = CONFIG_PHASE_PENDING;
    record.generation = store->record.generation + 1;
    record.revision   = revision;
    snprintf(record.requestId, sizeof(record.requestId), "%s", requestId);

    result = commit_record(store, &record);

done:
    pthread_mutex_unlock(&store->mutex);
    return result;
}

int config_store_track_app(ConfigStateStore *store, const char *app) {

    ConfigRecord record;
    int result = STORE_CONFLICT;

    if (!valid_app(app)) {
        return STORE_BAD_REQUEST;
    }

    pthread_mutex_lock(&store->mutex);
    record = store->record;

    if (record.mode != CONFIG_MODE_CONFIG ||
        record.phase != CONFIG_PHASE_PENDING || record.error[0] != '\0') {
        goto done;
    }

    if (find_app(&record, app) >= 0) {
        result = STORE_OK;
        goto done;
    }

    if (record.appCount >= CONFIG_STATE_APPS) {
        goto done;
    }

    snprintf(record.apps[record.appCount], CONFIG_STATE_APP_SIZE, "%s", app);
    record.appCount++;
    result = commit_record(store, &record);

done:
    pthread_mutex_unlock(&store->mutex);
    return result;
}

static int sync_active_config(ConfigStateStore *store) {

    int fd;
    int result;

    /* CONFIG already activates files/symlinks in this filesystem. Flush those
     * writes before committing a status that will survive the same reboot. */
    fd = open(store->configRoot, O_RDONLY | O_DIRECTORY | O_CLOEXEC);
    if (fd < 0) {
        return -1;
    }

    result = syncfs(fd);
    if (close(fd) != 0) {
        result = -1;
    }

    return result;
}

int config_store_finish(ConfigStateStore *store, int success) {

    ConfigRecord record;
    int result = STORE_CONFLICT;

    pthread_mutex_lock(&store->mutex);
    record = store->record;

    if (record.mode != CONFIG_MODE_CONFIG ||
        record.phase != CONFIG_PHASE_PENDING || record.error[0] != '\0') {
        goto done;
    }

    if (success && !active_config_matches(store)) {
        success = 0;
    }

    if (success && sync_active_config(store) != 0) {
        success = 0;
    }

    record.phase = CONFIG_PHASE_FAILED;
    if (success) {
        record.phase = CONFIG_PHASE_COMPLETED;
    }

    result = commit_record(store, &record);
    if (result == STORE_OK && !success) {
        result = STORE_CONFLICT;
    }

done:
    pthread_mutex_unlock(&store->mutex);
    return result;
}
