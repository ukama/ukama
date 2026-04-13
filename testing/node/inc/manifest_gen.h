/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef MANIFEST_GEN_H
#define MANIFEST_GEN_H

#include "config.h"

int create_manifest_config(Configs *configs);
void purge_manifest_config(const char *fileName);

#endif
