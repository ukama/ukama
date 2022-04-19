/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_TIMER_H
#define USYS_TIMER_H

#include "usys_types.h"

/**
 * @fn bool usys_timer(uint32_t, void(*)())
 * @brief
 *
 * @param  resolution
 * @param  tick_handler
 * @return On Success true
 *         On Failure false
 */
bool usys_timer(uint32_t resolution, void (*tick_handler)());

#endif /* USYS_TIMER_H */
