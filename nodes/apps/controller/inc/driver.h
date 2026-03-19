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

typedef struct {
    double      pv_voltage_v;       /* Panel voltage (V) */
    double      pv_current_a;       /* Panel current (A) */
    double      pv_power_w;         /* Panel power (W) */
    double      yield_today_kwh;    /* Daily yield (kWh) */
    double      yield_total_kwh;    /* Total yield (kWh) */

    double      batt_voltage_v;     /* Battery voltage (V) */
    double      batt_current_a;     /* Battery current (A) */
    int         batt_soc_pct;       /* State of charge (%), -1 if unavailable */
    ChargeState charge_state;       /* Current charge state */

    char        firmware[32];       /* Firmware version string */
    char        serial[32];         /* Serial number */
    char        product_id[16];     /* Product ID (hex string) */
    double      temperature_c;      /* Internal temperature (°C), NAN if unavailable */
    uint32_t    error_code;         /* Error/fault code */

    double      absorption_voltage_v;
    double      float_voltage_v;

    bool        load_output_available;
    bool        load_output_state;
    double      load_current_a;

    bool        relay_available;
    bool        relay_state;

    uint64_t    timestamp_ms;

    bool        comm_ok;
    int         comm_errors;
} ControllerData;

typedef struct ControllerDriver {
    const char *name;
    const char *description;

    int  (*open)(void *ctx, const char *port, int baud);
    void (*close)(void *ctx);

    int  (*read_data)(void *ctx, ControllerData *out);

    int  (*set_absorption_voltage)(void *ctx, double voltage_v);
    int  (*set_float_voltage)(void *ctx, double voltage_v);
    int  (*set_charge_mode)(void *ctx, int mode);

    int  (*set_relay)(void *ctx, bool state);
    int  (*set_load_output)(void *ctx, bool state);

    size_t ctx_size;
} ControllerDriver;

const ControllerDriver *driver_find(const char *name);
void driver_list_available(void);

extern const ControllerDriver victron_driver;
/* Future: extern const ControllerDriver epever_driver; */

const char *charge_state_str(ChargeState state);
const char *error_code_str(uint32_t code);

#endif /* DRIVER_H */
