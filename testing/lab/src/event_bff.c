/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "event.h"

int event_bff(event_ctx_t *ctx, const event_spec_t *event,
              ulab_error_t *err) {
    (void)ctx;
    (void)event;
    snprintf(err->msg, sizeof(err->msg),
             "create_ues is defined but not enabled in v1.0 runtime path");
    return ULAB_ERR;
}
