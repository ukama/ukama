/*
 * Minimal stubs for Ukama platform (libusys) functions.
 * Used for local macOS development builds only — not for production.
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdarg.h>
#include <time.h>

/* log.h is the underlying logger that usys_log_* macros call into. */
#include "log.h"

/* ---- log_* implementations (subset used by controller.d) -------------- */

static int   s_level   = 0;
static char  s_service[64] = "controllerd";

const char *log_level_string(int level) {
    static const char *strings[] = { "TRACE","DEBUG","INFO","WARN","ERROR","FATAL" };
    if (level < 0 || level > 5) return "UNKNOWN";
    return strings[level];
}

void log_set_level(int level)         { s_level = level; }
void log_set_service(char *service)   { strncpy(s_service, service, sizeof(s_service) - 1); }
void log_remote_init(char *service)   { (void)service; }

void log_log(int level, const char *file, int line, const char *fmt, ...) {
    if (level < s_level) return;

    time_t t = time(NULL);
    struct tm *lt = localtime(&t);
    char ts[16];
    strftime(ts, sizeof(ts), "%H:%M:%S", lt);

    static const char *colors[] = {
        "\x1b[94m", "\x1b[36m", "\x1b[32m", "\x1b[33m", "\x1b[31m", "\x1b[35m"
    };
    const char *reset = "\x1b[0m";

    fprintf(stderr, "%s %s%-5s%s %s:%d: ", ts,
            colors[level], log_level_string(level), reset, file, line);

    va_list ap;
    va_start(ap, fmt);
    vfprintf(stderr, fmt, ap);
    va_end(ap);
    fprintf(stderr, "\n");
}

/* ---- usys_find_service_port stub --------------------------------------- */

int usys_find_service_port(char *serviceName) {
    (void)serviceName;
    return 0;
}

