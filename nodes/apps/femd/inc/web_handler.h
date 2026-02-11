/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WEB_HANDLERS_H
#define WEB_HANDLERS_H

#include <ulfius.h>

#include "web_service.h"

int web_cb_default(const URequest *request, UResponse *response, void *user_data);

int web_cb_get_op(const URequest *request, UResponse *response, void *user_data);

int web_cb_get_ctrl_snapshot(const URequest *request, UResponse *response, void *user_data);
int web_cb_post_ctrl_sample(const URequest *request, UResponse *response, void *user_data);

int web_cb_get_fem_snapshot(const URequest *request, UResponse *response, void *user_data);
int web_cb_post_fem_sample(const URequest *request, UResponse *response, void *user_data);

#endif /* WEB_HANDLERS_H */
