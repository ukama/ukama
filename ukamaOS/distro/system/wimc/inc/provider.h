/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WIMC_PROVIDER_H
#define WIMC_PROVIDER_H

#include <stdio.h>
#include <string.h>

int get_service_url_from_provider(WimcCfg *cfg, char *name, char *tag,
				  ServiceURL **urls, int *count);
#endif /* WIMC_PROVIDER_H */
