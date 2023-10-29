/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */


#ifndef UTILS_TIMER_H_
#define UTILS_TIMER_H_

#include <stdlib.h>
#include <pthread.h>

typedef enum
{
TIMER_SINGLE_SHOT = 0, /*Periodic Timer*/
TIMER_PERIODIC         /*Single Shot Timer*/
} t_timer;

typedef void (*time_handler)(size_t timer_id, void * user_data);

int initialize(pthread_t* threadid);
size_t start_timer(unsigned int interval, time_handler handler, t_timer type, void * user_data);
void stop_timer(size_t timer_id);
void finalize(pthread_t threadid);

#endif /* UTILS_TIMER_H_ */
