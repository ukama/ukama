/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <jansson.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

typedef struct LogSpool LogSpool;

typedef struct {
    uint64_t trace;
    uint64_t debug;
    uint64_t info;
    uint64_t warn;
    uint64_t error;
    uint64_t critical;
} LogDropCounts;

typedef bool (*LogSpoolSendFn)(json_t *frame, void *data);

LogSpool *log_spool_open(const char *directory, size_t maxBytes);
void log_spool_close(LogSpool *spool);

int log_spool_append(LogSpool *spool, json_t *frame);
bool log_spool_has_backlog(LogSpool *spool);
int log_spool_replay(LogSpool *spool,
                     size_t maxRecords,
                     LogSpoolSendFn sendFn,
                     void *data);

void log_spool_get_dropped(LogSpool *spool, LogDropCounts *counts);
void log_spool_clear_dropped(LogSpool *spool);
size_t log_spool_bytes(LogSpool *spool);
