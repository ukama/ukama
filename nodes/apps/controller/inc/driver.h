/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef DRIVER_H
#define DRIVER_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#include "controllerd.h"

/*
 * Generic charge controller data structure.
 * All drivers must populate this structure.
 */
typedef struct {
    /* Solar panel measurements */
    double      pv_voltage_v;       /* Panel voltage (V) */
    double      pv_current_a;       /* Panel current (A) */
    double      pv_power_w;         /* Panel power (W) */
    double      yield_today_kwh;    /* Daily yield (kWh) */
    double      yield_total_kwh;    /* Total yield (kWh) */

    /* Battery measurements */
    double      batt_voltage_v;     /* Battery voltage (V) */
    double      batt_current_a;     /* Battery current (A) */
    int         batt_soc_pct;       /* State of charge (%), -1 if unavailable */
    ChargeState charge_state;       /* Current charge state */

    /* Controller info */
    char        firmware[32];       /* Firmware version string */
    char        serial[32];         /* Serial number */
    char        product_id[16];     /* Product ID (hex string) */
    double      temperature_c;      /* Internal temperature (°C), NAN if unavailable */
    uint32_t    error_code;         /* Error/fault code */

    /* Current settings (read from device) */
    double      absorption_voltage_v;
    double      float_voltage_v;

    /* Load output (if available) */
    bool        load_output_available;
    bool        load_output_state;
    double      load_current_a;

    /* Relay (if available) */
    bool        relay_available;
    bool        relay_state;

    /* Timestamp */
    uint64_t    timestamp_ms;

    /* Communication status */
    bool        comm_ok;
    int         comm_errors;
} ControllerData;

/*
 * Driver interface (vtable pattern).
 * Each charge controller vendor implements this interface.
 */
typedef struct ControllerDriver {
    const char *name;
    const char *description;

    /* Lifecycle */
    int  (*open)(void *ctx, const char *port, int baud);
    void (*close)(void *ctx);

    /* Data acquisition */
    int  (*read_data)(void *ctx, ControllerData *out);

    /* Configuration (optional, return -1 if not supported) */
    int  (*set_absorption_voltage)(void *ctx, double voltage_v);
    int  (*set_float_voltage)(void *ctx, double voltage_v);
    int  (*set_charge_mode)(void *ctx, int mode);

    /* Control (optional, return -1 if not supported) */
    int  (*set_relay)(void *ctx, bool state);
    int  (*set_load_output)(void *ctx, bool state);

    /* Context size for allocation */
    size_t ctx_size;
} ControllerDriver;

/* Driver registry functions */
const ControllerDriver *driver_find(const char *name);
void driver_list_available(void);

/* Built-in drivers */
extern const ControllerDriver victron_driver;
/* Future: extern const ControllerDriver epever_driver; */

/* Helper to get charge state string */
const char *charge_state_str(ChargeState state);

/* Helper to get error code string */
const char *error_code_str(uint32_t code);

#endif /* DRIVER_H */
