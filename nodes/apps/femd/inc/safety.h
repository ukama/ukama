/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SAFETY_H
#define SAFETY_H

#include <stdbool.h>

#include "femd.h"
typedef struct {
    float maxTemperatureC;
    float maxReversePowerDbm;
    float maxPaCurrentA;
} SafetyConfig;

typedef struct {
    SafetyConfig cfg;
    bool         paDisabled[3];
    bool         initialized;
} SafetyState;

typedef struct {
    struct Jobs         *jobs;
    struct SnapshotStore *snap;
    SafetyState          st;
} Safety;

int  safety_init(Safety *s, struct Jobs *jobs, struct SnapshotStore *snap, const SafetyConfig *cfg);
void safety_cleanup(Safety *s);

int  safety_tick(Safety *s, FemUnit unit);

#endif /* SAFETY_H */
