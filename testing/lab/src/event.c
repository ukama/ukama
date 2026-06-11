/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "log.h"
#include "event.h"

int event_traffic(event_ctx_t *ctx, const event_spec_t *event,
                  ulab_error_t *err);
int event_runtime(event_ctx_t *ctx, const event_spec_t *event,
                  ulab_error_t *err);
int event_bff(event_ctx_t *ctx, const event_spec_t *event,
              ulab_error_t *err);

int event_run(event_ctx_t *ctx, const event_spec_t *event,
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

void event_list_supported(void) {
    scenario_list_events();
}
