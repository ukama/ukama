/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "service.h"

#include "notify_macros.h"
#include "web_service.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"


int service_at_exit() {
    int ret = STATUS_OK;

    return ret;
}

int service_init(int port) {
    int ret = STATUS_OK;

    ret = web_service_init(port);
    if (ret) {
        return ret;
    }
    return ret;
}

void service() {
    web_service_start();
}
