/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>

#include "control.h"

#include "deviced.h"
#include "http_status.h"

static ControlSubsysState *get_subsys(ControlCtx *ctx, ControlSubsystem subsystem) {

    ControlSubsysState *ss = NULL;

    if (!ctx) return NULL;

    switch (subsystem) {
    case CONTROL_SUBSYS_SERVICE:
        ss = &ctx->Service;
        break;
    case CONTROL_SUBSYS_RADIO:
        ss = &ctx->Radio;
        break;
    case CONTROL_SUBSYS_RESTART:
        ss = &ctx->Restart;
        break;
    default:
        ss = NULL;
        break;
    }

    return ss;
}

static bool is_active_phase(ControlPhase phase) {

    if (phase == CONTROL_PHASE_PENDING) return true;
    if (phase == CONTROL_PHASE_EXECUTING) return true;
    return false;
}

static bool is_pending_subsys(ControlSubsysState *ss) {

    if (!ss) return false;
    return ss->Phase == CONTROL_PHASE_PENDING;
}

static bool is_executing_subsys(ControlSubsysState *ss) {

    if (!ss) return false;
    return ss->Phase == CONTROL_PHASE_EXECUTING;
}

static void cancel_pending_subsys(ControlSubsysState *ss) {

    if (!ss) return;
    ss->Phase = CONTROL_PHASE_IDLE;
    ss->Token++;
}

ControlCtx *control_create(void) {

    ControlCtx *ctx = NULL;

    ctx = (ControlCtx *)calloc(1, sizeof(ControlCtx));
    if (!ctx) return NULL;

    pthread_mutex_init(&ctx->Lock, NULL);
    ctx->Active = CONTROL_SUBSYS_NONE;

    ctx->Service.Phase = CONTROL_PHASE_IDLE;
    ctx->Service.Current = CONTROL_STATE_OFF;
    ctx->Service.Desired = CONTROL_STATE_OFF;
    ctx->Service.Token = 1;

    ctx->Radio.Phase = CONTROL_PHASE_IDLE;
    ctx->Radio.Current = CONTROL_STATE_OFF;
    ctx->Radio.Desired = CONTROL_STATE_OFF;
    ctx->Radio.Token = 1;

    ctx->Restart.Phase = CONTROL_PHASE_IDLE;
    ctx->Restart.Current = CONTROL_STATE_OFF;
    ctx->Restart.Desired = CONTROL_STATE_OFF;
    ctx->Restart.Token = 1;

    return ctx;
}

void control_destroy(ControlCtx *ctx) {

    if (!ctx) return;
    pthread_mutex_destroy(&ctx->Lock);
    free(ctx);
}

bool control_is_busy(ControlCtx *ctx) {

    bool busy = false;

    if (!ctx) return false;

    pthread_mutex_lock(&ctx->Lock);
    busy = (ctx->Active != CONTROL_SUBSYS_NONE);
    pthread_mutex_unlock(&ctx->Lock);

    return busy;
}

int control_get_public_state(ControlCtx *ctx,
                             const char *nodeType,
                             char *outState,
                             size_t outStateSize) {

    ControlSubsysState *ss = NULL;
    const char *stateStr = NULL;

    if (!ctx || !nodeType || !outState || outStateSize == 0) return STATUS_NOK;

    pthread_mutex_lock(&ctx->Lock);

    if (strcmp(nodeType, UKAMA_TOWER_NODE) == 0) {
        ss = &ctx->Service;
    } else if (strcmp(nodeType, UKAMA_AMPLIFIER_NODE) == 0) {
        ss = &ctx->Radio;
    } else {
        pthread_mutex_unlock(&ctx->Lock);
        return STATUS_NOK;
    }

    if (ss->Phase == CONTROL_PHASE_FAULT) {
        stateStr = "fault";
    } else if (ss->Phase == CONTROL_PHASE_PENDING || ss->Phase == CONTROL_PHASE_EXECUTING) {
        stateStr = "transitioning";
    } else {
        stateStr = (ss->Current == CONTROL_STATE_ON) ? "on" : "off";
    }

    (void)snprintf(outState, outStateSize, "%s", stateStr);

    pthread_mutex_unlock(&ctx->Lock);

    return STATUS_OK;
}

