/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "app_runtime.h"

#include <errno.h>
#include <signal.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>

#include "usys_log.h"

static bool runtime_set_workdir(const char *workdir) {

    if (!workdir || !*workdir) return true;
    if (chdir(workdir) != 0) return false;
    return true;
}

bool app_runtime_start(Config *config, App *app, const char *execPath) {

    pid_t pid;
    char **argv;
    char **envp;

    if (!config || !app || !execPath) return false;

    argv = app->argv;
    envp = app->envp;

    pid = fork();
    if (pid < 0) {
        usys_log_error("runtime: fork failed for %s/%s", app->space, app->name);
        return false;
    }

    if (pid == 0) {

        if (setpgid(0, 0) != 0) {
            _exit(127);
        }

        if (!runtime_set_workdir(app->workdir)) {
            _exit(127);
        }

        if (envp && envp[0]) {
            execvpe(execPath, argv, envp);
        } else {
            execvp(execPath, argv);
        }

        _exit(127);
    }

    app->pid = pid;
    app->pgid = pid;
    return true;
}

static bool runtime_wait_pid(pid_t pid, int timeoutSec) {

    time_t start;
    int status;
    pid_t r;

    start = time(NULL);
    while (true) {
        r = waitpid(pid, &status, WNOHANG);
        if (r == pid) return true;
        if (r == -1 && errno == ECHILD) return true;

        if ((int)(time(NULL) - start) >= timeoutSec) break;
        usleep(100 * 1000);
    }

    return false;
}

bool app_runtime_stop(Config *config, App *app) {

    if (!config || !app) return false;

    if (app->pid <= 0 || app->pgid <= 0) return true;

    killpg(app->pgid, SIGTERM);

    if (!runtime_wait_pid(app->pid, config->termGraceSec)) {
        killpg(app->pgid, SIGKILL);
        runtime_wait_pid(app->pid, 2);
    }

    return true;
}

void app_runtime_note_exit(App *app, int status) {

    if (!app) return;

    if (WIFEXITED(status)) {
        app->lastExitCode = WEXITSTATUS(status);
        app->lastExitSignal = 0;
    } else if (WIFSIGNALED(status)) {
        app->lastExitCode = 0;
        app->lastExitSignal = WTERMSIG(status);
    } else {
        app->lastExitCode = 0;
        app->lastExitSignal = 0;
    }

    app->pid = -1;
    app->pgid = -1;
}
