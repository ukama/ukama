/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef DRV_VICTRON_H
#define DRV_VICTRON_H

#include <stdint.h>
#include <stdbool.h>
#include <pthread.h>

#include "driver.h"

/*
 * Victron VE.Direct Text-Mode Protocol Driver
 *
 * Serial configuration: 19200 baud, 8N1, no flow control
 * Data format: <Label>\t<Value>\r\n blocks every ~1 second
 * Checksum: Modulo-256 sum of all bytes in block = 0
 *
 * Reference: VE.Direct Protocol v3.34
 */

#define VICTRON_BAUD_RATE       19200
#define VICTRON_MAX_LABEL_LEN   16
#define VICTRON_MAX_VALUE_LEN   64
#define VICTRON_MAX_FIELDS      32
#define VICTRON_FRAME_TIMEOUT_MS 3000

/* VE.Direct field labels */
#define VE_LABEL_VOLTAGE        "V"      /* Battery voltage (mV) */
#define VE_LABEL_CURRENT        "I"      /* Battery current (mA) */
#define VE_LABEL_PV_VOLTAGE     "VPV"    /* Panel voltage (mV) */
#define VE_LABEL_PV_POWER       "PPV"    /* Panel power (W) */
#define VE_LABEL_CHARGE_STATE   "CS"     /* Charge state */
#define VE_LABEL_ERROR          "ERR"    /* Error code */
#define VE_LABEL_LOAD           "LOAD"   /* Load output state (ON/OFF) */
#define VE_LABEL_LOAD_CURRENT   "IL"     /* Load current (mA) */
#define VE_LABEL_YIELD_TOTAL    "H19"    /* Total yield (0.01 kWh) */
#define VE_LABEL_YIELD_TODAY    "H20"    /* Today's yield (0.01 kWh) */
#define VE_LABEL_MAX_POWER_TODAY "H21"   /* Max power today (W) */
#define VE_LABEL_YIELD_YESTERDAY "H22"   /* Yesterday's yield (0.01 kWh) */
#define VE_LABEL_MAX_POWER_YESTERDAY "H23" /* Max power yesterday (W) */
#define VE_LABEL_FIRMWARE       "FW"     /* Firmware version */
#define VE_LABEL_SERIAL         "SER#"   /* Serial number */
#define VE_LABEL_PRODUCT_ID     "PID"    /* Product ID (hex) */
#define VE_LABEL_RELAY          "Relay"  /* Relay state (ON/OFF) */
#define VE_LABEL_ALARM          "AR"     /* Alarm reason */
#define VE_LABEL_MPPT           "MPPT"   /* MPPT tracker state */
#define VE_LABEL_TEMP           "T"      /* Battery temperature (°C, only when ext sensor paired) */
#define VE_LABEL_CHECKSUM       "Checksum"

/* Victron VE.Direct charge state values (CS field) */
typedef enum {
    VCS_OFF                  = 0,
    VCS_LOW_POWER            = 1,
    VCS_FAULT                = 2,
    VCS_BULK                 = 3,
    VCS_ABSORPTION           = 4,
    VCS_FLOAT                = 5,
    VCS_STORAGE              = 6,
    VCS_EQUALIZE             = 7,
    VCS_INVERTING            = 9,
    VCS_POWER_SUPPLY         = 11,
    VCS_STARTING             = 245,
    VCS_REPEATED_ABSORPTION  = 246,
    VCS_AUTO_EQUALIZE        = 247,
    VCS_BATTERY_SAFE         = 248,
    VCS_EXTERNAL_CONTROL     = 252
} VictronChargeState;

/* Victron VE.Direct error codes (ERR field) */
typedef enum {
    VERR_NONE                    = 0,
    VERR_BATTERY_VOLTAGE_HIGH    = 2,
    VERR_CHARGER_TEMP_HIGH       = 17,
    VERR_CHARGER_OVERCURRENT     = 18,
    VERR_CHARGER_CURRENT_REVERSED = 19,
    VERR_BULK_TIME_LIMIT         = 20,
    VERR_CURRENT_SENSOR_FAIL     = 21,
    VERR_TERMINALS_OVERHEATED    = 26,
    VERR_CONVERTER_ISSUE         = 28,
    VERR_INPUT_VOLTAGE_HIGH      = 33,
    VERR_INPUT_CURRENT_HIGH      = 34,
    VERR_INPUT_SHUTDOWN_BATTERY  = 38,
    VERR_INPUT_SHUTDOWN_CURRENT  = 39,
    VERR_LOST_COMMUNICATION      = 65,
    VERR_SYNC_CHARGING_CONFIG    = 66,
    VERR_BMS_CONNECTION_LOST     = 67,
    VERR_NETWORK_MISCONFIGURED   = 68,
    VERR_FACTORY_CALIBRATION     = 116,
    VERR_INVALID_FIRMWARE        = 117,
    VERR_USER_SETTINGS_INVALID   = 119
} VictronErrorCode;

/* Parsed VE.Direct field */
typedef struct {
    char label[VICTRON_MAX_LABEL_LEN];
    char value[VICTRON_MAX_VALUE_LEN];
} VeDirectField;

/* Victron driver context */
typedef struct {
    int             fd;                 /* Serial port file descriptor */
    char            port[128];          /* Serial port path */
    int             baud;               /* Baud rate */

    /* Receive buffer */
    char            rx_buf[512];
    int             rx_len;

    /* Last parsed frame */
    VeDirectField   fields[VICTRON_MAX_FIELDS];
    int             field_count;
    bool            frame_valid;
    uint64_t        last_frame_ts;

    /* Cached data */
    ControllerData  cached_data;
    pthread_mutex_t lock;

    /* Statistics */
    uint32_t        frames_received;
    uint32_t        checksum_errors;
    uint32_t        parse_errors;
} VictronCtx;

/* Driver implementation */
int  victron_open(void *ctx, const char *port, int baud);
void victron_close(void *ctx);
int  victron_read_data(void *ctx, ControllerData *out);
int  victron_set_absorption_voltage(void *ctx, double voltage_v);
int  victron_set_float_voltage(void *ctx, double voltage_v);
int  victron_set_relay(void *ctx, bool state);

#endif /* DRV_VICTRON_H */
