/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "log_spool.h"

#include <dirent.h>
#include <errno.h>
#include <limits.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#define SPOOL_CURRENT_FILE "current.jsonl"
#define SPOOL_REPLAY_DIR   "replay"
#define SPOOL_STATE_FILE   "state.json"
#define SPOOL_FILE_PREFIX  "spool-"
#define SPOOL_FILE_SUFFIX  ".jsonl"
#define SPOOL_LINE_MAX     (128U * 1024U)
#define SPOOL_LOW_WATER    75U

struct LogSpool {
    char directory[PATH_MAX];
    char replayDir[PATH_MAX];
    char currentPath[PATH_MAX];
    char statePath[PATH_MAX];

    FILE *current;
    FILE *replay;
    char replayPath[PATH_MAX];

    size_t maxBytes;
    size_t totalBytes;
    size_t currentBytes;
    size_t replayBytes;
    uint64_t nextSegment;
    size_t replayCount;

    LogDropCounts dropped;
};

static int mkdir_one(const char *path) {
    if (mkdir(path, 0755) == 0 || errno == EEXIST) {
        return 0;
    }
    return -1;
}

static int mkdir_p(const char *path) {
    char tmp[PATH_MAX];
    char *cursor;

    if (!path || !*path || strlen(path) >= sizeof(tmp)) {
        return -1;
    }

    memcpy(tmp, path, strlen(path) + 1);
    for (cursor = tmp + 1; *cursor; cursor++) {
        if (*cursor != '/') continue;

        *cursor = '\0';
        if (mkdir_one(tmp) != 0) return -1;
        *cursor = '/';
    }

    return mkdir_one(tmp);
}

static uint64_t segment_number(const char *name) {
    unsigned long long value;
    char tail;

    if (!name) return 0;

    if (sscanf(name, SPOOL_FILE_PREFIX "%llu" SPOOL_FILE_SUFFIX "%c",
               &value, &tail) != 1) {
        return 0;
    }

    return (uint64_t)value;
}

static int write_state(LogSpool *spool) {
    char temporary[PATH_MAX];
    json_t *state;
    char *data;
    FILE *file;
    int rc;

    if (!spool) return -1;

    if (snprintf(temporary, sizeof(temporary), "%s.tmp",
                 spool->statePath) >= (int)sizeof(temporary)) {
        return -1;
    }

    state = json_pack("{s:s,s:I,s:I,s:I,s:I,s:I,s:I}",
                      "schema", "ukama.starterd.spool.v1",
                      "trace", (json_int_t)spool->dropped.trace,
                      "debug", (json_int_t)spool->dropped.debug,
                      "info", (json_int_t)spool->dropped.info,
                      "warn", (json_int_t)spool->dropped.warn,
                      "error", (json_int_t)spool->dropped.error,
                      "critical",
                      (json_int_t)spool->dropped.critical);
    if (!state) return -1;

    data = json_dumps(state, JSON_COMPACT | JSON_SORT_KEYS);
    json_decref(state);
    if (!data) return -1;

    file = fopen(temporary, "w");
    if (!file) {
        free(data);
        return -1;
    }

    rc = fprintf(file, "%s\n", data) < 0 ? -1 : 0;
    if (rc == 0 && fflush(file) != 0) rc = -1;
    if (rc == 0 && fsync(fileno(file)) != 0) rc = -1;
    fclose(file);
    free(data);

    if (rc == 0 && rename(temporary, spool->statePath) != 0) {
        rc = -1;
    }
    if (rc != 0) unlink(temporary);

    return rc;
}

static uint64_t state_integer(json_t *state, const char *key) {
    json_t *value;

    value = json_object_get(state, key);
    if (!json_is_integer(value)) return 0;
    return (uint64_t)json_integer_value(value);
}

static void load_state(LogSpool *spool) {
    json_error_t error;
    json_t *state;

    if (!spool) return;

    state = json_load_file(spool->statePath, 0, &error);
    if (!state || !json_is_object(state)) {
        if (state) json_decref(state);
        return;
    }

    spool->dropped.trace = state_integer(state, "trace");
    spool->dropped.debug = state_integer(state, "debug");
    spool->dropped.info = state_integer(state, "info");
    spool->dropped.warn = state_integer(state, "warn");
    spool->dropped.error = state_integer(state, "error");
    spool->dropped.critical = state_integer(state, "critical");

    json_decref(state);
}

static void count_drop(LogSpool *spool, json_t *frame) {
    json_t *record;
    const char *level;

    if (!spool) return;

    record = frame ? json_object_get(frame, "record") : NULL;
    level = json_string_value(record ?
                              json_object_get(record, "level") : NULL);

    if (level && strcmp(level, "trace") == 0) {
        spool->dropped.trace++;
    } else if (level && strcmp(level, "debug") == 0) {
        spool->dropped.debug++;
    } else if (level && strcmp(level, "warn") == 0) {
        spool->dropped.warn++;
    } else if (level && strcmp(level, "error") == 0) {
        spool->dropped.error++;
    } else if (level && strcmp(level, "critical") == 0) {
        spool->dropped.critical++;
    } else {
        spool->dropped.info++;
    }

    (void)write_state(spool);
}

