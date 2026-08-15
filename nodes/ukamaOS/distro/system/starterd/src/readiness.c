/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "app.h"
#include "readiness.h"
#include "starterd.h"
#include "web_client.h"

#include "usys_log.h"
#include "usys_services.h"

struct ReadinessMonitor {
    Config *config;
    Space *spaceList;
    StarterContext *ctx;

    pthread_t thread;
    pthread_mutex_t mutex;
    pthread_cond_t condition;
    bool running;
    bool enabled;
    bool bootEnabled;
    bool bootStarted;

    NodeReadinessState state;
    time_t readinessDeadline;
    char reason[STARTERD_READY_REASON_LEN];

    NodeReadinessState bootState;
    time_t bootDeadline;
    char bootReason[STARTERD_READY_REASON_LEN];

    bool meshKnown;
    bool meshConnected;
    time_t meshCheckedAt;
    char meshReason[STARTERD_READY_REASON_LEN];
};

static const char *app_readiness_str(AppReadinessState state) {

    switch (state) {
    case APP_READINESS_IGNORED: return "ignored";
    case APP_READINESS_PENDING: return "pending";
    case APP_READINESS_READY:   return "ready";
    case APP_READINESS_FAULTY:  return "faulty";
    default:                    return "unknown";
    }
}

static const char *node_readiness_str(NodeReadinessState state) {

    switch (state) {
    case NODE_READINESS_PENDING: return "pending";
    case NODE_READINESS_READY:   return "ready";
    case NODE_READINESS_FAULTY:  return "faulty";
    default:                     return "unknown";
    }
}

static bool app_is_mesh(const App *app) {

    if (!app || !app->service) return false;

    return strcmp(app->service, SERVICE_MESH) == 0;
}

static void copy_reason(char *dst, size_t size, const char *reason) {

    if (!dst || size == 0) return;

    snprintf(dst, size, "%s", (reason && *reason) ? reason : "none");
}

static void app_set_pending(App *app,
                            int status,
                            const char *reason,
                            time_t now) {

    if (app->readinessState != APP_READINESS_PENDING) {
        app->readinessSince = now;
    }

    app->readinessState = APP_READINESS_PENDING;
    app->readinessHttpStatus = status;
    app->readinessCheckedAt = now;
    copy_reason(app->readinessReason,
                sizeof(app->readinessReason),
                reason);
}

static void app_update(ReadinessMonitor *monitor,
                       App *app,
                       const AppReadyResponse *result,
                       bool responded,
                       uint32_t generation,
                       time_t now) {

    pthread_mutex_lock(&monitor->mutex);

    if (app->readinessGeneration != generation) {
        pthread_mutex_unlock(&monitor->mutex);
        return;
    }

    if (!responded) {
        app_set_pending(app,
                        0,
                        "ready endpoint unavailable",
                        now);
    } else if (result->status == 200 && result->ready) {
        if (app->readinessState != APP_READINESS_READY) {
            app->readinessSince = now;
        }
        app->readinessState = APP_READINESS_READY;
        app->readinessHttpStatus = result->status;
        app->readinessCheckedAt = now;
        copy_reason(app->readinessReason,
                    sizeof(app->readinessReason),
                    result->reason);
    } else if (result->status == 503) {
        if (app->readinessState != APP_READINESS_FAULTY) {
            app->readinessSince = now;
        }
        app->readinessState = APP_READINESS_FAULTY;
        app->readinessHttpStatus = result->status;
        app->readinessCheckedAt = now;
        copy_reason(app->readinessReason,
                    sizeof(app->readinessReason),
                    result->reason);
    } else {
        app_set_pending(app,
                        result->status,
                        result->reason,
                        now);
    }

    snprintf(app->readinessRequestId,
             sizeof(app->readinessRequestId),
             "%s",
             responded ? result->requestId : "");

    pthread_mutex_unlock(&monitor->mutex);
}

