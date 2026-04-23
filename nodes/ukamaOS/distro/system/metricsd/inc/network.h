/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef INC_NETWORK_H_
#define INC_NETWORK_H_

#include "web_service.h"

int start_admin_web_service(UInst *adminInst, int configPort);
void stop_admin_web_service(UInst *adminInst);

#endif /* INC_NETWORK_H_ */
