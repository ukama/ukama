/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <stdio.h>

#include <ulfius.h>
#include <jansson.h>

#include "notifier.h"
#include "jserdes.h"
#include "usys_log.h"

int notifier_init(Notifier *n, Config *cfg) {

    if (!n || !cfg) return STATUS_NOK;
    memset(n, 0, sizeof(*n));
    n->config = cfg;
    return STATUS_OK;
}

int notifier_send_pa_alarm(Notifier *n, int type, int *retCode) {

    char url[256];
    int rc;
    int status;
    json_t *json = NULL;
    char *body = NULL;

    struct _u_request req;
    struct _u_response resp;

    if (!n || !n->config) return STATUS_NOK;
    if (!n->config->enableNotify) return STATUS_OK;

    if (snprintf(url, sizeof(url), "http://%s:%d%s",
                 n->config->notifyHost,
                 n->config->notifyPort,
                 n->config->notifyPath) >= (int)sizeof(url)) {
        return STATUS_NOK;
    }

    if (json_serialize_pa_alarm_notification(&json, n->config, type) != USYS_TRUE || !json) {
        return STATUS_NOK;
    }

    body = json_dumps(json, 0);
    json_decref(json);
    json = NULL;

    if (!body) return STATUS_NOK;

    ulfius_init_request(&req);
    ulfius_init_response(&resp);

    req.http_url = url;
    req.http_verb = "POST";
    u_map_put(req.map_header, "Content-Type", "application/json");

    ulfius_set_string_body_request(&req, "application/json", body);

    rc = ulfius_send_http_request(&req, &resp);

    if (rc != U_OK) {
        usys_log_error("notify failed url=%s rc=%d", url, rc);
        status = STATUS_NOK;
    } else {
        if (retCode) *retCode = resp.status;
        status = (resp.status == HttpStatus_Accepted ||
                  resp.status == HttpStatus_Ok) ? STATUS_OK : STATUS_NOK;
        if (status != STATUS_OK) {
            usys_log_error("notify bad status url=%s code=%d", url, resp.status);
        }
    }

    ulfius_clean_response(&resp);
    ulfius_clean_request(&req);
    free(body);

    return status;
}