static void aggregate_boot(ReadinessMonitor *monitor) {

    Space *space;
    App *app;
    App *pendingApp;
    NodeReadinessState state;
    const char *reason;
    char appReason[STARTERD_READY_REASON_LEN];
    time_t now;

    state = NODE_READINESS_READY;
    reason = "ready";
    appReason[0] = '\0';
    pendingApp = NULL;
    now = time(NULL);

    pthread_mutex_lock(&monitor->mutex);

    if (!monitor->bootStarted) {
        state = NODE_READINESS_PENDING;
        reason = "starterd is starting boot applications";
        monitor->bootDeadline = 0;
    } else if (monitor->bootEnabled) {
        space = space_find(monitor->spaceList,
                           monitor->config->bootSpace);
        app = space ? space->appList : NULL;

        while (app && state != NODE_READINESS_FAULTY) {
            if (!app->readinessRequired) {
                app = app->next;
                continue;
            }

            if (app->readinessState == APP_READINESS_FAULTY) {
                state = NODE_READINESS_FAULTY;
                snprintf(appReason,
                         sizeof(appReason),
                         "%.48s: %.140s",
                         app->name,
                         app->readinessReason);
                reason = appReason;
                break;
            }

            if (app->readinessState != APP_READINESS_READY &&
                state == NODE_READINESS_READY) {
                state = NODE_READINESS_PENDING;
                pendingApp = app;
                snprintf(appReason,
                         sizeof(appReason),
                         "%.48s: %.140s",
                         app->name,
                         app->readinessReason);
                reason = appReason;
            }

            app = app->next;
        }

        if (state == NODE_READINESS_READY ||
            state == NODE_READINESS_FAULTY) {
            monitor->bootDeadline = 0;
        } else {
            if (monitor->bootDeadline == 0) {
                monitor->bootDeadline =
                    now + monitor->config->readyTimeoutSec;
            }

            if (now >= monitor->bootDeadline) {
                state = NODE_READINESS_FAULTY;
                snprintf(appReason,
                         sizeof(appReason),
                         "%.48s: readiness timeout after %d seconds",
                         pendingApp ? pendingApp->name : "app",
                         monitor->config->readyTimeoutSec);
                reason = appReason;
            }
        }
    } else {
        monitor->bootDeadline = 0;
    }

    monitor->bootState = state;
    copy_reason(monitor->bootReason,
                sizeof(monitor->bootReason),
                reason);
    pthread_mutex_unlock(&monitor->mutex);
}

static void aggregate(ReadinessMonitor *monitor) {

    Space *space;
    App *app;
    NodeReadinessState state;
    const char *reason;
    char appReason[STARTERD_READY_REASON_LEN];
    App *pendingApp;
    time_t now;

    state = NODE_READINESS_READY;
    reason = "ready";
    appReason[0] = '\0';
    pendingApp = NULL;
    now = time(NULL);

    pthread_mutex_lock(&monitor->mutex);

    if (!monitor->ctx->bootCompleted) {
        state = NODE_READINESS_PENDING;
        reason = "starterd is starting applications";
        monitor->readinessDeadline = 0;
    } else if (monitor->enabled) {
        space = monitor->spaceList;
        while (space && state != NODE_READINESS_FAULTY) {
            app = space->appList;
            while (app) {
                if (!app->readinessRequired) {
                    app = app->next;
                    continue;
                }

                if (app->readinessState == APP_READINESS_FAULTY) {
                    state = NODE_READINESS_FAULTY;
                    snprintf(appReason,
                             sizeof(appReason),
                             "%.48s: %.140s",
                             app->name,
                             app->readinessReason);
                    reason = appReason;
                    break;
                }

                if (app->readinessState != APP_READINESS_READY &&
                    state == NODE_READINESS_READY) {
                    state = NODE_READINESS_PENDING;
                    pendingApp = app;
                    snprintf(appReason,
                             sizeof(appReason),
                             "%.48s: %.140s",
                             app->name,
                             app->readinessReason);
                    reason = appReason;
                }
                app = app->next;
            }
            space = space->next;
        }

        if (state == NODE_READINESS_READY ||
            state == NODE_READINESS_FAULTY) {
            monitor->readinessDeadline = 0;
        } else {
            if (monitor->readinessDeadline == 0) {
                monitor->readinessDeadline =
                    now + monitor->config->readyTimeoutSec;
            }

            if (now >= monitor->readinessDeadline) {
                state = NODE_READINESS_FAULTY;
                snprintf(appReason,
                         sizeof(appReason),
                         "%.48s: readiness timeout after %d seconds",
                         pendingApp ? pendingApp->name : "app",
                         monitor->config->readyTimeoutSec);
                reason = appReason;
            }
        }
    }

    monitor->state = state;
    copy_reason(monitor->reason, sizeof(monitor->reason), reason);
    pthread_mutex_unlock(&monitor->mutex);
}

