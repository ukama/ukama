/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "log_store.h"

#include <dirent.h>
#include <errno.h>
#include <fcntl.h>
#include <limits.h>
#include <pthread.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <time.h>
#include <unistd.h>
#include <zstd.h>

#define DEFAULT_LOG_DIR        "/ukama/logs"
#define DEFAULT_NODE_ID        "unknown"
#define DEFAULT_ROTATE_BYTES   (16U * 1024U * 1024U)
#define DEFAULT_ROTATE_SECONDS (60 * 60)
#define DEFAULT_RETAIN_BYTES   (128U * 1024U * 1024U)
#define DEFAULT_RETAIN_DAYS    7
#define BOOT_ID_PATH           "/proc/sys/kernel/random/boot_id"
#define ACTIVE_FILE            "events-current.jsonl"
#define CURRENT_LINK           "current.jsonl"
#define STATE_FILE             "rlogd.json"
#define LINE_MAX_BYTES         (64U * 1024U)

struct LogStore {
    char logDir[PATH_MAX];
    char bootDir[PATH_MAX];
    char stateDir[PATH_MAX];
    char activePath[PATH_MAX];
    char currentLink[PATH_MAX];
    char statePath[PATH_MAX];
    char nodeId[128];
    char bootId[128];

    FILE *active;
    uint64_t currentSeq;
    uint64_t activeFirstSeq;
    size_t activeBytes;
    time_t activeOpenedAt;

    size_t rotateBytes;
    int rotateSeconds;
    size_t retainBytes;
    int retainDays;

    pthread_mutex_t mutex;
};

typedef struct {
    char path[PATH_MAX];
    time_t mtime;
    size_t size;
} RetainedFile;

static int mkdir_one(const char *path) {
    if (mkdir(path, 0755) == 0 || errno == EEXIST) {
        return 0;
    }
    return -1;
}

static int mkdir_p(const char *path) {
    char tmp[PATH_MAX];
    char *p;

    if (!path || !*path || strlen(path) >= sizeof(tmp)) {
        return -1;
    }

    strcpy(tmp, path);
    for (p = tmp + 1; *p; p++) {
        if (*p != '/') {
            continue;
        }
        *p = '\0';
        if (mkdir_one(tmp) != 0) {
            return -1;
        }
        *p = '/';
    }

    return mkdir_one(tmp);
}

static void trim_newline(char *s) {
    size_t len;

    if (!s) {
        return;
    }

    len = strlen(s);
    while (len > 0 && (s[len - 1] == '\n' || s[len - 1] == '\r')) {
        s[--len] = '\0';
    }
}

static void read_boot_id(char *buffer, size_t size) {
    FILE *fp;

    if (!buffer || size == 0) {
        return;
    }

    snprintf(buffer, size, "unknown-%ld", (long)getpid());

    fp = fopen(BOOT_ID_PATH, "r");
    if (!fp) {
        return;
    }

    if (fgets(buffer, (int)size, fp) != NULL) {
        trim_newline(buffer);
    }
    fclose(fp);
}

static void utc_timestamp(char *buffer, size_t size) {
    struct timespec now;
    struct tm utc;

    if (!buffer || size == 0) {
        return;
    }

    clock_gettime(CLOCK_REALTIME, &now);
    gmtime_r(&now.tv_sec, &utc);
    snprintf(buffer, size,
             "%04d-%02d-%02dT%02d:%02d:%02d.%03ldZ",
             utc.tm_year + 1900,
             utc.tm_mon + 1,
             utc.tm_mday,
             utc.tm_hour,
             utc.tm_min,
             utc.tm_sec,
             now.tv_nsec / 1000000L);
}

static uint64_t seq_from_record(json_t *record) {
    json_t *value;

    if (!record || !json_is_object(record)) {
        return 0;
    }

    value = json_object_get(record, "seq");
    if (!json_is_integer(value)) {
        return 0;
    }

    return (uint64_t)json_integer_value(value);
}

