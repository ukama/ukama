/**
 * Copyright (c) 2021-present, Ukama Inc.
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

#include "service.h"

/**
 * @fn     int service_init()
 * @brief  Perform the pre-requisites required for node service to run.
 * 		   1. Parsing of the property JSON
 * 		   2. Reading node schema from the inventory.
 * 		   3. Initializing ledger.
 *
 * @return On Success, USYS_OK (0)
 * 		   On Failure, USYS_NOK (-1)
 */
int service_init();

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
