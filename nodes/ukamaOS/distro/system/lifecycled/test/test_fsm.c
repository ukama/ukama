/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "fsm.h"
#include "state_store.h"

static int gTests;
static int gFailures;

#define EXPECT_TRUE(expression)                                          \
    do {                                                                 \
        gTests++;                                                        \
        if (!(expression)) {                                             \
            fprintf(stderr,                                              \
                    "%s:%d expected true: %s\n",                       \
                    __FILE__,                                            \
                    __LINE__,                                            \
                    #expression);                                        \
            gFailures++;                                                 \
        }                                                                \
    } while (0)

#define EXPECT_EQ(expected, actual)                                      \
    do {                                                                 \
        long long expectedValue = (long long)(expected);                 \
        long long actualValue = (long long)(actual);                     \
        gTests++;                                                        \
        if (expectedValue != actualValue) {                              \
            fprintf(stderr,                                              \
                    "%s:%d expected %lld, got %lld\n",                 \
                    __FILE__,                                            \
                    __LINE__,                                            \
                    expectedValue,                                       \
                    actualValue);                                        \
            gFailures++;                                                 \
        }                                                                \
    } while (0)

static StarterSnapshot starter_snapshot(StarterAggregateState aggregate,
                                        ConfigPhase configPhase) {

    StarterSnapshot starter;

    memset(&starter, 0, sizeof(starter));
    starter.available = true;
    starter.aggregate = aggregate;
    starter.configPhase = configPhase;
    snprintf(starter.aggregateReason,
             sizeof(starter.aggregateReason),
             "%s",
             starter_aggregate_str(aggregate));
    return starter;
}

static void reach_ready(LifecycleFsm *fsm) {

    StarterSnapshot starter;

    lifecycle_fsm_init(fsm, 10);
    EXPECT_TRUE(lifecycle_fsm_begin_check_in(fsm, true, 2, 1000, 11));

    starter = starter_snapshot(STARTER_AGGREGATE_READY,
                               CONFIG_PHASE_AWAITING);
    EXPECT_TRUE(!lifecycle_fsm_tick(fsm, &starter, 10, 2, 2999, 12));
    EXPECT_EQ(LIFECYCLE_STATE_CHECKING_IN, fsm->state);

    EXPECT_TRUE(lifecycle_fsm_tick(fsm, &starter, 10, 2, 3000, 13));
    EXPECT_EQ(LIFECYCLE_STATE_READY, fsm->state);
    EXPECT_TRUE(fsm->gateOpen);
}

static void test_initial_state(void) {

    LifecycleFsm fsm;

    lifecycle_fsm_init(&fsm, 100);
    EXPECT_EQ(LIFECYCLE_STATE_STARTING, fsm.state);
    EXPECT_EQ(1, fsm.sequence);
    EXPECT_EQ(100, fsm.stateSince);
}

static void test_ready_timeout_reaches_operational(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;

    reach_ready(&fsm);
    starter = starter_snapshot(STARTER_AGGREGATE_READY,
                               CONFIG_PHASE_AWAITING);

    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm,
                                    &starter,
                                    10,
                                    2,
                                    4999,
                                    14));
    EXPECT_EQ(LIFECYCLE_STATE_READY, fsm.state);

    EXPECT_TRUE(lifecycle_fsm_tick(&fsm,
                                   &starter,
                                   10,
                                   2,
                                   5000,
                                   15));
    EXPECT_EQ(LIFECYCLE_STATE_OPERATIONAL, fsm.state);
}

static void test_no_config_reaches_operational(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;
    LifecycleConfigureResult result;

    reach_ready(&fsm);
    starter = starter_snapshot(STARTER_AGGREGATE_READY,
                               CONFIG_PHASE_AWAITING);

    result = lifecycle_fsm_configure(&fsm,
                                     "request-1",
                                     "assignment-1",
                                     2,
                                     4000,
                                     20);
    EXPECT_EQ(LIFECYCLE_CONFIGURE_ACCEPTED, result);
    EXPECT_EQ(LIFECYCLE_STATE_CONFIGURING, fsm.state);

    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm, &starter, 10, 2, 5999, 21));
    EXPECT_EQ(LIFECYCLE_STATE_CONFIGURING, fsm.state);

    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 10, 2, 6000, 22));
    EXPECT_EQ(LIFECYCLE_STATE_OPERATIONAL, fsm.state);
}

