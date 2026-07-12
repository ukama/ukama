/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#define _GNU_SOURCE

#include <dirent.h>
#include <errno.h>
#include <getopt.h>
#include <jansson.h>
#include <limits.h>
#include <signal.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <sys/stat.h>
#include <time.h>
#include <unistd.h>
#include <zstd.h>

#include "version.h"

#define DEFAULT_LOG_DIR      "/ukama/logs"
#define CURRENT_LINK         "current.jsonl"
#define STATE_FILE           "state/rlogd.json"
#define BOOTS_DIR            "boots"
#define ACTIVE_FILE          "events-current.jsonl"
#define MAX_RECORD_BYTES     (256U * 1024U)
#define FOLLOW_INTERVAL_NS   200000000L
#define NAME_MAX_BYTES       128

#define EXIT_ERROR 1
#define EXIT_USAGE 2

typedef enum {
    FORMAT_TEXT = 0,
    FORMAT_JSONL,
    FORMAT_SUMMARY
} OutputFormat;

typedef struct {
    char **items;
    size_t count;
} StringList;

typedef struct {
    char logDir[PATH_MAX];
    char boot[NAME_MAX_BYTES];
    bool allBoots;
    bool follow;
    OutputFormat format;

    StringList levels;
    StringList apps;
    char event[NAME_MAX_BYTES];
    char session[NAME_MAX_BYTES];
    char trace[NAME_MAX_BYTES];

    bool hasSince;
    time_t since;

    bool hasSeqStart;
    bool hasSeqEnd;
    uint64_t seqStart;
    uint64_t seqEnd;
} Options;

typedef struct {
    char path[PATH_MAX];
    char id[NAME_MAX_BYTES];
    time_t mtime;
} BootDir;

typedef struct {
    char path[PATH_MAX];
    uint64_t firstSeq;
    bool active;
    bool compressed;
} LogFile;

typedef struct {
    char name[NAME_MAX_BYTES];
    uint64_t count;
} Counter;

typedef struct {
    uint64_t records;
    uint64_t invalid;
    Counter *levels;
    size_t levelCount;
    Counter *apps;
    size_t appCount;
    Counter *events;
    size_t eventCount;
} Summary;

typedef struct {
    Options *options;
    Summary *summary;
    char lastBootId[NAME_MAX_BYTES];
    uint64_t lastSeq;
    bool suppressDuplicates;
} Reader;

typedef struct {
    char *data;
    size_t length;
    size_t capacity;
    bool dropping;
    Reader *reader;
} LineAssembler;

static volatile sig_atomic_t gStop = 0;

static void usage(FILE *stream) {
    fprintf(stream, "ukama-log: read Ukama node logs\n");
    fprintf(stream, "Usage: ukama-log [options]\n\n");
    fprintf(stream, "Options:\n");
    fprintf(stream, "  -f, --follow              Follow the current log\n");
    fprintf(stream, "      --log-dir PATH        Log directory");
    fprintf(stream, " (default: %s)\n", DEFAULT_LOG_DIR);
    fprintf(stream, "      --boot ID             current, all, or a boot ID\n");
    fprintf(stream, "      --since DURATION      Records from the last ");
    fprintf(stream, "10s, 5m, 2h, or 3d\n");
    fprintf(stream, "      --level LIST          Comma-separated levels\n");
    fprintf(stream, "      --app NAME            Match one application\n");
    fprintf(stream, "      --apps LIST           Comma-separated ");
    fprintf(stream, "applications\n");
    fprintf(stream, "      --event NAME          Match event\n");
    fprintf(stream, "      --session ID          Match session_id\n");
    fprintf(stream, "      --trace ID            Match trace_id\n");
    fprintf(stream, "      --seq START[:END]     Match canonical sequence\n");
    fprintf(stream, "      --format FORMAT       text, jsonl, or summary\n");
    fprintf(stream, "  -v, --version             Show version\n");
    fprintf(stream, "  -h, --help                Show this help\n");
}

static void on_signal(int signalNumber) {
    (void)signalNumber;
    gStop = 1;
}

static const char *json_string_or(json_t *record, const char *key,
                                  const char *fallback) {
    json_t *value;
    const char *text;

    value = json_object_get(record, key);
    text = json_string_value(value);
    return text ? text : fallback;
}

static uint64_t json_u64(json_t *record, const char *key) {
    json_t *value;
    json_int_t number;

    value = json_object_get(record, key);
    if (!json_is_integer(value)) {
        return 0;
    }

    number = json_integer_value(value);
    return number > 0 ? (uint64_t)number : 0;
}

