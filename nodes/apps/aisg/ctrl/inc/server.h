/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SERVER_H_
#define SERVER_H_

#include "backend.h"
#include "config.h"

bool ctrl_server_run(Config *config, Backend *backend, volatile bool *running);

#endif /* SERVER_H_ */
