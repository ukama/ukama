/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "fsm.h"

static void copy_text(char *dst, size_t size, const char *src) {

    if (!dst || size == 0) return;
    snprintf(dst, size, "%s", (src && *src) ? src : "none");
}

static void copy_optional(char *dst, size_t size, const char *src) {

    if (!dst || size == 0) return;
    snprintf(dst, size, "%s", src ? src : "");
}

static bool request_matches(const LifecycleFsm *fsm,
                            const StarterSnapshot *starter) {

    if (!fsm || !starter) return false;
    if (starter->configRequestId[0] == '\0') return true;

    return strcmp(fsm->requestId, starter->configRequestId) == 0;
}

static bool transition(LifecycleFsm *fsm,
                       LifecycleState state,
                       const char *reason,
                       int64_t epochSec) {

    if (!fsm) return false;

    if (fsm->state == state) {
        copy_text(fsm->reason, sizeof(fsm->reason), reason);
        return false;
    }

    fsm->state = state;
    fsm->sequence++;
    fsm->stateSince = epochSec;
    copy_text(fsm->reason, sizeof(fsm->reason), reason);
    return true;
}

static bool enter_fault(LifecycleFsm *fsm,
                        LifecycleFault fault,
                        LifecycleState returnState,
                        const char *reason,
                        int64_t epochSec) {

    fsm->fault = fault;
    fsm->faultReturnState = returnState;

    return transition(fsm,
                      LIFECYCLE_STATE_FAULTY,
                      reason,
                      epochSec);
}

const char *lifecycle_state_str(LifecycleState state) {

    switch (state) {
    case LIFECYCLE_STATE_STARTING:     return "STARTING";
    case LIFECYCLE_STATE_CHECKING_IN:  return "CHECKING_IN";
    case LIFECYCLE_STATE_READY:        return "READY";
    case LIFECYCLE_STATE_CONFIGURING:  return "CONFIGURING";
    case LIFECYCLE_STATE_OPERATIONAL:  return "OPERATIONAL";
    case LIFECYCLE_STATE_FAULTY:       return "FAULTY";
    default:                           return "UNKNOWN";
    }
}

const char *starter_aggregate_str(StarterAggregateState state) {

    switch (state) {
    case STARTER_AGGREGATE_PENDING: return "pending";
    case STARTER_AGGREGATE_READY:   return "ready";
    case STARTER_AGGREGATE_FAULTY:  return "faulty";
    default:                        return "unknown";
    }
}

const char *config_phase_str(ConfigPhase phase) {

    switch (phase) {
    case CONFIG_PHASE_ABSENT:      return "absent";
    case CONFIG_PHASE_AWAITING:    return "awaiting";
    case CONFIG_PHASE_IN_PROGRESS: return "in_progress";
    case CONFIG_PHASE_APPLIED:     return "applied";
    case CONFIG_PHASE_FAILED:      return "failed";
    default:                       return "unknown";
    }
}

bool lifecycle_state_parse(const char *value, LifecycleState *state) {

    int i;

    if (!value || !state) return false;

    for (i = LIFECYCLE_STATE_STARTING;
         i <= LIFECYCLE_STATE_FAULTY;
         i++) {
        if (strcmp(value, lifecycle_state_str((LifecycleState)i)) == 0) {
            *state = (LifecycleState)i;
            return true;
        }
    }

    return false;
}

bool starter_aggregate_parse(const char *value,
                             StarterAggregateState *state) {

    if (!value || !state) return false;

    if (strcmp(value, "pending") == 0) {
        *state = STARTER_AGGREGATE_PENDING;
    } else if (strcmp(value, "ready") == 0) {
        *state = STARTER_AGGREGATE_READY;
    } else if (strcmp(value, "faulty") == 0) {
        *state = STARTER_AGGREGATE_FAULTY;
    } else {
        return false;
    }

    return true;
}

