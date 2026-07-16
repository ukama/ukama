/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef INGEST_H_
#define INGEST_H_

#include "log_store.h"

typedef struct IngestServer IngestServer;

IngestServer *ingest_start(const char *socketPath, LogStore *store);
void ingest_stop(IngestServer *server);
const char *ingest_socket_path(const IngestServer *server);

#endif /* INGEST_H_ */
