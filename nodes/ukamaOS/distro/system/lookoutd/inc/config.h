/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#include <stdbool.h>

typedef enum {
    LOOKOUT_APP_MANAGER_STARTERD = 0,
    LOOKOUT_APP_MANAGER_SUPERVISORD
} LookoutAppManager;

/* Service configuration */
typedef struct {

    int servicePort;
    int nodedPort;
    int starterdPort;

    char *nodeID;

    LookoutAppManager appManager;
    bool isTowerNode;
} Config;

#endif /* CONFIG_H_ */