static void scan_active(LogStore *store) {
    FILE *fp;
    char *line;
    size_t cap;
    ssize_t n;
    uint64_t seq;

    if (!store) {
        return;
    }

    fp = fopen(store->activePath, "r");
    if (!fp) {
        return;
    }

    line = NULL;
    cap = 0;
    while ((n = getline(&line, &cap, fp)) >= 0) {
        json_error_t error;
        json_t *record;

        if ((size_t)n > LINE_MAX_BYTES) {
            continue;
        }

        record = json_loadb(line, (size_t)n, 0, &error);
        if (!record) {
            continue;
        }

        seq = seq_from_record(record);
        if (seq > 0) {
            if (store->activeFirstSeq == 0) {
                store->activeFirstSeq = seq;
            }
            if (seq > store->currentSeq) {
                store->currentSeq = seq;
            }
        }
        json_decref(record);
    }

    free(line);
    fclose(fp);
}

static uint64_t rotated_last_seq(const char *name) {
    unsigned long long first;
    unsigned long long last;

    if (!name) {
        return 0;
    }

    if (sscanf(name, "events-%llu-%llu.jsonl.zst", &first, &last) == 2) {
        (void)first;
        return (uint64_t)last;
    }

    return 0;
}

static void scan_rotated(LogStore *store) {
    DIR *dir;
    struct dirent *entry;
    uint64_t last;

    if (!store) {
        return;
    }

    dir = opendir(store->bootDir);
    if (!dir) {
        return;
    }

    while ((entry = readdir(dir)) != NULL) {
        last = rotated_last_seq(entry->d_name);
        if (last > store->currentSeq) {
            store->currentSeq = last;
        }
    }

    closedir(dir);
}

static int write_state(LogStore *store) {
    char tmpPath[PATH_MAX];
    json_t *state;
    char *data;
    FILE *fp;
    int rc;

    if (!store) {
        return -1;
    }

    if (snprintf(tmpPath, sizeof(tmpPath), "%s.tmp",
                 store->statePath) >= (int)sizeof(tmpPath)) {
        return -1;
    }

    state = json_object();
    if (!state) {
        return -1;
    }

    json_object_set_new(state, "schema",
                        json_string("ukama.rlog.state.v1"));
    json_object_set_new(state, "boot_id", json_string(store->bootId));
    json_object_set_new(state, "node_id", json_string(store->nodeId));
    json_object_set_new(state, "current_seq",
                        json_integer((json_int_t)store->currentSeq));
    json_object_set_new(state, "active_first_seq",
                        json_integer((json_int_t)store->activeFirstSeq));
    json_object_set_new(state, "active_bytes",
                        json_integer((json_int_t)store->activeBytes));
    json_object_set_new(state, "active_file",
                        json_string(store->activePath));

    data = json_dumps(state, JSON_COMPACT | JSON_SORT_KEYS);
    json_decref(state);
    if (!data) {
        return -1;
    }

    fp = fopen(tmpPath, "w");
    if (!fp) {
        free(data);
        return -1;
    }

    rc = fprintf(fp, "%s\n", data) < 0 ? -1 : 0;
    if (fflush(fp) != 0) {
        rc = -1;
    }
    fclose(fp);
    free(data);

    if (rc == 0 && rename(tmpPath, store->statePath) != 0) {
        rc = -1;
    }
    if (rc != 0) {
        unlink(tmpPath);
    }

    return rc;
}

static int compress_file(const char *inputPath, const char *outputPath) {
    FILE *input;
    FILE *output;
    struct stat st;
    void *src;
    void *dst;
    size_t bound;
    size_t compressed;
    size_t readBytes;
    int rc;

    if (stat(inputPath, &st) != 0 || st.st_size < 0) {
        return -1;
    }

    input = fopen(inputPath, "rb");
    if (!input) {
        return -1;
    }

    src = malloc((size_t)st.st_size ? (size_t)st.st_size : 1U);
    if (!src) {
        fclose(input);
        return -1;
    }

    readBytes = fread(src, 1, (size_t)st.st_size, input);
    fclose(input);
    if (readBytes != (size_t)st.st_size) {
        free(src);
        return -1;
    }

    bound = ZSTD_compressBound(readBytes);
    dst = malloc(bound ? bound : 1U);
    if (!dst) {
        free(src);
        return -1;
    }

    compressed = ZSTD_compress(dst, bound, src, readBytes, 3);
    free(src);
    if (ZSTD_isError(compressed)) {
        free(dst);
        return -1;
    }

    output = fopen(outputPath, "wb");
    if (!output) {
        free(dst);
        return -1;
    }

    rc = fwrite(dst, 1, compressed, output) == compressed ? 0 : -1;
    if (fflush(output) != 0) {
        rc = -1;
    }
    fclose(output);
    free(dst);

    if (rc != 0) {
        unlink(outputPath);
    }

    return rc;
}

