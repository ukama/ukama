/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
