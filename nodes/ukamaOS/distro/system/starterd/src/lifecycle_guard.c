/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <unistd.h>

#include "service_api.h"
#include "lifecycle_guard.h"
#include "guard_client.h"
#include <stdint.h>
#include <time.h>
#include "usys_log.h"

#define GUARD_WAIT_MSEC     60000
#define GUARD_POLL_SEC      1
#define CLOCK_MSEC_PER_SEC  1000
#define CLOCK_NSEC_PER_MSEC 1000000
#define CLOCK_NSEC_PER_SEC  1000000000

static uint64_t guard_clock_msec(void) {

    struct timespec now;

    if (clock_gettime(CLOCK_MONOTONIC, &now) != 0) {
        return 0;
    }

    return (uint64_t)now.tv_sec * CLOCK_MSEC_PER_SEC +
           now.tv_nsec / CLOCK_NSEC_PER_MSEC;
}

static int notify_guard(App *app, bool active) {

    json_t *body;
    int code;

    body = json_pack("{s:s,s:b,s:I}", "app", app->name, "active", active,
                     "token", (json_int_t)app->guardToken);

    if (body == NULL) {
        return HTTP_TRANSPORT_ERROR;
    }

    code = guard_http(HTTP_POST, app->lifecycleGuard, body, NULL);
    json_decref(body);

    return code;
}

bool lifecycle_guard_begin(App *app) {

    uint64_t deadline;
    struct timespec now;
    int code;

    if (app->lifecycleGuard == NULL) {
        return true;
    }

    if (app->guardActive) {
        lifecycle_guard_end(app);

        if (app->guardActive) {
            return false;
        }
    }

    if (clock_gettime(CLOCK_MONOTONIC, &now) != 0) {
        return false;
    }

    app->guardToken  = (uint64_t)now.tv_sec * CLOCK_NSEC_PER_SEC + now.tv_nsec;
    app->guardActive = true;
    deadline         = guard_clock_msec() + GUARD_WAIT_MSEC;
    do {
        code = notify_guard(app, true);

        if (code == HttpStatus_OK) {
            return true;
        }

        if (code != HttpStatus_Accepted) {
            break;
        }

        sleep(GUARD_POLL_SEC);
    } while (guard_clock_msec() < deadline);
    usys_log_error("lifecycle guard did not release %s/%s", app->space,
                   app->name);
    lifecycle_guard_end(app);

    return false;
}

void lifecycle_guard_end(App *app) {

    if (app->lifecycleGuard == NULL || !app->guardActive) {
        return;
    }

    if (notify_guard(app, false) == HttpStatus_OK) {
        app->guardActive = false;
    }

    /* The supervisor retries a lost completion while idle. Never restore
     * the guarded service here; that decision belongs to its operator.
     */
}
