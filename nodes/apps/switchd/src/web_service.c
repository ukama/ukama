/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <stdlib.h>

#include <jansson.h>

#include "http_status.h"
#include "jserdes.h"
#include "json_types.h"
#include "policy.h"
#include "switchd.h"
#include "utils.h"
#include "web_service.h"

#include "usys_log.h"

#include "version.h"

static void ws_setup_unsupported_methods(UInst *serviceInst,
                                         const char *allowed,
                                         const char *prefix,
                                         const char *resource) {
    if (strcmp(allowed, "GET") != 0) {
        ulfius_add_endpoint_by_val(serviceInst, "GET", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowed);
    }
    if (strcmp(allowed, "POST") != 0) {
        ulfius_add_endpoint_by_val(serviceInst, "POST", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowed);
    }
    if (strcmp(allowed, "PUT") != 0) {
        ulfius_add_endpoint_by_val(serviceInst, "PUT", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowed);
    }
    if (strcmp(allowed, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(serviceInst, "DELETE", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowed);
    }
}

static int ws_reply_text(UResponse *response, int status, const char *body) {
    ulfius_set_string_body_response(response, status, body ? body : "");
    return U_CALLBACK_CONTINUE;
}

static int ws_reply_json(UResponse *response, int status, JsonObj *json) {
    char *body;

    if (json == NULL) {
        ulfius_add_header_to_response(response,
                                      "Content-Type",
                                      "application/json");
        ulfius_set_string_body_response(response, status, "{}");
        return U_CALLBACK_CONTINUE;
    }

    body = json_dumps(json, JSON_INDENT(2));
    json_free(&json);

    if (body == NULL) {
        ulfius_add_header_to_response(response,
                                      "Content-Type",
                                      "application/json");
        ulfius_set_string_body_response(response, status, "{}");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_add_header_to_response(response,
                                  "Content-Type",
                                  "application/json");
    ulfius_set_string_body_response(response, status, body);
    free(body);

    return U_CALLBACK_CONTINUE;
}

static SwitchdContext *ws_ctx(void *epConfig) {
    return (SwitchdContext *)epConfig;
}

static void ws_get_source(const URequest *request, char *source, size_t len) {
    JsonErrObj err;
    JsonObj *json;
    JsonObj *entry;
    const char *value;

    if (source == NULL || len == 0) {
        return;
    }

    source[0] = '\0';
    if (request == NULL || request->binary_body == NULL ||
        request->binary_body_length == 0) {
        return;
    }

    memset(&err, 0, sizeof(err));
    json = json_loadb((const char *)request->binary_body,
                      request->binary_body_length,
                      0,
                      &err);
    if (json == NULL) {
        return;
    }

    entry = json_object_get(json, JTAG_SOURCE);
    if (entry != NULL && json_is_string(entry)) {
        value = json_string_value(entry);
        snprintf(source, len, "%s", value ? value : "");
    }

    json_decref(json);
}

static int ws_policy_reply(UResponse *response,
                           int ret,
                           const char *detail) {
    JsonObj *json;

    json = json_object();
    json_object_set_new(json,
                        JTAG_ERROR,
                        json_string(switch_error_to_str(ret)));
    json_object_set_new(json, JTAG_DETAIL, json_string(detail ? detail : ""));

    if (ret == SWITCHD_ERR_AUTH) {
        return ws_reply_json(response, HttpStatus_Forbidden, json);
    }
    if (ret == SWITCHD_ERR_INVAL) {
        return ws_reply_json(response, HttpStatus_BadRequest, json);
    }
    if (ret == SWITCHD_ERR_NOTFOUND) {
        return ws_reply_json(response, HttpStatus_NotFound, json);
    }

    return ws_reply_json(response, HttpStatus_InternalServerError, json);
}

extern JsonObj *switchd_debug_status_json(SwitchdContext *ctx);

static int ws_get_port_id(const URequest *request, uint32_t *portId) {
    const char *value;

    value = u_map_get(request->map_url, "id");
    if (value == NULL || *value == '\0') {
        return STATUS_NOK;
    }

    *portId = (uint32_t)strtoul(value, NULL, 10);
    return (*portId > 0) ? STATUS_OK : STATUS_NOK;
}

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig) {
    (void)request;
    (void)epConfig;

    return ws_reply_text(response,
                         HttpStatus_OK,
                         HttpStatusStr(HttpStatus_OK));
}

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig) {
    (void)request;
    (void)epConfig;

    return ws_reply_text(response, HttpStatus_OK, VERSION);
}

