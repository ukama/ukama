/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <signal.h>
#include <stdbool.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#include "usys_log.h"

#include "exec.h"

bool exec_tool_exists(const char *cmd) {

    const char *path;
    char *copy;
    char *dir;
    char full[EXEC_MAX_PATH];
    bool found;

    if (cmd == NULL || cmd[0] == '\0') return false;

    if (strchr(cmd, '/') != NULL) return access(cmd, X_OK) == 0;

    path = getenv("PATH");
    if (path == NULL || path[0] == '\0') {
        path = "/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin";
    }

    copy = strdup(path);
    if (copy == NULL) return false;

    found = false;
    dir = strtok(copy, ":");

    while (dir != NULL) {
        snprintf(full, sizeof(full), "%s/%s", dir, cmd);
        if (access(full, X_OK) == 0) {
            found = true;
            break;
        }
        dir = strtok(NULL, ":");
    }

    free(copy);
    return found;
}

int exec_cmd_argv(int timeoutSec, char *const argv[]) {

    pid_t pid;
    int status;
    int waited;
    int rc;

    if (argv == NULL || argv[0] == NULL) return -1;

    pid = fork();
    if (pid < 0) {
        usys_log_error("fork failed for %s: %s", argv[0], strerror(errno));
        return -1;
    }

    if (pid == 0) {
        execvp(argv[0], argv);
        _exit(127);
    }

    waited = 0;
    status = 0;

    while (true) {
        rc = waitpid(pid, &status, WNOHANG);
        if (rc == pid) {
            if (WIFEXITED(status)) return WEXITSTATUS(status);
            if (WIFSIGNALED(status)) return 128 + WTERMSIG(status);
            return -1;
        }

        if (rc < 0) {
            if (errno == EINTR) continue;
            usys_log_error("waitpid failed for %s: %s", argv[0], strerror(errno));
            return -1;
        }

        if (waited >= timeoutSec) {
            usys_log_error("command timed out: %s", argv[0]);
            kill(pid, SIGTERM);
            sleep(1);
            kill(pid, SIGKILL);
            waitpid(pid, &status, 0);
            return -1;
        }

        sleep(1);
        waited++;
    }
}

int exec_cmd(int timeoutSec, const char *cmd, ...) {

    va_list args;
    const char *arg;
    char *argv[EXEC_MAX_ARGS];
    int argc;
    int rc;

    if (cmd == NULL) return -1;

    memset(argv, 0, sizeof(argv));

    argc = 0;
    argv[argc++] = (char *)cmd;

    va_start(args, cmd);
    while (argc < EXEC_MAX_ARGS - 1) {
        arg = va_arg(args, const char *);
        if (arg == NULL) break;
        argv[argc++] = (char *)arg;
    }
    va_end(args);

    argv[argc] = NULL;

    rc = exec_cmd_argv(timeoutSec, argv);
    if (rc != 0) {
        usys_log_error("command failed rc=%d cmd=%s", rc, cmd);
    }

    return rc;
}