static int string_list_add(StringList *list, const char *value) {
    char **items;
    char *copy;

    if (!list || !value || !*value) {
        return -1;
    }

    copy = strdup(value);
    if (!copy) {
        return -1;
    }

    items = realloc(list->items,
                    (list->count + 1) * sizeof(*list->items));
    if (!items) {
        free(copy);
        return -1;
    }

    list->items = items;
    list->items[list->count++] = copy;
    return 0;
}

static int string_list_add_csv(StringList *list, const char *csv) {
    char *copy;
    char *save;
    char *token;
    int rc;

    if (!list || !csv || !*csv) {
        return -1;
    }

    copy = strdup(csv);
    if (!copy) {
        return -1;
    }

    rc = 0;
    save = NULL;
    token = strtok_r(copy, ",", &save);
    while (token) {
        while (*token == ' ' || *token == '\t') {
            token++;
        }
        if (*token && string_list_add(list, token) != 0) {
            rc = -1;
            break;
        }
        token = strtok_r(NULL, ",", &save);
    }

    free(copy);
    return rc;
}

static bool level_list_contains(const StringList *list,
                                const char *value) {
    size_t idx;

    if (!list || list->count == 0) {
        return true;
    }
    if (!value) {
        return false;
    }

    for (idx = 0; idx < list->count; idx++) {
        if (strcasecmp(list->items[idx], value) == 0) {
            return true;
        }
    }

    return false;
}

static bool string_list_contains(const StringList *list,
                                 const char *value) {
    size_t idx;

    if (!list || list->count == 0) {
        return true;
    }
    if (!value) {
        return false;
    }

    for (idx = 0; idx < list->count; idx++) {
        if (strcmp(list->items[idx], value) == 0) {
            return true;
        }
    }

    return false;
}

static void string_list_free(StringList *list) {
    size_t idx;

    if (!list) {
        return;
    }

    for (idx = 0; idx < list->count; idx++) {
        free(list->items[idx]);
    }
    free(list->items);
    memset(list, 0, sizeof(*list));
}

static int parse_duration(const char *text, time_t *since) {
    unsigned long long value;
    unsigned long long multiplier;
    unsigned long long seconds;
    char *end;
    time_t now;

    if (!text || !*text || !since) {
        return -1;
    }

    errno = 0;
    value = strtoull(text, &end, 10);
    if (errno != 0 || end == text) {
        return -1;
    }

    multiplier = 1;
    if (*end != '\0') {
        if (end[1] != '\0') {
            return -1;
        }
        switch (*end) {
        case 's':
            multiplier = 1;
            break;
        case 'm':
            multiplier = 60;
            break;
        case 'h':
            multiplier = 60 * 60;
            break;
        case 'd':
            multiplier = 24 * 60 * 60;
            break;
        default:
            return -1;
        }
    }

    if (value > ULLONG_MAX / multiplier) {
        return -1;
    }

    seconds = value * multiplier;
    now = time(NULL);
    if (seconds > (unsigned long long)now) {
        *since = 0;
    } else {
        *since = now - (time_t)seconds;
    }

    return 0;
}

static int parse_u64(const char *text, uint64_t *value) {
    unsigned long long parsed;
    char *end;

    if (!text || !*text || !value) {
        return -1;
    }

    errno = 0;
    parsed = strtoull(text, &end, 10);
    if (errno != 0 || end == text || *end != '\0') {
        return -1;
    }

    *value = (uint64_t)parsed;
    return 0;
}

static int parse_seq_range(const char *text, Options *options) {
    char buffer[128];
    char *separator;

    if (!text || !options || strlen(text) >= sizeof(buffer)) {
        return -1;
    }

    strcpy(buffer, text);
    separator = strchr(buffer, ':');
    if (!separator) {
        if (parse_u64(buffer, &options->seqStart) != 0) {
            return -1;
        }
        options->seqEnd = options->seqStart;
        options->hasSeqStart = true;
        options->hasSeqEnd = true;
        return 0;
    }

    *separator = '\0';
    if (*buffer) {
        if (parse_u64(buffer, &options->seqStart) != 0) {
            return -1;
        }
        options->hasSeqStart = true;
    }

    separator++;
    if (*separator) {
        if (parse_u64(separator, &options->seqEnd) != 0) {
            return -1;
        }
        options->hasSeqEnd = true;
    }

    if (!options->hasSeqStart && !options->hasSeqEnd) {
        return -1;
    }
    if (options->hasSeqStart && options->hasSeqEnd &&
        options->seqStart > options->seqEnd) {
        return -1;
    }

    return 0;
}

static int parse_format(const char *text, OutputFormat *format) {
    if (!text || !format) {
        return -1;
    }

    if (strcmp(text, "text") == 0) {
        *format = FORMAT_TEXT;
    } else if (strcmp(text, "jsonl") == 0) {
        *format = FORMAT_JSONL;
    } else if (strcmp(text, "summary") == 0) {
        *format = FORMAT_SUMMARY;
    } else {
        return -1;
    }

    return 0;
}