int web_service_cb_get_metrics(const URequest *request,
                               UResponse *response,
                               void *epConfig) {
    SwitchdContext *ctx;
    JsonObj *json;
    int ret;

    (void)request;

    ctx = ws_ctx(epConfig);
    if (ctx == NULL) {
        return ws_reply_json(response, HttpStatus_OK, NULL);
    }

    /*
     * /v1/metrics is consumed by metrics.d, so it must return the current
     * switch view, not an empty cache. Keep this endpoint simple: refresh
     * aggregate KPIs and per-port state, then serialize the cached result.
     *
     * switchd_refresh_* already serializes access to the driver using
     * ctx->driverMutex, so this is safe even when the poller is running.
     */
    ret = switchd_refresh_kpis(ctx);
    if (ret != SWITCHD_OK) {
        usys_log_error("switchd: /v1/metrics refresh_kpis failed: %d", ret);
    }

    ret = switchd_refresh_ports(ctx);
    if (ret != SWITCHD_OK) {
        usys_log_error("switchd: /v1/metrics refresh_ports failed: %d", ret);
    }

    json = json_serialize_metrics(ctx);
    return ws_reply_json(response, HttpStatus_OK, json);
}

int web_service_cb_get_status(const URequest *request,
                             UResponse *response,
                             void *epConfig) {
    (void)request;

    return ws_reply_json(response,
                         HttpStatus_OK,
                         switchd_debug_status_json(ws_ctx(epConfig)));
}

int web_service_cb_get_switch(const URequest *request,
                              UResponse *response,
                              void *epConfig) {
    JsonObj *json;

    (void)request;
    json = json_serialize_switch_info(ws_ctx(epConfig));
    return ws_reply_json(response, HttpStatus_OK, json);
}

int web_service_cb_get_switch_health(const URequest *request,
                                     UResponse *response,
                                     void *epConfig) {
    JsonObj *json;

    (void)request;
    json = json_serialize_switch_health(ws_ctx(epConfig));
    return ws_reply_json(response, HttpStatus_OK, json);
}

int web_service_cb_get_switch_capabilities(const URequest *request,
                                           UResponse *response,
                                           void *epConfig) {
    JsonObj *json;

    (void)request;
    json = json_serialize_switch_capabilities(ws_ctx(epConfig));
    return ws_reply_json(response, HttpStatus_OK, json);
}

int web_service_cb_get_switch_alarms(const URequest *request,
                                     UResponse *response,
                                     void *epConfig) {
    JsonObj *json;

    (void)request;
    json = json_serialize_active_alarms(ws_ctx(epConfig));
    return ws_reply_json(response, HttpStatus_OK, json);
}

int web_service_cb_get_switch_kpis(const URequest *request,
                                   UResponse *response,
                                   void *epConfig) {
    JsonObj *json;

    (void)request;
    json = json_serialize_switch_kpis(ws_ctx(epConfig));
    return ws_reply_json(response, HttpStatus_OK, json);
}

int web_service_cb_get_ports_policy(const URequest *request,
                                    UResponse *response,
                                    void *epConfig) {
    (void)request;
    return ws_reply_json(response,
                         HttpStatus_OK,
                         policy_serialize(ws_ctx(epConfig)));
}

