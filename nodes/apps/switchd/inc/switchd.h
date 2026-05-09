/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SWITCHD_H_
#define SWITCHD_H_

#include <jansson.h>
#include <ulfius.h>

#include "types.h"

#include "usys_services.h"
#include "usys_types.h"

#ifndef SERVICE_SWITCH
#define SERVICE_SWITCH "switch.d"
#endif

#ifndef SERVICE_NOTIFY
#define SERVICE_NOTIFY "notify.d"
#endif

#define SERVICE_NAME SERVICE_SWITCH

#define STATUS_OK  (0)
#define STATUS_NOK (-1)

#define DEF_LOG_LEVEL                "INFO"
#define DEF_HTTP_HOST                "0.0.0.0"
#define DEF_SNMP_HOST                "127.0.0.1"
#define DEF_SNMP_COMMUNITY           "public"
#define DEF_TFTP_BIND_IP             "0.0.0.0"
#define DEF_TFTP_ROOT                "/tmp/switchd-tftp"
#define DEF_POLICY_PATH              "/ukama/configs/switchd/policy.json"
#define DEF_URL_PREFIX               "/v1"
#define DEF_DRIVER_NAME              "tycon_snmp"
#define DEF_NOTIFY_HOST              "127.0.0.1"
#define DEF_NOTIFY_EP                "/notify/v1/event/"

#define ENV_SWITCHD_LOG_LEVEL              "SWITCHD_LOG_LEVEL"
#define ENV_SWITCHD_DRIVER                 "SWITCHD_DRIVER"
#define ENV_SWITCHD_HTTP_HOST              "SWITCHD_HTTP_HOST"
#define ENV_SWITCHD_HTTP_PORT              "SWITCHD_HTTP_PORT"
#define ENV_SWITCHD_URL_PREFIX             "SWITCHD_URL_PREFIX"
#define ENV_SWITCHD_SNMP_HOST              "SWITCHD_SNMP_HOST"
#define ENV_SWITCHD_SNMP_PORT              "SWITCHD_SNMP_PORT"
#define ENV_SWITCHD_SNMP_COMMUNITY         "SWITCHD_SNMP_COMMUNITY"
#define ENV_SWITCHD_SNMP_VERSION           "SWITCHD_SNMP_VERSION"
#define ENV_SWITCHD_SNMP_TIMEOUT_MS        "SWITCHD_SNMP_TIMEOUT_MS"
#define ENV_SWITCHD_SNMP_RETRIES           "SWITCHD_SNMP_RETRIES"
#define ENV_SWITCHD_POLL_STATUS_SEC        "SWITCHD_POLL_STATUS_SEC"
#define ENV_SWITCHD_POLL_KPIS_SEC          "SWITCHD_POLL_KPIS_SEC"
#define ENV_SWITCHD_POLL_INFO_SEC          "SWITCHD_POLL_INFO_SEC"
#define ENV_SWITCHD_ALARM_SCAN_SEC         "SWITCHD_ALARM_SCAN_SEC"
#define ENV_SWITCHD_COMMAND_TIMEOUT_MS     "SWITCHD_COMMAND_TIMEOUT_MS"
#define ENV_SWITCHD_FIRMWARE_RECONNECT_SEC "SWITCHD_FIRMWARE_RECONNECT_SEC"
#define ENV_SWITCHD_FIRMWARE_VERIFY_SEC    "SWITCHD_FIRMWARE_VERIFY_SEC"
#define ENV_SWITCHD_POE_CYCLE_MS           "SWITCHD_POE_CYCLE_MS"
#define ENV_SWITCHD_NOTIFY_URL             "SWITCHD_NOTIFY_URL"
#define ENV_SWITCHD_NOTIFY_TIMEOUT_MS      "SWITCHD_NOTIFY_TIMEOUT_MS"
#define ENV_SWITCHD_NOTIFY_HOST            "SWITCHD_NOTIFY_HOST"
#define ENV_SWITCHD_NOTIFY_PORT            "SWITCHD_NOTIFY_PORT"
#define ENV_SWITCHD_TFTP_BIND_IP           "SWITCHD_TFTP_BIND_IP"
#define ENV_SWITCHD_TFTP_PORT              "SWITCHD_TFTP_PORT"
#define ENV_SWITCHD_TFTP_ROOT              "SWITCHD_TFTP_ROOT"
#define ENV_SWITCHD_STRICT_LINK_ALARMS     "SWITCHD_STRICT_LINK_ALARMS"
#define ENV_SWITCHD_SAVE_AFTER_WRITE       "SWITCHD_SAVE_AFTER_WRITE"
#define ENV_SWITCHD_POLICY_PATH            "SWITCHD_POLICY_PATH"

typedef struct _u_instance UInst;
typedef struct _u_request URequest;
typedef struct _u_response UResponse;
typedef json_t JsonObj;
typedef json_error_t JsonErrObj;

extern SwitchdContext gSwitchd;

int switchd_init(SwitchdContext *ctx);
int switchd_start(SwitchdContext *ctx);
void switchd_request_terminate(SwitchdContext *ctx);
void switchd_stop(SwitchdContext *ctx);
void switchd_cleanup(SwitchdContext *ctx);

int switchd_refresh_info(SwitchdContext *ctx);
int switchd_refresh_ports(SwitchdContext *ctx);
int switchd_refresh_kpis(SwitchdContext *ctx);

int switchd_set_port_admin(SwitchdContext *ctx, uint32_t portId, bool up);
int switchd_set_port_poe(SwitchdContext *ctx, uint32_t portId, bool on);
int switchd_cycle_port_poe(SwitchdContext *ctx, uint32_t portId, int offMs);
int switchd_stage_firmware(SwitchdContext *ctx,
                           const char *path,
                           const char *version,
                           const char *sha256);
int switchd_apply_firmware(SwitchdContext *ctx);

SwitchPortState *switchd_get_port(SwitchdContext *ctx, uint32_t portId);

#endif /* SWITCHD_H_ */
