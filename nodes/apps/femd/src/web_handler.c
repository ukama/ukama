/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <stdlib.h>

#include "web_handlers.h"
#include "jserdes.h"
#include "usys_log.h"

static int respond_text(UResponse *response, int status, const char *msg) {
    ulfius_set_string_body_response(response, status, msg ? msg : "");
    return U_CALLBACK_CONTINUE;
}

static int respond_json(UResponse *response, int status, json_t *json) {

    char *s;

    if (!json) return respond_text(response, HttpStatus_InternalServerError, "serialize");
    s = json_dumps(json, 0);
    if (!s) return respond_text(response, HttpStatus_InternalServerError, "serialize");

    ulfius_set_string_body_response(response, status, s);
    free(s);
    return U_CALLBACK_CONTINUE;
}

static int parse_fem_unit(const URequest *request, FemUnit *unit) {

    const char *id;

    if (!request || !unit) return STATUS_NOK;

    id = u_map_get(request->map_url, "femId");
    if (!id) return STATUS_NOK;

    if (!strcmp(id, "1")) { *unit = FEM_UNIT_1; return STATUS_OK; }
    if (!strcmp(id, "2")) { *unit = FEM_UNIT_2; return STATUS_OK; }

    return STATUS_NOK;
}

static uint64_t parse_op_id(const URequest *request) {

    const char *id;
    char *end;
    unsigned long long v;

    if (!request) return 0;

    id = u_map_get(request->map_url, "opId");
    if (!id) return 0;

    end = NULL;
    v = strtoull(id, &end, 10);
    if (!end || *end != '\0') return 0;

    return (uint64_t)v;
}

int web_cb_default(const URequest *request, UResponse *response, void *user_data) {
    (void)request;
    (void)user_data;
    return respond_text(response, HttpStatus_NotFound, HttpStatusStr(HttpStatus_NotFound));
}

int web_cb_get_op(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    uint64_t opId;
    OpStatus st;
    json_t *j = NULL;

    if (!ctx || !ctx->jobs) return respond_text(response, HttpStatus_InternalServerError, "internal");

    opId = parse_op_id(request);
    if (opId == 0) return respond_text(response, HttpStatus_BadRequest, "bad opId");

    if (jobs_get_op(ctx->jobs, opId, &st) != STATUS_OK) {
        return respond_text(response, HttpStatus_NotFound, "not found");
    }

    if (json_serialize_op_status(&j, &st) != USYS_TRUE || !j) {
        return respond_text(response, HttpStatus_InternalServerError, "serialize");
    }

    respond_json(response, HttpStatus_Ok, j);
    json_decref(j);
    return U_CALLBACK_CONTINUE;
}

int web_cb_get_ctrl_snapshot(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    CtrlSnapshot s;
    json_t *j = NULL;

    (void)request;

    if (!ctx || !ctx->snap) return respond_text(response, HttpStatus_InternalServerError, "internal");

    if (snapshot_get_ctrl(ctx->snap, &s) != STATUS_OK) {
        return respond_text(response, HttpStatus_NotFound, "no data");
    }

    if (json_serialize_ctrl_snapshot(&j, &s) != USYS_TRUE || !j) {
        return respond_text(response, HttpStatus_InternalServerError, "serialize");
    }

    respond_json(response, HttpStatus_Ok, j);
    json_decref(j);
    return U_CALLBACK_CONTINUE;
}

int web_cb_post_ctrl_sample(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    Job job;
    uint64_t opId;
    uint32_t nowMs;
    json_t *j = NULL;

    (void)request;

    if (!ctx || !ctx->jobs) return respond_text(response, HttpStatus_InternalServerError, "internal");

    memset(&job, 0, sizeof(job));
    job.lane = LaneCtrl;
    job.cmd  = JobCmdSampleCtrl;
    job.prio = JobPrioHi;

    nowMs = snapshot_now_ms();
    opId = jobs_enqueue(ctx->jobs, &job, nowMs);
    if (opId == 0) return respond_text(response, HttpStatus_ServiceUnavailable, "queue full");

    if (json_serialize_op_id(&j, opId) != USYS_TRUE || !j) {
        return respond_text(response, HttpStatus_InternalServerError, "serialize");
    }

    respond_json(response, HttpStatus_Accepted, j);
    json_decref(j);
    return U_CALLBACK_CONTINUE;
}

int web_cb_get_fem_snapshot(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    FemSnapshot s;
    json_t *j = NULL;

    if (!ctx || !ctx->snap) return respond_text(response, HttpStatus_InternalServerError, "internal");

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    if (snapshot_get_fem(ctx->snap, unit, &s) != STATUS_OK) {
        return respond_text(response, HttpStatus_NotFound, "no data");
    }

    if (json_serialize_fem_snapshot(&j, unit, &s) != USYS_TRUE || !j) {
        return respond_text(response, HttpStatus_InternalServerError, "serialize");
    }

    respond_json(response, HttpStatus_Ok, j);
    json_decref(j);
    return U_CALLBACK_CONTINUE;
}

int web_cb_post_fem_sample(const URequest *request, UResponse *response, void *user_data) {

    WebCtx *ctx = (WebCtx *)user_data;
    FemUnit unit;
    Job job;
    uint64_t opId;
    uint32_t nowMs;
    json_t *j = NULL;

    if (!ctx || !ctx->jobs) return respond_text(response, HttpStatus_InternalServerError, "internal");

    if (parse_fem_unit(request, &unit) != STATUS_OK) {
        return respond_text(response, HttpStatus_BadRequest, "bad femId");
    }

    memset(&job, 0, sizeof(job));
    job.lane    = (unit == FEM_UNIT_1) ? LaneFem1 : LaneFem2;
    job.femUnit = unit;
    job.cmd     = JobCmdSampleFem;
    job.prio    = JobPrioHi;

    nowMs = snapshot_now_ms();
    opId = jobs_enqueue(ctx->jobs, &job, nowMs);
    if (opId == 0) return respond_text(response, HttpStatus_ServiceUnavailable, "queue full");

    if (json_serialize_op_id(&j, opId) != USYS_TRUE || !j) {
        return respond_text(response, HttpStatus_InternalServerError, "serialize");
    }

    respond_json(response, HttpStatus_Accepted, j);
    json_decref(j);

    return U_CALLBACK_CONTINUE;
}
