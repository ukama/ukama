/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_SERVICE_H_
#define INC_SERVICE_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "service.h"

/**
 * @fn int service_init(char*, char*, char*)
 * @brief  Perform the pre-requisites required for node service to run.
 *         1. Parsing of the property JSON
 *         2. Reading node schema from the inventory.
 *         3. Initializing ledger.
 *         4. Set notification server url.
 *         5. Initialize web framework.
 *
 * @param invtDb
 * @param propCfg
 * @param notifServer
 * @return On Success, USYS_OK (0)
 *         On Failure, USYS_NOK (-1)
 */
int service_init(char *invtDb, char *propCfg, char* notifServer);

/**
 * @fn     int service_at_exit()
 * @brief  Perform the cleanup for clean exit of the node service.
 * 		   1. Closes ledger.
 * 		   2. Stop communication module.
 *
 * @return On Success, USYS_OK (0)
 * 		   On Failure, USYS_NOK (-1)
 */
int service_at_exit();

/**
 * @fn     void service()
 * @brief  Start NodeD service by starting communication module.
 *
 */
void service();

#ifdef __cplusplus
}
#endif

#endif /* INC_SERVICE_H_ */
