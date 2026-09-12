/* SPDX-License-Identifier: MPL-2.0 */
#include <assert.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "fsm.h"
#include "state_store.h"

static StarterSnapshot starter = {
    .available = true,
    .aggregate = STARTER_AGGREGATE_READY
};

static ConfigSnapshot configuration = {
    .available = true,
    .phase = CONFIG_PHASE_AWAITING,
    .mode = "NONE"
};

static void tick(LifecycleFsm *fsm, int64_t milliseconds) {

    lifecycle_fsm_tick(fsm, &starter, &configuration, 2, 2,
                       milliseconds, 1000 + milliseconds / 1000);
}

static void boot(LifecycleFsm *fsm) {

    lifecycle_fsm_init(fsm, 1000);
    assert(lifecycle_fsm_begin_check_in(fsm, true, 1, 100, 1000));
    tick(fsm, 1100);
    assert(fsm->state == LIFECYCLE_STATE_READY);
}

static void decision(const char *mode, const char *request, uint64_t generation) {

    configuration.phase = CONFIG_PHASE_APPLIED;
    configuration.generation = generation;
    configuration.revision = 0;
    snprintf(configuration.mode, sizeof(configuration.mode), "%s", mode);
    snprintf(configuration.requestId, sizeof(configuration.requestId), "%s", request);
}

static void test_flow(void) {

    LifecycleFsm fsm;
    uint64_t sequence;
    int index;

    boot(&fsm);
    for (index = 1; index <= 100; index++) {
        tick(&fsm, index * 60000);
        assert(fsm.state == LIFECYCLE_STATE_READY);
    }

    decision("NOCONFIG", "assignment-1", 1);
    tick(&fsm, 6000100);
    assert(fsm.state == LIFECYCLE_STATE_CONFIGURING);
    tick(&fsm, 6000200);
    assert(fsm.state == LIFECYCLE_STATE_OPERATIONAL);
    assert(fsm.confirmedGeneration == 1);
    sequence = fsm.sequence;
    tick(&fsm, 6000300);
    assert(fsm.sequence == sequence);

    configuration.generation = 2;
    starter.aggregate = STARTER_AGGREGATE_PENDING;
    tick(&fsm, 6000400);
    assert(fsm.sequence == sequence);
    starter.aggregate = STARTER_AGGREGATE_READY;
    tick(&fsm, 6000500);
    assert(fsm.state == LIFECYCLE_STATE_OPERATIONAL);
    assert(fsm.sequence == sequence + 1 && fsm.confirmedGeneration == 2);

    boot(&fsm);
    tick(&fsm, 1200);
    assert(fsm.state == LIFECYCLE_STATE_CONFIGURING);
    tick(&fsm, 1300);
    assert(fsm.state == LIFECYCLE_STATE_OPERATIONAL);
    puts("PASS: indefinite READY, completed-between-polls, repeat confirmation, reboot");

    decision("CONFIG", "config-1", 3);
    configuration.phase = CONFIG_PHASE_IN_PROGRESS;
    tick(&fsm, 1400);
    assert(fsm.state == LIFECYCLE_STATE_CONFIGURING);
    tick(&fsm, 9999999);
    assert(fsm.state == LIFECYCLE_STATE_CONFIGURING);
    configuration.phase = CONFIG_PHASE_FAILED;
    tick(&fsm, 10000000);
    assert(fsm.state == LIFECYCLE_STATE_FAULTY);
    configuration.phase = CONFIG_PHASE_APPLIED;
    tick(&fsm, 10000100);
    assert(fsm.state == LIFECYCLE_STATE_CONFIGURING);
    tick(&fsm, 10000200);
    assert(fsm.state == LIFECYCLE_STATE_OPERATIONAL);
    puts("PASS: CONFIG pending never succeeds on timeout; failure and recovery");

    configuration.available = false;
    tick(&fsm, 10000300);
    tick(&fsm, 10003000);
    assert(fsm.state == LIFECYCLE_STATE_FAULTY);
    configuration.available = true;
    tick(&fsm, 10003100);
    tick(&fsm, 10003200);
    assert(fsm.state == LIFECYCLE_STATE_OPERATIONAL);
    configuration.generation = 1;
    tick(&fsm, 10003300);
    assert(fsm.state == LIFECYCLE_STATE_FAULTY);
    tick(&fsm, 10003400);
    assert(fsm.state == LIFECYCLE_STATE_FAULTY);
    puts("PASS: configd outage and stale generation cannot establish Operational");

    lifecycle_fsm_init(&fsm, 1000);
    lifecycle_fsm_begin_check_in(&fsm, true, 1, 100, 1000);
    starter.available = false;
    tick(&fsm, 200);
    tick(&fsm, 2500);
    assert(fsm.state == LIFECYCLE_STATE_FAULTY);
    starter.available = true;
    starter.aggregate = STARTER_AGGREGATE_PENDING;
    tick(&fsm, 2600);
    assert(fsm.state == LIFECYCLE_STATE_CHECKING_IN);
    tick(&fsm, 2700);
    assert(fsm.gateOpen);
    puts("PASS: starter outage before gate does not deadlock startup");
}

