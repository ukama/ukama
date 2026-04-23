/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include "aggregator.h"
#include "config.h"

#include "usys_log.h"

typedef struct {
    char *buf;
    size_t len;
} MemBuf;

static pthread_once_t curlOnce = PTHREAD_ONCE_INIT;

static void curl_init_once(void) {

    curl_global_init(CURL_GLOBAL_ALL);
}

static size_t write_cb(void *contents, size_t size, size_t nmemb, void *data) {

    char *next = NULL;
    MemBuf *buf = NULL;
    size_t total = 0;

    buf = (MemBuf *)data;
    total = size * nmemb;

    next = realloc(buf->buf, buf->len + total + 1);
    if (next == NULL) {
        return 0;
    }

    buf->buf = next;
    memcpy(buf->buf + buf->len, contents, total);
    buf->len += total;
    buf->buf[buf->len] = '\0';

    return total;
}

static char *dup_string(const char *value) {

    char *copy = NULL;
    size_t len = 0;

    if (value == NULL) {
        return NULL;
    }

    len = strlen(value);
    copy = calloc(1, len + 1);
    if (copy == NULL) {
        return NULL;
    }

    memcpy(copy, value, len);
    copy[len] = '\0';

    return copy;
}

static int fetch_metrics(const char *url,
                         int timeoutMs,
                         int *httpCode,
                         char **body) {

    CURL *curl = NULL;
    CURLcode rc;
    MemBuf buf = {0};
    long responseCode = 0;

    if (url == NULL || httpCode == NULL || body == NULL) {
        return RETURN_NOTOK;
    }

    *httpCode = 0;
    *body = NULL;

    pthread_once(&curlOnce, curl_init_once);

    curl = curl_easy_init();
    if (curl == NULL) {
        return RETURN_NOTOK;
    }

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, (long)timeoutMs);
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, (long)timeoutMs);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &buf);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "aggregator.d");

    rc = curl_easy_perform(curl);
    if (rc != CURLE_OK) {
        usys_log_error("fetch failed url=%s rc=%d err=%s",
                       url, rc, curl_easy_strerror(rc));
        curl_easy_cleanup(curl);
        free(buf.buf);
        return RETURN_NOTOK;
    }

    rc = curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &responseCode);
    if (rc != CURLE_OK) {
        usys_log_error("getinfo failed url=%s rc=%d err=%s",
                       url, rc, curl_easy_strerror(rc));
        curl_easy_cleanup(curl);
        free(buf.buf);
        return RETURN_NOTOK;
    }

    curl_easy_cleanup(curl);

    *httpCode = (int)responseCode;

    if (*httpCode < 200 || *httpCode >= 300) {
        free(buf.buf);
        return RETURN_NOTOK;
    }

    if (buf.buf == NULL) {
        buf.buf = dup_string("");
        if (buf.buf == NULL) {
            return RETURN_NOTOK;
        }
    }

    *body = buf.buf;
    return RETURN_OK;
}

static int seen_metric(char **names, int count, const char *name) {

    int idx = 0;

    for (idx = 0; idx < count; idx++) {
        if (strcmp(names[idx], name) == 0) {
            return 1;
        }
    }

    return 0;
}

static int add_seen_metric(char ***names, int *count, const char *name) {

    char **next = NULL;
    char *copy = NULL;

    copy = dup_string(name);
    if (copy == NULL) {
        return RETURN_NOTOK;
    }

    next = realloc(*names, sizeof(char *) * (*count + 1));
    if (next == NULL) {
        free(copy);
        return RETURN_NOTOK;
    }

    *names = next;
    (*names)[*count] = copy;
    (*count)++;

    return RETURN_OK;
}

static void free_seen_metrics(char **names, int count) {

    int idx = 0;

    if (names == NULL) {
        return;
    }

    for (idx = 0; idx < count; idx++) {
        free(names[idx]);
    }

    free(names);
}

