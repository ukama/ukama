/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SAFETY_H
#define SAFETY_H

#include <pthread.h>
#include <stdbool.h>

#include "usys_types.h"
#include "gpio_controller.h"
#include "safety_config.h"

#include "jobs.h"
#include "snapshot.h"
#include "notifier.h"

typedef struct {
    pthread_t        thread;
    pthread_mutex_t  mu;

    Jobs            *jobs;
    SnapshotStore   *snap;
    Notifier        *notifier;

    SafetyConfig     cfg;
    bool             initialized;
    bool             running;

    bool             paDisabled[3];
} Safety;

int  safety_init(Safety *s, Jobs *jobs, SnapshotStore *snap, Notifier *notifier, const char *cfgPath);
void safety_cleanup(Safety *s);

int  safety_start(Safety *s);
int  safety_stop(Safety *s);

int  safety_tick(Safety *s, FemUnit unit);

int  safety_get_config(Safety *s, SafetyConfig *out);
int  safety_set_config(Safety *s, const SafetyConfig *in);

int  safety_force_restore(Safety *s, FemUnit unit);

#endif /* SAFETY_H */
