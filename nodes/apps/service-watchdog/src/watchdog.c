/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "watchdog.h"

json_t *
watchdog_observation(json_t *group,
                     const ServiceObservation members[SERVICE_MEMBER_COUNT]) {

    const char *phase   = json_string_value(json_object_get(group, "phase"));
    const char *current = json_string_value(json_object_get(group, "current"));
    json_t *token       = json_object_get(group, "token");
    json_t *remaining   = json_object_get(group, "remaining_ms");
    bool healthy        = true;
    bool consistent     = true;
    bool fault;
    size_t i;

    if (phase == NULL || current == NULL || !json_is_integer(token) ||
        !json_is_integer(remaining)) {
        return NULL;
    }

    if (strcmp(current, "on") && strcmp(current, "off")) {
        return NULL;
    }

    for (i = 0; i < SERVICE_MEMBER_COUNT; i++) {
        healthy    = healthy && members[i].reachable && members[i].ready;
        consistent = consistent && members[i].reachable &&
                     members[i].confirmed &&
                     members[i].on == !strcmp(current, "on");
    }

    if (!strcmp(phase, "changing")) {
        /* Mixed service flags are expected; a dead dependency is not. */
        fault = !healthy || json_integer_value(remaining) <= 0;
    } else if (!strcmp(phase, "idle")) {
        fault = !consistent || (!strcmp(current, "on") && !healthy);
    } else if (!strcmp(phase, "fault")) {
        fault = true;
    } else {
        return NULL;
    }

    return json_pack("{s:I,s:b,s:b}", "token", json_integer_value(token),
                     "healthy", healthy, "fault", fault);
}
