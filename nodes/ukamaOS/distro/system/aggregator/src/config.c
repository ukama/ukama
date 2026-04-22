/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "config.h"
#include "toml.h"
#include "aggregator.h"

#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"

static char *dup_string(const char *value) {

    char *copy = NULL;
    size_t len = 0;

    if (value == NULL) {
        return NULL;
    }

    len = strlen(value);
    copy = calloc(1, len + 1);
    if (copy == NULL) {
        usys_log_error("failed to allocate string");
        return NULL;
    }

    memcpy(copy, value, len);
    copy[len] = '\0';

    return copy;
}

static int read_int_or_default(toml_table_t *tab, const char *key, int def) {

    toml_datum_t val = toml_int_in(tab, key);

    if (!val.ok) {
        return def;
    }

    return (int)val.u.i;
}

static int source_append(Config *config,
                         const char *name,
                         const char *url,
                         int required) {

    SourceConfig *node = NULL;
    SourceConfig *tail = NULL;

    node = calloc(1, sizeof(SourceConfig));
    if (node == NULL) {
        return RETURN_NOTOK;
    }

    node->name     = dup_string(name);
    node->url      = dup_string(url);
    node->required = required;

    if (node->name == NULL || node->url == NULL) {
        usys_free(node->name);
        usys_free(node->url);
        usys_free(node);
        return RETURN_NOTOK;
    }

    if (config->sources == NULL) {
        config->sources = node;
    } else {
        tail = config->sources;
        while (tail->next != NULL) {
            tail = tail->next;
        }
        tail->next = node;
    }

    config->sourceCount++;
    return RETURN_OK;
}

int config_load(const char *fileName, Config *config) {

    FILE *fp = NULL;
    char errBuf[200]      = {0};
    toml_table_t *root    = NULL;
    toml_array_t *sources = NULL;
    int idx = 0;

    if (fileName == NULL || config == NULL) {
        return RETURN_NOTOK;
    }

    memset(config, 0, sizeof(Config));
    config->refreshIntervalSec = 5;
    config->requestTimeoutMs   = 2000;
    config->staleGraceSec      = 15;

    if (usys_file_exist((char *)fileName) != 1) {
        usys_log_error("missing config file: %s", fileName);
        return RETURN_NOTOK;
    }

    fp = fopen(fileName, "r");
    if (fp == NULL) {
        usys_log_error("failed to open config file: %s", fileName);
        return RETURN_NOTOK;
    }

    root = toml_parse_file(fp, errBuf, sizeof(errBuf));
    fclose(fp);
    if (root == NULL) {
        usys_log_error("failed to parse config: %s", errBuf);
        return RETURN_NOTOK;
    }

    config->refreshIntervalSec =
        read_int_or_default(root,
                            TAG_REFRESH_INTERVAL_SEC,
                            config->refreshIntervalSec);
    config->requestTimeoutMs =
        read_int_or_default(root,
                            TAG_REQUEST_TIMEOUT_MS,
                            config->requestTimeoutMs);
    config->staleGraceSec =
        read_int_or_default(root,
                            TAG_STALE_GRACE_SEC,
                            config->staleGraceSec);

    sources = toml_array_in(root, TAG_SOURCE);
    if (sources == NULL) {
        usys_log_error("config is missing [[source]] entries");
        toml_free(root);
        return RETURN_NOTOK;
    }

    for (idx = 0; idx < toml_array_nelem(sources); idx++) {

        toml_table_t *src = NULL;
        toml_datum_t name;
        toml_datum_t url;
        toml_datum_t required;
        int req = 0;

        src = toml_table_at(sources, idx);
        if (src == NULL) {
            continue;
        }

        name     = toml_string_in(src, TAG_NAME);
        url      = toml_string_in(src, TAG_URL);
        required = toml_bool_in(src,   TAG_REQUIRED);

        if (!name.ok || !url.ok) {
            if (name.ok) free(name.u.s);
            if (url.ok)  free(url.u.s);
            usys_log_error("source entry %d missing name/url", idx);
            toml_free(root);
            return RETURN_NOTOK;
        }

        if (required.ok) {
            req = required.u.b ? 1 : 0;
        }

        if (source_append(config, name.u.s, url.u.s, req) != RETURN_OK) {
            free(name.u.s);
            free(url.u.s);
            toml_free(root);
            return RETURN_NOTOK;
        }

        free(name.u.s);
        free(url.u.s);
    }

    toml_free(root);

    if (config->sourceCount <= 0) {
        usys_log_error("no sources configured");
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}

void config_free(Config *config) {

    SourceConfig *src = NULL;
    SourceConfig *next = NULL;

    if (config == NULL) return;

    src = config->sources;
    while (src != NULL) {
        next = src->next;
        usys_free(src->name);
        usys_free(src->url);
        usys_free(src);
        src = next;
    }

    memset(config, 0, sizeof(Config));
}
