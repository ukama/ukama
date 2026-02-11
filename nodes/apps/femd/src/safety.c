/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "safety.h"
#include "usys_log.h"

static int enqueue_disable(Safety *s, FemUnit unit) {

    Job j;
    uint32_t nowMs;

    memset(&j, 0, sizeof(j));
    j.lane = (unit == FEM_UNIT_1) ? LaneFem1 : LaneFem2;
    j.femUnit = unit;
    j.cmd = JobCmdDisablePa;
    j.prio = JobPrioHi;

    nowMs = snapshot_now_ms();
    return jobs_enqueue(s->jobs, &j, nowMs) ? STATUS_OK : STATUS_NOK;
}

static int enqueue_restore(Safety *s, FemUnit unit) {

    Job j;
    uint32_t nowMs;

    memset(&j, 0, sizeof(j));
    j.lane = (unit == FEM_UNIT_1) ? LaneFem1 : LaneFem2;
    j.femUnit = unit;
    j.cmd = JobCmdRestorePa;
    j.prio = JobPrioHi;

    nowMs = snapshot_now_ms();
    return jobs_enqueue(s->jobs, &j, nowMs) ? STATUS_OK : STATUS_NOK;
}

static bool health_ok(const YamlSafetyConfig *c, const FemSnapshot *sn) {

    if (!c || !sn) return false;
    if (!sn->present) return false;
    if (!sn->haveAdc || !sn->haveTemp) return false;

    if (sn->reversePowerDbm > c->max_reverse_power_dbm) return false;
    if (sn->forwardPowerDbm > c->max_forward_power_dbm) return false;
    if (sn->paCurrentA > c->max_pa_current_a) return false;

    if (sn->tempC > c->max_temperature_c) return false;
    if (sn->tempC < c->min_temperature_c) return false;

    return true;
}

int safety_init(Safety *s, Jobs *jobs, SnapshotStore *snap, Notifier *n, const char *yamlPath) {

    if (!s || !jobs || !snap || !yamlPath) return STATUS_NOK;

    memset(s, 0, sizeof(*s));
    s->jobs = jobs;
    s->snap = snap;
    s->notifier = n;

    yaml_config_set_defaults(&s->st.cfg);
    (void)yaml_config_load(yamlPath, &s->st.cfg);

    if (yaml_config_validate(&s->st.cfg) != STATUS_OK) {
        return STATUS_NOK;
    }

    yaml_config_print(&s->st.cfg);
    return STATUS_OK;
}

int safety_tick(Safety *s, FemUnit unit) {

    FemSnapshot sn;
    uint32_t nowMs;
    uint32_t elapsed;

    if (!s || !s->jobs || !s->snap) return STATUS_NOK;
    if (unit != FEM_UNIT_1 && unit != FEM_UNIT_2) return STATUS_NOK;

    if (!s->st.cfg.enabled) return STATUS_OK;

    if (snapshot_get_fem(s->snap, unit, &sn) != STATUS_OK) {
        return STATUS_OK;
    }

    nowMs = snapshot_now_ms();

    if (!health_ok(&s->st.cfg, &sn)) {
        if (s->st.violationCount[unit] < 0xFFFFFFFFu) s->st.violationCount[unit]++;

        if (!s->st.paShutdown[unit]) {
            if (enqueue_disable(s, unit) == STATUS_OK) {
                s->st.paShutdown[unit] = true;
                s->st.lastShutdownMs[unit] = nowMs;
                s->st.okStreak[unit] = 0;
                if (s->notifier) {
                    (void)notifier_send_pa_alarm(s->notifier, ALARM_TYPE_PA_OFF, NULL);
                }
                usys_log_warn("safety: fem=%d pa disabled", unit);
            }
        }
        return STATUS_OK;
    }

    s->st.violationCount[unit] = 0;

    if (!s->st.cfg.auto_restore_enabled) return STATUS_OK;
    if (!s->st.paShutdown[unit]) return STATUS_OK;

    elapsed = nowMs - s->st.lastShutdownMs[unit];
    if (elapsed < s->st.cfg.restore_cooldown_ms) return STATUS_OK;

    if (s->st.okStreak[unit] < 0xFFFFFFFFu) s->st.okStreak[unit]++;

    if (s->st.okStreak[unit] >= s->st.cfg.restore_ok_checks) {
        if (enqueue_restore(s, unit) == STATUS_OK) {
            s->st.paShutdown[unit] = false;
            s->st.okStreak[unit] = 0;
            if (s->notifier) {
                (void)notifier_send_pa_alarm(s->notifier, ALARM_TYPE_PA_ON, NULL);
            }
            usys_log_info("safety: fem=%d pa restore", unit);
        }
    }

    return STATUS_OK;
}
