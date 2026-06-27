/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdarg.h>
#include <string.h>
#include <time.h>

#include "log.h"

static int g_verbose;
static int g_quiet;

static void vlog(const char *lvl, const char *fmt, va_list ap) {
    time_t now;
    struct tm tmv;
    char ts[32];

    if (g_quiet && lvl[0] != 'E') {
        return;
    }

    now = time(NULL);
    localtime_r(&now, &tmv);
    strftime(ts, sizeof(ts), "%H:%M:%S", &tmv);
    fprintf(stderr, "%s %-5s ", ts, lvl);
    vfprintf(stderr, fmt, ap);
    fprintf(stderr, "\n");
}

void ulab_log_set_verbose(int verbose) {
    g_verbose = verbose;
}

void ulab_log_set_quiet(int quiet) {
    g_quiet = quiet;
}

void ulab_log_debug(const char *fmt, ...) {
    va_list ap;

    if (!g_verbose) {
        return;
    }
    va_start(ap, fmt);
    vlog("DEBUG", fmt, ap);
    va_end(ap);
}

void ulab_log_info(const char *fmt, ...) {
    va_list ap;

    va_start(ap, fmt);
    vlog("INFO", fmt, ap);
    va_end(ap);
}

void ulab_log_warn(const char *fmt, ...) {
    va_list ap;

    va_start(ap, fmt);
    vlog("WARN", fmt, ap);
    va_end(ap);
}

void ulab_log_error(const char *fmt, ...) {
    va_list ap;

    va_start(ap, fmt);
    vlog("ERROR", fmt, ap);
    va_end(ap);
}

void ulab_status(const char *state, const char *fmt, ...) {
    va_list ap;
    const char *color;

    if (g_quiet) {
        return;
    }

    if (strcmp(state, "FAIL") == 0 || strcmp(state, "ERROR") == 0) {
        color = "\033[1;31m";
    } else if (strcmp(state, "PASS") == 0) {
        color = "\033[1;32m";
    } else if (strcmp(state, "SKIP") == 0) {
        color = "\033[1;33m";
    } else {
        color = "\033[1;36m";
    }

    fprintf(stderr, "%s%-8s\033[0m ", color, state);

    va_start(ap, fmt);
    vfprintf(stderr, fmt, ap);
    va_end(ap);

    fprintf(stderr, "\n");
}

