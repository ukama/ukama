/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "log.h"

/*
 * Remote logging is intentionally disabled. Applications emit structured
 * records to stderr; starterd will capture and route them to rlog.d.
 * Keep these symbols temporarily so existing applications continue to link.
 */
void log_remote_init(char *serviceName) {
    (void)serviceName;
}

int log_rlogd(char *message) {
    (void)message;
    return 0;
}

int is_connect_with_rlogd(void) {
    return 0;
}

void log_enable_rlogd(int flag) {
    (void)flag;
}
