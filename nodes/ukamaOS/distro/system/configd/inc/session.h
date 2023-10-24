/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#ifndef INC_SESSION_H_
#define INC_SESSION_H_

#include "config_macros.h"
#include "usys_types.h"

typedef struct  {
	char *fileName;
	char *app;
	char *data;
	char *version;
	int timestamp;
} ConfigData;

typedef enum  {
	UNKNOWN = 0,
	UPDATE_AVAILABLE = 1,
	REQUESTED_REBOOT = 2,
	UPDATE_CONFIRMED = 3,
} ConfigState;

typedef struct {
	char *app;
	ConfigState state;
}AppState;

typedef struct  {
	AppState *apps[MAX_SERVICE_COUNT];
	uint32_t timestamp;
	char* version;
	uint32_t count;
}ConfigSession;

#endif
