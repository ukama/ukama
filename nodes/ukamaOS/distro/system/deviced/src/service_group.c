/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/file.h>
#include <sys/stat.h>
#include <unistd.h>

#include "service_api.h"
#include "service_client.h"
#include "service_group.h"
#include "service_lease.h"

#define TRANSITION_MSEC     45000
#define GROUP_POLL_SEC      1
#define CLOCK_NSEC_PER_SEC  1000000000
#define SERVICE_RUNTIME_DIR "/run/ukama"
#define WATCHDOG_MSEC       10000

/* Lock protects intent and observations. Only worker writes component
 * commands, and it checks the token between calls. OFF revokes permission
 * immediately, including while an ON command is in flight.
 */
struct ServiceGroup {
    Config *config;
    pthread_t worker;
    bool stopping;
    bool starting;
    bool healthy;
    bool held;
    uint64_t watchdogAt;
    uint64_t deadline;
    int lockFd;
    char reason[128];
    unsigned long long holdToken;
    char holdApp[64];
    char holdPath[512];
    char lockPath[512];
    ServicePeer peers[SERVICE_MEMBER_COUNT];
};

typedef struct ServiceGroup ServiceGroup;

static bool write_lease(void) {

    char path[512];
    char text[32];
    int fd;
    int length;
    bool ok;

    if (snprintf(path, sizeof(path), "%s.tmp", service_lease_path()) >=
        (int)sizeof(path)) {
        return false;
    }

    fd = open(path, O_WRONLY | O_CREAT | O_TRUNC | O_NOFOLLOW, 0600);

    if (fd < 0) {
        return false;
    }

    length =
        snprintf(text, sizeof(text), "%llu\n",
                 (unsigned long long)service_clock_msec() + SERVICE_LEASE_MSEC);
    ok = write(fd, text, length) == length;

    if (close(fd) != 0) {
        ok = false;
    }

    if (ok) {
        ok = rename(path, service_lease_path()) == 0;
    }

    if (!ok) {
        unlink(path);
    }

    return ok;
}

static bool write_hold(ServiceGroup *group, const char *app,
                       unsigned long long token) {

    char path[sizeof(group->holdPath) + 8];
    char text[128];
    int fd;
    int length;
    bool ok;

    snprintf(path, sizeof(path), "%s.tmp", group->holdPath);
    fd = open(path, O_WRONLY | O_CREAT | O_TRUNC | O_NOFOLLOW, 0600);

    if (fd < 0) {
        return false;
    }

    length = snprintf(text, sizeof(text), "%s %llu\n", app, token);
    ok     = write(fd, text, length) == length;

    if (close(fd) != 0) {
        ok = false;
    }

    if (ok) {
        ok = rename(path, group->holdPath) == 0;
    }

    if (!ok) {
        unlink(path);
    }

    return ok;
}

static bool watchdog_live(ServiceGroup *group) {

    uint64_t now = service_clock_msec();

    return group->watchdogAt != 0 && now >= group->watchdogAt &&
           now - group->watchdogAt < WATCHDOG_MSEC;
}

/* Called with control->Lock held. Use the same gate for status and ON. */
static bool ready_to_enable(ServiceGroup *group) {

    ControlCtx *control = group->config->control;

    return !group->held && group->healthy && watchdog_live(group) &&
           control->Service.Phase == CONTROL_PHASE_IDLE &&
           control->Active == CONTROL_SUBSYS_NONE &&
           control->Radio.Phase != CONTROL_PHASE_PENDING &&
           control->Radio.Phase != CONTROL_PHASE_EXECUTING &&
           control->Reboot.Phase != CONTROL_PHASE_PENDING &&
           control->Reboot.Phase != CONTROL_PHASE_EXECUTING;
}

/* Called with control->Lock held. Revocation must precede teardown. */
static void request_off(ServiceGroup *group, const char *reason) {

    ControlCtx *control       = group->config->control;
    ControlSubsysState *state = &control->Service;

    unlink(service_lease_path());
    group->starting = false;
    state->Desired = CONTROL_STATE_OFF;
    state->Phase   = CONTROL_PHASE_PENDING;
    state->Token++;

    if (control->Active == CONTROL_SUBSYS_NONE) {
        control->Active = CONTROL_SUBSYS_SERVICE;
    }

    group->deadline = service_clock_msec() + TRANSITION_MSEC;
    snprintf(group->reason, sizeof(group->reason), "%s", reason);
}

static bool may_enable(ServiceGroup *group, unsigned long long token) {

    ControlCtx *control = group->config->control;
    bool allowed;

    pthread_mutex_lock(&control->Lock);
    allowed = !group->stopping && !group->held && group->healthy &&
              watchdog_live(group) && control->Service.Token == token &&
              control->Service.Desired == CONTROL_STATE_ON &&
              service_clock_msec() < group->deadline;

    if (allowed) {
        allowed = write_lease();
    }

    pthread_mutex_unlock(&control->Lock);

    return allowed;
}

