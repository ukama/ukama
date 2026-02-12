/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "usys_log.h"

#include "safety.h"
#include "snapshot.h"
#include "jobs.h"

static uint32_t now_ms(void) {
    return snapshot_now_ms();
}

int safety_init(Safety *s, Jobs *jobs, SnapshotStore *snap, const SafetyConfig *cfg) {
    if (!s || !jobs || !snap || !cfg) return STATUS_NOK;

    memset(s, 0, sizeof(*s));
    s->jobs = jobs;
    s->snap = snap;
    s->st.cfg = *cfg;
    s->st.initialized = true;
    return STATUS_OK;
}

void safety_cleanup(Safety *s) {
    if (!s) return;
    memset(s, 0, sizeof(*s));
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

int safety_tick(Safety *s, FemUnit unit) {
    FemSnapshot sn;
    bool violation = false;
    const float maxT = s->st.cfg.maxTemperatureC;
    const float maxR = s->st.cfg.maxReversePowerDbm;
    const float maxI = s->st.cfg.maxPaCurrentA;

    if (!s || !s->st.initialized) return STATUS_NOK;

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

    LaneId lane = (unit == FEM_UNIT_1) ? LaneFem1 : LaneFem2;

    if (violation && !s->st.paDisabled[unit]) {
        s->st.paDisabled[unit] = true;
        enqueue_simple(s, lane, unit, JobCmdSafetyDisablePa);
        return STATUS_OK;
    }

    if (!violation && s->st.paDisabled[unit]) {
        s->st.paDisabled[unit] = false;
        enqueue_simple(s, lane, unit, JobCmdSafetyRestorePa);
        return STATUS_OK;
    }

    return STATUS_OK;
}
