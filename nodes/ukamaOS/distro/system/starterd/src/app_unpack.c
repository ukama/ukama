/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>
#include <errno.h>

#include "usys_log.h"

static bool unpack_validate_path(const char *p) {

    if (!p || !*p) return false;
    if (p[0] == '/') return false;
    if (strstr(p, "..") != NULL) return false;
    if (strstr(p, "\\") != NULL) return false;
    return true;
}

static bool unpack_list_and_validate(const char *tarPath) {

    int pipefd[2];
    pid_t pid;
    char buf[4096];
    ssize_t n;
    char line[1024];
    size_t li;
    int status;
    bool ok;

    ok = true;
    li = 0;

    if (pipe(pipefd) != 0) return false;

    pid = fork();
    if (pid < 0) {
        close(pipefd[0]);
        close(pipefd[1]);
        return false;
    }

    if (pid == 0) {
        dup2(pipefd[1], STDOUT_FILENO);
        dup2(pipefd[1], STDERR_FILENO);
        close(pipefd[0]);
        close(pipefd[1]);
        execlp("tar", "tar", "-tf", tarPath, (char *)NULL);
        _exit(127);
    }

    close(pipefd[1]);

    while ((n = read(pipefd[0], buf, sizeof(buf))) > 0) {

        for (ssize_t i = 0; i < n; i++) {
            if (buf[i] == '\n') {
                line[li] = '\0';
                if (li > 0) {
                    if (!unpack_validate_path(line)) {
                        ok = false;
                        break;
                    }
                }
                li = 0;
            } else {
                if (li + 1 < sizeof(line)) {
                    line[li++] = buf[i];
                }
            }
        }

        if (!ok) break;
    }

    close(pipefd[0]);

    waitpid(pid, &status, 0);
    if (!WIFEXITED(status) || WEXITSTATUS(status) != 0) ok = false;

    return ok;
}

bool app_unpack_package(const char *tarPath, const char *dstDir) {

    pid_t pid;
    int status;

    if (!tarPath || !dstDir) return false;

    if (!unpack_list_and_validate(tarPath)) {
        usys_log_error("unpack: validation failed for %s", tarPath);
        return false;
    }

    pid = fork();
    if (pid < 0) return false;

    if (pid == 0) {
        execlp("tar", "tar", "-xzf", tarPath, "-C", dstDir, (char *)NULL);
        _exit(127);
    }

    waitpid(pid, &status, 0);
    if (!WIFEXITED(status) || WEXITSTATUS(status) != 0) {
        usys_log_error("unpack: tar extract failed %s", tarPath);
        return false;
    }

    return true;
}
