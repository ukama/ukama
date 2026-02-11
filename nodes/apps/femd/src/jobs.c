/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "jobs.h"
#include "usys_log.h"

static inline bool valid_lane(LaneId lane) {
    return lane >= LaneCtrl && lane < LaneMax;
}

static int q_push_hi(LaneQueue *q, const Job *job) {

    if (q->hiCount >= (uint32_t)(sizeof(q->hiQ) / sizeof(q->hiQ[0]))) {
        return STATUS_NOK;
    }

    q->hiQ[q->hiTail] = *job;
    q->hiTail = (q->hiTail + 1U) % (uint32_t)(sizeof(q->hiQ) / sizeof(q->hiQ[0]));
    q->hiCount++;
    return STATUS_OK;
}

static int q_push_lo(LaneQueue *q, const Job *job) {

    if (q->loCount >= (uint32_t)(sizeof(q->loQ) / sizeof(q->loQ[0]))) {
        return STATUS_NOK;
    }

    q->loQ[q->loTail] = *job;
    q->loTail = (q->loTail + 1U) % (uint32_t)(sizeof(q->loQ) / sizeof(q->loQ[0]));
    q->loCount++;
    return STATUS_OK;
}

static int q_pop_hi(LaneQueue *q, Job *out) {

    if (q->hiCount == 0U) return STATUS_NOK;

    *out = q->hiQ[q->hiHead];
    q->hiHead = (q->hiHead + 1U) % (uint32_t)(sizeof(q->hiQ) / sizeof(q->hiQ[0]));
    q->hiCount--;
    return STATUS_OK;
}

static int q_pop_lo(LaneQueue *q, Job *out) {

    if (q->loCount == 0U) return STATUS_NOK;

    *out = q->loQ[q->loHead];
    q->loHead = (q->loHead + 1U) % (uint32_t)(sizeof(q->loQ) / sizeof(q->loQ[0]));
    q->loCount--;
    return STATUS_OK;
}

static OpStatus *op_slot(Jobs *jobs, uint64_t opId) {

    if (opId == 0) return NULL;

    for (uint32_t i = 0; i < (uint32_t)(sizeof(jobs->table) / sizeof(jobs->table[0])); i++) {
        if (jobs->table[i].opId == opId) return &jobs->table[i];
    }
    return NULL;
}

static OpStatus *op_insert(Jobs *jobs, const Job *job, uint32_t nowMs) {

    OpStatus *s = &jobs->table[jobs->tablePos];
    memset(s, 0, sizeof(*s));

    s->opId      = job->opId;
    s->lane      = job->lane;
    s->femUnit   = job->femUnit;
    s->cmd       = job->cmd;
    s->state     = OpStateQueued;
    s->result    = STATUS_OK;
    s->createdMs = nowMs;

    jobs->tablePos = (jobs->tablePos + 1U) % (uint32_t)(sizeof(jobs->table) / sizeof(jobs->table[0]));
    return s;
}

int jobs_init(Jobs *jobs) {

    if (!jobs) return STATUS_NOK;

    memset(jobs, 0, sizeof(*jobs));
    if (pthread_mutex_init(&jobs->mu, NULL) != 0) return STATUS_NOK;

    jobs->nextOpId = 1;

    for (int i = 0; i < (int)LaneMax; i++) {
        LaneQueue *q = &jobs->lane[i];
        if (pthread_mutex_init(&q->mu, NULL) != 0) return STATUS_NOK;
        if (pthread_cond_init(&q->cv, NULL) != 0) return STATUS_NOK;
        q->stop = false;
    }

    jobs->initialized = true;
    return STATUS_OK;
}

void jobs_cleanup(Jobs *jobs) {

    if (!jobs || !jobs->initialized) return;

    for (int i = 0; i < (int)LaneMax; i++) {
        LaneQueue *q = &jobs->lane[i];
        (void)pthread_cond_destroy(&q->cv);
        (void)pthread_mutex_destroy(&q->mu);
    }

    (void)pthread_mutex_destroy(&jobs->mu);
    memset(jobs, 0, sizeof(*jobs));
}