static int copy_option(char *destination, size_t size,
                       const char *value) {
    size_t length;

    if (!destination || size == 0 || !value) {
        return -1;
    }

    length = strlen(value);
    if (length >= size) {
        return -1;
    }

    memcpy(destination, value, length + 1);
    return 0;
}

static int parse_options(int argc, char **argv, Options *options) {
    enum {
        OPT_LOG_DIR = 1000,
        OPT_BOOT,
        OPT_SINCE,
        OPT_LEVEL,
        OPT_APP,
        OPT_APPS,
        OPT_EVENT,
        OPT_SESSION,
        OPT_TRACE,
        OPT_SEQ,
        OPT_FORMAT
    };
    static struct option longOptions[] = {
        {"follow", no_argument, NULL, 'f'},
        {"log-dir", required_argument, NULL, OPT_LOG_DIR},
        {"boot", required_argument, NULL, OPT_BOOT},
        {"since", required_argument, NULL, OPT_SINCE},
        {"level", required_argument, NULL, OPT_LEVEL},
        {"app", required_argument, NULL, OPT_APP},
        {"apps", required_argument, NULL, OPT_APPS},
        {"event", required_argument, NULL, OPT_EVENT},
        {"session", required_argument, NULL, OPT_SESSION},
        {"trace", required_argument, NULL, OPT_TRACE},
        {"seq", required_argument, NULL, OPT_SEQ},
        {"format", required_argument, NULL, OPT_FORMAT},
        {"version", no_argument, NULL, 'v'},
        {"help", no_argument, NULL, 'h'},
        {NULL, 0, NULL, 0}
    };
    int option;

    memset(options, 0, sizeof(*options));
    strcpy(options->logDir, DEFAULT_LOG_DIR);
    strcpy(options->boot, "current");
    options->format = FORMAT_TEXT;

    while ((option = getopt_long(argc, argv, "fhv", longOptions,
                                 NULL)) != -1) {
        switch (option) {
        case 'f':
            options->follow = true;
            break;
        case 'h':
            usage(stdout);
            exit(0);
        case 'v':
            printf("ukama-log - Version: %s\n", VERSION);
            exit(0);
        case OPT_LOG_DIR:
            if (copy_option(options->logDir,
                            sizeof(options->logDir), optarg) != 0) {
                return -1;
            }
            break;
        case OPT_BOOT:
            if (copy_option(options->boot,
                            sizeof(options->boot), optarg) != 0) {
                return -1;
            }
            options->allBoots = strcmp(optarg, "all") == 0;
            break;
        case OPT_SINCE:
            if (parse_duration(optarg, &options->since) != 0) {
                return -1;
            }
            options->hasSince = true;
            break;
        case OPT_LEVEL:
            if (string_list_add_csv(&options->levels, optarg) != 0) {
                return -1;
            }
            break;
        case OPT_APP:
            if (string_list_add(&options->apps, optarg) != 0) {
                return -1;
            }
            break;
        case OPT_APPS:
            if (string_list_add_csv(&options->apps, optarg) != 0) {
                return -1;
            }
            break;
        case OPT_EVENT:
            if (copy_option(options->event,
                            sizeof(options->event), optarg) != 0) {
                return -1;
            }
            break;
        case OPT_SESSION:
            if (copy_option(options->session,
                            sizeof(options->session), optarg) != 0) {
                return -1;
            }
            break;
        case OPT_TRACE:
            if (copy_option(options->trace,
                            sizeof(options->trace), optarg) != 0) {
                return -1;
            }
            break;
        case OPT_SEQ:
            if (parse_seq_range(optarg, options) != 0) {
                return -1;
            }
            break;
        case OPT_FORMAT:
            if (parse_format(optarg, &options->format) != 0) {
                return -1;
            }
            break;
        default:
            return -1;
        }
    }

    if (optind != argc) {
        return -1;
    }
    if (options->follow && options->allBoots) {
        fprintf(stderr, "ukama-log: --follow cannot use --boot all\n");
        return -1;
    }
    if (options->follow && options->format == FORMAT_SUMMARY) {
        fprintf(stderr,
                "ukama-log: summary format cannot be followed\n");
        return -1;
    }
    if (options->follow && strcmp(options->boot, "current") != 0) {
        fprintf(stderr,
                "ukama-log: --follow requires --boot current\n");
        return -1;
    }

    return 0;
}

static bool valid_boot_id(const char *bootId) {
    const unsigned char *cursor;

    if (!bootId || !*bootId || strlen(bootId) >= NAME_MAX_BYTES) {
        return false;
    }

    cursor = (const unsigned char *)bootId;
    while (*cursor) {
        if (*cursor == '/' || *cursor == '\\') {
            return false;
        }
        cursor++;
    }

    return true;
}

