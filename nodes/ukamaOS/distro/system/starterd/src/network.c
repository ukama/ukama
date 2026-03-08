/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "network.h"
#include "web_service.h"

#include <ulfius.h>

#include "usys_log.h"

bool network_init(StarterContext *ctx) {

    if (!ctx || !ctx->config) return false;

    ctx->uInstance = calloc(1, sizeof(struct _u_instance));
    if (!ctx->uInstance) return false;

    if (ulfius_init_instance(ctx->uInstance,
                             ctx->config->httpPort,
                             NULL,
                             NULL) != U_OK) {
        usys_log_error("network: init failed");
        free(ctx->uInstance);
        ctx->uInstance = NULL;
        return false;
    }

    return true;
}

void network_shutdown(StarterContext *ctx) {

    if (!ctx || !ctx->uInstance) return;

    ulfius_stop_framework(ctx->uInstance);
    ulfius_clean_instance(ctx->uInstance);
    free(ctx->uInstance);
    ctx->uInstance = NULL;
}
