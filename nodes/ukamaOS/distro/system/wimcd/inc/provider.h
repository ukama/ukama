/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef WIMC_PROVIDER_H
#define WIMC_PROVIDER_H

#include <stdio.h>
#include <string.h>

#include "wimc.h"

int get_service_url_from_provider(Config *cfg, char *name, char *tag,
                                  ServiceURL **urls, int *count);
#endif /* WIMC_PROVIDER_H */