static bool low_priority(json_t *frame) {
    json_t *record;
    const char *level;

    record = frame ? json_object_get(frame, "record") : NULL;
    level = json_string_value(record ?
                              json_object_get(record, "level") : NULL);

    return !level || strcmp(level, "trace") == 0 ||
           strcmp(level, "debug") == 0 || strcmp(level, "info") == 0;
}

static int scan_replay(LogSpool *spool) {
    DIR *directory;
    struct dirent *entry;

    directory = opendir(spool->replayDir);
    if (!directory) return -1;

    while ((entry = readdir(directory)) != NULL) {
        char path[PATH_MAX];
        struct stat info;
        uint64_t number;

        number = segment_number(entry->d_name);
        if (number == 0) continue;

        if (snprintf(path, sizeof(path), "%s/%s", spool->replayDir,
                     entry->d_name) >= (int)sizeof(path)) {
            continue;
        }
        if (stat(path, &info) != 0 || !S_ISREG(info.st_mode)) continue;

        spool->totalBytes += (size_t)info.st_size;
        spool->replayCount++;
        if (number >= spool->nextSegment) {
            spool->nextSegment = number + 1;
        }
    }

    closedir(directory);
    return 0;
}

static int open_current(LogSpool *spool) {
    struct stat info;

    spool->current = fopen(spool->currentPath, "a+");
    if (!spool->current) return -1;

    if (stat(spool->currentPath, &info) == 0) {
        spool->currentBytes = (size_t)info.st_size;
    } else {
        spool->currentBytes = 0;
    }

    return 0;
}

static int next_segment_path(LogSpool *spool,
                             char *path,
                             size_t pathSize) {
    uint64_t segment;

    segment = spool->nextSegment++;
    if (snprintf(path, pathSize, "%s/" SPOOL_FILE_PREFIX "%012llu"
                 SPOOL_FILE_SUFFIX, spool->replayDir,
                 (unsigned long long)segment) >= (int)pathSize) {
        return -1;
    }

    return 0;
}

static int seal_current(LogSpool *spool) {
    char path[PATH_MAX];

    if (!spool || !spool->current || spool->currentBytes == 0) {
        return 0;
    }

    if (next_segment_path(spool, path, sizeof(path)) != 0) return -1;

    if (fflush(spool->current) != 0 ||
        fsync(fileno(spool->current)) != 0) {
        return -1;
    }

    fclose(spool->current);
    spool->current = NULL;

    if (rename(spool->currentPath, path) != 0) {
        (void)open_current(spool);
        return -1;
    }

    spool->replayCount++;
    spool->currentBytes = 0;
    return open_current(spool);
}

static int seal_previous_current(LogSpool *spool) {
    struct stat info;
    char path[PATH_MAX];

    if (stat(spool->currentPath, &info) != 0 || info.st_size == 0) {
        return 0;
    }

    if (next_segment_path(spool, path, sizeof(path)) != 0) return -1;
    if (rename(spool->currentPath, path) != 0) return -1;

    spool->totalBytes += (size_t)info.st_size;
    spool->replayCount++;
    return 0;
}

static int find_oldest(LogSpool *spool,
                       char *path,
                       size_t pathSize,
                       size_t *bytesOut) {
    DIR *directory;
    struct dirent *entry;
    uint64_t oldest;
    char oldestName[NAME_MAX + 1];

    directory = opendir(spool->replayDir);
    if (!directory) return -1;

    oldest = 0;
    oldestName[0] = '\0';
    while ((entry = readdir(directory)) != NULL) {
        uint64_t number;

        number = segment_number(entry->d_name);
        if (number == 0 || (oldest != 0 && number >= oldest)) continue;

        oldest = number;
        snprintf(oldestName, sizeof(oldestName), "%s", entry->d_name);
    }
    closedir(directory);

    if (oldest == 0) return 0;
    if (snprintf(path, pathSize, "%s/%s", spool->replayDir,
                 oldestName) >= (int)pathSize) {
        return -1;
    }

    if (bytesOut) {
        struct stat info;

        *bytesOut = 0;
        if (stat(path, &info) == 0 && info.st_size >= 0) {
            *bytesOut = (size_t)info.st_size;
        }
    }

    return 1;
}

static int open_replay(LogSpool *spool) {
    int found;

    if (spool->replay) return 1;

    found = find_oldest(spool, spool->replayPath,
                        sizeof(spool->replayPath),
                        &spool->replayBytes);
    if (found <= 0) return found;

    spool->replay = fopen(spool->replayPath, "r");
    if (!spool->replay) return -1;

    return 1;
}

static void finish_replay(LogSpool *spool) {
    if (!spool || !spool->replay) return;

    fclose(spool->replay);
    spool->replay = NULL;

    if (unlink(spool->replayPath) == 0) {
        if (spool->totalBytes >= spool->replayBytes) {
            spool->totalBytes -= spool->replayBytes;
        } else {
            spool->totalBytes = 0;
        }
        if (spool->replayCount > 0) spool->replayCount--;
    }

    spool->replayPath[0] = '\0';
    spool->replayBytes = 0;
}

