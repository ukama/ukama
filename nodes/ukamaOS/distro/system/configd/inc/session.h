/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
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
	int reason;
	int fileCount;
} ConfigData;

typedef enum {
	CONFIG_UNKNOWN = 0,
	CONFIG_ADDED = 1,
	CONFIG_DELETED = 2,
	CONFIG_UPDATED = 3
} Reason;

typedef enum  {
	STATE_UNKNOWN = 0,
	STATE_UPDATE_AVAILABLE = 1,
	STATE_REQUESTED_REBOOT = 2,
	STATE_UPDATE_CONFIRMED = 3,
} ConfigState;

typedef struct {
	char *app;
	char *fileName;
	ConfigState state;
}AppState;

typedef struct  {
	AppState *apps[MAX_SERVICE_COUNT];
	uint32_t timestamp;
	char* version;
	uint32_t count;
	uint32_t expectedCount;
	bool configdVer;
	bool stored;
}ConfigSession;

#endif