int web_service_cb_put_ports_policy(const URequest *request,
                                    UResponse *response,
                                    void *epConfig) {
    SwitchdContext *ctx;
    char err[SWITCHD_OP_DETAIL_LEN];
    int ret;

    ctx = ws_ctx(epConfig);
    memset(err, 0, sizeof(err));

    if (ctx == NULL || request == NULL || request->binary_body == NULL ||
        request->binary_body_length == 0) {
        return ws_reply_text(response,
                             HttpStatus_BadRequest,
                             "empty policy body");
    }

    ret = policy_apply_body(ctx,
                            (const char *)request->binary_body,
                            request->binary_body_length,
                            err,
                            sizeof(err));
    if (ret != SWITCHD_OK) {
        return ws_policy_reply(response, ret, err);
    }

    return ws_reply_json(response,
                         HttpStatus_OK,
                         policy_serialize(ctx));
}

int web_service_cb_get_ports(const URequest *request,
                             UResponse *response,
                             void *epConfig) {
    JsonObj *json;

    (void)request;
    json = json_serialize_ports(ws_ctx(epConfig));
    return ws_reply_json(response, HttpStatus_OK, json);
}

int web_service_cb_get_port(const URequest *request,
                            UResponse *response,
                            void *epConfig) {
    SwitchdContext *ctx;
    SwitchPortState *port;
    uint32_t portId;
    JsonObj *json;

    ctx = ws_ctx(epConfig);
    if (ws_get_port_id(request, &portId) != STATUS_OK) {
        return ws_reply_text(response,
                             HttpStatus_BadRequest,
                             "bad port id");
    }

    port = switchd_get_port(ctx, portId);
    if (port == NULL) {
        return ws_reply_text(response,
                             HttpStatus_NotFound,
                             "port not found");
    }

    json = json_serialize_port_with_policy(ctx, port);
    return ws_reply_json(response, HttpStatus_OK, json);
}

int web_service_cb_get_port_kpis(const URequest *request,
                                 UResponse *response,
                                 void *epConfig) {
    return web_service_cb_get_port(request, response, epConfig);
}

int web_service_cb_post_port_admin(const URequest *request,
                                   UResponse *response,
                                   void *epConfig) {
    SwitchdContext *ctx;
    uint32_t portId;
    bool up;
    int ret;
    char source[SWITCHD_NAME_LEN];
    char err[SWITCHD_OP_DETAIL_LEN];

    ctx = ws_ctx(epConfig);
    if (ws_get_port_id(request, &portId) != STATUS_OK ||
        !json_deserialize_bool_request(request, JTAG_UP, &up)) {
        return ws_reply_text(response,
                             HttpStatus_BadRequest,
                             "expected {\"up\":true|false}");
    }

    ws_get_source(request, source, sizeof(source));
    memset(err, 0, sizeof(err));
    ret = policy_check_action(ctx,
                              portId,
                              up ? SWITCH_POLICY_ACTION_ADMIN_UP :
                                   SWITCH_POLICY_ACTION_ADMIN_DOWN,
                              source,
                              err,
                              sizeof(err));
    if (ret != SWITCHD_OK) {
        return ws_policy_reply(response, ret, err);
    }

    ret = switchd_set_port_admin(ctx, portId, up);
    if (ret == SWITCHD_ERR_BUSY) {
        return ws_reply_text(response,
                             HttpStatus_Conflict,
                             "operation in progress");
    } else if (ret == SWITCHD_ERR_NOTFOUND) {
        return ws_reply_text(response,
                             HttpStatus_NotFound,
                             "port not found");
    } else if (ret != SWITCHD_OK) {
        return ws_reply_text(response,
                             HttpStatus_InternalServerError,
                             switch_error_to_str(ret));
    }

    return ws_reply_text(response, HttpStatus_OK, HttpStatusStr(HttpStatus_OK));
}