static int retained_compare(const void *a, const void *b) {
    const RetainedFile *left;
    const RetainedFile *right;

    left = (const RetainedFile *)a;
    right = (const RetainedFile *)b;

    if (left->mtime < right->mtime) return -1;
    if (left->mtime > right->mtime) return 1;
    return strcmp(left->path, right->path);
}

static size_t collect_retained(LogStore *store,
                               RetainedFile **filesOut) {
    char bootsPath[PATH_MAX];
    DIR *boots;
    struct dirent *bootEntry;
    RetainedFile *files;
    size_t count;
    size_t capacity;

    if (!store || !filesOut) {
        return 0;
    }

    *filesOut = NULL;
    if (snprintf(bootsPath, sizeof(bootsPath), "%s/boots",
                 store->logDir) >= (int)sizeof(bootsPath)) {
        return 0;
    }

    boots = opendir(bootsPath);
    if (!boots) {
        return 0;
    }

    files = NULL;
    count = 0;
    capacity = 0;

    while ((bootEntry = readdir(boots)) != NULL) {
        char bootPath[PATH_MAX];
        DIR *bootDir;
        struct dirent *entry;

        if (bootEntry->d_name[0] == '.') {
            continue;
        }

        if (snprintf(bootPath, sizeof(bootPath), "%s/%s",
                     bootsPath, bootEntry->d_name) >=
            (int)sizeof(bootPath)) {
            continue;
        }

        bootDir = opendir(bootPath);
        if (!bootDir) {
            continue;
        }

        while ((entry = readdir(bootDir)) != NULL) {
            char path[PATH_MAX];
            struct stat st;
            size_t len;

            len = strlen(entry->d_name);
            if (len < 4 || strcmp(entry->d_name + len - 4, ".zst") != 0) {
                continue;
            }

            if (snprintf(path, sizeof(path), "%s/%s", bootPath,
                         entry->d_name) >= (int)sizeof(path)) {
                continue;
            }

            if (stat(path, &st) != 0 || !S_ISREG(st.st_mode)) {
                continue;
            }

            if (count == capacity) {
                size_t newCapacity;
                RetainedFile *tmp;

                newCapacity = capacity == 0 ? 16 : capacity * 2;
                tmp = realloc(files, newCapacity * sizeof(*files));
                if (!tmp) {
                    closedir(bootDir);
                    closedir(boots);
                    free(files);
                    return 0;
                }
                files = tmp;
                capacity = newCapacity;
            }

            snprintf(files[count].path, sizeof(files[count].path), "%s",
                     path);
            files[count].mtime = st.st_mtime;
            files[count].size = (size_t)st.st_size;
            count++;
        }
        closedir(bootDir);
    }

    closedir(boots);
    *filesOut = files;
    return count;
}

static void enforce_retention(LogStore *store) {
    RetainedFile *files;
    size_t count;
    size_t total;
    size_t idx;
    time_t now;
    time_t cutoff;

    if (!store) {
        return;
    }

    count = collect_retained(store, &files);
    if (count == 0 || !files) {
        free(files);
        return;
    }

    qsort(files, count, sizeof(*files), retained_compare);

    total = 0;
    for (idx = 0; idx < count; idx++) {
        total += files[idx].size;
    }

    now = time(NULL);
    cutoff = now - ((time_t)store->retainDays * 24 * 60 * 60);

    for (idx = 0; idx < count; idx++) {
        bool expired;
        bool overLimit;

        expired = store->retainDays > 0 && files[idx].mtime < cutoff;
        overLimit = store->retainBytes > 0 && total > store->retainBytes;
        if (!expired && !overLimit) {
            continue;
        }

        if (unlink(files[idx].path) == 0) {
            total -= files[idx].size;
        }
    }

    free(files);
}

