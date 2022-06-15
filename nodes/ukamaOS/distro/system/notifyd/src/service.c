/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "service.h"

#include "notification.h"
#include "web_service.h"
#include "web_client.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"


int service_at_exit() {
    int ret = STATUS_OK;

    /* Exit web service */
    web_service_exit();

    /* Notification Exit */
    ret = notification_exit();
    if (ret) {
        return ret;
    }
    return ret;
}

int service_init(Config *config) {
    int ret = STATUS_OK;

    char nodeId[32] = {0};
    char nodeType[32] = {0};

    /* Read Node Info from noded */
    ret = web_client_init(nodeId, nodeType, config);
    if (ret) {
        return ret;
    }

    /* Notification Init */
    ret = notification_init(nodeId, nodeType, config);
    if (ret) {
        return ret;
    }

    /* Initialize web server */
    ret = web_service_init(config->port);
    if (ret) {
        return ret;
    }
    return ret;
}

void service() {

    web_service_start();
}