int web_service_cb_post_port_poe(const URequest *request,
                                 UResponse *response,
                                 void *epConfig) {
    SwitchdContext *ctx;
    uint32_t portId;
    bool on;
    int ret;
    char source[SWITCHD_NAME_LEN];
    char err[SWITCHD_OP_DETAIL_LEN];

    ctx = ws_ctx(epConfig);
    if (ws_get_port_id(request, &portId) != STATUS_OK ||
        !json_deserialize_bool_request(request, JTAG_ON, &on)) {
        return ws_reply_text(response,
                             HttpStatus_BadRequest,
                             "expected {\"on\":true|false}");
    }

    ws_get_source(request, source, sizeof(source));
    memset(err, 0, sizeof(err));
    ret = policy_check_action(ctx,
                              portId,
                              on ? SWITCH_POLICY_ACTION_POE_ON :
                                   SWITCH_POLICY_ACTION_POE_OFF,
                              source,
                              err,
                              sizeof(err));
    if (ret != SWITCHD_OK) {
        return ws_policy_reply(response, ret, err);
    }

    ret = switchd_set_port_poe(ctx, portId, on);
    if (ret == SWITCHD_ERR_BUSY) {
        return ws_reply_text(response,
                             HttpStatus_Conflict,
                             "operation in progress");
    } else if (ret == SWITCHD_ERR_NOTFOUND) {
        return ws_reply_text(response,
                             HttpStatus_NotFound,
                             "port not found");
    } else if (ret == SWITCHD_ERR_UNSUPPORTED) {
        return ws_reply_text(response,
                             HttpStatus_NotImplemented,
                             "PoE not supported on this port");
    } else if (ret != SWITCHD_OK) {
        return ws_reply_text(response,
                             HttpStatus_InternalServerError,
                             switch_error_to_str(ret));
    }

    return ws_reply_text(response, HttpStatus_OK, HttpStatusStr(HttpStatus_OK));
}

int web_service_cb_post_port_poe_cycle(const URequest *request,
                                       UResponse *response,
                                       void *epConfig) {
    SwitchdContext *ctx;
    uint32_t portId;
    int offMs;
    int ret;
    char source[SWITCHD_NAME_LEN];
    char err[SWITCHD_OP_DETAIL_LEN];

    ctx = ws_ctx(epConfig);
    if (ws_get_port_id(request, &portId) != STATUS_OK) {
        return ws_reply_text(response,
                             HttpStatus_BadRequest,
                             "bad port id");
    }

    ws_get_source(request, source, sizeof(source));
    memset(err, 0, sizeof(err));
    ret = policy_check_action(ctx,
                              portId,
                              SWITCH_POLICY_ACTION_POE_CYCLE,
                              source,
                              err,
                              sizeof(err));
    if (ret != SWITCHD_OK) {
        return ws_policy_reply(response, ret, err);
    }

    offMs = ctx->config.poeCycleMs;
    (void)json_deserialize_int_request(request, JTAG_OFF_MS, &offMs);

    ret = switchd_cycle_port_poe(ctx, portId, offMs);
    if (ret == SWITCHD_ERR_BUSY) {
        return ws_reply_text(response,
                             HttpStatus_Conflict,
                             "operation in progress");
    } else if (ret == SWITCHD_ERR_NOTFOUND) {
        return ws_reply_text(response,
                             HttpStatus_NotFound,
                             "port not found");
    } else if (ret == SWITCHD_ERR_UNSUPPORTED) {
        return ws_reply_text(response,
                             HttpStatus_NotImplemented,
                             "PoE not supported on this port");
    } else if (ret != SWITCHD_OK) {
        return ws_reply_text(response,
                             HttpStatus_InternalServerError,
                             switch_error_to_str(ret));
    }

    return ws_reply_text(response, HttpStatus_OK, HttpStatusStr(HttpStatus_OK));
}

