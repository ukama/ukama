/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SAFETY_H
#define SAFETY_H

#include <stdint.h>
#include <stdbool.h>

#include "yaml_config.h"
#include "snapshot.h"
#include "jobs.h"
#include "notifier.h"

typedef struct {
    YamlSafetyConfig cfg;

    bool     paShutdown[3];
    uint32_t lastShutdownMs[3];
    uint32_t okStreak[3];
    uint32_t violationCount[3];
} SafetyState;

typedef struct {
    SafetyState   st;
    Jobs          *jobs;
    SnapshotStore *snap;
    Notifier      *notifier;
} Safety;

int  safety_init(Safety *s, Jobs *jobs, SnapshotStore *snap, Notifier *n, const char *yamlPath);
int  safety_tick(Safety *s, FemUnit unit);

#endif /* SAFETY_H */
