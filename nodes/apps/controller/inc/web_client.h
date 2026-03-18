/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WEB_CLIENT_H
#define WEB_CLIENT_H

#include "config.h"
#include "alarms.h"

int get_nodeid_and_type_from_noded(Config *config);

int alarms_send_notification(const Config *config,
                             AlarmType type,
                             Severity severity,
                             const char *message);

#endif /* WEB_CLIENT_H */
