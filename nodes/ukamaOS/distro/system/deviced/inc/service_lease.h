/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SERVICE_LEASE_H_
#define SERVICE_LEASE_H_

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#define SERVICE_LEASE_MSEC          30000
#define SERVICE_LEASE_ENV           "UKAMA_SERVICE_LEASE"
#define SERVICE_CLOCK_MSEC_PER_SEC  1000
#define SERVICE_CLOCK_NSEC_PER_MSEC 1000000
#define SERVICE_LEASE_PATH          "/run/ukama/service.lease"

/* All group members share the node's /run and monotonic clock. The lease
 * permits service; it never enables a member. Only device.d writes it.
 * A missing, malformed, expired or implausible lease always means OFF.
 */
static inline const char *service_lease_path(void) {

    const char *path = getenv(SERVICE_LEASE_ENV);

    return path && path[0] == '/' ? path : SERVICE_LEASE_PATH;
}

static inline uint64_t service_clock_msec(void) {

    struct timespec now;

    if (clock_gettime(CLOCK_MONOTONIC, &now) != 0) {
        return 0;
    }

    return (uint64_t)now.tv_sec * SERVICE_CLOCK_MSEC_PER_SEC +
           now.tv_nsec / SERVICE_CLOCK_NSEC_PER_MSEC;
}

static inline bool service_lease_valid(void) {

    FILE *file;
    unsigned long long deadline;
    uint64_t now;
    char extra;
    int fields;

    file = fopen(service_lease_path(), "r");

    if (file == NULL) {
        return false;
    }

    fields = fscanf(file, "%llu %c", &deadline, &extra);
    fclose(file);
    now = service_clock_msec();

    return fields == 1 && now != 0 && deadline > now &&
           deadline - now <= SERVICE_LEASE_MSEC;
}

#endif /* SERVICE_LEASE_H_ */
