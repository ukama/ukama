/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WEB_CLIENT_H_
#define WEB_CLIENT_H_

#include <stddef.h>

#include "types.h"

int web_client_post_json(const char *url,
                         const char *json,
                         long timeoutMs,
                         long *status,
                         char *response,
                         size_t responseLen);
int web_client_notify_alarm(const SwitchdConfig *cfg,
                            const SwitchAlarm *alarm,
                            bool clear);

#endif /* WEB_CLIENT_H_ */
