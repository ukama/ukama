/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <pthread.h>
#include <string.h>
#include <unistd.h>

#include "web_service.h"

#include "actions.h"
#include "config.h"
#include "control.h"
#include "deviced.h"
#include "http_status.h"
#include "web_client.h"

#include "usys_log.h"
#include "usys_mem.h"

#include "version.h"

typedef struct {
    Config *Config;
    ControlSubsystem Subsystem;
    ControlState Desired;
    bool Immediate;
    unsigned long long Token;
} WorkerArgs;

static int json_set_empty(UResponse *response, int status) {

    ulfius_set_empty_body_response(response, status);
    return U_CALLBACK_CONTINUE;
}

static bool _parse_state_request(const URequest *request,
                                 ControlState *desired,
                                 bool *force) {

    JsonErrObj err;
    JsonObj *json = NULL;
    json_t *jState = NULL;
    json_t *jForce = NULL;
    const char *stateStr = NULL;

    if (!request || !desired || !force) return false;

    *force = false;

    if (!request->binary_body || request->binary_body_length == 0) {
        return false;
    }

    memset(&err, 0, sizeof(err));

    json = json_loadb((const char *)request->binary_body,
                      request->binary_body_length,
                      0,
                      &err);
    if (!json) {
        return false;
    }

    jState = json_object_get(json, "state");
    if (!jState || !json_is_string(jState)) {
        json_decref(json);
        return false;
    }

    stateStr = json_string_value(jState);
    if (!stateStr) {
        json_decref(json);
        return false;
    }

    if (strcasecmp(stateStr, "on") == 0) {
        *desired = CONTROL_STATE_ON;
    } else if (strcasecmp(stateStr, "off") == 0) {
        *desired = CONTROL_STATE_OFF;
    } else {
        json_decref(json);
        return false;
    }

    jForce = json_object_get(json, "force");
    if (jForce) {
        if (!json_is_boolean(jForce)) {
            json_decref(json);
            return false;
        }
        *force = json_is_true(jForce) ? true : false;
    }

    json_decref(json);
    return true;
}

static void* _worker_run(void *arg) {

    WorkerArgs *args = NULL;
    Config *config = NULL;
    ControlCtx *control = NULL;
    ControlSubsysState *ss = NULL;
    int retCode = -1;
    int execRet = STATUS_NOK;
    int delay = 0;
    ControlState desired = CONTROL_STATE_OFF;

    args = (WorkerArgs *)arg;
    if (!args) {
        pthread_exit(NULL);
    }

    config = args->Config;
    if (!config || !config->control || !config->nodeType) {
        usys_free(args);
        pthread_exit(NULL);
    }

    control = config->control;

    delay = args->Immediate ? 0 : WAIT_BEFORE_REBOOT;
    if (delay > 0) {
        sleep(delay);
    }

    pthread_mutex_lock(&control->Lock);
    ss = NULL;
    if (args->Subsystem == CONTROL_SUBSYS_SERVICE) {
        ss = &control->Service;
    } else if (args->Subsystem == CONTROL_SUBSYS_RADIO) {
        ss = &control->Radio;
    } else if (args->Subsystem == CONTROL_SUBSYS_RESTART) {
        ss = &control->Restart;
    }

    if (!ss || ss->Phase != CONTROL_PHASE_PENDING || ss->Token != args->Token) {
        pthread_mutex_unlock(&control->Lock);
        usys_free(args);
        pthread_exit(NULL);
    }

    desired = ss->Desired;
    pthread_mutex_unlock(&control->Lock);

    if (!control_begin_execute(control, args->Subsystem, args->Token)) {
        usys_free(args);
        pthread_exit(NULL);
    }

    if (args->Subsystem == CONTROL_SUBSYS_RESTART) {
        (void)wc_send_action_alarm_to_notifyd(config,
                                              "restart",
                                              "Restarting the node",
                                              &retCode);
        execRet = actions_restart_apply(config);
        if (execRet != STATUS_OK) {
            control_mark_fault(control, args->Subsystem);
        }
        usys_free(args);
        pthread_exit(NULL);
    }

    if (args->Subsystem == CONTROL_SUBSYS_SERVICE) {
        (void)wc_send_action_alarm_to_notifyd(config,
                     (desired == CONTROL_STATE_ON) ? "service_on" : "service_off",
                                          (desired == CONTROL_STATE_ON) ?
                                          "Enabling cellular service" : "Disabling cellular service",
                                           &retCode);

        execRet = actions_service_apply(config, desired);
        if (execRet == STATUS_OK) {
            control_mark_done(control, args->Subsystem, desired);
        } else {
            control_mark_fault(control, args->Subsystem);
        }
        usys_free(args);
        pthread_exit(NULL);
    }

    if (args->Subsystem == CONTROL_SUBSYS_RADIO) {
        (void)wc_send_action_alarm_to_notifyd(config,
                                              (desired == CONTROL_STATE_ON) ?
                                              "radio_on" : "radio_off",
                                              (desired == CONTROL_STATE_ON) ?
                                              "Enabling radio" : "Disabling radio",
                                              &retCode);

        execRet = actions_radio_apply(config, desired);
        if (execRet == STATUS_OK) {
            control_mark_done(control, args->Subsystem, desired);
        } else {
            control_mark_fault(control, args->Subsystem);
        }
        usys_free(args);
        pthread_exit(NULL);
    }

    usys_free(args);
    pthread_exit(NULL);
}