static int boot_id_from_link(const Options *options, char *bootId,
                             size_t size) {
    char linkPath[PATH_MAX];
    char target[PATH_MAX];
    char *lastSlash;
    char *parentSlash;
    ssize_t length;

    if (snprintf(linkPath, sizeof(linkPath), "%s/%s",
                 options->logDir, CURRENT_LINK) >=
        (int)sizeof(linkPath)) {
        return -1;
    }

    length = readlink(linkPath, target, sizeof(target) - 1);
    if (length < 0) {
        return -1;
    }
    target[length] = '\0';

    lastSlash = strrchr(target, '/');
    if (!lastSlash) {
        return -1;
    }
    *lastSlash = '\0';
    parentSlash = strrchr(target, '/');
    if (!parentSlash || !parentSlash[1]) {
        return -1;
    }

    return copy_option(bootId, size, parentSlash + 1);
}

static int boot_id_from_state(const Options *options, char *bootId,
                              size_t size) {
    char path[PATH_MAX];
    json_error_t error;
    json_t *state;
    const char *value;
    int rc;

    if (snprintf(path, sizeof(path), "%s/%s",
                 options->logDir, STATE_FILE) >= (int)sizeof(path)) {
        return -1;
    }

    state = json_load_file(path, 0, &error);
    if (!state) {
        return -1;
    }

    value = json_string_value(json_object_get(state, "boot_id"));
    rc = value ? copy_option(bootId, size, value) : -1;
    json_decref(state);
    return rc;
}

static int current_boot_id(const Options *options, char *bootId,
                           size_t size) {
    if (boot_id_from_link(options, bootId, size) == 0) {
        return 0;
    }
    return boot_id_from_state(options, bootId, size);
}

static int boot_compare(const void *left, const void *right) {
    const BootDir *a;
    const BootDir *b;

    a = left;
    b = right;
    if (a->mtime < b->mtime) {
        return -1;
    }
    if (a->mtime > b->mtime) {
        return 1;
    }
    return strcmp(a->id, b->id);
}

static int add_boot(BootDir **boots, size_t *count,
                    const char *path, const char *id) {
    BootDir *items;
    struct stat status;

    if (stat(path, &status) != 0 || !S_ISDIR(status.st_mode)) {
        return -1;
    }

    items = realloc(*boots, (*count + 1) * sizeof(**boots));
    if (!items) {
        return -1;
    }
    *boots = items;

    memset(&items[*count], 0, sizeof(items[*count]));
    if (copy_option(items[*count].path,
                    sizeof(items[*count].path), path) != 0 ||
        copy_option(items[*count].id,
                    sizeof(items[*count].id), id) != 0) {
        return -1;
    }
    items[*count].mtime = status.st_mtime;
    (*count)++;
    return 0;
}

static int collect_boots(const Options *options, BootDir **bootsOut,
                         size_t *countOut) {
    char bootsPath[PATH_MAX];
    char bootPath[PATH_MAX];
    char bootId[NAME_MAX_BYTES];
    DIR *directory;
    struct dirent *entry;
    BootDir *boots;
    size_t count;

    boots = NULL;
    count = 0;

    if (snprintf(bootsPath, sizeof(bootsPath), "%s/%s",
                 options->logDir, BOOTS_DIR) >=
        (int)sizeof(bootsPath)) {
        return -1;
    }

    if (!options->allBoots) {
        if (strcmp(options->boot, "current") == 0) {
            if (current_boot_id(options, bootId,
                                sizeof(bootId)) != 0) {
                return -1;
            }
        } else {
            if (!valid_boot_id(options->boot)) {
                return -1;
            }
            strcpy(bootId, options->boot);
        }

        if (snprintf(bootPath, sizeof(bootPath), "%s/%s",
                     bootsPath, bootId) >= (int)sizeof(bootPath) ||
            add_boot(&boots, &count, bootPath, bootId) != 0) {
            free(boots);
            return -1;
        }
    } else {
        directory = opendir(bootsPath);
        if (!directory) {
            return -1;
        }

        while ((entry = readdir(directory)) != NULL) {
            if (entry->d_name[0] == '.' ||
                !valid_boot_id(entry->d_name)) {
                continue;
            }
            if (snprintf(bootPath, sizeof(bootPath), "%s/%s",
                         bootsPath, entry->d_name) >=
                (int)sizeof(bootPath)) {
                continue;
            }
            (void)add_boot(&boots, &count, bootPath,
                           entry->d_name);
        }
        closedir(directory);
        qsort(boots, count, sizeof(*boots), boot_compare);
    }

    *bootsOut = boots;
    *countOut = count;
    return count > 0 ? 0 : -1;
}