ConfigPhase config_phase_from_reason(const char *reason) {

    if (!reason || !*reason) return CONFIG_PHASE_UNKNOWN;

    if (strncmp(reason,
                "awaiting_configuration",
                strlen("awaiting_configuration")) == 0) {
        return CONFIG_PHASE_AWAITING;
    }

    if (strncmp(reason,
                "configuration_in_progress",
                strlen("configuration_in_progress")) == 0) {
        return CONFIG_PHASE_IN_PROGRESS;
    }

    if (strncmp(reason,
                "configuration_applied",
                strlen("configuration_applied")) == 0) {
        return CONFIG_PHASE_APPLIED;
    }

    if (strncmp(reason,
                "configuration_failed",
                strlen("configuration_failed")) == 0) {
        return CONFIG_PHASE_FAILED;
    }

    return CONFIG_PHASE_UNKNOWN;
}

void lifecycle_fsm_init(LifecycleFsm *fsm, int64_t epochSec) {

    if (!fsm) return;

    memset(fsm, 0, sizeof(*fsm));
    fsm->state = LIFECYCLE_STATE_STARTING;
    fsm->faultReturnState = LIFECYCLE_STATE_STARTING;
    fsm->fault = LIFECYCLE_FAULT_NONE;
    fsm->sequence = 1;
    fsm->stateSince = epochSec;
    copy_text(fsm->reason,
              sizeof(fsm->reason),
              "lifecycle manager started");
}

bool lifecycle_fsm_begin_check_in(LifecycleFsm *fsm,
                                  bool bootHealthy,
                                  int checkInTimeoutSec,
                                  int64_t nowMs,
                                  int64_t epochSec) {

    if (!fsm || checkInTimeoutSec <= 0) return false;

    if (fsm->state == LIFECYCLE_STATE_READY ||
        fsm->state == LIFECYCLE_STATE_CONFIGURING ||
        fsm->state == LIFECYCLE_STATE_OPERATIONAL) {
        return false;
    }

    if (!bootHealthy) {
        return enter_fault(fsm,
                           LIFECYCLE_FAULT_BOOT,
                           LIFECYCLE_STATE_CHECKING_IN,
                           "starter boot applications are degraded",
                           epochSec);
    }

    if (fsm->state == LIFECYCLE_STATE_CHECKING_IN) return false;

    if (fsm->state == LIFECYCLE_STATE_FAULTY &&
        fsm->fault != LIFECYCLE_FAULT_BOOT) {
        return false;
    }

    fsm->fault = LIFECYCLE_FAULT_NONE;
    fsm->gateOpen = false;
    fsm->checkInDeadlineMs = nowMs +
        ((int64_t)checkInTimeoutSec * 1000);

    return transition(fsm,
                      LIFECYCLE_STATE_CHECKING_IN,
                      "waiting for backend check-in window",
                      epochSec);
}

LifecycleConfigureResult lifecycle_fsm_configure(
    LifecycleFsm *fsm,
    const char *requestId,
    const char *assignmentId,
    int configTimeoutSec,
    int64_t nowMs,
    int64_t epochSec) {

    if (!fsm || !requestId || !*requestId || configTimeoutSec <= 0) {
        return LIFECYCLE_CONFIGURE_INVALID_REQUEST;
    }

    if (fsm->requestId[0] != '\0' &&
        strcmp(fsm->requestId, requestId) == 0) {
        return LIFECYCLE_CONFIGURE_DUPLICATE;
    }

    if (fsm->state == LIFECYCLE_STATE_CONFIGURING) {
        return LIFECYCLE_CONFIGURE_BUSY;
    }

    if (fsm->state != LIFECYCLE_STATE_READY &&
        fsm->state != LIFECYCLE_STATE_OPERATIONAL &&
        fsm->state != LIFECYCLE_STATE_FAULTY) {
        return LIFECYCLE_CONFIGURE_INVALID_STATE;
    }

    fsm->fault                = LIFECYCLE_FAULT_NONE;
    fsm->configurationSeen    = false;
    fsm->configurationApplied = false;
    fsm->configDeadlineMs     = nowMs +
        ((int64_t)configTimeoutSec * 1000);

    copy_optional(fsm->requestId, sizeof(fsm->requestId), requestId);
    copy_optional(fsm->assignmentId,
                  sizeof(fsm->assignmentId),
                  assignmentId);

    transition(fsm,
               LIFECYCLE_STATE_CONFIGURING,
               "backend configuration command accepted",
               epochSec);

    return LIFECYCLE_CONFIGURE_ACCEPTED;
}