static bool confirmed(ServiceGroup *group, bool enabled) {

    ServiceObservation observation;
    size_t i;
    bool ok = true;

    for (i = 0; i < SERVICE_MEMBER_COUNT; i++) {
        observation = service_peer_read(&group->peers[i]);

        if (!observation.reachable || !observation.confirmed ||
            observation.on != enabled || (enabled && !observation.ready)) {
            ok = false;
        }
    }

    return ok;
}

static bool disable_group(ServiceGroup *group) {

    size_t i;

    /* The network guard is first. Still try every member after a failure. */

    for (i = 0; i < SERVICE_MEMBER_COUNT; i++) {
        service_peer_set(&group->peers[i], false);
    }

    return confirmed(group, false);
}

static bool enable_group(ServiceGroup *group, unsigned long long token) {

    const ServiceMember order[] = {SERVICE_MEMBER_PCRF, SERVICE_MEMBER_EPC,
                                   SERVICE_MEMBER_NETWORK};
    ServiceObservation observation;
    size_t i;

    for (i = 0; i < SERVICE_MEMBER_COUNT; i++) {
        observation = service_peer_read(&group->peers[i]);

        if (!observation.reachable || !observation.ready ||
            !observation.confirmed || observation.on) {
            return false;
        }
    }

    for (i = 0; i < sizeof(order) / sizeof(order[0]); i++) {
        if (!may_enable(group, token) ||
            !service_peer_set(&group->peers[order[i]], true)) {
            return false;
        }
    }

    return may_enable(group, token) && confirmed(group, true);
}

static void *group_worker(void *arg) {

    ServiceGroup *group       = arg;
    ControlCtx *control       = group->config->control;
    ControlSubsysState *state = &control->Service;
    unsigned long long token;
    bool enabled;
    bool ok;

    for (;;) {
        pthread_mutex_lock(&control->Lock);

        if (group->stopping) {
            pthread_mutex_unlock(&control->Lock);
            break;
        }

        if (state->Desired == CONTROL_STATE_ON &&
            (!watchdog_live(group) || !group->healthy || group->held)) {
            request_off(group, "watchdog_unavailable_or_unhealthy");
        }

        if (state->Phase == CONTROL_PHASE_IDLE) {
            if (state->Current == CONTROL_STATE_ON && !write_lease()) {
                request_off(group, "service_permission_failed");
            } else {
                pthread_mutex_unlock(&control->Lock);
                sleep(GROUP_POLL_SEC);
                continue;
            }
        }

        token        = state->Token;
        enabled      = state->Desired == CONTROL_STATE_ON;

        if (!enabled && group->starting &&
            service_clock_msec() >= group->deadline) {
            group->starting = false;
            state->Phase    = CONTROL_PHASE_FAULT;

            if (control->Active == CONTROL_SUBSYS_SERVICE) {
                control->Active = CONTROL_SUBSYS_NONE;
            }
        } else if (state->Phase != CONTROL_PHASE_FAULT) {
            state->Phase = CONTROL_PHASE_EXECUTING;
        }

        pthread_mutex_unlock(&control->Lock);

        ok = enabled ? enable_group(group, token) : disable_group(group);

        pthread_mutex_lock(&control->Lock);

        if (state->Token == token) {
            if (ok) {
                group->starting = false;
                state->Current = enabled ? CONTROL_STATE_ON : CONTROL_STATE_OFF;
                state->Phase   = CONTROL_PHASE_IDLE;

                if (control->Active == CONTROL_SUBSYS_SERVICE) {
                    control->Active = CONTROL_SUBSYS_NONE;
                }

                group->deadline = 0;
            } else if (enabled) {
                request_off(group, "service_on_failed");
            } else if (!group->starting ||
                       service_clock_msec() >= group->deadline) {
                group->starting = false;
                state->Phase    = CONTROL_PHASE_FAULT;

                if (control->Active == CONTROL_SUBSYS_SERVICE) {
                    control->Active = CONTROL_SUBSYS_NONE;
                }
            }
        }

        pthread_mutex_unlock(&control->Lock);
        sleep(GROUP_POLL_SEC);
    }

    return NULL;
}