static void test_config_apply_reaches_operational(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;

    reach_ready(&fsm);
    EXPECT_EQ(LIFECYCLE_CONFIGURE_ACCEPTED,
              lifecycle_fsm_configure(&fsm,
                                      "request-2",
                                      "assignment-2",
                                      60,
                                      4000,
                                      20));

    starter = starter_snapshot(STARTER_AGGREGATE_PENDING,
                               CONFIG_PHASE_IN_PROGRESS);
    snprintf(starter.configRequestId,
             sizeof(starter.configRequestId),
             "request-2");
    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm, &starter, 10, 2, 4100, 21));
    EXPECT_TRUE(fsm.configurationSeen);

    starter.aggregate = STARTER_AGGREGATE_READY;
    starter.configPhase = CONFIG_PHASE_APPLIED;
    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 10, 2, 4200, 22));
    EXPECT_EQ(LIFECYCLE_STATE_OPERATIONAL, fsm.state);
    EXPECT_TRUE(fsm.configurationApplied);
}

static void test_config_in_progress_timeout_faults(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;

    reach_ready(&fsm);
    EXPECT_EQ(LIFECYCLE_CONFIGURE_ACCEPTED,
              lifecycle_fsm_configure(&fsm,
                                      "request-timeout",
                                      "assignment-timeout",
                                      2,
                                      4000,
                                      20));

    starter = starter_snapshot(STARTER_AGGREGATE_PENDING,
                               CONFIG_PHASE_IN_PROGRESS);
    snprintf(starter.configRequestId,
             sizeof(starter.configRequestId),
             "request-timeout");

    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm,
                                    &starter,
                                    10,
                                    2,
                                    4100,
                                    21));
    EXPECT_TRUE(fsm.configurationSeen);
    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm,
                                    &starter,
                                    10,
                                    2,
                                    5999,
                                    22));
    EXPECT_TRUE(lifecycle_fsm_tick(&fsm,
                                   &starter,
                                   10,
                                   2,
                                   6000,
                                   23));
    EXPECT_EQ(LIFECYCLE_STATE_FAULTY, fsm.state);
    EXPECT_EQ(LIFECYCLE_FAULT_CONFIGURATION, fsm.fault);
}

static void test_fast_matching_apply_is_not_missed(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;

    reach_ready(&fsm);
    EXPECT_EQ(LIFECYCLE_CONFIGURE_ACCEPTED,
              lifecycle_fsm_configure(&fsm,
                                      "request-fast",
                                      "assignment-fast",
                                      60,
                                      4000,
                                      20));

    starter = starter_snapshot(STARTER_AGGREGATE_READY,
                               CONFIG_PHASE_APPLIED);
    snprintf(starter.configRequestId,
             sizeof(starter.configRequestId),
             "request-fast");

    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 10, 2, 4100, 21));
    EXPECT_EQ(LIFECYCLE_STATE_OPERATIONAL, fsm.state);
}

static void test_config_failure_latches_fault(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;

    reach_ready(&fsm);
    EXPECT_EQ(LIFECYCLE_CONFIGURE_ACCEPTED,
              lifecycle_fsm_configure(&fsm,
                                      "request-3",
                                      "assignment-3",
                                      60,
                                      4000,
                                      20));

    starter = starter_snapshot(STARTER_AGGREGATE_FAULTY,
                               CONFIG_PHASE_FAILED);
    snprintf(starter.configReason,
             sizeof(starter.configReason),
             "configuration_failed");
    snprintf(starter.configRequestId,
             sizeof(starter.configRequestId),
             "request-3");

    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 10, 2, 4100, 21));
    EXPECT_EQ(LIFECYCLE_STATE_FAULTY, fsm.state);
    EXPECT_EQ(LIFECYCLE_FAULT_CONFIGURATION, fsm.fault);

    starter.aggregate = STARTER_AGGREGATE_READY;
    starter.configPhase = CONFIG_PHASE_APPLIED;
    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm, &starter, 10, 2, 4200, 22));
    EXPECT_EQ(LIFECYCLE_STATE_FAULTY, fsm.state);
}

static void test_stale_config_failure_is_ignored(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;

    reach_ready(&fsm);
    EXPECT_EQ(LIFECYCLE_CONFIGURE_ACCEPTED,
              lifecycle_fsm_configure(&fsm,
                                      "request-new",
                                      "assignment-new",
                                      2,
                                      4000,
                                      20));

    starter = starter_snapshot(STARTER_AGGREGATE_READY,
                               CONFIG_PHASE_FAILED);
    snprintf(starter.configReason,
             sizeof(starter.configReason),
             "configuration_failed");

    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm, &starter, 10, 2, 4100, 21));
    EXPECT_EQ(LIFECYCLE_STATE_CONFIGURING, fsm.state);
    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 10, 2, 6000, 22));
    EXPECT_EQ(LIFECYCLE_STATE_OPERATIONAL, fsm.state);
}