static int open_active(LogStore *store) {
    struct stat st;

    if (!store) {
        return -1;
    }

    store->active = fopen(store->activePath, "a+");
    if (!store->active) {
        return -1;
    }

    if (stat(store->activePath, &st) == 0) {
        store->activeBytes = (size_t)st.st_size;
    } else {
        store->activeBytes = 0;
    }
    store->activeOpenedAt = time(NULL);

    return 0;
}

static int rotate_active(LogStore *store) {
    char plainPath[PATH_MAX];
    char compressedPath[PATH_MAX];

    if (!store || !store->active || store->activeFirstSeq == 0 ||
        store->activeBytes == 0) {
        return 0;
    }

    if (snprintf(plainPath, sizeof(plainPath),
                 "%s/events-%012llu-%012llu.jsonl",
                 store->bootDir,
                 (unsigned long long)store->activeFirstSeq,
                 (unsigned long long)store->currentSeq) >=
        (int)sizeof(plainPath)) {
        return -1;
    }

    if (snprintf(compressedPath, sizeof(compressedPath), "%s.zst",
                 plainPath) >= (int)sizeof(compressedPath)) {
        return -1;
    }

    fflush(store->active);
    fclose(store->active);
    store->active = NULL;

    if (rename(store->activePath, plainPath) != 0) {
        open_active(store);
        return -1;
    }

    if (compress_file(plainPath, compressedPath) != 0) {
        rename(plainPath, store->activePath);
        open_active(store);
        return -1;
    }

    unlink(plainPath);
    store->activeFirstSeq = 0;
    store->activeBytes = 0;

    if (open_active(store) != 0) {
        return -1;
    }

    enforce_retention(store);
    write_state(store);
    return 0;
}

static bool rotation_due(LogStore *store, size_t nextBytes) {
    time_t now;

    if (!store || store->activeBytes == 0) {
        return false;
    }

    if (store->rotateBytes > 0 &&
        store->activeBytes + nextBytes > store->rotateBytes) {
        return true;
    }

    now = time(NULL);
    if (store->rotateSeconds > 0 &&
        now - store->activeOpenedAt >= store->rotateSeconds) {
        return true;
    }

    return false;
}

LogStore *log_store_open(const LogStoreConfig *config) {
    LogStore *store;
    const char *logDir;
    const char *nodeId;
    char bootsDir[PATH_MAX];

    store = calloc(1, sizeof(*store));
    if (!store) {
        return NULL;
    }

    logDir = config && config->logDir && *config->logDir ?
             config->logDir : DEFAULT_LOG_DIR;
    nodeId = config && config->nodeId && *config->nodeId ?
             config->nodeId : DEFAULT_NODE_ID;

    snprintf(store->logDir, sizeof(store->logDir), "%s", logDir);
    snprintf(store->nodeId, sizeof(store->nodeId), "%s", nodeId);
    read_boot_id(store->bootId, sizeof(store->bootId));

    store->rotateBytes = config && config->rotateBytes ?
                         config->rotateBytes : DEFAULT_ROTATE_BYTES;
    store->rotateSeconds = config && config->rotateSeconds > 0 ?
                           config->rotateSeconds :
                           DEFAULT_ROTATE_SECONDS;
    store->retainBytes = config && config->retainBytes ?
                         config->retainBytes : DEFAULT_RETAIN_BYTES;
    store->retainDays = config && config->retainDays > 0 ?
                        config->retainDays : DEFAULT_RETAIN_DAYS;

    if (snprintf(bootsDir, sizeof(bootsDir), "%s/boots",
                 store->logDir) >= (int)sizeof(bootsDir) ||
        snprintf(store->bootDir, sizeof(store->bootDir), "%s/%s",
                 bootsDir, store->bootId) >=
                 (int)sizeof(store->bootDir) ||
        snprintf(store->stateDir, sizeof(store->stateDir), "%s/state",
                 store->logDir) >= (int)sizeof(store->stateDir) ||
        snprintf(store->activePath, sizeof(store->activePath), "%s/%s",
                 store->bootDir, ACTIVE_FILE) >=
                 (int)sizeof(store->activePath) ||
        snprintf(store->currentLink, sizeof(store->currentLink), "%s/%s",
                 store->logDir, CURRENT_LINK) >=
                 (int)sizeof(store->currentLink) ||
        snprintf(store->statePath, sizeof(store->statePath), "%s/%s",
                 store->stateDir, STATE_FILE) >=
                 (int)sizeof(store->statePath)) {
        free(store);
        return NULL;
    }

    if (mkdir_p(store->bootDir) != 0 || mkdir_p(store->stateDir) != 0) {
        free(store);
        return NULL;
    }

    pthread_mutex_init(&store->mutex, NULL);
    scan_rotated(store);
    scan_active(store);

    if (open_active(store) != 0) {
        pthread_mutex_destroy(&store->mutex);
        free(store);
        return NULL;
    }

    unlink(store->currentLink);
    if (symlink(store->activePath, store->currentLink) != 0) {
        log_store_close(store);
        return NULL;
    }

    enforce_retention(store);
    write_state(store);
    return store;
}

