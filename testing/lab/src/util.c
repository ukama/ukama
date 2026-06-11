/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "util.h"
#include <ctype.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <unistd.h>

char *ulab_trim(char *s) {
    char *e;

    if (s == NULL) {
        return NULL;
    }
    while (isspace((unsigned char)*s)) {
        s++;
    }
    if (*s == '\0') {
        return s;
    }
    e = s + strlen(s) - 1;
    while (e > s && isspace((unsigned char)*e)) {
        *e-- = '\0';
    }
    if ((*s == '"' && *e == '"') || (*s == '\'' && *e == '\'')) {
        *e = '\0';
        s++;
    }
    return s;
}

int ulab_streq(const char *a, const char *b) {
    if (a == NULL || b == NULL) {
        return 0;
    }
    return strcmp(a, b) == 0;
}

int ulab_starts(const char *s, const char *prefix) {
    size_t n;

    if (s == NULL || prefix == NULL) {
        return 0;
    }
    n = strlen(prefix);
    return strncmp(s, prefix, n) == 0;
}

int ulab_ends(const char *s, const char *suffix) {
    size_t ns;
    size_t nf;

    if (s == NULL || suffix == NULL) {
        return 0;
    }
    ns = strlen(s);
    nf = strlen(suffix);
    if (nf > ns) {
        return 0;
    }
    return strcmp(s + ns - nf, suffix) == 0;
}

int ulab_parse_u32(const char *s, uint32_t *out) {
    char *end;
    unsigned long v;

    errno = 0;
    v = strtoul(s, &end, 10);
    if (errno != 0 || end == s || *ulab_trim(end) != '\0') {
        return ULAB_ERR;
    }
    *out = (uint32_t)v;
    return ULAB_OK;
}

int ulab_parse_u64(const char *s, uint64_t *out) {
    char *end;
    unsigned long long v;

    errno = 0;
    v = strtoull(s, &end, 10);
    if (errno != 0 || end == s || *ulab_trim(end) != '\0') {
        return ULAB_ERR;
    }
    *out = (uint64_t)v;
    return ULAB_OK;
}

int ulab_parse_double(const char *s, double *out) {
    char *end;
    double v;

    errno = 0;
    v = strtod(s, &end);
    if (errno != 0 || end == s || *ulab_trim(end) != '\0') {
        return ULAB_ERR;
    }
    *out = v;
    return ULAB_OK;
}

int ulab_copy(char *dst, size_t n, const char *src) {
    if (dst == NULL || n == 0 || src == NULL) {
        return ULAB_ERR;
    }
    if (snprintf(dst, n, "%s", src) >= (int)n) {
        return ULAB_ERR;
    }
    return ULAB_OK;
}

int ulab_mkdir_p(const char *path) {
    char tmp[ULAB_MAX_PATH];
    char *p;

    if (ulab_copy(tmp, sizeof(tmp), path) != ULAB_OK) {
        return ULAB_ERR;
    }
    for (p = tmp + 1; *p; p++) {
        if (*p == '/') {
            *p = '\0';
            if (mkdir(tmp, 0755) != 0 && errno != EEXIST) {
                return ULAB_ERR;
            }
            *p = '/';
        }
    }
    if (mkdir(tmp, 0755) != 0 && errno != EEXIST) {
        return ULAB_ERR;
    }
    return ULAB_OK;
}

uint32_t ulab_hash32(const char *s, uint32_t seed) {
    uint32_t h = 2166136261u ^ seed;

    while (*s) {
        h ^= (unsigned char)*s++;
        h *= 16777619u;
    }
    return h;
}

int ulab_within_pct(uint64_t expected, uint64_t actual, uint32_t pct) {
    uint64_t diff;
    uint64_t lim;

    if (expected == actual) {
        return 1;
    }
    diff = expected > actual ? expected - actual : actual - expected;
    lim = expected == 0 ? pct : (expected * pct) / 100;
    if (lim == 0 && pct > 0) {
        lim = 1;
    }
    return diff <= lim;
}

int ulab_run_cmd(const char *cmd, char *out, size_t out_len) {
    FILE *fp;
    int rc;
    size_t used = 0;

    if (out != NULL && out_len > 0) {
        out[0] = '\0';
    }
    fp = popen(cmd, "r");
    if (fp == NULL) {
        return ULAB_ERR;
    }
    if (out != NULL && out_len > 1) {
        while (fgets(out + used, out_len - used, fp) != NULL) {
            used = strlen(out);
            if (used + 1 >= out_len) {
                break;
            }
        }
    }
    rc = pclose(fp);
    if (rc == -1 || !WIFEXITED(rc) || WEXITSTATUS(rc) != 0) {
        return ULAB_ERR;
    }
    return ULAB_OK;
}

void ulab_json_escape(const char *in, char *out, size_t out_len) {
    size_t j = 0;

    if (out_len == 0) {
        return;
    }
    while (*in && j + 2 < out_len) {
        if (*in == '"' || *in == '\\') {
            out[j++] = '\\';
        }
        out[j++] = *in++;
    }
    out[j] = '\0';
}

const char *ulab_getenv_default(const char *name, const char *def) {
    const char *v = getenv(name);

    if (v == NULL || *v == '\0') {
        return def;
    }
    return v;
}
