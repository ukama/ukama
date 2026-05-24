/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef PCRF_H_
#define PCRF_H_

#include <stdbool.h>

#include "epcemu.h"

int pcrf_probe(EpcemuConfig *config, EpcemuStatus *status);
bool pcrf_is_ready(EpcemuConfig *config);
int pcrf_create_session(EpcemuConfig *config, const char *imsi,
                        const char *ip, const char *apn);
int pcrf_delete_session(EpcemuConfig *config, const char *imsi);

#endif /* PCRF_H_ */