int web_service_cb_get_firmware(const URequest *request,
                                UResponse *response,
                                void *epConfig) {
    JsonObj *json;

    (void)request;
    json = json_serialize_firmware(ws_ctx(epConfig));
    return ws_reply_json(response, HttpStatus_OK, json);
}

int web_service_cb_get_firmware_status(const URequest *request,
                                       UResponse *response,
                                       void *epConfig) {
    return web_service_cb_get_firmware(request, response, epConfig);
}

int web_service_cb_post_firmware_stage(const URequest *request,
                                       UResponse *response,
                                       void *epConfig) {
    SwitchdContext *ctx;
    char path[SWITCHD_STAGE_PATH_LEN];
    char version[SWITCHD_NAME_LEN];
    char sha256[SWITCHD_SHA256_LEN];
    int ret;

    ctx = ws_ctx(epConfig);
    if (!json_deserialize_firmware_stage_request(request,
                                                 path,
                                                 sizeof(path),
                                                 version,
                                                 sizeof(version),
                                                 sha256,
                                                 sizeof(sha256))) {
        return ws_reply_text(response,
                             HttpStatus_BadRequest,
                             "expected path, optional version/sha256");
    }

    ret = switchd_stage_firmware(ctx,
                                 path,
                                 version[0] ? version : NULL,
                                 sha256[0] ? sha256 : NULL);
    if (ret == SWITCHD_ERR_BUSY) {
        return ws_reply_text(response,
                             HttpStatus_Conflict,
                             "operation in progress");
    } else if (ret == SWITCHD_ERR_NOTFOUND) {
        return ws_reply_text(response,
                             HttpStatus_NotFound,
                             "firmware file not readable");
    } else if (ret != SWITCHD_OK) {
        return ws_reply_text(response,
                             HttpStatus_InternalServerError,
                             switch_error_to_str(ret));
    }

    return ws_reply_json(response,
                         HttpStatus_Accepted,
                         json_serialize_firmware(ctx));
}

