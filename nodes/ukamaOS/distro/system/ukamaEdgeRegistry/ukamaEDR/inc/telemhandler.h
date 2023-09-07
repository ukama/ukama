/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_TELEMHANDLER_H_
#define INC_TELEMHANDLER_H_

void telemhandler_init();
void telemhandler_start();
void telemhandler_stop(size_t timer);
void telemhandler_exit();
void telemhandler_service(size_t timer_id, void* data);

#endif /* INC_TELEMHANDLER_H_ */
