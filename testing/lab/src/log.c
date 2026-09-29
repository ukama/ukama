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
#include <unistd.h>
#include <sys/ioctl.h>

#include "log.h"

static int g_verbose;
static int g_quiet;
static int g_progress_visible;
static double g_progress_snapshot;

/* All access is protected by stderr's stream lock, also used by normal logs. */
static void clear_progress(void) {
    if (g_progress_visible) {
        fputs("\r\033[2K", stderr);
        g_progress_visible = 0;
    }
}

void ulab_progress_clear(void) {
    flockfile(stderr);
    clear_progress();
    fflush(stderr);
    funlockfile(stderr);
}

void ulab_progress(const char *fmt, ...) {
    va_list ap;
    char line[4096];
    struct timespec now;
    struct winsize size;
    double seconds;
    int width;

    if (g_quiet) return;
    va_start(ap, fmt);
    vsnprintf(line, sizeof(line), fmt, ap);
    va_end(ap);
    flockfile(stderr);
    if (isatty(fileno(stderr))) {
        width = ioctl(fileno(stderr), TIOCGWINSZ, &size) == 0 && size.ws_col > 1 ? size.ws_col : 80;
        clear_progress();
        /* Leave the final column empty to prevent terminal line wrapping. */
        if (width > 14) fprintf(stderr, "PROGRESS   %.*s", width - 12, line);
        else fprintf(stderr, "%.*s", width - 1, line);
        g_progress_visible = 1;
    } else {
        clock_gettime(CLOCK_MONOTONIC, &now);
        seconds = now.tv_sec + now.tv_nsec / 1000000000.0;
        if (!g_progress_snapshot || seconds - g_progress_snapshot >= 60) {
            fprintf(stderr, "PROGRESS   %s\n", line);
            g_progress_snapshot = seconds;
        }
    }
    fflush(stderr);
    funlockfile(stderr);
}

#define C_RESET  "\033[0m"
#define C_RED    "\033[1;31m"
#define C_GREEN  "\033[1;32m"
#define C_YELLOW "\033[1;33m"
#define C_CYAN   "\033[1;36m"

static const char *level_color(const char *lvl) {
    if (strcmp(lvl, "ERROR") == 0) {
        return C_RED;
    }

    return "";
}

static const char *state_color(const char *state) {
    if (strcmp(state, "FAIL") == 0 || strcmp(state, "ERROR") == 0) {
        return C_RED;
    }

    if (strcmp(state, "PASS") == 0) {
        return C_GREEN;
    }

    if (strcmp(state, "SKIP") == 0 || strcmp(state, "WARN") == 0) {
        return C_YELLOW;
    }

    return C_CYAN;
}

static void vlog(const char *lvl, const char *fmt, va_list ap) {
    time_t now;
    struct tm tmv;
    char ts[32];
    const char *color;

    if (g_quiet && lvl[0] != 'E') {
        return;
    }

    now = time(NULL);
    localtime_r(&now, &tmv);
    strftime(ts, sizeof(ts), "%H:%M:%S", &tmv);

    color = level_color(lvl);
    flockfile(stderr);
    clear_progress();
    if (color[0] != '\0') {
        fprintf(stderr, "%s %s%-5s%s ", ts, color, lvl, C_RESET);
    } else {
        fprintf(stderr, "%s %-5s ", ts, lvl);
    }

    vfprintf(stderr, fmt, ap);
    fprintf(stderr, "\n");
    funlockfile(stderr);
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

    color = state_color(state);
    flockfile(stderr);
    clear_progress();
    fprintf(stderr, "%s%-10s%s ", color, state, C_RESET);

    va_start(ap, fmt);
    vfprintf(stderr, fmt, ap);
    va_end(ap);

    fprintf(stderr, "\n");
    funlockfile(stderr);
}