static bool parse_rotated_name(const char *name, uint64_t *firstSeq) {
    unsigned long long first;
    unsigned long long last;
    char extra;

    if (!name || !firstSeq) {
        return false;
    }

    if (sscanf(name, "events-%llu-%llu.jsonl.zst%c",
               &first, &last, &extra) != 2) {
        return false;
    }

    (void)last;
    *firstSeq = (uint64_t)first;
    return true;
}

static int log_file_compare(const void *left, const void *right) {
    const LogFile *a;
    const LogFile *b;

    a = left;
    b = right;
    if (a->active != b->active) {
        return a->active ? 1 : -1;
    }
    if (a->firstSeq < b->firstSeq) {
        return -1;
    }
    if (a->firstSeq > b->firstSeq) {
        return 1;
    }
    return strcmp(a->path, b->path);
}

static int add_log_file(LogFile **files, size_t *count,
                        const char *path, uint64_t firstSeq,
                        bool active, bool compressed) {
    LogFile *items;

    items = realloc(*files, (*count + 1) * sizeof(**files));
    if (!items) {
        return -1;
    }
    *files = items;

    memset(&items[*count], 0, sizeof(items[*count]));
    if (copy_option(items[*count].path,
                    sizeof(items[*count].path), path) != 0) {
        return -1;
    }
    items[*count].firstSeq = firstSeq;
    items[*count].active = active;
    items[*count].compressed = compressed;
    (*count)++;
    return 0;
}

static int collect_log_files(const BootDir *boot, LogFile **filesOut,
                             size_t *countOut) {
    DIR *directory;
    struct dirent *entry;
    char path[PATH_MAX];
    uint64_t firstSeq;
    LogFile *files;
    size_t count;

    directory = opendir(boot->path);
    if (!directory) {
        return -1;
    }

    files = NULL;
    count = 0;
    while ((entry = readdir(directory)) != NULL) {
        bool active;
        bool compressed;

        active = strcmp(entry->d_name, ACTIVE_FILE) == 0;
        compressed = parse_rotated_name(entry->d_name, &firstSeq);
        if (!active && !compressed) {
            continue;
        }

        if (snprintf(path, sizeof(path), "%s/%s",
                     boot->path, entry->d_name) >=
            (int)sizeof(path)) {
            continue;
        }

        if (add_log_file(&files, &count, path,
                         active ? UINT64_MAX : firstSeq,
                         active, compressed) != 0) {
            closedir(directory);
            free(files);
            return -1;
        }
    }
    closedir(directory);

    qsort(files, count, sizeof(*files), log_file_compare);
    *filesOut = files;
    *countOut = count;
    return 0;
}

static bool parse_timestamp(const char *text, time_t *value) {
    char date[20];
    struct tm utc;
    char *end;

    if (!text || strlen(text) < 19 || !value) {
        return false;
    }

    memcpy(date, text, 19);
    date[19] = '\0';
    memset(&utc, 0, sizeof(utc));

    end = strptime(date, "%Y-%m-%dT%H:%M:%S", &utc);
    if (!end || *end != '\0') {
        return false;
    }

    *value = timegm(&utc);
    return true;
}

static bool record_matches(Reader *reader, json_t *record) {
    Options *options;
    const char *timestamp;
    time_t recordTime;
    uint64_t seq;

    options = reader->options;

    if (!level_list_contains(&options->levels,
                             json_string_or(record, "level", NULL)) ||
        !string_list_contains(&options->apps,
                              json_string_or(record, "app", NULL))) {
        return false;
    }

    if (*options->event &&
        strcmp(options->event,
               json_string_or(record, "event", "")) != 0) {
        return false;
    }
    if (*options->session &&
        strcmp(options->session,
               json_string_or(record, "session_id", "")) != 0) {
        return false;
    }
    if (*options->trace &&
        strcmp(options->trace,
               json_string_or(record, "trace_id", "")) != 0) {
        return false;
    }

    seq = json_u64(record, "seq");
    if (options->hasSeqStart && seq < options->seqStart) {
        return false;
    }
    if (options->hasSeqEnd && seq > options->seqEnd) {
        return false;
    }

    if (options->hasSince) {
        timestamp = json_string_or(record, "received_ts", NULL);
        if (!timestamp) {
            timestamp = json_string_or(record, "ts", NULL);
        }
        if (!parse_timestamp(timestamp, &recordTime) ||
            recordTime < options->since) {
            return false;
        }
    }

    return true;
}

