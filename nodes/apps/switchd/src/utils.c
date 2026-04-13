/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <ctype.h>
#include <errno.h>
#include <fcntl.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <time.h>
#include <unistd.h>

#include "utils.h"
#include "types.h"

char *trim_ws(char *s) {
    char *end;

    if (s == NULL) {
        return NULL;
    }

    while (*s && isspace((unsigned char)*s)) {
        s++;
    }

    if (*s == '\0') {
        return s;
    }

    end = s + strlen(s) - 1;
    while (end > s && isspace((unsigned char)*end)) {
        *end = '\0';
        end--;
    }

    return s;
}

int safe_snprintf(char *dst, size_t dstLen, const char *fmt, ...) {
    va_list ap;
    int written;

    if (dst == NULL || dstLen == 0) {
        return -1;
    }

    va_start(ap, fmt);
    written = vsnprintf(dst, dstLen, fmt, ap);
    va_end(ap);

    if (written < 0 || (size_t)written >= dstLen) {
        dst[dstLen - 1] = '\0';
        return -1;
    }

    return written;
}

double parse_prefixed_double(const char *s) {
    char *end;

    if (s == NULL) {
        return 0.0;
    }

    while (*s && !(isdigit((unsigned char)*s) || *s == '-' || *s == '+')) {
        s++;
    }

    return strtod(s, &end);
}

uint64_t monotonic_msec(void) {
    struct timespec ts;

    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (uint64_t)ts.tv_sec * 1000ULL +
           (uint64_t)(ts.tv_nsec / 1000000ULL);
}

int mkdir_p(const char *path, int mode) {
    char tmp[512];
    char *p;
    size_t len;

    if (path == NULL || *path == '\0') {
        return -1;
    }

    snprintf(tmp, sizeof(tmp), "%s", path);
    len = strlen(tmp);
    if (tmp[len - 1] == '/') {
        tmp[len - 1] = '\0';
    }

    for (p = tmp + 1; *p; p++) {
        if (*p == '/') {
            *p = '\0';
            if (mkdir(tmp, (mode_t)mode) != 0 && errno != EEXIST) {
                return -1;
            }
            *p = '/';
        }
    }

    if (mkdir(tmp, (mode_t)mode) != 0 && errno != EEXIST) {
        return -1;
    }

    return 0;
}

int copy_file(const char *src, const char *dst) {
    int inFd;
    int outFd;
    ssize_t bytesRead;
    char buffer[8192];

    inFd = open(src, O_RDONLY);
    if (inFd < 0) {
        return -1;
    }

    outFd = open(dst, O_WRONLY | O_CREAT | O_TRUNC, 0644);
    if (outFd < 0) {
        close(inFd);
        return -1;
    }

    while ((bytesRead = read(inFd, buffer, sizeof(buffer))) > 0) {
        if (write(outFd, buffer, (size_t)bytesRead) != bytesRead) {
            close(inFd);
            close(outFd);
            return -1;
        }
    }

    close(inFd);
    close(outFd);
    return (bytesRead < 0) ? -1 : 0;
}

const char *state_to_str(int state) {
    switch (state) {
    case SWITCHD_STATE_INIT:
        return "init";
    case SWITCHD_STATE_READY:
        return "ready";
    case SWITCHD_STATE_BUSY:
        return "busy";
    case SWITCHD_STATE_DEGRADED:
        return "degraded";
    case SWITCHD_STATE_UPDATING:
        return "updating";
    case SWITCHD_STATE_RECOVERING:
        return "recovering";
    case SWITCHD_STATE_ERROR:
        return "error";
    case SWITCHD_STATE_TERMINATING:
        return "terminating";
    default:
        return "unknown";
    }
}

const char *op_type_to_str(int type) {
    switch (type) {
    case SWITCHD_OP_NONE:
        return "none";
    case SWITCHD_OP_PORT_ADMIN_SET:
        return "setPortAdmin";
    case SWITCHD_OP_PORT_POE_SET:
        return "setPortPoe";
    case SWITCHD_OP_PORT_POE_CYCLE:
        return "cyclePortPoe";
    case SWITCHD_OP_SWITCH_REBOOT:
        return "rebootSwitch";
    case SWITCHD_OP_FW_STAGE:
        return "firmwareStage";
    case SWITCHD_OP_FW_APPLY:
        return "firmwareApply";
    default:
        return "unknown";
    }
}

const char *op_state_to_str(int state) {
    switch (state) {
    case SWITCHD_OP_STATE_IDLE:
        return "idle";
    case SWITCHD_OP_STATE_RUNNING:
        return "running";
    case SWITCHD_OP_STATE_DONE:
        return "done";
    case SWITCHD_OP_STATE_FAILED:
        return "failed";
    default:
        return "unknown";
    }
}

const char *fw_state_to_str(int state) {
    switch (state) {
    case SWITCHD_FW_IDLE:
        return "idle";
    case SWITCHD_FW_STAGED:
        return "staged";
    case SWITCHD_FW_APPLYING:
        return "applying";
    case SWITCHD_FW_REBOOTING:
        return "rebooting";
    case SWITCHD_FW_RECONNECTING:
        return "reconnecting";
    case SWITCHD_FW_VERIFYING:
        return "verifying";
    case SWITCHD_FW_DONE:
        return "done";
    case SWITCHD_FW_FAILED:
        return "failed";
    default:
        return "unknown";
    }
}

const char *alarm_severity_to_str(int severity) {
    switch (severity) {
    case SWITCHD_ALARM_SEV_INFO:
        return "info";
    case SWITCHD_ALARM_SEV_WARNING:
        return "warning";
    case SWITCHD_ALARM_SEV_CRITICAL:
        return "critical";
    default:
        return "unknown";
    }
}

const char *switch_error_to_str(int code) {
    switch (code) {
    case SWITCHD_OK:
        return "ok";
    case SWITCHD_ERR_NOMEM:
        return "out of memory";
    case SWITCHD_ERR_INVAL:
        return "invalid argument";
    case SWITCHD_ERR_IO:
        return "I/O error";
    case SWITCHD_ERR_TIMEOUT:
        return "timeout";
    case SWITCHD_ERR_BUSY:
        return "busy";
    case SWITCHD_ERR_NOTFOUND:
        return "not found";
    case SWITCHD_ERR_UNSUPPORTED:
        return "unsupported";
    case SWITCHD_ERR_SNMP:
        return "SNMP error";
    case SWITCHD_ERR_PROTOCOL:
        return "protocol error";
    case SWITCHD_ERR_STATE:
        return "invalid state";
    case SWITCHD_ERR_AUTH:
        return "authentication error";
    case SWITCHD_ERR_INTERNAL:
        return "internal error";
    default:
        return "unknown error";
    }
}