LogSpool *log_spool_open(const char *directory, size_t maxBytes) {
    LogSpool *spool;

    if (!directory || !*directory || strlen(directory) >= PATH_MAX) {
        errno = EINVAL;
        return NULL;
    }

    spool = calloc(1, sizeof(*spool));
    if (!spool) return NULL;

    memcpy(spool->directory, directory, strlen(directory) + 1);
    spool->maxBytes = maxBytes;
    spool->nextSegment = 1;

    if (snprintf(spool->replayDir, sizeof(spool->replayDir), "%s/%s",
                 directory, SPOOL_REPLAY_DIR) >=
        (int)sizeof(spool->replayDir) ||
        snprintf(spool->currentPath, sizeof(spool->currentPath), "%s/%s",
                 directory, SPOOL_CURRENT_FILE) >=
        (int)sizeof(spool->currentPath) ||
        snprintf(spool->statePath, sizeof(spool->statePath), "%s/%s",
                 directory, SPOOL_STATE_FILE) >=
        (int)sizeof(spool->statePath)) {
        free(spool);
        return NULL;
    }

    if (mkdir_p(spool->replayDir) != 0 || scan_replay(spool) != 0 ||
        seal_previous_current(spool) != 0 || open_current(spool) != 0) {
        log_spool_close(spool);
        return NULL;
    }

    load_state(spool);
    return spool;
}

void log_spool_close(LogSpool *spool) {
    if (!spool) return;

    if (spool->current) {
        fflush(spool->current);
        fsync(fileno(spool->current));
        fclose(spool->current);
    }
    if (spool->replay) fclose(spool->replay);

    (void)write_state(spool);
    free(spool);
}

int log_spool_append(LogSpool *spool, json_t *frame) {
    char *line;
    size_t lineBytes;
    size_t lowWater;
    int rc;

    if (!spool || !frame || !json_is_object(frame)) return -1;

    line = json_dumps(frame, JSON_COMPACT | JSON_ENSURE_ASCII);
    if (!line) return -1;

    lineBytes = strlen(line) + 1U;
    lowWater = (spool->maxBytes * SPOOL_LOW_WATER) / 100U;

    if ((spool->maxBytes > 0 &&
         spool->totalBytes + lineBytes > spool->maxBytes) ||
        (spool->maxBytes > 0 && spool->totalBytes >= lowWater &&
         low_priority(frame))) {
        count_drop(spool, frame);
        free(line);
        errno = ENOSPC;
        return 1;
    }

    rc = fprintf(spool->current, "%s\n", line) < 0 ? -1 : 0;
    if (rc == 0 && fflush(spool->current) != 0) rc = -1;
    free(line);

    if (rc == 0) {
        spool->currentBytes += lineBytes;
        spool->totalBytes += lineBytes;
    }

    return rc;
}

bool log_spool_has_backlog(LogSpool *spool) {
    if (!spool) return false;
    return spool->replay != NULL || spool->replayCount > 0 ||
           spool->currentBytes > 0;
}

int log_spool_replay(LogSpool *spool,
                     size_t maxRecords,
                     LogSpoolSendFn sendFn,
                     void *data) {
    size_t sent;
    char *line;
    size_t capacity;

    if (!spool || !sendFn || maxRecords == 0) return -1;

    sent = 0;
    line = NULL;
    capacity = 0;

    while (sent < maxRecords) {
        off_t offset;
        ssize_t size;
        json_error_t error;
        json_t *frame;
        int opened;

        opened = open_replay(spool);
        if (opened < 0) {
            free(line);
            return -1;
        }
        if (opened == 0) {
            if (spool->currentBytes == 0) break;
            if (seal_current(spool) != 0) {
                free(line);
                return -1;
            }
            continue;
        }

        offset = ftello(spool->replay);
        size = getline(&line, &capacity, spool->replay);
        if (size < 0) {
            if (feof(spool->replay)) {
                finish_replay(spool);
                continue;
            }
            free(line);
            return -1;
        }

        if ((size_t)size > SPOOL_LINE_MAX) continue;

        frame = json_loadb(line, (size_t)size, 0, &error);
        if (!frame || !json_is_object(frame)) {
            if (frame) json_decref(frame);
            continue;
        }

        if (!sendFn(frame, data)) {
            json_decref(frame);
            clearerr(spool->replay);
            if (fseeko(spool->replay, offset, SEEK_SET) != 0) {
                free(line);
                return -1;
            }
            break;
        }

        json_decref(frame);
        sent++;
    }

    free(line);
    return (int)sent;
}

void log_spool_get_dropped(LogSpool *spool, LogDropCounts *counts) {
    if (!counts) return;
    memset(counts, 0, sizeof(*counts));
    if (!spool) return;
    *counts = spool->dropped;
}

void log_spool_clear_dropped(LogSpool *spool) {
    if (!spool) return;
    memset(&spool->dropped, 0, sizeof(spool->dropped));
    (void)write_state(spool);
}

size_t log_spool_bytes(LogSpool *spool) {
    return spool ? spool->totalBytes : 0;
}
