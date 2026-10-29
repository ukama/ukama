/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef PDN_H_
#define PDN_H_

#include "epcemu.h"

#define PDN_ATTACH_TIMEOUT_SEC 30

int pdn_init(ServiceContext *ctx);
void pdn_poll(ServiceContext *ctx);
void pdn_shutdown(ServiceContext *ctx);
int pdn_attach(ServiceContext *ctx, JsonObj *body, JsonObj **reply,
                char *reason, size_t reasonLen);
int pdn_detach(ServiceContext *ctx, JsonObj *body,
                char *reason, size_t reasonLen);

#endif /* PDN_H_ */
