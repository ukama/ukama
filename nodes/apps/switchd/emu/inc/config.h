/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_H
#define CONFIG_H

#include "types.h"

void config_init_default(EmuConfig *cfg);
int config_load(EmuConfig *cfg, int argc, char **argv);
void config_usage(void);
void config_version(void);
int config_parse_log_level(const char *slevel);

#endif /* CONFIG_H */
