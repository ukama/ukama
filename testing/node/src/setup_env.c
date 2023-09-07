/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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



