/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */


#include "sys_error.h"
#include "sys_string.h"
#include "sys_sync.h"
#include "sys_thread.h"
#include "sys_types.h"
#include "usys_log.h"

int main() {
    usys_log_set_level(LOG_TRACE);

    usys_log_info("Starting test app for platform.");

    return 0;
}