static void test_configure_idempotency(void) {

    LifecycleFsm fsm;
    int64_t deadline;

    reach_ready(&fsm);
    EXPECT_EQ(LIFECYCLE_CONFIGURE_ACCEPTED,
              lifecycle_fsm_configure(&fsm,
                                      "request-4",
                                      "assignment-4",
                                      2,
                                      4000,
                                      20));
    deadline = fsm.configDeadlineMs;

    EXPECT_EQ(LIFECYCLE_CONFIGURE_DUPLICATE,
              lifecycle_fsm_configure(&fsm,
                                      "request-4",
                                      "assignment-4",
                                      2,
                                      5000,
                                      21));
    EXPECT_EQ(deadline, fsm.configDeadlineMs);

    EXPECT_EQ(LIFECYCLE_CONFIGURE_BUSY,
              lifecycle_fsm_configure(&fsm,
                                      "request-other",
                                      "assignment-4",
                                      2,
                                      5000,
                                      21));
}

static void test_starter_fault_recovers(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;

    reach_ready(&fsm);
    EXPECT_EQ(LIFECYCLE_CONFIGURE_ACCEPTED,
              lifecycle_fsm_configure(&fsm,
                                      "request-5",
                                      "assignment-5",
                                      1,
                                      4000,
                                      20));
    starter = starter_snapshot(STARTER_AGGREGATE_READY,
                               CONFIG_PHASE_AWAITING);
    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 10, 2, 5000, 21));
    EXPECT_EQ(LIFECYCLE_STATE_OPERATIONAL, fsm.state);

    starter.aggregate = STARTER_AGGREGATE_FAULTY;
    snprintf(starter.aggregateReason,
             sizeof(starter.aggregateReason),
             "required app timeout");
    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 10, 2, 5100, 22));
    EXPECT_EQ(LIFECYCLE_STATE_FAULTY, fsm.state);

    starter.aggregate = STARTER_AGGREGATE_READY;
    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 10, 2, 5200, 23));
    EXPECT_EQ(LIFECYCLE_STATE_OPERATIONAL, fsm.state);
}

static void test_starter_unavailable_timeout(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;

    reach_ready(&fsm);
    memset(&starter, 0, sizeof(starter));

    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm, &starter, 2, 2, 4000, 20));
    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm, &starter, 2, 2, 5999, 21));
    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 2, 2, 6000, 22));
    EXPECT_EQ(LIFECYCLE_STATE_FAULTY, fsm.state);
}

static void test_starting_recovers_to_starting(void) {

    LifecycleFsm fsm;
    StarterSnapshot starter;

    lifecycle_fsm_init(&fsm, 10);
    memset(&starter, 0, sizeof(starter));

    EXPECT_TRUE(!lifecycle_fsm_tick(&fsm, &starter, 1, 2, 1000, 11));
    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 1, 2, 2000, 12));
    EXPECT_EQ(LIFECYCLE_STATE_FAULTY, fsm.state);

    starter = starter_snapshot(STARTER_AGGREGATE_READY,
                               CONFIG_PHASE_AWAITING);
    EXPECT_TRUE(lifecycle_fsm_tick(&fsm, &starter, 1, 2, 2100, 13));
    EXPECT_EQ(LIFECYCLE_STATE_STARTING, fsm.state);
}

static void test_state_store_is_boot_scoped(void) {

    LifecycleFsm fsm;
    LifecycleFsm loaded;
    char path[] = "/tmp/lifecycled-state-XXXXXX";
    int fd;

    fd = mkstemp(path);
    EXPECT_TRUE(fd >= 0);
    if (fd < 0) return;
    close(fd);

    reach_ready(&fsm);
    EXPECT_TRUE(state_store_save(path, "boot-one", &fsm));
    memset(&loaded, 0, sizeof(loaded));
    EXPECT_TRUE(state_store_load(path, "boot-one", &loaded));
    EXPECT_EQ(fsm.state, loaded.state);
    EXPECT_EQ(fsm.sequence, loaded.sequence);

    memset(&loaded, 0, sizeof(loaded));
    EXPECT_TRUE(!state_store_load(path, "boot-two", &loaded));
    unlink(path);
}

int main(void) {

    test_initial_state();
    test_ready_timeout_reaches_operational();
    test_no_config_reaches_operational();
    test_config_apply_reaches_operational();
    test_config_in_progress_timeout_faults();
    test_fast_matching_apply_is_not_missed();
    test_config_failure_latches_fault();
    test_stale_config_failure_is_ignored();
    test_configure_idempotency();
    test_starter_fault_recovers();
    test_starter_unavailable_timeout();
    test_starting_recovers_to_starting();
    test_state_store_is_boot_scoped();

    if (gFailures) {
        fprintf(stderr,
                "FAILED: %d of %d assertions failed\n",
                gFailures,
                gTests);
        return 1;
    }

    printf("PASS: %d lifecycle FSM assertions\n", gTests);
    return 0;
}
