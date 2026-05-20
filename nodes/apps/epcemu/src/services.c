/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "epcemu.h"
#include "services.h"

#include "usys_services.h"
#include "usys_file.h"

static int resolve_port(const char *name, int *port) {

    int value;

    if (name == NULL || port == NULL) return USYS_FALSE;

    value = usys_find_service_port((char *)name);
    if (value <= 0) {
        usys_log_error("service %s not found in /etc/services", name);
        return USYS_FALSE;
    }

    *port = value;
    return USYS_TRUE;
}

int services_resolve(EpcemuConfig *config, EpcemuStatus *status) {

    if (config == NULL || status == NULL) return USYS_FALSE;

    status_set(status, EpcemuStateResolvingServices,
               "resolving local services");

    if (!resolve_port(EPCEMU_SERVICE_NAME, &config->servicePort)) {
        status_fail(status, "failed to resolve epcemu service port");
        return USYS_FALSE;
    }

    if (!resolve_port(EPCEMU_PCRF_SERVICE, &config->pcrfPort)) {
        status_fail(status, "failed to resolve pcrf service port");
        return USYS_FALSE;
    }

    if (!resolve_port(EPCEMU_INITNET_SERVICE, &config->initNetworkPort)) {
        status_fail(status, "failed to resolve init-network service port");
        return USYS_FALSE;
    }

    snprintf(config->pcrfUrl, sizeof(config->pcrfUrl),
             "http://localhost:%d", config->pcrfPort);

    snprintf(config->initNetworkUrl, sizeof(config->initNetworkUrl),
             "http://localhost:%d", config->initNetworkPort);

    usys_log_debug("resolved services epcemu=%d pcrf=%d init-network=%d",
                   config->servicePort,
                   config->pcrfPort,
                   config->initNetworkPort);

    return USYS_TRUE;
}
