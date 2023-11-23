/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

/*
 * Prepare build environment for virtual node
 * 1. Checks if running on local machine or container
 * 2. Extract the source code and for building
 *    virtual node.
 * 3. Set UAKAM_OS env variable.
 */

#include "setup_env.h"

#define SETUP_SCRIPT     "./scripts/mk_env.sh"
#define MAX_BUFFER 1024

int prepare_env_for_creating_virtual_node(){

	char runMe[MAX_BUFFER] = {0};
	sprintf(runMe, "%s", SETUP_SCRIPT);
	if (system(runMe) != 0) return FALSE;

	return TRUE;
}