static int counter_add(Counter **counters, size_t *count,
                       const char *name) {
    Counter *items;
    size_t idx;

    if (!name || !*name) {
        name = "unknown";
    }

    for (idx = 0; idx < *count; idx++) {
        if (strcmp((*counters)[idx].name, name) == 0) {
            (*counters)[idx].count++;
            return 0;
        }
    }

    items = realloc(*counters, (*count + 1) * sizeof(**counters));
    if (!items) {
        return -1;
    }
    *counters = items;

    memset(&items[*count], 0, sizeof(items[*count]));
    if (copy_option(items[*count].name,
                    sizeof(items[*count].name), name) != 0) {
        strcpy(items[*count].name, "truncated");
    }
    items[*count].count = 1;
    (*count)++;
    return 0;
}

static void summary_add(Summary *summary, json_t *record) {
    if (!summary) {
        return;
    }

    summary->records++;
    (void)counter_add(&summary->levels, &summary->levelCount,
                      json_string_or(record, "level", "unknown"));
    (void)counter_add(&summary->apps, &summary->appCount,
                      json_string_or(record, "app", "unknown"));
    (void)counter_add(&summary->events, &summary->eventCount,
                      json_string_or(record, "event", "unknown"));
}

static int counter_compare(const void *left, const void *right) {
    const Counter *a;
    const Counter *b;

    a = left;
    b = right;
    if (a->count > b->count) {
        return -1;
    }
    if (a->count < b->count) {
        return 1;
    }
    return strcmp(a->name, b->name);
}

static void print_counter(const char *title, Counter *counters,
                          size_t count) {
    size_t idx;

    qsort(counters, count, sizeof(*counters), counter_compare);
    printf("%s:\n", title);
    for (idx = 0; idx < count; idx++) {
        printf("  %-28s %llu\n", counters[idx].name,
               (unsigned long long)counters[idx].count);
    }
}

static void print_summary(Summary *summary) {
    printf("records: %llu\n",
           (unsigned long long)summary->records);
    printf("invalid: %llu\n",
           (unsigned long long)summary->invalid);
    print_counter("levels", summary->levels, summary->levelCount);
    print_counter("apps", summary->apps, summary->appCount);
    print_counter("events", summary->events, summary->eventCount);
}

static void summary_free(Summary *summary) {
    if (!summary) {
        return;
    }
    free(summary->levels);
    free(summary->apps);
    free(summary->events);
    memset(summary, 0, sizeof(*summary));
}

static void uppercase_level(const char *level, char *buffer,
                            size_t size) {
    size_t idx;

    if (!buffer || size == 0) {
        return;
    }

    if (!level) {
        level = "unknown";
    }

    for (idx = 0; idx + 1 < size && level[idx]; idx++) {
        char value;

        value = level[idx];
        if (value >= 'a' && value <= 'z') {
            value = (char)(value - 'a' + 'A');
        }
        buffer[idx] = value;
    }
    buffer[idx] = '\0';
}

static void print_message(const char *message) {
    const unsigned char *cursor;

    cursor = (const unsigned char *)message;
    while (*cursor) {
        if (*cursor == '\n' || *cursor == '\r' || *cursor == '\t') {
            putchar(' ');
        } else {
            putchar((int)*cursor);
        }
        cursor++;
    }
}

static void print_text(json_t *record) {
    const char *timestamp;
    const char *level;
    const char *app;
    const char *component;
    const char *event;
    const char *message;
    char upperLevel[16];
    uint64_t seq;

    timestamp = json_string_or(record, "ts", NULL);
    if (!timestamp) {
        timestamp = json_string_or(record, "received_ts", "-");
    }
    level = json_string_or(record, "level", "unknown");
    app = json_string_or(record, "app", "unknown");
    component = json_string_or(record, "component", "general");
    event = json_string_or(record, "event", "log_message");
    message = json_string_or(record, "msg", "");
    seq = json_u64(record, "seq");
    uppercase_level(level, upperLevel, sizeof(upperLevel));

    printf("%s %-8s %s/%s %s", timestamp, upperLevel,
           app, component, event);
    if (seq > 0) {
        printf(" seq=%llu", (unsigned long long)seq);
    }
    if (*message) {
        printf(": ");
        print_message(message);
    }
    putchar('\n');
}

static void print_jsonl(json_t *record) {
    char *line;

    line = json_dumps(record, JSON_COMPACT | JSON_ENSURE_ASCII);
    if (!line) {
        return;
    }
    printf("%s\n", line);
    free(line);
}

