/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#include "types.h"

int config_load(SwitchdConfig *cfg);
void config_usage(void);
void config_print_env_help(void);

#endif /* CONFIG_H_ */
