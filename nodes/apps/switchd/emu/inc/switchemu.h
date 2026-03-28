/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SWITCHEMU_H
#define SWITCHEMU_H

#define SERVICE_NAME              "switchemu.d"
#define SERVICE_SWITCH_EMU        "switch-emu"
#define SERVICE_NOTIFY            "notify"

#define DEF_BIND_ADDRESS          "0.0.0.0"
#define DEF_LOG_LEVEL             "INFO"
#define DEF_SCENARIO              "normal"
#define DEF_STATE_FILE            ""
#define DEF_NOTIFY_HOST           "127.0.0.1"
#define DEF_NOTIFY_PATH           "/v1/alarms"

#define DEF_HTTP_PORT             18088
#define DEF_SNMP_PORT             1161
#define DEF_TFTP_PORT             1069
#define DEF_NOTIFY_PORT           9094

#define EMU_MAX_PORTS             10
#define EMU_MAX_STR               128
#define EMU_MAX_PATH              256
#define EMU_MAX_ALARMS            32
#define EMU_HTTP_REQ_BUF          4096
#define EMU_HTTP_RESP_BUF         8192
#define EMU_SNMP_BUF              1024
#define EMU_TFTP_BUF              1024

#define STATUS_OK                 0
#define STATUS_NOK               -1

#endif /* SWITCHEMU_H */
