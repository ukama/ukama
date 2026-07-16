/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <jansson.h>

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "hub.h"
#include "util.h"

typedef struct {
    char *data;
    size_t len;
} hub_buffer_t;

static size_t write_response(void *ptr, size_t size, size_t count, void *arg) {
    hub_buffer_t *buffer;
    char *next;
    size_t bytes;

    buffer = arg;
    bytes = size * count;
    next = realloc(buffer->data, buffer->len + bytes + 1);
    if (next == NULL) {
        return 0;
    }

    buffer->data = next;
    memcpy(buffer->data + buffer->len, ptr, bytes);
    buffer->len += bytes;
    buffer->data[buffer->len] = '\0';

    return bytes;
}

static int safe_component(const char *value) {
    const unsigned char *p;

    if (value == NULL || value[0] == '\0') {
        return 0;
    }

    p = (const unsigned char *)value;
    while (*p != '\0') {
        if (!( (*p >= 'a' && *p <= 'z') ||
               (*p >= 'A' && *p <= 'Z') ||
               (*p >= '0' && *p <= '9') ||
               *p == '-' || *p == '_' || *p == '.' || *p == '+')) {
            return 0;
        }
        p++;
    }

    return 1;
}

static int json_u64(json_t *value, uint64_t *out) {
    const char *text;
    char *end;
    unsigned long long number;

    if (value == NULL || out == NULL) {
        return ULAB_ERR;
    }

    if (json_is_integer(value)) {
        if (json_integer_value(value) < 0) {
            return ULAB_ERR;
        }
        *out = (uint64_t)json_integer_value(value);
        return ULAB_OK;
    }

    if (!json_is_string(value)) {
        return ULAB_ERR;
    }

    text = json_string_value(value);
    if (text == NULL || text[0] == '\0') {
        return ULAB_ERR;
    }

    errno = 0;
    end = NULL;
    number = strtoull(text, &end, 10);
    if (errno != 0 || end == text || *end != '\0') {
        return ULAB_ERR;
    }

    *out = (uint64_t)number;
    return ULAB_OK;
}

static int request_app(hub_client_t *hub,
                       const char *app,
                       json_t **response,
                       int *not_found,
                       ulab_error_t *err) {
    CURL *curl;
    CURLcode curl_rc;
    struct curl_slist *headers;
    hub_buffer_t body;
    char url[ULAB_MAX_URL * 2];
    long status;
    int n;
    json_error_t json_error;

    if (hub == NULL || response == NULL || not_found == NULL) {
        snprintf(err->msg, sizeof(err->msg), "invalid Hub request");
        return ULAB_ERR;
    }

    *response = NULL;
    *not_found = 0;
    body.data = NULL;
    body.len = 0;
    headers = NULL;
    status = 0;

    n = snprintf(url, sizeof(url), "%s/v1/hub/app/%s", hub->url, app);
    if (n < 0 || (size_t)n >= sizeof(url)) {
        snprintf(err->msg, sizeof(err->msg), "Hub URL too long");
        return ULAB_ERR;
    }

    curl = curl_easy_init();
    if (curl == NULL) {
        snprintf(err->msg, sizeof(err->msg), "Hub curl init failed");
        return ULAB_ERR;
    }

    headers = curl_slist_append(headers, "Accept: application/json");
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_response);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &body);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT, 10L);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 30L);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "ukama-lab");

    curl_rc = curl_easy_perform(curl);
    if (curl_rc == CURLE_OK) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &status);
    }

    if (hub->logf != NULL) {
        fprintf(hub->logf, "GET %s\n", url);
        fprintf(hub->logf, "HTTP %ld\n", status);
        fprintf(hub->logf, "%s\n\n", body.data ? body.data : "");
        fflush(hub->logf);
    }

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);

    if (curl_rc != CURLE_OK) {
        snprintf(err->msg, sizeof(err->msg),
                 "Hub request failed: %s", curl_easy_strerror(curl_rc));
        free(body.data);
        return ULAB_ERR;
    }

    if (status == 404) {
        *not_found = 1;
        free(body.data);
        return ULAB_OK;
    }

    if (status < 200 || status >= 300) {
        snprintf(err->msg, sizeof(err->msg),
                 "Hub returned HTTP %ld for app %.128s: %.384s",
                 status, app, body.data ? body.data : "");
        free(body.data);
        return ULAB_ERR;
    }

    *response = json_loads(body.data ? body.data : "", 0, &json_error);
    free(body.data);
    if (*response == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "Hub returned invalid JSON for app %.128s: %.256s",
                 app, json_error.text);
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int hub_init(hub_client_t *hub,
             const char *url,
             const char *run_dir,
             ulab_error_t *err) {
    char path[ULAB_MAX_PATH * 2];
    size_t len;

    if (hub == NULL || url == NULL || url[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "Hub URL is required");
        return ULAB_ERR;
    }

    memset(hub, 0, sizeof(*hub));
    if (!ulab_starts(url, "http://") && !ulab_starts(url, "https://")) {
        snprintf(err->msg, sizeof(err->msg), "invalid Hub URL: %s", url);
        return ULAB_ERR;
    }

    if (ulab_copy(hub->url, sizeof(hub->url), url)) {
        snprintf(err->msg, sizeof(err->msg), "Hub URL too long");
        return ULAB_ERR;
    }

    len = strlen(hub->url);
    while (len > 0 && hub->url[len - 1] == '/') {
        hub->url[--len] = '\0';
    }

    if (curl_global_init(CURL_GLOBAL_ALL) != CURLE_OK) {
        snprintf(err->msg, sizeof(err->msg), "Hub curl global init failed");
        return ULAB_ERR;
    }
    hub->curl_ready = 1;

    if (run_dir != NULL && run_dir[0] != '\0') {
        snprintf(path, sizeof(path), "%s/hub.log", run_dir);
        hub->logf = fopen(path, "w");
    }

    return ULAB_OK;
}

