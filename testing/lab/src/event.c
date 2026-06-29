/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <stdio.h>

#include "log.h"
#include "event.h"
#include "util.h"

int event_traffic(event_ctx_t *ctx, const event_spec_t *event,
                  ulab_error_t *err);
int event_runtime(event_ctx_t *ctx, const event_spec_t *event,
                  ulab_error_t *err);
int event_bff(event_ctx_t *ctx, const event_spec_t *event,
              ulab_error_t *err);

static int event_run_actual(event_ctx_t *ctx, const event_spec_t *event,
                            ulab_error_t *err) {
    ulab_status("EVENT", "%s", scenario_event_name(event->type));
    switch (event->type) {
    case EVT_TRAFFIC:
    case EVT_TRAFFIC_BY_PROFILE:
        return event_traffic(ctx, event, err);
    case EVT_CREATE_UES:
        return event_bff(ctx, event, err);
    case EVT_START_UES:
    case EVT_WAIT_UES_ATTACHED:
    case EVT_RESTART_NODES:
    case EVT_WAIT_NODES_READY:
        return event_runtime(ctx, event, err);
    case EVT_CHECK:
        return ULAB_OK;
    default:
        snprintf(err->msg, sizeof(err->msg), "unsupported event");
        return ULAB_ERR;
    }
}

int event_run(event_ctx_t *ctx, const event_spec_t *event,
              ulab_error_t *err) {
    ulab_error_t actual_err;
    int expect_failure;
    int rc;

    memset(&actual_err, 0, sizeof(actual_err));
    expect_failure = ulab_streq(event->expect_result, "failure") ||
                     ulab_streq(event->expect_result, "fail");

    rc = event_run_actual(ctx, event, &actual_err);

    if (!expect_failure) {
        if (rc != ULAB_OK) {
            *err = actual_err;
        }
        return rc;
    }

    if (rc == ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg),
                 "expected event failure but event succeeded");
        return ULAB_ERR;
    }

    if ((ulab_streq(event->expect_result, "failure") ||
         ulab_streq(event->expect_result, "fail")) &&
        event->error_contains[0] &&
        strstr(actual_err.msg, event->error_contains) == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "expected failure containing '%.256s', got '%.512s'",
                 event->error_contains, actual_err.msg);
        return ULAB_ERR;
    }

    return ULAB_OK;
}

void event_list_supported(void) {
    scenario_list_events();
}
