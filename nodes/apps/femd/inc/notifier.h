/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef NOTIFIER_H
#define NOTIFIER_H

#include <stdint.h>

#include "config.h"

typedef struct {
    Config *config;
} Notifier;

int notifier_init(Notifier *n, Config *cfg);
int notifier_send_pa_alarm(Notifier *n, int type, int *retCode);

#endif /* NOTIFIER_H */
