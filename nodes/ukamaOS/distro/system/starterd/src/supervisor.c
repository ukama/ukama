/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "supervisor.h"
#include "space.h"
#include "app.h"
#include "installer.h"
#include "state_store.h"
#include "web_client.h"
#include "restart_policy.h"
#include "app_runtime.h"
#include "starterd.h"

#include <errno.h>
#include <pthread.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <sys/stat.h>
#include <unistd.h>
#include <time.h>

#include "usys_log.h"

struct Supervisor {
    Config         *config;
    Space          *spaceList;
    ActionQueue    *queue;
    StarterContext *ctx;

    pthread_t       thread;
    pthread_mutex_t mu;
    pthread_cond_t  cv;

    bool running;
    bool bootDone;
};

static void action_free(Action *a) {

    if (!a) {
        return;
    }

    free(a->space);
    free(a->name);
    free(a->tag);
    free(a);
}

static void ready_touch(Config *config) {

    FILE *f;
    char *p;
    char *slash;

    if (!config || !config->readyFile) {
        return;
    }

    p = strdup(config->readyFile);
    if (!p) {
        return;
    }

    slash = strrchr(p, '/');
    if (slash) {
        *slash = '\0';
        mkdir(p, 0755);
    }
    free(p);

    f = fopen(config->readyFile, "w");
    if (f) {
        fputs("ready\n", f);
        fclose(f);
    }
}

static char* app_exec_path(Config *config, App *app) {

    char *p;

    if (!config || !app || !app->cmd) {
        return NULL;
    }

    if (app->cmd[0] == '/') {
        return strdup(app->cmd);
    }

    p = NULL;
    if (asprintf(&p, "%s/%s/%s/current/%s",
                 config->appsRoot,
                 app->space,
                 app->name,
                 app->cmd) < 0) {
        p = NULL;
    }

    return p;
}

static bool app_start(Config *config, App *app) {

    char *execPath;
    bool ok;
    time_t now;

    execPath = NULL;
    ok = false;

    if (!config || !app) {
        return false;
    }

    execPath = app_exec_path(config, app);
    if (!execPath) {
        return false;
    }

    app->state = APP_STATE_STARTING;

    now = time(NULL);
    restart_policy_on_start(config, app, now);

    ok = app_runtime_start(config, app, execPath);
    if (!ok) {
        app->state = APP_STATE_FAILED;
    } else {
        app->state = APP_STATE_RUNNING;
    }

    free(execPath);
    return ok;
}

static bool app_stop(Config *config, App *app) {

    if (!config || !app) {
        return false;
    }

    app->state = APP_STATE_STOPPING;
    app_runtime_stop(config, app);
    app->state = APP_STATE_STOPPED;

    return true;
}

static bool app_is_self(const App *app) {

    if (!app || !app->name) {
        return false;
    }

    return strcmp(app->name, STARTERD_SERVICE_NAME) == 0;
}

static void supervisor_stop_all_apps(Supervisor *s) {

    Space *sp;
    App *a;

    if (!s) {
        return;
    }

    sp = s->spaceList;
    while (sp) {
        a = sp->appList;
        while (a) {
            if (!app_is_self(a) &&
                (a->state == APP_STATE_RUNNING ||
                 a->state == APP_STATE_STARTING ||
                 a->state == APP_STATE_FAILED)) {
                usys_log_info("shutdown: stopping %s/%s", a->space, a->name);
                app_stop(s->config, a);
            }
            a = a->next;
        }
        sp = sp->next;
    }

    state_store_save(s->config, s->spaceList);
}

static void supervisor_reap(Supervisor *s) {

    int status;
    pid_t pid;
    Space *sp;
    App *a;
    time_t now;
    int delay;

    if (!s) {
        return;
    }

    while ((pid = waitpid(-1, &status, WNOHANG)) > 0) {

        sp = s->spaceList;
        while (sp) {
            a = sp->appList;
            while (a) {
                if (a->pid == pid) {
                    app_runtime_note_exit(a, status);
                    now = time(NULL);
                    restart_policy_on_exit(s->config, a, now);

                    if (a->state == APP_STATE_STOPPING ||
                        a->state == APP_STATE_STOPPED) {
                        a->state = APP_STATE_STOPPED;
                        break;
                    }

                    a->state = APP_STATE_FAILED;

                    delay = restart_policy_next_delay(s->config, a, now);
                    usys_log_error("app: exited %s/%s restart in %d sec",
                                   a->space, a->name, delay);
                    sleep(delay);
                    app_start(s->config, a);
                    state_store_save(s->config, s->spaceList);
                    break;
                }
                a = a->next;
            }
            sp = sp->next;
        }
    }
}