void log_store_close(LogStore *store) {
    if (!store) {
        return;
    }

    pthread_mutex_lock(&store->mutex);
    write_state(store);
    if (store->active) {
        fflush(store->active);
        fclose(store->active);
        store->active = NULL;
    }
    pthread_mutex_unlock(&store->mutex);

    pthread_mutex_destroy(&store->mutex);
    free(store);
}

int log_store_append(LogStore *store, json_t *record, uint64_t *seqOut) {
    json_t *copy;
    char timestamp[40];
    char *line;
    size_t lineBytes;
    uint64_t seq;
    int rc;

    if (!store || !record || !json_is_object(record)) {
        return -1;
    }

    copy = json_deep_copy(record);
    if (!copy) {
        return -1;
    }

    pthread_mutex_lock(&store->mutex);

    seq = store->currentSeq + 1;
    utc_timestamp(timestamp, sizeof(timestamp));

    json_object_set_new(copy, "node_id", json_string(store->nodeId));
    json_object_set_new(copy, "boot_id", json_string(store->bootId));
    json_object_set_new(copy, "seq", json_integer((json_int_t)seq));
    json_object_set_new(copy, "received_ts", json_string(timestamp));

    line = json_dumps(copy, JSON_COMPACT | JSON_ENSURE_ASCII);
    json_decref(copy);
    if (!line) {
        pthread_mutex_unlock(&store->mutex);
        return -1;
    }

    lineBytes = strlen(line) + 1;
    if (rotation_due(store, lineBytes) && rotate_active(store) != 0) {
        free(line);
        pthread_mutex_unlock(&store->mutex);
        return -1;
    }

    if (store->activeFirstSeq == 0) {
        store->activeFirstSeq = seq;
    }

    rc = fprintf(store->active, "%s\n", line) < 0 ? -1 : 0;
    if (rc == 0 && fflush(store->active) != 0) {
        rc = -1;
    }
    free(line);

    if (rc == 0) {
        store->currentSeq = seq;
        store->activeBytes += lineBytes;
        if (seqOut) {
            *seqOut = seq;
        }
    }

    pthread_mutex_unlock(&store->mutex);
    return rc;
}

const char *log_store_boot_id(const LogStore *store) {
    return store ? store->bootId : NULL;
}

const char *log_store_node_id(const LogStore *store) {
    return store ? store->nodeId : NULL;
}

const char *log_store_active_path(const LogStore *store) {
    return store ? store->activePath : NULL;
}

const char *log_store_state_dir(const LogStore *store) {
    return store ? store->stateDir : NULL;
}

uint64_t log_store_current_seq(const LogStore *store) {
    return store ? store->currentSeq : 0;
}

size_t log_store_active_bytes(const LogStore *store) {
    return store ? store->activeBytes : 0;
}