void hub_close(hub_client_t *hub) {
    if (hub == NULL) {
        return;
    }

    if (hub->logf != NULL) {
        fclose(hub->logf);
        hub->logf = NULL;
    }

    if (hub->curl_ready) {
        curl_global_cleanup();
        hub->curl_ready = 0;
    }
}

int hub_tar_version_exists(hub_client_t *hub,
                           const char *app,
                           const char *version,
                           int *exists,
                           uint64_t *size_bytes,
                           char *artifact_url,
                           size_t artifact_url_len,
                           ulab_error_t *err) {
    json_t *root;
    json_t *versions;
    json_t *version_item;
    json_t *formats;
    json_t *format;
    const char *item_version;
    const char *type;
    const char *url;
    uint64_t size;
    size_t i;
    size_t j;
    int not_found;
    int version_found;

    if (exists == NULL || size_bytes == NULL || artifact_url == NULL ||
        artifact_url_len == 0 || !safe_component(app) ||
        !safe_component(version)) {
        snprintf(err->msg, sizeof(err->msg),
                 "invalid Hub app/version request");
        return ULAB_ERR;
    }

    *exists = 0;
    *size_bytes = 0;
    artifact_url[0] = '\0';
    root = NULL;
    not_found = 0;
    version_found = 0;

    if (request_app(hub, app, &root, &not_found, err)) {
        return ULAB_ERR;
    }
    if (not_found) {
        return ULAB_OK;
    }

    versions = json_object_get(root, "versions");
    if (versions == NULL || !json_is_array(versions)) {
        snprintf(err->msg, sizeof(err->msg),
                 "Hub response missing versions for app %s", app);
        json_decref(root);
        return ULAB_ERR;
    }

    for (i = 0; i < json_array_size(versions); i++) {
        version_item = json_array_get(versions, i);
        if (version_item == NULL || !json_is_object(version_item)) {
            continue;
        }

        item_version = json_string_value(
            json_object_get(version_item, "version"));
        if (item_version == NULL || !ulab_streq(item_version, version)) {
            continue;
        }

        version_found = 1;
        formats = json_object_get(version_item, "FormatInfo");
        if (formats == NULL || !json_is_array(formats)) {
            break;
        }

        for (j = 0; j < json_array_size(formats); j++) {
            format = json_array_get(formats, j);
            if (format == NULL || !json_is_object(format)) {
                continue;
            }

            type = json_string_value(json_object_get(format, "Type"));
            url = json_string_value(json_object_get(format, "Url"));
            if (type == NULL || url == NULL ||
                !ulab_streq(type, "tar.gz") || !ulab_ends(url, ".tar.gz") ||
                json_u64(json_object_get(format, "Size"), &size) ||
                size == 0) {
                continue;
            }

            if (ulab_copy(artifact_url, artifact_url_len, url)) {
                snprintf(err->msg, sizeof(err->msg),
                         "Hub artifact URL too long for %s:%s", app, version);
                json_decref(root);
                return ULAB_ERR;
            }

            *size_bytes = size;
            *exists = 1;
            json_decref(root);
            return ULAB_OK;
        }
        break;
    }

    json_decref(root);

    if (version_found) {
        snprintf(err->msg, sizeof(err->msg),
                 "Hub app %s version %s has no valid tar.gz artifact",
                 app, version);
        return ULAB_ERR;
    }

    return ULAB_OK;
}
