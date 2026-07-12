/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef LOG_STORE_H_
#define LOG_STORE_H_

#include <jansson.h>
#include <stddef.h>
#include <stdint.h>

typedef struct LogStore LogStore;

typedef struct {
    const char *logDir;
    const char *nodeId;
    size_t rotateBytes;
    int rotateSeconds;
    size_t retainBytes;
    int retainDays;
} LogStoreConfig;

LogStore *log_store_open(const LogStoreConfig *config);
void log_store_close(LogStore *store);

int log_store_append(LogStore *store, json_t *record, uint64_t *seqOut);

const char *log_store_boot_id(const LogStore *store);
const char *log_store_node_id(const LogStore *store);
const char *log_store_active_path(const LogStore *store);
const char *log_store_state_dir(const LogStore *store);
uint64_t log_store_current_seq(const LogStore *store);
size_t log_store_active_bytes(const LogStore *store);

#endif /* LOG_STORE_H_ */