bool control_set_pending(ControlCtx *ctx,
                         ControlSubsystem subsystem,
                         ControlState desired,
                         bool force,
                         int *httpStatus,
                         bool *runImmediate,
                         unsigned long long *outToken) {

    ControlSubsysState *ss = NULL;
    bool allowed = false;
    bool immediate = false;
    int status = HttpStatus_Conflict;

    if (!ctx || !httpStatus || !runImmediate || !outToken) return false;

    pthread_mutex_lock(&ctx->Lock);

    ss = get_subsys(ctx, subsystem);
    if (!ss) {
        status = HttpStatus_BadRequest;
        goto done;
    }

    if (ctx->Active != CONTROL_SUBSYS_NONE && ctx->Active != subsystem) {

        ControlSubsysState *activeSs = NULL;

        activeSs = get_subsys(ctx, ctx->Active);
        if (!force) {
            status = HttpStatus_Conflict;
            goto done;
        }

        if (is_executing_subsys(activeSs)) {
            status = HttpStatus_Conflict;
            goto done;
        }

        if (is_pending_subsys(activeSs)) {
            cancel_pending_subsys(activeSs);
            ctx->Active = CONTROL_SUBSYS_NONE;
        }
    }

    if (ss->Phase == CONTROL_PHASE_EXECUTING) {
        status = HttpStatus_Conflict;
        goto done;
    }

    if (ss->Phase == CONTROL_PHASE_PENDING) {
        if (!force) {
            status = HttpStatus_Conflict;
            goto done;
        }

        ss->Desired = desired;
        ss->Token++;
        immediate = true;
        ctx->Active = subsystem;
        allowed = true;
        status = HttpStatus_Accepted;
        goto done;
    }

    if (ss->Phase == CONTROL_PHASE_FAULT) {
        if (!force) {
            status = HttpStatus_Conflict;
            goto done;
        }

        ss->Phase = CONTROL_PHASE_IDLE;
    }

    if (ss->Phase == CONTROL_PHASE_IDLE && ss->Current == desired) {
        status = HttpStatus_OK;
        allowed = false;
        goto done;
    }

    ss->Desired = desired;
    ss->Phase = CONTROL_PHASE_PENDING;
    ss->Token++;
    ctx->Active = subsystem;
    allowed = true;
    status = HttpStatus_Accepted;
    immediate = force ? true : false;

done:
    *httpStatus = status;
    *runImmediate = immediate;
    *outToken = ss ? ss->Token : 0;
    pthread_mutex_unlock(&ctx->Lock);
    return allowed;
}

bool control_begin_execute(ControlCtx *ctx, ControlSubsystem subsystem, unsigned long long token) {

    ControlSubsysState *ss = NULL;
    bool ok = false;

    if (!ctx) return false;

    pthread_mutex_lock(&ctx->Lock);

    if (ctx->Active != subsystem) {
        goto done;
    }

    ss = get_subsys(ctx, subsystem);
    if (!ss) {
        goto done;
    }

    if (ss->Phase != CONTROL_PHASE_PENDING) {
        goto done;
    }

    if (token != 0 && ss->Token != token) {
        goto done;
    }

    ss->Phase = CONTROL_PHASE_EXECUTING;
    ok = true;

done:
    pthread_mutex_unlock(&ctx->Lock);
    return ok;
}

void control_mark_fault(ControlCtx *ctx, ControlSubsystem subsystem) {

    ControlSubsysState *ss = NULL;

    if (!ctx) return;

    pthread_mutex_lock(&ctx->Lock);

    ss = get_subsys(ctx, subsystem);
    if (ss) {
        ss->Phase = CONTROL_PHASE_FAULT;
        ss->Token++;
    }
    if (ctx->Active == subsystem) {
        ctx->Active = CONTROL_SUBSYS_NONE;
    }

    pthread_mutex_unlock(&ctx->Lock);
}

void control_mark_done(ControlCtx *ctx,
                       ControlSubsystem subsystem,
                       ControlState finalState) {

    ControlSubsysState *ss = NULL;

    if (!ctx) return;

    pthread_mutex_lock(&ctx->Lock);

    ss = get_subsys(ctx, subsystem);
    if (ss) {
        ss->Current = finalState;
        ss->Phase = CONTROL_PHASE_IDLE;
        ss->Token++;
    }

    if (ctx->Active == subsystem) {
        ctx->Active = CONTROL_SUBSYS_NONE;
    }

    pthread_mutex_unlock(&ctx->Lock);
}

void control_mark_restart_done(ControlCtx *ctx) {
    control_mark_done(ctx, CONTROL_SUBSYS_RESTART, CONTROL_STATE_OFF);
}

int control_request(ControlCtx *ctx,
                    const char *nodeType,
                    ControlSubsystem subsystem,
                    ControlState desired,
                    bool force,
                    int *httpStatus) {

    bool allowed = false;
    bool runImmediate = false;
    unsigned long long token = 0;

    if (!ctx || !nodeType || !httpStatus) return STATUS_NOK;

    allowed = control_set_pending(ctx, subsystem, desired, force, httpStatus, &runImmediate, &token);
    if (!allowed) return STATUS_OK;

    (void)runImmediate;
    (void)token;

    return STATUS_OK;
}
