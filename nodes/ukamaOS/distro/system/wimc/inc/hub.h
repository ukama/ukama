/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef WIMC_HUB_H
#define WIMC_HUB_H

#include <stdio.h>
#include <string.h>
#include <curl/curl.h>

#include "wimc.h"

void free_artifact(Artifact *artifact);
int get_artifacts_info_from_hub(Artifact *artifact, Config *config,
				char *name, char *tag, CURLcode *curlCode);

#endif /* WIMC_HUB_H */
