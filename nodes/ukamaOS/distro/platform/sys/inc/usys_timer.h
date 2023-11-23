/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
