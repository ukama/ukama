/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONTROL_H_
#define CONTROL_H_

#include <pthread.h>

#include "usys_types.h"

typedef enum {
    CONTROL_SUBSYS_NONE = 0,
    CONTROL_SUBSYS_SERVICE,
    CONTROL_SUBSYS_RADIO,
    CONTROL_SUBSYS_RESTART,
} ControlSubsystem;

typedef enum {
    CONTROL_PHASE_IDLE = 0,
    CONTROL_PHASE_PENDING,
    CONTROL_PHASE_EXECUTING,
    CONTROL_PHASE_FAULT,
} ControlPhase;

typedef enum {
    CONTROL_STATE_OFF = 0,
    CONTROL_STATE_ON,
} ControlState;

typedef struct {
    ControlPhase       Phase;
    ControlState       Current;
    ControlState       Desired;
    unsigned long long Token;
} ControlSubsysState;

typedef struct {
    pthread_mutex_t Lock;

    ControlSubsystem   Active;
    ControlSubsysState Service;
    ControlSubsysState Radio;
    ControlSubsysState Restart;
} ControlCtx;

ControlCtx *control_create(void);
void control_destroy(ControlCtx *ctx);

bool control_is_busy(ControlCtx *ctx);

int control_get_public_state(ControlCtx *ctx,
                             const char *nodeType,
                             char *outState,
                             size_t outStateSize);

int control_request(ControlCtx *ctx,
                    const char *nodeType,
                    ControlSubsystem subsystem,
                    ControlState desired,
                    bool force,
                    int *httpStatus);

void control_mark_fault(ControlCtx *ctx, ControlSubsystem subsystem);

void control_mark_done(ControlCtx *ctx,
                       ControlSubsystem subsystem,
                       ControlState finalState);

void control_mark_restart_done(ControlCtx *ctx);

bool control_begin_execute(ControlCtx *ctx,
                           ControlSubsystem subsystem,
                           unsigned long long token);

bool control_set_pending(ControlCtx *ctx,
                         ControlSubsystem subsystem,
                         ControlState desired,
                         bool force,
                         int *httpStatus,
                         bool *runImmediate,
                         unsigned long long *outToken);

#endif
