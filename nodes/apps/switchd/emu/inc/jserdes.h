/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSERDES_H
#define JSERDES_H

#include <stddef.h>

#include "types.h"

int json_serialize_state(const EmuModel *model, char *buf, size_t len);
int json_serialize_ports(const EmuModel *model, char *buf, size_t len);
int json_serialize_firmware(const EmuModel *model, char *buf, size_t len);
int json_serialize_result_ok(char *buf, size_t len);
int json_serialize_error(const char *err, char *buf, size_t len);

#endif /* JSERDES_H */