static bool handle_starter_fault(LifecycleFsm *fsm,
                                 const StarterSnapshot *starter,
                                 int64_t epochSec) {

    char reason[LIFECYCLE_REASON_LEN];

    if (starter->aggregate != STARTER_AGGREGATE_FAULTY) {
        return false;
    }

    snprintf(reason,
             sizeof(reason),
             "starter readiness faulty: %.150s",
             starter->aggregateReason[0] ?
             starter->aggregateReason : "unknown reason");

    return enter_fault(fsm,
                       LIFECYCLE_FAULT_STARTER,
                       fsm->state,
                       reason,
                       epochSec);
}

static bool handle_unavailable_starter(
    LifecycleFsm *fsm,
    const StarterSnapshot *starter,
    int starterUnavailableTimeoutSec,
    int64_t nowMs,
    int64_t epochSec) {

    if (starter->available) {
        fsm->starterUnavailableSinceMs = 0;
        return false;
    }

    if (fsm->starterUnavailableSinceMs == 0) {
        fsm->starterUnavailableSinceMs = nowMs;
        return false;
    }

    if (nowMs - fsm->starterUnavailableSinceMs <
        ((int64_t)starterUnavailableTimeoutSec * 1000)) {
        return false;
    }

    if (fsm->state == LIFECYCLE_STATE_FAULTY) {
        return false;
    }

    return enter_fault(fsm,
                       LIFECYCLE_FAULT_STARTER,
                       fsm->state,
                       "starter status unavailable",
                       epochSec);
}

static bool tick_checking_in(LifecycleFsm *fsm,
                             const StarterSnapshot *starter,
                             int64_t nowMs,
                             int64_t epochSec) {

    if (!fsm->gateOpen && nowMs >= fsm->checkInDeadlineMs) {
        fsm->gateOpen = true;
    }

    if (!fsm->gateOpen) return false;

    if (handle_starter_fault(fsm, starter, epochSec)) {
        return true;
    }

    if (starter->aggregate == STARTER_AGGREGATE_READY) {
        return transition(fsm,
                          LIFECYCLE_STATE_READY,
                          "check-in complete; required applications ready",
                          epochSec);
    }

    return false;
}

static bool tick_configuring(LifecycleFsm *fsm,
                             const StarterSnapshot *starter,
                             int64_t nowMs,
                             int64_t epochSec) {

    bool matchingRequest;

    matchingRequest = request_matches(fsm, starter);

    if (starter->configPhase == CONFIG_PHASE_IN_PROGRESS &&
        matchingRequest) {
        fsm->configurationSeen = true;
        fsm->configurationApplied = false;
    }

    if (starter->configPhase == CONFIG_PHASE_APPLIED &&
        matchingRequest &&
        (fsm->configurationSeen ||
         starter->configRequestId[0] != '\0')) {
        fsm->configurationSeen = true;
        fsm->configurationApplied = true;
    }

    if (starter->configPhase == CONFIG_PHASE_FAILED &&
        matchingRequest &&
        (fsm->configurationSeen ||
         starter->configRequestId[0] != '\0')) {
        return enter_fault(fsm,
                           LIFECYCLE_FAULT_CONFIGURATION,
                           LIFECYCLE_STATE_CONFIGURING,
                           starter->configReason[0] ?
                           starter->configReason :
                           "configuration failed",
                           epochSec);
    }

    if (handle_starter_fault(fsm, starter, epochSec)) {
        return true;
    }

    if (fsm->configurationSeen &&
        fsm->configurationApplied &&
        starter->aggregate == STARTER_AGGREGATE_READY) {
        fsm->configDeadlineMs = 0;
        return transition(fsm,
                          LIFECYCLE_STATE_OPERATIONAL,
                          "configuration applied; applications ready",
                          epochSec);
    }

    if (!fsm->configurationSeen &&
        nowMs >= fsm->configDeadlineMs &&
        starter->aggregate == STARTER_AGGREGATE_READY) {
        fsm->configDeadlineMs = 0;
        return transition(fsm,
                          LIFECYCLE_STATE_OPERATIONAL,
                          "configuration window completed; no config received",
                          epochSec);
    }

    return false;
}

