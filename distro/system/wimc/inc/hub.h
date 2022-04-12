/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WIMC_HUB_H
#define WIMC_HUB_H

#include <stdio.h>
#include <string.h>
#include <curl/curl.h>

#include "wimc.h"

void free_artifact(Artifact *artifact);
int get_artifacts_info_from_hub(Artifact *artifact, WimcCfg *cfg,
				char *name, char *tag, CURLcode *curlCode);

#endif /* WIMC_HUB_H */