static void poll_apps(ReadinessMonitor *monitor) {

    Space *space;
    App *app;
    AppReadyResponse result;
    bool responded;
    bool bootStarted;
    bool pollSpace;
    bool appRunning;
    uint32_t generation;
    time_t now;

    pthread_mutex_lock(&monitor->mutex);
    bootStarted = monitor->bootStarted;
    pthread_mutex_unlock(&monitor->mutex);

    if (!bootStarted && !monitor->ctx->bootCompleted) return;

    space = monitor->spaceList;
    while (space) {
        pollSpace = monitor->ctx->bootCompleted ||
                    (bootStarted &&
                     strcmp(space->name,
                            monitor->config->bootSpace) == 0);

        if (!pollSpace) {
            space = space->next;
            continue;
        }

        app = space->appList;
        while (app) {
            if (app->readinessRequired) {
                pthread_mutex_lock(&monitor->mutex);
                generation = app->readinessGeneration;
                appRunning = app->state == APP_STATE_RUNNING &&
                             app->pid > 0;
                pthread_mutex_unlock(&monitor->mutex);

                if (!appRunning) {
                    app = app->next;
                    continue;
                }

                memset(&result, 0, sizeof(result));
                responded = wc_app_ready(monitor->config,
                                         app,
                                         &result);
                now = time(NULL);
                app_update(monitor,
                           app,
                           &result,
                           responded,
                           generation,
                           now);
            }
            app = app->next;
        }
        space = space->next;
    }
}

static void poll_mesh(ReadinessMonitor *monitor) {

    Space *space;
    App *app;
    bool connected;
    bool known;
    char reason[STARTERD_READY_REASON_LEN];

    known = false;
    connected = false;
    snprintf(reason, sizeof(reason), "meshd is not in the manifest");

    space = monitor->spaceList;
    while (space && !known) {
        app = space->appList;
        while (app) {
            if (app_is_mesh(app)) {
                known = wc_mesh_status(monitor->config,
                                       app,
                                       &connected,
                                       reason,
                                       sizeof(reason));
                if (!known) {
                    snprintf(reason,
                             sizeof(reason),
                             "meshd status unavailable");
                }
                break;
            }
            app = app->next;
        }
        space = space->next;
    }

    pthread_mutex_lock(&monitor->mutex);
    monitor->meshKnown = known;
    monitor->meshConnected = connected;
    monitor->meshCheckedAt = time(NULL);
    copy_reason(monitor->meshReason,
                sizeof(monitor->meshReason),
                reason);
    pthread_mutex_unlock(&monitor->mutex);
}

static bool wait_for_next_poll(ReadinessMonitor *monitor) {

    struct timespec deadline;
    bool running;

    clock_gettime(CLOCK_REALTIME, &deadline);
    deadline.tv_sec += monitor->config->readyPollSec;

    pthread_mutex_lock(&monitor->mutex);
    if (monitor->running) {
        pthread_cond_timedwait(&monitor->condition,
                               &monitor->mutex,
                               &deadline);
    }
    running = monitor->running;
    pthread_mutex_unlock(&monitor->mutex);

    return running;
}

static void *readiness_thread(void *arg) {

    ReadinessMonitor *monitor;

    monitor = (ReadinessMonitor *)arg;

    do {
        poll_apps(monitor);
        aggregate_boot(monitor);
        poll_mesh(monitor);
        aggregate(monitor);
    } while (wait_for_next_poll(monitor));

    return NULL;
}