static int process_record(Reader *reader, json_t *record) {
    const char *bootId;
    uint64_t seq;

    bootId = json_string_or(record, "boot_id", "");
    seq = json_u64(record, "seq");

    if (reader->suppressDuplicates && *bootId) {
        if (strcmp(reader->lastBootId, bootId) != 0) {
            if (copy_option(reader->lastBootId,
                            sizeof(reader->lastBootId), bootId) != 0) {
                reader->lastBootId[0] = '\0';
            }
            reader->lastSeq = 0;
        }
        if (seq > 0 && seq <= reader->lastSeq) {
            return 0;
        }
        if (seq > reader->lastSeq) {
            reader->lastSeq = seq;
        }
    }

    if (!record_matches(reader, record)) {
        return 0;
    }

    switch (reader->options->format) {
    case FORMAT_JSONL:
        print_jsonl(record);
        break;
    case FORMAT_SUMMARY:
        summary_add(reader->summary, record);
        break;
    case FORMAT_TEXT:
    default:
        print_text(record);
        break;
    }

    if (reader->options->follow) {
        fflush(stdout);
    }
    return 0;
}

static int process_line(Reader *reader, const char *line,
                        size_t length) {
    json_error_t error;
    json_t *record;

    while (length > 0 &&
           (line[length - 1] == '\n' || line[length - 1] == '\r')) {
        length--;
    }
    if (length == 0) {
        return 0;
    }
    if (length > MAX_RECORD_BYTES) {
        reader->summary->invalid++;
        return -1;
    }

    record = json_loadb(line, length, 0, &error);
    if (!record || !json_is_object(record)) {
        if (record) {
            json_decref(record);
        }
        reader->summary->invalid++;
        fprintf(stderr, "ukama-log: invalid JSON record at line %d: %s\n",
                error.line, error.text);
        return -1;
    }

    process_record(reader, record);
    json_decref(record);
    return 0;
}

static int process_plain_file(const char *path, Reader *reader) {
    FILE *file;
    char *line;
    size_t capacity;
    ssize_t length;

    file = fopen(path, "r");
    if (!file) {
        fprintf(stderr, "ukama-log: unable to open %s: %s\n",
                path, strerror(errno));
        return -1;
    }

    line = NULL;
    capacity = 0;
    while ((length = getline(&line, &capacity, file)) >= 0) {
        (void)process_line(reader, line, (size_t)length);
    }

    free(line);
    fclose(file);
    return 0;
}

static int assembler_reserve(LineAssembler *assembler, size_t needed) {
    char *data;
    size_t capacity;

    if (needed <= assembler->capacity) {
        return 0;
    }

    capacity = assembler->capacity ? assembler->capacity : 4096;
    while (capacity < needed) {
        if (capacity >= MAX_RECORD_BYTES) {
            return -1;
        }
        capacity *= 2;
        if (capacity > MAX_RECORD_BYTES) {
            capacity = MAX_RECORD_BYTES;
        }
    }

    data = realloc(assembler->data, capacity);
    if (!data) {
        return -1;
    }
    assembler->data = data;
    assembler->capacity = capacity;
    return 0;
}

static void assembler_consume(LineAssembler *assembler,
                              const char *data, size_t size) {
    size_t offset;

    for (offset = 0; offset < size; offset++) {
        char value;

        value = data[offset];
        if (assembler->dropping) {
            if (value == '\n') {
                assembler->dropping = false;
                assembler->length = 0;
            }
            continue;
        }

        if (value == '\n') {
            (void)process_line(assembler->reader,
                               assembler->data,
                               assembler->length);
            assembler->length = 0;
            continue;
        }

        if (assembler_reserve(assembler,
                              assembler->length + 1) != 0) {
            assembler->reader->summary->invalid++;
            assembler->dropping = true;
            assembler->length = 0;
            continue;
        }
        assembler->data[assembler->length++] = value;
    }
}

static int process_zstd_file(const char *path, Reader *reader) {
    FILE *file;
    ZSTD_DStream *stream;
    void *inputData;
    void *outputData;
    size_t inputSize;
    size_t outputSize;
    LineAssembler assembler;
    int rc;

    file = fopen(path, "rb");
    if (!file) {
        fprintf(stderr, "ukama-log: unable to open %s: %s\n",
                path, strerror(errno));
        return -1;
    }

    stream = ZSTD_createDStream();
    inputSize = ZSTD_DStreamInSize();
    outputSize = ZSTD_DStreamOutSize();
    inputData = malloc(inputSize);
    outputData = malloc(outputSize);
    memset(&assembler, 0, sizeof(assembler));
    assembler.reader = reader;
    rc = 0;

    if (!stream || !inputData || !outputData ||
        ZSTD_isError(ZSTD_initDStream(stream))) {
        rc = -1;
        goto done;
    }

    while (!feof(file)) {
        ZSTD_inBuffer input;
        size_t readBytes;

        readBytes = fread(inputData, 1, inputSize, file);
        if (readBytes == 0) {
            if (ferror(file)) {
                rc = -1;
            }
            break;
        }

        input.src = inputData;
        input.size = readBytes;
        input.pos = 0;

        while (input.pos < input.size) {
            ZSTD_outBuffer output;
            size_t status;

            output.dst = outputData;
            output.size = outputSize;
            output.pos = 0;

            status = ZSTD_decompressStream(stream, &output, &input);
            if (ZSTD_isError(status)) {
                fprintf(stderr,
                        "ukama-log: invalid zstd file %s: %s\n",
                        path, ZSTD_getErrorName(status));
                rc = -1;
                goto done;
            }

            assembler_consume(&assembler, outputData, output.pos);
        }
    }

    if (assembler.length > 0 && !assembler.dropping) {
        (void)process_line(reader, assembler.data,
                           assembler.length);
    }

done:
    free(assembler.data);
    free(inputData);
    free(outputData);
    ZSTD_freeDStream(stream);
    fclose(file);
    return rc;
}