uint64_t jobs_enqueue(Jobs *jobs, const Job *in, uint32_t nowMs) {

    Job job;
    LaneQueue *q;
    uint64_t opId = 0;

    if (!jobs || !jobs->initialized || !in) return 0;
    if (!valid_lane(in->lane)) return 0;

    q = &jobs->lane[in->lane];
    job = *in;

    if (pthread_mutex_lock(&jobs->mu) != 0) return 0;
    opId = jobs->nextOpId++;
    if (opId == 0) opId = jobs->nextOpId++;
    job.opId = opId;
    (void)op_insert(jobs, &job, nowMs);
    (void)pthread_mutex_unlock(&jobs->mu);

    if (pthread_mutex_lock(&q->mu) != 0) return 0;

    if (q->stop) {
        (void)pthread_mutex_unlock(&q->mu);
        (void)jobs_set_op_state(jobs, opId, OpStateCanceled, STATUS_NOK, nowMs);
        return 0;
    }

    if ((job.prio == JobPrioHi ? q_push_hi(q, &job) : q_push_lo(q, &job)) != STATUS_OK) {
        (void)pthread_mutex_unlock(&q->mu);
        (void)jobs_set_op_state(jobs, opId, OpStateFailed, STATUS_NOK, nowMs);
        return 0;
    }

    (void)pthread_cond_signal(&q->cv);
    (void)pthread_mutex_unlock(&q->mu);

    return opId;
}

int jobs_dequeue(Jobs *jobs, LaneId lane, Job *out, uint32_t nowMs) {

    LaneQueue *q;
    Job job;

    if (!jobs || !jobs->initialized || !out) return STATUS_NOK;
    if (!valid_lane(lane)) return STATUS_NOK;

    q = &jobs->lane[lane];

    if (pthread_mutex_lock(&q->mu) != 0) return STATUS_NOK;

    while (!q->stop && q->hiCount == 0U && q->loCount == 0U) {
        (void)pthread_cond_wait(&q->cv, &q->mu);
    }

    if (q->stop && q->hiCount == 0U && q->loCount == 0U) {
        (void)pthread_mutex_unlock(&q->mu);
        return STATUS_NOK;
    }

    if (q_pop_hi(q, &job) != STATUS_OK) {
        if (q_pop_lo(q, &job) != STATUS_OK) {
            (void)pthread_mutex_unlock(&q->mu);
            return STATUS_NOK;
        }
    }

    (void)pthread_mutex_unlock(&q->mu);

    *out = job;
    (void)jobs_set_op_state(jobs, job.opId, OpStateRunning, STATUS_OK, nowMs);

    return STATUS_OK;
}

int jobs_set_op_state(Jobs *jobs, uint64_t opId, OpState state, int result, uint32_t nowMs) {

    OpStatus *s;

    if (!jobs || !jobs->initialized || opId == 0) return STATUS_NOK;

    if (pthread_mutex_lock(&jobs->mu) != 0) return STATUS_NOK;

    s = op_slot(jobs, opId);
    if (!s) {
        (void)pthread_mutex_unlock(&jobs->mu);
        return STATUS_NOK;
    }

    s->state  = state;
    s->result = result;

    if (state == OpStateRunning && s->startedMs == 0U) {
        s->startedMs = nowMs;
    }
    if ((state == OpStateDone ||
         state == OpStateFailed ||
         state == OpStateCanceled) && s->endedMs == 0U) {
        s->endedMs = nowMs;
    }

    (void)pthread_mutex_unlock(&jobs->mu);
    return STATUS_OK;
}

int jobs_get_op(Jobs *jobs, uint64_t opId, OpStatus *out) {

    OpStatus *s;

    if (!jobs || !jobs->initialized || !out || opId == 0) return STATUS_NOK;

    if (pthread_mutex_lock(&jobs->mu) != 0) return STATUS_NOK;

    s = op_slot(jobs, opId);
    if (!s) {
        (void)pthread_mutex_unlock(&jobs->mu);
        return STATUS_NOK;
    }

    *out = *s;

    (void)pthread_mutex_unlock(&jobs->mu);
    return STATUS_OK;
}

int jobs_request_lane_stop(Jobs *jobs, LaneId lane) {

    LaneQueue *q;

    if (!jobs || !jobs->initialized) return STATUS_NOK;
    if (!valid_lane(lane)) return STATUS_NOK;

    q = &jobs->lane[lane];

    if (pthread_mutex_lock(&q->mu) != 0) return STATUS_NOK;

    q->stop = true;
    (void)pthread_cond_broadcast(&q->cv);

    (void)pthread_mutex_unlock(&q->mu);

    usys_log_info("jobs: lane stop requested lane=%d", (int)lane);
    return STATUS_OK;
}
