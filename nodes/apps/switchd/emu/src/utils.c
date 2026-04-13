/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "utils.h"

int write_all(int fd, const void *buf, size_t len) {
    const char *ptr = (const char *)buf;

    while (len > 0) {
        ssize_t written = write(fd, ptr, len);
        if (written <= 0) {
            return STATUS_NOK;
        }

        ptr += written;
        len -= (size_t)written;
    }

    return STATUS_OK;
}

static int find_key(const char *json, const char *key, const char **out) {
    char pattern[128] = {0};

    (void)snprintf(pattern, sizeof(pattern), "\"%s\"", key);
    *out = strstr(json, pattern);

    return (*out != NULL) ? STATUS_OK : STATUS_NOK;
}

int json_get_string_field(const char *json, const char *key,
                          char *buf, size_t bufLen) {
    const char *start = NULL;
    const char *end   = NULL;
    size_t len        = 0;

    if (json == NULL || key == NULL || buf == NULL || bufLen == 0) {
        return STATUS_NOK;
    }

    if (find_key(json, key, &start) != STATUS_OK) {
        return STATUS_NOK;
    }

    start = strchr(start, ':');
    if (start == NULL) {
        return STATUS_NOK;
    }

    while (*start != '\0' && *start != '"') {
        start++;
    }

    if (*start != '"') {
        return STATUS_NOK;
    }

    start++;
    end = strchr(start, '"');
    if (end == NULL) {
        return STATUS_NOK;
    }

    len = (size_t)(end - start);
    if (len >= bufLen) {
        len = bufLen - 1;
    }

    memcpy(buf, start, len);
    buf[len] = '\0';

    return STATUS_OK;
}

int json_get_int_field(const char *json, const char *key, int *value) {
    const char *start = NULL;

    if (json == NULL || key == NULL || value == NULL) {
        return STATUS_NOK;
    }

    if (find_key(json, key, &start) != STATUS_OK) {
        return STATUS_NOK;
    }

    start = strchr(start, ':');
    if (start == NULL) {
        return STATUS_NOK;
    }

    start++;
    while (*start != '\0' && isspace((unsigned char)*start)) {
        start++;
    }

    *value = atoi(start);

    return STATUS_OK;
}

int json_get_bool_field(const char *json, const char *key, int *value) {
    const char *start = NULL;

    if (json == NULL || key == NULL || value == NULL) {
        return STATUS_NOK;
    }

    if (find_key(json, key, &start) != STATUS_OK) {
        return STATUS_NOK;
    }

    start = strchr(start, ':');
    if (start == NULL) {
        return STATUS_NOK;
    }

    start++;
    while (*start != '\0' && isspace((unsigned char)*start)) {
        start++;
    }

    *value = (strncmp(start, "true", 4) == 0) ? 1 : 0;

    return STATUS_OK;
}

const char *bool_to_json(int value) {
    return value ? "true" : "false";
}