static int _schedule_worker(Config *config,
                            ControlSubsystem subsystem,
                            bool immediate,
                            unsigned long long token) {

    pthread_t thread;
    WorkerArgs *args = NULL;
    int ret = 0;

    if (!config || !config->control) return STATUS_NOK;

    args = (WorkerArgs *)usys_malloc(sizeof(WorkerArgs));
    if (!args) return STATUS_NOK;

    memset(args, 0, sizeof(*args));
    args->Config = config;
    args->Subsystem = subsystem;
    args->Immediate = immediate;
    args->Token = token;

    ret = pthread_create(&thread, NULL, _worker_run, (void *)args);
    if (ret != 0) {
        usys_free(args);
        return STATUS_NOK;
    }

    pthread_detach(thread);
    return STATUS_OK;
}

static int _post_state_change(const URequest *request,
                              UResponse *response,
                              Config *config,
                              ControlSubsystem subsystem) {

    ControlState desired = CONTROL_STATE_OFF;
    bool force = false;
    int httpStatus = HttpStatus_BadRequest;
    bool allowed = false;
    bool immediate = false;
    unsigned long long token = 0;

    if (!config || !config->control) {
        return json_set_empty(response, HttpStatus_InternalServerError);
    }

    if (!_parse_state_request(request, &desired, &force)) {
        return json_set_empty(response, HttpStatus_BadRequest);
    }

    allowed = control_set_pending(config->control,
                                 subsystem,
                                 desired,
                                 force,
                                 &httpStatus,
                                 &immediate,
                                 &token);

    if (httpStatus == HttpStatus_OK) {
        return json_set_empty(response, HttpStatus_OK);
    }

    if (!allowed) {
        return json_set_empty(response, httpStatus);
    }

    if (_schedule_worker(config, subsystem, immediate, token) != STATUS_OK) {
        control_mark_fault(config->control, subsystem);
        return json_set_empty(response, HttpStatus_InternalServerError);
    }

    return json_set_empty(response, HttpStatus_Accepted);
}

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

    (void)request;
    (void)epConfig;
    return json_set_empty(response, HttpStatus_OK);
}

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    (void)request;
    (void)epConfig;

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_state(const URequest *request,
                         UResponse *response,
                         void *epConfig) {

    Config *config = NULL;
    char state[32];
    JsonObj *json = NULL;
    time_t now = 0;
    long uptime = 0;

    (void)request;

    config = (Config *)epConfig;
    if (!config || !config->control || !config->nodeType) {
        return json_set_empty(response, HttpStatus_InternalServerError);
    }

    memset(state, 0, sizeof(state));
    if (control_get_public_state(config->control,
                                 config->nodeType,
                                 state,
                                 sizeof(state)) != STATUS_OK) {
        return json_set_empty(response, HttpStatus_InternalServerError);
    }

    now = time(NULL);
    if (config->startTime > 0 && now >= config->startTime) {
        uptime = (long)(now - config->startTime);
    }

    json = json_object();
    if (!json) {
        return json_set_empty(response, HttpStatus_InternalServerError);
    }

    if (strcmp(config->nodeType, UKAMA_TOWER_NODE) == 0) {
        json_object_set_new(json, "service", json_string(state));
    } else if (strcmp(config->nodeType, UKAMA_AMPLIFIER_NODE) == 0) {
        json_object_set_new(json, "radio", json_string(state));
    } else {
        json_decref(json);
        return json_set_empty(response, HttpStatus_BadRequest);
    }

    json_object_set_new(json, "uptime_s", json_integer(uptime));

    ulfius_set_json_body_response(response, HttpStatus_OK, json);
    json_decref(json);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_service(const URequest *request,
                               UResponse *response,
                               void *epConfig) {

    Config *config = NULL;

    config = (Config *)epConfig;
    if (!config || !config->nodeType) {
        return json_set_empty(response, HttpStatus_InternalServerError);
    }

    if (strcmp(config->nodeType, UKAMA_TOWER_NODE) != 0) {
        return json_set_empty(response, HttpStatus_BadRequest);
    }

    return _post_state_change(request, response, config, CONTROL_SUBSYS_SERVICE);
}

int web_service_cb_post_radio(const URequest *request,
                             UResponse *response,
                             void *epConfig) {

    Config *config = NULL;

    config = (Config *)epConfig;
    if (!config || !config->nodeType) {
        return json_set_empty(response, HttpStatus_InternalServerError);
    }

    if (strcmp(config->nodeType, UKAMA_AMPLIFIER_NODE) != 0) {
        return json_set_empty(response, HttpStatus_BadRequest);
    }

    return _post_state_change(request, response, config, CONTROL_SUBSYS_RADIO);
}

int web_service_cb_post_restart(const URequest *request,
                                UResponse *response,
                                void *epConfig) {

    Config *config = NULL;
    ControlState desired = CONTROL_STATE_OFF;
    bool force = false;
    int httpStatus = HttpStatus_BadRequest;
    bool allowed = false;
    bool immediate = false;
    unsigned long long token = 0;
    JsonErrObj err;
    JsonObj *json = NULL;
    json_t *jForce = NULL;

    (void)desired;

    config = (Config *)epConfig;
    if (!config || !config->control) {
        return json_set_empty(response, HttpStatus_InternalServerError);
    }

    force = false;
    if (request->binary_body && request->binary_body_length > 0) {
        memset(&err, 0, sizeof(err));
        json = json_loadb((const char *)request->binary_body,
                          request->binary_body_length,
                          0,
                          &err);
        if (!json) {
            return json_set_empty(response, HttpStatus_BadRequest);
        }
        jForce = json_object_get(json, "force");
        if (jForce) {
            if (!json_is_boolean(jForce)) {
                json_decref(json);
                return json_set_empty(response, HttpStatus_BadRequest);
            }
            force = json_is_true(jForce) ? true : false;
        }
        json_decref(json);
    }

    allowed = control_set_pending(config->control,
                                 CONTROL_SUBSYS_RESTART,
                                 CONTROL_STATE_OFF,
                                 force,
                                 &httpStatus,
                                 &immediate,
                                 &token);

    if (httpStatus == HttpStatus_OK) {
        return json_set_empty(response, HttpStatus_OK);
    }

    if (!allowed) {
        return json_set_empty(response, httpStatus);
    }

    if (_schedule_worker(config, CONTROL_SUBSYS_RESTART, immediate, token) != STATUS_OK) {
        control_mark_fault(config->control, CONTROL_SUBSYS_RESTART);
        return json_set_empty(response, HttpStatus_InternalServerError);
    }

    return json_set_empty(response, HttpStatus_Accepted);
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    (void)request;
    (void)epConfig;
    return json_set_empty(response, HttpStatus_NotFound);
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data) {

    (void)request;
    (void)user_data;
    return json_set_empty(response, HttpStatus_MethodNotAllowed);
}
