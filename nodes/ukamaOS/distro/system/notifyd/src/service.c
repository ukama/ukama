/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdlib.h>

#include "service.h"

#include "notification.h"
#include "web_service.h"
#include "web_client.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

int service_init(Config *config) {

    char nodeId[32] = {0};
    
    /* Read Node Info from noded */
    if (getenv(ENV_NOTIFY_DEBUG_MODE)) {
        strcpy(nodeId, DEF_NODE_ID);
        usys_log_info("notify.d: Using default Node ID: %s", nodeId);
    } else {
        if (web_client_init(nodeId, config) == STATUS_NOK) {
            return STATUS_NOK;
        }
    }

    /* Notification Init */
    if (notification_init(nodeId, config) == STATUS_NOK) {
        return STATUS_NOK;
    }

    /* Initialize web server */
    if (web_service_init(config->servicePort) == STATUS_NOK) {
        return STATUS_NOK;
    }
        
    return STATUS_OK;
}

void service_start() {

    web_service_start();
}