static bool app_wait_commit(Config *config, App *app) {

    time_t start;

    if (!config || !app) {
        return false;
    }

    start = time(NULL);
    while (true) {

        if (wc_app_ping(config, app) &&
            wc_app_version_matches(config, app, app->tag)) {
            return true;
        }

        if ((int)(time(NULL) - start) >= config->commitTimeoutSec) {
            break;
        }

        usleep(200 * 1000);
    }

    return false;
}

static bool run_space(Config *config,
                      Space *spaceList,
                      const char *spaceName,
                      bool gate) {

    Space *s;
    App *a;

    s = space_find(spaceList, spaceName);
    if (!s) {
        return true;
    }

    a = s->appList;
    while (a) {

        if (!installer_ensure_installed(config, a, NULL)) {
            return false;
        }

        if (!installer_switch_current(config, a)) {
            return false;
        }

        if (!app_start(config, a)) {
            return false;
        }

        if (gate) {
            if (!app_wait_commit(config, a)) {
                usys_log_error("boot: gate failed %s/%s", a->space, a->name);
                return false;
            }
        }

        a = a->next;
    }

    return true;
}

static bool update_self(Supervisor *s,
                        App *a,
                        const char *tag,
                        const char *hub) {

    char *oldTag;
    char *oldLastGood;

    if (!s || !s->config || !a || !tag) {
        return false;
    }

    oldTag = strdup(a->tag ? a->tag : "");
    oldLastGood = strdup(a->lastGoodTag ? a->lastGoodTag : "");
    if (!oldTag || !oldLastGood) {
        free(oldTag);
        free(oldLastGood);
        return false;
    }

    usys_log_info("self-update: %s/%s -> %s", a->space, a->name, tag);

    free(a->tag);
    a->tag = strdup(tag);
    if (!a->tag) {
        a->tag = oldTag;
        free(oldLastGood);
        return false;
    }

    if (!installer_ensure_installed(s->config, a, hub)) {
        usys_log_error("self-update: install failed %s/%s", a->space, a->name);
        free(a->tag);
        a->tag = oldTag;
        free(oldLastGood);
        return false;
    }

    if (!installer_switch_current(s->config, a)) {
        usys_log_error("self-update: switch failed %s/%s", a->space, a->name);
        free(a->tag);
        a->tag = oldTag;
        free(oldLastGood);
        return false;
    }

    free(a->lastGoodTag);
    a->lastGoodTag = strdup(tag);
    if (!a->lastGoodTag) {
        a->lastGoodTag = oldLastGood;
        free(oldTag);
        return false;
    }

    state_store_save(s->config, s->spaceList);

    if (s->ctx) {
        s->ctx->switchRequested = 1;
        s->ctx->exitCode = 77;
    }

    free(oldTag);
    free(oldLastGood);

    usys_log_info("self-update: staged successfully, switch requested");
    return true;
}

static bool update_app(Supervisor *s,
                       const char *space,
                       const char *name,
                       const char *tag,
                       const char *hub) {

    App *a;
    char *oldTag;
    char *oldLastGood;

    if (!s || !space || !name || !tag) {
        return false;
    }

    a = app_find(s->spaceList, space, name);
    if (!a) {
        return false;
    }

    if (app_is_self(a)) {
        return update_self(s, a, tag, hub);
    }

    oldTag = strdup(a->tag ? a->tag : "");
    oldLastGood = strdup(a->lastGoodTag ? a->lastGoodTag : "");
    if (!oldTag || !oldLastGood) {
        free(oldTag);
        free(oldLastGood);
        return false;
    }

    usys_log_info("update: %s/%s -> %s", space, name, tag);

    app_stop(s->config, a);

    free(a->tag);
    a->tag = strdup(tag);
    if (!a->tag) {
        a->tag = oldTag;
        free(oldLastGood);
        return false;
    }

    if (!installer_ensure_installed(s->config, a, hub)) {
        usys_log_error("update: install failed %s/%s", space, name);
        free(a->tag);
        a->tag = oldTag;
        app_start(s->config, a);
        free(oldLastGood);
        return false;
    }

    if (!installer_switch_current(s->config, a)) {
        usys_log_error("update: switch failed %s/%s", space, name);
        free(a->tag);
        a->tag = oldTag;
        app_start(s->config, a);
        free(oldLastGood);
        return false;
    }

    if (!app_start(s->config, a)) {
        usys_log_error("update: start failed %s/%s", space, name);
        installer_revert_to_last_good(s->config, a);
        free(a->tag);
        a->tag = oldTag;
        app_start(s->config, a);
        free(oldLastGood);
        return false;
    }

    if (!app_wait_commit(s->config, a)) {
        usys_log_error("update: commit failed %s/%s auto-revert", space, name);
        app_stop(s->config, a);
        free(a->tag);
        a->tag = oldTag;
        installer_switch_current(s->config, a);
        app_start(s->config, a);
        free(oldLastGood);
        return false;
    }

    free(a->lastGoodTag);
    a->lastGoodTag = strdup(tag);

    state_store_save(s->config, s->spaceList);

    free(oldTag);
    free(oldLastGood);
    return true;
}

