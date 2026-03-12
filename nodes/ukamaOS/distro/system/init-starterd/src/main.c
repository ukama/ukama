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
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>
#include <sys/stat.h>

#include "usys_log.h"

#define INIT_SERVICE_NAME        "init-starter"
#define STARTER_SERVICE_NAME     "starter.d"
#define INIT_DEFAULT_ROOT        "/ukama/init/starter"
#define INIT_DEFAULT_READY_FILE  "/ukama/init/starter/ready"
#define INIT_DEFAULT_READY_TO    20
#define INIT_DEFAULT_TERM_GRACE  5
#define INIT_EXIT_SWITCH         77 /* starter exist with this code 
                                     * to trigger switch by init-starter */

static volatile sig_atomic_t gTerminate = 0;

static void on_signal(int sig) {

    (void)sig;
    gTerminate = 1;
}

static void setup_signals(void) {

    struct sigaction sa;

    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = on_signal;
    sigaction(SIGTERM, &sa, NULL);
    sigaction(SIGINT, &sa, NULL);
}

static char *get_str(const char *name, const char *defVal) {

    const char *v;

    v = getenv(name);
    if (v && *v) {
        return strdup(v);
    }

    return defVal ? strdup(defVal) : NULL;
}

static int get_int(const char *name, int defVal) {

    const char *v;
    long x;
    char *end;

    v = getenv(name);
    if (!v || !*v) {
        return defVal;
    }

    errno = 0;
    x = strtol(v, &end, 10);
    if (errno != 0 || end == v || *end != '\0') {
        return defVal;
    }

    if (x < 0) {
        return defVal;
    }

    if (x > 3600) {
        return 3600;
    }

    return (int)x;
}

static bool path_exists(const char *p) {

    struct stat st;

    if (!p) {
        return false;
    }

    return stat(p, &st) == 0;
}

static bool symlink_atomic(const char *linkPath, const char *target) {

    char tmp[512];
    int rc;

    snprintf(tmp, sizeof(tmp), "%s.tmp.%d", linkPath, getpid());
    unlink(tmp);

    rc = symlink(target, tmp);
    if (rc != 0) {
        return false;
    }

    rc = rename(tmp, linkPath);
    if (rc != 0) {
        unlink(tmp);
        return false;
    }

    return true;
}

static char *readlink_dup(const char *p) {

    char buf[512];
    ssize_t n;

    n = readlink(p, buf, sizeof(buf) - 1);
    if (n <= 0) {
        return NULL;
    }

    buf[n] = '\0';
    return strdup(buf);
}

static bool is_allowed_slot_target(const char *t) {

    if (!t || !*t) {
        return false;
    }

    if (strcmp(t, "slots/A") == 0) return true;
    if (strcmp(t, "slots/B") == 0) return true;

    return false;
}

static bool flip_slot_with_prev(const char *root,
                                char **oldTargetOut,
                                char **newTargetOut) {

    char curPath[512];
    char prevPath[512];
    char nextPath[512];
    char *oldTarget;
    char *newTarget;
    bool ok;

    ok = false;
    oldTarget = NULL;
    newTarget = NULL;

    snprintf(curPath,  sizeof(curPath),  "%s/current", root);
    snprintf(prevPath, sizeof(prevPath), "%s/prev",    root);
    snprintf(nextPath, sizeof(nextPath), "%s/next",    root);

    newTarget = readlink_dup(nextPath);
    if (!newTarget) {
        goto out;
    }

    if (!is_allowed_slot_target(newTarget)) {
        usys_log_error("invalid next target '%s'", newTarget);
        goto out;
    }

    oldTarget = readlink_dup(curPath);
    if (!oldTarget) {
        usys_log_error("current is not a symlink or unreadable");
        goto out;
    }

    if (!symlink_atomic(prevPath, oldTarget)) {
        usys_log_error("failed to set prev -> %s", oldTarget);
        goto out;
    }

    if (!symlink_atomic(curPath, newTarget)) {
        usys_log_error("failed to set current -> %s", newTarget);
        symlink_atomic(curPath, oldTarget);
        goto out;
    }

    unlink(nextPath);

    if (oldTargetOut) *oldTargetOut = oldTarget;
    if (newTargetOut) *newTargetOut = newTarget;

    oldTarget = NULL;
    newTarget = NULL;
    ok = true;

out:
    if (oldTarget) free(oldTarget);
    if (newTarget) free(newTarget);

    return ok;
}

static bool rollback_to_prev(const char *root) {

    char curPath[512];
    char prevPath[512];
    char *prevTarget;
    bool ok;

    ok         = false;
    prevTarget = NULL;

    snprintf(curPath,  sizeof(curPath),  "%s/current", root);
    snprintf(prevPath, sizeof(prevPath), "%s/prev",    root);

    prevTarget = readlink_dup(prevPath);
    if (!prevTarget) {
        usys_log_error("rollback failed, prev not set");
        goto out;
    }

    if (!is_allowed_slot_target(prevTarget)) {
        usys_log_error("rollback failed, invalid prev target '%s'", prevTarget);
        goto out;
    }

    if (!symlink_atomic(curPath, prevTarget)) {
        usys_log_error("rollback failed, cannot set current -> %s", prevTarget);
        goto out;
    }

    ok = true;

out:
    if (prevTarget) free(prevTarget);

    return ok;
}

