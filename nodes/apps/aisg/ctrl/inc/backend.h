/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef BACKEND_H_
#define BACKEND_H_

#include "ctrl.h"
#include "request.h"
#include "response.h"

typedef struct Backend Backend;

typedef struct {
    bool (*open)(Backend *backend);
    void (*close)(Backend *backend);
    bool (*execute)(Backend *backend,
                    CtrlRequest *request,
                    CtrlResponse *response);
} BackendOps;

struct Backend {
    Config *config;
    BackendOps ops;
    void *priv;
};

bool backend_init(Backend *backend, Config *config);
bool backend_open(Backend *backend);
void backend_close(Backend *backend);
bool backend_execute(Backend *backend,
                     CtrlRequest *request,
                     CtrlResponse *response);

#endif /* BACKEND_H_ */
