/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "backend.h"
#include "backend_raw_rs485.h"
#include "backend_stm_uart.h"

bool backend_init(Backend *backend, Config *config) {
    memset(backend, 0, sizeof(Backend));
    backend->config = config;

    if (config->backendType == BackendTypeStmUart) {
        return backend_stm_uart_init(backend, config);
    }

    return backend_raw_rs485_init(backend, config);
}


bool backend_open(Backend *backend) {
    if (backend == NULL || backend->ops.open == NULL) {
        return false;
    }

    return backend->ops.open(backend);
}

void backend_close(Backend *backend) {
    if (backend != NULL && backend->ops.close != NULL) {
        backend->ops.close(backend);
    }
}

bool backend_execute(Backend *backend,
                     CtrlRequest *request,
                     CtrlResponse *response) {
    if (backend == NULL || backend->ops.execute == NULL) {
        return ctrl_response_set_error(response,
                                       CtrlCodeUnsupportedProcedure,
                                       "backend not available");
    }

    return backend->ops.execute(backend, request, response);
}
