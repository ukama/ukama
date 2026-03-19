/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>

#include "driver.h"

int tycon_driver_attach(SwitchdContext *ctx);

int driver_init(SwitchdContext *ctx) {
    if (ctx == NULL) {
        return SWITCHD_ERR_INVAL;
    }

    if (strcmp(ctx->config.driverName, "tycon_snmp") == 0 ||
        strcmp(ctx->config.driverName, "tycon") == 0) {
        return tycon_driver_attach(ctx);
    }

    return SWITCHD_ERR_UNSUPPORTED;
}

void driver_cleanup(SwitchdContext *ctx) {
    if (ctx == NULL || ctx->driver == NULL) {
        return;
    }

    if (ctx->driver->ops.cleanup != NULL) {
        ctx->driver->ops.cleanup(ctx);
    }

    free(ctx->driver);
    ctx->driver = NULL;
}
