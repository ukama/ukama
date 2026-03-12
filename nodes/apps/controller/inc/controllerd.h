/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONTROLLERD_H
#define CONTROLLERD_H

#define SERVICE_NAME        "controller.d"
#define DEF_LOG_LEVEL       "DEBUG"

/* REST API */
#define URL_PREFIX          "/v1/controller"
#define API_RES_EP(ep)      "/" ep

/* Default configuration */
#define DEF_LISTEN_ADDR     "0.0.0.0"
#define DEF_LISTEN_PORT     8095
#define DEF_SAMPLE_MS       1000
#define DEF_SERIAL_PORT     "/dev/ttyUSB0"
#define DEF_BAUD_RATE       19200

/* Notify.d integration */
#define DEF_NOTIFY_HOST     "127.0.0.1"
#define DEF_NOTIFY_EP       "/notify/v1/event/"

/* Alarm thresholds (48V system defaults) */
#define DEF_LOW_VOLT_WARN   46.0
#define DEF_LOW_VOLT_CRIT   44.0
#define DEF_HIGH_TEMP_WARN  55.0
#define DEF_HIGH_TEMP_CRIT  65.0

/* Generic charge states - driver-independent */
typedef enum {
    CHARGE_STATE_OFF        = 0,
    CHARGE_STATE_FAULT      = 1,
    CHARGE_STATE_BULK       = 2,
    CHARGE_STATE_ABSORPTION = 3,
    CHARGE_STATE_FLOAT      = 4,
    CHARGE_STATE_STORAGE    = 5,
    CHARGE_STATE_EQUALIZE   = 6,
    CHARGE_STATE_UNKNOWN    = 7
} ChargeState;

/* Status codes */
#define STATUS_OK   0
#define STATUS_NOK  (-1)

#endif /* CONTROLLERD_H */
