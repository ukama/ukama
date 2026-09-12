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

    if (!dst || size == 0) {
        return;
    }
    snprintf(dst, size, "%s", (src && *src) ? src : "none");
}

static void copy_optional(char *dst, size_t size, const char *src) {

    if (!dst || size == 0) {
        return;
    }
    snprintf(dst, size, "%s", src ? src : "");
}

static bool transition(LifecycleFsm *fsm,
                       LifecycleState state,
                       const char *reason,
                       int64_t epochSec) {

    if (!fsm) {
        return false;
    }

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

    if (!value || !state) {
        return false;
    }

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

    if (!value || !state) {
        return false;
    }

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

void lifecycle_fsm_init(LifecycleFsm *fsm, int64_t epochSec) {

    if (!fsm) {
        return;
    }

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

    if (!fsm || checkInTimeoutSec <= 0) {
        return false;
    }

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

    if (fsm->state == LIFECYCLE_STATE_CHECKING_IN) {
        return false;
    }

    if (fsm->state == LIFECYCLE_STATE_FAULTY &&
        fsm->fault != LIFECYCLE_FAULT_BOOT && fsm->gateOpen) {
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

static bool applications_ready(const StarterSnapshot *starter) {

    return starter->available && starter->aggregate == STARTER_AGGREGATE_READY;
}

static bool configuration_matches(const LifecycleFsm *fsm,
                                   const ConfigSnapshot *configuration) {

    return fsm->configGeneration == configuration->generation &&
        strcmp(fsm->requestId, configuration->requestId) == 0 &&
        strcmp(fsm->configMode, configuration->mode) == 0;
}

static bool start_configuration(LifecycleFsm *fsm,
                                 const ConfigSnapshot *configuration,
                                 int64_t epochSec) {

    fsm->fault = LIFECYCLE_FAULT_NONE;
    fsm->configurationSeen = true;
    fsm->configurationApplied = false;
    fsm->configGeneration = configuration->generation;
    copy_optional(fsm->requestId, sizeof(fsm->requestId),
                   configuration->requestId);
    copy_optional(fsm->configMode, sizeof(fsm->configMode),
                   configuration->mode);

    return transition(fsm, LIFECYCLE_STATE_CONFIGURING,
                       "configuration decision observed", epochSec);
}

static bool configuration_fault(LifecycleFsm *fsm, const char *reason,
                                 int64_t epochSec) {

    return enter_fault(fsm, LIFECYCLE_FAULT_CONFIGURATION,
                        LIFECYCLE_STATE_CONFIGURING, reason, epochSec);
}

static bool assignment_removed(const LifecycleFsm *fsm,
                                 const ConfigSnapshot *configuration) {

    return configuration->available &&
        configuration->phase == CONFIG_PHASE_AWAITING &&
        strcmp(configuration->mode, "NONE") == 0 &&
        configuration->generation > 0 &&
        configuration->generation >= fsm->configGeneration &&
        (!fsm->configurationSeen ||
         configuration->generation > fsm->configGeneration);
}

static bool clear_configuration(LifecycleFsm *fsm,
                                 const ConfigSnapshot *configuration,
                                 int64_t epochSec) {

    fsm->configurationSeen = false;
    fsm->configurationApplied = false;
    fsm->configGeneration = configuration->generation;
    fsm->confirmedGeneration = 0;
    fsm->requestId[0] = '\0';
    copy_optional(fsm->configMode, sizeof(fsm->configMode), "NONE");
    fsm->fault = LIFECYCLE_FAULT_NONE;
    fsm->faultReturnState = LIFECYCLE_STATE_READY;

    return transition(fsm, LIFECYCLE_STATE_READY,
                       "assignment removed; awaiting configuration", epochSec);
}

static bool tick_configuration(LifecycleFsm *fsm,
                                const StarterSnapshot *starter,
                                const ConfigSnapshot *configuration,
                                int64_t epochSec) {

    bool matches;
    bool repeatConfirmation;

    if (!configuration->available) {
        return false;
    }

    if (configuration->phase == CONFIG_PHASE_FAILED) {
        if (configuration->generation >= fsm->configGeneration &&
            configuration->requestId[0] != '\0') {
            fsm->configurationSeen = true;
            fsm->configGeneration = configuration->generation;
            copy_optional(fsm->requestId,
                          sizeof(fsm->requestId),
                          configuration->requestId);
            copy_optional(fsm->configMode,
                          sizeof(fsm->configMode),
                          configuration->mode);
        }
        return configuration_fault(fsm, "configd reports configuration failure",
                                    epochSec);
    }

    if (configuration->phase == CONFIG_PHASE_AWAITING) {
        if (assignment_removed(fsm, configuration)) {
            if (!applications_ready(starter)) {
                return false;
            }
            return clear_configuration(fsm, configuration, epochSec);
        }
        if (fsm->configurationSeen ||
            configuration->generation < fsm->configGeneration) {
            return configuration_fault(fsm, "configd lost its configuration record",
                                        epochSec);
        }
        return false;
    }

    if (configuration->generation < fsm->configGeneration) {
        return configuration_fault(fsm, "configd generation moved backwards",
                                    epochSec);
    }

    matches = configuration_matches(fsm, configuration);
    if (!matches && fsm->configGeneration != 0 &&
        configuration->generation == fsm->configGeneration) {
        return configuration_fault(fsm, "configd changed identity without a new generation",
                                    epochSec);
    }

    repeatConfirmation = fsm->state == LIFECYCLE_STATE_OPERATIONAL &&
        strcmp(configuration->mode, "NOCONFIG") == 0 &&
        strcmp(fsm->requestId, configuration->requestId) == 0;

    if (fsm->state == LIFECYCLE_STATE_READY ||
        (!matches && !repeatConfirmation)) {
        return start_configuration(fsm, configuration, epochSec);
    }

    if (configuration->phase != CONFIG_PHASE_APPLIED ||
        !applications_ready(starter)) {
        return false;
    }

    fsm->configurationApplied = true;
    fsm->configGeneration = configuration->generation;

    if (fsm->state == LIFECYCLE_STATE_CONFIGURING) {
        fsm->confirmedGeneration = configuration->generation;
        return transition(fsm, LIFECYCLE_STATE_OPERATIONAL,
                           "configuration complete; required applications ready",
                           epochSec);
    }

    if (fsm->state == LIFECYCLE_STATE_OPERATIONAL &&
        fsm->confirmedGeneration < configuration->generation) {
        fsm->confirmedGeneration = configuration->generation;
        fsm->sequence++;
        fsm->stateSince = epochSec;
        copy_text(fsm->reason, sizeof(fsm->reason),
                   "configuration confirmation repeated");
        return true;
    }

    return false;
}

static bool tick_faulty(LifecycleFsm *fsm,
                        const StarterSnapshot *starter,
                        const ConfigSnapshot *configuration,
                        int64_t epochSec) {

    if (fsm->fault == LIFECYCLE_FAULT_BOOT) {
        /* Starter must repeat its existing boot check-in after repair. */
        return false;
    }

    if (!fsm->gateOpen) {
        /* Services cannot be ready until starter can pass this gate. */
        fsm->fault = LIFECYCLE_FAULT_NONE;
        return transition(fsm, fsm->faultReturnState,
                           "starter reachable; resuming startup", epochSec);
    }

    if (!applications_ready(starter) || !configuration->available ||
        configuration->phase == CONFIG_PHASE_FAILED) {
        return false;
    }

    if (configuration->phase == CONFIG_PHASE_AWAITING) {
        if (assignment_removed(fsm, configuration)) {
            return clear_configuration(fsm, configuration, epochSec);
        }
        if (fsm->configurationSeen ||
            configuration->generation < fsm->configGeneration) {
            return false;
        }
        fsm->fault = LIFECYCLE_FAULT_NONE;
        return transition(fsm, LIFECYCLE_STATE_READY,
                           "applications recovered; awaiting assignment", epochSec);
    }

    if (configuration->generation < fsm->configGeneration) {
        return false;
    }

    if (configuration->generation == fsm->configGeneration &&
        !configuration_matches(fsm, configuration)) {
        return false;
    }

    return start_configuration(fsm, configuration, epochSec);
}

static bool service_unavailable(int64_t *sinceMs, int timeoutSec,
                                 int64_t nowMs) {

    if (*sinceMs == 0) {
        *sinceMs = nowMs;
        return false;
    }

    return nowMs - *sinceMs >= (int64_t)timeoutSec * 1000;
}

bool lifecycle_fsm_tick(LifecycleFsm *fsm,
                        const StarterSnapshot *starter,
                        const ConfigSnapshot *configuration,
                        int starterUnavailableTimeoutSec,
                        int configUnavailableTimeoutSec,
                        int64_t nowMs,
                        int64_t epochSec) {

    if (fsm == NULL || starter == NULL || configuration == NULL) {
        return false;
    }

    if (!starter->available) {
        if (service_unavailable(&fsm->starterUnavailableSinceMs,
                                 starterUnavailableTimeoutSec, nowMs) &&
            fsm->state != LIFECYCLE_STATE_FAULTY) {
            return enter_fault(fsm, LIFECYCLE_FAULT_STARTER, fsm->state,
                                "starter status unavailable", epochSec);
        }
        return false;
    }
    fsm->starterUnavailableSinceMs = 0;

    if (fsm->state == LIFECYCLE_STATE_FAULTY) {
        return tick_faulty(fsm, starter, configuration, epochSec);
    }

    if (fsm->state == LIFECYCLE_STATE_STARTING) {
        return false;
    }

    if (fsm->state == LIFECYCLE_STATE_CHECKING_IN) {
        if (nowMs >= fsm->checkInDeadlineMs) {
            fsm->gateOpen = true;
        }
        if (!fsm->gateOpen) {
            return false;
        }
    }

    if (starter->aggregate == STARTER_AGGREGATE_FAULTY) {
        return enter_fault(fsm, LIFECYCLE_FAULT_STARTER, fsm->state,
                            "required application readiness failed", epochSec);
    }

    if (fsm->state == LIFECYCLE_STATE_CHECKING_IN) {
        if (applications_ready(starter)) {
            return transition(fsm, LIFECYCLE_STATE_READY,
                               "required applications ready; awaiting assignment",
                               epochSec);
        }
        return false;
    }

    if (!configuration->available) {
        if (service_unavailable(&fsm->configUnavailableSinceMs,
                                 configUnavailableTimeoutSec, nowMs)) {
            return configuration_fault(fsm, "configd status unavailable", epochSec);
        }
        return false;
    }
    fsm->configUnavailableSinceMs = 0;

    return tick_configuration(fsm, starter, configuration, epochSec);
}

HttpStatus lifecycle_fsm_gate_status(const LifecycleFsm *fsm,
                                     int64_t nowMs,
                                     int *remainingSec) {

    int64_t remainingMs;

    if (remainingSec) {
        *remainingSec = 0;
    }
    if (!fsm) {
        return HttpStatus_ServiceUnavailable;
    }

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
    if (remainingMs < 0) {
        remainingMs = 0;
    }

    if (remainingSec) {
        *remainingSec = (int)((remainingMs + 999) / 1000);
    }

    return HttpStatus_Accepted;
}
