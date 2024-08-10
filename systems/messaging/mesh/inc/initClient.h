/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef NODE_INFO_H
#define NODE_INFO_H

#include <jansson.h>

#include "config.h"

#define TRUE 1
#define FALSE 0

#define MAX_URL_LEN 1024

/* overall node info */
typedef struct systemInfo_ {

	char *systemName;
	char *systemId;
	char *certificate;
	char *ip;
	char *port;
	char *health;
} SystemInfo;

int get_systemInfo_from_initClient(char *systemName,
                                   char **systemHost,
                                   char **systemPort);
void free_system_info(SystemInfo *systemInfo);

#endif /* NODE_INFO_H */
