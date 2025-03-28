/*
 * Copyright (c) 2020 rxi
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 */

#include "log.h"

#include <stdbool.h>

#define MAX_CALLBACKS 32

typedef struct {
    log_LogFn fn;
    void *udata;
    int level;
} Callback;

static struct {
    void *udata;
    log_LockFn lock;
    int level;
    bool quiet;
    char *service;
    int  rlogdEnable;
    Callback callbacks[MAX_CALLBACKS];
} l;

static const char *levelStrings[] = { "TRACE", "DEBUG", "INFO",
                                      "WARN",  "ERROR", "FATAL" };

#ifdef LOG_USE_COLOR
static const char *level_colors[] = { "\x1b[94m", "\x1b[36m", "\x1b[32m",
                                      "\x1b[33m", "\x1b[31m", "\x1b[35m" };
#endif

static void stdout_callback(log_Event *ev) {
    char buf[16];
    buf[strftime(buf, sizeof(buf), "%H:%M:%S", ev->time)] = '\0';
#ifdef LOG_USE_COLOR
    fprintf(ev->udata, "%s %s %s%-5s\x1b[0m \x1b[90m%s:%d:\x1b[0m ",
            l.service, buf,
            level_colors[ev->level], levelStrings[ev->level], ev->file,
            ev->line);
#else
    fprintf(ev->udata, "%s %s %-5s %s:%d: ", l.service, buf,
            levelStrings[ev->level], ev->file, ev->line);
#endif
    vfprintf(ev->udata, ev->fmt, ev->ap);
    fprintf(ev->udata, "\n");
    fflush(ev->udata);
}

static void rlogd_callback(log_Event *ev) {
    char buf[16];
    char msg[512] = {0};
    buf[strftime(buf, sizeof(buf), "%H:%M:%S", ev->time)] = '\0';
    sprintf(ev->udata, "%s %s %-5s %s:%d: ", l.service, buf,
            levelStrings[ev->level], ev->file, ev->line);
    vsprintf(&msg[0], ev->fmt, ev->ap);
    sprintf(ev->udata, "%s %s\n", ev->udata, &msg[0]);
}

static void file_callback(log_Event *ev) {
    char buf[64];
    buf[strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", ev->time)] = '\0';
    fprintf(ev->udata, "%s %s %-5s %s:%d: ", l.service, buf,
            levelStrings[ev->level], ev->file, ev->line);
    vfprintf(ev->udata, ev->fmt, ev->ap);
    fprintf(ev->udata, "\n");
    fflush(ev->udata);
}

static void lock(void) {
    if (l.lock) {
        l.lock(true, l.udata);
    }
}

static void unlock(void) {
    if (l.lock) {
        l.lock(false, l.udata);
    }
}

const char *log_level_string(int level) {
    return levelStrings[level];
}

void log_set_lock(log_LockFn fn, void *udata) {
    l.lock = fn;
    l.udata = udata;
}

void log_set_level(int level) {
    l.level = level;
}

void log_set_quiet(bool enable) {
    l.quiet = enable;
}

void log_set_service(char *service) {
    l.service = service;
    l.rlogdEnable = 0;
}

void log_enable_rlogd(int flag) {
    l.rlogdEnable = flag;
}

int log_add_callback(log_LogFn fn, void *udata, int level) {
    for (int i = 0; i < MAX_CALLBACKS; i++) {
        if (!l.callbacks[i].fn) {
            l.callbacks[i] = (Callback){ fn, udata, level };
            return 0;
        }
    }
    return -1;
}

int log_add_fp(FILE *fp, int level) {
    return log_add_callback(file_callback, fp, level);
}

static void init_event(log_Event *ev, void *udata) {
    if (!ev->time) {
        time_t t = time(NULL);
        ev->time = localtime(&t);
    }
    ev->udata = udata;
}

void log_log(int level, const char *file, int line, const char *fmt, ...) {
    log_Event ev = {
        .fmt = fmt,
        .file = file,
        .line = line,
        .level = level,
    };

    if (!is_connect_with_rlogd() && l.rlogdEnable) {
        log_remote_init(l.service);
    }

    if (is_connect_with_rlogd()) {
        char buf[512] = {0};

        lock();
        init_event(&ev, &buf[0]);
        va_start(ev.ap, fmt);
        rlogd_callback(&ev);
        va_end(ev.ap);
        log_rlogd(&buf[0]);

        unlock();
        return;
    }

    lock();

    if (!l.quiet && level >= l.level) {
        init_event(&ev, stderr);
        va_start(ev.ap, fmt);
        stdout_callback(&ev);
        va_end(ev.ap);
    }

    for (int i = 0; i < MAX_CALLBACKS && l.callbacks[i].fn; i++) {
        Callback *cb = &l.callbacks[i];
        if (level >= cb->level) {
            init_event(&ev, cb->udata);
            va_start(ev.ap, fmt);
            cb->fn(&ev);
            va_end(ev.ap);
        }
    }

    unlock();
}