int web_service_cb_post_firmware_apply(const URequest *request,
                                       UResponse *response,
                                       void *epConfig) {
    SwitchdContext *ctx;
    int ret;

    (void)request;
    ctx = ws_ctx(epConfig);
    ret = switchd_apply_firmware(ctx);
    if (ret == SWITCHD_ERR_BUSY) {
        return ws_reply_text(response,
                             HttpStatus_Conflict,
                             "operation in progress");
    } else if (ret == SWITCHD_ERR_STATE) {
        return ws_reply_text(response,
                             HttpStatus_Conflict,
                             "firmware is not staged");
    } else if (ret != SWITCHD_OK) {
        return ws_reply_text(response,
                             HttpStatus_InternalServerError,
                             switch_error_to_str(ret));
    }

    return ws_reply_json(response,
                         HttpStatus_Accepted,
                         json_serialize_firmware(ctx));
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {
    (void)request;
    (void)epConfig;

    return ws_reply_text(response,
                         HttpStatus_NotFound,
                         HttpStatusStr(HttpStatus_NotFound));
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *epConfig) {
    (void)request;
    (void)epConfig;

    return ws_reply_text(response,
                         HttpStatus_MethodNotAllowed,
                         HttpStatusStr(HttpStatus_MethodNotAllowed));
}

int web_service_start(SwitchdContext *ctx, UInst *serviceInst) {
    if (ulfius_init_instance(serviceInst,
                             ctx->config.httpPort,
                             NULL,
                             NULL) != U_OK) {
        return STATUS_NOK;
    }

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/ping",
                               NULL,
                               0,
                               &web_service_cb_ping,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1", "ping");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/version",
                               NULL,
                               0,
                               &web_service_cb_version,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1", "version");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/metrics",
                               NULL,
                               0,
                               &web_service_cb_get_metrics,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1", "metrics");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/status",
                               NULL,
                               0,
                               &web_service_cb_get_status,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1", "status");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/switch",
                               NULL,
                               0,
                               &web_service_cb_get_switch,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1", "switch");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/switch/health",
                               NULL,
                               0,
                               &web_service_cb_get_switch_health,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1/switch", "health");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/switch/capabilities",
                               NULL,
                               0,
                               &web_service_cb_get_switch_capabilities,
                               ctx);
    ws_setup_unsupported_methods(serviceInst,
                                 "GET",
                                 "/v1/switch",
                                 "capabilities");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/switch/alarms",
                               NULL,
                               0,
                               &web_service_cb_get_switch_alarms,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1/switch", "alarms");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/switch/kpis",
                               NULL,
                               0,
                               &web_service_cb_get_switch_kpis,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1/switch", "kpis");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/ports/policy",
                               NULL,
                               0,
                               &web_service_cb_get_ports_policy,
                               ctx);

    ulfius_add_endpoint_by_val(serviceInst,
                               "PUT",
                               "/v1/ports/policy",
                               NULL,
                               0,
                               &web_service_cb_put_ports_policy,
                               ctx);

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/ports",
                               NULL,
                               0,
                               &web_service_cb_get_ports,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1", "ports");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/ports/:id",
                               NULL,
                               0,
                               &web_service_cb_get_port,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1", "ports/:id");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/ports/:id/kpis",
                               NULL,
                               0,
                               &web_service_cb_get_port_kpis,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1", "ports/:id/kpis");

    ulfius_add_endpoint_by_val(serviceInst,
                               "POST",
                               "/v1/ports/:id/admin",
                               NULL,
                               0,
                               &web_service_cb_post_port_admin,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "POST", "/v1", "ports/:id/admin");

    ulfius_add_endpoint_by_val(serviceInst,
                               "POST",
                               "/v1/ports/:id/poe",
                               NULL,
                               0,
                               &web_service_cb_post_port_poe,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "POST", "/v1", "ports/:id/poe");

    ulfius_add_endpoint_by_val(serviceInst,
                               "POST",
                               "/v1/ports/:id/poe/cycle",
                               NULL,
                               0,
                               &web_service_cb_post_port_poe_cycle,
                               ctx);
    ws_setup_unsupported_methods(serviceInst,
                                 "POST",
                                 "/v1",
                                 "ports/:id/poe/cycle");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/firmware",
                               NULL,
                               0,
                               &web_service_cb_get_firmware,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1", "firmware");

    ulfius_add_endpoint_by_val(serviceInst,
                               "GET",
                               "/v1/firmware/status",
                               NULL,
                               0,
                               &web_service_cb_get_firmware_status,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "GET", "/v1/firmware", "status");

    ulfius_add_endpoint_by_val(serviceInst,
                               "POST",
                               "/v1/firmware/stage",
                               NULL,
                               0,
                               &web_service_cb_post_firmware_stage,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "POST", "/v1/firmware", "stage");

    ulfius_add_endpoint_by_val(serviceInst,
                               "POST",
                               "/v1/firmware/apply",
                               NULL,
                               0,
                               &web_service_cb_post_firmware_apply,
                               ctx);
    ws_setup_unsupported_methods(serviceInst, "POST", "/v1/firmware", "apply");

    if (ulfius_set_default_endpoint(serviceInst,
                                    &web_service_cb_default,
                                    ctx) != U_OK) {
        usys_log_error("Failed to set default endpoint");
        ulfius_clean_instance(serviceInst);
        return STATUS_NOK;
    }

    if (ulfius_start_framework(serviceInst) != U_OK) {
        ulfius_clean_instance(serviceInst);
        return STATUS_NOK;
    }

    usys_log_debug("%s web service listening on %s:%d",
                   SERVICE_NAME,
                   ctx->config.httpHost,
                   ctx->config.httpPort);

    return STATUS_OK;
}

void web_service_stop(UInst *serviceInst) {
    if (serviceInst == NULL) {
        return;
    }

    ulfius_stop_framework(serviceInst);
    ulfius_clean_instance(serviceInst);
}
