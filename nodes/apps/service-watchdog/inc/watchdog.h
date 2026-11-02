/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WATCHDOG_H_
#define WATCHDOG_H_

#include <stdint.h>

#include "service_client.h"

#define DEF_WATCHDOG_PORT     18092
#define DEF_DEVICE_URL_FORMAT "http://127.0.0.1:%d/v1"
#define ENV_DEVICE_URL        "UKAMA_DEVICE_URL"
#define ENV_WATCHDOG_PORT     "SERVICE_WATCHDOG_PORT"
#define DEVICE_PATH_LEN       32
#define DEVICE_API_STATE      "/state"
#define DEVICE_API_WATCHDOG   "/service/watchdog"

#define WATCHDOG_REPORT_MSEC 10000
#define WATCHDOG_POLL_SEC    1
#define CLOCK_MSEC_PER_SEC   1000
#define CLOCK_NSEC_PER_MSEC  1000000
#define MIN_SERVICE_PORT     1
#define MAX_SERVICE_PORT     65535
#define DECIMAL_BASE         10

#define ARG_VERSION_LONG   "--version"
#define ARG_VERSION_SHORT  "-v"
#define JTAG_SERVICE_GROUP "service_group"
#define JTAG_REASON        "reason"

/* Return a token-bound observation. No member command is issued here. */
json_t *
watchdog_observation(json_t *group,
                     const ServiceObservation members[SERVICE_MEMBER_COUNT]);

#endif /* WATCHDOG_H_ */