static int process_boot(const BootDir *boot, Reader *reader) {
    LogFile *files;
    size_t count;
    size_t idx;
    int rc;

    files = NULL;
    count = 0;
    if (collect_log_files(boot, &files, &count) != 0) {
        return -1;
    }

    rc = 0;
    for (idx = 0; idx < count; idx++) {
        int fileRc;

        if (files[idx].compressed) {
            fileRc = process_zstd_file(files[idx].path, reader);
        } else {
            fileRc = process_plain_file(files[idx].path, reader);
        }
        if (fileRc != 0) {
            rc = -1;
        }
    }

    free(files);
    return rc;
}

static bool same_file(const struct stat *left,
                      const struct stat *right) {
    return left->st_dev == right->st_dev &&
           left->st_ino == right->st_ino;
}

static int open_follow_file(const char *path, FILE **fileOut,
                            struct stat *statusOut) {
    FILE *file;

    file = fopen(path, "r");
    if (!file) {
        return -1;
    }
    if (fstat(fileno(file), statusOut) != 0) {
        fclose(file);
        return -1;
    }

    *fileOut = file;
    return 0;
}

static void follow_sleep(void) {
    struct timespec delay;

    delay.tv_sec = 0;
    delay.tv_nsec = FOLLOW_INTERVAL_NS;
    nanosleep(&delay, NULL);
}

static int follow_current(const Options *options, Reader *reader) {
    char path[PATH_MAX];
    FILE *file;
    struct stat openedStatus;
    char *line;
    size_t capacity;

    if (snprintf(path, sizeof(path), "%s/%s",
                 options->logDir, CURRENT_LINK) >= (int)sizeof(path)) {
        return -1;
    }

    file = NULL;
    line = NULL;
    capacity = 0;

    while (!gStop) {
        struct stat currentStatus;
        ssize_t length;

        if (!file) {
            if (open_follow_file(path, &file, &openedStatus) != 0) {
                follow_sleep();
                continue;
            }
        }

        while ((length = getline(&line, &capacity, file)) >= 0) {
            (void)process_line(reader, line, (size_t)length);
        }
        clearerr(file);

        if (stat(path, &currentStatus) == 0 &&
            !same_file(&openedStatus, &currentStatus)) {
            fclose(file);
            file = NULL;
            continue;
        }

        follow_sleep();
    }

    free(line);
    if (file) {
        fclose(file);
    }
    return 0;
}

int main(int argc, char **argv) {
    Options options;
    BootDir *boots;
    size_t bootCount;
    size_t idx;
    Summary summary;
    Reader reader;
    struct sigaction action;
    int rc;

    if (parse_options(argc, argv, &options) != 0) {
        usage(stderr);
        return EXIT_USAGE;
    }

    memset(&summary, 0, sizeof(summary));
    memset(&reader, 0, sizeof(reader));
    reader.options = &options;
    reader.summary = &summary;
    reader.suppressDuplicates = true;

    boots = NULL;
    bootCount = 0;
    if (collect_boots(&options, &boots, &bootCount) != 0) {
        fprintf(stderr, "ukama-log: no matching log boots in %s\n",
                options.logDir);
        rc = EXIT_ERROR;
        goto done;
    }

    rc = 0;
    for (idx = 0; idx < bootCount; idx++) {
        if (process_boot(&boots[idx], &reader) != 0) {
            rc = EXIT_ERROR;
        }
    }

    if (options.follow && rc == 0) {
        memset(&action, 0, sizeof(action));
        action.sa_handler = on_signal;
        sigemptyset(&action.sa_mask);
        sigaction(SIGINT, &action, NULL);
        sigaction(SIGTERM, &action, NULL);
        if (follow_current(&options, &reader) != 0) {
            rc = EXIT_ERROR;
        }
    }

    if (options.format == FORMAT_SUMMARY) {
        print_summary(&summary);
    }

done:
    free(boots);
    summary_free(&summary);
    string_list_free(&options.levels);
    string_list_free(&options.apps);
    return rc;
}
