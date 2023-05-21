/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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

int get_systemInfo_from_initClient(Config *config, char *systemName,
								   char **host, char **port);
void free_system_info(SystemInfo *systemInfo);

#endif /* NODE_INFO_H */
