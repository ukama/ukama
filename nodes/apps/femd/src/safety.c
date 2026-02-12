/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <unistd.h>

#include "safety.h"
#include "yaml_config.h"
#include "usys_log.h"

static uint32_t now_ms(void) {
    return snapshot_now_ms();
}

static void enqueue_simple(Safety *s, LaneId lane, FemUnit unit, JobCmd cmd) {

    Job j;

    memset(&j, 0, sizeof(j));
    j.lane = lane;
    j.femUnit = unit;
    j.prio = JobPrioHi;
    j.cmd = cmd;

    (void)jobs_enqueue(s->jobs, &j, now_ms());
}

static void* safety_thread_main(void *arg) {

    Safety *s = (Safety *)arg;

    while (1) {

        pthread_mutex_lock(&s->mu);
        if (!s->running) {
            pthread_mutex_unlock(&s->mu);
            break;
        }
        pthread_mutex_unlock(&s->mu);

        (void)safety_tick(s, FEM_UNIT_1);
        (void)safety_tick(s, FEM_UNIT_2);

        usleep(250 * 1000);
    }

    return NULL;
}

int safety_init(Safety *s, Jobs *jobs, SnapshotStore *snap, Notifier *notifier, const char *cfgPath) {

    YamlSafetyConfig y;
    int ok = 0;

    if (!s || !jobs || !snap) return STATUS_NOK;

    memset(s, 0, sizeof(*s));
    pthread_mutex_init(&s->mu, NULL);

    s->jobs = jobs;
    s->snap = snap;
    s->notifier = notifier;

    memset(&s->cfg, 0, sizeof(s->cfg));
    s->cfg.max_temperature_c     = 95.0f;
    s->cfg.max_reverse_power_dbm = 10.0f;
    s->cfg.max_pa_current_a      = 5.0f;

    memset(&y, 0, sizeof(y));
    if (cfgPath && cfgPath[0]) {
        if (yaml_config_load(&y, cfgPath) == STATUS_OK && y.enabled) {
            //            s->cfg = y.safety; XXXX
            ok = 1;
        }
        // yaml_config_cleanup(&y); XXXX
    }

    if (!ok) {
        usys_log_debug("Safety: using defaults (cfg not loaded or missing safety section)");
    }

    s->initialized = true;
    return STATUS_OK;
}

void safety_cleanup(Safety *s) {

    if (!s) return;

    (void)safety_stop(s);

    pthread_mutex_destroy(&s->mu);
    memset(s, 0, sizeof(*s));
}

int safety_start(Safety *s) {

    if (!s || !s->initialized) return STATUS_NOK;

    pthread_mutex_lock(&s->mu);
    if (s->running) {
        pthread_mutex_unlock(&s->mu);
        return STATUS_OK;
    }
    s->running = true;
    pthread_mutex_unlock(&s->mu);

    if (pthread_create(&s->thread, NULL, safety_thread_main, s) != 0) {
        pthread_mutex_lock(&s->mu);
        s->running = false;
        pthread_mutex_unlock(&s->mu);
        return STATUS_NOK;
    }

    return STATUS_OK;
}

int safety_stop(Safety *s) {

    if (!s || !s->initialized) return STATUS_NOK;

    pthread_mutex_lock(&s->mu);
    if (!s->running) {
        pthread_mutex_unlock(&s->mu);
        return STATUS_OK;
    }
    s->running = false;
    pthread_mutex_unlock(&s->mu);

    (void)pthread_join(s->thread, NULL);
    return STATUS_OK;
}

int safety_get_config(Safety *s, SafetyConfig *out) {

    if (!s || !out || !s->initialized) return STATUS_NOK;

    pthread_mutex_lock(&s->mu);
    *out = s->cfg;
    pthread_mutex_unlock(&s->mu);

    return STATUS_OK;
}

int safety_set_config(Safety *s, const SafetyConfig *in) {

    if (!s || !in || !s->initialized) return STATUS_NOK;

    pthread_mutex_lock(&s->mu);
    s->cfg = *in;
    pthread_mutex_unlock(&s->mu);

    return STATUS_OK;
}

int safety_force_restore(Safety *s, FemUnit unit) {

    LaneId lane;

    if (!s || !s->initialized) return STATUS_NOK;
    if (unit != FEM_UNIT_1 && unit != FEM_UNIT_2) return STATUS_NOK;

    lane = (unit == FEM_UNIT_1) ? LaneFem1 : LaneFem2;

    pthread_mutex_lock(&s->mu);
    s->paDisabled[unit] = false;
    pthread_mutex_unlock(&s->mu);

    enqueue_simple(s, lane, unit, JobCmdSafetyRestorePa);
    return STATUS_OK;
}

int safety_tick(Safety *s, FemUnit unit) {

    FemSnapshot sn;
    bool violation = false;
    LaneId lane;

    float maxT;
    float maxR;
    float maxI;

    if (!s || !s->initialized) return STATUS_NOK;
    if (unit != FEM_UNIT_1 && unit != FEM_UNIT_2) return STATUS_NOK;

    lane = (unit == FEM_UNIT_1) ? LaneFem1 : LaneFem2;

    pthread_mutex_lock(&s->mu);
    maxT = s->cfg.max_temperature_c;
    maxR = s->cfg.max_reverse_power_dbm;
    maxI = s->cfg.max_pa_current_a;
    pthread_mutex_unlock(&s->mu);

    memset(&sn, 0, sizeof(sn));
    if (snapshot_get_fem(s->snap, unit, &sn) != STATUS_OK || !sn.present) {
        return STATUS_OK;
    }

    if (sn.haveTemp && sn.tempC > maxT) {
        violation = true;
        usys_log_warn("Safety: FEM%d temp=%.1fC > %.1fC", unit, sn.tempC, maxT);
    }

    if (sn.haveAdc && sn.reversePowerDbm > maxR) {
        violation = true;
        usys_log_warn("Safety: FEM%d reverseP=%.1fdBm > %.1fdBm", unit, sn.reversePowerDbm, maxR);
    }

    if (sn.haveAdc && sn.paCurrentA > maxI) {
        violation = true;
        usys_log_warn("Safety: FEM%d current=%.2fA > %.2fA", unit, sn.paCurrentA, maxI);
    }

    pthread_mutex_lock(&s->mu);

    if (violation && !s->paDisabled[unit]) {
        s->paDisabled[unit] = true;
        pthread_mutex_unlock(&s->mu);
        enqueue_simple(s, lane, unit, JobCmdSafetyDisablePa);
        return STATUS_OK;
    }

    if (!violation && s->paDisabled[unit]) {
        s->paDisabled[unit] = false;
        pthread_mutex_unlock(&s->mu);
        enqueue_simple(s, lane, unit, JobCmdSafetyRestorePa);
        return STATUS_OK;
    }

    pthread_mutex_unlock(&s->mu);
    return STATUS_OK;
}
