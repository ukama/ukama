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
#include <time.h>

#include "usys_log.h"

extern char **environ;

static bool runtime_set_workdir(const char *workdir) {

    if (!workdir || !*workdir) return true;
    if (chdir(workdir) != 0) return false;
    return true;
}

static int runtime_env_count(char **envp) {

    int n = 0;

    if (!envp) {
        return 0;
    }

    while (envp[n]) {
        n++;
    }

    return n;
}

static size_t runtime_env_key_len(const char *kv) {

    const char *eq;

    if (!kv) {
        return 0;
    }

    eq = strchr(kv, '=');
    if (!eq) {
        return strlen(kv);
    }

    return (size_t)(eq - kv);
}

static bool runtime_env_same_key(const char *a, const char *b) {

    size_t alen;
    size_t blen;

    if (!a || !b) {
        return false;
    }

    alen = runtime_env_key_len(a);
    blen = runtime_env_key_len(b);

    if (alen != blen) {
        return false;
    }

    return strncmp(a, b, alen) == 0;
}

static void runtime_free_env_array(char **envp) {

    int i;

    if (!envp) {
        return;
    }

    for (i = 0; envp[i]; i++) {
        free(envp[i]);
    }

    free(envp);
}

static char **runtime_dup_env_array(char **src) {

    int i;
    int n;
    char **dst;

    n = runtime_env_count(src);

    dst = (char **)calloc(n + 1, sizeof(char *));
    if (!dst) {
        return NULL;
    }

    for (i = 0; i < n; i++) {
        dst[i] = strdup(src[i]);
        if (!dst[i]) {
            runtime_free_env_array(dst);
            return NULL;
        }
    }

    dst[n] = NULL;
    return dst;
}

static char **runtime_env_add(char **envp, const char *kv) {

    int n;
    char **newEnvp;

    if (!kv) {
        return envp;
    }

    n = runtime_env_count(envp);

    newEnvp = (char **)realloc(envp, sizeof(char *) * (n + 2));
    if (!newEnvp) {
        runtime_free_env_array(envp);
        return NULL;
    }

    newEnvp[n] = strdup(kv);
    if (!newEnvp[n]) {
        runtime_free_env_array(newEnvp);
        return NULL;
    }

    newEnvp[n + 1] = NULL;
    return newEnvp;
}

static bool runtime_env_set(char ***envpRef, const char *kv) {

    int i;
    char **envp;

    if (!envpRef || !kv) {
        return false;
    }

    envp = *envpRef;
    for (i = 0; envp && envp[i]; i++) {
        if (runtime_env_same_key(envp[i], kv)) {
            char *dup = strdup(kv);
            if (!dup) {
                return false;
            }

            free(envp[i]);
            envp[i] = dup;
            return true;
        }
    }

    envp = runtime_env_add(envp, kv);
    if (!envp) {
        return false;
    }

    *envpRef = envp;
    return true;
}

static char **runtime_build_child_envp(char **appEnvp) {

    int i;
    char **merged;

    merged = runtime_dup_env_array(environ);
    if (!merged) {
        return NULL;
    }

    for (i = 0; appEnvp && appEnvp[i]; i++) {
        if (!runtime_env_set(&merged, appEnvp[i])) {
            runtime_free_env_array(merged);
            return NULL;
        }
    }

    return merged;
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
        char **childEnvp = NULL;

        if (setpgid(0, 0) != 0) {
            _exit(127);
        }

        if (!runtime_set_workdir(app->workdir)) {
            _exit(127);
        }

        if (envp && envp[0]) {
            childEnvp = runtime_build_child_envp(envp);
            if (!childEnvp) {
                _exit(127);
            }

            execvpe(execPath, argv, childEnvp);
            runtime_free_env_array(childEnvp);
        } else {
            execvp(execPath, argv);
        }

        _exit(127);
    }

    app->pid = pid;
    app->pgid = pid;
    return true;
}

static int runtime_wait_pid(pid_t pid, int timeoutSec, int *statusOut) {

    time_t start;
    int status;
    pid_t r;

    start = time(NULL);
    while (true) {
        r = waitpid(pid, &status, WNOHANG);
        if (r == pid) {
            if (statusOut) {
                *statusOut = status;
            }
            return 1;
        }

        if (r == -1 && errno == ECHILD) {
            return -1;
        }

        if ((int)(time(NULL) - start) >= timeoutSec) {
            break;
        }

        usleep(100 * 1000);
    }

    return 0;
}

bool app_runtime_stop(Config *config, App *app) {

    int status;
    int rc;

    if (!config || !app) return false;

    if (app->pid <= 0 || app->pgid <= 0) return true;

    killpg(app->pgid, SIGTERM);

    status = 0;
    rc = runtime_wait_pid(app->pid, config->termGraceSec, &status);
    if (rc == 1) {
        app_runtime_note_exit(app, status);
        return true;
    }

    if (rc == -1) {
        app->lastPid = app->pid;
        app->lastPgid = app->pgid;
        app->pid = 0;
        app->pgid = 0;
        return true;
    }

    killpg(app->pgid, SIGKILL);

    status = 0;
    rc = runtime_wait_pid(app->pid, 2, &status);
    if (rc == 1) {
        app_runtime_note_exit(app, status);
        return true;
    }

    if (rc == -1) {
        app->lastPid = app->pid;
        app->lastPgid = app->pgid;
        app->pid = 0;
        app->pgid = 0;
        return true;
    }

    return false;
}

void app_runtime_note_exit(App *app, int status) {

    if (!app) return;

    app->lastPid = app->pid;
    app->lastPgid = app->pgid;

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

    app->pid = 0;
    app->pgid = 0;
}