static void kill_tree(pid_t pid, int graceSec) {

    int i;
    int rc;

    if (pid <= 0) return;

    rc = killpg(pid, SIGTERM);
    if (rc != 0) {
        kill(pid, SIGTERM);
    }

    for (i = 0; i < graceSec * 10; i++) {
        rc = kill(pid, 0);
        if (rc != 0) {
            return;
        }

        usleep(100000);
    }

    rc = killpg(pid, SIGKILL);
    if (rc != 0) {
        kill(pid, SIGKILL);
    }
}

static bool wait_ready_or_exit(pid_t pid, const char *readyFile, int timeoutSec) {

    int i;
    int status;
    pid_t w;

    if (timeoutSec <= 0) {
        timeoutSec = 1;
    }

    for (i = 0; i < timeoutSec * 10; i++) {

        if (gTerminate) {
            return false;
        }

        if (readyFile && path_exists(readyFile)) {
            return true;
        }

        status = 0;
        w = waitpid(pid, &status, WNOHANG);
        if (w == pid) {
            return false;
        }

        usleep(100000);
    }

    return false;
}

static int wait_child(pid_t pid) {

    int status;
    pid_t w;

    status = 0;

    while (1) {
        w = waitpid(pid, &status, 0);
        if (w == pid) break;
        if (w < 0 && errno == EINTR) {
            if (gTerminate) {
                return 128 + SIGTERM;
            }
            continue;
        }

        return 127;
    }

    if (WIFEXITED(status))   return WEXITSTATUS(status);
    if (WIFSIGNALED(status)) return 128 + WTERMSIG(status);

    return 127;
}

static int run_once(const char *root,
                    const char *readyFile,
                    int readyTimeoutSec,
                    int termGraceSec) {

    char cur[512];
    char bin[1024];
    pid_t pid;
    bool readyOk;
    int code;

    snprintf(cur, sizeof(cur), "%s/current", root);
    snprintf(bin, sizeof(bin), "%s/%s", cur, STARTER_SERVICE_NAME);

    if (!path_exists(bin)) {
        usys_log_error("bootstrap: missing %s", bin);
        return 127;
    }

    if (readyFile && *readyFile) {
        unlink(readyFile);
    }

    pid = fork();
    if (pid < 0) {
        usys_log_error("bootstrap: fork failed: %s", strerror(errno));
        return 127;
    }

    if (pid == 0) {
        setpgid(0, 0);
        execl(bin, STARTER_SERVICE_NAME, (char *)NULL);
        _exit(127);
    }

    readyOk = wait_ready_or_exit(pid, readyFile, readyTimeoutSec);
    if (!readyOk) {
        usys_log_error("bootstrap: starterd not ready, killing");
        kill_tree(pid, termGraceSec);
        code = wait_child(pid);
        (void)code;
        return 125;
    }

    return wait_child(pid);
}

int main(int argc, char **argv) {

    char *root;
    char *readyFile;
    int readyTimeoutSec;
    int termGraceSec;
    int code;
    int backoff;
    char *oldTarget;
    char *newTarget;

    (void)argc;
    (void)argv;

    setup_signals();

    usys_log_set_service(INIT_SERVICE_NAME);
    usys_log_set_level(USYS_LOG_INFO);

    root            = get_str("STARTER_INIT_ROOT",              INIT_DEFAULT_ROOT);
    readyFile       = get_str("STARTER_INIT_READY_FILE",        INIT_DEFAULT_READY_FILE);
    readyTimeoutSec = get_int("STARTER_INIT_READY_TIMEOUT_SEC", INIT_DEFAULT_READY_TO);
    termGraceSec    = get_int("STARTER_INIT_TERM_GRACE_SEC",    INIT_DEFAULT_TERM_GRACE);

    if (!root)      return 1;
    if (!readyFile) return 1;

    backoff = 1;

    while (!gTerminate) {

        code = run_once(root, readyFile, readyTimeoutSec, termGraceSec);

        if (code == 0) {
            backoff = 1;
            continue;
        }

        if (code == INIT_EXIT_SWITCH) {
            usys_log_info("bootstrap: switch requested");
            oldTarget = NULL;
            newTarget = NULL;

            if (!flip_slot_with_prev(root, &oldTarget, &newTarget)) {
                usys_log_error("bootstrap: switch failed");
                backoff = 1;
            } else {
                usys_log_info("bootstrap: switched current %s -> %s",
                              oldTarget ? oldTarget : "(null)",
                              newTarget ? newTarget : "(null)");

                free(oldTarget);
                free(newTarget);

                code = run_once(root, readyFile, readyTimeoutSec, termGraceSec);
                if (code != 0) {
                    usys_log_error("bootstrap: new starterd failed (%d), rolling back", code);
                    rollback_to_prev(root);
                }
                backoff = 1;
            }

            continue;
        }

        usys_log_error("bootstrap: starterd exit %d", code);
        sleep(backoff);
        backoff *= 2;
        if (backoff > 30) backoff = 30;
    }

    free(root);
    free(readyFile);

    return 0;
}