static int append_text(char **buf, size_t *len, const char *text) {

    char *next = NULL;
    size_t add = 0;

    if (text == NULL) {
        return RETURN_OK;
    }

    add = strlen(text);
    next = realloc(*buf, *len + add + 1);
    if (next == NULL) {
        return RETURN_NOTOK;
    }

    *buf = next;
    memcpy(*buf + *len, text, add);
    *len += add;
    (*buf)[*len] = '\0';

    return RETURN_OK;
}

static int append_line(char **buf, size_t *len, const char *line) {

    if (append_text(buf, len, line) != RETURN_OK) {
        return RETURN_NOTOK;
    }

    if (append_text(buf, len, "\n") != RETURN_OK) {
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}

static int merge_source_body(char **buf,
                             size_t *len,
                             const char *body,
                             char ***metaNames,
                             int *metaCount) {

    const char *cur = NULL;

    if (body == NULL || *body == '\0') {
        return RETURN_OK;
    }

    cur = body;
    while (*cur != '\0') {
        const char *end = strchr(cur, '\n');
        char line[4096] = {0};
        size_t lineLen = 0;

        if (end == NULL) {
            lineLen = strlen(cur);
        } else {
            lineLen = (size_t)(end - cur);
        }

        if (lineLen >= sizeof(line)) {
            lineLen = sizeof(line) - 1;
        }

        memcpy(line, cur, lineLen);
        line[lineLen] = '\0';

        if (strncmp(line, "# HELP ", 7) == 0 || strncmp(line, "# TYPE ", 7) == 0) {
            char metricName[256] = {0};
            const char *p = strchr(line + 7, ' ');
            size_t nameLen = 0;

            if (p != NULL) {
                nameLen = (size_t)(p - (line + 7));
                if (nameLen >= sizeof(metricName)) {
                    nameLen = sizeof(metricName) - 1;
                }
                memcpy(metricName, line + 7, nameLen);
                metricName[nameLen] = '\0';

                if (!seen_metric(*metaNames, *metaCount, metricName)) {
                    if (add_seen_metric(metaNames, metaCount, metricName) != RETURN_OK) {
                        return RETURN_NOTOK;
                    }
                    if (append_line(buf, len, line) != RETURN_OK) {
                        return RETURN_NOTOK;
                    }
                }
            }
        } else if (append_line(buf, len, line) != RETURN_OK) {
            return RETURN_NOTOK;
        }

        if (end == NULL) {
            break;
        }

        cur = end + 1;
    }

    return RETURN_OK;
}

static char *build_health_metrics_locked(AppState *state) {

    char *buf = NULL;
    size_t len = 0;
    time_t now = 0;
    int idx = 0;
    char line[512] = {0};

    now = time(NULL);

    append_line(&buf, &len, "# HELP aggregator_up Aggregator process health.");
    append_line(&buf, &len, "# TYPE aggregator_up gauge");
    append_line(&buf, &len, "aggregator_up 1");
    append_line(&buf, &len, "# HELP aggregator_snapshot_age_seconds Age of current merged snapshot.");
    append_line(&buf, &len, "# TYPE aggregator_snapshot_age_seconds gauge");
    snprintf(line, sizeof(line), "aggregator_snapshot_age_seconds %ld",
             (long)((state->snapshotAt > 0) ? (now - state->snapshotAt) : -1));
    append_line(&buf, &len, line);

    append_line(&buf, &len, "# HELP aggregator_source_up Source scrape success, 1=up 0=down.");
    append_line(&buf, &len, "# TYPE aggregator_source_up gauge");
    append_line(&buf, &len, "# HELP aggregator_source_age_seconds Seconds since last successful scrape by source.");
    append_line(&buf, &len, "# TYPE aggregator_source_age_seconds gauge");
    append_line(&buf, &len, "# HELP aggregator_source_errors_total Total scrape failures by source.");
    append_line(&buf, &len, "# TYPE aggregator_source_errors_total counter");

    for (idx = 0; idx < state->sourceCount; idx++) {
        long age = -1;

        if (state->sources[idx].lastSuccess > 0) {
            age = now - state->sources[idx].lastSuccess;
        }

        snprintf(line, sizeof(line),
                 "aggregator_source_up{source=\"%s\"} %d",
                 state->sources[idx].name,
                 state->sources[idx].up);
        append_line(&buf, &len, line);

        snprintf(line, sizeof(line),
                 "aggregator_source_age_seconds{source=\"%s\"} %ld",
                 state->sources[idx].name,
                 age);
        append_line(&buf, &len, line);

        snprintf(line, sizeof(line),
                 "aggregator_source_errors_total{source=\"%s\"} %d",
                 state->sources[idx].name,
                 state->sources[idx].errorCount);
        append_line(&buf, &len, line);
    }

    return buf;
}

static char *build_snapshot_locked(AppState *state) {

    char *buf = NULL;
    char *health = NULL;
    char **metaNames = NULL;
    size_t len = 0;
    int metaCount = 0;
    int idx = 0;
    time_t now = 0;

    now = time(NULL);

    for (idx = 0; idx < state->sourceCount; idx++) {
        long age = -1;

        if (state->sources[idx].body == NULL) {
            continue;
        }

        if (state->sources[idx].lastSuccess > 0) {
            age = now - state->sources[idx].lastSuccess;
        }

        if (age < 0 || age > state->staleGraceSec) {
            continue;
        }

        if (merge_source_body(&buf, &len, state->sources[idx].body,
                              &metaNames, &metaCount) != RETURN_OK) {
            free(buf);
            free_seen_metrics(metaNames, metaCount);
            return NULL;
        }
    }

    health = build_health_metrics_locked(state);
    if (health == NULL) {
        free(buf);
        free_seen_metrics(metaNames, metaCount);
        return NULL;
    }

    if (append_text(&buf, &len, health) != RETURN_OK) {
        free(health);
        free(buf);
        free_seen_metrics(metaNames, metaCount);
        return NULL;
    }

    free(health);
    free_seen_metrics(metaNames, metaCount);
    return buf;
}

static void refresh_once(AppState *state) {

    int idx = 0;
    time_t now = 0;

    now = time(NULL);

    usys_log_info("refresh_once: starting");

    for (idx = 0; idx < state->sourceCount; idx++) {
        int httpCode = 0;
        char *body = NULL;
        int rc = RETURN_NOTOK;

        usys_log_info("refresh_once: fetching source=%s url=%s",
                      state->sources[idx].name,
                      state->sources[idx].url);

        rc = fetch_metrics(state->sources[idx].url,
                           state->requestTimeoutMs,
                           &httpCode,
                           &body);

        pthread_mutex_lock(&state->mutex);
        state->sources[idx].lastAttempt = now;
        state->sources[idx].lastHttpCode = httpCode;

        if (rc == RETURN_OK) {
            free(state->sources[idx].body);
            state->sources[idx].body = body;
            state->sources[idx].up = 1;
            state->sources[idx].lastSuccess = now;
            usys_log_info("refresh_once: source=%s ok http=%d",
                          state->sources[idx].name, httpCode);
        } else {
            state->sources[idx].up = 0;
            state->sources[idx].errorCount++;
            free(body);
            usys_log_error("refresh_once: source=%s failed http=%d",
                           state->sources[idx].name, httpCode);
        }
        pthread_mutex_unlock(&state->mutex);
    }

    pthread_mutex_lock(&state->mutex);
    free(state->snapshot);
    state->snapshot = build_snapshot_locked(state);
    state->snapshotAt = now;
    pthread_mutex_unlock(&state->mutex);

    usys_log_info("refresh_once: snapshot updated");
}

static void *refresh_thread(void *arg) {

    AppState *state = NULL;

    state = (AppState *)arg;
    if (state == NULL) {
        return NULL;
    }

    usys_log_info("refresh thread started");

    while (state->running == 1) {
        refresh_once(state);
        sleep(state->refreshIntervalSec);
    }

    usys_log_info("refresh thread stopped");

    return NULL;
}

int app_state_init(AppState *state, const Config *config) {

    SourceConfig *src = NULL;
    int idx = 0;

    if (state == NULL || config == NULL || config->sourceCount <= 0) {
        return RETURN_NOTOK;
    }

    memset(state, 0, sizeof(AppState));
    pthread_mutex_init(&state->mutex, NULL);

    state->sources = calloc(config->sourceCount, sizeof(SourceState));
    if (state->sources == NULL) {
        pthread_mutex_destroy(&state->mutex);
        return RETURN_NOTOK;
    }

    state->sourceCount = config->sourceCount;
    state->refreshIntervalSec = config->refreshIntervalSec;
    state->requestTimeoutMs = config->requestTimeoutMs;
    state->staleGraceSec = config->staleGraceSec;
    state->running = 0;

    src = config->sources;
    while (src != NULL && idx < config->sourceCount) {
        snprintf(state->sources[idx].name,
                 sizeof(state->sources[idx].name),
                 "%s",
                 src->name);
        state->sources[idx].url = dup_string(src->url);
        state->sources[idx].required = src->required;
        if (state->sources[idx].url == NULL) {
            app_state_cleanup(state);
            return RETURN_NOTOK;
        }
        src = src->next;
        idx++;
    }

    return RETURN_OK;
}

void app_state_cleanup(AppState *state) {

    int idx = 0;

    if (state == NULL) {
        return;
    }

    for (idx = 0; idx < state->sourceCount; idx++) {
        free(state->sources[idx].url);
        free(state->sources[idx].body);
    }

    free(state->sources);
    free(state->snapshot);
    pthread_mutex_destroy(&state->mutex);
    memset(state, 0, sizeof(AppState));
}

int app_state_start(AppState *state) {

    if (state == NULL) {
        return RETURN_NOTOK;
    }

    state->running = 1;

    if (pthread_create(&state->thread, NULL, refresh_thread, state) != 0) {
        state->running = 0;
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}

void app_state_stop(AppState *state) {

    if (state == NULL) {
        return;
    }

    if (state->running == 1) {
        state->running = 0;
        pthread_join(state->thread, NULL);
    }
}

char *app_state_dup_snapshot(AppState *state) {

    char *copy = NULL;

    if (state == NULL) {
        return NULL;
    }

    pthread_mutex_lock(&state->mutex);
    if (state->snapshot != NULL) {
        copy = dup_string(state->snapshot);
    } else {
        copy = dup_string("# aggregator snapshot pending\naggregator_up 1\n");
    }
    pthread_mutex_unlock(&state->mutex);

    return copy;
}

char *app_state_status_json(AppState *state) {

    char *buf = NULL;
    size_t len = 0;
    char line[512] = {0};
    int idx = 0;
    time_t now = 0;

    if (state == NULL) {
        return NULL;
    }

    now = time(NULL);

    pthread_mutex_lock(&state->mutex);
    append_text(&buf, &len, "{\"status\":\"ok\",\"sources\":[");

    for (idx = 0; idx < state->sourceCount; idx++) {
        long age = -1;

        if (idx > 0) {
            append_text(&buf, &len, ",");
        }

        if (state->sources[idx].lastSuccess > 0) {
            age = now - state->sources[idx].lastSuccess;
        }

        snprintf(line, sizeof(line),
                 "{\"name\":\"%s\",\"up\":%d,\"required\":%d,\"httpCode\":%d,\"ageSec\":%ld,\"errors\":%d}",
                 state->sources[idx].name,
                 state->sources[idx].up,
                 state->sources[idx].required,
                 state->sources[idx].lastHttpCode,
                 age,
                 state->sources[idx].errorCount);
        append_text(&buf, &len, line);
    }

    append_text(&buf, &len, "]}");
    pthread_mutex_unlock(&state->mutex);

    return buf;
}
