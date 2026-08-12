/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#define CONFIG_REQUEST_ID_LEN 96

typedef enum {
    CONFIG_APPLY_AWAITING = 0,
    CONFIG_APPLY_IN_PROGRESS,
    CONFIG_APPLY_APPLIED,
    CONFIG_APPLY_FAILED
} ConfigApplyState;

/* Service configuration */
typedef struct {
    char *serviceName;
    int servicePort;
    char *nodedHost;
    int nodedPort;
    char *nodedEP;
    char *starterHost;
    int starterPort;
    char *starterEP;
    char *nodeId;
    void *updateSession;

    ConfigApplyState applyState;
    char requestId[CONFIG_REQUEST_ID_LEN];
} Config;

#endif /* CONFIG_H_ */
