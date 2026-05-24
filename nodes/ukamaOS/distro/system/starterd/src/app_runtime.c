/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "app_runtime.h"

#include <errno.h>
#include <fcntl.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>
#include <time.h>

#include "usys_log.h"

extern char **environ;

#define RUNTIME_LOADER      "lib64/ld-linux-x86-64.so.2"
#define RUNTIME_LIB_PATH    "lib:lib/x86_64-linux-gnu:usr/lib:" \
                            "usr/lib/x86_64-linux-gnu"

typedef struct RuntimeLaunch {
    char  *appRoot;
    char  *loaderPath;
    char  *libraryPath;
    char **argv;
    bool   contained;
} RuntimeLaunch;

static bool runtime_set_workdir(const char *workdir, const char *appRoot) {

    if (workdir && *workdir) {
        return chdir(workdir) == 0;
    }

    if (appRoot && *appRoot) {
        return chdir(appRoot) == 0;
    }

    return true;
}

static int runtime_env_count(char **envp) {

    int n;

    n = 0;
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
    char *dup;

    if (!envpRef || !kv) {
        return false;
    }

    envp = *envpRef;
    for (i = 0; envp && envp[i]; i++) {
        if (runtime_env_same_key(envp[i], kv)) {
            dup = strdup(kv);
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

static bool runtime_env_setf(char ***envpRef,
                             const char *key,
                             const char *value) {

    char *kv;
    bool ok;

    if (!envpRef || !key || !value) {
        return false;
    }

    kv = NULL;
    if (asprintf(&kv, "%s=%s", key, value) < 0) {
        return false;
    }

    ok = runtime_env_set(envpRef, kv);
    free(kv);

    return ok;
}

static char **runtime_build_child_envp(App *app,
                                       RuntimeLaunch *launch,
                                       char **appEnvp) {

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

    if (launch && launch->appRoot) {
        if (!runtime_env_setf(&merged, "UKAMA_APP_ROOT",
                              launch->appRoot)) {
            runtime_free_env_array(merged);
            return NULL;
        }
    }

    if (app && app->name) {
        if (!runtime_env_setf(&merged, "UKAMA_APP_NAME", app->name)) {
            runtime_free_env_array(merged);
            return NULL;
        }
    }

    if (app && app->space) {
        if (!runtime_env_setf(&merged, "UKAMA_APP_SPACE", app->space)) {
            runtime_free_env_array(merged);
            return NULL;
        }
    }

    if (app && app->tag) {
        if (!runtime_env_setf(&merged, "UKAMA_APP_TAG", app->tag)) {
            runtime_free_env_array(merged);
            return NULL;
        }
    }

    if (launch && launch->contained) {
        if (!runtime_env_set(&merged, "UKAMA_RUNTIME=contained")) {
            runtime_free_env_array(merged);
            return NULL;
        }
    }

    return merged;
}

static char *runtime_path_join(const char *a, const char *b) {

    char *p;

    if (!a || !b) {
        return NULL;
    }

    p = NULL;
    if (asprintf(&p, "%s/%s", a, b) < 0) {
        return NULL;
    }

    return p;
}

static char *runtime_app_root_from_exec(const App *app,
                                        const char *execPath) {

    size_t execLen;
    size_t cmdLen;
    size_t rootLen;
    char *root;

    if (!app || !app->cmd || !execPath) {
        return NULL;
    }

    if (app->cmd[0] == '/') {
        return NULL;
    }

    execLen = strlen(execPath);
    cmdLen  = strlen(app->cmd);

    if (execLen <= cmdLen) {
        return NULL;
    }

    if (strcmp(execPath + execLen - cmdLen, app->cmd) != 0) {
        return NULL;
    }

    rootLen = execLen - cmdLen;
    while (rootLen > 0 && execPath[rootLen - 1] == '/') {
        rootLen--;
    }

    if (rootLen == 0) {
        return NULL;
    }

    root = (char *)calloc(rootLen + 1, 1);
    if (!root) {
        return NULL;
    }

    memcpy(root, execPath, rootLen);
    root[rootLen] = '\0';

    return root;
}

static char *runtime_make_library_path(const char *appRoot) {

    char *path;

    if (!appRoot) {
        return NULL;
    }

    path = NULL;
    if (asprintf(&path, "%s/lib:%s/lib/x86_64-linux-gnu:"
                        "%s/usr/lib:%s/usr/lib/x86_64-linux-gnu",
                 appRoot, appRoot, appRoot, appRoot) < 0) {
        return NULL;
    }

    return path;
}

static int runtime_argv_count(char **argv) {

    int n;

    n = 0;
    while (argv && argv[n]) {
        n++;
    }

    return n;
}

static void runtime_free_argv_array(char **argv) {

    int i;

    if (!argv) {
        return;
    }

    for (i = 0; argv[i]; i++) {
        free(argv[i]);
    }

    free(argv);
}

static char **runtime_build_loader_argv(const char *loaderPath,
                                        const char *libraryPath,
                                        const char *execPath,
                                        char **appArgv) {

    int appArgc;
    int i;
    int j;
    int srcStart;
    char **argv;

    if (!loaderPath || !libraryPath || !execPath) {
        return NULL;
    }

    appArgc  = runtime_argv_count(appArgv);
    srcStart = appArgc > 0 ? 1 : 0;

    argv = (char **)calloc(4 + (appArgc - srcStart) + 1,
                           sizeof(char *));
    if (!argv) {
        return NULL;
    }

    i = 0;
    argv[i++] = strdup(loaderPath);
    argv[i++] = strdup("--library-path");
    argv[i++] = strdup(libraryPath);
    argv[i++] = strdup(execPath);

    for (j = srcStart; j < appArgc; j++) {
        argv[i++] = strdup(appArgv[j]);
    }

    argv[i] = NULL;

    for (j = 0; j < i; j++) {
        if (!argv[j]) {
            runtime_free_argv_array(argv);
            return NULL;
        }
    }

    return argv;
}

static void runtime_launch_init(RuntimeLaunch *launch) {

    if (!launch) {
        return;
    }

    memset(launch, 0, sizeof(RuntimeLaunch));
}

static void runtime_launch_free(RuntimeLaunch *launch) {

    if (!launch) {
        return;
    }

    free(launch->appRoot);
    free(launch->loaderPath);
    free(launch->libraryPath);

    if (launch->contained) {
        runtime_free_argv_array(launch->argv);
    }

    memset(launch, 0, sizeof(RuntimeLaunch));
}

static bool runtime_launch_build(App *app,
                                 const char *execPath,
                                 RuntimeLaunch *launch) {

    char *loaderPath;
    char *appRoot;
    char *libraryPath;
    char **loaderArgv;

    if (!app || !execPath || !launch) {
        return false;
    }

    runtime_launch_init(launch);

    appRoot = runtime_app_root_from_exec(app, execPath);
    if (!appRoot) {
        launch->argv      = app->argv;
        launch->contained = false;
        return true;
    }

    loaderPath = runtime_path_join(appRoot, RUNTIME_LOADER);
    if (!loaderPath) {
        free(appRoot);
        return false;
    }

    if (access(loaderPath, X_OK) != 0) {
        free(loaderPath);
        launch->appRoot   = appRoot;
        launch->argv      = app->argv;
        launch->contained = false;
        return true;
    }

    libraryPath = runtime_make_library_path(appRoot);
    if (!libraryPath) {
        free(appRoot);
        free(loaderPath);
        return false;
    }

    loaderArgv = runtime_build_loader_argv(loaderPath,
                                           libraryPath,
                                           execPath,
                                           app->argv);
    if (!loaderArgv) {
        free(appRoot);
        free(loaderPath);
        free(libraryPath);
        return false;
    }

    launch->appRoot     = appRoot;
    launch->loaderPath  = loaderPath;
    launch->libraryPath = libraryPath;
    launch->argv        = loaderArgv;
    launch->contained   = true;

    return true;
}

static const char *runtime_exec_path(const char *execPath,
                                     RuntimeLaunch *launch) {

    if (launch && launch->contained && launch->loaderPath) {
        return launch->loaderPath;
    }

    return execPath;
}

static void runtime_child_fail(int fd, int err) {

    if (fd >= 0) {
        (void)write(fd, &err, sizeof(err));
    }

    _exit(127);
}

static bool runtime_set_cloexec(int fd) {

    int flags;

    flags = fcntl(fd, F_GETFD);
    if (flags < 0) {
        return false;
    }

    if (fcntl(fd, F_SETFD, flags | FD_CLOEXEC) < 0) {
        return false;
    }

    return true;
}

bool app_runtime_start(Config *config, App *app, const char *execPath) {

    pid_t pid;
    char **envp;
    int errPipe[2];
    int childErr;
    ssize_t n;

    if (!config || !app || !execPath) {
        return false;
    }

    envp = app->envp;
    errPipe[0] = -1;
    errPipe[1] = -1;

    if (pipe(errPipe) != 0) {
        usys_log_error("runtime: pipe failed for %s/%s",
                       app->space, app->name);
        return false;
    }

    if (!runtime_set_cloexec(errPipe[1])) {
        close(errPipe[0]);
        close(errPipe[1]);
        usys_log_error("runtime: cloexec failed for %s/%s",
                       app->space, app->name);
        return false;
    }

    pid = fork();
    if (pid < 0) {
        close(errPipe[0]);
        close(errPipe[1]);
        usys_log_error("runtime: fork failed for %s/%s",
                       app->space, app->name);
        return false;
    }

    if (pid == 0) {
        RuntimeLaunch launch;
        char **childEnvp;
        const char *realExec;

        close(errPipe[0]);

        runtime_launch_init(&launch);

        if (!runtime_launch_build(app, execPath, &launch)) {
            runtime_child_fail(errPipe[1], errno ? errno : EINVAL);
        }

        if (setpgid(0, 0) != 0) {
            runtime_launch_free(&launch);
            runtime_child_fail(errPipe[1], errno);
        }

        if (!runtime_set_workdir(app->workdir, launch.appRoot)) {
            runtime_launch_free(&launch);
            runtime_child_fail(errPipe[1], errno);
        }

        childEnvp = runtime_build_child_envp(app, &launch, envp);
        if (!childEnvp) {
            runtime_launch_free(&launch);
            runtime_child_fail(errPipe[1], ENOMEM);
        }

        realExec = runtime_exec_path(execPath, &launch);

        execve(realExec, launch.argv, childEnvp);

        childErr = errno;
        runtime_free_env_array(childEnvp);
        runtime_launch_free(&launch);

        runtime_child_fail(errPipe[1], childErr);
    }

    close(errPipe[1]);

    childErr = 0;
    n = read(errPipe[0], &childErr, sizeof(childErr));
    close(errPipe[0]);

    if (n == sizeof(childErr)) {
        (void)waitpid(pid, NULL, 0);
        usys_log_error("runtime: exec failed for %s/%s: %s",
                       app->space, app->name, strerror(childErr));
        return false;
    }

    if (n < 0) {
        (void)waitpid(pid, NULL, 0);
        usys_log_error("runtime: exec status read failed for %s/%s",
                       app->space, app->name);
        return false;
    }

    app->pid  = pid;
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

    if (!config || !app) {
        return false;
    }

    if (app->pid <= 0 || app->pgid <= 0) {
        return true;
    }

    killpg(app->pgid, SIGTERM);

    status = 0;
    rc = runtime_wait_pid(app->pid, config->termGraceSec, &status);
    if (rc == 1) {
        app_runtime_note_exit(app, status);
        return true;
    }

    if (rc == -1) {
        app->lastPid  = app->pid;
        app->lastPgid = app->pgid;
        app->pid  = 0;
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
        app->lastPid  = app->pid;
        app->lastPgid = app->pgid;
        app->pid  = 0;
        app->pgid = 0;
        return true;
    }

    return false;
}

void app_runtime_note_exit(App *app, int status) {

    if (!app) {
        return;
    }

    app->lastPid  = app->pid;
    app->lastPgid = app->pgid;

    if (WIFEXITED(status)) {
        app->lastExitCode   = WEXITSTATUS(status);
        app->lastExitSignal = 0;
    } else if (WIFSIGNALED(status)) {
        app->lastExitCode   = 0;
        app->lastExitSignal = WTERMSIG(status);
    } else {
        app->lastExitCode   = 0;
        app->lastExitSignal = 0;
    }

    app->pid  = 0;
    app->pgid = 0;
}