ReadinessMonitor *readiness_start(Config *config,
                                  Space *spaceList,
                                  StarterContext *ctx) {

    ReadinessMonitor *monitor;
    Space *space;
    App *app;

    if (!config || !spaceList || !ctx) return NULL;

    monitor = calloc(1, sizeof(*monitor));
    if (!monitor) return NULL;

    monitor->config = config;
    monitor->spaceList = spaceList;
    monitor->ctx = ctx;
    monitor->state = NODE_READINESS_PENDING;
    monitor->bootState = NODE_READINESS_PENDING;
    monitor->running = true;
    copy_reason(monitor->reason,
                sizeof(monitor->reason),
                "starterd is starting applications");
    copy_reason(monitor->bootReason,
                sizeof(monitor->bootReason),
                "starterd is starting boot applications");
    copy_reason(monitor->meshReason,
                sizeof(monitor->meshReason),
                "meshd status unavailable");

    space = spaceList;
    while (space) {
        app = space->appList;
        while (app) {
            if (app->readinessRequired) {
                monitor->enabled = true;
                if (strcmp(space->name, config->bootSpace) == 0) {
                    monitor->bootEnabled = true;
                }
            }
            app = app->next;
        }
        space = space->next;
    }

    pthread_mutex_init(&monitor->mutex, NULL);
    pthread_cond_init(&monitor->condition, NULL);

    if (pthread_create(&monitor->thread,
                       NULL,
                       readiness_thread,
                       monitor) != 0) {
        pthread_cond_destroy(&monitor->condition);
        pthread_mutex_destroy(&monitor->mutex);
        free(monitor);
        return NULL;
    }

    usys_log_info("readiness: monitor %s, timeout=%d sec poll=%d sec",
                  monitor->enabled ? "enabled" : "dormant",
                  config->readyTimeoutSec,
                  config->readyPollSec);
    return monitor;
}

void readiness_stop(ReadinessMonitor *monitor) {

    if (!monitor) return;

    pthread_mutex_lock(&monitor->mutex);
    monitor->running = false;
    pthread_cond_signal(&monitor->condition);
    pthread_mutex_unlock(&monitor->mutex);

    pthread_join(monitor->thread, NULL);
    pthread_cond_destroy(&monitor->condition);
    pthread_mutex_destroy(&monitor->mutex);
    free(monitor);
}

NodeReadinessState readiness_get(ReadinessMonitor *monitor,
                                 char *reason,
                                 size_t reasonSize) {

    NodeReadinessState state;

    if (!monitor) {
        copy_reason(reason, reasonSize, "readiness monitor unavailable");
        return NODE_READINESS_FAULTY;
    }

    pthread_mutex_lock(&monitor->mutex);
    state = monitor->state;
    copy_reason(reason, reasonSize, monitor->reason);
    pthread_mutex_unlock(&monitor->mutex);

    return state;
}

NodeReadinessState readiness_get_boot(ReadinessMonitor *monitor,
                                      char *reason,
                                      size_t reasonSize) {

    NodeReadinessState state;

    if (!monitor) {
        copy_reason(reason, reasonSize, "readiness monitor unavailable");
        return NODE_READINESS_FAULTY;
    }

    pthread_mutex_lock(&monitor->mutex);
    state = monitor->bootState;
    copy_reason(reason, reasonSize, monitor->bootReason);
    pthread_mutex_unlock(&monitor->mutex);

    return state;
}

void readiness_boot_started(ReadinessMonitor *monitor) {

    if (!monitor) return;

    pthread_mutex_lock(&monitor->mutex);
    monitor->bootStarted = true;
    monitor->bootDeadline = 0;
    pthread_cond_signal(&monitor->condition);
    pthread_mutex_unlock(&monitor->mutex);
}

void readiness_app_started(ReadinessMonitor *monitor, App *app) {

    time_t now;

    if (!monitor || !app || !app->readinessRequired) return;

    now = time(NULL);

    pthread_mutex_lock(&monitor->mutex);
    app->readinessGeneration = app->generation;
    app_set_pending(app, 0, "waiting for startup", now);
    app->readinessRequestId[0] = '\0';
    pthread_cond_signal(&monitor->condition);
    pthread_mutex_unlock(&monitor->mutex);
}