static bool tick_faulty(LifecycleFsm *fsm,
                        const StarterSnapshot *starter,
                        int64_t epochSec) {

    LifecycleState target;

    if (fsm->fault != LIFECYCLE_FAULT_STARTER ||
        !starter->available ||
        starter->aggregate != STARTER_AGGREGATE_READY) {
        return false;
    }

    target = fsm->faultReturnState;
    fsm->fault = LIFECYCLE_FAULT_NONE;

    if (target == LIFECYCLE_STATE_CONFIGURING) {
        return transition(fsm,
                          target,
                          "starter readiness recovered during configuration",
                          epochSec);
    }

    return transition(fsm,
                      target,
                      "starter readiness recovered",
                      epochSec);
}

bool lifecycle_fsm_tick(LifecycleFsm *fsm,
                        const StarterSnapshot *starter,
                        int starterUnavailableTimeoutSec,
                        int64_t nowMs,
                        int64_t epochSec) {

    if (!fsm || !starter || starterUnavailableTimeoutSec <= 0) {
        return false;
    }

    if (handle_unavailable_starter(fsm,
                                   starter,
                                   starterUnavailableTimeoutSec,
                                   nowMs,
                                   epochSec)) {
        return true;
    }

    if (!starter->available) return false;

    switch (fsm->state) {
    case LIFECYCLE_STATE_CHECKING_IN:
        return tick_checking_in(fsm, starter, nowMs, epochSec);

    case LIFECYCLE_STATE_READY:
    case LIFECYCLE_STATE_OPERATIONAL:
        return handle_starter_fault(fsm, starter, epochSec);

    case LIFECYCLE_STATE_CONFIGURING:
        return tick_configuring(fsm, starter, nowMs, epochSec);

    case LIFECYCLE_STATE_FAULTY:
        return tick_faulty(fsm, starter, epochSec);

    case LIFECYCLE_STATE_STARTING:
    default:
        return false;
    }
}

HttpStatus lifecycle_fsm_gate_status(const LifecycleFsm *fsm,
                                     int64_t nowMs,
                                     int *remainingSec) {

    int64_t remainingMs;

    if (remainingSec) *remainingSec = 0;
    if (!fsm) return HttpStatus_ServiceUnavailable;

    if (fsm->state == LIFECYCLE_STATE_FAULTY) {
        return HttpStatus_ServiceUnavailable;
    }

    if (fsm->state == LIFECYCLE_STATE_STARTING) {
        return HttpStatus_Accepted;
    }

    if (fsm->gateOpen ||
        fsm->state == LIFECYCLE_STATE_READY ||
        fsm->state == LIFECYCLE_STATE_CONFIGURING ||
        fsm->state == LIFECYCLE_STATE_OPERATIONAL) {
        return HttpStatus_OK;
    }

    remainingMs = fsm->checkInDeadlineMs - nowMs;
    if (remainingMs < 0) remainingMs = 0;

    if (remainingSec) {
        *remainingSec = (int)((remainingMs + 999) / 1000);
    }

    return HttpStatus_Accepted;
}