bool service_group_start(Config *config) {

    ServiceGroup *group;
    FILE *hold;
    struct timespec now;

    if (mkdir(SERVICE_RUNTIME_DIR, 0755) != 0 && errno != EEXIST) {
        return false;
    }

    group = calloc(1, sizeof(*group));

    if (group == NULL) {
        return false;
    }

    if (snprintf(group->holdPath, sizeof(group->holdPath), "%s.hold",
                 service_lease_path()) >= (int)sizeof(group->holdPath) ||
        snprintf(group->lockPath, sizeof(group->lockPath), "%s.lock",
                 service_lease_path()) >= (int)sizeof(group->lockPath)) {
        free(group);

        return false;
    }

    group->lockFd = open(group->lockPath, O_CREAT | O_RDWR | O_CLOEXEC, 0600);

    if (group->lockFd < 0 || flock(group->lockFd, LOCK_EX | LOCK_NB) != 0 ||
        !service_peers_init(group->peers)) {
        if (group->lockFd >= 0) {
            close(group->lockFd);
        }

        free(group);

        return false;
    }

    group->config = config;
    group->held   = access(group->holdPath, F_OK) == 0;

    if (group->held) {
        hold = fopen(group->holdPath, "r");

        if (hold != NULL) {
            if (fscanf(hold, "%63s %llu", group->holdApp, &group->holdToken) !=
                2) {
                group->holdApp[0] = '\0';
            }

            fclose(hold);
        }
    }

    config->serviceGroup         = group;
    config->control->ServiceHeld = group->held;

    if (clock_gettime(CLOCK_MONOTONIC, &now) != 0) {
        close(group->lockFd);
        free(group);
        config->serviceGroup = NULL;

        return false;
    }

    config->control->Service.Token =
        (uint64_t)now.tv_sec * CLOCK_NSEC_PER_SEC + now.tv_nsec;
    request_off(group, "startup");
    group->starting = true;

    if (pthread_create(&group->worker, NULL, group_worker, group) != 0) {
        config->serviceGroup = NULL;
        close(group->lockFd);
        free(group);

        return false;
    }

    return true;
}

void service_group_stop(Config *config) {

    ServiceGroup *group = config->serviceGroup;

    if (group == NULL) {
        return;
    }

    pthread_mutex_lock(&config->control->Lock);
    group->stopping = true;
    request_off(group, "controller_shutdown");
    pthread_mutex_unlock(&config->control->Lock);

    pthread_join(group->worker, NULL);
    disable_group(group);

    close(group->lockFd);
    free(group);
    config->serviceGroup = NULL;
}

JsonObj *service_group_state(Config *config) {

    ServiceGroup *group       = config->serviceGroup;
    ControlSubsysState *state = &config->control->Service;
    JsonObj *json;
    const char *phase;
    const char *current;
    const char *desired;
    uint64_t now;
    json_int_t remaining;

    if (group == NULL) {
        return NULL;
    }

    pthread_mutex_lock(&config->control->Lock);
    if (state->Phase == CONTROL_PHASE_IDLE) {
        phase = "idle";
    } else if (state->Phase == CONTROL_PHASE_FAULT) {
        phase = "fault";
    } else {
        phase = "changing";
    }

    current = state->Current == CONTROL_STATE_ON ? SERVICE_STATE_ON :
                                                  SERVICE_STATE_OFF;
    desired = state->Desired == CONTROL_STATE_ON ? SERVICE_STATE_ON :
                                                  SERVICE_STATE_OFF;
    now       = service_clock_msec();
    remaining = group->deadline > now ? group->deadline - now : 0;

    json = json_pack("{s:s,s:s,s:s,s:I,s:I,s:s,s:b,s:b}",
                     "current", current,
                     "desired", desired,
                     "phase", phase,
                     "token", (json_int_t)state->Token,
                     "remaining_ms", remaining,
                     "reason", group->reason,
                     "maintenance", group->held,
                     "ready", ready_to_enable(group));
    pthread_mutex_unlock(&config->control->Lock);

    return json;
}

int service_group_command(const URequest *request, UResponse *response,
                          void *data) {

    Config *config      = data;
    ServiceGroup *group = config->serviceGroup;
    ControlCtx *control = config->control;
    JsonObj *body       = ulfius_get_json_body_request(request, NULL);
    const char *desired = json_string_value(json_object_get(body, "state"));
    JsonObj *force      = json_object_get(body, "force");
    int code            = HttpStatus_Accepted;

    if ((force != NULL && !json_is_boolean(force)) || desired == NULL ||
        (strcasecmp(desired, "on") && strcasecmp(desired, "off"))) {
        code = HttpStatus_BadRequest;
    } else if (group == NULL) {
        code = HttpStatus_ServiceUnavailable;
    } else {
        pthread_mutex_lock(&control->Lock);

        if (!strcasecmp(desired, "off")) {
            request_off(group, "operator_off");
        } else if (!ready_to_enable(group)) {
            code = HttpStatus_Conflict;
        } else if (control->Service.Current != CONTROL_STATE_ON) {
            control->Service.Desired = CONTROL_STATE_ON;
            control->Service.Phase   = CONTROL_PHASE_PENDING;
            control->Service.Token++;
            control->Active = CONTROL_SUBSYS_SERVICE;
            group->deadline = service_clock_msec() + TRANSITION_MSEC;
            snprintf(group->reason, sizeof(group->reason), "operator_on");
        }

        pthread_mutex_unlock(&control->Lock);
    }

    json_decref(body);
    ulfius_set_empty_body_response(response, code);

    return U_CALLBACK_CONTINUE;
}