void readiness_app_exited(ReadinessMonitor *monitor, App *app) {

    time_t now;

    if (!monitor || !app || !app->readinessRequired) return;

    now = time(NULL);

    pthread_mutex_lock(&monitor->mutex);
    if (app->readinessState != APP_READINESS_FAULTY) {
        app->readinessSince = now;
    }
    app->readinessState = APP_READINESS_FAULTY;
    app->readinessHttpStatus = 0;
    app->readinessCheckedAt = now;
    copy_reason(app->readinessReason,
                sizeof(app->readinessReason),
                "application exited; restart pending");
    app->readinessRequestId[0] = '\0';
    pthread_cond_signal(&monitor->condition);
    pthread_mutex_unlock(&monitor->mutex);
}

json_t *readiness_status_json(ReadinessMonitor *monitor) {

    json_t *root;
    json_t *apps;
    json_t *entry;
    Space *space;
    App *app;

    root = json_object();
    apps = json_array();
    if (!root || !apps) {
        json_decref(root);
        json_decref(apps);
        return NULL;
    }

    if (!monitor) {
        json_object_set_new(root, "enabled", json_false());
        json_object_set_new(root, "state", json_string("faulty"));
        json_object_set_new(root,
                            "reason",
                            json_string("readiness monitor unavailable"));
        json_object_set_new(root, "apps", apps);
        return root;
    }

    pthread_mutex_lock(&monitor->mutex);

    json_object_set_new(root,
                        "enabled",
                        json_boolean(monitor->enabled));
    json_object_set_new(root,
                        "state",
                        json_string(node_readiness_str(monitor->state)));
    json_object_set_new(root,
                        "reason",
                        json_string(monitor->reason));
    json_object_set_new(root,
                        "timeoutSec",
                        json_integer(monitor->config->readyTimeoutSec));
    json_object_set_new(root,
                        "deadline",
                        monitor->readinessDeadline > 0 ?
                        json_integer(monitor->readinessDeadline) :
                        json_null());
    space = monitor->spaceList;
    while (space) {
        app = space->appList;
        while (app) {
            if (!app->readinessRequired) {
                app = app->next;
                continue;
            }

            entry = json_object();
            json_object_set_new(entry,
                                "space",
                                json_string(app->space));
            json_object_set_new(entry,
                                "name",
                                json_string(app->name));
            json_object_set_new(entry,
                                "service",
                                json_string(app->service));
            json_object_set_new(entry,
                                "state",
                                json_string(app_readiness_str(
                                    app->readinessState)));
            json_object_set_new(entry,
                                "httpStatus",
                                json_integer(app->readinessHttpStatus));
            json_object_set_new(entry,
                                "reason",
                                json_string(app->readinessReason));
            json_object_set_new(entry,
                                "requestId",
                                app->readinessRequestId[0] ?
                                json_string(app->readinessRequestId) :
                                json_null());
            json_object_set_new(entry,
                                "since",
                                json_integer(app->readinessSince));
            json_object_set_new(entry,
                                "checkedAt",
                                json_integer(app->readinessCheckedAt));
            json_array_append_new(apps, entry);
            app = app->next;
        }
        space = space->next;
    }

    json_object_set_new(root, "apps", apps);
    pthread_mutex_unlock(&monitor->mutex);

    return root;
}

json_t *readiness_connectivity_json(ReadinessMonitor *monitor) {

    json_t *root;

    root = json_object();
    if (!root) return NULL;

    if (!monitor) {
        json_object_set_new(root, "known", json_false());
        json_object_set_new(root, "connected", json_false());
        json_object_set_new(root,
                            "reason",
                            json_string("readiness monitor unavailable"));
        return root;
    }

    pthread_mutex_lock(&monitor->mutex);
    json_object_set_new(root,
                        "known",
                        json_boolean(monitor->meshKnown));
    json_object_set_new(root,
                        "connected",
                        json_boolean(monitor->meshConnected));
    json_object_set_new(root,
                        "reason",
                        json_string(monitor->meshReason));
    json_object_set_new(root,
                        "checkedAt",
                        json_integer(monitor->meshCheckedAt));
    pthread_mutex_unlock(&monitor->mutex);

    return root;
}
