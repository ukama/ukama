/* Copyright 2019 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef __CROS_EC_BASEBOARD_CHARGER_MT6370_H
#define __CROS_EC_BASEBOARD_CHARGER_MT6370_H

#include "charge_state.h"

void mt6370_charger_profile_override(struct charge_state_data *curr);

#endif
