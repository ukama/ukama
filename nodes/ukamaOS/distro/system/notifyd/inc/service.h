/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef INC_SERVICE_H_
#define INC_SERVICE_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "config.h"

/**
 * @fn     int service_init(char*, char*, char*)
 * @brief  Perform the pre-requisites required for node service to run.
 *         1. Initialize web framework.
 *
 * @param  cfg
 * @return On Success, USYS_OK (0)
 *         On Failure, USYS_NOK (-1)
 */
int service_init(Config* cfg);

/**
 * @fn     int service_at_exit()
 * @brief  Perform the cleanup for clean exit of the node service.
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