static void* supervisor_thread(void *arg) {

    Supervisor *s;
    Action *a;

    s = (Supervisor *)arg;

    pthread_mutex_lock(&s->mu);

    while (s->running) {

        supervisor_reap(s);

        a = actions_dequeue(s->queue);
        if (!a) {
            struct timespec ts;
            clock_gettime(CLOCK_REALTIME, &ts);
            ts.tv_sec += 1;
            pthread_cond_timedwait(&s->cv, &s->mu, &ts);
            continue;
        }

        pthread_mutex_unlock(&s->mu);

        if (a->type == ACTION_RUN_BOOT) {
            if (!run_space(s->config, s->spaceList, s->config->bootSpace, true)) {
                usys_log_error("boot: failed");
            } else {
                s->bootDone = true;
                ready_touch(s->config);
                usys_log_info("boot: ready");
            }
        } else if (a->type == ACTION_RUN_ALL) {
            Space *sp;

            sp = s->spaceList;
            while (sp) {
                if (strcmp(sp->name, s->config->bootSpace) != 0) {
                    run_space(s->config, s->spaceList, sp->name, false);
                }
                sp = sp->next;
            }
        } else if (a->type == ACTION_TERMINATE_APP) {
            App *app;

            app = app_find(s->spaceList, a->space, a->name);
            if (app) {
                usys_log_info("terminate: %s/%s", a->space, a->name);
                app_stop(s->config, app);
                state_store_save(s->config, s->spaceList);
            }
        } else if (a->type == ACTION_UPDATE_APP) {
            if (!update_app(s, a->space, a->name, a->tag, a->hub)) {
                usys_log_error("update: failed %s/%s -> %s",
                               a->space ? a->space : "(null)",
                               a->name ? a->name :   "(null)",
                               a->tag ? a->tag :     "(null)",
                               a->hub ? a->hub :     "(null)");
            }
            if (s->ctx) {
                s->ctx->updateInProgress = 0;
            }
        }

        action_free(a);

        pthread_mutex_lock(&s->mu);
    }

    pthread_mutex_unlock(&s->mu);
    return NULL;
}

Supervisor* supervisor_start(Config *config,
                             Space *spaceList,
                             ActionQueue *queue,
                             StarterContext *ctx) {

    Supervisor *s;

    if (!config || !spaceList || !queue || !ctx) {
        return NULL;
    }

    s = calloc(1, sizeof(*s));
    if (!s) {
        return NULL;
    }

    s->config = config;
    s->spaceList = spaceList;
    s->queue = queue;
    s->ctx = ctx;

    pthread_mutex_init(&s->mu, NULL);
    pthread_cond_init(&s->cv, NULL);

    s->running = true;
    s->bootDone = false;

    if (pthread_create(&s->thread, NULL, supervisor_thread, s) != 0) {
        pthread_mutex_destroy(&s->mu);
        pthread_cond_destroy(&s->cv);
        free(s);
        return NULL;
    }

    return s;
}

void supervisor_stop(Supervisor *s) {

    if (!s) {
        return;
    }

    pthread_mutex_lock(&s->mu);
    s->running = false;
    pthread_cond_broadcast(&s->cv);
    pthread_mutex_unlock(&s->mu);

    pthread_join(s->thread, NULL);

    supervisor_stop_all_apps(s);

    pthread_mutex_destroy(&s->mu);
    pthread_cond_destroy(&s->cv);

    free(s);
}

bool supervisor_signal(Supervisor *s) {

    if (!s) {
        return false;
    }

    pthread_mutex_lock(&s->mu);
    pthread_cond_broadcast(&s->cv);
    pthread_mutex_unlock(&s->mu);

    return true;
}
