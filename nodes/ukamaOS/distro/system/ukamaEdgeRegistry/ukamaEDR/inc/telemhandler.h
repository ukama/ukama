/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_TELEMHANDLER_H_
#define INC_TELEMHANDLER_H_

void telemhandler_init();
void telemhandler_start();
void telemhandler_stop(size_t timer);
void telemhandler_exit();
void telemhandler_service(size_t timer_id, void* data);

#endif /* INC_TELEMHANDLER_H_ */