int service_group_watchdog(const URequest *request, UResponse *response,
                           void *data) {

    Config *config            = data;
    ServiceGroup *group       = config->serviceGroup;
    ControlSubsysState *state = &config->control->Service;
    JsonObj *body             = ulfius_get_json_body_request(request, NULL);
    JsonObj *token            = json_object_get(body, "token");
    JsonObj *healthy          = json_object_get(body, "healthy");
    JsonObj *fault            = json_object_get(body, "fault");
    int code                  = HttpStatus_Accepted;

    if (group == NULL || !json_is_integer(token) || !json_is_boolean(healthy) ||
        !json_is_boolean(fault)) {
        code = HttpStatus_BadRequest;
    } else {
        pthread_mutex_lock(&config->control->Lock);

        if ((unsigned long long)json_integer_value(token) != state->Token) {
            code = HttpStatus_Conflict;
        } else {
            group->watchdogAt = service_clock_msec();
            group->healthy    = json_is_true(healthy);

            if (json_is_true(fault) && (state->Desired == CONTROL_STATE_ON ||
                                        state->Phase == CONTROL_PHASE_IDLE)) {
                request_off(group, "component_fault");
            }
        }

        pthread_mutex_unlock(&config->control->Lock);
    }

    json_decref(body);
    ulfius_set_empty_body_response(response, code);

    return U_CALLBACK_CONTINUE;
}

/* The supervisor reports app names; membership belongs to the service layer. */
static bool critical_app(const char *name) {

    const char *members[] = {"pcrf",         "epcemu",  "epc",
                             "init-network", "deviced", "service-watchdog",
                             "starterd"};
    size_t i;

    for (i = 0; i < sizeof(members) / sizeof(members[0]); i++) {
        if (!strcmp(name, members[i])) {
            return true;
        }
    }

    return false;
}

int service_group_maintenance(const URequest *request, UResponse *response,
                              void *data) {

    Config *config                = data;
    ServiceGroup *group           = config->serviceGroup;
    JsonObj *body                 = ulfius_get_json_body_request(request, NULL);
    const char *app = json_string_value(json_object_get(body, "app"));
    JsonObj *active               = json_object_get(body, "active");
    JsonObj *token                = json_object_get(body, "token");
    unsigned long long ownerToken = json_integer_value(token);
    int code                      = HttpStatus_OK;

    if (group == NULL || app == NULL || !json_is_boolean(active) ||
        !json_is_integer(token) || json_integer_value(token) < 0 ||
        (json_is_true(active) && ownerToken == 0)) {
        code = HttpStatus_BadRequest;
    } else if (critical_app(app)) {
        pthread_mutex_lock(&config->control->Lock);

        if (json_is_true(active)) {
            if (group->held && (strcmp(group->holdApp, app) ||
                                group->holdToken != ownerToken)) {
                code = HttpStatus_Conflict;
            } else {
                if (!group->held && !write_hold(group, app, ownerToken)) {
                    code = HttpStatus_ServiceUnavailable;
                } else {
                    if (!group->held) {
                        request_off(group, "maintenance");
                    }

                    group->held                  = true;
                    config->control->ServiceHeld = true;
                    group->holdToken             = ownerToken;
                    snprintf(group->holdApp, sizeof(group->holdApp), "%s", app);

                    if (code == HttpStatus_OK &&
                        (config->control->Service.Phase != CONTROL_PHASE_IDLE ||
                         config->control->Service.Current !=
                             CONTROL_STATE_OFF)) {
                        code = HttpStatus_Accepted;
                    }
                }
            }
        } else if (group->held && !strcmp(group->holdApp, app) &&
                   (ownerToken == group->holdToken || ownerToken == 0)) {
            /* Token zero is startup recovery by a new supervisor. It may
             * release an interrupted operation only after confirmed OFF. */

            if (ownerToken == 0 &&
                (config->control->Service.Phase != CONTROL_PHASE_IDLE ||
                 config->control->Service.Current != CONTROL_STATE_OFF)) {
                code = HttpStatus_Accepted;
            } else if (unlink(group->holdPath) == 0 || errno == ENOENT) {
                group->held                  = false;
                config->control->ServiceHeld = false;
                group->holdToken             = 0;
                group->holdApp[0]            = '\0';
            } else {
                code = HttpStatus_ServiceUnavailable;
            }
        }

        pthread_mutex_unlock(&config->control->Lock);
    }

    json_decref(body);
    ulfius_set_empty_body_response(response, code);

    return U_CALLBACK_CONTINUE;
}
