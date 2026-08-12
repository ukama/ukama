/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "http_status.h"

#define LIFECYCLE_REASON_LEN 192
#define LIFECYCLE_ID_LEN      96

typedef enum {
    LIFECYCLE_STATE_STARTING = 0,
    LIFECYCLE_STATE_CHECKING_IN,
    LIFECYCLE_STATE_READY,
    LIFECYCLE_STATE_CONFIGURING,
    LIFECYCLE_STATE_OPERATIONAL,
    LIFECYCLE_STATE_FAULTY
} LifecycleState;

typedef enum {
    STARTER_AGGREGATE_UNKNOWN = 0,
    STARTER_AGGREGATE_PENDING,
    STARTER_AGGREGATE_READY,
    STARTER_AGGREGATE_FAULTY
} StarterAggregateState;

typedef enum {
    CONFIG_PHASE_ABSENT = 0,
    CONFIG_PHASE_AWAITING,
    CONFIG_PHASE_IN_PROGRESS,
    CONFIG_PHASE_APPLIED,
    CONFIG_PHASE_FAILED,
    CONFIG_PHASE_UNKNOWN
} ConfigPhase;

typedef enum {
    LIFECYCLE_FAULT_NONE = 0,
    LIFECYCLE_FAULT_BOOT,
    LIFECYCLE_FAULT_STARTER,
    LIFECYCLE_FAULT_CONFIGURATION
} LifecycleFault;

typedef enum {
    LIFECYCLE_CONFIGURE_ACCEPTED = 0,
    LIFECYCLE_CONFIGURE_DUPLICATE,
    LIFECYCLE_CONFIGURE_BUSY,
    LIFECYCLE_CONFIGURE_INVALID_STATE,
    LIFECYCLE_CONFIGURE_INVALID_REQUEST
} LifecycleConfigureResult;

typedef struct {
    bool available;
    StarterAggregateState aggregate;
    ConfigPhase configPhase;
    char aggregateReason[LIFECYCLE_REASON_LEN];
    char configReason[LIFECYCLE_REASON_LEN];
    char configRequestId[LIFECYCLE_ID_LEN];
} StarterSnapshot;

typedef struct {
    LifecycleState state;
    LifecycleState faultReturnState;
    LifecycleFault fault;

    uint64_t sequence;
    int64_t stateSince;
    int64_t checkInDeadlineMs;
    int64_t configDeadlineMs;
    int64_t starterUnavailableSinceMs;

    bool gateOpen;
    bool configurationSeen;
    bool configurationApplied;

    char requestId[LIFECYCLE_ID_LEN];
    char assignmentId[LIFECYCLE_ID_LEN];
    char reason[LIFECYCLE_REASON_LEN];
} LifecycleFsm;

const char *lifecycle_state_str(LifecycleState state);
const char *starter_aggregate_str(StarterAggregateState state);
const char *config_phase_str(ConfigPhase phase);

bool lifecycle_state_parse(const char *value, LifecycleState *state);
bool starter_aggregate_parse(const char *value,
                             StarterAggregateState *state);
ConfigPhase config_phase_from_reason(const char *reason);

void lifecycle_fsm_init(LifecycleFsm *fsm, int64_t epochSec);

bool lifecycle_fsm_begin_check_in(LifecycleFsm *fsm,
                                  bool bootHealthy,
                                  int checkInTimeoutSec,
                                  int64_t nowMs,
                                  int64_t epochSec);

LifecycleConfigureResult lifecycle_fsm_configure(
    LifecycleFsm *fsm,
    const char *requestId,
    const char *assignmentId,
    int configTimeoutSec,
    int64_t nowMs,
    int64_t epochSec);

bool lifecycle_fsm_tick(LifecycleFsm *fsm,
                        const StarterSnapshot *starter,
                        int starterUnavailableTimeoutSec,
                        int64_t nowMs,
                        int64_t epochSec);

HttpStatus lifecycle_fsm_gate_status(const LifecycleFsm *fsm,
                                     int64_t nowMs,
                                     int *remainingSec);