static void test_delete(void) {

    LifecycleFsm fsm;
    uint64_t sequence;
    int index;

    starter.available = true;
    starter.aggregate = STARTER_AGGREGATE_READY;
    configuration.available = true;
    decision("NOCONFIG", "assignment-delete", 1);
    boot(&fsm);
    tick(&fsm, 1200);
    tick(&fsm, 1300);
    assert(fsm.state == LIFECYCLE_STATE_OPERATIONAL);

    configuration.phase = CONFIG_PHASE_AWAITING;
    snprintf(configuration.mode, sizeof(configuration.mode), "NONE");
    configuration.generation = 2;
    tick(&fsm, 1400);
    assert(fsm.state == LIFECYCLE_STATE_READY);
    assert(!fsm.configurationSeen && !fsm.configurationApplied);
    assert(fsm.configGeneration == 2 && fsm.confirmedGeneration == 0);
    assert(fsm.requestId[0] == '\0');
    sequence = fsm.sequence;
    for (index = 1; index <= 100; index++) {
        tick(&fsm, index * 60000);
        assert(fsm.state == LIFECYCLE_STATE_READY && fsm.sequence == sequence);
    }

    boot(&fsm);
    tick(&fsm, 1500);
    assert(fsm.state == LIFECYCLE_STATE_READY && fsm.configGeneration == 2);
    configuration.generation = 0;
    configuration.requestId[0] = '\0';
    tick(&fsm, 1600);
    assert(fsm.state == LIFECYCLE_STATE_FAULTY);
    configuration.generation = 2;
    snprintf(configuration.requestId, sizeof(configuration.requestId), "assignment-delete");
    starter.aggregate = STARTER_AGGREGATE_FAULTY;
    tick(&fsm, 1700);
    assert(fsm.state == LIFECYCLE_STATE_FAULTY);
    starter.aggregate = STARTER_AGGREGATE_READY;
    tick(&fsm, 1800);
    assert(fsm.state == LIFECYCLE_STATE_READY);

    decision("CONFIG", "config-delete", 3);
    configuration.phase = CONFIG_PHASE_IN_PROGRESS;
    tick(&fsm, 1900);
    assert(fsm.state == LIFECYCLE_STATE_CONFIGURING);
    configuration.phase = CONFIG_PHASE_AWAITING;
    snprintf(configuration.mode, sizeof(configuration.mode), "NONE");
    configuration.generation = 4;
    tick(&fsm, 2000);
    assert(fsm.state == LIFECYCLE_STATE_READY);
    decision("NOCONFIG", "new-assignment", 5);
    tick(&fsm, 2100);
    tick(&fsm, 2200);
    assert(fsm.state == LIFECYCLE_STATE_OPERATIONAL);
    puts("PASS: DELETE returns to READY, survives reboot, preserves faults and permits a new decision");
}

static void test_checkpoint(void) {

    LifecycleContext source = {0};
    LifecycleContext restored = {0};
    Config config = {0};
    char directory[] = "/tmp/lifecycle-checkpoint-XXXXXX";
    char path[512];
    FILE *file;

    assert(mkdtemp(directory));
    snprintf(path, sizeof(path), "%s/checkpoint", directory);
    config.stateFile = path;
    source.config = &config;
    restored.config = &config;
    snprintf(source.bootId, sizeof(source.bootId), "boot-1");
    snprintf(restored.bootId, sizeof(restored.bootId), "boot-1");
    lifecycle_fsm_init(&source.fsm, 1000);
    source.eventCount = 1;
    source.events[0].state = LIFECYCLE_STATE_STARTING;
    source.events[0].sequence = 1;
    source.events[0].occurredAt = 1000;
    snprintf(source.events[0].bootId, sizeof(source.events[0].bootId), "boot-1");
    assert(state_store_save(&source));
    assert(state_store_load(&restored));
    assert(restored.eventCount == 1 && restored.events[0].sequence == 1);
    snprintf(restored.bootId, sizeof(restored.bootId), "boot-2");
    assert(!state_store_load(&restored) && errno == ESTALE);
    file = fopen(path, "w");
    assert(file);
    fputs("{broken", file);
    fclose(file);
    assert(!state_store_load(&restored) && errno == EINVAL);
    unlink(path);
    rmdir(directory);
    puts("PASS: checkpoint retains events; new boot resets; corrupt state is rejected");
}

int main(void) {

    test_flow();
    test_checkpoint();
    test_delete();
    return 0;
}
