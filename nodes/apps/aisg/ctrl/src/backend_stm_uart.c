/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>

#include "backend_stm_uart.h"

typedef struct {
    bool opened;
} StmBackend;

static bool stm_open(Backend *backend)
{
    StmBackend *ctx = NULL;

    if (backend == NULL || backend->priv == NULL) {
        return false;
    }

    ctx = backend->priv;
    ctx->opened = true;

    return true;
}

static void stm_close(Backend *backend)
{
    if (backend == NULL) {
        return;
    }

    free(backend->priv);
    backend->priv = NULL;
}

static bool stm_execute(Backend *backend,
                        CtrlRequest *request,
                        CtrlResponse *response)
{
    (void)backend;
    (void)request;

    return ctrl_response_set_error(response,
                                   CtrlCodeUnsupportedProcedure,
                                   "stm-uart backend reserved for firmware");
}

bool backend_stm_uart_init(Backend *backend, Config *config)
{
    StmBackend *ctx = NULL;

    if (backend == NULL || config == NULL) {
        return false;
    }

    ctx = calloc(1, sizeof(StmBackend));
    if (ctx == NULL) {
        return false;
    }

    backend->config      = config;
    backend->priv        = ctx;
    backend->ops.open    = stm_open;
    backend->ops.close   = stm_close;
    backend->ops.execute = stm_execute;

    return true;
}
